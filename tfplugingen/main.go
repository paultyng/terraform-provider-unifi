package main

import (
	"context"
	"fmt"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"golang.org/x/tools/go/packages"
)

func main() {
	err := run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, args []string) error {
	dir := "../internal/provider"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("unable to read dir %s: %w", dir, err)
	}

	conf := config{}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".hcl" {
			continue
		}

		fn := filepath.Join(dir, f.Name())

		log.Printf("decoding %s", fn)

		var fileConf config
		err := hclsimple.DecodeFile(fn, nil, &fileConf)
		if err != nil {
			return fmt.Errorf("unable to decode file %s: %w", fn, err)
		}

		if fileConf.Provider != nil {
			conf.Provider = fileConf.Provider
		}
		conf.Resources = append(conf.Resources, fileConf.Resources...)
	}

	fmt.Printf("%#v\n%#v\n", conf.Provider, conf.Resources)

	// TODO: validate
	// - resources should have names

	err = collectSDKInformation(ctx, dir, conf)
	if err != nil {
		return fmt.Errorf("unable to build model: %w", err)
	}

	return nil
}

func coallesce(s ...string) string {
	for _, v := range s {
		if v != "" {
			return v
		}
	}
	return ""
}

func uniq(s ...string) []string {
	n := []string{}
	for _, v := range s {
		if v != "" {
			n = append(n, v)
		}
	}
	return n
}

func collectSDKInformation(ctx context.Context, dir string, conf config) error {
	pkgConfig := &packages.Config{
		Mode:    packages.NeedTypes | packages.NeedName | packages.NeedImports | packages.NeedDeps,
		Context: ctx,
		Logf:    log.Printf,
		Dir:     dir,
		// TODO: Env, BuildFlags,
	}

	// build list of package names
	pkgPaths := []string{
		conf.Provider.DefaultClientPackage,
		conf.Provider.DefaultResourcePackage,
		conf.Provider.DefaultSDKPackage,
	}
	for _, resConf := range conf.Resources {
		pkgPaths = append(pkgPaths, resConf.Package)
	}

	log.Println(pkgPaths)
	pkgPaths = uniq(pkgPaths...)
	log.Println(pkgPaths)

	provPkgs, err := packages.Load(pkgConfig)
	if err != nil {
		return fmt.Errorf("unable to load SDK packages: %w", err)
	}

	// TODO: check for errors on provider pkg

	pkgMap := map[string]*packages.Package{}
	resources := map[string]*types.Struct{}

	for _, pkg := range provPkgs[0].Imports {
		if len(pkg.Errors) > 0 {
			return fmt.Errorf("unable to load package %s: %w", pkg.Name, pkg.Errors[0])
		}

		fmt.Println(pkg.PkgPath)

		for _, n := range pkgPaths {
			log.Printf("testing %q against %q", pkg.PkgPath, n)
			if pkg.PkgPath == n {
				log.Printf("found %s", n)
				pkgMap[n] = pkg
			}
		}
	}

	for _, resConf := range conf.Resources {
		pkgName := coallesce(resConf.Package, conf.Provider.DefaultResourcePackage, conf.Provider.DefaultSDKPackage)
		pkg := pkgMap[pkgName]
		if pkg == nil {
			return fmt.Errorf("package %s not loaded", pkgName)
		}

		typeName := coallesce(resConf.Type, conf.Provider.DefaultResourceType(resConf.Name))
		tn := pkg.Types.Scope().Lookup(typeName)
		if tn == nil {
			return fmt.Errorf("unable to find resource %s type %s in package %s", resConf.Name, typeName, pkgName)
		}
		tn, ok := tn.(*types.TypeName)
		if !ok {
			return fmt.Errorf("object %s in package %s is not a TypeName", typeName, pkgName)
		}
		next := tn.Type()
		var st *types.Struct
		for next != nil {
			switch ty := next.(type) {
			default:
				next = nil
			case *types.Named:
				next = ty.Underlying()
			case *types.Struct:
				st = ty
				next = nil
			}
		}
		if st == nil {
			return fmt.Errorf("unable to find underlying struct for %s in package %s", typeName, pkgName)
		}
		resources[resConf.Name] = st
	}
	panic("not implemented")
}

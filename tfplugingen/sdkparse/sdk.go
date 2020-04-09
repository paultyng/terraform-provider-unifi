package sdkparse

import (
	"context"
	"fmt"
	"go/types"
	"log"

	"golang.org/x/tools/go/packages"

	"github.com/paultyng/tfplugingen/config"
)

const (
	defaultReadFunc   = config.ResourceTemplate("Get{{ .Type }}")
	defaultCreateFunc = config.ResourceTemplate("Create{{ .Type }}")
	defaultUpdateFunc = config.ResourceTemplate("Update{{ .Type }}")
	defaultDeleteFunc = config.ResourceTemplate("Delete{{ .Type }}")
)

type ResourceInfo struct {
	name string

	conf config.Resource

	typeName   *types.TypeName
	structType *types.Struct

	clientTypeName *types.TypeName
	readFunc       *types.Func
	createFunc     *types.Func
	updateFunc     *types.Func
	deleteFunc     *types.Func

	schema *config.SchemaProperties
}

func loadPackages(ctx context.Context, dir string, logf func(format string, args ...interface{}), paths ...string) (map[string]*packages.Package, error) {
	pkgConfig := &packages.Config{
		Mode:    packages.NeedTypes | packages.NeedName | packages.NeedImports | packages.NeedDeps,
		Context: ctx,
		Logf:    logf,
		Dir:     dir,
		// TODO: Env, BuildFlags,
	}

	provPkgs, err := packages.Load(pkgConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK packages: %w", err)
	}

	// TODO: check for errors on provider pkg

	pkgMap := map[string]*packages.Package{}

	for _, pkg := range provPkgs[0].Imports {
		if len(pkg.Errors) > 0 {
			return nil, fmt.Errorf("unable to load package %s: %w", pkg.Name, pkg.Errors[0])
		}

		for _, n := range paths {
			if pkg.PkgPath == n {
				pkgMap[n] = pkg
			}
		}
	}

	return pkgMap, nil
}

func Resources(ctx context.Context, dir string, conf config.File) (map[string]ResourceInfo, error) {
	// build list of package names
	pkgPaths := []string{
		conf.Provider.DefaultClientPackage,
		conf.Provider.DefaultResourcePackage,
		conf.Provider.DefaultSDKPackage,
	}
	for _, resConf := range conf.Resources {
		pkgPaths = append(pkgPaths, resConf.Package)
	}

	pkgPaths = uniq(pkgPaths...)

	pkgMap, err := loadPackages(ctx, dir, log.Printf, pkgPaths...)
	if err != nil {
		return nil, fmt.Errorf("unable to load packages: %w", err)
	}

	resources := map[string]ResourceInfo{}

	for _, resConf := range conf.Resources {
		res, err := collectResourceInfo(ctx, resConf, *conf.Provider, pkgMap)
		if err != nil {
			return nil, fmt.Errorf("unable to collect resource %q info: %w", resConf.Name, err)
		}

		resources[resConf.Name] = res
	}

	return resources, nil
}

func findResourceType(ctx context.Context, resConf config.Resource, provConf config.Provider, pkg *packages.Package) (*types.TypeName, *types.Struct, error) {
	candidateNames := []string{
		resConf.Type,
		provConf.DefaultResourceType(resConf.Name),
	}
	return findNamedStructTypeName(ctx, candidateNames, pkg)
}

func findNamedStructTypeName(ctx context.Context, names []string, pkg *packages.Package) (*types.TypeName, *types.Struct, error) {
	name := coallesce(names...)
	obj := pkg.Types.Scope().Lookup(name)
	if obj == nil {
		return nil, nil, fmt.Errorf("unable to find type %q in package %q", name, pkg.PkgPath)
	}

	tn, ok := obj.(*types.TypeName)
	if !ok {
		return nil, nil, fmt.Errorf("object %q in package %q is not a TypeName", name, pkg.PkgPath)
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
		return nil, nil, fmt.Errorf("unable to find underlying struct for %q in package %q", tn.Name(), pkg.PkgPath)
	}

	return tn, st, nil
}

func findResourceReadFunc(ctx context.Context, resConf config.Resource, provConf config.Provider, clientTypeName *types.TypeName, resourceTypeName *types.TypeName) (*types.Func, error) {
	templateInfo := config.ResourceTemplateInfo{
		Name: resConf.Name,
		Type: resourceTypeName.Name(),
	}

	candidateNames := []string{
		resConf.ReadFunc,
		provConf.DefaultReadFunc.MustRender(templateInfo),
		defaultReadFunc.MustRender(templateInfo),
	}

	return findResourceFunc(ctx, candidateNames, clientTypeName, true)
}

func findResourceCreateFunc(ctx context.Context, resConf config.Resource, provConf config.Provider, clientTypeName *types.TypeName, resourceTypeName *types.TypeName) (*types.Func, error) {
	templateInfo := config.ResourceTemplateInfo{
		Name: resConf.Name,
		Type: resourceTypeName.Name(),
	}

	candidateNames := []string{
		resConf.CreateFunc,
		provConf.DefaultCreateFunc.MustRender(templateInfo),
		defaultCreateFunc.MustRender(templateInfo),
	}

	return findResourceFunc(ctx, candidateNames, clientTypeName, true)
}

func findResourceUpdateFunc(ctx context.Context, resConf config.Resource, provConf config.Provider, clientTypeName *types.TypeName, resourceTypeName *types.TypeName) (*types.Func, error) {
	templateInfo := config.ResourceTemplateInfo{
		Name: resConf.Name,
		Type: resourceTypeName.Name(),
	}

	candidateNames := []string{
		resConf.UpdateFunc,
		provConf.DefaultUpdateFunc.MustRender(templateInfo),
		defaultUpdateFunc.MustRender(templateInfo),
	}

	return findResourceFunc(ctx, candidateNames, clientTypeName, false)
}

func findResourceDeleteFunc(ctx context.Context, resConf config.Resource, provConf config.Provider, clientTypeName *types.TypeName, resourceTypeName *types.TypeName) (*types.Func, error) {
	templateInfo := config.ResourceTemplateInfo{
		Name: resConf.Name,
		Type: resourceTypeName.Name(),
	}

	candidateNames := []string{
		resConf.DeleteFunc,
		provConf.DefaultDeleteFunc.MustRender(templateInfo),
		defaultDeleteFunc.MustRender(templateInfo),
	}

	return findResourceFunc(ctx, candidateNames, clientTypeName, true)
}

func findResourceFunc(ctx context.Context, names []string, tn *types.TypeName, mustExist bool) (*types.Func, error) {
	name := coallesce(names...)
	obj, _, _ := types.LookupFieldOrMethod(tn.Type(), true, tn.Pkg(), name)
	if obj == nil {
		if mustExist {
			return nil, fmt.Errorf("unable to find resource func %q", name)
		}
		return nil, nil
	}
	switch obj := obj.(type) {
	default:
		return nil, fmt.Errorf("unexpected object type for func %q: %T", name, obj)
	case *types.Func:
		return obj, nil
	}
}

func collectResourceInfo(ctx context.Context, resConf config.Resource, provConf config.Provider, pkgMap map[string]*packages.Package) (ResourceInfo, error) {

	pkgName := coallesce(resConf.Package, provConf.DefaultResourcePackage, provConf.DefaultSDKPackage)
	pkg := pkgMap[pkgName]
	if pkg == nil {
		return ResourceInfo{}, fmt.Errorf("package %q not loaded", pkgName)
	}

	tn, st, err := findResourceType(ctx, resConf, provConf, pkg)
	if err != nil {
		return ResourceInfo{}, fmt.Errorf("unable to find resource type in package %q: %w", pkgName, err)
	}

	// TODO: make this a special cased findClientType?
	ctn, _, err := findNamedStructTypeName(ctx, []string{provConf.DefaultClientType, "Client"}, pkg)
	if err != nil {
		return ResourceInfo{}, fmt.Errorf("unable to find resource client type in package %q: %w", pkgName, err)
	}

	readFunc, err := findResourceReadFunc(ctx, resConf, provConf, ctn, tn)
	if err != nil {
		return ResourceInfo{}, fmt.Errorf("unable to find read func: %w", err)
	}

	createFunc, err := findResourceCreateFunc(ctx, resConf, provConf, ctn, tn)
	if err != nil {
		return ResourceInfo{}, fmt.Errorf("unable to find create func: %w", err)
	}

	updateFunc, err := findResourceUpdateFunc(ctx, resConf, provConf, ctn, tn)
	if err != nil {
		return ResourceInfo{}, fmt.Errorf("unable to find update func: %w", err)
	}

	deleteFunc, err := findResourceDeleteFunc(ctx, resConf, provConf, ctn, tn)
	if err != nil {
		return ResourceInfo{}, fmt.Errorf("unable to find delete func: %w", err)
	}

	props := resConf.SchemaProperties

	schema, err := schemaFromType(st)
	if err != nil {
		return ResourceInfo{}, fmt.Errorf("unable to determine schema from type %s: %w", tn.Name(), err)
	}

	if schema != nil {
		if *schema.Type != "object" {
			return ResourceInfo{}, fmt.Errorf("expected object schema type, got %q", *schema.Type)
		}
		props = schema.SchemaProperties
	}

	return ResourceInfo{
		conf: resConf,
		name: resConf.Name,

		typeName:   tn,
		structType: st,

		clientTypeName: ctn,
		readFunc:       readFunc,
		createFunc:     createFunc,
		updateFunc:     updateFunc,
		deleteFunc:     deleteFunc,

		schema: props,
	}, nil
}

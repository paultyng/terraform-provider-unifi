package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl/v2/hclsimple"

	"github.com/paultyng/tfplugingen/config"
	"github.com/paultyng/tfplugingen/sdkparse"
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

	conf := config.File{}

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".hcl" {
			continue
		}

		fn := filepath.Join(dir, f.Name())

		log.Printf("decoding %s", fn)

		var fileConf config.File
		err := hclsimple.DecodeFile(fn, nil, &fileConf)
		if err != nil {
			return fmt.Errorf("unable to decode file %s: %w", fn, err)
		}

		if fileConf.Provider != nil {
			conf.Provider = fileConf.Provider
		}
		conf.Resources = append(conf.Resources, fileConf.Resources...)
	}

	// TODO: validate
	// - resources should have names

	resources, err := sdkparse.Resources(ctx, dir, conf)
	if err != nil {
		return fmt.Errorf("unable to build model: %w", err)
	}

	err = generateModel(ctx, resources)
	if err != nil {
		return fmt.Errorf("unable to generate model: %w", err)
	}

	return nil
}

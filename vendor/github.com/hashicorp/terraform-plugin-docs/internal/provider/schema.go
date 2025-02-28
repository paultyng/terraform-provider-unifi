// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/checkpoint"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
)

func TerraformProviderSchemaFromTerraform(ctx context.Context, providerName, providerDir, tfVersion string, l *Logger) (*tfjson.ProviderSchema, error) {
	var err error

	shortName := providerShortName(providerName)

	tmpDir, err := os.MkdirTemp("", "tfws")
	if err != nil {
		return nil, fmt.Errorf("unable to create temporary provider install directory %q: %w", tmpDir, err)
	}
	defer os.RemoveAll(tmpDir)

	l.infof("compiling provider %q", shortName)
	providerPath := fmt.Sprintf("plugins/registry.terraform.io/hashicorp/%s/0.0.1/%s_%s", shortName, runtime.GOOS, runtime.GOARCH)
	outFile := filepath.Join(tmpDir, providerPath, fmt.Sprintf("terraform-provider-%s", shortName))
	switch runtime.GOOS {
	case "windows":
		outFile = outFile + ".exe"
	}
	buildCmd := exec.Command("go", "build", "-o", outFile)
	buildCmd.Dir = providerDir
	// TODO: constrain env here to make it a little safer?
	_, err = runCmd(buildCmd)
	if err != nil {
		return nil, fmt.Errorf("unable to execute go build command: %w", err)
	}

	err = writeFile(filepath.Join(tmpDir, "provider.tf"), fmt.Sprintf(`
provider %[1]q {
}
`, shortName))
	if err != nil {
		return nil, fmt.Errorf("unable to write provider.tf file: %w", err)
	}

	i := install.NewInstaller()
	var sources []src.Source
	if tfVersion != "" {
		l.infof("downloading Terraform CLI binary version from releases.hashicorp.com: %s", tfVersion)
		sources = []src.Source{
			&releases.ExactVersion{
				Product:    product.Terraform,
				Version:    version.Must(version.NewVersion(tfVersion)),
				InstallDir: tmpDir,
			},
		}
	} else {
		l.infof("using Terraform CLI binary from PATH if available, otherwise downloading latest Terraform CLI binary")
		sources = []src.Source{
			&fs.AnyVersion{
				Product: &product.Terraform,
			},
			&checkpoint.LatestVersion{
				InstallDir: tmpDir,
				Product:    product.Terraform,
			},
		}
	}

	tfBin, err := i.Ensure(context.Background(), sources)
	if err != nil {
		return nil, fmt.Errorf("unable to download Terraform binary: %w", err)
	}

	tf, err := tfexec.NewTerraform(tmpDir, tfBin)
	if err != nil {
		return nil, fmt.Errorf("unable to create new terraform exec instance: %w", err)
	}

	l.infof("running terraform init")
	err = tf.Init(ctx, tfexec.Get(false), tfexec.PluginDir("./plugins"))
	if err != nil {
		return nil, fmt.Errorf("unable to run terraform init on provider: %w", err)
	}

	l.infof("getting provider schema")
	schemas, err := tf.ProvidersSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve provider schema from terraform exec: %w", err)
	}

	if ps, ok := schemas.Schemas[shortName]; ok {
		return ps, nil
	}

	if ps, ok := schemas.Schemas["registry.terraform.io/hashicorp/"+shortName]; ok {
		return ps, nil
	}

	return nil, fmt.Errorf("unable to find schema in JSON for provider %q", shortName)
}

func TerraformProviderSchemaFromFile(providerName, providersSchemaPath string, l *Logger) (*tfjson.ProviderSchema, error) {
	var err error

	shortName := providerShortName(providerName)

	l.infof("getting provider schema")
	schemas, err := extractSchemaFromFile(providersSchemaPath)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve provider schema from JSON file: %w", err)
	}

	if ps, ok := schemas.Schemas[shortName]; ok {
		return ps, nil
	}

	if ps, ok := schemas.Schemas["registry.terraform.io/hashicorp/"+shortName]; ok {
		return ps, nil
	}

	return nil, fmt.Errorf("unable to find schema in JSON for provider %q", shortName)
}

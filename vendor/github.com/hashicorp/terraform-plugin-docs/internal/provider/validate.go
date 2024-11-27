// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hashicorp/cli"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-docs/internal/check"
)

const (
	FileExtensionHtmlMarkdown = `.html.markdown`
	FileExtensionHtmlMd       = `.html.md`
	FileExtensionMarkdown     = `.markdown`
	FileExtensionMd           = `.md`

	DocumentationGlobPattern    = `{docs/index.*,docs/{,cdktf/}{data-sources,ephemeral-resources,guides,resources,functions}/**/*,website/docs/**/*}`
	DocumentationDirGlobPattern = `{docs/{,cdktf/}{data-sources,ephemeral-resources,guides,resources,functions}{,/*},website/docs/**/*}`
)

var ValidLegacyFileExtensions = []string{
	FileExtensionHtmlMarkdown,
	FileExtensionHtmlMd,
	FileExtensionMarkdown,
	FileExtensionMd,
}

var ValidRegistryFileExtensions = []string{
	FileExtensionMd,
}

var LegacyFrontMatterOptions = &check.FrontMatterOptions{
	NoSidebarCurrent:   true,
	RequireDescription: true,
	RequireLayout:      true,
	RequirePageTitle:   true,
}

var LegacyIndexFrontMatterOptions = &check.FrontMatterOptions{
	NoSidebarCurrent:   true,
	NoSubcategory:      true,
	RequireDescription: true,
	RequireLayout:      true,
	RequirePageTitle:   true,
}

var RegistryFrontMatterOptions = &check.FrontMatterOptions{
	NoLayout:         true,
	NoSidebarCurrent: true,
}

var RegistryIndexFrontMatterOptions = &check.FrontMatterOptions{
	NoLayout:         true,
	NoSidebarCurrent: true,
	NoSubcategory:    true,
}

var RegistryGuideFrontMatterOptions = &check.FrontMatterOptions{
	NoLayout:         true,
	NoSidebarCurrent: true,
	RequirePageTitle: true,
}

type validator struct {
	providerName        string
	providerDir         string
	providerFS          fs.FS
	providersSchemaPath string

	tfVersion      string
	providerSchema *tfjson.ProviderSchema

	logger *Logger
}

func Validate(ui cli.Ui, providerDir, providerName, providersSchemaPath, tfversion string) error {
	// Ensure provider directory is resolved absolute path
	if providerDir == "" {
		wd, err := os.Getwd()

		if err != nil {
			return fmt.Errorf("error getting working directory: %w", err)
		}

		providerDir = wd
	} else {
		absProviderDir, err := filepath.Abs(providerDir)

		if err != nil {
			return fmt.Errorf("error getting absolute path with provider directory %q: %w", providerDir, err)
		}

		providerDir = absProviderDir
	}

	// Verify provider directory
	providerDirFileInfo, err := os.Stat(providerDir)

	if err != nil {
		return fmt.Errorf("error getting information for provider directory %q: %w", providerDir, err)
	}

	if !providerDirFileInfo.IsDir() {
		return fmt.Errorf("expected %q to be a directory", providerDir)
	}

	providerFs := os.DirFS(providerDir)

	v := &validator{
		providerName:        providerName,
		providerDir:         providerDir,
		providerFS:          providerFs,
		providersSchemaPath: providersSchemaPath,
		tfVersion:           tfversion,

		logger: NewLogger(ui),
	}

	ctx := context.Background()

	return v.validate(ctx)
}

func (v *validator) validate(ctx context.Context) error {
	var result error

	var err error

	if v.providerName == "" {
		v.providerName = filepath.Base(v.providerDir)
	}

	if v.providersSchemaPath == "" {
		v.logger.infof("exporting schema from Terraform")
		v.providerSchema, err = TerraformProviderSchemaFromTerraform(ctx, v.providerName, v.providerDir, v.tfVersion, v.logger)
		if err != nil {
			return fmt.Errorf("error exporting provider schema from Terraform: %w", err)
		}
	} else {
		v.logger.infof("exporting schema from JSON file")
		v.providerSchema, err = TerraformProviderSchemaFromFile(v.providerName, v.providersSchemaPath, v.logger)
		if err != nil {
			return fmt.Errorf("error exporting provider schema from JSON file: %w", err)
		}
	}

	files, globErr := doublestar.Glob(v.providerFS, DocumentationGlobPattern)
	if globErr != nil {
		return fmt.Errorf("error finding documentation files: %w", err)
	}

	log.Printf("[DEBUG] Found documentation files %v", files)

	v.logger.infof("running mixed directories check")
	err = check.MixedDirectoriesCheck(files)
	result = errors.Join(result, err)

	if dirExists(v.providerFS, "docs") {
		v.logger.infof("detected static docs directory, running checks")
		err = v.validateStaticDocs()
		result = errors.Join(result, err)

	}
	if dirExists(v.providerFS, "website/docs") {
		v.logger.infof("detected legacy website directory, running checks")
		err = v.validateLegacyWebsite()
		result = errors.Join(result, err)
	}

	return result
}

func (v *validator) validateStaticDocs() error {
	dir := "docs"
	var result error

	options := &check.ProviderFileOptions{
		FileOptions:     &check.FileOptions{BasePath: v.providerDir},
		FrontMatter:     RegistryFrontMatterOptions,
		ValidExtensions: ValidRegistryFileExtensions,
	}

	var files []string

	err := fs.WalkDir(v.providerFS, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory %q: %w", dir, err)
		}
		if d.IsDir() {
			match, err := doublestar.Match(DocumentationDirGlobPattern, path)
			if err != nil {
				return err
			}
			if !match {
				return nil // skip valid non-documentation directories
			}

			v.logger.infof("running invalid directories check on %s", path)
			result = errors.Join(result, check.InvalidDirectoriesCheck(path))
			return nil
		}
		match, err := doublestar.Match(DocumentationGlobPattern, path)
		if err != nil {
			return err
		}
		if !match {
			return nil // skip valid non-documentation files
		}

		// Configure FrontMatterOptions based on file type
		if removeAllExt(d.Name()) == "index" {
			options.FrontMatter = RegistryIndexFrontMatterOptions
		} else if _, relErr := filepath.Rel(dir+"/guides", path); relErr == nil {
			options.FrontMatter = RegistryGuideFrontMatterOptions
		} else {
			options.FrontMatter = RegistryFrontMatterOptions
		}
		v.logger.infof("running file checks on %s", path)
		result = errors.Join(result, check.NewProviderFileCheck(v.providerFS, options).Run(path))

		files = append(files, path)
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking directory %q: %w", dir, err)
	}

	mismatchOpt := &check.FileMismatchOptions{
		ProviderShortName: providerShortName(v.providerName),
		Schema:            v.providerSchema,
	}

	if dirExists(v.providerFS, dir+"/data-sources") {
		dataSourceFiles, _ := fs.ReadDir(v.providerFS, dir+"/data-sources")
		mismatchOpt.DatasourceEntries = dataSourceFiles
	}
	if dirExists(v.providerFS, dir+"/resources") {
		resourceFiles, _ := fs.ReadDir(v.providerFS, dir+"/resources")
		mismatchOpt.ResourceEntries = resourceFiles
	}
	if dirExists(v.providerFS, dir+"/functions") {
		functionFiles, _ := fs.ReadDir(v.providerFS, dir+"/functions")
		mismatchOpt.FunctionEntries = functionFiles
	}
	if dirExists(v.providerFS, dir+"/ephemeral-resources") {
		ephemeralResourceFiles, _ := fs.ReadDir(v.providerFS, dir+"/ephemeral-resources")
		mismatchOpt.EphemeralResourceEntries = ephemeralResourceFiles
	}

	v.logger.infof("running file mismatch check")
	if err := check.NewFileMismatchCheck(mismatchOpt).Run(); err != nil {
		result = errors.Join(result, err)
	}

	return result
}

func (v *validator) validateLegacyWebsite() error {
	dir := "website/docs"
	var result error

	options := &check.ProviderFileOptions{
		FileOptions:     &check.FileOptions{BasePath: v.providerDir},
		FrontMatter:     LegacyFrontMatterOptions,
		ValidExtensions: ValidLegacyFileExtensions,
	}

	var files []string
	err := fs.WalkDir(v.providerFS, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory %q: %w", dir, err)
		}
		if d.IsDir() {
			match, err := doublestar.Match(DocumentationDirGlobPattern, path)
			if err != nil {
				return err
			}
			if !match {
				return nil // skip valid non-documentation directories
			}

			v.logger.infof("running invalid directories check on %s", path)
			result = errors.Join(result, check.InvalidDirectoriesCheck(path))
			return nil
		}

		match, err := doublestar.Match(DocumentationGlobPattern, path)
		if err != nil {
			return err
		}
		if !match {
			return nil // skip non-documentation files
		}

		// Configure FrontMatterOptions based on file type
		if removeAllExt(d.Name()) == "index" {
			options.FrontMatter = LegacyIndexFrontMatterOptions
		} else {
			options.FrontMatter = LegacyFrontMatterOptions
		}
		v.logger.infof("running file checks on %s", path)
		result = errors.Join(result, check.NewProviderFileCheck(v.providerFS, options).Run(path))

		files = append(files, path)
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking directory %q: %w", dir, err)
	}

	mismatchOpt := &check.FileMismatchOptions{
		ProviderShortName: providerShortName(v.providerName),
		Schema:            v.providerSchema,
	}

	if dirExists(v.providerFS, dir+"/d") {
		dataSourceFiles, _ := fs.ReadDir(v.providerFS, dir+"/d")
		mismatchOpt.DatasourceEntries = dataSourceFiles
	}
	if dirExists(v.providerFS, dir+"/r") {
		resourceFiles, _ := fs.ReadDir(v.providerFS, dir+"/r")
		mismatchOpt.ResourceEntries = resourceFiles
	}
	if dirExists(v.providerFS, dir+"/functions") {
		functionFiles, _ := fs.ReadDir(v.providerFS, dir+"/functions")
		mismatchOpt.FunctionEntries = functionFiles
	}
	if dirExists(v.providerFS, dir+"/ephemeral-resources") {
		ephemeralResourceFiles, _ := fs.ReadDir(v.providerFS, dir+"/ephemeral-resources")
		mismatchOpt.EphemeralResourceEntries = ephemeralResourceFiles
	}

	v.logger.infof("running file mismatch check")
	if err := check.NewFileMismatchCheck(mismatchOpt).Run(); err != nil {
		result = errors.Join(result, err)
	}

	return result
}

func dirExists(fileSys fs.FS, name string) bool {
	if file, err := fs.Stat(fileSys, name); err != nil {
		return false
	} else if !file.IsDir() {
		return false
	}

	return true
}

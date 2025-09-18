// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
)

type ProviderFileOptions struct {
	*FileOptions

	FrontMatter     *FrontMatterOptions
	ValidExtensions []string
}

type ProviderFileCheck struct {
	Options    *ProviderFileOptions
	ProviderFs fs.FS
}

func NewProviderFileCheck(providerFs fs.FS, opts *ProviderFileOptions) *ProviderFileCheck {
	check := &ProviderFileCheck{
		Options:    opts,
		ProviderFs: providerFs,
	}

	if check.Options == nil {
		check.Options = &ProviderFileOptions{}
	}

	if check.Options.FileOptions == nil {
		check.Options.FileOptions = &FileOptions{}
	}

	if check.Options.FrontMatter == nil {
		check.Options.FrontMatter = &FrontMatterOptions{}
	}

	return check
}

func (check *ProviderFileCheck) Run(path string) error {
	fullpath := check.Options.FullPath(path)

	log.Printf("[DEBUG] Checking file: %s", fullpath)

	if err := FileExtensionCheck(path, check.Options.ValidExtensions); err != nil {
		return fmt.Errorf("%s: error checking file extension: %w", filepath.FromSlash(path), err)
	}

	if err := FileSizeCheck(check.ProviderFs, path); err != nil {
		return fmt.Errorf("%s: error checking file size: %w", filepath.FromSlash(path), err)
	}

	content, err := fs.ReadFile(check.ProviderFs, path)

	if err != nil {
		return fmt.Errorf("%s: error reading file: %w", filepath.FromSlash(path), err)
	}

	if err := NewFrontMatterCheck(check.Options.FrontMatter).Run(content); err != nil {
		return fmt.Errorf("%s: error checking file frontmatter: %w", filepath.FromSlash(path), err)
	}

	return nil
}

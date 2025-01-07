// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
)

type FileOptions struct {
	BasePath string
}

func (opts *FileOptions) FullPath(path string) string {
	if opts.BasePath != "" {
		return filepath.Join(opts.BasePath, path)
	}

	return path
}

// FileSizeCheck verifies that documentation file is below the Terraform Registry storage limit.
func FileSizeCheck(providerFs fs.FS, path string) error {
	fi, err := fs.Stat(providerFs, path)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] File %s size: %d (limit: %d)", path, fi.Size(), RegistryMaximumSizeOfFile)
	if fi.Size() >= int64(RegistryMaximumSizeOfFile) {
		return fmt.Errorf("exceeded maximum (%d) size of documentation file for Terraform Registry: %d", RegistryMaximumSizeOfFile, fi.Size())
	}

	return nil
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	FileExtensionHtmlMarkdown = `.html.markdown`
	FileExtensionHtmlMd       = `.html.md`
	FileExtensionMarkdown     = `.markdown`
	FileExtensionMd           = `.md`
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

// FileExtensionCheck checks if the file extension of the given path is valid.
func FileExtensionCheck(path string, validExtensions []string) error {
	if !FilePathEndsWithExtensionFrom(path, validExtensions) {
		return fmt.Errorf("file does not end with a valid extension, valid extensions: %v", validExtensions)
	}

	return nil
}

func FilePathEndsWithExtensionFrom(path string, validExtensions []string) bool {
	for _, validExtension := range validExtensions {
		if strings.HasSuffix(path, validExtension) {
			return true
		}
	}

	return false
}

// TrimFileExtension removes file extensions including those with multiple periods.
func TrimFileExtension(path string) string {
	filename := filepath.Base(path)

	if filename == "." {
		return ""
	}

	dotIndex := strings.IndexByte(filename, '.')

	if dotIndex > 0 {
		return filename[:dotIndex]
	}

	return filename
}

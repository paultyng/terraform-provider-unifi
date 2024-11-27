// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
)

const (
	CdktfIndexDirectory = `cdktf`

	LegacyIndexDirectory              = `website/docs`
	LegacyDataSourcesDirectory        = `d`
	LegacyEphemeralResourcesDirectory = `ephemeral-resources`
	LegacyGuidesDirectory             = `guides`
	LegacyResourcesDirectory          = `r`
	LegacyFunctionsDirectory          = `functions`

	RegistryIndexDirectory              = `docs`
	RegistryDataSourcesDirectory        = `data-sources`
	RegistryEphemeralResourcesDirectory = `ephemeral-resources`
	RegistryGuidesDirectory             = `guides`
	RegistryResourcesDirectory          = `resources`
	RegistryFunctionsDirectory          = `functions`

	// Terraform Registry Storage Limits
	// https://www.terraform.io/docs/registry/providers/docs.html#storage-limits
	RegistryMaximumNumberOfFiles = 2000
	RegistryMaximumSizeOfFile    = 500000 // 500KB

)

var ValidLegacyDirectories = []string{
	LegacyIndexDirectory,
	LegacyIndexDirectory + "/" + LegacyDataSourcesDirectory,
	LegacyIndexDirectory + "/" + LegacyEphemeralResourcesDirectory,
	LegacyIndexDirectory + "/" + LegacyGuidesDirectory,
	LegacyIndexDirectory + "/" + LegacyResourcesDirectory,
	LegacyIndexDirectory + "/" + LegacyFunctionsDirectory,
}

var ValidRegistryDirectories = []string{
	RegistryIndexDirectory,
	RegistryIndexDirectory + "/" + RegistryDataSourcesDirectory,
	RegistryIndexDirectory + "/" + RegistryEphemeralResourcesDirectory,
	RegistryIndexDirectory + "/" + RegistryGuidesDirectory,
	RegistryIndexDirectory + "/" + RegistryResourcesDirectory,
	RegistryIndexDirectory + "/" + RegistryFunctionsDirectory,
}

var ValidCdktfLanguages = []string{
	"csharp",
	"go",
	"java",
	"python",
	"typescript",
}

var ValidLegacySubdirectories = []string{
	LegacyIndexDirectory,
	LegacyDataSourcesDirectory,
	LegacyEphemeralResourcesDirectory,
	LegacyGuidesDirectory,
	LegacyResourcesDirectory,
}

var ValidRegistrySubdirectories = []string{
	RegistryIndexDirectory,
	RegistryDataSourcesDirectory,
	RegistryEphemeralResourcesDirectory,
	RegistryGuidesDirectory,
	RegistryResourcesDirectory,
}

func InvalidDirectoriesCheck(dirPath string) error {
	if IsValidRegistryDirectory(dirPath) {
		return nil
	}

	if IsValidLegacyDirectory(dirPath) {
		return nil
	}

	if IsValidCdktfDirectory(dirPath) {
		return nil
	}

	return fmt.Errorf("invalid Terraform Provider documentation directory found: %s", filepath.FromSlash(dirPath))

}

func MixedDirectoriesCheck(docFiles []string) error {
	var legacyDirectoryFound bool
	var registryDirectoryFound bool
	err := fmt.Errorf("mixed Terraform Provider documentation directory layouts found, must use only legacy or registry layout")

	for _, file := range docFiles {
		directory := path.Dir(file)
		log.Printf("[DEBUG] Found directory: %s", directory)

		// Allow docs/ with other files
		if IsValidRegistryDirectory(directory) && directory != RegistryIndexDirectory {
			registryDirectoryFound = true

			if legacyDirectoryFound {
				log.Printf("[DEBUG] Found mixed directories")
				return err
			}
		}

		if IsValidLegacyDirectory(directory) {
			legacyDirectoryFound = true

			if registryDirectoryFound {
				log.Printf("[DEBUG] Found mixed directories")
				return err
			}
		}
	}

	return nil
}

func IsValidLegacyDirectory(directory string) bool {
	for _, validLegacyDirectory := range ValidLegacyDirectories {
		if directory == validLegacyDirectory {
			return true
		}
	}

	return false
}

func IsValidRegistryDirectory(directory string) bool {
	for _, validRegistryDirectory := range ValidRegistryDirectories {
		if directory == validRegistryDirectory {
			return true
		}
	}

	return false
}

func IsValidCdktfDirectory(directory string) bool {
	if directory == fmt.Sprintf("%s/%s", LegacyIndexDirectory, CdktfIndexDirectory) {
		return true
	}

	if directory == fmt.Sprintf("%s/%s", RegistryIndexDirectory, CdktfIndexDirectory) {
		return true
	}

	for _, validCdktfLanguage := range ValidCdktfLanguages {

		if directory == fmt.Sprintf("%s/%s/%s", LegacyIndexDirectory, CdktfIndexDirectory, validCdktfLanguage) {
			return true
		}

		if directory == fmt.Sprintf("%s/%s/%s", RegistryIndexDirectory, CdktfIndexDirectory, validCdktfLanguage) {
			return true
		}

		for _, validLegacySubdirectory := range ValidLegacySubdirectories {
			if directory == fmt.Sprintf("%s/%s/%s/%s", LegacyIndexDirectory, CdktfIndexDirectory, validCdktfLanguage, validLegacySubdirectory) {
				return true
			}
		}

		for _, validRegistrySubdirectory := range ValidRegistrySubdirectories {
			if directory == fmt.Sprintf("%s/%s/%s/%s", RegistryIndexDirectory, CdktfIndexDirectory, validCdktfLanguage, validRegistrySubdirectory) {
				return true
			}
		}
	}

	return false
}

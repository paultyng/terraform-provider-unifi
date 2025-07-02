// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sort"

	tfjson "github.com/hashicorp/terraform-json"
)

type FileMismatchOptions struct {
	*FileOptions

	IgnoreFileMismatch []string

	IgnoreFileMissing []string

	ProviderShortName string

	DatasourceEntries []os.DirEntry

	ResourceEntries []os.DirEntry

	FunctionEntries []os.DirEntry

	EphemeralResourceEntries []os.DirEntry

	Schema *tfjson.ProviderSchema
}

type FileMismatchCheck struct {
	Options *FileMismatchOptions
}

func NewFileMismatchCheck(opts *FileMismatchOptions) *FileMismatchCheck {
	check := &FileMismatchCheck{
		Options: opts,
	}

	if check.Options == nil {
		check.Options = &FileMismatchOptions{}
	}

	if check.Options.FileOptions == nil {
		check.Options.FileOptions = &FileOptions{}
	}

	return check
}

func (check *FileMismatchCheck) Run() error {
	var result error

	if check.Options.Schema == nil {
		log.Printf("[DEBUG] Skipping file mismatch checks due to missing provider schema")
		return nil
	}

	if check.Options.ResourceEntries != nil {
		err := check.ResourceFileMismatchCheck(check.Options.ResourceEntries, "resource", check.Options.Schema.ResourceSchemas)
		result = errors.Join(result, err)
	}

	if check.Options.DatasourceEntries != nil {
		err := check.ResourceFileMismatchCheck(check.Options.DatasourceEntries, "datasource", check.Options.Schema.DataSourceSchemas)
		result = errors.Join(result, err)
	}

	if check.Options.FunctionEntries != nil {
		err := check.FunctionFileMismatchCheck(check.Options.FunctionEntries, check.Options.Schema.Functions)
		result = errors.Join(result, err)
	}

	if check.Options.EphemeralResourceEntries != nil {
		err := check.ResourceFileMismatchCheck(check.Options.EphemeralResourceEntries, "ephemeral resource", check.Options.Schema.EphemeralResourceSchemas)
		result = errors.Join(result, err)
	}

	return result
}

// ResourceFileMismatchCheck checks for mismatched files, either missing or extraneous, against the resource/datasource schema
func (check *FileMismatchCheck) ResourceFileMismatchCheck(files []os.DirEntry, resourceType string, schemas map[string]*tfjson.Schema) error {
	if len(files) == 0 {
		log.Printf("[DEBUG] Skipping %s file mismatch checks due to missing file list", resourceType)
		return nil
	}

	if len(schemas) == 0 {
		log.Printf("[DEBUG] Skipping %s file mismatch checks due to missing schemas", resourceType)
		return nil
	}

	var extraFiles []string
	var missingFiles []string

	for _, file := range files {
		log.Printf("[DEBUG] Found file %s", file.Name())
		if fileHasResource(schemas, check.Options.ProviderShortName, file.Name()) {
			continue
		}

		if check.IgnoreFileMismatch(file.Name()) {
			continue
		}

		log.Printf("[DEBUG] Found extraneous file %s", file.Name())
		extraFiles = append(extraFiles, file.Name())
	}

	for _, resourceName := range resourceNames(schemas) {
		log.Printf("[DEBUG] Found %s %s", resourceType, resourceName)
		if resourceHasFile(files, check.Options.ProviderShortName, resourceName) {
			continue
		}

		if check.IgnoreFileMissing(resourceName) {
			continue
		}

		log.Printf("[DEBUG] Missing file for %s %s", resourceType, resourceName)
		missingFiles = append(missingFiles, resourceName)
	}

	var result error

	for _, extraFile := range extraFiles {
		err := fmt.Errorf("matching %s for documentation file (%s) not found, file is extraneous or incorrectly named", resourceType, extraFile)
		result = errors.Join(result, err)
	}

	for _, missingFile := range missingFiles {
		err := fmt.Errorf("missing documentation file for %s: %s", resourceType, missingFile)
		result = errors.Join(result, err)
	}

	return result

}

// FunctionFileMismatchCheck checks for mismatched files, either missing or extraneous, against the function signature
func (check *FileMismatchCheck) FunctionFileMismatchCheck(files []os.DirEntry, functions map[string]*tfjson.FunctionSignature) error {
	if len(files) == 0 {
		log.Printf("[DEBUG] Skipping function file mismatch checks due to missing file list")
		return nil
	}

	if len(functions) == 0 {
		log.Printf("[DEBUG] Skipping function file mismatch checks due to missing schemas")
		return nil
	}

	var extraFiles []string
	var missingFiles []string

	for _, file := range files {
		if fileHasFunction(functions, file.Name()) {
			continue
		}

		if check.IgnoreFileMismatch(file.Name()) {
			continue
		}

		extraFiles = append(extraFiles, file.Name())
	}

	for _, functionName := range functionNames(functions) {
		if functionHasFile(files, functionName) {
			continue
		}

		if check.IgnoreFileMissing(functionName) {
			continue
		}

		missingFiles = append(missingFiles, functionName)
	}

	var result error

	for _, extraFile := range extraFiles {
		err := fmt.Errorf("matching function for documentation file (%s) not found, file is extraneous or incorrectly named", extraFile)
		result = errors.Join(result, err)
	}

	for _, missingFile := range missingFiles {
		err := fmt.Errorf("missing documentation file for function: %s", missingFile)
		result = errors.Join(result, err)
	}

	return result

}

func (check *FileMismatchCheck) IgnoreFileMismatch(file string) bool {
	for _, ignoreResourceName := range check.Options.IgnoreFileMismatch {
		if ignoreResourceName == fileResourceNameWithProvider(check.Options.ProviderShortName, file) {
			return true
		} else if ignoreResourceName == TrimFileExtension(file) {
			// While uncommon, it is valid for a resource type to be named the same as the provider itself.
			// https://github.com/hashicorp/terraform-plugin-docs/issues/419
			return true
		}
	}

	return false
}

func (check *FileMismatchCheck) IgnoreFileMissing(resourceName string) bool {
	for _, ignoreResourceName := range check.Options.IgnoreFileMissing {
		if ignoreResourceName == resourceName {
			return true
		}
	}

	return false
}

func fileHasResource(schemaResources map[string]*tfjson.Schema, providerName, file string) bool {
	if _, ok := schemaResources[fileResourceNameWithProvider(providerName, file)]; ok {
		return true
	}

	// While uncommon, it is valid for a resource type to be named the same as the provider itself.
	// https://github.com/hashicorp/terraform-plugin-docs/issues/419
	if _, ok := schemaResources[TrimFileExtension(file)]; ok {
		return true
	}

	return false
}

func fileHasFunction(functions map[string]*tfjson.FunctionSignature, file string) bool {
	if _, ok := functions[TrimFileExtension(file)]; ok {
		return true
	}

	return false
}

func fileResourceNameWithProvider(providerName, fileName string) string {
	resourceSuffix := TrimFileExtension(fileName)

	return fmt.Sprintf("%s_%s", providerName, resourceSuffix)
}

func resourceHasFile(files []os.DirEntry, providerName, resourceName string) bool {
	var found bool

	for _, file := range files {
		if fileResourceNameWithProvider(providerName, file.Name()) == resourceName {
			found = true
			break
		} else if TrimFileExtension(file.Name()) == resourceName {
			// While uncommon, it is valid for a resource type to be named the same as the provider itself.
			// https://github.com/hashicorp/terraform-plugin-docs/issues/419
			found = true
			break
		}
	}

	return found
}

func functionHasFile(files []os.DirEntry, functionName string) bool {
	var found bool

	for _, file := range files {
		if TrimFileExtension(file.Name()) == functionName {
			found = true
			break
		}
	}

	return found
}

func resourceNames(resources map[string]*tfjson.Schema) []string {
	names := make([]string, 0, len(resources))

	for name := range resources {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

func functionNames(functions map[string]*tfjson.FunctionSignature) []string {
	names := make([]string, 0, len(functions))

	for name := range functions {
		names = append(names, name)
	}

	sort.Strings(names)

	return names
}

// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Kunde21/markdownfmt/v3/markdown"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

func providerShortName(n string) string {
	return strings.TrimPrefix(n, "terraform-provider-")
}

func resourceShortName(name, providerName string) string {
	psn := providerShortName(providerName)
	return strings.TrimPrefix(name, psn+"_")
}

func copyFile(srcPath, dstPath string, mode os.FileMode) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Ensure destination path exists for file creation
	err = os.MkdirAll(filepath.Dir(dstPath), 0755)
	if err != nil {
		return err
	}

	// If the destination file already exists, we shouldn't blow it away
	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		// If the file already exists, we can skip it without returning an error.
		if errors.Is(err, os.ErrExist) {
			return nil
		}
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func removeAllExt(file string) string {
	for {
		ext := filepath.Ext(file)
		if ext == "" || ext == file {
			return file
		}
		file = strings.TrimSuffix(file, ext)
	}
}

// resourceSchema determines whether there is a schema in the supplied schemas map which
// has either the providerShortName or the providerShortName concatenated with the
// templateFileName (stripped of file extension.
func resourceSchema(schemas map[string]*tfjson.Schema, providerShortName, templateFileName string) (*tfjson.Schema, string) {
	resName := providerShortName + "_" + removeAllExt(templateFileName)
	if schema, ok := schemas[resName]; ok {
		return schema, resName
	}

	if schema, ok := schemas[providerShortName]; ok {
		return schema, providerShortName
	}

	return nil, resName
}

func writeFile(path string, data string) error {
	dir, _ := filepath.Split(path)

	var err error
	if dir != "" {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return fmt.Errorf("unable to make dir %q: %w", dir, err)
		}
	}

	err = os.WriteFile(path, []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("unable to write file %q: %w", path, err)
	}

	return nil
}

//nolint:unparam
func runCmd(cmd *exec.Cmd) ([]byte, error) {
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("error executing %q, %v", cmd.Path, cmd.Args)
		log.Print(string(output))
		return nil, fmt.Errorf("error executing %q: %w", cmd.Path, err)
	}
	return output, nil
}

func cp(srcDir, dstDir string) error {
	err := filepath.Walk(srcDir, func(srcPath string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		switch mode := f.Mode(); {
		case mode.IsDir():
			if err := os.Mkdir(dstPath, f.Mode()); err != nil && !os.IsExist(err) {
				return err
			}
		case mode.IsRegular():
			if err := copyFile(srcPath, dstPath, mode); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unknown file type (%d / %s) for %s", f.Mode(), f.Mode().String(), srcPath)
		}

		return nil
	})
	return err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func extractSchemaFromFile(path string) (*tfjson.ProviderSchemas, error) {
	schemajson, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read file %q: %w", path, err)
	}

	schemas := &tfjson.ProviderSchemas{
		FormatVersion: "",
		Schemas:       nil,
	}
	err = schemas.UnmarshalJSON(schemajson)
	if err != nil {
		return nil, err
	}

	return schemas, nil
}

func newMarkdownRenderer() goldmark.Markdown {
	mr := markdown.NewRenderer()
	extensions := []goldmark.Extender{
		extension.GFM,
		meta.Meta, // We need this to skip YAML frontmatter when parsing.
	}
	parserOptions := []parser.Option{
		parser.WithAttribute(), // We need this to enable # headers {#custom-ids}.
	}

	gm := goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(parserOptions...),
		goldmark.WithRenderer(mr),
	)
	return gm
}

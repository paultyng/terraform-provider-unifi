// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/cli"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type migrator struct {
	// providerDir is the absolute path to the root provider directory
	providerDir string

	websiteDir   string
	templatesDir string
	examplesDir  string

	providerName string

	ui cli.Ui
}

func (m *migrator) infof(format string, a ...interface{}) {
	m.ui.Info(fmt.Sprintf(format, a...))
}

func (m *migrator) warnf(format string, a ...interface{}) {
	m.ui.Warn(fmt.Sprintf(format, a...))
}

func Migrate(ui cli.Ui, providerDir string, templatesDir string, examplesDir string, providerName string) error {
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

	// Default providerName to provider directory name
	if providerName == "" {
		providerName = filepath.Base(providerDir)
	}

	// Determine website directory
	websiteDir, err := determineWebsiteDir(providerDir)
	if err != nil {
		return err
	}

	m := &migrator{
		providerDir:  providerDir,
		templatesDir: templatesDir,
		examplesDir:  examplesDir,
		websiteDir:   websiteDir,
		providerName: providerName,
		ui:           ui,
	}

	return m.Migrate()
}

func (m *migrator) Migrate() error {
	m.infof("migrating website from %q to %q", m.ProviderWebsiteDir(), m.ProviderTemplatesDir())

	err := filepath.WalkDir(m.ProviderWebsiteDir(), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("unable to walk path %q: %w", path, err)
		}

		if d.IsDir() {
			switch d.Name() {
			case "d", "data-sources": //data-sources
				m.infof("migrating datasources directory: %s", d.Name())
				err := filepath.WalkDir(path, m.MigrateTemplate("data-sources"))
				if err != nil {
					return err
				}
				return filepath.SkipDir
			case "r", "resources": //resources
				m.infof("migrating resources directory: %s", d.Name())
				err := filepath.WalkDir(path, m.MigrateTemplate("resources"))
				if err != nil {
					return err
				}
				return filepath.SkipDir
			case "functions":
				m.infof("migrating functions directory: %s", d.Name())
				err := filepath.WalkDir(path, m.MigrateTemplate("functions"))
				if err != nil {
					return err
				}
				return filepath.SkipDir
			case "ephemeral-resources":
				m.infof("migrating ephemeral resources directory: %s", d.Name())
				err := filepath.WalkDir(path, m.MigrateTemplate("ephemeral-resources"))
				if err != nil {
					return err
				}
				return filepath.SkipDir
			case "guides":
				m.infof("copying guides directory: %s", d.Name())
				err := cp(path, filepath.Join(m.ProviderTemplatesDir(), "guides"))
				if err != nil {
					return fmt.Errorf("unable to copy guides directory %q: %w", path, err)
				}
				return filepath.SkipDir
			}
		} else {
			switch {
			case regexp.MustCompile(`index.*`).MatchString(d.Name()): //index file
				m.infof("migrating provider index: %s", d.Name())
				err := filepath.WalkDir(path, m.MigrateTemplate(""))
				if err != nil {
					return err
				}
				return nil
			default:
				//skip non-index files
				return nil
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to migrate website: %w", err)
	}

	//remove legacy website directory
	err = os.RemoveAll(filepath.Join(m.providerDir, "website"))
	if err != nil {
		return fmt.Errorf("unable to remove legacy website directory: %w", err)
	}

	return nil
}

func (m *migrator) MigrateTemplate(relDir string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			//skip processing directories
			return nil
		}

		m.infof("migrating file %q", d.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read file %q: %w", d.Name(), err)
		}

		baseName, _, _ := strings.Cut(d.Name(), ".")
		shortName := providerShortName(m.providerName)
		fileName := strings.TrimPrefix(baseName, shortName+"_")

		var exampleRelDir string
		if fileName == "index" {
			exampleRelDir = relDir
		} else {
			exampleRelDir = filepath.Join(relDir, fileName)
		}
		templateFilePath := filepath.Join(m.ProviderTemplatesDir(), relDir, fileName+".md.tmpl")

		err = os.MkdirAll(filepath.Dir(templateFilePath), 0755)
		if err != nil {
			return fmt.Errorf("unable to create directory %q: %w", templateFilePath, err)
		}

		templateFile, err := os.OpenFile(templateFilePath, os.O_WRONLY|os.O_CREATE, 0600)

		if err != nil {
			return fmt.Errorf("unable to open file %q: %w", templateFilePath, err)
		}

		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				m.warnf("unable to close file %q: %q", f.Name(), err)
			}
		}(templateFile)

		m.infof("extracting YAML frontmatter to %q", templateFilePath)
		err = m.ExtractFrontMatter(data, relDir, templateFile)
		if err != nil {
			return fmt.Errorf("unable to extract front matter to %q: %w", templateFilePath, err)
		}

		m.infof("extracting code examples from %q", d.Name())
		err = m.ExtractCodeExamples(data, exampleRelDir, templateFile)
		if err != nil {
			return fmt.Errorf("unable to extract code examples from %q: %w", templateFilePath, err)
		}

		return nil
	}

}

func (m *migrator) ExtractFrontMatter(content []byte, relDir string, templateFile *os.File) error {
	fileScanner := bufio.NewScanner(bytes.NewReader(content))
	fileScanner.Split(bufio.ScanLines)

	hasFirstLine := fileScanner.Scan()
	if !hasFirstLine || fileScanner.Text() != "---" {
		m.warnf("no frontmatter found in %q", templateFile.Name())
		return nil
	}
	_, err := templateFile.WriteString(fileScanner.Text() + "\n")
	if err != nil {
		return fmt.Errorf("unable to append frontmatter to %q: %w", templateFile.Name(), err)
	}
	exited := false
	for fileScanner.Scan() {
		if strings.Contains(fileScanner.Text(), "layout:") {
			// skip layout front matter
			continue
		}
		_, err = templateFile.WriteString(fileScanner.Text() + "\n")
		if err != nil {
			return fmt.Errorf("unable to append frontmatter to %q: %w", templateFile.Name(), err)
		}
		if fileScanner.Text() == "---" {
			exited = true
			break
		}
	}

	if !exited {
		return fmt.Errorf("cannot find ending of frontmatter block in %q", templateFile.Name())
	}

	// add comment to end of front matter briefly explaining template functionality
	if relDir == "functions" {
		_, err = templateFile.WriteString(migrateFunctionTemplateComment + "\n")
	} else {
		_, err = templateFile.WriteString(migrateProviderTemplateComment + "\n")
	}
	if err != nil {
		return fmt.Errorf("unable to append template comment to %q: %w", templateFile.Name(), err)
	}

	return nil
}

func (m *migrator) ExtractCodeExamples(content []byte, newRelDir string, templateFile *os.File) error {
	md := newMarkdownRenderer()
	p := md.Parser()
	root := p.Parse(text.NewReader(content))

	exampleCount := 0
	importCount := 0

	err := ast.Walk(root, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		// skip the root node
		if !enter || node.Type() == ast.TypeDocument {
			return ast.WalkContinue, nil
		}

		if fencedNode, isFenced := node.(*ast.FencedCodeBlock); isFenced && fencedNode.Info != nil {
			var ext, exampleName, examplePath, template string

			lang := string(fencedNode.Info.Text(content)[:])
			switch lang {
			case "hcl", "terraform":
				exampleCount++
				ext = ".tf"
				exampleName = "example_" + strconv.Itoa(exampleCount) + ext
				examplePath = filepath.Join(m.examplesDir, newRelDir, exampleName)
				template = fmt.Sprintf("{{tffile \"%s\"}}", examplePath)
				m.infof("creating example file %q", filepath.Join(m.providerDir, examplePath))
			case "console":
				importCount++
				ext = ".sh"
				exampleName = "import_" + strconv.Itoa(importCount) + ext
				examplePath = filepath.Join(m.examplesDir, newRelDir, exampleName)
				template = fmt.Sprintf("{{codefile \"shell\" \"%s\"}}", examplePath)
				m.infof("creating import file %q", filepath.Join(m.providerDir, examplePath))
			default:
				// Render node as is
				m.infof("skipping code block with unknown language %q", lang)
				err := md.Renderer().Render(templateFile, content, node)
				if err != nil {
					return ast.WalkStop, fmt.Errorf("unable to render node: %w", err)
				}
				return ast.WalkSkipChildren, nil
			}

			// add code block text to buffer
			codeBuf := bytes.Buffer{}
			for i := 0; i < node.Lines().Len(); i++ {
				line := node.Lines().At(i)
				_, _ = codeBuf.Write(line.Value(content))
			}

			// create example file from code block
			err := writeFile(examplePath, codeBuf.String())
			if err != nil {
				return ast.WalkStop, fmt.Errorf("unable to write file %q: %w", examplePath, err)
			}

			// replace original code block with tfplugindocs template
			_, err = templateFile.WriteString("\n\n" + template)
			if err != nil {
				return ast.WalkStop, fmt.Errorf("unable to write to template %q: %w", template, err)
			}

			return ast.WalkSkipChildren, nil
		}

		// Render non-code nodes as is
		err := md.Renderer().Render(templateFile, content, node)
		if err != nil {
			return ast.WalkStop, fmt.Errorf("unable to render node: %w", err)
		}
		if node.HasChildren() {
			return ast.WalkSkipChildren, nil
		}

		return ast.WalkContinue, nil
	})
	if err != nil {
		return fmt.Errorf("unable to walk AST: %w", err)
	}

	_, err = templateFile.WriteString("\n")
	if err != nil {
		return fmt.Errorf("unable to write to template %q: %w", templateFile.Name(), err)
	}
	m.infof("finished creating template %q", templateFile.Name())

	return nil
}

// ProviderWebsiteDir returns the absolute path to the joined provider and
// the website directory that templates will be migrated from, which defaults to either "website/docs/" or "docs".
func (m *migrator) ProviderWebsiteDir() string {
	return filepath.Join(m.providerDir, m.websiteDir)
}

// ProviderTemplatesDir returns the absolute path to the joined provider and
// given new templates directory, which defaults to "templates".
func (m *migrator) ProviderTemplatesDir() string {
	return filepath.Join(m.providerDir, m.templatesDir)
}

// ProviderExamplesDir returns the absolute path to the joined provider and
// given examples directory, which defaults to "examples".
func (m *migrator) ProviderExamplesDir() string {
	return filepath.Join(m.providerDir, m.examplesDir)
}

func determineWebsiteDir(providerDir string) (string, error) {
	// Check for legacy website directory
	providerWebsiteDirFileInfo, err := os.Stat(filepath.Join(providerDir, "website/docs"))

	if err != nil {
		if os.IsNotExist(err) {
			// Legacy website directory does not exist, check for docs directory
		} else {
			return "", fmt.Errorf("error getting information for provider website directory %q: %w", providerDir, err)
		}
	} else if providerWebsiteDirFileInfo.IsDir() {
		return "website/docs", nil
	}

	// Check for docs directory
	providerDocsDirFileInfo, err := os.Stat(filepath.Join(providerDir, "docs"))

	if err != nil {
		return "", fmt.Errorf("error getting information for provider docs directory %q: %w", providerDir, err)
	}

	if providerDocsDirFileInfo.IsDir() {
		return "docs", nil
	}

	return "", fmt.Errorf("unable to determine website directory for provider %q", providerDir)

}

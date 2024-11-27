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
	"strings"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/go-version"
	install "github.com/hashicorp/hc-install"
	"github.com/hashicorp/hc-install/checkpoint"
	"github.com/hashicorp/hc-install/fs"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/hc-install/src"
	"github.com/hashicorp/terraform-exec/tfexec"
	tfjson "github.com/hashicorp/terraform-json"
	"golang.org/x/exp/slices"
)

var (
	websiteResourceFile                 = "resources/%s.md.tmpl"
	websiteResourceFallbackFile         = "resources.md.tmpl"
	websiteResourceFileStaticCandidates = []string{
		"resources/%s.md",
		"resources/%s.markdown",
		"resources/%s.html.markdown",
		"resources/%s.html.md",
		"r/%s.markdown",
		"r/%s.md",
		"r/%s.html.markdown",
		"r/%s.html.md",
	}
	websiteDataSourceFile                 = "data-sources/%s.md.tmpl"
	websiteDataSourceFallbackFile         = "data-sources.md.tmpl"
	websiteDataSourceFileStaticCandidates = []string{
		"data-sources/%s.md",
		"data-sources/%s.markdown",
		"data-sources/%s.html.markdown",
		"data-sources/%s.html.md",
		"d/%s.markdown",
		"d/%s.md",
		"d/%s.html.markdown",
		"d/%s.html.md",
	}
	websiteFunctionFile                 = "functions/%s.md.tmpl"
	websiteFunctionFallbackFile         = "functions.md.tmpl"
	websiteFunctionFileStaticCandidates = []string{
		"functions/%s.md",
		"functions/%s.markdown",
		"functions/%s.html.markdown",
		"functions/%s.html.md",
	}
	websiteEphemeralResourceFile                 = "ephemeral-resources/%s.md.tmpl"
	websiteEphemeralResourceFallbackFile         = "ephemeral-resources.md.tmpl"
	websiteEphemeralResourceFileStaticCandidates = []string{
		"ephemeral-resources/%s.md",
		"ephemeral-resources/%s.markdown",
		"ephemeral-resources/%s.html.markdown",
		"ephemeral-resources/%s.html.md",
	}
	websiteProviderFile                 = "index.md.tmpl"
	websiteProviderFileStaticCandidates = []string{
		"index.markdown",
		"index.md",
		"index.html.markdown",
		"index.html.md",
	}

	managedWebsiteSubDirectories = []string{
		"data-sources",
		"guides",
		"resources",
		"functions",
		"ephemeral-resources",
	}

	managedWebsiteFiles = []string{
		"index.md",
	}
)

type generator struct {
	ignoreDeprecated bool
	tfVersion        string

	// providerDir is the absolute path to the root provider directory
	providerDir string

	providerName         string
	providersSchemaPath  string
	renderedProviderName string
	renderedWebsiteDir   string
	examplesDir          string
	templatesDir         string
	websiteTmpDir        string

	ui cli.Ui
}

func (g *generator) infof(format string, a ...interface{}) {
	g.ui.Info(fmt.Sprintf(format, a...))
}

func (g *generator) warnf(format string, a ...interface{}) {
	g.ui.Warn(fmt.Sprintf(format, a...))
}

func Generate(ui cli.Ui, providerDir, providerName, providersSchemaPath, renderedProviderName, renderedWebsiteDir, examplesDir, websiteTmpDir, templatesDir, tfVersion string, ignoreDeprecated bool) error {
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

	g := &generator{
		ignoreDeprecated: ignoreDeprecated,
		tfVersion:        tfVersion,

		providerDir:          providerDir,
		providerName:         providerName,
		providersSchemaPath:  providersSchemaPath,
		renderedProviderName: renderedProviderName,
		renderedWebsiteDir:   renderedWebsiteDir,
		examplesDir:          examplesDir,
		templatesDir:         templatesDir,
		websiteTmpDir:        websiteTmpDir,

		ui: ui,
	}

	ctx := context.Background()

	return g.Generate(ctx)
}

func (g *generator) Generate(ctx context.Context) error {
	var err error

	if g.providerName == "" {
		g.providerName = filepath.Base(g.providerDir)
	}

	if g.renderedProviderName == "" {
		g.renderedProviderName = g.providerName
	}

	g.infof("rendering website for provider %q (as %q)", g.providerName, g.renderedProviderName)

	switch {
	case g.websiteTmpDir == "":
		g.websiteTmpDir, err = os.MkdirTemp("", "tfws")
		if err != nil {
			return fmt.Errorf("error creating temporary website directory: %w", err)
		}
		defer os.RemoveAll(g.websiteTmpDir)
	default:
		g.infof("cleaning tmp dir %q", g.websiteTmpDir)
		err = os.RemoveAll(g.websiteTmpDir)
		if err != nil {
			return fmt.Errorf("error removing temporary website directory %q: %w", g.websiteTmpDir, err)
		}

		g.infof("creating tmp dir %q", g.websiteTmpDir)
		err = os.MkdirAll(g.websiteTmpDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating temporary website directory %q: %w", g.websiteTmpDir, err)
		}
	}

	templatesDirInfo, err := os.Stat(g.ProviderTemplatesDir())
	switch {
	case os.IsNotExist(err):
		// do nothing, no template dir
	case err != nil:
		return fmt.Errorf("error getting information for provider templates directory %q: %w", g.ProviderTemplatesDir(), err)
	default:
		if !templatesDirInfo.IsDir() {
			return fmt.Errorf("template path is not a directory: %s", g.ProviderTemplatesDir())
		}

		g.infof("copying any existing content to tmp dir")
		err = cp(g.ProviderTemplatesDir(), g.TempTemplatesDir())
		if err != nil {
			return fmt.Errorf("error copying exiting content to temporary directory %q: %w", g.TempTemplatesDir(), err)
		}
	}

	var providerSchema *tfjson.ProviderSchema

	if g.providersSchemaPath == "" {
		g.infof("exporting schema from Terraform")
		providerSchema, err = g.terraformProviderSchemaFromTerraform(ctx)
		if err != nil {
			return fmt.Errorf("error exporting provider schema from Terraform: %w", err)
		}
	} else {
		g.infof("exporting schema from JSON file")
		providerSchema, err = g.terraformProviderSchemaFromFile()
		if err != nil {
			return fmt.Errorf("error exporting provider schema from JSON file: %w", err)
		}
	}

	g.infof("generating missing templates")
	err = g.generateMissingTemplates(providerSchema)
	if err != nil {
		return fmt.Errorf("error generating missing templates: %w", err)
	}

	g.infof("rendering static website")
	err = g.renderStaticWebsite(providerSchema)
	if err != nil {
		return fmt.Errorf("error rendering static website: %w", err)
	}

	return nil
}

// ProviderDocsDir returns the absolute path to the joined provider and
// given website documentation directory, which defaults to "docs".
func (g *generator) ProviderDocsDir() string {
	return filepath.Join(g.providerDir, g.renderedWebsiteDir)
}

// ProviderExamplesDir returns the absolute path to the joined provider and
// given examples directory, which defaults to "examples".
func (g *generator) ProviderExamplesDir() string {
	return filepath.Join(g.providerDir, g.examplesDir)
}

// ProviderTemplatesDir returns the absolute path to the joined provider and
// given templates directory, which defaults to "templates".
func (g *generator) ProviderTemplatesDir() string {
	return filepath.Join(g.providerDir, g.templatesDir)
}

// TempTemplatesDir returns the absolute path to the joined temporary and
// hardcoded "templates" subdirectory, which is where provider templates are
// copied.
func (g *generator) TempTemplatesDir() string {
	return filepath.Join(g.websiteTmpDir, "templates")
}

func (g *generator) generateMissingResourceTemplate(resourceName string) error {
	templatePath := fmt.Sprintf(websiteResourceFile, resourceShortName(resourceName, g.providerName))
	templatePath = filepath.Join(g.TempTemplatesDir(), templatePath)
	if fileExists(templatePath) {
		g.infof("resource %q template exists, skipping", resourceName)
		return nil
	}

	fallbackTemplatePath := filepath.Join(g.TempTemplatesDir(), websiteResourceFallbackFile)
	if fileExists(fallbackTemplatePath) {
		g.infof("resource %q fallback template exists, creating template", resourceName)
		err := cp(fallbackTemplatePath, templatePath)
		if err != nil {
			return fmt.Errorf("unable to copy fallback template for %q: %w", resourceName, err)
		}
		return nil
	}

	for _, candidate := range websiteResourceFileStaticCandidates {
		candidatePath := fmt.Sprintf(candidate, resourceShortName(resourceName, g.providerName))
		candidatePath = filepath.Join(g.TempTemplatesDir(), candidatePath)
		if fileExists(candidatePath) {
			g.infof("resource %q static file exists, skipping", resourceName)
			return nil
		}
	}

	g.infof("generating new template for %q", resourceName)
	err := writeFile(templatePath, string(defaultResourceTemplate))
	if err != nil {
		return fmt.Errorf("unable to write template for %q: %w", resourceName, err)
	}

	return nil
}

func (g *generator) generateMissingDataSourceTemplate(datasourceName string) error {
	templatePath := fmt.Sprintf(websiteDataSourceFile, resourceShortName(datasourceName, g.providerName))
	templatePath = filepath.Join(g.TempTemplatesDir(), templatePath)
	if fileExists(templatePath) {
		g.infof("data-source %q template exists, skipping", datasourceName)
		return nil
	}

	fallbackTemplatePath := filepath.Join(g.TempTemplatesDir(), websiteDataSourceFallbackFile)
	if fileExists(fallbackTemplatePath) {
		g.infof("data-source %q fallback template exists, creating template", datasourceName)
		err := cp(fallbackTemplatePath, templatePath)
		if err != nil {
			return fmt.Errorf("unable to copy fallback template for %q: %w", datasourceName, err)
		}
		return nil
	}

	for _, candidate := range websiteDataSourceFileStaticCandidates {
		candidatePath := fmt.Sprintf(candidate, resourceShortName(datasourceName, g.providerName))
		candidatePath = filepath.Join(g.TempTemplatesDir(), candidatePath)
		if fileExists(candidatePath) {
			g.infof("data-source %q static file exists, skipping", datasourceName)
			return nil
		}
	}

	g.infof("generating new template for data-source %q", datasourceName)
	err := writeFile(templatePath, string(defaultResourceTemplate))
	if err != nil {
		return fmt.Errorf("unable to write template for %q: %w", datasourceName, err)
	}

	return nil
}

func (g *generator) generateMissingFunctionTemplate(functionName string) error {
	templatePath := fmt.Sprintf(websiteFunctionFile, resourceShortName(functionName, g.providerName))
	templatePath = filepath.Join(g.TempTemplatesDir(), templatePath)
	if fileExists(templatePath) {
		g.infof("function %q template exists, skipping", functionName)
		return nil
	}

	fallbackTemplatePath := filepath.Join(g.TempTemplatesDir(), websiteFunctionFallbackFile)
	if fileExists(fallbackTemplatePath) {
		g.infof("function %q fallback template exists, creating template", functionName)
		err := cp(fallbackTemplatePath, templatePath)
		if err != nil {
			return fmt.Errorf("unable to copy fallback template for %q: %w", functionName, err)
		}
		return nil
	}

	for _, candidate := range websiteFunctionFileStaticCandidates {
		candidatePath := fmt.Sprintf(candidate, resourceShortName(functionName, g.providerName))
		candidatePath = filepath.Join(g.TempTemplatesDir(), candidatePath)
		if fileExists(candidatePath) {
			g.infof("function %q static file exists, skipping", functionName)
			return nil
		}
	}

	g.infof("generating new template for function %q", functionName)
	err := writeFile(templatePath, string(defaultFunctionTemplate))
	if err != nil {
		return fmt.Errorf("unable to write template for %q: %w", functionName, err)
	}

	return nil
}

func (g *generator) generateMissingEphemeralResourceTemplate(resourceName string) error {
	templatePath := fmt.Sprintf(websiteEphemeralResourceFile, resourceShortName(resourceName, g.providerName))
	templatePath = filepath.Join(g.TempTemplatesDir(), templatePath)
	if fileExists(templatePath) {
		g.infof("ephemeral resource %q template exists, skipping", resourceName)
		return nil
	}

	fallbackTemplatePath := filepath.Join(g.TempTemplatesDir(), websiteEphemeralResourceFallbackFile)
	if fileExists(fallbackTemplatePath) {
		g.infof("ephemeral resource %q fallback template exists, creating template", resourceName)
		err := cp(fallbackTemplatePath, templatePath)
		if err != nil {
			return fmt.Errorf("unable to copy fallback template for %q: %w", resourceName, err)
		}
		return nil
	}

	for _, candidate := range websiteEphemeralResourceFileStaticCandidates {
		candidatePath := fmt.Sprintf(candidate, resourceShortName(resourceName, g.providerName))
		candidatePath = filepath.Join(g.TempTemplatesDir(), candidatePath)
		if fileExists(candidatePath) {
			g.infof("ephemeral resource %q static file exists, skipping", resourceName)
			return nil
		}
	}

	g.infof("generating new template for %q", resourceName)
	err := writeFile(templatePath, string(defaultResourceTemplate))
	if err != nil {
		return fmt.Errorf("unable to write template for %q: %w", resourceName, err)
	}

	return nil
}

func (g *generator) generateMissingProviderTemplate() error {
	templatePath := filepath.Join(g.TempTemplatesDir(), websiteProviderFile)
	if fileExists(templatePath) {
		g.infof("provider %q template exists, skipping", g.providerName)
		return nil
	}

	for _, candidate := range websiteProviderFileStaticCandidates {
		candidatePath := filepath.Join(g.TempTemplatesDir(), candidate)
		if fileExists(candidatePath) {
			g.infof("provider %q static file exists, skipping", g.providerName)
			return nil
		}
	}

	g.infof("generating new template for %q", g.providerName)
	err := writeFile(templatePath, string(defaultProviderTemplate))
	if err != nil {
		return fmt.Errorf("unable to write template for %q: %w", g.providerName, err)
	}

	return nil
}

func (g *generator) generateMissingTemplates(providerSchema *tfjson.ProviderSchema) error {
	g.infof("generating missing resource content")
	for name, schema := range providerSchema.ResourceSchemas {
		if g.ignoreDeprecated && schema.Block.Deprecated {
			continue
		}

		err := g.generateMissingResourceTemplate(name)
		if err != nil {
			return fmt.Errorf("unable to generate template for resource %q: %w", name, err)
		}
	}

	g.infof("generating missing data source content")
	for name, schema := range providerSchema.DataSourceSchemas {
		if g.ignoreDeprecated && schema.Block.Deprecated {
			continue
		}

		err := g.generateMissingDataSourceTemplate(name)
		if err != nil {
			return fmt.Errorf("unable to generate template for data-source %q: %w", name, err)
		}
	}

	g.infof("generating missing function content")
	for name, signature := range providerSchema.Functions {
		if g.ignoreDeprecated && signature.DeprecationMessage != "" {
			continue
		}

		err := g.generateMissingFunctionTemplate(name)
		if err != nil {
			return fmt.Errorf("unable to generate template for function %q: %w", name, err)
		}
	}

	g.infof("generating missing ephemeral resource content")
	for name, schema := range providerSchema.EphemeralResourceSchemas {
		if g.ignoreDeprecated && schema.Block.Deprecated {
			continue
		}

		err := g.generateMissingEphemeralResourceTemplate(name)
		if err != nil {
			return fmt.Errorf("unable to generate template for ephemeral resource %q: %w", name, err)
		}
	}

	g.infof("generating missing provider content")
	err := g.generateMissingProviderTemplate()
	if err != nil {
		return fmt.Errorf("unable to generate template for provider: %w", err)
	}

	return nil
}

func (g *generator) renderStaticWebsite(providerSchema *tfjson.ProviderSchema) error {
	g.infof("cleaning rendered website dir")
	dirEntry, err := os.ReadDir(g.ProviderDocsDir())
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to read rendered website directory %q: %w", g.ProviderDocsDir(), err)
	}

	for _, file := range dirEntry {

		// Remove subdirectories managed by tfplugindocs
		if file.IsDir() && slices.Contains(managedWebsiteSubDirectories, file.Name()) {
			g.infof("removing directory: %q", file.Name())
			err = os.RemoveAll(filepath.Join(g.ProviderDocsDir(), file.Name()))
			if err != nil {
				return fmt.Errorf("unable to remove directory %q from rendered website directory: %w", file.Name(), err)
			}
			continue
		}

		// Remove files managed by tfplugindocs
		if !file.IsDir() && slices.Contains(managedWebsiteFiles, file.Name()) {
			g.infof("removing file: %q", file.Name())
			err = os.RemoveAll(filepath.Join(g.ProviderDocsDir(), file.Name()))
			if err != nil {
				return fmt.Errorf("unable to remove file %q from rendered website directory: %w", file.Name(), err)
			}
			continue
		}
	}

	shortName := providerShortName(g.providerName)

	g.infof("rendering templated website to static markdown")

	err = filepath.WalkDir(g.websiteTmpDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("unable to walk path %q: %w", path, err)
		}
		if d.IsDir() {
			// skip directories
			return nil
		}

		rel, err := filepath.Rel(filepath.Join(g.TempTemplatesDir()), path)
		if err != nil {
			return fmt.Errorf("unable to retrieve the relative path of basepath %q and targetpath %q: %w",
				filepath.Join(g.TempTemplatesDir()), path, err)
		}

		relDir, relFile := filepath.Split(rel)
		relDir = filepath.ToSlash(relDir)

		// skip special top-level generic resource, data source, function, and ephemeral resource templates
		if relDir == "" && (relFile == "resources.md.tmpl" ||
			relFile == "data-sources.md.tmpl" ||
			relFile == "functions.md.tmpl" ||
			relFile == "ephemeral-resources.md.tmpl") {
			return nil
		}

		renderedPath := filepath.Join(g.ProviderDocsDir(), rel)
		err = os.MkdirAll(filepath.Dir(renderedPath), 0755)
		if err != nil {
			return fmt.Errorf("unable to create rendered website subdirectory %q: %w", renderedPath, err)
		}

		ext := filepath.Ext(path)
		if ext != ".tmpl" {
			g.infof("copying non-template file: %q", rel)
			return cp(path, renderedPath)
		}

		renderedPath = strings.TrimSuffix(renderedPath, ext)

		tmplData, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read file %q: %w", rel, err)
		}

		out, err := os.Create(renderedPath)
		if err != nil {
			return fmt.Errorf("unable to create file %q: %w", renderedPath, err)
		}
		defer out.Close()

		g.infof("rendering %q", rel)
		switch relDir {
		case "data-sources/":
			resSchema, resName := resourceSchema(providerSchema.DataSourceSchemas, shortName, relFile)
			exampleFilePath := filepath.Join(g.ProviderExamplesDir(), "data-sources", resName, "data-source.tf")

			if resSchema != nil {
				tmpl := resourceTemplate(tmplData)
				render, err := tmpl.Render(g.providerDir, resName, g.providerName, g.renderedProviderName, "Data Source", exampleFilePath, "", resSchema)
				if err != nil {
					return fmt.Errorf("unable to render data source template %q: %w", rel, err)
				}
				_, err = out.WriteString(render)
				if err != nil {
					return fmt.Errorf("unable to write rendered string: %w", err)
				}
				return nil
			}
			g.warnf("data source entitled %q, or %q does not exist", shortName, resName)
		case "resources/":
			resSchema, resName := resourceSchema(providerSchema.ResourceSchemas, shortName, relFile)
			exampleFilePath := filepath.Join(g.ProviderExamplesDir(), "resources", resName, "resource.tf")
			importFilePath := filepath.Join(g.ProviderExamplesDir(), "resources", resName, "import.sh")

			if resSchema != nil {
				tmpl := resourceTemplate(tmplData)
				render, err := tmpl.Render(g.providerDir, resName, g.providerName, g.renderedProviderName, "Resource", exampleFilePath, importFilePath, resSchema)
				if err != nil {
					return fmt.Errorf("unable to render resource template %q: %w", rel, err)
				}
				_, err = out.WriteString(render)
				if err != nil {
					return fmt.Errorf("unable to write rendered string: %w", err)
				}
				return nil
			}
			g.warnf("resource entitled %q, or %q does not exist", shortName, resName)
		case "functions/":
			funcName := removeAllExt(relFile)
			if signature, ok := providerSchema.Functions[funcName]; ok {
				exampleFilePath := filepath.Join(g.ProviderExamplesDir(), "functions", funcName, "function.tf")

				tmpl := functionTemplate(tmplData)
				render, err := tmpl.Render(g.providerDir, funcName, g.providerName, g.renderedProviderName, "function", exampleFilePath, signature)
				if err != nil {
					return fmt.Errorf("unable to render function template %q: %w", rel, err)
				}
				_, err = out.WriteString(render)
				if err != nil {
					return fmt.Errorf("unable to write rendered string: %w", err)
				}
				return nil
			}

			g.warnf("function entitled %q does not exist", funcName)
		case "ephemeral-resources/":
			resSchema, resName := resourceSchema(providerSchema.EphemeralResourceSchemas, shortName, relFile)
			exampleFilePath := filepath.Join(g.ProviderExamplesDir(), "ephemeral-resources", resName, "ephemeral-resource.tf")

			if resSchema != nil {
				tmpl := resourceTemplate(tmplData)
				render, err := tmpl.Render(g.providerDir, resName, g.providerName, g.renderedProviderName, "Ephemeral Resource", exampleFilePath, "", resSchema)
				if err != nil {
					return fmt.Errorf("unable to render ephemeral resource template %q: %w", rel, err)
				}
				_, err = out.WriteString(render)
				if err != nil {
					return fmt.Errorf("unable to write rendered string: %w", err)
				}
				return nil
			}
			g.warnf("ephemeral resource entitled %q, or %q does not exist", shortName, resName)
		case "": // provider
			if relFile == "index.md.tmpl" {
				tmpl := providerTemplate(tmplData)
				exampleFilePath := filepath.Join(g.ProviderExamplesDir(), "provider", "provider.tf")
				render, err := tmpl.Render(g.providerDir, g.providerName, g.renderedProviderName, exampleFilePath, providerSchema.ConfigSchema)
				if err != nil {
					return fmt.Errorf("unable to render provider template %q: %w", rel, err)
				}
				_, err = out.WriteString(render)
				if err != nil {
					return fmt.Errorf("unable to write rendered string: %w", err)
				}
				return nil
			}
		}

		tmpl := docTemplate(tmplData)
		err = tmpl.Render(g.providerDir, out)
		if err != nil {
			return fmt.Errorf("unable to render template %q: %w", rel, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("unable to render templated website to static markdown: %w", err)
	}

	return nil
}

func (g *generator) terraformProviderSchemaFromTerraform(ctx context.Context) (*tfjson.ProviderSchema, error) {
	var err error

	shortName := providerShortName(g.providerName)

	tmpDir, err := os.MkdirTemp("", "tfws")
	if err != nil {
		return nil, fmt.Errorf("unable to create temporary provider install directory %q: %w", tmpDir, err)
	}
	defer os.RemoveAll(tmpDir)

	g.infof("compiling provider %q", shortName)
	providerPath := fmt.Sprintf("plugins/registry.terraform.io/hashicorp/%s/0.0.1/%s_%s", shortName, runtime.GOOS, runtime.GOARCH)
	outFile := filepath.Join(tmpDir, providerPath, fmt.Sprintf("terraform-provider-%s", shortName))
	switch runtime.GOOS {
	case "windows":
		outFile = outFile + ".exe"
	}
	buildCmd := exec.Command("go", "build", "-o", outFile)
	buildCmd.Dir = g.providerDir
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
	if g.tfVersion != "" {
		g.infof("downloading Terraform CLI binary version from releases.hashicorp.com: %s", g.tfVersion)
		sources = []src.Source{
			&releases.ExactVersion{
				Product:    product.Terraform,
				Version:    version.Must(version.NewVersion(g.tfVersion)),
				InstallDir: tmpDir,
			},
		}
	} else {
		g.infof("using Terraform CLI binary from PATH if available, otherwise downloading latest Terraform CLI binary")
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

	g.infof("running terraform init")
	err = tf.Init(ctx, tfexec.Get(false), tfexec.PluginDir("./plugins"))
	if err != nil {
		return nil, fmt.Errorf("unable to run terraform init on provider: %w", err)
	}

	g.infof("getting provider schema")
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

func (g *generator) terraformProviderSchemaFromFile() (*tfjson.ProviderSchema, error) {
	var err error

	shortName := providerShortName(g.providerName)

	g.infof("getting provider schema")
	schemas, err := extractSchemaFromFile(g.providersSchemaPath)
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

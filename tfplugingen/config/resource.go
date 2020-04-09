package config

import (
	"bytes"
	"fmt"
	"html/template"
)

type ResourceTemplateInfo struct {
	Name string
	Type string
}

type ResourceTemplate string

func (t ResourceTemplate) Render(res ResourceTemplateInfo) (string, error) {
	s := string(t)
	if s == "" {
		return "", nil
	}

	tmpl, err := template.New("ResourceTemplate").Parse(s)
	if err != nil {
		return "", fmt.Errorf("unable to parse template %q: %w", s, err)
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, res)
	if err != nil {
		return "", fmt.Errorf("unable to execute template: %w", err)
	}

	return buf.String(), nil
}

func (t ResourceTemplate) MustRender(res ResourceTemplateInfo) string {
	s, err := t.Render(res)
	if err != nil {
		panic(err)
	}
	return s
}

type Resource struct {
	Name string `hcl:",label"`

	Package string `hcl:"package,optional"`
	Type    string `hcl:"type,optional"` // TODO: infer default from name?

	ReadFunc   string `hcl:"read_func,optional"`
	CreateFunc string `hcl:"create_func,optional"`
	UpdateFunc string `hcl:"update_func,optional"`
	DeleteFunc string `hcl:"delete_func,optional"`

	*SchemaProperties
}

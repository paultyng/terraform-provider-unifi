// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package check

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

type FrontMatterCheck struct {
	Options *FrontMatterOptions
}

// FrontMatterData represents the YAML frontmatter of Terraform Provider documentation.
type FrontMatterData struct {
	Description    *string `yaml:"description,omitempty"`
	Layout         *string `yaml:"layout,omitempty"`
	PageTitle      *string `yaml:"page_title,omitempty"`
	SidebarCurrent *string `yaml:"sidebar_current,omitempty"`
	Subcategory    *string `yaml:"subcategory,omitempty"`
}

// FrontMatterOptions represents configuration options for FrontMatter.
type FrontMatterOptions struct {
	NoLayout           bool
	NoPageTitle        bool
	NoSidebarCurrent   bool
	NoSubcategory      bool
	RequireDescription bool
	RequireLayout      bool
	RequirePageTitle   bool
}

func NewFrontMatterCheck(opts *FrontMatterOptions) *FrontMatterCheck {
	check := &FrontMatterCheck{
		Options: opts,
	}

	if check.Options == nil {
		check.Options = &FrontMatterOptions{}
	}

	return check
}

func (check *FrontMatterCheck) Run(src []byte) error {
	frontMatter := FrontMatterData{}

	md := goldmark.New(
		goldmark.WithExtensions(&frontmatter.Extender{}),
	)

	ctx := parser.NewContext()
	var buff bytes.Buffer

	err := md.Convert(src, &buff, parser.WithContext(ctx))
	if err != nil {
		return err
	}
	d := frontmatter.Get(ctx)
	if d == nil {
		return fmt.Errorf("no frontmatter found")
	}

	err = d.Decode(&frontMatter)
	if err != nil {
		return fmt.Errorf("error parsing YAML frontmatter: %w", err)
	}

	if check.Options.NoLayout && frontMatter.Layout != nil {
		return fmt.Errorf("YAML frontmatter should not contain layout")
	}

	if check.Options.NoPageTitle && frontMatter.PageTitle != nil {
		return fmt.Errorf("YAML frontmatter should not contain page_title")
	}

	if check.Options.NoSidebarCurrent && frontMatter.SidebarCurrent != nil {
		return fmt.Errorf("YAML frontmatter should not contain sidebar_current")
	}

	if check.Options.NoSubcategory && frontMatter.Subcategory != nil {
		return fmt.Errorf("YAML frontmatter should not contain subcategory")
	}

	if check.Options.RequireDescription && frontMatter.Description == nil {
		return fmt.Errorf("YAML frontmatter missing required description")
	}

	if check.Options.RequireLayout && frontMatter.Layout == nil {
		return fmt.Errorf("YAML frontmatter missing required layout")
	}

	if check.Options.RequirePageTitle && frontMatter.PageTitle == nil {
		return fmt.Errorf("YAML frontmatter missing required page_title")
	}

	return nil
}

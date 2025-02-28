// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mdplain

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// Clean runs a VERY naive cleanup of markdown text to make it more palatable as plain text.
func PlainMarkdown(markdown string) (string, error) {
	var buf bytes.Buffer
	extensions := []goldmark.Extender{
		extension.Linkify,
	}
	md := goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithRenderer(NewTextRenderer()),
	)
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

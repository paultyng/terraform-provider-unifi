// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mdplain

import (
	"bytes"
	"io"

	"github.com/yuin/goldmark/ast"
	extAST "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
)

type TextRender struct{}

func NewTextRenderer() *TextRender {
	return &TextRender{}
}

func (r *TextRender) Render(w io.Writer, source []byte, n ast.Node) error {
	out := bytes.NewBuffer([]byte{})
	err := ast.Walk(n, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering || node.Type() == ast.TypeDocument {
			return ast.WalkContinue, nil
		}

		switch node := node.(type) {
		case *ast.Blockquote, *ast.Heading:
			doubleSpace(out)
			out.Write(node.Text(source))
			return ast.WalkSkipChildren, nil
		case *ast.ThematicBreak:
			doubleSpace(out)
			return ast.WalkSkipChildren, nil
		case *ast.CodeBlock:
			doubleSpace(out)
			for i := 0; i < node.Lines().Len(); i++ {
				line := node.Lines().At(i)
				out.Write(line.Value(source))
			}
			return ast.WalkSkipChildren, nil
		case *ast.FencedCodeBlock:
			doubleSpace(out)
			doubleSpace(out)
			for i := 0; i < node.Lines().Len(); i++ {
				line := node.Lines().At(i)
				_, _ = out.Write(line.Value(source))
			}
			return ast.WalkSkipChildren, nil
		case *ast.List:
			doubleSpace(out)
			return ast.WalkContinue, nil
		case *ast.Paragraph:
			doubleSpace(out)
			if node.Text(source)[0] == '|' { // Write tables as-is.
				for i := 0; i < node.Lines().Len(); i++ {
					line := node.Lines().At(i)
					out.Write(line.Value(source))
				}
				return ast.WalkSkipChildren, nil
			}
			return ast.WalkContinue, nil
		case *extAST.Strikethrough:
			out.Write(node.Text(source))
			return ast.WalkContinue, nil
		case *ast.AutoLink:
			out.Write(node.URL(source))
			return ast.WalkSkipChildren, nil
		case *ast.CodeSpan:
			out.Write(node.Text(source))
			return ast.WalkSkipChildren, nil
		case *ast.Link:
			_, err := out.Write(node.Text(source))
			if !isRelativeLink(node.Destination) {
				out.WriteString(" ")
				out.Write(node.Destination)
			}
			return ast.WalkSkipChildren, err
		case *ast.Text:
			out.Write(node.Text(source))
			if node.SoftLineBreak() {
				doubleSpace(out)
			}
			return ast.WalkContinue, nil
		case *ast.Image:
			return ast.WalkSkipChildren, nil

		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return err
	}
	_, err = w.Write(out.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (r *TextRender) AddOptions(...renderer.Option) {}

func doubleSpace(out *bytes.Buffer) {
	if out.Len() > 0 {
		out.WriteByte('\n')
	}
}

func isRelativeLink(link []byte) (yes bool) {
	yes = false

	// a tag begin with '#'
	if link[0] == '#' {
		yes = true
	}

	// link begin with '/' but not '//', the second maybe a protocol relative link
	if len(link) >= 2 && link[0] == '/' && link[1] != '/' {
		yes = true
	}

	// only the root '/'
	if len(link) == 1 && link[0] == '/' {
		yes = true
	}
	return
}

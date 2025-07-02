package markdown

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strconv"
	"unicode/utf8"
	"unsafe"

	"github.com/yuin/goldmark/ast"
	extAST "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/text"
)

var (
	newLineChar             = []byte{'\n'}
	spaceChar               = []byte{' '}
	strikeThroughChars      = []byte("~~")
	thematicBreakChars      = []byte("---")
	blockquoteChars         = []byte{'>', ' '}
	codeBlockChars          = []byte("```")
	tableHeaderColChar      = []byte{'-'}
	tableHeaderAlignColChar = []byte{':'}
	heading1UnderlineChar   = []byte{'='}
	heading2UnderlineChar   = []byte{'-'}
	fourSpacesChars         = bytes.Repeat([]byte{' '}, 4)
)

// Ensure compatibility with Goldmark parser.
var _ renderer.Renderer = &Renderer{}

// Renderer allows to render markdown AST into markdown bytes in consistent format.
// Render is reusable across Renders, it holds configuration only.
type Renderer struct {
	underlineHeadings bool
	softWraps         bool
	emphToken         []byte
	strongToken       []byte // if nil, use emphToken*2
	listIndentStyle   ListIndentStyle

	// language name => format function
	formatters map[string]func([]byte) []byte
}

// AddOptions pulls Markdown renderer specific options from the given list,
// and applies them to the renderer.
func (mr *Renderer) AddOptions(opts ...renderer.Option) {
	mdopts := make([]Option, 0, len(opts))
	for _, o := range opts {
		if mo, ok := o.(Option); ok {
			mdopts = append(mdopts, mo)
		}
	}
	mr.AddMarkdownOptions(mdopts...)
}

// AddMarkdownOptions modifies the Renderer with the given options.
func (mr *Renderer) AddMarkdownOptions(opts ...Option) {
	for _, o := range opts {
		o.apply(mr)
	}
}

// Option customizes the behavior of the markdown renderer.
type Option interface {
	renderer.Option

	apply(r *Renderer)
}

type optionFunc func(*Renderer)

func (f optionFunc) SetConfig(*renderer.Config) {}

func (f optionFunc) apply(r *Renderer) {
	f(r)
}

// WithUnderlineHeadings configures the renderer to use
// Setext-style headers (=== and ---).
func WithUnderlineHeadings() Option {
	return optionFunc(func(r *Renderer) {
		r.underlineHeadings = true
	})
}

// WithSoftWraps allows you to wrap lines even on soft line breaks.
func WithSoftWraps() Option {
	return optionFunc(func(r *Renderer) {
		r.softWraps = true
	})
}

// WithEmphasisToken specifies the character used to wrap emphasised text.
// Per the CommonMark spec, valid values are '*' and '_'.
//
// Defaults to '*'.
func WithEmphasisToken(c rune) Option {
	return optionFunc(func(r *Renderer) {
		buf := make([]byte, 4) // enough to encode any utf8 rune
		n := utf8.EncodeRune(buf, c)
		r.emphToken = buf[:n]
	})
}

// WithStrongToken specifies the string used to wrap bold text.
// Per the CommonMark spec, valid values are '**' and '__'.
//
// Defaults to repeating the emphasis token twice.
// See [WithEmphasisToken] for how to change that.
func WithStrongToken(s string) Option {
	return optionFunc(func(r *Renderer) {
		r.strongToken = []byte(s)
	})
}

// ListIndentStyle specifies how items nested inside lists
// should be indented.
type ListIndentStyle int

const (
	// ListIndentAligned specifies that items inside a list item
	// should be aligned to the content in the first item.
	//
	//	- First paragraph.
	//
	//	  Second paragraph aligned with the first.
	//
	// This applies to ordered lists too.
	//
	//	1. First paragraph.
	//
	//	   Second paragraph aligned with the first.
	//
	//	...
	//
	//	10. Contents.
	//
	//	    Long lists indent content further.
	//
	// This is the default.
	ListIndentAligned ListIndentStyle = iota

	// ListIndentUniform specifies that items inside a list item
	// should be aligned uniformly with 4 spaces.
	//
	// For example:
	//
	//	- First paragraph.
	//
	//	    Second paragraph indented 4 spaces.
	//
	// For ordered lists:
	//
	//	1. First paragraph.
	//
	//	    Second paragraph indented 4 spaces.
	//
	//	...
	//
	//	10. Contents.
	//
	//	    Always indented 4 spaces.
	ListIndentUniform
)

// WithListIndentStyle specifies how contents nested under a list item
// should be indented.
//
// Defaults to [ListIndentAligned].
func WithListIndentStyle(style ListIndentStyle) Option {
	return optionFunc(func(r *Renderer) {
		r.listIndentStyle = style
	})
}

// CodeFormatter reformats code samples found in the document,
// matching them by name.
type CodeFormatter struct {
	// Name of the language.
	Name string

	// Aliases for the language, if any.
	Aliases []string

	// Function to format the code snippet.
	// In case of errors, format functions should typically return
	// the original string unchanged.
	Format func([]byte) []byte
}

// GoCodeFormatter is a [CodeFormatter] that reformats Go source code inside
// fenced code blocks tagged with 'go' or 'Go'.
//
//	```go
//	func main() {
//	}
//	```
//
// Supply it to the renderer with [WithCodeFormatters].
var GoCodeFormatter = CodeFormatter{
	Name:    "go",
	Aliases: []string{"Go"},
	Format:  formatGo,
}

func formatGo(src []byte) []byte {
	gofmt, err := format.Source(src)
	if err != nil {
		// We don't handle gofmt errors.
		// If code is not compilable we just
		// don't format it without any warning.
		return src
	}
	return gofmt
}

// WithCodeFormatters changes the functions used to reformat code blocks found
// in the original file.
//
//	formatters := []markdown.CodeFormatter{
//		markdown.GoCodeFormatter,
//		// ...
//	}
//	r := NewRenderer()
//	r.AddMarkdownOptions(WithCodeFormatters(formatters...))
//
// Defaults to empty.
func WithCodeFormatters(fs ...CodeFormatter) Option {
	return optionFunc(func(r *Renderer) {
		formatters := make(map[string]func([]byte) []byte, len(fs))
		for _, f := range fs {
			formatters[f.Name] = f.Format
			for _, alias := range f.Aliases {
				formatters[alias] = f.Format
			}
		}
		r.formatters = formatters
	})
}

// NewRenderer builds a new Markdown renderer with default settings.
// To use this with goldmark.Markdown, use the goldmark.WithRenderer option.
//
//	r := markdown.NewRenderer()
//	md := goldmark.New(goldmark.WithRenderer(r))
//	md.Convert(src, w)
//
// Alternatively, you can call [Renderer.Render] directly.
//
//	r := markdown.NewRenderer()
//	r.Render(w, src, node)
//
// Use [Renderer.AddMarkdownOptions] to customize the output of the renderer.
func NewRenderer() *Renderer {
	return &Renderer{
		emphToken: []byte{'*'},
		// Leave strongToken as nil by default.
		// At render time, we'll use what was specified,
		// or repeat emphToken twice to get the strong token.
	}
}

// render represents a single markdown rendering operation.
type render struct {
	mr *Renderer

	emphToken   []byte
	strongToken []byte

	// TODO(bwplotka): Wrap it with something that catch errors.
	w      *lineIndentWriter
	source []byte
}

func (mr *Renderer) newRender(w io.Writer, source []byte) *render {
	strongToken := mr.strongToken
	if len(strongToken) == 0 {
		strongToken = bytes.Repeat(mr.emphToken, 2)
	}

	return &render{
		mr:          mr,
		w:           wrapWithLineIndentWriter(w),
		source:      source,
		strongToken: strongToken,
		emphToken:   mr.emphToken,
	}
}

// Render renders the given AST node to the given writer,
// given the original source from which the node was parsed.
//
// NOTE: This is the entry point used by Goldmark.
func (mr *Renderer) Render(w io.Writer, source []byte, node ast.Node) error {
	// Perform DFS.
	return ast.Walk(node, mr.newRender(w, source).renderNode)
}

func (r *render) renderNode(node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering && node.PreviousSibling() != nil {
		switch node.(type) {
		// All Block types (except few) usually have 2x new lines before itself when they are non-first siblings.
		case *ast.Paragraph, *ast.Heading, *ast.FencedCodeBlock,
			*ast.CodeBlock, *ast.ThematicBreak, *extAST.Table,
			*ast.Blockquote:
			_, _ = r.w.Write(newLineChar)
			_, _ = r.w.Write(newLineChar)
		case *ast.List, *ast.HTMLBlock:
			_, _ = r.w.Write(newLineChar)
			if node.HasBlankPreviousLines() {
				_, _ = r.w.Write(newLineChar)
			}
		case *ast.ListItem:
			// TODO(bwplotka): Handle tight/loose rule explicitly.
			// See: https://github.github.com/gfm/#loose
			if node.HasBlankPreviousLines() {
				_, _ = r.w.Write(newLineChar)
			}
		}
	}

	switch tnode := node.(type) {
	case *ast.Document:
		if entering {
			break
		}

		_, _ = r.w.Write(newLineChar)

	// Spans, meaning no newlines before or after.
	case *ast.Text:
		if entering {
			text := tnode.Segment.Value(r.source)
			_ = writeClean(r.w, text)
			break
		}

		if tnode.SoftLineBreak() {
			char := spaceChar
			if r.mr.softWraps {
				char = newLineChar
			}
			_, _ = r.w.Write(char)
		}

		if tnode.HardLineBreak() {
			if tnode.SoftLineBreak() {
				_, _ = r.w.Write(spaceChar)
			}
			_, _ = r.w.Write(newLineChar)
		}
	case *ast.String:
		if entering {
			_, _ = r.w.Write(tnode.Value)
		}
	case *ast.AutoLink:
		// We treat autolink as normal string.
		if entering {
			_, _ = r.w.Write(tnode.Label(r.source))
		}
	case *extAST.TaskCheckBox:
		if !entering {
			break
		}
		if tnode.IsChecked {
			_, _ = r.w.Write([]byte("[X] "))
			break
		}
		_, _ = r.w.Write([]byte("[ ] "))
	case *ast.CodeSpan:
		if entering {
			_, _ = r.w.Write([]byte{'`'})
			break
		}

		_, _ = r.w.Write([]byte{'`'})
	case *extAST.Strikethrough:
		return r.wrapNonEmptyContentWith(strikeThroughChars, entering), nil
	case *ast.Emphasis:
		var emWrapper []byte
		switch tnode.Level {
		case 1:
			emWrapper = r.emphToken
		case 2:
			emWrapper = r.strongToken
		default:
			emWrapper = bytes.Repeat(r.emphToken, tnode.Level)
		}
		return r.wrapNonEmptyContentWith(emWrapper, entering), nil
	case *ast.Link:
		if entering {
			r.w.AddIndentOnFirstWrite([]byte("["))
			break
		}

		_, _ = fmt.Fprintf(r.w, "](%s", tnode.Destination)
		if len(tnode.Title) > 0 {
			_, _ = fmt.Fprintf(r.w, ` "%s"`, tnode.Title)
		}
		_, _ = r.w.Write([]byte{')'})
	case *ast.Image:
		if entering {
			r.w.AddIndentOnFirstWrite([]byte("!["))
			break
		}

		_, _ = fmt.Fprintf(r.w, "](%s", tnode.Destination)
		if len(tnode.Title) > 0 {
			_, _ = fmt.Fprintf(r.w, ` "%s"`, tnode.Title)
		}
		_, _ = r.w.Write([]byte{')'})
	case *ast.RawHTML:
		if !entering {
			break
		}

		for i := 0; i < tnode.Segments.Len(); i++ {
			segment := tnode.Segments.At(i)
			_, _ = r.w.Write(segment.Value(r.source))
		}
		return ast.WalkSkipChildren, nil

	// Blocks.
	case *ast.Paragraph, *ast.TextBlock, *ast.List, *extAST.TableCell:
		// Things that has no content, just children elements, go there.
		break
	case *ast.Heading:
		if !entering {
			break
		}

		// Render it straight away. No nested headings are supported and we expect
		// headings to have limited content, so limit WALK.
		if err := r.renderHeading(tnode); err != nil {
			return ast.WalkStop, fmt.Errorf("rendering heading: %w", err)
		}
		return ast.WalkSkipChildren, nil
	case *ast.HTMLBlock:
		if !entering {
			break
		}

		var segments []text.Segment
		for i := 0; i < node.Lines().Len(); i++ {
			segments = append(segments, node.Lines().At(i))
		}

		if tnode.ClosureLine.Len() != 0 {
			segments = append(segments, tnode.ClosureLine)
		}
		for i, s := range segments {
			o := s.Value(r.source)
			if i == len(segments)-1 {
				o = bytes.TrimSuffix(o, []byte("\n"))
			}
			_, _ = r.w.Write(o)
		}
		return ast.WalkSkipChildren, nil
	case *ast.CodeBlock, *ast.FencedCodeBlock:
		if !entering {
			break
		}

		_, _ = r.w.Write(codeBlockChars)

		var lang []byte
		if fencedNode, isFenced := node.(*ast.FencedCodeBlock); isFenced && fencedNode.Info != nil {
			lang = fencedNode.Info.Text(r.source)
			_, _ = r.w.Write(lang)
			for _, elt := range bytes.Fields(lang) {
				elt = bytes.TrimSpace(bytes.TrimLeft(elt, ". "))
				if len(elt) == 0 {
					continue
				}
				lang = elt
				break
			}
		}

		_, _ = r.w.Write(newLineChar)
		codeBuf := bytes.Buffer{}
		for i := 0; i < tnode.Lines().Len(); i++ {
			line := tnode.Lines().At(i)
			_, _ = codeBuf.Write(line.Value(r.source))
		}

		if formatCode, ok := r.mr.formatters[noAllocString(lang)]; ok {
			code := formatCode(codeBuf.Bytes())
			if !bytes.HasSuffix(code, newLineChar) {
				// Ensure code sample ends with a newline.
				code = append(code, newLineChar...)
			}
			_, _ = r.w.Write(code)
		} else {
			_, _ = r.w.Write(codeBuf.Bytes())
		}

		_, _ = r.w.Write(codeBlockChars)
		return ast.WalkSkipChildren, nil
	case *ast.ThematicBreak:
		if !entering {
			break
		}

		_, _ = r.w.Write(thematicBreakChars)
	case *ast.Blockquote:
		if entering {
			r.w.PushIndent(blockquoteChars)
			if node.Parent() != nil && node.Parent().Kind() == ast.KindListItem &&
				node.PreviousSibling() == nil {
				_, _ = r.w.Write(blockquoteChars)
			}
		} else {
			r.w.PopIndent()
		}

	case *ast.ListItem:
		if entering {
			liMarker := listItemMarkerChars(tnode)
			_, _ = r.w.Write(liMarker)
			if r.mr.listIndentStyle == ListIndentUniform &&
				// We can use 4 spaces for indentation only if
				// that would still qualify as part of the list
				// item text. e.g., given "123. foo",
				// for content to be part of that list item,
				// it must be indented 5 spaces.
				//
				//	123. foo
				//
				//	     bar
				len(liMarker) <= len(fourSpacesChars) {
				r.w.PushIndent(fourSpacesChars)
			} else {
				r.w.PushIndent(bytes.Repeat(spaceChar, len(liMarker)))
			}
		} else {
			if tnode.NextSibling() != nil && tnode.NextSibling().Kind() == ast.KindListItem {
				// Newline after list item.
				_, _ = r.w.Write(newLineChar)
			}
			r.w.PopIndent()
		}

	case *extAST.Table:
		if !entering {
			break
		}

		// Render it straight away. No nested tables are supported and we expect
		// tables to have limited content, so limit WALK.
		if err := r.renderTable(tnode); err != nil {
			return ast.WalkStop, fmt.Errorf("rendering table: %w", err)
		}
		return ast.WalkSkipChildren, nil
	case *extAST.TableRow, *extAST.TableHeader:
		return ast.WalkStop, fmt.Errorf("%v element detected, but table should be rendered in renderTable instead", tnode.Kind())
	default:
		return ast.WalkStop, fmt.Errorf("detected unexpected tree type %v", tnode.Kind())
	}
	return ast.WalkContinue, nil
}

func (r *render) wrapNonEmptyContentWith(b []byte, entering bool) ast.WalkStatus {
	if entering {
		r.w.AddIndentOnFirstWrite(b)
		return ast.WalkContinue
	}

	if r.w.WasIndentOnFirstWriteWritten() {
		_, _ = r.w.Write(b)
		return ast.WalkContinue
	}
	r.w.DelIndentOnFirstWrite(b)
	return ast.WalkContinue
}

func listItemMarkerChars(tnode *ast.ListItem) []byte {
	parList := tnode.Parent().(*ast.List)
	if parList.IsOrdered() {
		cnt := 1
		if parList.Start != 0 {
			cnt = parList.Start
		}
		s := tnode.PreviousSibling()
		for s != nil {
			cnt++
			s = s.PreviousSibling()
		}
		return append(strconv.AppendInt(nil, int64(cnt), 10), parList.Marker, ' ')
	}
	return []byte{parList.Marker, spaceChar[0]}
}

func noAllocString(buf []byte) string {
	return *(*string)(unsafe.Pointer(&buf))
}

// writeClean writes the given byte slice to the writer
// replacing consecutive spaces, newlines, and tabs
// with single spaces.
func writeClean(w io.Writer, bs []byte) error {
	// This works by scanning the byte slice,
	// and writing sub-slices of bs
	// as we see and skip blank sections.

	var (
		// Start of the current sub-slice to be written.
		startIdx int
		// Normalized last character we saw:
		// for whitespace, this is ' ',
		// for everything else, it's left as-is.
		p byte
	)

	for idx, q := range bs {
		if q == '\n' || q == '\r' || q == '\t' {
			q = ' '
		}

		if q == ' ' {
			if p != ' ' {
				// Going from non-blank to blank.
				// Write the current sub-slice and the blank.
				if _, err := w.Write(bs[startIdx:idx]); err != nil {
					return err
				}
				if _, err := w.Write(spaceChar); err != nil {
					return err
				}
			}
			startIdx = idx + 1
		} else if p == ' ' {
			// Going from blank to non-blank.
			// Start a new sub-slice.
			startIdx = idx
		}
		p = q
	}

	_, err := w.Write(bs[startIdx:])
	return err
}

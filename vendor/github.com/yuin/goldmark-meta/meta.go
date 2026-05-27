// package meta is a extension for the goldmark(http://github.com/yuin/goldmark).
//
// This extension parses YAML metadata blocks and store metadata to a
// parser.Context.
package meta

import (
	"bytes"
	"fmt"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	east "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"

	"gopkg.in/yaml.v2"
)

type data struct {
	Map   map[string]interface{}
	Items yaml.MapSlice
	Error error
	Node  gast.Node
}

var contextKey = parser.NewContextKey()

// Option interface sets options for this extension.
type Option interface {
	metaOption()
}

// Get returns a YAML metadata.
func Get(pc parser.Context) map[string]interface{} {
	v := pc.Get(contextKey)
	if v == nil {
		return nil
	}
	d := v.(*data)
	return d.Map
}

// TryGet tries to get a YAML metadata.
// If there are YAML parsing errors, then nil and error are returned
func TryGet(pc parser.Context) (map[string]interface{}, error) {
	dtmp := pc.Get(contextKey)
	if dtmp == nil {
		return nil, nil
	}
	d := dtmp.(*data)
	if d.Error != nil {
		return nil, d.Error
	}
	return d.Map, nil
}

// GetItems returns a YAML metadata.
// GetItems preserves defined key order.
func GetItems(pc parser.Context) yaml.MapSlice {
	v := pc.Get(contextKey)
	if v == nil {
		return nil
	}
	d := v.(*data)
	return d.Items
}

// TryGetItems returns a YAML metadata.
// TryGetItems preserves defined key order.
// If there are YAML parsing errors, then nil and erro are returned.
func TryGetItems(pc parser.Context) (yaml.MapSlice, error) {
	dtmp := pc.Get(contextKey)
	if dtmp == nil {
		return nil, nil
	}
	d := dtmp.(*data)
	if d.Error != nil {
		return nil, d.Error
	}
	return d.Items, nil
}

type metaParser struct {
}

var defaultParser = &metaParser{}

// NewParser returns a BlockParser that can parse YAML metadata blocks.
func NewParser() parser.BlockParser {
	return defaultParser
}

func isSeparator(line []byte) bool {
	line = util.TrimRightSpace(util.TrimLeftSpace(line))
	for i := 0; i < len(line); i++ {
		if line[i] != '-' {
			return false
		}
	}
	return true
}

func (b *metaParser) Trigger() []byte {
	return []byte{'-'}
}

func (b *metaParser) Open(parent gast.Node, reader text.Reader, pc parser.Context) (gast.Node, parser.State) {
	linenum, _ := reader.Position()
	if linenum != 0 {
		return nil, parser.NoChildren
	}
	line, _ := reader.PeekLine()
	if isSeparator(line) {
		return gast.NewTextBlock(), parser.NoChildren
	}
	return nil, parser.NoChildren
}

func (b *metaParser) Continue(node gast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	if isSeparator(line) && !util.IsBlank(line) {
		reader.Advance(segment.Len())
		return parser.Close
	}
	node.Lines().Append(segment)
	return parser.Continue | parser.NoChildren
}

func (b *metaParser) Close(node gast.Node, reader text.Reader, pc parser.Context) {
	lines := node.Lines()
	var buf bytes.Buffer
	for i := 0; i < lines.Len(); i++ {
		segment := lines.At(i)
		buf.Write(segment.Value(reader.Source()))
	}
	d := &data{}
	d.Node = node
	meta := map[string]interface{}{}
	if err := yaml.Unmarshal(buf.Bytes(), &meta); err != nil {
		d.Error = err
	} else {
		d.Map = meta
	}

	metaMapSlice := yaml.MapSlice{}
	if err := yaml.Unmarshal(buf.Bytes(), &metaMapSlice); err != nil {
		d.Error = err
	} else {
		d.Items = metaMapSlice
	}

	pc.Set(contextKey, d)

	if d.Error == nil {
		node.Parent().RemoveChild(node.Parent(), node)
	}
}

func (b *metaParser) CanInterruptParagraph() bool {
	return false
}

func (b *metaParser) CanAcceptIndentedLine() bool {
	return false
}

type astTransformer struct {
	transformerConfig
}

type transformerConfig struct {
	// Renders metadata as an html table.
	Table bool

	// Stores metadata in ast.Document.Meta().
	StoresInDocument bool
}

type transformerOption interface {
	Option

	// SetMetaOption sets options for the metadata parser.
	SetMetaOption(*transformerConfig)
}

var _ transformerOption = &withTable{}

type withTable struct {
	value bool
}

func (o *withTable) metaOption() {}

func (o *withTable) SetMetaOption(m *transformerConfig) {
	m.Table = o.value
}

// WithTable is a functional option that renders a YAML metadata as a table.
func WithTable() Option {
	return &withTable{
		value: true,
	}
}

var _ transformerOption = &withStoresInDocument{}

type withStoresInDocument struct {
	value bool
}

func (o *withStoresInDocument) metaOption() {}

func (o *withStoresInDocument) SetMetaOption(c *transformerConfig) {
	c.StoresInDocument = o.value
}

// WithStoresInDocument is a functional option that parser will store YAML meta in ast.Document.Meta().
func WithStoresInDocument() Option {
	return &withStoresInDocument{
		value: true,
	}
}

func newTransformer(opts ...transformerOption) parser.ASTTransformer {
	p := &astTransformer{
		transformerConfig: transformerConfig{
			Table:            false,
			StoresInDocument: false,
		},
	}
	for _, o := range opts {
		o.SetMetaOption(&p.transformerConfig)
	}
	return p
}

func (a *astTransformer) Transform(node *gast.Document, reader text.Reader, pc parser.Context) {
	dtmp := pc.Get(contextKey)
	if dtmp == nil {
		return
	}
	d := dtmp.(*data)
	if d.Error != nil {
		msg := gast.NewString([]byte(fmt.Sprintf("<!-- %s -->", d.Error)))
		msg.SetCode(true)
		d.Node.AppendChild(d.Node, msg)
		return
	}

	if a.Table {
		meta := GetItems(pc)
		if meta == nil {
			return
		}
		table := east.NewTable()
		alignments := []east.Alignment{}
		for range meta {
			alignments = append(alignments, east.AlignNone)
		}
		row := east.NewTableRow(alignments)
		for _, item := range meta {
			cell := east.NewTableCell()
			cell.AppendChild(cell, gast.NewString([]byte(fmt.Sprintf("%v", item.Key))))
			row.AppendChild(row, cell)
		}
		table.AppendChild(table, east.NewTableHeader(row))

		row = east.NewTableRow(alignments)
		for _, item := range meta {
			cell := east.NewTableCell()
			cell.AppendChild(cell, gast.NewString([]byte(fmt.Sprintf("%v", item.Value))))
			row.AppendChild(row, cell)
		}
		table.AppendChild(table, row)
		node.InsertBefore(node, node.FirstChild(), table)
	}

	if a.StoresInDocument {
		for k, v := range d.Map {
			node.AddMeta(k, v)
		}
	}
}

type meta struct {
	options []Option
}

// Meta is a extension for the goldmark.
var Meta = &meta{}

// New returns a new Meta extension.
func New(opts ...Option) goldmark.Extender {
	e := &meta{
		options: opts,
	}
	return e
}

// Extend implements goldmark.Extender.
func (e *meta) Extend(m goldmark.Markdown) {
	topts := []transformerOption{}
	for _, opt := range e.options {
		if topt, ok := opt.(transformerOption); ok {
			topts = append(topts, topt)
		}
	}
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(NewParser(), 0),
		),
	)
	m.Parser().AddOptions(
		parser.WithASTTransformers(
			util.Prioritized(newTransformer(topts...), 0),
		),
	)
}

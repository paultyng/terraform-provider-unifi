package main

import (
	"context"

	"github.com/paultyng/tfplugingen/sdkparse"
)

type AttributeType interface{}

type Attribute struct {
	Name        string
	Description string
	Deprecated  bool

	Type AttributeType

	Required bool
	Optional bool
	Computed bool

	Sensitive bool
}

type NestingMode int

const (
	InvalidNesting NestingMode = iota
	SingleNesting
	ListNesting
	SetNesting
	MapNesting
	GroupNesting
)

type NestedBlock struct {
	Name string
	Block

	MinItems int
	MaxItems int

	Nesting NestingMode
}

type Block struct {
	// TODO: version?

	Description string
	Deprecated  bool

	Attributes []Attribute
	Blocks     []NestedBlock
}

type Resource struct {
	Name string

	Block
}

func generateModel(ctx context.Context, resources map[string]sdkparse.ResourceInfo) error {
	panic("not implemented")
}

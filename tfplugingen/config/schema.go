package config

type Value interface{}

type SchemaProperties struct {
	Required   []string          `hcl:"required,optional"`
	Properties map[string]Schema `hcl:"property,block"`
}

// type NamedSchema struct {
// 	Name string `hcl:",label"`

// 	Schema
// }

type Schema struct {
	// TODO: NamedSchema references?

	Type             *string `hcl:"type,optional"` // explicitly does not support an array of types
	Format           *string `hcl:"format,optional"`
	Pattern          *string `hcl:"pattern,optional"`
	MaxLength        *int    `hcl:"max_length,optional"`
	MinLength        *int    `hcl:"min_length,optional"`
	Enum             []Value `hcl:"enum,optional"`
	Default          *Value  `hcl:"default,optional"`
	UniqueItems      *bool   `hcl:"unique_items,optional"`
	Maximum          *Value  `hcl:"maximum,optional"`
	ExclusiveMaximum *Value  `hcl:"exclusive_maximum,optional"`
	Minimum          *Value  `hcl:"minimum,optional"`
	ExclusiveMinimum *Value  `hcl:"exclusive_minimum,optional"`

	MaxItems *int    `hcl:"max_items,optional"`
	MinItems *int    `hcl:"min_items,optional"`
	Items    *Schema `hcl:"item,block"`

	*SchemaProperties
}

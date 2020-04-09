package config

type File struct {
	Provider  *Provider  `hcl:"provider,block"`
	Resources []Resource `hcl:"resource,block"`
	// TODO: Schemas   []NamedSchema `hcl:"schema,block"`
}

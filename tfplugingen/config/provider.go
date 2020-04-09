package config

import (
	"strings"

	"github.com/iancoleman/strcase"
)

type Provider struct {
	Name string `hcl:",label"`

	DefaultSDKPackage      string           `hcl:"sdk_package,optional"`
	DefaultClientPackage   string           `hcl:"client_package,optional"`
	DefaultClientType      string           `hcl:"client_type,optional"`
	DefaultResourcePackage string           `hcl:"resource_package,optional"`
	DefaultReadFunc        ResourceTemplate `hcl:"read_func,optional"`
	DefaultCreateFunc      ResourceTemplate `hcl:"create_func,optional"`
	DefaultUpdateFunc      ResourceTemplate `hcl:"update_func,optional"`
	DefaultDeleteFunc      ResourceTemplate `hcl:"delete_func,optional"`
}

func (p *Provider) DefaultResourceType(n string) string {
	// TODO: this should do some case changes and stuff to massage a
	// resource name to type name, should possibly be configurable?
	// maybe supply a list of common replacements or case corrections?

	// strip the provider name prefix
	parts := strings.SplitN(n, "_", 2)
	if len(parts) != 2 {
		return ""
	}
	return strcase.ToCamel(parts[1])
}

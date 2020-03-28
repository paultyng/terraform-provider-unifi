package main

type resourceTemplate string

type resource struct {
	Name    string `hcl:",label"`
	Package string `hcl:"package,optional"`
	Type    string `hcl:"type,optional"` // TODO: infer default from name?
}

type provider struct {
	Name string `hcl:",label"`

	DefaultSDKPackage      string           `hcl:"sdk_package,optional"`
	DefaultClientPackage   string           `hcl:"client_package,optional"`
	DefaultClientType      string           `hcl:"client_type,optional"`
	DefaultResourcePackage string           `hcl:"resource_package,optional"`
	DefaultReadFunc        resourceTemplate `hcl:"read_func,optional"`
	DefaultCreateFunc      resourceTemplate `hcl:"create_func,optional"`
	DefaultUpdateFunc      resourceTemplate `hcl:"update_func,optional"`
	DefaultDeleteFunc      resourceTemplate `hcl:"delete_func,optional"`
}

func (p *provider) DefaultResourceType(n string) string {
	// TODO: this should do some case changes and stuff to massage a
	// resource name to type name, should possibly be configurable?
	return ""
}

type config struct {
	Provider  *provider  `hcl:"provider,block"`
	Resources []resource `hcl:"resource,block"`
}

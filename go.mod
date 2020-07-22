module github.com/paultyng/terraform-provider-unifi

go 1.14

// replace github.com/paultyng/go-unifi => ../go-unifi
// replace github.com/hashicorp/terraform-plugin-docs => ../../hashicorp/terraform-plugin-docs

require (
	github.com/go-test/deep v1.0.4 // indirect
	github.com/hashicorp/go-hclog v0.10.1 // indirect
	github.com/hashicorp/terraform-plugin-docs v0.1.3
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.1
	github.com/hashicorp/yamux v0.0.0-20190923154419-df201c70410d // indirect
	github.com/paultyng/go-unifi v1.3.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

module github.com/paultyng/terraform-provider-unifi

go 1.15

// replace github.com/paultyng/go-unifi => ../go-unifi
// replace github.com/hashicorp/terraform-plugin-docs => ../../hashicorp/terraform-plugin-docs

require (
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/go-test/deep v1.0.4 // indirect
	github.com/hashicorp/go-hclog v0.14.1 // indirect
	github.com/hashicorp/terraform-plugin-docs v0.2.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.3
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/paultyng/go-unifi v1.7.0
	github.com/posener/complete v1.2.1 // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible // indirect
)

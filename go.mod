module github.com/paultyng/terraform-provider-unifi

go 1.16

// replace github.com/paultyng/go-unifi => ../go-unifi
// replace github.com/hashicorp/terraform-plugin-docs => ../../hashicorp/terraform-plugin-docs
// replace github.com/hashicorp/terraform-plugin-sdk/v2 => ../../hashicorp/terraform-plugin-sdk

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/go-test/deep v1.0.4 // indirect
	github.com/hashicorp/go-version v1.5.0
	github.com/hashicorp/terraform-plugin-docs v0.10.1
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.16.0
	github.com/hashicorp/yamux v0.0.0-20200609203250-aecfd211c9ce // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/paultyng/go-unifi v1.25.3
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210126160654-44e461bb6506 // indirect
)

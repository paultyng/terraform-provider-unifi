module github.com/paultyng/terraform-provider-unifi

go 1.15

// replace github.com/paultyng/go-unifi => ../go-unifi
// replace github.com/hashicorp/terraform-plugin-docs => ../../hashicorp/terraform-plugin-docs
replace github.com/hashicorp/terraform-plugin-sdk/v2 => ../../hashicorp/terraform-plugin-sdk


require (
	github.com/hashicorp/terraform-plugin-docs v0.1.4
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.0.3
	github.com/paultyng/go-unifi v1.6.0
)

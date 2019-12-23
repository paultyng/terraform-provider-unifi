package main // import "github.com/paultyng/terraform-provider-unifi"

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"

	"github.com/paultyng/terraform-provider-unifi/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}

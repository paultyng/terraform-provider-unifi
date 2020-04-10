package main // import "github.com/paultyng/terraform-provider-unifi"

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/paultyng/terraform-provider-unifi/internal/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}

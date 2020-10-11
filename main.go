package main // import "github.com/jalfresi/terraform-provider-unifi"

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/jalfresi/terraform-provider-unifi/internal/provider"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.New,
	})
}

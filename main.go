package main // import "github.com/paultyng/terraform-provider-unifi"

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/sayedh/terraform-provider-unifi/internal/provider"
)

// Generate docs for website
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"

	// goreleaser can also pass the specific commit if you want
	// commit  string = ""
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: provider.New(version)}

	if debugMode {
		opts.Debug = true
		opts.ProviderAddr = "registry.terraform.io/sayedh/unifi"
	}

	plugin.Serve(opts)
}

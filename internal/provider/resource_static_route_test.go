package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStaticRoute_nextHop(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRouteConfig_nextHop,
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckFirewallGroupExists(t, "name"),
				// ),
			},
			importStep("unifi_static_route.test"),
		},
	})
}

func TestAccStaticRoute_blackhole(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRouteConfig_blackhole,
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckFirewallGroupExists(t, "name"),
				// ),
			},
			importStep("unifi_static_route.test"),
		},
	})
}

func TestAccStaticRoute_interface(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRouteConfig_interface,
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckFirewallGroupExists(t, "name"),
				// ),
			},
			importStep("unifi_static_route.test"),
		},
	})
}

const testAccStaticRouteConfig_nextHop = `
resource "unifi_static_route" "test" {
	type     = "nexthop-route"
	network  = "172.17.0.0/16"
	name     = "tf-acc basic nexthop"
	distance = 1
	next_hop = "172.16.0.1"
}
`

const testAccStaticRouteConfig_blackhole = `
resource "unifi_static_route" "test" {
	type     = "blackhole"
	network  = "172.17.0.0/16"
	name     = "tf-acc basic blackhole"
	distance = 1
}
`

const testAccStaticRouteConfig_interface = `
resource "unifi_static_route" "test" {
	type      = "interface-route"
	network   = "172.17.0.0/16"
	name      = "tf-acc basic interface"
	distance  = 1
	interface = "WAN2"
}
`

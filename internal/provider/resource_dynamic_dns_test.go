package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDynamicDNS_dyndns(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDynamicDNSConfig,
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckFirewallGroupExists(t, "name"),
				// ),
			},
			importStep("unifi_dynamic_dns.test"),
		},
	})
}

const testAccDynamicDNSConfig = `
resource "unifi_dynamic_dns" "test" {
	service = "dyndns"
	
	host_name = "test.example.com"

	server   = "dyndns.example.com"
	login    = "testuser"
	password = "password"
}
`

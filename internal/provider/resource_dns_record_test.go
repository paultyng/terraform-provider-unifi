package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDNSRecord_default(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSRecordConfig,
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckFirewallGroupExists(t, "name"),
				// ),
			},
			importStep("unifi_dns_record.test"),
		},
	})
}

const testAccDNSRecordConfig = `
resource "unifi_dns_record" "test" {
	service = "default"
	
	host_name = "test.example.com"

	server   = "default.example.com"
	login    = "testuser"
	password = "password"
}
`

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccWLAN_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			// importStep("unifi_wlan.test"),
		},
	})
}

const testAccWLANConfig = `
resource "unifi_wlan" "test" {
	name       = "foo"
	vlan_id    = 202
	passphrase = "12345678"
}
`

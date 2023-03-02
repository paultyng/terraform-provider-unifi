package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPortProfile_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccPortProfileConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_port_profile.test", "poe_mode", "off"),
				),
			},
			importStep("unifi_port_profile.test"),
		},
	})
}

const testAccPortProfileConfig = `
resource "unifi_port_profile" "test" {
	name = "provider created"

	poe_mode	  = "off"
	speed 		  = 1000
	stp_port_mode = false
}
`

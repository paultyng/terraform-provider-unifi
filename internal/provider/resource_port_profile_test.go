package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPortProfile_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckVersionConstraint(t, "< 7.4")
		},
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccPortProfileConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_port_profile.test", "poe_mode", "off"),
					resource.TestCheckResourceAttr("unifi_port_profile.test", "name", name),
				),
			},
			importStep("unifi_port_profile.test"),
		},
	})
}

func testAccPortProfileConfig(name string) string {
	return fmt.Sprintf(`
resource "unifi_port_profile" "test" {
	name = "%s"

	poe_mode	  = "off"
	speed 		  = 1000
	stp_port_mode = false
}
`, name)
}

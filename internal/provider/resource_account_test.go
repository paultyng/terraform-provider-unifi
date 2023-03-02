package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccount_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountConfig("tfacc", "secure"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_account.test", "name", "tfacc"),
				),
			},
			importStep("unifi_account.test"),
		},
	})
}

func TestAccAccount_mac(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountConfig("00B0D06FC226", "00B0D06FC226"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_account.test", "name", "00B0D06FC226"),
					resource.TestCheckResourceAttr("unifi_account.test", "password", "00B0D06FC226"),
				),
			},
			importStep("unifi_account.test"),
		},
	})
}

func testAccAccountConfig(name, password string) string {
	return fmt.Sprintf(`
resource "unifi_account" "test" {
	name = "%s"
	password = "%s"
}
`, name, password)
}

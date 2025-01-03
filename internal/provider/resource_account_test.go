package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAccount_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountConfig(name, "secure"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_account.test", "name", name),
				),
			},
			importStep("unifi_account.test"),
		},
	})
}

func TestAccAccount_mac(t *testing.T) {
	mac, unallocateMac := allocateTestMac(t)
	defer unallocateMac()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccAccountConfig(mac, mac),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_account.test", "name", mac),
					resource.TestCheckResourceAttr("unifi_account.test", "password", mac),
				),
			},
			importStep("unifi_account.test"),
		},
	})
}

func testAccAccountConfig(name, password string) string {
	return fmt.Sprintf(`
resource "unifi_account" "test" {
	name = "%[1]s"
	password = "%[2]s"
}
`, name, password)
}

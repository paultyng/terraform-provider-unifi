package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataAccount_default(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDataAccountConfig("tfusertest", "secure_1234"),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func TestAccDataAccount_mac(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDataMacAccountConfig("00B0D06FC226"),
				Check:  resource.ComposeTestCheckFunc(),
			},
		},
	})
}

func testAccDataAccountConfig(name, password string) string {
	return fmt.Sprintf(`
resource "unifi_account" "test" {
	name = "%s"
	password = "%s"
}

data "unifi_account" "test" {
	name = "%s"
depends_on = [
    unifi_account.test
  ]
}
`, name, password, name)
}

func testAccDataMacAccountConfig(mac string) string {
	return fmt.Sprintf(`
resource "unifi_account" "test" {
	name = "%s"
	password = "%s"
}

data "unifi_account" "test" {
	name = "%s"
depends_on = [
    unifi_account.test
  ]
}
`, mac, mac, mac)
}

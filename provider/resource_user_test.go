package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccUser_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { preCheck(t) },
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig("00:00:5E:00:53:00", "tfacc", "tfacc note"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "note", "tfacc note"),
				),
			},
			importStep("unifi_user.test", "allow_existing"),
			{
				Config: testAccUserConfig("00:00:5E:00:53:00", "tfacc-2", "tfacc note 2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "note", "tfacc note 2"),
				),
			},
			importStep("unifi_user.test", "allow_existing"),
		},
	})
}

func TestAccUser_fixed_ip(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { preCheck(t) },
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig("00:00:5E:00:53:10", "tfacc", "tfacc fixed ip"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "fixed_ip", ""),
				),
			},
			importStep("unifi_user.test", "allow_existing"),
			{
				Config: testAccUserConfig_fixedIP("00:00:5E:00:53:10"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "fixed_ip", "10.1.10.50"),
				),
			},
			importStep("unifi_user.test", "allow_existing"),
			{
				// this passes the network again even though its not used
				// to avoid a destroy order of operations issue, can
				// maybe work it out some other way
				Config: testAccUserConfig_network + testAccUserConfig("00:00:5E:00:53:10", "tfacc", "tfacc fixed ip"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "fixed_ip", ""),
				),
			},
			importStep("unifi_user.test", "allow_existing"),
		},
	})
}

func TestAccUser_blocking(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { preCheck(t) },
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_block("00:00:5E:00:53:20", false),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "blocked", "false"),
				),
			},
			importStep("unifi_user.test", "allow_existing"),
			{
				Config: testAccUserConfig_block("00:00:5E:00:53:20", true),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "blocked", "true"),
				),
			},
			importStep("unifi_user.test", "allow_existing"),
			{
				Config: testAccUserConfig_block("00:00:5E:00:53:20", false),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "blocked", "false"),
				),
			},
			importStep("unifi_user.test", "allow_existing"),
		},
	})
}

// for test MAC addresses, see https://tools.ietf.org/html/rfc7042#section-2.1.2
// func TestAccUser_existing_mac_allow(t *testing.T) {
// func TestAccUser_existing_mac_deny(t *testing.T) {

func testAccUserConfig(mac, name, note string) string {
	return fmt.Sprintf(`
resource "unifi_user" "test" {
	mac  = "%s"
	name = "%s"
	note = "%s"
}
`, mac, name, note)
}

const testAccUserConfig_network = `
variable "subnet" {
	default = "10.1.10.1/24"
}

resource "unifi_network" "test" {
	name    = "tfaccfixedip"
	purpose = "corporate"

	vlan_id      = 66
	subnet       = var.subnet
	dhcp_start   = cidrhost(var.subnet, 6)
	dhcp_stop    = cidrhost(var.subnet, 254)
	dhcp_enabled = true
}
`

func testAccUserConfig_fixedIP(mac string) string {
	return fmt.Sprintf(testAccUserConfig_network+`
resource "unifi_user" "test" {
	mac  = "%s"
	name = "tfacc"
	note = "tfacc fixed ip"

	fixed_ip   = "10.1.10.50"
	network_id = unifi_network.test.id
}
`, mac)
}

func testAccUserConfig_block(mac string, blocked bool) string {
	return fmt.Sprintf(`
resource "unifi_user" "test" {
	mac  = "%s"
	name = "tfacc"
	note = "tfacc block %t"

	blocked = %t
}
`, mac, blocked, blocked)
}

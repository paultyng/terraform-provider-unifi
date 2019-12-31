package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/paultyng/terraform-provider-unifi/unifi"
)

func userImportStep(name string) resource.TestStep {
	return importStep(name, "allow_existing", "skip_forget_on_destroy")
}

// for test MAC addresses, see https://tools.ietf.org/html/rfc7042#section-2.1.2

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
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig("00:00:5E:00:53:00", "tfacc-2", "tfacc note 2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "note", "tfacc note 2"),
				),
			},
			userImportStep("unifi_user.test"),
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
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_fixedIP("00:00:5E:00:53:10"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "fixed_ip", "10.1.10.50"),
				),
			},
			userImportStep("unifi_user.test"),
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
			userImportStep("unifi_user.test"),
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
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_block("00:00:5E:00:53:20", true),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "blocked", "true"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_block("00:00:5E:00:53:20", false),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "blocked", "false"),
				),
			},
			userImportStep("unifi_user.test"),
		},
	})
}

func TestAccUser_existing_mac_allow(t *testing.T) {
	testMAC := "00:00:5e:00:53:30"

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		PreCheck: func() {
			preCheck(t)

			_, err := testClient.CreateUser("default", &unifi.User{
				MAC:  testMAC,
				Name: "tfacc-existing",
				Note: "tfacc-existing",
			})
			if err != nil {
				t.Fatal(err)
			}
		},
		CheckDestroy: func(*terraform.State) error {
			// TODO: CheckDestroy: ,

			return testClient.DeleteUserByMAC("default", testMAC)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_existing(testMAC, "tfacc", "tfacc note", true, true),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "note", "tfacc note"),
				),
			},
			userImportStep("unifi_user.test"),
		},
	})
}

func TestAccUser_existing_mac_deny(t *testing.T) {
	testMAC := "00:00:5e:00:53:40"

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		PreCheck: func() {
			preCheck(t)

			_, err := testClient.CreateUser("default", &unifi.User{
				MAC:  testMAC,
				Name: "tfacc-existing",
				Note: "tfacc-existing",
			})
			if err != nil {
				t.Fatal(err)
			}
		},
		CheckDestroy: func(*terraform.State) error {
			// TODO: CheckDestroy: ,

			return testClient.DeleteUserByMAC("default", testMAC)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccUserConfig_existing(testMAC, "tfacc", "tfacc note", false, false),
				ExpectError: regexp.MustCompile("api\\.err\\.MacUsed"),
			},
		},
	})
}

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

func testAccUserConfig_existing(mac, name, note string, allow, skip bool) string {
	return fmt.Sprintf(`
resource "unifi_user" "test" {
	mac  = "%s"
	name = "%s"
	note = "%s"

	allow_existing         = %t
	skip_forget_on_destroy = %t
}
`, mac, name, note, allow, skip)
}

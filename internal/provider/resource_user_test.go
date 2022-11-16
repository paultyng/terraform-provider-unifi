package provider

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/paultyng/go-unifi/unifi"
)

func userImportStep(name string) resource.TestStep {
	return importStep(name, "allow_existing", "skip_forget_on_destroy")
}

// for test MAC addresses, see https://tools.ietf.org/html/rfc7042#section-2.1.2
func generateTestMac() string {
	mac := net.HardwareAddr{0x00, 0x00, 0x5e, 0x00, 0x53, byte(rand.Intn(256))}
	return mac.String()
}

func TestAccUser_basic(t *testing.T) {
	mac := generateTestMac()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig(mac, "tfacc", "tfacc note"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "note", "tfacc note"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig(mac, "tfacc-2", "tfacc note 2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "note", "tfacc note 2"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig(strings.ReplaceAll(strings.ToLower(mac), ":", "-"), "tfacc-2", "tfacc note 2 dash and lower"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "note", "tfacc note 2 dash and lower"),
				),
			},
			userImportStep("unifi_user.test"),
		},
	})
}

func TestAccUser_fixed_ip(t *testing.T) {
	mac := generateTestMac()
	vlanID := 301

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig(mac, "tfacc", "tfacc fixed ip"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "fixed_ip", ""),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_fixedIP(vlanID, mac),
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
				Config: testAccUserConfig_network(vlanID) + testAccUserConfig(mac, "tfacc", "tfacc fixed ip"),
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
	mac := generateTestMac()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_block(mac, false),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "blocked", "false"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_block(mac, true),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "blocked", "true"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_block(mac, false),
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
	mac := generateTestMac()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)

			_, err := testClient.CreateUser(context.Background(), "default", &unifi.User{
				MAC:  mac,
				Name: "tfacc-existing",
				Note: "tfacc-existing",
			})
			if err != nil {
				t.Fatal(err)
			}
		},
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: CheckDestroy: ,

			return testClient.DeleteUserByMAC(context.Background(), "default", mac)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_existing(mac, "tfacc", "tfacc note", true, true),
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
	mac := generateTestMac()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)

			_, err := testClient.CreateUser(context.Background(), "default", &unifi.User{
				MAC:  mac,
				Name: "tfacc-existing",
				Note: "tfacc-existing",
			})
			if err != nil {
				t.Fatal(err)
			}
		},
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: CheckDestroy: ,

			return testClient.DeleteUserByMAC(context.Background(), "default", mac)
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccUserConfig_existing(mac, "tfacc", "tfacc note", false, false),
				ExpectError: regexp.MustCompile(`api\.err\.MacUsed`),
			},
		},
	})
}

func TestAccUser_fingerprint(t *testing.T) {
	testMAC := generateTestMac()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_fingerprint(testMAC, "tfacc", 123),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "dev_id_override", "123"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_fingerprint(testMAC, "tfacc", 456),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "dev_id_override", "456"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig(testMAC, "tfacc", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "dev_id_override", "0"),
				),
			},
			userImportStep("unifi_user.test"),
		},
	})
}

func TestAccUser_localdns(t *testing.T) {
	testMAC := generateTestMac()
	vlanID := 301

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckVersionConstraint(t, ">= 7.2.91")
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig(testMAC, "tfacc", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "local_dns_record", ""),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_localdns(vlanID, testMAC, "tfacc", "resource.example.com"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "local_dns_record", "resource.example.com"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_localdns(vlanID, testMAC, "tfacc", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "local_dns_record", ""),
				),
			},
			userImportStep("unifi_user.test"),
		},
	})
}

func testCheckUserDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "unifi_user" {
			continue
		}

		_, err := testClient.GetUser(context.Background(), "default", rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("User still exists: %s", rs.Primary.ID)
		}

		if _, ok := err.(*unifi.NotFoundError); ok {
			continue
		}

		return err
	}

	return nil
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

func testAccUserConfig_network(vlanID int) string {
	return fmt.Sprintf(`
variable "subnet" {
	default = "10.1.10.1/24"
}

resource "unifi_network" "test" {
	name    = "tfaccfixedip"
	purpose = "corporate"

	vlan_id      = %d
	subnet       = var.subnet
	dhcp_start   = cidrhost(var.subnet, 6)
	dhcp_stop    = cidrhost(var.subnet, 254)
	dhcp_enabled = true
}
`, vlanID)
}

func testAccUserConfig_fixedIP(vlanID int, mac string) string {
	return fmt.Sprintf(testAccUserConfig_network(vlanID)+`
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

func testAccUserConfig_fingerprint(mac, name string, devIdOverride int) string {
	return fmt.Sprintf(`
resource "unifi_user" "test" {
	mac             = "%s"
	name            = "%s"
	dev_id_override = %d
}
`, mac, name, devIdOverride)
}

func testAccUserConfig_localdns(vlanID int, mac, name string, localDnsRecord string) string {
	return fmt.Sprintf(testAccUserConfig_network(vlanID)+`
resource "unifi_user" "test" {
	mac             = "%s"
	name            = "%s"

	fixed_ip   = "10.1.10.50"
	network_id = unifi_network.test.id
	local_dns_record = "%s"
}
`, mac, name, localDnsRecord)
}

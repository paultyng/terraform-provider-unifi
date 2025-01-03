package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"net"
	"regexp"
	"strings"
	"testing"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/paultyng/go-unifi/unifi"
)

func userImportStep(name string) resource.TestStep {
	return importStep(name, "allow_existing", "skip_forget_on_destroy")
}

func TestAccUser_basic(t *testing.T) {
	mac, unallocateTestMac := allocateTestMac(t)
	defer unallocateTestMac()

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
	mac, unallocateTestMac := allocateTestMac(t)
	defer unallocateTestMac()
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	ip, err := cidr.Host(subnet, 1)
	if err != nil {
		t.Error(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig(mac, name, "tfacc fixed ip"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "fixed_ip", ""),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_fixedIP(name, subnet, vlan, mac, &ip),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "fixed_ip", ip.String()),
				),
			},
			userImportStep("unifi_user.test"),
			{
				// this passes the network again even though its not used
				// to avoid a destroy order of operations issue, can
				// maybe work it out some other way
				Config: testAccUserConfig_network(name, subnet, vlan) + testAccUserConfig(mac, name, "tfacc fixed ip"),
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
	mac, unallocateTestMac := allocateTestMac(t)
	defer unallocateTestMac()
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_block(mac, name, false),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "blocked", "false"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_block(mac, name, true),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "blocked", "true"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_block(mac, name, false),
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
	mac, unallocateTestMac := allocateTestMac(t)
	defer unallocateTestMac()
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)

			_, err := testClient.CreateUser(context.Background(), "default", &unifi.User{
				MAC:  mac,
				Name: name,
				Note: "tfacc-existing",
			})
			if err != nil && strings.Contains(err.Error(), "api.Err.MacUsed") {
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
				Config: testAccUserConfig_existing(mac, name, "tfacc note", true, true),
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
	mac, unallocateTestMac := allocateTestMac(t)
	name := acctest.RandomWithPrefix("tfacc")

	_, err := testClient.CreateUser(context.Background(), "default", &unifi.User{
		MAC:  mac,
		Name: name,
		Note: "tfacc-existing",
	})
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		err := testClient.DeleteUserByMAC(context.Background(), "default", mac)
		if err != nil {
			t.Fatal(err)
		}

		unallocateTestMac()
	}()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccUserConfig_existing(mac, name, "tfacc note", false, false),
				ExpectError: regexp.MustCompile(`api\.err\.MacUsed`),
			},
		},
	})
}

func TestAccUser_fingerprint(t *testing.T) {
	mac, unallocateTestMac := allocateTestMac(t)
	defer unallocateTestMac()
	name := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_fingerprint(mac, name, 123),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "dev_id_override", "123"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_fingerprint(mac, name, 456),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "dev_id_override", "456"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig(mac, name, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "dev_id_override", "0"),
				),
			},
			userImportStep("unifi_user.test"),
		},
	})
}

func TestAccUser_localdns(t *testing.T) {
	mac, unallocateTestMac := allocateTestMac(t)
	defer unallocateTestMac()
	name := acctest.RandomWithPrefix("tfacc")

	subnet, vlan := getTestVLAN(t)

	ip, err := cidr.Host(subnet, 1)
	if err != nil {
		t.Error(err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckVersionConstraint(t, ">= 7.2.91")
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig(mac, name, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "local_dns_record", ""),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_localdns(subnet, vlan, mac, name, "resource.example.com", &ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "local_dns_record", "resource.example.com"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_localdns(subnet, vlan, mac, name, "", &ip),
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

func testAccUserConfig_network(name string, subnet *net.IPNet, vlan int) string {
	return fmt.Sprintf(`
resource "unifi_network" "test" {
	name    = "%[1]s"
	purpose = "corporate"

	vlan_id      = %[3]d
	subnet       = "%[2]s"
	dhcp_start   = cidrhost("%[2]s", 6)
	dhcp_stop    = cidrhost("%[2]s", 254)
	dhcp_enabled = true
}
`, name, subnet, vlan)
}

func testAccUserConfig_fixedIP(name string, subnet *net.IPNet, vlan int, mac string, ip *net.IP) string {
	return fmt.Sprintf(testAccUserConfig_network(name, subnet, vlan)+`
resource "unifi_user" "test" {
	mac  = "%[1]s"
	name = "%[2]s"
	note = "%[2]s fixed ip"

	fixed_ip   = "%[3]s"
	network_id = unifi_network.test.id
}
`, mac, name, ip)
}

func testAccUserConfig_block(mac, name string, blocked bool) string {
	return fmt.Sprintf(`
resource "unifi_user" "test" {
	mac  = "%[1]s"
	name = "%[2]s"
	note = "%[2]s block %[3]t"

	blocked = %[3]t
}
`, mac, name, blocked)
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

func testAccUserConfig_localdns(subnet *net.IPNet, vlan int, mac, name string, localDnsRecord string, ip *net.IP) string {
	return fmt.Sprintf(testAccUserConfig_network(name, subnet, vlan)+`
resource "unifi_user" "test" {
	mac             = "%[1]s"
	name            = "%[2]s"

	fixed_ip   = "%[4]s"
	network_id = unifi_network.test.id
	local_dns_record = "%[3]s"
}
`, mac, name, localDnsRecord, ip)
}

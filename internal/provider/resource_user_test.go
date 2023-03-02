package provider

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/paultyng/go-unifi/unifi"
)

func userImportStep(name string) resource.TestStep {
	return importStep(name, "allow_existing", "skip_forget_on_destroy")
}

var (
	testMacLock       sync.Mutex
	testMacsAvailable []*net.HardwareAddr
)

// for test MAC addresses, see https://tools.ietf.org/html/rfc7042#section-2.1.2
func init() {
	testMacsAvailable = make([]*net.HardwareAddr, 256)

	for i := 0; i < 256; i++ {
		testMacsAvailable[i] = &net.HardwareAddr{0x00, 0x00, 0x5e, 0x00, 0x53, byte(i)}
	}
}

func allocateTestMac(t *testing.T) (string, func()) {
	testMacLock.Lock()
	defer testMacLock.Unlock()

	if len(testMacsAvailable) == 0 {
		t.Fatal("Unable to allocate test MAC")
	}

	var mac *net.HardwareAddr
	mac, testMacsAvailable = testMacsAvailable[0], testMacsAvailable[1:]

	unallocate := func() {
		testMacLock.Lock()
		defer testMacLock.Unlock()

		testMacsAvailable = append(testMacsAvailable, mac) //nolint:makezero
	}

	return mac.String(), unallocate
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
				Config: testAccUserConfig(mac, "tfacc", "tfacc fixed ip"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_user.test", "fixed_ip", ""),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_fixedIP(subnet, vlan, mac, &ip),
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
				Config: testAccUserConfig_network(subnet, vlan) + testAccUserConfig(mac, "tfacc", "tfacc fixed ip"),
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
	mac, unallocateTestMac := allocateTestMac(t)
	defer unallocateTestMac()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)

			_, err := testClient.CreateUser(context.Background(), "default", &unifi.User{
				MAC:  mac,
				Name: "tfacc-existing",
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
	mac, unallocateTestMac := allocateTestMac(t)

	_, err := testClient.CreateUser(context.Background(), "default", &unifi.User{
		MAC:  mac,
		Name: "tfacc-existing",
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
				Config:      testAccUserConfig_existing(mac, "tfacc", "tfacc note", false, false),
				ExpectError: regexp.MustCompile(`api\.err\.MacUsed`),
			},
		},
	})
}

func TestAccUser_fingerprint(t *testing.T) {
	mac, unallocateTestMac := allocateTestMac(t)
	defer unallocateTestMac()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testCheckUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccUserConfig_fingerprint(mac, "tfacc", 123),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "dev_id_override", "123"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_fingerprint(mac, "tfacc", 456),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "dev_id_override", "456"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig(mac, "tfacc", ""),
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
				Config: testAccUserConfig(mac, "tfacc", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "local_dns_record", ""),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_localdns(subnet, vlan, mac, "tfacc", "resource.example.com", &ip),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_user.test", "local_dns_record", "resource.example.com"),
				),
			},
			userImportStep("unifi_user.test"),
			{
				Config: testAccUserConfig_localdns(subnet, vlan, mac, "tfacc", "", &ip),
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

func testAccUserConfig_network(subnet *net.IPNet, vlan int) string {
	return fmt.Sprintf(`
resource "unifi_network" "test" {
	name    = "tfaccfixedip"
	purpose = "corporate"

	vlan_id      = %[2]d
	subnet       = "%[1]s"
	dhcp_start   = cidrhost("%[1]s", 6)
	dhcp_stop    = cidrhost("%[1]s", 254)
	dhcp_enabled = true
}
`, subnet, vlan)
}

func testAccUserConfig_fixedIP(subnet *net.IPNet, vlan int, mac string, ip *net.IP) string {
	return fmt.Sprintf(testAccUserConfig_network(subnet, vlan)+`
resource "unifi_user" "test" {
	mac  = "%[1]s"
	name = "tfacc"
	note = "tfacc fixed ip"

	fixed_ip   = "%[2]s"
	network_id = unifi_network.test.id
}
`, mac, ip)
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

func testAccUserConfig_localdns(subnet *net.IPNet, vlan int, mac, name string, localDnsRecord string, ip *net.IP) string {
	return fmt.Sprintf(testAccUserConfig_network(subnet, vlan)+`
resource "unifi_user" "test" {
	mac             = "%[1]s"
	name            = "%[2]s"

	fixed_ip   = "%[4]s"
	network_id = unifi_network.test.id
	local_dns_record = "%[3]s"
}
`, mac, name, localDnsRecord, ip)
}

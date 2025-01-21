package provider

import (
	"fmt"
	"net"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccWLAN_wpapsk(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wpapsk(name, subnet, vlan, "disabled"),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_open(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_open(name, subnet, vlan),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_open_mac_filter(name, subnet, vlan),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_open(name, subnet, vlan),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_change_security_and_pmf(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wpapsk(name, subnet, vlan, "disabled"),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_open(name, subnet, vlan),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_wpapsk(name, subnet, vlan, "optional"),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_wpapsk(name, subnet, vlan, "required"),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_wpapsk(name, subnet, vlan, "disabled"),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_schedule(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_schedule(name, subnet, vlan),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			// remove schedule
			{
				Config: testAccWLANConfig_open(name, subnet, vlan),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_wpaeap(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wpaeap(name, subnet, vlan),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_wlan_band(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wlan_band(name, subnet, vlan),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_no2ghz_oui(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_no2ghz_oui(name, subnet, vlan),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_proxy_arp(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_proxy_arp(name, subnet, vlan, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_wlan.test", "proxy_arp", "true"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_bss_transition(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_bss_transition(name, subnet, vlan, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_wlan.test", "bss_transition", "false"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_uapsd(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_uapsd(name, subnet, vlan),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_fast_roaming_enabled(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_fast_roaming_enabled(name, subnet, vlan, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_wlan.test", "fast_roaming_enabled", "true"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_wpa3(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckMinVersion(t, controllerVersionWPA3)
		},
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wpa3(name, subnet, vlan, false, "required"),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_wpa3(name, subnet, vlan, true, "optional"),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_wpa3(name, subnet, vlan, false, "required"),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_minimum_data_rate(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	subnet, vlan := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_minimum_data_rate(name, subnet, vlan, 5500, 18000),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_minimum_data_rate(name, subnet, vlan, 1000, 18000),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_minimum_data_rate(name, subnet, vlan, 0, 0),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_minimum_data_rate(name, subnet, vlan, 6000, 9000),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_minimum_data_rate(name, subnet, vlan, 18000, 6000),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func testAccWLANBaseConfig(name string, subnet *net.IPNet, vlan int) string {
	return fmt.Sprintf(`
data "unifi_ap_group" "default" {}

data "unifi_user_group" "default" {}

resource "unifi_network" "test" {
	name    = "%[1]s"
	purpose = "corporate"
	subnet  = "%[2]s"
    vlan_id = "%[3]d"
}`, name, subnet, vlan)
}

func testAccWLANConfig_wpapsk(name string, subnet *net.IPNet, vlan int, pmf string) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
resource "unifi_wlan" "test" {
	name          = "%[1]s-wpapsk"
	network_id    = unifi_network.test.id
	passphrase    = "12345678"
	ap_group_ids  = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "wpapsk"
	
	multicast_enhance = true

	pmf_mode = %[2]q
}
`, name, pmf)
}

func testAccWLANConfig_wpaeap(name string, subnet *net.IPNet, vlan int) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
data "unifi_radius_profile" "default" {}

resource "unifi_setting_radius" "this" {
	enabled = true
	secret  = "securepw"
}

resource "unifi_wlan" "test" {
	name          = "%[1]s-wpapsk"
	network_id    = unifi_network.test.id
	passphrase    = "12345678"
	ap_group_ids  = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "wpaeap"

	radius_profile_id = data.unifi_radius_profile.default.id
}
`, name)
}

func testAccWLANConfig_open(name string, subnet *net.IPNet, vlan int) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
resource "unifi_wlan" "test" {
	name          = "%[1]s-open"
	network_id    = unifi_network.test.id
	ap_group_ids  = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "open"
}
`, name)
}

func testAccWLANConfig_schedule(name string, subnet *net.IPNet, vlan int) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
resource "unifi_wlan" "test" {
	name          = "%[1]s-sched"
	network_id    = unifi_network.test.id
	ap_group_ids  = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "open"

	schedule {
		day_of_week = "mon"
		start_hour  = 3
		duration    = 60*6
	}

	schedule {
		day_of_week  = "wed"
		start_hour   = 13
		start_minute = 30
		duration     = (60*3)+30
		name         = "minute"
	}

	schedule {
		day_of_week = "thu"
		start_hour  = 19
		duration    = 60*1
	}

	schedule {
		day_of_week = "fri"
		start_hour  = 19
		duration    = 60*1
	}
}
`, name)
}

func testAccWLANConfig_open_mac_filter(name string, subnet *net.IPNet, vlan int) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
resource "unifi_wlan" "test" {
	name          = "%[1]s-open"
	network_id    = unifi_network.test.id
	ap_group_ids  = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "open"

	mac_filter_enabled = true
	mac_filter_list    = ["ab:cd:ef:12:34:56"]
	mac_filter_policy  = "allow"
}
`, name)
}

func testAccWLANConfig_wlan_band(name string, subnet *net.IPNet, vlan int) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
resource "unifi_wlan" "test" {
	name              = "%[1]s-wpapsk"
	network_id        = unifi_network.test.id
	passphrase        = "12345678"
	ap_group_ids      = [data.unifi_ap_group.default.id]
	user_group_id     = data.unifi_user_group.default.id
	security          = "wpapsk"
	wlan_band         = "5g"
	multicast_enhance = true
}
`, name)
}

func testAccWLANConfig_no2ghz_oui(name string, subnet *net.IPNet, vlan int) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
resource "unifi_wlan" "test" {
	name              = "%[1]s-wpapsk"
	network_id        = unifi_network.test.id
	passphrase        = "12345678"
	ap_group_ids      = [data.unifi_ap_group.default.id]
	user_group_id     = data.unifi_user_group.default.id
	security          = "wpapsk"
	no2ghz_oui        = false
	multicast_enhance = true
}
`, name)
}

func testAccWLANConfig_proxy_arp(name string, subnet *net.IPNet, vlan int, proxyArp bool) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
resource "unifi_wlan" "test" {
	name          = "%[1]s"
	network_id    = unifi_network.test.id
	passphrase    = "12345678"
    ap_group_ids  = [data.unifi_ap_group.default.id]
    user_group_id = data.unifi_user_group.default.id
	security      = "wpapsk"
	proxy_arp     = %[2]t
}
`, name, proxyArp)
}

func testAccWLANConfig_bss_transition(name string, subnet *net.IPNet, vlan int, bssTransition bool) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
resource "unifi_wlan" "test" {
	name           = "%[1]s"
	network_id     = unifi_network.test.id
	passphrase     = "12345678"
    ap_group_ids   = [data.unifi_ap_group.default.id]
    user_group_id  = data.unifi_user_group.default.id
	security       = "wpapsk"
	bss_transition = %[2]t
}
`, name, bssTransition)
}

func testAccWLANConfig_uapsd(name string, subnet *net.IPNet, vlan int) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
resource "unifi_wlan" "test" {
	name          = "%[1]s-wpapsk"
	network_id    = unifi_network.test.id
	passphrase    = "12345678"
	ap_group_ids  = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "wpapsk"
	uapsd         = true
}
`, name)
}

func testAccWLANConfig_fast_roaming_enabled(name string, subnet *net.IPNet, vlan int, fastRoamingEnabled bool) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
resource "unifi_wlan" "test" {
	name                 = "%[1]s"
	network_id           = unifi_network.test.id
	passphrase           = "12345678"
    ap_group_ids         = [data.unifi_ap_group.default.id]
    user_group_id        = data.unifi_user_group.default.id
	security             = "wpapsk"
	fast_roaming_enabled = %[2]t
}
`, name, fastRoamingEnabled)
}

func testAccWLANConfig_wpa3(name string, subnet *net.IPNet, vlan int, wpa3Transition bool, pmf string) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
resource "unifi_wlan" "test" {
	name          = "%[1]s-wpapsk"
	network_id    = unifi_network.test.id
	passphrase    = "12345678"
	ap_group_ids  = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "wpapsk"

	wpa3_support    = true
	wpa3_transition = %[2]t
	pmf_mode        = %[3]q
}
`, name, wpa3Transition, pmf)
}

func testAccWLANConfig_minimum_data_rate(name string, subnet *net.IPNet, vlan int, min2g int, min5g int) string {
	return testAccWLANBaseConfig(name, subnet, vlan) + fmt.Sprintf(`
resource "unifi_wlan" "test" {
	name          = "%[1]s-wpapsk"
	network_id    = unifi_network.test.id
	passphrase    = "12345678"
	ap_group_ids  = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "wpapsk"

	multicast_enhance = true

	minimum_data_rate_2g_kbps = %[2]d
	minimum_data_rate_5g_kbps = %[3]d
}
`, name, min2g, min5g)
}

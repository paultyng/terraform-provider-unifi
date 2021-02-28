package provider

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TODO: test changing security, see https://github.com/paultyng/terraform-provider-unifi/issues/32

// there is a max of 4 SSID's at once, and if you are running this on a
// controller with existing SSID's, you may want to limit the concurrency.
var wlanConcurrency chan struct{}

func init() {
	wcs := os.Getenv("UNIFI_ACC_WLAN_CONCURRENCY")
	if wcs == "" {
		// default concurrent SSIDs
		wcs = "1"
	}
	wc, err := strconv.Atoi(wcs)
	if err != nil {
		panic(err)
	}
	wlanConcurrency = make(chan struct{}, wc)
}

func wlanPreCheck(t *testing.T) {
	if cap(wlanConcurrency) == 0 {
		t.Skip("concurrency for WLAN testing set to 0")
	}

	wlanConcurrency <- struct{}{}
}

func TestAccWLAN_wpapsk(t *testing.T) {
	vlanID := getTestVLAN(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckV6Only(t)
			wlanPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			<-wlanConcurrency
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wpapsk(vlanID),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_open(t *testing.T) {
	vlanID := getTestVLAN(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckV6Only(t)
			wlanPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			<-wlanConcurrency
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_open(vlanID),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_open_mac_filter(vlanID),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_open(vlanID),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_change_security(t *testing.T) {
	vlanID := getTestVLAN(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckV6Only(t)
			wlanPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			<-wlanConcurrency
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wpapsk(vlanID),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_open(vlanID),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_wpapsk(vlanID),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
		},
	})
}

func TestAccWLAN_schedule(t *testing.T) {
	vlanID := getTestVLAN(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckV6Only(t)
			wlanPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			<-wlanConcurrency
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_schedule(vlanID),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_wpaeap(t *testing.T) {
	if os.Getenv("UNIFI_TEST_RADIUS") == "" {
		t.Skip("UNIFI_TEST_RADIUS not set, skipping RADIUS test")
	}

	vlanID := getTestVLAN(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckV6Only(t)
			wlanPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			<-wlanConcurrency
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wpaeap(vlanID),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_wlan_band(t *testing.T) {
	vlanID := getTestVLAN(t)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckV6Only(t)
			wlanPreCheck(t)
		},
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			<-wlanConcurrency
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wlan_band(vlanID),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func testAccWLANConfig_wpapsk(vlanID int) string {
	return fmt.Sprintf(`
data "unifi_ap_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_network" "test" {
	name    = "tfacc"
	purpose = "corporate"

	subnet        = cidrsubnet("10.0.0.0/8", 4, %[1]d)
	vlan_id       = %[1]d
}

resource "unifi_wlan" "test" {
	name          = "tfacc-wpapsk"
	network_id    = unifi_network.test.id
	passphrase    = "12345678"
	ap_group_ids = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "wpapsk"
	
	multicast_enhance = true
}
`, vlanID)
}

func testAccWLANConfig_wpaeap(vlanID int) string {
	return fmt.Sprintf(`
data "unifi_ap_group" "default" {
}

data "unifi_user_group" "default" {
}

data "unifi_radius_profile" "default" {
}

resource "unifi_network" "test" {
	name    = "tfacc"
	purpose = "corporate"

	subnet        = cidrsubnet("10.0.0.0/8", 4, %[1]d)
	vlan_id       = %[1]d
}

resource "unifi_wlan" "test" {
	name          = "tfacc-wpapsk"
	network_id    = unifi_network.test.id
	passphrase    = "12345678"
	ap_group_ids = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "wpaeap"

	radius_profile_id = data.unifi_radius_profile.default.id
}
`, vlanID)
}

func testAccWLANConfig_open(vlanID int) string {
	return fmt.Sprintf(`
data "unifi_ap_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_network" "test" {
	name    = "tfacc"
	purpose = "corporate"

	subnet        = cidrsubnet("10.0.0.0/8", 4, %[1]d)
	vlan_id       = %[1]d
}

resource "unifi_wlan" "test" {
	name          = "tfacc-open"
	network_id    = unifi_network.test.id
	ap_group_ids  = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "open"
}
`, vlanID)
}

func testAccWLANConfig_schedule(vlanID int) string {
	return fmt.Sprintf(`
data "unifi_ap_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_network" "test" {
	name    = "tfacc"
	purpose = "corporate"

	subnet        = cidrsubnet("10.0.0.0/8", 4, %[1]d)
	vlan_id       = %[1]d
}

resource "unifi_wlan" "test" {
	name          = "tfacc-open-schedule"
	network_id    = unifi_network.test.id
	ap_group_ids = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "open"

	schedule {
		day_of_week = "mon"
		block_start = "03:00"
		block_end   = "9:00"
	}

	schedule {
		day_of_week = "wed"
		block_start = "13:00"
		block_end   = "17:00"
	}
}
`, vlanID)
}

func testAccWLANConfig_open_mac_filter(vlanID int) string {
	return fmt.Sprintf(`
data "unifi_ap_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_network" "test" {
	name    = "tfacc"
	purpose = "corporate"

	subnet        = cidrsubnet("10.0.0.0/8", 4, %[1]d)
	vlan_id       = %[1]d
}

resource "unifi_wlan" "test" {
	name          = "tfacc-open"
	network_id    = unifi_network.test.id
	ap_group_ids  = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "open"

	mac_filter_enabled = true
	mac_filter_list    = ["ab:cd:ef:12:34:56"]
	mac_filter_policy  = "allow"
}
`, vlanID)
}

func testAccWLANConfig_wlan_band(vlanID int) string {
	return fmt.Sprintf(`
data "unifi_ap_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_network" "test" {
	name    = "tfacc"
	purpose = "corporate"

	subnet        = cidrsubnet("10.0.0.0/8", 4, %[1]d)
	vlan_id       = %[1]d
}

resource "unifi_wlan" "test" {
	name          = "tfacc-wpapsk"
	network_id    = unifi_network.test.id
	passphrase    = "12345678"
	ap_group_ids = [data.unifi_ap_group.default.id]
	user_group_id = data.unifi_user_group.default.id
	security      = "wpapsk"
	wlan_band = "5g"
	
	multicast_enhance = true
}
`, vlanID)
}

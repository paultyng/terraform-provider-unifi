package provider

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

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

func wlanPreCheck(t *testing.T) func() {
	return func() {
		if cap(wlanConcurrency) == 0 {
			t.Skip("concurrency for WLAN testing set to 0")
		}

		preCheck(t)

		wlanConcurrency <- struct{}{}
	}
}

func TestAccWLAN_wpapsk(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: wlanPreCheck(t),
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			<-wlanConcurrency
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wpapsk,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_open(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: wlanPreCheck(t),
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			<-wlanConcurrency
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_open,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_open_mac_filter,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_open,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_change_security(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: wlanPreCheck(t),
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			<-wlanConcurrency
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wpapsk,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_open,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_wpapsk,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
		},
	})
}

func TestAccWLAN_schedule(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: wlanPreCheck(t),
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			<-wlanConcurrency
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_schedule,
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

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: wlanPreCheck(t),
		CheckDestroy: func(*terraform.State) error {
			// TODO: actual CheckDestroy

			<-wlanConcurrency
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wpaeap,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

const testAccWLANConfig_wpapsk = `
data "unifi_wlan_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_wlan" "test" {
	name          = "tfacc-wpapsk"
	vlan_id       = 202
	passphrase    = "12345678"
	wlan_group_id = data.unifi_wlan_group.default.id
	user_group_id = data.unifi_user_group.default.id
	security      = "wpapsk"
	
	multicast_enhance = true
}
`

const testAccWLANConfig_wpaeap = `
data "unifi_wlan_group" "default" {
}

data "unifi_user_group" "default" {
}

data "unifi_radius_profile" "default" {
}

resource "unifi_wlan" "test" {
	name          = "tfacc-wpapsk"
	vlan_id       = 202
	passphrase    = "12345678"
	wlan_group_id = data.unifi_wlan_group.default.id
	user_group_id = data.unifi_user_group.default.id
	security      = "wpaeap"

	radius_profile_id = data.unifi_radius_profile.default.id
}
`

const testAccWLANConfig_open = `
data "unifi_wlan_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_wlan" "test" {
	name          = "tfacc-open"
	vlan_id       = 202
	wlan_group_id = data.unifi_wlan_group.default.id
	user_group_id = data.unifi_user_group.default.id
	security      = "open"
}
`

const testAccWLANConfig_schedule = `
data "unifi_wlan_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_wlan" "test" {
	name          = "tfacc-open-schedule"
	vlan_id       = 202
	wlan_group_id = data.unifi_wlan_group.default.id
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
`

const testAccWLANConfig_open_mac_filter = `
data "unifi_wlan_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_wlan" "test" {
	name          = "tfacc-open"
	vlan_id       = 202
	wlan_group_id = data.unifi_wlan_group.default.id
	user_group_id = data.unifi_user_group.default.id
	security      = "open"

	mac_filter_enabled = true
	mac_filter_list    = ["ab:cd:ef:12:34:56"]
	mac_filter_policy  = "allow"
}
`

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccWLAN_v5_wpapsk(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckV5Only(t)
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
				Config: testAccWLANConfig_v5_wpapsk,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_v5_open(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckV5Only(t)
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
				Config: testAccWLANConfig_v5_open,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_v5_open_mac_filter,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_v5_open,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_v5_change_security(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckV5Only(t)
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
				Config: testAccWLANConfig_v5_wpapsk,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_v5_open,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
			{
				Config: testAccWLANConfig_v5_wpapsk,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
		},
	})
}

func TestAccWLAN_v5_schedule(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckV5Only(t)
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
				Config: testAccWLANConfig_v5_schedule,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

func TestAccWLAN_v5_wpaeap(t *testing.T) {
	if os.Getenv("UNIFI_TEST_RADIUS") == "" {
		t.Skip("UNIFI_TEST_RADIUS not set, skipping RADIUS test")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckV5Only(t)
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
				Config: testAccWLANConfig_v5_wpaeap,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			importStep("unifi_wlan.test"),
		},
	})
}

const testAccWLANConfig_v5_wpapsk = `
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

const testAccWLANConfig_v5_wpaeap = `
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

const testAccWLANConfig_v5_open = `
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

const testAccWLANConfig_v5_schedule = `
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

const testAccWLANConfig_v5_open_mac_filter = `
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

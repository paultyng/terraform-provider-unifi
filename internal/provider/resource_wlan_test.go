package provider

import (
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

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

func TestAccWLAN_wpapsk(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		PreCheck: func() {
			preCheck(t)

			wlanConcurrency <- struct{}{}
		},
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
		Providers: providers,
		PreCheck: func() {
			preCheck(t)

			wlanConcurrency <- struct{}{}
		},
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
		},
	})
}

func TestAccWLAN_change_security(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		PreCheck: func() {
			preCheck(t)

			wlanConcurrency <- struct{}{}
		},
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

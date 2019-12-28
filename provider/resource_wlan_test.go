package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func wlanImportStep() resource.TestStep {
	return importStep("unifi_wlan.test",
		"name", "passphrase", "vlan_id", "wlan_group_id",
		"user_group_id", "security",
	)
}

func TestAccWLAN_wpapsk(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preCheck(t) },
		Providers: providers,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_wpapsk,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			wlanImportStep(),
		},
	})
}

func TestAccWLAN_open(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { preCheck(t) },
		Providers: providers,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccWLANConfig_open,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
			wlanImportStep(),
		},
	})
}

const testAccWLANConfig_wpapsk = `
data "unifi_wlan_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_wlan" "test" {
	name          = "foowpapsk"
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
	name          = "fooopen"
	vlan_id       = 202
	wlan_group_id = data.unifi_wlan_group.default.id
	user_group_id = data.unifi_user_group.default.id
	security      = "open"
}
`

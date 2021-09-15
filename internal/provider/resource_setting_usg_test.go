package provider

import (
	"fmt"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// using an additional lock to the one around the resource to avoid deadlocking accidentally
var settingUsgLock = sync.Mutex{}

func TestAccSettingUsg_mdns(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			settingUsgLock.Lock()
			t.Cleanup(func() {
				settingUsgLock.Unlock()
			})
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingUsgConfig_mdns(true),
				Check:  resource.ComposeTestCheckFunc(),
			},
			importStep("unifi_setting_usg.test"),
			{
				Config: testAccSettingUsgConfig_mdns(false),
				Check:  resource.ComposeTestCheckFunc(),
			},
			importStep("unifi_setting_usg.test"),
			{
				Config: testAccSettingUsgConfig_mdns(true),
				Check:  resource.ComposeTestCheckFunc(),
			},
			importStep("unifi_setting_usg.test"),
		},
	})
}

func TestAccSettingUsg_dhcpRelay(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			settingUsgLock.Lock()
			t.Cleanup(func() {
				settingUsgLock.Unlock()
			})
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingUsgConfig_dhcpRelay(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			importStep("unifi_setting_usg.test"),
		},
	})
}

func TestAccSettingUsg_site(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			settingUsgLock.Lock()
			t.Cleanup(func() {
				settingUsgLock.Unlock()
			})
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingUsgConfig_site(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      "unifi_setting_usg.test",
				ImportState:       true,
				ImportStateIdFunc: siteAndIDImportStateIDFunc("unifi_setting_usg.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSettingUsgConfig_mdns(mdns bool) string {
	return fmt.Sprintf(`
resource "unifi_setting_usg" "test" {
	multicast_dns_enabled = %t
}
`, mdns)
}

func testAccSettingUsgConfig_dhcpRelay() string {
	return `
resource "unifi_setting_usg" "test" {
	dhcp_relay_servers = [
		"10.1.2.3",
		"10.1.2.4",
	]
}
`
}

func testAccSettingUsgConfig_site() string {
	return `
resource "unifi_site" "test" {
	description = "test"
}

resource "unifi_setting_usg" "test" {
	site = unifi_site.test.name

	multicast_dns_enabled = true
}
`
}

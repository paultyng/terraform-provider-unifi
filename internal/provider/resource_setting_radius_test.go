package provider

import (
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var settingRadiusLock = sync.Mutex{}

func TestAccSettingRadius_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			settingRadiusLock.Lock()
			t.Cleanup(func() {
				settingRadiusLock.Unlock()
			})
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingRadiusConfig_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			importStep("unifi_setting_radius.test"),
		},
	})
}

func TestAccSettingRadius_site(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			settingRadiusLock.Lock()
			t.Cleanup(func() {
				settingRadiusLock.Unlock()
			})
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingRadiusConfig_site(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      "unifi_setting_radius.test",
				ImportState:       true,
				ImportStateIdFunc: siteAndIDImportStateIDFunc("unifi_setting_radius.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSettingRadius_full(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			settingRadiusLock.Lock()
			t.Cleanup(func() {
				settingRadiusLock.Unlock()
			})
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingRadiusConfig_full(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      "unifi_setting_radius.test",
				ImportState:       true,
				ImportStateIdFunc: siteAndIDImportStateIDFunc("unifi_setting_radius.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSettingRadius_vlan(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			settingRadiusLock.Lock()
			t.Cleanup(func() {
				settingRadiusLock.Unlock()
			})
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingRadiusConfig_vlan(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			importStep("unifi_setting_radius.test"),
		},
	})
}

func testAccSettingRadiusConfig_basic() string {
	return `
resource "unifi_setting_radius" "test" {
	enabled = true
	secret = "securepw"
}
`
}

func testAccSettingRadiusConfig_site() string {
	return `
resource "unifi_site" "test" {
	description = "test"
}

resource "unifi_setting_radius" "test" {
	site = unifi_site.test.name
	enabled = true
	secret = "securepw"
}
`
}

func testAccSettingRadiusConfig_full() string {
	return `
resource "unifi_setting_radius" "test" {
	enabled = true
	secret = "securepw"
	accounting_port = "9999"
	auth_port = "8888"
}
`
}

func testAccSettingRadiusConfig_vlan() string {
	return `
resource "unifi_setting_radius" "test" {
	enabled = true
	secret = "securepw"
	accounting_enabled = true
}
`
}

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSettingMgmt_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingMgmtConfig_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			importStep("unifi_setting_mgmt.test"),
		},
	})
}

func TestAccSettingMgmt_site(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSettingMgmtConfig_site(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      "unifi_setting_mgmt.test",
				ImportState:       true,
				ImportStateIdFunc: siteAndIDImportStateIDFunc("unifi_setting_mgmt.test"),
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSettingMgmtConfig_basic() string {
	return `
resource "unifi_setting_mgmt" "test" {
	auto_upgrade = true
}
`
}

func testAccSettingMgmtConfig_site() string {
	return `
resource "unifi_site" "test" {
	description = "test"
}

resource "unifi_setting_mgmt" "test" {
	site = unifi_site.test.name
	auto_upgrade = true
}
`
}

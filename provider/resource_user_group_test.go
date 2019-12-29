package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccUserGroup_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { preCheck(t) },
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupConfig,
				// Check:  resource.ComposeTestCheckFunc(
				// // testCheckUserGroupExists(t, "name"),
				// ),
			},
			{
				Config: testAccUserGroupConfig_qos,
			},
			importStep("unifi_user_group.test"),
			{
				Config: testAccUserGroupConfig,
			},
			importStep("unifi_user_group.test"),
		},
	})
}

const testAccUserGroupConfig = `
resource "unifi_user_group" "test" {
	name = "tfacc"
}
`

const testAccUserGroupConfig_qos = `
resource "unifi_user_group" "test" {
	name = "tfacc"

	qos_rate_max_up   = 2000
	qos_rate_max_down = 50
}
`

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataUserGroup_default(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { preCheck(t) },
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDataUserGroupConfig_default,
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
		},
	})
}

const testAccDataUserGroupConfig_default = `
data "unifi_user_group" "default" {
}
`

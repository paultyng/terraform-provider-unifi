package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataUserGroup_default(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
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

func TestAccDataUserGroup_multiple_providers(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { preCheck(t) },
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"unifi2": func() (*schema.Provider, error) {
				return New("acctest")(), nil
			},
			"unifi3": func() (*schema.Provider, error) {
				return New("acctest")(), nil
			},
		},
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: `
				data "unifi_user_group" "unifi2" {
					provider = "unifi2"
				}
				data "unifi_user_group" "unifi3" {
					provider = "unifi3"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
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

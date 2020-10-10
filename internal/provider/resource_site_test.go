package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSite_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckSiteResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConfig,
			},
		},
	})
}

func testAccCheckSiteResourceDestroy(s *terraform.State) error {
	sites, err := testClient.ListSites(context.Background())
	if err != nil {
		return err
	}
	for _, site := range sites {
		if site.Description == "tfacc" {
			return fmt.Errorf("site tfacc not destroyed")
		}
	}
	return nil
}

const testAccSiteConfig = `
resource "unifi_site" "test" {
	description = "tfacc"
}
`

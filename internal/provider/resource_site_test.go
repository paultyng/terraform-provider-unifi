package provider

import (
	"context"
	"fmt"
	"strings"
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
				Config: testAccSiteConfig("tfacc-desc1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_site.test", "description", "tfacc-desc1"),
				),
			},
			importStep("unifi_site.test"),
			{
				Config: testAccSiteConfig("tfacc-desc2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_site.test", "description", "tfacc-desc2"),
				),
			},
			importStep("unifi_site.test"),
		},
	})
}

func testAccCheckSiteResourceDestroy(s *terraform.State) error {
	sites, err := testClient.ListSites(context.Background())
	if err != nil {
		return err
	}
	for _, site := range sites {
		if strings.HasPrefix(site.Description, "tfacc-") {
			return fmt.Errorf("site not destroyed")
		}
	}
	return nil
}

func testAccSiteConfig(desc string) string {
	return fmt.Sprintf(`
resource "unifi_site" "test" {
	description = %q
}
`, desc)
}

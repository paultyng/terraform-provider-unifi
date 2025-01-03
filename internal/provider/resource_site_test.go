package provider

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSite_basic(t *testing.T) {
	var siteName string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// FIXME causes flaky tests. See: https://github.com/paultyng/terraform-provider-unifi/issues/480
		//CheckDestroy:      testAccCheckSiteResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConfig("tfacc-desc1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_site.test", "description", "tfacc-desc1"),

					// extract siteName for future use
					func(s *terraform.State) error {
						siteName = s.RootModule().Resources["unifi_site.test"].Primary.Attributes["name"]
						return nil
					},
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

			// test importing from name, not id
			{
				ResourceName: "unifi_site.test",
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return siteName, nil
				},
				ImportState:       true,
				ImportStateVerify: true,
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

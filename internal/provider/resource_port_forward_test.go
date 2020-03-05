package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	//	"github.com/paultyng/go-unifi/unifi"
)

func TestAccPortForward_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers,
		PreCheck:  func() { preCheck(t) },
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccPortForwardConfig("22", false, "1.1.1.1", "22", "ssh fwd name"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_port_forward.test", "dst_port", "22"),
				),
			},
			{
				Config: testAccPortForwardConfig("22", false, "1.1.1.1", "8022", "ssh fwd name"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_port_forward.test", "fwd_port", "8022"),
				),
			},
			{
				Config: testAccPortForwardConfig("22", false, "1.1.1.1", "22", "ssh fwd name 2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_port_forward.test", "name", "ssh fwd name 2"),
				),
			},
		},
	})
}

func testAccPortForwardConfig(dst_port string, enabled bool, fwd string, fwd_port string, name string) string {
	return fmt.Sprintf(`
resource "unifi_port_forward" "test" {
	dst_port  = %s
	enabled = %t
	fwd = "%s"
	fwd_port = "%s"
	name = "%s"
}
`, dst_port, enabled, fwd, fwd_port, name)
}

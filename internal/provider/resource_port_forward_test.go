package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPortForward_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	name2 := acctest.RandomWithPrefix("tfacc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccPortForwardConfig("22", false, "10.1.1.1", "22", name),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_port_forward.test", "dst_port", "22"),
				),
			},
			importStep("unifi_port_forward.test"),
			{
				Config: testAccPortForwardConfig("22", false, "10.1.1.2", "8022", name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_port_forward.test", "fwd_port", "8022"),
					resource.TestCheckResourceAttr("unifi_port_forward.test", "fwd_ip", "10.1.1.2"),
				),
			},
			importStep("unifi_port_forward.test"),
			{
				Config: testAccPortForwardConfig("22", false, "10.1.1.1", "22", name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_port_forward.test", "name", name2),
				),
			},
			importStep("unifi_port_forward.test"),
		},
	})
}

func TestAccPortForward_src_ip(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccPortForwardConfigSrc("22", false, "10.1.1.1", "22", name, "192.168.1.0"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_port_forward.test", "dst_port", "22"),
				),
			},
			importStep("unifi_port_forward.test"),
		},
	})
}

func TestAccPortForward_src_cidr(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccPortForwardConfigSrc("22", false, "10.1.1.1", "22", name, "192.168.1.0/20"),
				Check: resource.ComposeTestCheckFunc(
					// testCheckNetworkExists(t, "name"),
					resource.TestCheckResourceAttr("unifi_port_forward.test", "dst_port", "22"),
				),
			},
			importStep("unifi_port_forward.test"),
		},
	})
}

func testAccPortForwardConfig(dstPort string, enabled bool, fwdIP, fwdPort, name string) string {
	return fmt.Sprintf(`
resource "unifi_port_forward" "test" {
	dst_port = %q
	enabled  = %t
	fwd_ip   = %q
	fwd_port = %q
	name     = %q
}
`, dstPort, enabled, fwdIP, fwdPort, name)
}

func testAccPortForwardConfigSrc(dstPort string, enabled bool, fwdIP, fwdPort, name, src string) string {
	return fmt.Sprintf(`
resource "unifi_port_forward" "test" {
	dst_port = %q
	enabled  = %t
	fwd_ip   = %q
	fwd_port = %q
	name     = %q
	src_ip   = %q
}
`, dstPort, enabled, fwdIP, fwdPort, name, src)
}

package provider

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStaticRoute_nextHop(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	network := &net.IPNet{
		IP:   net.IPv4(172, 17, 0, 0).To4(),
		Mask: net.IPv4Mask(255, 255, 0, 0),
	}
	distance := 1
	nextHop := net.IPv4(172, 16, 0, 1).To4()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRouteConfig_nextHop(name, network, distance, &nextHop),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_static_route.test", "type", "nexthop-route"),
					resource.TestCheckResourceAttr("unifi_static_route.test", "network", network.String()),
					resource.TestCheckResourceAttr("unifi_static_route.test", "name", name),
					resource.TestCheckResourceAttr("unifi_static_route.test", "distance", strconv.Itoa(distance)),
					resource.TestCheckResourceAttr("unifi_static_route.test", "next_hop", nextHop.String()),
				),
			},
			importStep("unifi_static_route.test"),
		},
	})
}

func TestAccStaticRoute_nextHop_ipv6(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	network := &net.IPNet{
		IP:   net.IP{0xfd, 0x6a, 0x37, 0xbe, 0xe3, 0x62, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
		Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}
	distance := 1
	nextHop := net.IP{0xfd, 0x6a, 0x37, 0xbe, 0xe3, 0x62, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRouteConfig_nextHop(name, network, distance, &nextHop),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_static_route.test", "type", "nexthop-route"),
					resource.TestCheckResourceAttr("unifi_static_route.test", "network", network.String()),
					resource.TestCheckResourceAttr("unifi_static_route.test", "name", name),
					resource.TestCheckResourceAttr("unifi_static_route.test", "distance", strconv.Itoa(distance)),
					resource.TestCheckResourceAttr("unifi_static_route.test", "next_hop", nextHop.String()),
				),
			},
			importStep("unifi_static_route.test"),
		},
	})
}

func TestAccStaticRoute_blackhole(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	network := &net.IPNet{
		IP:   net.IPv4(172, 18, 0, 0).To4(),
		Mask: net.IPv4Mask(255, 255, 0, 0),
	}
	distance := 1

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRouteConfig_blackhole(name, network, distance),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_static_route.test", "type", "blackhole"),
					resource.TestCheckResourceAttr("unifi_static_route.test", "network", network.String()),
					resource.TestCheckResourceAttr("unifi_static_route.test", "name", name),
					resource.TestCheckResourceAttr("unifi_static_route.test", "distance", strconv.Itoa(distance)),
				),
			},
			importStep("unifi_static_route.test"),
		},
	})
}

func TestAccStaticRoute_blackhole_ipv6(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	network := &net.IPNet{
		IP:   net.IP{0xfd, 0x6a, 0x37, 0xbe, 0xe3, 0x62, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
		Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}
	distance := 1

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRouteConfig_blackhole(name, network, distance),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_static_route.test", "type", "blackhole"),
					resource.TestCheckResourceAttr("unifi_static_route.test", "network", network.String()),
					resource.TestCheckResourceAttr("unifi_static_route.test", "name", name),
					resource.TestCheckResourceAttr("unifi_static_route.test", "distance", strconv.Itoa(distance)),
				),
			},
			importStep("unifi_static_route.test"),
		},
	})
}

func TestAccStaticRoute_interface(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	network := &net.IPNet{
		IP:   net.IPv4(172, 19, 0, 0).To4(),
		Mask: net.IPv4Mask(255, 255, 0, 0),
	}
	distance := 1
	networkInterface := "WAN2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRouteConfig_interface(name, network, distance, networkInterface),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_static_route.test", "type", "interface-route"),
					resource.TestCheckResourceAttr("unifi_static_route.test", "network", network.String()),
					resource.TestCheckResourceAttr("unifi_static_route.test", "name", name),
					resource.TestCheckResourceAttr("unifi_static_route.test", "distance", strconv.Itoa(distance)),
					resource.TestCheckResourceAttr("unifi_static_route.test", "interface", networkInterface),
				),
			},
			importStep("unifi_static_route.test"),
		},
	})
}

func TestAccStaticRoute_interface_ipv6(t *testing.T) {
	name := acctest.RandomWithPrefix("tfacc")
	network := &net.IPNet{
		IP:   net.IP{0xfd, 0x6a, 0x37, 0xbe, 0xe3, 0x62, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
		Mask: net.IPMask{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}
	distance := 1
	networkInterface := "WAN2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRouteConfig_interface(name, network, distance, networkInterface),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("unifi_static_route.test", "type", "interface-route"),
					resource.TestCheckResourceAttr("unifi_static_route.test", "network", network.String()),
					resource.TestCheckResourceAttr("unifi_static_route.test", "name", name),
					resource.TestCheckResourceAttr("unifi_static_route.test", "distance", strconv.Itoa(distance)),
					resource.TestCheckResourceAttr("unifi_static_route.test", "interface", networkInterface),
				),
			},
			importStep("unifi_static_route.test"),
		},
	})
}

func testAccStaticRouteConfig_nextHop(name string, network *net.IPNet, distance int, nextHop *net.IP) string {
	return fmt.Sprintf(`
resource "unifi_static_route" "test" {
	type     = "nexthop-route"
	network  = "%[2]s"
	name     = "%[1]s"
	distance = %[3]d
	next_hop = "%[4]s"
}
`, name, network, distance, nextHop)
}

func testAccStaticRouteConfig_blackhole(name string, network *net.IPNet, distance int) string {
	return fmt.Sprintf(`
resource "unifi_static_route" "test" {
	type     = "blackhole"
	network  = "%[2]s"
	name     = "%[1]s"
	distance = %[3]d
}
`, name, network, distance)
}

func testAccStaticRouteConfig_interface(name string, network *net.IPNet, distance int, networkInterface string) string {
	return fmt.Sprintf(`
resource "unifi_static_route" "test" {
	type      = "interface-route"
	network   = "%[2]s"
	name      = "%[1]s"
	distance  = %[3]d
	interface = "%[4]s"
}
`, name, network, distance, networkInterface)
}

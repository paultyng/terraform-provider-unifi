package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataNetwork_byName(t *testing.T) {
	defaultName := "Default"
	v, err := version.NewVersion(testClient.Version())
	if err != nil {
		t.Fatalf("error parsing version: %s", err)
	}
	if v.LessThan(controllerV7) {
		defaultName = "LAN"
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNetworkConfig_byName(defaultName),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
		},
	})
}

func TestAccDataNetwork_byID(t *testing.T) {
	defaultName := "Default"
	v, err := version.NewVersion(testClient.Version())
	if err != nil {
		t.Fatalf("error parsing version: %s", err)
	}
	if v.LessThan(controllerV7) {
		defaultName = "LAN"
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
		},
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config: testAccDataNetworkConfig_byID(defaultName),
				Check:  resource.ComposeTestCheckFunc(
				// testCheckNetworkExists(t, "name"),
				),
			},
		},
	})
}

func testAccDataNetworkConfig_byName(name string) string {
	return fmt.Sprintf(`
data "unifi_network" "lan" {
	name = %q
}
`, name)
}

func testAccDataNetworkConfig_byID(name string) string {
	return fmt.Sprintf(`
data "unifi_network" "lan" {
	name = %q
}

data "unifi_network" "lan_id" {
	id = data.unifi_network.lan.id
}
`, name)
}

package provider

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func preCheckSwitch(t *testing.T) {
	switchID := os.Getenv("UNIFI_TEST_SWITCH_ID")
	if switchID == "" {
		t.Skipf("UNIFI_TEST_SWITCH_ID not set")
	}
	switchMAC := os.Getenv("UNIFI_TEST_SWITCH_MAC")
	if switchMAC == "" {
		t.Skipf("UNIFI_TEST_SWITCH_MAC not set")
	}
	poePortList := os.Getenv("UNIFI_TEST_SWITCH_PORT_NUMBERS")
	if poePortList == "" {
		t.Skipf("UNIFI_TEST_SWITCH_PORT_NUMBERS is not set")
	}
	poePorts := strings.Split(poePortList, ",")
	if len(poePorts) < 2 {
		t.Skipf("At least 2 ports are required for testing.")
	}
}

func TestAccDevice_switch_basic(t *testing.T) {
	//switchID := os.Getenv("UNIFI_TEST_SWITCH_ID")
	switchMAC := os.Getenv("UNIFI_TEST_SWITCH_MAC")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckSwitch(t)
		},
		ProviderFactories: providerFactories,
		// TODO: CheckDestroy: ,
		Steps: []resource.TestStep{
			{
				Config:      testAccDeviceConfigEmpty(),
				ExpectError: regexp.MustCompile("no MAC address specified, please import the device using terraform import"),
			},

			{
				Config: testAccDeviceConfig(switchMAC),
				Check:  resource.ComposeTestCheckFunc(
				// TODO:
				),

				// this plan will be non-empty since ports will already be configured most likely
				ExpectNonEmptyPlan: true,
			},

			// import with ID
			importStep("unifi_device.test"),

			// import with mac
			{
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     switchMAC,
				ResourceName:      "unifi_device.test",
			},

			// TODO: update switch
			// TODO: test port overrides
		},
	})
}

func testAccDeviceConfigEmpty() string {
	return fmt.Sprintf(`
resource "unifi_device" "test" {
}
`)
}

func testAccDeviceConfig(mac string) string {
	return fmt.Sprintf(`
resource "unifi_device" "test" {
	mac = %q
}
`, mac)
}

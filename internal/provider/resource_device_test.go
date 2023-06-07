package provider

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/paultyng/go-unifi/unifi"
)

var (
	deviceInit sync.Once
	devicePool mapset.Set[*unifi.Device] = mapset.NewSet[*unifi.Device]()
)

func allocateDevice(t *testing.T) (*unifi.Device, func()) {
	ctx := context.Background()

	deviceInit.Do(func() {
		// The demo devices don't appear instantly when the controller starts.
		err := resource.RetryContext(ctx, 1*time.Minute, func() *resource.RetryError {
			devices, err := testClient.ListDevice(ctx, "default")
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("Error listing devices: %w", err))
			}

			if len(devices) == 0 {
				return resource.RetryableError(fmt.Errorf("No devices found"))
			}

			for _, device := range devices {
				if device.Type != "usw" {
					continue
				}

				// These devices aren't really switches.
				if device.Model == "USPRPS" || device.Model == "USPRPSP" || device.Model == "USPPDUHD" || device.Model == "USPPDUP" {
					continue
				}

				// The USW-Leaf is an EOL product and consistently fails to be adopted.
				if device.Model == "UDC48X6" {
					continue
				}

				// Only switches with these chipsets support both port mirroring ang aggregation.
				if !(isBroadcomSwitch(device) || isMicrosemiSwitch(device) || isNephosSwitch(device)) {
					continue
				}

				d := device
				if ok := devicePool.Add(&d); !ok {
					return resource.NonRetryableError(fmt.Errorf("Failed to add device to pool"))
				}
			}

			return nil
		})

		if err != nil {
			t.Fatal(err)
		}
	})

	var device *unifi.Device

	err := resource.RetryContext(ctx, 1*time.Minute, func() *resource.RetryError {
		var ok bool
		device, ok = devicePool.Pop()

		if device == nil || !ok {
			return resource.RetryableError(fmt.Errorf("Unable to allocate test device"))
		}

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	unallocate := func() {
		if ok := devicePool.Add(device); !ok {
			t.Fatal("Failed to add device to pool")
		}
	}

	return device, unallocate
}

func isBroadcomSwitch(device unifi.Device) bool {
	if device.Type != "usw" {
		return false
	}

	switch device.Model {
	// US-8 variants
	case "US8", "US8P60", "US8P150", "S28150":
		return true

	// US-16 variants
	case "US16P150", "S216150", "USXG":
		return true

	// US-24 variants
	case "US24", "US24P250", "S224250", "US24P500", "S224500", "US24PL2":
		return true

	// US-48 variants
	case "US48", "US48P500", "S248500", "US48P750", "S248750", "US48PL2":
		return true

	// USW-Pro
	case "US24PRO", "US24PRO2", "US48PRO", "US48PRO2", "USAGGPRO":
		return true

		// USW-Enterprise
	case "US624P", "US648P", "USXG24":
		return true

	// US-XG-6PoE
	case "US6XG150":
		return true
	}

	return false
}

func isMicrosemiSwitch(device unifi.Device) bool {
	if device.Type != "usw" {
		return false
	}

	switch device.Model {
	// US-8 variants
	case "USC8", "USC8P60", "USC8P150":
		return true

	// USW-Industrial
	case "USC8P450":
		return true
	}

	return false
}

func isNephosSwitch(device unifi.Device) bool {
	if device.Type != "usw" {
		return false
	}

	switch device.Model {
	// USW-Leaf
	case "UDC48X6":
		return true
	}

	return false
}

func preCheckDeviceExists(t *testing.T, site, mac string) {
	_, err := testClient.GetDeviceByMAC(context.Background(), site, mac)

	if _, ok := err.(*unifi.NotFoundError); ok {
		t.Fatal("Test device not found")
	}
}

func TestAccDevice_empty(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDeviceConfigEmpty(),
				ExpectError: regexp.MustCompile(`no MAC address specified, please import the device using terraform import`),
			},
		},
	})
}

func TestAccDevice_switch_basic(t *testing.T) {
	resourceName := "unifi_device.test"
	site := "default"

	device, unallocateDevice := allocateDevice(t)
	defer unallocateDevice()

	importStateVerifyIgnore := []string{"allow_adoption", "forget_on_destroy"}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckDeviceExists(t, site, device.MAC)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceConfig(device.MAC),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site", site),
					resource.TestCheckResourceAttr(resourceName, "mac", device.MAC),
					resource.TestCheckResourceAttr(resourceName, "name", ""),
				),
			},

			// Import with ID
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
			},

			// Import with MAC
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateId:           device.MAC,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
			},

			{
				Config: testAccDeviceConfig_withName(device.MAC, "Test Switch"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", "Test Switch"),
				),
			},
		},
	})
}

func TestAccDevice_switch_portOverrides(t *testing.T) {
	resourceName := "unifi_device.test"
	site := "default"

	device, unallocateDevice := allocateDevice(t)
	defer unallocateDevice()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckDeviceExists(t, site, device.MAC)
			preCheckVersionConstraint(t, "< 7.4")
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceConfig_withPortOverrides(device.MAC),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "port_override.#", "3"),

					// TODO: Why are these out of order?
					resource.TestCheckResourceAttr(resourceName, "port_override.0.number", "3"),
					resource.TestCheckResourceAttr(resourceName, "port_override.0.name", ""),
					resource.TestCheckResourceAttr(resourceName, "port_override.0.port_profile_id", ""),
					resource.TestCheckResourceAttr(resourceName, "port_override.0.op_mode", "aggregate"),
					resource.TestCheckResourceAttr(resourceName, "port_override.0.aggregate_num_ports", "2"),

					resource.TestCheckResourceAttr(resourceName, "port_override.1.number", "1"),
					resource.TestCheckResourceAttr(resourceName, "port_override.1.name", "Port 1"),
					resource.TestCheckResourceAttr(resourceName, "port_override.1.port_profile_id", ""),
					//resource.TestCheckResourceAttr(resourceName, "port_override.1.op_mode", "switch"),

					resource.TestCheckResourceAttr(resourceName, "port_override.2.number", "2"),
					resource.TestCheckResourceAttr(resourceName, "port_override.2.name", "Port 2"),
					//resource.TestCheckResourceAttr(resourceName, "port_override.2.port_profile_id", ""),
					//resource.TestCheckResourceAttr(resourceName, "port_override.2.op_mode", "switch"),
				),
			},
			{
				Config: testAccDeviceConfig(device.MAC),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "port_override.#", "0"),
				),
			},
		},
	})
}

func testAccDeviceConfigEmpty() string {
	return `
resource "unifi_device" "test" {}
`
}

func testAccDeviceConfig(mac string) string {
	return fmt.Sprintf(`
resource "unifi_device" "test" {
	mac = %q
}
`, mac)
}

func testAccDeviceConfig_withName(mac, name string) string {
	return fmt.Sprintf(`
resource "unifi_device" "test" {
	mac  = %q
	name = %q
}
`, mac, name)
}

func testAccDeviceConfig_withPortOverrides(mac string) string {
	return fmt.Sprintf(`
data "unifi_port_profile" "all" {}

resource "unifi_device" "test" {
	mac = %q

	port_override {
		number = 1
		name   = "Port 1"
	}

	port_override {
		number          = 2
		name            = "Port 2"
		port_profile_id = data.unifi_port_profile.all.id
		op_mode         = "switch"
	}

	port_override {
		number              = 3
		op_mode             = "aggregate"
		aggregate_num_ports = 2
	}
}
`, mac)
}

func testAccCheckDeviceDestroy(s *terraform.State) error {
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "unifi_device" {
			continue
		}

		device, err := testClient.GetDevice(ctx, rs.Primary.Attributes["site"], rs.Primary.ID)
		if device != nil {
			return fmt.Errorf("Device still exists with ID %v", rs.Primary.ID)
		}
		if _, ok := err.(*unifi.NotFoundError); !ok {
			return err
		}
	}

	return nil
}

func testAccCheckDeviceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		id := rs.Primary.ID
		site := rs.Primary.Attributes["site"]

		device, err := testClient.GetDevice(context.Background(), site, id)
		if device == nil {
			return fmt.Errorf("Device not found with ID %v", id)
		}
		if _, ok := err.(*unifi.NotFoundError); !ok {
			return err
		}

		return nil
	}
}

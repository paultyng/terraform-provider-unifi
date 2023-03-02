package provider

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/paultyng/go-unifi/unifi"
)

var (
	deviceInit  sync.Once
	deviceMutex sync.Mutex
	devicePool  []string = []string{}
)

func allocateDevice(t *testing.T) (string, func()) {
	ctx := context.Background()

	deviceInit.Do(func() {
		devices, err := testClient.ListDevice(ctx, "default")
		if err != nil {
			t.Fatalf("Error listing devices: %s", err)
		}

		for _, device := range devices {
			if device.Type != "usw" {
				continue
			}

			// These devices aren't really switches.
			if device.Model == "USPRPS" || device.Model == "USPRPSP" || device.Model == "USPPDUHD" || device.Model == "USPPDUP" {
				continue
			}

			t.Logf("Discovered device %s\n", device.MAC)
			devicePool = append(devicePool, device.MAC)
		}
	})

	var device string

	err := resource.RetryContext(ctx, 1*time.Minute, func() *resource.RetryError {
		deviceMutex.Lock()
		defer deviceMutex.Unlock()

		device, devicePool = devicePool[0], devicePool[1:]

    if len(devicePool) == 0 {
      return resource.RetryableError(fmt.Errorf("Unable to allocate test device"))
    }

		t.Logf("Allocated device %s. Device pool = #%v\n", device, devicePool)
		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	unallocate := func() {
		deviceMutex.Lock()
		defer deviceMutex.Unlock()
		devicePool = append(devicePool, device)
		t.Logf("Unallocated device %s\n", device)
	}

	return device, unallocate
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

	switchMAC, unallocateDevice := allocateDevice(t)
	defer unallocateDevice()

	importStateVerifyIgnore := []string{"allow_adoption", "forget_on_destroy"}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckDeviceExists(t, site, switchMAC)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceConfig(switchMAC),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "site", site),
					resource.TestCheckResourceAttr(resourceName, "mac", switchMAC),
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
				ImportStateId:           switchMAC,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: importStateVerifyIgnore,
			},

			{
				Config: testAccDeviceConfig_withName(switchMAC, "Test Switch"),
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

	switchMAC, unallocateDevice := allocateDevice(t)
	defer unallocateDevice()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckDeviceExists(t, site, switchMAC)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceConfig_withPortOverrides(switchMAC),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDeviceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "port_override.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "port_override.0.number", "1"),
					resource.TestCheckResourceAttr(resourceName, "port_override.0.name", "Port 1"),
					resource.TestCheckResourceAttr(resourceName, "port_override.1.number", "2"),
					resource.TestCheckResourceAttr(resourceName, "port_override.1.name", "Port 2"),
				),
			},
		},
	})
}

func TestAccDevice_remove_portOverrides(t *testing.T) {
	resourceName := "unifi_device.test"
	site := "default"

	switchMAC, unallocateDevice := allocateDevice(t)
	defer unallocateDevice()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			preCheckDeviceExists(t, site, switchMAC)
		},
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckDeviceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceConfig_withPortOverrides(switchMAC),
			},
			{
				Config: testAccDeviceConfig(switchMAC),
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
resource "unifi_device" "test" {
	mac = %q

	port_override {
		number = 1
		name   = "Port 1"
	}

	port_override {
		number = 2
		name   = "Port 2"
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

package provider

import (
	mapset "github.com/deckarep/golang-set/v2"
	"net"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var macAddressRegexp = regexp.MustCompile("^([0-9a-fA-F][0-9a-fA-F][-:]){5}([0-9a-fA-F][0-9a-fA-F])$")

func cleanMAC(mac string) string {
	return strings.TrimSpace(strings.ReplaceAll(strings.ToLower(mac), "-", ":"))
}

func macDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	old = cleanMAC(old)
	new = cleanMAC(new)
	return old == new
}

var (
	macInit sync.Once
	macPool = mapset.NewSet[*net.HardwareAddr]()
)

func allocateTestMac(t *testing.T) (string, func()) {
	macInit.Do(func() {
		// for test MAC addresses, see https://tools.ietf.org/html/rfc7042#section-2.1.
		for i := 0; i < 512; i++ {
			mac := net.HardwareAddr{0x00, 0x00, 0x5e, 0x00, 0x53, byte(i)}
			if ok := macPool.Add(&mac); !ok {
				t.Fatal("Failed to add MAC to pool")
			}
		}
	})

	mac, ok := macPool.Pop()
	if mac == nil || !ok {
		t.Fatal("Unable to allocate test MAC")
	}

	unallocate := func() {
		if ok := macPool.Add(mac); !ok {
			t.Fatal("Failed to add MAC to pool")
		}
	}

	return mac.String(), unallocate
}

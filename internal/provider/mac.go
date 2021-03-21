package provider

import (
	"regexp"
	"strings"

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

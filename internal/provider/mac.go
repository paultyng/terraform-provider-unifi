package provider

import (
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var macAddressRegexp = regexp.MustCompile("^([0-9a-fA-F][0-9a-fA-F][-:]){5}([0-9a-fA-F][0-9a-fA-F])$")

func macDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	old = strings.TrimSpace(strings.ReplaceAll(strings.ToLower(old), "-", ":"))
	new = strings.TrimSpace(strings.ReplaceAll(strings.ToLower(new), "-", ":"))
	return old == new
}

package provider

import (
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var timeOfDayRegexp = regexp.MustCompile("^\\d{1,2}:\\d{2}$")

func timeOfDayDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	return timeFromConfig(old) == timeFromConfig(new)
}

func timeFromConfig(t string) string {
	if len(t) == 0 {
		return ""
	}
	s := "0" + strings.ReplaceAll(t, ":", "")
	return s[len(s)-4:]
}

func timeFromUnifi(t string) string {
	i := len(t) - 2
	return strings.TrimPrefix(t[0:i], "0") + ":" + t[i:]
}

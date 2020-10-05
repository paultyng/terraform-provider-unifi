package provider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	wanUsernameRegexp   = regexp.MustCompile("[^\"' ]+")
	validateWANUsername = validation.StringMatch(wanUsernameRegexp, "invalid WAN username")
)


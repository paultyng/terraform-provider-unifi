package provider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	wanPasswordRegexp   = regexp.MustCompile("[^\"' ]+")
	validateWANPassword = validation.StringMatch(wanPasswordRegexp, "invalid WAN password")
)


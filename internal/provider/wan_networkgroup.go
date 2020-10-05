package provider

import (
"regexp"

"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	wanNetworkGroupRegexp   = regexp.MustCompile("WAN[2]?|WAN_LTE_FAILOVER")
	validateWANNetworkGroup = validation.StringMatch(wanNetworkGroupRegexp, "invalid WAN network group")
)

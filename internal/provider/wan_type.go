package provider

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	wanTypeRegexp   = regexp.MustCompile("disabled|dhcp|static|pppoe")
	validateWANType = validation.StringMatch(wanTypeRegexp, "invalid WAN password")
)


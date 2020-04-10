package provider

import (
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func cidrDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	_, oldNet, err := net.ParseCIDR(old)
	if err != nil {
		return false
	}

	_, newNet, err := net.ParseCIDR(new)
	if err != nil {
		return false
	}

	return oldNet.String() == newNet.String()
}

func cidrZeroBased(cidr string) string {
	_, cidrNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return ""
	}
	if len(cidrNet.Mask) == net.IPv6len {
		return ""
	}

	return cidrNet.String()
}

func cidrOneBased(cidr string) string {
	_, cidrNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return ""
	}
	if len(cidrNet.Mask) == net.IPv6len {
		return ""
	}

	cidrNet.IP[3]++

	return cidrNet.String()
}

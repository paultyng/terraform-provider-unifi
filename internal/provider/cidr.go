package provider

import (
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func cidrValidate(raw any, key string) ([]string, []error) {
	v, ok := raw.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected string, got %T", raw)}
	}

	_, _, err := net.ParseCIDR(v)
	if err != nil {
		return nil, []error{err}
	}

	return nil, nil
}

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

	return cidrNet.String()
}

func cidrOneBased(cidr string) string {
	_, cidrNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return ""
	}

	cidrNet.IP[3]++

	return cidrNet.String()
}

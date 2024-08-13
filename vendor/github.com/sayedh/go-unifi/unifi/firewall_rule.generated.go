// Code generated from ace.jar fields *.json files
// DO NOT EDIT.

package unifi

import (
	"context"
	"encoding/json"
	"fmt"
)

// just to fix compile issues with the import
var (
	_ context.Context
	_ fmt.Formatter
	_ json.Marshaler
)

type FirewallRule struct {
	ID     string `json:"_id,omitempty"`
	SiteID string `json:"site_id,omitempty"`

	Hidden   bool   `json:"attr_hidden,omitempty"`
	HiddenID string `json:"attr_hidden_id,omitempty"`
	NoDelete bool   `json:"attr_no_delete,omitempty"`
	NoEdit   bool   `json:"attr_no_edit,omitempty"`

	Action                string   `json:"action,omitempty"` // drop|reject|accept
	Contiguous            bool     `json:"contiguous"`
	DstAddress            string   `json:"dst_address,omitempty"`
	DstAddressIPV6        string   `json:"dst_address_ipv6,omitempty"`
	DstFirewallGroupIDs   []string `json:"dst_firewallgroup_ids,omitempty"` // [\d\w]+
	DstNetworkID          string   `json:"dst_networkconf_id"`              // [\d\w]+|^$
	DstNetworkType        string   `json:"dst_networkconf_type,omitempty"`  // ADDRv4|NETv4
	DstPort               string   `json:"dst_port,omitempty"`
	Enabled               bool     `json:"enabled"`
	ICMPTypename          string   `json:"icmp_typename"`   // ^$|address-mask-reply|address-mask-request|any|communication-prohibited|destination-unreachable|echo-reply|echo-request|fragmentation-needed|host-precedence-violation|host-prohibited|host-redirect|host-unknown|host-unreachable|ip-header-bad|network-prohibited|network-redirect|network-unknown|network-unreachable|parameter-problem|port-unreachable|precedence-cutoff|protocol-unreachable|redirect|required-option-missing|router-advertisement|router-solicitation|source-quench|source-route-failed|time-exceeded|timestamp-reply|timestamp-request|TOS-host-redirect|TOS-host-unreachable|TOS-network-redirect|TOS-network-unreachable|ttl-zero-during-reassembly|ttl-zero-during-transit
	ICMPv6Typename        string   `json:"icmpv6_typename"` // ^$|address-unreachable|bad-header|beyond-scope|communication-prohibited|destination-unreachable|echo-reply|echo-request|failed-policy|neighbor-advertisement|neighbor-solicitation|no-route|packet-too-big|parameter-problem|port-unreachable|redirect|reject-route|router-advertisement|router-solicitation|time-exceeded|ttl-zero-during-reassembly|ttl-zero-during-transit|unknown-header-type|unknown-option
	IPSec                 string   `json:"ipsec"`           // match-ipsec|match-none|^$
	Logging               bool     `json:"logging"`
	MonthDays             string   `json:"monthdays"` // ^$|^(([1-9]|[12][0-9]|3[01])(,([1-9]|[12][0-9]|3[01])){0,30})$
	MonthDaysNegate       bool     `json:"monthdays_negate"`
	Name                  string   `json:"name,omitempty"` // .{1,128}
	Protocol              string   `json:"protocol"`       // ^$|all|([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])|tcp_udp|ah|ax.25|dccp|ddp|egp|eigrp|encap|esp|etherip|fc|ggp|gre|hip|hmp|icmp|idpr-cmtp|idrp|igmp|igp|ip|ipcomp|ipencap|ipip|ipv6|ipv6-frag|ipv6-icmp|ipv6-nonxt|ipv6-opts|ipv6-route|isis|iso-tp4|l2tp|manet|mobility-header|mpls-in-ip|ospf|pim|pup|rdp|rohc|rspf|rsvp|sctp|shim6|skip|st|tcp|udp|udplite|vmtp|vrrp|wesp|xns-idp|xtp
	ProtocolMatchExcepted bool     `json:"protocol_match_excepted"`
	ProtocolV6            string   `json:"protocol_v6"`                  // ^$|([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])|ah|all|dccp|eigrp|esp|gre|icmpv6|ipcomp|ipv6|ipv6-frag|ipv6-icmp|ipv6-nonxt|ipv6-opts|ipv6-route|isis|l2tp|manet|mobility-header|mpls-in-ip|ospf|pim|rsvp|sctp|shim6|tcp|tcp_udp|udp|vrrp
	RuleIndex             int      `json:"rule_index,omitempty"`         // 2[0-9]{3,4}|4[0-9]{3,4}
	Ruleset               string   `json:"ruleset,omitempty"`            // WAN_IN|WAN_OUT|WAN_LOCAL|LAN_IN|LAN_OUT|LAN_LOCAL|GUEST_IN|GUEST_OUT|GUEST_LOCAL|WANv6_IN|WANv6_OUT|WANv6_LOCAL|LANv6_IN|LANv6_OUT|LANv6_LOCAL|GUESTv6_IN|GUESTv6_OUT|GUESTv6_LOCAL
	SettingPreference     string   `json:"setting_preference,omitempty"` // auto|manual
	SrcAddress            string   `json:"src_address,omitempty"`
	SrcAddressIPV6        string   `json:"src_address_ipv6,omitempty"`
	SrcFirewallGroupIDs   []string `json:"src_firewallgroup_ids,omitempty"` // [\d\w]+
	SrcMACAddress         string   `json:"src_mac_address"`                 // ^([0-9A-Fa-f]{2}:){5}([0-9A-Fa-f]{2})$|^$
	SrcNetworkID          string   `json:"src_networkconf_id"`              // [\d\w]+|^$
	SrcNetworkType        string   `json:"src_networkconf_type,omitempty"`  // ADDRv4|NETv4
	SrcPort               string   `json:"src_port,omitempty"`
	StartDate             string   `json:"startdate"` // ^$|^(20[0-9]{2}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])T([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])$
	StartTime             string   `json:"starttime"` // ^$|^(([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])$
	StateEstablished      bool     `json:"state_established"`
	StateInvalid          bool     `json:"state_invalid"`
	StateNew              bool     `json:"state_new"`
	StateRelated          bool     `json:"state_related"`
	StopDate              string   `json:"stopdate"` // ^$|^(20[0-9]{2}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])T([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])$
	StopTime              string   `json:"stoptime"` // ^$|^(([01][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9])$
	UTC                   bool     `json:"utc"`
	Weekdays              string   `json:"weekdays"` // ^$|^((Mon|Tue|Wed|Thu|Fri|Sat|Sun)(,(Mon|Tue|Wed|Thu|Fri|Sat|Sun)){0,6})$
	WeekdaysNegate        bool     `json:"weekdays_negate"`
}

func (dst *FirewallRule) UnmarshalJSON(b []byte) error {
	type Alias FirewallRule
	aux := &struct {
		RuleIndex emptyStringInt `json:"rule_index"`

		*Alias
	}{
		Alias: (*Alias)(dst),
	}

	err := json.Unmarshal(b, &aux)
	if err != nil {
		return fmt.Errorf("unable to unmarshal alias: %w", err)
	}
	dst.RuleIndex = int(aux.RuleIndex)

	return nil
}

func (c *Client) listFirewallRule(ctx context.Context, site string) ([]FirewallRule, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []FirewallRule `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/firewallrule", site), nil, &respBody)
	if err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (c *Client) getFirewallRule(ctx context.Context, site, id string) (*FirewallRule, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []FirewallRule `json:"data"`
	}

	err := c.do(ctx, "GET", fmt.Sprintf("s/%s/rest/firewallrule/%s", site, id), nil, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	d := respBody.Data[0]
	return &d, nil
}

func (c *Client) deleteFirewallRule(ctx context.Context, site, id string) error {
	err := c.do(ctx, "DELETE", fmt.Sprintf("s/%s/rest/firewallrule/%s", site, id), struct{}{}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) createFirewallRule(ctx context.Context, site string, d *FirewallRule) (*FirewallRule, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []FirewallRule `json:"data"`
	}

	err := c.do(ctx, "POST", fmt.Sprintf("s/%s/rest/firewallrule", site), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

func (c *Client) updateFirewallRule(ctx context.Context, site string, d *FirewallRule) (*FirewallRule, error) {
	var respBody struct {
		Meta meta           `json:"meta"`
		Data []FirewallRule `json:"data"`
	}

	err := c.do(ctx, "PUT", fmt.Sprintf("s/%s/rest/firewallrule/%s", site, d.ID), d, &respBody)
	if err != nil {
		return nil, err
	}

	if len(respBody.Data) != 1 {
		return nil, &NotFoundError{}
	}

	new := respBody.Data[0]

	return &new, nil
}

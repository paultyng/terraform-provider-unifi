---
page_title: "unifi_firewall_rule Resource - terraform-provider-unifi"
subcategory: ""
description: |-
  unifi_firewall_rule manages an individual firewall rule on the gateway.
---

# Resource `unifi_firewall_rule`

`unifi_firewall_rule` manages an individual firewall rule on the gateway.

## Example Usage

```terraform
variable "ip_address" {
  type = string
}

resource "unifi_firewall_rule" "drop_all" {
  name    = "drop all"
  action  = "drop"
  ruleset = "LAN_IN"

  rule_index = 2011

  protocol = "all"

  dst_address = var.ip_address
}
```

## Schema

### Required

- **action** (String, Required) The action of the firewall rule. Must be one of `drop`, `accept`, or `reject`.
- **name** (String, Required) The name of the firewall rule.
- **protocol** (String, Required) The protocol of the rule.
- **rule_index** (Number, Required) The index of the rule. Must be >= 2000 < 3000 or >= 4000 < 5000.
- **ruleset** (String, Required) The ruleset for the rule. This is from the perspective of the security gateway. Must be one of `WAN_IN`, `WAN_OUT`, `WAN_LOCAL`, `LAN_IN`, `LAN_OUT`, `LAN_LOCAL`, `GUEST_IN`, `GUEST_OUT`, `GUEST_LOCAL`, `WANv6_IN`, `WANv6_OUT`, `WANv6_LOCAL`, `LANv6_IN`, `LANv6_OUT`, `LANv6_LOCAL`, `GUESTv6_IN`, `GUESTv6_OUT`, or `GUESTv6_LOCAL`.

### Optional

- **dst_address** (String, Optional) The destination address of the firewall rule.
- **dst_firewall_group_ids** (Set of String, Optional) The destination firewall group IDs of the firewall rule.
- **dst_network_id** (String, Optional) The destination network ID of the firewall rule.
- **dst_network_type** (String, Optional) The destination network type of the firewall rule. Can be one of `ADDRv4` or `NETv4`. Defaults to `NETv4`.
- **id** (String, Optional) The ID of this resource.
- **ip_sec** (String, Optional) Specify whether the rule matches on IPsec packets. Can be one of `match-ipset` or `match-none`.
- **logging** (Boolean, Optional) Enable logging for the firewall rule.
- **src_address** (String, Optional) The source address for the firewall rule.
- **src_firewall_group_ids** (Set of String, Optional) The source firewall group IDs for the firewall rule.
- **src_mac** (String, Optional) The source MAC address of the firewall rule.
- **src_network_id** (String, Optional) The source network ID for the firewall rule.
- **src_network_type** (String, Optional) The source network type of the firewall rule. Can be one of `ADDRv4` or `NETv4`. Defaults to `NETv4`.
- **state_established** (Boolean, Optional) Match where the state is established.
- **state_invalid** (Boolean, Optional) Match where the state is invalid.
- **state_new** (Boolean, Optional) Match where the state is new.
- **state_related** (Boolean, Optional) Match where the state is related.



---
subcategory: ""
layout: ""
page_title: "terraform-provider-unifi: unifi_firewall_group"
description: |-
  unifi_firewall_group manages groups of addresses or ports for use in firewall rules (unifi_firewall_rule).
---

# Resource: `unifi_firewall_group`

unifi_firewall_group manages groups of addresses or ports for use in firewall rules (unifi_firewall_rule).



## Schema

### Required

- **members** (Set of String, Required) The members of the firewall group.
- **name** (String, Required) The name of the firewall group.
- **type** (String, Required) The type of the firewall group. Must be one of: `address-group`, `port-group`, or `ipv6-address-group`.

### Optional

- **id** (String, Optional)



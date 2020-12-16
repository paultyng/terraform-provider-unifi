---
page_title: "unifi_firewall_group Resource - terraform-provider-unifi"
subcategory: ""
description: |-
  unifi_firewall_group manages groups of addresses or ports for use in firewall rules (unifi_firewall_rule).
---

# Resource `unifi_firewall_group`

`unifi_firewall_group` manages groups of addresses or ports for use in firewall rules (`unifi_firewall_rule`).

## Example Usage

```terraform
variable "laptop_ips" {
  type = list(string)
}

resource "unifi_firewall_group" "can_print" {
  name = "can-print"
  type = "address-group"

  members = var.laptop_ips
}
```

## Schema

### Required

- **members** (Set of String) The members of the firewall group.
- **name** (String) The name of the firewall group.
- **type** (String) The type of the firewall group. Must be one of: `address-group`, `port-group`, or `ipv6-address-group`.

### Optional

- **site** (String) The name of the site to associate the firewall group with.

### Read-only

- **id** (String) The ID of the firewall group.



---
page_title: "unifi_network Resource - terraform-provider-unifi"
subcategory: ""
description: |-
  unifi_network manages LAN/VLAN networks.
---

# Resource `unifi_network`

`unifi_network` manages LAN/VLAN networks.

## Example Usage

```terraform
variable "vlan_id" {
  default = 10
}

resource "unifi_network" "vlan" {
  name    = "wifi-vlan"
  purpose = "corporate"

  subnet       = "10.0.0.1/24"
  vlan_id      = var.vlan_id
  dhcp_start   = "10.0.0.6"
  dhcp_stop    = "10.0.0.254"
  dhcp_enabled = true
}
```

## Schema

### Required

- **name** (String, Required) The name of the network.
- **purpose** (String, Required) The purpose of the network. Must be one of `corporate`, `guest`, or `vlan-only`.

### Optional

- **dhcp_dns** (List of String, Optional) Specifies the IPv4 addresses for the DNS server to be returned from the DHCP server. Leave blank to disable this feature.
- **dhcp_enabled** (Boolean, Optional) Specifies whether DHCP is enabled or not on this network.
- **dhcp_lease** (Number, Optional) Specifies the lease time for DHCP addresses. Defaults to `86400`.
- **dhcp_start** (String, Optional) The IPv4 address where the DHCP range of addresses starts.
- **dhcp_stop** (String, Optional) The IPv4 address where the DHCP range of addresses stops.
- **domain_name** (String, Optional) The domain name of this network.
- **igmp_snooping** (Boolean, Optional) Specifies whether IGMP snooping is enabled or not.
- **ipv6_interface_type** (String, Optional) Specifies which type of IPv6 connection to use. Defaults to `none`.
- **ipv6_pd_interface** (String, Optional) Specifies which WAN interface to use for IPv6 PD.
- **ipv6_pd_prefixid** (String, Optional) Specifies the IPv6 Prefix ID.
- **ipv6_ra_enable** (Boolean, Optional) Specifies whether to enable router advertisements or not.
- **ipv6_static_subnet** (String, Optional) Specifies the static IPv6 subnet when ipv6_interface_type is 'static'.
- **network_group** (String, Optional) The group of the network. Defaults to `LAN`.
- **subnet** (String, Optional) The subnet of the network. Must be a valid CIDR address.
- **vlan_id** (Number, Optional) The VLAN ID of the network.

### Read-only

- **id** (String, Read-only) The ID of the network.



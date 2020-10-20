---
page_title: "unifi_wlan Resource - terraform-provider-unifi"
subcategory: ""
description: |-
  unifi_wlan manages a WiFi network / SSID.
---

# Resource `unifi_wlan`

`unifi_wlan` manages a WiFi network / SSID.

## Example Usage

```terraform
variable "vlan_id" {
  default = 10
}

data "unifi_ap_group" "default" {
}

data "unifi_user_group" "default" {
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

resource "unifi_wlan" "wifi" {
  name       = "myssid"
  passphrase = "12345678"
  security   = "wpapsk"

  network_id    = unifi_network.vlan.id
  ap_group_ids  = [data.unifi_ap_group.default.id]
  user_group_id = data.unifi_user_group.default.id
}
```

## Schema

### Required

- **name** (String, Required) The SSID of the network.
- **security** (String, Required) The type of WiFi security for this network. Valid values are: `wpapsk`, `wpaeap`, and `open`.
- **user_group_id** (String, Required) ID of the user group to use for this network.

### Optional

- **ap_group_ids** (Set of String, Optional) IDs of the AP groups to use for this network.
- **hide_ssid** (Boolean, Optional) Indicates whether or not to hide the SSID from broadcast.
- **is_guest** (Boolean, Optional) Indicates that this is a guest WLAN and should use guest behaviors.
- **mac_filter_enabled** (Boolean, Optional) Indicates whether or not the MAC filter is turned of for the network.
- **mac_filter_list** (Set of String, Optional) List of MAC addresses to filter (only valid if `mac_filter_enabled` is `true`).
- **mac_filter_policy** (String, Optional) MAC address filter policy (only valid if `mac_filter_enabled` is `true`). Defaults to `deny`.
- **multicast_enhance** (Boolean, Optional) Indicates whether or not Multicast Enhance is turned of for the network.
- **network_id** (String, Optional) ID of the network for this SSID
- **passphrase** (String, Optional) The passphrase for the network, this is only required if `security` is not set to `open`.
- **radius_profile_id** (String, Optional) ID of the RADIUS profile to use when security `wpaeap`. You can query this via the `unifi_radius_profile` data source.
- **schedule** (Block List) Start and stop schedules for the WLAN (see [below for nested schema](#nestedblock--schedule))
- **vlan_id** (Number, Optional, Deprecated) VLAN ID for the network. Set network_id instead of vlan_id for controller version >= 6.
- **wlan_group_id** (String, Optional, Deprecated) ID of the WLAN group to use for this network. Set ap_group_ids instead of wlan_group_id for controller version >= 6.

### Read-only

- **id** (String, Read-only) The ID of the network.

<a id="nestedblock--schedule"></a>
### Nested Schema for `schedule`

Required:

- **block_end** (String, Required) Time of day to end the block.
- **block_start** (String, Required) Time of day to start the block.
- **day_of_week** (String, Required) Day of week for the block. Valid values are `sun`, `mon`, `tue`, `wed`, `thu`, `fri`, `sat`.



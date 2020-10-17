---
page_title: "unifi_user Resource - terraform-provider-unifi"
subcategory: ""
description: |-
  unifi_user manages a user (or "client" in the UI) of the network, these are identified by unique MAC addresses.
  Users are created in the controller when observed on the network, so the resource defaults to allowing itself to just take over management of a MAC address, but this can be turned off.
---

# Resource `unifi_user`

`unifi_user` manages a user (or "client" in the UI) of the network, these are identified by unique MAC addresses.

Users are created in the controller when observed on the network, so the resource defaults to allowing itself to just take over management of a MAC address, but this can be turned off.

## Example Usage

```terraform
resource "unifi_user" "test" {
  mac  = "01:23:45:67:89:AB"
  name = "some client"
  note = "my note"

  fixed_ip   = "10.1.10.50"
  network_id = unifi_network.my_vlan.id
}
```

## Schema

### Required

- **mac** (String, Required) The MAC address of the user.
- **name** (String, Required) The name of the user.

### Optional

- **allow_existing** (Boolean, Optional) Specifies whether this resource should just take over control of an existing user. Defaults to `true`.
- **blocked** (Boolean, Optional) Specifies whether this user should be blocked from the network.
- **fixed_ip** (String, Optional) A fixed IPv4 address for this user.
- **network_id** (String, Optional) The network ID for this user.
- **note** (String, Optional) A note with additional information for the user.
- **skip_forget_on_destroy** (Boolean, Optional) Specifies whether this resource should tell the controller to "forget" the user on destroy. Defaults to `false`.
- **user_group_id** (String, Optional) The user group ID for the user.

### Read-only

- **hostname** (String, Read-only) The hostname of the user.
- **id** (String, Read-only) The ID of the user.
- **ip** (String, Read-only) The IP address of the user.



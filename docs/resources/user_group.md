---
page_title: "unifi_user_group Resource - terraform-provider-unifi"
subcategory: ""
description: |-
  unifi_user_group manages a user group (called "client group" in the UI), which can be used to limit bandwidth for groups of users.
---

# Resource `unifi_user_group`

`unifi_user_group` manages a user group (called "client group" in the UI), which can be used to limit bandwidth for groups of users.

## Example Usage

```terraform
resource "unifi_user_group" "wifi" {
  name = "wifi"

  qos_rate_max_down = 2000 # 2mbps
  qos_rate_max_up   = 10   # 10kbps
}
```

## Schema

### Required

- **name** (String, Required) The name of the user group.

### Optional

- **id** (String, Optional)
- **qos_rate_max_down** (Number, Optional) The QOS maximum download rate. Defaults to `-1`.
- **qos_rate_max_up** (Number, Optional) The QOS maximum upload rate. Defaults to `-1`.



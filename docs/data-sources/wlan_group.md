---
page_title: "unifi_wlan_group Data Source - terraform-provider-unifi"
subcategory: ""
description: |-
  unifi_wlan_group data source can be used to retrieve the ID for a WLAN group by name.
---

# Data Source `unifi_wlan_group`

`unifi_wlan_group` data source can be used to retrieve the ID for a WLAN group by name.

## Example Usage

```terraform
data "unifi_wlan_group" "default" {
}
```

## Schema

### Optional

- **id** (String, Optional) The ID of this resource.
- **name** (String, Optional) The name of the WLAN group to look up. Defaults to `Default`.



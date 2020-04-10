---
subcategory: ""
layout: ""
page_title: "terraform-provider-unifi: unifi_wlan_group"
description: |-
  unifi_wlan_group data source can be used to retrieve the ID for a WLAN group by name.
---

# Resource: `unifi_wlan_group`

unifi_wlan_group data source can be used to retrieve the ID for a WLAN group by name.

## Example Usage

```terraform
data "unifi_wlan_group" "default" {
}
```

## Schema

### Optional

- **id** (String, Optional)
- **name** (String, Optional) The name of the WLAN group to look up. Defaults to `Default`.



---
page_title: "unifi_wlan_group Data Source - terraform-provider-unifi"
subcategory: ""
description: |-
  unifi_wlan_group data source can be used to retrieve the ID for a WLAN group by name.
  Please note that WLAN Groups are deprecated in v6 of the controller.
---

# Data Source `unifi_wlan_group`

`unifi_wlan_group` data source can be used to retrieve the ID for a WLAN group by name.

Please note that WLAN Groups are deprecated in v6 of the controller.



## Schema

### Optional

- **name** (String) The name of the WLAN group to look up. Defaults to `Default`.
- **site** (String) The name of the site the wlan group is associated with.

### Read-only

- **id** (String) The ID of this AP group.



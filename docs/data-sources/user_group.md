---
subcategory: ""
layout: ""
page_title: "unifi_user_group Data Source - terraform-provider-unifi"
description: |-
  unifi_user_group data source can be used to retrieve the ID for a user group by name.
---

# Data Source `unifi_user_group`

`unifi_user_group` data source can be used to retrieve the ID for a user group by name.



## Schema

### Optional

- **id** (String, Optional)
- **name** (String, Optional) The name of the user group to look up. Defaults to `Default`.

### Read-only

- **qos_rate_max_down** (Number, Read-only)
- **qos_rate_max_up** (Number, Read-only)



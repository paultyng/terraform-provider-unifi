---
page_title: "unifi_radius_profile Data Source - terraform-provider-unifi"
subcategory: ""
description: |-
  unifi_radius_profile data source can be used to retrieve the ID for a RADIUS profile by name.
---

# Data Source `unifi_radius_profile`

`unifi_radius_profile` data source can be used to retrieve the ID for a RADIUS profile by name.



## Schema

### Optional

- **name** (String) The name of the RADIUS profile to look up. Defaults to `Default`.
- **site** (String) The name of the site the radius profile is associated with.

### Read-only

- **id** (String) The ID of this AP group.



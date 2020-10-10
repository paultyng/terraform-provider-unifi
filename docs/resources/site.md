---
page_title: "unifi_site Resource - terraform-provider-unifi"
subcategory: ""
description: |-
  unifi_site manages Unifi sites
---

# Resource `unifi_site`

`unifi_site` manages Unifi sites

## Example Usage

```terraform
resource "unifi_site" "mysite" {
  description = "mysite"
}
```

## Schema

### Required

- **description** (String, Required) The description of the site.

### Read-only

- **id** (String, Read-only) The ID of the site.
- **name** (String, Read-only) The name of the site.



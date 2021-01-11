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

- **description** (String) The description of the site.

### Read-only

- **id** (String) The ID of the site.
- **name** (String) The name of the site.

## Import

Import is supported using the following syntax:

```shell
# import using the ID
terraform import unifi_site.mysite 5fe6261995fe130013456a36
```

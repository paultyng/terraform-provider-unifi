---
layout: ""
page_title: "Provider: Unifi"
description: |-
  The Unifi provider provides resources to interact with a Unifi controller API.
---

# Unifi Provider

The Unifi provider provides resources to interact with a Unifi controller API.

## Example Usage

```terraform
provider "unifi" {
  username = var.username # optionally use UNIFI_USERNAME env var
  password = var.password # optionally use UNIFI_PASSWORD env var
  api_url  = var.api_url  # optionally use UNIFI_API env var

  # if you are not configuring the default site, you can change the site
  # site = "foo" or optionally use UNIFI_SITE env var
}
```

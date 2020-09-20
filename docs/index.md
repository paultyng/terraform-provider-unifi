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

  # you may need to allow insecure TLS communications unless you have configured
  # certificates for your controller
  allow_insecure = var.insecure # optionally use UNIFI_INSECURE env var

  # if you are not configuring the default site, you can change the site
  # site = "foo" or optionally use UNIFI_SITE env var
}
```

## Schema

### Optional

- **allow_insecure** (Boolean, Optional) Skip verification of TLS certificates of API requests. You may need to set this to `true` if you are using your local API without setting up a signed certificate. Can be specified with the `UNIFI_INSECURE` environment variable.
- **api_url** (String, Optional) URL of the controller API. Can be specified with the `UNIFI_API` environment variable.
- **password** (String, Optional) Password for the user accessing the API. Can be specified with the `UNIFI_PASSWORD` environment variable.
- **site** (String, Optional) The site in the Unifi controller this provider will manage. Can be specified with the `UNIFI_SITE` environment variable. Default: `default`
- **username** (String, Optional) Local user name for the Unifi controller API. Can be specified with the `UNIFI_USERNAME` environment variable.

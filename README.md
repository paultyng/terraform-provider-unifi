> [!IMPORTANT]
> **This repository is archived and no longer maintained.**
>
> See [#461](https://github.com/paultyng/terraform-provider-unifi/issues/461) for context. There is no single successor; the following forks are currently published to the Terraform Registry:
>
> - [filipowm/unifi](https://registry.terraform.io/providers/filipowm/unifi)
> - [ubiquiti-community/unifi](https://registry.terraform.io/providers/ubiquiti-community/unifi)
> - [akerl/unifi](https://registry.terraform.io/providers/akerl/unifi)
>
> Pick whichever fork best fits your setup.
>
> ## Migrating
>
> Update your `required_providers` block to point at one of the forks above. For example:
>
> ```hcl
> terraform {
>   required_providers {
>     unifi = {
>       # was: source = "paultyng/unifi"
>       source  = "filipowm/unifi"  # or "ubiquiti-community/unifi", "akerl/unifi"
>       version = "~> 1.0"
>     }
>   }
> }
> ```
>
> Then run `terraform init -upgrade`.
>
> **State compatibility:** `ubiquiti-community` and `akerl` continue the v0.41.x lineage and should be drop-in. `filipowm` has rebased to v1.0.0 with schema changes and may require state edits. **Always run `terraform plan` before applying.**

---


![Acceptance Tests](https://github.com/paultyng/terraform-provider-unifi/workflows/Acceptance%20Tests/badge.svg?event=push)

# Unifi Terraform Provider (terraform-provider-unifi)

**Note** You can't (for obvious reasons) configure your network while connected to something that may disconnect (like the WiFi). Use a hard-wired connection to your controller to use this provider.

Functionality first needs to be added to the [go-unifi](https://github.com/paultyng/go-unifi) SDK.

## Documentation

You can browse documentation on the [Terraform provider registry](https://registry.terraform.io/providers/paultyng/unifi/latest/docs).

## Supported Unifi Controller Versions

As of version [v0.34](https://github.com/paultyng/terraform-provider-unifi/releases/tag/v0.34.0), this provider only supports version 6 of the Unifi controller software. If you need v5 support, you can pin an older version of the provider.

The docker, UDM, and UDM-Pro versions are slightly different (the API is proxied a little differently) but for the most part should all be supported. Individual patch versions of the controller are generally not tested for compatibility, just the latest stable versions.

## Using the Provider

### Terraform 1.0 and above

You can use the provider via the [Terraform provider registry](https://registry.terraform.io/providers/paultyng/unifi).

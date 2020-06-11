# Unifi Terraform Provider (terraform-provider-unifi)

**Note** You can't (for obvious reasons) configure your network while connected to something that may disconnect (like the WiFi). Use a hard-wired connection to your controller to use this provider.

Functionality first needs to be added to the [go-unifi](https://github.com/paultyng/go-unifi) SDK.

## Documentation

You can browse documentation on the [Terraform provider registry](https://registry.terraform.io/providers/paultyng/unifi).

## Using the Provider

### Terraform 0.13 and above

You can use the provider via the [Terraform provider registry](https://registry.terraform.io/providers/paultyng/unifi).

### Terraform 0.12 or manual installation

You can download a pre-built binary from the [releases](https://github.com/paultyng/terraform-provider-unifi/releases) page, these are built using [goreleaser](https://goreleaser.com/) (the [configuration](.goreleaser.yml) is in the repo). You can verify the signature and my [key ownership via Keybase](https://keybase.io/paultyng).

If you want to build from source, you can simply use `go build` in the root of the repository.

To use an external provider binary with Terraform ([until the provider registry is GA](https://www.hashicorp.com/blog/announcing-providers-in-the-new-terraform-registry/)), you need to place the binary in a [plugin location for Terraform](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) to find it.

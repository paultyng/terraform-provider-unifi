<!-- ![Acceptance Tests](https://github.com/sayedh/terraform-provider-unifi/workflows/Acceptance%20Tests/badge.svg?event=push) -->

# Unifi Terraform Provider (terraform-provider-unifi)

This project is a fork of [paultyng/terraform-provider-unifi](https://github.com/paultyng/terraform-provider-unifi), updated to support the latest UniFi Network versions.

## Documentation

Browse the official provider documentation on the [Terraform provider registry](https://registry.terraform.io/providers/sayedh/unifi/latest/docs).

## Supported UniFi Controller Versions

As of **version v1.0.1**, this provider supports **UniFi Controller v8.4.59**. Earlier versions, up to **v7.4.162**, were supported in the original [paultyng](https://github.com/paultyng/terraform-provider-unifi) release.


## Development Status

This provider is currently under active development, with efforts to update compatibility with the latest UniFi Network Controller version **v8.4.59**. While significant progress has been made, **only a few features are fully tested and confirmed to work**:

- **Resource: Network**
- **Resource: Port Profile**

Other resources and data sources **may not yet be fully functional** on the latest UniFi Network version, and testing is ongoing. Users should be aware that these features might work inconsistently or could break entirely when working with newer versions of the UniFi controller.

## Fork History

- This Terraform provider was originally created by [Paul Tyng](https://github.com/paultyng) to support UniFi Network versions up to 7.4.162.
- In August 2024, this fork was initiated to extend support to newer versions of the UniFi Network software, specifically **v8.4.59**.
- Ongoing efforts are being made to thoroughly test and refine all functionality.

## Using the Provider

### Terraform 1.0 and Above

You can use the provider via the [Terraform provider registry](https://registry.terraform.io/providers/sayedh/unifi/latest).

**Note**: When using this provider, ensure you're connected via a hard-wired connection to the UniFi Controller rather than WiFi, as configuring your network over a connection that could disconnect (like WiFi) is risky and may result in issues.

### Terraform Configuration Example

```hcl
provider "unifi" {
  controller = "https://<unifi-controller-url>"
  username   = "admin"
  password   = "password"
}

resource "unifi_network" "example" {
  name = "Example Network"
  # Add relevant configuration options here
}
```

## Versioning

The provider has been versioned with the goal of maintaining backward compatibility where possible. However, many changes involve breaking adjustments. Updates will continue to follow semantic versioning.

## Note on UniFi Go SDK

This provider relies on the [go-unifi SDK](https://github.com/sayedh/go-unifi), also forked from [Paul Tyng's original SDK](https://github.com/paultyng/go-unifi), to interact with the UniFi Controller. Code generation and updates to the SDK are performed to support the latest UniFi Controller features.

## Contributing

Contributions are highly appreciated to ensure better compatibility with the latest UniFi Network versions. Please submit issues and pull requests to the [GitHub repository](https://github.com/sayedh/terraform-provider-unifi).

## Disclaimer

This project is in an **experimental state**, and users should proceed with caution when using it in production environments. Expect rapid changes and potentially breaking updates as development progresses.


# terraform-provider-unifi

This is very much WIP, just adding functionality as I need it for my local setup.

## unifi_network

Example:

```terraform
resource "unifi_network" "test" {
	name    = "foo"
	purpose = "corporate"

	subnet       = "10.0.202.1/24"
	vlan_id      = 202
	dhcp_start   = "10.0.202.6"
	dhcp_stop    = "10.0.202.254"
	dhcp_enabled = true
}
```

# terraform-provider-unifi

This is very much WIP, just adding functionality as I need it for my local setup.

**Note** You can't (for obvious reasons) configure your network while connected to something that may disconnect (like the WiFi). Use a hard-wired connection to your controller to use this provider.

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

## unifi_wlan

Example:

```terraform
data "unifi_wlan_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_wlan" "test" {
	name          = "foo"
	vlan_id       = 202
	passphrase    = "12345678"
	wlan_group_id = data.unifi_wlan_group.default.id
	user_group_id = data.unifi_user_group.default.id
	security      = "wpapsk"
}
```

## TODO

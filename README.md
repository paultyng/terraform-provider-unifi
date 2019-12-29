# terraform-provider-unifi

This is very much WIP, just adding functionality as I need it for my local setup.

**Note** You can't (for obvious reasons) configure your network while connected to something that may disconnect (like the WiFi). Use a hard-wired connection to your controller to use this provider.

## Provider configuration

```terraform
provider "unifi" {
	username = "user" // optionally use UNIFI_USERNAME env var
	password = "pass" // optionally use UNIFI_PASSWORD env var
	api_url  = "https://localhost:8443/api/" // optionally use UNIFI_API env var

	// if you are not configuring the default site, you can change the site
	site = "foo" // optionally use UNIFI_SITE env var
}
```

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

## unifi_user_group

Example:

```terraform
resource "unifi_user_group" "test" {
	name = "foo"

	qos_rate_max_down = 2000 # 2mbps
	qos_rate_max_up   = 10   # 10kbps
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

* [ ] automatically fixup subnet cidrs from .0 to .1?
* [ ] update support for wlan
* [ ] update support for network

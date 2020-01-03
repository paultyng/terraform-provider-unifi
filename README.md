# terraform-provider-unifi

This is very much WIP, just adding functionality as I need it for my local setup.

**Note** You can't (for obvious reasons) configure your network while connected to something that may disconnect (like the WiFi). Use a hard-wired connection to your controller to use this provider.

## Examples

### Clients (Users) with optional fixed IPs from CSV

```terraform
locals {
  userscsv = csvdecode(file("${path.module}/users.csv"))
  users    = { for user in local.userscsv : user.mac => user }
}

resource "unifi_user" "user" {
  for_each = local.users

  mac  = each.key
  name = each.value.name
  # append an optional additional note
  note = trimspace("${each.value.note}\n\nmanaged by TF")

  fixed_ip   = each.value.fixed_ip
  # this assumes there is a unifi_network for_each with names
  network_id = each.value.network != "" ? unifi_network.vlan[each.value.network].id : ""

  allow_existing         = true
  skip_forget_on_destroy = true
}
```

```csv
mac,name,note,network,fixed_ip
00:00:00:00:00:00,My Device,custom note,,,
00:00:00:00:00:00,My Device,custom note,network name to lookup,10.0.3.4
```

### WIFI (WLAN) and Network for VLAN

This example sets up a WIFI SSID and VLAN with bandwidth throttling:

```terraform
variable "vlan_id" {
	default = 10
}

data "unifi_wlan_group" "default" {
}

resource "unifi_user_group" "wifi" {
	name = "wifi"

	qos_rate_max_down = 10000000
	qos_rate_max_up   = 2000000
}

resource "unifi_wlan" "wifi" {
	name          = "myssid"
	vlan_id       = var.vlan_id
	passphrase    = "12345678"
	wlan_group_id = data.unifi_wlan_group.default.id
	user_group_id = unifi_user_group.wifi.id
	security      = "wpapsk"
}

resource "unifi_network" "vlan" {
	name = "wifi-vlan"
	purpose = "corporate"

	subnet       = "10.0.0.1/24"
	vlan_id      = var.vlan_id
	dhcp_start   = "10.0.0.6"
	dhcp_stop    = "10.0.0.254"
	dhcp_enabled = true
}
```

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

## Resources

### unifi_network

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

### unifi_user

User's (called "Clients" in the UI), are unique as they are "created" when observed, so the resource defaults to allowing itself to just take over management of a MAC address, but this can be turned off.

Example:

```terraform
resource "unifi_user" "test" {
	mac  = "00:00:5e:00:53:10"
	name = "some client"
	note = "my note"

	fixed_ip   = "10.1.10.50"
	network_id = unifi_network.my_vlan.id
}
```

### unifi_user_group

Example:

```terraform
resource "unifi_user_group" "test" {
	name = "foo"

	qos_rate_max_down = 2000 # 2mbps
	qos_rate_max_up   = 10   # 10kbps
}
```

### unifi_wlan

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

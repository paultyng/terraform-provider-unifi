variable "vlan_id" {
  default = 10
}

data "unifi_ap_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_network" "vlan" {
  name    = "wifi-vlan"
  purpose = "corporate"

  subnet       = "10.0.0.1/24"
  vlan_id      = var.vlan_id
  dhcp_start   = "10.0.0.6"
  dhcp_stop    = "10.0.0.254"
  dhcp_enabled = true
}

resource "unifi_wlan" "wifi" {
  name       = "myssid"
  passphrase = "12345678"
  security   = "wpapsk"

  # enable WPA2/WPA3 support
  wpa3_support    = true
  wpa3_transition = true

  network_id    = unifi_network.vlan.id
  ap_group_ids  = [data.unifi_ap_group.default.id]
  user_group_id = data.unifi_user_group.default.id
}

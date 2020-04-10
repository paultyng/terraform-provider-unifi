data "unifi_wlan_group" "default" {
}

data "unifi_user_group" "default" {
}

resource "unifi_wlan" "wifi" {
  name          = "myssid"
  vlan_id       = 10
  passphrase    = "12345678"
  wlan_group_id = data.unifi_wlan_group.default.id
  user_group_id = data.unifi_user_group.default.id
  security      = "wpapsk"
}

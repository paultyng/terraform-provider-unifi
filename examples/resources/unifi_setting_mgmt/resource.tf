resource "unifi_site" "example" {
  description = "example"
}

resource "unifi_setting_mgmt" "example" {
  site         = unifi_site.example.name
  auto_upgrade = true
}

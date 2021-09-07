#retrieve network data by unifi network name
data "unifi_network" "lan_network" {
  name = "LAN"
}

#retrieve network data from user record
data "unifi_user" "my_device" {
  mac = "01:23:45:67:89:ab"
}
data "unifi_network" "my_network" {
  id = data.unifi_user.my_device.network_id
}

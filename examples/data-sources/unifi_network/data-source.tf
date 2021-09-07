#retrieve network data by unifi network name
data "unifi_network" "byName" {
  name="LAN"
}

#retrieve network data from user record
data "unifi_user" "byMac" {
  mac = "01:23:45:67:89:ab"
}
data "unifi_network" "byID" {
  id = data.unifi_user.byMac.network_id
}

resource "unifi_user" "test" {
  mac  = "01:23:45:67:89:AB"
  name = "some client"
  note = "my note"

  fixed_ip   = "10.0.0.50"
  network_id = unifi_network.my_vlan.id
}
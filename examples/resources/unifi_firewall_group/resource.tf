variable "laptop_ips" {
  type = list(string)
}

resource "unifi_firewall_group" "can_print" {
  name = "can-print"
  type = "address-group"

  members = var.laptop_ips
}
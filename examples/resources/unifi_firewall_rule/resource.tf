variable "ip_address" {
  type = string
}

resource "unifi_firewall_rule" "drop_all" {
  name    = "drop all"
  action  = "drop"
  ruleset = "LAN_IN"

  rule_index = 2011

  protocol = "all"

  dst_address = var.ip_address
}
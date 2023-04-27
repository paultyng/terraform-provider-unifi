variable "tgt_ip_address" {
  type = string
}

variable "src_ip_address" {
  type = string
}

resource "unifi_firewall_rule" "allow_from" {
  name    = "drop all"
  action  = "drop"
  ruleset = "LAN_IN"

  protocol = "all"

  src_address = var.src_ip_address
  dst_address = var.tgt_ip_address
}

resource "unifi_firewall_rule" "drop_all" {
  name    = "drop all"
  action  = "drop"
  ruleset = "LAN_IN"

  protocol = "all"

  dst_address = var.tgt_ip_address
}

resource "unifi_firewall_ruleset" "lan_in" {
  ruleset = "LAN_IN"
  before_predefined = [
    unifi_firewall_rule.allow_from.id,
    unifi_firewall_rule.drop_all.id,
  ]
}

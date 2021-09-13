resource "unifi_firewall_rule" "rule" {
  # list of sites
  for_each = toset(["default", "vq98kwez", "bfa2l6i7"])
  # use the key of the list as the site value
  site = each.key

  name    = "drop all"
  action  = "drop"
  ruleset = "LAN_IN"

  rule_index = 2011

  protocol = "all"

  dst_address = var.ip_address
}

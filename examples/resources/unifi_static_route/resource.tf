resource "unifi_static_route" "nexthop" {
  type     = "nexthop-route"
  network  = "172.17.0.0/16"
  name     = "basic nexthop"
  distance = 1
  next_hop = "172.16.0.1"
}

resource "unifi_static_route" "blackhole" {
  type     = "blackhole"
  network  = var.blackhole_cidr
  name     = "blackhole traffice to cidr"
  distance = 1
}

resource "unifi_static_route" "interface" {
  type      = "interface-route"
  network   = var.wan2_cidr
  name      = "send traffic over wan2"
  distance  = 1
  interface = "WAN2"
}

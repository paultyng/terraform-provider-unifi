data "unifi_port_profile" "disabled" {
  # look up the built-in disabled port profile
  name = "Disabled"
}

resource "unifi_port_profile" "poe" {
  name    = "poe"
  forward = "customize"

  native_networkconf_id = var.native_network_id

  poe_mode = "auto"
}

resource "unifi_device" "us_24_poe" {
  # optionally specify MAC address to skip manually importing
  # manual import is the safest way to add a device
  mac = "01:23:45:67:89:AB"

  name = "Switch with POE"

  port_override {
    number          = 1
    name            = "port w/ poe"
    port_profile_id = unifi_port_profile.poe.id
  }

  port_override {
    number          = 2
    name            = "disabled"
    port_profile_id = data.unifi_port_profile.disabled.id
  }

  # port aggregation for ports 11 and 12
  port_override {
    number              = 11
    op_mode             = "aggregate"
    aggregate_num_ports = 2
  }
}

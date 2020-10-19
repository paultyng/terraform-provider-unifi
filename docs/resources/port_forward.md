---
page_title: "unifi_port_forward Resource - terraform-provider-unifi"
subcategory: ""
description: |-
  unifi_port_forward manages a port forwarding rule on the gateway.
---

# Resource `unifi_port_forward`

`unifi_port_forward` manages a port forwarding rule on the gateway.



## Schema

### Optional

- **dst_port** (String, Optional) The destination port for the forwarding.
- **enabled** (Boolean, Optional, Deprecated) Specifies whether the port forwarding rule is enabled or not. This will attribute will be removed in a future release. Instead of disabling a port forwarding rule you can remove it from your configuration.
- **fwd_ip** (String, Optional) The IPv4 address to forward traffic to.
- **fwd_port** (String, Optional) The port to forward traffic to.
- **log** (Boolean, Optional) Specifies whether to log forwarded traffic or not. Defaults to `false`.
- **name** (String, Optional) The name of the port forwarding rule.
- **port_forward_interface** (String, Optional) The port forwarding interface. Can be `wan`, `wan2`, or `both`.
- **protocol** (String, Optional) The protocol for the port forwarding rule. Can be `tcp`, `udp`, or `tcp_udp`. Defaults to `tcp_udp`.
- **src_ip** (String, Optional) The source IPv4 address (or CIDR) of the port forwarding rule. For all traffic, specify `any`. Defaults to `any`.

### Read-only

- **id** (String, Read-only) The ID of the port forwarding rule.



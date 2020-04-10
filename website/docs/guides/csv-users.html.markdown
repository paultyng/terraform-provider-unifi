---
subcategory: ""
page_title: "Manage Users/Clients in a CSV - Unifi Provider"
description: |-
    An example of using a CSV to manage all of your users of your network.
---

# Manage Users in a CSV

Given a CSV file with the following content:

```csv
mac,name,note,network,fixed_ip
00:00:00:00:00:00,My Device,custom note,,,
00:00:00:00:00:00,My Device,custom note,network name to lookup,10.0.3.4
```

You could create/manage a `unifi_user` for every row/MAC address in the CSV with the following config:

```terraform
locals {
  userscsv = csvdecode(file("${path.module}/users.csv"))
  users    = { for user in local.userscsv : user.mac => user }
}

resource "unifi_user" "user" {
  for_each = local.users

  mac  = each.key
  name = each.value.name
  # append an optional additional note
  note = trimspace("${each.value.note}\n\nmanaged by TF")

  fixed_ip = each.value.fixed_ip
  # this assumes there is a unifi_network for_each with names
  network_id = each.value.network != "" ? unifi_network.vlan[each.value.network].id : ""

  allow_existing         = true
  skip_forget_on_destroy = true
}
```

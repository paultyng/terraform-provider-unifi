---
subcategory: ""
page_title: "Manage Users/Clients in a CSV - Unifi Provider"
description: |-
    An example of using a CSV to manage all of your users of your network.
---

# Manage Users in a CSV

Given a CSV file with the following content:

```csv
mac,name,note
01:23:45:67:89:AB,My Device,custom note
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

  allow_existing         = true
  skip_forget_on_destroy = true
}
```

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
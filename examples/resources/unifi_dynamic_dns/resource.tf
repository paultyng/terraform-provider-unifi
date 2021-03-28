resource "unifi_dynamic_dns" "test" {
  service = "dyndns"

  host_name = "my-network.example.com"

  server   = "domains.google.com"
  login    = var.dns_login
  password = var.dns_password
}

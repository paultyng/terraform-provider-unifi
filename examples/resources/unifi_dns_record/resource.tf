resource "unifi_dns_record" "test" {
  name        = "my-network"
  enabled     = true
  port        = 0
  priority    = 10
  record_type = "A"
  ttl         = 300
  value       = "my-network.example.com"
}

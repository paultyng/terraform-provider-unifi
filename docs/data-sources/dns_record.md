---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "unifi_dns_record Data Source - terraform-provider-unifi"
subcategory: ""
description: |-
  unifi_dns_record data source can be used to retrieve the ID for an DNS record by name.
---

# unifi_dns_record (Data Source)

`unifi_dns_record` data source can be used to retrieve the ID for an DNS record by name.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `name` (String) The name of the DNS record to look up, leave blank to look up the default DNS record.
- `port` (Number) The port of the DNS record.
- `priority` (Number) The priority of the DNS record.
- `record_type` (String) The type of the DNS record.
- `site` (String) The name of the site the DNS record is associated with.
- `ttl` (Number) The TTL of the DNS record.
- `value` (String) The value of the DNS record.
- `weight` (Number) The weight of the DNS record.

### Read-Only

- `id` (String) The ID of this DNS record.
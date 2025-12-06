---
page_title: "ipnetwork Provider"
description: |-
  IP Network functions
---

# ipnetwork Provider

Provides functions to manipulate IP Network data (address, cidr, ...)

## Example Usage

```terraform
terraform {
  required_providers {
    ipnetwork = {
      source = "jeremmfr/ipnetwork"
    }
  }
  required_version = ">= 1.8.0"
}

output "ptr" {
  value = provider::ipnetwork::ptr("192.0.2.128")
}
# result: "128.2.0.192.in-addr.arpa."

output "contain" {
  value = provider::ipnetwork::contain("192.0.2.0/23", "192.0.2.64/22")
}
# result: false

output "summarize" {
  value = provider::ipnetwork::summarize(toset([
    "192.0.2.1",
    "192.0.2.0/24",
    "192.0.2.128/25",
  ]))
}
# result: ["192.0.2.0/24"]

output "ipv6_eui64" {
  value = provider::ipnetwork::generate6_eui64("fe80::", "00:00:5e:53:53:00")
}
# result: "fe80::200:5eff:fe53:5300"

output "is_private" {
  value = provider::ipnetwork::is_private("192.168.0.0/10")
}
# result: false
```

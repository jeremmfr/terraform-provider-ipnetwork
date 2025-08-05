---
page_title: "contain function - ipnetwork"
description: |-
  contain function
---

# function: contain

Reports whether a prefix (`container`) contains

- an address if `address` is not in CIDR format
- all addresses of address block if `address` is in CIDR format

Reports `false` if container and address have different IP version

## Example Usage

```terraform
output "ip_v4" {
  value = provider::ipnetwork::contain("192.0.2.0/23", "192.0.3.64")
}
# result: true

output "ip_v4_cidr" {
  value = provider::ipnetwork::contain("192.0.2.0/23", "192.0.2.64/22")
}
# result: false

output "ip_v4_ip_container" {
  value = provider::ipnetwork::contain("192.0.3.128/23", "192.0.2.64/23")
}
# result: true

output "ip_v6" {
  value = provider::ipnetwork::contain("2001:db8::ffff/64", "2001:db8::a:ffff")
}
# result: true

output "ip_v6_cidr" {
  value = provider::ipnetwork::contain("2001:db8::ffff/64", "2001:db8::/65")
}
# result: true
```

## Signature

```text
contain(container string, address string) boolean
```

## Arguments

1. `container` (String) Container address to parse
1. `address` (String) Included address(es) to parse

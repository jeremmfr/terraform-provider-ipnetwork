---
page_title: "equal_prefix function - ipnetwork"
description: |-
  equal_prefix function
---

# function: equal_prefix

Compare two CIDR addresses if they are in the same prefix.

## Example Usage

```terraform
output "ip_v4_match" {
  value = provider::ipnetwork::equal_prefix("192.0.3.128/23", "192.0.2.64/23")
}
# result: true

output "ip_v4_not_match" {
  value = provider::ipnetwork::equal_prefix("192.0.2.64/24", "192.0.2.64/23")
}
# result: false

output "ip_v6_match" {
  value = provider::ipnetwork::equal_prefix("2001:db8::ffff/64", "2001:db8::a:ffff/64")
}
# result: true
```

## Signature

```text
equal_prefix(address_x string, address_y string) boolean
```

## Arguments

1. `address_x` (String) First address to parse
1. `address_y` (String) Second address to parse

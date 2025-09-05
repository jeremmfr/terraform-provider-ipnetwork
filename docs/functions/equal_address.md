---
page_title: "equal_address function - ipnetwork"
description: |-
  equal_address function
---

# function: equal_address

Compare two address if there are equal regardless of format: CIDR or not, IPv6 expanded or not.

## Example Usage

```terraform
output "ip_v4" {
  value = provider::ipnetwork::equal_address("192.0.2.128/24", "192.0.2.128")
}
# result: true

output "ip_v6_match" {
  value = provider::ipnetwork::equal_address("2001:0db8:0:0:0:0:0:f", "2001:db8::f/64")
}
# result: true
```

## Signature

```text
equal_address(address_x string, address_y string) boolean
```

## Arguments

1. `address_x` (String) First address to parse
2. `address_y` (String) Second address to parse

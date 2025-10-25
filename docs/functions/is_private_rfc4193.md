---
page_title: "is_private_rfc4193 function - ipnetwork"
description: |-
  is_private_rfc4193 function
---

# function: is_private_rfc4193

Reports whether an address or prefix is in RFC4193 Unique Local Address (ULA) space.

For single addresses, checks if the address is in the ULA range.

For prefixes (CIDR notation), checks if the **entire prefix** is contained within the ULA range.  
A prefix is considered in RFC4193 space only if it is entirely contained within the range.
Prefixes that partially overlap with the ULA range (e.g., larger prefixes containing
both ULA and non-ULA addresses) return `false`.

Returns `true` for `fc00::/7`.

## Example Usage

```terraform
# ULA addresses
output "ula_address" {
  value = provider::ipnetwork::is_private_rfc4193("fd00::1")
}
# result: true

output "ula_prefix" {
  value = provider::ipnetwork::is_private_rfc4193("fc00::/7")
}
# result: true

# Non-ULA
output "public_ipv6" {
  value = provider::ipnetwork::is_private_rfc4193("2001:4860:4860::8888")
}
# result: false

output "ipv4_address" {
  value = provider::ipnetwork::is_private_rfc4193("10.0.0.1")
}
# result: false (RFC4193 is IPv6-specific)
```

## Signature

```text
is_private_rfc4193(input string) boolean
```

## Arguments

1. `input` (String) Address or prefix to parse

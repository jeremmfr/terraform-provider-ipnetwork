---
page_title: "is_private_rfc6598 function - ipnetwork"
description: |-
  is_private_rfc6598 function
---

# function: is_private_rfc6598

Reports whether an address or prefix is in RFC6598 Shared Address Space.

For single addresses, checks if the address is in the Shared Address Space range.

For prefixes (CIDR notation), checks if the **entire prefix** is contained within
the Shared Address Space range.  
A prefix is considered in RFC6598 space only if it is entirely contained within the range.
Prefixes that partially overlap with the Shared Address Space (e.g., larger prefixes containing
both Shared Address Space and non-Shared Address Space addresses) return `false`.

Returns `true` for `100.64.0.0/10`.

-> **Note:**
  IPv6 address/prefix in `::ffff:0:0/96` is unmap to IPv4 version
  (unmap the prefix mask by subtracting 96)

## Example Usage

```terraform
# RFC6598 Shared Address Space
output "cgn_address" {
  value = provider::ipnetwork::is_private_rfc6598("100.64.0.1")
}
# result: true

output "cgn_prefix" {
  value = provider::ipnetwork::is_private_rfc6598("100.64.0.0/10")
}
# result: true

# Non-RFC6598
output "rfc1918_private" {
  value = provider::ipnetwork::is_private_rfc6598("10.0.0.1")
}
# result: false (this is RFC1918, not RFC6598)

output "ipv6_address" {
  value = provider::ipnetwork::is_private_rfc6598("fc00::1")
}
# result: false (RFC6598 is IPv4-specific)
```

## Signature

```text
is_private_rfc6598(input string) boolean
```

## Arguments

1. `input` (String) Address or prefix to parse

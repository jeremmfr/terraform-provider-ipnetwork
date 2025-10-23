---
page_title: "is_private_rfc1918 function - ipnetwork"
description: |-
  is_private_rfc1918 function
---

# function: is_private_rfc1918

Reports whether an address or prefix is in RFC1918 private address space.

For single addresses, checks if the address is in RFC1918 ranges.

For prefixes (CIDR notation), checks if the **entire prefix** is contained within RFC1918 ranges.  
A prefix is considered in RFC1918 space only if it is entirely contained within one of the ranges.
Prefixes that partially overlap with RFC1918 ranges (e.g., larger prefixes containing
both RFC1918 and non-RFC1918 addresses) return `false`.

Returns `true` for:

- `10.0.0.0/8`
- `172.16.0.0/12`
- `192.168.0.0/16`

-> **Note:**
  IPv6 address/prefix in `::ffff:0:0/96` is unmap to IPv4 version
  (unmap the prefix mask by subtracting 96)

## Example Usage

```terraform
# RFC1918 addresses
output "private_10" {
  value = provider::ipnetwork::is_private_rfc1918("10.0.0.1")
}
# result: true

output "private_prefix" {
  value = provider::ipnetwork::is_private_rfc1918("192.168.0.0/16")
}
# result: true

# Non-RFC1918
output "public_ipv4" {
  value = provider::ipnetwork::is_private_rfc1918("8.8.8.8")
}
# result: false

output "cgn_rfc6598" {
  value = provider::ipnetwork::is_private_rfc1918("100.64.0.1")
}
# result: false (this is RFC6598, not RFC1918)

output "ipv6_address" {
  value = provider::ipnetwork::is_private_rfc1918("fc00::1")
}
# result: false (RFC1918 is IPv4-specific)
```

## Signature

```text
is_private_rfc1918(input string) boolean
```

## Arguments

1. `input` (String) Address or prefix to parse

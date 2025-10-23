---
page_title: "is_private function - ipnetwork"
description: |-
  is_private function
---

# function: is_private

Reports whether an address or prefix is private (internally routable).

For single addresses, checks if the address is private.

For prefixes (CIDR notation), checks if the **entire prefix** contains only private addresses.  
A prefix is considered private only if it is entirely contained within a defined private range.
Prefixes that partially overlap with private ranges (e.g., larger prefixes containing
both private and public addresses) return `false`.

Returns `true` for:

- Private-Use addresses (`10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`, `fc00::/7`)
- Shared Address Space (`100.64.0.0/10`)
- Benchmarking (`198.18.0.0/15`, `2001:2::/48`)
- Discard prefix (`100::/64`)
- IPv4-IPv6 Translation (`64:ff9b:1::/48`)
- Segment Routing (SRv6) SIDs (`5f00::/16`)

-> **Note:**
  IPv6 address/prefix in `::ffff:0:0/96` is unmap to IPv4 version
  (unmap the prefix mask by subtracting 96)

~> **Note:**
  This function specifically identifies addresses designated for private/internal use.
  It does **not** include loopback addresses, link-local addresses, multicast addresses,
  documentation ranges, or other special-use addresses.
  For comprehensive non-public address detection, use the `is_public` function.

## Example Usage

```terraform
# Single addresses
output "private_ipv4" {
  value = provider::ipnetwork::is_private("10.0.0.1")
}
# result: true

output "public_ipv4" {
  value = provider::ipnetwork::is_private("8.8.8.8")
}
# result: false

# Private prefixes
output "private_prefix_ipv4" {
  value = provider::ipnetwork::is_private("192.168.0.0/16")
}
# result: true

output "private_prefix_ipv6" {
  value = provider::ipnetwork::is_private("fc00::/7")
}
# result: true

# Non-private addresses/prefixes
output "public_prefix" {
  value = provider::ipnetwork::is_private("1.1.1.0/24")
}
# result: false

output "link_local_ipv6" {
  value = provider::ipnetwork::is_private("fe80::1")
}
# result: false (link-local is not considered private)

output "prefix_contains_private" {
  value = provider::ipnetwork::is_private("192.0.0.0/8")
}
# result: false (partially overlaps with 192.168.0.0/16 but also contains public addresses)
```

## Signature

```text
is_private(input string) boolean
```

## Arguments

1. `input` (String) Address or prefix to parse

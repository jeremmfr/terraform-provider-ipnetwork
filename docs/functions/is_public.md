---
page_title: "is_public function - ipnetwork"
description: |-
  is_public function
---

# function: is_public

Reports whether an address or prefix is public (globally routable).

For single addresses, checks if the address is public.

For prefixes (CIDR notation), checks if the **entire prefix** contains only public addresses.  
A prefix is considered non-public if it overlaps with any private, reserved, documentation,
or special-use ranges.

Returns `false` for:

- Private addresses (`10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16`, `fc00::/7`)
- Shared Address Space (`100.64.0.0/10`)
- Loopback addresses (`127.0.0.0/8`, `::1/128`)
- Link-local addresses (`169.254.0.0/16`, `fe80::/10`)
- Multicast addresses (`224.0.0.0/4`, `ff00::/8`)
- "This network" (`0.0.0.0/8`) & Unspecified addresses (`::/128`)
- IETF Protocol Assignments (`192.0.0.0/24`)
- Benchmarking (`198.18.0.0/15`, `2001:2::/48`)
- Reserved addresses (`240.0.0.0/4`), including broadcast address
- Documentation ranges (`192.0.2.0/24`, `198.51.100.0/24`, `203.0.113.0/24`, `2001:db8::/32`, `3fff::/20`)
- Discard prefix (`100::/64`)
- Segment Routing (SRv6) SIDs (`5f00::/16`)
- IPv4/IPv6 Translation (`64:ff9b:1::/48`)

-> **Note:**
  IPv6 address/prefix in `::ffff:0:0/96` is unmap to IPv4 version
  (unmap the prefix mask by subtracting 96)

## Example Usage

```terraform
# Single addresses
output "public_ipv4" {
  value = provider::ipnetwork::is_public("8.8.8.8")
}
# result: true

output "private_rfc1918" {
  value = provider::ipnetwork::is_public("192.168.1.1")
}
# result: false

# Public prefixes
output "public_prefix_ipv4" {
  value = provider::ipnetwork::is_public("1.1.1.0/24")
}
# result: true

output "public_prefix_ipv6" {
  value = provider::ipnetwork::is_public("2001:4860::/48")
}
# result: true

# Non-public prefixes
output "private_prefix" {
  value = provider::ipnetwork::is_public("192.168.0.0/16")
}
# result: false

output "prefix_contains_private" {
  value = provider::ipnetwork::is_public("192.0.0.0/20")
}
# result: false (contains 192.0.0.0/24 and 192.0.2.0/24)
```

## Signature

```text
is_public(input string) boolean
```

## Arguments

1. `input` (String) Address or prefix to parse

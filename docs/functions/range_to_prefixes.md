---
page_title: "range_to_prefixes function - ipnetwork"
description: |-
  range_to_prefixes function
---

# function: range_to_prefixes

Convert a range of IP addresses defined by a start and end address into the
minimal list of CIDR prefixes that exactly cover the range.

The function:

- Accepts IPv4 or IPv6 addresses (both must be the same version)
- Returns the smallest possible list of CIDR prefixes covering the range
- Start address must be less than or equal to end address

## Example Usage

```terraform
# Aligned range produces a single prefix
output "aligned" {
  value = provider::ipnetwork::range_to_prefixes("10.0.0.0", "10.0.0.255")
}
# result: ["10.0.0.0/24"]

# Single address
output "single" {
  value = provider::ipnetwork::range_to_prefixes("10.0.0.1", "10.0.0.1")
}
# result: ["10.0.0.1/32"]

# Unaligned range splits into multiple prefixes
output "unaligned" {
  value = provider::ipnetwork::range_to_prefixes("10.0.0.5", "10.0.0.20")
}
# result: ["10.0.0.5/32", "10.0.0.6/31", "10.0.0.8/29", "10.0.0.16/30", "10.0.0.20/32"]

# IPv6 range
output "ipv6" {
  value = provider::ipnetwork::range_to_prefixes("2001:db8::", "2001:db8::ff")
}
# result: ["2001:db8::/120"]
```

## Signature

```text
range_to_prefixes(start string, end string) list of string
```

## Arguments

1. `start` (String) Start address of the range
2. `end` (String) End address of the range

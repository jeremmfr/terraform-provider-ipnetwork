---
page_title: "summarize function - ipnetwork"
description: |-
  summarize function
---

# function: summarize

Summarize a set of IP addresses and prefixes into the smallest possible list of
prefixes that cover the same addresses.

The function:

- Converts standalone IP addresses to host prefixes (`/32` for IPv4, `/128` for IPv6)
- Merges adjacent prefixes of the same size into larger blocks
- Removes overlapping prefixes (keeps only the largest covering prefix)
- Processes IPv4 and IPv6 addresses separately
- Returns results sorted by address

## Example Usage

```terraform
# Merge two adjacent /24 networks into a /23
output "adjacent_merge" {
  value = provider::ipnetwork::summarize(toset([
    "192.0.2.0/24",
    "192.0.3.0/24",
  ]))
}
# result: ["192.0.2.0/23"]

# Convert addresses to /32 prefixes
output "addresses" {
  value = provider::ipnetwork::summarize(toset([
    "192.0.2.2",
    "192.0.2.3",
  ]))
}
# result: ["192.0.2.2/31"]

# Remove overlapping prefixes
output "overlapping" {
  value = provider::ipnetwork::summarize(toset([
    "192.0.2.0/24",
    "192.0.2.0/25",
    "192.0.2.128/25",
  ]))
}
# result: ["192.0.2.0/24"]

# Process IPv4 and IPv6 separately
output "mixed_families" {
  value = provider::ipnetwork::summarize(toset([
    "192.0.2.0/24",
    "192.0.3.0/24",
    "2001:db8::/64",
    "2001:db8:0:1::/64",
  ]))
}
# result: ["192.0.2.0/23", "2001:db8::/63"]

# Mix addresses and prefixes
output "mixed" {
  value = provider::ipnetwork::summarize(toset([
    "192.0.2.1",
    "192.0.2.0/25",
    "192.0.3.0/24",
  ]))
}
# result: ["192.0.2.0/25", "192.0.3.0/24"]

# IPv6 example
output "ipv6_adjacent" {
  value = provider::ipnetwork::summarize(toset([
    "2001:db8::/64",
    "2001:db8:0:1::/64",
    "2001:db8:0:2::/64",
    "2001:db8:0:3::/64",
  ]))
}
# result: ["2001:db8::/62"]
```

## Signature

```text
summarize(inputs set of string) list of string
```

## Arguments

1. `inputs` (Set of String) Set of IP addresses and prefixes to summarize

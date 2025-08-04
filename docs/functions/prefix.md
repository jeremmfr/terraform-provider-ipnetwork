---
page_title: "prefix function - ipnetwork"
description: |-
  prefix function
---

# function: prefix

Canonicalize CIDR address to obtain the 'network' address (prefix) of the address block.

## Example Usage

```terraform
output "ip_v4" {
  value = provider::ipnetwork::prefix("192.0.3.128/23")
}
# result: "192.0.2.0/23"

output "ip_v6" {
  value = provider::ipnetwork::prefix("2001:db8::1:2:3:4/64")
}
# result: "2001:db8::/64"
```

## Signature

```text
prefix(input string) string
```

## Arguments

1. `input` (String) Address to parse

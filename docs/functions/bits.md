---
page_title: "bits function - ipnetwork"
description: |-
  bits function
---

# function: bits

Return the prefix length (mask in bits) of a CIDR address.

## Example Usage

```terraform
output "ip_v4" {
  value = provider::ipnetwork::bits("192.0.2.128/24")
}
# result: 24

output "ip_v6" {
  value = provider::ipnetwork::bits("2001:db8::1:2:3:4/64")
}
# result: 64
```

## Signature

```text
bits(input string) number
```

## Arguments

1. `input` (String) Address to parse

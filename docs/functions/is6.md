---
page_title: "is6 function - ipnetwork"
description: |-
  is6 function
---

# function: is6

Reports whether an address is an IPv6 address.

Trim mask if input is in CIDR format.

## Example Usage

```terraform
output "ip_v4" {
  value = provider::ipnetwork::is6("192.0.2.128")
}
# result: false

output "ip_v6" {
  value = provider::ipnetwork::is6("2001:db8::1:2:3:4")
}
# result: true
```

## Signature

```text
is6(input string) boolean
```

## Arguments

1. `input` (String) Address to parse

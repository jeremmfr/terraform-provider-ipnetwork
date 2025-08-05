---
page_title: "is4 function - ipnetwork"
description: |-
  is4 function
---

# function: is4

Reports whether an address is an IPv4 address.

Trim mask if input is in CIDR format.

## Example Usage

```terraform
output "ip_v4" {
  value = provider::ipnetwork::is4("192.0.2.128")
}
# result: true

output "ip_v6" {
  value = provider::ipnetwork::is4("2001:db8::1:2:3:4")
}
# result: false
```

## Signature

```text
is4(input string) boolean
```

## Arguments

1. `input` (String) Address to parse

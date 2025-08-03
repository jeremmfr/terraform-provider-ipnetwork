---
page_title: "expand6 function - ipnetwork"
description: |-
  expand6 function
---

# function: expand6

Expand IPv6 address, with CIDR format or not, to long format (leading zeroes and no '::' compression).

No change with an IPv4 address.

## Example Usage

```terraform
output "ip" {
  value = provider::ipnetwork::expand6("2001:db8::1")
}
# result: "2001:0db8:0000:0000:0000:0000:0000:0001"

output "ip_cidr" {
  value = provider::ipnetwork::expand6("2001:db8::1/64")
}
# result: "2001:0db8:0000:0000:0000:0000:0000:0001/64"
```

## Signature

```text
expand6(input string) string
```

## Arguments

1. `input` (String) Address to parse

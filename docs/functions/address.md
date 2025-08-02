---
page_title: "address function - ipnetwork"
description: |-
  address function
---

# function: address

Validate an address with completion and cleanup of unwanted data and then proper format it.

completion, cleanup list:

- remove potential mask from CIDR format
- remove potential leading and trailing white space
- remove potential scoped zone
- add `0` decimal if missing one, two or three decimal(s) in IPv4 address

## Example Usage

```terraform
output "ip_cidr" {
  value = provider::ipnetwork::address("192.0.2.1/24")
}
# result: "192.0.2.1"

output "short_ip" {
  value = provider::ipnetwork::address("10")
}
# result: "10.0.0.0"

output "expanded_uppercase_ipv6" {
  value = provider::ipnetwork::address("2001:0DB8:0000:0000:0000:0000:0000:0000")
}
# result: "2001:db8::"
```

## Signature

```text
address(input string) string
```

## Arguments

1. `input` (String) Address to parse

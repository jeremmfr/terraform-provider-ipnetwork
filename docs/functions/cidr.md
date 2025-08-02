---
page_title: "cidr function - ipnetwork"
description: |-
  cidr function
---

# function: cidr

Validate a CIDR address with completion and replacement/cleanup of incorrect/unwanted data and
then proper format it.

completion, replacement/cleanup list:

- remove potential leading and trailing white space
- remove potential scoped zone
- add mask if missing:
  - `/0` for `0.0.0.0` and `::` address
  - `/32` for other IPv4 address
  - `/128` for other IPv6 address
- add `0` decimal if missing one, two or three decimal(s) in IPv4 address
- replace potential decimal notation mask to bits value for IPv4 address  
  (example `255.255.255.0` to `24`)

## Example Usage

```terraform
output "ip_only" {
  value = provider::ipnetwork::cidr("192.0.2.1")
}
# result: "192.0.2.1/32"

output "short_route" {
  value = provider::ipnetwork::cidr("10/8")
}
# result: "10.0.0.0/8"

output "expanded_uppercase_ipv6" {
  value = provider::ipnetwork::cidr("2001:0DB8:0000:0000:0000:0000:0000:0000/64")
}
# result: "2001:db8::/64"
```

## Signature

```text
cidr(input string) string
```

## Arguments

1. `input` (String) Address to parse

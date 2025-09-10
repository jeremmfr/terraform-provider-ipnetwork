---
page_title: "address_port function - ipnetwork"
description: |-
  address_port function
---

# function: address_port

Generate an ip:port string representation from IP address and port (add square brackets for IPv6 address).

Trim mask if `address` is in CIDR format and trim potential scoped zone for IPv6 address.

## Example Usage

```terraform
output "ip_v4" {
  value = provider::ipnetwork::address_port("192.0.2.128/24", 80)
}
# result: "192.0.2.128:80"

output "ip_v6" {
  value = provider::ipnetwork::address_port("2001:db8::0:1:2:3", 443)
}
# result: "[2001:db8::1:2:3]:443"
```

## Signature

```text
address_port(address string, port number) string
```

## Arguments

1. `address` (String) Address to parse
2. `port` (Number) Port to parse

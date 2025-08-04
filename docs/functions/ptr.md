---
page_title: "ptr function - ipnetwork"
description: |-
  ptr function
---

# function: ptr

Generate the PTR name from an address.

Trim mask if input is in CIDR format.  
Output string have `in-addr.arpa.` suffix for IPv4 address and `ip6.arpa.` suffix for IPv6 address.

## Example Usage

```terraform
output "ip_v4" {
  value = provider::ipnetwork::ptr("192.0.2.128")
}
# result: "128.2.0.192.in-addr.arpa."

output "ip_v6" {
  value = provider::ipnetwork::ptr("2001:db8::1:2:3:4")
}
# result: "4.0.0.0.3.0.0.0.2.0.0.0.1.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa."
```

## Signature

```text
ptr(input string) string
```

## Arguments

1. `input` (String) Address to parse

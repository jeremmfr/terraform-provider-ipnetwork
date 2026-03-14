---
page_title: "sort function - ipnetwork"
description: |-
  sort function
---

# function: sort

Sort a list of IP addresses with or without mask in numerical order.

When two entries share the same address, the address without mask comes first,
then CIDR addresses are sorted by mask length (shortest first).

## Example Usage

```terraform
# Sort IPv4 addresses numerically
output "ipv4" {
  value = provider::ipnetwork::sort(["10.0.0.1", "2.0.0.1", "192.168.0.1"])
}
# result: ["2.0.0.1", "10.0.0.1", "192.168.0.1"]

# Sort CIDR addresses by address then by mask length
output "cidr" {
  value = provider::ipnetwork::sort(["10.0.0.0/24", "10.0.0.0/8", "10.0.0.0/16"])
}
# result: ["10.0.0.0/8", "10.0.0.0/16", "10.0.0.0/24"]

# Mix addresses and CIDR addresses
output "mixed" {
  value = provider::ipnetwork::sort(["192.168.1.0/24", "10.0.0.1", "172.16.0.0/12"])
}
# result: ["10.0.0.1", "172.16.0.0/12", "192.168.1.0/24"]

# Sort IPv6 addresses
output "ipv6" {
  value = provider::ipnetwork::sort(["2001:db8::3", "2001:db8::1", "2001:db8::2"])
}
# result: ["2001:db8::1", "2001:db8::2", "2001:db8::3"]
```

## Signature

```text
sort(inputs list of string) list of string
```

## Arguments

1. `inputs` (List of String) List of IP addresses to sort

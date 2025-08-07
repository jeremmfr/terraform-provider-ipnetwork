---
page_title: "translate_4to6 function - ipnetwork"
description: |-
  translate_4to6 function
---

# function: translate_4to6

Translate an IPv4 address to an IPv6 address using an IPv6 prefix, as defined in [RFC 6052 section 2.2](https://tools.ietf.org/html/rfc6052#section-2.2).

Trim mask if `address` is in CIDR format.  
Mask of `prefix` address determines how the IPv4 address is embedded:

- `no mask` or `>64`: the IPv4 address is encoded in positions 96 to 127.
- `>56` and `<=64`: the IPv4 address is encoded in positions 72 to 103.
- `>48` and `<=56`: bits of the IPv4 address are encoded in positions 56 to 63,
  with the remaining 24 bits in position 72 to 95.
- `>40` and `<=48`: 16 bits of the IPv4 address are encoded in positions 48 to 63,
  with the remaining 16 bits in position 72 to 87.
- `>32` and `<=40`: 24 bits of the IPv4 address are encoded in positions 40 to 63,
  with the remaining 8 bits in position 72 to 79.
- `<=32`: the IPv4 address is encoded in positions 32 to 63.

## Example Usage

```terraform
output "nat46_96" {
  value = provider::ipnetwork::translate_4to6("192.0.2.33", "2001:db8:122:344::")
}
# result: "2001:db8:122:344::c000:221"

output "nat46_64" {
  value = provider::ipnetwork::translate_4to6("192.0.2.33", "2001:db8:122:344::/64")
}
# result: "2001:db8:122:344:c0:2:2100:0"
```

## Signature

```text
translate_4to6(address string, prefix string) string
```

## Arguments

1. `address` (String) Address to parse
1. `prefix` (String) Prefix address to parse

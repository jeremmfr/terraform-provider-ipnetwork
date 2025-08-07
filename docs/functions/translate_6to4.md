---
page_title: "translate_6to4 function - ipnetwork"
description: |-
  translate_6to4 function
---

# function: translate_6to4

Translate an IPv6 address to an IPv4 address, as defined in [RFC 6052 section 2.2](https://tools.ietf.org/html/rfc6052#section-2.2).

Mask of address determines how the IPv4 address is embedded:

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
output "nat64_96" {
  value = provider::ipnetwork::translate_6to4("2001:db8:122:344::c000:221")
}
# result: "192.0.2.33"

output "nat64_64" {
  value = provider::ipnetwork::translate_6to4("2001:db8:122:344:c0:2:2100:0/64")
}
# result: "192.0.2.33"
```

## Signature

```text
translate_6to4(input string) string
```

## Arguments

1. `input` (String) Address to parse

---
page_title: "generate6_eui64 function - ipnetwork"
description: |-
  generate6_eui64 function
---

# function: generate6_eui64

Generate an IPv6 address from MAC address with the modified EUI-64 format,
as defined in [RFC 4291 section 2.5.1](https://tools.ietf.org/html/rfc4291#section-2.5.1).

Trim mask if `prefix` is in CIDR format and trim potential scoped zone.  
If the latest 64 bits of `prefix` is not zero, they are still
overwrite by MAC address in modified EUI-64 format.

## Example Usage

```terraform
output "with_colon" {
  value = provider::ipnetwork::generate6_eui64("fe80::", "00:00:5e:53:53:00")
}
# result: "fe80::200:5eff:fe53:5300"

output "with_dash" {
  value = provider::ipnetwork::generate6_eui64("fe80::1234", "00-00-5E-00-53-00")
}
# result: "fe80::200:5eff:fe00:5300"
```

## Signature

```text
generate6_eui64(prefix string, mac string) string
```

## Arguments

1. `prefix` (String) IPv6 prefix address to parse
2. `mac` (String) MAC address to parse

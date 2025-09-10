---
page_title: "generate6_opaque function - ipnetwork"
description: |-
  generate6_opaque function
---

# function: generate6_opaque

Generate an IPv6 address with an opaque interface identifier,
as defined in [RFC 7217 section 5](https://tools.ietf.org/html/rfc7217#section-5).

Trim mask if `prefix` is in CIDR format and trim potential scoped zone.  
If the latest 64 bits of `prefix` is not zero, they are still
overwrite by the generated interface identifier.  
Use SHA256 as the pseudorandom function.

## Example Usage

```terraform
output "opaque" {
  value = provider::ipnetwork::generate6_opaque(
    "fe80::", "00-00-5E-00-53-00", null, null, "secret_key-secret_key",
  )
}
# result: "fe80::374e:8e0a:5de9:71cc"

output "opaque_1" {
  value = provider::ipnetwork::generate6_opaque(
    "fe80::", "00-00-5E-00-53-00", null, 1, "secret_key-secret_key",
  )
}
# result: "fe80::65e9:568f:1d6:f0a4"

output "other_net_iface" {
  value = provider::ipnetwork::generate6_opaque(
    "fe80::", "eth0", null, null, "secret_key-secret_key",
  )
}
# result: "fe80::85a5:13f3:2427:4229"

output "with_network_id" {
  value = provider::ipnetwork::generate6_opaque(
    "fe80::", "eth0", "123", null, "secret_key-secret_key",
  )
}
# result: "fe80::8707:476:3661:d360"
```

## Signature

```text
generate6_opaque(prefix string, net_iface string, network_id string, dad_counter number, secret_key string) string
```

## Arguments

1. `prefix` (String) IPv6 prefix address to parse
2. `net_iface` (String) Interface identifier
3. `network_id` (String) Network subnet identifier  
    allow `null` and consider as an empty string
4. `dad_counter` (Number) Counter to resolve DAD conflict  
    allow `null` and consider as 0
5. `secret_key` (String) Secret key

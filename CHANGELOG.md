<!-- markdownlint-disable-file MD013 -->
# changelog

## [v1.0.0](https://github.com/jeremmfr/terraform-provider-ipnetwork/tree/v1.0.0) (2025-08-09)

First release with this functions:

- `address(input string) string`: Validate an address.
- `bits(input string) number`: Return the prefix length of a CIDR address.
- `cidr(input string) string`: Validate a CIDR address.
- `contain(container string, address string) boolean`: Reports whether a prefix contains address(es).
- `equal_address(address_x string, address_y string) boolean`: Compare two address if there are equal.
- `equal_prefix(address_x string, address_y string) boolean`: Compare two CIDR addresses if they are in the same prefix.
- `expand6(input string) string`: Expand IPv6 address.
- `is4(input string) boolean`: Reports whether an address is an IPv4 address.
- `is6(input string) boolean`: Reports whether an address is an IPv6 address.
- `prefix(input string) string`: Canonicalize CIDR address.
- `ptr(input string) string`: Generate the PTR name from an address.
- `translate_4to6(address string, prefix string) string`: Translate an IPv4 address to an IPv6 address.
- `translate_6to4(input string) string`: Translate an IPv6 address to an IPv4 address.

<!-- markdownlint-disable-file MD013 -->
# changelog

## [v1.2.1](https://github.com/jeremmfr/terraform-provider-ipnetwork/tree/v1.2.1) (2025-10-26)

BUG FIXES:

* **function/is_public**: fix missing exclusion of `100:0:0:1::/64` (Dummy IPv6 Prefix - RFC9780) from public addresses

## [v1.2.0](https://github.com/jeremmfr/terraform-provider-ipnetwork/tree/v1.2.0) (2025-10-25)

FEATURES:

* add new functions:
  * `is_private(input string) boolean`: reports whether an address or prefix is private (internally routable).
  * `is_private_rfc1918(input string) boolean`: reports whether an address or prefix is in RFC1918 private address space.
  * `is_private_rfc4193(input string) boolean`: reports whether an address or prefix is in RFC4193 Unique Local Address (ULA) space.
  * `is_private_rfc6598(input string) boolean`: reports whether an address or prefix is in RFC6598 Shared Address Space.
  * `is_public(input string) boolean`: reports whether an address or prefix is public (globally routable).
  * `summarize(inputs set of string) list of string`: summarize IP prefixes.

## [v1.1.0](https://github.com/jeremmfr/terraform-provider-ipnetwork/tree/v1.1.0) (2025-09-10)

FEATURES:

* add new functions:
  * `address_port(address string, port number) string`: generate an ip:port string representation.
  * `generate6_eui64(prefix string, mac string) string`: generate an IPv6 address from MAC address with the modified EUI-64 format.
  * `generate6_opaque(prefix string, net_iface string, network_id string, dad_counter number, secret_key string) string`: generate an IPv6 address with an opaque interface identifier.

ENHANCEMENTS:

* release now with Go 1.25.

## [v1.0.0](https://github.com/jeremmfr/terraform-provider-ipnetwork/tree/v1.0.0) (2025-08-09)

First release with this functions:

* `address(input string) string`: Validate an address.
* `bits(input string) number`: Return the prefix length of a CIDR address.
* `cidr(input string) string`: Validate a CIDR address.
* `contain(container string, address string) boolean`: Reports whether a prefix contains address(es).
* `equal_address(address_x string, address_y string) boolean`: Compare two address if there are equal.
* `equal_prefix(address_x string, address_y string) boolean`: Compare two CIDR addresses if they are in the same prefix.
* `expand6(input string) string`: Expand IPv6 address.
* `is4(input string) boolean`: Reports whether an address is an IPv4 address.
* `is6(input string) boolean`: Reports whether an address is an IPv6 address.
* `prefix(input string) string`: Canonicalize CIDR address.
* `ptr(input string) string`: Generate the PTR name from an address.
* `translate_4to6(address string, prefix string) string`: Translate an IPv4 address to an IPv6 address.
* `translate_6to4(input string) string`: Translate an IPv6 address to an IPv4 address.

<!-- markdownlint-disable-file MD013 MD041 -->
FEATURES:

* add `generate6_eui64(prefix string, mac string) string` function to generate an IPv6 address from MAC address with the modified EUI-64 format
* add `generate6_opaque(prefix string, net_iface string, network_id string, dad_counter number, secret_key string) string` function to generate an IPv6 address with an opaque interface identifier

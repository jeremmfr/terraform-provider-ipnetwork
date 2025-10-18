package provider

import (
	"net/netip"
	"testing"
)

func TestAddressIsPrivateRFC4193(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input  string
		expect bool
	}{
		"public_ipv4": {
			input:  "8.8.8.8",
			expect: false,
		},
		"private_ipv4": {
			input:  "192.168.1.1",
			expect: false,
		},
		"public_ipv6_google": {
			input:  "2001:4860:4860::8888",
			expect: false,
		},
		"public_ipv6_cloudflare": {
			input:  "2606:4700:4700::1111",
			expect: false,
		},
		"ipv6_ula_fc_start": {
			input:  "fc00::1",
			expect: true,
		},
		"ipv6_ula_fc_mid": {
			input:  "fc80:1234:5678:9abc:def0:1234:5678:9abc",
			expect: true,
		},
		"ipv6_ula_fd_start": {
			input:  "fd00::1",
			expect: true,
		},
		"ipv6_ula_fd_mid": {
			input:  "fd12:3456:789a:bcde::1",
			expect: true,
		},
		"ipv6_ula_fd_end": {
			input:  "fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			expect: true,
		},
		"not_ula_adjacent_before": {
			input:  "fbff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			expect: false,
		},
		"not_ula_adjacent_after": {
			input:  "fe00::1",
			expect: false,
		},
		"not_ula_link_local": {
			input:  "fe80::1",
			expect: false,
		},
		"not_ula_loopback": {
			input:  "::1",
			expect: false,
		},
		"not_ula_unspecified": {
			input:  "::",
			expect: false,
		},
		"not_ula_multicast": {
			input:  "ff02::1",
			expect: false,
		},
		"not_ula_documentation": {
			input:  "2001:db8::1",
			expect: false,
		},
		"ipv4_mapped_ipv4": {
			input:  "::ffff:192.168.1.1",
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := addressIsPrivateRFC4193(netip.MustParseAddr(test.input))
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

func TestPrefixIsPrivateRFC4193(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input  string
		expect bool
	}{
		"public_ipv4": {
			input:  "8.8.8.0/24",
			expect: false,
		},
		"private_ipv4": {
			input:  "192.168.0.0/16",
			expect: false,
		},
		"public_ipv6": {
			input:  "2001:4860::/48",
			expect: false,
		},
		"ipv6_ula_7": {
			input:  "fc00::/7",
			expect: true,
		},
		"ipv6_ula_fc_8": {
			input:  "fc00::/8",
			expect: true,
		},
		"ipv6_ula_fd_8": {
			input:  "fd00::/8",
			expect: true,
		},
		"ipv6_ula_48": {
			input:  "fd12:3456:789a::/48",
			expect: true,
		},
		"ipv6_ula_64": {
			input:  "fd12:3456:789a:bcde::/64",
			expect: true,
		},
		"ipv6_ula_128": {
			input:  "fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128",
			expect: true,
		},
		"not_ula_6": {
			input:  "fc00::/6",
			expect: false,
		},
		"not_ula_adjacent_before": {
			input:  "fbff::/16",
			expect: false,
		},
		"not_ula_adjacent_after": {
			input:  "fe00::/16",
			expect: false,
		},
		"not_ula_link_local": {
			input:  "fe80::/10",
			expect: false,
		},
		"not_ula_loopback": {
			input:  "::1/128",
			expect: false,
		},
		"not_ula_multicast": {
			input:  "ff00::/8",
			expect: false,
		},
		"not_ula_documentation": {
			input:  "2001:db8::/32",
			expect: false,
		},
		"ipv4_mapped": {
			input:  "::ffff:192.168.0.0/112",
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := prefixIsPrivateRFC4193(netip.MustParsePrefix(test.input))
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

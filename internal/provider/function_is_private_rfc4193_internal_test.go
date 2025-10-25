package provider

import (
	"net/netip"
	"testing"
)

func TestAddressIsPrivateRFC4193(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  netip.Addr
		expect bool
	}

	tests := map[string]testCase{
		"public_ipv4": {
			input:  netip.MustParseAddr("8.8.8.8"),
			expect: false,
		},
		"private_ipv4": {
			input:  netip.MustParseAddr("192.168.1.1"),
			expect: false,
		},
		"public_ipv6_google": {
			input:  netip.MustParseAddr("2001:4860:4860::8888"),
			expect: false,
		},
		"public_ipv6_cloudflare": {
			input:  netip.MustParseAddr("2606:4700:4700::1111"),
			expect: false,
		},
		"ipv6_ula_fc_start": {
			input:  netip.MustParseAddr("fc00::1"),
			expect: true,
		},
		"ipv6_ula_fc_mid": {
			input:  netip.MustParseAddr("fc80:1234:5678:9abc:def0:1234:5678:9abc"),
			expect: true,
		},
		"ipv6_ula_fd_start": {
			input:  netip.MustParseAddr("fd00::1"),
			expect: true,
		},
		"ipv6_ula_fd_mid": {
			input:  netip.MustParseAddr("fd12:3456:789a:bcde::1"),
			expect: true,
		},
		"ipv6_ula_fd_end": {
			input:  netip.MustParseAddr("fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
			expect: true,
		},
		"not_ula_adjacent_before": {
			input:  netip.MustParseAddr("fbff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
			expect: false,
		},
		"not_ula_adjacent_after": {
			input:  netip.MustParseAddr("fe00::1"),
			expect: false,
		},
		"not_ula_link_local": {
			input:  netip.MustParseAddr("fe80::1"),
			expect: false,
		},
		"not_ula_loopback": {
			input:  netip.MustParseAddr("::1"),
			expect: false,
		},
		"not_ula_unspecified": {
			input:  netip.MustParseAddr("::"),
			expect: false,
		},
		"not_ula_multicast": {
			input:  netip.MustParseAddr("ff02::1"),
			expect: false,
		},
		"not_ula_documentation": {
			input:  netip.MustParseAddr("2001:db8::1"),
			expect: false,
		},
		"ipv4_mapped_ipv4": {
			input:  netip.MustParseAddr("::ffff:192.168.1.1"),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := addressIsPrivateRFC4193(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

func TestPrefixIsPrivateRFC4193(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  netip.Prefix
		expect bool
	}

	tests := map[string]testCase{
		"public_ipv4": {
			input:  netip.MustParsePrefix("8.8.8.0/24"),
			expect: false,
		},
		"private_ipv4": {
			input:  netip.MustParsePrefix("192.168.0.0/16"),
			expect: false,
		},
		"public_ipv6": {
			input:  netip.MustParsePrefix("2001:4860::/48"),
			expect: false,
		},
		"ipv6_ula_7": {
			input:  netip.MustParsePrefix("fc00::/7"),
			expect: true,
		},
		"ipv6_ula_fc_8": {
			input:  netip.MustParsePrefix("fc00::/8"),
			expect: true,
		},
		"ipv6_ula_fd_8": {
			input:  netip.MustParsePrefix("fd00::/8"),
			expect: true,
		},
		"ipv6_ula_48": {
			input:  netip.MustParsePrefix("fd12:3456:789a::/48"),
			expect: true,
		},
		"ipv6_ula_64": {
			input:  netip.MustParsePrefix("fd12:3456:789a:bcde::/64"),
			expect: true,
		},
		"ipv6_ula_128": {
			input:  netip.MustParsePrefix("fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128"),
			expect: true,
		},
		"not_ula_6": {
			input:  netip.MustParsePrefix("fc00::/6"),
			expect: false,
		},
		"not_ula_adjacent_before": {
			input:  netip.MustParsePrefix("fbff::/16"),
			expect: false,
		},
		"not_ula_adjacent_after": {
			input:  netip.MustParsePrefix("fe00::/16"),
			expect: false,
		},
		"not_ula_link_local": {
			input:  netip.MustParsePrefix("fe80::/10"),
			expect: false,
		},
		"not_ula_loopback": {
			input:  netip.MustParsePrefix("::1/128"),
			expect: false,
		},
		"not_ula_multicast": {
			input:  netip.MustParsePrefix("ff00::/8"),
			expect: false,
		},
		"not_ula_documentation": {
			input:  netip.MustParsePrefix("2001:db8::/32"),
			expect: false,
		},
		"ipv4_mapped": {
			input:  netip.MustParsePrefix("::ffff:192.168.0.0/112"),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := prefixIsPrivateRFC4193(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

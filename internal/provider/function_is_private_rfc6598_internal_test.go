package provider

import (
	"net/netip"
	"testing"
)

func TestAddressIsPrivateRFC6598(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  netip.Addr
		expect bool
	}

	tests := map[string]testCase{
		"public_ipv4_google": {
			input:  netip.MustParseAddr("8.8.8.8"),
			expect: false,
		},
		"public_ipv4_cloudflare": {
			input:  netip.MustParseAddr("1.1.1.1"),
			expect: false,
		},
		"cgn_rfc6598_start": {
			input:  netip.MustParseAddr("100.64.0.0"),
			expect: true,
		},
		"cgn_rfc6598_mid": {
			input:  netip.MustParseAddr("100.100.50.25"),
			expect: true,
		},
		"cgn_rfc6598_end": {
			input:  netip.MustParseAddr("100.127.255.255"),
			expect: true,
		},
		"not_cgn_adjacent_before": {
			input:  netip.MustParseAddr("100.63.255.255"),
			expect: false,
		},
		"not_cgn_adjacent_after": {
			input:  netip.MustParseAddr("100.128.0.0"),
			expect: false,
		},
		"not_cgn_private_10": {
			input:  netip.MustParseAddr("10.0.0.1"),
			expect: false,
		},
		"not_cgn_private_172": {
			input:  netip.MustParseAddr("172.16.0.1"),
			expect: false,
		},
		"not_cgn_private_192": {
			input:  netip.MustParseAddr("192.168.1.1"),
			expect: false,
		},
		"not_cgn_loopback": {
			input:  netip.MustParseAddr("127.0.0.1"),
			expect: false,
		},
		"not_cgn_link_local": {
			input:  netip.MustParseAddr("169.254.1.1"),
			expect: false,
		},
		"public_ipv6_google": {
			input:  netip.MustParseAddr("2001:4860:4860::8888"),
			expect: false,
		},
		"ipv6_ula": {
			input:  netip.MustParseAddr("fd00::1"),
			expect: false,
		},
		"ipv4_mapped_cgn": {
			input:  netip.MustParseAddr("::ffff:100.64.0.1"),
			expect: true,
		},
		"ipv4_mapped_cgn_mid": {
			input:  netip.MustParseAddr("::ffff:100.100.50.25"),
			expect: true,
		},
		"ipv4_mapped_public": {
			input:  netip.MustParseAddr("::ffff:8.8.8.8"),
			expect: false,
		},
		"ipv4_mapped_private": {
			input:  netip.MustParseAddr("::ffff:192.168.1.1"),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := addressIsPrivateRFC6598(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

func TestPrefixIsPrivateRFC6598(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  netip.Prefix
		expect bool
	}

	tests := map[string]testCase{
		"public_ipv4_8": {
			input:  netip.MustParsePrefix("8.0.0.0/8"),
			expect: false,
		},
		"public_ipv4_24": {
			input:  netip.MustParsePrefix("1.1.1.0/24"),
			expect: false,
		},
		"cgn_rfc6598_10": {
			input:  netip.MustParsePrefix("100.64.0.0/10"),
			expect: true,
		},
		"cgn_rfc6598_16": {
			input:  netip.MustParsePrefix("100.64.0.0/16"),
			expect: true,
		},
		"cgn_rfc6598_24": {
			input:  netip.MustParsePrefix("100.100.50.0/24"),
			expect: true,
		},
		"cgn_rfc6598_32": {
			input:  netip.MustParsePrefix("100.127.255.255/32"),
			expect: true,
		},
		"not_cgn_9": {
			input:  netip.MustParsePrefix("100.64.0.0/9"),
			expect: false,
		},
		"not_cgn_adjacent_before": {
			input:  netip.MustParsePrefix("100.0.0.0/10"),
			expect: false,
		},
		"not_cgn_adjacent_after": {
			input:  netip.MustParsePrefix("100.128.0.0/10"),
			expect: false,
		},
		"not_cgn_private_10": {
			input:  netip.MustParsePrefix("10.0.0.0/8"),
			expect: false,
		},
		"not_cgn_private_172": {
			input:  netip.MustParsePrefix("172.16.0.0/12"),
			expect: false,
		},
		"not_cgn_private_192": {
			input:  netip.MustParsePrefix("192.168.0.0/16"),
			expect: false,
		},
		"not_cgn_loopback": {
			input:  netip.MustParsePrefix("127.0.0.0/8"),
			expect: false,
		},
		"not_cgn_link_local": {
			input:  netip.MustParsePrefix("169.254.0.0/16"),
			expect: false,
		},
		"public_ipv6": {
			input:  netip.MustParsePrefix("2001:4860::/48"),
			expect: false,
		},
		"ipv6_ula": {
			input:  netip.MustParsePrefix("fd00::/7"),
			expect: false,
		},
		"ipv4_mapped_cgn": {
			input:  netip.MustParsePrefix("::ffff:100.64.0.0/106"),
			expect: true,
		},
		"ipv4_mapped_cgn_16": {
			input:  netip.MustParsePrefix("::ffff:100.100.0.0/112"),
			expect: true,
		},
		"ipv4_mapped_public": {
			input:  netip.MustParsePrefix("::ffff:8.8.8.8/128"),
			expect: false,
		},
		"ipv4_mapped_private": {
			input:  netip.MustParsePrefix("::ffff:192.168.0.0/112"),
			expect: false,
		},
		"ipv4_mapped_too_short": {
			input:  netip.MustParsePrefix("::ffff:100.64.0.0/95"),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := prefixIsPrivateRFC6598(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

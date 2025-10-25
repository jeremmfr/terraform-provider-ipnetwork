package provider

import (
	"net/netip"
	"testing"
)

func TestAddressIsPrivateRFC1918(t *testing.T) {
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
		"private_10_start": {
			input:  netip.MustParseAddr("10.0.0.0"),
			expect: true,
		},
		"private_10_mid": {
			input:  netip.MustParseAddr("10.128.64.32"),
			expect: true,
		},
		"private_10_end": {
			input:  netip.MustParseAddr("10.255.255.255"),
			expect: true,
		},
		"private_172_16_start": {
			input:  netip.MustParseAddr("172.16.0.0"),
			expect: true,
		},
		"private_172_16_mid": {
			input:  netip.MustParseAddr("172.20.5.10"),
			expect: true,
		},
		"private_172_31_end": {
			input:  netip.MustParseAddr("172.31.255.255"),
			expect: true,
		},
		"private_192_168_start": {
			input:  netip.MustParseAddr("192.168.0.0"),
			expect: true,
		},
		"private_192_168_mid": {
			input:  netip.MustParseAddr("192.168.128.64"),
			expect: true,
		},
		"private_192_168_end": {
			input:  netip.MustParseAddr("192.168.255.255"),
			expect: true,
		},
		"not_private_adjacent_10_before": {
			input:  netip.MustParseAddr("9.255.255.255"),
			expect: false,
		},
		"not_private_adjacent_10_after": {
			input:  netip.MustParseAddr("11.0.0.0"),
			expect: false,
		},
		"not_private_172_15": {
			input:  netip.MustParseAddr("172.15.255.255"),
			expect: false,
		},
		"not_private_172_32": {
			input:  netip.MustParseAddr("172.32.0.0"),
			expect: false,
		},
		"not_private_adjacent_192_168_before": {
			input:  netip.MustParseAddr("192.167.255.255"),
			expect: false,
		},
		"not_private_adjacent_192_168_after": {
			input:  netip.MustParseAddr("192.169.0.0"),
			expect: false,
		},
		"not_private_cgn_rfc6598": {
			input:  netip.MustParseAddr("100.64.0.1"),
			expect: false,
		},
		"not_private_loopback": {
			input:  netip.MustParseAddr("127.0.0.1"),
			expect: false,
		},
		"not_private_unspecified": {
			input:  netip.MustParseAddr("0.0.0.0"),
			expect: false,
		},
		"not_private_link_local": {
			input:  netip.MustParseAddr("169.254.1.1"),
			expect: false,
		},
		"not_private_multicast": {
			input:  netip.MustParseAddr("224.0.0.1"),
			expect: false,
		},
		"not_private_broadcast": {
			input:  netip.MustParseAddr("255.255.255.255"),
			expect: false,
		},
		"not_private_documentation_testnet1": {
			input:  netip.MustParseAddr("192.0.2.1"),
			expect: false,
		},
		"not_private_documentation_testnet2": {
			input:  netip.MustParseAddr("198.51.100.1"),
			expect: false,
		},
		"not_private_documentation_testnet3": {
			input:  netip.MustParseAddr("203.0.113.1"),
			expect: false,
		},
		"not_private_benchmarking": {
			input:  netip.MustParseAddr("198.18.1.1"),
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
		"ipv6_ula": {
			input:  netip.MustParseAddr("fd00::1"),
			expect: false,
		},
		"ipv6_loopback": {
			input:  netip.MustParseAddr("::1"),
			expect: false,
		},
		"ipv6_unspecified": {
			input:  netip.MustParseAddr("::"),
			expect: false,
		},
		"ipv6_link_local": {
			input:  netip.MustParseAddr("fe80::1"),
			expect: false,
		},
		"ipv6_multicast": {
			input:  netip.MustParseAddr("ff02::1"),
			expect: false,
		},
		"ipv4_mapped_private_10": {
			input:  netip.MustParseAddr("::ffff:10.0.0.1"),
			expect: true,
		},
		"ipv4_mapped_private_172": {
			input:  netip.MustParseAddr("::ffff:172.16.0.1"),
			expect: true,
		},
		"ipv4_mapped_private_192": {
			input:  netip.MustParseAddr("::ffff:192.168.1.1"),
			expect: true,
		},
		"ipv4_mapped_public": {
			input:  netip.MustParseAddr("::ffff:8.8.8.8"),
			expect: false,
		},
		"ipv4_mapped_cgn": {
			input:  netip.MustParseAddr("::ffff:100.64.0.1"),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := addressIsPrivateRFC1918(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

func TestPrefixIsPrivateRFC1918(t *testing.T) {
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
		"private_10_8": {
			input:  netip.MustParsePrefix("10.0.0.0/8"),
			expect: true,
		},
		"private_10_16": {
			input:  netip.MustParsePrefix("10.0.0.0/16"),
			expect: true,
		},
		"private_10_24": {
			input:  netip.MustParsePrefix("10.5.10.0/24"),
			expect: true,
		},
		"private_10_32": {
			input:  netip.MustParsePrefix("10.128.64.32/32"),
			expect: true,
		},
		"private_172_16_12": {
			input:  netip.MustParsePrefix("172.16.0.0/12"),
			expect: true,
		},
		"private_172_16_16": {
			input:  netip.MustParsePrefix("172.16.0.0/16"),
			expect: true,
		},
		"private_172_24": {
			input:  netip.MustParsePrefix("172.20.5.0/24"),
			expect: true,
		},
		"private_172_31_24": {
			input:  netip.MustParsePrefix("172.31.255.0/24"),
			expect: true,
		},
		"private_192_168_16": {
			input:  netip.MustParsePrefix("192.168.0.0/16"),
			expect: true,
		},
		"private_192_168_24": {
			input:  netip.MustParsePrefix("192.168.1.0/24"),
			expect: true,
		},
		"private_192_168_32": {
			input:  netip.MustParsePrefix("192.168.128.64/32"),
			expect: true,
		},
		"not_private_10_7": {
			input:  netip.MustParsePrefix("10.0.0.0/7"),
			expect: false,
		},
		"not_private_172_11": {
			input:  netip.MustParsePrefix("172.16.0.0/11"),
			expect: false,
		},
		"not_private_192_168_15": {
			input:  netip.MustParsePrefix("192.168.0.0/15"),
			expect: false,
		},
		"not_private_adjacent_10_before": {
			input:  netip.MustParsePrefix("9.0.0.0/8"),
			expect: false,
		},
		"not_private_adjacent_10_after": {
			input:  netip.MustParsePrefix("11.0.0.0/8"),
			expect: false,
		},
		"not_private_172_15": {
			input:  netip.MustParsePrefix("172.15.0.0/16"),
			expect: false,
		},
		"not_private_172_32": {
			input:  netip.MustParsePrefix("172.32.0.0/12"),
			expect: false,
		},
		"not_private_adjacent_192_168_before": {
			input:  netip.MustParsePrefix("192.167.0.0/16"),
			expect: false,
		},
		"not_private_adjacent_192_168_after": {
			input:  netip.MustParsePrefix("192.169.0.0/16"),
			expect: false,
		},
		"not_private_cgn_rfc6598": {
			input:  netip.MustParsePrefix("100.64.0.0/10"),
			expect: false,
		},
		"not_private_loopback": {
			input:  netip.MustParsePrefix("127.0.0.0/8"),
			expect: false,
		},
		"not_private_link_local": {
			input:  netip.MustParsePrefix("169.254.0.0/16"),
			expect: false,
		},
		"not_private_multicast": {
			input:  netip.MustParsePrefix("224.0.0.0/4"),
			expect: false,
		},
		"not_private_documentation_testnet1": {
			input:  netip.MustParsePrefix("192.0.2.0/24"),
			expect: false,
		},
		"not_private_benchmarking": {
			input:  netip.MustParsePrefix("198.18.0.0/15"),
			expect: false,
		},
		"contains_private_10": {
			input:  netip.MustParsePrefix("8.0.0.0/6"),
			expect: false,
		},
		"contains_private_172": {
			input:  netip.MustParsePrefix("172.0.0.0/9"),
			expect: false,
		},
		"contains_private_192_168": {
			input:  netip.MustParsePrefix("192.0.0.0/8"),
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
		"ipv4_mapped_private_10": {
			input:  netip.MustParsePrefix("::ffff:10.0.0.0/104"),
			expect: true,
		},
		"ipv4_mapped_private_172": {
			input:  netip.MustParsePrefix("::ffff:172.16.0.0/108"),
			expect: true,
		},
		"ipv4_mapped_private_192": {
			input:  netip.MustParsePrefix("::ffff:192.168.0.0/112"),
			expect: true,
		},
		"ipv4_mapped_public": {
			input:  netip.MustParsePrefix("::ffff:8.8.8.8/128"),
			expect: false,
		},
		"ipv4_mapped_cgn": {
			input:  netip.MustParsePrefix("::ffff:100.64.0.0/106"),
			expect: false,
		},
		"ipv4_mapped_too_short": {
			input:  netip.MustParsePrefix("::ffff:10.0.0.0/95"),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := prefixIsPrivateRFC1918(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

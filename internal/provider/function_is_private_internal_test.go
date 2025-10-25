package provider

import (
	"net/netip"
	"testing"
)

func TestAddressV4IsPrivate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  netip.Addr
		expect bool
	}

	tests := map[string]testCase{
		"public_google_dns": {
			input:  netip.MustParseAddr("8.8.8.8"),
			expect: false,
		},
		"public_cloudflare_dns": {
			input:  netip.MustParseAddr("1.1.1.1"),
			expect: false,
		},
		"private_10": {
			input:  netip.MustParseAddr("10.0.0.1"),
			expect: true,
		},
		"private_10_end": {
			input:  netip.MustParseAddr("10.255.255.254"),
			expect: true,
		},
		"private_172_16": {
			input:  netip.MustParseAddr("172.16.0.1"),
			expect: true,
		},
		"private_172_mid": {
			input:  netip.MustParseAddr("172.20.5.10"),
			expect: true,
		},
		"private_172_end": {
			input:  netip.MustParseAddr("172.31.255.254"),
			expect: true,
		},
		"private_192_168": {
			input:  netip.MustParseAddr("192.168.1.1"),
			expect: true,
		},
		"private_192_168_end": {
			input:  netip.MustParseAddr("192.168.255.254"),
			expect: true,
		},
		"cgn": {
			input:  netip.MustParseAddr("100.64.0.1"),
			expect: true,
		},
		"cgn_mid": {
			input:  netip.MustParseAddr("100.100.50.25"),
			expect: true,
		},
		"cgn_end": {
			input:  netip.MustParseAddr("100.127.255.254"),
			expect: true,
		},
		"benchmarking": {
			input:  netip.MustParseAddr("198.18.0.1"),
			expect: true,
		},
		"benchmarking_mid": {
			input:  netip.MustParseAddr("198.19.100.50"),
			expect: true,
		},
		"benchmarking_end": {
			input:  netip.MustParseAddr("198.19.255.254"),
			expect: true,
		},
		"not_private_loopback": {
			input:  netip.MustParseAddr("127.0.0.1"),
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
		"not_private_doc_testnet1": {
			input:  netip.MustParseAddr("192.0.2.1"),
			expect: false,
		},
		"not_private_doc_testnet2": {
			input:  netip.MustParseAddr("198.51.100.1"),
			expect: false,
		},
		"adjacent_to_10_before": {
			input:  netip.MustParseAddr("9.255.255.255"),
			expect: false,
		},
		"adjacent_to_10_after": {
			input:  netip.MustParseAddr("11.0.0.0"),
			expect: false,
		},
		"adjacent_to_172_before": {
			input:  netip.MustParseAddr("172.15.255.255"),
			expect: false,
		},
		"adjacent_to_172_after": {
			input:  netip.MustParseAddr("172.32.0.0"),
			expect: false,
		},
		"adjacent_to_192_168_after": {
			input:  netip.MustParseAddr("192.169.0.0"),
			expect: false,
		},
		"adjacent_to_cgn_before": {
			input:  netip.MustParseAddr("100.63.255.255"),
			expect: false,
		},
		"adjacent_to_cgn_after": {
			input:  netip.MustParseAddr("100.128.0.0"),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := addressV4IsPrivate(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

func TestAddressV6IsPrivate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  netip.Addr
		expect bool
	}

	tests := map[string]testCase{
		"public_ipv6_google": {
			input:  netip.MustParseAddr("2001:4860:4860::8888"),
			expect: false,
		},
		"public_ipv6_cloudflare": {
			input:  netip.MustParseAddr("2606:4700:4700::1111"),
			expect: false,
		},
		"ipv6_ula_fc": {
			input:  netip.MustParseAddr("fc00::1"),
			expect: true,
		},
		"ipv6_ula_fd": {
			input:  netip.MustParseAddr("fd00::1"),
			expect: true,
		},
		"ipv6_ula_expanded": {
			input:  netip.MustParseAddr("fd12:3456:789a:bcde::1"),
			expect: true,
		},
		"ipv6_discard_prefix": {
			input:  netip.MustParseAddr("100::"),
			expect: true,
		},
		"ipv6_discard_prefix_2": {
			input:  netip.MustParseAddr("100::1:2:3:4"),
			expect: true,
		},
		"ipv6_discard_end": {
			input:  netip.MustParseAddr("100::ffff:ffff:ffff:ffff"),
			expect: true,
		},
		"ipv6_translation": {
			input:  netip.MustParseAddr("64:ff9b:1::192.0.2.1"),
			expect: true,
		},
		"ipv6_translation_expanded": {
			input:  netip.MustParseAddr("64:ff9b:1:ffff:ffff:ffff:ffff:ffff"),
			expect: true,
		},
		"ipv6_srv6": {
			input:  netip.MustParseAddr("5f00::1"),
			expect: true,
		},
		"ipv6_srv6_end": {
			input:  netip.MustParseAddr("5f00:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
			expect: true,
		},
		"ipv6_benchmarking": {
			input:  netip.MustParseAddr("2001:2::1"),
			expect: true,
		},
		"ipv6_benchmarking_expanded": {
			input:  netip.MustParseAddr("2001:2:0:ffff:ffff:ffff:ffff:ffff"),
			expect: true,
		},
		"not_private_loopback": {
			input:  netip.MustParseAddr("::1"),
			expect: false,
		},
		"not_private_unspecified": {
			input:  netip.MustParseAddr("::"),
			expect: false,
		},
		"not_private_link_local": {
			input:  netip.MustParseAddr("fe80::1"),
			expect: false,
		},
		"not_private_multicast": {
			input:  netip.MustParseAddr("ff02::1"),
			expect: false,
		},
		"not_private_doc_db8": {
			input:  netip.MustParseAddr("2001:db8::1"),
			expect: false,
		},
		"ipv4_mapped_private": {
			input:  netip.MustParseAddr("::ffff:192.168.1.1"),
			expect: true,
		},
		"ipv4_mapped_cgn": {
			input:  netip.MustParseAddr("::ffff:100.64.0.1"),
			expect: true,
		},
		"ipv4_mapped_public": {
			input:  netip.MustParseAddr("::ffff:8.8.8.8"),
			expect: false,
		},
		"adjacent_to_discard_before": {
			input:  netip.MustParseAddr("ff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
			expect: false,
		},
		"adjacent_to_discard_after": {
			input:  netip.MustParseAddr("100:0:1::"),
			expect: false,
		},
		"adjacent_to_ula_before": {
			input:  netip.MustParseAddr("fbff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"),
			expect: false,
		},
		"adjacent_to_ula_after": {
			input:  netip.MustParseAddr("fe00::"),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := addressV6IsPrivate(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

func TestPrefixV4IsPrivate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  netip.Prefix
		expect bool
	}

	tests := map[string]testCase{
		"public_24": {
			input:  netip.MustParsePrefix("1.1.1.0/24"),
			expect: false,
		},
		"public_8": {
			input:  netip.MustParsePrefix("8.0.0.0/8"),
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
		"private_192_168_16": {
			input:  netip.MustParsePrefix("192.168.0.0/16"),
			expect: true,
		},
		"private_192_168_24": {
			input:  netip.MustParsePrefix("192.168.1.0/24"),
			expect: true,
		},
		"private_192_168_15": {
			input:  netip.MustParsePrefix("192.168.0.0/15"),
			expect: false,
		},
		"cgn_10": {
			input:  netip.MustParsePrefix("100.64.0.0/10"),
			expect: true,
		},
		"cgn_16": {
			input:  netip.MustParsePrefix("100.64.0.0/16"),
			expect: true,
		},
		"cgn_24": {
			input:  netip.MustParsePrefix("100.100.50.0/24"),
			expect: true,
		},
		"benchmarking_15": {
			input:  netip.MustParsePrefix("198.18.0.0/15"),
			expect: true,
		},
		"benchmarking_16": {
			input:  netip.MustParsePrefix("198.19.0.0/16"),
			expect: true,
		},
		"benchmarking_24": {
			input:  netip.MustParsePrefix("198.19.100.0/24"),
			expect: true,
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
			input:  netip.MustParsePrefix("192.160.0.0/12"),
			expect: false,
		},
		"contains_private_192_0_0_8": {
			input:  netip.MustParsePrefix("192.0.0.0/8"),
			expect: false,
		},
		"not_private_loopback_8": {
			input:  netip.MustParsePrefix("127.0.0.0/8"),
			expect: false,
		},
		"not_private_link_local_16": {
			input:  netip.MustParsePrefix("169.254.0.0/16"),
			expect: false,
		},
		"not_private_multicast_4": {
			input:  netip.MustParsePrefix("224.0.0.0/4"),
			expect: false,
		},
		"not_private_doc_testnet1": {
			input:  netip.MustParsePrefix("192.0.2.0/24"),
			expect: false,
		},
		"adjacent_to_10_public": {
			input:  netip.MustParsePrefix("11.0.0.0/8"),
			expect: false,
		},
		"adjacent_to_172_public": {
			input:  netip.MustParsePrefix("172.32.0.0/12"),
			expect: false,
		},
		"adjacent_to_192_168_public": {
			input:  netip.MustParsePrefix("192.169.0.0/16"),
			expect: false,
		},
		"adjacent_to_cgn_before": {
			input:  netip.MustParsePrefix("100.0.0.0/10"),
			expect: false,
		},
		"adjacent_to_cgn_after": {
			input:  netip.MustParsePrefix("100.128.0.0/10"),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := prefixV4IsPrivate(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

func TestPrefixV6IsPrivate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  netip.Prefix
		expect bool
	}

	tests := map[string]testCase{
		"public_google_48": {
			input:  netip.MustParsePrefix("2001:4860::/48"),
			expect: false,
		},
		"public_cloudflare_32": {
			input:  netip.MustParsePrefix("2606:4700::/32"),
			expect: false,
		},
		"ula_7": {
			input:  netip.MustParsePrefix("fc00::/7"),
			expect: true,
		},
		"ula_fd_8": {
			input:  netip.MustParsePrefix("fd00::/8"),
			expect: true,
		},
		"ula_48": {
			input:  netip.MustParsePrefix("fd12:3456:789a::/48"),
			expect: true,
		},
		"ula_64": {
			input:  netip.MustParsePrefix("fd12:3456:789a:bcde::/64"),
			expect: true,
		},
		"discard_64": {
			input:  netip.MustParsePrefix("100::/64"),
			expect: true,
		},
		"translation_48": {
			input:  netip.MustParsePrefix("64:ff9b:1::/48"),
			expect: true,
		},
		"translation_64": {
			input:  netip.MustParsePrefix("64:ff9b:1:1234::/64"),
			expect: true,
		},
		"srv6_16": {
			input:  netip.MustParsePrefix("5f00::/16"),
			expect: true,
		},
		"srv6_32": {
			input:  netip.MustParsePrefix("5f00:1234::/32"),
			expect: true,
		},
		"srv6_64": {
			input:  netip.MustParsePrefix("5f00:1234:5678:abcd::/64"),
			expect: true,
		},
		"benchmarking_48": {
			input:  netip.MustParsePrefix("2001:2::/48"),
			expect: true,
		},
		"benchmarking_64": {
			input:  netip.MustParsePrefix("2001:2:0:1234::/64"),
			expect: true,
		},
		"contains_ula": {
			input:  netip.MustParsePrefix("fc00::/6"),
			expect: false,
		},
		"contains_benchmarking": {
			input:  netip.MustParsePrefix("2001::/16"),
			expect: false,
		},
		"not_private_loopback_128": {
			input:  netip.MustParsePrefix("::1/128"),
			expect: false,
		},
		"not_private_link_local_10": {
			input:  netip.MustParsePrefix("fe80::/10"),
			expect: false,
		},
		"not_private_multicast_8": {
			input:  netip.MustParsePrefix("ff00::/8"),
			expect: false,
		},
		"not_private_doc_db8_32": {
			input:  netip.MustParsePrefix("2001:db8::/32"),
			expect: false,
		},
		"ipv4_mapped_private_10": {
			input:  netip.MustParsePrefix("::ffff:0a00:0001/120"),
			expect: true,
		},
		"ipv4_mapped_private_cgn": {
			input:  netip.MustParsePrefix("::ffff:6440:0001/120"),
			expect: true,
		},
		"ipv4_mapped_public": {
			input:  netip.MustParsePrefix("::ffff:0808:0808/104"),
			expect: false,
		},
		"adjacent_to_ula_before": {
			input:  netip.MustParsePrefix("fbff::/16"),
			expect: false,
		},
		"adjacent_to_ula_after": {
			input:  netip.MustParsePrefix("fe00::/16"),
			expect: false,
		},
		"adjacent_to_benchmarking_before": {
			input:  netip.MustParsePrefix("2001:1::/48"),
			expect: false,
		},
		"adjacent_to_benchmarking_after": {
			input:  netip.MustParsePrefix("2001:3::/48"),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := prefixV6IsPrivate(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

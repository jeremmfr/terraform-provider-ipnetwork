package provider

import (
	"net/netip"
	"testing"
)

func TestAddressV4IsPublic(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  netip.Addr
		expect bool
	}

	tests := map[string]testCase{
		"public_google_dns": {
			input:  netip.MustParseAddr("8.8.8.8"),
			expect: true,
		},
		"public_cloudflare_dns": {
			input:  netip.MustParseAddr("1.1.1.1"),
			expect: true,
		},
		"private_10": {
			input:  netip.MustParseAddr("10.0.0.1"),
			expect: false,
		},
		"private_172_16": {
			input:  netip.MustParseAddr("172.16.0.1"),
			expect: false,
		},
		"private_192_168": {
			input:  netip.MustParseAddr("192.168.1.1"),
			expect: false,
		},
		"cgn": {
			input:  netip.MustParseAddr("100.64.0.1"),
			expect: false,
		},
		"this_network": {
			input:  netip.MustParseAddr("0.1.2.3"),
			expect: false,
		},
		"loopback": {
			input:  netip.MustParseAddr("127.0.0.1"),
			expect: false,
		},
		"link_local": {
			input:  netip.MustParseAddr("169.254.1.1"),
			expect: false,
		},
		"ietf_protocol": {
			input:  netip.MustParseAddr("192.0.0.1"),
			expect: false,
		},
		"testnet1": {
			input:  netip.MustParseAddr("192.0.2.1"),
			expect: false,
		},
		"testnet2": {
			input:  netip.MustParseAddr("198.51.100.1"),
			expect: false,
		},
		"testnet3": {
			input:  netip.MustParseAddr("203.0.113.1"),
			expect: false,
		},
		"multicast": {
			input:  netip.MustParseAddr("224.0.0.1"),
			expect: false,
		},
		"reserved_240": {
			input:  netip.MustParseAddr("240.0.0.1"),
			expect: false,
		},
		"benchmarking": {
			input:  netip.MustParseAddr("198.18.1.2"),
			expect: false,
		},
		"benchmarking_end": {
			input:  netip.MustParseAddr("198.19.255.255"),
			expect: false,
		},
		"broadcast": {
			input:  netip.MustParseAddr("255.255.255.255"),
			expect: false,
		},
		"adjacent_to_10_after": {
			input:  netip.MustParseAddr("11.0.0.0"),
			expect: true,
		},
		"adjacent_to_172_after": {
			input:  netip.MustParseAddr("172.32.0.0"),
			expect: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := addressV4IsPublic(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

func TestAddressV6IsPublic(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  netip.Addr
		expect bool
	}

	tests := map[string]testCase{
		"public_ipv6_google": {
			input:  netip.MustParseAddr("2001:4860:4860::8888"),
			expect: true,
		},
		"public_ipv6_cloudflare": {
			input:  netip.MustParseAddr("2606:4700:4700::1111"),
			expect: true,
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
		"ipv6_doc_db8": {
			input:  netip.MustParseAddr("2001:db8::1"),
			expect: false,
		},
		"ipv6_doc_3fff": {
			input:  netip.MustParseAddr("3fff::1"),
			expect: false,
		},
		"ipv6_doc_3fff_end": {
			input:  netip.MustParseAddr("3fff:0fff:ffff:ffff:ffff:ffff:ffff:ffff"),
			expect: false,
		},
		"ipv6_discard_prefix": {
			input:  netip.MustParseAddr("100::"),
			expect: false,
		},
		"ipv6_discard_prefix_2": {
			input:  netip.MustParseAddr("100::1:2:3:4"),
			expect: false,
		},
		"ipv6_benchmarking": {
			input:  netip.MustParseAddr("2001:2::1"),
			expect: false,
		},
		"ipv6_benchmarking_end": {
			input:  netip.MustParseAddr("2001:2:0:ffff:ffff:ffff:ffff:ffff"),
			expect: false,
		},
		"ipv4_mapped_public": {
			input:  netip.MustParseAddr("::ffff:8.8.8.8"),
			expect: true,
		},
		"ipv4_mapped_private": {
			input:  netip.MustParseAddr("::ffff:192.168.1.1"),
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := addressV6IsPublic(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

func TestPrefixV4IsPublic(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  netip.Prefix
		expect bool
	}

	tests := map[string]testCase{
		"public_24": {
			input:  netip.MustParsePrefix("1.1.1.0/24"),
			expect: true,
		},
		"public_20": {
			input:  netip.MustParsePrefix("8.8.0.0/20"),
			expect: true,
		},
		"public_16": {
			input:  netip.MustParsePrefix("8.8.0.0/16"),
			expect: true,
		},
		"public_8": {
			input:  netip.MustParsePrefix("8.0.0.0/8"),
			expect: true,
		},
		"private_10_8": {
			input:  netip.MustParsePrefix("10.0.0.0/8"),
			expect: false,
		},
		"private_10_16": {
			input:  netip.MustParsePrefix("10.0.0.0/16"),
			expect: false,
		},
		"private_172_16_12": {
			input:  netip.MustParsePrefix("172.16.0.0/12"),
			expect: false,
		},
		"private_172_16_16": {
			input:  netip.MustParsePrefix("172.16.0.0/16"),
			expect: false,
		},
		"private_192_168_16": {
			input:  netip.MustParsePrefix("192.168.0.0/16"),
			expect: false,
		},
		"private_192_168_24": {
			input:  netip.MustParsePrefix("192.168.1.0/24"),
			expect: false,
		},
		"contains__private_10_8": {
			input:  netip.MustParsePrefix("8.0.0.0/6"),
			expect: false,
		},
		"contains_private_172_16_12": {
			input:  netip.MustParsePrefix("172.0.0.0/9"),
			expect: false,
		},
		"contains_private_192_168_16": {
			input:  netip.MustParsePrefix("192.160.0.0/12"),
			expect: false,
		},
		"cgn_10": {
			input:  netip.MustParsePrefix("100.64.0.0/10"),
			expect: false,
		},
		"cgn_16": {
			input:  netip.MustParsePrefix("100.64.0.0/16"),
			expect: false,
		},
		"this_network_8": {
			input:  netip.MustParsePrefix("0.0.0.0/8"),
			expect: false,
		},
		"this_network_24": {
			input:  netip.MustParsePrefix("0.1.2.0/24"),
			expect: false,
		},
		"loopback_8": {
			input:  netip.MustParsePrefix("127.0.0.0/8"),
			expect: false,
		},
		"loopback_16": {
			input:  netip.MustParsePrefix("127.0.0.0/16"),
			expect: false,
		},
		"link_local_16": {
			input:  netip.MustParsePrefix("169.254.0.0/16"),
			expect: false,
		},
		"link_local_24": {
			input:  netip.MustParsePrefix("169.254.1.0/24"),
			expect: false,
		},
		"ietf_protocol_24": {
			input:  netip.MustParsePrefix("192.0.0.0/24"),
			expect: false,
		},
		"testnet1_24": {
			input:  netip.MustParsePrefix("192.0.2.0/24"),
			expect: false,
		},
		"testnet2_24": {
			input:  netip.MustParsePrefix("198.51.100.0/24"),
			expect: false,
		},
		"testnet3_24": {
			input:  netip.MustParsePrefix("203.0.113.0/24"),
			expect: false,
		},
		"contains_testnet1_23": {
			input:  netip.MustParsePrefix("192.0.2.0/23"),
			expect: false,
		},
		"contains_testnet1_22": {
			input:  netip.MustParsePrefix("192.0.0.0/22"),
			expect: false,
		},
		"contains_testnet1_20": {
			input:  netip.MustParsePrefix("192.0.0.0/20"),
			expect: false,
		},
		"contains_ietf_and_testnet1_20": {
			input:  netip.MustParsePrefix("192.0.0.0/20"),
			expect: false,
		},
		"benchmarking_15": {
			input:  netip.MustParsePrefix("198.18.0.0/15"),
			expect: false,
		},
		"benchmarking_16": {
			input:  netip.MustParsePrefix("198.19.0.0/16"),
			expect: false,
		},
		"multicast_4": {
			input:  netip.MustParsePrefix("224.0.0.0/4"),
			expect: false,
		},
		"multicast_8": {
			input:  netip.MustParsePrefix("224.0.0.0/8"),
			expect: false,
		},
		"reserved_240_4": {
			input:  netip.MustParsePrefix("240.0.0.0/4"),
			expect: false,
		},
		"reserved_255_8": {
			input:  netip.MustParsePrefix("255.0.0.0/8"),
			expect: false,
		},
		"adjacent_to_10_public": {
			input:  netip.MustParsePrefix("11.0.0.0/8"),
			expect: true,
		},
		"adjacent_to_172_public": {
			input:  netip.MustParsePrefix("172.32.0.0/12"),
			expect: true,
		},
		"adjacent_to_192_168_public": {
			input:  netip.MustParsePrefix("192.169.0.0/16"),
			expect: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := prefixV4IsPublic(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

func TestPrefixV6IsPublic(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  netip.Prefix
		expect bool
	}

	tests := map[string]testCase{
		"public_google_48": {
			input:  netip.MustParsePrefix("2001:4860::/48"),
			expect: true,
		},
		"public_cloudflare_32": {
			input:  netip.MustParsePrefix("2606:4700::/32"),
			expect: true,
		},
		"public_64": {
			input:  netip.MustParsePrefix("2001:4860:4860::/64"),
			expect: true,
		},
		"public_56": {
			input:  netip.MustParsePrefix("2001:4860:4860::/56"),
			expect: true,
		},
		"ula_7": {
			input:  netip.MustParsePrefix("fc00::/7"),
			expect: false,
		},
		"ula_fd_8": {
			input:  netip.MustParsePrefix("fd00::/8"),
			expect: false,
		},
		"ula_48": {
			input:  netip.MustParsePrefix("fd12:3456:789a::/48"),
			expect: false,
		},
		"ula_64": {
			input:  netip.MustParsePrefix("fd12:3456:789a:bcde::/64"),
			expect: false,
		},
		"loopback_128": {
			input:  netip.MustParsePrefix("::1/128"),
			expect: false,
		},
		"unspecified_128": {
			input:  netip.MustParsePrefix("::/128"),
			expect: false,
		},
		"unspecified_10": {
			input:  netip.MustParsePrefix("::/10"),
			expect: false,
		},
		"link_local_10": {
			input:  netip.MustParsePrefix("fe80::/10"),
			expect: false,
		},
		"link_local_64": {
			input:  netip.MustParsePrefix("fe80::/64"),
			expect: false,
		},
		"multicast_8": {
			input:  netip.MustParsePrefix("ff00::/8"),
			expect: false,
		},
		"multicast_16": {
			input:  netip.MustParsePrefix("ff00::/16"),
			expect: false,
		},
		"multicast_interface_local": {
			input:  netip.MustParsePrefix("ff01::/16"),
			expect: false,
		},
		"multicast_link_local": {
			input:  netip.MustParsePrefix("ff02::/16"),
			expect: false,
		},
		"documentation_db8_32": {
			input:  netip.MustParsePrefix("2001:db8::/32"),
			expect: false,
		},
		"documentation_db8_48": {
			input:  netip.MustParsePrefix("2001:db8:1234::/48"),
			expect: false,
		},
		"documentation_3fff_20": {
			input:  netip.MustParsePrefix("3fff::/20"),
			expect: false,
		},
		"documentation_3fff_subset_32": {
			input:  netip.MustParsePrefix("3fff:0f00::/32"),
			expect: false,
		},
		"contains_db8_16": {
			input:  netip.MustParsePrefix("2001::/16"),
			expect: false,
		},
		"contains_db8_24": {
			input:  netip.MustParsePrefix("2001:db8::/24"),
			expect: false,
		},
		"ipv4_mapped_8.0.0.0_8": {
			input:  netip.MustParsePrefix("::ffff:0808:0808/104"),
			expect: true,
		},
		"ipv4_mapped_0_0": {
			input:  netip.MustParsePrefix("::ffff:0:0/96"),
			expect: false,
		},
		"ipv4_mapped_1.1.1.0_1": {
			input:  netip.MustParsePrefix("::ffff:0101:0101/97"),
			expect: false,
		},
		"ipv4_mapped_10.0.0.0_24": {
			input:  netip.MustParsePrefix("::ffff:0a00:0001/120"),
			expect: false,
		},
		"discard_only_64": {
			input:  netip.MustParsePrefix("100::/64"),
			expect: false,
		},
		"benchmarking_48": {
			input:  netip.MustParsePrefix("2001:2::/48"),
			expect: false,
		},
		"benchmarking_64": {
			input:  netip.MustParsePrefix("2001:2:0:1234::/64"),
			expect: false,
		},
		"adjacent_to_ula_overlaps_link_local": {
			input:  netip.MustParsePrefix("fb00::/7"),
			expect: true,
		},
		"adjacent_to_db8_public": {
			input:  netip.MustParsePrefix("2001:dc8::/32"),
			expect: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := prefixV6IsPublic(test.input)
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

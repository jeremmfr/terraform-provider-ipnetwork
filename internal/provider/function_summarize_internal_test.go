package provider

import (
	"net/netip"
	"testing"
)

func TestPrefixesSummarize(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input  []netip.Prefix
		output []netip.Prefix
	}

	tests := map[string]testCase{
		"empty": {
			input:  []netip.Prefix{},
			output: nil,
		},
		"adjacent_slash_8": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/8"),
				netip.MustParsePrefix("11.0.0.0/8"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/7"),
			},
		},
		"adjacent_with_gap": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/24"),
				netip.MustParsePrefix("10.0.1.0/24"),
				netip.MustParsePrefix("10.0.3.0/24"),
				netip.MustParsePrefix("10.0.4.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/23"),
				netip.MustParsePrefix("10.0.3.0/24"),
				netip.MustParsePrefix("10.0.4.0/24"),
			},
		},
		"all_private_rfc1918": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/8"),
				netip.MustParsePrefix("172.16.0.0/12"),
				netip.MustParsePrefix("192.168.0.0/16"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/8"),
				netip.MustParsePrefix("172.16.0.0/12"),
				netip.MustParsePrefix("192.168.0.0/16"),
			},
		},
		"completely_overlapping_hierarchy": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/16"),
				netip.MustParsePrefix("10.0.0.0/20"),
				netip.MustParsePrefix("10.0.0.0/24"),
				netip.MustParsePrefix("10.0.0.0/28"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/16"),
			},
		},
		"ipv4_mapped_ipv6": {
			input: []netip.Prefix{
				netip.MustParsePrefix("::ffff:192.0.2.0/120"),
				netip.MustParsePrefix("::ffff:192.0.3.0/120"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("::ffff:192.0.2.0/119"),
			},
		},
		"max_prefix_ipv4": {
			input: []netip.Prefix{
				netip.MustParsePrefix("0.0.0.0/0"),
				netip.MustParsePrefix("192.0.2.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("0.0.0.0/0"),
			},
		},
		"max_prefix_ipv6": {
			input: []netip.Prefix{
				netip.MustParsePrefix("::/0"),
				netip.MustParsePrefix("2001:db8::/32"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("::/0"),
			},
		},
		"mix_slash_31_and_32": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/31"),
				netip.MustParsePrefix("10.0.0.2/32"),
				netip.MustParsePrefix("10.0.0.3/32"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/30"),
			},
		},
		"mixed_ipv4_ipv6_extensive": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/24"),
				netip.MustParsePrefix("10.0.1.0/24"),
				netip.MustParsePrefix("2001:db8::/64"),
				netip.MustParsePrefix("2001:db8:0:1::/64"),
				netip.MustParsePrefix("192.168.0.0/24"),
				netip.MustParsePrefix("192.168.1.0/24"),
				netip.MustParsePrefix("fd00::/64"),
				netip.MustParsePrefix("fd00:0:0:1::/64"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/23"),
				netip.MustParsePrefix("192.168.0.0/23"),
				netip.MustParsePrefix("2001:db8::/63"),
				netip.MustParsePrefix("fd00::/63"),
			},
		},
		"multiple_groups_adjacent": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/24"),
				netip.MustParsePrefix("10.0.1.0/24"),
				netip.MustParsePrefix("10.0.4.0/24"),
				netip.MustParsePrefix("10.0.5.0/24"),
				netip.MustParsePrefix("10.0.8.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/23"),
				netip.MustParsePrefix("10.0.4.0/23"),
				netip.MustParsePrefix("10.0.8.0/24"),
			},
		},
		"consecutive_pairs": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/32"),
				netip.MustParsePrefix("192.0.2.1/32"),
				netip.MustParsePrefix("192.0.2.2/32"),
				netip.MustParsePrefix("192.0.2.3/32"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/30"),
			},
		},
		"non_consecutive_ips": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/32"),
				netip.MustParsePrefix("192.0.2.2/32"),
				netip.MustParsePrefix("192.0.2.4/32"),
				netip.MustParsePrefix("192.0.2.6/32"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/32"),
				netip.MustParsePrefix("192.0.2.2/32"),
				netip.MustParsePrefix("192.0.2.4/32"),
				netip.MustParsePrefix("192.0.2.6/32"),
			},
		},
		"reverse_order_input": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.0.5.0/24"),
				netip.MustParsePrefix("192.0.4.0/24"),
				netip.MustParsePrefix("192.0.3.0/24"),
				netip.MustParsePrefix("192.0.2.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/23"),
				netip.MustParsePrefix("192.0.4.0/23"),
			},
		},
		"single_prefix": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/24"),
			},
		},
		"subnet_of_24_into_23": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/25"),
				netip.MustParsePrefix("192.168.0.128/25"),
				netip.MustParsePrefix("192.168.1.0/25"),
				netip.MustParsePrefix("192.168.1.128/25"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.168.0.0/23"),
			},
		},
		"three_way_merge": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/26"),
				netip.MustParsePrefix("10.0.0.64/26"),
				netip.MustParsePrefix("10.0.0.128/26"),
				netip.MustParsePrefix("10.0.0.192/26"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/24"),
			},
		},
		"two_adjacent_ipv4": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/24"),
				netip.MustParsePrefix("192.0.3.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/23"),
			},
		},
		"four_adjacent_ipv4": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/24"),
				netip.MustParsePrefix("10.0.1.0/24"),
				netip.MustParsePrefix("10.0.2.0/24"),
				netip.MustParsePrefix("10.0.3.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/22"),
			},
		},
		"non_adjacent_ipv4": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/24"),
				netip.MustParsePrefix("192.0.4.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/24"),
				netip.MustParsePrefix("192.0.4.0/24"),
			},
		},
		"overlapping_ipv4": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/24"),
				netip.MustParsePrefix("192.0.2.0/25"),
				netip.MustParsePrefix("192.0.2.128/25"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/24"),
			},
		},
		"duplicate_ipv4": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/24"),
				netip.MustParsePrefix("192.0.2.0/24"),
				netip.MustParsePrefix("192.0.2.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/24"),
			},
		},
		"eight_adjacent_ipv4": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/24"),
				netip.MustParsePrefix("10.0.1.0/24"),
				netip.MustParsePrefix("10.0.2.0/24"),
				netip.MustParsePrefix("10.0.3.0/24"),
				netip.MustParsePrefix("10.0.4.0/24"),
				netip.MustParsePrefix("10.0.5.0/24"),
				netip.MustParsePrefix("10.0.6.0/24"),
				netip.MustParsePrefix("10.0.7.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/21"),
			},
		},
		"host_addresses_ipv4": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.1/32"),
				netip.MustParsePrefix("192.0.2.2/32"),
				netip.MustParsePrefix("192.0.2.3/32"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.1/32"),
				netip.MustParsePrefix("192.0.2.2/31"),
			},
		},
		"two_adjacent_ipv6": {
			input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/64"),
				netip.MustParsePrefix("2001:db8:0:1::/64"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/63"),
			},
		},
		"four_adjacent_ipv6": {
			input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/64"),
				netip.MustParsePrefix("2001:db8:0:1::/64"),
				netip.MustParsePrefix("2001:db8:0:2::/64"),
				netip.MustParsePrefix("2001:db8:0:3::/64"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/62"),
			},
		},
		"non_adjacent_ipv6": {
			input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/64"),
				netip.MustParsePrefix("2001:db8:0:2::/64"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/64"),
				netip.MustParsePrefix("2001:db8:0:2::/64"),
			},
		},
		"overlapping_ipv6": {
			input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
				netip.MustParsePrefix("2001:db8::/64"),
				netip.MustParsePrefix("2001:db8:0:1::/64"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/48"),
			},
		},
		"host_addresses_ipv6": {
			input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::1/128"),
				netip.MustParsePrefix("2001:db8::2/128"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::1/128"),
				netip.MustParsePrefix("2001:db8::2/128"),
			},
		},
		"ipv6_link_local": {
			input: []netip.Prefix{
				netip.MustParsePrefix("fe80::/64"),
				netip.MustParsePrefix("fe80:0:0:1::/64"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("fe80::/63"),
			},
		},
		"ipv6_unique_local": {
			input: []netip.Prefix{
				netip.MustParsePrefix("fd00::/64"),
				netip.MustParsePrefix("fd00:0:0:1::/64"),
				netip.MustParsePrefix("fd00:0:0:2::/64"),
				netip.MustParsePrefix("fd00:0:0:3::/64"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("fd00::/62"),
			},
		},
		"eight_adjacent_ipv6": {
			input: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/64"),
				netip.MustParsePrefix("2001:db8:0:1::/64"),
				netip.MustParsePrefix("2001:db8:0:2::/64"),
				netip.MustParsePrefix("2001:db8:0:3::/64"),
				netip.MustParsePrefix("2001:db8:0:4::/64"),
				netip.MustParsePrefix("2001:db8:0:5::/64"),
				netip.MustParsePrefix("2001:db8:0:6::/64"),
				netip.MustParsePrefix("2001:db8:0:7::/64"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("2001:db8::/61"),
			},
		},
		"mixed_families": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/24"),
				netip.MustParsePrefix("192.0.3.0/24"),
				netip.MustParsePrefix("2001:db8::/64"),
				netip.MustParsePrefix("2001:db8:0:1::/64"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/23"),
				netip.MustParsePrefix("2001:db8::/63"),
			},
		},
		"mixed_sizes_ipv4": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/24"),
				netip.MustParsePrefix("10.0.1.0/24"),
				netip.MustParsePrefix("10.0.4.0/22"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/23"),
				netip.MustParsePrefix("10.0.4.0/22"),
			},
		},
		"unsorted_input": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.0.4.0/24"),
				netip.MustParsePrefix("192.0.2.0/24"),
				netip.MustParsePrefix("192.0.3.0/24"),
				netip.MustParsePrefix("192.0.5.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.0.2.0/23"),
				netip.MustParsePrefix("192.0.4.0/23"),
			},
		},
		"complex_merge": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/25"),
				netip.MustParsePrefix("10.0.0.128/25"),
				netip.MustParsePrefix("10.0.1.0/25"),
				netip.MustParsePrefix("10.0.1.128/25"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/23"),
			},
		},
		"partial_overlap": {
			input: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/23"),
				netip.MustParsePrefix("10.0.1.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("10.0.0.0/23"),
			},
		},
		"large_block_with_small": {
			input: []netip.Prefix{
				netip.MustParsePrefix("192.0.0.0/16"),
				netip.MustParsePrefix("192.0.2.0/24"),
			},
			output: []netip.Prefix{
				netip.MustParsePrefix("192.0.0.0/16"),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := prefixesSummarize(test.input)

			// Check length
			if len(result) != len(test.output) {
				t.Errorf("got unexpected number of prefixes: want %d, got %d", len(test.output), len(result))
				t.Errorf("expected: %v", test.output)
				t.Errorf("got: %v", result)

				return
			}

			// Check each prefix
			for i, expected := range test.output {
				if result[i] != expected {
					t.Errorf("at index %d: want %s, got %s", i, expected, result[i])
				}
			}
		})
	}
}

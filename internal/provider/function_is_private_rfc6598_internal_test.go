package provider

import (
	"net/netip"
	"testing"
)

func TestAddressIsPrivateRFC6598(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input  string
		expect bool
	}{
		"public_ipv4_google": {
			input:  "8.8.8.8",
			expect: false,
		},
		"public_ipv4_cloudflare": {
			input:  "1.1.1.1",
			expect: false,
		},
		"cgn_rfc6598_start": {
			input:  "100.64.0.0",
			expect: true,
		},
		"cgn_rfc6598_mid": {
			input:  "100.100.50.25",
			expect: true,
		},
		"cgn_rfc6598_end": {
			input:  "100.127.255.255",
			expect: true,
		},
		"not_cgn_adjacent_before": {
			input:  "100.63.255.255",
			expect: false,
		},
		"not_cgn_adjacent_after": {
			input:  "100.128.0.0",
			expect: false,
		},
		"not_cgn_private_10": {
			input:  "10.0.0.1",
			expect: false,
		},
		"not_cgn_private_172": {
			input:  "172.16.0.1",
			expect: false,
		},
		"not_cgn_private_192": {
			input:  "192.168.1.1",
			expect: false,
		},
		"not_cgn_loopback": {
			input:  "127.0.0.1",
			expect: false,
		},
		"not_cgn_link_local": {
			input:  "169.254.1.1",
			expect: false,
		},
		"public_ipv6_google": {
			input:  "2001:4860:4860::8888",
			expect: false,
		},
		"ipv6_ula": {
			input:  "fd00::1",
			expect: false,
		},
		"ipv4_mapped_cgn": {
			input:  "::ffff:100.64.0.1",
			expect: true,
		},
		"ipv4_mapped_cgn_mid": {
			input:  "::ffff:100.100.50.25",
			expect: true,
		},
		"ipv4_mapped_public": {
			input:  "::ffff:8.8.8.8",
			expect: false,
		},
		"ipv4_mapped_private": {
			input:  "::ffff:192.168.1.1",
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := addressIsPrivateRFC6598(netip.MustParseAddr(test.input))
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

func TestPrefixIsPrivateRFC6598(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input  string
		expect bool
	}{
		"public_ipv4_8": {
			input:  "8.0.0.0/8",
			expect: false,
		},
		"public_ipv4_24": {
			input:  "1.1.1.0/24",
			expect: false,
		},
		"cgn_rfc6598_10": {
			input:  "100.64.0.0/10",
			expect: true,
		},
		"cgn_rfc6598_16": {
			input:  "100.64.0.0/16",
			expect: true,
		},
		"cgn_rfc6598_24": {
			input:  "100.100.50.0/24",
			expect: true,
		},
		"cgn_rfc6598_32": {
			input:  "100.127.255.255/32",
			expect: true,
		},
		"not_cgn_9": {
			input:  "100.64.0.0/9",
			expect: false,
		},
		"not_cgn_adjacent_before": {
			input:  "100.0.0.0/10",
			expect: false,
		},
		"not_cgn_adjacent_after": {
			input:  "100.128.0.0/10",
			expect: false,
		},
		"not_cgn_private_10": {
			input:  "10.0.0.0/8",
			expect: false,
		},
		"not_cgn_private_172": {
			input:  "172.16.0.0/12",
			expect: false,
		},
		"not_cgn_private_192": {
			input:  "192.168.0.0/16",
			expect: false,
		},
		"not_cgn_loopback": {
			input:  "127.0.0.0/8",
			expect: false,
		},
		"not_cgn_link_local": {
			input:  "169.254.0.0/16",
			expect: false,
		},
		"public_ipv6": {
			input:  "2001:4860::/48",
			expect: false,
		},
		"ipv6_ula": {
			input:  "fd00::/7",
			expect: false,
		},
		"ipv4_mapped_cgn": {
			input:  "::ffff:100.64.0.0/106",
			expect: true,
		},
		"ipv4_mapped_cgn_16": {
			input:  "::ffff:100.100.0.0/112",
			expect: true,
		},
		"ipv4_mapped_public": {
			input:  "::ffff:8.8.8.8/128",
			expect: false,
		},
		"ipv4_mapped_private": {
			input:  "::ffff:192.168.0.0/112",
			expect: false,
		},
		"ipv4_mapped_too_short": {
			input:  "::ffff:100.64.0.0/95",
			expect: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := prefixIsPrivateRFC6598(netip.MustParsePrefix(test.input))
			if resp != test.expect {
				t.Errorf("got unexpected resp: want %t, got %t", test.expect, resp)
			}
		})
	}
}

package provider

import (
	"net/netip"
	"testing"
)

func TestIpAddrToMaskBits(t *testing.T) {
	t.Parallel()

	type testCase struct {
		mask       netip.Addr
		expectOk   bool
		expectBits int
	}

	tests := map[string]testCase{
		"0": {
			mask:       netip.MustParseAddr("0.0.0.0"),
			expectOk:   true,
			expectBits: 0,
		},
		"8": {
			mask:       netip.MustParseAddr("255.0.0.0"),
			expectOk:   true,
			expectBits: 8,
		},
		"~7": {
			mask:       netip.MustParseAddr("253.0.0.0"),
			expectOk:   false,
			expectBits: 0,
		},
		"~9": {
			mask:       netip.MustParseAddr("255.64.0.0"),
			expectOk:   false,
			expectBits: 0,
		},
		"15": {
			mask:       netip.MustParseAddr("255.254.0.0"),
			expectOk:   true,
			expectBits: 15,
		},
		"~15": {
			mask:       netip.MustParseAddr("254.254.0.0"),
			expectOk:   false,
			expectBits: 0,
		},
		"16": {
			mask:       netip.MustParseAddr("255.255.0.0"),
			expectOk:   true,
			expectBits: 16,
		},
		"~16": {
			mask:       netip.MustParseAddr("254.255.0.0"),
			expectOk:   false,
			expectBits: 0,
		},
		"~23": {
			mask:       netip.MustParseAddr("255.255.200.0"),
			expectOk:   false,
			expectBits: 0,
		},
		"24": {
			mask:       netip.MustParseAddr("255.255.255.0"),
			expectOk:   true,
			expectBits: 24,
		},
		"~24": {
			mask:       netip.MustParseAddr("255.0.255.0"),
			expectOk:   false,
			expectBits: 0,
		},
		"28": {
			mask:       netip.MustParseAddr("255.255.255.240"),
			expectOk:   true,
			expectBits: 28,
		},
		"~29": {
			mask:       netip.MustParseAddr("255.255.128.248"),
			expectOk:   false,
			expectBits: 0,
		},
		"29": {
			mask:       netip.MustParseAddr("255.255.255.248"),
			expectOk:   true,
			expectBits: 29,
		},
		"~31": {
			mask:       netip.MustParseAddr("255.255.255.251"),
			expectOk:   false,
			expectBits: 0,
		},
		"32": {
			mask:       netip.MustParseAddr("255.255.255.255"),
			expectOk:   true,
			expectBits: 32,
		},
		"~32": {
			mask:       netip.MustParseAddr("255.128.255.255"),
			expectOk:   false,
			expectBits: 0,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			bits, ok := ipAddrToMaskBits(test.mask)

			if test.expectBits != bits {
				t.Errorf("got unexpected bits: want %d, got %d", test.expectBits, bits)
			}
			if test.expectOk != ok {
				t.Errorf("got unexpected ok: want %t, got %t", test.expectOk, ok)
			}
		})
	}
}

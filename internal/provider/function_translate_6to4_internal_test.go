package provider

import (
	"net/netip"
	"testing"
)

func TestTranslateAddress6to4(t *testing.T) {
	t.Parallel()

	type testCase struct {
		address    netip.Prefix
		expectAddr netip.Addr
	}

	tests := map[string]testCase{
		"min_31": {
			address:    netip.MustParsePrefix("3fff:0:fffe:fdfc:0:0:0:0/31"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_31": {
			address:    netip.MustParsePrefix("3fff:fffe:aabb:ccdd:0:0:0:0/31"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_32": {
			address:    netip.MustParsePrefix("3fff:0:fffe:fdfc:0:0:0:0/32"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_32": {
			address:    netip.MustParsePrefix("3fff:ffff:aabb:ccdd:0:0:0:0/32"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_33": {
			address:    netip.MustParsePrefix("3fff:0:ff:fefd:fc:0:0:0/33"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_33": {
			address:    netip.MustParsePrefix("3fff:ffff:80aa:bbcc:dd:0:0:0/33"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_39": {
			address:    netip.MustParsePrefix("3fff:0:ff:fefd:fc:0:0:0/39"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_39": {
			address:    netip.MustParsePrefix("3fff:ffff:feaa:bbcc:dd:0:0:0/39"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_40": {
			address:    netip.MustParsePrefix("3fff:0:ff:fefd:fc:0:0:0/39"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_40": {
			address:    netip.MustParsePrefix("3fff:ffff:feaa:bbcc:dd:0:0:0/40"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_41": {
			address:    netip.MustParsePrefix("3fff:0:0:fffe:fd:fc00:0:0/41"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_41": {
			address:    netip.MustParsePrefix("3fff:ffff:ff80:aabb:cc:dd00:0:0/41"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_47": {
			address:    netip.MustParsePrefix("3fff:0:0:fffe:fd:fc00:0:0/47"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_47": {
			address:    netip.MustParsePrefix("3fff:ffff:fffe:aabb:cc:dd00:0:0/47"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_48": {
			address:    netip.MustParsePrefix("3fff:0:0:fffe:fd:fc00:0:0/48"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_48": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:aabb:cc:dd00:0:0/48"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_49": {
			address:    netip.MustParsePrefix("3fff:0:0:ff:fe:fdfc:0:0/49"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_49": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:80aa:bb:ccdd:0:0/49"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_55": {
			address:    netip.MustParsePrefix("3fff:0:0:ff:fe:fdfc:0:0/55"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_55": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:feaa:bb:ccdd:0:0/55"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_56": {
			address:    netip.MustParsePrefix("3fff:0:0:ff:fe:fdfc:0:0/56"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_56": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:ffaa:bb:ccdd:0:0/56"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_57": {
			address:    netip.MustParsePrefix("3fff:0:0:0:ff:fefd:fc00:0/57"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_57": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:ff80:aa:bbcc:dd00:0/57"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_63": {
			address:    netip.MustParsePrefix("3fff:0:0:0:ff:fefd:fc00:0/63"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_63": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:fffe:aa:bbcc:dd00:0/63"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_64": {
			address:    netip.MustParsePrefix("3fff:0:0:0:ff:fefd:fc00:0/64"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_64": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:ffff:aa:bbcc:dd00:0/64"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_65": {
			address:    netip.MustParsePrefix("3fff:0:0:0:0:0:fffe:fdfc/65"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_65": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:ffff:8000:0:aabb:ccdd/65"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_96": {
			address:    netip.MustParsePrefix("3fff:0:0:0:0:0:fffe:fdfc/96"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_96": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:aabb:ccdd/96"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_128": {
			address:    netip.MustParsePrefix("3fff:0:0:0:0:0:fffe:fdfc/128"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_128": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:aabb:ccdd/128"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"rfc_2.4_example_32": {
			address:    netip.MustParsePrefix("2001:db8:c000:221::/32"),
			expectAddr: netip.MustParseAddr("192.0.2.33"),
		},
		"rfc_2.4_example_40": {
			address:    netip.MustParsePrefix("2001:db8:1c0:2:21::/40"),
			expectAddr: netip.MustParseAddr("192.0.2.33"),
		},
		"rfc_2.4_example_48": {
			address:    netip.MustParsePrefix("2001:db8:122:c000:2:2100::/48"),
			expectAddr: netip.MustParseAddr("192.0.2.33"),
		},
		"rfc_2.4_example_56": {
			address:    netip.MustParsePrefix("2001:db8:122:3c0:0:221::/56"),
			expectAddr: netip.MustParseAddr("192.0.2.33"),
		},
		"rfc_2.4_example_64": {
			address:    netip.MustParsePrefix("2001:db8:122:344:c0:2:2100::/64"),
			expectAddr: netip.MustParseAddr("192.0.2.33"),
		},
		"rfc_2.4_example_96": {
			address:    netip.MustParsePrefix("2001:db8:122:344::192.0.2.33/96"),
			expectAddr: netip.MustParseAddr("192.0.2.33"),
		},
		"rfc_2.4_example_well_known": {
			address:    netip.MustParsePrefix("64:ff9b::192.0.2.33/96"),
			expectAddr: netip.MustParseAddr("192.0.2.33"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resp := translateAddress6to4(test.address)
			if resp != test.expectAddr {
				t.Errorf("got unexpected resp: want %q, got %q", test.expectAddr, resp)
			}
		})
	}
}

func BenchmarkTranslateAddress6to4_32(b *testing.B) {
	benchmarkTranslateAddress6to4(b,
		netip.MustParsePrefix("2001:db8:c000:221::/32"),
	)
}

func BenchmarkTranslateAddress6to4_40(b *testing.B) {
	benchmarkTranslateAddress6to4(b,
		netip.MustParsePrefix("2001:db8:1c0:2:21::/40"),
	)
}

func BenchmarkTranslateAddress6to4_48(b *testing.B) {
	benchmarkTranslateAddress6to4(b,
		netip.MustParsePrefix("2001:db8:122:c000:2:2100::/48"),
	)
}

func BenchmarkTranslateAddress6to4_56(b *testing.B) {
	benchmarkTranslateAddress6to4(b,
		netip.MustParsePrefix("2001:db8:122:3c0:0:221::/56"),
	)
}

func BenchmarkTranslateAddress6to4_64(b *testing.B) {
	benchmarkTranslateAddress6to4(b,
		netip.MustParsePrefix("2001:db8:122:344:c0:2:2100::/64"),
	)
}

func BenchmarkTranslateAddress6to4_96(b *testing.B) {
	benchmarkTranslateAddress6to4(b,
		netip.MustParsePrefix("2001:db8:122:344::192.0.2.33/96"),
	)
}

func benchmarkTranslateAddress6to4(b *testing.B,
	address netip.Prefix,
) {
	b.Helper()

	for b.Loop() {
		translateAddress6to4(address)
	}
}

package provider

import (
	"net/netip"
	"testing"
)

func TestPtrNameFromIP(t *testing.T) {
	t.Parallel()

	type testCase struct {
		ip            netip.Addr
		expectPtrName string
	}

	tests := map[string]testCase{
		"null": {
			expectPtrName: "",
		},
		"ipv4_1": {
			ip:            netip.MustParseAddr("192.0.200.255"),
			expectPtrName: "255.200.0.192.in-addr.arpa.",
		},
		"ipv6_1": {
			ip:            netip.MustParseAddr("2001:db8:0:a9ba:a::1"),
			expectPtrName: "1.0.0.0.0.0.0.0.0.0.0.0.a.0.0.0.a.b.9.a.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa.",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resp := ptrNameFromIP(test.ip)
			if resp != test.expectPtrName {
				t.Errorf("got unexpected resp: want %q, got %q", test.expectPtrName, resp)
			}
		})
	}
}

func BenchmarkPtrNameFromIP_ipv4(b *testing.B) {
	benchmarkPtrNameFromIP(b, netip.MustParseAddr("255.255.255.255"))
}

func BenchmarkPtrNameFromIP_ipv6(b *testing.B) {
	benchmarkPtrNameFromIP(b, netip.MustParseAddr("2001:db8:0:a9ba:a::1"))
}

func benchmarkPtrNameFromIP(b *testing.B, ip netip.Addr) {
	b.Helper()

	for b.Loop() {
		ptrNameFromIP(ip)
	}
}

func TestTranslateAddress4to6(t *testing.T) {
	t.Parallel()

	type testCase struct {
		address    netip.Addr
		prefix     netip.Prefix
		expectAddr netip.Addr
	}

	tests := map[string]testCase{
		"min_/31": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/31"),
			expectAddr: netip.MustParseAddr("3fff:0:fffe:fdfc:0:0:0:0"),
		},
		"max_/31": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/31"),
			expectAddr: netip.MustParseAddr("3fff:fffe:aabb:ccdd:0:0:0:0"),
		},
		"min_/32": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/32"),
			expectAddr: netip.MustParseAddr("3fff:0:fffe:fdfc:0:0:0:0"),
		},
		"max_/32": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/32"),
			expectAddr: netip.MustParseAddr("3fff:ffff:aabb:ccdd:0:0:0:0"),
		},
		"min_/33": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/33"),
			expectAddr: netip.MustParseAddr("3fff:0:ff:fefd:fc:0:0:0"),
		},
		"max_/33": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/33"),
			expectAddr: netip.MustParseAddr("3fff:ffff:80aa:bbcc:dd:0:0:0"),
		},
		"min_/39": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/39"),
			expectAddr: netip.MustParseAddr("3fff:0:ff:fefd:fc:0:0:0"),
		},
		"max_/39": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/39"),
			expectAddr: netip.MustParseAddr("3fff:ffff:feaa:bbcc:dd:0:0:0"),
		},
		"min_/40": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/40"),
			expectAddr: netip.MustParseAddr("3fff:0:ff:fefd:fc:0:0:0"),
		},
		"max_/40": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/40"),
			expectAddr: netip.MustParseAddr("3fff:ffff:ffaa:bbcc:dd:0:0:0"),
		},
		"min_/41": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/41"),
			expectAddr: netip.MustParseAddr("3fff:0:0:fffe:fd:fc00:0:0"),
		},
		"max_/41": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/41"),
			expectAddr: netip.MustParseAddr("3fff:ffff:ff80:aabb:cc:dd00:0:0"),
		},
		"min_/47": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/47"),
			expectAddr: netip.MustParseAddr("3fff:0:0:fffe:fd:fc00:0:0"),
		},
		"max_/47": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/47"),
			expectAddr: netip.MustParseAddr("3fff:ffff:fffe:aabb:cc:dd00:0:0"),
		},
		"min_/48": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/48"),
			expectAddr: netip.MustParseAddr("3fff:0:0:fffe:fd:fc00:0:0"),
		},
		"max_/48": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/48"),
			expectAddr: netip.MustParseAddr("3fff:ffff:ffff:aabb:cc:dd00:0:0"),
		},
		"min_/49": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/49"),
			expectAddr: netip.MustParseAddr("3fff:0:0:ff:fe:fdfc:0:0"),
		},
		"max_/49": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/49"),
			expectAddr: netip.MustParseAddr("3fff:ffff:ffff:80aa:bb:ccdd:0:0"),
		},
		"min_/55": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/55"),
			expectAddr: netip.MustParseAddr("3fff:0:0:ff:fe:fdfc:0:0"),
		},
		"max_/55": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/55"),
			expectAddr: netip.MustParseAddr("3fff:ffff:ffff:feaa:bb:ccdd:0:0"),
		},
		"min_/56": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/56"),
			expectAddr: netip.MustParseAddr("3fff:0:0:ff:fe:fdfc:0:0"),
		},
		"max_/56": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/56"),
			expectAddr: netip.MustParseAddr("3fff:ffff:ffff:ffaa:bb:ccdd:0:0"),
		},
		"min_/57": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/57"),
			expectAddr: netip.MustParseAddr("3fff:0:0:0:ff:fefd:fc00:0"),
		},
		"max_/57": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/57"),
			expectAddr: netip.MustParseAddr("3fff:ffff:ffff:ff80:aa:bbcc:dd00:0"),
		},
		"min_/63": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/63"),
			expectAddr: netip.MustParseAddr("3fff:0:0:0:ff:fefd:fc00:0"),
		},
		"max_/63": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/63"),
			expectAddr: netip.MustParseAddr("3fff:ffff:ffff:fffe:aa:bbcc:dd00:0"),
		},
		"min_/64": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/64"),
			expectAddr: netip.MustParseAddr("3fff:0:0:0:ff:fefd:fc00:0"),
		},
		"max_/64": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/64"),
			expectAddr: netip.MustParseAddr("3fff:ffff:ffff:ffff:aa:bbcc:dd00:0"),
		},
		"min_/65": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/65"),
			expectAddr: netip.MustParseAddr("3fff:0:0:0:0:0:fffe:fdfc"),
		},
		"max_/65": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/65"),
			expectAddr: netip.MustParseAddr("3fff:ffff:ffff:ffff:8000:0:aabb:ccdd"),
		},
		"min_/96": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/96"),
			expectAddr: netip.MustParseAddr("3fff:0:0:0:0:0:fffe:fdfc"),
		},
		"max_/96": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/96"),
			expectAddr: netip.MustParseAddr("3fff:ffff:ffff:ffff:ffff:ffff:aabb:ccdd"),
		},
		"min_/128": {
			address:    netip.MustParseAddr("255.254.253.252"),
			prefix:     netip.MustParsePrefix("3fff::/128"),
			expectAddr: netip.MustParseAddr("3fff:0:0:0:0:0:fffe:fdfc"),
		},
		"max_/128": {
			address:    netip.MustParseAddr("170.187.204.221"),
			prefix:     netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:ffff:ffff/128"),
			expectAddr: netip.MustParseAddr("3fff:ffff:ffff:ffff:ffff:ffff:aabb:ccdd"),
		},
		"rfc_2.4_example_32": {
			address:    netip.MustParseAddr("192.0.2.33"),
			prefix:     netip.MustParsePrefix("2001:db8::/32"),
			expectAddr: netip.MustParseAddr("2001:db8:c000:221::"),
		},
		"rfc_2.4_example_40": {
			address:    netip.MustParseAddr("192.0.2.33"),
			prefix:     netip.MustParsePrefix("2001:db8:100::/40"),
			expectAddr: netip.MustParseAddr("2001:db8:1c0:2:21::"),
		},
		"rfc_2.4_example_48": {
			address:    netip.MustParseAddr("192.0.2.33"),
			prefix:     netip.MustParsePrefix("2001:db8:122::/48"),
			expectAddr: netip.MustParseAddr("2001:db8:122:c000:2:2100::"),
		},
		"rfc_2.4_example_56": {
			address:    netip.MustParseAddr("192.0.2.33"),
			prefix:     netip.MustParsePrefix("2001:db8:122:300::/56"),
			expectAddr: netip.MustParseAddr("2001:db8:122:3c0:0:221::"),
		},
		"rfc_2.4_example_64": {
			address:    netip.MustParseAddr("192.0.2.33"),
			prefix:     netip.MustParsePrefix("2001:db8:122:344::/64"),
			expectAddr: netip.MustParseAddr("2001:db8:122:344:c0:2:2100::"),
		},
		"rfc_2.4_example_96": {
			address:    netip.MustParseAddr("192.0.2.33"),
			prefix:     netip.MustParsePrefix("2001:db8:122:344::/96"),
			expectAddr: netip.MustParseAddr("2001:db8:122:344::192.0.2.33"),
		},
		"rfc_2.4_example_well_known": {
			address:    netip.MustParseAddr("192.0.2.33"),
			prefix:     netip.MustParsePrefix("64:ff9b::/96"),
			expectAddr: netip.MustParseAddr("64:ff9b::192.0.2.33"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			resp := translateAddress4to6(test.address, test.prefix)
			if resp != test.expectAddr {
				t.Errorf("got unexpected resp: want %q, got %q", test.expectAddr, resp)
			}
		})
	}
}

func BenchmarkTranslateAddress4to6_32(b *testing.B) {
	benchmarkTranslateAddress4to6(b,
		netip.MustParseAddr("192.0.2.33"),
		netip.MustParsePrefix("2001:db8::/32"),
	)
}

func BenchmarkTranslateAddress4to6_40(b *testing.B) {
	benchmarkTranslateAddress4to6(b,
		netip.MustParseAddr("192.0.2.33"),
		netip.MustParsePrefix("2001:db8:100::/40"),
	)
}

func BenchmarkTranslateAddress4to6_48(b *testing.B) {
	benchmarkTranslateAddress4to6(b,
		netip.MustParseAddr("192.0.2.33"),
		netip.MustParsePrefix("2001:db8:122::/48"),
	)
}

func BenchmarkTranslateAddress4to6_56(b *testing.B) {
	benchmarkTranslateAddress4to6(b,
		netip.MustParseAddr("192.0.2.33"),
		netip.MustParsePrefix("2001:db8:122:300::/56"),
	)
}

func BenchmarkTranslateAddress4to6_64(b *testing.B) {
	benchmarkTranslateAddress4to6(b,
		netip.MustParseAddr("192.0.2.33"),
		netip.MustParsePrefix("2001:db8:122:344::/64"),
	)
}

func BenchmarkTranslateAddress4to6_96(b *testing.B) {
	benchmarkTranslateAddress4to6(b,
		netip.MustParseAddr("192.0.2.33"),
		netip.MustParsePrefix("2001:db8:122:344::/96"),
	)
}

func benchmarkTranslateAddress4to6(b *testing.B,
	address netip.Addr,
	prefix netip.Prefix,
) {
	b.Helper()

	for b.Loop() {
		translateAddress4to6(address, prefix)
	}
}

func TestTranslateAddress6to4(t *testing.T) {
	t.Parallel()

	type testCase struct {
		address    netip.Prefix
		expectAddr netip.Addr
	}

	tests := map[string]testCase{
		"min_/31": {
			address:    netip.MustParsePrefix("3fff:0:fffe:fdfc:0:0:0:0/31"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/31": {
			address:    netip.MustParsePrefix("3fff:fffe:aabb:ccdd:0:0:0:0/31"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/32": {
			address:    netip.MustParsePrefix("3fff:0:fffe:fdfc:0:0:0:0/32"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/32": {
			address:    netip.MustParsePrefix("3fff:ffff:aabb:ccdd:0:0:0:0/32"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/33": {
			address:    netip.MustParsePrefix("3fff:0:ff:fefd:fc:0:0:0/33"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/33": {
			address:    netip.MustParsePrefix("3fff:ffff:80aa:bbcc:dd:0:0:0/33"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/39": {
			address:    netip.MustParsePrefix("3fff:0:ff:fefd:fc:0:0:0/39"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/39": {
			address:    netip.MustParsePrefix("3fff:ffff:feaa:bbcc:dd:0:0:0/39"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/40": {
			address:    netip.MustParsePrefix("3fff:0:ff:fefd:fc:0:0:0/39"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/40": {
			address:    netip.MustParsePrefix("3fff:ffff:feaa:bbcc:dd:0:0:0/40"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/41": {
			address:    netip.MustParsePrefix("3fff:0:0:fffe:fd:fc00:0:0/41"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/41": {
			address:    netip.MustParsePrefix("3fff:ffff:ff80:aabb:cc:dd00:0:0/41"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/47": {
			address:    netip.MustParsePrefix("3fff:0:0:fffe:fd:fc00:0:0/47"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/47": {
			address:    netip.MustParsePrefix("3fff:ffff:fffe:aabb:cc:dd00:0:0/47"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/48": {
			address:    netip.MustParsePrefix("3fff:0:0:fffe:fd:fc00:0:0/48"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/48": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:aabb:cc:dd00:0:0/48"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/49": {
			address:    netip.MustParsePrefix("3fff:0:0:ff:fe:fdfc:0:0/49"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/49": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:80aa:bb:ccdd:0:0/49"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/55": {
			address:    netip.MustParsePrefix("3fff:0:0:ff:fe:fdfc:0:0/55"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/55": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:feaa:bb:ccdd:0:0/55"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/56": {
			address:    netip.MustParsePrefix("3fff:0:0:ff:fe:fdfc:0:0/56"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/56": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:ffaa:bb:ccdd:0:0/56"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/57": {
			address:    netip.MustParsePrefix("3fff:0:0:0:ff:fefd:fc00:0/57"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/57": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:ff80:aa:bbcc:dd00:0/57"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/63": {
			address:    netip.MustParsePrefix("3fff:0:0:0:ff:fefd:fc00:0/63"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/63": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:fffe:aa:bbcc:dd00:0/63"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/64": {
			address:    netip.MustParsePrefix("3fff:0:0:0:ff:fefd:fc00:0/64"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/64": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:ffff:aa:bbcc:dd00:0/64"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/65": {
			address:    netip.MustParsePrefix("3fff:0:0:0:0:0:fffe:fdfc/65"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/65": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:ffff:8000:0:aabb:ccdd/65"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/96": {
			address:    netip.MustParsePrefix("3fff:0:0:0:0:0:fffe:fdfc/96"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/96": {
			address:    netip.MustParsePrefix("3fff:ffff:ffff:ffff:ffff:ffff:aabb:ccdd/96"),
			expectAddr: netip.MustParseAddr("170.187.204.221"),
		},
		"min_/128": {
			address:    netip.MustParsePrefix("3fff:0:0:0:0:0:fffe:fdfc/128"),
			expectAddr: netip.MustParseAddr("255.254.253.252"),
		},
		"max_/128": {
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

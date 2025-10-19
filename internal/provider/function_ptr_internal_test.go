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

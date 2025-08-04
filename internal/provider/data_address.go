package provider

import (
	"net/netip"
	"slices"
	"strconv"
	"strings"
)

const hexDigits = "0123456789abcdef"

func completeIPv4Address(input string) string {
	switch strings.Count(input, ".") {
	case 3:
		return input
	case 2:
		return input + ".0"
	case 1:
		return input + ".0.0"
	case 0:
		return input + ".0.0.0"
	default:
		return input
	}
}

func ptrNameFromIP(ip netip.Addr) string {
	ipOcts := ip.AsSlice()
	switch {
	case ip.Is4():
		b := strings.Builder{}
		b.Grow(len(ipOcts)*4 + len("in-addr.arpa."))

		for _, v := range slices.Backward(ipOcts) {
			_, _ = b.WriteString(strconv.FormatUint(uint64(v), 10))
			_, _ = b.WriteRune('.')
		}
		_, _ = b.WriteString("in-addr.arpa.")

		return b.String()

	case ip.Is6():
		b := strings.Builder{}
		b.Grow(len(ipOcts)*4 + len("ip6.arpa."))

		for _, v := range slices.Backward(ipOcts) {
			_ = b.WriteByte(hexDigits[v&0xF])
			_, _ = b.WriteRune('.')
			_ = b.WriteByte(hexDigits[v>>4])
			_, _ = b.WriteRune('.')
		}
		_, _ = b.WriteString("ip6.arpa.")

		return b.String()
	}

	return ""
}

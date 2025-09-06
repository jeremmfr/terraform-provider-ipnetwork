package provider

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"net"
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

func translateAddress4to6(address netip.Addr, prefix netip.Prefix) netip.Addr {
	if !address.IsValid() || !address.Is4() {
		return netip.Addr{}
	}
	if !prefix.IsValid() || !prefix.Addr().Is6() {
		return netip.Addr{}
	}

	prefixOcts := prefix.Masked().Addr().AsSlice()

	switch bits := prefix.Bits(); {
	case bits <= 32:
		result, _ := netip.AddrFromSlice(append(append(
			prefixOcts[:4],
			address.AsSlice()...),
			[]byte{0, 0, 0, 0, 0, 0, 0, 0}...,
		))

		return result
	case bits > 32 && bits <= 40:
		result, _ := netip.AddrFromSlice(append(append(append(append(
			prefixOcts[:5],
			address.AsSlice()[:3]...),
			byte(0)),
			address.AsSlice()[3:]...),
			[]byte{0, 0, 0, 0, 0, 0}...,
		))

		return result
	case bits > 40 && bits <= 48:
		result, _ := netip.AddrFromSlice(append(append(append(append(
			prefixOcts[:6],
			address.AsSlice()[:2]...),
			byte(0)),
			address.AsSlice()[2:]...),
			[]byte{0, 0, 0, 0, 0}...,
		))

		return result
	case bits > 48 && bits <= 56:
		result, _ := netip.AddrFromSlice(append(append(append(append(
			prefixOcts[:7],
			address.AsSlice()[:1]...),
			byte(0)),
			address.AsSlice()[1:]...),
			[]byte{0, 0, 0, 0}...,
		))

		return result
	case bits > 56 && bits <= 64:
		result, _ := netip.AddrFromSlice(append(append(append(
			prefixOcts[:8],
			byte(0)),
			address.AsSlice()...),
			[]byte{0, 0, 0}...,
		))

		return result
	default:
		result, _ := netip.AddrFromSlice(append(
			prefixOcts[:12],
			address.AsSlice()...),
		)

		return result
	}
}

func translateAddress6to4(address netip.Prefix) netip.Addr {
	if !address.IsValid() || !address.Addr().Is6() {
		return netip.Addr{}
	}

	switch bits := address.Bits(); {
	case bits <= 32:
		result, _ := netip.AddrFromSlice(address.Addr().AsSlice()[4:8])

		return result
	case bits > 32 && bits <= 40:
		result, _ := netip.AddrFromSlice(append(
			address.Addr().AsSlice()[5:8],
			address.Addr().AsSlice()[9:10]...,
		))

		return result
	case bits > 40 && bits <= 48:
		result, _ := netip.AddrFromSlice(append(
			address.Addr().AsSlice()[6:8],
			address.Addr().AsSlice()[9:11]...,
		))

		return result
	case bits > 48 && bits <= 56:
		result, _ := netip.AddrFromSlice(append(
			address.Addr().AsSlice()[7:8],
			address.Addr().AsSlice()[9:12]...,
		))

		return result
	case bits > 56 && bits <= 64:
		result, _ := netip.AddrFromSlice(address.Addr().AsSlice()[9:13])

		return result
	default:
		result, _ := netip.AddrFromSlice(address.Addr().AsSlice()[12:])

		return result
	}
}

func computeIPv6AddressEUI64(prefix netip.Addr, mac net.HardwareAddr) netip.Addr {
	if !prefix.Is6() || len(mac) != 6 {
		return netip.Addr{}
	}

	newAddress := prefix.AsSlice()

	// copy first part of mac
	copy(newAddress[8:11], mac[0:3])
	// revert the "u" bit
	newAddress[8] ^= 0x02
	// insert FFFE hexadecimal
	newAddress[11] = 0xff
	newAddress[12] = 0xfe
	// copy second part of mac
	copy(newAddress[13:16], mac[3:6])

	newAddr, _ := netip.AddrFromSlice(newAddress)

	return newAddr
}

func computeIPv6AddressOpaque(
	prefix netip.Addr,
	netIface []byte,
	networkID []byte,
	dadCounter uint32,
	secretKey []byte,
) netip.Addr {
	if !prefix.Is6() || len(netIface) == 0 || len(secretKey) < 16 {
		return netip.Addr{}
	}

	newAddress := prefix.AsSlice()

	hash := sha256.New()
	_, _ = hash.Write(newAddress[0:8]) // It never returns an error.
	_, _ = hash.Write(netIface)        // It never returns an error.
	_, _ = hash.Write(networkID)       // It never returns an error.
	if err := binary.Write(hash, binary.LittleEndian, dadCounter); err != nil {
		return netip.Addr{}
	}
	_, _ = hash.Write(secretKey) // It never returns an error.

	// compute a random identifier and limit to 64bit
	iid := hash.Sum(nil)[0:8]
	newAddr, _ := netip.AddrFromSlice(append(newAddress[0:8], iid...))

	// check colision with reserved IPv6 interface identifiers
	// cf https://www.iana.org/assignments/ipv6-interface-ids/ipv6-interface-ids.xhtml
	if bytes.Equal(iid,
		[]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, // 0000:0000:0000:0000
	) {
		goto colision
	}
	if first, last := bytes.Compare(iid,
		[]byte{0x02, 0x0, 0x5e, 0xff, 0xfe, 0x0, 0x0, 0x0}, // 0200:5EFF:FE00:0000
	), bytes.Compare(iid,
		[]byte{0x02, 0x0, 0x5e, 0xff, 0xfe, 0x0, 0x52, 0x12}, // 0200:5EFF:FE00:5212
	); first >= 0 && last <= 0 {
		goto colision
	}
	if bytes.Equal(iid,
		[]byte{0x02, 0x0, 0x5e, 0xff, 0xfe, 0x0, 0x52, 0x13}, // 0200:5EFF:FE00:5213
	) {
		goto colision
	}
	if first, last := bytes.Compare(iid,
		[]byte{0x02, 0x0, 0x5e, 0xff, 0xfe, 0x0, 0x52, 0x14}, // 0200:5EFF:FE00:5214
	), bytes.Compare(iid,
		[]byte{0x02, 0x0, 0x5e, 0xff, 0xfe, 0xff, 0xff, 0xff}, // 0200:5EFF:FEFF:FFFF
	); first >= 0 && last <= 0 {
		goto colision
	}
	if first, last := bytes.Compare(iid,
		[]byte{0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x80}, // FDFF:FFFF:FFFF:FF80
	), bytes.Compare(iid,
		[]byte{0xfd, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, // FDFF:FFFF:FFFF:FFFF
	); first >= 0 && last <= 0 {
		goto colision
	}

	return newAddr

colision: // retry with DAD_counter+1
	return computeIPv6AddressOpaque(prefix, netIface, networkID, dadCounter+1, secretKey)
}

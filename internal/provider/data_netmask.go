package provider

import (
	"math/bits"
	"net/netip"
)

//nolint:nakedret
func ipAddrToMaskBits(mask netip.Addr) (_ int, _ bool) {
	if !mask.IsValid() || !mask.Is4() {
		return
	}

	maskOcts := mask.AsSlice()
	switch {
	case maskOcts[3] != 0:
		if bits.Len8(bits.Reverse8(maskOcts[3])) != bits.OnesCount8(maskOcts[3]) {
			return
		}
		if maskOcts[0] != 255 || maskOcts[1] != 255 || maskOcts[2] != 255 {
			return
		}
	case maskOcts[2] != 0:
		if bits.Len8(bits.Reverse8(maskOcts[2])) != bits.OnesCount8(maskOcts[2]) {
			return
		}
		if maskOcts[0] != 255 || maskOcts[1] != 255 {
			return
		}
	case maskOcts[1] != 0:
		if bits.Len8(bits.Reverse8(maskOcts[1])) != bits.OnesCount8(maskOcts[1]) {
			return
		}
		if maskOcts[0] != 255 {
			return
		}
	case maskOcts[0] != 0:
		if bits.Len8(bits.Reverse8(maskOcts[0])) != bits.OnesCount8(maskOcts[0]) {
			return
		}
	}

	return bits.OnesCount8(maskOcts[0]) +
			bits.OnesCount8(maskOcts[1]) +
			bits.OnesCount8(maskOcts[2]) +
			bits.OnesCount8(maskOcts[3]),
		true
}

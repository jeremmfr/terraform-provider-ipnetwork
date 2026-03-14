package provider

import (
	"context"
	"encoding/binary"
	"math/bits"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = rangeToPrefixesFunction{}

func newRangeToPrefixesFunction() function.Function {
	return rangeToPrefixesFunction{}
}

type rangeToPrefixesFunction struct{}

func (f rangeToPrefixesFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "range_to_prefixes"
}

func (f rangeToPrefixesFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Convert an IP range to a list of CIDR prefixes.",
		Description: "Convert a range of IP addresses defined by a start and end address " +
			"into the minimal list of CIDR prefixes that exactly cover the range.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "start",
				Description: "Start address of the range",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			function.StringParameter{
				Name:        "end",
				Description: "End address of the range",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Return: function.ListReturn{
			ElementType: types.StringType,
		},
	}
}

func (f rangeToPrefixesFunction) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var inputStart, inputEnd string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputStart, &inputEnd))
	if resp.Error != nil {
		return
	}

	start, err := netip.ParseAddr(strings.TrimSpace(inputStart))
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid address"),
			function.NewFuncError("unable to parse start address input: "+err.Error()),
		)

		return
	}

	end, err := netip.ParseAddr(strings.TrimSpace(inputEnd))
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(1, "Invalid address"),
			function.NewFuncError("unable to parse end address input: "+err.Error()),
		)

		return
	}

	if start.Is4() != end.Is4() {
		resp.Error = function.ConcatFuncErrors(
			function.NewFuncError("start and end addresses must be the same IP version"),
		)

		return
	}

	if end.Less(start) {
		resp.Error = function.ConcatFuncErrors(
			function.NewFuncError("start address must be less than or equal to end address"),
		)

		return
	}

	result := rangeToPrefixes(start, end)

	resultStrings := make([]string, len(result))
	for i, p := range result {
		resultStrings[i] = p.String()
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, resultStrings))
}

// rangeToPrefixes converts an IP range [start, end] into the minimal list of
// CIDR prefixes that exactly cover that range.
func rangeToPrefixes(start, end netip.Addr) []netip.Prefix {
	prefixes := make([]netip.Prefix, 0)
	maxBits := start.BitLen()

	current := start
	for current.Compare(end) <= 0 {
		// Largest possible aligned prefix based on trailing zero bits
		bestBits := maxBits - addrTrailingZeros(current)

		// Shrink until the prefix doesn't exceed end
		for bestBits < maxBits {
			prefix := netip.PrefixFrom(current, bestBits)
			lastAddr := prefixLastAddr(prefix)

			if lastAddr.Compare(end) <= 0 {
				break
			}

			bestBits++
		}

		prefix := netip.PrefixFrom(current, bestBits)
		prefixes = append(prefixes, prefix)

		// Move to the next address after this prefix
		lastAddr := prefixLastAddr(prefix)
		next := lastAddr.Next()
		if !next.IsValid() {
			break
		}

		current = next
	}

	return prefixes
}

// addrTrailingZeros returns the number of trailing zero bits in an address.
func addrTrailingZeros(addr netip.Addr) int {
	b := addr.As16()

	if addr.Is4() {
		return bits.TrailingZeros32(binary.BigEndian.Uint32(b[12:16]))
	}

	lo := binary.BigEndian.Uint64(b[8:16])
	if lo != 0 {
		return bits.TrailingZeros64(lo)
	}

	return bits.TrailingZeros64(binary.BigEndian.Uint64(b[0:8])) + 64
}

// prefixLastAddr returns the last address in a prefix.
func prefixLastAddr(prefix netip.Prefix) netip.Addr {
	addr := prefix.Masked().Addr()
	b := addr.As16()

	// For IPv4-in-IPv6 representation, IPv4 bytes start at offset 12
	byteOffset := 0
	if addr.Is4() {
		byteOffset = 12
	}

	// Set all host bits to 1
	for i := prefix.Bits(); i < addr.BitLen(); i++ {
		byteIndex := byteOffset + i/8
		bitIndex := 7 - (i % 8)
		b[byteIndex] |= 1 << bitIndex
	}

	if addr.Is4() {
		return netip.AddrFrom16(b).Unmap()
	}

	return netip.AddrFrom16(b)
}

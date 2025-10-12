package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = isPrivateFunction{}

func newIsPrivateFunction() function.Function {
	return isPrivateFunction{}
}

type isPrivateFunction struct{}

func (f isPrivateFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "is_private"
}

func (f isPrivateFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Reports whether an address or prefix is private (internally routable).",
		Description: "Reports whether an address or prefix is private (internally routable). " +
			"For single addresses, checks if the address is private. " +
			"For prefixes (CIDR notation), checks if the entire prefix contains only private addresses. " +
			"Returns true for RFC1918, Shared Address Space (RFC6598), Unique Local Addresses (RFC4193), " +
			"and other internally routable ranges as defined by various RFCs. ",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "input",
				Description: "Address or prefix to parse",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f isPrivateFunction) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var input string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &input))
	if resp.Error != nil {
		return
	}

	// Check if input contains a slash (CIDR notation)
	switch strings.Contains(input, "/") {
	case true:
		prefix, err := netip.ParsePrefix(input)
		switch {
		case err != nil:
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(0, "Invalid CIDR address"),
				function.NewFuncError("unable to parse prefix input: "+err.Error()),
			)
		case prefix.Addr().Is4():
			resp.Error = function.ConcatFuncErrors(
				resp.Result.Set(ctx, prefixV4IsPrivate(prefix)),
			)
		case prefix.Addr().Is6():
			resp.Error = function.ConcatFuncErrors(
				resp.Result.Set(ctx, prefixV6IsPrivate(prefix)),
			)
		}
	case false:
		address, err := netip.ParseAddr(input)
		switch {
		case err != nil:
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(0, "Invalid address"),
				function.NewFuncError("unable to parse address input: "+err.Error()),
			)
		case address.Is4():
			resp.Error = function.ConcatFuncErrors(
				resp.Result.Set(ctx, addressV4IsPrivate(address)),
			)
		case address.Is6():
			resp.Error = function.ConcatFuncErrors(
				resp.Result.Set(ctx, addressV6IsPrivate(address)),
			)
		}
	}
}

// addressV4IsPrivate checks if a single IPv4 address is private (internally routable).
// Returns true for private networks including RFC1918, Shared Address Space, and Benchmarking ranges.
// Uses netip.IsPrivate() for RFC1918 ranges and adds additional private-use ranges.
func addressV4IsPrivate(address netip.Addr) bool {
	if !address.IsValid() || !address.Is4() {
		return false
	}

	addressBytes := address.As4()
	switch {
	case addressBytes[0] == 100 && (addressBytes[1]&0xC0) == 64:
		// Check for RFC 6598 Shared Address Space (100.64.0.0/10)
		// Used for Carrier-Grade NAT
		return true
	case addressBytes[0] == 198 && (addressBytes[1]&0xFE) == 18:
		// Check for Benchmarking (198.18.0.0/15) - RFC2544
		return true
	case address.IsPrivate():
		// Use standard IsPrivate for RFC1918: 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
		return true
	default:
		return false
	}
}

// addressV6IsPrivate checks if a single IPv6 address is private (internally routable).
// Returns true for private networks including ULA, Discard-Only, IPv4/IPv6 Translation, SRv6, and Benchmarking ranges.
// Uses netip.IsPrivate() for RFC4193 ULA and adds additional private-use ranges.
func addressV6IsPrivate(address netip.Addr) bool {
	if !address.IsValid() || !address.Is6() {
		return false
	}
	if address.Is4In6() {
		return addressV4IsPrivate(address.Unmap())
	}

	addressBytes := address.As16()
	switch {
	case addressBytes[0] == 0x01 && addressBytes[1] == 0x00 &&
		addressBytes[2] == 0x00 && addressBytes[3] == 0x00 &&
		addressBytes[4] == 0x00 && addressBytes[5] == 0x00 &&
		addressBytes[6] == 0x00 && addressBytes[7] == 0x00:
		// Check for Discard prefix 100::/64 - RFC6666
		return true
	case addressBytes[0] == 0x00 && addressBytes[1] == 0x64 &&
		addressBytes[2] == 0xff && addressBytes[3] == 0x9b &&
		addressBytes[4] == 0x00 && addressBytes[5] == 0x01:
		// Check for Local-Use IPv4/IPv6 Translation prefix 64:ff9b:1::/48 - RFC8215
		return true
	case addressBytes[0] == 0x5f && addressBytes[1] == 0x00:
		// Check for Segment Routing (SRv6) SIDs prefix 5f00::/16 - RFC9602
		return true
	case addressBytes[0] == 0x20 && addressBytes[1] == 0x01 &&
		addressBytes[2] == 0x00 && addressBytes[3] == 0x02 &&
		addressBytes[4] == 0x00 && addressBytes[5] == 0x00:
		// Check for Benchmarking prefix 2001:2::/48 - RFC5180
		return true
	case address.IsPrivate():
		// Use standard IsPrivate for RFC4193: fc00::/7 (Unique Local Addresses)
		return true
	default:
		return false
	}
}

// prefixV4IsPrivate checks if an IPv4 prefix contains only private (internally routable) addresses.
// It checks if the prefix is entirely contained within a private range (RFC1918, RFC6598, RFC2544).
// Returns true only if ALL addresses in the prefix are private.
func prefixV4IsPrivate(prefix netip.Prefix) bool {
	if !prefix.IsValid() || !prefix.Addr().Is4() {
		return false
	}

	addressBytes := prefix.Masked().Addr().As4()
	switch {
	case addressBytes[0] == 10 &&
		prefix.Bits() >= 8:
		// Check for Private-Use 10.0.0.0/8 - RFC1918
		return true
	case addressBytes[0] == 100 && (addressBytes[1]&0xC0) == 64 &&
		prefix.Bits() >= 10:
		// Check for RFC 6598 Shared Address Space (100.64.0.0/10)
		// Used for Carrier-Grade NAT
		return true
	case addressBytes[0] == 172 && (addressBytes[1]&0xF0) == 16 &&
		prefix.Bits() >= 12:
		// Check for Private-Use 172.16.0.0/12 - RFC1918
		return true
	case addressBytes[0] == 192 && addressBytes[1] == 168 &&
		prefix.Bits() >= 16:
		// Check for Private-Use 192.168.0.0/16 - RFC1918
		return true
	case addressBytes[0] == 198 && (addressBytes[1]&0xFE) == 18 &&
		prefix.Bits() >= 15:
		// Check for Benchmarking (198.18.0.0/15) - RFC2544
		return true
	default:
		return false
	}
}

// prefixV6IsPrivate checks if an IPv6 prefix contains only private (internally routable) addresses.
// It checks if the prefix is entirely contained within a private range (RFC4193, RFC6666, RFC8215, RFC9602, RFC5180).
// Returns true only if ALL addresses in the prefix are private.
func prefixV6IsPrivate(prefix netip.Prefix) bool {
	if !prefix.IsValid() || !prefix.Addr().Is6() {
		return false
	}
	if prefix.Addr().Is4In6() {
		if prefix.Bits() < 96 {
			return false
		}

		return prefixV4IsPrivate(netip.PrefixFrom(prefix.Addr().Unmap(), prefix.Bits()-96))
	}

	addressBytes := prefix.Masked().Addr().As16()
	switch {
	case addressBytes[0] == 0x01 && addressBytes[1] == 0x00 &&
		addressBytes[2] == 0x00 && addressBytes[3] == 0x00 &&
		addressBytes[4] == 0x00 && addressBytes[5] == 0x00 &&
		addressBytes[6] == 0x00 && addressBytes[7] == 0x00 &&
		prefix.Bits() >= 64:
		// Check for Discard prefix 100::/64 - RFC6666
		return true
	case addressBytes[0] == 0x00 && addressBytes[1] == 0x64 &&
		addressBytes[2] == 0xff && addressBytes[3] == 0x9b &&
		addressBytes[4] == 0x00 && addressBytes[5] == 0x01 &&
		prefix.Bits() >= 48:
		// Check for Local-Use IPv4/IPv6 Translation prefix 64:ff9b:1::/48 - RFC8215
		return true
	case addressBytes[0] == 0x5f && addressBytes[1] == 0x00 &&
		prefix.Bits() >= 16:
		// Check for Segment Routing (SRv6) SIDs prefix 5f00::/16 - RFC9602
		return true
	case addressBytes[0] == 0x20 && addressBytes[1] == 0x01 &&
		addressBytes[2] == 0x00 && addressBytes[3] == 0x02 &&
		addressBytes[4] == 0x00 && addressBytes[5] == 0x00 &&
		prefix.Bits() >= 48:
		// Check for Benchmarking prefix 2001:2::/48 - RFC5180
		return true
	case addressBytes[0]&0xFE == 0xfc &&
		prefix.Bits() >= 7:
		// Check for Unique-Local fc00::/7 - RFC4193
		return true
	default:
		return false
	}
}

package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = isPublicFunction{}

func newIsPublicFunction() function.Function {
	return isPublicFunction{}
}

type isPublicFunction struct{}

func (f isPublicFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "is_public"
}

func (f isPublicFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Reports whether an address or prefix is public (globally routable).",
		Description: "Reports whether an address or prefix is public (globally routable). " +
			"For single addresses, checks if the address is public. " +
			"For prefixes (CIDR notation), checks if the entire prefix contains only public addresses. " +
			"Returns false for private, reserved, documentation, or special-use addresses. ",
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

func (f isPublicFunction) Run(
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
				resp.Result.Set(ctx, prefixV4IsPublic(prefix)),
			)
		case prefix.Addr().Is6():
			resp.Error = function.ConcatFuncErrors(
				resp.Result.Set(ctx, prefixV6IsPublic(prefix)),
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
				resp.Result.Set(ctx, addressV4IsPublic(address)),
			)
		case address.Is6():
			resp.Error = function.ConcatFuncErrors(
				resp.Result.Set(ctx, addressV6IsPublic(address)),
			)
		}
	}
}

// addressV4IsPublic checks if a single IPv4 address is public (globally routable).
// Returns false for private, reserved, documentation, and special-use addresses.
func addressV4IsPublic(address netip.Addr) bool {
	if !address.IsValid() || !address.Is4() {
		return false
	}

	addressBytes := address.As4()
	switch {
	case addressBytes[0] == 0:
		// Check for "This network" (0.0.0.0/8) - RFC791
		return false
	case addressBytes[0] == 100 && (addressBytes[1]&0xC0) == 64:
		// Check for RFC 6598 Shared Address Space (100.64.0.0/10)
		// Used for Carrier-Grade NAT
		return false
	case addressBytes[0] == 192 && addressBytes[1] == 0 && addressBytes[2] == 0:
		// Check for IETF Protocol Assignments (192.0.0.0/24) - RFC6890
		return false
	case addressBytes[0] == 192 && addressBytes[1] == 0 && addressBytes[2] == 2,
		addressBytes[0] == 198 && addressBytes[1] == 51 && addressBytes[2] == 100,
		addressBytes[0] == 203 && addressBytes[1] == 0 && addressBytes[2] == 113:
		// Check for Documentation ranges (RFC5737)
		// TEST-NET-1 (192.0.2.0/24), TEST-NET-2 (198.51.100.0/24), TEST-NET-3 (203.0.113.0/24)
		return false
	case addressBytes[0] == 198 && (addressBytes[1]&0xFE) == 18:
		// Check for Benchmarking (198.18.0.0/15) - RFC2544
		return false
	case addressBytes[0]&0xF0 == 240:
		// Check for Reserved (240.0.0.0/4) - RFC1112 Section 4
		return false
	case address.IsPrivate(),
		address.IsLoopback(),
		address.IsLinkLocalUnicast(),
		address.IsMulticast(),
		address.IsUnspecified():
		// An address is NOT public if it's any of these special-use addresses
		return false
	default:
		return true
	}
}

// addressV6IsPublic checks if a single IPv6 address is public (globally routable).
// Returns false for private, reserved, documentation, and special-use addresses.
func addressV6IsPublic(address netip.Addr) bool {
	if !address.IsValid() || !address.Is6() {
		return false
	}
	if address.Is4In6() {
		return addressV4IsPublic(address.Unmap())
	}

	addressBytes := address.As16()
	switch {
	case addressBytes[0] == 0x00 && addressBytes[1] == 0x64 &&
		addressBytes[2] == 0xff && addressBytes[3] == 0x9b &&
		addressBytes[4] == 0x00 && addressBytes[5] == 0x01:
		// Check for Local-Use IPv4/IPv6 Translation prefix 64:ff9b:1::/48 - RFC8215
		return false
	case addressBytes[0] == 0x01 && addressBytes[1] == 0x00 &&
		addressBytes[2] == 0x00 && addressBytes[3] == 0x00 &&
		addressBytes[4] == 0x00 && addressBytes[5] == 0x00 &&
		addressBytes[6] == 0x00 && addressBytes[7] == 0x00:
		// Check for Discard prefix 100::/64 - RFC6666
		return false
	case addressBytes[0] == 0x20 && addressBytes[1] == 0x01 &&
		addressBytes[2] == 0x00 && addressBytes[3] == 0x02 &&
		addressBytes[4] == 0x00 && addressBytes[5] == 0x00:
		// Check for Benchmarking prefix 2001:2::/48 - RFC5180
		return false
	case addressBytes[0] == 0x20 && addressBytes[1] == 0x01 &&
		addressBytes[2] == 0x0d && addressBytes[3] == 0xb8:
		// Check for Documentation prefix 2001:db8::/32 - RFC3849
		return false
	case addressBytes[0] == 0x3f && addressBytes[1] == 0xff && (addressBytes[2]&0xF0) == 0x00:
		// Check for Documentation prefix 3fff::/20 - RFC9637
		// 3fff::/20 covers 3fff:0000:: to 3fff:0fff:ffff:...:ffff
		return false
	case addressBytes[0] == 0x5f && addressBytes[1] == 0x00:
		// Check for Segment Routing (SRv6) SIDs prefix 5f00::/16 - RFC9602
		return false
	case address.IsPrivate(),
		address.IsLoopback(),
		address.IsLinkLocalUnicast(),
		address.IsMulticast(),
		address.IsUnspecified():
		// An address is NOT public if it's any of these special-use addresses
		return false
	default:
		return true
	}
}

// prefixV4IsPublic checks if an IPv4 prefix contains only public addresses.
// It checks if the prefix overlaps with (contains or is contained in) any private, reserved, or documentation ranges.
func prefixV4IsPublic(prefix netip.Prefix) bool {
	if !prefix.IsValid() || !prefix.Addr().Is4() {
		return false
	}

	// Define all non-public IPv4 ranges
	nonPublicRanges := []netip.Prefix{
		netip.MustParsePrefix("0.0.0.0/8"),       // "This network" - RFC791
		netip.MustParsePrefix("10.0.0.0/8"),      // Private - RFC1918
		netip.MustParsePrefix("100.64.0.0/10"),   // Shared Address Space - RFC6598
		netip.MustParsePrefix("127.0.0.0/8"),     // Loopback
		netip.MustParsePrefix("169.254.0.0/16"),  // Link-local
		netip.MustParsePrefix("172.16.0.0/12"),   // Private - RFC1918
		netip.MustParsePrefix("192.0.0.0/24"),    // IETF Protocol Assignments - RFC6890
		netip.MustParsePrefix("192.0.2.0/24"),    // TEST-NET-1 - RFC5737
		netip.MustParsePrefix("192.168.0.0/16"),  // Private - RFC1918
		netip.MustParsePrefix("198.18.0.0/15"),   // Benchmarking - RFC2544
		netip.MustParsePrefix("198.51.100.0/24"), // TEST-NET-2 - RFC5737
		netip.MustParsePrefix("203.0.113.0/24"),  // TEST-NET-3 - RFC5737
		netip.MustParsePrefix("224.0.0.0/4"),     // Multicast
		netip.MustParsePrefix("240.0.0.0/4"),     // Reserved (includes broadcast)
	}

	// Check if prefix overlaps with any non-public range
	for _, nonPublic := range nonPublicRanges {
		if prefix.Overlaps(nonPublic) {
			return false
		}
	}

	return true
}

// prefixV6IsPublic checks if an IPv6 prefix contains only public addresses.
// It checks if the prefix overlaps with (contains or is contained in) any private, reserved, or documentation ranges.
func prefixV6IsPublic(prefix netip.Prefix) bool {
	if !prefix.IsValid() || !prefix.Addr().Is6() {
		return false
	}
	if prefix.Addr().Is4In6() {
		if prefix.Bits() < 96 {
			return false
		}

		return prefixV4IsPublic(netip.PrefixFrom(prefix.Addr().Unmap(), prefix.Bits()-96))
	}

	// Define all non-public IPv6 ranges
	nonPublicRanges := []netip.Prefix{
		netip.MustParsePrefix("::/128"),         // Unspecified
		netip.MustParsePrefix("::1/128"),        // Loopback
		netip.MustParsePrefix("64:ff9b:1::/48"), // Local-Use IPv4/IPv6 Translation - RFC8215
		netip.MustParsePrefix("100::/64"),       // Discard-Only - RFC6666
		netip.MustParsePrefix("2001:2::/48"),    // Benchmarking - RFC5180
		netip.MustParsePrefix("2001:db8::/32"),  // Documentation - RFC3849
		netip.MustParsePrefix("3fff::/20"),      // Documentation - RFC9637
		netip.MustParsePrefix("5f00::/16"),      // Segment Routing (SRv6) SIDs
		netip.MustParsePrefix("fc00::/7"),       // Unique Local Addresses (ULA)
		netip.MustParsePrefix("fe80::/10"),      // Link-local
		netip.MustParsePrefix("ff00::/8"),       // Multicast
	}

	// Check if prefix overlaps with any non-public range
	for _, nonPublic := range nonPublicRanges {
		if prefix.Overlaps(nonPublic) {
			return false
		}
	}

	return true
}

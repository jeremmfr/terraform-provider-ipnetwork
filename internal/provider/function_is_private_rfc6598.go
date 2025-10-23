package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = isPrivateRFC6598Function{}

func newIsPrivateRFC6598Function() function.Function {
	return isPrivateRFC6598Function{}
}

type isPrivateRFC6598Function struct{}

func (f isPrivateRFC6598Function) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "is_private_rfc6598"
}

func (f isPrivateRFC6598Function) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Reports whether an address or prefix is in RFC6598 Shared Address Space.",
		Description: "Reports whether an address or prefix is in RFC6598 Shared Address Space. " +
			"For single addresses, checks if the address is in the Shared Address Space range. " +
			"For prefixes (CIDR notation), checks if the entire prefix is contained within the Shared Address Space range. " +
			"RFC6598 defines the 100.64.0.0/10 block for Carrier-Grade NAT (CGN). ",
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

func (f isPrivateRFC6598Function) Run(
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
		if err != nil {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(0, "Invalid CIDR address"),
				function.NewFuncError("unable to parse prefix input: "+err.Error()),
			)

			return
		}
		resp.Error = function.ConcatFuncErrors(
			resp.Result.Set(ctx, prefixIsPrivateRFC6598(prefix)),
		)
	case false:
		address, err := netip.ParseAddr(input)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(0, "Invalid address"),
				function.NewFuncError("unable to parse address input: "+err.Error()),
			)

			return
		}
		resp.Error = function.ConcatFuncErrors(
			resp.Result.Set(ctx, addressIsPrivateRFC6598(address)),
		)
	}
}

// addressIsPrivateRFC6598 checks if an IPv4 address is in RFC6598 Shared Address Space.
func addressIsPrivateRFC6598(address netip.Addr) bool {
	if !address.IsValid() || (address.Is6() && !address.Is4In6()) {
		return false
	}

	if address.Is4In6() {
		address = address.Unmap()
	}

	addressBytes := address.As4()

	// 100.64.0.0/10
	return addressBytes[0] == 100 && (addressBytes[1]&0xC0) == 64
}

// prefixIsPrivateRFC6598 checks if an IPv4 prefix is entirely contained within RFC6598 Shared Address Space.
func prefixIsPrivateRFC6598(prefix netip.Prefix) bool {
	if !prefix.IsValid() || (prefix.Addr().Is6() && !prefix.Addr().Is4In6()) {
		return false
	}

	if prefix.Addr().Is4In6() {
		if prefix.Bits() < 96 {
			return false
		}

		return prefixIsPrivateRFC6598(netip.PrefixFrom(prefix.Addr().Unmap(), prefix.Bits()-96))
	}

	addressBytes := prefix.Masked().Addr().As4()

	// 100.64.0.0/10
	return addressBytes[0] == 100 && (addressBytes[1]&0xC0) == 64 &&
		prefix.Bits() >= 10
}

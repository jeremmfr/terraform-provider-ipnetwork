package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = isPrivateRFC4193Function{}

func newIsPrivateRFC4193Function() function.Function {
	return isPrivateRFC4193Function{}
}

type isPrivateRFC4193Function struct{}

func (f isPrivateRFC4193Function) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "is_private_rfc4193"
}

func (f isPrivateRFC4193Function) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Reports whether an address or prefix is in RFC4193 Unique Local Address (ULA) space.",
		Description: "Reports whether an address or prefix is in RFC4193 Unique Local Address (ULA) space. " +
			"For single addresses, checks if the address is in the ULA range. " +
			"For prefixes (CIDR notation), checks if the entire prefix is contained within the ULA range. " +
			"RFC4193 defines the fc00::/7 block for Unique Local Addresses. ",
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

func (f isPrivateRFC4193Function) Run(
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
			resp.Result.Set(ctx, prefixIsPrivateRFC4193(prefix)),
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
			resp.Result.Set(ctx, addressIsPrivateRFC4193(address)),
		)
	}
}

// addressIsPrivateRFC4193 checks if an IPv6 address is in RFC4193 ULA space.
func addressIsPrivateRFC4193(address netip.Addr) bool {
	if !address.IsValid() || !address.Is6() {
		return false
	}

	addressBytes := address.As16()

	// fc00::/7
	return addressBytes[0]&0xFE == 0xfc
}

// prefixIsPrivateRFC4193 checks if an IPv6 prefix is entirely contained within RFC4193 ULA space.
func prefixIsPrivateRFC4193(prefix netip.Prefix) bool {
	if !prefix.IsValid() || !prefix.Addr().Is6() {
		return false
	}

	addressBytes := prefix.Masked().Addr().As16()
	// fc00::/7
	return addressBytes[0]&0xFE == 0xfc &&
		prefix.Bits() >= 7
}

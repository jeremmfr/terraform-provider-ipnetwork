package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = isPrivateRFC1918Function{}

func newIsPrivateRFC1918Function() function.Function {
	return isPrivateRFC1918Function{}
}

type isPrivateRFC1918Function struct{}

func (f isPrivateRFC1918Function) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "is_private_rfc1918"
}

func (f isPrivateRFC1918Function) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Reports whether an address or prefix is in RFC1918 private address space.",
		Description: "Reports whether an address or prefix is in RFC1918 private address space. " +
			"For single addresses, checks if the address is in RFC1918 ranges. " +
			"For prefixes (CIDR notation), checks if the entire prefix is contained within RFC1918 ranges. " +
			"RFC1918 defines three blocks: 10.0.0.0/8, 172.16.0.0/12, and 192.168.0.0/16. ",
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

func (f isPrivateRFC1918Function) Run(
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
			resp.Result.Set(ctx, prefixIsPrivateRFC1918(prefix)),
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
			resp.Result.Set(ctx, addressIsPrivateRFC1918(address)),
		)
	}
}

// addressIsPrivateRFC1918 checks if an IPv4 address is in RFC1918 private address space.
func addressIsPrivateRFC1918(address netip.Addr) bool {
	if !address.IsValid() || (address.Is6() && !address.Is4In6()) {
		return false
	}

	if address.Is4In6() {
		address = address.Unmap()
	}

	addressBytes := address.As4()
	switch {
	case addressBytes[0] == 10:
		// 10.0.0.0/8
		return true
	case addressBytes[0] == 172 && (addressBytes[1]&0xF0) == 16:
		// 172.16.0.0/12
		return true
	case addressBytes[0] == 192 && addressBytes[1] == 168:
		// 192.168.0.0/16
		return true
	default:
		return false
	}
}

// prefixIsPrivateRFC1918 checks if an IPv4 prefix is entirely contained within RFC1918 private address space.
func prefixIsPrivateRFC1918(prefix netip.Prefix) bool {
	if !prefix.IsValid() || (prefix.Addr().Is6() && !prefix.Addr().Is4In6()) {
		return false
	}

	if prefix.Addr().Is4In6() {
		if prefix.Bits() < 96 {
			return false
		}

		return prefixIsPrivateRFC1918(netip.PrefixFrom(prefix.Addr().Unmap(), prefix.Bits()-96))
	}

	addressBytes := prefix.Masked().Addr().As4()
	switch {
	case addressBytes[0] == 10 &&
		prefix.Bits() >= 8:
		// 10.0.0.0/8
		return true
	case addressBytes[0] == 172 && (addressBytes[1]&0xF0) == 16 &&
		prefix.Bits() >= 12:
		// 172.16.0.0/12
		return true
	case addressBytes[0] == 192 && addressBytes[1] == 168 &&
		prefix.Bits() >= 16:
		// 192.168.0.0/16
		return true
	default:
		return false
	}
}

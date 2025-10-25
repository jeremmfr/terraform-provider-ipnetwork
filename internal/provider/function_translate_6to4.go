package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = translate6to4Function{}

func newTranslate6to4Function() function.Function {
	return translate6to4Function{}
}

type translate6to4Function struct{}

func (f translate6to4Function) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "translate_6to4"
}

func (f translate6to4Function) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Translate an IPv6 address to an IPv4 address.",
		Description: "Translate an IPv6 address to an IPv4 address," +
			" as defined in RFC 6052 section 2.2.\n" +
			" Mask of address determines how the IPv4 address is embedded.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "input",
				Description: "Address to parse",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Return: function.StringReturn{},
	}
}

func (f translate6to4Function) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var input string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &input))
	if resp.Error != nil {
		return
	}

	var address netip.Prefix
	switch strings.Contains(input, "/") {
	case true:
		var err error
		address, err = netip.ParsePrefix(input)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(1, "Invalid address"),
				function.NewFuncError("unable to parse address input: "+err.Error()),
			)

			return
		}
		if !address.Addr().Is6() {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(1, "Invalid address"),
				function.NewFuncError("must be an IPv6 address"),
			)

			return
		}
	case false:
		onlyAddress, err := netip.ParseAddr(input)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(1, "Invalid address"),
				function.NewFuncError("unable to parse address input: "+err.Error()),
			)

			return
		}
		if !onlyAddress.Is6() {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(1, "Invalid address"),
				function.NewFuncError("must be an IPv6 address"),
			)

			return
		}

		address = netip.PrefixFrom(onlyAddress, 96)
	}

	output := translateAddress6to4(address)
	if !output.IsValid() {
		// if happen, it's a bug
		resp.Error = function.NewFuncError("Internal Error," +
			" this is a bug in the provider, which should be reported in the provider's own issue tracker.")

		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, output.String()))
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

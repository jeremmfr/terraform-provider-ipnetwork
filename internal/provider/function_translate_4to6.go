package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = translate4to6Function{}

func newTranslate4to6Function() function.Function {
	return translate4to6Function{}
}

type translate4to6Function struct{}

func (f translate4to6Function) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "translate_4to6"
}

func (f translate4to6Function) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Translate an IPv4 address to an IPv6 address.",
		Description: "Translate an IPv4 address to an IPv6 address using an IPv6 prefix," +
			" as defined in RFC 6052 section 2.2.\n" +
			" Mask of prefix determines how the IPv4 address is embedded.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "address",
				Description: "Address to parse",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			function.StringParameter{
				Name:        "prefix",
				Description: "Prefix address to parse",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Return: function.StringReturn{},
	}
}

func (f translate4to6Function) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var inputAddress, inputPrefix string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputAddress, &inputPrefix))
	if resp.Error != nil {
		return
	}

	// remove potential mask
	inputAddress, _, _ = strings.Cut(inputAddress, "/")

	address, err := netip.ParseAddr(inputAddress)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid address"),
			function.NewFuncError("unable to parse address input: "+err.Error()),
		)

		return
	}
	if !address.Is4() {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid address"),
			function.NewFuncError("must be an IPv4 address"),
		)

		return
	}

	var prefix netip.Prefix
	switch strings.Contains(inputPrefix, "/") {
	case true:
		var err error
		prefix, err = netip.ParsePrefix(inputPrefix)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(1, "Invalid prefix address"),
				function.NewFuncError("unable to parse prefix address input: "+err.Error()),
			)

			return
		}
		if !prefix.Addr().Is6() {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(1, "Invalid prefix address"),
				function.NewFuncError("must be an IPv6 address"),
			)

			return
		}
	case false:
		prefixAddress, err := netip.ParseAddr(inputPrefix)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(1, "Invalid prefix address"),
				function.NewFuncError("unable to parse prefix address input: "+err.Error()),
			)

			return
		}
		if !prefixAddress.Is6() {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(1, "Invalid prefix address"),
				function.NewFuncError("must be an IPv6 address"),
			)

			return
		}

		prefix = netip.PrefixFrom(prefixAddress, 96)
	}

	output := translateAddress4to6(address, prefix)
	if !output.IsValid() {
		// if happen, it's a bug
		resp.Error = function.NewFuncError("Internal Error," +
			" this is a bug in the provider, which should be reported in the provider's own issue tracker.")

		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, output.String()))
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

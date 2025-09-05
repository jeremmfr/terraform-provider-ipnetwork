package provider

import (
	"context"
	"net"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = generate6EUI64Function{}

func newGenerate6EUI64Function() function.Function {
	return generate6EUI64Function{}
}

type generate6EUI64Function struct{}

func (f generate6EUI64Function) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "generate6_eui64"
}

func (f generate6EUI64Function) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Generate an IPv6 address from MAC address with the modified EUI-64 format.",
		Description: "Generate an IPv6 address from MAC address with the modified EUI-64 format," +
			" as defined in RFC 4291 section 2.5.1.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "prefix",
				Description: "IPv6 prefix address to parse",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			function.StringParameter{
				Name:        "mac",
				Description: "MAC address to parse",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Return: function.StringReturn{},
	}
}

func (f generate6EUI64Function) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var inputPrefix, inputMac string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputPrefix, &inputMac))
	if resp.Error != nil {
		return
	}

	// remove potential mask
	inputPrefix, _, _ = strings.Cut(inputPrefix, "/")
	// remove potential scoped zone
	inputPrefix, _, _ = strings.Cut(inputPrefix, "%")

	prefix, err := netip.ParseAddr(inputPrefix)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid prefix"),
			function.NewFuncError("unable to parse prefix address input: "+err.Error()),
		)

		return
	}
	if !prefix.Is6() {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid prefix"),
			function.NewFuncError("prefix address must be an IPv6 address"),
		)

		return
	}

	mac, err := net.ParseMAC(inputMac)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid MAC"),
			function.NewFuncError("unable to parse MAC address input: "+err.Error()),
		)

		return
	}
	if len(mac) != 6 {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid MAC"),
			function.NewFuncError("MAC address must be in EUI-48 format"),
		)

		return
	}

	output := computeIPv6AddressEUI64(prefix, mac)
	if !output.IsValid() {
		// if happen, it's a bug
		resp.Error = function.NewFuncError("Internal Error," +
			" this is a bug in the provider, which should be reported in the provider's own issue tracker.")

		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, output.String()))
}

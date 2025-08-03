package provider

import (
	"context"
	"net/netip"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = expand6Function{}

func newExpand6Function() function.Function {
	return expand6Function{}
}

type expand6Function struct{}

func (f expand6Function) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "expand6"
}

func (f expand6Function) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Expand IPv6 address (CIDR or not).",
		Description: "Expand IPv6 address, with CIDR format or not, to" +
			" long format (leading zeroes and no '::' compression).",
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

func (f expand6Function) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var input string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &input))
	if resp.Error != nil {
		return
	}

	switch strings.Contains(input, "/") {
	case true:
		output, err := netip.ParsePrefix(input)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(0, "Invalid CIDR address"),
				function.NewFuncError("unable to parse address input: "+err.Error()),
			)

			return
		}

		resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx,
			output.Addr().StringExpanded()+"/"+strconv.Itoa(output.Bits())))

	case false:
		output, err := netip.ParseAddr(input)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(0, "Invalid address"),
				function.NewFuncError("unable to parse address input: "+err.Error()),
			)

			return
		}

		resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, output.StringExpanded()))
	}
}

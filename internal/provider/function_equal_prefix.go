package provider

import (
	"context"
	"net/netip"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = equalPrefixFunction{}

func newEqualPrefixFunction() function.Function {
	return equalPrefixFunction{}
}

type equalPrefixFunction struct{}

func (f equalPrefixFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "equal_prefix"
}

func (f equalPrefixFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary:     "Compare two CIDR addresses if they are in the same prefix.",
		Description: "Compare two CIDR addresses if they are in the same prefix.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "address_x",
				Description: "First address to parse",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			function.StringParameter{
				Name:        "address_y",
				Description: "Second address to parse",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f equalPrefixFunction) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var inputAddressX, inputAddressY string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputAddressX, &inputAddressY))
	if resp.Error != nil {
		return
	}

	addressX, err := netip.ParsePrefix(inputAddressX)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid CIDR address"),
			function.NewFuncError("unable to parse address input: "+err.Error()),
		)

		return
	}
	addressY, err := netip.ParsePrefix(inputAddressY)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(1, "Invalid CIDR address"),
			function.NewFuncError("unable to parse address input: "+err.Error()),
		)

		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, addressX.Masked() == addressY.Masked()))
}

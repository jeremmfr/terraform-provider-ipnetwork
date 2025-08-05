package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = equalAddressFunction{}

func newEqualAddressFunction() function.Function {
	return equalAddressFunction{}
}

type equalAddressFunction struct{}

func (f equalAddressFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "equal_address"
}

func (f equalAddressFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Compare two address if there are equal.",
		Description: "Compare two address if there are equal" +
			" regardless of format: CIDR or not, IPv6 expanded or not.",
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

func (f equalAddressFunction) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var inputAddressX, inputAddressY string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputAddressX, &inputAddressY))
	if resp.Error != nil {
		return
	}

	// remove potential mask
	inputAddressXCut, _, _ := strings.Cut(inputAddressX, "/")
	inputAddressYCut, _, _ := strings.Cut(inputAddressY, "/")

	// remove potential scoped zone
	inputAddressXCut, _, _ = strings.Cut(inputAddressXCut, "%")
	inputAddressYCut, _, _ = strings.Cut(inputAddressYCut, "%")

	addressX, err := netip.ParseAddr(inputAddressXCut)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid address"),
			function.NewFuncError("unable to parse address input: "+err.Error()),
		)

		return
	}
	addressY, err := netip.ParseAddr(inputAddressYCut)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(1, "Invalid address"),
			function.NewFuncError("unable to parse address input: "+err.Error()),
		)

		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, addressX == addressY))
}

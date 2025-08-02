package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = addressFunction{}

func newAddressFunction() function.Function {
	return addressFunction{}
}

type addressFunction struct{}

func (f addressFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "address"
}

func (f addressFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Validate an address.",
		Description: "Validate an address" +
			" with completion and cleanup of unwanted data and then proper format it.",
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

func (f addressFunction) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var input string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &input))
	if resp.Error != nil {
		return
	}

	// remove potential mask
	inputAddress, _, _ := strings.Cut(input, "/")

	// clean potential leading or trailing white space
	inputAddress = strings.TrimSpace(inputAddress)
	if len(inputAddress) == 0 {
		resp.Error = function.NewArgumentFuncError(0, "String only with space character(s)")

		return
	}

	// remove potential scoped zone
	inputAddress, _, _ = strings.Cut(inputAddress, "%")

	// complete IPv4 address if missing a part
	if !strings.Contains(inputAddress, ":") && strings.Count(inputAddress, ".") != 3 {
		inputAddress = completeIPv4Address(inputAddress)
	}

	// read address
	output, err := netip.ParseAddr(inputAddress)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid address"),
			function.NewFuncError("unable to parse address input: "+err.Error()),
		)

		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, output.String()))
}

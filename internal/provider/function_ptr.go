package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = ptrFunction{}

func newPtrFunction() function.Function {
	return ptrFunction{}
}

type ptrFunction struct{}

func (f ptrFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "ptr"
}

func (f ptrFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Generate the PTR name from an address.",
		Description: "Generate the PTR name from an address.\n" +
			" Output have 'in-addr.arpa.' suffix for IPv4 address and 'ip6.arpa.' suffix for IPv6 address.",
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

func (f ptrFunction) Run(
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
	input, _, _ = strings.Cut(input, "/")

	address, err := netip.ParseAddr(input)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid address"),
			function.NewFuncError("unable to parse address input: "+err.Error()),
		)

		return
	}

	ptr := ptrNameFromIP(address)
	if ptr == "" {
		// if happen, it's a bug
		resp.Error = function.NewFuncError("Internal Error," +
			" this is a bug in the provider, which should be reported in the provider's own issue tracker.")

		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, ptr))
}

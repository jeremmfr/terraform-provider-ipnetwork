package provider

import (
	"context"
	"fmt"
	"math"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = addressPortFunction{}

func newAddressPortFunction() function.Function {
	return addressPortFunction{}
}

type addressPortFunction struct{}

func (f addressPortFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "address_port"
}

func (f addressPortFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Generate an ip:port string representation from IP address and port.",
		Description: "Generate an ip:port string representation from IP address and port" +
			" (add square brackets for IPv6 address).",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "address",
				Description: "Address to parse",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			function.Int32Parameter{
				Name:        "port",
				Description: "Port to parse",
				Validators: []function.Int32ParameterValidator{
					int32validator.Between(0, math.MaxUint16),
				},
			},
		},
		Return: function.StringReturn{},
	}
}

func (f addressPortFunction) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var (
		inputAddress string
		inputPort    int32
	)
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputAddress, &inputPort))
	if resp.Error != nil {
		return
	}

	if inputPort < 0 || inputPort > math.MaxUint16 {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(1, "Invalid port"),
			function.NewFuncError(fmt.Sprintf("port must be between %d and %d: ", 0, math.MaxUint16)),
		)

		return
	}

	// remove potential mask
	inputAddress, _, _ = strings.Cut(inputAddress, "/")
	// remove potential scoped zone
	inputAddress, _, _ = strings.Cut(inputAddress, "%")

	address, err := netip.ParseAddr(inputAddress)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid address"),
			function.NewFuncError("unable to parse address input: "+err.Error()),
		)

		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, netip.AddrPortFrom(address, uint16(inputPort)).String()))
}

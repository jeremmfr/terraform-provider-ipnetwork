package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = containFunction{}

func newContainFunction() function.Function {
	return containFunction{}
}

type containFunction struct{}

func (f containFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "contain"
}

func (f containFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Reports whether a prefix contains address(es).",
		Description: "Reports whether a prefix (container) contains" +
			" an address if address is not in CIDR format or" +
			" all addresses of address block if address is in CIDR format.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "container",
				Description: "Container address to parse",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			function.StringParameter{
				Name:        "address",
				Description: "Included address(es) to parse",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
		Return: function.BoolReturn{},
	}
}

func (f containFunction) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var inputContainer, inputAddress string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputContainer, &inputAddress))
	if resp.Error != nil {
		return
	}

	container, err := netip.ParsePrefix(inputContainer)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid CIDR address"),
			function.NewFuncError("unable to parse container address input: "+err.Error()),
		)

		return
	}

	switch strings.Contains(inputAddress, "/") {
	case true:
		address, err := netip.ParsePrefix(inputAddress)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(1, "Invalid CIDR address"),
				function.NewFuncError("unable to parse address input: "+err.Error()),
			)

			return
		}

		switch {
		case container.Addr().BitLen() != address.Addr().BitLen():
			// reports false if container and address have different IP version
			resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, false))
		case container.Bits() > address.Bits():
			// container smaller than address block
			resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, false))
		case container.Bits() <= address.Bits() && container.Contains(address.Addr()):
			resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, true))
		default:
			resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, false))
		}
	case false:
		address, err := netip.ParseAddr(inputAddress)
		if err != nil {
			resp.Error = function.ConcatFuncErrors(
				function.NewArgumentFuncError(1, "Invalid address"),
				function.NewFuncError("unable to parse address input: "+err.Error()),
			)

			return
		}

		resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, container.Contains(address)))
	}
}

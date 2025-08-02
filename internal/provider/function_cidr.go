package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var _ function.Function = cidrFunction{}

func newCidrFunction() function.Function {
	return cidrFunction{}
}

type cidrFunction struct{}

func (f cidrFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "cidr"
}

func (f cidrFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Validate a CIDR address.",
		Description: "Validate a CIDR address" +
			" with completion and replacement/cleanup of incorrect/unwanted data and then proper format it.",
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

func (f cidrFunction) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var input string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &input))
	if resp.Error != nil {
		return
	}

	// split address and mask fields
	inputAddress, inputMask, _ := strings.Cut(input, "/")

	// clean potential leading or trailing white space
	inputAddress = strings.TrimSpace(inputAddress)
	inputMask = strings.TrimSpace(inputMask)
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

	// read address part without mask
	netAddress, err := netip.ParseAddr(inputAddress)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid CIDR address"),
			function.NewFuncError("unable to parse address field: "+err.Error()),
		)

		return
	}

	var output netip.Prefix
	switch {
	case netAddress.Is4():
		switch {
		case inputMask == "":
			// no mask, so add it
			switch {
			case netAddress == netip.MustParseAddr("0.0.0.0"):
				output = netip.PrefixFrom(netAddress, 0)
			default:
				output = netip.PrefixFrom(netAddress, 32)
			}
		case strings.Count(inputMask, ".") == 3:
			// mask in potential address format
			maskAddr, err := netip.ParseAddr(inputMask)
			if err != nil {
				resp.Error = function.ConcatFuncErrors(
					function.NewArgumentFuncError(0, "Invalid CIDR address"),
					function.NewFuncError("unable to parse mask field in decimal format: "+err.Error()),
				)

				return
			}
			maskBits, ok := ipAddrToMaskBits(maskAddr)
			if !ok {
				resp.Error = function.ConcatFuncErrors(
					function.NewArgumentFuncError(0, "Invalid CIDR address"),
					function.NewFuncError("unable to parse mask field in decimal format: invalid octet"),
				)

				return
			}
			output = netip.PrefixFrom(netAddress, maskBits)

		default:
			var err error
			output, err = netip.ParsePrefix(netAddress.String() + "/" + inputMask)
			if err != nil {
				resp.Error = function.ConcatFuncErrors(
					function.NewArgumentFuncError(0, "Invalid CIDR address"),
					function.NewFuncError("unable to parse CIDR address due to mask field: "+err.Error()),
				)

				return
			}
		}
	case netAddress.Is6():
		switch {
		case inputMask == "":
			// no mask, so add it
			switch {
			case netAddress == netip.MustParseAddr("::"):
				output = netip.PrefixFrom(netAddress, 0)
			default:
				output = netip.PrefixFrom(netAddress, 128)
			}
		default:
			output, err = netip.ParsePrefix(netAddress.String() + "/" + inputMask)
			if err != nil {
				resp.Error = function.ConcatFuncErrors(
					function.NewArgumentFuncError(0, "Invalid CIDR address"),
					function.NewFuncError("unable to parse CIDR address due to mask field: "+err.Error()),
				)

				return
			}
		}
	default:
		// if happen, it's a bug
		resp.Error = function.NewFuncError("Internal Error," +
			" this is a bug in the provider, which should be reported in the provider's own issue tracker.")

		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, output.String()))
}

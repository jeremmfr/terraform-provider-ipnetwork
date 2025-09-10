package provider

import (
	"context"
	"net/netip"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/int32validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = generate6OpaqueFunction{}

func newGenerate6OpaqueFunction() function.Function {
	return generate6OpaqueFunction{}
}

type generate6OpaqueFunction struct{}

func (f generate6OpaqueFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "generate6_opaque"
}

func (f generate6OpaqueFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Generate an IPv6 address with an opaque interface identifier.",
		Description: "Generate an IPv6 address with an opaque interface identifier," +
			" as defined in RFC 7217 section 5.",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "prefix",
				Description: "IPv6 prefix address to parse",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			function.StringParameter{
				Name:        "net_iface",
				Description: "Interface identifier",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(1),
				},
			},
			function.StringParameter{
				Name:           "network_id",
				Description:    "(Optional) Network subnet identifier",
				AllowNullValue: true,
			},
			function.Int32Parameter{
				Name:           "dad_counter",
				Description:    "(Optional) Counter to resolve DAD conflict",
				AllowNullValue: true,
				Validators: []function.Int32ParameterValidator{
					int32validator.AtLeast(0),
				},
			},
			function.StringParameter{
				Name:        "secret_key",
				Description: "Secret key",
				Validators: []function.StringParameterValidator{
					stringvalidator.LengthAtLeast(16), // at least 128 bits in UTF8 encoding
				},
			},
		},
		Return: function.StringReturn{},
	}
}

func (f generate6OpaqueFunction) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var (
		inputPrefix, inputNetIface, inputSecretKey, networkID string
		inputNetworkID                                        types.String
		inputDADCounter                                       types.Int32
		dadCounter                                            int32
	)
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx,
		&inputPrefix,
		&inputNetIface,
		&inputNetworkID,
		&inputDADCounter,
		&inputSecretKey,
	))
	if resp.Error != nil {
		return
	}

	if !inputNetworkID.IsNull() {
		networkID = inputNetworkID.ValueString()
	}
	if !inputDADCounter.IsNull() {
		dadCounter = inputDADCounter.ValueInt32()
	}

	// remove potential mask
	inputPrefix, _, _ = strings.Cut(inputPrefix, "/")
	// remove potential scoped zone
	inputPrefix, _, _ = strings.Cut(inputPrefix, "%")

	prefix, err := netip.ParseAddr(inputPrefix)
	if err != nil {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid Prefix"),
			function.NewFuncError("unable to parse prefix address input: "+err.Error()),
		)

		return
	}
	if !prefix.Is6() {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(0, "Invalid Prefix"),
			function.NewFuncError("prefix address must be an IPv6 address"),
		)

		return
	}

	if inputNetIface == "" {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(1, "Invalid Net_Iface"),
			function.NewFuncError("value is empty"),
		)

		return
	}
	if dadCounter < 0 {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(3, "Invalid DAD_Counter"),
			function.NewFuncError("must be at least 0"),
		)

		return
	}
	if len(inputSecretKey) < 16 {
		resp.Error = function.ConcatFuncErrors(
			function.NewArgumentFuncError(4, "Invalid secret_key"),
			function.NewFuncError("value is too small, must be at least 128 bits in UTF8 encoding"),
		)

		return
	}

	output := computeIPv6AddressOpaque(
		prefix,
		[]byte(inputNetIface),
		[]byte(networkID),
		uint32(dadCounter),
		[]byte(inputSecretKey),
	)
	if !output.IsValid() {
		// if happen, it's a bug
		resp.Error = function.NewFuncError("Internal Error," +
			" this is a bug in the provider, which should be reported in the provider's own issue tracker.")

		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, output.String()))
}

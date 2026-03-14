package provider

import (
	"context"
	"net/netip"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = sortFunction{}

func newSortFunction() function.Function {
	return sortFunction{}
}

type sortFunction struct{}

func (f sortFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "sort"
}

func (f sortFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Sort IP addresses with or without mask in numerical order.",
		Description: "Sort a list of IP addresses with or without mask in numerical order. " +
			"When two entries share the same address, the address without mask comes first, " +
			"then CIDR addresses are sorted by mask length (shortest first).",
		Parameters: []function.Parameter{
			function.ListParameter{
				ElementType: types.StringType,
				Name:        "inputs",
				Description: "List of IP addresses to sort",
			},
		},
		Return: function.ListReturn{
			ElementType: types.StringType,
		},
	}
}

func (f sortFunction) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var inputs []string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputs))
	if resp.Error != nil {
		return
	}

	type entry struct {
		raw  string
		addr netip.Addr
		bits int
	}

	entries := make([]entry, 0, len(inputs))
	for _, item := range inputs {
		if strings.Contains(item, "/") {
			prefix, err := netip.ParsePrefix(item)
			if err != nil {
				resp.Error = function.ConcatFuncErrors(
					function.NewArgumentFuncError(0, "Invalid CIDR address"),
					function.NewFuncError("unable to parse address input: "+err.Error()),
				)

				return
			}
			entries = append(entries, entry{raw: item, addr: prefix.Addr(), bits: prefix.Bits()})
		} else {
			addr, err := netip.ParseAddr(item)
			if err != nil {
				resp.Error = function.ConcatFuncErrors(
					function.NewArgumentFuncError(0, "Invalid address"),
					function.NewFuncError("unable to parse address input: "+err.Error()),
				)

				return
			}
			entries = append(entries, entry{raw: item, addr: addr, bits: -1})
		}
	}

	slices.SortStableFunc(entries, func(a, b entry) int {
		if cmp := a.addr.Compare(b.addr); cmp != 0 {
			return cmp
		}

		return a.bits - b.bits
	})

	result := make([]string, len(entries))
	for i, e := range entries {
		result[i] = e.raw
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result))
}

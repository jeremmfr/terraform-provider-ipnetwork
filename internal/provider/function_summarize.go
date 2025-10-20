package provider

import (
	"context"
	"net/netip"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ function.Function = summarizeFunction{}

func newSummarizeFunction() function.Function {
	return summarizeFunction{}
}

type summarizeFunction struct{}

func (f summarizeFunction) Metadata(
	_ context.Context,
	_ function.MetadataRequest,
	resp *function.MetadataResponse,
) {
	resp.Name = "summarize"
}

func (f summarizeFunction) Definition(
	_ context.Context,
	_ function.DefinitionRequest,
	resp *function.DefinitionResponse,
) {
	resp.Definition = function.Definition{
		Summary: "Summarize IP prefixes.",
		Description: "Summarize a set of IP addresses and prefixes into " +
			"the smallest possible list of prefixes that cover the same addresses.",
		Parameters: []function.Parameter{
			function.SetParameter{
				ElementType: types.StringType,
				Name:        "inputs",
				Description: "Set of IP addresses and prefixes to summarize",
			},
		},
		Return: function.ListReturn{
			ElementType: types.StringType,
		},
	}
}

func (f summarizeFunction) Run(
	ctx context.Context,
	req function.RunRequest,
	resp *function.RunResponse,
) {
	var inputs []string
	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &inputs))
	if resp.Error != nil {
		return
	}

	// Convert all inputs to prefixes
	prefixes := make([]netip.Prefix, 0, len(inputs))
	for _, item := range inputs {
		switch strings.Contains(item, "/") {
		case true:
			prefix, err := netip.ParsePrefix(item)
			if err != nil {
				resp.Error = function.ConcatFuncErrors(
					function.NewArgumentFuncError(1, "Invalid CIDR address"),
					function.NewFuncError("unable to parse address input: "+err.Error()),
				)

				return
			}
			prefixes = append(prefixes, prefix)
		case false:
			address, err := netip.ParseAddr(item)
			if err != nil {
				resp.Error = function.ConcatFuncErrors(
					function.NewArgumentFuncError(1, "Invalid address"),
					function.NewFuncError("unable to parse address input: "+err.Error()),
				)

				return
			}
			switch {
			case address.Is4():
				prefixes = append(prefixes, netip.PrefixFrom(address, 32))
			case address.Is6():
				prefixes = append(prefixes, netip.PrefixFrom(address, 128))
			}
		}
	}

	// Summarize the prefixes
	summarized := prefixesSummarize(prefixes)

	// Convert back to strings
	result := make([]string, len(summarized))
	for i, p := range summarized {
		result[i] = p.String()
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, result))
}

// prefixesSummarize takes a slice of IP prefixes and returns the minimal list
// that covers the same IP space by merging adjacent and overlapping prefixes.
func prefixesSummarize(prefixes []netip.Prefix) []netip.Prefix {
	if len(prefixes) == 0 {
		return []netip.Prefix{}
	}

	// canonicalize all prefixes to their network address for easier comparison
	for i := range prefixes {
		prefixes[i] = prefixes[i].Masked()
	}

	// sort prefixes for easier comparison
	sort.Slice(prefixes, func(i, j int) bool {
		return prefixes[i].Addr().Less(prefixes[j].Addr())
	})

	for {
		newPrefixes := make([]netip.Prefix, 0, len(prefixes))
		nextMergeInPrevious := false

		for i, prefix := range prefixes {
			switch {
			case nextMergeInPrevious:
				// skip merged item in previous
				nextMergeInPrevious = false
			case i+1 == len(prefixes):
				// last item
				newPrefixes = append(newPrefixes, prefix)
			case prefix.Addr().Is4() != prefixes[i+1].Addr().Is4():
				// don't merge IPv4 and IPv6 prefixes
				newPrefixes = append(newPrefixes, prefix)
			case prefix.Overlaps(prefixes[i+1]):
				nextMergeInPrevious = true
				newPrefixes = append(newPrefixes, prefix)
			case prefix.Bits() == prefixes[i+1].Bits():
				newPrefix := netip.PrefixFrom(prefix.Addr(), prefix.Bits()-1).Masked()
				comparePrefix := netip.PrefixFrom(prefixes[i+1].Addr(), prefix.Bits()-1).Masked()
				if newPrefix == comparePrefix {
					nextMergeInPrevious = true
					newPrefixes = append(newPrefixes, newPrefix)

					continue
				}

				fallthrough
			default:
				newPrefixes = append(newPrefixes, prefix)
			}
		}

		if len(newPrefixes) == len(prefixes) {
			// no more merges possible
			return newPrefixes
		}

		prefixes = newPrefixes
	}
}

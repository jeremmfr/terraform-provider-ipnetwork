package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionSummarize(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input       string
		expectError *regexp.Regexp
		output      []string
	}

	tests := map[string]testCase{
		"adjacent_ipv4": {
			input: `toset([
				"192.0.2.0/24",
				"192.0.3.0/24",
			])`,
			output: []string{
				"192.0.2.0/23",
			},
		},
		"adjacent_ipv6": {
			input: `toset([
				"2001:db8::/64",
				"2001:db8:0:1::/64",
			])`,
			output: []string{
				"2001:db8::/63",
			},
		},
		"adjacent_slash_8": {
			input: `toset([
				"10.0.0.0/8",
				"11.0.0.0/8",
			])`,
			output: []string{
				"10.0.0.0/7",
			},
		},
		"addresses_ipv4": {
			input: `toset([
				"192.0.2.1",
				"192.0.2.2",
			])`,
			output: []string{
				"192.0.2.1/32",
				"192.0.2.2/32",
			},
		},
		"overlap_ipv4": {
			input: `toset([
				"192.0.2.0/24",
				"192.0.2.0/25",
				"192.0.2.128/25",
			])`,
			output: []string{
				"192.0.2.0/24",
			},
		},
		"addresses_ipv6": {
			input: `toset([
				"2001:db8::1",
				"2001:db8::2",
			])`,
			output: []string{
				"2001:db8::1/128",
				"2001:db8::2/128",
			},
		},
		"mixed_families": {
			input: `toset([
				"192.0.2.0/24",
				"192.0.3.0/24",
				"2001:db8::/64",
				"2001:db8:0:1::/64",
			])`,
			output: []string{
				"192.0.2.0/23",
				"2001:db8::/63",
			},
		},
		"four_adjacent": {
			input: `toset([
				"10.0.0.0/24",
				"10.0.1.0/24",
				"10.0.2.0/24",
				"10.0.3.0/24",
			])`,
			output: []string{
				"10.0.0.0/22",
			},
		},
		"duplicate": {
			input: `toset([
				"192.0.2.0/24",
				"192.0.2.0/24",
			])`,
			output: []string{
				"192.0.2.0/24",
			},
		},
		"single_prefix": {
			input: `toset([
				"192.0.2.0/24",
			])`,
			output: []string{
				"192.0.2.0/24",
			},
		},
		"single_address": {
			input: `toset([
				"192.0.2.1",
			])`,
			output: []string{
				"192.0.2.1/32",
			},
		},
		"non_adjacent": {
			input: `toset([
				"192.0.2.0/24",
				"192.0.4.0/24",
			])`,
			output: []string{
				"192.0.2.0/24",
				"192.0.4.0/24",
			},
		},
		"all_private_rfc1918": {
			input: `toset([
				"10.0.0.0/8",
				"172.16.0.0/12",
				"192.168.0.0/16",
			])`,
			output: []string{
				"10.0.0.0/8",
				"172.16.0.0/12",
				"192.168.0.0/16",
			},
		},
		"completely_overlapping": {
			input: `toset([
				"10.0.0.0/16",
				"10.0.0.0/20",
				"10.0.0.0/24",
			])`,
			output: []string{
				"10.0.0.0/16",
			},
		},
		"eight_adjacent_ipv4": {
			input: `toset([
				"10.0.0.0/24",
				"10.0.1.0/24",
				"10.0.2.0/24",
				"10.0.3.0/24",
				"10.0.4.0/24",
				"10.0.5.0/24",
				"10.0.6.0/24",
				"10.0.7.0/24",
			])`,
			output: []string{
				"10.0.0.0/21",
			},
		},
		"eight_adjacent_ipv6": {
			input: `toset([
				"2001:db8::/64",
				"2001:db8:0:1::/64",
				"2001:db8:0:2::/64",
				"2001:db8:0:3::/64",
				"2001:db8:0:4::/64",
				"2001:db8:0:5::/64",
				"2001:db8:0:6::/64",
				"2001:db8:0:7::/64",
			])`,
			output: []string{
				"2001:db8::/61",
			},
		},
		"empty_set": {
			input:  `toset([])`,
			output: []string{},
		},
		"ipv4_mapped_ipv6": {
			input: `toset([
				"::ffff:192.0.2.0/120",
				"::ffff:192.0.3.0/120",
			])`,
			output: []string{
				"::ffff:192.0.2.0/119",
			},
		},
		"ipv6_link_local": {
			input: `toset([
				"fe80::/64",
				"fe80:0:0:1::/64",
			])`,
			output: []string{
				"fe80::/63",
			},
		},
		"max_prefix_ipv4": {
			input: `toset([
				"0.0.0.0/0",
				"192.0.2.0/24",
			])`,
			output: []string{
				"0.0.0.0/0",
			},
		},
		"max_prefix_ipv6": {
			input: `toset([
				"::/0",
				"2001:db8::/32",
			])`,
			output: []string{
				"::/0",
			},
		},
		"mix_slash_31_and_32": {
			input: `toset([
				"10.0.0.0/31",
				"10.0.0.2",
				"10.0.0.3",
			])`,
			output: []string{
				"10.0.0.0/30",
			},
		},
		"mixed_ipv4": {
			input: `toset([
				"192.0.2.1",
				"192.0.2.0/25",
				"192.0.3.0/24",
			])`,
			output: []string{
				"192.0.2.0/25",
				"192.0.3.0/24",
			},
		},
		"multiple_groups": {
			input: `toset([
				"10.0.0.0/24",
				"10.0.1.0/24",
				"10.0.4.0/24",
				"10.0.5.0/24",
				"10.0.8.0/24",
			])`,
			output: []string{
				"10.0.0.0/23",
				"10.0.4.0/23",
				"10.0.8.0/24",
			},
		},
		"reverse_order": {
			input: `toset([
				"192.0.5.0/24",
				"192.0.4.0/24",
				"192.0.3.0/24",
				"192.0.2.0/24",
			])`,
			output: []string{
				"192.0.2.0/23",
				"192.0.4.0/23",
			},
		},
		"subnet_merge": {
			input: `toset([
				"192.168.0.0/25",
				"192.168.0.128/25",
				"192.168.1.0/25",
				"192.168.1.128/25",
			])`,
			output: []string{
				"192.168.0.0/23",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if test.expectError != nil {
				resource.UnitTest(t, resource.TestCase{
					TerraformVersionChecks: []tfversion.TerraformVersionCheck{
						tfversion.SkipBelow(tfversion.Version1_8_0),
					},
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					Steps: []resource.TestStep{
						{
							Config: `
							output "test" {
								value = provider::ipnetwork::summarize(` + test.input + `)
							}
							`,
							ExpectError: test.expectError,
						},
					},
				})
			} else {
				expectedValues := make([]knownvalue.Check, len(test.output))
				for i, v := range test.output {
					expectedValues[i] = knownvalue.StringExact(v)
				}

				resource.UnitTest(t, resource.TestCase{
					TerraformVersionChecks: []tfversion.TerraformVersionCheck{
						tfversion.SkipBelow(tfversion.Version1_8_0),
					},
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					Steps: []resource.TestStep{
						{
							Config: `
							output "test" {
								value = provider::ipnetwork::summarize(` + test.input + `)
							}
							`,
							ConfigStateChecks: []statecheck.StateCheck{
								statecheck.ExpectKnownOutputValue(
									"test",
									knownvalue.SetExact(expectedValues),
								),
							},
						},
					},
				})
			}
		})
	}
}

package provider_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionSort(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input       []string
		expectError *regexp.Regexp
		output      []string
	}

	tests := map[string]testCase{
		"ipv4_unsorted": {
			input: []string{
				"10.0.0.1",
				"2.0.0.1",
				"192.168.0.1",
			},
			output: []string{
				"2.0.0.1",
				"10.0.0.1",
				"192.168.0.1",
			},
		},
		"ipv4_already_sorted": {
			input: []string{
				"10.0.0.1",
				"10.0.0.2",
				"10.0.0.3",
			},
			output: []string{
				"10.0.0.1",
				"10.0.0.2",
				"10.0.0.3",
			},
		},
		"ipv4_cidr": {
			input: []string{
				"10.0.0.0/24",
				"10.0.0.0/8",
				"10.0.0.0/16",
			},
			output: []string{
				"10.0.0.0/8",
				"10.0.0.0/16",
				"10.0.0.0/24",
			},
		},
		"ipv4_mixed_addresses_and_cidr": {
			input: []string{
				"192.168.1.0/24",
				"10.0.0.1",
				"172.16.0.0/12",
			},
			output: []string{
				"10.0.0.1",
				"172.16.0.0/12",
				"192.168.1.0/24",
			},
		},
		"address_before_cidr_same_ip": {
			input: []string{
				"10.0.0.0/24",
				"10.0.0.0/8",
				"10.0.0.0",
			},
			output: []string{
				"10.0.0.0",
				"10.0.0.0/8",
				"10.0.0.0/24",
			},
		},
		"ipv4_cidr_unmasked": {
			input: []string{
				"10.0.0.3/24",
				"10.0.0.2/23",
				"10.0.0.1/25",
			},
			output: []string{
				"10.0.0.1/25",
				"10.0.0.2/23",
				"10.0.0.3/24",
			},
		},
		"ipv6_unsorted": {
			input: []string{
				"2001:db8::3",
				"2001:db8::1",
				"2001:db8::2",
			},
			output: []string{
				"2001:db8::1",
				"2001:db8::2",
				"2001:db8::3",
			},
		},
		"ipv6_cidr": {
			input: []string{
				"2001:db8::/64",
				"2001:db8::/32",
				"2001:db8::/48",
			},
			output: []string{
				"2001:db8::/32",
				"2001:db8::/48",
				"2001:db8::/64",
			},
		},
		"single_element": {
			input: []string{
				"10.0.0.1",
			},
			output: []string{
				"10.0.0.1",
			},
		},
		"empty_list": {
			input:  []string{},
			output: []string{},
		},
		"invalid_address": {
			input: []string{
				"10.0.0.1",
				"invalid",
			},
			expectError: regexp.MustCompile("Invalid address"),
		},
		"invalid_cidr": {
			input: []string{
				"10.0.0.1",
				"10.0.0.0/99",
			},
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"address_before_cidrs_mixed_families": {
			input: []string{
				"192.0.2.2",
				"192.0.2.1/24",
				"192.0.2.1",
				"192.0.2.2/24",
				"192.0.2.2/25",
				"::10",
			},
			output: []string{
				"192.0.2.1",
				"192.0.2.1/24",
				"192.0.2.2",
				"192.0.2.2/24",
				"192.0.2.2/25",
				"::10",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			quotedInput := make([]string, len(test.input))
			for i, v := range test.input {
				quotedInput[i] = fmt.Sprintf("%q", v)
			}

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
								value = provider::ipnetwork::sort([` + strings.Join(quotedInput, ", ") + `])
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
								value = provider::ipnetwork::sort([` + strings.Join(quotedInput, ", ") + `])
							}
							`,
							ConfigStateChecks: []statecheck.StateCheck{
								statecheck.ExpectKnownOutputValue(
									"test",
									knownvalue.ListExact(expectedValues),
								),
							},
						},
					},
				})
			}
		})
	}
}

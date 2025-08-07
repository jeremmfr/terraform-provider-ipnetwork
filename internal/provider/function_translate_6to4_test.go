package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionTranslate6to4(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input       string
		expectError *regexp.Regexp
		output      string
	}

	tests := map[string]testCase{
		"empty_address": {
			input:       "",
			expectError: regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space_address": {
			input:       " ",
			expectError: regexp.MustCompile("unable to parse address input"),
		},
		"ipv4": {
			input:       "192.0.2.33",
			expectError: regexp.MustCompile("must be an IPv6 address"),
		},
		"valid_64": {
			input:  "2001:db8:122:344:c0:2:2100:0/64",
			output: "192.0.2.33",
		},
		"missing_mask": {
			input:  "2001:db8:122:344::c000:221",
			output: "192.0.2.33",
		},
		"invalid": {
			input:       "2001:db8:::c000:22122/64",
			expectError: regexp.MustCompile("unable to parse address"),
		},
		"valid_32": {
			input:  "2001:db8:c000:221::/32",
			output: "192.0.2.33",
		},
		"valid_40": {
			input:  "2001:db8:1c0:2:21::/40",
			output: "192.0.2.33",
		},
		"valid_48": {
			input:  "2001:db8:122:c000:2:2100::/48",
			output: "192.0.2.33",
		},
		"valid_56": {
			input:  "2001:db8:122:3c0:0:221::/56",
			output: "192.0.2.33",
		},
		"valid_96": {
			input:  "2001:db8:122:344::c000:221/96",
			output: "192.0.2.33",
		},
		"well_known": {
			input:  "64:ff9b::c000:221",
			output: "192.0.2.33",
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
								value = provider::ipnetwork::translate_6to4("` + test.input + `")
							}
							`,
							ExpectError: test.expectError,
						},
					},
				})
			} else {
				resource.UnitTest(t, resource.TestCase{
					TerraformVersionChecks: []tfversion.TerraformVersionCheck{
						tfversion.SkipBelow(tfversion.Version1_8_0),
					},
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					Steps: []resource.TestStep{
						{
							Config: `
							output "test" {
								value = provider::ipnetwork::translate_6to4("` + test.input + `")
							}
							`,
							ConfigStateChecks: []statecheck.StateCheck{
								statecheck.ExpectKnownOutputValue(
									"test",
									knownvalue.StringExact(test.output),
								),
							},
						},
					},
				})
			}
		})
	}
}

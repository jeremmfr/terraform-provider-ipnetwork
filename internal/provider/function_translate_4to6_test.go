package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionTranslate4to6(t *testing.T) {
	t.Parallel()

	type testCase struct {
		inputAddress string
		inputPrefix  string
		expectError  *regexp.Regexp
		output       string
	}

	tests := map[string]testCase{
		"empty_address": {
			inputAddress: "",
			inputPrefix:  "2001:db8:122:344::/64",
			expectError:  regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space_address": {
			inputAddress: " ",
			inputPrefix:  "2001:db8:122:344::/64",
			expectError:  regexp.MustCompile("Invalid address"),
		},
		"empty_prefix": {
			inputAddress: "192.0.2.33",
			inputPrefix:  "",
			expectError:  regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space_prefix": {
			inputAddress: "192.0.2.33",
			inputPrefix:  " ",
			expectError:  regexp.MustCompile("Invalid prefix address"),
		},
		"ipv4_prefix": {
			inputAddress: "192.0.2.33",
			inputPrefix:  "192.0.2.34",
			expectError:  regexp.MustCompile("Invalid prefix address"),
		},
		"ipv6_address": {
			inputAddress: "2001:db8:122:345::",
			inputPrefix:  "2001:db8:122:344::/64",
			expectError:  regexp.MustCompile("Invalid address"),
		},
		"valid_64": {
			inputAddress: "192.0.2.33",
			inputPrefix:  "2001:db8:122:344::/64",
			output:       "2001:db8:122:344:c0:2:2100:0",
		},
		"mask_address": {
			inputAddress: "192.0.2.33/24",
			inputPrefix:  "2001:db8:122:344::/64",
			output:       "2001:db8:122:344:c0:2:2100:0",
		},
		"missing_mask_prefix": {
			inputAddress: "192.0.2.33",
			inputPrefix:  "2001:db8:122:344::",
			output:       "2001:db8:122:344::c000:221",
		},
		"invalid_address": {
			inputAddress: "192.0.2.a",
			inputPrefix:  "2001:db8:122:344::/64",
			expectError:  regexp.MustCompile("Invalid address"),
		},
		"invalid_prefix": {
			inputAddress: "192.0.2.33",
			inputPrefix:  "2001:db8:122:344::g/64",
			expectError:  regexp.MustCompile("Invalid prefix address"),
		},
		"valid_32": {
			inputAddress: "192.0.2.33",
			inputPrefix:  "2001:db8:122:344::/32",
			output:       "2001:db8:c000:221::",
		},
		"valid_40": {
			inputAddress: "192.0.2.33",
			inputPrefix:  "2001:db8:122:344::/40",
			output:       "2001:db8:1c0:2:21::",
		},
		"valid_48": {
			inputAddress: "192.0.2.33",
			inputPrefix:  "2001:db8:122:344::/48",
			output:       "2001:db8:122:c000:2:2100::",
		},
		"valid_56": {
			inputAddress: "192.0.2.33",
			inputPrefix:  "2001:db8:122:344::/56",
			output:       "2001:db8:122:3c0:0:221::",
		},
		"valid_96": {
			inputAddress: "192.0.2.33",
			inputPrefix:  "2001:db8:122:344::/96",
			output:       "2001:db8:122:344::c000:221",
		},
		"well_known": {
			inputAddress: "192.0.2.33",
			inputPrefix:  "64:ff9b::/96",
			output:       "64:ff9b::c000:221",
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
								value = provider::ipnetwork::translate_4to6("` + test.inputAddress + `", "` + test.inputPrefix + `")
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
								value = provider::ipnetwork::translate_4to6("` + test.inputAddress + `", "` + test.inputPrefix + `")
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

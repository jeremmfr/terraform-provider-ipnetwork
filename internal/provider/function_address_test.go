package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionAddress(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input       string
		expectError *regexp.Regexp
		output      string
	}

	tests := map[string]testCase{
		"empty": {
			input:       "",
			expectError: regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space": {
			input:       " ",
			expectError: regexp.MustCompile(`String only with space character\(s\)`),
		},
		"valid_ipv4": {
			input:  "192.0.2.1",
			output: "192.0.2.1",
		},
		"invalid_ipv4": {
			input:       "192.0.2.a",
			expectError: regexp.MustCompile("Invalid address"),
		},
		"address_cidr_ipv4": {
			input:  "192.0.2.2/24",
			output: "192.0.2.2",
		},
		"address_0_ipv4": {
			input:  "0",
			output: "0.0.0.0",
		},
		"missing_ipv4": {
			input:  "10",
			output: "10.0.0.0",
		},
		"missing_ipv4_2": {
			input:  "10.20",
			output: "10.20.0.0",
		},
		"valid_ipv6": {
			input:  "2001:DB8::",
			output: "2001:db8::",
		},
		"valid_ipv6_expanded": {
			input:  "2001:0db8:0000:0000:0000:0000:0000:0000",
			output: "2001:db8::",
		},
		"valid_ipv6_short_expanded": {
			input:  "2001:0db8:0:0:0:0:0:0",
			output: "2001:db8::",
		},
		"invalid_ipv6": {
			input:       "2001:db8::h",
			expectError: regexp.MustCompile("Invalid address"),
		},
		"address_cidr_ipv6": {
			input:  "2001:db8::1/64",
			output: "2001:db8::1",
		},
		"address_0_ipv6": {
			input:  "::",
			output: "::",
		},
		"address_scoped": {
			input:  "fe80::1cc0:3e8c:119f:c2e1%ens18",
			output: "fe80::1cc0:3e8c:119f:c2e1",
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
								value = provider::ipnetwork::address("` + test.input + `")
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
								value = provider::ipnetwork::address("` + test.input + `")
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

package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionIs6(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input       string
		expectError *regexp.Regexp
		output      bool
	}

	tests := map[string]testCase{
		"empty": {
			input:       "",
			expectError: regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space": {
			input:       " ",
			expectError: regexp.MustCompile("Invalid address"),
		},
		"valid_ipv4": {
			input:  "192.0.2.1",
			output: false,
		},
		"invalid_ipv4": {
			input:       "192.0.2.a",
			expectError: regexp.MustCompile("Invalid address"),
		},
		"address_cidr_ipv4": {
			input:  "192.0.2.2/24",
			output: false,
		},
		"address_0_ipv4": {
			input:  "0.0.0.0/0",
			output: false,
		},
		"valid_ipv6": {
			input:  "2001:db8::fedc:1234",
			output: true,
		},
		"valid_ipv6_expanded": {
			input:  "2001:0db8:0000:0000:0000:9876:0000:1234",
			output: true,
		},
		"valid_ipv6_short_expanded": {
			input:  "2001:0db8:0:0:4:3:2:1",
			output: true,
		},
		"invalid_ipv6": {
			input:       "2001:db8::h",
			expectError: regexp.MustCompile("Invalid address"),
		},
		"address_cidr_ipv6": {
			input:  "2001:db8::1/64",
			output: true,
		},
		"address_0_ipv6": {
			input:  "::",
			output: true,
		},
		"address_scoped": {
			input:  "fe80::1cc0:3e8c:119f:c2e1%ens18",
			output: true,
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
								value = provider::ipnetwork::is6("` + test.input + `")
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
								value = provider::ipnetwork::is6("` + test.input + `")
							}
							`,
							ConfigStateChecks: []statecheck.StateCheck{
								statecheck.ExpectKnownOutputValue(
									"test",
									knownvalue.Bool(test.output),
								),
							},
						},
					},
				})
			}
		})
	}
}

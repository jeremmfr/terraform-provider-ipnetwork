package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionBits(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input       string
		expectError *regexp.Regexp
		output      int32
	}

	tests := map[string]testCase{
		"empty": {
			input:       "",
			expectError: regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space": {
			input:       " ",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"valid_ipv4": {
			input:  "192.0.3.1/23",
			output: 23,
		},
		"invalid_ipv4": {
			input:       "192.0.2.a/24",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"invalid_ipv4_mask": {
			input:       "192.0.2.2/33",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"address_only_ipv4": {
			input:       "192.0.2.2",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"valid_ipv4_32": {
			input:  "192.0.2.3/32",
			output: 32,
		},
		"valid_ipv4_1": {
			input:  "192.0.3.1/1",
			output: 1,
		},
		"valid_ipv4_0": {
			input:  "192.0.3.1/0",
			output: 0,
		},
		"valid_ipv6": {
			input:  "2001:db8::ffff/64",
			output: 64,
		},
		"valid_ipv6_expanded": {
			input:  "2001:0db8:0000:0000:0000:0000:0000:ffff/64",
			output: 64,
		},
		"valid_ipv6_short_expanded": {
			input:  "2001:0db8:0:0:0:0:0:f/64",
			output: 64,
		},
		"invalid_ipv6": {
			input:       "2001:db8::h/64",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"invalid_ipv6_mask": {
			input:       "2001:db8::/130",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"address_only_ipv6": {
			input:       "2001:db8::1",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"address_scoped": {
			input:       "fe80::1cc0:3e8c:119f:c2e1%ens18/64",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"valid_ipv6_128": {
			input:  "2001:db8::ffff/128",
			output: 128,
		},
		"valid_ipv6_127": {
			input:  "2001:db8::ffff/127",
			output: 127,
		},
		"valid_ipv6_1": {
			input:  "fe80:db8::ffff/1",
			output: 1,
		},
		"valid_ipv6_0": {
			input:  "2001:db8::ffff/0",
			output: 0,
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
								value = provider::ipnetwork::bits("` + test.input + `")
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
								value = provider::ipnetwork::bits("` + test.input + `")
							}
							`,
							ConfigStateChecks: []statecheck.StateCheck{
								statecheck.ExpectKnownOutputValue(
									"test",
									knownvalue.Int32Exact(test.output),
								),
							},
						},
					},
				})
			}
		})
	}
}

package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionCidr(t *testing.T) {
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
			input:  "192.0.2.1/24",
			output: "192.0.2.1/24",
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
			input:  "192.0.2.2",
			output: "192.0.2.2/32",
		},
		"address_only_0_ipv4": {
			input:  "0",
			output: "0.0.0.0/0",
		},
		"missing_ipv4": {
			input:  "10/8",
			output: "10.0.0.0/8",
		},
		"missing_ipv4_2": {
			input:  "10.20/24",
			output: "10.20.0.0/24",
		},
		"decimal_netmask": {
			input:  "192.0.2.3/255.255.255.0",
			output: "192.0.2.3/24",
		},
		"decimal_netmask_0": {
			input:  "192.0.2.3/0.0.0.0",
			output: "192.0.2.3/0",
		},
		"decimal_netmask_invalid": {
			input:       "192.0.2.3/255.254.255.0",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"decimal_netmask_invalid2": {
			input:       "192.0.2.3/255.255.253.0",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"valid_ipv6": {
			input:  "2001:DB8::/32",
			output: "2001:db8::/32",
		},
		"valid_ipv6_expanded": {
			input:  "2001:0db8:0000:0000:0000:0000:0000:0000/48",
			output: "2001:db8::/48",
		},
		"valid_ipv6_short_expanded": {
			input:  "2001:0db8:0:0:0:0:0:0/48",
			output: "2001:db8::/48",
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
			input:  "2001:db8::1",
			output: "2001:db8::1/128",
		},
		"address_only_0_ipv6": {
			input:  "::",
			output: "::/0",
		},
		"address_scoped": {
			input:  "fe80::1cc0:3e8c:119f:c2e1%ens18/64",
			output: "fe80::1cc0:3e8c:119f:c2e1/64",
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
								value = provider::ipnetwork::cidr("` + test.input + `")
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
								value = provider::ipnetwork::cidr("` + test.input + `")
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

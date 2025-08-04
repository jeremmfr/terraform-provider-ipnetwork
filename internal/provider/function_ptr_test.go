package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionPtr(t *testing.T) {
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
			expectError: regexp.MustCompile("Invalid address"),
		},
		"valid_ipv4": {
			input:  "192.0.2.1",
			output: "1.2.0.192.in-addr.arpa.",
		},
		"invalid_ipv4": {
			input:       "192.0.2.a",
			expectError: regexp.MustCompile("Invalid address"),
		},
		"address_cidr_ipv4": {
			input:  "192.0.2.2/24",
			output: "2.2.0.192.in-addr.arpa.",
		},
		"address_0_ipv4": {
			input:  "0.0.0.0/0",
			output: "0.0.0.0.in-addr.arpa.",
		},
		"valid_ipv6": {
			input:  "2001:db8::fedc:1234",
			output: "4.3.2.1.c.d.e.f.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa.",
		},
		"valid_ipv6_expanded": {
			input:  "2001:0db8:0000:0000:0000:9876:0000:1234",
			output: "4.3.2.1.0.0.0.0.6.7.8.9.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa.",
		},
		"valid_ipv6_short_expanded": {
			input:  "2001:0db8:0:0:4:3:2:1",
			output: "1.0.0.0.2.0.0.0.3.0.0.0.4.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa.",
		},
		"invalid_ipv6": {
			input:       "2001:db8::h",
			expectError: regexp.MustCompile("Invalid address"),
		},
		"address_cidr_ipv6": {
			input:  "2001:db8::1/64",
			output: "1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.8.b.d.0.1.0.0.2.ip6.arpa.",
		},
		"address_0_ipv6": {
			input:  "::",
			output: "0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.ip6.arpa.",
		},
		"address_scoped": {
			input:  "fe80::1cc0:3e8c:119f:c2e1%ens18",
			output: "1.e.2.c.f.9.1.1.c.8.e.3.0.c.c.1.0.0.0.0.0.0.0.0.0.0.0.0.0.8.e.f.ip6.arpa.",
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
								value = provider::ipnetwork::ptr("` + test.input + `")
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
								value = provider::ipnetwork::ptr("` + test.input + `")
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

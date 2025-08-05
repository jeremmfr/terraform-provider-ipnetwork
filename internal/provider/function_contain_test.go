package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionEqualContains(t *testing.T) {
	t.Parallel()

	type testCase struct {
		container   string
		address     string
		expectError *regexp.Regexp
		output      bool
	}

	tests := map[string]testCase{
		"empty_container": {
			container:   "",
			address:     "192.0.2.0/24",
			expectError: regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space_container": {
			container:   " ",
			address:     "192.0.2.0/24",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"empty_address": {
			container:   "192.0.2.0/24",
			address:     "",
			expectError: regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space_y": {
			container:   "192.0.2.0/24",
			address:     " ",
			expectError: regexp.MustCompile("Invalid address"),
		},
		"valid_ipv4": {
			container: "192.0.3.1/23",
			address:   "192.0.3.2/24",
			output:    true,
		},
		"valid_ipv4_mask": {
			container: "192.0.3.1/23",
			address:   "192.0.3.2/23",
			output:    true,
		},
		"valid_ipv4_not": {
			container: "192.0.2.1/25",
			address:   "192.0.2.129/25",
			output:    false,
		},
		"valid_ipv4_not_mask": {
			container: "192.0.2.1/25",
			address:   "192.0.2.1/24",
			output:    false,
		},
		"valid_ipv4_addr": {
			container: "192.0.3.1/23",
			address:   "192.0.3.2",
			output:    true,
		},
		"invalid_ipv4_x": {
			container:   "192.0.2.a/24",
			address:     "192.0.2.2/23",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"invalid_ipv4_y": {
			container:   "192.0.2.2/20",
			address:     "192.0.2.b/23",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"invalid_ipv4_y2": {
			container:   "192.0.2.2/20",
			address:     "192.0.2.b",
			expectError: regexp.MustCompile("Invalid address"),
		},
		"address_only_ipv4": {
			container:   "192.0.2.2",
			address:     "192.0.3.2/23",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"valid_ipv4_1": {
			container: "192.0.3.1/1",
			address:   "128.0.3.0/2",
			output:    true,
		},
		"valid_ipv4_0": {
			container: "192.0.3.1/0",
			address:   "1.0.3.0/0",
			output:    true,
		},
		"valid_ipv6": {
			container: "2001:db8::ffff/60",
			address:   "2001:db8::a:ffff/64",
			output:    true,
		},
		"valid_ipv6_not": {
			container: "2001:db8::ffff/64",
			address:   "2001:db8:c::a:ffff/64",
			output:    false,
		},
		"valid_ipv6_mask": {
			container: "2001:db8::ffff/64",
			address:   "2001:db8::b:ffff/64",
			output:    true,
		},
		"valid_ipv6_not_mask": {
			container: "2001:db8::ffff/64",
			address:   "2001:db8::/63",
			output:    false,
		},
		"valid_ipv6_addr": {
			container: "2001:db8::ffff/64",
			address:   "2001:db8::",
			output:    true,
		},
		"valid_ipv6_expanded": {
			container: "2001:0db8:0000:0000:0000:0000:0000:ffff/64",
			address:   "2001:db8::a:ffff/64",
			output:    true,
		},
		"valid_ipv6_short_expanded": {
			container: "2001:0db8:0:0:0:0:0:f/64",
			address:   "2001:db8::a:ffff/64",
			output:    true,
		},
		"invalid_ipv6": {
			container:   "2001:db8::h/64",
			address:     "2001:db8::a:ffff/64",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"address_only_ipv6": {
			container:   "2001:db8::1",
			address:     "2001:db8::a:ffff/64",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"address_scoped": {
			container:   "fe80::1cc0:3e8c:119f:c2e1%ens18/64",
			address:     "2001:db8::a:ffff/64",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"valid_ipv6_1": {
			container: "fe80:db8::ffff/1",
			address:   "8000::ffff/1",
			output:    true,
		},
		"valid_ipv6_0": {
			container: "2001:db8::ffff/0",
			address:   "db8::",
			output:    true,
		},
		"ipv4_ipv6": {
			container: "2001:db8::ffff/64",
			address:   "192.0.2.0/24",
			output:    false,
		},
		"ipv4_in_ipv6": {
			container: "::ffff:c000:0200/96",
			address:   "192.0.2.0",
			output:    false,
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
								value = provider::ipnetwork::contain("` + test.container + `","` + test.address + `")
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
								value = provider::ipnetwork::contain("` + test.container + `","` + test.address + `")
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

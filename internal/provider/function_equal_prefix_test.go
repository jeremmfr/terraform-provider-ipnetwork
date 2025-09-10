package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionEqualPrefix(t *testing.T) {
	t.Parallel()

	type testCase struct {
		addressX    string
		addressY    string
		expectError *regexp.Regexp
		output      bool
	}

	tests := map[string]testCase{
		"empty_x": {
			addressX:    "",
			addressY:    "192.0.2.0/24",
			expectError: regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space_x": {
			addressX:    " ",
			addressY:    "192.0.2.0/24",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"empty_y": {
			addressX:    "192.0.2.0/24",
			addressY:    "",
			expectError: regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space_y": {
			addressX:    "192.0.2.0/24",
			addressY:    " ",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"valid_ipv4": {
			addressX: "192.0.3.1/23",
			addressY: "192.0.2.2/23",
			output:   true,
		},
		"valid_ipv4_not": {
			addressX: "192.0.2.1/25",
			addressY: "192.0.2.129/25",
			output:   false,
		},
		"valid_ipv4_not_mask": {
			addressX: "192.0.2.1/25",
			addressY: "192.0.2.1/24",
			output:   false,
		},
		"invalid_ipv4_x": {
			addressX:    "192.0.2.a/24",
			addressY:    "192.0.2.2/23",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"invalid_ipv4_y": {
			addressX:    "192.0.2.2/24",
			addressY:    "192.0.2.b/23",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"address_only_ipv4": {
			addressX:    "192.0.2.2",
			addressY:    "192.0.3.2/23",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"valid_ipv4_1": {
			addressX: "192.0.3.1/1",
			addressY: "128.0.3.0/1",
			output:   true,
		},
		"valid_ipv4_0": {
			addressX: "192.0.3.1/0",
			addressY: "1.0.3.0/0",
			output:   true,
		},
		"valid_ipv6": {
			addressX: "2001:db8::ffff/64",
			addressY: "2001:db8::a:ffff/64",
			output:   true,
		},
		"valid_ipv6_not": {
			addressX: "2001:db8::ffff/64",
			addressY: "2001:db8:c::a:ffff/64",
			output:   false,
		},
		"valid_ipv6_not_mask": {
			addressX: "2001:db8::ffff/64",
			addressY: "2001:db8::ffff/63",
			output:   false,
		},
		"valid_ipv6_expanded": {
			addressX: "2001:0db8:0000:0000:0000:0000:0000:ffff/64",
			addressY: "2001:db8::a:ffff/64",
			output:   true,
		},
		"valid_ipv6_short_expanded": {
			addressX: "2001:0db8:0:0:0:0:0:f/64",
			addressY: "2001:db8::a:ffff/64",
			output:   true,
		},
		"invalid_ipv6": {
			addressX:    "2001:db8::h/64",
			addressY:    "2001:db8::a:ffff/64",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"address_only_ipv6": {
			addressX:    "2001:db8::1",
			addressY:    "2001:db8::a:ffff/64",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"address_scoped": {
			addressX:    "fe80::1cc0:3e8c:119f:c2e1%ens18/64",
			addressY:    "2001:db8::a:ffff/64",
			expectError: regexp.MustCompile("Invalid CIDR address"),
		},
		"valid_ipv6_1": {
			addressX: "fe80:db8::ffff/1",
			addressY: "8000::ffff/1",
			output:   true,
		},
		"valid_ipv6_0": {
			addressX: "2001:db8::ffff/0",
			addressY: "::/0",
			output:   true,
		},
		"ipv4_ipv6": {
			addressX: "192.0.2.0/24",
			addressY: "2001:db8::ffff/64",
			output:   false,
		},
		"ipv4_in_ipv6": {
			addressX: "192.0.2.0/24",
			addressY: "::ffff:c000:0200/96",
			output:   false,
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
								value = provider::ipnetwork::equal_prefix("` + test.addressX + `", "` + test.addressY + `")
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
								value = provider::ipnetwork::equal_prefix("` + test.addressX + `", "` + test.addressY + `")
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

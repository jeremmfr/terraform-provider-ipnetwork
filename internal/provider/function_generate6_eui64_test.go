package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionGenerate6EUI64(t *testing.T) {
	t.Parallel()

	type testCase struct {
		inputPrefix string
		inputMac    string
		expectError *regexp.Regexp
		output      string
	}

	tests := map[string]testCase{
		"empty_prefix": {
			inputPrefix: "",
			inputMac:    "00-00-5E-00-53-00",
			expectError: regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"empty_mac": {
			inputPrefix: "2001:db8::/64",
			inputMac:    "",
			expectError: regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space_prefix": {
			inputPrefix: " ",
			inputMac:    "00-00-5E-00-53-00",
			expectError: regexp.MustCompile(`Invalid prefix`),
		},
		"space_mac": {
			inputPrefix: "2001:db8::/64",
			inputMac:    " ",
			expectError: regexp.MustCompile(`Invalid MAC`),
		},
		"valid": {
			inputPrefix: "2001:db8::",
			inputMac:    "00-00-5E-00-53-00",
			output:      "2001:db8::200:5eff:fe00:5300",
		},
		"invalid_prefix": {
			inputPrefix: "2001:db8::h",
			inputMac:    "00-00-5E-00-53-00",
			expectError: regexp.MustCompile("Invalid prefix"),
		},
		"invalid_mac": {
			inputPrefix: "2001:db8::",
			inputMac:    "00-00-5E-00-53-100",
			expectError: regexp.MustCompile("Invalid MAC"),
		},
		"prefix_cidr": {
			inputPrefix: "2001:db8::1/64",
			inputMac:    "02-00-5E-00-53-00",
			output:      "2001:db8::5eff:fe00:5300",
		},
		"local": {
			inputPrefix: "fe80::",
			inputMac:    "00-00-5E-00-53-00",
			output:      "fe80::200:5eff:fe00:5300",
		},
		"mac_colon": {
			inputPrefix: "fe80::",
			inputMac:    "00:00:5e:00:53:00",
			output:      "fe80::200:5eff:fe00:5300",
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
								value = provider::ipnetwork::generate6_eui64("` + test.inputPrefix + `", "` + test.inputMac + `")
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
								value = provider::ipnetwork::generate6_eui64("` + test.inputPrefix + `", "` + test.inputMac + `")
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

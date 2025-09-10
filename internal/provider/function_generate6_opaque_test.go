package provider_test

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionGenerate6Opaque(t *testing.T) {
	t.Parallel()

	const secretKeyTest = "secret_key-secret_key"
	networkID := "id"
	dadCounter1 := int32(1)

	type testCase struct {
		inputPrefix     string
		inputNetIface   string
		inputNetworkID  *string
		inputDADCounter *int32
		inputSecretKey  string
		expectError     *regexp.Regexp
		output          string
	}

	tests := map[string]testCase{
		"empty_prefix": {
			inputPrefix:    "",
			inputNetIface:  "00-00-5E-00-53-00",
			inputSecretKey: secretKeyTest,
			expectError:    regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"empty_net_iface": {
			inputPrefix:    "2001:db8::/64",
			inputNetIface:  "",
			inputSecretKey: secretKeyTest,
			expectError:    regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space_prefix": {
			inputPrefix:    " ",
			inputNetIface:  "00-00-5E-00-53-00",
			inputSecretKey: secretKeyTest,
			expectError:    regexp.MustCompile(`Invalid Prefix`),
		},
		"too_small_secret_key": {
			inputPrefix:    "2001:db8::/64",
			inputNetIface:  "00-00-5E-00-53-00",
			inputSecretKey: string([]byte(secretKeyTest)[0:8]),
			expectError:    regexp.MustCompile(`Invalid Parameter Value Length`),
		},
		"valid": {
			inputPrefix:    "2001:db8::",
			inputNetIface:  "00-00-5E-00-53-00",
			inputSecretKey: secretKeyTest,
			output:         "2001:db8::e919:9c5c:f8ab:26e2",
		},
		"invalid_prefix": {
			inputPrefix:    "2001:db8::h",
			inputNetIface:  "00-00-5E-00-53-00",
			inputSecretKey: secretKeyTest,
			expectError:    regexp.MustCompile("Invalid Prefix"),
		},
		"valid+1": {
			inputPrefix:     "2001:db8::",
			inputNetIface:   "00-00-5E-00-53-00",
			inputDADCounter: &dadCounter1,
			inputSecretKey:  secretKeyTest,
			output:          "2001:db8::a2b0:4703:b842:e2fa",
		},
		"valid_network_id": {
			inputPrefix:    "2001:db8::",
			inputNetIface:  "00-00-5E-00-53-00",
			inputNetworkID: &networkID,
			inputSecretKey: secretKeyTest,
			output:         "2001:db8::3f76:53c0:c91:29b4",
		},
		"prefix_cidr": {
			inputPrefix:    "2001:db8::1/64",
			inputNetIface:  "00-00-5E-00-53-00",
			inputSecretKey: secretKeyTest,
			output:         "2001:db8::e919:9c5c:f8ab:26e2",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			arguments := `"` + test.inputPrefix + `", "` + test.inputNetIface + `"`
			if test.inputNetworkID != nil {
				arguments += `, "` + *test.inputNetworkID + `"`
			} else {
				arguments += `, null`
			}
			if test.inputDADCounter != nil {
				arguments += `, ` + strconv.Itoa(int(*test.inputDADCounter))
			} else {
				arguments += `, null`
			}
			arguments += `, "` + test.inputSecretKey + `"`

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
								value = provider::ipnetwork::generate6_opaque(` + arguments + `)
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
								value = provider::ipnetwork::generate6_opaque(` + arguments + `)
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

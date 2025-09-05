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

func TestAccFunctionAddressPort(t *testing.T) {
	t.Parallel()

	type testCase struct {
		inputAddress string
		inputPort    int32
		expectError  *regexp.Regexp
		output       string
	}

	tests := map[string]testCase{
		"empty": {
			inputAddress: "",
			inputPort:    0,
			expectError:  regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"space": {
			inputAddress: " ",
			inputPort:    0,
			expectError:  regexp.MustCompile("Invalid address"),
		},
		"valid_ipv4": {
			inputAddress: "192.0.2.1",
			inputPort:    80,
			output:       "192.0.2.1:80",
		},
		"invalid_ipv4": {
			inputAddress: "192.0.2.a",
			inputPort:    80,
			expectError:  regexp.MustCompile("Invalid address"),
		},
		"invalid_port": {
			inputAddress: "192.0.2.1",
			inputPort:    -80,
			expectError:  regexp.MustCompile("Invalid Parameter Value"),
		},
		"invalid_port2": {
			inputAddress: "192.0.2.1",
			inputPort:    88888,
			expectError:  regexp.MustCompile("Invalid Parameter Value"),
		},
		"address_cidr_ipv4": {
			inputAddress: "192.0.2.2/24",
			inputPort:    443,
			output:       "192.0.2.2:443",
		},
		"valid_ipv6": {
			inputAddress: "2001:db8::fedc:1234",
			inputPort:    22,
			output:       "[2001:db8::fedc:1234]:22",
		},
		"valid_ipv6_expanded": {
			inputAddress: "2001:0db8:0000:0000:0000:9876:0000:1234",
			inputPort:    5000,
			output:       "[2001:db8::9876:0:1234]:5000",
		},
		"valid_ipv6_short_expanded": {
			inputAddress: "2001:0db8:0:0:4:3:2:1",
			inputPort:    65000,
			output:       "[2001:db8::4:3:2:1]:65000",
		},
		"invalid_ipv6": {
			inputAddress: "2001:db8::h",
			inputPort:    80,
			expectError:  regexp.MustCompile("Invalid address"),
		},
		"address_cidr_ipv6": {
			inputAddress: "2001:db8::1/64",
			inputPort:    444,
			output:       "[2001:db8::1]:444",
		},
		"address_scoped": {
			inputAddress: "fe80::1cc0:3e8c:119f:c2e1%ens18",
			inputPort:    445,
			output:       "[fe80::1cc0:3e8c:119f:c2e1]:445",
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
								value = provider::ipnetwork::address_port("` +
								test.inputAddress + `", ` + strconv.FormatInt(int64(test.inputPort), 10) +
								`)
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
								value = provider::ipnetwork::address_port("` +
								test.inputAddress + `", ` + strconv.FormatInt(int64(test.inputPort), 10) +
								`)
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

package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionIsPrivateRFC6598(t *testing.T) {
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
		"invalid_address": {
			input:       "192.0.2.a",
			expectError: regexp.MustCompile("Invalid address"),
		},
		"public_ipv4": {
			input:  "8.8.8.8",
			output: false,
		},
		"public_ipv4_cloudflare": {
			input:  "1.0.0.1",
			output: false,
		},
		"private_ipv4_10": {
			input:  "10.0.0.1",
			output: false, // Not RFC6598
		},
		"private_ipv4_172": {
			input:  "172.16.0.1",
			output: false, // Not RFC6598
		},
		"private_ipv4_192": {
			input:  "192.168.1.1",
			output: false, // Not RFC6598
		},
		"cgn_ipv4_start": {
			input:  "100.64.0.0",
			output: true,
		},
		"cgn_ipv4_100_64_1": {
			input:  "100.64.0.1",
			output: true,
		},
		"cgn_ipv4_mid_100_100": {
			input:  "100.100.50.25",
			output: true,
		},
		"cgn_ipv4_end_100_127": {
			input:  "100.127.255.254",
			output: true,
		},
		"cgn_ipv4_end_100_127_255": {
			input:  "100.127.255.255",
			output: true,
		},
		"not_cgn_100_63": {
			input:  "100.63.255.255",
			output: false, // Just before CGN range
		},
		"not_cgn_100_128": {
			input:  "100.128.0.0",
			output: false, // Just after CGN range
		},
		"not_cgn_100_0": {
			input:  "100.0.0.1",
			output: false, // 100.x but not in 100.64-127 range
		},
		"not_cgn_100_200": {
			input:  "100.200.0.1",
			output: false, // 100.x but not in 100.64-127 range
		},
		"loopback_ipv4": {
			input:  "127.0.0.1",
			output: false,
		},
		"unspecified_ipv4": {
			input:  "0.0.0.0",
			output: false,
		},
		"link_local_ipv4": {
			input:  "169.254.1.1",
			output: false,
		},
		"multicast_ipv4": {
			input:  "224.0.0.1",
			output: false,
		},
		"public_ipv6_google": {
			input:  "2001:4860:4860::8888",
			output: false,
		},
		"private_ipv6_ula": {
			input:  "fc00::1",
			output: false,
		},
		"ipv4_mapped_cgn": {
			input:  "::ffff:100.64.0.1",
			output: true,
		},
		"prefix_ipv4_public_24": {
			input:  "1.1.1.0/24",
			output: false,
		},
		"prefix_ipv4_private_10_8": {
			input:  "10.0.0.0/8",
			output: false, // Not RFC6598
		},
		"prefix_ipv4_cgn_10": {
			input:  "100.64.0.0/10",
			output: true,
		},
		"prefix_ipv4_cgn_16": {
			input:  "100.64.0.0/16",
			output: true,
		},
		"prefix_ipv4_cgn_24": {
			input:  "100.100.50.0/24",
			output: true,
		},
		"prefix_ipv4_cgn_100_127_24": {
			input:  "100.127.0.0/24",
			output: true,
		},
		"prefix_ipv4_cgn_too_broad": {
			input:  "100.64.0.0/9",
			output: false, // Too broad, extends beyond 100.64.0.0/10
		},
		"prefix_ipv4_not_cgn_100_0": {
			input:  "100.0.0.0/16",
			output: false, // Not in CGN range
		},
		"prefix_ipv4_not_cgn_100_128": {
			input:  "100.128.0.0/16",
			output: false, // Not in CGN range
		},
		"prefix_ipv4_contains_cgn": {
			input:  "100.0.0.0/8",
			output: false, // Contains but not entirely within CGN
		},
		"prefix_ipv6_public_48": {
			input:  "2001:4860::/48",
			output: false, // IPv6 not RFC6598
		},
		"prefix_ipv6_ula_7": {
			input:  "fc00::/7",
			output: false, // IPv6 not RFC6598
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
								value = provider::ipnetwork::is_private_rfc6598("` + test.input + `")
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
								value = provider::ipnetwork::is_private_rfc6598("` + test.input + `")
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

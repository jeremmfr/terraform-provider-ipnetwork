package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionIsPrivateRFC1918(t *testing.T) {
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
			output: true,
		},
		"private_ipv4_10_end": {
			input:  "10.255.255.254",
			output: true,
		},
		"private_ipv4_172": {
			input:  "172.16.0.1",
			output: true,
		},
		"private_ipv4_172_mid": {
			input:  "172.20.5.10",
			output: true,
		},
		"private_ipv4_172_end": {
			input:  "172.31.255.254",
			output: true,
		},
		"private_ipv4_192": {
			input:  "192.168.1.1",
			output: true,
		},
		"private_ipv4_192_end": {
			input:  "192.168.255.254",
			output: true,
		},
		"not_private_ipv4_172_before": {
			input:  "172.15.255.254",
			output: false,
		},
		"not_private_ipv4_172_after": {
			input:  "172.32.0.1",
			output: false,
		},
		"private_ipv4_cgn_rfc6598": {
			input:  "100.64.0.1",
			output: false, // Not RFC1918
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
		"broadcast_ipv4": {
			input:  "255.255.255.255",
			output: false,
		},
		"documentation_testnet1": {
			input:  "192.0.2.1",
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
		"ipv4_mapped_private": {
			input:  "::ffff:192.168.1.1",
			output: true,
		},
		"prefix_ipv4_public_24": {
			input:  "1.1.1.0/24",
			output: false,
		},
		"prefix_ipv4_private_10_8": {
			input:  "10.0.0.0/8",
			output: true,
		},
		"prefix_ipv4_private_10_24": {
			input:  "10.5.10.0/24",
			output: true,
		},
		"prefix_ipv4_private_172_12": {
			input:  "172.16.0.0/12",
			output: true,
		},
		"prefix_ipv4_private_172_24": {
			input:  "172.20.5.0/24",
			output: true,
		},
		"prefix_ipv4_private_192_16": {
			input:  "192.168.0.0/16",
			output: true,
		},
		"prefix_ipv4_private_192_24": {
			input:  "192.168.1.0/24",
			output: true,
		},
		"prefix_ipv4_cgn_10": {
			input:  "100.64.0.0/10",
			output: false, // Not RFC1918
		},
		"prefix_ipv4_contains_private": {
			input:  "192.0.0.0/8",
			output: false, // Contains but not entirely within RFC1918
		},
		"prefix_ipv4_10_too_broad": {
			input:  "10.0.0.0/7",
			output: false, // Too broad, extends beyond 10.0.0.0/8
		},
		"prefix_ipv4_172_too_broad": {
			input:  "172.16.0.0/11",
			output: false, // Too broad, extends beyond 172.16.0.0/12
		},
		"prefix_ipv4_192_too_broad": {
			input:  "192.168.0.0/15",
			output: false, // Too broad, extends beyond 192.168.0.0/16
		},
		"prefix_ipv6_public_48": {
			input:  "2001:4860::/48",
			output: false, // IPv6 not RFC1918
		},
		"prefix_ipv6_ula_7": {
			input:  "fc00::/7",
			output: false, // IPv6 not RFC1918
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
								value = provider::ipnetwork::is_private_rfc1918("` + test.input + `")
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
								value = provider::ipnetwork::is_private_rfc1918("` + test.input + `")
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

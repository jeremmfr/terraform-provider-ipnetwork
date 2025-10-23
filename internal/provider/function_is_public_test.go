package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionIsPublic(t *testing.T) {
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
			output: true,
		},
		"public_ipv4_cidr": {
			input:  "1.1.1.1/24",
			output: true,
		},
		"public_ipv4_cloudflare": {
			input:  "1.0.0.1",
			output: true,
		},
		"private_ipv4_10": {
			input:  "10.0.0.1",
			output: false,
		},
		"private_ipv4_cgn_rfc6598": {
			input:  "100.64.0.1",
			output: false,
		},
		"private_ipv4_cgn_rfc6598_end": {
			input:  "100.127.255.254",
			output: false,
		},
		"private_ipv4_172": {
			input:  "172.16.0.1",
			output: false,
		},
		"private_ipv4_192": {
			input:  "192.168.1.1",
			output: false,
		},
		"private_ipv4_192_cidr": {
			input:  "192.168.1.1/24",
			output: false,
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
		"this_network_ipv4": {
			input:  "0.1.2.3",
			output: false,
		},
		"ietf_protocol_ipv4": {
			input:  "192.0.0.1",
			output: false,
		},
		"documentation_testnet1": {
			input:  "192.0.2.1",
			output: false,
		},
		"documentation_testnet2": {
			input:  "198.51.100.1",
			output: false,
		},
		"documentation_testnet3": {
			input:  "203.0.113.1",
			output: false,
		},
		"reserved_240": {
			input:  "240.0.0.1",
			output: false,
		},
		"reserved_255": {
			input:  "255.0.0.1",
			output: false,
		},
		"public_ipv6_google": {
			input:  "2001:4860:4860::8888",
			output: true,
		},
		"public_ipv6_cloudflare": {
			input:  "2606:4700:4700::1111",
			output: true,
		},
		"public_ipv6_cidr": {
			input:  "2001:4860:4860::8888/64",
			output: true,
		},
		"private_ipv6_ula": {
			input:  "fd00::1",
			output: false,
		},
		"private_ipv6_ula_expanded": {
			input:  "fd12:3456:789a:bcde::1",
			output: false,
		},
		"loopback_ipv6": {
			input:  "::1",
			output: false,
		},
		"unspecified_ipv6": {
			input:  "::",
			output: false,
		},
		"link_local_ipv6": {
			input:  "fe80::1",
			output: false,
		},
		"link_local_ipv6_scoped": {
			input:  "fe80::1cc0:3e8c:119f:c2e1%eth0",
			output: false,
		},
		"multicast_ipv6": {
			input:  "ff02::1",
			output: false,
		},
		"multicast_ipv6_interface_local": {
			input:  "ff01::1",
			output: false,
		},
		"multicast_ipv6_link_local": {
			input:  "ff02::2",
			output: false,
		},
		"benchmarking_ipv6_2001_2": {
			input:  "2001:2::1",
			output: false,
		},
		"benchmarking_ipv6_2001_2_end": {
			input:  "2001:2:0:ffff:ffff:ffff:ffff:ffff",
			output: false,
		},
		"documentation_ipv6_db8": {
			input:  "2001:db8::1",
			output: false,
		},
		"documentation_ipv6_db8_end": {
			input:  "2001:db8:ffff:ffff:ffff:ffff:ffff:ffff",
			output: false,
		},
		"documentation_ipv6_3fff": {
			input:  "3fff::1",
			output: false,
		},
		"documentation_ipv6_3fff:fff": {
			input:  "3fff:fff::1",
			output: false,
		},
		"adjacent_to_documentation_ipv6_3fff": {
			input:  "3fff:ffff::1",
			output: true,
		},
		"ipv4_mapped_public": {
			input:  "::ffff:8.8.8.8",
			output: true,
		},
		"ipv4_mapped_private": {
			input:  "::ffff:192.168.1.1",
			output: false,
		},
		"prefix_ipv4_public_24": {
			input:  "1.1.1.0/24",
			output: true,
		},
		"prefix_ipv4_public_20": {
			input:  "8.8.0.0/20",
			output: true,
		},
		"prefix_ipv4_private_10_8": {
			input:  "10.0.0.0/8",
			output: false,
		},
		"prefix_ipv4_private_192_16": {
			input:  "192.168.0.0/16",
			output: false,
		},
		"prefix_ipv4_contains_private": {
			input:  "192.0.0.0/20",
			output: false,
		},
		"prefix_ipv4_contains_testnet1": {
			input:  "192.0.0.0/22",
			output: false,
		},
		"prefix_ipv4_testnet2_24": {
			input:  "198.51.100.0/24",
			output: false,
		},
		"prefix_ipv4_cgn_10": {
			input:  "100.64.0.0/10",
			output: false,
		},
		"prefix_ipv4_multicast_4": {
			input:  "224.0.0.0/4",
			output: false,
		},
		"prefix_ipv6_public_48": {
			input:  "2001:4860::/48",
			output: true,
		},
		"prefix_ipv6_public_32": {
			input:  "2606:4700::/32",
			output: true,
		},
		"prefix_ipv6_ula_7": {
			input:  "fc00::/7",
			output: false,
		},
		"prefix_ipv6_ula_48": {
			input:  "fd12:3456:789a::/48",
			output: false,
		},
		"prefix_ipv6_link_local_10": {
			input:  "fe80::/10",
			output: false,
		},
		"prefix_ipv6_benchmarking_2001_2": {
			input:  "2001:2::/48",
			output: false,
		},
		"prefix_ipv6_benchmarking_2001_2_subnet": {
			input:  "2001:2:0:abcd::/64",
			output: false,
		},
		"prefix_ipv6_documentation_db8_32": {
			input:  "2001:db8::/32",
			output: false,
		},
		"prefix_ipv6_contains_db8": {
			input:  "2001::/16",
			output: false,
		},
		"prefix_ipv6_multicast_8": {
			input:  "ff00::/8",
			output: false,
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
								value = provider::ipnetwork::is_public("` + test.input + `")
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
								value = provider::ipnetwork::is_public("` + test.input + `")
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

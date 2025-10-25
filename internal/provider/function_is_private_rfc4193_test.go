package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionIsPrivateRFC4193(t *testing.T) {
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
			output: false, // IPv4 not RFC4193
		},
		"private_ipv4_10": {
			input:  "10.0.0.1",
			output: false, // IPv4 not RFC4193
		},
		"public_ipv6_google": {
			input:  "2001:4860:4860::8888",
			output: false,
		},
		"public_ipv6_cloudflare": {
			input:  "2606:4700:4700::1111",
			output: false,
		},
		"private_ipv6_ula_fc00": {
			input:  "fc00::1",
			output: true,
		},
		"private_ipv6_ula_fd00": {
			input:  "fd00::1",
			output: true,
		},
		"private_ipv6_ula_fd_expanded": {
			input:  "fd12:3456:789a:bcde::1",
			output: true,
		},
		"private_ipv6_ula_fc_end": {
			input:  "fcff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			output: true,
		},
		"private_ipv6_ula_fd_end": {
			input:  "fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			output: true,
		},
		"not_ula_fe00": {
			input:  "fe00::1",
			output: false,
		},
		"not_ula_fb00": {
			input:  "fb00::1",
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
		"multicast_ipv6": {
			input:  "ff02::1",
			output: false,
		},
		"documentation_ipv6_db8": {
			input:  "2001:db8::1",
			output: false,
		},
		"ipv4_mapped_public": {
			input:  "::ffff:8.8.8.8",
			output: false, // IPv4-mapped not ULA
		},
		"ipv4_mapped_private": {
			input:  "::ffff:192.168.1.1",
			output: false, // IPv4-mapped not ULA
		},
		"prefix_ipv4_private_10": {
			input:  "10.0.0.0/8",
			output: false, // IPv4 not RFC4193
		},
		"prefix_ipv6_public_48": {
			input:  "2001:4860::/48",
			output: false,
		},
		"prefix_ipv6_ula_fc_7": {
			input:  "fc00::/7",
			output: true,
		},
		"prefix_ipv6_ula_fc_8": {
			input:  "fc00::/8",
			output: true,
		},
		"prefix_ipv6_ula_fd_8": {
			input:  "fd00::/8",
			output: true,
		},
		"prefix_ipv6_ula_fc_48": {
			input:  "fc12:3456:789a::/48",
			output: true,
		},
		"prefix_ipv6_ula_fd_48": {
			input:  "fd12:3456:789a::/48",
			output: true,
		},
		"prefix_ipv6_ula_fc_64": {
			input:  "fc00:1234:5678:abcd::/64",
			output: true,
		},
		"prefix_ipv6_link_local_10": {
			input:  "fe80::/10",
			output: false,
		},
		"prefix_ipv6_documentation_db8_32": {
			input:  "2001:db8::/32",
			output: false,
		},
		"prefix_ipv6_multicast_8": {
			input:  "ff00::/8",
			output: false,
		},
		"prefix_ipv6_ula_too_broad": {
			input:  "fc00::/6",
			output: false, // Too broad, extends beyond fc00::/7
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
								value = provider::ipnetwork::is_private_rfc4193("` + test.input + `")
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
								value = provider::ipnetwork::is_private_rfc4193("` + test.input + `")
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

package provider_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFunctionRangeToPrefixes(t *testing.T) {
	t.Parallel()

	type testCase struct {
		start       string
		end         string
		expectError *regexp.Regexp
		output      []string
	}

	tests := map[string]testCase{
		"aligned_ipv4_24": {
			start: "10.0.0.0",
			end:   "10.0.0.255",
			output: []string{
				"10.0.0.0/24",
			},
		},
		"single_address_ipv4": {
			start: "10.0.0.1",
			end:   "10.0.0.1",
			output: []string{
				"10.0.0.1/32",
			},
		},
		"unaligned_small_range": {
			start: "10.0.0.5",
			end:   "10.0.0.20",
			output: []string{
				"10.0.0.5/32",
				"10.0.0.6/31",
				"10.0.0.8/29",
				"10.0.0.16/30",
				"10.0.0.20/32",
			},
		},
		"two_slash_24": {
			start: "10.0.0.0",
			end:   "10.0.1.255",
			output: []string{
				"10.0.0.0/23",
			},
		},
		"aligned_ipv6_120": {
			start: "2001:db8::",
			end:   "2001:db8::ff",
			output: []string{
				"2001:db8::/120",
			},
		},
		"single_address_ipv6": {
			start: "2001:db8::1",
			end:   "2001:db8::1",
			output: []string{
				"2001:db8::1/128",
			},
		},
		"unaligned_ipv6": {
			start: "2001:db8::1",
			end:   "2001:db8::6",
			output: []string{
				"2001:db8::1/128",
				"2001:db8::2/127",
				"2001:db8::4/127",
				"2001:db8::6/128",
			},
		},
		"full_ipv4_range": {
			start: "0.0.0.0",
			end:   "255.255.255.255",
			output: []string{
				"0.0.0.0/0",
			},
		},
		"mixed_families": {
			start:       "10.0.0.0",
			end:         "2001:db8::ff",
			expectError: regexp.MustCompile(`start and\s+end addresses must be the same IP version`),
		},
		"mixed_families_reverse": {
			start:       "2001:db8::1",
			end:         "10.0.0.1",
			expectError: regexp.MustCompile(`start and\s+end addresses must be the same IP version`),
		},
		"start_greater_than_end": {
			start:       "10.0.0.10",
			end:         "10.0.0.1",
			expectError: regexp.MustCompile(`start\s+address must be less than or equal to end address`),
		},
		"invalid_start": {
			start:       "invalid",
			end:         "10.0.0.1",
			expectError: regexp.MustCompile("Invalid address"),
		},
		"invalid_end": {
			start:       "10.0.0.1",
			end:         "invalid",
			expectError: regexp.MustCompile("Invalid address"),
		},
		"empty_start": {
			start:       "",
			end:         "10.0.0.1",
			expectError: regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"empty_end": {
			start:       "10.0.0.1",
			end:         "",
			expectError: regexp.MustCompile("Invalid Parameter Value Length"),
		},
		"three_addresses": {
			start: "10.0.0.1",
			end:   "10.0.0.3",
			output: []string{
				"10.0.0.1/32",
				"10.0.0.2/31",
			},
		},
		"cross_boundary": {
			start: "10.0.0.254",
			end:   "10.0.1.1",
			output: []string{
				"10.0.0.254/31",
				"10.0.1.0/31",
			},
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
								value = provider::ipnetwork::range_to_prefixes("` + test.start + `", "` + test.end + `")
							}
							`,
							ExpectError: test.expectError,
						},
					},
				})
			} else {
				expectedValues := make([]knownvalue.Check, len(test.output))
				for i, v := range test.output {
					expectedValues[i] = knownvalue.StringExact(v)
				}

				resource.UnitTest(t, resource.TestCase{
					TerraformVersionChecks: []tfversion.TerraformVersionCheck{
						tfversion.SkipBelow(tfversion.Version1_8_0),
					},
					ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
					Steps: []resource.TestStep{
						{
							Config: `
							output "test" {
								value = provider::ipnetwork::range_to_prefixes("` + test.start + `", "` + test.end + `")
							}
							`,
							ConfigStateChecks: []statecheck.StateCheck{
								statecheck.ExpectKnownOutputValue(
									"test",
									knownvalue.ListExact(expectedValues),
								),
							},
						},
					},
				})
			}
		})
	}
}

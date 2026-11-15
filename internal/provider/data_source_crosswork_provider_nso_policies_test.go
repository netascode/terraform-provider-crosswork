// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package provider

// Section below is generated&owned by "gen/generator.go". //template:begin imports
import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSource
func TestAccDataSourceCwProviderNSOPolicies(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "provider_name", "nso"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "match", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "match_rule", "*"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "onboard_to_nso", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "onboard_to_nso_rule", "*"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "onboard_from_nso", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "onboard_from_nso_rule", "*"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "sync_from_nso", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "sync_from_nso_rule", "*"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "check_sync", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "check_sync_rule", "*"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "neds.0.ned", "cisco-iosxr-cli-7.76"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_provider_nso_policies.test", "neds.0.rule", "product_info.software_type = 'IOS XR*'"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwProviderNSOPoliciesConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwProviderNSOPoliciesConfig() string {
	config := `resource "crosswork_provider_nso_policies" "test" {` + "\n"
	config += `	provider_name = "nso"` + "\n"
	config += `	match = true` + "\n"
	config += `	match_rule = "*"` + "\n"
	config += `	onboard_to_nso = true` + "\n"
	config += `	onboard_to_nso_rule = "*"` + "\n"
	config += `	onboard_from_nso = false` + "\n"
	config += `	onboard_from_nso_rule = "*"` + "\n"
	config += `	sync_from_nso = true` + "\n"
	config += `	sync_from_nso_rule = "*"` + "\n"
	config += `	check_sync = true` + "\n"
	config += `	check_sync_rule = "*"` + "\n"
	config += `	neds = [{` + "\n"
	config += `	  ned = "cisco-iosxr-cli-7.76"` + "\n"
	config += `	  rule = "product_info.software_type = 'IOS XR*'"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_provider_nso_policies" "test" {
			id = crosswork_provider_nso_policies.test.id
			provider_name = "nso"
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

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
func TestAccDataSourceCwL2VPNTagSets(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_tag_sets.test", "name", "cwtest-l2-tags"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_tag_sets.test", "co_flag", "10"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwL2VPNTagSetsConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwL2VPNTagSetsConfig() string {
	config := `resource "crosswork_l2vpn_tag_sets" "test" {` + "\n"
	config += `	name = "cwtest-l2-tags"` + "\n"
	config += `	tag_value = [4812]` + "\n"
	config += `	co_flag = "10"` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_l2vpn_tag_sets" "test" {
			id = crosswork_l2vpn_tag_sets.test.id
			name = "cwtest-l2-tags"
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

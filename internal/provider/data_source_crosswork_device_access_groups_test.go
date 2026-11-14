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
func TestAccDataSourceCwDeviceAccessGroups(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_device_access_groups.test", "name", "cwtest-access-group"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_device_access_groups.test", "description", "Test access group"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwDeviceAccessGroupsPrerequisitesConfig + testAccDataSourceCwDeviceAccessGroupsConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccDataSourceCwDeviceAccessGroupsPrerequisitesConfig = `
data "crosswork_groupings" "root" {
  classifier = "DeviceAccess"
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwDeviceAccessGroupsConfig() string {
	config := `resource "crosswork_device_access_groups" "test" {` + "\n"
	config += `	name = "cwtest-access-group"` + "\n"
	config += `	description = "Test access group"` + "\n"
	config += `	parent_uuid = [for group in data.crosswork_groupings.root.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_device_access_groups" "test" {
			id = crosswork_device_access_groups.test.id
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

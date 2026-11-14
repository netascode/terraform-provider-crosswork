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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc
func TestAccCwDeviceAccessGroups(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_device_access_groups.test", "name", "cwtest-access-group"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_device_access_groups.test", "description", "Test access group"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccCwDeviceAccessGroupsPrerequisitesConfig + testAccCwDeviceAccessGroupsConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccCwDeviceAccessGroupsPrerequisitesConfig + testAccCwDeviceAccessGroupsConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_device_access_groups.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// End of section. //template:end testAcc

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccCwDeviceAccessGroupsPrerequisitesConfig = `
data "crosswork_groupings" "root" {
  classifier = "DeviceAccess"
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwDeviceAccessGroupsConfig_minimum() string {
	config := `resource "crosswork_device_access_groups" "test" {` + "\n"
	config += `	name = "cwtest-access-group"` + "\n"
	config += `	parent_uuid = [for group in data.crosswork_groupings.root.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwDeviceAccessGroupsConfig_all() string {
	config := `resource "crosswork_device_access_groups" "test" {` + "\n"
	config += `	name = "cwtest-access-group"` + "\n"
	config += `	description = "Test access group"` + "\n"
	config += `	parent_uuid = [for group in data.crosswork_groupings.root.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

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
func TestAccCwDeviceGroups(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_device_groups.test", "name", "cwtest-device-group"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_device_groups.test", "description", "Test device group"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccCwDeviceGroupsPrerequisitesConfig + testAccCwDeviceGroupsConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccCwDeviceGroupsPrerequisitesConfig + testAccCwDeviceGroupsConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_device_groups.test",
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
const testAccCwDeviceGroupsPrerequisitesConfig = `
data "crosswork_groupings" "anchor" {
  classifier = "LocationDevices"
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwDeviceGroupsConfig_minimum() string {
	config := `resource "crosswork_device_groups" "test" {` + "\n"
	config += `	name = "cwtest-device-group"` + "\n"
	config += `	parent_uuid = [for group in data.crosswork_groupings.anchor.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwDeviceGroupsConfig_all() string {
	config := `resource "crosswork_device_groups" "test" {` + "\n"
	config += `	name = "cwtest-device-group"` + "\n"
	config += `	description = "Test device group"` + "\n"
	config += `	parent_uuid = [for group in data.crosswork_groupings.anchor.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

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

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc
func TestAccCwAssignDevicesToAccessGroup(t *testing.T) {
	var checks []resource.TestCheckFunc

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwAssignDevicesToAccessGroupPrerequisitesConfig + testAccCwAssignDevicesToAccessGroupConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_assign_devices_to_access_group.test",
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
const testAccCwAssignDevicesToAccessGroupPrerequisitesConfig = `
data "crosswork_groupings" "assign_root" {
  classifier = "DeviceAccess"
}

resource "crosswork_device_access_groups" "assign_test" {
  name        = "cwtest-gen-assign-group"
  parent_uuid = [for group in data.crosswork_groupings.assign_root.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]
}

resource "crosswork_credential_profiles" "assign_test" {
  profile = "cwtest-gen-assign-profile"
  credentials = [
    {
      type      = "ROBOT_USERPASS_SSH"
      user_name = "cwtest"
      password  = "Cwtest-Password-1"
    }
  ]
}

resource "crosswork_network_devices" "assign_test" {
  name               = "cwtest-gen-assign-device"
  host_name          = "cwtest-gen-assign-device"
  credential_profile = crosswork_credential_profiles.assign_test.id
  software_type      = "IOS XR"
  capabilities       = ["YANG_CLI"]
  connectivity = [{
    protocol      = "ROBOT_MSVC_TRANS_SSH"
    ip_address    = "192.0.2.241"
    prefix_length = 24
    port          = 22
  }]

  # Crosswork can move the admin state itself: an embedded data gateway pool attaches a new device and brings it up
  lifecycle {
    ignore_changes = [admin_state]
  }
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwAssignDevicesToAccessGroupConfig_minimum() string {
	config := `resource "crosswork_assign_devices_to_access_group" "test" {` + "\n"
	config += `	group_id = crosswork_device_access_groups.assign_test.id` + "\n"
	config += `	devices = [crosswork_network_devices.assign_test.id]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwAssignDevicesToAccessGroupConfig_all() string {
	config := `resource "crosswork_assign_devices_to_access_group" "test" {` + "\n"
	config += `	group_id = crosswork_device_access_groups.assign_test.id` + "\n"
	config += `	devices = [crosswork_network_devices.assign_test.id]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

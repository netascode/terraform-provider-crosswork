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
func TestAccCwAssignPortsToPortGroup(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_assign_ports_to_port_group.test", "ports.0.device", "build-router-1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_assign_ports_to_port_group.test", "ports.0.port_name", "GigabitEthernet0/0/0/0"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwAssignPortsToPortGroupPrerequisitesConfig + testAccCwAssignPortsToPortGroupConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_assign_ports_to_port_group.test",
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
const testAccCwAssignPortsToPortGroupPrerequisitesConfig = `
data "crosswork_groupings" "assign_ports_root" {
  classifier = "UserDefinedPorts"
}

resource "crosswork_port_groups" "assign_ports_test" {
  name        = "cwtest-gen-assign-ports"
  parent_uuid = [for group in data.crosswork_groupings.assign_ports_root.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwAssignPortsToPortGroupConfig_minimum() string {
	config := `resource "crosswork_assign_ports_to_port_group" "test" {` + "\n"
	config += `	group_id = crosswork_port_groups.assign_ports_test.id` + "\n"
	config += `	ports = [{` + "\n"
	config += `	  device = "build-router-1"` + "\n"
	config += `	  port_name = "GigabitEthernet0/0/0/0"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwAssignPortsToPortGroupConfig_all() string {
	config := `resource "crosswork_assign_ports_to_port_group" "test" {` + "\n"
	config += `	group_id = crosswork_port_groups.assign_ports_test.id` + "\n"
	config += `	ports = [{` + "\n"
	config += `	  device = "build-router-1"` + "\n"
	config += `	  port_name = "GigabitEthernet0/0/0/0"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

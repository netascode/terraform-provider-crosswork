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
func TestAccCwPortGroupsRules(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_port_groups_rules.test", "conditions.0.attribute", "name"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_port_groups_rules.test", "conditions.0.operator", "starts_with"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_port_groups_rules.test", "conditions.0.value", "Four"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwPortGroupsRulesPrerequisitesConfig + testAccCwPortGroupsRulesConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_port_groups_rules.test",
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
const testAccCwPortGroupsRulesPrerequisitesConfig = `
data "crosswork_groupings" "rules_root" {
  classifier = "UserDefinedPorts"
}

resource "crosswork_port_groups" "rules_test" {
  name        = "cwtest-gen-rules"
  parent_uuid = [for group in data.crosswork_groupings.rules_root.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwPortGroupsRulesConfig_minimum() string {
	config := `resource "crosswork_port_groups_rules" "test" {` + "\n"
	config += `	group_id = crosswork_port_groups.rules_test.id` + "\n"
	config += `	name = crosswork_port_groups.rules_test.name` + "\n"
	config += `	conditions = [{` + "\n"
	config += `	  attribute = "name"` + "\n"
	config += `	  operator = "starts_with"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwPortGroupsRulesConfig_all() string {
	config := `resource "crosswork_port_groups_rules" "test" {` + "\n"
	config += `	group_id = crosswork_port_groups.rules_test.id` + "\n"
	config += `	name = crosswork_port_groups.rules_test.name` + "\n"
	config += `	conditions = [{` + "\n"
	config += `	  attribute = "name"` + "\n"
	config += `	  operator = "starts_with"` + "\n"
	config += `	  value = "Four"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

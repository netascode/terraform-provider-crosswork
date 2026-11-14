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

// Needs the two devices that Crosswork has collected: the management port of each has no speed, and every other port is a
// GigabitEthernet port at 1 Gbps.

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/tidwall/sjson"
)

const ruleMgmtPort = "MgmtEth0/RP0/CPU0/0"

const rulePrerequisites = `
data "crosswork_groupings" "rules_members_root" {
  classifier = "UserDefinedPorts"
}

resource "crosswork_port_groups" "rules_members" {
  name        = "cwtest-gen-rules-members"
  parent_uuid = [for group in data.crosswork_groupings.rules_members_root.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]
}
`

func ruleConfig(conditions string) string {
	return rulePrerequisites + `
resource "crosswork_port_groups_rules" "test" {
  group_id   = crosswork_port_groups.rules_members.id
  name       = crosswork_port_groups.rules_members.name
  conditions = [` + conditions + `]
}
`
}

func checkRuleMembers(t *testing.T, want int64) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		group := s.RootModule().Resources["crosswork_port_groups.rules_members"].Primary.ID
		res, err := testDeviceClient(t).Get(fmt.Sprintf("%s/port/%s?start=0&end=500", groupingPath, group))
		if err != nil {
			return err
		}
		if got := res.Get("total").Int(); got != want {
			return fmt.Errorf("the group holds %d ports, want %d", got, want)
		}
		return nil
	}
}

func ruleInTheUI(t *testing.T, group string) func() {
	return func() {
		c := testDeviceClient(t)
		res, err := c.Get(fmt.Sprintf("%s/rule/group/%s", groupingPath, group))
		if err != nil {
			t.Fatal(err)
		}
		body := res.Get("rule").Raw
		text := ruleConditionsText([]PortGroupsRulesConditions{cond("name", "equals", "GigabitEthernet0/0/0/1", "")})
		body, _ = sjson.Set(body, "conditions", text)
		if r, err := c.Put(fmt.Sprintf("%s/rule/%s", groupingPath, res.Get("rule.uuid").String()), body); err != nil {
			t.Fatalf("editing the rule: %v, %s", err, r.String())
		}
	}
}

func TestAccCwPortGroupsRulesMembers(t *testing.T) {
	mgmt := `{ attribute = "name", operator = "equals", value = "` + ruleMgmtPort + `" }`
	// AND binds tighter than OR: a management port or (GigabitEthernet0/0/0/0 at 1 Gbps), which are 2 and 2 ports on the two
	// devices; read from left to right it would be the 2 ports at 1 Gbps only, because the management port has no speed
	precedence := `{ attribute = "name", operator = "equals", value = "` + ruleMgmtPort + `" },
    { attribute = "name", operator = "equals", value = "GigabitEthernet0/0/0/0", logical_operator = "or" },
    { attribute = "speed", operator = "equals", value = "1 Gbps", logical_operator = "and" }`
	var group string
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t); skipWithoutPortDevices(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: ruleConfig(mgmt),
				Check: resource.ComposeTestCheckFunc(checkRuleMembers(t, 2), func(s *terraform.State) error {
					group = s.RootModule().Resources["crosswork_port_groups.rules_members"].Primary.ID
					return nil
				}),
			},
			{
				Config: ruleConfig(precedence),
				Check:  checkRuleMembers(t, 4),
			},
			{
				PreConfig:          func() { ruleInTheUI(t, group)() },
				Config:             ruleConfig(precedence),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: ruleConfig(precedence),
				Check:  checkRuleMembers(t, 4),
			},
			{
				ResourceName:      "crosswork_port_groups_rules.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// A port type that Crosswork does not list is refused before it is sent: Crosswork would accept it
				Config:      ruleConfig(`{ attribute = "type", operator = "equals", value = "No Such Type" }`),
				ExpectError: regexp.MustCompile(`does not list the port type`),
			},
			{
				Config:      ruleConfig(`{ attribute = "speed", operator = "equals", value = "100Gbps" }`),
				ExpectError: regexp.MustCompile(`a whole number, a space and a unit`),
			},
			{
				Config: ruleConfig(precedence),
				Check:  checkRuleMembers(t, 4),
			},
		},
	})
}

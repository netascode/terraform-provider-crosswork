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

// Needs two onboarded, collected IOS XR devices with ports of the same names.

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/tidwall/sjson"
)

const (
	portDeviceA = "build-router-1"
	portDeviceB = "build-router-2"
	portZero    = "GigabitEthernet0/0/0/0"
	portOne     = "GigabitEthernet0/0/0/10"
)

const portMembersPrerequisites = `
data "crosswork_groupings" "port_members_root" {
  classifier = "UserDefinedPorts"
}

resource "crosswork_port_groups" "members" {
  name        = "cwtest-gen-ports-members"
  parent_uuid = [for group in data.crosswork_groupings.port_members_root.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]
}
`

var portMembersGroup string

func portConfig(ports ...[2]string) string {
	rows := make([]string, len(ports))
	for i, p := range ports {
		rows[i] = fmt.Sprintf("{ device = %q, port_name = %q }", p[0], p[1])
	}
	return portMembersPrerequisites + `
resource "crosswork_assign_ports_to_port_group" "test" {
  group_id = crosswork_port_groups.members.id
  ports    = [` + strings.Join(rows, ", ") + `]
}
`
}

func skipWithoutPortDevices(t *testing.T) {
	c := testDeviceClient(t)
	if _, err := deviceUUIDs(c, []string{portDeviceA, portDeviceB}); err != nil {
		t.Skipf("the cluster has no %s and %s: %v", portDeviceA, portDeviceB, err)
	}
}

func checkPortMembers(t *testing.T, want ...[2]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		portMembersGroup = s.RootModule().Resources["crosswork_port_groups.members"].Primary.ID
		got, err := portMembers(testDeviceClient(t), portMembersGroup)
		if err != nil {
			return err
		}
		var have, wanted []string
		for _, p := range got {
			have = append(have, p.Device.ValueString()+"|"+p.PortName.ValueString())
		}
		for _, p := range want {
			wanted = append(wanted, p[0]+"|"+p[1])
		}
		sort.Strings(have)
		sort.Strings(wanted)
		if strings.Join(have, ",") != strings.Join(wanted, ",") {
			return fmt.Errorf("the group holds %v, want %v", have, wanted)
		}
		return nil
	}
}

func portsInTheUI(t *testing.T, add, remove [][2]string) func() {
	return func() {
		c := testDeviceClient(t)
		devices, err := deviceUUIDs(c, []string{portDeviceA, portDeviceB})
		if err != nil {
			t.Fatal(err)
		}
		if len(add) > 0 {
			all, err := listPorts(c, portMembersGroup, "Add")
			if err != nil {
				t.Fatal(err)
			}
			byKey := portsByKey(all)
			var refs []map[string]string
			for _, p := range add {
				cand := byKey[portKey{devices[p[0]], p[1]}]
				refs = append(refs, map[string]string{"portUuid": cand.Get("portUuid").String(), "deviceUuid": devices[p[0]]})
			}
			body, _ := sjson.Set("", "ports", refs)
			if res, err := c.Post(fmt.Sprintf("%s/port/%s", groupingPath, portMembersGroup), body); err != nil {
				t.Fatalf("adding ports: %v, %s", err, res.String())
			}
		}
		if len(remove) > 0 {
			manual, err := listPorts(c, portMembersGroup, "Manual")
			if err != nil {
				t.Fatal(err)
			}
			byKey := portsByKey(manual)
			var objs []interface{}
			for _, p := range remove {
				objs = append(objs, byKey[portKey{devices[p[0]], p[1]}].Value())
			}
			body, _ := sjson.Set("", "ports", objs)
			if res, err := c.DeleteBody(fmt.Sprintf("%s/port/%s", groupingPath, portMembersGroup), body); err != nil {
				t.Fatalf("removing ports: %v, %s", err, res.String())
			}
		}
	}
}

func TestAccCwAssignPortsToPortGroupMembers(t *testing.T) {
	a0, a1, b0 := [2]string{portDeviceA, portZero}, [2]string{portDeviceA, portOne}, [2]string{portDeviceB, portZero}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t); skipWithoutPortDevices(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: portConfig(a0),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("crosswork_assign_ports_to_port_group.test", "ports.#", "1"), checkPortMembers(t, a0)),
			},
			{
				// The same port name on the second device is another port.
				Config: portConfig(a0, b0),
				Check:  checkPortMembers(t, a0, b0),
			},
			{
				Config: portConfig(a1, b0),
				Check:  checkPortMembers(t, a1, b0),
			},
			{
				PreConfig:          portsInTheUI(t, [][2]string{a0}, nil),
				Config:             portConfig(a1, b0),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: portConfig(a1, b0),
				Check:  checkPortMembers(t, a1, b0),
			},
			{
				PreConfig:          portsInTheUI(t, nil, [][2]string{b0}),
				Config:             portConfig(a1, b0),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: portConfig(a1, b0),
				Check:  checkPortMembers(t, a1, b0),
			},
			{
				// A port that does not exist is refused, and the group is left as it was.
				Config:      portConfig(a1, b0, [2]string{portDeviceA, "NoSuchPort9/9/9"}),
				ExpectError: regexp.MustCompile(`does not know the port`),
			},
			{
				Config: portConfig(),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("crosswork_assign_ports_to_port_group.test", "ports.#", "0"), checkPortMembers(t)),
			},
		},
	})
}

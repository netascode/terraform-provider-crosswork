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

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

const placeMembersPrerequisites = `
data "crosswork_groupings" "place_location" {
  classifier = "LocationDevices"
}

resource "crosswork_device_groups" "place_one" {
  name        = "cwtest-gen-place-one"
  parent_uuid = [for group in data.crosswork_groupings.place_location.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]
}

resource "crosswork_device_groups" "place_two" {
  name        = "cwtest-gen-place-two"
  parent_uuid = [for group in data.crosswork_groupings.place_location.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]
}

resource "crosswork_credential_profiles" "place" {
  profile = "cwtest-gen-place-profile"
  credentials = [
    {
      type      = "ROBOT_USERPASS_SSH"
      user_name = "cwtest"
      password  = "Cwtest-Password-1"
    }
  ]
}

resource "crosswork_network_devices" "place_a" {
  name               = "cwtest-gen-place-a"
  host_name          = "cwtest-gen-place-a"
  credential_profile = crosswork_credential_profiles.place.id
  software_type      = "IOS XR"
  capabilities       = ["YANG_CLI"]
  connectivity = [{
    protocol      = "ROBOT_MSVC_TRANS_SSH"
    ip_address    = "192.0.2.245"
    prefix_length = 24
    port          = 22
  }]

  # Crosswork can move the admin state itself: an embedded data gateway pool attaches a new device and brings it up
  lifecycle {
    ignore_changes = [admin_state]
  }
}

resource "crosswork_network_devices" "place_b" {
  name               = "cwtest-gen-place-b"
  host_name          = "cwtest-gen-place-b"
  credential_profile = crosswork_credential_profiles.place.id
  software_type      = "IOS XR"
  capabilities       = ["YANG_CLI"]
  connectivity = [{
    protocol      = "ROBOT_MSVC_TRANS_SSH"
    ip_address    = "192.0.2.246"
    prefix_length = 24
    port          = 22
  }]

  # Crosswork can move the admin state itself: an embedded data gateway pool attaches a new device and brings it up
  lifecycle {
    ignore_changes = [admin_state]
  }
}
`

const (
	placeA = "cwtest-gen-place-a"
	placeB = "cwtest-gen-place-b"
)

var placeOne, placeTwo string

func placeConfig(one []string, two []string) string {
	refs := map[string]string{placeA: "crosswork_network_devices.place_a.id", placeB: "crosswork_network_devices.place_b.id"}
	list := func(names []string) string {
		out := make([]string, len(names))
		for i, n := range names {
			out[i] = refs[n]
		}
		return strings.Join(out, ", ")
	}
	config := placeMembersPrerequisites + `
resource "crosswork_assign_devices_to_device_group" "one" {
  group_id = crosswork_device_groups.place_one.id
  devices  = [` + list(one) + `]
}
`
	if two != nil {
		config += `
resource "crosswork_assign_devices_to_device_group" "two" {
  group_id = crosswork_device_groups.place_two.id
  devices  = [` + list(two) + `]
}
`
	}
	return config
}

func unassignedGroup(t *testing.T) string {
	c := testDeviceClient(t)
	roots, err := c.Get(groupingPath+"/group/root/LocationDevices/uuid", cw.ReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	tree, err := c.Get(fmt.Sprintf("%s/groups/%s?brief=true&direct=false", groupingPath, roots.Array()[0].String()), cw.ReadOnly)
	if err != nil {
		t.Fatal(err)
	}
	var found string
	var walk func(n gjson.Result)
	walk = func(n gjson.Result) {
		if n.Get("name").String() == "Unassigned Devices" {
			found = n.Get("uuid").String()
		}
		for _, child := range n.Get("children").Array() {
			walk(child)
		}
	}
	for _, n := range tree.Array() {
		walk(n)
	}
	if found == "" {
		t.Fatal("Unassigned Devices not found")
	}
	return found
}

func checkPlace(t *testing.T, group func() string, want ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		placeOne = s.RootModule().Resources["crosswork_device_groups.place_one"].Primary.ID
		placeTwo = s.RootModule().Resources["crosswork_device_groups.place_two"].Primary.ID
		got, err := groupDevices(testDeviceClient(t), group())
		if err != nil {
			return err
		}
		// the built-in group also holds the devices of the user, so only the devices of the test are compared
		var mine []string
		for _, n := range got {
			if strings.HasPrefix(n, "cwtest-gen-place-") {
				mine = append(mine, n)
			}
		}
		sort.Strings(want)
		if strings.Join(mine, ",") != strings.Join(want, ",") {
			return fmt.Errorf("the group holds %v, want %v", mine, want)
		}
		return nil
	}
}

func inOne() string { return placeOne }
func inTwo() string { return placeTwo }

func placeInTheUI(t *testing.T, move, remove []string) func() {
	return func() {
		c := testDeviceClient(t)
		if len(move) > 0 {
			uuids, err := deviceUUIDs(c, move)
			if err != nil {
				t.Fatal(err)
			}
			if err := moveMembers(c, placeOne, values(uuids)); err != nil {
				t.Fatal(err)
			}
		}
		if len(remove) > 0 {
			uuids, err := deviceUUIDs(c, remove)
			if err != nil {
				t.Fatal(err)
			}
			if err := removeMembers(c, placeOne, values(uuids)); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestAccCwAssignDevicesToDeviceGroupMembers(t *testing.T) {
	const one = "crosswork_assign_devices_to_device_group.one"
	unassigned := func() string { return unassignedGroup(t) }
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: placeConfig([]string{placeA}, nil),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(one, "devices.#", "1"), checkPlace(t, inOne, placeA), checkPlace(t, unassigned, placeB)),
			},
			{
				Config: placeConfig([]string{placeA, placeB}, nil),
				Check:  checkPlace(t, inOne, placeA, placeB),
			},
			{
				Config: placeConfig([]string{placeB}, nil),
				Check:  resource.ComposeTestCheckFunc(checkPlace(t, inOne, placeB), checkPlace(t, unassigned, placeA)),
			},
			{
				PreConfig:          placeInTheUI(t, []string{placeA}, nil),
				Config:             placeConfig([]string{placeB}, nil),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: placeConfig([]string{placeB}, nil),
				Check:  resource.ComposeTestCheckFunc(checkPlace(t, inOne, placeB), checkPlace(t, unassigned, placeA)),
			},
			{
				PreConfig:          placeInTheUI(t, nil, []string{placeB}),
				Config:             placeConfig([]string{placeB}, nil),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: placeConfig([]string{placeB}, nil),
				Check:  checkPlace(t, inOne, placeB),
			},
			{
				Config: placeConfig([]string{placeB}, []string{placeA}),
				Check:  resource.ComposeTestCheckFunc(checkPlace(t, inOne, placeB), checkPlace(t, inTwo, placeA)),
			},
			{
				Config: placeConfig([]string{placeA}, []string{placeB}),
				Check:  resource.ComposeTestCheckFunc(checkPlace(t, inOne, placeA), checkPlace(t, inTwo, placeB)),
			},
			{
				// Empty lists put every device back into Unassigned Devices.
				Config: placeConfig([]string{}, []string{}),
				Check:  resource.ComposeTestCheckFunc(checkPlace(t, inOne), checkPlace(t, inTwo), checkPlace(t, unassigned, placeA, placeB)),
			},
		},
	})
}

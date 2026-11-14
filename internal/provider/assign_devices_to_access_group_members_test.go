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
	"github.com/tidwall/sjson"
)

const assignMembersPrerequisites = `
data "crosswork_groupings" "members_access" {
  classifier = "DeviceAccess"
}

data "crosswork_groupings" "members_location" {
  classifier = "LocationDevices"
}

resource "crosswork_device_access_groups" "members" {
  name        = "cwtest-gen-assign-members-group"
  parent_uuid = [for group in data.crosswork_groupings.members_access.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]
}

resource "crosswork_device_groups" "members" {
  name        = "cwtest-gen-assign-members-place"
  parent_uuid = [for group in data.crosswork_groupings.members_location.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]
}

resource "crosswork_credential_profiles" "members" {
  profile = "cwtest-gen-assign-members-profile"
  credentials = [
    {
      type      = "ROBOT_USERPASS_SSH"
      user_name = "cwtest"
      password  = "Cwtest-Password-1"
    }
  ]
}

resource "crosswork_network_devices" "members_a" {
  name               = "cwtest-gen-assign-members-a"
  host_name          = "cwtest-gen-assign-members-a"
  credential_profile = crosswork_credential_profiles.members.id
  software_type      = "IOS XR"
  capabilities       = ["YANG_CLI"]
  connectivity = [{
    protocol      = "ROBOT_MSVC_TRANS_SSH"
    ip_address    = "192.0.2.242"
    prefix_length = 24
    port          = 22
  }]

  # Crosswork can move the admin state itself: an embedded data gateway pool attaches a new device and brings it up
  lifecycle {
    ignore_changes = [admin_state]
  }
}

resource "crosswork_network_devices" "members_b" {
  name               = "cwtest-gen-assign-members-b"
  host_name          = "cwtest-gen-assign-members-b"
  credential_profile = crosswork_credential_profiles.members.id
  software_type      = "IOS XR"
  capabilities       = ["YANG_CLI"]
  connectivity = [{
    protocol      = "ROBOT_MSVC_TRANS_SSH"
    ip_address    = "192.0.2.243"
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
	memberA = "cwtest-gen-assign-members-a"
	memberB = "cwtest-gen-assign-members-b"
)

var assignMembersGroup, assignMembersPlace string

func assignMembersConfig(devices ...string) string {
	// The devices are referred to by their resources, so that they are created first and deleted last.
	refs := map[string]string{memberA: "crosswork_network_devices.members_a.id", memberB: "crosswork_network_devices.members_b.id"}
	quoted := make([]string, len(devices))
	for i, d := range devices {
		quoted[i] = refs[d]
	}
	return assignMembersPrerequisites + `
resource "crosswork_assign_devices_to_access_group" "test" {
  group_id = crosswork_device_access_groups.members.id
  devices  = [` + strings.Join(quoted, ", ") + `]
}
`
}

func checkMembers(t *testing.T, want ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		assignMembersGroup = s.RootModule().Resources["crosswork_device_access_groups.members"].Primary.ID
		assignMembersPlace = s.RootModule().Resources["crosswork_device_groups.members"].Primary.ID
		c := testDeviceClient(t)
		got, err := groupDevices(c, assignMembersGroup)
		if err != nil {
			return err
		}
		sort.Strings(want)
		if strings.Join(got, ",") != strings.Join(want, ",") {
			return fmt.Errorf("the group has the devices %v, want %v", got, want)
		}
		return nil
	}
}

func inTheUI(t *testing.T, add, remove []string) func() {
	return func() {
		c := testDeviceClient(t)
		if len(add) > 0 {
			uuids, err := deviceUUIDs(c, add)
			if err != nil {
				t.Fatal(err)
			}
			if err := addMembers(c, assignMembersGroup, values(uuids)); err != nil {
				t.Fatal(err)
			}
		}
		if len(remove) > 0 {
			uuids, err := deviceUUIDs(c, remove)
			if err != nil {
				t.Fatal(err)
			}
			if err := removeMembers(c, assignMembersGroup, values(uuids)); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func moveToPlace(t *testing.T, name string) func() {
	return func() {
		c := testDeviceClient(t)
		uuids, err := deviceUUIDs(c, []string{name})
		if err != nil {
			t.Fatal(err)
		}
		uuid := uuids[name]
		sources, err := sourceGroups(c, []string{uuid})
		if err != nil || sources[uuid] == "" {
			t.Fatalf("finding the location group of %s: %v", name, err)
		}
		body, _ := sjson.Set("", "groupUuid", sources[uuid])
		body, _ = sjson.Set(body, "newGroupUuid", assignMembersPlace)
		body, _ = sjson.Set(body, "references", []string{uuid})
		if res, err := c.Post(groupingPath+"/group/member", body); err != nil {
			t.Fatalf("moving %s: %v, %s", name, err, res.String())
		}
	}
}

func values(m map[string]string) []string {
	out := make([]string, 0, len(m))
	for _, v := range m {
		out = append(out, v)
	}
	return out
}

var _ = cw.ReadOnly

func TestAccCwAssignDevicesToAccessGroupMembers(t *testing.T) {
	const res = "crosswork_assign_devices_to_access_group.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: assignMembersConfig(memberA),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(res, "devices.#", "1"), checkMembers(t, memberA)),
			},
			{
				Config: assignMembersConfig(memberA, memberB),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(res, "devices.#", "2"), checkMembers(t, memberA, memberB)),
			},
			{
				Config: assignMembersConfig(memberB),
				Check:  checkMembers(t, memberB),
			},
			{
				PreConfig:          inTheUI(t, []string{memberA}, nil),
				Config:             assignMembersConfig(memberB),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: assignMembersConfig(memberB),
				Check:  checkMembers(t, memberB),
			},
			{
				PreConfig:          inTheUI(t, nil, []string{memberB}),
				Config:             assignMembersConfig(memberB),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: assignMembersConfig(memberB),
				Check:  checkMembers(t, memberB),
			},
			{
				Config: assignMembersConfig(),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(res, "devices.#", "0"), checkMembers(t)),
			},
			{
				// A device that sits in another location group than Unassigned Devices is found and copied from there.
				PreConfig: moveToPlace(t, memberA),
				Config:    assignMembersConfig(memberA, memberB),
				Check:     checkMembers(t, memberA, memberB),
			},
		},
	})
}

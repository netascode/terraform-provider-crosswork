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

// Uses real devices because a throwaway device is copied into NSO by Crosswork and comes back after the delete.
// Needs a Crosswork with no pool, a free Data Gateway instance (the `CLUSTER` tag) and two onboarded devices that are not
// attached to a data gateway. Attaching a device brings it admin up; the destroy detaches it and takes it admin down.

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	cw "github.com/netascode/go-crosswork"
)

const (
	cdgDeviceA = "build-router-1"
	cdgDeviceB = "build-router-2"
)

const cdgMembersPrerequisites = `
data "crosswork_data_gateway_instances" "members" {}

resource "crosswork_data_gateway_pools" "members" {
  name          = "cwtest-gen-attach-members"
  virtual_ips   = ["172.16.34.198"]
  prefix_length = 24
  instances = [[for instance in data.crosswork_data_gateway_instances.members.items : instance.id if instance.name == sort([for i in data.crosswork_data_gateway_instances.members.items : i.name])[0]][0]]
}
`

func cdgConfig(devices ...string) string {
	quoted := make([]string, len(devices))
	for i, d := range devices {
		quoted[i] = fmt.Sprintf("%q", d)
	}
	return cdgMembersPrerequisites + `
resource "crosswork_assign_devices_to_data_gateway" "test" {
  data_gateway_id = crosswork_data_gateway_pools.members.data_gateway_ids[0]
  devices         = [` + strings.Join(quoted, ", ") + `]
}
`
}

func skipWithoutCdgDevices(t *testing.T) {
	if os.Getenv("CLUSTER") == "" {
		t.Skip("skipping test, set environment variable CLUSTER")
	}
	c := testDeviceClient(t)
	found, err := lookupDevices(c, []string{cdgDeviceA, cdgDeviceB})
	if err != nil {
		t.Skipf("the cluster has no %s and %s: %v", cdgDeviceA, cdgDeviceB, err)
	}
	for _, d := range found {
		if d.Gateway != "" {
			t.Skipf("%s is attached to a data gateway already", d.Name)
		}
	}
}

func cdgMembers(want ...string) resource.TestCheckFunc {
	sort.Strings(want)
	return func(s *terraform.State) error {
		pool := s.RootModule().Resources["crosswork_data_gateway_pools.members"]
		if pool == nil {
			return fmt.Errorf("the pool is not in the state")
		}
		members, err := attachedDevices(testDeviceClientNoT(), pool.Primary.Attributes["data_gateway_ids.0"])
		if err != nil {
			return err
		}
		var got []string
		for n := range members {
			got = append(got, n)
		}
		sort.Strings(got)
		if strings.Join(got, ",") != strings.Join(want, ",") {
			return fmt.Errorf("the data gateway holds %v, want %v", got, want)
		}
		return nil
	}
}

func TestAccCwAssignDevicesToDataGatewayMembers(t *testing.T) {
	skipWithoutCdgDevices(t)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: cdgConfig(cdgDeviceA),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("crosswork_assign_devices_to_data_gateway.test", "devices.#", "1"),
					cdgMembers(cdgDeviceA),
				),
			},
			{
				Config: cdgConfig(cdgDeviceA, cdgDeviceB),
				Check:  cdgMembers(cdgDeviceA, cdgDeviceB),
			},
			{
				Config: cdgConfig(cdgDeviceB),
				Check:  cdgMembers(cdgDeviceB),
			},
		},
	})
}

// testDeviceClientNoT is testDeviceClient for a check function, which has no *testing.T.
func testDeviceClientNoT() *cw.Client {
	c, err := registeredClient(os.Getenv("CW_URL"), os.Getenv("CW_USERNAME"), os.Getenv("CW_PASSWORD"), true, 3, "test")
	if err != nil {
		panic(err)
	}
	return c
}

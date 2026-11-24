// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License Version 2.0 (the "License"),
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://mozilla.org/MPL/2.0/
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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// The name carries a space and a colon: the controller reads a space written as "+" literally, so a policy of this name
// could not be updated or deleted before the path was escaped.
const testAccCwDeviceNotificationPoliciesName = "cwtest-dnp a:b"

const testAccCwDeviceNotificationPoliciesGroupConfig = `
data "crosswork_groupings" "anchor" {
  classifier = "LocationDevices"
}

resource "crosswork_device_groups" "dnp_test" {
  name        = "cwtest-dnp-group"
  parent_uuid = [for group in data.crosswork_groupings.anchor.items : group.id if group.discovery_type == "StaticSystem" && group.subgroups_allowed][0]
}
`

func testAccCwDeviceNotificationPoliciesLifecycleConfig(body string) string {
	return testAccCwDeviceNotificationPoliciesPrerequisitesConfig + testAccCwDeviceNotificationPoliciesGroupConfig +
		`resource "crosswork_device_notification_policies" "test" {` + "\n" +
		`  name = "` + testAccCwDeviceNotificationPoliciesName + `"` + "\n" + body + "}\n"
}

func TestAccCwDeviceNotificationPoliciesLifecycle(t *testing.T) {
	const address = "crosswork_device_notification_policies.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create: two filters in the reverse of the order the controller stores them, a group and a trap receiver.
				Config: testAccCwDeviceNotificationPoliciesLifecycleConfig(`
  description = "first"
  device_group = crosswork_device_groups.dnp_test.id
  filters = [
    { category = "BGP", severities = ["major", "minor"] },
    { severities = ["critical"] },
  ]
  trap_receivers = [{ name = crosswork_trap_destination.dnp_test.name }]
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(address, "id", testAccCwDeviceNotificationPoliciesName),
					resource.TestCheckResourceAttrSet(address, "device_group"),
					resource.TestCheckResourceAttr(address, "filters.#", "2"),
					resource.TestCheckResourceAttr(address, "filters.0.category", "BGP"),
				),
			},
			{
				// The refresh after the create reads the same policy: no difference, and the order of the configuration stays.
				Config:   testAccCwDeviceNotificationPoliciesLifecycleConfig("  description = \"first\"\n  device_group = crosswork_device_groups.dnp_test.id\n  filters = [\n    { category = \"BGP\", severities = [\"major\", \"minor\"] },\n    { severities = [\"critical\"] },\n  ]\n  trap_receivers = [{ name = crosswork_trap_destination.dnp_test.name }]\n"),
				PlanOnly: true,
			},
			{
				// Update in place: another filter, no group, no description.
				Config: testAccCwDeviceNotificationPoliciesLifecycleConfig(`
  filters = [{ events = ["OPTICS RX LOL LANE-0"] }]
  trap_receivers = [{ name = crosswork_trap_destination.dnp_test.name }]
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(address, "device_group"),
					resource.TestCheckNoResourceAttr(address, "description"),
					resource.TestCheckResourceAttr(address, "filters.0.events.#", "1"),
				),
			},
			{
				// A receiver change replaces the policy; the syslog receiver is added.
				Config: testAccCwDeviceNotificationPoliciesLifecycleConfig(`
  filters = [{ severities = ["warning", "information"] }]
  trap_receivers = [{ name = crosswork_trap_destination.dnp_test.name }]
  syslog_receivers = [{ name = crosswork_syslog_destination.dnp_test.name }]
`),
				Check: resource.TestCheckResourceAttr(address, "syslog_receivers.#", "1"),
			},
			{
				ResourceName:      address,
				ImportState:       true,
				ImportStateId:     testAccCwDeviceNotificationPoliciesName,
				ImportStateVerify: true,
			},
		},
	})
}

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

// Covers the updates of crosswork_network_devices that the controller refuses while a device is admin up (the update runs
// it as down, change, up) and the ones it accepts.

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

type fieldsDevice struct {
	name      string
	admin     string
	reach     string
	sshPort   int
	ip        string
	snmp      bool
	caps      string
	extraAttr string
}

func (d fieldsDevice) config() string {
	rows := fmt.Sprintf(`{ protocol = "ROBOT_MSVC_TRANS_SSH", ip_address = %q, prefix_length = 24, port = %d, timeout = 60 }`, d.ip, d.sshPort)
	if d.snmp {
		rows += fmt.Sprintf(`, { protocol = "ROBOT_MSVC_TRANS_SNMP", ip_address = %q, prefix_length = 24, port = 161, timeout = 60 }`, d.ip)
	}
	var b strings.Builder
	b.WriteString(`resource "crosswork_network_devices" "test" {` + "\n")
	b.WriteString(fmt.Sprintf("\tname = %q\n\thost_name = \"cwtest-fields\"\n\tsoftware_type = \"IOS XR\"\n", d.name))
	b.WriteString("\tcredential_profile = crosswork_credential_profiles.device_test.id\n")
	b.WriteString(fmt.Sprintf("\tadmin_state = %q\n\treachability_check = %q\n\tcapabilities = [%s]\n", d.admin, d.reach, d.caps))
	b.WriteString("\tconnectivity = [" + rows + "]\n")
	b.WriteString(d.extraAttr)
	b.WriteString("}\n")
	return b.String()
}

func TestAccCwNetworkDevicesFields(t *testing.T) {
	const res = "crosswork_network_devices.test"
	const up, down = "ROBOT_ADMIN_STATE_UP", "ROBOT_ADMIN_STATE_DOWN"
	base := fieldsDevice{name: "cwtest-gen-device-fields", admin: up, reach: "REACH_CHECK_ENABLE", sshPort: 22, ip: "192.0.2.233", snmp: true, caps: `"SNMP", "YANG_CLI"`}
	with := func(f func(*fieldsDevice)) string {
		d := base
		f(&d)
		return testAccCwNetworkDevicesPrerequisitesConfig + d.config()
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: with(func(d *fieldsDevice) {}),
				Check:  resource.TestCheckResourceAttr(res, "connectivity.#", "2"),
			},
			{
				// The reachability check is refused while the device is admin up.
				Config: with(func(d *fieldsDevice) { d.reach = "REACH_CHECK_DISABLE" }),
				Check:  resource.TestCheckResourceAttr(res, "reachability_check", "REACH_CHECK_DISABLE"),
			},
			{
				Config: with(func(d *fieldsDevice) { d.reach = "REACH_CHECK_DISABLE"; d.sshPort = 2222; d.ip = "192.0.2.234" }),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "connectivity.0.port", "2222"),
					resource.TestCheckResourceAttr(res, "connectivity.0.ip_address", "192.0.2.234"),
				),
			},
			{
				// A row and a capability removed: the SNMP capability needs an SNMP row.
				Config: with(func(d *fieldsDevice) {
					d.reach = "REACH_CHECK_DISABLE"
					d.sshPort = 2222
					d.ip = "192.0.2.234"
					d.snmp = false
					d.caps = `"YANG_CLI"`
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "connectivity.#", "1"),
					resource.TestCheckResourceAttr(res, "capabilities.#", "1"),
				),
			},
			{
				// The row and the capability back, with routing ids (the TE router id is refused while up) and a location.
				Config: with(func(d *fieldsDevice) {
					d.reach = "REACH_CHECK_DISABLE"
					d.sshPort = 2222
					d.ip = "192.0.2.234"
					d.extraAttr = "\tte_router_id = \"192.0.2.9\"\n\tospf_router_id = \"192.0.2.8\"\n\tlocation_city = \"Amsterdam\"\n\tlocation_latitude = 52.37\n\tlocation_longitude = 4.89\n"
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "connectivity.#", "2"),
					resource.TestCheckResourceAttr(res, "te_router_id", "192.0.2.9"),
					resource.TestCheckResourceAttr(res, "location_city", "Amsterdam"),
				),
			},
			{
				Config: with(func(d *fieldsDevice) { d.reach = "REACH_CHECK_DISABLE"; d.sshPort = 2222; d.ip = "192.0.2.234" }),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(res, "te_router_id"),
					resource.TestCheckNoResourceAttr(res, "location_city"),
				),
			},
			{
				// Admin down together with a port change: the device is taken down and stays down.
				Config: with(func(d *fieldsDevice) {
					d.admin = down
					d.reach = "REACH_CHECK_DISABLE"
					d.sshPort = 2223
					d.ip = "192.0.2.234"
				}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "admin_state", down),
					resource.TestCheckResourceAttr(res, "connectivity.0.port", "2223"),
				),
			},
			{
				// A refused field while the device is already down is a plain change.
				Config: with(func(d *fieldsDevice) {
					d.admin = down
					d.reach = "REACH_CHECK_ENABLE"
					d.sshPort = 2223
					d.ip = "192.0.2.234"
				}),
				Check: resource.TestCheckResourceAttr(res, "reachability_check", "REACH_CHECK_ENABLE"),
			},
			{
				Config: with(func(d *fieldsDevice) { d.admin = up; d.sshPort = 2223; d.ip = "192.0.2.234" }),
				Check:  resource.TestCheckResourceAttr(res, "admin_state", up),
			},
			{
				Config: with(func(d *fieldsDevice) { d.name = "cwtest-gen-device-fields-2"; d.sshPort = 2223; d.ip = "192.0.2.234" }),
				Check:  resource.TestCheckResourceAttr(res, "id", "cwtest-gen-device-fields-2"),
			},
		},
	})
}

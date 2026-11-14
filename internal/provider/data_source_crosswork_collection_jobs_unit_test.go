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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCollectionJobsDataSource(t *testing.T) {
	s := newCollectionStub(t)
	s.seed("app-one", "ctx-cli", "CLI_COLLECTOR", `{"device_set":{"devices":{"device_ids":["uuid-1","uuid-2"]}}}`, `{"cli_sensor":{"command":"show version"}}`)
	s.seed("app-one", "ctx-snmp", "SNMP_COLLECTOR", `{"device_set":{"device_group":"cwtest-tag"}}`,
		`{"snmp_sensor":{"snmp_mib":{"snmp_operation":"WALK","oids":["1.3.6.1.2.1.1.5"]}}}`)
	// a job of a collector type the resource does not manage is read as the controller names it
	s.seed("app-one", "ctx-gnmi", "GNMI_COLLECTOR", `{"device_set":{"device_group":"cwtest-tag"}}`, `{"gnmi_sensor":{"mode":"SAMPLE"}}`)
	read := func(ctx string) string {
		return s.providerBlock(false) + `
data "crosswork_collection_jobs" "test" {
  application_id = "app-one"
  context_id     = "` + ctx + `"
}
`
	}
	const ds = "data.crosswork_collection_jobs.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: read("ctx-cli"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ds, "id", "app-one,ctx-cli"),
					resource.TestCheckResourceAttr(ds, "collector_type", "CLI"),
					resource.TestCheckResourceAttr(ds, "devices.#", "2"),
					resource.TestCheckTypeSetElemAttr(ds, "devices.*", "dev-1"),
					resource.TestCheckResourceAttr(ds, "destination", "cwtest-dest-kafka"),
					resource.TestCheckResourceAttr(ds, "cli_sensors.0.command", "show version"),
					resource.TestCheckResourceAttr(ds, "cli_sensors.0.cadence", "60"),
					resource.TestCheckResourceAttr(ds, "cli_sensors.0.topic", "topic-seed"),
					resource.TestCheckResourceAttr(ds, "snmp_sensors.#", "0"),
					resource.TestCheckResourceAttr(ds, "status", "READY"),
					resource.TestCheckResourceAttr(ds, "phase", "ACTIVE"),
					resource.TestCheckResourceAttr(ds, "life_cycle_state", "SUCCESS_LIFE_CYCLE_STATE"),
				),
			},
			{
				Config: read("ctx-snmp"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ds, "collector_type", "SNMP"),
					resource.TestCheckResourceAttr(ds, "device_tag", "cwtest-tag"),
					resource.TestCheckResourceAttr(ds, "devices.#", "0"),
					resource.TestCheckResourceAttr(ds, "snmp_sensors.0.oid", "1.3.6.1.2.1.1.5"),
					resource.TestCheckResourceAttr(ds, "snmp_sensors.0.operation", "WALK"),
					resource.TestCheckResourceAttr(ds, "cli_sensors.#", "0"),
				),
			},
			{
				Config: read("ctx-gnmi"),
				Check:  resource.TestCheckResourceAttr(ds, "collector_type", "GNMI_COLLECTOR"),
			},
			{
				Config:      read("no-such-context"),
				ExpectError: regexp.MustCompile(`(?s)There is no collection job with the application id app-one`),
			},
		},
	})
}

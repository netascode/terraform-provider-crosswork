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

// Needs a cluster with Crosswork Data Gateway for the destination and the two onboarded devices build-router-1
// and build-router-2, which the test neither creates nor deletes. A job is deleted before the destination and the tag, because
// Crosswork does not refuse to delete those while a job uses them.

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccCwCollectionJobsConfig(command string, cadence int, topic string) string {
	return fmt.Sprintf(`
resource "crosswork_data_gateway_destinations" "collection_test" {
  name               = "cwtest-gen-collection-destination"
  data_source        = "data-gateway"
  server_type        = "kafka"
  encoding           = "gpbkv"
  compression_method = "snappy"
  max_message_size   = 100000000
  buffer_memory      = 52428800
  batch_size         = 1048576
  linger_time        = 2000
  request_timeout    = 30
  connections = [
    {
      address       = "192.0.2.40"
      prefix_length = 24
      port          = 9192
    }
  ]
}

resource "crosswork_tags" "collection_test" {
  name     = "cwtest-gen-collection-tag"
  category = "cwtest-category"
}

resource "crosswork_collection_jobs" "cli" {
  application_id = "cwtest-gen-app"
  context_id     = "cwtest-gen-cli"
  collector_type = "CLI"
  devices        = ["build-router-1", "build-router-2"]
  destination    = crosswork_data_gateway_destinations.collection_test.name
  cli_sensors = [
    {
      command = %q
      cadence = %d
      topic   = %q
    }
  ]
}

resource "crosswork_collection_jobs" "snmp" {
  application_id = "cwtest-gen-app"
  context_id     = "cwtest-gen-snmp"
  collector_type = "SNMP"
  device_tag     = crosswork_tags.collection_test.name
  destination    = crosswork_data_gateway_destinations.collection_test.name
  snmp_sensors = [
    {
      oid       = "1.3.6.1.2.1.1.5"
      operation = "WALK"
      cadence   = 60
      topic     = "cwtest-topic"
    }
  ]
}

data "crosswork_collection_jobs" "cli" {
  application_id = crosswork_collection_jobs.cli.application_id
  context_id     = crosswork_collection_jobs.cli.context_id
}
`, command, cadence, topic)
}

func TestAccCwCollectionJobs(t *testing.T) {
	const cli, snmp, ds = "crosswork_collection_jobs.cli", "crosswork_collection_jobs.snmp", "data.crosswork_collection_jobs.cli"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCwCollectionJobsConfig("show version", 60, "cwtest-topic"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(cli, "id", "cwtest-gen-app,cwtest-gen-cli"),
					resource.TestCheckResourceAttr(cli, "devices.#", "2"),
					resource.TestCheckTypeSetElemAttr(cli, "devices.*", "build-router-1"),
					resource.TestCheckTypeSetElemAttr(cli, "devices.*", "build-router-2"),
					resource.TestCheckResourceAttr(cli, "destination", "cwtest-gen-collection-destination"),
					resource.TestCheckResourceAttr(cli, "cli_sensors.0.command", "show version"),
					resource.TestCheckResourceAttr(cli, "cli_sensors.0.cadence", "60"),
					resource.TestCheckResourceAttr(cli, "cli_sensors.0.topic", "cwtest-topic"),
					resource.TestCheckResourceAttr(snmp, "device_tag", "cwtest-gen-collection-tag"),
					resource.TestCheckResourceAttr(snmp, "snmp_sensors.0.oid", "1.3.6.1.2.1.1.5"),
					resource.TestCheckResourceAttr(snmp, "snmp_sensors.0.operation", "WALK"),
					resource.TestCheckResourceAttr(snmp, "snmp_sensors.0.cadence", "60"),
					resource.TestCheckResourceAttr(ds, "collector_type", "CLI"),
					resource.TestCheckResourceAttr(ds, "devices.#", "2"),
					resource.TestCheckResourceAttr(ds, "cli_sensors.0.command", "show version"),
					resource.TestCheckResourceAttrSet(ds, "status"),
					resource.TestCheckResourceAttrSet(ds, "phase"),
				),
			},
			{
				Config: testAccCwCollectionJobsConfig("show clock", 120, "cwtest-topic-2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(cli, "cli_sensors.0.command", "show clock"),
					resource.TestCheckResourceAttr(cli, "cli_sensors.0.cadence", "120"),
					resource.TestCheckResourceAttr(cli, "cli_sensors.0.topic", "cwtest-topic-2"),
				),
			},
			{
				ResourceName:      cli,
				ImportState:       true,
				ImportStateId:     "cwtest-gen-app,cwtest-gen-cli",
				ImportStateVerify: true,
			},
			{
				ResourceName:      snmp,
				ImportState:       true,
				ImportStateId:     "cwtest-gen-app,cwtest-gen-snmp",
				ImportStateVerify: true,
			},
		},
	})
}

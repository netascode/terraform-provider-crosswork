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
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const collectionJobRes = "crosswork_collection_jobs.test"

func cliJobConfig(devices, destination, sensors string) string {
	return `
resource "crosswork_collection_jobs" "test" {
  application_id = "cwtest-app"
  context_id     = "cwtest-ctx"
  collector_type = "CLI"
  devices        = ` + devices + `
  destination    = "` + destination + `"
  cli_sensors    = [` + sensors + `]
}
`
}

const (
	cliSensorVersion = `{ command = "show version", cadence = 60, topic = "topic-a" }`
	cliSensorClock   = `{ command = "show clock", cadence = 30, topic = "topic-b" }`
)

func TestCollectionJobsLifecycle(t *testing.T) {
	fastCollectionWait(t)
	s := newCollectionStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// the cadence goes out in milliseconds; the devices come back in reverse order, which is no difference
				Config: s.providerBlock(false) + cliJobConfig(`["dev-1", "dev-2"]`, "cwtest-dest-kafka", cliSensorVersion+", "+cliSensorClock),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(collectionJobRes, "id", "cwtest-app,cwtest-ctx"),
					resource.TestCheckResourceAttr(collectionJobRes, "devices.#", "2"),
					resource.TestCheckTypeSetElemAttr(collectionJobRes, "devices.*", "dev-1"),
					resource.TestCheckResourceAttr(collectionJobRes, "destination", "cwtest-dest-kafka"),
					resource.TestCheckResourceAttr(collectionJobRes, "cli_sensors.#", "2"),
					resource.TestCheckResourceAttr(collectionJobRes, "cli_sensors.1.command", "show clock"),
					resource.TestCheckResourceAttr(collectionJobRes, "cli_sensors.1.cadence", "30"),
					resource.TestCheckResourceAttr(collectionJobRes, "cli_sensors.1.topic", "topic-b"),
					resource.TestCheckNoResourceAttr(collectionJobRes, "snmp_sensors"),
					resource.TestCheckNoResourceAttr(collectionJobRes, "device_tag"),
					func(*terraform.State) error {
						body := s.jobs["cwtest-app/cwtest-ctx"].body
						if got := body.Get("sensor_input_configs.0.cadence_in_millisec").String(); got != "60000" {
							t.Errorf("the cadence goes out in milliseconds, got %q", got)
						}
						if got := body.Get("sensor_output_configs.1.destination.context_id").String(); got != "topic-b" {
							t.Errorf("the topic is the context id of the destination of the paired output, got %q", got)
						}
						if got := body.Get("sensor_output_configs.0.destination.destination_id").String(); got != "dest-id-1" {
							t.Errorf("the destination is sent with its id, got %q", got)
						}
						return nil
					},
				),
			},
			{
				// the sensors and the destination change in place: one PUT, no delete
				PreConfig: s.resetCalls,
				Config:    s.providerBlock(false) + cliJobConfig(`["dev-2", "dev-1"]`, "cwtest-dest-plain", cliSensorClock),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(collectionJobRes, plancheck.ResourceActionUpdate),
				}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(collectionJobRes, "destination", "cwtest-dest-plain"),
					resource.TestCheckResourceAttr(collectionJobRes, "cli_sensors.#", "1"),
					resource.TestCheckResourceAttr(collectionJobRes, "cli_sensors.0.command", "show clock"),
					func(*terraform.State) error {
						if got := strings.Join(s.callLog(), ","); got != "PUT cwtest-app/cwtest-ctx" {
							t.Errorf("an in-place update is one PUT, got %s", got)
						}
						return nil
					},
				),
			},
			{
				ResourceName:      collectionJobRes,
				ImportState:       true,
				ImportStateId:     "cwtest-app,cwtest-ctx",
				ImportStateVerify: true,
			},
			{
				// a change of the devices replaces the job: the delete comes first and waits for the old job to leave the lists
				PreConfig: s.resetCalls,
				Config:    s.providerBlock(false) + cliJobConfig(`["dev-3"]`, "cwtest-dest-plain", cliSensorClock),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(collectionJobRes, plancheck.ResourceActionDestroyBeforeCreate),
				}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(collectionJobRes, "devices.#", "1"),
					func(*terraform.State) error {
						if got := strings.Join(s.callLog(), ","); got != "DELETE cwtest-app/cwtest-ctx,PUT cwtest-app/cwtest-ctx" {
							t.Errorf("a replace deletes and then creates, got %s", got)
						}
						return nil
					},
				),
			},
		},
	})
	if job := s.jobs["cwtest-app/cwtest-ctx"]; job == nil || !job.deleted {
		t.Errorf("the destroy deletes the job")
	}
}

func TestCollectionJobsSnmpByTag(t *testing.T) {
	fastCollectionWait(t)
	s := newCollectionStub(t)
	config := func(oid, op string) string {
		return s.providerBlock(false) + `
resource "crosswork_collection_jobs" "test" {
  application_id = "cwtest-app"
  context_id     = "cwtest-snmp"
  collector_type = "SNMP"
  device_tag     = "cwtest-tag"
  destination    = "cwtest-dest-kafka"
  snmp_sensors = [
    { oid = "` + oid + `", operation = "` + op + `", cadence = 120, topic = "topic-a" }
  ]
}
`
	}
	const res = "crosswork_collection_jobs.test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("1.3.6.1.2.1.1.5", "WALK"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "device_tag", "cwtest-tag"),
					resource.TestCheckNoResourceAttr(res, "devices"),
					resource.TestCheckResourceAttr(res, "snmp_sensors.0.oid", "1.3.6.1.2.1.1.5"),
					resource.TestCheckResourceAttr(res, "snmp_sensors.0.operation", "WALK"),
					resource.TestCheckResourceAttr(res, "snmp_sensors.0.cadence", "120"),
					func(*terraform.State) error {
						in := s.jobs["cwtest-app/cwtest-snmp"].body.Get("sensor_input_configs.0.sensor_data.snmp_sensor.snmp_mib")
						if got := in.Get("oids").Raw; got != `["1.3.6.1.2.1.1.5"]` {
							t.Errorf("one OID goes out as a list of one, got %s", got)
						}
						return nil
					},
				),
			},
			{
				Config: config("1.3.6.1.2.1.1.3.0", "SCALAR"),
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(res, plancheck.ResourceActionUpdate),
				}},
				Check: resource.TestCheckResourceAttr(res, "snmp_sensors.0.operation", "SCALAR"),
			},
		},
	})
}

func TestCollectionJobsRefusals(t *testing.T) {
	fastCollectionWait(t)
	s := newCollectionStub(t)
	unknownTag := `
resource "crosswork_collection_jobs" "test" {
  application_id = "cwtest-app"
  context_id     = "cwtest-ctx"
  collector_type = "CLI"
  device_tag     = "no-such-tag"
  destination    = "cwtest-dest-kafka"
  cli_sensors    = [` + cliSensorVersion + `]
}
`
	wrongSensor := `
resource "crosswork_collection_jobs" "test" {
  application_id = "cwtest-app"
  context_id     = "cwtest-ctx"
  collector_type = "SNMP"
  devices        = ["dev-1"]
  destination    = "cwtest-dest-kafka"
  cli_sensors    = [` + cliSensorVersion + `]
}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      s.providerBlock(false) + unknownTag,
				ExpectError: regexp.MustCompile(`(?s)group/tag no-such-tag not found in inventory or invalid`),
			},
			{
				Config:      s.providerBlock(false) + wrongSensor,
				ExpectError: regexp.MustCompile(`(?s)paths need to be of CLI_COLLECTOR`),
			},
			{
				Config:      s.providerBlock(false) + cliJobConfig(`["dev-1", "no-such-device"]`, "cwtest-dest-kafka", cliSensorVersion),
				ExpectError: regexp.MustCompile(`(?s)the device\s+no-such-device does not exist`),
			},
			{
				Config:      s.providerBlock(false) + cliJobConfig(`["dev-1"]`, "no-such-destination", cliSensorVersion),
				ExpectError: regexp.MustCompile(`(?s)There is no destination named no-such-destination`),
			},
		},
	})
	if len(s.jobs) != 0 {
		t.Errorf("a refused job leaves nothing behind, got %d", len(s.jobs))
	}
}

func TestCollectionJobsExistingJob(t *testing.T) {
	fastCollectionWait(t)
	s := newCollectionStub(t)
	s.seed("cwtest-app", "cwtest-ctx", "CLI_COLLECTOR", `{"device_set":{"devices":{"device_ids":["uuid-1"]}}}`, `{"cli_sensor":{"command":"show inventory"}}`)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// the same PUT would replace the job, so a create stops at a job it did not make
				Config:      s.providerBlock(false) + cliJobConfig(`["dev-1"]`, "cwtest-dest-kafka", cliSensorVersion),
				ExpectError: regexp.MustCompile(`(?s)A collection job cwtest-app,cwtest-ctx already exists`),
			},
			{
				Config: s.providerBlock(true) + cliJobConfig(`["dev-1"]`, "cwtest-dest-kafka", cliSensorVersion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(collectionJobRes, "cli_sensors.0.command", "show version"),
				),
			},
		},
	})
}

func TestCollectionJobsGoneOutsideTerraform(t *testing.T) {
	fastCollectionWait(t)
	s := newCollectionStub(t)
	config := s.providerBlock(false) + cliJobConfig(`["dev-1"]`, "cwtest-dest-kafka", cliSensorVersion)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{
				// a job deleted on the controller leaves the state on the next read and is created again
				PreConfig: func() { s.remove("cwtest-app", "cwtest-ctx") },
				Config:    config,
				ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{
					plancheck.ExpectResourceAction(collectionJobRes, plancheck.ResourceActionCreate),
				}},
			},
			{
				PreConfig: func() { s.remove("cwtest-app", "cwtest-ctx") },
				Config:    config,
				Destroy:   true,
			},
		},
	})
}

func TestCollectionJobsDeleteWaitsForTheListsToClear(t *testing.T) {
	fastCollectionWait(t)
	s := newCollectionStub(t)
	s.stale = 3
	s.seed("cwtest-app", "cwtest-ctx", "CLI_COLLECTOR", `{"device_set":{"device_group":"cwtest-tag"}}`, `{"cli_sensor":{"command":"show version"}}`)
	c, _ := registeredClient(s.URL, "admin", "stub-password", true, 0, "test")

	gone, err := deleteCollectionJob(context.Background(), c, "cwtest-app", "cwtest-ctx")
	if err != nil || !gone {
		t.Fatalf("gone %v, err %v", gone, err)
	}
	if left := s.jobs["cwtest-app/cwtest-ctx"].staleReads; left != 0 {
		t.Errorf("the delete returned while the configuration query still listed the job, %d stale reads left", left)
	}
	if gone, err := deleteCollectionJob(context.Background(), c, "cwtest-app", "cwtest-ctx"); err != nil || !gone {
		t.Errorf("deleting a missing job: gone %v, err %v", gone, err)
	}

	// a controller that keeps listing the job ends the wait without an error
	s.mu.Lock()
	s.stale = 1 << 30
	s.mu.Unlock()
	s.seed("cwtest-app", "cwtest-two", "CLI_COLLECTOR", `{"device_set":{"device_group":"cwtest-tag"}}`, `{"cli_sensor":{"command":"show version"}}`)
	if gone, err := deleteCollectionJob(withShortWait(), c, "cwtest-app", "cwtest-two"); err != nil || gone {
		t.Errorf("a job that stays listed: gone %v, err %v", gone, err)
	}
}

func withShortWait() context.Context {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	time.AfterFunc(time.Second, cancel)
	return ctx
}

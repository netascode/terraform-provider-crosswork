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
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/tidwall/gjson"
)

const odnRes = "crosswork_srte_odn_templates.test"

func odnConfig(s *odnStub, maxTimeout int, name, body string) string {
	return s.providerBlock(maxTimeout) + fmt.Sprintf(`
resource "crosswork_srte_odn_templates" "test" {
  name  = %q
%s
}
`, name, body)
}

const odnFull = `
  color             = 4242
  maximum_sid_depth = 8
  head_end = [{ name = "PE-2" }, { name = "PE-1" }]
  dynamic_pce         = true
  dynamic_metric_type = "igp"
  dynamic_affinity_rules = [
    { action = "include-any", color = ["red", "blue"] },
    { action = "exclude-any", color = ["green"] },
  ]
  delay_measurement_profile                = "delay-profile-1"
  delay_measurement_logging_delay_exceeded = true
`

func TestSRTEODNTemplatesLifecycle(t *testing.T) {
	s := newODNStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: odnConfig(s, 180, "cwtest-odn", odnFull),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(odnRes, "id", "cwtest-odn"),
					resource.TestCheckResourceAttr(odnRes, "color", "4242"),
					resource.TestCheckResourceAttr(odnRes, "head_end.#", "2"),
					resource.TestCheckResourceAttr(odnRes, "dynamic_pce", "true"),
					resource.TestCheckResourceAttr(odnRes, "delay_measurement_logging_delay_exceeded", "true"),
					resource.TestCheckResourceAttr(odnRes, "dynamic_affinity_rules.#", "2"),
					func(*terraform.State) error {
						body := gjson.Parse(s.bodies[0]).Get("cisco-sr-te-cfp-sr-odn:odn-template.0")
						check := func(path, want string) {
							if got := body.Get(path).Raw; got != want {
								t.Errorf("POST body %s = %s, want %s\n%s", path, got, want, s.bodies[0])
							}
						}
						check("name", `"cwtest-odn"`)
						check("dynamic.pce", `{}`)                                                          // a presence container
						check("performance-measurement.delay-measurement.logging.delay-exceeded", `[null]`) // an empty leaf
						check("color", `4242`)
						if got := s.planReadCount["cwtest-odn"]; got < 1 {
							t.Errorf("the create did not read the plan")
						}
						return nil
					},
				),
			},
			{
				// no change after the apply: the set of head ends, the sorted colors and the flags read back as they were sent
				Config:   odnConfig(s, 180, "cwtest-odn", odnFull),
				PlanOnly: true,
			},
			{
				// a head end is dropped, the color changes and a flag is switched off: one PUT with the whole object
				Config: odnConfig(s, 180, "cwtest-odn", strings.NewReplacer(`{ name = "PE-2" }, `, "", "4242", "4243", "  delay_measurement_logging_delay_exceeded = true", "").Replace(odnFull)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(odnRes, "color", "4243"),
					resource.TestCheckResourceAttr(odnRes, "head_end.#", "1"),
					resource.TestCheckNoResourceAttr(odnRes, "delay_measurement_logging_delay_exceeded"),
					func(*terraform.State) error {
						body := gjson.Parse(s.bodies[len(s.bodies)-1]).Get("cisco-sr-te-cfp-sr-odn:odn-template.0")
						if body.Get("name").String() != "cwtest-odn" || body.Get("head-end.#").Int() != 1 || body.Get("performance-measurement.delay-measurement.logging").Exists() {
							t.Errorf("the PUT is not the whole object: %s", s.bodies[len(s.bodies)-1])
						}
						return nil
					},
				),
			},
		},
		CheckDestroy: func(*terraform.State) error {
			if _, ok := s.objects["cwtest-odn"]; ok {
				return fmt.Errorf("the template was not deleted")
			}
			if got := s.planReadCount["cwtest-odn"]; got < 4 {
				return fmt.Errorf("the delete did not wait for the plan to go: %d plan reads in all", got)
			}
			return nil
		},
	})
}

func TestSRTEODNTemplatesKeyIsEscapedAsPathSegment(t *testing.T) {
	s := newODNStub(t)
	const name = "cwtest,odn%1"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: odnConfig(s, 180, name, "  color = 4242"),
			Check: func(*terraform.State) error {
				if _, ok := s.objects[name]; !ok {
					return fmt.Errorf("the template is not stored under its name, objects: %v", s.objects)
				}
				var got []string
				for _, c := range s.calls {
					if strings.Contains(c, "odn-template") {
						got = append(got, c)
					}
				}
				for _, c := range got {
					if !strings.Contains(c, "=cwtest%2Codn%251") {
						return fmt.Errorf("the key is not escaped as a path segment in %q", c)
					}
				}
				if len(got) == 0 {
					return fmt.Errorf("no call to the template")
				}
				return nil
			},
		}},
	})
}

func TestSRTEODNTemplatesDeviceRefusalFailsTheApply(t *testing.T) {
	s := newODNStub(t)
	s.refuse = "External error in the NED implementation for device PE-1: segment-routing traffic-eng on-demand color 4242 srv6 !!% Invalid argument: Unsupported behavior: 71"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      odnConfig(s, 180, "cwtest-odn", "  color = 4242"),
			ExpectError: regexp.MustCompile(`Unsupported behavior: 71`),
		}},
	})
	if len(s.objects) != 0 {
		t.Errorf("nothing is created when the device refuses: %v", s.objects)
	}
}

func TestSRTEODNTemplatesWaitTimeout(t *testing.T) {
	s := newODNStub(t)
	s.neverReady = true
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      odnConfig(s, 1, "cwtest-odn", "  color = 4242"),
			ExpectError: regexp.MustCompile(`(?s)timed out waiting for the NSO\s+plan`),
		}},
		// the template was created before the wait, so it is in the state (tainted) and the destroy removes it
		CheckDestroy: func(*terraform.State) error {
			if _, ok := s.objects["cwtest-odn"]; ok {
				return fmt.Errorf("the template stayed in NSO: it was not kept in the state")
			}
			return nil
		},
	})
}

func TestSRTEODNTemplatesFailedRemovalIsAnError(t *testing.T) {
	s := newODNStub(t)
	s.failRemoval = true
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: odnConfig(s, 180, "cwtest-odn", "  color = 4242")},
			{
				Config:      odnConfig(s, 180, "cwtest-odn", "  color = 4242"),
				Destroy:     true,
				ExpectError: regexp.MustCompile(`(?s)NSO plan failed.*TSDN-SR-301.*device PE-1.*recommended action`),
			},
		},
	})
}

func TestSRTEODNTemplatesGoneObjectIsRecreated(t *testing.T) {
	s := newODNStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: odnConfig(s, 180, "cwtest-odn", "  color = 4242")},
			{
				PreConfig: func() {
					s.mu.Lock()
					delete(s.objects, "cwtest-odn")
					s.mu.Unlock()
				},
				Config: odnConfig(s, 180, "cwtest-odn", "  color = 4242"),
				Check: func(*terraform.State) error {
					if _, ok := s.objects["cwtest-odn"]; !ok {
						return fmt.Errorf("the template was not created again")
					}
					return nil
				},
			},
		},
	})
}

// A head end that NSO does not know yet is waited for before anything is written; with the timeout 0 of the stub the
// create fails at once, names the device and writes nothing, and the next apply goes through when the device is there.
func TestSRTEODNTemplatesWaitsForTheHeadEnd(t *testing.T) {
	s := newODNStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      odnConfig(s, 180, "cwtest-odn", "  color = 4242\n  head_end = [{ name = \"PE-3\" }]"),
				ExpectError: regexp.MustCompile(`PE-3: not in the inventory`),
				Check: func(*terraform.State) error {
					if len(s.bodies) != 0 {
						return fmt.Errorf("something was written before the device was ready: %v", s.bodies)
					}
					return nil
				},
			},
			{
				PreConfig: func() {
					s.mu.Lock()
					s.devices["PE-3"] = "SYNC_FROM_STARTED"
					s.mu.Unlock()
				},
				Config:      odnConfig(s, 180, "cwtest-odn", "  color = 4242\n  head_end = [{ name = \"PE-3\" }]"),
				ExpectError: regexp.MustCompile(`PE-3: reachability CONN_STATE_REACHABLE, nso SYNC_FROM_STARTED`),
			},
			{
				PreConfig: func() {
					s.mu.Lock()
					s.devices["PE-3"] = nsoSynced
					s.mu.Unlock()
				},
				Config: odnConfig(s, 180, "cwtest-odn", "  color = 4242\n  head_end = [{ name = \"PE-3\" }]"),
				Check:  resource.TestCheckResourceAttr(odnRes, "head_end.#", "1"),
			},
		},
	})
}

// An update waits only for the head ends it adds: a change of the color goes through while the device of an
// existing head end is out of sync, and a new head end that NSO does not have synced stops the update.
func TestSRTEODNTemplatesUpdateWaitsOnlyForAddedHeadEnds(t *testing.T) {
	s := newODNStub(t)
	setState := func(name, state string) func() {
		return func() {
			s.mu.Lock()
			s.devices[name] = state
			s.mu.Unlock()
		}
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: odnConfig(s, 180, "cwtest-odn", "  color = 4242\n  head_end = [{ name = \"PE-1\" }]"),
			},
			{
				PreConfig: setState("PE-1", "NOT_SYNCED"),
				Config:    odnConfig(s, 180, "cwtest-odn", "  color = 4243\n  head_end = [{ name = \"PE-1\" }]"),
				Check:     resource.TestCheckResourceAttr(odnRes, "color", "4243"),
			},
			{
				PreConfig:   setState("PE-2", "NOT_SYNCED"),
				Config:      odnConfig(s, 180, "cwtest-odn", "  color = 4243\n  head_end = [{ name = \"PE-1\" }, { name = \"PE-2\" }]"),
				ExpectError: regexp.MustCompile(`PE-2: reachability CONN_STATE_REACHABLE, nso NOT_SYNCED`),
			},
		},
	})
}

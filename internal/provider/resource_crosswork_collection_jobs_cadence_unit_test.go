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
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const cadenceRes = "crosswork_collection_jobs_cadence.test"

func cadenceConfig(values string) string {
	return `
resource "crosswork_collection_jobs_cadence" "test" {
` + values + `
}
`
}

func TestCollectionJobsCadenceLifecycle(t *testing.T) {
	fastCadenceWait(t)
	s := newCadenceStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// two values are set: the update still names all four, the other two with the values the cluster holds, and
				// the state keeps only the two that the configuration names
				Config: s.providerBlock() + cadenceConfig("  snmp_check = 900\n  show_clock = 1200"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(cadenceRes, "id", "collection_jobs_cadence"),
					resource.TestCheckResourceAttr(cadenceRes, "snmp_check", "900"),
					resource.TestCheckResourceAttr(cadenceRes, "show_clock", "1200"),
					resource.TestCheckNoResourceAttr(cadenceRes, "reachability_check"),
					resource.TestCheckNoResourceAttr(cadenceRes, "te_tunnel_id"),
					func(*terraform.State) error {
						var puts []string
						for _, c := range s.callLog() {
							if strings.HasPrefix(c, "PUT ") {
								puts = append(puts, c)
							}
						}
						want := `PUT {"reach-check":600,"snmp":900,"te-tunnel-id":1200,"show-clock":1200}`
						if len(puts) != 1 || puts[0] != want {
							t.Errorf("one update with all four values, got %v, want %s", puts, want)
						}
						return nil
					},
				),
			},
			{
				// a value the configuration leaves out is not a difference when the cluster holds another one, and a change of
				// one value is one update that keeps the others
				PreConfig: func() { s.mu.Lock(); s.values["reach-check"] = 900; s.mu.Unlock(); s.resetCalls() },
				Config:    s.providerBlock() + cadenceConfig("  snmp_check = 1500\n  show_clock = 1200"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(cadenceRes, "snmp_check", "1500"),
					resource.TestCheckNoResourceAttr(cadenceRes, "reachability_check"),
					func(*terraform.State) error {
						if got, want := s.value("snmp"), int64(1500); got != want {
							t.Errorf("snmp is %d, want %d", got, want)
						}
						if got, want := s.value("reach-check"), int64(900); got != want {
							t.Errorf("a value the configuration leaves out keeps the value of the cluster, got %d, want %d", got, want)
						}
						return nil
					},
				),
			},
			{
				ResourceName:  cadenceRes,
				ImportState:   true,
				ImportStateId: "collection_jobs_cadence",
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 {
						t.Fatalf("one imported state, got %d", len(states))
					}
					a := states[0].Attributes
					for attr, want := range map[string]string{"reachability_check": "900", "snmp_check": "1500", "te_tunnel_id": "1200", "show_clock": "1200"} {
						if a[attr] != want {
							t.Errorf("%s is %q, want %s", attr, a[attr], want)
						}
					}
					return nil
				},
			},
		},
	})
}

func TestCollectionJobsCadenceNamesAllFourValues(t *testing.T) {
	fastCadenceWait(t)
	s := newCadenceStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: s.providerBlock() + cadenceConfig("  reachability_check = 300\n  snmp_check = 600\n  te_tunnel_id = 86400\n  show_clock = 7200"),
				Check: func(*terraform.State) error {
					for name, want := range map[string]int64{"reach-check": 300, "snmp": 600, "te-tunnel-id": 86400, "show-clock": 7200} {
						if got := s.value(name); got != want {
							t.Errorf("%s is %d, want %d", name, got, want)
						}
					}
					return nil
				},
			},
		},
	})
}

func TestCollectionJobsCadenceDestroyLeavesTheValues(t *testing.T) {
	fastCadenceWait(t)
	s := newCadenceStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: s.providerBlock() + cadenceConfig("  snmp_check = 900"),
			},
			{
				PreConfig: s.resetCalls,
				Config:    s.providerBlock() + cadenceConfig("  snmp_check = 900"),
				Destroy:   true,
				Check: func(*terraform.State) error {
					for _, c := range s.callLog() {
						if strings.HasPrefix(c, "PUT ") {
							t.Errorf("a destroy sends no update, got %s", c)
						}
					}
					if got := s.value("snmp"); got != 900 {
						t.Errorf("the value stays on the cluster after a destroy, got %d", got)
					}
					return nil
				},
			},
		},
	})
}

func TestCollectionJobsCadenceRefusedValue(t *testing.T) {
	fastCadenceWait(t)
	s := newCadenceStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// the controller answers HTTP 200 with a failed job; its text is returned as it is
				PreConfig: func() {
					s.mu.Lock()
					s.reject = "Invalid cadence value for job(s). For reachability job, the cadence must be between 300 seconds"
					s.mu.Unlock()
				},
				Config:      s.providerBlock() + cadenceConfig("  reachability_check = 300"),
				ExpectError: regexp.MustCompile(`Invalid cadence value for job\(s\)\. For reachability job`),
			},
		},
	})
}

func TestCollectionJobsCadenceValueShowsLate(t *testing.T) {
	fastCadenceWait(t)
	s := newCadenceStub(t)
	s.lag = 3

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// the new value shows after several reads, which the update waits for
				Config: s.providerBlock() + cadenceConfig("  show_clock = 3600"),
				Check:  resource.TestCheckResourceAttr(cadenceRes, "show_clock", "3600"),
			},
		},
	})
}

func TestCollectionJobsCadenceValueNeverShows(t *testing.T) {
	fastCadenceWait(t)
	s := newCadenceStub(t)
	s.frozen = true

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      s.providerBlock() + cadenceConfig("  show_clock = 3600"),
				ExpectError: regexp.MustCompile(`(?s)still shows.*show-clock=1800.*show-clock=3600`),
			},
		},
	})
}

func TestCollectionJobsCadenceRanges(t *testing.T) {
	s := newCadenceStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// the reachability check may be as short as 5 minutes, and the others only 10
				Config:      s.providerBlock() + cadenceConfig("  snmp_check = 300"),
				ExpectError: regexp.MustCompile(`(?s)snmp_check.*between 600 and 86400`),
			},
			{
				Config:      s.providerBlock() + cadenceConfig("  reachability_check = 299"),
				ExpectError: regexp.MustCompile(`(?s)reachability_check.*between 300 and 86400`),
			},
			{
				Config:      s.providerBlock() + cadenceConfig("  show_clock = 86401"),
				ExpectError: regexp.MustCompile(`(?s)show_clock.*between 600 and 86400`),
			},
		},
	})
}

func TestCollectionJobsCadenceImportIdentifier(t *testing.T) {
	s := newCadenceStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:        s.providerBlock() + cadenceConfig(""),
				ResourceName:  cadenceRes,
				ImportState:   true,
				ImportStateId: "something-else",
				ExpectError:   regexp.MustCompile(`Expected import identifier with format: collection_jobs_cadence`),
			},
		},
	})
}

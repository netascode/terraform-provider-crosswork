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
)

const lcmThrRes = "crosswork_local_congestion_mitigation_interface_thresholds.test"

const (
	n1i14 = "cwtest-node-1 GigabitEthernet0/0/0/14"
	n1i15 = "cwtest-node-1 GigabitEthernet0/0/0/15"
	n2i14 = "cwtest-node-2 GigabitEthernet0/0/0/14"
)

func lcmThrConfig(domain string, rows ...string) string {
	var b strings.Builder
	for _, r := range rows {
		f := strings.Fields(r)
		fmt.Fprintf(&b, "    { node = %q, interface = %q, utilization_threshold = %s },\n", f[0], f[1], f[2])
	}
	return fmt.Sprintf("\nresource \"crosswork_local_congestion_mitigation_interface_thresholds\" \"test\" {\n  domain_id = %q\n  interface_thresholds = [\n%s  ]\n}\n", domain, b.String())
}

func lcmThrKeys(rows []lcmThreshold) map[string]int64 {
	m := map[string]int64{}
	for _, r := range rows {
		m[r.node+" "+r.iface] = r.threshold
	}
	return m
}

func lcmThrCheckRows(s *lcmThresholdStub, want map[string]int64) resource.TestCheckFunc {
	return func(*terraform.State) error {
		got := lcmThrKeys(s.held("32"))
		if len(got) != len(want) {
			return fmt.Errorf("the controller holds %v, want %v", got, want)
		}
		for k, v := range want {
			if got[k] != v {
				return fmt.Errorf("the controller holds %v, want %v", got, want)
			}
		}
		return nil
	}
}

func TestLocalCongestionMitigationInterfaceThresholdsLifecycle(t *testing.T) {
	s := newLcmThresholdStub(t, "32")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             lcmThrCheckRows(s, map[string]int64{}),
		Steps: []resource.TestStep{
			{
				Config: s.providerBlock() + lcmThrConfig("32", n1i14+" 75", n1i15+" 80"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(lcmThrRes, "id", "32"),
					resource.TestCheckResourceAttr(lcmThrRes, "interface_thresholds.#", "2"),
					lcmThrCheckRows(s, map[string]int64{n1i14: 75, n1i15: 80}),
					func(*terraform.State) error {
						// the replace flag and the empty errors list of every row are sent
						call := s.callLog()[0]
						if !strings.Contains(call, `"merge":false`) || !strings.Contains(call, `"errors":[]`) {
							return fmt.Errorf("unexpected body: %s", call)
						}
						return nil
					},
				),
			},
			{
				// a changed threshold moves its row to the end on the controller; the plan stays empty
				Config: s.providerBlock() + lcmThrConfig("32", n1i14+" 90", n1i15+" 80", n2i14+" 60"),
				Check:  lcmThrCheckRows(s, map[string]int64{n1i14: 90, n1i15: 80, n2i14: 60}),
			},
			{
				// a row that is not configured any more is removed by the replacing write
				Config: s.providerBlock() + lcmThrConfig("32", n1i15+" 80"),
				Check:  lcmThrCheckRows(s, map[string]int64{n1i15: 80}),
			},
			{
				// an empty list cannot be written: the rows read are deleted
				Config: s.providerBlock() + lcmThrConfig("32"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(lcmThrRes, "interface_thresholds.#", "0"),
					lcmThrCheckRows(s, map[string]int64{}),
					func(*terraform.State) error {
						if last := s.callLog()[len(s.callLog())-1]; !strings.HasPrefix(last, "DELETE ") {
							return fmt.Errorf("the last call was %s", last)
						}
						return nil
					},
				),
			},
			{
				Config: s.providerBlock() + lcmThrConfig("32", n1i14+" 75"),
				Check:  lcmThrCheckRows(s, map[string]int64{n1i14: 75}),
			},
			{
				ResourceName:                         lcmThrRes,
				ImportState:                          true,
				ImportStateId:                        "32",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "domain_id",
			},
		},
	})
}

func TestLocalCongestionMitigationInterfaceThresholdsRefusedRows(t *testing.T) {
	s := newLcmThresholdStub(t, "32")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// the texts of the refused rows are returned, and nothing is written, not even the good row
				Config:      s.providerBlock() + lcmThrConfig("32", n1i14+" 75", "cwtest-node-1 GigabitEthernet0/0/0/0 75", n1i15+" 5"),
				ExpectError: regexp.MustCompile(`(?s)REJECTED.*not\s+found.*hold\s+margin`),
			},
		},
	})
	if got := s.held("32"); len(got) != 0 {
		t.Fatalf("a refused write changed the list: %v", got)
	}
}

func TestLocalCongestionMitigationInterfaceThresholdsUnknownDomain(t *testing.T) {
	s := newLcmThresholdStub(t, "32")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      s.providerBlock() + lcmThrConfig("999", n1i14+" 75"),
				ExpectError: regexp.MustCompile(`Request timeout`),
			},
		},
	})
}

func TestLocalCongestionMitigationInterfaceThresholdsDrift(t *testing.T) {
	s := newLcmThresholdStub(t, "32")
	config := s.providerBlock() + lcmThrConfig("32", n1i14+" 75", n1i15+" 80")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{
				// a row added in the UI is seen
				PreConfig: func() {
					s.set("32", lcmThreshold{"cwtest-node-1", "GigabitEthernet0/0/0/14", 75}, lcmThreshold{"cwtest-node-1", "GigabitEthernet0/0/0/15", 80}, lcmThreshold{"cwtest-node-2", "GigabitEthernet0/0/0/14", 50})
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// a row removed in the UI is seen
				PreConfig: func() {
					s.set("32", lcmThreshold{"cwtest-node-1", "GigabitEthernet0/0/0/14", 75})
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				// the apply restores the list
				Config: config,
				Check:  lcmThrCheckRows(s, map[string]int64{n1i14: 75, n1i15: 80}),
			},
		},
	})
}

func TestLocalCongestionMitigationInterfaceThresholdsRemovedOutsideTerraform(t *testing.T) {
	s := newLcmThresholdStub(t, "32")
	config := s.providerBlock() + lcmThrConfig("32", n1i14+" 75")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{
				// every row deleted in the UI: the list of the state is no longer the configured one
				PreConfig:          func() { s.set("32") },
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

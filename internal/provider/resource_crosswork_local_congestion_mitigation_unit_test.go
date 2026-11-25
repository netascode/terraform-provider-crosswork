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

const lcmRes = "crosswork_local_congestion_mitigation.test"

func lcmConfig(domain, body string) string {
	return fmt.Sprintf("\nresource \"crosswork_local_congestion_mitigation\" \"test\" {\n  domain_id = %q\n%s}\n", domain, body)
}

func TestLocalCongestionMitigationLifecycle(t *testing.T) {
	s := newLcmStub(t, "200")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			// a destroy writes the defaults of the controller for the domain
			c := s.config("200")
			if c.Get("enable").Bool() || c.Get("color").Int() != 2000 || len(c.Get("nodeTagMaxSidDepth").Array()) != 0 ||
				len(c.Get("affinityConstraint.includeAny").Array()) != 0 || c.Get("description").String() != "LCM startup config" {
				return fmt.Errorf("the destroy did not write the defaults: %s", c.Raw)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: s.providerBlock() + lcmConfig("200", `  color = 3000
  affinity_include_any = ["30"]
  affinity_include_all = ["30", "29"]
  node_tag_max_sid_depth = { GRPC_MSL = "10", RTP = "5" }
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(lcmRes, "id", "200"),
					resource.TestCheckResourceAttr(lcmRes, "domain_id", "200"),
					resource.TestCheckResourceAttr(lcmRes, "color", "3000"),
					resource.TestCheckResourceAttr(lcmRes, "utilization_threshold", "100"),
					resource.TestCheckResourceAttr(lcmRes, "affinity_include_any.#", "1"),
					resource.TestCheckResourceAttr(lcmRes, "affinity_include_any.0", "30"),
					resource.TestCheckResourceAttr(lcmRes, "affinity_include_all.1", "29"),
					resource.TestCheckResourceAttr(lcmRes, "node_tag_max_sid_depth.GRPC_MSL", "10"),
					resource.TestCheckResourceAttr(lcmRes, "node_tag_max_sid_depth.RTP", "5"),
					func(*terraform.State) error {
						// the hops are written as one element that holds every tag, the lists as bit strings
						c := s.config("200")
						if len(c.Get("nodeTagMaxSidDepth").Array()) != 1 || c.Get("nodeTagMaxSidDepth.0.GRPC_MSL").String() != "10" ||
							c.Get("affinityConstraint.includeAll.1").String() != "29" {
							return fmt.Errorf("unexpected body: %s", c.Raw)
						}
						return nil
					},
				),
			},
			{
				// the unset lists and hops are written empty, and the plan is empty afterwards
				Config: s.providerBlock() + lcmConfig("200", "  color = 3000\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(lcmRes, "node_tag_max_sid_depth.GRPC_MSL"),
					func(*terraform.State) error {
						c := s.config("200")
						if len(c.Get("nodeTagMaxSidDepth").Array()) != 0 || len(c.Get("affinityConstraint.includeAny").Array()) != 0 {
							return fmt.Errorf("the empty blocks were not written empty: %s", c.Raw)
						}
						return nil
					},
				),
			},
			{
				ResourceName:                         lcmRes,
				ImportState:                          true,
				ImportStateId:                        "200",
				ImportStateVerify:                    true,
				ImportStateVerifyIdentifierAttribute: "domain_id",
			},
		},
	})
}

func TestLocalCongestionMitigationRefusedWrite(t *testing.T) {
	s := newLcmStub(t, "200")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// a refusal is a failure status inside HTTP 200; the text of the controller is returned
				Config:      s.providerBlock() + lcmConfig("200", "  utilization_threshold = 50\n  utilization_hold_margin = 60\n"),
				ExpectError: regexp.MustCompile(`utilization threshold must be greater than or equal to the hold margin`),
			},
		},
	})
}

func TestLocalCongestionMitigationUnknownDomain(t *testing.T) {
	s := newLcmStub(t, "200")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      s.providerBlock() + lcmConfig("999", ""),
				ExpectError: regexp.MustCompile(`msg.{1,4}is null`),
			},
		},
	})
	if puts := s.putLog(); len(puts) == 0 || !strings.HasPrefix(puts[0], "999 ") {
		t.Fatalf("the write was not tried for the unknown domain: %v", puts)
	}
}

func TestLocalCongestionMitigationDecimalIsDrift(t *testing.T) {
	s := newLcmStub(t, "200")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: s.providerBlock() + lcmConfig("200", "  overprovision_factor = 3\n"),
				Check:  resource.TestCheckResourceAttr(lcmRes, "overprovision_factor", "3"),
			},
			{
				// a decimal set in the UI reads as its whole part: the plan shows it and the apply writes the whole number
				PreConfig: func() {
					s.mu.Lock()
					s.fraction["200"] = 2.5
					s.mu.Unlock()
				},
				Config:             s.providerBlock() + lcmConfig("200", "  overprovision_factor = 3\n"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: s.providerBlock() + lcmConfig("200", "  overprovision_factor = 3\n"),
				Check:  resource.TestCheckResourceAttr(lcmRes, "overprovision_factor", "3"),
			},
		},
	})
}

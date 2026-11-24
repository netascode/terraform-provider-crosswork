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

// Section below is generated&owned by "gen/generator.go". //template:begin imports
import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc
func TestAccCwBandwidthOnDemand(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "enable", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "primary_objective", "MIN_METRIC"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "link_utilization_threshold", "85"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "reoptimization_interval", "61"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "metric_reoptimization_hours", "2"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "metric_reoptimization_minutes", "15"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "sr_policy_traffic", "MEASURED"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "model_update_throttle", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "optimizer_event_threshold", "50"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "optimizer_run_time_limit", "70"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "allowed_policy_violations", "STRICT"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "prefer_strict_sids", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "debug_optimizer", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_bandwidth_on_demand.test", "debug_optimizer_max_files", "31"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwBandwidthOnDemandConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_bandwidth_on_demand.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{},
		ImportStateId:           "bandwidth_on_demand",
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// End of section. //template:end testAcc

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwBandwidthOnDemandConfig_minimum() string {
	config := `resource "crosswork_bandwidth_on_demand" "test" {` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwBandwidthOnDemandConfig_all() string {
	config := `resource "crosswork_bandwidth_on_demand" "test" {` + "\n"
	config += `	enable = false` + "\n"
	config += `	primary_objective = "MIN_METRIC"` + "\n"
	config += `	link_utilization_threshold = 85` + "\n"
	config += `	reoptimization_interval = 61` + "\n"
	config += `	metric_reoptimization_hours = 2` + "\n"
	config += `	metric_reoptimization_minutes = 15` + "\n"
	config += `	sr_policy_traffic = "MEASURED"` + "\n"
	config += `	model_update_throttle = 5` + "\n"
	config += `	optimizer_event_threshold = 50` + "\n"
	config += `	optimizer_run_time_limit = 70` + "\n"
	config += `	allowed_policy_violations = "STRICT"` + "\n"
	config += `	prefer_strict_sids = true` + "\n"
	config += `	debug_optimizer = true` + "\n"
	config += `	debug_optimizer_max_files = 31` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

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
func TestAccCwLocalCongestionMitigation(t *testing.T) {
	if os.Getenv("OPTIMIZATION") == "" {
		t.Skip("skipping test, set environment variable OPTIMIZATION")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "enable", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "color", "2000"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "utilization_threshold", "100"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "utilization_hold_margin", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "delete_lsps", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "profile_id", "0"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "congestion_check_interval", "900"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "max_tte_group", "8"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "include_all_interfaces", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "description", "LCM startup config"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "auto_repair_solution", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "stay_in_area", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "adjacency_hop_type", "UNPROTECTED"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "operation_mode", "MANUAL"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "objective", "MIN_IGP_METRIC"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "pce_dispatch_timeout", "180"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "congestion_check_suspension_interval", "600"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "overprovision_factor", "0"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "uneven_ecmp_traffic_threshold", "0"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "throttle_mode_threshold", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "history_retention_time", "30"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "debug_opt", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "debug_opt_max_files", "40"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "affinity_include_any.0", "30"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "affinity_exclude_any.0", "30"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_local_congestion_mitigation.test", "affinity_include_all.0", "30"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwLocalCongestionMitigationConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_local_congestion_mitigation.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"domain_id"},
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
func testAccCwLocalCongestionMitigationConfig_minimum() string {
	config := `resource "crosswork_local_congestion_mitigation" "test" {` + "\n"
	config += `	domain_id = "200"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwLocalCongestionMitigationConfig_all() string {
	config := `resource "crosswork_local_congestion_mitigation" "test" {` + "\n"
	config += `	domain_id = "200"` + "\n"
	config += `	enable = false` + "\n"
	config += `	color = 2000` + "\n"
	config += `	utilization_threshold = 100` + "\n"
	config += `	utilization_hold_margin = 5` + "\n"
	config += `	delete_lsps = false` + "\n"
	config += `	profile_id = 0` + "\n"
	config += `	congestion_check_interval = 900` + "\n"
	config += `	max_tte_group = 8` + "\n"
	config += `	include_all_interfaces = false` + "\n"
	config += `	description = "LCM startup config"` + "\n"
	config += `	auto_repair_solution = false` + "\n"
	config += `	stay_in_area = false` + "\n"
	config += `	adjacency_hop_type = "UNPROTECTED"` + "\n"
	config += `	operation_mode = "MANUAL"` + "\n"
	config += `	objective = "MIN_IGP_METRIC"` + "\n"
	config += `	pce_dispatch_timeout = 180` + "\n"
	config += `	congestion_check_suspension_interval = 600` + "\n"
	config += `	overprovision_factor = 0` + "\n"
	config += `	uneven_ecmp_traffic_threshold = 0` + "\n"
	config += `	throttle_mode_threshold = 5` + "\n"
	config += `	history_retention_time = 30` + "\n"
	config += `	debug_opt = false` + "\n"
	config += `	debug_opt_max_files = 40` + "\n"
	config += `	affinity_include_any = ["30"]` + "\n"
	config += `	affinity_exclude_any = ["30"]` + "\n"
	config += `	affinity_include_all = ["30"]` + "\n"
	config += `	node_tag_max_sid_depth = { GRPC_MSL = "10" }` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

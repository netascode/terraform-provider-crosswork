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
func TestAccCwSRTEODNTemplates(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "name", "cwtest-odn"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "color", "4242"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "head_end.0.name", "build-router-1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "maximum_sid_depth", "8"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "pce_group", "cwtest-pce-group"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "bandwidth", "1000"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "source_address", "192.0.2.1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "dynamic_pce", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "dynamic_flex_alg", "129"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "dynamic_metric_type", "igp"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "dynamic_metric_margin_absolute", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "dynamic_affinity_rules.0.action", "include-any"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "dynamic_segments_sid_algorithm", "129"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "dynamic_disjoint_path_type", "node"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "dynamic_disjoint_path_group_id", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "dynamic_disjoint_path_sub_id", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "delay_measurement_profile", "delay-profile-1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "delay_measurement_logging_delay_exceeded", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_srte_odn_templates.test", "performance_measurement_reverse_path_label", "2000"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccCwSRTEODNTemplatesConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccCwSRTEODNTemplatesConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
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
func testAccCwSRTEODNTemplatesConfig_minimum() string {
	config := `resource "crosswork_srte_odn_templates" "test" {` + "\n"
	config += `	name = "cwtest-odn"` + "\n"
	config += `	color = 4242` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwSRTEODNTemplatesConfig_all() string {
	config := `resource "crosswork_srte_odn_templates" "test" {` + "\n"
	config += `	name = "cwtest-odn"` + "\n"
	config += `	color = 4242` + "\n"
	config += `	head_end = [{` + "\n"
	config += `	  name = "build-router-1"` + "\n"
	config += `	}]` + "\n"
	config += `	maximum_sid_depth = 8` + "\n"
	config += `	pce_group = "cwtest-pce-group"` + "\n"
	config += `	bandwidth = 1000` + "\n"
	config += `	source_address = "192.0.2.1"` + "\n"
	config += `	dynamic_pce = true` + "\n"
	config += `	dynamic_flex_alg = 129` + "\n"
	config += `	dynamic_metric_type = "igp"` + "\n"
	config += `	dynamic_metric_margin_absolute = 5` + "\n"
	config += `	dynamic_affinity_rules = [{` + "\n"
	config += `	  action = "include-any"` + "\n"
	config += `	  color = ["red"]` + "\n"
	config += `	}]` + "\n"
	config += `	dynamic_segments_sid_algorithm = 129` + "\n"
	config += `	dynamic_disjoint_path_type = "node"` + "\n"
	config += `	dynamic_disjoint_path_group_id = 5` + "\n"
	config += `	dynamic_disjoint_path_sub_id = 1` + "\n"
	config += `	delay_measurement_profile = "delay-profile-1"` + "\n"
	config += `	delay_measurement_logging_delay_exceeded = true` + "\n"
	config += `	performance_measurement_reverse_path_label = 2000` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

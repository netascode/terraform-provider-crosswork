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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSource
func TestAccDataSourceCwSRTEPolicies(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "name", "cwtest-policy"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "head_end", "build-router-1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "tail_end", "192.0.2.61"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "color", "999"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "binding_sid", "24001"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "bandwidth", "1000"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "source_address", "192.0.2.1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "delay_measurement_profile", "delay-profile-1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "delay_measurement_logging_delay_exceeded", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "performance_measurement_reverse_path_label", "2000"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "auto_route_metric_relative", "-5"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "auto_route_include_prefixes.0.prefix_address", "10.250.0.0/16"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "auto_route_force_sr_include", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "auto_route_forward_class", "2"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "paths.0.preference", "100"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "paths.0.explicit_sid_lists.0.weight", "10"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "paths.0.explicit_affinity_rules.0.action", "include-any"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "paths.0.explicit_segments_sid_algorithm", "128"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "paths.0.explicit_disjoint_path_type", "link"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "paths.0.explicit_disjoint_path_group_id", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_srte_policies.test", "paths.0.explicit_disjoint_path_sub_id", "1"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwSRTEPoliciesPrerequisitesConfig + testAccDataSourceCwSRTEPoliciesConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccDataSourceCwSRTEPoliciesPrerequisitesConfig = `
resource "crosswork_srte_sid_lists" "test" {
  name = "cwtest-policy-sid-list"
  sids = [{ index = 1, label = 16001 }]
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwSRTEPoliciesConfig() string {
	config := `resource "crosswork_srte_policies" "test" {` + "\n"
	config += `	name = "cwtest-policy"` + "\n"
	config += `	head_end = "build-router-1"` + "\n"
	config += `	tail_end = "192.0.2.61"` + "\n"
	config += `	color = 999` + "\n"
	config += `	binding_sid = 24001` + "\n"
	config += `	bandwidth = 1000` + "\n"
	config += `	source_address = "192.0.2.1"` + "\n"
	config += `	delay_measurement_profile = "delay-profile-1"` + "\n"
	config += `	delay_measurement_logging_delay_exceeded = true` + "\n"
	config += `	performance_measurement_reverse_path_label = 2000` + "\n"
	config += `	auto_route_metric_relative = -5` + "\n"
	config += `	auto_route_include_prefixes = [{` + "\n"
	config += `	  prefix_address = "10.250.0.0/16"` + "\n"
	config += `	}]` + "\n"
	config += `	auto_route_force_sr_include = true` + "\n"
	config += `	auto_route_forward_class = 2` + "\n"
	config += `	paths = [{` + "\n"
	config += `	  preference = 100` + "\n"
	config += `	  explicit_sid_lists = [{` + "\n"
	config += `		name = crosswork_srte_sid_lists.test.name` + "\n"
	config += `		weight = 10` + "\n"
	config += `	}]` + "\n"
	config += `	  explicit_affinity_rules = [{` + "\n"
	config += `		action = "include-any"` + "\n"
	config += `		colors = ["red"]` + "\n"
	config += `	}]` + "\n"
	config += `	  explicit_segments_sid_algorithm = 128` + "\n"
	config += `	  explicit_disjoint_path_type = "link"` + "\n"
	config += `	  explicit_disjoint_path_group_id = 1` + "\n"
	config += `	  explicit_disjoint_path_sub_id = 1` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_srte_policies" "test" {
			id = crosswork_srte_policies.test.id
			name = "cwtest-policy"
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

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

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSource
func TestAccDataSourceCwDemandMatrixGlobalConfiguration(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_demand_matrix_global_configuration.test", "computation_cadence", "30"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_demand_matrix_global_configuration.test", "minimum_demand_traffic", "100"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_demand_matrix_global_configuration.test", "debug_mode", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_demand_matrix_global_configuration.test", "debug_max_plan_files", "5"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwDemandMatrixGlobalConfigurationConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwDemandMatrixGlobalConfigurationConfig() string {
	config := `resource "crosswork_demand_matrix_global_configuration" "test" {` + "\n"
	config += `	computation_cadence = 30` + "\n"
	config += `	minimum_demand_traffic = 100` + "\n"
	config += `	debug_mode = true` + "\n"
	config += `	debug_max_plan_files = 5` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_demand_matrix_global_configuration" "test" {
			depends_on = [crosswork_demand_matrix_global_configuration.test]
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

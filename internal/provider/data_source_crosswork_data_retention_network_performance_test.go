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
func TestAccDataSourceCwDataRetentionNetworkPerformance(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_retention_network_performance.test", "collect_lsp_state_change", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_retention_network_performance.test", "collect_lsp_path_change", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_retention_network_performance.test", "collect_lsp_utilization", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_retention_network_performance.test", "collect_lsp_pm", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_retention_network_performance.test", "collect_link_pm", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_retention_network_performance.test", "raw_hours", "12"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_retention_network_performance.test", "hourly_days", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_retention_network_performance.test", "daily_days", "20"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_retention_network_performance.test", "weekly_weeks", "30"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwDataRetentionNetworkPerformanceConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwDataRetentionNetworkPerformanceConfig() string {
	config := `resource "crosswork_data_retention_network_performance" "test" {` + "\n"
	config += `	collect_lsp_state_change = false` + "\n"
	config += `	collect_lsp_path_change = false` + "\n"
	config += `	collect_lsp_utilization = false` + "\n"
	config += `	collect_lsp_pm = true` + "\n"
	config += `	collect_link_pm = true` + "\n"
	config += `	raw_hours = 12` + "\n"
	config += `	hourly_days = 5` + "\n"
	config += `	daily_days = 20` + "\n"
	config += `	weekly_weeks = 30` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_data_retention_network_performance" "test" {
			depends_on = [crosswork_data_retention_network_performance.test]
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

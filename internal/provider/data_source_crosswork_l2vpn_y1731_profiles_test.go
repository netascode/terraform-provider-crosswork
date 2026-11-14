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
func TestAccDataSourceCwL2VPNY1731Profiles(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "name", "cwtest-y1731"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "type", "delay"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "schedule_interval", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "schedule_duration", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "probe_type", "burst"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "probe_message_count", "10"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "probe_message_period", "100"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "probe_measurement_interval", "60"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "probe_frame_size", "1000"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "probe_priority", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "delay_statistics.0.type", "delay-two-way"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "bucket_size", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_y1731_profiles.test", "bucket_archive", "20"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwL2VPNY1731ProfilesConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwL2VPNY1731ProfilesConfig() string {
	config := `resource "crosswork_l2vpn_y1731_profiles" "test" {` + "\n"
	config += `	name = "cwtest-y1731"` + "\n"
	config += `	type = "delay"` + "\n"
	config += `	schedule_interval = 5` + "\n"
	config += `	schedule_duration = "5"` + "\n"
	config += `	probe_type = "burst"` + "\n"
	config += `	probe_message_count = 10` + "\n"
	config += `	probe_message_period = 100` + "\n"
	config += `	probe_measurement_interval = "60"` + "\n"
	config += `	probe_frame_size = 1000` + "\n"
	config += `	probe_priority = 5` + "\n"
	config += `	delay_statistics = [{` + "\n"
	config += `	  type = "delay-two-way"` + "\n"
	config += `	}]` + "\n"
	config += `	bucket_size = 1` + "\n"
	config += `	bucket_archive = 20` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_l2vpn_y1731_profiles" "test" {
			id = crosswork_l2vpn_y1731_profiles.test.id
			name = "cwtest-y1731"
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

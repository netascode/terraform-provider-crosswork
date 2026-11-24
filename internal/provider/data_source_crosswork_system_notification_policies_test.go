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
func TestAccDataSourceCwSystemNotificationPolicies(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_system_notification_policies.test", "name", "cwtest-policy"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_system_notification_policies.test", "description", "Forward all alarms"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_system_notification_policies.test", "criteria", "*"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwSystemNotificationPoliciesPrerequisitesConfig + testAccDataSourceCwSystemNotificationPoliciesConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccDataSourceCwSystemNotificationPoliciesPrerequisitesConfig = `
resource "crosswork_trap_destination" "anp_test" {
  name                 = "cwtest-anp-trap"
  address              = "192.0.2.131"
  community_wo         = "cwtest-community"
  community_wo_version = 1
}

resource "crosswork_syslog_destination" "anp_test" {
  name    = "cwtest-anp-syslog"
  address = "192.0.2.132"
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwSystemNotificationPoliciesConfig() string {
	config := `resource "crosswork_system_notification_policies" "test" {` + "\n"
	config += `	name = "cwtest-policy"` + "\n"
	config += `	description = "Forward all alarms"` + "\n"
	config += `	criteria = "*"` + "\n"
	config += `	trap_receivers = [{` + "\n"
	config += `	  name = crosswork_trap_destination.anp_test.name` + "\n"
	config += `	}]` + "\n"
	config += `	syslog_receivers = [{` + "\n"
	config += `	  name = crosswork_syslog_destination.anp_test.name` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_system_notification_policies" "test" {
			id = crosswork_system_notification_policies.test.id
			name = "cwtest-policy"
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

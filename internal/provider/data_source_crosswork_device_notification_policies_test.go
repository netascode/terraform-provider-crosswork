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
func TestAccDataSourceCwDeviceNotificationPolicies(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_device_notification_policies.test", "name", "cwtest-dnp"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_device_notification_policies.test", "description", "Critical alarms"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwDeviceNotificationPoliciesPrerequisitesConfig + testAccDataSourceCwDeviceNotificationPoliciesConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccDataSourceCwDeviceNotificationPoliciesPrerequisitesConfig = `
resource "crosswork_trap_destination" "dnp_test" {
  name                 = "cwtest-dnp-trap"
  address              = "192.0.2.133"
  community_wo         = "cwtest-community"
  community_wo_version = 1
}

resource "crosswork_syslog_destination" "dnp_test" {
  name    = "cwtest-dnp-syslog"
  address = "192.0.2.134"
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwDeviceNotificationPoliciesConfig() string {
	config := `resource "crosswork_device_notification_policies" "test" {` + "\n"
	config += `	name = "cwtest-dnp"` + "\n"
	config += `	description = "Critical alarms"` + "\n"
	config += `	filters = [{` + "\n"
	config += `	  severities = ["critical"]` + "\n"
	config += `	}]` + "\n"
	config += `	trap_receivers = [{` + "\n"
	config += `	  name = crosswork_trap_destination.dnp_test.name` + "\n"
	config += `	}]` + "\n"
	config += `	syslog_receivers = [{` + "\n"
	config += `	  name = crosswork_syslog_destination.dnp_test.name` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_device_notification_policies" "test" {
			id = crosswork_device_notification_policies.test.id
			name = "cwtest-dnp"
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

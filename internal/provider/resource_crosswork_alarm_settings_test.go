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
func TestAccCwAlarmSettings(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_alarm_settings.test", "security_device_alarms_cleared_days", "31"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_alarm_settings.test", "non_security_device_alarms_cleared_days", "8"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_alarm_settings.test", "all_events_days", "31"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_alarm_settings.test", "device_alarms_days", "61"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_alarm_settings.test", "system_alarms_days", "62"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_alarm_settings.test", "network_alarms_days", "63"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_alarm_settings.test", "audit_logs_days", "31"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_alarm_settings.test", "trap_processing", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_alarm_settings.test", "syslog_processing", "false"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwAlarmSettingsConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_alarm_settings.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{},
		ImportStateId:           "alarm_settings",
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
func testAccCwAlarmSettingsConfig_minimum() string {
	config := `resource "crosswork_alarm_settings" "test" {` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwAlarmSettingsConfig_all() string {
	config := `resource "crosswork_alarm_settings" "test" {` + "\n"
	config += `	security_device_alarms_cleared_days = 31` + "\n"
	config += `	non_security_device_alarms_cleared_days = 8` + "\n"
	config += `	all_events_days = 31` + "\n"
	config += `	device_alarms_days = 61` + "\n"
	config += `	system_alarms_days = 62` + "\n"
	config += `	network_alarms_days = 63` + "\n"
	config += `	audit_logs_days = 31` + "\n"
	config += `	trap_processing = false` + "\n"
	config += `	syslog_processing = false` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

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
func TestAccCwAlarmSeverities(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_alarm_severities.test", "severities.0.name", "L2-ELMI-4-AIPC"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_alarm_severities.test", "severities.0.severity", "minor"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_alarm_severities.test", "severities.0.auto_clear_minutes", "1440"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwAlarmSeveritiesConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_alarm_severities.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"severities"},
		ImportStateId:           "alarm_severities",
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
func testAccCwAlarmSeveritiesConfig_minimum() string {
	config := `resource "crosswork_alarm_severities" "test" {` + "\n"
	config += `	severities = [{` + "\n"
	config += `	  name = "L2-ELMI-4-AIPC"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwAlarmSeveritiesConfig_all() string {
	config := `resource "crosswork_alarm_severities" "test" {` + "\n"
	config += `	severities = [{` + "\n"
	config += `	  name = "L2-ELMI-4-AIPC"` + "\n"
	config += `	  severity = "minor"` + "\n"
	config += `	  auto_clear_minutes = 1440` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

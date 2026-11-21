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
func TestAccCwAAASessionSettings(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_session_settings.test", "fallback_type", "FB_ON_NO_RESPONSE"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_session_settings.test", "idle_session_timeout", "60"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_session_settings.test", "parallel_sessions", "100"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_session_settings.test", "parallel_sessions_per_user", "100"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_session_settings.test", "source_ip_auditing", "true"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwAAASessionSettingsConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_aaa_session_settings.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{},
		ImportStateId:           "aaa_session_settings",
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
func testAccCwAAASessionSettingsConfig_minimum() string {
	config := `resource "crosswork_aaa_session_settings" "test" {` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwAAASessionSettingsConfig_all() string {
	config := `resource "crosswork_aaa_session_settings" "test" {` + "\n"
	config += `	fallback_type = "FB_ON_NO_RESPONSE"` + "\n"
	config += `	idle_session_timeout = 60` + "\n"
	config += `	parallel_sessions = 100` + "\n"
	config += `	parallel_sessions_per_user = 100` + "\n"
	config += `	source_ip_auditing = true` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

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
func TestAccDataSourceCwAAAPasswordPolicy(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "min_password_length", "9"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "no_username", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "no_cisco_variant", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "no_char_repetition", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "num_changed_chars_enabled", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "num_changed_chars", "2"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "num_reuse_limit_enabled", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "num_reuse_limit", "6"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "password_reuse_days_enabled", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "password_reuse_days", "16"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "failed_logins_before_lockout_enabled", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "failed_logins_before_lockout", "6"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "lock_out_user_time_enabled", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "lock_out_user_time", "7"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "password_expiry_days_enabled", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "password_expiry_days", "61"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "days_for_warning_enabled", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "days_for_warning", "16"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_password_policy.test", "change_password_on_first_login", "true"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwAAAPasswordPolicyConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwAAAPasswordPolicyConfig() string {
	config := `resource "crosswork_aaa_password_policy" "test" {` + "\n"
	config += `	min_password_length = 9` + "\n"
	config += `	no_username = true` + "\n"
	config += `	no_cisco_variant = true` + "\n"
	config += `	no_char_repetition = true` + "\n"
	config += `	num_changed_chars_enabled = true` + "\n"
	config += `	num_changed_chars = 2` + "\n"
	config += `	num_reuse_limit_enabled = true` + "\n"
	config += `	num_reuse_limit = 6` + "\n"
	config += `	password_reuse_days_enabled = true` + "\n"
	config += `	password_reuse_days = 16` + "\n"
	config += `	failed_logins_before_lockout_enabled = true` + "\n"
	config += `	failed_logins_before_lockout = 6` + "\n"
	config += `	lock_out_user_time_enabled = true` + "\n"
	config += `	lock_out_user_time = 7` + "\n"
	config += `	password_expiry_days_enabled = true` + "\n"
	config += `	password_expiry_days = 61` + "\n"
	config += `	days_for_warning_enabled = true` + "\n"
	config += `	days_for_warning = 16` + "\n"
	config += `	change_password_on_first_login = true` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_aaa_password_policy" "test" {
			depends_on = [crosswork_aaa_password_policy.test]
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

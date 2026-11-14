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

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc
func TestAccCwLockUsers(t *testing.T) {
	var checks []resource.TestCheckFunc

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwLockUsersPrerequisitesConfig + testAccCwLockUsersConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_lock_users.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// End of section. //template:end testAcc

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccCwLockUsersPrerequisitesConfig = `
resource "crosswork_users" "test" {
  username            = "cwtest-lockuser"
  password_wo         = "CwtestPass123!"
  password_wo_version = 1
  policy_id           = "admin"
  first_name          = "Test"
  last_name           = "User"
  device_access_groups = [{
    uuid        = "b84e1a37-ef0d-4b79-984c-b162249d2e23"
    domain_name = "ALL-ACCESS"
  }]
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwLockUsersConfig_minimum() string {
	config := `resource "crosswork_lock_users" "test" {` + "\n"
	config += `	username = crosswork_users.test.username` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwLockUsersConfig_all() string {
	config := `resource "crosswork_lock_users" "test" {` + "\n"
	config += `	username = crosswork_users.test.username` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

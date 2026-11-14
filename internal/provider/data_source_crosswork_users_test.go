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
func TestAccDataSourceCwUsers(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_users.test", "username", "cwtest-gen-user"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_users.test", "policy_id", "admin"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_users.test", "first_name", "Test"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_users.test", "last_name", "User"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_users.test", "device_access_groups.0.uuid", "b84e1a37-ef0d-4b79-984c-b162249d2e23"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_users.test", "device_access_groups.0.domain_name", "ALL-ACCESS"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwUsersConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwUsersConfig() string {
	config := `resource "crosswork_users" "test" {` + "\n"
	config += `	username = "cwtest-gen-user"` + "\n"
	config += `	password_wo = "CwtestPass123!"` + "\n"
	config += `	password_wo_version = 1` + "\n"
	config += `	policy_id = "admin"` + "\n"
	config += `	first_name = "Test"` + "\n"
	config += `	last_name = "User"` + "\n"
	config += `	device_access_groups = [{` + "\n"
	config += `	  uuid = "b84e1a37-ef0d-4b79-984c-b162249d2e23"` + "\n"
	config += `	  domain_name = "ALL-ACCESS"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_users" "test" {
			id = crosswork_users.test.id
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

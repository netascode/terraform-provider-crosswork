// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceCwRoles(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: testAccCwRoleConfig(false) + `
data "crosswork_roles" "test" {
	id = crosswork_roles.test.id
}
`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("data.crosswork_roles.test", "name", "cwtest-role"),
				resource.TestCheckResourceAttr("data.crosswork_roles.test", "api_permissions.cw-fault-events-api", "RW"),
				resource.TestCheckResourceAttr("data.crosswork_roles.test", "api_permissions.inventory_cwinventory", "R"),
				resource.TestCheckResourceAttr("data.crosswork_roles.test", "task_permissions.#", "3"),
			),
		}},
	})
}

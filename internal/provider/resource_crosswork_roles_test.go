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

func TestAccCwRoles(t *testing.T) {
	const res = "crosswork_roles.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{ // read-only grants are stored with POST next to GET; they must still read back as R
				Config: testAccCwRoleConfig(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "name", "cwtest-role"),
					resource.TestCheckResourceAttr(res, "api_permissions.%", "3"),
					resource.TestCheckResourceAttr(res, "api_permissions.cw-fault-events-api", "RW"),
					resource.TestCheckResourceAttr(res, "api_permissions.inventory_cwinventory", "R"),
					resource.TestCheckResourceAttr(res, "task_permissions.#", "3"),
				),
			},
			{
				Config: testAccCwRoleConfig(true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "api_permissions.inventory_cwinventory", "RW"),
					resource.TestCheckResourceAttr(res, "task_permissions.#", "4"),
				),
			},
			{ // back to the first set: the task goes first, then the API is narrowed
				Config: testAccCwRoleConfig(false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "api_permissions.inventory_cwinventory", "R"),
					resource.TestCheckResourceAttr(res, "task_permissions.#", "3"),
				),
			},
			{
				ResourceName:  res,
				ImportState:   true,
				ImportStateId: "cwtest-role",
			},
		},
	})
}

func testAccCwRoleConfig(more bool) string {
	inventory := "R"
	tasks := `"id_export_audit_logs_access", "id_view_audit_logs_access", "id_nso_fp_deployment_management"`
	if more {
		inventory = "RW"
		tasks += `, "id_provisioning"`
	}
	config := `resource "crosswork_roles" "test" {` + "\n"
	config += `	name = "cwtest-role"` + "\n"
	config += `	api_permissions = {` + "\n"
	config += `		"cw-fault-events-api"   = "RW"` + "\n"
	config += `		"nso-fp-dep-mngr"       = "RWD"` + "\n"
	config += `		"inventory_cwinventory" = "` + inventory + `"` + "\n"
	config += `	}` + "\n"
	config += `	task_permissions = [` + tasks + `]` + "\n"
	config += `}` + "\n"
	return config
}

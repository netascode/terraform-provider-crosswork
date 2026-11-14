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

// The catalog depends on what is installed on the cluster, so the test checks only entries that every
// cluster has.
func TestAccDataSourceCwRoleCatalog(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: `data "crosswork_role_catalog" "test" {}`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckTypeSetElemNestedAttrs("data.crosswork_role_catalog.test", "apis.*", map[string]string{
					"id": "cw-fault-events-api", "category": "Alarms and Events",
				}),
				resource.TestCheckTypeSetElemNestedAttrs("data.crosswork_role_catalog.test", "tasks.*", map[string]string{
					"id": "id_view_audit_logs_access", "category": "Audit Logs",
				}),
				resource.TestCheckTypeSetElemNestedAttrs("data.crosswork_role_catalog.test", "task_api_permissions.*", map[string]string{
					"task_id": "id_view_audit_logs_access", "api_id": "cw-fault-events-api", "access": "R",
				}),
			),
		}},
	})
}

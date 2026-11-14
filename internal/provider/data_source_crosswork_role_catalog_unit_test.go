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
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// API categories are an object of category name to a list of {api_id, name}; tasks are a list of groups with items; the
// task-to-API map lists, per task, the API ids with their access letters.
type catalogStub struct {
	*ssoStub
	failAPIs bool
	calls    []string
}

const catalogAPIsBody = `{
  "Alarms and Events": [
    {"api_id": "cw-fault-events-api", "name": "Alarms & Events"},
    {"api_id": "alarm-rest-service", "name": "Alarms & Events"}
  ],
  "Inventory": [{"api_id": "inventory_cwinventory", "name": "Inventory APIs"}],
  "Workflow Manager": [{"api_id": "cwm_adapter", "name": "Adapter Handler"}]
}`

const catalogTasksBody = `[
  {"id": "audit_logs", "name": "Audit Logs", "items": [
    {"id": "id_export_audit_logs_access", "name": "Export Audit Logs", "enabled": true},
    {"id": "id_view_audit_logs_access", "name": "View Audit Logs", "enabled": true}
  ]},
  {"id": "coe", "name": "Crosswork Optimization Engine", "items": [
    {"id": "id_bwod_config", "name": "Bandwidth on Demand Configuration", "enabled": true}
  ]}
]`

const catalogTaskAPIsBody = `{
  "id_export_audit_logs_access": {"apiIds": {"cw-fault-events-api": ["W", "R"]}},
  "id_view_audit_logs_access": {"apiIds": {"cw-fault-events-api": ["R"]}}
}`

func newCatalogStub(t *testing.T) *catalogStub {
	t.Helper()
	s := &catalogStub{ssoStub: newSSOStub(t)}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/crosswork/aaa/") {
			sso.ServeHTTP(w, r)
			return
		}
		s.calls = append(s.calls, r.Method+" "+r.URL.Path)
		switch r.URL.Path {
		case "/crosswork/aaa/v2/api":
			if s.failAPIs {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, `{"error":"boom"}`)
				return
			}
			fmt.Fprint(w, catalogAPIsBody)
		case "/crosswork/aaa/v1/task/admin":
			fmt.Fprint(w, catalogTasksBody)
		case "/crosswork/aaa/v1/taskAPIPermission/admin":
			fmt.Fprint(w, catalogTaskAPIsBody)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	t.Cleanup(Shutdown)
	return s
}

func (s *catalogStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

func TestRoleCatalogDataSource(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newCatalogStub(t)
	const ds = "data.crosswork_role_catalog.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: s.providerBlock() + "\ndata \"crosswork_role_catalog\" \"test\" {}\n",
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(ds, "id", "role_catalog"),
				resource.TestCheckResourceAttr(ds, "apis.#", "4"),
				resource.TestCheckTypeSetElemNestedAttrs(ds, "apis.*", map[string]string{
					"id": "cw-fault-events-api", "name": "Alarms & Events", "category": "Alarms and Events",
				}),
				resource.TestCheckTypeSetElemNestedAttrs(ds, "apis.*", map[string]string{
					"id": "inventory_cwinventory", "name": "Inventory APIs", "category": "Inventory",
				}),
				resource.TestCheckResourceAttr(ds, "tasks.#", "3"),
				resource.TestCheckTypeSetElemNestedAttrs(ds, "tasks.*", map[string]string{
					"id": "id_view_audit_logs_access", "name": "View Audit Logs", "category": "Audit Logs",
				}),
				resource.TestCheckTypeSetElemNestedAttrs(ds, "tasks.*", map[string]string{
					"id": "id_bwod_config", "name": "Bandwidth on Demand Configuration", "category": "Crosswork Optimization Engine",
				}),
				// only tasks that need APIs are listed, and the letters come back in the order R, W, D
				resource.TestCheckResourceAttr(ds, "task_api_permissions.#", "2"),
				resource.TestCheckTypeSetElemNestedAttrs(ds, "task_api_permissions.*", map[string]string{
					"task_id": "id_export_audit_logs_access", "api_id": "cw-fault-events-api", "access": "RW",
				}),
				resource.TestCheckTypeSetElemNestedAttrs(ds, "task_api_permissions.*", map[string]string{
					"task_id": "id_view_audit_logs_access", "api_id": "cw-fault-events-api", "access": "R",
				}),
			),
		}},
	})
}

func TestRoleCatalogDataSourceError(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newCatalogStub(t)
	s.failAPIs = true

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      s.providerBlock() + "\ndata \"crosswork_role_catalog\" \"test\" {}\n",
			ExpectError: regexp.MustCompile(`Failed to retrieve the APIs`),
		}},
	})
}

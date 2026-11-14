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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Rules the tests rely on: base APIs added to every role, POST stored next to GET for read-only grants,
// a task rejected unless the role holds the APIs it needs, and an update that drops an API a task needs accepted.
type roleStub struct {
	*ssoStub

	mu      sync.Mutex
	roles   map[string]map[string]map[string]any
	tasks   map[string]map[string]bool
	calls   []string
	creates int
}

var roleStubAPIs = map[string]string{
	"cw-fault-events-api":   "Fault Events",
	"nso-fp-dep-mngr":       "NSO Deployment Manager",
	"inventory_cwinventory": "Inventory",
	"aaa_cw_role_read":      "Know my role",
	"aaa_cwpassword":        "Change password",
	"aaa_selected_pref":     "Preferences",
}

var roleStubTasks = map[string]map[string]string{
	"id_export_audit_logs_access":     {"cw-fault-events-api": "RW"},
	"id_view_audit_logs_access":       {"cw-fault-events-api": "R"},
	"id_nso_fp_deployment_management": {"nso-fp-dep-mngr": "RWD"},
	"id_provisioning":                 {"inventory_cwinventory": "RW"},
	"id_bwod_config":                  {},
}

func newRoleStub(t *testing.T) *roleStub {
	t.Helper()
	s := &roleStub{ssoStub: newSSOStub(t), roles: map[string]map[string]map[string]any{}, tasks: map[string]map[string]bool{}}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/crosswork/aaa/v1/") {
			s.serve(w, r)
			return
		}
		sso.ServeHTTP(w, r)
	})
	t.Cleanup(Shutdown)
	return s
}

func roleEntry(id string, methods ...string) map[string]any {
	ms := make([]any, 0, len(methods))
	for _, m := range methods {
		ms = append(ms, m)
	}
	return map[string]any{"api_name": roleStubAPIs[id], "api_id": id, "versions": []any{},
		"allowed_urls": []any{map[string]any{"url": "/.*", "methods": ms}}}
}

func holds(entry map[string]any, level string) bool {
	have := map[string]bool{}
	for _, u := range entry["allowed_urls"].([]any) {
		for _, m := range u.(map[string]any)["methods"].([]any) {
			have[m.(string)] = true
		}
	}
	for _, m := range methodsFor(level) {
		if !have[m] {
			return false
		}
	}
	return true
}

var postWithGet = map[string]bool{"cw-fault-events-api": true, "nso-fp-dep-mngr": true, "inventory_cwinventory": true}

func withBase(sent map[string]any) map[string]map[string]any {
	out := map[string]map[string]any{}
	for id, e := range sent {
		entry := e.(map[string]any)
		if postWithGet[id] && holds(entry, "R") && !holds(entry, "W") {
			url0 := entry["allowed_urls"].([]any)[0].(map[string]any)
			url0["methods"] = append(url0["methods"].([]any), "POST")
		}
		out[id] = entry
	}
	out["aaa_cw_role_read"] = roleEntry("aaa_cw_role_read", "GET")
	out["aaa_cwpassword"] = roleEntry("aaa_cwpassword", "PUT")
	out["aaa_selected_pref"] = roleEntry("aaa_selected_pref", "GET", "PUT")
	return out
}

func (s *roleStub) serve(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	s.mu.Lock()
	defer s.mu.Unlock()
	rest := strings.TrimPrefix(r.URL.Path, "/crosswork/aaa/v1/")
	if (r.Method == http.MethodPost || r.Method == http.MethodPut) && r.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	name := ""
	if i := strings.Index(rest, "/"); i >= 0 {
		name = rest[i+1:]
	}
	json200 := func(v any) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}

	switch {
	case r.Method == http.MethodGet && rest == "api":
		ids := make([]string, 0, len(roleStubAPIs))
		for id := range roleStubAPIs {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		list := []any{}
		for _, id := range ids {
			list = append(list, map[string]any{"api_id": id, "name": roleStubAPIs[id], "listen_path": "/x"})
		}
		json200(list)
	case r.Method == http.MethodGet && rest == "role":
		out := map[string]any{}
		for n, rights := range s.roles {
			out[n] = map[string]any{"id": n, "access_rights": rights}
		}
		json200(out)
	case r.Method == http.MethodPost && rest == "role":
		s.calls = append(s.calls, "POST role")
		var wrapped map[string]map[string]any
		if err := json.Unmarshal(raw, &wrapped); err != nil || len(wrapped) != 1 {
			userFail(w, 422, "Failed to create role: json: cannot unmarshal string into Go value of type user.Policy")
			return
		}
		for n, body := range wrapped {
			rights, _ := body["access_rights"].(map[string]any)
			switch {
			case s.roles[n] != nil:
				userFail(w, 500, "Role with same name already exist.")
			case len(rights) == 0:
				userFail(w, 500, "Role permissions cannot be empty.")
			default:
				s.roles[n] = withBase(rights)
				s.creates++
				w.WriteHeader(http.StatusCreated)
			}
		}
	case r.Method == http.MethodPut && strings.HasPrefix(rest, "role/"):
		s.calls = append(s.calls, "PUT role")
		var body map[string]any
		_ = json.Unmarshal(raw, &body)
		rights, _ := body["access_rights"].(map[string]any)
		if s.roles[name] == nil {
			userFail(w, 500, "Failed to update role: Role doesn't exist.")
			return
		}
		if len(rights) == 0 {
			userFail(w, 500, "Failed to update role: Role permissions cannot be empty")
			return
		}
		next := withBase(rights)
		s.roles[name] = next
		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodDelete && strings.HasPrefix(rest, "role/"):
		s.calls = append(s.calls, "DELETE role")
		if s.roles[name] == nil {
			userFail(w, 500, "Failed to delete role: Role doesn't exist.")
			return
		}
		delete(s.roles, name)
		delete(s.tasks, name)
		w.WriteHeader(http.StatusNoContent)
	case r.Method == http.MethodGet && strings.HasPrefix(rest, "task/"):
		group := func(id, label string, items ...string) any {
			list := []any{}
			for _, item := range items {
				list = append(list, map[string]any{"id": item, "name": item, "enabled": s.tasks[name][item]})
			}
			return map[string]any{"id": id, "name": label, "items": list}
		}
		json200([]any{
			group("audit_logs", "Audit Logs", "id_export_audit_logs_access", "id_view_audit_logs_access"),
			group("coe", "Optimization", "id_bwod_config"),
			group("crosswork_network_controller", "Network Controller", "id_provisioning"),
			group("nso_management", "NSO", "id_nso_fp_deployment_management"),
		})
	case r.Method == http.MethodPost && strings.HasPrefix(rest, "task/"):
		s.calls = append(s.calls, "POST task")
		var list []map[string]string
		_ = json.Unmarshal(raw, &list)
		set := map[string]bool{}
		for _, item := range list {
			id := item["taskId"]
			needs, known := roleStubTasks[id]
			if !known {
				userFail(w, 500, "Failed to update Taskpq: insert or update on table \"policy_task\" violates foreign key constraint \"fk_task\"")
				return
			}
			for api, level := range needs {
				if e, ok := s.roles[name][api]; !ok || !holds(e, level) {
					userFail(w, 500, fmt.Sprintf("Failed to update Task: %s. Associated APIs %s - %s [%s] do not have sufficient permissions mandated by the task", id, id, api, level))
					return
				}
			}
			set[id] = true
		}
		if set["id_export_audit_logs_access"] {
			set["id_view_audit_logs_access"] = true
		}
		s.tasks[name] = set
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *roleStub) enabledTasks(name string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []string{}
	for id, on := range s.tasks[name] {
		if on {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}

func (s *roleStub) roleAPIs(name string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []string{}
	for id := range s.roles[name] {
		out = append(out, id)
	}
	sort.Strings(out)
	return out
}

func (s *roleStub) callLog() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return strings.Join(s.calls, ",")
}

const cwtestRole = "cwtest-unit-role"

func (s *roleStub) roleConfig(perms string, tasks ...string) string {
	quoted := make([]string, 0, len(tasks))
	for _, t := range tasks {
		quoted = append(quoted, fmt.Sprintf("%q", t))
	}
	return s.providerBlock() + fmt.Sprintf(`
resource "crosswork_roles" "test" {
  name            = %q
  api_permissions = {%s}
  task_permissions = [%s]
}
`, cwtestRole, perms, strings.Join(quoted, ", "))
}

func (s *roleStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

func TestRolesLifecycle(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newRoleStub(t)
	const res = "crosswork_roles.test"
	emptyPlan := resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}}
	const base = `"cw-fault-events-api" = "RW", "nso-fp-dep-mngr" = "RWD", "inventory_cwinventory" = "R"`
	const export, view, nso, prov = "id_export_audit_logs_access", "id_view_audit_logs_access", "id_nso_fp_deployment_management", "id_provisioning"

	logIs := func(want string) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if got := s.callLog(); got != want {
				return fmt.Errorf("calls were %q, want %q", got, want)
			}
			return nil
		}
	}
	equal := func(what string, got func() []string, want ...string) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if g := got(); strings.Join(g, ",") != strings.Join(want, ",") {
				return fmt.Errorf("%s are %v, want %v", what, g, want)
			}
			return nil
		}
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			if len(s.roleAPIs(cwtestRole)) != 0 {
				return fmt.Errorf("role %s still exists on the stub", cwtestRole)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{ // create: the base APIs and the task the server enables by itself cause no difference
				Config: s.roleConfig(base, export, nso),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", cwtestRole),
					resource.TestCheckResourceAttr(res, "api_permissions.%", "3"),
					resource.TestCheckResourceAttr(res, "api_permissions.inventory_cwinventory", "R"),
					resource.TestCheckResourceAttr(res, "api_permissions.cw-fault-events-api", "RW"),
					resource.TestCheckResourceAttr(res, "api_permissions.nso-fp-dep-mngr", "RWD"),
					resource.TestCheckResourceAttr(res, "task_permissions.#", "2"),
					equal("APIs on the stub", func() []string { return s.roleAPIs(cwtestRole) },
						"aaa_cw_role_read", "aaa_cwpassword", "aaa_selected_pref", "cw-fault-events-api", "inventory_cwinventory", "nso-fp-dep-mngr"),
					equal("tasks on the stub", func() []string { return s.enabledTasks(cwtestRole) }, export, nso, view),
					logIs("POST role,POST task"),
				),
				ConfigPlanChecks: emptyPlan,
			},
			{ // widen an API and add a task that needs it: the API first, then the tasks
				Config: s.roleConfig(`"cw-fault-events-api" = "RW", "nso-fp-dep-mngr" = "RWD", "inventory_cwinventory" = "RW"`, export, nso, prov),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "task_permissions.#", "3"),
					logIs("POST role,POST task,PUT role,POST task"),
				),
				ConfigPlanChecks: emptyPlan,
			},
			{ // revoke an API and the tasks that need it: the tasks first, then the API (Crosswork does not enforce this order)
				Config: s.roleConfig(`"cw-fault-events-api" = "RW"`, export),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "api_permissions.%", "1"),
					resource.TestCheckResourceAttr(res, "task_permissions.#", "1"),
					equal("tasks on the stub", func() []string { return s.enabledTasks(cwtestRole) }, export, view),
					logIs("POST role,POST task,PUT role,POST task,POST task,PUT role"),
				),
				ConfigPlanChecks: emptyPlan,
			},
			{ // import by name reads what is beyond the base APIs, which Crosswork adds by itself
				ResourceName:  res,
				ImportState:   true,
				ImportStateId: cwtestRole,
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					a := states[0].Attributes
					if a["name"] != cwtestRole || a["api_permissions.cw-fault-events-api"] != "RW" || a["api_permissions.%"] != "1" {
						return fmt.Errorf("unexpected imported state: %v", a)
					}
					return nil
				},
			},
			{ // deleted on the controller: the read finds no role, so the apply creates it again
				PreConfig: func() {
					s.mu.Lock()
					delete(s.roles, cwtestRole)
					delete(s.tasks, cwtestRole)
					s.mu.Unlock()
				},
				Config:           s.roleConfig(`"cw-fault-events-api" = "RW"`, export),
				ConfigPlanChecks: emptyPlan,
			},
		},
	})
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.creates != 2 {
		t.Errorf("stub saw %d creates; want 2 (the second follows the out of band delete)", s.creates)
	}
}

func TestRolesUnknownAPI(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newRoleStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{{Config: s.roleConfig(`"no-such-api" = "R"`), ExpectError: regexp.MustCompile(`unknown API id "no-such-api"`)}},
	})
	if got := s.callLog(); got != "" {
		t.Errorf("calls were %q; nothing should have been sent", got)
	}
}

func TestRolesTaskNeedsAPI(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newRoleStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      s.roleConfig(`"cw-fault-events-api" = "RW"`, "id_nso_fp_deployment_management"),
			ExpectError: regexp.MustCompile(`do not have\s+sufficient permissions`),
		}},
	})
}

func TestRolesValidation(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newRoleStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: s.roleConfig(``), ExpectError: regexp.MustCompile(`(?i)at least 1`)},
			{Config: s.roleConfig(`"cw-fault-events-api" = "X"`), ExpectError: regexp.MustCompile(`(?s)must be one of`)},
		},
	})
}

func TestRolesCreateExisting(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newRoleStub(t)
	s.roles[cwtestRole] = withBase(map[string]any{"cw-fault-events-api": roleEntry("cw-fault-events-api", "GET")})
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{{Config: s.roleConfig(`"cw-fault-events-api" = "RW"`), ExpectError: regexp.MustCompile(`same name already exist`)}},
	})
}

func TestRolesDataSource(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newRoleStub(t)
	s.roles[cwtestRole] = withBase(map[string]any{"cw-fault-events-api": roleEntry("cw-fault-events-api", "GET", "PATCH", "POST", "PUT")})
	s.tasks[cwtestRole] = map[string]bool{"id_export_audit_logs_access": true, "id_view_audit_logs_access": true}
	const ds = "data.crosswork_roles.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: s.providerBlock() + fmt.Sprintf("\ndata \"crosswork_roles\" \"test\" {\n  id = %q\n}\n", cwtestRole),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ds, "name", cwtestRole),
					resource.TestCheckResourceAttr(ds, "api_permissions.%", "4"),
					resource.TestCheckResourceAttr(ds, "api_permissions.cw-fault-events-api", "RW"),
					resource.TestCheckResourceAttr(ds, "api_permissions.aaa_cw_role_read", "R"),
					resource.TestCheckResourceAttr(ds, "task_permissions.#", "2"),
				),
			},
			{
				Config:      s.providerBlock() + "\ndata \"crosswork_roles\" \"test\" {\n  id = \"cwtest-missing\"\n}\n",
				ExpectError: regexp.MustCompile(`Failed to find role`),
			},
		},
	})
}

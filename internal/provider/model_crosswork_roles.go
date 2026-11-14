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

// Written by hand: a role needs a body that the generator cannot produce (a JSON object keyed by API id, built from
// lookups) and a second call for the task permissions.

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

const (
	rolePath = "/crosswork/aaa/v1/role"
	taskPath = "/crosswork/aaa/v1/task"
	apiPath  = "/crosswork/aaa/v1/api"
)

type Roles struct {
	Id              types.String `tfsdk:"id"`
	Instance        types.String `tfsdk:"instance"`
	Name            types.String `tfsdk:"name"`
	ApiPermissions  types.Map    `tfsdk:"api_permissions"`
	TaskPermissions types.Set    `tfsdk:"task_permissions"`
}

var accessMethods = map[rune][]string{
	'R': {"GET"},
	'W': {"PATCH", "POST", "PUT"},
	'D': {"DELETE"},
}

func methodsFor(level string) []string {
	var methods []string
	for _, letter := range level {
		methods = append(methods, accessMethods[letter]...)
	}
	sort.Strings(methods)
	return methods
}

// Crosswork adds these levels to every role by itself on three APIs (know my role, change the password, preferences). They come back on top of
// what was asked for, so they are not drift: a role that lists aaa_cwpassword with "R" reads back "RW" however many times it is written.
var baseGrants = map[string]string{
	"aaa_cw_role_read":  "R",
	"aaa_cwpassword":    "RW",
	"aaa_selected_pref": "RW",
}

func withoutBase(id, level, want string) string {
	base := baseGrants[id]
	out := ""
	for _, letter := range level {
		if strings.ContainsRune(want, letter) || !strings.ContainsRune(base, letter) {
			out += string(letter)
		}
	}
	return out
}

// PUT or PATCH counts as write; POST does not, because Crosswork stores POST next to GET for most read-only grants (a grant of "R" comes back
// as GET and POST), so POST alone would turn "R" into "RW" on every read.
func levelFor(methods []string) string {
	var read, write, del bool
	for _, m := range methods {
		switch m {
		case "GET":
			read = true
		case "PUT", "PATCH":
			write = true
		case "DELETE":
			del = true
		}
	}
	level := ""
	if read {
		level += "R"
	}
	if write {
		level += "W"
	}
	if del {
		level += "D"
	}
	return level
}

func unionLevels(a, b map[string]string) map[string]string {
	out := map[string]string{}
	for _, m := range []map[string]string{a, b} {
		for api, level := range m {
			out[api] = mergeLevel(out[api], level)
		}
	}
	return out
}

func mergeLevel(a, b string) string {
	var letters []string
	for _, l := range []string{"R", "W", "D"} {
		if strings.Contains(a, l) || strings.Contains(b, l) {
			letters = append(letters, l)
		}
	}
	return strings.Join(letters, "")
}

func equalPermissions(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func (data Roles) permissions(ctx context.Context) (map[string]string, diag.Diagnostics) {
	out := map[string]string{}
	if data.ApiPermissions.IsNull() || data.ApiPermissions.IsUnknown() {
		return out, nil
	}
	diags := data.ApiPermissions.ElementsAs(ctx, &out, false)
	return out, diags
}

func (data Roles) taskIds(ctx context.Context) ([]string, diag.Diagnostics) {
	var out []string
	if data.TaskPermissions.IsNull() || data.TaskPermissions.IsUnknown() {
		return out, nil
	}
	diags := data.TaskPermissions.ElementsAs(ctx, &out, false)
	sort.Strings(out)
	return out, diags
}

// An id missing from names is an error, so a typo is reported before anything is sent.
func roleBody(perms map[string]string, names map[string]string) (map[string]any, error) {
	rights := map[string]any{}
	for id, level := range perms {
		name, ok := names[id]
		if !ok {
			return nil, fmt.Errorf("unknown API id %q: it is not in the list returned by %s", id, apiPath)
		}
		rights[id] = map[string]any{
			"api_name":     name,
			"api_id":       id,
			"versions":     []string{},
			"allowed_urls": []any{map[string]any{"url": "/.*", "methods": methodsFor(level)}},
		}
	}
	return map[string]any{
		"_id": "", "id": "", "org_id": "1", "rate": 1000, "per": 60,
		"quota_max": -1, "quota_renewal_rate": 60, "key_expires_in": -1,
		"hmac_enabled": false, "active": true, "is_inactive": false,
		"tags":          []string{},
		"partitions":    map[string]bool{"quota": false, "rate_limit": false, "acl": false},
		"access_rights": rights,
	}, nil
}

// The create call wraps the role object under the role name; the update call takes the role object itself.
func createBody(name string, role map[string]any) (string, error) {
	b, err := json.Marshal(map[string]any{name: role})
	return string(b), err
}

func flatBody(role map[string]any) (string, error) {
	b, err := json.Marshal(role)
	return string(b), err
}

func tasksBody(ids []string) string {
	list := make([]map[string]string, 0, len(ids))
	for _, id := range ids {
		list = append(list, map[string]string{"taskId": id})
	}
	b, _ := json.Marshal(list)
	return string(b)
}

// Crosswork adds a few base APIs to every role and can enable tasks the role depends on, so once the configuration is known only the APIs and
// tasks in it are read back. A null api_permissions means an import, which reads everything; hideBuiltIn leaves out what Crosswork adds by itself.
func (data *Roles) fromRole(ctx context.Context, role gjson.Result, enabled []string, hideBuiltIn bool) diag.Diagnostics {
	var diags diag.Diagnostics
	importing := data.ApiPermissions.IsNull()
	wantPerms, d := data.permissions(ctx)
	diags.Append(d...)

	perms := map[string]string{}
	role.Get("access_rights").ForEach(func(key, value gjson.Result) bool {
		if !importing {
			if _, keep := wantPerms[key.String()]; !keep {
				return true
			}
		}
		// Only the grant of the whole API counts. Crosswork adds grants limited to a URL on a few APIs (for example
		// the delete of the user's own layouts on cw-topo-visualization) to every role, so they are not read back.
		var methods []string
		value.Get("allowed_urls").ForEach(func(_, u gjson.Result) bool {
			if u.Get("url").String() != "/.*" {
				return true
			}
			for _, m := range u.Get("methods").Array() {
				methods = append(methods, m.String())
			}
			return true
		})
		level := levelFor(methods)
		if !importing || hideBuiltIn {
			// On an import nothing was asked for, so only what is beyond the built-in grants is read.
			level = withoutBase(key.String(), level, wantPerms[key.String()])
		}
		if level != "" {
			perms[key.String()] = level
		}
		return true
	})
	var pd diag.Diagnostics
	data.ApiPermissions, pd = types.MapValueFrom(ctx, types.StringType, perms)
	diags.Append(pd...)

	if importing || !data.TaskPermissions.IsNull() {
		wantTasks, d := data.taskIds(ctx)
		diags.Append(d...)
		keep := map[string]bool{}
		for _, id := range wantTasks {
			keep[id] = true
		}
		tasks := []string{}
		for _, id := range enabled {
			if importing || keep[id] {
				tasks = append(tasks, id)
			}
		}
		var td diag.Diagnostics
		if importing && hideBuiltIn && len(tasks) == 0 {
			// A role without tasks reads as unset, as a configuration without task_permissions does.
			data.TaskPermissions = types.SetNull(types.StringType)
		} else {
			data.TaskPermissions, td = types.SetValueFrom(ctx, types.StringType, tasks)
		}
		diags.Append(td...)
	}
	return diags
}

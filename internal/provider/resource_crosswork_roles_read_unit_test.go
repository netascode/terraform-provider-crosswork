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
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

// A stored grant is a list of {url, methods}; the whole-API grant has the url "/.*".
const roleJSON = `{"access_rights": {
  "aaa_cwpassword":        {"allowed_urls": [{"url": "/.*", "methods": ["GET", "POST", "PUT"]}]},
  "aaa_selected_pref":     {"allowed_urls": [{"url": "/.*", "methods": ["GET", "PUT"]}]},
  "aaa_cw_role_read":      {"allowed_urls": [{"url": "/.*", "methods": ["GET", "POST"]}]},
  "cw-topo-visualization": {"allowed_urls": [
      {"url": "/.+(?:/group/layout/user/.*)$", "methods": ["DELETE"]},
      {"url": "/.*", "methods": ["GET", "POST"]}]},
  "cw-geo-service":        {"allowed_urls": [{"url": "/.*", "methods": ["GET", "PUT"]}]},
  "inventory_cwinventory": {"allowed_urls": [{"url": "/.*", "methods": ["GET", "POST", "DELETE"]}]}
}}`

func readRole(t *testing.T, want map[string]attr.Value, hideBuiltIn ...bool) map[string]string {
	t.Helper()
	ctx := context.Background()
	data := Roles{ApiPermissions: types.MapNull(types.StringType), TaskPermissions: types.SetNull(types.StringType)}
	if want != nil {
		data.ApiPermissions = types.MapValueMust(types.StringType, want)
	}
	if diags := data.fromRole(ctx, gjson.Parse(roleJSON), nil, len(hideBuiltIn) > 0 && hideBuiltIn[0]); diags.HasError() {
		t.Fatalf("fromRole: %v", diags)
	}
	got, diags := data.permissions(ctx)
	if diags.HasError() {
		t.Fatalf("permissions: %v", diags)
	}
	return got
}

func TestRolesReadIgnoresWhatCrossworkAddsItself(t *testing.T) {
	got := readRole(t, map[string]attr.Value{
		"aaa_cwpassword":        types.StringValue("R"), // Crosswork adds write to it
		"aaa_selected_pref":     types.StringValue("R"), // the same
		"aaa_cw_role_read":      types.StringValue("R"),
		"cw-topo-visualization": types.StringValue("R"), // Crosswork adds a delete limited to a url
	})
	for id, want := range map[string]string{
		"aaa_cwpassword": "R", "aaa_selected_pref": "R", "aaa_cw_role_read": "R", "cw-topo-visualization": "R",
	} {
		if got[id] != want {
			t.Errorf("%s: got %q, want %q", id, got[id], want)
		}
	}
}

func TestRolesReadWriteOnlyOnBaseAPIs(t *testing.T) {
	// Asking for write only on a built-in API still stores read next to it, because Crosswork adds it.
	const written = `{"access_rights": {
	  "aaa_cwpassword":    {"allowed_urls": [{"url": "/.*", "methods": ["GET", "PATCH", "POST", "PUT"]}]},
	  "aaa_selected_pref": {"allowed_urls": [{"url": "/.*", "methods": ["GET", "PATCH", "POST", "PUT"]}]},
	  "aaa_cw_role_read":  {"allowed_urls": [{"url": "/.*", "methods": ["GET", "PATCH", "POST", "PUT"]}]}
	}}`
	ctx := context.Background()
	data := Roles{
		ApiPermissions: types.MapValueMust(types.StringType, map[string]attr.Value{
			"aaa_cwpassword": types.StringValue("W"), "aaa_selected_pref": types.StringValue("W"), "aaa_cw_role_read": types.StringValue("W"),
		}),
		TaskPermissions: types.SetNull(types.StringType),
	}
	if diags := data.fromRole(ctx, gjson.Parse(written), nil, true); diags.HasError() {
		t.Fatalf("fromRole: %v", diags)
	}
	got, _ := data.permissions(ctx)
	for id, level := range got {
		if level != "W" {
			t.Errorf("%s: got %q, want W", id, level)
		}
	}
}

func TestRolesReadKeepsWhatWasAsked(t *testing.T) {
	got := readRole(t, map[string]attr.Value{
		"aaa_cwpassword":    types.StringValue("RW"),
		"aaa_selected_pref": types.StringValue("RW"),
	})
	if got["aaa_cwpassword"] != "RW" || got["aaa_selected_pref"] != "RW" {
		t.Errorf("got %v", got)
	}
}

func TestRolesReadStillReportsRealDrift(t *testing.T) {
	got := readRole(t, map[string]attr.Value{
		"cw-geo-service":        types.StringValue("R"), // write added by someone: not a base API
		"inventory_cwinventory": types.StringValue("R"), // delete on the whole API: not a base grant
		"aaa_cwpassword":        types.StringValue("R"),
	})
	if got["cw-geo-service"] != "RW" {
		t.Errorf("cw-geo-service: got %q, want RW", got["cw-geo-service"])
	}
	if got["inventory_cwinventory"] != "RD" {
		t.Errorf("inventory_cwinventory: got %q, want RD", got["inventory_cwinventory"])
	}
	if got["aaa_cwpassword"] != "R" {
		t.Errorf("aaa_cwpassword: got %q, want R", got["aaa_cwpassword"])
	}
}

func TestRolesImportOfTheResourceSkipsWhatCrossworkAddsItself(t *testing.T) {
	got := readRole(t, nil, true)
	// aaa_cwpassword is stored as GET, POST and PUT and aaa_selected_pref as GET and PUT: only the built-in grants
	if _, ok := got["aaa_cwpassword"]; ok {
		t.Errorf("aaa_cwpassword holds only the built-in grants and should not be read: %v", got)
	}
	if _, ok := got["aaa_selected_pref"]; ok {
		t.Errorf("aaa_selected_pref holds only the built-in grants and should not be read: %v", got)
	}
	// everything else is read as stored, and the delete limited to a url is not counted
	if got["cw-topo-visualization"] != "R" || got["cw-geo-service"] != "RW" || got["inventory_cwinventory"] != "RD" {
		t.Errorf("unexpected levels: %v", got)
	}
}

func TestRolesDataSourceReadsEverything(t *testing.T) {
	got := readRole(t, nil, false)
	if got["aaa_cwpassword"] != "RW" || got["aaa_selected_pref"] != "RW" {
		t.Errorf("a lookup shows the built-in grants too: %v", got)
	}
	if len(got) != 6 {
		t.Errorf("a lookup should read all 6 APIs, got %d", len(got))
	}
}

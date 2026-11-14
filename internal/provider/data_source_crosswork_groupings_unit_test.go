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

// The root call returns a plain list of ids; the group call returns the group and its direct children, or the whole
// hierarchy below it with direct=false.
type groupingStub struct {
	*ssoStub
	calls []string
	fail  bool
}

type stubGroup struct {
	name      string
	children  []string
	discovery string
	subgroups bool
}

var groupingTree = map[string]stubGroup{
	"id-root":    {"ALL-ACCESS", []string{"id-ops", "id-eng"}, "StaticSystem", true},
	"id-ops":     {"ops", nil, "Static", true},
	"id-eng":     {"eng", []string{"id-eng-sub"}, "Static", true},
	"id-eng-sub": {"eng-sub", nil, "Static", false},
}

func stubGroupJSON(id string, whole bool) string {
	g := groupingTree[id]
	children := []string{}
	for _, c := range g.children {
		if whole {
			children = append(children, stubGroupJSON(c, true))
		} else {
			children = append(children, fmt.Sprintf(`{"uuid":%q,"name":%q}`, c, groupingTree[c].name))
		}
	}
	operations := `{}`
	if g.subgroups {
		operations = `{"subGrp":true}`
	}
	return fmt.Sprintf(`{"uuid":%q,"name":%q,"classifier":"DeviceAccess","discoveryType":%q,"operations":%s,"children":[%s]}`, id, g.name, g.discovery, operations, strings.Join(children, ","))
}

func newGroupingStub(t *testing.T) *groupingStub {
	t.Helper()
	s := &groupingStub{ssoStub: newSSOStub(t)}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "/crosswork/grouping/v1/grouping/"
		if !strings.HasPrefix(r.URL.Path, prefix) {
			sso.ServeHTTP(w, r)
			return
		}
		rest := strings.TrimPrefix(r.URL.Path, prefix)
		s.calls = append(s.calls, rest+"?"+r.URL.RawQuery)
		switch {
		case s.fail && strings.HasPrefix(rest, "groups/"):
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, `{"error":"boom"}`)
		case rest == "group/root/DeviceAccess/uuid":
			fmt.Fprint(w, `["id-root"]`)
		case strings.HasPrefix(rest, "group/root/"):
			fmt.Fprint(w, `[]`)
		case strings.HasPrefix(rest, "groups/"):
			id := strings.TrimPrefix(rest, "groups/")
			_, ok := groupingTree[id]
			if !ok {
				fmt.Fprint(w, `[]`)
				return
			}
			fmt.Fprintf(w, `[%s]`, stubGroupJSON(id, r.URL.Query().Get("direct") == "false"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	t.Cleanup(Shutdown)
	return s
}

func (s *groupingStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

func TestGroupingsDataSourceWalksTheTree(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGroupingStub(t)
	const ds = "data.crosswork_groupings.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: s.providerBlock() + "\ndata \"crosswork_groupings\" \"test\" {\n  classifier = \"DeviceAccess\"\n}\n",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ds, "id", "DeviceAccess"),
					resource.TestCheckResourceAttr(ds, "items.#", "4"),
					resource.TestCheckTypeSetElemNestedAttrs(ds, "items.*", map[string]string{"id": "id-root", "name": "ALL-ACCESS", "path": "ALL-ACCESS", "discovery_type": "StaticSystem", "subgroups_allowed": "true"}),
					resource.TestCheckTypeSetElemNestedAttrs(ds, "items.*", map[string]string{"id": "id-ops", "name": "ops", "parent_id": "id-root", "path": "ALL-ACCESS/ops"}),
					resource.TestCheckTypeSetElemNestedAttrs(ds, "items.*", map[string]string{"id": "id-eng", "name": "eng", "parent_id": "id-root", "path": "ALL-ACCESS/eng"}),
					resource.TestCheckTypeSetElemNestedAttrs(ds, "items.*", map[string]string{"id": "id-eng-sub", "name": "eng-sub", "parent_id": "id-eng", "path": "ALL-ACCESS/eng/eng-sub", "discovery_type": "Static", "subgroups_allowed": "false"}),
				),
			},
			{
				Config:      s.providerBlock() + "\ndata \"crosswork_groupings\" \"test\" {\n  classifier = \"NoSuchClassifier\"\n}\n",
				ExpectError: regexp.MustCompile(`(?s)must be one of`),
			},
			{
				Config: s.providerBlock() + "\ndata \"crosswork_groupings\" \"test\" {\n  classifier = \"LocationDevices\"\n}\n",
				Check:  resource.TestCheckResourceAttr(ds, "items.#", "0"),
			},
		},
	})
	// The hierarchy is read with one call for the root, not one call for each group.
	for _, call := range s.calls {
		if strings.HasPrefix(call, "groups/") && call != "groups/id-root?direct=false" {
			t.Errorf("unexpected group call %q: the whole hierarchy comes with the call for the root", call)
		}
	}
}

func TestGroupingsDataSourceError(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGroupingStub(t)
	s.fail = true

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      s.providerBlock() + "\ndata \"crosswork_groupings\" \"test\" {\n  classifier = \"DeviceAccess\"\n}\n",
			ExpectError: regexp.MustCompile(`Failed to retrieve the groups`),
		}},
	})
}

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

// NSO proxy behaviour the tests rely on: a create is a POST on the parent container answering 201 with an empty body, an update is a PUT
// answering 204, a duplicate is 409, a missing object is 404 with an invalid-value error, a presence container is {} and an empty leaf [null],
// affinity colors come back sorted, and a delete answers 204 at once while the plan shows the removal until it answers 404.

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	odnStubBase = "/crosswork/proxy/nso/restconf/data/cisco-sr-te-cfp:sr-te/cisco-sr-te-cfp-sr-odn:odn"
	odnStubKey  = "cisco-sr-te-cfp-sr-odn:odn-template"
)

type odnStub struct {
	*ssoStub
	mu            sync.Mutex
	objects       map[string]string
	plans         map[string]int // name -> reads of the plan that still show the removal; absent means ready
	failedPlans   map[string]bool
	refuse        string // when set, a write answers 500 with this device text
	failRemoval   bool
	neverReady    bool
	removalReads  int
	calls         []string
	bodies        []string
	planReadCount map[string]int
	// The head ends of the tests, PE-1 and PE-2, are synced.
	devices map[string]string
}

func newODNStub(t *testing.T) *odnStub {
	t.Helper()
	clearEnv(t)
	Shutdown()
	old := planPollInterval
	planPollInterval = 2 * time.Millisecond
	t.Cleanup(func() { planPollInterval = old })
	s := &odnStub{ssoStub: newSSOStub(t), objects: map[string]string{}, plans: map[string]int{}, failedPlans: map[string]bool{},
		planReadCount: map[string]int{}, removalReads: 2, devices: map[string]string{"PE-1": nsoSynced, "PE-2": nsoSynced}}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == nsoNodesQuery {
			s.serveInventory(w, r)
			return
		}
		if !strings.HasPrefix(r.URL.Path, "/crosswork/proxy/nso/") {
			sso.ServeHTTP(w, r)
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		raw, _ := io.ReadAll(r.Body)
		s.calls = append(s.calls, r.Method+" "+r.URL.EscapedPath())
		if r.Method != http.MethodGet {
			s.bodies = append(s.bodies, string(raw))
		}
		s.serve(w, r, raw)
	})
	return s
}

func (s *odnStub) serveInventory(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	name := gjson.GetBytes(raw, "filter.inventory_id").String()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = append(s.calls, "INVENTORY "+name)
	rows := []any{}
	for _, d := range sortedKeys(s.devices) {
		if name == "" || d == name {
			rows = append(rows, map[string]any{"inventory_id": d, "uuid": "uuid-" + d, "reachability_state": nsoReady, "nso_state": s.devices[d], "admin_state": nsoAdminUp})
		}
	}
	write(w, 200, map[string]any{"data": rows})
}

func (s *odnStub) serve(w http.ResponseWriter, r *http.Request, raw []byte) {
	path := r.URL.EscapedPath()
	switch {
	case r.Method == http.MethodPost && path == odnStubBase:
		if s.refused(w) {
			return
		}
		entry := gjson.GetBytes(raw, odnStubKey+".0")
		name := entry.Get("name").String()
		if _, ok := s.objects[name]; ok {
			odnError(w, 409, "data-exists", "object already exists: "+name)
			return
		}
		s.objects[name] = nsoShape(entry.Raw)
		w.Header().Set("Location", "http://enso.default.svc.cluster.local:8080/restconf/data/x")
		w.WriteHeader(201)
	case strings.HasPrefix(path, odnStubBase+"/odn-template-plan="):
		s.servePlan(w, unescape(strings.TrimPrefix(path, odnStubBase+"/odn-template-plan=")))
	case strings.HasPrefix(path, odnStubBase+"/odn-template="):
		s.serveObject(w, r, unescape(strings.TrimPrefix(path, odnStubBase+"/odn-template=")), raw)
	default:
		w.WriteHeader(404)
	}
}

func (s *odnStub) refused(w http.ResponseWriter) bool {
	if s.refuse == "" {
		return false
	}
	odnError(w, 500, "operation-failed", s.refuse)
	return true
}

func (s *odnStub) serveObject(w http.ResponseWriter, r *http.Request, name string, raw []byte) {
	entry, ok := s.objects[name]
	switch r.Method {
	case http.MethodGet:
		if !ok {
			odnError(w, 404, "invalid-value", "uri keypath not found")
			return
		}
		out, _ := sjson.SetRaw("", odnStubKey, "["+entry+"]")
		w.Header().Set("Content-Type", "application/yang-data+json")
		fmt.Fprint(w, out)
	case http.MethodPut:
		if s.refused(w) {
			return
		}
		e := gjson.GetBytes(raw, odnStubKey+".0")
		if e.Get("name").String() != name {
			odnError(w, 400, "invalid-value", "the key in the body does not match the key in the path")
			return
		}
		s.objects[name] = nsoShape(e.Raw)
		w.WriteHeader(204)
	case http.MethodDelete:
		if !ok {
			odnError(w, 404, "invalid-value", "uri keypath not found")
			return
		}
		delete(s.objects, name)
		s.plans[name] = s.removalReads
		s.failedPlans[name] = s.failRemoval
		w.WriteHeader(204)
	}
}

func odnError(w http.ResponseWriter, status int, tag, msg string) {
	w.Header().Set("Content-Type", "application/yang-data+json")
	w.WriteHeader(status)
	b, _ := sjson.Set("", "ietf-restconf:errors.error.0.error-tag", tag)
	b, _ = sjson.Set(b, "ietf-restconf:errors.error.0.error-type", "application")
	b, _ = sjson.Set(b, "ietf-restconf:errors.error.0.error-message", msg)
	fmt.Fprint(w, b)
}

func unescape(s string) string {
	out, err := url.PathUnescape(s)
	if err != nil {
		return s
	}
	return out
}

// Presence containers and empty leaves keep their encoding, affinity colors are sorted, lists are kept as sent.
func nsoShape(entry string) string {
	gjson.Get(entry, "dynamic.affinity.rule").ForEach(func(k, rule gjson.Result) bool {
		var colors []string
		for _, c := range rule.Get("color").Array() {
			colors = append(colors, c.String())
		}
		sort.Strings(colors)
		entry, _ = sjson.Set(entry, fmt.Sprintf("dynamic.affinity.rule.%d.color", k.Int()), colors)
		return true
	})
	return entry
}

func odnPlanBody(name, status string, failed bool) string {
	b, _ := sjson.Set("", "cisco-sr-te-cfp-sr-odn:odn-template-plan.0.name", name)
	p := "cisco-sr-te-cfp-sr-odn:odn-template-plan.0.plan."
	b, _ = sjson.Set(b, p+"component.0.type", "tailf-ncs:self")
	b, _ = sjson.Set(b, p+"component.0.name", "self")
	b, _ = sjson.Set(b, p+"component.0.state.0.name", "tailf-ncs:init")
	b, _ = sjson.Set(b, p+"component.0.state.0.status", status)
	b, _ = sjson.Set(b, p+"component.0.state.1.name", "tailf-ncs:ready")
	b, _ = sjson.Set(b, p+"component.0.state.1.status", status)
	if failed {
		b, _ = sjson.SetRaw(b, p+"failed", "[null]")
		b, _ = sjson.Set(b, p+"error-info.message", "Failed to connect to device PE-1: connection refused")
		b, _ = sjson.Set(b, p+"status-code-detail.0.code", "TSDN-SR-301")
		b, _ = sjson.Set(b, p+"status-code-detail.0.severity", "ERROR")
		b, _ = sjson.Set(b, p+"status-code-detail.0.impacted-device", "PE-1")
		b, _ = sjson.Set(b, p+"status-code-detail.0.recommended-action", "Check device connectivity from NSO and perform recovery steps.")
		b, _ = sjson.Set(b, p+"status-code-detail.0.context.0.context-name", "Device unreachable")
		b, _ = sjson.Set(b, p+"status-code-detail.0.context.0.context-msg", "Failed to connect to device PE-1")
	}
	return b
}

func (s *odnStub) servePlan(w http.ResponseWriter, name string) {
	s.planReadCount[name]++
	w.Header().Set("Content-Type", "application/yang-data+json")
	if left, removing := s.plans[name]; removing {
		switch {
		case s.failedPlans[name]:
			fmt.Fprint(w, odnPlanBody(name, "failed", true))
		case left > 0:
			s.plans[name] = left - 1
			fmt.Fprint(w, odnPlanBody(name, "not-reached", false))
		default:
			delete(s.plans, name)
			odnError(w, 404, "invalid-value", "uri keypath not found")
		}
		return
	}
	if _, ok := s.objects[name]; !ok {
		odnError(w, 404, "invalid-value", "uri keypath not found")
		return
	}
	if s.neverReady {
		fmt.Fprint(w, odnPlanBody(name, "not-reached", false))
		return
	}
	fmt.Fprint(w, odnPlanBody(name, "reached", false))
}

func (s *odnStub) providerBlock(maxTimeout int) string {
	return fmt.Sprintf(`provider "crosswork" {
  url         = %q
  username    = "admin"
  password    = "stub-password"
  insecure    = true
  retries     = 0
  max_timeout = %d
  device_ready_timeout = 0
}
`, s.URL, maxTimeout)
}

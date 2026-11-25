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

// The list of one LCM domain as the controller answered when probed: a write replaces the whole list and needs merge, an
// empty list cannot be written, a delete that names a missing row is refused whole, a refusal is HTTP 200 with a top-level
// failure and the texts in failedData, a threshold is read back as a decimal string, and the order of the rows changes when a
// row is updated.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/tidwall/gjson"
)

const (
	lcmThresholdsReadPath   = "/crosswork/optima/api/v1/functionPacks/lcm/domain/interfaceThresholds"
	lcmThresholdsWritePath  = "/crosswork/optima/api/v1/functionPacks/lcm/domain/updateInterfaceThresholds"
	lcmThresholdsDeletePath = "/crosswork/optima/api/v1/functionPacks/lcm/domain/deleteInterfaceThresholds"
	lcmHoldMargin           = 10
)

type lcmThresholdStub struct {
	*ssoStub
	mu    sync.Mutex
	rows  map[string][]lcmThreshold // domain id to its rows, in the order the controller answers
	valid map[string]bool           // "node interface" of the IGP links of the domain
	calls []string                  // "PUT <body>", "DELETE <body>", in order
}

func newLcmThresholdStub(t *testing.T, domain string) *lcmThresholdStub {
	t.Helper()
	clearEnv(t)
	Shutdown()
	s := &lcmThresholdStub{ssoStub: newSSOStub(t), rows: map[string][]lcmThreshold{domain: nil}, valid: map[string]bool{}}
	for _, node := range []string{"cwtest-node-1", "cwtest-node-2"} {
		for _, iface := range []string{"GigabitEthernet0/0/0/14", "GigabitEthernet0/0/0/15"} {
			s.valid[node+" "+iface] = true
		}
	}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case lcmThresholdsReadPath, lcmThresholdsWritePath, lcmThresholdsDeletePath:
		default:
			sso.ServeHTTP(w, r)
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		raw, _ := io.ReadAll(r.Body)
		body := gjson.ParseBytes(raw)
		id := body.Get("id").String()
		if r.URL.Path != lcmThresholdsReadPath {
			s.calls = append(s.calls, r.Method+" "+compactJSON(string(raw)))
		}
		if _, known := s.rows[id]; !known {
			fmt.Fprint(w, `{"message":"Request timeout.Please check the logs for more details.","status":"failure"}`)
			return
		}
		switch r.URL.Path {
		case lcmThresholdsReadPath:
			s.read(w, id)
		case lcmThresholdsWritePath:
			s.write(w, id, body)
		default:
			s.remove(w, id, body)
		}
	})
	t.Cleanup(Shutdown)
	return s
}

func (s *lcmThresholdStub) read(w http.ResponseWriter, id string) {
	list := []string{}
	for _, r := range s.rows[id] {
		list = append(list, fmt.Sprintf(`{"node":%q,"interface":%q,"threshold":"%d.0","errors":[]}`, r.node, r.iface, r.threshold))
	}
	fmt.Fprintf(w, `{"data":[%s],"message":"GET successful","status":"success"}`, strings.Join(list, ","))
}

func (s *lcmThresholdStub) write(w http.ResponseWriter, id string, body gjson.Result) {
	if !body.Get("merge").Exists() {
		fmt.Fprint(w, `{"message":"Cannot invoke \"java.lang.Boolean.booleanValue()\" because \"isMerge\" is null","status":"failure"}`)
		return
	}
	add := body.Get("addInts").Array()
	if len(add) == 0 {
		fmt.Fprint(w, `{"message":"Empty CSV file imported. Please add data and try again.","status":"failure"}`)
		return
	}
	var wanted []lcmThreshold
	var failed []string
	for i, a := range add {
		row := lcmThreshold{a.Get("node").String(), a.Get("interface").String(), a.Get("threshold").Int()}
		var errs []string
		if !s.valid[row.node+" "+row.iface] {
			errs = append(errs, fmt.Sprintf("Interface %s on node %s not found", row.iface, row.node))
		}
		if row.threshold < lcmHoldMargin {
			errs = append(errs, fmt.Sprintf("Threshold %d.0 must be greater than utilization hold margin %d.0", row.threshold, lcmHoldMargin))
		}
		if len(errs) > 0 {
			list, _ := json.Marshal(errs)
			failed = append(failed, fmt.Sprintf(`{"node":%q,"interface":%q,"rowNum":"%d","errors":%s}`, row.node, row.iface, i, list))
			continue
		}
		wanted = append(wanted, row)
	}
	if len(failed) > 0 {
		fmt.Fprintf(w, `{"message":"REJECTED","status":"failure","failedData":[%s]}`, strings.Join(failed, ","))
		return
	}
	if body.Get("merge").Bool() {
		wanted = append(s.rows[id], wanted...)
	}
	s.rows[id] = wanted
	fmt.Fprint(w, `{"status":"success"}`)
}

func (s *lcmThresholdStub) remove(w http.ResponseWriter, id string, body gjson.Result) {
	data := body.Get("data").Array()
	if len(data) == 0 {
		fmt.Fprint(w, `{"message":"","status":"failure"}`)
		return
	}
	drop := map[string]bool{}
	for _, d := range data {
		key := d.Get("node").String() + " " + d.Get("interface").String()
		found := false
		for _, r := range s.rows[id] {
			found = found || r.node+" "+r.iface == key
		}
		if !found {
			fmt.Fprint(w, `{"message":"REJECTED","status":"failure"}`)
			return
		}
		drop[key] = true
	}
	var kept []lcmThreshold
	for _, r := range s.rows[id] {
		if !drop[r.node+" "+r.iface] {
			kept = append(kept, r)
		}
	}
	s.rows[id] = kept
	fmt.Fprint(w, `{"status":"success"}`)
}

// set changes the rows the way an edit in the UI would, outside Terraform.
func (s *lcmThresholdStub) set(id string, rows ...lcmThreshold) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rows[id] = rows
}

func (s *lcmThresholdStub) held(id string) []lcmThreshold {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]lcmThreshold(nil), s.rows[id]...)
}

func (s *lcmThresholdStub) callLog() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.calls...)
}

func (s *lcmThresholdStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

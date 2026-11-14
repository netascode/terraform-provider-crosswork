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

// The body must name all four jobs, each within a range; a refused value is a failed job inside HTTP 200.

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/tidwall/gjson"
)

type cadenceStub struct {
	*ssoStub
	mu     sync.Mutex
	values map[string]int64 // the four jobs, in seconds
	calls  []string         // "PUT {...}" and "QUERY", in order
	// a new value shows after this many queries
	lag    int
	reject string
	frozen bool

	pending map[string]int64
	queries int
}

func newCadenceStub(t *testing.T) *cadenceStub {
	t.Helper()
	clearEnv(t)
	Shutdown()
	s := &cadenceStub{ssoStub: newSSOStub(t), values: map[string]int64{"reach-check": 600, "show-clock": 1800, "snmp": 1200, "te-tunnel-id": 1200}}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != cadencePath && r.URL.Path != cadenceQueryPath {
			sso.ServeHTTP(w, r)
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		raw, _ := io.ReadAll(r.Body)
		if r.URL.Path == cadenceQueryPath {
			s.query(w)
			return
		}
		s.update(w, gjson.ParseBytes(raw))
	})
	t.Cleanup(Shutdown)
	return s
}

func (s *cadenceStub) query(w http.ResponseWriter) {
	s.calls = append(s.calls, "QUERY")
	if s.pending != nil {
		s.queries++
		if s.queries > s.lag {
			for k, v := range s.pending {
				s.values[k] = v
			}
			s.pending = nil
		}
	}
	fmt.Fprintf(w, `{"JobToCadence":{"reach-check":%d,"show-clock":%d,"snmp":%d,"te-tunnel-id":%d}}`,
		s.values["reach-check"], s.values["show-clock"], s.values["snmp"], s.values["te-tunnel-id"])
}

func (s *cadenceStub) failed(w http.ResponseWriter, text string) {
	fmt.Fprintf(w, `{"job_id":"job-failed","state":"JOB_FAILED","type":"Replace device job cadence","error":%q}`, text)
}

func (s *cadenceStub) update(w http.ResponseWriter, body gjson.Result) {
	s.calls = append(s.calls, "PUT "+compactJSON(body.Get("JobToCadence").Raw))
	jobs := body.Get("JobToCadence").Map()
	if s.reject != "" {
		s.failed(w, s.reject)
		return
	}
	if len(jobs) != 4 {
		s.failed(w, "All the jobs & cadence values are mandatory")
		return
	}
	next := map[string]int64{}
	for name, v := range jobs {
		if _, known := s.values[name]; !known {
			s.failed(w, "Invalid job name for job:  "+name)
			return
		}
		low := int64(600)
		if name == "reach-check" {
			low = 300
		}
		if v.Int() < low || v.Int() > 86400 {
			s.failed(w, "Invalid cadence value for job(s). For reachability job, the cadence must be between 300 seconds (5 minutes) and 86400 seconds (24 hours). For all other jobs, the cadence must be between 600 seconds (10 minutes) and 86400 seconds (24 hours).")
			return
		}
		next[name] = v.Int()
	}
	if !s.frozen {
		s.pending, s.queries = next, 0
	}
	fmt.Fprint(w, `{"job_id":"job-ok","state":"JOB_ACCEPTED","type":"Replace device job cadence"}`)
}

// gjson keeps the order the provider wrote, which is the order of cadenceKeys.
func compactJSON(raw string) string {
	out := make([]byte, 0, len(raw))
	for i := 0; i < len(raw); i++ {
		if raw[i] != ' ' && raw[i] != '\n' && raw[i] != '\t' {
			out = append(out, raw[i])
		}
	}
	return string(out)
}

func (s *cadenceStub) callLog() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.calls...)
}

func (s *cadenceStub) resetCalls() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = nil
}

func (s *cadenceStub) value(name string) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.values[name]
}

func (s *cadenceStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

func fastCadenceWait(t *testing.T) {
	old := cadenceShownInterval
	cadenceShownInterval = time.Millisecond
	t.Cleanup(func() { cadenceShownInterval = old })
}

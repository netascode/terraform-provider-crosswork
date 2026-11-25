// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License Version 2.0 (the "License"),
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

// The device alarm manager as the controller answered when probed: one flat object with a key "alarmManager/<model>" for each
// device model, a write merges the keys it sends, a key without the prefix gets it and then does not match a model, a key that is
// not a model is ignored with the answer 200, and the answer of a write holds only the keys that were sent, with the value
// the controller holds now.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"testing"
)

var deviceAlarmManagerStubModels = []string{"Cisco ASR 9010 Router", "Cisco ASR 9006 Router", "Cisco NCS 5501", "Cisco 8201 Router", "Cisco NCS 540X-ACC-SYS Router"}

type deviceAlarmManagerStub struct {
	*ssoStub
	mu    sync.Mutex
	held  map[string]bool // device model, without the prefix
	calls []string        // "GET", "POST <body>", in order
}

func newDeviceAlarmManagerStub(t *testing.T) *deviceAlarmManagerStub {
	t.Helper()
	clearEnv(t)
	Shutdown()
	s := &deviceAlarmManagerStub{ssoStub: newSSOStub(t), held: map[string]bool{}}
	for _, name := range deviceAlarmManagerStubModels {
		s.held[name] = false
	}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != deviceAlarmManagerPath {
			sso.ServeHTTP(w, r)
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		switch r.Method {
		case http.MethodGet:
			s.calls = append(s.calls, "GET")
			s.write(w, s.names())
		case http.MethodPost:
			raw, _ := io.ReadAll(r.Body)
			s.calls = append(s.calls, "POST "+string(raw))
			s.post(w, raw)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	t.Cleanup(Shutdown)
	return s
}

func (s *deviceAlarmManagerStub) names() []string {
	names := make([]string, 0, len(s.held))
	for name := range s.held {
		names = append(names, deviceAlarmManagerPrefix+name)
	}
	sort.Strings(names)
	return names
}

func (s *deviceAlarmManagerStub) write(w http.ResponseWriter, keys []string) {
	out := map[string]bool{}
	for _, key := range keys {
		out[key] = s.held[strings.TrimPrefix(key, deviceAlarmManagerPrefix)]
	}
	body, _ := json.Marshal(out)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(body))
}

func (s *deviceAlarmManagerStub) post(w http.ResponseWriter, raw []byte) {
	var in map[string]any
	_ = json.Unmarshal(raw, &in)
	var keys []string
	for key, value := range in {
		if !strings.HasPrefix(key, deviceAlarmManagerPrefix) {
			key = deviceAlarmManagerPrefix + key
		}
		name := strings.TrimPrefix(key, deviceAlarmManagerPrefix)
		if _, known := s.held[name]; known {
			if enabled, isBool := value.(bool); isBool {
				s.held[name] = enabled
			}
		}
		keys = append(keys, key)
	}
	s.write(w, keys)
}

// set changes a model the way an edit in the UI would, outside Terraform.
func (s *deviceAlarmManagerStub) set(name string, enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.held[name] = enabled
}

func (s *deviceAlarmManagerStub) state() map[string]bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := map[string]bool{}
	for name, enabled := range s.held {
		out[name] = enabled
	}
	return out
}

func (s *deviceAlarmManagerStub) callLog() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.calls...)
}

func (s *deviceAlarmManagerStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

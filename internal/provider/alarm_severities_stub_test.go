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

// The alarm catalog as the controller answers it when probed: the read without the query leaves out the hidden types, a name is the
// key of a write (a type whose name differs from its event type is refused by the event type), a call is applied whole or not at all
// and an unknown name is a plain-text 400 that names nothing.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
)

type alarmStubRow struct {
	name      string
	eventType string
	severity  string // as the controller reads it: capitalised
	revert    string // minutes, empty when none
	hidden    bool
}

type alarmStub struct {
	*ssoStub
	mu    sync.Mutex
	rows  []*alarmStubRow
	calls []string // "GET", "GET full", "POST scc minor [a b]", in order
}

func newAlarmStub(t *testing.T) *alarmStub {
	t.Helper()
	clearEnv(t)
	Shutdown()
	s := &alarmStub{ssoStub: newSSOStub(t), rows: []*alarmStubRow{
		{name: "AUTO-AMPLI-MISMATCH", eventType: "AUTO-AMPLI-MISMATCH", severity: "Information"},
		{name: "L2-ELMI-4-AIPC", eventType: "L2-ELMI-4-AIPC", severity: "Minor", revert: "1440"},
		{name: "DWDM-5-G709ALARM description", eventType: "DWDM-5-G709ALARM", severity: "Information"},
		{name: "Link down(admin down)", eventType: "LINK_DOWN_ADMIN", severity: "Major"},
		{name: "ADHOC_ROGUE_AUTO_CONTAINED", eventType: "ADHOC_ROGUE_AUTO_CONTAINED", severity: "Major", hidden: true},
	}}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != alarmSeveritiesPath {
			sso.ServeHTTP(w, r)
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if r.Method == http.MethodGet {
			s.read(w, r.URL.Query().Get("completeEventList") == "true")
			return
		}
		raw, _ := io.ReadAll(r.Body)
		s.write(w, raw)
	})
	t.Cleanup(Shutdown)
	return s
}

func (s *alarmStub) read(w http.ResponseWriter, full bool) {
	if full {
		s.calls = append(s.calls, "GET full")
	} else {
		s.calls = append(s.calls, "GET")
	}
	items := []map[string]string{}
	for _, row := range s.rows {
		if row.hidden && !full {
			continue
		}
		item := map[string]string{"name": row.name, "eventTypeName": row.eventType, "severity": row.severity, "defaultCategory": "Test"}
		if row.revert != "" {
			item["revert"] = row.revert
		}
		items = append(items, item)
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
}

func (s *alarmStub) refuse(w http.ResponseWriter, text string) {
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, text)
}

func (s *alarmStub) write(w http.ResponseWriter, raw []byte) {
	var body struct {
		SourceType  string   `json:"sourceType"`
		SourceValue string   `json:"sourceValue"`
		EventTypes  []string `json:"eventTypes"`
	}
	_ = json.Unmarshal(raw, &body)
	s.calls = append(s.calls, fmt.Sprintf("POST %s %s [%s]", body.SourceType, body.SourceValue, strings.Join(body.EventTypes, "|")))
	if body.SourceType != "scc" && body.SourceType != "aac" {
		s.refuse(w, "Invalid sourceType")
		return
	}
	if len(body.EventTypes) == 0 {
		s.refuse(w, "eventTypes must not be empty")
		return
	}
	var targets []*alarmStubRow
	for _, name := range body.EventTypes {
		found := false
		for _, row := range s.rows {
			if row.name == name {
				targets, found = append(targets, row), true
			}
		}
		if !found {
			s.refuse(w, "Invalid eventType")
			return
		}
	}
	switch body.SourceType {
	case "scc":
		severity := map[string]string{"critical": "Critical", "major": "Major", "minor": "Minor", "warning": "Warning", "information": "Information"}[strings.ToLower(body.SourceValue)]
		if severity == "" {
			s.refuse(w, "Invalid severity")
			return
		}
		for _, row := range targets {
			if !row.hidden { // a hidden type answers success and keeps its severity
				row.severity = severity
			}
		}
	case "aac":
		minutes, err := strconv.Atoi(body.SourceValue)
		if err != nil || minutes < 5 || minutes > 599940 || (minutes <= 55 && minutes%5 != 0) || (minutes >= 60 && minutes%60 != 0) || (minutes > 55 && minutes < 60) {
			s.refuse(w, "Please enter valid integer value in the range of 5 to 599940. If value is less than or equal to 55, then it should be in multiples of 5. If the value is greater than or equal to 60, then it should be in multiples of 60.")
			return
		}
		for _, row := range targets {
			row.revert = strconv.Itoa(minutes)
		}
	}
	fmt.Fprint(w, `{"status":"success","message":"Alarm autoclear update:success"}`)
}

func (s *alarmStub) row(name string) alarmStubRow {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, row := range s.rows {
		if row.name == name {
			return *row
		}
	}
	return alarmStubRow{}
}

func (s *alarmStub) set(name, severity, revert string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, row := range s.rows {
		if row.name == name {
			row.severity, row.revert = severity, revert
		}
	}
}

func (s *alarmStub) callLog() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.calls...)
}

// posts returns only the writes, in order.
func (s *alarmStub) posts() []string {
	var out []string
	for _, c := range s.callLog() {
		if strings.HasPrefix(c, "POST ") {
			out = append(out, c)
		}
	}
	return out
}

func (s *alarmStub) resetCalls() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = nil
}

func (s *alarmStub) names() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []string
	for _, row := range s.rows {
		out = append(out, row.name)
	}
	sort.Strings(out)
	return out
}

func (s *alarmStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

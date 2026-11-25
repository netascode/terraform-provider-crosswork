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

// The controller holds one configuration for each LCM domain it created; a write replaces the whole configuration. A refused
// write is HTTP 200 with a failure status, and an id that is not a domain fails with a message about a missing NATS reply.

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	lcmPath       = "/crosswork/optima/api/v1/functionPacks/lcm"
	lcmConfigPath = "/crosswork/optima/api/v1/functionPacks/lcm/domain/config"
)

var lcmBitNames = map[string]string{"28": "IMPAIRED", "29": "ENCRYPTED", "30": "PLANE-A", "31": "PLANE-B"}

type lcmStub struct {
	*ssoStub
	mu      sync.Mutex
	configs map[string]string // domain id to the configuration as it was written
	puts    []string          // "<id> <config>", in order
	// a value the controller holds that the provider cannot write, set by a test
	fraction map[string]float64
}

func newLcmStub(t *testing.T, domains ...string) *lcmStub {
	t.Helper()
	clearEnv(t)
	Shutdown()
	s := &lcmStub{ssoStub: newSSOStub(t), configs: map[string]string{}, fraction: map[string]float64{}}
	for _, d := range domains {
		s.configs[d] = `{}`
	}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != lcmPath && r.URL.Path != lcmConfigPath {
			sso.ServeHTTP(w, r)
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		raw, _ := io.ReadAll(r.Body)
		body := gjson.ParseBytes(raw)
		id := body.Get("id").String()
		if _, known := s.configs[id]; !known {
			msg := `Cannot invoke \"io.nats.client.Message.getData()\" because \"msg\" is null`
			if r.URL.Path == lcmConfigPath {
				// the read has no list of errors: the library does not see this answer as a failure
				fmt.Fprintf(w, `{"data":{"error":"%s","message":"%s","status":"failure"}}`, msg, msg)
			} else {
				s.puts = append(s.puts, id+" "+compactJSON(body.Get("config").Raw))
				fmt.Fprintf(w, `{"data":{"errors":[{"field":null,"errorDescription":"%s","errorCode":null}],"status":"failure"}}`, msg)
			}
			return
		}
		if r.URL.Path == lcmConfigPath {
			fmt.Fprint(w, s.read(id))
			return
		}
		s.write(w, id, body.Get("config"))
	})
	t.Cleanup(Shutdown)
	return s
}

func (s *lcmStub) write(w http.ResponseWriter, id string, cfg gjson.Result) {
	s.puts = append(s.puts, id+" "+compactJSON(cfg.Raw))
	if cfg.Get("utilThreshold").Float() < cfg.Get("utilHoldMargin").Float() {
		fmt.Fprint(w, `{"data":{"status":"failure","errors":[{"field":"utilThreshold","errorCode":"THRESHMARGINERROR","errorDescription":"The utilization threshold must be greater than or equal to the hold margin"}]}}`)
		return
	}
	s.configs[id] = cfg.Raw
	delete(s.fraction, id)
	fmt.Fprint(w, `{"data":{"status":"success"}}`)
}

// read answers the way the controller does: whole numbers as decimals, affinities as objects with the name and the bit,
// and the hops as one object for each tag.
func (s *lcmStub) read(id string) string {
	cfg := gjson.Parse(s.configs[id])
	out := cfg.Raw
	for _, k := range []string{"utilThreshold", "utilHoldMargin", "overprovisionFactor", "unevenEcmpTrafficThreshold"} {
		if cfg.Get(k).Exists() {
			out, _ = sjson.SetRaw(out, k, fmt.Sprintf("%d.0", cfg.Get(k).Int()))
		}
	}
	if f, ok := s.fraction[id]; ok {
		out, _ = sjson.SetRaw(out, "overprovisionFactor", strconv.FormatFloat(f, 'f', -1, 64))
	}
	for _, k := range []string{"includeAny", "excludeAny", "includeAll"} {
		list := []string{}
		for _, bit := range cfg.Get("affinityConstraint." + k).Array() {
			list = append(list, fmt.Sprintf(`{"id":1,"name":%q,"bitPosition":%s}`, lcmBitNames[bit.String()], bit.String()))
		}
		if cfg.Get("affinityConstraint").Exists() {
			out, _ = sjson.SetRaw(out, "affinityConstraint."+k, "["+strings.Join(list, ",")+"]")
		}
	}
	hops := []string{}
	for _, group := range cfg.Get("nodeTagMaxSidDepth").Array() {
		keys := []string{}
		for tag := range group.Map() {
			keys = append(keys, tag)
		}
		sort.Strings(keys)
		for _, tag := range keys {
			hops = append(hops, fmt.Sprintf(`{%q:%d}`, tag, group.Get(tag).Int()))
		}
	}
	out, _ = sjson.SetRaw(out, "nodeTagMaxSidDepth", "["+strings.Join(hops, ",")+"]")
	return `{"data":{"config":` + out + `,"tooltips":{}}}`
}

func (s *lcmStub) config(id string) gjson.Result {
	s.mu.Lock()
	defer s.mu.Unlock()
	return gjson.Parse(s.configs[id])
}

func (s *lcmStub) putLog() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.puts...)
}

func (s *lcmStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

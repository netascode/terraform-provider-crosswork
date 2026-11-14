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
	"reflect"
	"strings"
	"testing"
)

type assertions struct{}

var assert assertions

func fail(t *testing.T, format string, a ...any) {
	t.Helper()
	t.Errorf(format, a...)
}

func (assertions) Equal(t *testing.T, want, got any, msg ...any) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		fail(t, "want %v, got %v %s", want, got, fmt.Sprint(msg...))
	}
}

func (assertions) NotEqual(t *testing.T, a, b any, msg ...any) {
	t.Helper()
	if reflect.DeepEqual(a, b) {
		fail(t, "both are %v %s", a, fmt.Sprint(msg...))
	}
}

func (assertions) True(t *testing.T, v bool, msg ...any) {
	t.Helper()
	if !v {
		fail(t, "want true %s", fmt.Sprint(msg...))
	}
}

func (assertions) False(t *testing.T, v bool, msg ...any) {
	t.Helper()
	if v {
		fail(t, "want false %s", fmt.Sprint(msg...))
	}
}

func (assertions) Contains(t *testing.T, s, sub string) {
	t.Helper()
	if !strings.Contains(s, sub) {
		fail(t, "%q does not contain %q", s, sub)
	}
}

func (assertions) NotContains(t *testing.T, s, sub string, msg ...any) {
	t.Helper()
	if strings.Contains(s, sub) {
		fail(t, "%q contains %q %s", s, sub, fmt.Sprint(msg...))
	}
}

func (assertions) Len(t *testing.T, obj any, n int, msg ...any) {
	t.Helper()
	if l := reflect.ValueOf(obj).Len(); l != n {
		fail(t, "length is %d, want %d %s", l, n, fmt.Sprint(msg...))
	}
}

func TestGcCanonical(t *testing.T) {
	assert.Equal(t, `{"a":2,"b":1}`, gcCanonical("JSON", `{ "b": 1, "a": 2 }`), "keys sorted, no spaces")
	assert.Equal(t, `{ "b": 1 }`, gcCanonical("STRING", `{ "b": 1 }`), "a string is left alone")
	assert.Equal(t, "not json", gcCanonical("JSON", "not json"))
	assert.True(t, gcSameValue("JSON", `{"a":1,"b":2}`, `{"b": 2, "a": 1}`))
	assert.False(t, gcSameValue("JSON", `{"a":1}`, `{"a":2}`))
	assert.True(t, gcValidJSON(`[1,2]`))
	assert.False(t, gcValidJSON("{"))
	assert.Equal(t, map[string]any{"a": float64(1)}, gcBodyValue("JSON", `{"a":1}`))
	assert.Equal(t, "text", gcBodyValue("STRING", "text"))
}

func TestGcMakeGroups(t *testing.T) {
	plain := func(v string) []gcVar { return []gcVar{{Name: "ntp", Type: "STRING", Value: v}} }
	secret := func(v string, version int64) []gcVar {
		return []gcVar{{Name: "key", Type: "STRING", Value: v, Secret: true, Version: version}}
	}

	t.Run("equal values share a job, the first group keeps the template name", func(t *testing.T) {
		groups := gcMakeGroups("tpl", []string{"d3", "d1", "d2"}, map[string][]gcVar{"d1": nil, "d2": plain("x"), "d3": nil})
		assert.Len(t, groups, 2)
		assert.Equal(t, []string{"d1", "d3"}, groups[0].Devices)
		assert.Equal(t, "tpl", groups[0].Name)
		assert.Equal(t, []string{"d2"}, groups[1].Devices)
		assert.True(t, strings.HasPrefix(groups[1].Name, "tpl-") && len(groups[1].Name) == len("tpl-")+gcHashLen, groups[1].Name)
	})

	t.Run("two devices with different secrets of the same version never share a job", func(t *testing.T) {
		groups := gcMakeGroups("tpl", []string{"d1", "d2"}, map[string][]gcVar{"d1": secret("one", 1), "d2": secret("two", 1)})
		assert.Len(t, groups, 2)
	})

	t.Run("the name does not depend on a secret", func(t *testing.T) {
		a := gcMakeGroups("tpl", []string{"d0", "d1"}, map[string][]gcVar{"d0": nil, "d1": secret("one", 1)})
		b := gcMakeGroups("tpl", []string{"d0", "d1"}, map[string][]gcVar{"d0": nil, "d1": secret("another", 1)})
		assert.Equal(t, a[1].Name, b[1].Name, "the hash covers the version, so a guessed secret cannot be tested against the name")
		c := gcMakeGroups("tpl", []string{"d0", "d1"}, map[string][]gcVar{"d0": nil, "d1": secret("one", 2)})
		assert.NotEqual(t, a[1].Name, c[1].Name, "another version is another name")
	})

	t.Run("a job name stays within the 64 characters the controller accepts", func(t *testing.T) {
		long := strings.Repeat("n", 64)
		groups := gcMakeGroups(long, []string{"d1", "d2"}, map[string][]gcVar{"d1": nil, "d2": plain("x")})
		assert.Equal(t, long, groups[0].Name)
		assert.Len(t, groups[1].Name, gcJobNameMax)
		assert.True(t, strings.HasPrefix(groups[1].Name, strings.Repeat("n", 55)+"-"))
	})

	t.Run("the order of the variables does not matter", func(t *testing.T) {
		ab := []gcVar{{Name: "a", Type: "STRING", Value: "1"}, {Name: "b", Type: "STRING", Value: "2"}}
		ba := []gcVar{ab[1], ab[0]}
		assert.Len(t, gcMakeGroups("tpl", []string{"d1", "d2"}, map[string][]gcVar{"d1": ab, "d2": ba}), 1)
	})
}

func TestGcVarKeys(t *testing.T) {
	json1 := gcVar{Name: "o", Type: "JSON", Value: `{"b":1,"a":2}`}
	json2 := gcVar{Name: "o", Type: "JSON", Value: `{ "a": 2, "b": 1 }`}
	assert.Equal(t, json1.nameKey(), json2.nameKey(), "a JSON value is keyed in canonical form")
	s1 := gcVar{Name: "k", Type: "STRING", Value: "one", Secret: true, Version: 1}
	s2 := gcVar{Name: "k", Type: "STRING", Value: "two", Secret: true, Version: 1}
	assert.Equal(t, s1.nameKey(), s2.nameKey(), "a secret is keyed by its version")
	assert.NotEqual(t, s1.groupKey(), s2.groupKey(), "and grouped by its value")
}

func TestGcJobBody(t *testing.T) {
	nodes := map[string]nsoNode{"d1": {Name: "d1", UUID: "u1", HostName: "real-host-name"}}
	b := int64(2)
	body := gcJobBody(gcJobParams{ConfigID: 7, TemplateName: "tpl", Vendor: "Cisco Systems", Series: "S", Mop: "m", MopID: 4, Batch: &b},
		gcGroup{Name: "tpl", Devices: []string{"d1"}, Vars: []gcVar{{Name: "v", Type: "STRING", Value: "x"}}}, nodes)
	for _, want := range []string{`"app":"Golden Config"`, `"mopId":4`, `"tcId":7`, `"uuid":"u1"`, `"hostName":"d1"`, `"batchSize":2`, `"productSeries":"S"`} {
		assert.Contains(t, body, want)
	}
	assert.NotContains(t, body, "real-host-name", "hostName is sent as the inventory id")
	assert.NotContains(t, body, "errorBudget", "an unset budget is omitted")
	assert.NotContains(t, body, "jobActions", "the MOP supplies its own actions")
	assert.NotContains(t, body, `"source"`, "a job value is sent without a source so it beats the template")
}

func TestGcNodeReady(t *testing.T) {
	up := nsoNode{Name: "d", Reachability: nsoReady, NSO: nsoSynced, Admin: nsoAdminUp}
	down := nsoNode{Name: "d", Reachability: "CONN_STATE_UNKNOWN", NSO: nsoSynced, Admin: "ROBOT_ADMIN_STATE_DOWN"}
	unreachable := nsoNode{Name: "d", Reachability: "CONN_STATE_UNREACHABLE", NSO: nsoSynced, Admin: nsoAdminUp}
	notSynced := nsoNode{Name: "d", Reachability: nsoReady, NSO: "FETCH_SSH_KEYS_FAILED", Admin: nsoAdminUp}
	for _, tc := range []struct {
		name         string
		node         nsoNode
		synced, full bool
	}{
		{"reachable, synced and admin up", up, true, true},
		{"admin down and not reachable, but synced", down, true, false},
		{"unreachable but synced", unreachable, true, false},
		{"reachable but not synced", notSynced, false, false},
	} {
		assert.Equal(t, tc.synced, tc.node.ready(nsoModeSynced), "synced: "+tc.name)
		assert.Equal(t, tc.full, tc.node.ready(nsoModeFull), "full: "+tc.name)
	}
	assert.Equal(t, down.ready(nsoModeSynced), down.ready(""), "an empty mode is the synced mode")
}

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

import (
	"testing"

	"github.com/tidwall/gjson"
)

func TestLcmThresholdsBodies(t *testing.T) {
	rows := []lcmThreshold{{"cwtest-node-1", "GigabitEthernet0/0/0/14", 75}}
	want := `{"id":"32","merge":false,"addInts":[{"errors":[],"interface":"GigabitEthernet0/0/0/14","node":"cwtest-node-1","threshold":75}]}`
	if got := lcmThresholdsPutBody("32", rows); got != want {
		t.Fatalf("put body:\n got %s\nwant %s", got, want)
	}
	want = `{"id":"32","data":[{"errors":[],"interface":"GigabitEthernet0/0/0/14","node":"cwtest-node-1","threshold":75}]}`
	if got := lcmThresholdsDeleteBody("32", rows); got != want {
		t.Fatalf("delete body:\n got %s\nwant %s", got, want)
	}
}

func TestLcmThresholdsFromBody(t *testing.T) {
	res := gjson.Parse(`{"data":[{"node":"b","interface":"Gi2","threshold":"80.0"},{"node":"a","interface":"Gi9","threshold":"75.5"},{"node":"a","interface":"Gi1","threshold":"10.0"}],"status":"success"}`)
	got := lcmThresholdsFromBody(res)
	want := []lcmThreshold{{"a", "Gi1", 10}, {"a", "Gi9", 75}, {"b", "Gi2", 80}}
	if len(got) != len(want) {
		t.Fatalf("got %v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("row %d: got %v, want %v", i, got[i], want[i])
		}
	}
}

func TestLcmThresholdItemsEmpty(t *testing.T) {
	if items := lcmThresholdItems(nil, nil); items != nil {
		t.Fatalf("no prior list must stay null, got %v", items)
	}
	prior := []LocalCongestionMitigationInterfaceThresholdsInterfaceThresholds{}
	if items := lcmThresholdItems(nil, prior); items == nil || len(items) != 0 {
		t.Fatalf("an empty prior list must stay empty, got %v", items)
	}
}

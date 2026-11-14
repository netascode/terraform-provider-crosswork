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
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/tidwall/gjson"
)

type lookupStub struct {
	*ssoStub
	nodes             int // number of nodes of the cluster, named dev-<n> with the uuid uuid-<n>
	single, listCalls atomic.Int32
}

func newLookupStub(t *testing.T, nodes int) *lookupStub {
	t.Helper()
	clearEnv(t)
	Shutdown()
	s := &lookupStub{ssoStub: newSSOStub(t), nodes: nodes}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case groupingPath + "/group/root/LocationDevices/uuid":
			fmt.Fprint(w, `["g-root"]`)
			return
		case groupingPath + "/groups/g-root":
			// the whole hierarchy, as direct=false returns it: the root and All Locations hold no devices, the groups below
			// do, and one group reports no operations at all
			fmt.Fprint(w, `[{"uuid":"g-root","operations":{"cpf":true},"children":[
				{"uuid":"g-all","operations":{"cpf":true,"subGrp":true},"children":[
					{"uuid":"g-unassigned","operations":{"showMem":true,"cpf":true,"addMem":true,"mv":1}},
					{"uuid":"g-custom","operations":{"showMem":true,"addMem":true,"subGrp":true},"children":[
						{"uuid":"g-deep","operations":{"showMem":true,"addMem":true}}]},
					{"uuid":"g-unreported"}]}]}]`)
			return
		case nodesQueryPath:
		default:
			sso.ServeHTTP(w, r)
			return
		}
		raw, _ := io.ReadAll(r.Body)
		body := gjson.ParseBytes(raw)
		node := func(n int) string {
			return fmt.Sprintf(`{"inventory_id":"dev-%d","uuid":"uuid-%d"}`, n, n)
		}
		var data []string
		switch {
		case body.Get("filter.inventory_id").Exists():
			s.single.Add(1)
			for n := 0; n < s.nodes; n++ {
				// the controller matches the inventory id without regard to case
				if strings.EqualFold(body.Get("filter.inventory_id").String(), fmt.Sprintf("dev-%d", n)) {
					data = append(data, node(n))
				}
			}
		case body.Get("filter.uuid").Exists():
			s.single.Add(1)
			for n := 0; n < s.nodes; n++ {
				if body.Get("filter.uuid").String() == fmt.Sprintf("uuid-%d", n) {
					data = append(data, node(n))
				}
			}
		default:
			s.listCalls.Add(1)
			size, page := int(body.Get("filterData.PageSize").Int()), int(body.Get("filterData.PageNum").Int())
			for n := page * size; n < s.nodes && n < (page+1)*size; n++ {
				data = append(data, node(n))
			}
		}
		fmt.Fprintf(w, `{"data":[%s],"total_count":%d,"result_count":%d}`, strings.Join(data, ","), s.nodes, len(data))
	})
	t.Cleanup(Shutdown)
	return s
}

func names(count int) []string {
	out := make([]string, count)
	for i := range out {
		out[i] = fmt.Sprintf("dev-%d", i)
	}
	return out
}

func uuids(count int) []string {
	out := make([]string, count)
	for i := range out {
		out[i] = fmt.Sprintf("uuid-%d", i)
	}
	return out
}

func TestDeviceLookupFewDevicesAskOneByOne(t *testing.T) {
	s := newLookupStub(t, 30)
	c, _ := registeredClient(s.URL, "admin", "stub-password", true, 0, "test")
	got, err := deviceUUIDs(c, names(bulkLookupMin-1))
	if err != nil || len(got) != bulkLookupMin-1 {
		t.Fatalf("got %v, %v", got, err)
	}
	if s.listCalls.Load() != 0 || int(s.single.Load()) != bulkLookupMin-1 {
		t.Errorf("a few devices are asked for one by one: list calls %d, single calls %d", s.listCalls.Load(), s.single.Load())
	}
}

func TestDeviceLookupManyDevicesUseTheList(t *testing.T) {
	s := newLookupStub(t, 30)
	c, _ := registeredClient(s.URL, "admin", "stub-password", true, 0, "test")

	got, err := deviceUUIDs(c, names(20))
	if err != nil || len(got) != 20 || got["dev-7"] != "uuid-7" {
		t.Fatalf("got %v, %v", got, err)
	}
	back, err := inventoryIDs(c, uuids(20))
	if err != nil || len(back) != 20 || back["uuid-7"] != "dev-7" {
		t.Fatalf("got %v, %v", back, err)
	}
	if s.single.Load() != 0 || s.listCalls.Load() != 2 {
		t.Errorf("twenty devices each way take one list call each way: list calls %d, single calls %d", s.listCalls.Load(), s.single.Load())
	}
}

func TestDeviceLookupFallsBackForWhatTheListLacks(t *testing.T) {
	s := newLookupStub(t, 30)
	c, _ := registeredClient(s.URL, "admin", "stub-password", true, 0, "test")

	// a name in another case is found by the controller's own filter, not by the list
	ask := append(names(9), "DEV-9")
	got, err := deviceUUIDs(c, ask)
	if err != nil || got["DEV-9"] != "uuid-9" {
		t.Fatalf("got %v, %v", got, err)
	}
	if s.listCalls.Load() != 1 || s.single.Load() != 1 {
		t.Errorf("one list call and one single call: list calls %d, single calls %d", s.listCalls.Load(), s.single.Load())
	}

	if _, err := deviceUUIDs(c, append(names(9), "no-such-device")); err == nil || !strings.Contains(err.Error(), "no-such-device does not exist") {
		t.Errorf("a missing device is an error, got %v", err)
	}

	// a uuid with no node is left out, and a complete list needs no call of its own for it
	before := s.single.Load()
	back, err := inventoryIDs(c, append(uuids(9), "uuid-gone"))
	if err != nil || len(back) != 9 || s.single.Load() != before {
		t.Errorf("got %v, %v, single calls %d", back, err, s.single.Load()-before)
	}
}

func TestDeviceLookupReadsOnlyAsManyPagesAsItSaves(t *testing.T) {
	// 2500 nodes are three pages. Ten devices may read one page, because the pages must not cost more calls than
	// the lookups they replace; the devices not on it are asked for one by one.
	s := newLookupStub(t, 2500)
	c, _ := registeredClient(s.URL, "admin", "stub-password", true, 0, "test")
	ask := []string{"dev-1", "dev-2", "dev-3", "dev-1500", "dev-1501", "dev-2400", "dev-2401", "dev-2402", "dev-2403", "dev-2404"}
	got, err := deviceUUIDs(c, ask)
	if err != nil || len(got) != len(ask) {
		t.Fatalf("got %v, %v", got, err)
	}
	if s.listCalls.Load() != 1 || s.single.Load() != 7 {
		t.Errorf("one page and seven single calls: list calls %d, single calls %d", s.listCalls.Load(), s.single.Load())
	}
	if total := int(s.listCalls.Load() + s.single.Load()); total > len(ask) {
		t.Errorf("%d calls for %d devices is more than one call each", total, len(ask))
	}

	s.listCalls.Store(0)
	s.single.Store(0)
	many := names(2500)[:24]
	many = append(many, "dev-2499")
	got, err = deviceUUIDs(c, many)
	if err != nil || len(got) != 25 {
		t.Fatalf("got %d, %v", len(got), err)
	}
	if s.listCalls.Load() != 3 || s.single.Load() != 0 {
		t.Errorf("three pages and no single call: list calls %d, single calls %d", s.listCalls.Load(), s.single.Load())
	}
}

func TestLocationGroupsSkipsGroupsThatHoldNoDevices(t *testing.T) {
	s := newLookupStub(t, 1)
	c, _ := registeredClient(s.URL, "admin", "stub-password", true, 0, "test")
	got, err := locationGroups(c)
	if err != nil {
		t.Fatal(err)
	}
	want := "g-unassigned,g-custom,g-deep,g-unreported"
	if strings.Join(got, ",") != want {
		t.Errorf("got %v, want the groups that show members, at any depth, and the one that reports no operations: %s", got, want)
	}
}

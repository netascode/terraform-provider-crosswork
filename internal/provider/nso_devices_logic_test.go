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
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
)

// Can be told to ignore the filter.
type invStub struct {
	*ssoStub
	mu           sync.Mutex
	nodes        map[string]*stubNode
	full         atomic.Int32
	byName       atomic.Int32
	ignoreFilter bool
}

func newInvStub(t *testing.T, names ...string) (*invStub, *cw.Client) {
	t.Helper()
	s := &invStub{ssoStub: newSSOStub(t), nodes: map[string]*stubNode{}}
	for _, n := range names {
		s.nodes[n] = &stubNode{Name: n, UUID: "uuid-" + n, Host: "host-" + n, Reach: nsoReady, NSO: nsoSynced, Admin: nsoAdminUp}
	}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != nsoNodesQuery {
			sso.ServeHTTP(w, r)
			return
		}
		raw, _ := io.ReadAll(r.Body)
		var body struct {
			Filter map[string]string `json:"filter"`
		}
		_ = json.Unmarshal(raw, &body)
		name := body.Filter["inventory_id"]
		if name == "" || s.ignoreFilter {
			s.full.Add(1)
			name = ""
		} else {
			s.byName.Add(1)
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		rows := []any{}
		for _, n := range sortedKeys(s.nodes) {
			x := s.nodes[n]
			if name == "" || x.Name == name {
				rows = append(rows, map[string]any{"inventory_id": x.Name, "uuid": x.UUID, "host_name": x.Host, "reachability_state": x.Reach, "nso_state": x.NSO, "admin_state": x.Admin})
			}
		}
		write(w, 200, map[string]any{"data": rows})
	})
	c, err := registeredClient(s.URL, "admin", "stub-password", true, 0, "test")
	if err != nil {
		t.Fatal(err)
	}
	oldInterval := nsoReadinessInterval
	nsoReadinessInterval = 5 * time.Millisecond
	t.Cleanup(func() {
		nsoReadinessInterval = oldInterval
		Shutdown()
	})
	return s, c
}

func (s *invStub) set(name, nso string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes[name].NSO = nso
}

func devs(names ...string) nsoDeviceSet {
	set := nsoDeviceSet{}
	for _, n := range names {
		set[n] = struct{}{}
	}
	return set
}

func cfgOf(mode string, timeout time.Duration) nsoDeviceConfig {
	return nsoDeviceConfig{Mode: mode, Timeout: timeout, Store: newNsoSnapshots()}
}

func noErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func TestNSODevicesReadyFromSnapshot(t *testing.T) {
	s, c := newInvStub(t, "d1", "d2")
	cfg := cfgOf(nsoModeSynced, 0)
	nodes, notReady, err := waitNSODevices(context.Background(), c, "", cfg, devs("d1", "d2"))
	noErr(t, err)
	assert.Len(t, notReady, 0)
	assert.Equal(t, "uuid-d1", nodes["d1"].UUID)
	assert.Equal(t, "host-d2", nodes["d2"].HostName)
	_, notReady, err = waitNSODevices(context.Background(), c, "", cfg, devs("d1"))
	noErr(t, err)
	assert.Len(t, notReady, 0)
	assert.Equal(t, int32(1), s.full.Load(), "one read of the inventory")
	assert.Equal(t, int32(0), s.byName.Load(), "no device was read by name")

	now := time.Now()
	nsoNow = func() time.Time { return now.Add(nsoSnapshotTTL + time.Second) }
	t.Cleanup(func() { nsoNow = time.Now })
	_, _, err = waitNSODevices(context.Background(), c, "", cfg, devs("d1"))
	noErr(t, err)
	assert.Equal(t, int32(2), s.full.Load())
}

func TestNSODevicesNotReadyIsReadByName(t *testing.T) {
	s, c := newInvStub(t, "d1", "d2")
	s.set("d2", "NOT_SYNCED")
	go func() {
		time.Sleep(40 * time.Millisecond)
		s.set("d2", nsoSynced)
	}()
	nodes, notReady, err := waitNSODevices(context.Background(), c, "", cfgOf(nsoModeSynced, 5*time.Second), devs("d1", "d2"))
	noErr(t, err)
	assert.Len(t, notReady, 0)
	assert.Equal(t, nsoSynced, nodes["d2"].NSO)
	assert.Equal(t, int32(1), s.full.Load())
	assert.True(t, s.byName.Load() > 1, "d2 was read by name at each interval")
}

func TestNSODevicesSnapshotIsNotTrustedForNotReady(t *testing.T) {
	s, c := newInvStub(t, "d1")
	s.set("d1", "NOT_SYNCED")
	cfg := cfgOf(nsoModeSynced, 0)
	_, notReady, err := waitNSODevices(context.Background(), c, "", cfg, devs("d1"))
	noErr(t, err)
	assert.Len(t, notReady, 1)
	s.set("d1", nsoSynced)
	_, notReady, err = waitNSODevices(context.Background(), c, "", cfg, devs("d1"))
	noErr(t, err)
	assert.Len(t, notReady, 0, "the snapshot is still fresh and shows NOT_SYNCED; the read by name shows SYNCED")
	assert.Equal(t, int32(1), s.full.Load())
}

// A timeout of 0 checks once.
func TestNSODevicesTimeoutDescribesTheDevices(t *testing.T) {
	s, c := newInvStub(t, "d1", "d2")
	s.set("d2", "SYNC_FROM_STARTED")
	_, notReady, err := waitNSODevices(context.Background(), c, "", cfgOf(nsoModeSynced, 0), devs("d1", "d2", "d3"))
	noErr(t, err)
	assert.Equal(t, []string{
		"d2: reachability CONN_STATE_REACHABLE, nso SYNC_FROM_STARTED, admin ROBOT_ADMIN_STATE_UP, waited 0s",
		"d3: not in the inventory, waited 0s",
	}, notReady)
}

// The stub ignores the filter, so the row is matched by name.
func TestNSODevicesFullModeAndFilterIgnoredByTheServer(t *testing.T) {
	s, c := newInvStub(t, "d1", "d2")
	s.ignoreFilter = true
	s.mu.Lock()
	s.nodes["d2"].Admin = "ROBOT_ADMIN_STATE_DOWN"
	s.mu.Unlock()
	_, notReady, err := waitNSODevices(context.Background(), c, "", cfgOf(nsoModeFull, 0), devs("d1", "d2"))
	noErr(t, err)
	assert.Len(t, notReady, 1)
	assert.Contains(t, notReady[0], "d2: reachability CONN_STATE_REACHABLE, nso SYNCED, admin ROBOT_ADMIN_STATE_DOWN")
	_, notReady, err = waitNSODevices(context.Background(), c, "", cfgOf(nsoModeSynced, 0), devs("d1", "d2"))
	noErr(t, err)
	assert.Len(t, notReady, 0, "synced takes a device that is admin down")
}

func TestNSODevicesNoneAndPresent(t *testing.T) {
	s, c := newInvStub(t, "d1")
	s.set("d1", "NOT_SYNCED")
	nodes, notReady, err := waitNSODevices(context.Background(), c, "", cfgOf(nsoModeNone, 0), devs("d1", "missing"))
	noErr(t, err)
	assert.Len(t, notReady, 0)
	assert.Len(t, nodes, 0)
	assert.Equal(t, int32(0), s.full.Load()+s.byName.Load(), "mode none reads nothing")
	nodes, notReady, err = waitNSODevices(context.Background(), c, "", cfgOf(nsoModePresent, 0), devs("d1", "missing"))
	noErr(t, err)
	assert.Equal(t, "uuid-d1", nodes["d1"].UUID)
	assert.Equal(t, []string{"missing: not in the inventory, waited 0s"}, notReady)
}

func TestNSODevicesConcurrentCallersShareOneRead(t *testing.T) {
	s, c := newInvStub(t, "d1", "d2", "d3")
	cfg := cfgOf(nsoModeSynced, 0)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, notReady, err := waitNSODevices(context.Background(), c, "", cfg, devs("d1", "d2", "d3"))
			if err != nil || len(notReady) != 0 {
				t.Errorf("err %v, not ready %v", err, notReady)
			}
		}()
	}
	wg.Wait()
	assert.Equal(t, int32(1), s.full.Load())
}

func TestNSODevicesSnapshotIsPerInstance(t *testing.T) {
	s, c := newInvStub(t, "d1")
	cfg := cfgOf(nsoModeSynced, 0)
	for _, instance := range []string{"", "instance2", ""} {
		_, _, err := waitNSODevices(context.Background(), c, instance, cfg, devs("d1"))
		noErr(t, err)
	}
	assert.Equal(t, int32(2), s.full.Load(), "one read for each instance")
}

func TestNSODevicesContextCancelled(t *testing.T) {
	s, c := newInvStub(t, "d1")
	s.set("d1", "NOT_SYNCED")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(30 * time.Millisecond)
		cancel()
	}()
	_, _, err := waitNSODevices(ctx, c, "", cfgOf(nsoModeSynced, time.Minute), devs("d1"))
	assert.True(t, err == context.Canceled, "the wait ends with the cancelled context")
}

func TestWaitForNSODevicesDiagnostics(t *testing.T) {
	s, c := newInvStub(t, "d1")
	s.set("d1", "FETCH_SSH_KEYS_FAILED")
	var diags diag.Diagnostics
	ok := waitForNSODevices(context.Background(), c, "", cfgOf(nsoModeSynced, 0), devs("d1"), &diags)
	assert.False(t, ok)
	assert.Len(t, diags, 1)
	assert.Contains(t, diags[0].Detail(), `device_ready_mode "synced"`)
	assert.Contains(t, diags[0].Detail(), "d1: reachability CONN_STATE_REACHABLE, nso FETCH_SSH_KEYS_FAILED")
	diags = nil
	assert.True(t, waitForNSODevices(context.Background(), c, "", cfgOf(nsoModeSynced, 0), nsoDeviceSet{}, &diags), "no device, no wait")
	assert.True(t, waitForNSODevices(context.Background(), c, "", cfgOf(nsoModeNone, 0), devs("d1"), &diags), "mode none")
	assert.Len(t, diags, 0)
}

func TestNSODeviceSets(t *testing.T) {
	set := nsoDeviceSet{}
	nsoAddDevice(set, types.StringValue("d1"))
	nsoAddDevice(set, types.StringValue(""))
	nsoAddDevice(set, types.StringNull())
	nsoAddDevice(set, types.StringUnknown())
	nsoAddDevice(set, types.StringValue("d2"))
	assert.Equal(t, []string{"d1", "d2"}, set.sorted())
	assert.Equal(t, []string{"d3"}, nsoDevicesAdded(devs("d1", "d3"), devs("d1", "d2")).sorted())
	assert.Len(t, nsoDevicesAdded(devs("d1"), devs("d1", "d2")), 0)
}

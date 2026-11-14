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
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// A resource that names a device (an attribute with reference_nso_device in its definition) waits, before it is
// created or updated, until NSO knows the device: Crosswork accepts a device at once, but NSO onboards it later
// and refuses a reference to it until then. The wait reads the Crosswork inventory.

const (
	nsoNodesQuery = "/crosswork/inventory/v1/nodes/query"

	nsoReady   = "CONN_STATE_REACHABLE"
	nsoSynced  = "SYNCED"
	nsoAdminUp = "ROBOT_ADMIN_STATE_UP"
)

const (
	// The values of device_ready_mode.
	nsoModeNone = "none"
	// synced is needed to push a service or a job.
	nsoModeSynced = "synced"
	nsoModeFull   = "full"
	// The golden-config jobs use it for mode none, because a job cannot run without the uuid of the device.
	nsoModePresent = "present"
)

const (
	nsoDefaultTimeout = 300 * time.Second
	nsoSnapshotTTL    = 15 * time.Second
)

// The waits are variables so the tests can shorten them.
var (
	nsoReadinessInterval = 10 * time.Second
	nsoNow               = time.Now
)

type nsoNode struct {
	Name         string
	UUID         string
	HostName     string
	Reachability string
	NSO          string
	Admin        string
}

func (n nsoNode) ready(mode string) bool {
	switch mode {
	case nsoModePresent:
		return true
	case nsoModeFull:
		return n.Reachability == nsoReady && n.NSO == nsoSynced && n.Admin == nsoAdminUp
	}
	return n.NSO == nsoSynced
}

func (n nsoNode) describe() string {
	return fmt.Sprintf("%s: reachability %s, nso %s, admin %s", n.Name, n.Reachability, n.NSO, n.Admin)
}

func nsoNodeFrom(r interface{ Get(string) gjson.Result }) nsoNode {
	return nsoNode{Name: r.Get("inventory_id").String(), UUID: r.Get("uuid").String(), HostName: r.Get("host_name").String(),
		Reachability: r.Get("reachability_state").String(), NSO: r.Get("nso_state").String(), Admin: r.Get("admin_state").String()}
}

type nsoDeviceSet map[string]struct{}

func nsoAddDevice(set nsoDeviceSet, v types.String) {
	if v.IsNull() || v.IsUnknown() || v.ValueString() == "" {
		return
	}
	set[v.ValueString()] = struct{}{}
}

// An update waits only for the devices it adds.
func nsoDevicesAdded(plan, state nsoDeviceSet) nsoDeviceSet {
	added := nsoDeviceSet{}
	for name := range plan {
		if _, ok := state[name]; !ok {
			added[name] = struct{}{}
		}
	}
	return added
}

func (s nsoDeviceSet) sorted() []string {
	names := make([]string, 0, len(s))
	for name := range s {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// The zero value (a resource that was never configured) means the defaults.
type nsoDeviceConfig struct {
	Mode    string
	Timeout time.Duration
	Store   *nsoSnapshots
}

func (c nsoDeviceConfig) withDefaults() nsoDeviceConfig {
	if c.Mode == "" {
		c.Mode = nsoModeSynced
		c.Timeout = nsoDefaultTimeout
	}
	if c.Store == nil {
		c.Store = newNsoSnapshots()
	}
	return c
}

// Each provider process has its own, so nothing outlives a plan or an apply.
type nsoSnapshots struct {
	mu      sync.Mutex
	entries map[string]*nsoSnapshot
}

type nsoSnapshot struct {
	mu    sync.Mutex // held while the inventory is read, so concurrent callers share one read
	at    time.Time
	nodes map[string]nsoNode
}

func newNsoSnapshots() *nsoSnapshots {
	return &nsoSnapshots{entries: map[string]*nsoSnapshot{}}
}

func (s *nsoSnapshots) entry(instance string) *nsoSnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[instance]
	if !ok {
		e = &nsoSnapshot{}
		s.entries[instance] = e
	}
	return e
}

// The caller must not change the result.
func (s *nsoSnapshots) get(c *cw.Client, instance string) (map[string]nsoNode, error) {
	e := s.entry(instance)
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.nodes != nil && nsoNow().Sub(e.at) < nsoSnapshotTTL {
		return e.nodes, nil
	}
	rows, err := c.Paginate(cw.NewCursorPager(nsoNodesQuery, "{}", 100))
	if err != nil {
		return nil, err
	}
	nodes := make(map[string]nsoNode, len(rows))
	for _, r := range rows {
		n := nsoNodeFrom(r)
		nodes[n.Name] = n
	}
	e.nodes, e.at = nodes, nsoNow()
	return nodes, nil
}

// Never reads from the snapshot. The filter is case-insensitive and the row is matched by the name as written.
func nsoReadNode(c *cw.Client, name string) (nsoNode, bool, error) {
	body, _ := sjson.Set("{}", "filter.inventory_id", name)
	rows, err := c.Paginate(cw.NewCursorPager(nsoNodesQuery, body, 100))
	if err != nil {
		return nsoNode{}, false, err
	}
	for _, r := range rows {
		if n := nsoNodeFrom(r); n.Name == name {
			return n, true, nil
		}
	}
	return nsoNode{}, false, nil
}

// A zero timeout checks once. Mode none makes no check and reads nothing.
//
// A device that the snapshot shows as ready is accepted. A device that is not ready, unknown or missing is never
// judged by the snapshot: it is read again by name, at each poll interval, until it is ready or the time is over.
func waitNSODevices(ctx context.Context, c *cw.Client, instance string, cfg nsoDeviceConfig, devices nsoDeviceSet) (map[string]nsoNode, []string, error) {
	cfg = cfg.withDefaults()
	nodes := map[string]nsoNode{}
	if cfg.Mode == nsoModeNone || len(devices) == 0 {
		return nodes, nil, nil
	}
	snapshot, err := cfg.Store.get(c, instance)
	if err != nil {
		return nil, nil, err
	}
	pending := devices.sorted()
	last := map[string]string{}
	deadline := nsoNow().Add(cfg.Timeout)
	for first := true; ; first = false {
		var still []string
		for _, name := range pending {
			if first {
				if n, ok := snapshot[name]; ok && n.ready(cfg.Mode) {
					nodes[name] = n
					continue
				}
			}
			n, ok, err := nsoReadNode(c, name)
			if err != nil {
				return nil, nil, err
			}
			if ok && n.ready(cfg.Mode) {
				nodes[name] = n
				continue
			}
			if ok {
				last[name] = n.describe()
			} else {
				last[name] = name + ": not in the inventory"
			}
			still = append(still, name)
		}
		if len(still) == 0 {
			return nodes, nil, nil
		}
		pending = still
		wait := nsoReadinessInterval
		if left := deadline.Sub(nsoNow()); left < wait {
			wait = left
		}
		if wait <= 0 {
			notReady := make([]string, 0, len(pending))
			for _, name := range pending {
				notReady = append(notReady, fmt.Sprintf("%s, waited %ds", last[name], int(cfg.Timeout.Seconds())))
			}
			return nodes, notReady, nil
		}
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		case <-time.After(wait):
		}
	}
}

// Returns false when the resource must not be written.
func waitForNSODevices(ctx context.Context, c *cw.Client, instance string, cfg nsoDeviceConfig, devices nsoDeviceSet, diags *diag.Diagnostics) bool {
	if len(devices) == 0 {
		return true
	}
	_, notReady, err := waitNSODevices(ctx, c, instance, cfg, devices)
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Failed to read the inventory (POST), got error: %s", err))
		return false
	}
	if len(notReady) > 0 {
		cfg = cfg.withDefaults()
		diags.AddError("Devices not ready",
			fmt.Sprintf("NSO does not have these devices as the provider setting device_ready_mode %q needs (the wait is limited by device_ready_timeout, now %d seconds; the mode none turns it off):\n%s",
				cfg.Mode, int(cfg.Timeout.Seconds()), strings.Join(notReady, "\n")))
		return false
	}
	return true
}

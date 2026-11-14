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

// Hand-written operations of crosswork_data_gateway_pools, which the generator cannot express: two key casings in
// one body, numbers sent as strings, an id that comes back as "<uuid> <name>", the data gateways of a pool read
// from a second call, and a change that is repeated while the controller settles.

import (
	"context"
	"fmt"
	"net/netip"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	hapoolPath        = "/crosswork/dg-manager/v2/hapool"
	hapoolQueryPath   = "/crosswork/dg-manager/v2/hapool/query"
	vdgQueryPath      = "/crosswork/dg-manager/v2/vdg/query"
	instanceQueryPath = "/crosswork/dg-manager/v1/dg/query"

	defaultCdgReadinessTimeout = 60
)

// A variable so that tests can shorten it.
var poolRetryInterval = 2 * time.Second

// poolTransient are the messages of the controller that mean "not now": another change of the pool is running, an
// instance has not settled after it left a pool, or the instances are not ready. Any other message is final.
var poolTransient = []string{
	"is currently locked",
	"is not in Not_Ready state",
	"number of DGs in Operational NOT_READY state",
}

func poolTransientError(message string) bool {
	for _, m := range poolTransient {
		if strings.Contains(message, m) {
			return true
		}
	}
	return false
}

func poolTimeout(data DataGatewayPools) time.Duration {
	seconds := int64(defaultCdgReadinessTimeout)
	if !data.CdgReadinessTimeout.IsNull() && !data.CdgReadinessTimeout.IsUnknown() && data.CdgReadinessTimeout.ValueInt64() > 0 {
		seconds = data.CdgReadinessTimeout.ValueInt64()
	}
	return time.Duration(seconds) * time.Second
}

func sleepContext(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func retryPool(ctx context.Context, timeout time.Duration, op func() (cw.Res, error)) (cw.Res, error) {
	deadline := time.Now().Add(timeout)
	for {
		res, err := op()
		if err == nil || !poolTransientError(err.Error()+res.String()) || !time.Now().Add(poolRetryInterval).Before(deadline) {
			return res, err
		}
		if !sleepContext(ctx, poolRetryInterval) {
			return res, err
		}
	}
}

func listOfStrings(l types.List) []string {
	if l.IsNull() || l.IsUnknown() {
		return nil
	}
	out := make([]string, 0, len(l.Elements()))
	for _, e := range l.Elements() {
		if s, ok := e.(types.String); ok {
			out = append(out, s.ValueString())
		}
	}
	return out
}

func stringList(values []string) types.List {
	elems := make([]attr.Value, len(values))
	for i, v := range values {
		elems[i] = types.StringValue(v)
	}
	return types.ListValueMust(types.StringType, elems)
}

// The controller keeps an address as it was written, so two spellings of one address are compared in this form.
func normalizeAddress(s string) string {
	a, err := netip.ParseAddr(strings.TrimSpace(s))
	if err != nil {
		return s
	}
	return a.String()
}

func addressFamily(s string) string {
	a, err := netip.ParseAddr(strings.TrimSpace(s))
	switch {
	case err != nil:
		return ""
	case a.Is4():
		return "ROBOT_INET_ADDR_TYPE_v4"
	default:
		return "ROBOT_INET_ADDR_TYPE_v6"
	}
}

// The numbers are strings, and the keys are camelCase except those of the addresses and enable_fqdn.
func poolRequestBody(data DataGatewayPools, id string) string {
	pool := ""
	if id != "" {
		pool, _ = sjson.Set(pool, "puuid", id)
	}
	pool, _ = sjson.Set(pool, "name", data.Name.ValueString())
	if !data.Description.IsNull() && !data.Description.IsUnknown() && data.Description.ValueString() != "" {
		pool, _ = sjson.Set(pool, "description", data.Description.ValueString())
	}
	pool, _ = sjson.Set(pool, "haStrategy", "VIP_BASED")
	mask := strconv.FormatInt(data.PrefixLength.ValueInt64(), 10)
	for i, a := range listOfStrings(data.VirtualIps) {
		pool, _ = sjson.Set(pool, fmt.Sprintf("ipaddrs.%d.inetaddrs.0.inet_af", i), addressFamily(a))
		pool, _ = sjson.Set(pool, fmt.Sprintf("ipaddrs.%d.inetaddrs.0.inet_addr", i), strings.TrimSpace(a))
		pool, _ = sjson.Set(pool, fmt.Sprintf("ipaddrs.%d.inetaddrs.0.mask", i), mask)
	}
	pool, _ = sjson.Set(pool, "pdgUuids", listOfStrings(data.Instances))
	pool, _ = sjson.Set(pool, "spares", strconv.FormatInt(data.Spares.ValueInt64(), 10))
	pool, _ = sjson.Set(pool, "subnetMask", mask)
	pool, _ = sjson.Set(pool, "enable_fqdn", false)
	return `{"pools":[` + pool + `]}`
}

// poolIdFromAnswer reads the id from the answer of a create: `{"puuid":["<uuid> <name>"]}`.
func poolIdFromAnswer(res cw.Res) string {
	fields := strings.Fields(res.Get("puuid.0").String())
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func findPool(res cw.Res, id, name string) gjson.Result {
	var found gjson.Result
	res.Get("data").ForEach(func(_, pool gjson.Result) bool {
		if (id != "" && pool.Get("puuid").String() == id) || (id == "" && name != "" && pool.Get("name").String() == name) {
			found = pool
			return false
		}
		return true
	})
	return found
}

func poolVirtualIps(pool gjson.Result) []string {
	var out []string
	pool.Get("ipaddrs").ForEach(func(_, entry gjson.Result) bool {
		entry.Get("inetaddrs").ForEach(func(_, a gjson.Result) bool {
			out = append(out, a.Get("inetAddr").String())
			return true
		})
		return true
	})
	return out
}

func sameAddresses(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if normalizeAddress(a[i]) != normalizeAddress(b[i]) {
			return false
		}
	}
	return true
}

// gatewaySuffix is the number n of a data gateway named "<pool>-<n>", 0 when the name has no such suffix.
func gatewaySuffix(name string) int {
	i := strings.LastIndex(name, "-")
	if i < 0 {
		return 0
	}
	n, err := strconv.Atoi(name[i+1:])
	if err != nil {
		return 0
	}
	return n
}

func poolGatewayIds(res cw.Res, poolId string) []string {
	type gateway struct {
		n    int
		name string
		id   string
	}
	var gateways []gateway
	res.Get("data").ForEach(func(_, g gjson.Result) bool {
		if g.Get("poolId").String() == poolId {
			gateways = append(gateways, gateway{gatewaySuffix(g.Get("name").String()), g.Get("name").String(), g.Get("vdgUuid").String()})
		}
		return true
	})
	sort.Slice(gateways, func(i, j int) bool {
		if gateways[i].n != gateways[j].n {
			return gateways[i].n < gateways[j].n
		}
		return gateways[i].name < gateways[j].name
	})
	ids := make([]string, len(gateways))
	for i, g := range gateways {
		ids[i] = g.id
	}
	return ids
}

func listPools(c *cw.Client) (cw.Res, error) {
	return c.Post(hapoolQueryPath, `{}`, cw.ReadOnly)
}

func listPoolGateways(c *cw.Client, poolId string) ([]string, error) {
	res, err := c.Post(vdgQueryPath, `{}`, cw.ReadOnly)
	if err != nil {
		return nil, err
	}
	return poolGatewayIds(res, poolId), nil
}

// The controller lists the data gateways after it created the pool.
func waitPoolGateways(ctx context.Context, c *cw.Client, poolId string, want int, timeout time.Duration) ([]string, error) {
	deadline := time.Now().Add(timeout)
	for {
		ids, err := listPoolGateways(c, poolId)
		if err != nil {
			return ids, err
		}
		if len(ids) >= want {
			return ids, nil
		}
		if !time.Now().Add(poolRetryInterval).Before(deadline) || !sleepContext(ctx, poolRetryInterval) {
			return ids, fmt.Errorf("the pool exists, but %d of its %d data gateways are listed after %s", len(ids), want, timeout)
		}
	}
}

func waitPoolGone(ctx context.Context, c *cw.Client, id string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		res, err := listPools(c)
		if err != nil {
			return err
		}
		if !findPool(res, id, "").Exists() {
			return nil
		}
		if !time.Now().Add(poolRetryInterval).Before(deadline) || !sleepContext(ctx, poolRetryInterval) {
			return fmt.Errorf("the pool %s is still listed after %s", id, timeout)
		}
	}
}

// The id is returned even when the wait fails: the pool exists and has to be in the state.
func createPool(ctx context.Context, c *cw.Client, plan DataGatewayPools) (string, []string, error) {
	timeout := poolTimeout(plan)
	body := poolRequestBody(plan, "")
	res, err := retryPool(ctx, timeout, func() (cw.Res, error) { return c.Post(hapoolPath, body) })
	if err != nil {
		return "", nil, fmt.Errorf("%w, %s", err, res.String())
	}
	id := poolIdFromAnswer(res)
	if id == "" {
		return "", nil, fmt.Errorf("the answer of the controller has no pool id: %s", res.String())
	}
	ids, err := waitPoolGateways(ctx, c, id, len(listOfStrings(plan.VirtualIps)), timeout)
	return id, ids, err
}

func updatePool(ctx context.Context, c *cw.Client, plan DataGatewayPools, id string) error {
	body := poolRequestBody(plan, id)
	res, err := retryPool(ctx, poolTimeout(plan), func() (cw.Res, error) { return c.Put(hapoolPath, body) })
	if err != nil {
		return fmt.Errorf("%w, %s", err, res.String())
	}
	return nil
}

// The controller refuses while a data gateway of the pool has devices attached.
func deletePool(ctx context.Context, c *cw.Client, state DataGatewayPools) error {
	timeout := poolTimeout(state)
	body, _ := sjson.Set(`{"pools":[{}]}`, "pools.0.name", state.Name.ValueString())
	body, _ = sjson.Set(body, "pools.0.puuid", state.Id.ValueString())
	res, err := retryPool(ctx, timeout, func() (cw.Res, error) { return c.DeleteBody(hapoolPath, body) })
	if err != nil {
		return fmt.Errorf("%w, %s", err, res.String())
	}
	return waitPoolGone(ctx, c, state.Id.ValueString(), timeout)
}

func poolIdByNameOrId(c *cw.Client, key string) (string, error) {
	res, err := listPools(c)
	if err != nil {
		return "", err
	}
	if pool := findPool(res, key, ""); pool.Exists() {
		return pool.Get("puuid").String(), nil
	}
	return findPool(res, "", key).Get("puuid").String(), nil
}

// The controller returns the subnet mask as a string.
func poolPrefixLength(pool gjson.Result) int64 {
	if n, err := strconv.ParseInt(pool.Get("subnetMask").String(), 10, 64); err == nil {
		return n
	}
	if n, err := strconv.ParseInt(pool.Get("ipaddrs.0.inetaddrs.0.mask").String(), 10, 64); err == nil {
		return n
	}
	return 0
}

func poolInstanceIds(pool gjson.Result) []string {
	var out []string
	pool.Get("pdgUuids").ForEach(func(_, v gjson.Result) bool {
		out = append(out, v.String())
		return true
	})
	return out
}

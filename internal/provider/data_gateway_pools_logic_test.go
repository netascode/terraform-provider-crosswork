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
	"errors"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

func testPool(addresses []string, prefix, spares int64, instances ...string) DataGatewayPools {
	return DataGatewayPools{
		Name:         types.StringValue("cw-dg-pool-1"),
		Description:  types.StringValue("Example pool"),
		VirtualIps:   stringList(addresses),
		PrefixLength: types.Int64Value(prefix),
		Spares:       types.Int64Value(spares),
		Instances:    stringList(instances),
	}
}

func TestPoolRequestBodyIPv4(t *testing.T) {
	body := gjson.Parse(poolRequestBody(testPool([]string{"192.0.2.195"}, 24, 1, "dg-1", "dg-2"), ""))
	pool := body.Get("pools.0")
	assert.Equal(t, "cw-dg-pool-1", pool.Get("name").String())
	assert.Equal(t, "VIP_BASED", pool.Get("haStrategy").String())
	assert.Equal(t, "ROBOT_INET_ADDR_TYPE_v4", pool.Get("ipaddrs.0.inetaddrs.0.inet_af").String())
	assert.Equal(t, "192.0.2.195", pool.Get("ipaddrs.0.inetaddrs.0.inet_addr").String())
	// the controller wants numbers as strings
	assert.Equal(t, gjson.String, pool.Get("spares").Type)
	assert.Equal(t, "1", pool.Get("spares").String())
	assert.Equal(t, gjson.String, pool.Get("subnetMask").Type)
	assert.Equal(t, "24", pool.Get("subnetMask").String())
	assert.Equal(t, "24", pool.Get("ipaddrs.0.inetaddrs.0.mask").String())
	assert.Equal(t, 2, len(pool.Get("pdgUuids").Array()))
	assert.False(t, pool.Get("enable_fqdn").Bool())
	assert.False(t, pool.Get("puuid").Exists(), "a create has no id")
}

func TestPoolRequestBodyUpdateAndDescription(t *testing.T) {
	pool := testPool([]string{"192.0.2.195"}, 24, 0, "dg-1")
	assert.Equal(t, "abc-1", gjson.Get(poolRequestBody(pool, "abc-1"), "pools.0.puuid").String())
	pool.Description = types.StringNull()
	assert.False(t, gjson.Get(poolRequestBody(pool, ""), "pools.0.description").Exists(), "no description is sent when there is none")
}

func TestPoolRequestBodyIPv6TwoAddresses(t *testing.T) {
	pool := gjson.Parse(poolRequestBody(testPool([]string{"fd00:34::10", "fd00:34::11"}, 64, 0, "dg-1", "dg-2"), "")).Get("pools.0")
	assert.Equal(t, "ROBOT_INET_ADDR_TYPE_v6", pool.Get("ipaddrs.0.inetaddrs.0.inet_af").String())
	assert.Equal(t, "fd00:34::11", pool.Get("ipaddrs.1.inetaddrs.0.inet_addr").String(), "one entry for each virtual IP, in order")
	assert.Equal(t, "64", pool.Get("subnetMask").String())
}

func TestNormalizeAndSameAddresses(t *testing.T) {
	assert.Equal(t, "fd00:34::10", normalizeAddress("FD00:34:0:0:0:0:0:10"))
	assert.Equal(t, "192.0.2.195", normalizeAddress(" 192.0.2.195 "))
	assert.Equal(t, "not-an-address", normalizeAddress("not-an-address"))
	assert.True(t, sameAddresses([]string{"FD00:34:0:0:0:0:0:10"}, []string{"fd00:34::10"}))
	assert.False(t, sameAddresses([]string{"fd00:34::10"}, []string{"fd00:34::11"}))
	assert.False(t, sameAddresses([]string{"a", "b"}, []string{"a"}))
	assert.Equal(t, "ROBOT_INET_ADDR_TYPE_v6", addressFamily("fd00::1"))
	assert.Equal(t, "", addressFamily("fd00:::1"))
}

func TestPoolAnswers(t *testing.T) {
	assert.Equal(t, "9f5c5691-2e2f-47da-982f-3b0732e0725e", poolIdFromAnswer(gjson.Parse(`{"puuid":["9f5c5691-2e2f-47da-982f-3b0732e0725e cw-dg-pool-1"]}`)))
	assert.Equal(t, "", poolIdFromAnswer(gjson.Parse(`{}`)))
	list := gjson.Parse(`{"data":[{"puuid":"p1","name":"a"},{"puuid":"p2","name":"b"}],"totalCount":2}`)
	assert.Equal(t, "b", findPool(list, "p2", "").Get("name").String())
	assert.Equal(t, "p1", findPool(list, "", "a").Get("puuid").String())
	assert.False(t, findPool(list, "", "").Exists(), "neither an id nor a name finds nothing")
	assert.False(t, findPool(list, "p9", "").Exists())
}

func TestPoolFromBodyAbsentSparesAndSpelling(t *testing.T) {
	res := gjson.Parse(`{"puuid":"p1","name":"cw-dg-pool-1","ipaddrs":[{"inetaddrs":[{"inetAf":"ROBOT_INET_ADDR_TYPE_v6","inetAddr":"FD00:34:0:0:0:0:0:10","mask":"64"}]}],
	  "pdgUuids":["dg-1"],"protectionStatus":"NOT_PLANNED","subnetMask":"64","haStrategy":"VIP_BASED"}`)
	var data DataGatewayPools
	data.fromBody(context.Background(), res)
	assert.Equal(t, int64(0), data.Spares.ValueInt64(), "the controller leaves spares out at 0")
	assert.Equal(t, int64(64), data.PrefixLength.ValueInt64())
	assert.Equal(t, []string{"fd00:34::10"}, listOfStrings(data.VirtualIps), "an import reads the canonical address")
	assert.True(t, data.Description.IsNull())
	assert.Equal(t, int64(60), data.CdgReadinessTimeout.ValueInt64(), "the wait defaults when it is not in the state")

	// a read keeps the spelling of the configuration when it is the same address
	data = testPool([]string{"FD00:34:0:0:0:0:0:10"}, 64, 0, "dg-1")
	data.updateFromBody(context.Background(), res)
	assert.Equal(t, []string{"FD00:34:0:0:0:0:0:10"}, listOfStrings(data.VirtualIps))
	data.VirtualIps = stringList([]string{"fd00:34::99"})
	data.updateFromBody(context.Background(), res)
	assert.Equal(t, []string{"fd00:34::10"}, listOfStrings(data.VirtualIps), "another address is a difference")
}

func TestPoolGatewayIdsOrder(t *testing.T) {
	res := gjson.Parse(`{"data":[
	  {"name":"pool-10","vdgUuid":"g10","poolId":"p1"},
	  {"name":"other-1","vdgUuid":"o1","poolId":"p2"},
	  {"name":"pool-2","vdgUuid":"g2","poolId":"p1"},
	  {"name":"pool-1","vdgUuid":"g1","poolId":"p1"}]}`)
	assert.Equal(t, []string{"g1", "g2", "g10"}, poolGatewayIds(res, "p1"), "by the number of the name, not as text, and only the pool's own")
	assert.Equal(t, 0, len(poolGatewayIds(res, "p9")))
	assert.Equal(t, 1, gatewaySuffix("cw-dg-pool-1-1"))
	assert.Equal(t, 12, gatewaySuffix("a-12"))
	assert.Equal(t, 0, gatewaySuffix("nosuffix"))
}

func TestPoolTransientError(t *testing.T) {
	for _, m := range []string{
		"update HA Pool cannot be processed since pool: cwtest is currently locked",
		"The High Availability (HA) pool update failed because the data gateway cw72-dg-1a is not in Not_Ready state. ",
		"number of DGs in Operational NOT_READY state: 0 are not enough to satisfy number of DGs needed: 1",
	} {
		assert.True(t, poolTransientError(m), m)
	}
	for _, m := range []string{
		"HA Pool with name cwtest already exists",
		"Virtual Data Gateway g in the pool p has devices attached to it. Please detach devices before deleting the pool.",
		"",
	} {
		assert.False(t, poolTransientError(m), m)
	}
}

func TestPoolTimeout(t *testing.T) {
	assert.Equal(t, 60*time.Second, poolTimeout(DataGatewayPools{}))
	assert.Equal(t, 90*time.Second, poolTimeout(DataGatewayPools{CdgReadinessTimeout: types.Int64Value(90)}))
}

func TestRetryPool(t *testing.T) {
	old := poolRetryInterval
	poolRetryInterval = time.Millisecond
	t.Cleanup(func() { poolRetryInterval = old })
	locked := errors.New("HTTP Request failed: StatusCode 500: update HA Pool cannot be processed since pool: x is currently locked")

	calls := 0
	_, err := retryPool(context.Background(), time.Second, func() (cw.Res, error) {
		calls++
		if calls < 3 {
			return cw.Res{}, locked
		}
		return cw.Res{}, nil
	})
	assert.True(t, err == nil, "a transient refusal is repeated until it succeeds")
	assert.Equal(t, 3, calls)

	calls = 0
	final := errors.New("HTTP Request failed: StatusCode 500: HA Pool with name x already exists")
	_, err = retryPool(context.Background(), time.Second, func() (cw.Res, error) { calls++; return cw.Res{}, final })
	assert.Equal(t, 1, calls, "any other message ends at once")
	assert.Equal(t, final, err, "and is returned as the controller wrote it")

	calls = 0
	_, err = retryPool(context.Background(), 20*time.Millisecond, func() (cw.Res, error) { calls++; return cw.Res{}, locked })
	assert.True(t, calls > 1, "a transient refusal is repeated while the time lasts")
	assert.Equal(t, locked, err, "and the last message is returned unchanged when it runs out")

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls = 0
	_, err = retryPool(ctx, time.Minute, func() (cw.Res, error) { calls++; return cw.Res{}, locked })
	assert.Equal(t, 1, calls, "a cancelled run does not wait")
	assert.Equal(t, locked, err)
}

func TestPoolReadMatchesConfig(t *testing.T) {
	res := gjson.Parse(`{"puuid":"00000000-0000-4000-8000-000000000001","name":"cw-dg-pool-1",
	  "ipaddrs":[{"inetaddrs":[{"inetAf":"ROBOT_INET_ADDR_TYPE_v4","inetAddr":"192.0.2.195","mask":"24"}]}],
	  "pdgUuids":["00000000-0000-4000-8000-000000000002","00000000-0000-4000-8000-000000000003"],"spares":"1","protectionStatus":"PROTECTED",
	  "description":"cw-dg-pool-1 description","subnetMask":"24","haStrategy":"VIP_BASED","balanced":true}`)
	want := DataGatewayPools{
		Name:                types.StringValue("cw-dg-pool-1"),
		Description:         types.StringValue("cw-dg-pool-1 description"),
		VirtualIps:          stringList([]string{"192.0.2.195"}),
		PrefixLength:        types.Int64Value(24),
		Spares:              types.Int64Value(1),
		Instances:           stringList([]string{"00000000-0000-4000-8000-000000000002", "00000000-0000-4000-8000-000000000003"}),
		CdgReadinessTimeout: types.Int64Value(60),
	}
	got := want
	got.updateFromBody(context.Background(), res)
	assert.True(t, got.Name.Equal(want.Name) && got.Description.Equal(want.Description) && got.VirtualIps.Equal(want.VirtualIps) &&
		got.PrefixLength.Equal(want.PrefixLength) && got.Spares.Equal(want.Spares) && got.Instances.Equal(want.Instances), "a read changes nothing")
	var imported DataGatewayPools
	imported.fromBody(context.Background(), res)
	assert.Equal(t, "00000000-0000-4000-8000-000000000001", imported.Id.ValueString())
	assert.True(t, imported.VirtualIps.Equal(want.VirtualIps) && imported.Instances.Equal(want.Instances) && imported.Spares.Equal(want.Spares), "an import reads the same values")
}

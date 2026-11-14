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

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

func gjsonRes(s string) cw.Res { return gjson.Parse(s) }

func TestDataGatewayInstancesDataSource(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newPoolStub(t)
	const ds = "data.crosswork_data_gateway_instances.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: s.providerBlock() + "\ndata \"crosswork_data_gateway_instances\" \"test\" {\n}\n",
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(ds, "items.#", "2"),
				// an unassigned instance has no pool and no data gateway
				resource.TestCheckTypeSetElemNestedAttrs(ds, "items.*", map[string]string{"id": "dg-1", "name": "cw72-dg-1a", "role": "UNASSIGNED", "oper_state": "OS_NOT_READY", "admin_state": "AS_UP"}),
				resource.TestCheckTypeSetElemNestedAttrs(ds, "items.*", map[string]string{"id": "dg-2", "name": "cw72-dg-2a", "role": "ASSIGNED", "pool_id": "pool-id-1", "vdg_id": "vdg-pool-id-1-1", "oper_state": "OS_UP"}),
			),
		}},
	})
}

func TestParseInstancesSortsByNameAndLeavesTheAbsentNull(t *testing.T) {
	items := parseInstances(gjsonRes(`{"data":[
	  {"duuid":"b","name":"zeta","configData":{"role":"SPARE","poolId":"p1"},"operationalData":{"operState":"OS_NOT_READY"}},
	  {"duuid":"a","name":"alpha","configData":{"adminState":"AS_UP"}}]}`))
	assert.Equal(t, 2, len(items))
	assert.Equal(t, "alpha", items[0].Name.ValueString())
	assert.True(t, items[0].Role.IsNull() && items[0].PoolId.IsNull() && items[0].VdgId.IsNull() && items[0].OperState.IsNull())
	assert.Equal(t, "p1", items[1].PoolId.ValueString())
	assert.True(t, items[1].VdgId.IsNull(), "a spare has no data gateway")
}

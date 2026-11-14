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

// Written by hand: the instances are virtual machines that enrol themselves with Crosswork, so Terraform only reads them.

import (
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

func optionalString(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

func parseInstances(res cw.Res) []DataGatewayInstancesItems {
	items := []DataGatewayInstancesItems{}
	res.Get("data").ForEach(func(_, d gjson.Result) bool {
		items = append(items, DataGatewayInstancesItems{
			Id:         types.StringValue(d.Get("duuid").String()),
			Name:       types.StringValue(d.Get("name").String()),
			Role:       optionalString(d.Get("configData.role").String()),
			PoolId:     optionalString(d.Get("configData.poolId").String()),
			VdgId:      optionalString(d.Get("configData.vdgUuid").String()),
			OperState:  optionalString(d.Get("operationalData.operState").String()),
			AdminState: optionalString(d.Get("configData.adminState").String()),
		})
		return true
	})
	sort.Slice(items, func(i, j int) bool {
		if items[i].Name.ValueString() != items[j].Name.ValueString() {
			return items[i].Name.ValueString() < items[j].Name.ValueString()
		}
		return items[i].Id.ValueString() < items[j].Id.ValueString()
	})
	return items
}

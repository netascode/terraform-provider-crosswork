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

// Section below is generated&owned by "gen/generator.go". //template:begin imports
import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type DataGatewayInstances struct {
	Instance types.String                `tfsdk:"instance"`
	Id       types.String                `tfsdk:"id"`
	Items    []DataGatewayInstancesItems `tfsdk:"items"`
}

type DataGatewayInstancesItems struct {
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Role       types.String `tfsdk:"role"`
	PoolId     types.String `tfsdk:"pool_id"`
	VdgId      types.String `tfsdk:"vdg_id"`
	OperState  types.String `tfsdk:"oper_state"`
	AdminState types.String `tfsdk:"admin_state"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data DataGatewayInstances) getPath() string {
	return "/crosswork/dg-manager/v1/dg/query"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data DataGatewayInstances) toBody(ctx context.Context, state DataGatewayInstances) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if len(data.Items) > 0 {
		body, _ = sjson.Set(body, "items", []interface{}{})
		for _, item := range data.Items {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueString())
			}
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if !item.Role.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "role", item.Role.ValueString())
			}
			if !item.PoolId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "poolId", item.PoolId.ValueString())
			}
			if !item.VdgId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "vdgId", item.VdgId.ValueString())
			}
			if !item.OperState.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "operState", item.OperState.ValueString())
			}
			if !item.AdminState.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "adminState", item.AdminState.ValueString())
			}
			body, _ = sjson.SetRaw(body, "items.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *DataGatewayInstances) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get(""); value.Exists() {
		data.Id = types.StringValue(value.String())
	} else {
		data.Id = types.StringNull()
	}
	if value := res.Get("items"); value.Exists() && len(value.Array()) > 0 {
		data.Items = make([]DataGatewayInstancesItems, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := DataGatewayInstancesItems{}
			if cValue := v.Get("id"); cValue.Exists() {
				item.Id = types.StringValue(cValue.String())
			} else {
				item.Id = types.StringNull()
			}
			if cValue := v.Get("name"); cValue.Exists() {
				item.Name = types.StringValue(cValue.String())
			} else {
				item.Name = types.StringNull()
			}
			if cValue := v.Get("role"); cValue.Exists() {
				item.Role = types.StringValue(cValue.String())
			} else {
				item.Role = types.StringNull()
			}
			if cValue := v.Get("poolId"); cValue.Exists() {
				item.PoolId = types.StringValue(cValue.String())
			} else {
				item.PoolId = types.StringNull()
			}
			if cValue := v.Get("vdgId"); cValue.Exists() {
				item.VdgId = types.StringValue(cValue.String())
			} else {
				item.VdgId = types.StringNull()
			}
			if cValue := v.Get("operState"); cValue.Exists() {
				item.OperState = types.StringValue(cValue.String())
			} else {
				item.OperState = types.StringNull()
			}
			if cValue := v.Get("adminState"); cValue.Exists() {
				item.AdminState = types.StringValue(cValue.String())
			} else {
				item.AdminState = types.StringNull()
			}
			data.Items = append(data.Items, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *DataGatewayInstances) updateFromBody(ctx context.Context, res gjson.Result) {
	for i := range data.Items {
		keys := [...]string{"id", "name", "role", "poolId", "vdgId", "operState", "adminState"}
		keyValues := [...]string{data.Items[i].Id.ValueString(), data.Items[i].Name.ValueString(), data.Items[i].Role.ValueString(), data.Items[i].PoolId.ValueString(), data.Items[i].VdgId.ValueString(), data.Items[i].OperState.ValueString(), data.Items[i].AdminState.ValueString()}

		var r gjson.Result
		res.Get("items").ForEach(
			func(_, v gjson.Result) bool {
				found := false
				for ik := range keys {
					if v.Get(keys[ik]).String() == keyValues[ik] {
						found = true
						continue
					}
					found = false
					break
				}
				if found {
					r = v
					return false
				}
				return true
			},
		)
		if value := r.Get("id"); value.Exists() && !data.Items[i].Id.IsNull() {
			data.Items[i].Id = types.StringValue(value.String())
		} else {
			data.Items[i].Id = types.StringNull()
		}
		if value := r.Get("name"); value.Exists() && !data.Items[i].Name.IsNull() {
			data.Items[i].Name = types.StringValue(value.String())
		} else {
			data.Items[i].Name = types.StringNull()
		}
		if value := r.Get("role"); value.Exists() && !data.Items[i].Role.IsNull() {
			data.Items[i].Role = types.StringValue(value.String())
		} else {
			data.Items[i].Role = types.StringNull()
		}
		if value := r.Get("poolId"); value.Exists() && !data.Items[i].PoolId.IsNull() {
			data.Items[i].PoolId = types.StringValue(value.String())
		} else {
			data.Items[i].PoolId = types.StringNull()
		}
		if value := r.Get("vdgId"); value.Exists() && !data.Items[i].VdgId.IsNull() {
			data.Items[i].VdgId = types.StringValue(value.String())
		} else {
			data.Items[i].VdgId = types.StringNull()
		}
		if value := r.Get("operState"); value.Exists() && !data.Items[i].OperState.IsNull() {
			data.Items[i].OperState = types.StringValue(value.String())
		} else {
			data.Items[i].OperState = types.StringNull()
		}
		if value := r.Get("adminState"); value.Exists() && !data.Items[i].AdminState.IsNull() {
			data.Items[i].AdminState = types.StringValue(value.String())
		} else {
			data.Items[i].AdminState = types.StringNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *DataGatewayInstances) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.Items) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

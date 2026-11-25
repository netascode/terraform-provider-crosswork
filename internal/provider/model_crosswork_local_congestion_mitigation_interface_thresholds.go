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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type LocalCongestionMitigationInterfaceThresholds struct {
	Instance            types.String                                                      `tfsdk:"instance"`
	Id                  types.String                                                      `tfsdk:"id"`
	DomainId            types.String                                                      `tfsdk:"domain_id"`
	InterfaceThresholds []LocalCongestionMitigationInterfaceThresholdsInterfaceThresholds `tfsdk:"interface_thresholds"`
}

type LocalCongestionMitigationInterfaceThresholdsInterfaceThresholds struct {
	Node                 types.String `tfsdk:"node"`
	Interface            types.String `tfsdk:"interface"`
	UtilizationThreshold types.Int64  `tfsdk:"utilization_threshold"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data LocalCongestionMitigationInterfaceThresholds) getPath() string {
	return "/crosswork/optima/api/v1/functionPacks/lcm/domain/updateInterfaceThresholds"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

func (data LocalCongestionMitigationInterfaceThresholds) getPathDelete() string {
	return "/crosswork/optima/api/v1/functionPacks/lcm/domain/deleteInterfaceThresholds"
}

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data LocalCongestionMitigationInterfaceThresholds) getPathGet() string {
	return "/crosswork/optima/api/v1/functionPacks/lcm/domain/interfaceThresholds"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// Escapes the id so it cannot change the structure of a request body.
func (data LocalCongestionMitigationInterfaceThresholds) idJSON() string {
	id := data.Id.ValueString()
	if id == "" {
		id = fmt.Sprint(data.DomainId.ValueString())
	}
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}

func (data LocalCongestionMitigationInterfaceThresholds) readGet(c *cw.Client, path string) (cw.Res, error) {
	return c.Post(path, strings.ReplaceAll(`{"id":"{id}"}`, "{id}", data.idJSON()), cw.ReadOnly)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data LocalCongestionMitigationInterfaceThresholds) toBody(ctx context.Context, state LocalCongestionMitigationInterfaceThresholds) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.DomainId.IsNull() {
		body, _ = sjson.Set(body, "id", data.DomainId.ValueString())
	}
	if len(data.InterfaceThresholds) > 0 {
		body, _ = sjson.Set(body, "addInts", []interface{}{})
		for _, item := range data.InterfaceThresholds {
			itemBody := ""
			if !item.Node.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "node", item.Node.ValueString())
			}
			if !item.Interface.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "interface", item.Interface.ValueString())
			}
			if !item.UtilizationThreshold.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "threshold", item.UtilizationThreshold.ValueInt64())
			}
			body, _ = sjson.SetRaw(body, "addInts.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *LocalCongestionMitigationInterfaceThresholds) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("data"); value.Exists() && len(value.Array()) > 0 {
		data.InterfaceThresholds = make([]LocalCongestionMitigationInterfaceThresholdsInterfaceThresholds, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := LocalCongestionMitigationInterfaceThresholdsInterfaceThresholds{}
			if cValue := v.Get("node"); cValue.Exists() {
				item.Node = types.StringValue(cValue.String())
			} else {
				item.Node = types.StringNull()
			}
			if cValue := v.Get("interface"); cValue.Exists() {
				item.Interface = types.StringValue(cValue.String())
			} else {
				item.Interface = types.StringNull()
			}
			if cValue := v.Get("threshold"); cValue.Exists() {
				item.UtilizationThreshold = types.Int64Value(cValue.Int())
			} else {
				item.UtilizationThreshold = types.Int64Null()
			}
			data.InterfaceThresholds = append(data.InterfaceThresholds, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *LocalCongestionMitigationInterfaceThresholds) updateFromBody(ctx context.Context, res gjson.Result) {
	for i := range data.InterfaceThresholds {
		keys := [...]string{"node", "interface"}
		keyValues := [...]string{data.InterfaceThresholds[i].Node.ValueString(), data.InterfaceThresholds[i].Interface.ValueString()}

		var r gjson.Result
		res.Get("data").ForEach(
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
		if value := r.Get("node"); value.Exists() && !data.InterfaceThresholds[i].Node.IsNull() {
			data.InterfaceThresholds[i].Node = types.StringValue(value.String())
		} else {
			data.InterfaceThresholds[i].Node = types.StringNull()
		}
		if value := r.Get("interface"); value.Exists() && !data.InterfaceThresholds[i].Interface.IsNull() {
			data.InterfaceThresholds[i].Interface = types.StringValue(value.String())
		} else {
			data.InterfaceThresholds[i].Interface = types.StringNull()
		}
		if value := r.Get("threshold"); value.Exists() && !data.InterfaceThresholds[i].UtilizationThreshold.IsNull() {
			data.InterfaceThresholds[i].UtilizationThreshold = types.Int64Value(value.Int())
		} else {
			data.InterfaceThresholds[i].UtilizationThreshold = types.Int64Null()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *LocalCongestionMitigationInterfaceThresholds) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.InterfaceThresholds) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

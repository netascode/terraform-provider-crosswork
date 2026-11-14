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
type AssignPortsToPortGroup struct {
	Instance types.String                  `tfsdk:"instance"`
	Id       types.String                  `tfsdk:"id"`
	GroupId  types.String                  `tfsdk:"group_id"`
	Ports    []AssignPortsToPortGroupPorts `tfsdk:"ports"`
}

type AssignPortsToPortGroupPorts struct {
	Device   types.String `tfsdk:"device"`
	PortName types.String `tfsdk:"port_name"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data AssignPortsToPortGroup) getPath() string {
	return "/crosswork/grouping/v1/grouping/port"
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
func (data AssignPortsToPortGroup) toBody(ctx context.Context, state AssignPortsToPortGroup) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.GroupId.IsNull() {
		body, _ = sjson.Set(body, "groupUuid", data.GroupId.ValueString())
	}
	if len(data.Ports) > 0 {
		body, _ = sjson.Set(body, "ports", []interface{}{})
		for _, item := range data.Ports {
			itemBody := ""
			if !item.Device.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "deviceUuid", item.Device.ValueString())
			}
			if !item.PortName.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "portName", item.PortName.ValueString())
			}
			body, _ = sjson.SetRaw(body, "ports.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *AssignPortsToPortGroup) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("groupUuid"); value.Exists() {
		data.GroupId = types.StringValue(value.String())
	} else {
		data.GroupId = types.StringNull()
	}
	if value := res.Get("ports"); value.Exists() && len(value.Array()) > 0 {
		data.Ports = make([]AssignPortsToPortGroupPorts, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := AssignPortsToPortGroupPorts{}
			if cValue := v.Get("deviceUuid"); cValue.Exists() {
				item.Device = types.StringValue(cValue.String())
			} else {
				item.Device = types.StringNull()
			}
			if cValue := v.Get("portName"); cValue.Exists() {
				item.PortName = types.StringValue(cValue.String())
			} else {
				item.PortName = types.StringNull()
			}
			data.Ports = append(data.Ports, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *AssignPortsToPortGroup) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("groupUuid"); value.Exists() && !data.GroupId.IsNull() {
		data.GroupId = types.StringValue(value.String())
	} else {
		data.GroupId = types.StringNull()
	}
	for i := range data.Ports {
		keys := [...]string{"deviceUuid", "portName"}
		keyValues := [...]string{data.Ports[i].Device.ValueString(), data.Ports[i].PortName.ValueString()}

		var r gjson.Result
		res.Get("ports").ForEach(
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
		if value := r.Get("deviceUuid"); value.Exists() && !data.Ports[i].Device.IsNull() {
			data.Ports[i].Device = types.StringValue(value.String())
		} else {
			data.Ports[i].Device = types.StringNull()
		}
		if value := r.Get("portName"); value.Exists() && !data.Ports[i].PortName.IsNull() {
			data.Ports[i].PortName = types.StringValue(value.String())
		} else {
			data.Ports[i].PortName = types.StringNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *AssignPortsToPortGroup) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.Ports) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

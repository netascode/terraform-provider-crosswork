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

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type DataGatewayResourceConfiguration struct {
	Instance      types.String                                 `tfsdk:"instance"`
	Id            types.String                                 `tfsdk:"id"`
	UpdateTimeout types.Int64                                  `tfsdk:"update_timeout"`
	Collectors    []DataGatewayResourceConfigurationCollectors `tfsdk:"collectors"`
}

type DataGatewayResourceConfigurationCollectors struct {
	Name    types.String `tfsdk:"name"`
	Memory  types.Int64  `tfsdk:"memory"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data DataGatewayResourceConfiguration) getPath() string {
	return "/crosswork/dg-manager/v2/oam/container/resources"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data DataGatewayResourceConfiguration) getPathGet() string {
	return "/crosswork/dg-manager/v2/oam/container/resources/query"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Written by hand: the controller lists the containers of the collectors, which are read as the collectors of the resource.
func (data DataGatewayResourceConfiguration) readGet(c *cw.Client, path string) (cw.Res, error) {
	res, err := c.Post(path, `{}`, cw.ReadOnly)
	if err != nil {
		return res, err
	}
	return flattenCollectors(res), nil
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data DataGatewayResourceConfiguration) toBody(ctx context.Context, state DataGatewayResourceConfiguration) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.UpdateTimeout.IsNull() {
		body, _ = sjson.Set(body, "update_timeout", data.UpdateTimeout.ValueInt64())
	}
	if len(data.Collectors) > 0 {
		body, _ = sjson.Set(body, "collectors", []interface{}{})
		for _, item := range data.Collectors {
			itemBody := ""
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if !item.Memory.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "memory", item.Memory.ValueInt64())
			}
			if !item.Enabled.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "enabled", item.Enabled.ValueBool())
			}
			body, _ = sjson.SetRaw(body, "collectors.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *DataGatewayResourceConfiguration) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("data_gateway_resource_configuration")
	if value := res.Get("update_timeout"); value.Exists() {
		data.UpdateTimeout = types.Int64Value(value.Int())
	} else {
		data.UpdateTimeout = types.Int64Value(120)
	}
	if value := res.Get("collectors"); value.Exists() && len(value.Array()) > 0 {
		data.Collectors = make([]DataGatewayResourceConfigurationCollectors, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := DataGatewayResourceConfigurationCollectors{}
			if cValue := v.Get("name"); cValue.Exists() {
				item.Name = types.StringValue(cValue.String())
			} else {
				item.Name = types.StringNull()
			}
			if cValue := v.Get("memory"); cValue.Exists() {
				item.Memory = types.Int64Value(cValue.Int())
			} else {
				item.Memory = types.Int64Null()
			}
			if cValue := v.Get("enabled"); cValue.Exists() {
				item.Enabled = types.BoolValue(cValue.Bool())
			} else {
				item.Enabled = types.BoolNull()
			}
			data.Collectors = append(data.Collectors, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Written by hand: the collectors are matched by name, since the controller does not keep the order of the configuration, and
// a memory or a state that the configuration leaves out stays unset.
func (data *DataGatewayResourceConfiguration) updateFromBody(ctx context.Context, res gjson.Result) {
	for i := range data.Collectors {
		for _, item := range res.Get("collectors").Array() {
			if item.Get("name").String() != data.Collectors[i].Name.ValueString() {
				continue
			}
			if value := item.Get("memory"); value.Exists() && !data.Collectors[i].Memory.IsNull() {
				data.Collectors[i].Memory = types.Int64Value(value.Int())
			}
			if value := item.Get("enabled"); value.Exists() && !data.Collectors[i].Enabled.IsNull() {
				data.Collectors[i].Enabled = types.BoolValue(value.Bool())
			}
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *DataGatewayResourceConfiguration) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.UpdateTimeout.IsNull() {
		return false
	}
	if len(data.Collectors) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

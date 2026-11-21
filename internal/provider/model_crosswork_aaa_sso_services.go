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
type AAASSOServices struct {
	Instance types.String             `tfsdk:"instance"`
	Id       types.String             `tfsdk:"id"`
	Services []AAASSOServicesServices `tfsdk:"services"`
}

type AAASSOServicesServices struct {
	Name             types.String `tfsdk:"name"`
	EvaluationOrder  types.Int64  `tfsdk:"evaluation_order"`
	Metadata         types.String `tfsdk:"metadata"`
	MetadataLocation types.String `tfsdk:"metadata_location"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data AAASSOServices) getPath() string {
	return "/crosswork/authconfig/v1/sso/idp/services"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

func (data AAASSOServices) getPathDelete() string {
	return "/crosswork/authconfig/v1/sso/idp/services"
}

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data AAASSOServices) getPathGet() string {
	return "/crosswork/authconfig/v1/sso/idp/services"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data AAASSOServices) toBody(ctx context.Context, state AAASSOServices) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if true {
		body, _ = sjson.Set(body, "registeredServiceProvider", []interface{}{})
		for _, item := range data.Services {
			itemBody := ""
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if !item.EvaluationOrder.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "evaluationOrder", item.EvaluationOrder.ValueInt64())
			}
			if !item.Metadata.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "metadata", item.Metadata.ValueString())
			}
			if !item.MetadataLocation.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "metadataLocation", item.MetadataLocation.ValueString())
			}
			body, _ = sjson.SetRaw(body, "registeredServiceProvider.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *AAASSOServices) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("aaa_sso_services")
	if value := res.Get("registeredServiceProvider"); value.Exists() && len(value.Array()) > 0 {
		data.Services = make([]AAASSOServicesServices, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := AAASSOServicesServices{}
			if cValue := v.Get("name"); cValue.Exists() {
				item.Name = types.StringValue(cValue.String())
			} else {
				item.Name = types.StringNull()
			}
			if cValue := v.Get("evaluationOrder"); cValue.Exists() {
				item.EvaluationOrder = types.Int64Value(cValue.Int())
			} else {
				item.EvaluationOrder = types.Int64Null()
			}
			item.Metadata = types.StringValue("")
			if cValue := v.Get("metadataLocation"); cValue.Exists() {
				item.MetadataLocation = types.StringValue(cValue.String())
			} else {
				item.MetadataLocation = types.StringNull()
			}
			data.Services = append(data.Services, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *AAASSOServices) updateFromBody(ctx context.Context, res gjson.Result) {
	for i := range data.Services {
		keys := [...]string{"name"}
		keyValues := [...]string{data.Services[i].Name.ValueString()}

		var r gjson.Result
		res.Get("registeredServiceProvider").ForEach(
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
		if value := r.Get("name"); value.Exists() && !data.Services[i].Name.IsNull() {
			data.Services[i].Name = types.StringValue(value.String())
		} else {
			data.Services[i].Name = types.StringNull()
		}
		if value := r.Get("evaluationOrder"); value.Exists() && !data.Services[i].EvaluationOrder.IsNull() {
			data.Services[i].EvaluationOrder = types.Int64Value(value.Int())
		} else {
			data.Services[i].EvaluationOrder = types.Int64Null()
		}
		if value := r.Get("metadataLocation"); value.Exists() && !data.Services[i].MetadataLocation.IsNull() {
			data.Services[i].MetadataLocation = types.StringValue(value.String())
		} else {
			data.Services[i].MetadataLocation = types.StringNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *AAASSOServices) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.Services) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

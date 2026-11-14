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
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type Groupings struct {
	Instance   types.String     `tfsdk:"instance"`
	Id         types.String     `tfsdk:"id"`
	Classifier types.String     `tfsdk:"classifier"`
	Items      []GroupingsItems `tfsdk:"items"`
}

type GroupingsItems struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	ParentId         types.String `tfsdk:"parent_id"`
	Path             types.String `tfsdk:"path"`
	DiscoveryType    types.String `tfsdk:"discovery_type"`
	SubgroupsAllowed types.Bool   `tfsdk:"subgroups_allowed"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data Groupings) getPath() string {
	return "/crosswork/grouping/v1/grouping/groups"
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
func (data Groupings) toBody(ctx context.Context, state Groupings) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Classifier.IsNull() {
		body, _ = sjson.Set(body, "classifier", data.Classifier.ValueString())
	}
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
			if !item.ParentId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "parentId", item.ParentId.ValueString())
			}
			if !item.Path.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "path", item.Path.ValueString())
			}
			if !item.DiscoveryType.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "discoveryType", item.DiscoveryType.ValueString())
			}
			if !item.SubgroupsAllowed.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "subgroupsAllowed", item.SubgroupsAllowed.ValueBool())
			}
			body, _ = sjson.SetRaw(body, "items.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *Groupings) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get(""); value.Exists() {
		data.Id = types.StringValue(value.String())
	} else {
		data.Id = types.StringNull()
	}
	if value := res.Get("classifier"); value.Exists() {
		data.Classifier = types.StringValue(value.String())
	} else {
		data.Classifier = types.StringNull()
	}
	if value := res.Get("items"); value.Exists() && len(value.Array()) > 0 {
		data.Items = make([]GroupingsItems, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := GroupingsItems{}
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
			if cValue := v.Get("parentId"); cValue.Exists() {
				item.ParentId = types.StringValue(cValue.String())
			} else {
				item.ParentId = types.StringNull()
			}
			if cValue := v.Get("path"); cValue.Exists() {
				item.Path = types.StringValue(cValue.String())
			} else {
				item.Path = types.StringNull()
			}
			if cValue := v.Get("discoveryType"); cValue.Exists() {
				item.DiscoveryType = types.StringValue(cValue.String())
			} else {
				item.DiscoveryType = types.StringNull()
			}
			if cValue := v.Get("subgroupsAllowed"); cValue.Exists() {
				item.SubgroupsAllowed = types.BoolValue(cValue.Bool())
			} else {
				item.SubgroupsAllowed = types.BoolNull()
			}
			data.Items = append(data.Items, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *Groupings) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("classifier"); value.Exists() && !data.Classifier.IsNull() {
		data.Classifier = types.StringValue(value.String())
	} else {
		data.Classifier = types.StringNull()
	}
	for i := range data.Items {
		keys := [...]string{"id", "name", "parentId", "path", "discoveryType", "subgroupsAllowed"}
		keyValues := [...]string{data.Items[i].Id.ValueString(), data.Items[i].Name.ValueString(), data.Items[i].ParentId.ValueString(), data.Items[i].Path.ValueString(), data.Items[i].DiscoveryType.ValueString(), strconv.FormatBool(data.Items[i].SubgroupsAllowed.ValueBool())}

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
		if value := r.Get("parentId"); value.Exists() && !data.Items[i].ParentId.IsNull() {
			data.Items[i].ParentId = types.StringValue(value.String())
		} else {
			data.Items[i].ParentId = types.StringNull()
		}
		if value := r.Get("path"); value.Exists() && !data.Items[i].Path.IsNull() {
			data.Items[i].Path = types.StringValue(value.String())
		} else {
			data.Items[i].Path = types.StringNull()
		}
		if value := r.Get("discoveryType"); value.Exists() && !data.Items[i].DiscoveryType.IsNull() {
			data.Items[i].DiscoveryType = types.StringValue(value.String())
		} else {
			data.Items[i].DiscoveryType = types.StringNull()
		}
		if value := r.Get("subgroupsAllowed"); value.Exists() && !data.Items[i].SubgroupsAllowed.IsNull() {
			data.Items[i].SubgroupsAllowed = types.BoolValue(value.Bool())
		} else {
			data.Items[i].SubgroupsAllowed = types.BoolNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *Groupings) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Classifier.IsNull() {
		return false
	}
	if len(data.Items) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

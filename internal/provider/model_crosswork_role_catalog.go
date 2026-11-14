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
type RoleCatalog struct {
	Instance           types.String                    `tfsdk:"instance"`
	Id                 types.String                    `tfsdk:"id"`
	Apis               []RoleCatalogApis               `tfsdk:"apis"`
	Tasks              []RoleCatalogTasks              `tfsdk:"tasks"`
	TaskApiPermissions []RoleCatalogTaskApiPermissions `tfsdk:"task_api_permissions"`
}

type RoleCatalogApis struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Category types.String `tfsdk:"category"`
}

type RoleCatalogTasks struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Category types.String `tfsdk:"category"`
}

type RoleCatalogTaskApiPermissions struct {
	TaskId types.String `tfsdk:"task_id"`
	ApiId  types.String `tfsdk:"api_id"`
	Access types.String `tfsdk:"access"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data RoleCatalog) getPath() string {
	return "/crosswork/aaa/v1/api"
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
func (data RoleCatalog) toBody(ctx context.Context, state RoleCatalog) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if len(data.Apis) > 0 {
		body, _ = sjson.Set(body, "apis", []interface{}{})
		for _, item := range data.Apis {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueString())
			}
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if !item.Category.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "category", item.Category.ValueString())
			}
			body, _ = sjson.SetRaw(body, "apis.-1", itemBody)
		}
	}
	if len(data.Tasks) > 0 {
		body, _ = sjson.Set(body, "tasks", []interface{}{})
		for _, item := range data.Tasks {
			itemBody := ""
			if !item.Id.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "id", item.Id.ValueString())
			}
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if !item.Category.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "category", item.Category.ValueString())
			}
			body, _ = sjson.SetRaw(body, "tasks.-1", itemBody)
		}
	}
	if len(data.TaskApiPermissions) > 0 {
		body, _ = sjson.Set(body, "taskApiPermissions", []interface{}{})
		for _, item := range data.TaskApiPermissions {
			itemBody := ""
			if !item.TaskId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "taskId", item.TaskId.ValueString())
			}
			if !item.ApiId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "apiId", item.ApiId.ValueString())
			}
			if !item.Access.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "access", item.Access.ValueString())
			}
			body, _ = sjson.SetRaw(body, "taskApiPermissions.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *RoleCatalog) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get(""); value.Exists() {
		data.Id = types.StringValue(value.String())
	} else {
		data.Id = types.StringNull()
	}
	if value := res.Get("apis"); value.Exists() && len(value.Array()) > 0 {
		data.Apis = make([]RoleCatalogApis, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := RoleCatalogApis{}
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
			if cValue := v.Get("category"); cValue.Exists() {
				item.Category = types.StringValue(cValue.String())
			} else {
				item.Category = types.StringNull()
			}
			data.Apis = append(data.Apis, item)
			return true
		})
	}
	if value := res.Get("tasks"); value.Exists() && len(value.Array()) > 0 {
		data.Tasks = make([]RoleCatalogTasks, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := RoleCatalogTasks{}
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
			if cValue := v.Get("category"); cValue.Exists() {
				item.Category = types.StringValue(cValue.String())
			} else {
				item.Category = types.StringNull()
			}
			data.Tasks = append(data.Tasks, item)
			return true
		})
	}
	if value := res.Get("taskApiPermissions"); value.Exists() && len(value.Array()) > 0 {
		data.TaskApiPermissions = make([]RoleCatalogTaskApiPermissions, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := RoleCatalogTaskApiPermissions{}
			if cValue := v.Get("taskId"); cValue.Exists() {
				item.TaskId = types.StringValue(cValue.String())
			} else {
				item.TaskId = types.StringNull()
			}
			if cValue := v.Get("apiId"); cValue.Exists() {
				item.ApiId = types.StringValue(cValue.String())
			} else {
				item.ApiId = types.StringNull()
			}
			if cValue := v.Get("access"); cValue.Exists() {
				item.Access = types.StringValue(cValue.String())
			} else {
				item.Access = types.StringNull()
			}
			data.TaskApiPermissions = append(data.TaskApiPermissions, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *RoleCatalog) updateFromBody(ctx context.Context, res gjson.Result) {
	for i := range data.Apis {
		keys := [...]string{"id", "name", "category"}
		keyValues := [...]string{data.Apis[i].Id.ValueString(), data.Apis[i].Name.ValueString(), data.Apis[i].Category.ValueString()}

		var r gjson.Result
		res.Get("apis").ForEach(
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
		if value := r.Get("id"); value.Exists() && !data.Apis[i].Id.IsNull() {
			data.Apis[i].Id = types.StringValue(value.String())
		} else {
			data.Apis[i].Id = types.StringNull()
		}
		if value := r.Get("name"); value.Exists() && !data.Apis[i].Name.IsNull() {
			data.Apis[i].Name = types.StringValue(value.String())
		} else {
			data.Apis[i].Name = types.StringNull()
		}
		if value := r.Get("category"); value.Exists() && !data.Apis[i].Category.IsNull() {
			data.Apis[i].Category = types.StringValue(value.String())
		} else {
			data.Apis[i].Category = types.StringNull()
		}
	}
	for i := range data.Tasks {
		keys := [...]string{"id", "name", "category"}
		keyValues := [...]string{data.Tasks[i].Id.ValueString(), data.Tasks[i].Name.ValueString(), data.Tasks[i].Category.ValueString()}

		var r gjson.Result
		res.Get("tasks").ForEach(
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
		if value := r.Get("id"); value.Exists() && !data.Tasks[i].Id.IsNull() {
			data.Tasks[i].Id = types.StringValue(value.String())
		} else {
			data.Tasks[i].Id = types.StringNull()
		}
		if value := r.Get("name"); value.Exists() && !data.Tasks[i].Name.IsNull() {
			data.Tasks[i].Name = types.StringValue(value.String())
		} else {
			data.Tasks[i].Name = types.StringNull()
		}
		if value := r.Get("category"); value.Exists() && !data.Tasks[i].Category.IsNull() {
			data.Tasks[i].Category = types.StringValue(value.String())
		} else {
			data.Tasks[i].Category = types.StringNull()
		}
	}
	for i := range data.TaskApiPermissions {
		keys := [...]string{"taskId", "apiId", "access"}
		keyValues := [...]string{data.TaskApiPermissions[i].TaskId.ValueString(), data.TaskApiPermissions[i].ApiId.ValueString(), data.TaskApiPermissions[i].Access.ValueString()}

		var r gjson.Result
		res.Get("taskApiPermissions").ForEach(
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
		if value := r.Get("taskId"); value.Exists() && !data.TaskApiPermissions[i].TaskId.IsNull() {
			data.TaskApiPermissions[i].TaskId = types.StringValue(value.String())
		} else {
			data.TaskApiPermissions[i].TaskId = types.StringNull()
		}
		if value := r.Get("apiId"); value.Exists() && !data.TaskApiPermissions[i].ApiId.IsNull() {
			data.TaskApiPermissions[i].ApiId = types.StringValue(value.String())
		} else {
			data.TaskApiPermissions[i].ApiId = types.StringNull()
		}
		if value := r.Get("access"); value.Exists() && !data.TaskApiPermissions[i].Access.IsNull() {
			data.TaskApiPermissions[i].Access = types.StringValue(value.String())
		} else {
			data.TaskApiPermissions[i].Access = types.StringNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *RoleCatalog) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.Apis) > 0 {
		return false
	}
	if len(data.Tasks) > 0 {
		return false
	}
	if len(data.TaskApiPermissions) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

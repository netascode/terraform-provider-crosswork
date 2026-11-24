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
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type SystemNotificationPolicies struct {
	Instance        types.String                                `tfsdk:"instance"`
	Id              types.String                                `tfsdk:"id"`
	Name            types.String                                `tfsdk:"name"`
	Description     types.String                                `tfsdk:"description"`
	Criteria        types.String                                `tfsdk:"criteria"`
	TrapReceivers   []SystemNotificationPoliciesTrapReceivers   `tfsdk:"trap_receivers"`
	SyslogReceivers []SystemNotificationPoliciesSyslogReceivers `tfsdk:"syslog_receivers"`
}

type SystemNotificationPoliciesTrapReceivers struct {
	Name types.String `tfsdk:"name"`
}

type SystemNotificationPoliciesSyslogReceivers struct {
	Name types.String `tfsdk:"name"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data SystemNotificationPolicies) getPath() string {
	return fmt.Sprintf("/crosswork/alarm/v1/notificationpolicy/%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data SystemNotificationPolicies) getPathGet() string {
	return "/crosswork/alarm/v1/notificationpolicy"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

func (data SystemNotificationPolicies) getPathPost() string {
	return "/crosswork/alarm/v1/notificationpolicy"
}

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

func (data SystemNotificationPolicies) getPathPut() string {
	return "/crosswork/alarm/v1/notificationpolicy"
}

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data SystemNotificationPolicies) toBody(ctx context.Context, state SystemNotificationPolicies) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "groupName", data.Name.ValueString())
	}
	if !data.Description.IsNull() {
		body, _ = sjson.Set(body, "description", data.Description.ValueString())
	}
	body, _ = sjson.Set(body, "alarmType", "System/Network")
	if !data.Criteria.IsNull() {
		body, _ = sjson.Set(body, "criteria", data.Criteria.ValueString())
	}
	if len(data.TrapReceivers) > 0 {
		body, _ = sjson.Set(body, "trapDestinations", []interface{}{})
		for _, item := range data.TrapReceivers {
			itemBody := ""
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			itemBody, _ = sjson.Set(itemBody, "action", "associate")
			body, _ = sjson.SetRaw(body, "trapDestinations.-1", itemBody)
		}
	}
	if len(data.SyslogReceivers) > 0 {
		body, _ = sjson.Set(body, "syslogDestinations", []interface{}{})
		for _, item := range data.SyslogReceivers {
			itemBody := ""
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			itemBody, _ = sjson.Set(itemBody, "action", "associate")
			body, _ = sjson.SetRaw(body, "syslogDestinations.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *SystemNotificationPolicies) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("description"); value.Exists() && value.String() != "" {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	if value := res.Get("criteria"); value.Exists() {
		data.Criteria = types.StringValue(value.String())
	} else {
		data.Criteria = types.StringNull()
	}
	if value := res.Get("trapDestinations"); value.Exists() && len(value.Array()) > 0 {
		data.TrapReceivers = make([]SystemNotificationPoliciesTrapReceivers, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := SystemNotificationPoliciesTrapReceivers{}
			if cValue := v.Get("name"); cValue.Exists() {
				item.Name = types.StringValue(cValue.String())
			} else {
				item.Name = types.StringNull()
			}
			data.TrapReceivers = append(data.TrapReceivers, item)
			return true
		})
	}
	if value := res.Get("syslogDestinations"); value.Exists() && len(value.Array()) > 0 {
		data.SyslogReceivers = make([]SystemNotificationPoliciesSyslogReceivers, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := SystemNotificationPoliciesSyslogReceivers{}
			if cValue := v.Get("name"); cValue.Exists() {
				item.Name = types.StringValue(cValue.String())
			} else {
				item.Name = types.StringNull()
			}
			data.SyslogReceivers = append(data.SyslogReceivers, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *SystemNotificationPolicies) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("description"); value.Exists() && value.String() != "" && !data.Description.IsNull() {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	if value := res.Get("criteria"); value.Exists() && !data.Criteria.IsNull() {
		data.Criteria = types.StringValue(value.String())
	} else {
		data.Criteria = types.StringNull()
	}
	for i := range data.TrapReceivers {
		keys := [...]string{"name"}
		keyValues := [...]string{data.TrapReceivers[i].Name.ValueString()}

		var r gjson.Result
		res.Get("trapDestinations").ForEach(
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
		if value := r.Get("name"); value.Exists() && !data.TrapReceivers[i].Name.IsNull() {
			data.TrapReceivers[i].Name = types.StringValue(value.String())
		} else {
			data.TrapReceivers[i].Name = types.StringNull()
		}
	}
	for i := range data.SyslogReceivers {
		keys := [...]string{"name"}
		keyValues := [...]string{data.SyslogReceivers[i].Name.ValueString()}

		var r gjson.Result
		res.Get("syslogDestinations").ForEach(
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
		if value := r.Get("name"); value.Exists() && !data.SyslogReceivers[i].Name.IsNull() {
			data.SyslogReceivers[i].Name = types.StringValue(value.String())
		} else {
			data.SyslogReceivers[i].Name = types.StringNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *SystemNotificationPolicies) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Description.IsNull() {
		return false
	}
	if !data.Criteria.IsNull() {
		return false
	}
	if len(data.TrapReceivers) > 0 {
		return false
	}
	if len(data.SyslogReceivers) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

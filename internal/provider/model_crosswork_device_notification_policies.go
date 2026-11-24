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
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type DeviceNotificationPolicies struct {
	Instance        types.String                                `tfsdk:"instance"`
	Id              types.String                                `tfsdk:"id"`
	Name            types.String                                `tfsdk:"name"`
	Description     types.String                                `tfsdk:"description"`
	DeviceGroup     types.String                                `tfsdk:"device_group"`
	Filters         []DeviceNotificationPoliciesFilters         `tfsdk:"filters"`
	TrapReceivers   []DeviceNotificationPoliciesTrapReceivers   `tfsdk:"trap_receivers"`
	SyslogReceivers []DeviceNotificationPoliciesSyslogReceivers `tfsdk:"syslog_receivers"`
}

type DeviceNotificationPoliciesFilters struct {
	Severities types.Set    `tfsdk:"severities"`
	Category   types.String `tfsdk:"category"`
	Events     types.Set    `tfsdk:"events"`
}

type DeviceNotificationPoliciesTrapReceivers struct {
	Name types.String `tfsdk:"name"`
}

type DeviceNotificationPoliciesSyslogReceivers struct {
	Name types.String `tfsdk:"name"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data DeviceNotificationPolicies) getPath() string {
	return fmt.Sprintf("/crosswork/alarm-forwarding/v1/alarm-forwarding/policy/Alarms/group/%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data DeviceNotificationPolicies) getPathGet() string {
	return "/crosswork/alarm-forwarding/v1/alarm-forwarding/policy/Alarms"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

func (data DeviceNotificationPolicies) getPathPost() string {
	return "/crosswork/alarm-forwarding/v1/alarm-forwarding/policy/Alarms/group"
}

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

func (data DeviceNotificationPolicies) getPathPut() string {
	return fmt.Sprintf("/crosswork/alarm-forwarding/v1/alarm-forwarding/policy/Alarms/group/%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Written by hand: the filters and the group become the criteria (Lisp) and every receiver one rule, see device_notification_policies_logic.go.
// A value that the controller would not read back (a quote in a name) is sent as it is; the datamodel refuses it.
func (data DeviceNotificationPolicies) toBody(ctx context.Context, state DeviceNotificationPolicies) string {
	group := ""
	if !data.DeviceGroup.IsNull() && !data.DeviceGroup.IsUnknown() {
		group = data.DeviceGroup.ValueString()
	}
	criteria := deviceCriteria(deviceFiltersOf(ctx, data.Filters), group)
	var traps, syslogs []string
	for _, item := range data.TrapReceivers {
		traps = append(traps, item.Name.ValueString())
	}
	for _, item := range data.SyslogReceivers {
		syslogs = append(syslogs, item.Name.ValueString())
	}
	sort.Strings(traps)
	sort.Strings(syslogs)
	body := ""
	body, _ = sjson.Set(body, "groupName", data.Name.ValueString())
	body, _ = sjson.Set(body, "alarmType", "Device")
	if !data.Description.IsNull() && data.Description.ValueString() != "" {
		body, _ = sjson.Set(body, "description", data.Description.ValueString())
	}
	body, _ = sjson.Set(body, "policyRules", deviceRules(data.Name.ValueString(), criteria, traps, syslogs))
	return body
}

// End of section. //template:end toBody

// Written by hand: the filters and the group come from the criteria of the first rule, the receivers from the actions of all rules.
// An absent value reads as null, never as an empty set or string. When the criteria cannot be read, or the rules differ, filters read as null.
func (data *DeviceNotificationPolicies) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("description"); value.Exists() && value.String() != "" {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	// The controller takes the rules as "policyRules" and returns them as "rules".
	rules := res.Get("rules").Array()
	if len(rules) == 0 {
		rules = res.Get("policyRules").Array()
	}
	var actions []string
	same := len(rules) > 0
	for _, rule := range rules {
		for _, action := range rule.Get("actions").Array() {
			actions = append(actions, action.String())
		}
		if rule.Get("criteria").String() != rules[0].Get("criteria").String() {
			same = false
		}
	}
	data.DeviceGroup = types.StringNull()
	data.Filters = nil
	if len(rules) > 0 {
		filters, group, ok := parseDeviceCriteria(rules[0].Get("criteria").String())
		if group != "" {
			data.DeviceGroup = types.StringValue(group)
		}
		if ok && same {
			data.Filters = deviceFiltersToModel(filters)
		}
	}
	traps, syslogs := receiversFromActions(actions)
	data.TrapReceivers = nil
	for _, name := range traps {
		data.TrapReceivers = append(data.TrapReceivers, DeviceNotificationPoliciesTrapReceivers{Name: types.StringValue(name)})
	}
	data.SyslogReceivers = nil
	for _, name := range syslogs {
		data.SyslogReceivers = append(data.SyslogReceivers, DeviceNotificationPoliciesSyslogReceivers{Name: types.StringValue(name)})
	}
}

// End of section. //template:end fromBody

// Written by hand: a refresh reads every attribute, so that a change made in the UI shows as a difference. The filters are a list, and
// the UI reorders them, so the order of the state is kept when the controller holds the same filters.
func (data *DeviceNotificationPolicies) updateFromBody(ctx context.Context, res gjson.Result) {
	before := data.Filters
	data.fromBody(ctx, res)
	if data.Filters != nil && sameDeviceFilters(deviceFiltersOf(ctx, before), deviceFiltersOf(ctx, data.Filters)) {
		data.Filters = before
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *DeviceNotificationPolicies) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Description.IsNull() {
		return false
	}
	if !data.DeviceGroup.IsNull() {
		return false
	}
	if len(data.Filters) > 0 {
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

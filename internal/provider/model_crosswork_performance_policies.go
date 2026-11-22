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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type PerformancePolicies struct {
	Instance        types.String                             `tfsdk:"instance"`
	Id              types.String                             `tfsdk:"id"`
	Name            types.String                             `tfsdk:"name"`
	PolicyTemplate  types.String                             `tfsdk:"policy_template"`
	Description     types.String                             `tfsdk:"description"`
	Contact         types.String                             `tfsdk:"contact"`
	Active          types.Bool                               `tfsdk:"active"`
	Devices         types.Set                                `tfsdk:"devices"`
	DeviceGroups    types.Set                                `tfsdk:"device_groups"`
	PortGroups      types.Set                                `tfsdk:"port_groups"`
	SchemaIntervals []PerformancePoliciesSchemaIntervals     `tfsdk:"schema_intervals"`
	Thresholds      map[string]PerformancePoliciesThresholds `tfsdk:"thresholds"`
}

type PerformancePoliciesSchemaIntervals struct {
	Schema  types.String `tfsdk:"schema"`
	Seconds types.Int64  `tfsdk:"seconds"`
}

type PerformancePoliciesThresholds struct {
	Conditions []PerformancePoliciesThresholdsConditions `tfsdk:"conditions"`
}

type PerformancePoliciesThresholdsConditions struct {
	Relation types.String `tfsdk:"relation"`
	Value    types.String `tfsdk:"value"`
	Count    types.Int64  `tfsdk:"count"`
	Severity types.String `tfsdk:"severity"`
	MatchCos types.Set    `tfsdk:"match_cos"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data PerformancePolicies) getPath() string {
	return "/crosswork/performance/v1/policies"
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

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Written by hand: the targets are comma-joined strings, the intervals are a map by schema and the thresholds a map of lists with
// numbered conditions. The devices of the plan are already uuids here: the resource turns the names into uuids on a copy.
func (data PerformancePolicies) toBody(ctx context.Context, state PerformancePolicies) string {
	body := ""
	if !state.Id.IsNull() && state.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "id", state.Id.ValueString())
	}
	body, _ = sjson.Set(body, "name", data.Name.ValueString())
	body, _ = sjson.Set(body, "policyTemplate", data.PolicyTemplate.ValueString())
	body, _ = sjson.Set(body, "description", data.Description.ValueString())
	body, _ = sjson.Set(body, "tag", tagFromContact(data.Contact))
	body, _ = sjson.Set(body, "active", data.Active.ValueBool())
	body, _ = sjson.Set(body, "devices", strings.Join(setStrings(ctx, data.Devices), ","))
	body, _ = sjson.Set(body, "deviceGroups", strings.Join(setStrings(ctx, data.DeviceGroups), ","))
	body, _ = sjson.Set(body, "portGroups", strings.Join(setStrings(ctx, data.PortGroups), ","))
	body, _ = sjson.SetRaw(body, "schemasInterval", "{}")
	for _, item := range data.SchemaIntervals {
		body, _ = sjson.Set(body, "schemasInterval."+escapeJSONKey(item.Schema.ValueString()), item.Seconds.ValueInt64())
	}
	body, _ = sjson.SetRaw(body, "thresholds", "{}")
	for metric, threshold := range data.Thresholds {
		conditions := "[]"
		for i, condition := range threshold.Conditions {
			conditions, _ = sjson.SetRaw(conditions, "-1", conditionBody(ctx, i+1, condition))
		}
		body, _ = sjson.SetRaw(body, "thresholds."+escapeJSONKey(metric), conditions)
	}
	return body
}

// End of section. //template:end toBody

// Written by hand: see toBody. The answer of a create, a read of one policy and an item of the list all hold the policy in
// "monitoringPolicy".
func (data *PerformancePolicies) fromBody(ctx context.Context, res gjson.Result) {
	policy := res.Get("monitoringPolicy")
	if !policy.Exists() {
		policy = res
	}
	if value := policy.Get("id"); value.Exists() {
		data.Id = types.StringValue(value.String())
	}
	data.Name = types.StringValue(policy.Get("name").String())
	data.PolicyTemplate = types.StringValue(policy.Get("policyTemplate").String())
	if value := policy.Get("description"); value.Exists() && value.String() != "" {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	data.Contact = contactFromTag(policy.Get("tag").String())
	data.Active = types.BoolValue(policy.Get("active").Bool())
	data.Devices = targetSet(splitTargets(policy.Get("devices").String()))
	data.DeviceGroups = targetSet(splitTargets(policy.Get("deviceGroups").String()))
	data.PortGroups = targetSet(splitTargets(policy.Get("portGroups").String()))
	data.SchemaIntervals = nil
	policy.Get("schemasInterval").ForEach(func(key, value gjson.Result) bool {
		data.SchemaIntervals = append(data.SchemaIntervals, PerformancePoliciesSchemaIntervals{
			Schema:  types.StringValue(key.String()),
			Seconds: types.Int64Value(value.Int()),
		})
		return true
	})
	data.Thresholds = nil
	policy.Get("thresholds").ForEach(func(key, value gjson.Result) bool {
		threshold := PerformancePoliciesThresholds{Conditions: []PerformancePoliciesThresholdsConditions{}}
		for _, condition := range value.Array() {
			threshold.Conditions = append(threshold.Conditions, conditionFromBody(condition))
		}
		if data.Thresholds == nil {
			data.Thresholds = map[string]PerformancePoliciesThresholds{}
		}
		data.Thresholds[key.String()] = threshold
		return true
	})
}

// End of section. //template:end fromBody

// Written by hand: a refresh reads every attribute, so that a change made in the UI shows as a difference. A target set or a
// thresholds map that the state holds empty stays empty when the controller holds nothing, as the plan has it.
func (data *PerformancePolicies) updateFromBody(ctx context.Context, res gjson.Result) {
	keepDevices := !data.Devices.IsNull() && len(data.Devices.Elements()) == 0
	keepDeviceGroups := !data.DeviceGroups.IsNull() && len(data.DeviceGroups.Elements()) == 0
	keepPortGroups := !data.PortGroups.IsNull() && len(data.PortGroups.Elements()) == 0
	keepThresholds := data.Thresholds != nil && len(data.Thresholds) == 0
	data.fromBody(ctx, res)
	empty, _ := types.SetValueFrom(ctx, types.StringType, []string{})
	if keepDevices && data.Devices.IsNull() {
		data.Devices = empty
	}
	if keepDeviceGroups && data.DeviceGroups.IsNull() {
		data.DeviceGroups = empty
	}
	if keepPortGroups && data.PortGroups.IsNull() {
		data.PortGroups = empty
	}
	if keepThresholds && data.Thresholds == nil {
		data.Thresholds = map[string]PerformancePoliciesThresholds{}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *PerformancePolicies) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Name.IsNull() {
		return false
	}
	if !data.PolicyTemplate.IsNull() {
		return false
	}
	if !data.Description.IsNull() {
		return false
	}
	if !data.Contact.IsNull() {
		return false
	}
	if !data.Active.IsNull() {
		return false
	}
	if !data.Devices.IsNull() {
		return false
	}
	if !data.DeviceGroups.IsNull() {
		return false
	}
	if !data.PortGroups.IsNull() {
		return false
	}
	if len(data.SchemaIntervals) > 0 {
		return false
	}
	if len(data.Thresholds) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

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
type PortGroupsRules struct {
	Instance   types.String                `tfsdk:"instance"`
	Id         types.String                `tfsdk:"id"`
	GroupId    types.String                `tfsdk:"group_id"`
	Name       types.String                `tfsdk:"name"`
	Conditions []PortGroupsRulesConditions `tfsdk:"conditions"`
}

type PortGroupsRulesConditions struct {
	Attribute       types.String `tfsdk:"attribute"`
	Operator        types.String `tfsdk:"operator"`
	Value           types.String `tfsdk:"value"`
	LogicalOperator types.String `tfsdk:"logical_operator"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data PortGroupsRules) getPath() string {
	return "/crosswork/grouping/v1/grouping/rule"
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

// Written by hand: Crosswork wants the conditions as one text, and the update names the rule.
func (data PortGroupsRules) toBody(ctx context.Context, state PortGroupsRules) string {
	body := ""
	body, _ = sjson.Set(body, "classifier", "UserDefinedPorts")
	body, _ = sjson.Set(body, "active", true)
	body, _ = sjson.Set(body, "ordering", 0)
	if state.Id.ValueString() != "" {
		body, _ = sjson.Set(body, "uuid", state.Id.ValueString())
	}
	body, _ = sjson.Set(body, "targetGroupUuid", data.GroupId.ValueString())
	body, _ = sjson.Set(body, "name", data.Name.ValueString())
	body, _ = sjson.Set(body, "conditions", ruleConditionsText(data.Conditions))
	return body
}

// End of section. //template:end toBody

// Written by hand: the conditions come back as one text.
func (data *PortGroupsRules) fromBody(ctx context.Context, res gjson.Result) {
	data.GroupId = types.StringNull()
	if value := res.Get("targetGroupUuid"); value.Exists() {
		data.GroupId = types.StringValue(value.String())
	}
	data.Name = types.StringNull()
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	}
	data.Conditions = parseRuleConditions(res.Get("conditions").String())
}

// End of section. //template:end fromBody

// Written by hand: the conditions are the whole list that Crosswork holds, so any added or changed row is a difference.
func (data *PortGroupsRules) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("targetGroupUuid"); value.Exists() && !data.GroupId.IsNull() {
		data.GroupId = types.StringValue(value.String())
	}
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	}
	data.Conditions = parseRuleConditions(res.Get("conditions").String())
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *PortGroupsRules) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.GroupId.IsNull() {
		return false
	}
	if !data.Name.IsNull() {
		return false
	}
	if len(data.Conditions) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

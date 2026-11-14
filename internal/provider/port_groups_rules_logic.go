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

// Written by hand: Crosswork stores the conditions as one text (a JSON document inside a string) and reads it back
// reformatted, and it accepts a type that does not exist or a speed in a wrong form and then fails to list the ports of
// the group, so the checks live here.

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

const portRuleConditionsPath = "/crosswork/grouping/v1/grouping/ports/rule/conditions"

var stringOperators = map[string]string{
	"equals": "SO_Equals", "not_equals": "SO_NotEquals", "matches": "SO_Matches", "not_matches": "SO_NotMatches",
	"contains": "SO_Contains", "not_contains": "SO_NotContains", "starts_with": "SO_StartWith", "ends_with": "SO_EndWith",
	"blank": "SO_Blank", "not_blank": "SO_NotBlank",
}

var numberOperators = map[string]string{
	"equals": "NO_Equals", "not_equals": "NO_NotEquals", "greater_than": "NO_GreaterThan",
	"greater_than_or_equals": "NO_GreaterThanOrEquals", "less_than": "NO_LessThan", "less_than_or_equals": "NO_LessThanOrEquals",
}

var speedValue = regexp.MustCompile(`^[0-9]+ (bps|Kbps|Mbps|Gbps|Tbps)$`)

func operatorsOf(attribute string) map[string]string {
	switch attribute {
	case "speed":
		return numberOperators
	case "name":
		return stringOperators
	default:
		return map[string]string{"equals": "SO_Equals", "not_equals": "SO_NotEquals"}
	}
}

func hasNoValue(operator string) bool { return operator == "blank" || operator == "not_blank" }

func validateRuleConditions(conds []PortGroupsRulesConditions) []string {
	var problems []string
	for i, c := range conds {
		at := fmt.Sprintf("condition %d", i+1)
		if c.Attribute.IsUnknown() || c.Operator.IsUnknown() || c.Value.IsUnknown() || c.LogicalOperator.IsUnknown() {
			continue
		}
		if i == 0 && !c.LogicalOperator.IsNull() {
			problems = append(problems, at+": the first condition has no logical_operator")
		}
		if i > 0 && c.LogicalOperator.IsNull() {
			problems = append(problems, at+": logical_operator is required after the first condition")
		}
		attribute, operator, value := c.Attribute.ValueString(), c.Operator.ValueString(), c.Value.ValueString()
		if _, ok := operatorsOf(attribute)[operator]; !ok {
			problems = append(problems, fmt.Sprintf("%s: the operator %q is not valid for %s", at, operator, attribute))
			continue
		}
		if hasNoValue(operator) {
			if !c.Value.IsNull() {
				problems = append(problems, fmt.Sprintf("%s: %s takes no value", at, operator))
			}
			continue
		}
		if c.Value.IsNull() || value == "" {
			problems = append(problems, fmt.Sprintf("%s: %s needs a value", at, operator))
			continue
		}
		switch attribute {
		case "admin_status", "oper_status":
			if value != "UP" && value != "DOWN" {
				problems = append(problems, fmt.Sprintf("%s: the value of %s is UP or DOWN, not %q", at, attribute, value))
			}
		case "speed":
			if !speedValue.MatchString(value) {
				problems = append(problems, fmt.Sprintf("%s: the value of speed is a whole number, a space and a unit (bps, Kbps, Mbps, Gbps or Tbps), as \"100 Gbps\", not %q", at, value))
			}
		}
	}
	return problems
}

// checkRuleTypes asks Crosswork which port types exist: a type that is not listed is accepted by Crosswork, matches nothing
// and makes the group fail to list its ports.
func checkRuleTypes(c *cw.Client, conds []PortGroupsRulesConditions) error {
	var wanted []string
	for _, cond := range conds {
		if cond.Attribute.ValueString() == "type" && !cond.Value.IsNull() {
			wanted = append(wanted, cond.Value.ValueString())
		}
	}
	if len(wanted) == 0 {
		return nil
	}
	res, err := c.Get(portRuleConditionsPath, cw.ReadOnly)
	if err != nil {
		return fmt.Errorf("reading the port types of Crosswork: %w", err)
	}
	known := map[string]bool{}
	for _, cond := range res.Get("conditions").Array() {
		if cond.Get("attributeName").String() == "type" {
			for _, v := range cond.Get("operators.0.preDefinedValues").Array() {
				known[v.String()] = true
			}
		}
	}
	var unknown []string
	for _, w := range wanted {
		if !known[w] {
			unknown = append(unknown, w)
		}
	}
	if len(unknown) > 0 {
		sort.Strings(unknown)
		return fmt.Errorf("Crosswork does not list the port type(s) %q; it would accept them and then fail to list the ports of the group", unknown)
	}
	return nil
}

func ruleConditionsText(conds []PortGroupsRulesConditions) string {
	rows := make([]map[string]any, 0, len(conds))
	for i, c := range conds {
		attribute, operator := c.Attribute.ValueString(), c.Operator.ValueString()
		row := map[string]any{"order": i + 1, "attributeName": attribute}
		if !c.Value.IsNull() {
			row["value"] = c.Value.ValueString()
		}
		if !c.LogicalOperator.IsNull() {
			row["logicalOperator"] = "LO_" + strings.ToUpper(c.LogicalOperator.ValueString()[:1]) + c.LogicalOperator.ValueString()[1:]
		}
		kind := "stringCondition"
		if attribute == "speed" {
			kind = "numericCondition"
		}
		row[kind] = map[string]string{"operator": operatorsOf(attribute)[operator]}
		rows = append(rows, row)
	}
	text, _ := json.Marshal(map[string]any{"conditions": rows})
	return string(text)
}

func parseRuleConditions(text string) []PortGroupsRulesConditions {
	var out []PortGroupsRulesConditions
	gjson.Parse(text).Get("conditions").ForEach(func(_, row gjson.Result) bool {
		attribute := row.Get("attributeName").String()
		cond := PortGroupsRulesConditions{Attribute: types.StringValue(attribute), Value: types.StringNull(), LogicalOperator: types.StringNull()}
		operator := row.Get("stringCondition.operator")
		if !operator.Exists() {
			operator = row.Get("numericCondition.operator")
		}
		for name, controller := range operatorsOf(attribute) {
			if controller == operator.String() {
				cond.Operator = types.StringValue(name)
			}
		}
		if cond.Operator.IsNull() || cond.Operator.IsUnknown() {
			cond.Operator = types.StringValue(operator.String())
		}
		if v := row.Get("value"); v.Exists() {
			cond.Value = types.StringValue(v.String())
		}
		if v := row.Get("logicalOperator"); v.Exists() {
			cond.LogicalOperator = types.StringValue(strings.ToLower(strings.TrimPrefix(v.String(), "LO_")))
		}
		out = append(out, cond)
		return true
	})
	return out
}

func checkRule(c *cw.Client, conds []PortGroupsRulesConditions, diags *diag.Diagnostics) bool {
	if problems := validateRuleConditions(conds); len(problems) > 0 {
		diags.AddError("Invalid conditions", strings.Join(problems, "; "))
		return false
	}
	if err := checkRuleTypes(c, conds); err != nil {
		diags.AddError("Invalid conditions", err.Error())
		return false
	}
	return true
}

// The conditions are unknown while Terraform still builds them from other resources; there is then nothing to check.
func (r *PortGroupsRulesResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var list types.List
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("conditions"), &list)...)
	if resp.Diagnostics.HasError() || list.IsNull() || list.IsUnknown() {
		return
	}
	for _, element := range list.Elements() {
		if element.IsUnknown() {
			return
		}
	}
	var conds []PortGroupsRulesConditions
	resp.Diagnostics.Append(list.ElementsAs(ctx, &conds, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if problems := validateRuleConditions(conds); len(problems) > 0 {
		resp.Diagnostics.AddError("Invalid conditions", strings.Join(problems, "; "))
	}
}

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
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

func cond(attribute, operator, value, logical string) PortGroupsRulesConditions {
	c := PortGroupsRulesConditions{Attribute: types.StringValue(attribute), Operator: types.StringValue(operator), Value: types.StringNull(), LogicalOperator: types.StringNull()}
	if value != "" {
		c.Value = types.StringValue(value)
	}
	if logical != "" {
		c.LogicalOperator = types.StringValue(logical)
	}
	return c
}

// a rule with one condition that has no value, and as Crosswork returns it, reformatted
var uiRule = []PortGroupsRulesConditions{
	cond("name", "matches", "Four*", ""), cond("admin_status", "equals", "UP", "and"), cond("type", "equals", "Ethernet CSMACD", "or"),
	cond("oper_status", "not_equals", "DOWN", "and"), cond("speed", "greater_than_or_equals", "100 Gbps", "and"), cond("name", "not_blank", "", "or"),
}

const controllerText = `{
  "conditions": [{
    "attributeName": "name", "value": "Four*", "order": 1, "stringCondition": {"operator": "SO_Matches"}
  }, {
    "attributeName": "admin_status", "value": "UP", "order": 2, "logicalOperator": "LO_And", "stringCondition": {"operator": "SO_Equals"}
  }, {
    "attributeName": "type", "value": "Ethernet CSMACD", "order": 3, "logicalOperator": "LO_Or", "stringCondition": {"operator": "SO_Equals"}
  }, {
    "attributeName": "oper_status", "value": "DOWN", "order": 4, "logicalOperator": "LO_And", "stringCondition": {"operator": "SO_NotEquals"}
  }, {
    "attributeName": "speed", "value": "100 Gbps", "order": 5, "logicalOperator": "LO_And", "numericCondition": {"operator": "NO_GreaterThanOrEquals"}
  }, {
    "attributeName": "name", "order": 6, "logicalOperator": "LO_Or", "stringCondition": {"operator": "SO_NotBlank"}
  }]
}`

func TestPortRuleConditionsText(t *testing.T) {
	rows := gjson.Parse(ruleConditionsText(uiRule)).Get("conditions").Array()
	if len(rows) != 6 {
		t.Fatalf("got %d rows, want 6", len(rows))
	}
	check := func(ok bool, what string) {
		t.Helper()
		if !ok {
			t.Error(what)
		}
	}
	check(rows[0].Get("stringCondition.operator").String() == "SO_Matches", "matches is SO_Matches")
	check(!rows[0].Get("logicalOperator").Exists(), "the first condition has no connector")
	check(rows[1].Get("logicalOperator").String() == "LO_And", "and is LO_And")
	check(rows[2].Get("logicalOperator").String() == "LO_Or", "or is LO_Or")
	check(rows[4].Get("numericCondition.operator").String() == "NO_GreaterThanOrEquals", "speed is numeric")
	check(!rows[5].Get("value").Exists(), "not_blank sends no value")
	check(rows[5].Get("order").Int() == 6, "the order counts from 1")
}

func TestPortRuleConditionsParse(t *testing.T) {
	if got := parseRuleConditions(controllerText); !reflect.DeepEqual(got, uiRule) {
		t.Errorf("the text of Crosswork reads as %v", got)
	}
	if got := parseRuleConditions(ruleConditionsText(uiRule)); !reflect.DeepEqual(got, uiRule) {
		t.Errorf("what is sent reads back as %v", got)
	}
	if got := parseRuleConditions(""); len(got) != 0 {
		t.Errorf("an empty text has conditions: %v", got)
	}
}

func TestPortRuleValidate(t *testing.T) {
	if got := validateRuleConditions(uiRule); len(got) != 0 {
		t.Errorf("a valid rule is refused: %v", got)
	}
	bad := map[string][]PortGroupsRulesConditions{
		"the first condition has no logical_operator": {cond("name", "equals", "a", "and")},
		"logical_operator is required":                {cond("name", "equals", "a", ""), cond("name", "equals", "b", "")},
		"is not valid for speed":                      {cond("speed", "contains", "1 Gbps", "")},
		"is not valid for type":                       {cond("type", "starts_with", "x", "")},
		"takes no value":                              {cond("name", "blank", "x", "")},
		"needs a value":                               {cond("name", "starts_with", "", "")},
		"is UP or DOWN":                               {cond("admin_status", "equals", "up", "")},
		"a whole number, a space and a unit":          {cond("speed", "equals", "100Gbps", "")},
	}
	for want, conds := range bad {
		if got := strings.Join(validateRuleConditions(conds), ";"); !strings.Contains(got, want) {
			t.Errorf("want %q in %q", want, got)
		}
	}
	for _, v := range []string{"1.5 Gbps", "100", "100 G", "100 gbps", " 100 Gbps"} {
		if len(validateRuleConditions([]PortGroupsRulesConditions{cond("speed", "equals", v, "")})) == 0 {
			t.Errorf("speed %q is accepted", v)
		}
	}
	for _, v := range []string{"500 bps", "1 Kbps", "10 Mbps", "100 Gbps", "1 Tbps"} {
		if got := validateRuleConditions([]PortGroupsRulesConditions{cond("speed", "equals", v, "")}); len(got) != 0 {
			t.Errorf("speed %q is refused: %v", v, got)
		}
	}
}

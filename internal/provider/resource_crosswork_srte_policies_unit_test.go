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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

func TestSRTEPoliciesBodyAndReadBack(t *testing.T) {
	ctx := context.Background()
	colors, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("red"), types.StringValue("blue")})
	const root = "cisco-sr-te-cfp-sr-policies:policy.0."
	model := SRTEPolicies{
		Name: types.StringValue("cwtest-pol"), HeadEnd: types.StringValue("cwtest-device-1"), TailEnd: types.StringValue("192.0.2.61"), Color: types.Int64Value(999),
		AutoRouteMetricRelative:  types.Int64Value(-5),
		AutoRouteIncludePrefixes: []SRTEPoliciesAutoRouteIncludePrefixes{{PrefixAddress: types.StringValue("10.250.0.0/16")}},
		Paths: []SRTEPoliciesPaths{
			{
				Preference:            types.Int64Value(300),
				ExplicitSidLists:      []SRTEPoliciesPathsExplicitSidLists{{Name: types.StringValue("sl-a"), Weight: types.Int64Value(10)}, {Name: types.StringValue("sl-b")}},
				ExplicitAffinityRules: []SRTEPoliciesPathsExplicitAffinityRules{{Action: types.StringValue("include-any"), Colors: colors}},
			},
			{Preference: types.Int64Value(100), DynamicPce: types.BoolValue(true), DynamicMetricType: types.StringValue("igp"), DynamicSidLimit: types.Int64Value(6)},
		},
	}
	body := model.toBody(ctx, SRTEPolicies{})
	for path, want := range map[string]string{
		root + "head-end.0.name": "cwtest-device-1",
		root + "head-end.#":      "1",
		root + "auto-route.auto-route-metric.metric-relative-value":          "-5",
		root + "auto-route.include-prefixes.include-prefix.0.prefix-address": "10.250.0.0/16",
		root + "path.0.explicit.sid-list.0.weight":                           "10",
		root + "path.0.explicit.constraints.affinity.rule.0.color.#":         "2",
		root + "path.1.dynamic.pce":                                          "{}",
		root + "path.1.dynamic.constraints.sid-limit":                        "6",
	} {
		if got := gjson.Get(body, path).String(); got != want {
			t.Errorf("%s = %q, want %q", path, got, want)
		}
	}
	for _, path := range []string{root + "path.0.explicit.sid-list.1.weight", root + "path.0.dynamic", root + "path.1.explicit", root + "srv6", root + "performance-measurement", root + "auto-route.auto-route-metric.metric-constant-value"} {
		if gjson.Get(body, path).Exists() {
			t.Errorf("%s was sent although it is not set", path)
		}
	}

	// NSO returns the paths sorted by preference and the colors sorted; the state keeps the configured order
	resp := gjson.Parse(`{"cisco-sr-te-cfp-sr-policies:policy":[{"name":"cwtest-pol","head-end":[{"name":"cwtest-device-1"}],"tail-end":"192.0.2.61","color":999,
	  "auto-route":{"auto-route-metric":{"metric-relative-value":-5},"include-prefixes":{"include-prefix":[{"prefix-address":"10.250.0.0/16"}]}},
	  "path":[{"preference":100,"dynamic":{"pce":{},"metric-type":"igp","constraints":{"sid-limit":6}}},
	          {"preference":300,"explicit":{"sid-list":[{"name":"sl-a","weight":10},{"name":"sl-b"}],"constraints":{"affinity":{"rule":[{"action":"include-any","color":["blue","red"]}]}}}}]}]}`)
	state := model
	state.updateFromBody(ctx, resp)
	if state.Paths[0].Preference.ValueInt64() != 300 || len(state.Paths[0].ExplicitSidLists) != 2 || state.Paths[0].ExplicitSidLists[0].Weight.ValueInt64() != 10 || !state.Paths[0].ExplicitSidLists[1].Weight.IsNull() {
		t.Errorf("the explicit path was not matched by its preference: %+v", state.Paths[0])
	}
	if !state.Paths[1].DynamicPce.ValueBool() || state.Paths[1].DynamicMetricType.ValueString() != "igp" || state.Paths[1].DynamicSidLimit.ValueInt64() != 6 {
		t.Errorf("the dynamic path was not matched by its preference: %+v", state.Paths[1])
	}
	if !state.Paths[0].ExplicitAffinityRules[0].Colors.Equal(colors) {
		t.Errorf("the colors differ: %v", state.Paths[0].ExplicitAffinityRules[0].Colors)
	}
	if state.HeadEnd.ValueString() != "cwtest-device-1" || state.AutoRouteMetricRelative.ValueInt64() != -5 || !state.AutoRouteMetricConstant.IsNull() {
		t.Errorf("the top level was not read: %+v", state)
	}
	var imported SRTEPolicies
	imported.fromBody(ctx, resp)
	if imported.HeadEnd.ValueString() != "cwtest-device-1" || len(imported.Paths) != 2 || len(imported.AutoRouteIncludePrefixes) != 1 {
		t.Errorf("fromBody (import) read it wrong: %+v", imported)
	}
}

func TestSRTESIDListsBodyAndReadBack(t *testing.T) {
	ctx := context.Background()
	model := SRTESIDLists{Name: types.StringValue("cwtest-sl"), Sids: []SRTESIDListsSids{
		{Index: types.Int64Value(1), Label: types.Int64Value(16001)},
		{Index: types.Int64Value(2), Address: types.StringValue("192.0.2.10")},
	}}
	body := model.toBody(ctx, SRTESIDLists{})
	const root = "cisco-sr-te-cfp-sr-policies:sid-list.0."
	if gjson.Get(body, root+"sid.0.mpls.label").Int() != 16001 || gjson.Get(body, root+"sid.1.ipv4.address").String() != "192.0.2.10" {
		t.Errorf("the entries are wrong: %s", body)
	}
	if gjson.Get(body, root+"sid.0.ipv4").Exists() || gjson.Get(body, root+"sid.1.mpls").Exists() {
		t.Errorf("the unset side of an entry was sent: %s", body)
	}
	resp := gjson.Parse(`{"cisco-sr-te-cfp-sr-policies:sid-list":[{"name":"cwtest-sl","sid":[{"index":1,"mpls":{"label":16001}},{"index":2,"ipv4":{"address":"192.0.2.10"}}]}]}`)
	state := model
	state.updateFromBody(ctx, resp)
	if state.Sids[0].Label.ValueInt64() != 16001 || !state.Sids[0].Address.IsNull() || state.Sids[1].Address.ValueString() != "192.0.2.10" || !state.Sids[1].Label.IsNull() {
		t.Errorf("the entries were read wrong: %+v", state.Sids)
	}
}

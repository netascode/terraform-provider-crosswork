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
type L3VPNRoutePolicies struct {
	Instance   types.String                   `tfsdk:"instance"`
	Id         types.String                   `tfsdk:"id"`
	Name       types.String                   `tfsdk:"name"`
	Statements []L3VPNRoutePoliciesStatements `tfsdk:"statements"`
}

type L3VPNRoutePoliciesStatements struct {
	Name                 types.String `tfsdk:"name"`
	MatchSourcePrefixSet types.String `tfsdk:"match_source_prefix_set"`
	MatchDestPrefixSet   types.String `tfsdk:"match_dest_prefix_set"`
	PolicyResult         types.String `tfsdk:"policy_result"`
	SetExtCommunity      types.String `tfsdk:"set_ext_community"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data L3VPNRoutePolicies) getPath() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/cisco-l3vpn-routing-policy:l3vpn-routing-policy/policy-definitions/policy-definition=%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

func (data L3VPNRoutePolicies) getPathPost() string {
	return "/crosswork/proxy/nso/restconf/data/cisco-l3vpn-routing-policy:l3vpn-routing-policy/policy-definitions"
}

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
func (data L3VPNRoutePolicies) toBody(ctx context.Context, state L3VPNRoutePolicies) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "cisco-l3vpn-routing-policy:policy-definition.0.name", data.Name.ValueString())
	}
	if len(data.Statements) > 0 {
		body, _ = sjson.Set(body, "cisco-l3vpn-routing-policy:policy-definition.0.statements.statement", []interface{}{})
		for _, item := range data.Statements {
			itemBody := ""
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if !item.MatchSourcePrefixSet.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "conditions.match-source-prefix-set.prefix-set", item.MatchSourcePrefixSet.ValueString())
			}
			if !item.MatchDestPrefixSet.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "conditions.match-dest-prefix-set.prefix-set", item.MatchDestPrefixSet.ValueString())
			}
			if !item.PolicyResult.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "actions.policy-result", item.PolicyResult.ValueString())
			}
			if !item.SetExtCommunity.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "actions.bgp-actions.set-ext-community.ext-community-set-ref", item.SetExtCommunity.ValueString())
			}
			body, _ = sjson.SetRaw(body, "cisco-l3vpn-routing-policy:policy-definition.0.statements.statement.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *L3VPNRoutePolicies) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-l3vpn-routing-policy:policy-definition.0.statements.statement"); value.Exists() && len(value.Array()) > 0 {
		data.Statements = make([]L3VPNRoutePoliciesStatements, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := L3VPNRoutePoliciesStatements{}
			if cValue := v.Get("name"); cValue.Exists() {
				item.Name = types.StringValue(cValue.String())
			} else {
				item.Name = types.StringNull()
			}
			if cValue := v.Get("conditions.match-source-prefix-set.prefix-set"); cValue.Exists() {
				item.MatchSourcePrefixSet = types.StringValue(cValue.String())
			} else {
				item.MatchSourcePrefixSet = types.StringNull()
			}
			if cValue := v.Get("conditions.match-dest-prefix-set.prefix-set"); cValue.Exists() {
				item.MatchDestPrefixSet = types.StringValue(cValue.String())
			} else {
				item.MatchDestPrefixSet = types.StringNull()
			}
			if cValue := v.Get("actions.policy-result"); cValue.Exists() {
				item.PolicyResult = types.StringValue(cValue.String())
			} else {
				item.PolicyResult = types.StringNull()
			}
			if cValue := v.Get("actions.bgp-actions.set-ext-community.ext-community-set-ref"); cValue.Exists() {
				item.SetExtCommunity = types.StringValue(cValue.String())
			} else {
				item.SetExtCommunity = types.StringNull()
			}
			data.Statements = append(data.Statements, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *L3VPNRoutePolicies) updateFromBody(ctx context.Context, res gjson.Result) {
	for i := range data.Statements {
		keys := [...]string{"name"}
		keyValues := [...]string{data.Statements[i].Name.ValueString()}

		var r gjson.Result
		res.Get("cisco-l3vpn-routing-policy:policy-definition.0.statements.statement").ForEach(
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
		if value := r.Get("name"); value.Exists() && !data.Statements[i].Name.IsNull() {
			data.Statements[i].Name = types.StringValue(value.String())
		} else {
			data.Statements[i].Name = types.StringNull()
		}
		if value := r.Get("conditions.match-source-prefix-set.prefix-set"); value.Exists() && !data.Statements[i].MatchSourcePrefixSet.IsNull() {
			data.Statements[i].MatchSourcePrefixSet = types.StringValue(value.String())
		} else {
			data.Statements[i].MatchSourcePrefixSet = types.StringNull()
		}
		if value := r.Get("conditions.match-dest-prefix-set.prefix-set"); value.Exists() && !data.Statements[i].MatchDestPrefixSet.IsNull() {
			data.Statements[i].MatchDestPrefixSet = types.StringValue(value.String())
		} else {
			data.Statements[i].MatchDestPrefixSet = types.StringNull()
		}
		if value := r.Get("actions.policy-result"); value.Exists() && !data.Statements[i].PolicyResult.IsNull() {
			data.Statements[i].PolicyResult = types.StringValue(value.String())
		} else {
			data.Statements[i].PolicyResult = types.StringNull()
		}
		if value := r.Get("actions.bgp-actions.set-ext-community.ext-community-set-ref"); value.Exists() && !data.Statements[i].SetExtCommunity.IsNull() {
			data.Statements[i].SetExtCommunity = types.StringValue(value.String())
		} else {
			data.Statements[i].SetExtCommunity = types.StringNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *L3VPNRoutePolicies) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.Statements) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

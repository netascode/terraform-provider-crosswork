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

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type ProviderNSOPolicies struct {
	Instance           types.String              `tfsdk:"instance"`
	Id                 types.String              `tfsdk:"id"`
	ProviderName       types.String              `tfsdk:"provider_name"`
	Match              types.Bool                `tfsdk:"match"`
	MatchRule          types.String              `tfsdk:"match_rule"`
	OnboardToNso       types.Bool                `tfsdk:"onboard_to_nso"`
	OnboardToNsoRule   types.String              `tfsdk:"onboard_to_nso_rule"`
	OnboardFromNso     types.Bool                `tfsdk:"onboard_from_nso"`
	OnboardFromNsoRule types.String              `tfsdk:"onboard_from_nso_rule"`
	SyncFromNso        types.Bool                `tfsdk:"sync_from_nso"`
	SyncFromNsoRule    types.String              `tfsdk:"sync_from_nso_rule"`
	CheckSync          types.Bool                `tfsdk:"check_sync"`
	CheckSyncRule      types.String              `tfsdk:"check_sync_rule"`
	Neds               []ProviderNSOPoliciesNeds `tfsdk:"neds"`
}

type ProviderNSOPoliciesNeds struct {
	Ned  types.String `tfsdk:"ned"`
	Rule types.String `tfsdk:"rule"`
}

// End of section. //template:end types

// Written by hand: one policy document holds every provider, so the path has no id.
func (data ProviderNSOPolicies) getPath() string {
	return providerNSOPolicyPath
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data ProviderNSOPolicies) getPathGet() string {
	return "/crosswork/inventory/v1/nso/policy/query"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Written by hand: the controller answers with the policy of all NSO providers, and the resource reads the entry of its provider.
func (data ProviderNSOPolicies) readGet(c *cw.Client, path string) (cw.Res, error) {
	res, err := c.Post(path, `{}`, cw.ReadOnly)
	if err != nil {
		return res, err
	}
	name := data.ProviderName.ValueString()
	if name == "" {
		name = data.Id.ValueString()
	}
	return flattenProviderNSOPolicyEntry(res, name)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data ProviderNSOPolicies) toBody(ctx context.Context, state ProviderNSOPolicies) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Match.IsNull() {
		body, _ = sjson.Set(body, "match", data.Match.ValueBool())
	}
	if !data.MatchRule.IsNull() {
		body, _ = sjson.Set(body, "matchRule", data.MatchRule.ValueString())
	}
	if !data.OnboardToNso.IsNull() {
		body, _ = sjson.Set(body, "onboardTo", data.OnboardToNso.ValueBool())
	}
	if !data.OnboardToNsoRule.IsNull() {
		body, _ = sjson.Set(body, "onboardToRule", data.OnboardToNsoRule.ValueString())
	}
	if !data.OnboardFromNso.IsNull() {
		body, _ = sjson.Set(body, "onboardFrom", data.OnboardFromNso.ValueBool())
	}
	if !data.OnboardFromNsoRule.IsNull() {
		body, _ = sjson.Set(body, "onboardFromRule", data.OnboardFromNsoRule.ValueString())
	}
	if !data.SyncFromNso.IsNull() {
		body, _ = sjson.Set(body, "syncFrom", data.SyncFromNso.ValueBool())
	}
	if !data.SyncFromNsoRule.IsNull() {
		body, _ = sjson.Set(body, "syncFromRule", data.SyncFromNsoRule.ValueString())
	}
	if !data.CheckSync.IsNull() {
		body, _ = sjson.Set(body, "checkSync", data.CheckSync.ValueBool())
	}
	if !data.CheckSyncRule.IsNull() {
		body, _ = sjson.Set(body, "checkSyncRule", data.CheckSyncRule.ValueString())
	}
	if len(data.Neds) > 0 {
		body, _ = sjson.Set(body, "neds", []interface{}{})
		for _, item := range data.Neds {
			itemBody := ""
			if !item.Ned.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "ned", item.Ned.ValueString())
			}
			if !item.Rule.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "rule", item.Rule.ValueString())
			}
			body, _ = sjson.SetRaw(body, "neds.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *ProviderNSOPolicies) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("match"); value.Exists() {
		data.Match = types.BoolValue(value.Bool())
	} else {
		data.Match = types.BoolNull()
	}
	if value := res.Get("matchRule"); value.Exists() {
		data.MatchRule = types.StringValue(value.String())
	} else {
		data.MatchRule = types.StringNull()
	}
	if value := res.Get("onboardTo"); value.Exists() {
		data.OnboardToNso = types.BoolValue(value.Bool())
	} else {
		data.OnboardToNso = types.BoolNull()
	}
	if value := res.Get("onboardToRule"); value.Exists() {
		data.OnboardToNsoRule = types.StringValue(value.String())
	} else {
		data.OnboardToNsoRule = types.StringNull()
	}
	if value := res.Get("onboardFrom"); value.Exists() {
		data.OnboardFromNso = types.BoolValue(value.Bool())
	} else {
		data.OnboardFromNso = types.BoolNull()
	}
	if value := res.Get("onboardFromRule"); value.Exists() {
		data.OnboardFromNsoRule = types.StringValue(value.String())
	} else {
		data.OnboardFromNsoRule = types.StringNull()
	}
	if value := res.Get("syncFrom"); value.Exists() {
		data.SyncFromNso = types.BoolValue(value.Bool())
	} else {
		data.SyncFromNso = types.BoolNull()
	}
	if value := res.Get("syncFromRule"); value.Exists() {
		data.SyncFromNsoRule = types.StringValue(value.String())
	} else {
		data.SyncFromNsoRule = types.StringNull()
	}
	if value := res.Get("checkSync"); value.Exists() {
		data.CheckSync = types.BoolValue(value.Bool())
	} else {
		data.CheckSync = types.BoolNull()
	}
	if value := res.Get("checkSyncRule"); value.Exists() {
		data.CheckSyncRule = types.StringValue(value.String())
	} else {
		data.CheckSyncRule = types.StringNull()
	}
	if value := res.Get("neds"); value.Exists() && len(value.Array()) > 0 {
		data.Neds = make([]ProviderNSOPoliciesNeds, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := ProviderNSOPoliciesNeds{}
			if cValue := v.Get("ned"); cValue.Exists() {
				item.Ned = types.StringValue(cValue.String())
			} else {
				item.Ned = types.StringNull()
			}
			if cValue := v.Get("rule"); cValue.Exists() {
				item.Rule = types.StringValue(cValue.String())
			} else {
				item.Rule = types.StringNull()
			}
			data.Neds = append(data.Neds, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *ProviderNSOPolicies) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("match"); value.Exists() && !data.Match.IsNull() {
		data.Match = types.BoolValue(value.Bool())
	} else {
		data.Match = types.BoolNull()
	}
	if value := res.Get("matchRule"); value.Exists() && !data.MatchRule.IsNull() {
		data.MatchRule = types.StringValue(value.String())
	} else {
		data.MatchRule = types.StringNull()
	}
	if value := res.Get("onboardTo"); value.Exists() && !data.OnboardToNso.IsNull() {
		data.OnboardToNso = types.BoolValue(value.Bool())
	} else {
		data.OnboardToNso = types.BoolNull()
	}
	if value := res.Get("onboardToRule"); value.Exists() && !data.OnboardToNsoRule.IsNull() {
		data.OnboardToNsoRule = types.StringValue(value.String())
	} else {
		data.OnboardToNsoRule = types.StringNull()
	}
	if value := res.Get("onboardFrom"); value.Exists() && !data.OnboardFromNso.IsNull() {
		data.OnboardFromNso = types.BoolValue(value.Bool())
	} else {
		data.OnboardFromNso = types.BoolNull()
	}
	if value := res.Get("onboardFromRule"); value.Exists() && !data.OnboardFromNsoRule.IsNull() {
		data.OnboardFromNsoRule = types.StringValue(value.String())
	} else {
		data.OnboardFromNsoRule = types.StringNull()
	}
	if value := res.Get("syncFrom"); value.Exists() && !data.SyncFromNso.IsNull() {
		data.SyncFromNso = types.BoolValue(value.Bool())
	} else {
		data.SyncFromNso = types.BoolNull()
	}
	if value := res.Get("syncFromRule"); value.Exists() && !data.SyncFromNsoRule.IsNull() {
		data.SyncFromNsoRule = types.StringValue(value.String())
	} else {
		data.SyncFromNsoRule = types.StringNull()
	}
	if value := res.Get("checkSync"); value.Exists() && !data.CheckSync.IsNull() {
		data.CheckSync = types.BoolValue(value.Bool())
	} else {
		data.CheckSync = types.BoolNull()
	}
	if value := res.Get("checkSyncRule"); value.Exists() && !data.CheckSyncRule.IsNull() {
		data.CheckSyncRule = types.StringValue(value.String())
	} else {
		data.CheckSyncRule = types.StringNull()
	}
	for i := range data.Neds {
		keys := [...]string{"ned", "rule"}
		keyValues := [...]string{data.Neds[i].Ned.ValueString(), data.Neds[i].Rule.ValueString()}

		var r gjson.Result
		res.Get("neds").ForEach(
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
		if value := r.Get("ned"); value.Exists() && !data.Neds[i].Ned.IsNull() {
			data.Neds[i].Ned = types.StringValue(value.String())
		} else {
			data.Neds[i].Ned = types.StringNull()
		}
		if value := r.Get("rule"); value.Exists() && !data.Neds[i].Rule.IsNull() {
			data.Neds[i].Rule = types.StringValue(value.String())
		} else {
			data.Neds[i].Rule = types.StringNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *ProviderNSOPolicies) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Match.IsNull() {
		return false
	}
	if !data.MatchRule.IsNull() {
		return false
	}
	if !data.OnboardToNso.IsNull() {
		return false
	}
	if !data.OnboardToNsoRule.IsNull() {
		return false
	}
	if !data.OnboardFromNso.IsNull() {
		return false
	}
	if !data.OnboardFromNsoRule.IsNull() {
		return false
	}
	if !data.SyncFromNso.IsNull() {
		return false
	}
	if !data.SyncFromNsoRule.IsNull() {
		return false
	}
	if !data.CheckSync.IsNull() {
		return false
	}
	if !data.CheckSyncRule.IsNull() {
		return false
	}
	if len(data.Neds) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

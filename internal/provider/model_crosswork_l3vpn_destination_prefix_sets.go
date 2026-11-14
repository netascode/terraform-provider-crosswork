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
type L3VPNDestinationPrefixSets struct {
	Instance   types.String                           `tfsdk:"instance"`
	Id         types.String                           `tfsdk:"id"`
	Name       types.String                           `tfsdk:"name"`
	Mode       types.String                           `tfsdk:"mode"`
	PrefixList []L3VPNDestinationPrefixSetsPrefixList `tfsdk:"prefix_list"`
}

type L3VPNDestinationPrefixSetsPrefixList struct {
	IpPrefix types.String `tfsdk:"ip_prefix"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data L3VPNDestinationPrefixSets) getPath() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/cisco-l3vpn-routing-policy:l3vpn-routing-policy/defined-sets/prefix-sets/dest-prefix-set=%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

func (data L3VPNDestinationPrefixSets) getPathPost() string {
	return "/crosswork/proxy/nso/restconf/data/cisco-l3vpn-routing-policy:l3vpn-routing-policy/defined-sets/prefix-sets"
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
func (data L3VPNDestinationPrefixSets) toBody(ctx context.Context, state L3VPNDestinationPrefixSets) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "cisco-l3vpn-routing-policy:dest-prefix-set.0.name", data.Name.ValueString())
	}
	if !data.Mode.IsNull() {
		body, _ = sjson.Set(body, "cisco-l3vpn-routing-policy:dest-prefix-set.0.mode", data.Mode.ValueString())
	}
	if len(data.PrefixList) > 0 {
		body, _ = sjson.Set(body, "cisco-l3vpn-routing-policy:dest-prefix-set.0.prefixes.prefix-list", []interface{}{})
		for _, item := range data.PrefixList {
			itemBody := ""
			if !item.IpPrefix.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "ip-prefix", item.IpPrefix.ValueString())
			}
			body, _ = sjson.SetRaw(body, "cisco-l3vpn-routing-policy:dest-prefix-set.0.prefixes.prefix-list.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *L3VPNDestinationPrefixSets) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-l3vpn-routing-policy:dest-prefix-set.0.mode"); value.Exists() {
		data.Mode = types.StringValue(value.String())
	} else {
		data.Mode = types.StringNull()
	}
	if value := res.Get("cisco-l3vpn-routing-policy:dest-prefix-set.0.prefixes.prefix-list"); value.Exists() && len(value.Array()) > 0 {
		data.PrefixList = make([]L3VPNDestinationPrefixSetsPrefixList, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := L3VPNDestinationPrefixSetsPrefixList{}
			if cValue := v.Get("ip-prefix"); cValue.Exists() {
				item.IpPrefix = types.StringValue(cValue.String())
			} else {
				item.IpPrefix = types.StringNull()
			}
			data.PrefixList = append(data.PrefixList, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *L3VPNDestinationPrefixSets) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-l3vpn-routing-policy:dest-prefix-set.0.mode"); value.Exists() && !data.Mode.IsNull() {
		data.Mode = types.StringValue(value.String())
	} else {
		data.Mode = types.StringNull()
	}
	for i := range data.PrefixList {
		keys := [...]string{"ip-prefix"}
		keyValues := [...]string{data.PrefixList[i].IpPrefix.ValueString()}

		var r gjson.Result
		res.Get("cisco-l3vpn-routing-policy:dest-prefix-set.0.prefixes.prefix-list").ForEach(
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
		if value := r.Get("ip-prefix"); value.Exists() && !data.PrefixList[i].IpPrefix.IsNull() {
			data.PrefixList[i].IpPrefix = types.StringValue(value.String())
		} else {
			data.PrefixList[i].IpPrefix = types.StringNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *L3VPNDestinationPrefixSets) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Mode.IsNull() {
		return false
	}
	if len(data.PrefixList) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

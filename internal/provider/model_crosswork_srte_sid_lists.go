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
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type SRTESIDLists struct {
	Instance types.String       `tfsdk:"instance"`
	Id       types.String       `tfsdk:"id"`
	Name     types.String       `tfsdk:"name"`
	Sids     []SRTESIDListsSids `tfsdk:"sids"`
}

type SRTESIDListsSids struct {
	Index   types.Int64  `tfsdk:"index"`
	Label   types.Int64  `tfsdk:"label"`
	Address types.String `tfsdk:"address"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data SRTESIDLists) getPath() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/cisco-sr-te-cfp:sr-te/cisco-sr-te-cfp-sr-policies:policies/sid-list=%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

func (data SRTESIDLists) getPathPost() string {
	return "/crosswork/proxy/nso/restconf/data/cisco-sr-te-cfp:sr-te/cisco-sr-te-cfp-sr-policies:policies"
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
func (data SRTESIDLists) toBody(ctx context.Context, state SRTESIDLists) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:sid-list.0.name", data.Name.ValueString())
	}
	if len(data.Sids) > 0 {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:sid-list.0.sid", []interface{}{})
		for _, item := range data.Sids {
			itemBody := ""
			if !item.Index.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "index", item.Index.ValueInt64())
			}
			if !item.Label.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "mpls.label", item.Label.ValueInt64())
			}
			if !item.Address.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "ipv4.address", item.Address.ValueString())
			}
			body, _ = sjson.SetRaw(body, "cisco-sr-te-cfp-sr-policies:sid-list.0.sid.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *SRTESIDLists) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-sr-te-cfp-sr-policies:sid-list.0.sid"); value.Exists() && len(value.Array()) > 0 {
		data.Sids = make([]SRTESIDListsSids, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := SRTESIDListsSids{}
			if cValue := v.Get("index"); cValue.Exists() {
				item.Index = types.Int64Value(cValue.Int())
			} else {
				item.Index = types.Int64Null()
			}
			if cValue := v.Get("mpls.label"); cValue.Exists() {
				item.Label = types.Int64Value(cValue.Int())
			} else {
				item.Label = types.Int64Null()
			}
			if cValue := v.Get("ipv4.address"); cValue.Exists() {
				item.Address = types.StringValue(cValue.String())
			} else {
				item.Address = types.StringNull()
			}
			data.Sids = append(data.Sids, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *SRTESIDLists) updateFromBody(ctx context.Context, res gjson.Result) {
	for i := range data.Sids {
		keys := [...]string{"index"}
		keyValues := [...]string{strconv.FormatInt(data.Sids[i].Index.ValueInt64(), 10)}

		var r gjson.Result
		res.Get("cisco-sr-te-cfp-sr-policies:sid-list.0.sid").ForEach(
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
		if value := r.Get("index"); value.Exists() && !data.Sids[i].Index.IsNull() {
			data.Sids[i].Index = types.Int64Value(value.Int())
		} else {
			data.Sids[i].Index = types.Int64Null()
		}
		if value := r.Get("mpls.label"); value.Exists() && !data.Sids[i].Label.IsNull() {
			data.Sids[i].Label = types.Int64Value(value.Int())
		} else {
			data.Sids[i].Label = types.Int64Null()
		}
		if value := r.Get("ipv4.address"); value.Exists() && !data.Sids[i].Address.IsNull() {
			data.Sids[i].Address = types.StringValue(value.String())
		} else {
			data.Sids[i].Address = types.StringNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *SRTESIDLists) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.Sids) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

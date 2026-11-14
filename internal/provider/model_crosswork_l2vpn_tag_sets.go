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

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type L2VPNTagSets struct {
	Instance types.String `tfsdk:"instance"`
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	TagValue types.Set    `tfsdk:"tag_value"`
	CoFlag   types.String `tfsdk:"co_flag"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data L2VPNTagSets) getPath() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/cisco-l2vpn-routing-policy:l2vpn-routing-policy/defined-sets/tag-sets/tag-set=%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

func (data L2VPNTagSets) getPathPost() string {
	return "/crosswork/proxy/nso/restconf/data/cisco-l2vpn-routing-policy:l2vpn-routing-policy/defined-sets/tag-sets"
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
func (data L2VPNTagSets) toBody(ctx context.Context, state L2VPNTagSets) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-routing-policy:tag-set.0.name", data.Name.ValueString())
	}
	if !data.TagValue.IsNull() {
		var values []int64
		data.TagValue.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "cisco-l2vpn-routing-policy:tag-set.0.tag-value", values)
	}
	if !data.CoFlag.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-routing-policy:tag-set.0.co-flag", data.CoFlag.ValueString())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *L2VPNTagSets) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-l2vpn-routing-policy:tag-set.0.tag-value"); value.Exists() && len(value.Array()) > 0 {
		data.TagValue = helpers.GetInt64Set(value.Array())
	} else {
		data.TagValue = types.SetNull(types.Int64Type)
	}
	if value := res.Get("cisco-l2vpn-routing-policy:tag-set.0.co-flag"); value.Exists() {
		data.CoFlag = types.StringValue(value.String())
	} else {
		data.CoFlag = types.StringNull()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *L2VPNTagSets) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-l2vpn-routing-policy:tag-set.0.tag-value"); value.Exists() && !data.TagValue.IsNull() {
		data.TagValue = helpers.GetInt64Set(value.Array())
	} else {
		data.TagValue = types.SetNull(types.Int64Type)
	}
	if value := res.Get("cisco-l2vpn-routing-policy:tag-set.0.co-flag"); value.Exists() && !data.CoFlag.IsNull() {
		data.CoFlag = types.StringValue(value.String())
	} else {
		data.CoFlag = types.StringNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *L2VPNTagSets) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.TagValue.IsNull() {
		return false
	}
	if !data.CoFlag.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

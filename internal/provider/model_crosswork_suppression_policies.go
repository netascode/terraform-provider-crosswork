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
type SuppressionPolicies struct {
	Instance     types.String `tfsdk:"instance"`
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Action       types.String `tfsdk:"action"`
	DeviceGroups types.Set    `tfsdk:"device_groups"`
	Criteria     types.String `tfsdk:"criteria"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data SuppressionPolicies) getPath() string {
	return fmt.Sprintf("/crosswork/alarm/v1/suppressionpolicy/%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data SuppressionPolicies) getPathGet() string {
	return "/crosswork/alarm/v1/suppressionpolicy"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

func (data SuppressionPolicies) getPathPost() string {
	return "/crosswork/alarm/v1/suppressionpolicy"
}

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

func (data SuppressionPolicies) getPathPut() string {
	return "/crosswork/alarm/v1/suppressionpolicy"
}

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data SuppressionPolicies) toBody(ctx context.Context, state SuppressionPolicies) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "policyname", data.Name.ValueString())
	}
	if !data.Description.IsNull() {
		body, _ = sjson.Set(body, "description", data.Description.ValueString())
	}
	if !data.Action.IsNull() {
		body, _ = sjson.Set(body, "action", data.Action.ValueString())
	}
	if !data.DeviceGroups.IsNull() {
		var values []string
		data.DeviceGroups.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "deviceGroups", values)
	}
	if !data.Criteria.IsNull() {
		body, _ = sjson.Set(body, "criteria", data.Criteria.ValueString())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *SuppressionPolicies) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("description"); value.Exists() && value.String() != "" {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	if value := res.Get("action"); value.Exists() {
		data.Action = types.StringValue(value.String())
	} else {
		data.Action = types.StringNull()
	}
	if value := res.Get("deviceGroups"); value.Exists() && len(value.Array()) > 0 {
		data.DeviceGroups = helpers.GetStringSet(value.Array())
	} else {
		data.DeviceGroups = types.SetNull(types.StringType)
	}
	if value := res.Get("criteria"); value.Exists() {
		data.Criteria = types.StringValue(value.String())
	} else {
		data.Criteria = types.StringNull()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *SuppressionPolicies) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("description"); value.Exists() && value.String() != "" && !data.Description.IsNull() {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	if value := res.Get("action"); value.Exists() && !data.Action.IsNull() {
		data.Action = types.StringValue(value.String())
	} else {
		data.Action = types.StringNull()
	}
	if value := res.Get("deviceGroups"); value.Exists() && !data.DeviceGroups.IsNull() {
		data.DeviceGroups = helpers.GetStringSet(value.Array())
	} else {
		data.DeviceGroups = types.SetNull(types.StringType)
	}
	if value := res.Get("criteria"); value.Exists() && !data.Criteria.IsNull() {
		data.Criteria = types.StringValue(value.String())
	} else {
		data.Criteria = types.StringNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *SuppressionPolicies) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Description.IsNull() {
		return false
	}
	if !data.Action.IsNull() {
		return false
	}
	if !data.DeviceGroups.IsNull() {
		return false
	}
	if !data.Criteria.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

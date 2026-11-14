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

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type AssignDevicesToDeviceGroup struct {
	Instance types.String `tfsdk:"instance"`
	Id       types.String `tfsdk:"id"`
	GroupId  types.String `tfsdk:"group_id"`
	Devices  types.Set    `tfsdk:"devices"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data AssignDevicesToDeviceGroup) getPath() string {
	return "/crosswork/grouping/v1/grouping/group/member"
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

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data AssignDevicesToDeviceGroup) toBody(ctx context.Context, state AssignDevicesToDeviceGroup) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.GroupId.IsNull() {
		body, _ = sjson.Set(body, "newGroupUuid", data.GroupId.ValueString())
	}
	if !data.Devices.IsNull() {
		var values []string
		data.Devices.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "references", values)
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *AssignDevicesToDeviceGroup) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("newGroupUuid"); value.Exists() {
		data.GroupId = types.StringValue(value.String())
	} else {
		data.GroupId = types.StringNull()
	}
	if value := res.Get("references"); value.Exists() && len(value.Array()) > 0 {
		data.Devices = helpers.GetStringSet(value.Array())
	} else {
		data.Devices = types.SetNull(types.StringType)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *AssignDevicesToDeviceGroup) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("newGroupUuid"); value.Exists() && !data.GroupId.IsNull() {
		data.GroupId = types.StringValue(value.String())
	} else {
		data.GroupId = types.StringNull()
	}
	if value := res.Get("references"); value.Exists() && !data.Devices.IsNull() {
		data.Devices = helpers.GetStringSet(value.Array())
	} else {
		data.Devices = types.SetNull(types.StringType)
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *AssignDevicesToDeviceGroup) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Devices.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

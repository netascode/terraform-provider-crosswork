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

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type TEGeneralSettings struct {
	Instance                             types.String `tfsdk:"instance"`
	Id                                   types.String `tfsdk:"id"`
	PolicyUnusedThreshold                types.Int64  `tfsdk:"policy_unused_threshold"`
	PolicyProvisioningTimeout            types.Int64  `tfsdk:"policy_provisioning_timeout"`
	BandwidthOnDemandProvisioningTimeout types.Int64  `tfsdk:"bandwidth_on_demand_provisioning_timeout"`
	IgpPathRequestTimeout                types.Int64  `tfsdk:"igp_path_request_timeout"`
	PolicyDetailsRequestTimeout          types.Int64  `tfsdk:"policy_details_request_timeout"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data TEGeneralSettings) getPath() string {
	return "/crosswork/optima/api/v1/setTimeouts"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data TEGeneralSettings) getPathGet() string {
	return "/crosswork/optima/api/v1/getTimeouts"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

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
func (data TEGeneralSettings) toBody(ctx context.Context, state TEGeneralSettings) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.PolicyUnusedThreshold.IsNull() {
		body, _ = sjson.Set(body, "OPTIMA_POLICY_UNUSED_THRESHOLD", fmt.Sprint(data.PolicyUnusedThreshold.ValueInt64()))
	}
	if !data.PolicyProvisioningTimeout.IsNull() {
		body, _ = sjson.Set(body, "PCE_DISPATCH_NON_BWOD_TIMEOUT", fmt.Sprint(data.PolicyProvisioningTimeout.ValueInt64()))
	}
	if !data.BandwidthOnDemandProvisioningTimeout.IsNull() {
		body, _ = sjson.Set(body, "PCE_DISPATCH_BWOD_TIMEOUT", fmt.Sprint(data.BandwidthOnDemandProvisioningTimeout.ValueInt64()))
	}
	if !data.IgpPathRequestTimeout.IsNull() {
		body, _ = sjson.Set(body, "OPTIMA_ENG_IGP_PATH_TIMEOUT", fmt.Sprint(data.IgpPathRequestTimeout.ValueInt64()))
	}
	if !data.PolicyDetailsRequestTimeout.IsNull() {
		body, _ = sjson.Set(body, "OPTIMA_ENG_UTILIZATION_TIMEOUT", fmt.Sprint(data.PolicyDetailsRequestTimeout.ValueInt64()))
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *TEGeneralSettings) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("te_general_settings")
	if value := res.Get("OPTIMA_POLICY_UNUSED_THRESHOLD"); value.Exists() {
		data.PolicyUnusedThreshold = types.Int64Value(value.Int())
	} else {
		data.PolicyUnusedThreshold = types.Int64Value(1000)
	}
	if value := res.Get("PCE_DISPATCH_NON_BWOD_TIMEOUT"); value.Exists() {
		data.PolicyProvisioningTimeout = types.Int64Value(value.Int())
	} else {
		data.PolicyProvisioningTimeout = types.Int64Value(60)
	}
	if value := res.Get("PCE_DISPATCH_BWOD_TIMEOUT"); value.Exists() {
		data.BandwidthOnDemandProvisioningTimeout = types.Int64Value(value.Int())
	} else {
		data.BandwidthOnDemandProvisioningTimeout = types.Int64Value(90)
	}
	if value := res.Get("OPTIMA_ENG_IGP_PATH_TIMEOUT"); value.Exists() {
		data.IgpPathRequestTimeout = types.Int64Value(value.Int())
	} else {
		data.IgpPathRequestTimeout = types.Int64Value(180)
	}
	if value := res.Get("OPTIMA_ENG_UTILIZATION_TIMEOUT"); value.Exists() {
		data.PolicyDetailsRequestTimeout = types.Int64Value(value.Int())
	} else {
		data.PolicyDetailsRequestTimeout = types.Int64Value(180)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *TEGeneralSettings) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("OPTIMA_POLICY_UNUSED_THRESHOLD"); value.Exists() && !data.PolicyUnusedThreshold.IsNull() {
		data.PolicyUnusedThreshold = types.Int64Value(value.Int())
	} else if data.PolicyUnusedThreshold.ValueInt64() != 1000 {
		data.PolicyUnusedThreshold = types.Int64Null()
	}
	if value := res.Get("PCE_DISPATCH_NON_BWOD_TIMEOUT"); value.Exists() && !data.PolicyProvisioningTimeout.IsNull() {
		data.PolicyProvisioningTimeout = types.Int64Value(value.Int())
	} else if data.PolicyProvisioningTimeout.ValueInt64() != 60 {
		data.PolicyProvisioningTimeout = types.Int64Null()
	}
	if value := res.Get("PCE_DISPATCH_BWOD_TIMEOUT"); value.Exists() && !data.BandwidthOnDemandProvisioningTimeout.IsNull() {
		data.BandwidthOnDemandProvisioningTimeout = types.Int64Value(value.Int())
	} else if data.BandwidthOnDemandProvisioningTimeout.ValueInt64() != 90 {
		data.BandwidthOnDemandProvisioningTimeout = types.Int64Null()
	}
	if value := res.Get("OPTIMA_ENG_IGP_PATH_TIMEOUT"); value.Exists() && !data.IgpPathRequestTimeout.IsNull() {
		data.IgpPathRequestTimeout = types.Int64Value(value.Int())
	} else if data.IgpPathRequestTimeout.ValueInt64() != 180 {
		data.IgpPathRequestTimeout = types.Int64Null()
	}
	if value := res.Get("OPTIMA_ENG_UTILIZATION_TIMEOUT"); value.Exists() && !data.PolicyDetailsRequestTimeout.IsNull() {
		data.PolicyDetailsRequestTimeout = types.Int64Value(value.Int())
	} else if data.PolicyDetailsRequestTimeout.ValueInt64() != 180 {
		data.PolicyDetailsRequestTimeout = types.Int64Null()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *TEGeneralSettings) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.PolicyUnusedThreshold.IsNull() {
		return false
	}
	if !data.PolicyProvisioningTimeout.IsNull() {
		return false
	}
	if !data.BandwidthOnDemandProvisioningTimeout.IsNull() {
		return false
	}
	if !data.IgpPathRequestTimeout.IsNull() {
		return false
	}
	if !data.PolicyDetailsRequestTimeout.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

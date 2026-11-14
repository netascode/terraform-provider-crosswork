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

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type L2VPNIDPools struct {
	Instance            types.String `tfsdk:"instance"`
	Id                  types.String `tfsdk:"id"`
	EviIdPool           types.String `tfsdk:"evi_id_pool"`
	EviSourceTargetPool types.String `tfsdk:"evi_source_target_pool"`
	RtPool              types.String `tfsdk:"rt_pool"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data L2VPNIDPools) getPath() string {
	return "/crosswork/proxy/nso/restconf/data/ietf-l2vpn-ntw:l2vpn-ntw/cisco-l2vpn-ntw:id-pools"
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

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data L2VPNIDPools) toBody(ctx context.Context, state L2VPNIDPools) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.EviIdPool.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:id-pools.evi-id-pool-name", data.EviIdPool.ValueString())
	}
	if !data.EviSourceTargetPool.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:id-pools.evi-source-target-pool-name", data.EviSourceTargetPool.ValueString())
	}
	if !data.RtPool.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:id-pools.rt-pool-name", data.RtPool.ValueString())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *L2VPNIDPools) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("l2vpn_id_pools")
	if value := res.Get("cisco-l2vpn-ntw:id-pools.evi-id-pool-name"); value.Exists() {
		data.EviIdPool = types.StringValue(value.String())
	} else {
		data.EviIdPool = types.StringNull()
	}
	if value := res.Get("cisco-l2vpn-ntw:id-pools.evi-source-target-pool-name"); value.Exists() {
		data.EviSourceTargetPool = types.StringValue(value.String())
	} else {
		data.EviSourceTargetPool = types.StringNull()
	}
	if value := res.Get("cisco-l2vpn-ntw:id-pools.rt-pool-name"); value.Exists() {
		data.RtPool = types.StringValue(value.String())
	} else {
		data.RtPool = types.StringNull()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *L2VPNIDPools) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-l2vpn-ntw:id-pools.evi-id-pool-name"); value.Exists() && !data.EviIdPool.IsNull() {
		data.EviIdPool = types.StringValue(value.String())
	} else {
		data.EviIdPool = types.StringNull()
	}
	if value := res.Get("cisco-l2vpn-ntw:id-pools.evi-source-target-pool-name"); value.Exists() && !data.EviSourceTargetPool.IsNull() {
		data.EviSourceTargetPool = types.StringValue(value.String())
	} else {
		data.EviSourceTargetPool = types.StringNull()
	}
	if value := res.Get("cisco-l2vpn-ntw:id-pools.rt-pool-name"); value.Exists() && !data.RtPool.IsNull() {
		data.RtPool = types.StringValue(value.String())
	} else {
		data.RtPool = types.StringNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *L2VPNIDPools) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.EviIdPool.IsNull() {
		return false
	}
	if !data.EviSourceTargetPool.IsNull() {
		return false
	}
	if !data.RtPool.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

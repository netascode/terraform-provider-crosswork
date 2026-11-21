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
type AAASessionSettings struct {
	Instance                types.String `tfsdk:"instance"`
	Id                      types.String `tfsdk:"id"`
	FallbackType            types.String `tfsdk:"fallback_type"`
	IdleSessionTimeout      types.Int64  `tfsdk:"idle_session_timeout"`
	ParallelSessions        types.Int64  `tfsdk:"parallel_sessions"`
	ParallelSessionsPerUser types.Int64  `tfsdk:"parallel_sessions_per_user"`
	SourceIpAuditing        types.Bool   `tfsdk:"source_ip_auditing"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data AAASessionSettings) getPath() string {
	return "/crosswork/aaa/v1/sessionconfig"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

func (data AAASessionSettings) getPathDelete() string {
	return "/crosswork/aaa/v1/sessionconfig"
}

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data AAASessionSettings) getPathGet() string {
	return "/crosswork/aaa/v1/sessionconfig"
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
func (data AAASessionSettings) toBody(ctx context.Context, state AAASessionSettings) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.FallbackType.IsNull() {
		body, _ = sjson.Set(body, "FallbackType", data.FallbackType.ValueString())
	}
	if !data.IdleSessionTimeout.IsNull() {
		body, _ = sjson.Set(body, "IdleSessionTimeout", data.IdleSessionTimeout.ValueInt64())
	}
	if !data.ParallelSessions.IsNull() {
		body, _ = sjson.Set(body, "NumParallelSessions", data.ParallelSessions.ValueInt64())
	}
	if !data.ParallelSessionsPerUser.IsNull() {
		body, _ = sjson.Set(body, "NumParallelSessionsPerUser", data.ParallelSessionsPerUser.ValueInt64())
	}
	if !data.SourceIpAuditing.IsNull() {
		body, _ = sjson.Set(body, "IsEnabled", data.SourceIpAuditing.ValueBool())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *AAASessionSettings) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("aaa_session_settings")
	if value := res.Get("FallbackType"); value.Exists() {
		data.FallbackType = types.StringValue(value.String())
	} else {
		data.FallbackType = types.StringValue("FB_ON_NORESP_AUTH_FAILURE")
	}
	if value := res.Get("IdleSessionTimeout"); value.Exists() {
		data.IdleSessionTimeout = types.Int64Value(value.Int())
	} else {
		data.IdleSessionTimeout = types.Int64Value(30)
	}
	if value := res.Get("NumParallelSessions"); value.Exists() {
		data.ParallelSessions = types.Int64Value(value.Int())
	} else {
		data.ParallelSessions = types.Int64Value(100)
	}
	if value := res.Get("NumParallelSessionsPerUser"); value.Exists() {
		data.ParallelSessionsPerUser = types.Int64Value(value.Int())
	} else {
		data.ParallelSessionsPerUser = types.Int64Value(100)
	}
	if value := res.Get("IsEnabled"); value.Exists() {
		data.SourceIpAuditing = types.BoolValue(value.Bool())
	} else {
		data.SourceIpAuditing = types.BoolValue(false)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *AAASessionSettings) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("FallbackType"); value.Exists() && !data.FallbackType.IsNull() {
		data.FallbackType = types.StringValue(value.String())
	} else if data.FallbackType.ValueString() != "FB_ON_NORESP_AUTH_FAILURE" {
		data.FallbackType = types.StringNull()
	}
	if value := res.Get("IdleSessionTimeout"); value.Exists() && !data.IdleSessionTimeout.IsNull() {
		data.IdleSessionTimeout = types.Int64Value(value.Int())
	} else if data.IdleSessionTimeout.ValueInt64() != 30 {
		data.IdleSessionTimeout = types.Int64Null()
	}
	if value := res.Get("NumParallelSessions"); value.Exists() && !data.ParallelSessions.IsNull() {
		data.ParallelSessions = types.Int64Value(value.Int())
	} else if data.ParallelSessions.ValueInt64() != 100 {
		data.ParallelSessions = types.Int64Null()
	}
	if value := res.Get("NumParallelSessionsPerUser"); value.Exists() && !data.ParallelSessionsPerUser.IsNull() {
		data.ParallelSessionsPerUser = types.Int64Value(value.Int())
	} else if data.ParallelSessionsPerUser.ValueInt64() != 100 {
		data.ParallelSessionsPerUser = types.Int64Null()
	}
	if value := res.Get("IsEnabled"); value.Exists() && !data.SourceIpAuditing.IsNull() {
		data.SourceIpAuditing = types.BoolValue(value.Bool())
	} else if data.SourceIpAuditing.ValueBool() != false {
		data.SourceIpAuditing = types.BoolNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *AAASessionSettings) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.FallbackType.IsNull() {
		return false
	}
	if !data.IdleSessionTimeout.IsNull() {
		return false
	}
	if !data.ParallelSessions.IsNull() {
		return false
	}
	if !data.ParallelSessionsPerUser.IsNull() {
		return false
	}
	if !data.SourceIpAuditing.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

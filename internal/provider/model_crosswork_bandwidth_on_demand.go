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
type BandwidthOnDemand struct {
	Instance                    types.String `tfsdk:"instance"`
	Id                          types.String `tfsdk:"id"`
	Enable                      types.Bool   `tfsdk:"enable"`
	PrimaryObjective            types.String `tfsdk:"primary_objective"`
	LinkUtilizationThreshold    types.Int64  `tfsdk:"link_utilization_threshold"`
	ReoptimizationInterval      types.Int64  `tfsdk:"reoptimization_interval"`
	MetricReoptimizationHours   types.Int64  `tfsdk:"metric_reoptimization_hours"`
	MetricReoptimizationMinutes types.Int64  `tfsdk:"metric_reoptimization_minutes"`
	SrPolicyTraffic             types.String `tfsdk:"sr_policy_traffic"`
	ModelUpdateThrottle         types.Int64  `tfsdk:"model_update_throttle"`
	OptimizerEventThreshold     types.Int64  `tfsdk:"optimizer_event_threshold"`
	OptimizerRunTimeLimit       types.Int64  `tfsdk:"optimizer_run_time_limit"`
	AllowedPolicyViolations     types.String `tfsdk:"allowed_policy_violations"`
	PreferStrictSids            types.Bool   `tfsdk:"prefer_strict_sids"`
	DebugOptimizer              types.Bool   `tfsdk:"debug_optimizer"`
	DebugOptimizerMaxFiles      types.Int64  `tfsdk:"debug_optimizer_max_files"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data BandwidthOnDemand) getPath() string {
	return "/crosswork/optima/api/v1/functionPacks/bwod"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

func (data BandwidthOnDemand) getPathDelete() string {
	return "/crosswork/optima/api/v1/functionPacks/bwod"
}

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data BandwidthOnDemand) getPathGet() string {
	return "/crosswork/optima/api/v1/functionPacks/bwod/request"
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
func (data BandwidthOnDemand) toBody(ctx context.Context, state BandwidthOnDemand) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Enable.IsNull() {
		body, _ = sjson.Set(body, "config.enable", data.Enable.ValueBool())
	}
	if !data.PrimaryObjective.IsNull() {
		body, _ = sjson.Set(body, "config.primaryObjective", data.PrimaryObjective.ValueString())
	}
	if !data.LinkUtilizationThreshold.IsNull() {
		body, _ = sjson.Set(body, "config.utilThreshold", data.LinkUtilizationThreshold.ValueInt64())
	}
	if !data.ReoptimizationInterval.IsNull() {
		body, _ = sjson.Set(body, "config.reoptInterval", data.ReoptimizationInterval.ValueInt64())
	}
	if !data.MetricReoptimizationHours.IsNull() {
		body, _ = sjson.Set(body, "config.metricReoptHours", data.MetricReoptimizationHours.ValueInt64())
	}
	if !data.MetricReoptimizationMinutes.IsNull() {
		body, _ = sjson.Set(body, "config.metricReoptMins", data.MetricReoptimizationMinutes.ValueInt64())
	}
	if !data.SrPolicyTraffic.IsNull() {
		body, _ = sjson.Set(body, "config.lspTraffic", data.SrPolicyTraffic.ValueString())
	}
	if !data.ModelUpdateThrottle.IsNull() {
		body, _ = sjson.Set(body, "config.modelUpdateThrottle", data.ModelUpdateThrottle.ValueInt64())
	}
	if !data.OptimizerEventThreshold.IsNull() {
		body, _ = sjson.Set(body, "config.optTimeEventThreshold", data.OptimizerEventThreshold.ValueInt64())
	}
	if !data.OptimizerRunTimeLimit.IsNull() {
		body, _ = sjson.Set(body, "config.optTimeLimit", data.OptimizerRunTimeLimit.ValueInt64())
	}
	if !data.AllowedPolicyViolations.IsNull() {
		body, _ = sjson.Set(body, "config.allowedPolicyViolations", data.AllowedPolicyViolations.ValueString())
	}
	if !data.PreferStrictSids.IsNull() {
		body, _ = sjson.Set(body, "config.strictSids", data.PreferStrictSids.ValueBool())
	}
	if !data.DebugOptimizer.IsNull() {
		body, _ = sjson.Set(body, "config.debugOpt", data.DebugOptimizer.ValueBool())
	}
	if !data.DebugOptimizerMaxFiles.IsNull() {
		body, _ = sjson.Set(body, "config.debugOptMaxFiles", data.DebugOptimizerMaxFiles.ValueInt64())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *BandwidthOnDemand) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("bandwidth_on_demand")
	if value := res.Get("data.config.enable"); value.Exists() {
		data.Enable = types.BoolValue(value.Bool())
	} else {
		data.Enable = types.BoolValue(false)
	}
	if value := res.Get("data.config.primaryObjective"); value.Exists() {
		data.PrimaryObjective = types.StringValue(value.String())
	} else {
		data.PrimaryObjective = types.StringValue("MAX_AVAIL_BW")
	}
	if value := res.Get("data.config.utilThreshold"); value.Exists() {
		data.LinkUtilizationThreshold = types.Int64Value(value.Int())
	} else {
		data.LinkUtilizationThreshold = types.Int64Value(80)
	}
	if value := res.Get("data.config.reoptInterval"); value.Exists() {
		data.ReoptimizationInterval = types.Int64Value(value.Int())
	} else {
		data.ReoptimizationInterval = types.Int64Value(60)
	}
	if value := res.Get("data.config.metricReoptHours"); value.Exists() {
		data.MetricReoptimizationHours = types.Int64Value(value.Int())
	} else {
		data.MetricReoptimizationHours = types.Int64Value(1)
	}
	if value := res.Get("data.config.metricReoptMins"); value.Exists() {
		data.MetricReoptimizationMinutes = types.Int64Value(value.Int())
	} else {
		data.MetricReoptimizationMinutes = types.Int64Value(30)
	}
	if value := res.Get("data.config.lspTraffic"); value.Exists() {
		data.SrPolicyTraffic = types.StringValue(value.String())
	} else {
		data.SrPolicyTraffic = types.StringValue("MAX_MEASURED_REQUESTED")
	}
	if value := res.Get("data.config.modelUpdateThrottle"); value.Exists() {
		data.ModelUpdateThrottle = types.Int64Value(value.Int())
	} else {
		data.ModelUpdateThrottle = types.Int64Value(0)
	}
	if value := res.Get("data.config.optTimeEventThreshold"); value.Exists() {
		data.OptimizerEventThreshold = types.Int64Value(value.Int())
	} else {
		data.OptimizerEventThreshold = types.Int64Value(45)
	}
	if value := res.Get("data.config.optTimeLimit"); value.Exists() {
		data.OptimizerRunTimeLimit = types.Int64Value(value.Int())
	} else {
		data.OptimizerRunTimeLimit = types.Int64Value(60)
	}
	if value := res.Get("data.config.allowedPolicyViolations"); value.Exists() {
		data.AllowedPolicyViolations = types.StringValue(value.String())
	} else {
		data.AllowedPolicyViolations = types.StringValue("LOOSE")
	}
	if value := res.Get("data.config.strictSids"); value.Exists() {
		data.PreferStrictSids = types.BoolValue(value.Bool())
	} else {
		data.PreferStrictSids = types.BoolValue(false)
	}
	if value := res.Get("data.config.debugOpt"); value.Exists() {
		data.DebugOptimizer = types.BoolValue(value.Bool())
	} else {
		data.DebugOptimizer = types.BoolValue(false)
	}
	if value := res.Get("data.config.debugOptMaxFiles"); value.Exists() {
		data.DebugOptimizerMaxFiles = types.Int64Value(value.Int())
	} else {
		data.DebugOptimizerMaxFiles = types.Int64Value(30)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *BandwidthOnDemand) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("data.config.enable"); value.Exists() && !data.Enable.IsNull() {
		data.Enable = types.BoolValue(value.Bool())
	} else if data.Enable.ValueBool() != false {
		data.Enable = types.BoolNull()
	}
	if value := res.Get("data.config.primaryObjective"); value.Exists() && !data.PrimaryObjective.IsNull() {
		data.PrimaryObjective = types.StringValue(value.String())
	} else if data.PrimaryObjective.ValueString() != "MAX_AVAIL_BW" {
		data.PrimaryObjective = types.StringNull()
	}
	if value := res.Get("data.config.utilThreshold"); value.Exists() && !data.LinkUtilizationThreshold.IsNull() {
		data.LinkUtilizationThreshold = types.Int64Value(value.Int())
	} else if data.LinkUtilizationThreshold.ValueInt64() != 80 {
		data.LinkUtilizationThreshold = types.Int64Null()
	}
	if value := res.Get("data.config.reoptInterval"); value.Exists() && !data.ReoptimizationInterval.IsNull() {
		data.ReoptimizationInterval = types.Int64Value(value.Int())
	} else if data.ReoptimizationInterval.ValueInt64() != 60 {
		data.ReoptimizationInterval = types.Int64Null()
	}
	if value := res.Get("data.config.metricReoptHours"); value.Exists() && !data.MetricReoptimizationHours.IsNull() {
		data.MetricReoptimizationHours = types.Int64Value(value.Int())
	} else if data.MetricReoptimizationHours.ValueInt64() != 1 {
		data.MetricReoptimizationHours = types.Int64Null()
	}
	if value := res.Get("data.config.metricReoptMins"); value.Exists() && !data.MetricReoptimizationMinutes.IsNull() {
		data.MetricReoptimizationMinutes = types.Int64Value(value.Int())
	} else if data.MetricReoptimizationMinutes.ValueInt64() != 30 {
		data.MetricReoptimizationMinutes = types.Int64Null()
	}
	if value := res.Get("data.config.lspTraffic"); value.Exists() && !data.SrPolicyTraffic.IsNull() {
		data.SrPolicyTraffic = types.StringValue(value.String())
	} else if data.SrPolicyTraffic.ValueString() != "MAX_MEASURED_REQUESTED" {
		data.SrPolicyTraffic = types.StringNull()
	}
	if value := res.Get("data.config.modelUpdateThrottle"); value.Exists() && !data.ModelUpdateThrottle.IsNull() {
		data.ModelUpdateThrottle = types.Int64Value(value.Int())
	} else if data.ModelUpdateThrottle.ValueInt64() != 0 {
		data.ModelUpdateThrottle = types.Int64Null()
	}
	if value := res.Get("data.config.optTimeEventThreshold"); value.Exists() && !data.OptimizerEventThreshold.IsNull() {
		data.OptimizerEventThreshold = types.Int64Value(value.Int())
	} else if data.OptimizerEventThreshold.ValueInt64() != 45 {
		data.OptimizerEventThreshold = types.Int64Null()
	}
	if value := res.Get("data.config.optTimeLimit"); value.Exists() && !data.OptimizerRunTimeLimit.IsNull() {
		data.OptimizerRunTimeLimit = types.Int64Value(value.Int())
	} else if data.OptimizerRunTimeLimit.ValueInt64() != 60 {
		data.OptimizerRunTimeLimit = types.Int64Null()
	}
	if value := res.Get("data.config.allowedPolicyViolations"); value.Exists() && !data.AllowedPolicyViolations.IsNull() {
		data.AllowedPolicyViolations = types.StringValue(value.String())
	} else if data.AllowedPolicyViolations.ValueString() != "LOOSE" {
		data.AllowedPolicyViolations = types.StringNull()
	}
	if value := res.Get("data.config.strictSids"); value.Exists() && !data.PreferStrictSids.IsNull() {
		data.PreferStrictSids = types.BoolValue(value.Bool())
	} else if data.PreferStrictSids.ValueBool() != false {
		data.PreferStrictSids = types.BoolNull()
	}
	if value := res.Get("data.config.debugOpt"); value.Exists() && !data.DebugOptimizer.IsNull() {
		data.DebugOptimizer = types.BoolValue(value.Bool())
	} else if data.DebugOptimizer.ValueBool() != false {
		data.DebugOptimizer = types.BoolNull()
	}
	if value := res.Get("data.config.debugOptMaxFiles"); value.Exists() && !data.DebugOptimizerMaxFiles.IsNull() {
		data.DebugOptimizerMaxFiles = types.Int64Value(value.Int())
	} else if data.DebugOptimizerMaxFiles.ValueInt64() != 30 {
		data.DebugOptimizerMaxFiles = types.Int64Null()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *BandwidthOnDemand) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Enable.IsNull() {
		return false
	}
	if !data.PrimaryObjective.IsNull() {
		return false
	}
	if !data.LinkUtilizationThreshold.IsNull() {
		return false
	}
	if !data.ReoptimizationInterval.IsNull() {
		return false
	}
	if !data.MetricReoptimizationHours.IsNull() {
		return false
	}
	if !data.MetricReoptimizationMinutes.IsNull() {
		return false
	}
	if !data.SrPolicyTraffic.IsNull() {
		return false
	}
	if !data.ModelUpdateThrottle.IsNull() {
		return false
	}
	if !data.OptimizerEventThreshold.IsNull() {
		return false
	}
	if !data.OptimizerRunTimeLimit.IsNull() {
		return false
	}
	if !data.AllowedPolicyViolations.IsNull() {
		return false
	}
	if !data.PreferStrictSids.IsNull() {
		return false
	}
	if !data.DebugOptimizer.IsNull() {
		return false
	}
	if !data.DebugOptimizerMaxFiles.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

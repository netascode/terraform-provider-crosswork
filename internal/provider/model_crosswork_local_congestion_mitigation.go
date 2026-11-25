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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type LocalCongestionMitigation struct {
	Instance                          types.String `tfsdk:"instance"`
	Id                                types.String `tfsdk:"id"`
	DomainId                          types.String `tfsdk:"domain_id"`
	Enable                            types.Bool   `tfsdk:"enable"`
	Color                             types.Int64  `tfsdk:"color"`
	UtilizationThreshold              types.Int64  `tfsdk:"utilization_threshold"`
	UtilizationHoldMargin             types.Int64  `tfsdk:"utilization_hold_margin"`
	DeleteLsps                        types.Bool   `tfsdk:"delete_lsps"`
	ProfileId                         types.Int64  `tfsdk:"profile_id"`
	CongestionCheckInterval           types.Int64  `tfsdk:"congestion_check_interval"`
	MaxTteGroup                       types.Int64  `tfsdk:"max_tte_group"`
	IncludeAllInterfaces              types.Bool   `tfsdk:"include_all_interfaces"`
	Description                       types.String `tfsdk:"description"`
	AutoRepairSolution                types.Bool   `tfsdk:"auto_repair_solution"`
	StayInArea                        types.Bool   `tfsdk:"stay_in_area"`
	AdjacencyHopType                  types.String `tfsdk:"adjacency_hop_type"`
	OperationMode                     types.String `tfsdk:"operation_mode"`
	Objective                         types.String `tfsdk:"objective"`
	PceDispatchTimeout                types.Int64  `tfsdk:"pce_dispatch_timeout"`
	CongestionCheckSuspensionInterval types.Int64  `tfsdk:"congestion_check_suspension_interval"`
	OverprovisionFactor               types.Int64  `tfsdk:"overprovision_factor"`
	UnevenEcmpTrafficThreshold        types.Int64  `tfsdk:"uneven_ecmp_traffic_threshold"`
	ThrottleModeThreshold             types.Int64  `tfsdk:"throttle_mode_threshold"`
	HistoryRetentionTime              types.Int64  `tfsdk:"history_retention_time"`
	DebugOpt                          types.Bool   `tfsdk:"debug_opt"`
	DebugOptMaxFiles                  types.Int64  `tfsdk:"debug_opt_max_files"`
	AffinityIncludeAny                types.List   `tfsdk:"affinity_include_any"`
	AffinityExcludeAny                types.List   `tfsdk:"affinity_exclude_any"`
	AffinityIncludeAll                types.List   `tfsdk:"affinity_include_all"`
	NodeTagMaxSidDepth                types.Map    `tfsdk:"node_tag_max_sid_depth"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data LocalCongestionMitigation) getPath() string {
	return "/crosswork/optima/api/v1/functionPacks/lcm"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

func (data LocalCongestionMitigation) getPathDelete() string {
	return "/crosswork/optima/api/v1/functionPacks/lcm"
}

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data LocalCongestionMitigation) getPathGet() string {
	return "/crosswork/optima/api/v1/functionPacks/lcm/domain/config"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// Escapes the id so it cannot change the structure of a request body.
func (data LocalCongestionMitigation) idJSON() string {
	id := data.Id.ValueString()
	if id == "" {
		id = fmt.Sprint(data.DomainId.ValueString())
	}
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}

func (data LocalCongestionMitigation) readGet(c *cw.Client, path string) (cw.Res, error) {
	return c.Post(path, strings.ReplaceAll(`{"id":"{id}","default":false}`, "{id}", data.idJSON()), cw.ReadOnly)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data LocalCongestionMitigation) toBody(ctx context.Context, state LocalCongestionMitigation) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.DomainId.IsNull() {
		body, _ = sjson.Set(body, "id", data.DomainId.ValueString())
	}
	if !data.Enable.IsNull() {
		body, _ = sjson.Set(body, "config.enable", data.Enable.ValueBool())
	}
	if !data.Color.IsNull() {
		body, _ = sjson.Set(body, "config.color", data.Color.ValueInt64())
	}
	if !data.UtilizationThreshold.IsNull() {
		body, _ = sjson.Set(body, "config.utilThreshold", data.UtilizationThreshold.ValueInt64())
	}
	if !data.UtilizationHoldMargin.IsNull() {
		body, _ = sjson.Set(body, "config.utilHoldMargin", data.UtilizationHoldMargin.ValueInt64())
	}
	if !data.DeleteLsps.IsNull() {
		body, _ = sjson.Set(body, "config.delLsps", data.DeleteLsps.ValueBool())
	}
	if !data.ProfileId.IsNull() {
		body, _ = sjson.Set(body, "config.profileId", data.ProfileId.ValueInt64())
	}
	if !data.CongestionCheckInterval.IsNull() {
		body, _ = sjson.Set(body, "config.congestionCheckInterval", data.CongestionCheckInterval.ValueInt64())
	}
	if !data.MaxTteGroup.IsNull() {
		body, _ = sjson.Set(body, "config.maxTteGroup", data.MaxTteGroup.ValueInt64())
	}
	if !data.IncludeAllInterfaces.IsNull() {
		body, _ = sjson.Set(body, "config.includeAllInterfaces", data.IncludeAllInterfaces.ValueBool())
	}
	if !data.Description.IsNull() {
		body, _ = sjson.Set(body, "config.description", data.Description.ValueString())
	}
	if !data.AutoRepairSolution.IsNull() {
		body, _ = sjson.Set(body, "config.autoRepairSolution", data.AutoRepairSolution.ValueBool())
	}
	if !data.StayInArea.IsNull() {
		body, _ = sjson.Set(body, "config.stayInArea", data.StayInArea.ValueBool())
	}
	if !data.AdjacencyHopType.IsNull() {
		body, _ = sjson.Set(body, "config.adjacencyHopType", data.AdjacencyHopType.ValueString())
	}
	if !data.OperationMode.IsNull() {
		body, _ = sjson.Set(body, "config.operationMode", data.OperationMode.ValueString())
	}
	if !data.Objective.IsNull() {
		body, _ = sjson.Set(body, "config.objective", data.Objective.ValueString())
	}
	if !data.PceDispatchTimeout.IsNull() {
		body, _ = sjson.Set(body, "config.pceDispatchTimeout", data.PceDispatchTimeout.ValueInt64())
	}
	if !data.CongestionCheckSuspensionInterval.IsNull() {
		body, _ = sjson.Set(body, "config.congestionCheckSuspensionInterval", data.CongestionCheckSuspensionInterval.ValueInt64())
	}
	if !data.OverprovisionFactor.IsNull() {
		body, _ = sjson.Set(body, "config.overprovisionFactor", data.OverprovisionFactor.ValueInt64())
	}
	if !data.UnevenEcmpTrafficThreshold.IsNull() {
		body, _ = sjson.Set(body, "config.unevenEcmpTrafficThreshold", data.UnevenEcmpTrafficThreshold.ValueInt64())
	}
	if !data.ThrottleModeThreshold.IsNull() {
		body, _ = sjson.Set(body, "config.throttleModeThreshold", data.ThrottleModeThreshold.ValueInt64())
	}
	if !data.HistoryRetentionTime.IsNull() {
		body, _ = sjson.Set(body, "config.historyRetentionTime", data.HistoryRetentionTime.ValueInt64())
	}
	if !data.DebugOpt.IsNull() {
		body, _ = sjson.Set(body, "config.debugOpt", data.DebugOpt.ValueBool())
	}
	if !data.DebugOptMaxFiles.IsNull() {
		body, _ = sjson.Set(body, "config.debugOptMaxFiles", data.DebugOptMaxFiles.ValueInt64())
	}
	if !data.AffinityIncludeAny.IsNull() {
		var values []string
		data.AffinityIncludeAny.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "config.affinityConstraint.includeAny", values)
	}
	if !data.AffinityExcludeAny.IsNull() {
		var values []string
		data.AffinityExcludeAny.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "config.affinityConstraint.excludeAny", values)
	}
	if !data.AffinityIncludeAll.IsNull() {
		var values []string
		data.AffinityIncludeAll.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "config.affinityConstraint.includeAll", values)
	}
	if !data.NodeTagMaxSidDepth.IsNull() {
		var values map[string]string
		data.NodeTagMaxSidDepth.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "config.nodeTagMaxSidDepth.0", values)
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *LocalCongestionMitigation) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("data.config.enable"); value.Exists() {
		data.Enable = types.BoolValue(value.Bool())
	} else {
		data.Enable = types.BoolValue(false)
	}
	if value := res.Get("data.config.color"); value.Exists() {
		data.Color = types.Int64Value(value.Int())
	} else {
		data.Color = types.Int64Value(2000)
	}
	if value := res.Get("data.config.utilThreshold"); value.Exists() {
		data.UtilizationThreshold = types.Int64Value(value.Int())
	} else {
		data.UtilizationThreshold = types.Int64Value(100)
	}
	if value := res.Get("data.config.utilHoldMargin"); value.Exists() {
		data.UtilizationHoldMargin = types.Int64Value(value.Int())
	} else {
		data.UtilizationHoldMargin = types.Int64Value(5)
	}
	if value := res.Get("data.config.delLsps"); value.Exists() {
		data.DeleteLsps = types.BoolValue(value.Bool())
	} else {
		data.DeleteLsps = types.BoolValue(false)
	}
	if value := res.Get("data.config.profileId"); value.Exists() {
		data.ProfileId = types.Int64Value(value.Int())
	} else {
		data.ProfileId = types.Int64Value(0)
	}
	if value := res.Get("data.config.congestionCheckInterval"); value.Exists() {
		data.CongestionCheckInterval = types.Int64Value(value.Int())
	} else {
		data.CongestionCheckInterval = types.Int64Value(900)
	}
	if value := res.Get("data.config.maxTteGroup"); value.Exists() {
		data.MaxTteGroup = types.Int64Value(value.Int())
	} else {
		data.MaxTteGroup = types.Int64Value(8)
	}
	if value := res.Get("data.config.includeAllInterfaces"); value.Exists() {
		data.IncludeAllInterfaces = types.BoolValue(value.Bool())
	} else {
		data.IncludeAllInterfaces = types.BoolValue(false)
	}
	if value := res.Get("data.config.description"); value.Exists() {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringValue("LCM startup config")
	}
	if value := res.Get("data.config.autoRepairSolution"); value.Exists() {
		data.AutoRepairSolution = types.BoolValue(value.Bool())
	} else {
		data.AutoRepairSolution = types.BoolValue(false)
	}
	if value := res.Get("data.config.stayInArea"); value.Exists() {
		data.StayInArea = types.BoolValue(value.Bool())
	} else {
		data.StayInArea = types.BoolValue(false)
	}
	if value := res.Get("data.config.adjacencyHopType"); value.Exists() {
		data.AdjacencyHopType = types.StringValue(value.String())
	} else {
		data.AdjacencyHopType = types.StringValue("UNPROTECTED")
	}
	if value := res.Get("data.config.operationMode"); value.Exists() {
		data.OperationMode = types.StringValue(value.String())
	} else {
		data.OperationMode = types.StringValue("MANUAL")
	}
	if value := res.Get("data.config.objective"); value.Exists() {
		data.Objective = types.StringValue(value.String())
	} else {
		data.Objective = types.StringValue("MIN_IGP_METRIC")
	}
	if value := res.Get("data.config.pceDispatchTimeout"); value.Exists() {
		data.PceDispatchTimeout = types.Int64Value(value.Int())
	} else {
		data.PceDispatchTimeout = types.Int64Value(180)
	}
	if value := res.Get("data.config.congestionCheckSuspensionInterval"); value.Exists() {
		data.CongestionCheckSuspensionInterval = types.Int64Value(value.Int())
	} else {
		data.CongestionCheckSuspensionInterval = types.Int64Value(600)
	}
	if value := res.Get("data.config.overprovisionFactor"); value.Exists() {
		data.OverprovisionFactor = types.Int64Value(value.Int())
	} else {
		data.OverprovisionFactor = types.Int64Value(0)
	}
	if value := res.Get("data.config.unevenEcmpTrafficThreshold"); value.Exists() {
		data.UnevenEcmpTrafficThreshold = types.Int64Value(value.Int())
	} else {
		data.UnevenEcmpTrafficThreshold = types.Int64Value(0)
	}
	if value := res.Get("data.config.throttleModeThreshold"); value.Exists() {
		data.ThrottleModeThreshold = types.Int64Value(value.Int())
	} else {
		data.ThrottleModeThreshold = types.Int64Value(5)
	}
	if value := res.Get("data.config.historyRetentionTime"); value.Exists() {
		data.HistoryRetentionTime = types.Int64Value(value.Int())
	} else {
		data.HistoryRetentionTime = types.Int64Value(30)
	}
	if value := res.Get("data.config.debugOpt"); value.Exists() {
		data.DebugOpt = types.BoolValue(value.Bool())
	} else {
		data.DebugOpt = types.BoolValue(false)
	}
	if value := res.Get("data.config.debugOptMaxFiles"); value.Exists() {
		data.DebugOptMaxFiles = types.Int64Value(value.Int())
	} else {
		data.DebugOptMaxFiles = types.Int64Value(40)
	}
	if value := res.Get("data.config.affinityConstraint.includeAny.#.bitPosition"); value.Exists() && len(value.Array()) > 0 {
		data.AffinityIncludeAny = helpers.GetStringList(value.Array())
	} else {
		data.AffinityIncludeAny = types.ListNull(types.StringType)
	}
	if value := res.Get("data.config.affinityConstraint.excludeAny.#.bitPosition"); value.Exists() && len(value.Array()) > 0 {
		data.AffinityExcludeAny = helpers.GetStringList(value.Array())
	} else {
		data.AffinityExcludeAny = types.ListNull(types.StringType)
	}
	if value := res.Get("data.config.affinityConstraint.includeAll.#.bitPosition"); value.Exists() && len(value.Array()) > 0 {
		data.AffinityIncludeAll = helpers.GetStringList(value.Array())
	} else {
		data.AffinityIncludeAll = types.ListNull(types.StringType)
	}
	if value := res.Get("data.config.nodeTagMaxSidDepth|@join"); value.Exists() && len(value.Map()) > 0 {
		data.NodeTagMaxSidDepth = helpers.GetStringMap(value.Map())
	} else {
		data.NodeTagMaxSidDepth = types.MapNull(types.StringType)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *LocalCongestionMitigation) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("data.config.enable"); value.Exists() && !data.Enable.IsNull() {
		data.Enable = types.BoolValue(value.Bool())
	} else if data.Enable.ValueBool() != false {
		data.Enable = types.BoolNull()
	}
	if value := res.Get("data.config.color"); value.Exists() && !data.Color.IsNull() {
		data.Color = types.Int64Value(value.Int())
	} else if data.Color.ValueInt64() != 2000 {
		data.Color = types.Int64Null()
	}
	if value := res.Get("data.config.utilThreshold"); value.Exists() && !data.UtilizationThreshold.IsNull() {
		data.UtilizationThreshold = types.Int64Value(value.Int())
	} else if data.UtilizationThreshold.ValueInt64() != 100 {
		data.UtilizationThreshold = types.Int64Null()
	}
	if value := res.Get("data.config.utilHoldMargin"); value.Exists() && !data.UtilizationHoldMargin.IsNull() {
		data.UtilizationHoldMargin = types.Int64Value(value.Int())
	} else if data.UtilizationHoldMargin.ValueInt64() != 5 {
		data.UtilizationHoldMargin = types.Int64Null()
	}
	if value := res.Get("data.config.delLsps"); value.Exists() && !data.DeleteLsps.IsNull() {
		data.DeleteLsps = types.BoolValue(value.Bool())
	} else if data.DeleteLsps.ValueBool() != false {
		data.DeleteLsps = types.BoolNull()
	}
	if value := res.Get("data.config.profileId"); value.Exists() && !data.ProfileId.IsNull() {
		data.ProfileId = types.Int64Value(value.Int())
	} else if data.ProfileId.ValueInt64() != 0 {
		data.ProfileId = types.Int64Null()
	}
	if value := res.Get("data.config.congestionCheckInterval"); value.Exists() && !data.CongestionCheckInterval.IsNull() {
		data.CongestionCheckInterval = types.Int64Value(value.Int())
	} else if data.CongestionCheckInterval.ValueInt64() != 900 {
		data.CongestionCheckInterval = types.Int64Null()
	}
	if value := res.Get("data.config.maxTteGroup"); value.Exists() && !data.MaxTteGroup.IsNull() {
		data.MaxTteGroup = types.Int64Value(value.Int())
	} else if data.MaxTteGroup.ValueInt64() != 8 {
		data.MaxTteGroup = types.Int64Null()
	}
	if value := res.Get("data.config.includeAllInterfaces"); value.Exists() && !data.IncludeAllInterfaces.IsNull() {
		data.IncludeAllInterfaces = types.BoolValue(value.Bool())
	} else if data.IncludeAllInterfaces.ValueBool() != false {
		data.IncludeAllInterfaces = types.BoolNull()
	}
	if value := res.Get("data.config.description"); value.Exists() && !data.Description.IsNull() {
		data.Description = types.StringValue(value.String())
	} else if data.Description.ValueString() != "LCM startup config" {
		data.Description = types.StringNull()
	}
	if value := res.Get("data.config.autoRepairSolution"); value.Exists() && !data.AutoRepairSolution.IsNull() {
		data.AutoRepairSolution = types.BoolValue(value.Bool())
	} else if data.AutoRepairSolution.ValueBool() != false {
		data.AutoRepairSolution = types.BoolNull()
	}
	if value := res.Get("data.config.stayInArea"); value.Exists() && !data.StayInArea.IsNull() {
		data.StayInArea = types.BoolValue(value.Bool())
	} else if data.StayInArea.ValueBool() != false {
		data.StayInArea = types.BoolNull()
	}
	if value := res.Get("data.config.adjacencyHopType"); value.Exists() && !data.AdjacencyHopType.IsNull() {
		data.AdjacencyHopType = types.StringValue(value.String())
	} else if data.AdjacencyHopType.ValueString() != "UNPROTECTED" {
		data.AdjacencyHopType = types.StringNull()
	}
	if value := res.Get("data.config.operationMode"); value.Exists() && !data.OperationMode.IsNull() {
		data.OperationMode = types.StringValue(value.String())
	} else if data.OperationMode.ValueString() != "MANUAL" {
		data.OperationMode = types.StringNull()
	}
	if value := res.Get("data.config.objective"); value.Exists() && !data.Objective.IsNull() {
		data.Objective = types.StringValue(value.String())
	} else if data.Objective.ValueString() != "MIN_IGP_METRIC" {
		data.Objective = types.StringNull()
	}
	if value := res.Get("data.config.pceDispatchTimeout"); value.Exists() && !data.PceDispatchTimeout.IsNull() {
		data.PceDispatchTimeout = types.Int64Value(value.Int())
	} else if data.PceDispatchTimeout.ValueInt64() != 180 {
		data.PceDispatchTimeout = types.Int64Null()
	}
	if value := res.Get("data.config.congestionCheckSuspensionInterval"); value.Exists() && !data.CongestionCheckSuspensionInterval.IsNull() {
		data.CongestionCheckSuspensionInterval = types.Int64Value(value.Int())
	} else if data.CongestionCheckSuspensionInterval.ValueInt64() != 600 {
		data.CongestionCheckSuspensionInterval = types.Int64Null()
	}
	if value := res.Get("data.config.overprovisionFactor"); value.Exists() && !data.OverprovisionFactor.IsNull() {
		data.OverprovisionFactor = types.Int64Value(value.Int())
	} else if data.OverprovisionFactor.ValueInt64() != 0 {
		data.OverprovisionFactor = types.Int64Null()
	}
	if value := res.Get("data.config.unevenEcmpTrafficThreshold"); value.Exists() && !data.UnevenEcmpTrafficThreshold.IsNull() {
		data.UnevenEcmpTrafficThreshold = types.Int64Value(value.Int())
	} else if data.UnevenEcmpTrafficThreshold.ValueInt64() != 0 {
		data.UnevenEcmpTrafficThreshold = types.Int64Null()
	}
	if value := res.Get("data.config.throttleModeThreshold"); value.Exists() && !data.ThrottleModeThreshold.IsNull() {
		data.ThrottleModeThreshold = types.Int64Value(value.Int())
	} else if data.ThrottleModeThreshold.ValueInt64() != 5 {
		data.ThrottleModeThreshold = types.Int64Null()
	}
	if value := res.Get("data.config.historyRetentionTime"); value.Exists() && !data.HistoryRetentionTime.IsNull() {
		data.HistoryRetentionTime = types.Int64Value(value.Int())
	} else if data.HistoryRetentionTime.ValueInt64() != 30 {
		data.HistoryRetentionTime = types.Int64Null()
	}
	if value := res.Get("data.config.debugOpt"); value.Exists() && !data.DebugOpt.IsNull() {
		data.DebugOpt = types.BoolValue(value.Bool())
	} else if data.DebugOpt.ValueBool() != false {
		data.DebugOpt = types.BoolNull()
	}
	if value := res.Get("data.config.debugOptMaxFiles"); value.Exists() && !data.DebugOptMaxFiles.IsNull() {
		data.DebugOptMaxFiles = types.Int64Value(value.Int())
	} else if data.DebugOptMaxFiles.ValueInt64() != 40 {
		data.DebugOptMaxFiles = types.Int64Null()
	}
	if value := res.Get("data.config.affinityConstraint.includeAny.#.bitPosition"); value.Exists() && !data.AffinityIncludeAny.IsNull() {
		data.AffinityIncludeAny = helpers.GetStringList(value.Array())
	} else {
		data.AffinityIncludeAny = types.ListNull(types.StringType)
	}
	if value := res.Get("data.config.affinityConstraint.excludeAny.#.bitPosition"); value.Exists() && !data.AffinityExcludeAny.IsNull() {
		data.AffinityExcludeAny = helpers.GetStringList(value.Array())
	} else {
		data.AffinityExcludeAny = types.ListNull(types.StringType)
	}
	if value := res.Get("data.config.affinityConstraint.includeAll.#.bitPosition"); value.Exists() && !data.AffinityIncludeAll.IsNull() {
		data.AffinityIncludeAll = helpers.GetStringList(value.Array())
	} else {
		data.AffinityIncludeAll = types.ListNull(types.StringType)
	}
	if value := res.Get("data.config.nodeTagMaxSidDepth|@join"); value.Exists() && !data.NodeTagMaxSidDepth.IsNull() {
		data.NodeTagMaxSidDepth = helpers.GetStringMap(value.Map())
	} else {
		data.NodeTagMaxSidDepth = types.MapNull(types.StringType)
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *LocalCongestionMitigation) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Enable.IsNull() {
		return false
	}
	if !data.Color.IsNull() {
		return false
	}
	if !data.UtilizationThreshold.IsNull() {
		return false
	}
	if !data.UtilizationHoldMargin.IsNull() {
		return false
	}
	if !data.DeleteLsps.IsNull() {
		return false
	}
	if !data.ProfileId.IsNull() {
		return false
	}
	if !data.CongestionCheckInterval.IsNull() {
		return false
	}
	if !data.MaxTteGroup.IsNull() {
		return false
	}
	if !data.IncludeAllInterfaces.IsNull() {
		return false
	}
	if !data.Description.IsNull() {
		return false
	}
	if !data.AutoRepairSolution.IsNull() {
		return false
	}
	if !data.StayInArea.IsNull() {
		return false
	}
	if !data.AdjacencyHopType.IsNull() {
		return false
	}
	if !data.OperationMode.IsNull() {
		return false
	}
	if !data.Objective.IsNull() {
		return false
	}
	if !data.PceDispatchTimeout.IsNull() {
		return false
	}
	if !data.CongestionCheckSuspensionInterval.IsNull() {
		return false
	}
	if !data.OverprovisionFactor.IsNull() {
		return false
	}
	if !data.UnevenEcmpTrafficThreshold.IsNull() {
		return false
	}
	if !data.ThrottleModeThreshold.IsNull() {
		return false
	}
	if !data.HistoryRetentionTime.IsNull() {
		return false
	}
	if !data.DebugOpt.IsNull() {
		return false
	}
	if !data.DebugOptMaxFiles.IsNull() {
		return false
	}
	if !data.AffinityIncludeAny.IsNull() {
		return false
	}
	if !data.AffinityExcludeAny.IsNull() {
		return false
	}
	if !data.AffinityIncludeAll.IsNull() {
		return false
	}
	if !data.NodeTagMaxSidDepth.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type DemandMatrixGlobalConfiguration struct {
	Instance             types.String  `tfsdk:"instance"`
	Id                   types.String  `tfsdk:"id"`
	ComputationCadence   types.Int64   `tfsdk:"computation_cadence"`
	MinimumDemandTraffic types.Float64 `tfsdk:"minimum_demand_traffic"`
	DebugMode            types.Bool    `tfsdk:"debug_mode"`
	DebugMaxPlanFiles    types.Int64   `tfsdk:"debug_max_plan_files"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data DemandMatrixGlobalConfiguration) getPath() string {
	return "/crosswork/optima/api/v1/setDdmGlobalConfigs"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data DemandMatrixGlobalConfiguration) getPathGet() string {
	return "/crosswork/optima/api/v1/getDdmGlobalConfigs"
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
func (data DemandMatrixGlobalConfiguration) idJSON() string {
	id := data.Id.ValueString()
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}

func (data DemandMatrixGlobalConfiguration) readGet(c *cw.Client, path string) (cw.Res, error) {
	return c.Post(path, strings.ReplaceAll(`{}`, "{id}", data.idJSON()), cw.ReadOnly)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data DemandMatrixGlobalConfiguration) toBody(ctx context.Context, state DemandMatrixGlobalConfiguration) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.ComputationCadence.IsNull() {
		body, _ = sjson.Set(body, "configs.computeCadence", data.ComputationCadence.ValueInt64())
	}
	if !data.MinimumDemandTraffic.IsNull() {
		body, _ = sjson.Set(body, "configs.minDemandTraffic", data.MinimumDemandTraffic.ValueFloat64())
	}
	if !data.DebugMode.IsNull() {
		body, _ = sjson.Set(body, "configs.debug", data.DebugMode.ValueBool())
	}
	if !data.DebugMaxPlanFiles.IsNull() {
		body, _ = sjson.Set(body, "configs.maxPlanFiles", data.DebugMaxPlanFiles.ValueInt64())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *DemandMatrixGlobalConfiguration) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("demand_matrix_global_configuration")
	if value := res.Get("configs.computeCadence"); value.Exists() {
		data.ComputationCadence = types.Int64Value(value.Int())
	} else {
		data.ComputationCadence = types.Int64Value(15)
	}
	if value := res.Get("configs.minDemandTraffic"); value.Exists() {
		data.MinimumDemandTraffic = types.Float64Value(value.Float())
	} else {
		data.MinimumDemandTraffic = types.Float64Value(0)
	}
	if value := res.Get("configs.debug"); value.Exists() {
		data.DebugMode = types.BoolValue(value.Bool())
	} else {
		data.DebugMode = types.BoolValue(false)
	}
	if value := res.Get("configs.maxPlanFiles"); value.Exists() {
		data.DebugMaxPlanFiles = types.Int64Value(value.Int())
	} else {
		data.DebugMaxPlanFiles = types.Int64Value(1)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *DemandMatrixGlobalConfiguration) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("configs.computeCadence"); value.Exists() && !data.ComputationCadence.IsNull() {
		data.ComputationCadence = types.Int64Value(value.Int())
	} else if data.ComputationCadence.ValueInt64() != 15 {
		data.ComputationCadence = types.Int64Null()
	}
	if value := res.Get("configs.minDemandTraffic"); value.Exists() && !data.MinimumDemandTraffic.IsNull() {
		data.MinimumDemandTraffic = types.Float64Value(value.Float())
	} else if data.MinimumDemandTraffic.ValueFloat64() != 0 {
		data.MinimumDemandTraffic = types.Float64Null()
	}
	if value := res.Get("configs.debug"); value.Exists() && !data.DebugMode.IsNull() {
		data.DebugMode = types.BoolValue(value.Bool())
	} else if data.DebugMode.ValueBool() != false {
		data.DebugMode = types.BoolNull()
	}
	if value := res.Get("configs.maxPlanFiles"); value.Exists() && !data.DebugMaxPlanFiles.IsNull() {
		data.DebugMaxPlanFiles = types.Int64Value(value.Int())
	} else if data.DebugMaxPlanFiles.ValueInt64() != 1 {
		data.DebugMaxPlanFiles = types.Int64Null()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *DemandMatrixGlobalConfiguration) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.ComputationCadence.IsNull() {
		return false
	}
	if !data.MinimumDemandTraffic.IsNull() {
		return false
	}
	if !data.DebugMode.IsNull() {
		return false
	}
	if !data.DebugMaxPlanFiles.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

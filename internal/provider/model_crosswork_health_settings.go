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
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type HealthSettings struct {
	Instance types.String                  `tfsdk:"instance"`
	Id       types.String                  `tfsdk:"id"`
	Rows     map[string]HealthSettingsRows `tfsdk:"rows"`
}

type HealthSettingsRows struct {
	Unit       types.String                   `tfsdk:"unit"`
	Categories []HealthSettingsRowsCategories `tfsdk:"categories"`
}

type HealthSettingsRowsCategories struct {
	Level  types.String  `tfsdk:"level"`
	Min    types.Float64 `tfsdk:"min"`
	Max    types.Float64 `tfsdk:"max"`
	Values types.List    `tfsdk:"values"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data HealthSettings) getPath() string {
	return "/crosswork/performance/v1/dashboards/healthsettings"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data HealthSettings) getPathGet() string {
	return "/crosswork/performance/v1/dashboards/healthsettings/all"
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

// Written by hand: a write reads the rows of the controller first (healthWriteBody), so the model has no body of its own.
func (data HealthSettings) toBody(ctx context.Context, state HealthSettings) string {
	return "{}"
}

// End of section. //template:end toBody

// Written by hand: every row of the answer is read (a data source).
func (data *HealthSettings) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("health_settings")
	data.Rows = make(map[string]HealthSettingsRows)
	res.ForEach(func(_, rows gjson.Result) bool {
		rows.ForEach(func(key, row gjson.Result) bool {
			data.Rows[key.String()] = healthRowFromBody(row, true)
			return true
		})
		return true
	})
}

// End of section. //template:end fromBody

// Written by hand: only the rows of the configuration are read back, and the unit only for a row that names it.
func (data *HealthSettings) updateFromBody(ctx context.Context, res gjson.Result) {
	for key, planned := range data.Rows {
		if _, row, found := healthRow(res, key); found {
			data.Rows[key] = healthRowFromBody(row, !planned.Unit.IsNull())
		} else {
			delete(data.Rows, key)
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *HealthSettings) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.Rows) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

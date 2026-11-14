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
type DataRetentionDevicePerformance struct {
	Instance types.String                                  `tfsdk:"instance"`
	Id       types.String                                  `tfsdk:"id"`
	Rows     map[string]DataRetentionDevicePerformanceRows `tfsdk:"rows"`
}

type DataRetentionDevicePerformanceRows struct {
	RawHours    types.Int64 `tfsdk:"raw_hours"`
	HourlyHours types.Int64 `tfsdk:"hourly_hours"`
	DailyHours  types.Int64 `tfsdk:"daily_hours"`
	WeeklyHours types.Int64 `tfsdk:"weekly_hours"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data DataRetentionDevicePerformance) getPath() string {
	return "/crosswork/performance/v1/dataretention"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

func (data DataRetentionDevicePerformance) getPathDelete() string {
	return "/crosswork/performance/v1/dataretention/reset"
}

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data DataRetentionDevicePerformance) getPathGet() string {
	return "/crosswork/performance/v1/dataretention/all"
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

// Written by hand: the controller reads and writes one JSON object keyed by the policy type, at the root of the body.
func (data DataRetentionDevicePerformance) toBody(ctx context.Context, state DataRetentionDevicePerformance) string {
	body := "{}"
	for key, row := range data.Rows {
		item := ""
		item, _ = sjson.Set(item, "rawDataRetentionPeriod", row.RawHours.ValueInt64())
		item, _ = sjson.Set(item, "hourlyDataRetentionPeriod", row.HourlyHours.ValueInt64())
		item, _ = sjson.Set(item, "dailyDataRetentionPeriod", row.DailyHours.ValueInt64())
		item, _ = sjson.Set(item, "weeklyDataRetentionPeriod", row.WeeklyHours.ValueInt64())
		body, _ = sjson.SetRaw(body, escapeJSONKey(key), item)
	}
	return body
}

// End of section. //template:end toBody

// Written by hand: every policy type of the answer becomes a row (a data source and an import read all of them).
func (data *DataRetentionDevicePerformance) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("data_retention_device_performance")
	data.Rows = make(map[string]DataRetentionDevicePerformanceRows)
	res.ForEach(func(key, value gjson.Result) bool {
		if value.IsObject() {
			data.Rows[key.String()] = retentionRowFromBody(value)
		}
		return true
	})
}

// End of section. //template:end fromBody

// Written by hand: only the rows of the configuration are read back. A policy type that the controller does not have is dropped.
func (data *DataRetentionDevicePerformance) updateFromBody(ctx context.Context, res gjson.Result) {
	for key := range data.Rows {
		if value := res.Get(escapeJSONKey(key)); value.IsObject() {
			data.Rows[key] = retentionRowFromBody(value)
		} else {
			delete(data.Rows, key)
		}
	}
}

func retentionRowFromBody(value gjson.Result) DataRetentionDevicePerformanceRows {
	return DataRetentionDevicePerformanceRows{
		RawHours:    types.Int64Value(value.Get("rawDataRetentionPeriod").Int()),
		HourlyHours: types.Int64Value(value.Get("hourlyDataRetentionPeriod").Int()),
		DailyHours:  types.Int64Value(value.Get("dailyDataRetentionPeriod").Int()),
		WeeklyHours: types.Int64Value(value.Get("weeklyDataRetentionPeriod").Int()),
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *DataRetentionDevicePerformance) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.Rows) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

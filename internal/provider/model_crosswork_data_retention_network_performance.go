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
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type DataRetentionNetworkPerformance struct {
	Instance              types.String `tfsdk:"instance"`
	Id                    types.String `tfsdk:"id"`
	CollectLspStateChange types.Bool   `tfsdk:"collect_lsp_state_change"`
	CollectLspPathChange  types.Bool   `tfsdk:"collect_lsp_path_change"`
	CollectLspUtilization types.Bool   `tfsdk:"collect_lsp_utilization"`
	CollectLspPm          types.Bool   `tfsdk:"collect_lsp_pm"`
	CollectLinkPm         types.Bool   `tfsdk:"collect_link_pm"`
	RawHours              types.Int64  `tfsdk:"raw_hours"`
	HourlyDays            types.Int64  `tfsdk:"hourly_days"`
	DailyDays             types.Int64  `tfsdk:"daily_days"`
	WeeklyWeeks           types.Int64  `tfsdk:"weekly_weeks"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data DataRetentionNetworkPerformance) getPath() string {
	return "/crosswork/optima-analytics/api/v1/config/settings"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data DataRetentionNetworkPerformance) getPathGet() string {
	return "/crosswork/optima-analytics/api/v1/config/settings"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Written by hand: the settings are the collection settings and the retention periods, two calls. The answer is one object with both.
func (data DataRetentionNetworkPerformance) readGet(c *cw.Client, path string) (cw.Res, error) {
	settings, err := c.Get(networkSettingsPath)
	if err != nil {
		return settings, err
	}
	periods, err := c.Get(networkRetentionPath)
	if err != nil {
		return periods, err
	}
	body, _ := sjson.SetRaw("{}", "settings", settings.Raw)
	retention := "{}"
	periods.Get("data").ForEach(func(_, v gjson.Result) bool {
		retention, _ = sjson.Set(retention, v.Get("type").String(), v.Get("value").Int())
		return true
	})
	body, _ = sjson.SetRaw(body, "retention", retention)
	return gjson.Parse(body), nil
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

func (data DataRetentionNetworkPerformance) toBody(ctx context.Context, state DataRetentionNetworkPerformance) string {
	link := flagString(data.CollectLinkPm.ValueBool())
	return networkSettingsBody(map[string]string{
		"lsp_state_change_collection": flagString(data.CollectLspStateChange.ValueBool()),
		"lsp_path_change_collection":  flagString(data.CollectLspPathChange.ValueBool()),
		"lsp_metrics_collection":      flagString(data.CollectLspUtilization.ValueBool()),
		"lsp_pm_enabled":              flagString(data.CollectLspPm.ValueBool()),
		"link_pm":                     link,
		"interface_pm_enabled":        link,
	})
}

// End of section. //template:end toBody

// Written by hand: every value is read (a data source).
func (data *DataRetentionNetworkPerformance) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("data_retention_network_performance")
	data.readValues(res)
}

func (data *DataRetentionNetworkPerformance) readValues(res gjson.Result) {
	on := func(key string) types.Bool {
		value := res.Get("settings." + key).String()
		return types.BoolValue(value == "1" || value == "true")
	}
	data.CollectLspStateChange = on("lsp_state_change_collection")
	data.CollectLspPathChange = on("lsp_path_change_collection")
	data.CollectLspUtilization = on("lsp_metrics_collection")
	data.CollectLspPm = on("lsp_pm_enabled")
	data.CollectLinkPm = on("link_pm")
	data.RawHours = types.Int64Value(res.Get("retention.RAW").Int())
	data.HourlyDays = types.Int64Value(res.Get("retention.HOURLY").Int())
	data.DailyDays = types.Int64Value(res.Get("retention.DAILY").Int())
	data.WeeklyWeeks = types.Int64Value(res.Get("retention.WEEKLY").Int())
}

// End of section. //template:end fromBody

// Written by hand: the resource manages every value, so a read replaces all of them.
func (data *DataRetentionNetworkPerformance) updateFromBody(ctx context.Context, res gjson.Result) {
	data.readValues(res)
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *DataRetentionNetworkPerformance) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.CollectLspStateChange.IsNull() {
		return false
	}
	if !data.CollectLspPathChange.IsNull() {
		return false
	}
	if !data.CollectLspUtilization.IsNull() {
		return false
	}
	if !data.CollectLspPm.IsNull() {
		return false
	}
	if !data.CollectLinkPm.IsNull() {
		return false
	}
	if !data.RawHours.IsNull() {
		return false
	}
	if !data.HourlyDays.IsNull() {
		return false
	}
	if !data.DailyDays.IsNull() {
		return false
	}
	if !data.WeeklyWeeks.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

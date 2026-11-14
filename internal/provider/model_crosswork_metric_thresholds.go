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
type MetricThresholds struct {
	Instance   types.String                          `tfsdk:"instance"`
	Id         types.String                          `tfsdk:"id"`
	Thresholds map[string]MetricThresholdsThresholds `tfsdk:"thresholds"`
}

type MetricThresholdsThresholds struct {
	Levels []MetricThresholdsThresholdsLevels `tfsdk:"levels"`
}

type MetricThresholdsThresholdsLevels struct {
	Severity   types.String  `tfsdk:"severity"`
	UpperValue types.Float64 `tfsdk:"upper_value"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data MetricThresholds) getPath() string {
	return "/crosswork/topology/v1/topology-service/topology/settings/metrics/set"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data MetricThresholds) getPathGet() string {
	return "/crosswork/topology/v1/topology-service/topology/settings/metrics/get"
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

// Written by hand: the thresholds are one object keyed by the metric, each metric a list of levels.
func (data MetricThresholds) toBody(ctx context.Context, state MetricThresholds) string {
	return metricThresholdsBody(data.Thresholds)
}

// End of section. //template:end toBody

// Written by hand: every metric of the answer is read (a data source).
func (data *MetricThresholds) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("metric_thresholds")
	data.Thresholds = make(map[string]MetricThresholdsThresholds)
	res.Get("thresholds").ForEach(func(name, levels gjson.Result) bool {
		data.Thresholds[name.String()] = MetricThresholdsThresholds{Levels: metricLevelsFromBody(levels)}
		return true
	})
}

// End of section. //template:end fromBody

// Written by hand: only the metrics of the configuration are read back. The controller drops a critical level with a finite value,
// so a metric whose answer differs from the configuration shows a difference.
func (data *MetricThresholds) updateFromBody(ctx context.Context, res gjson.Result) {
	for name := range data.Thresholds {
		if levels := res.Get("thresholds." + escapeJSONKey(name)); levels.IsArray() {
			data.Thresholds[name] = MetricThresholdsThresholds{Levels: metricLevelsFromBody(levels)}
		} else {
			delete(data.Thresholds, name)
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *MetricThresholds) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.Thresholds) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License Version 2.0 (the "License"),
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

func thresholdLevel(severity string, upper float64) MetricThresholdsThresholdsLevels {
	return MetricThresholdsThresholdsLevels{Severity: types.StringValue(severity), UpperValue: types.Float64Value(upper)}
}

func TestMetricThresholdsBodyHoldsOnlyTheMetricsOfThePlan(t *testing.T) {
	plan := MetricThresholds{Thresholds: map[string]MetricThresholdsThresholds{
		"Utilization": {Levels: []MetricThresholdsThresholdsLevels{thresholdLevel("HEALTHY", 25), thresholdLevel("CRITICAL", 100)}},
		"Delay":       {Levels: []MetricThresholdsThresholdsLevels{thresholdLevel("HEALTHY", 150000), thresholdLevel("CRITICAL", -1)}},
	}}
	body := gjson.Parse(plan.toBody(context.Background(), MetricThresholds{}))
	assert.Equal(t, int64(2), body.Get("thresholds.Utilization.#").Int())
	assert.Equal(t, "CRITICAL", body.Get("thresholds.Utilization.1.severity").String())
	assert.Equal(t, float64(-1), body.Get("thresholds.Delay.1.upperValue").Float())
	assert.False(t, body.Get("thresholds.Loss").Exists(), "a metric that is not in the plan is not sent")
}

func TestMetricThresholdsReadsOnlyTheMetricsOfTheConfiguration(t *testing.T) {
	answer := gjson.Parse(`{"thresholds":{"Loss":[{"severity":"HEALTHY","upperValue":0.0},{"severity":"CRITICAL","upperValue":100.0}],"Utilization":[{"severity":"HEALTHY","upperValue":25.0},{"severity":"CRITICAL","upperValue":100.0}]}}`)
	state := MetricThresholds{Thresholds: map[string]MetricThresholdsThresholds{"Utilization": {Levels: []MetricThresholdsThresholdsLevels{thresholdLevel("HEALTHY", 1)}}, "Gone": {}}}
	state.updateFromBody(context.Background(), answer)
	assert.Equal(t, 1, len(state.Thresholds), "a metric that the controller does not have is dropped")
	assert.Equal(t, 2, len(state.Thresholds["Utilization"].Levels))
	assert.Equal(t, float64(25), state.Thresholds["Utilization"].Levels[0].UpperValue.ValueFloat64())

	var data MetricThresholds
	data.fromBody(context.Background(), answer)
	assert.Equal(t, 2, len(data.Thresholds))
	assert.Equal(t, "metric_thresholds", data.Id.ValueString())
}

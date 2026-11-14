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
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cw "github.com/netascode/go-crosswork"
)

const (
	metricThresholdsSetPath      = "/crosswork/topology/v1/topology-service/topology/settings/metrics/set"
	metricThresholdsDefaultsPath = "/crosswork/topology/v1/topology-service/topology/settings/metrics/defaults"
)

// metricThresholdsBody is the body of the set call: one object keyed by the metric under `thresholds`, each metric a list of
// levels. A metric that the body leaves out keeps its levels.
func metricThresholdsBody(metrics map[string]MetricThresholdsThresholds) string {
	body := `{"thresholds":{}}`
	for name, metric := range metrics {
		levels := "[]"
		for _, level := range metric.Levels {
			item, _ := sjson.Set("", "severity", level.Severity.ValueString())
			item, _ = sjson.Set(item, "upperValue", level.UpperValue.ValueFloat64())
			levels, _ = sjson.SetRaw(levels, "-1", item)
		}
		body, _ = sjson.SetRaw(body, "thresholds."+escapeJSONKey(name), levels)
	}
	return body
}

func metricLevelsFromBody(levels gjson.Result) []MetricThresholdsThresholdsLevels {
	var out []MetricThresholdsThresholdsLevels
	levels.ForEach(func(_, level gjson.Result) bool {
		out = append(out, MetricThresholdsThresholdsLevels{Severity: types.StringValue(level.Get("severity").String()), UpperValue: types.Float64Value(level.Get("upperValue").Float())})
		return true
	})
	return out
}

func resetMetricThresholds(c *cw.Client) error {
	defaults, err := c.Get(metricThresholdsDefaultsPath)
	if err != nil {
		return fmt.Errorf("the default thresholds: %s, %s", err, defaults.String())
	}
	if !defaults.Get("thresholds").IsObject() {
		return fmt.Errorf("the controller returned no default thresholds: %s", defaults.String())
	}
	if res, err := c.Post(metricThresholdsSetPath, defaults.Raw, cw.UseMutex); err != nil {
		return fmt.Errorf("%s, %s", err, res.String())
	}
	return nil
}

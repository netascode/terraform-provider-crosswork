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
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cw "github.com/netascode/go-crosswork"
)

// The network performance data retention settings are two calls of the optima analytics service: the collection
// settings (strings "0" and "1") and the retention periods (raw in hours, hourly and daily in days, weekly in weeks).

const (
	networkSettingsPath         = "/crosswork/optima-analytics/api/v1/config/settings"
	networkDefaultSettingsPath  = "/crosswork/optima-analytics/api/v1/config/settings/defaultSettings"
	networkRetentionPath        = "/crosswork/optima-analytics/api/v1/metrics/common-retention"
	networkDefaultRetentionPath = "/crosswork/optima-analytics/api/v1/metrics/default-retention"
)

var networkRetentionTypes = []string{"RAW", "HOURLY", "DAILY", "WEEKLY"}

// The keys of the collection settings that the resource sends. The others (the utilization threshold, the retention interval,
// the link traffic rate) are left as they are.
var networkManagedSettings = []string{"lsp_state_change_collection", "lsp_path_change_collection", "lsp_metrics_collection", "lsp_pm_enabled", "link_pm", "interface_pm_enabled"}

func networkRetentionBody(values map[string]int64) string {
	body := "[]"
	for _, kind := range networkRetentionTypes {
		item, _ := sjson.Set("", "type", kind)
		item, _ = sjson.Set(item, "value", values[kind])
		body, _ = sjson.SetRaw(body, "-1", item)
	}
	return body
}

func networkSettingsBody(flags map[string]string) string {
	body := "{}"
	for _, key := range networkManagedSettings {
		if value, ok := flags[key]; ok {
			body, _ = sjson.Set(body, "settings."+key, value)
		}
	}
	return body
}

// The retention periods go first: the controller checks them and refuses a wrong one, and the settings are not checked.
func writeNetworkPerformance(c *cw.Client, flags map[string]string, retention map[string]int64) error {
	if res, err := c.Put(networkRetentionPath, networkRetentionBody(retention), cw.UseMutex); err != nil {
		return fmt.Errorf("the retention periods: %s, %s", err, res.String())
	}
	if res, err := c.Put(networkSettingsPath, networkSettingsBody(flags), cw.UseMutex); err != nil {
		return fmt.Errorf("the collection settings: %s, %s", err, res.String())
	}
	return nil
}

func applyNetworkPerformance(c *cw.Client, plan *DataRetentionNetworkPerformance) error {
	link := flagString(plan.CollectLinkPm.ValueBool())
	flags := map[string]string{
		"lsp_state_change_collection": flagString(plan.CollectLspStateChange.ValueBool()),
		"lsp_path_change_collection":  flagString(plan.CollectLspPathChange.ValueBool()),
		"lsp_metrics_collection":      flagString(plan.CollectLspUtilization.ValueBool()),
		"lsp_pm_enabled":              flagString(plan.CollectLspPm.ValueBool()),
		"link_pm":                     link,
		"interface_pm_enabled":        link,
	}
	retention := map[string]int64{"RAW": plan.RawHours.ValueInt64(), "HOURLY": plan.HourlyDays.ValueInt64(), "DAILY": plan.DailyDays.ValueInt64(), "WEEKLY": plan.WeeklyWeeks.ValueInt64()}
	if err := writeNetworkPerformance(c, flags, retention); err != nil {
		return err
	}
	read, err := plan.readGet(c, networkSettingsPath)
	if err != nil {
		return err
	}
	plan.updateFromBody(context.Background(), read)
	return nil
}

func resetNetworkPerformance(c *cw.Client) error {
	defaults, err := c.Get(networkDefaultSettingsPath)
	if err != nil {
		return fmt.Errorf("the default settings: %s, %s", err, defaults.String())
	}
	flags := map[string]string{}
	for _, key := range networkManagedSettings {
		if value := defaults.Get(key); value.Exists() {
			flags[key] = value.String()
		}
	}
	retention := map[string]int64{}
	periods, err := c.Get(networkDefaultRetentionPath)
	if err != nil {
		return fmt.Errorf("the default retention: %s, %s", err, periods.String())
	}
	periods.Get("data").ForEach(func(_, v gjson.Result) bool {
		retention[v.Get("type").String()] = v.Get("value").Int()
		return true
	})
	return writeNetworkPerformance(c, flags, retention)
}

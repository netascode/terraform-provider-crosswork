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

func TestNetworkPerformanceBodiesSendStringsAndTheInterfaceFlagWithLinkPm(t *testing.T) {
	plan := DataRetentionNetworkPerformance{CollectLspStateChange: types.BoolValue(true), CollectLspPathChange: types.BoolValue(false), CollectLspUtilization: types.BoolValue(true),
		CollectLspPm: types.BoolValue(false), CollectLinkPm: types.BoolValue(true)}
	body := gjson.Parse(plan.toBody(context.Background(), DataRetentionNetworkPerformance{}))
	assert.Equal(t, "1", body.Get("settings.lsp_state_change_collection").String())
	assert.Equal(t, "0", body.Get("settings.lsp_path_change_collection").String())
	assert.Equal(t, "1", body.Get("settings.link_pm").String())
	assert.Equal(t, "1", body.Get("settings.interface_pm_enabled").String(), "the interface flag follows Link PM")
	assert.False(t, body.Get("settings.lsp_utilization_threshold").Exists(), "a key that is not on the page is not sent")
	retention := gjson.Parse(networkRetentionBody(map[string]int64{"RAW": 12, "HOURLY": 5, "DAILY": 20, "WEEKLY": 30}))
	assert.Equal(t, int64(4), retention.Get("#").Int())
	assert.Equal(t, "WEEKLY", retention.Get("3.type").String())
	assert.Equal(t, int64(30), retention.Get("3.value").Int())
}

func TestNetworkPerformanceReadsBothCalls(t *testing.T) {
	answer := gjson.Parse(`{"settings":{"lsp_state_change_collection":"1","lsp_path_change_collection":"0","lsp_metrics_collection":"1","lsp_pm_enabled":"true","link_pm":"0","lsp_utilization_threshold":"1000"},"retention":{"RAW":24,"HOURLY":7,"DAILY":31,"WEEKLY":52}}`)
	var data DataRetentionNetworkPerformance
	data.fromBody(context.Background(), answer)
	assert.True(t, data.CollectLspStateChange.ValueBool())
	assert.False(t, data.CollectLspPathChange.ValueBool())
	assert.True(t, data.CollectLspPm.ValueBool(), "the controller keeps a value as it was written, so true reads as on")
	assert.False(t, data.CollectLinkPm.ValueBool())
	assert.Equal(t, int64(52), data.WeeklyWeeks.ValueInt64())
	assert.Equal(t, "data_retention_network_performance", data.Id.ValueString())
}

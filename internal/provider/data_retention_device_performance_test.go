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

const retentionAnswer = `{"CEPM_GNSS":{"rawDataRetentionPeriod":24,"hourlyDataRetentionPeriod":0,"dailyDataRetentionPeriod":0,"weeklyDataRetentionPeriod":0,"policyType":"GNSS","schemaName":"CEPMGNSS","hasAggrOption":false},"CEPM_INTERFACE":{"rawDataRetentionPeriod":24,"hourlyDataRetentionPeriod":168,"dailyDataRetentionPeriod":744,"weeklyDataRetentionPeriod":9072,"policyType":"INTERFACE","schemaName":"CEPMINTERFACE","hasAggrOption":true}}`

func retentionRow(raw, hourly, daily, weekly int64) DataRetentionDevicePerformanceRows {
	return DataRetentionDevicePerformanceRows{RawHours: types.Int64Value(raw), HourlyHours: types.Int64Value(hourly), DailyHours: types.Int64Value(daily), WeeklyHours: types.Int64Value(weekly)}
}

func TestRetentionBodyIsOneObjectKeyedByThePolicyType(t *testing.T) {
	plan := DataRetentionDevicePerformance{Rows: map[string]DataRetentionDevicePerformanceRows{"CEPM_INTERFACE": retentionRow(48, 100, 400, 1000)}}
	body := gjson.Parse(plan.toBody(context.Background(), DataRetentionDevicePerformance{}))
	assert.Equal(t, int64(48), body.Get("CEPM_INTERFACE.rawDataRetentionPeriod").Int())
	assert.Equal(t, int64(1000), body.Get("CEPM_INTERFACE.weeklyDataRetentionPeriod").Int())
	assert.False(t, body.Get("CEPM_GNSS").Exists(), "a row that is not in the plan is not sent")
}

func TestRetentionReadsOnlyTheRowsOfTheConfiguration(t *testing.T) {
	state := DataRetentionDevicePerformance{Rows: map[string]DataRetentionDevicePerformanceRows{"CEPM_INTERFACE": retentionRow(1, 1, 1, 1), "CEPM_GONE": retentionRow(1, 1, 1, 1)}}
	state.updateFromBody(context.Background(), gjson.Parse(retentionAnswer))
	assert.Equal(t, 1, len(state.Rows), "a policy type that the controller does not have is dropped")
	assert.Equal(t, int64(168), state.Rows["CEPM_INTERFACE"].HourlyHours.ValueInt64())
	_, exists := state.Rows["CEPM_GNSS"]
	assert.False(t, exists, "a row that the configuration does not name is not read")
}

func TestRetentionReadsEveryRowOfADataSource(t *testing.T) {
	var data DataRetentionDevicePerformance
	data.fromBody(context.Background(), gjson.Parse(retentionAnswer))
	assert.Equal(t, 2, len(data.Rows))
	assert.Equal(t, int64(0), data.Rows["CEPM_GNSS"].DailyHours.ValueInt64())
	assert.Equal(t, "data_retention_device_performance", data.Id.ValueString())
}

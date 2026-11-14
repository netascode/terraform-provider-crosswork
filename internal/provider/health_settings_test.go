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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

const healthAnswer = `{"DEVICEHEALTH":{"CPU_cpuUtilization":{"metric":"cpuUtilization","schemaName":"CPU","policy":"DEVICEHEALTH","categories":[{"level":"HEALTHY","min":0.0,"max":25.0},{"level":"CRITICAL","min":75.0,"max":100.0}],"unit":"PERCENTAGE","possibleUnits":[],"min":0,"max":100,"categoryType":"RANGE","editable":true}},"GNSS":{"CEPMGNSS_gnssModuleLockStatus":{"metric":"x","policy":"GNSS","categories":[{"level":"HEALTHY","values":[2]}],"unit":"NUMBER","categoryType":"INTEGER_VALUES","editable":false}},"PTP":{"CEPMPTP_ptpClockState":{"metric":"ptpClockState","policy":"PTP","categories":[{"level":"HEALTHY","values":[4,5]},{"level":"MAJOR","values":[2]}],"unit":"NUMBER","categoryType":"TRANSLATED_INTEGER_VALUES","editable":true}}}`

func healthCategory(level string, min, max float64) HealthSettingsRowsCategories {
	return HealthSettingsRowsCategories{Level: types.StringValue(level), Min: types.Float64Value(min), Max: types.Float64Value(max), Values: types.ListNull(types.Int64Type)}
}

func TestHealthWriteBodyKeepsTheRowAndReplacesTheCategories(t *testing.T) {
	rows := map[string]HealthSettingsRows{"CPU_cpuUtilization": {Unit: types.StringNull(), Categories: []HealthSettingsRowsCategories{healthCategory("HEALTHY", 0, 40), healthCategory("CRITICAL", 40, 100)}}}
	body, err := healthWriteBody(gjson.Parse(healthAnswer), rows)
	assert.True(t, err == nil, err)
	row := gjson.Parse(body).Get("DEVICEHEALTH.CPU_cpuUtilization")
	assert.Equal(t, "cpuUtilization", row.Get("metric").String(), "the row is sent as the controller read it, under its policy")
	assert.Equal(t, float64(40), row.Get("categories.0.max").Float())
	assert.Equal(t, "PERCENTAGE", row.Get("unit").String(), "a unit that is not named is kept")
	assert.False(t, gjson.Parse(body).Get("PTP").Exists(), "a row that is not in the plan is not sent")
}

func TestHealthWriteBodySendsValuesAndRefusesWhatCannotBeChanged(t *testing.T) {
	values, _ := types.ListValue(types.Int64Type, []attr.Value{types.Int64Value(4)})
	rows := map[string]HealthSettingsRows{"CEPMPTP_ptpClockState": {Unit: types.StringNull(), Categories: []HealthSettingsRowsCategories{{Level: types.StringValue("HEALTHY"), Min: types.Float64Null(), Max: types.Float64Null(), Values: values}}}}
	body, err := healthWriteBody(gjson.Parse(healthAnswer), rows)
	assert.True(t, err == nil, err)
	assert.Equal(t, int64(4), gjson.Parse(body).Get("PTP.CEPMPTP_ptpClockState.categories.0.values.0").Int())
	assert.False(t, gjson.Parse(body).Get("PTP.CEPMPTP_ptpClockState.categories.0.min").Exists(), "an end that is not set is not sent")
	_, err = healthWriteBody(gjson.Parse(healthAnswer), map[string]HealthSettingsRows{"CEPMGNSS_gnssModuleLockStatus": {}})
	assert.True(t, err != nil, "a row that is not editable is refused")
	assert.Contains(t, err.Error(), "cannot be changed")
	_, err = healthWriteBody(gjson.Parse(healthAnswer), map[string]HealthSettingsRows{"NOPE": {}})
	assert.True(t, err != nil, "an unknown row is refused")
	assert.Contains(t, err.Error(), "no health setting")
}

func TestHealthReadsTheRowsOfTheConfigurationAndTheUnitOnlyWhenNamed(t *testing.T) {
	state := HealthSettings{Rows: map[string]HealthSettingsRows{
		"CPU_cpuUtilization":    {Unit: types.StringNull()},
		"CEPMPTP_ptpClockState": {Unit: types.StringValue("NUMBER")},
		"GONE":                  {},
	}}
	state.updateFromBody(context.Background(), gjson.Parse(healthAnswer))
	assert.Equal(t, 2, len(state.Rows), "a row that the controller does not have is dropped")
	assert.True(t, state.Rows["CPU_cpuUtilization"].Unit.IsNull(), "no unit when the configuration names none")
	assert.Equal(t, "NUMBER", state.Rows["CEPMPTP_ptpClockState"].Unit.ValueString())
	assert.Equal(t, 2, len(state.Rows["CEPMPTP_ptpClockState"].Categories))
	assert.Equal(t, 2, len(state.Rows["CEPMPTP_ptpClockState"].Categories[0].Values.Elements()))
	assert.True(t, state.Rows["CPU_cpuUtilization"].Categories[0].Values.IsNull())

	var data HealthSettings
	data.fromBody(context.Background(), gjson.Parse(healthAnswer))
	assert.Equal(t, 3, len(data.Rows))
}

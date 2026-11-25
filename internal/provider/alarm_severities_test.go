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

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

const alarmAnswer = `{"items":[
{"name":"A","eventTypeName":"A","severity":"Information","defaultCategory":"X"},
{"name":"B","eventTypeName":"B","severity":"Minor","revert":"1440","defaultCategory":"X"},
{"name":"C name","eventTypeName":"C","severity":"Major","defaultCategory":"X"}]}`

func alarmRow(name, severity string, minutes int64) AlarmSeveritiesSeverities {
	row := AlarmSeveritiesSeverities{Name: types.StringValue(name), Severity: types.StringNull(), AutoClearMinutes: types.Int64Null()}
	if severity != "" {
		row.Severity = types.StringValue(severity)
	}
	if minutes != 0 {
		row.AutoClearMinutes = types.Int64Value(minutes)
	}
	return row
}

func TestAlarmWritesGroupTheNamesByValueAndSkipWhatMatches(t *testing.T) {
	writes := alarmWrites(gjson.Parse(alarmAnswer), []AlarmSeveritiesSeverities{
		alarmRow("C name", "minor", 120),
		alarmRow("A", "minor", 0),
		alarmRow("B", "minor", 1440), // both already have these values
		alarmRow("A", "", 60),
	})
	var got []string
	for _, w := range writes {
		got = append(got, w.body())
	}
	assert.Equal(t, []string{
		`{"sourceType":"scc","sourceValue":"minor","eventTypes":["C name","A"]}`,
		`{"sourceType":"aac","sourceValue":"60","eventTypes":["A"]}`,
		`{"sourceType":"aac","sourceValue":"120","eventTypes":["C name"]}`,
	}, got, "the severity calls first, then the auto clear calls by minutes; names keep their order")
}

func TestAlarmWritesTakeOnlyTheDeclaredFields(t *testing.T) {
	assert.Equal(t, 0, len(alarmWrites(gjson.Parse(alarmAnswer), []AlarmSeveritiesSeverities{alarmRow("A", "", 0)})), "a row with no field writes nothing")
	assert.Equal(t, 0, len(alarmWrites(gjson.Parse(alarmAnswer), []AlarmSeveritiesSeverities{alarmRow("A", "information", 0)})), "the read is capitalised, the write is not")
	assert.Equal(t, 1, len(alarmWrites(gjson.Parse(alarmAnswer), []AlarmSeveritiesSeverities{alarmRow("A", "", 5)})), "a type with no auto clear gets one")
	assert.Equal(t, 0, len(alarmWrites(gjson.Parse(alarmAnswer), []AlarmSeveritiesSeverities{alarmRow("Z", "minor", 5)})), "an unknown name is the caller's to refuse")
}

func TestAlarmUnknownNamesAreSortedAndListedOnce(t *testing.T) {
	unknown := alarmUnknownNames(gjson.Parse(alarmAnswer), []AlarmSeveritiesSeverities{alarmRow("b", "", 5), alarmRow("A", "", 5), alarmRow("C", "", 5), alarmRow("b", "", 5), alarmRow("C name", "", 5)})
	assert.Equal(t, []string{"C", "b"}, unknown, "names are case-sensitive and the event type of a renamed row is not its name")
}

func TestAlarmSeveritiesReadsOnlyTheRowsAndFieldsOfTheConfiguration(t *testing.T) {
	state := AlarmSeverities{Severities: []AlarmSeveritiesSeverities{alarmRow("B", "minor", 0), alarmRow("A", "", 60), alarmRow("GONE", "minor", 0), alarmRow("C name", "major", 10)}}
	state.updateFromBody(context.Background(), gjson.Parse(alarmAnswer))
	assert.Equal(t, 3, len(state.Severities), "a row that Crosswork does not have is dropped")
	assert.Equal(t, "minor", state.Severities[0].Severity.ValueString(), "the severity is lowercase")
	assert.True(t, state.Severities[0].AutoClearMinutes.IsNull(), "an auto clear that is not declared is not read")
	assert.True(t, state.Severities[1].Severity.IsNull(), "a severity that is not declared is not read")
	assert.True(t, state.Severities[1].AutoClearMinutes.IsNull(), "a declared auto clear that Crosswork does not have is read as none, so the difference shows")
	assert.Equal(t, "major", state.Severities[2].Severity.ValueString())
	assert.True(t, state.Severities[2].AutoClearMinutes.IsNull(), "a declared auto clear is read as none when there is none")

	var all AlarmSeverities
	all.fromBody(context.Background(), gjson.Parse(alarmAnswer))
	assert.Equal(t, 3, len(all.Severities))
	assert.Equal(t, int64(1440), all.Severities[1].AutoClearMinutes.ValueInt64())
	assert.True(t, all.Severities[0].AutoClearMinutes.IsNull(), "no auto clear is null")
	assert.Equal(t, "C name", all.Severities[2].Name.ValueString(), "the name, not the event type")
}

func TestAlarmHiddenSeverityNamesAreTheSeverityCallsOfTypesThePageDoesNotList(t *testing.T) {
	visible := gjson.Parse(`{"items":[{"name":"A","severity":"Minor"},{"name":"B","severity":"Minor"}]}`)
	writes := []alarmWrite{
		{sourceType: "scc", sourceValue: "major", names: []string{"B", "Z", "H"}},
		{sourceType: "aac", sourceValue: "60", names: []string{"Q"}},
	}
	assert.Equal(t, []string{"H", "Z"}, alarmHiddenSeverityNames(writes, visible), "an auto clear call is not checked")
}

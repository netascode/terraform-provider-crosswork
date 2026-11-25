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
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cw "github.com/netascode/go-crosswork"
)

const (
	alarmSeveritiesPath = "/crosswork/alarm/v1/autoclear"
	// The query returns the 1300 types that Crosswork has; without it the answer holds the 1268 that the page lists, and the other
	// 32 can be written but not read.
	alarmSeveritiesReadPath = alarmSeveritiesPath + "?completeEventList=true"
)

// A write is one call for one value and a list of names: `scc` sets a severity, `aac` an auto clear time in minutes. The controller
// takes the Name column of the page as the name of a type, applies a call whole or not at all, and answers an unknown name with a
// 400 that does not say which one.
type alarmWrite struct {
	sourceType  string
	sourceValue string
	names       []string
}

func (w alarmWrite) body() string {
	body, _ := sjson.Set("", "sourceType", w.sourceType)
	body, _ = sjson.Set(body, "sourceValue", w.sourceValue)
	body, _ = sjson.Set(body, "eventTypes", w.names)
	return body
}

// alarmSeverityIndex returns the rows of an answer by name. A name is unique and case-sensitive.
func alarmSeverityIndex(res gjson.Result) map[string]gjson.Result {
	rows := make(map[string]gjson.Result)
	res.Get("items").ForEach(func(_, row gjson.Result) bool {
		rows[row.Get("name").String()] = row
		return true
	})
	return rows
}

// alarmUnknownNames returns the names of the rows that Crosswork does not have, sorted and without repeats.
func alarmUnknownNames(res gjson.Result, rows []AlarmSeveritiesSeverities) []string {
	index := alarmSeverityIndex(res)
	seen := make(map[string]bool)
	var unknown []string
	for _, row := range rows {
		name := row.Name.ValueString()
		if _, found := index[name]; !found && !seen[name] {
			seen[name] = true
			unknown = append(unknown, name)
		}
	}
	sort.Strings(unknown)
	return unknown
}

// alarmWrites returns the calls that bring the declared fields of the rows to the values of the plan: the names are grouped by
// value, a field that a row leaves out is not written and a field that already has the value is not written either. The severity
// calls come first, then the auto clear calls, each ordered by value, so that the calls are the same on every run.
func alarmWrites(current gjson.Result, rows []AlarmSeveritiesSeverities) []alarmWrite {
	index := alarmSeverityIndex(current)
	severities := make(map[string][]string)
	minutes := make(map[int64][]string)
	for _, row := range rows {
		name := row.Name.ValueString()
		have, found := index[name]
		if !found {
			continue
		}
		if !row.Severity.IsNull() && !strings.EqualFold(have.Get("severity").String(), row.Severity.ValueString()) {
			severities[row.Severity.ValueString()] = append(severities[row.Severity.ValueString()], name)
		}
		if !row.AutoClearMinutes.IsNull() && (!have.Get("revert").Exists() || have.Get("revert").Int() != row.AutoClearMinutes.ValueInt64()) {
			minutes[row.AutoClearMinutes.ValueInt64()] = append(minutes[row.AutoClearMinutes.ValueInt64()], name)
		}
	}
	var writes []alarmWrite
	severityValues := make([]string, 0, len(severities))
	for value := range severities {
		severityValues = append(severityValues, value)
	}
	sort.Strings(severityValues)
	for _, value := range severityValues {
		writes = append(writes, alarmWrite{sourceType: "scc", sourceValue: value, names: severities[value]})
	}
	minuteValues := make([]int64, 0, len(minutes))
	for value := range minutes {
		minuteValues = append(minuteValues, value)
	}
	sort.Slice(minuteValues, func(i, j int) bool { return minuteValues[i] < minuteValues[j] })
	for _, value := range minuteValues {
		writes = append(writes, alarmWrite{sourceType: "aac", sourceValue: strconv.FormatInt(value, 10), names: minutes[value]})
	}
	return writes
}

// alarmHiddenSeverityNames returns the names of the severity calls whose type the page does not list. Crosswork answers a severity
// written to such a type with a success and changes nothing; an auto clear time is written.
func alarmHiddenSeverityNames(writes []alarmWrite, visible gjson.Result) []string {
	listed := alarmSeverityIndex(visible)
	var hidden []string
	for _, write := range writes {
		if write.sourceType != "scc" {
			continue
		}
		for _, name := range write.names {
			if _, found := listed[name]; !found {
				hidden = append(hidden, name)
			}
		}
	}
	sort.Strings(hidden)
	return hidden
}

func alarmSeverityWrites(writes []alarmWrite) []alarmWrite {
	var out []alarmWrite
	for _, write := range writes {
		if write.sourceType == "scc" {
			out = append(out, write)
		}
	}
	return out
}

// alarmSeverityFromBody reads one row. The severity is read lowercase, as it is written, and `revert` is the auto clear time in minutes
// that is absent when the type has none.
func alarmSeverityFromBody(row gjson.Result, name string, withSeverity, withAutoClear bool) AlarmSeveritiesSeverities {
	out := AlarmSeveritiesSeverities{Name: types.StringValue(name), Severity: types.StringNull(), AutoClearMinutes: types.Int64Null()}
	if withSeverity {
		out.Severity = types.StringValue(strings.ToLower(row.Get("severity").String()))
	}
	if withAutoClear && row.Get("revert").Exists() {
		out.AutoClearMinutes = types.Int64Value(row.Get("revert").Int())
	}
	return out
}

// applyAlarmSeverities reads the types of Crosswork, refuses the names it does not have before it sends anything, sends the calls
// that differ from what it read and reads the types again into the plan.
func applyAlarmSeverities(c *cw.Client, plan *AlarmSeverities) error {
	current, err := c.Get(alarmSeveritiesReadPath)
	if err != nil {
		return fmt.Errorf("%s, %s", err, current.String())
	}
	if unknown := alarmUnknownNames(current, plan.Severities); len(unknown) > 0 {
		return fmt.Errorf("unknown alarm or event types: %s; write each name as the Name column of the Severity and auto clear page shows it", strings.Join(unknown, ", "))
	}
	writes := alarmWrites(current, plan.Severities)
	if len(alarmSeverityWrites(writes)) > 0 {
		visible, err := c.Get(alarmSeveritiesPath)
		if err != nil {
			return fmt.Errorf("%s, %s", err, visible.String())
		}
		if hidden := alarmHiddenSeverityNames(writes, visible); len(hidden) > 0 {
			return fmt.Errorf("the severity of %s cannot be changed: Crosswork does not list these types on the Severity and auto clear page and ignores a severity written to them; leave severity out of these rows and set auto_clear_minutes only", strings.Join(hidden, ", "))
		}
	}
	for _, write := range writes {
		if res, err := c.Post(alarmSeveritiesPath, write.body(), cw.UseMutex); err != nil {
			return fmt.Errorf("setting %s to %s for %s: %s, %s", map[string]string{"scc": "the severity", "aac": "the auto clear minutes"}[write.sourceType], write.sourceValue, strings.Join(write.names, ", "), err, res.String())
		}
	}
	read, err := c.Get(alarmSeveritiesReadPath)
	if err != nil {
		return fmt.Errorf("%s, %s", err, read.String())
	}
	plan.updateFromBody(context.Background(), read)
	return nil
}

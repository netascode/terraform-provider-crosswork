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

// Written by hand: the generated body code that needs strconv and sjson is replaced by alarm_severities_logic.go.
import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type AlarmSeverities struct {
	Instance   types.String                `tfsdk:"instance"`
	Id         types.String                `tfsdk:"id"`
	Severities []AlarmSeveritiesSeverities `tfsdk:"severities"`
}

type AlarmSeveritiesSeverities struct {
	Name             types.String `tfsdk:"name"`
	Severity         types.String `tfsdk:"severity"`
	AutoClearMinutes types.Int64  `tfsdk:"auto_clear_minutes"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data AlarmSeverities) getPath() string {
	return "/crosswork/alarm/v1/autoclear"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data AlarmSeverities) getPathGet() string {
	return "/crosswork/alarm/v1/autoclear?completeEventList=true"
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

// Written by hand: a write reads the types of Crosswork first (alarmWrites), so the model has no body of its own.
func (data AlarmSeverities) toBody(ctx context.Context, state AlarmSeverities) string {
	return "{}"
}

// End of section. //template:end toBody

// Written by hand: every type of the answer is read (a data source and an import).
func (data *AlarmSeverities) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("alarm_severities")
	data.Severities = make([]AlarmSeveritiesSeverities, 0)
	res.Get("items").ForEach(func(_, row gjson.Result) bool {
		data.Severities = append(data.Severities, alarmSeverityFromBody(row, row.Get("name").String(), true, true))
		return true
	})
}

// End of section. //template:end fromBody

// Written by hand: only the rows of the configuration are read back, and a field only when the row declares it. A row that Crosswork does
// not have is dropped.
func (data *AlarmSeverities) updateFromBody(ctx context.Context, res gjson.Result) {
	index := alarmSeverityIndex(res)
	rows := make([]AlarmSeveritiesSeverities, 0, len(data.Severities))
	for _, planned := range data.Severities {
		if row, found := index[planned.Name.ValueString()]; found {
			rows = append(rows, alarmSeverityFromBody(row, planned.Name.ValueString(), !planned.Severity.IsNull(), !planned.AutoClearMinutes.IsNull()))
		}
	}
	data.Severities = rows
}

// End of section. //template:end updateFromBody

// Written by hand: an import reads no row. The state of an import has no rows, and the configuration then names the rows to manage
// (a state of every one of the 1300 types would make every plan remove most of them from the list).
func (data *AlarmSeverities) isNull(ctx context.Context, res gjson.Result) bool {
	return false
}

// End of section. //template:end isNull

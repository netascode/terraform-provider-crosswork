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
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type AlarmSettings struct {
	Instance                           types.String `tfsdk:"instance"`
	Id                                 types.String `tfsdk:"id"`
	SecurityDeviceAlarmsClearedDays    types.Int64  `tfsdk:"security_device_alarms_cleared_days"`
	NonSecurityDeviceAlarmsClearedDays types.Int64  `tfsdk:"non_security_device_alarms_cleared_days"`
	AllEventsDays                      types.Int64  `tfsdk:"all_events_days"`
	DeviceAlarmsDays                   types.Int64  `tfsdk:"device_alarms_days"`
	SystemAlarmsDays                   types.Int64  `tfsdk:"system_alarms_days"`
	NetworkAlarmsDays                  types.Int64  `tfsdk:"network_alarms_days"`
	AuditLogsDays                      types.Int64  `tfsdk:"audit_logs_days"`
	TrapProcessing                     types.Bool   `tfsdk:"trap_processing"`
	SyslogProcessing                   types.Bool   `tfsdk:"syslog_processing"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data AlarmSettings) getPath() string {
	return "/crosswork/alarm/v1/settings"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data AlarmSettings) getPathGet() string {
	return "/crosswork/alarm/v1/settings"
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

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data AlarmSettings) toBody(ctx context.Context, state AlarmSettings) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.SecurityDeviceAlarmsClearedDays.IsNull() {
		body, _ = sjson.Set(body, "securityAlertAgeout", data.SecurityDeviceAlarmsClearedDays.ValueInt64())
	}
	if !data.NonSecurityDeviceAlarmsClearedDays.IsNull() {
		body, _ = sjson.Set(body, "nonSecurityAlertAgeout", data.NonSecurityDeviceAlarmsClearedDays.ValueInt64())
	}
	if !data.AllEventsDays.IsNull() {
		body, _ = sjson.Set(body, "deleteAllEvents", data.AllEventsDays.ValueInt64())
	}
	if !data.DeviceAlarmsDays.IsNull() {
		body, _ = sjson.Set(body, "deleteOldAlertDays", data.DeviceAlarmsDays.ValueInt64())
	}
	if !data.SystemAlarmsDays.IsNull() {
		body, _ = sjson.Set(body, "systemAlertAgeout", data.SystemAlarmsDays.ValueInt64())
	}
	if !data.NetworkAlarmsDays.IsNull() {
		body, _ = sjson.Set(body, "networkAlertAgeout", data.NetworkAlarmsDays.ValueInt64())
	}
	if !data.AuditLogsDays.IsNull() {
		body, _ = sjson.Set(body, "auditAlertAgeout", data.AuditLogsDays.ValueInt64())
	}
	if !data.TrapProcessing.IsNull() {
		body, _ = sjson.Set(body, "trapCollectionJobEnable", data.TrapProcessing.ValueBool())
	}
	if !data.SyslogProcessing.IsNull() {
		body, _ = sjson.Set(body, "syslogCollectionJobEnable", data.SyslogProcessing.ValueBool())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *AlarmSettings) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("alarm_settings")
	if value := res.Get("securityAlertAgeout"); value.Exists() {
		data.SecurityDeviceAlarmsClearedDays = types.Int64Value(value.Int())
	} else {
		data.SecurityDeviceAlarmsClearedDays = types.Int64Value(30)
	}
	if value := res.Get("nonSecurityAlertAgeout"); value.Exists() {
		data.NonSecurityDeviceAlarmsClearedDays = types.Int64Value(value.Int())
	} else {
		data.NonSecurityDeviceAlarmsClearedDays = types.Int64Value(7)
	}
	if value := res.Get("deleteAllEvents"); value.Exists() {
		data.AllEventsDays = types.Int64Value(value.Int())
	} else {
		data.AllEventsDays = types.Int64Value(30)
	}
	if value := res.Get("deleteOldAlertDays"); value.Exists() {
		data.DeviceAlarmsDays = types.Int64Value(value.Int())
	} else {
		data.DeviceAlarmsDays = types.Int64Value(60)
	}
	if value := res.Get("systemAlertAgeout"); value.Exists() {
		data.SystemAlarmsDays = types.Int64Value(value.Int())
	} else {
		data.SystemAlarmsDays = types.Int64Value(60)
	}
	if value := res.Get("networkAlertAgeout"); value.Exists() {
		data.NetworkAlarmsDays = types.Int64Value(value.Int())
	} else {
		data.NetworkAlarmsDays = types.Int64Value(60)
	}
	if value := res.Get("auditAlertAgeout"); value.Exists() {
		data.AuditLogsDays = types.Int64Value(value.Int())
	} else {
		data.AuditLogsDays = types.Int64Value(30)
	}
	if value := res.Get("trapCollectionJobEnable"); value.Exists() {
		data.TrapProcessing = types.BoolValue(value.Bool())
	} else {
		data.TrapProcessing = types.BoolValue(true)
	}
	if value := res.Get("syslogCollectionJobEnable"); value.Exists() {
		data.SyslogProcessing = types.BoolValue(value.Bool())
	} else {
		data.SyslogProcessing = types.BoolValue(true)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *AlarmSettings) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("securityAlertAgeout"); value.Exists() && !data.SecurityDeviceAlarmsClearedDays.IsNull() {
		data.SecurityDeviceAlarmsClearedDays = types.Int64Value(value.Int())
	} else if data.SecurityDeviceAlarmsClearedDays.ValueInt64() != 30 {
		data.SecurityDeviceAlarmsClearedDays = types.Int64Null()
	}
	if value := res.Get("nonSecurityAlertAgeout"); value.Exists() && !data.NonSecurityDeviceAlarmsClearedDays.IsNull() {
		data.NonSecurityDeviceAlarmsClearedDays = types.Int64Value(value.Int())
	} else if data.NonSecurityDeviceAlarmsClearedDays.ValueInt64() != 7 {
		data.NonSecurityDeviceAlarmsClearedDays = types.Int64Null()
	}
	if value := res.Get("deleteAllEvents"); value.Exists() && !data.AllEventsDays.IsNull() {
		data.AllEventsDays = types.Int64Value(value.Int())
	} else if data.AllEventsDays.ValueInt64() != 30 {
		data.AllEventsDays = types.Int64Null()
	}
	if value := res.Get("deleteOldAlertDays"); value.Exists() && !data.DeviceAlarmsDays.IsNull() {
		data.DeviceAlarmsDays = types.Int64Value(value.Int())
	} else if data.DeviceAlarmsDays.ValueInt64() != 60 {
		data.DeviceAlarmsDays = types.Int64Null()
	}
	if value := res.Get("systemAlertAgeout"); value.Exists() && !data.SystemAlarmsDays.IsNull() {
		data.SystemAlarmsDays = types.Int64Value(value.Int())
	} else if data.SystemAlarmsDays.ValueInt64() != 60 {
		data.SystemAlarmsDays = types.Int64Null()
	}
	if value := res.Get("networkAlertAgeout"); value.Exists() && !data.NetworkAlarmsDays.IsNull() {
		data.NetworkAlarmsDays = types.Int64Value(value.Int())
	} else if data.NetworkAlarmsDays.ValueInt64() != 60 {
		data.NetworkAlarmsDays = types.Int64Null()
	}
	if value := res.Get("auditAlertAgeout"); value.Exists() && !data.AuditLogsDays.IsNull() {
		data.AuditLogsDays = types.Int64Value(value.Int())
	} else if data.AuditLogsDays.ValueInt64() != 30 {
		data.AuditLogsDays = types.Int64Null()
	}
	if value := res.Get("trapCollectionJobEnable"); value.Exists() && !data.TrapProcessing.IsNull() {
		data.TrapProcessing = types.BoolValue(value.Bool())
	} else if data.TrapProcessing.ValueBool() != true {
		data.TrapProcessing = types.BoolNull()
	}
	if value := res.Get("syslogCollectionJobEnable"); value.Exists() && !data.SyslogProcessing.IsNull() {
		data.SyslogProcessing = types.BoolValue(value.Bool())
	} else if data.SyslogProcessing.ValueBool() != true {
		data.SyslogProcessing = types.BoolNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *AlarmSettings) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.SecurityDeviceAlarmsClearedDays.IsNull() {
		return false
	}
	if !data.NonSecurityDeviceAlarmsClearedDays.IsNull() {
		return false
	}
	if !data.AllEventsDays.IsNull() {
		return false
	}
	if !data.DeviceAlarmsDays.IsNull() {
		return false
	}
	if !data.SystemAlarmsDays.IsNull() {
		return false
	}
	if !data.NetworkAlarmsDays.IsNull() {
		return false
	}
	if !data.AuditLogsDays.IsNull() {
		return false
	}
	if !data.TrapProcessing.IsNull() {
		return false
	}
	if !data.SyslogProcessing.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

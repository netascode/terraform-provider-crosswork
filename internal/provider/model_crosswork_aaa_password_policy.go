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
type AAAPasswordPolicy struct {
	Instance                         types.String `tfsdk:"instance"`
	Id                               types.String `tfsdk:"id"`
	MinPasswordLength                types.Int64  `tfsdk:"min_password_length"`
	NoUsername                       types.Bool   `tfsdk:"no_username"`
	NoCiscoVariant                   types.Bool   `tfsdk:"no_cisco_variant"`
	NoCharRepetition                 types.Bool   `tfsdk:"no_char_repetition"`
	NumChangedCharsEnabled           types.Bool   `tfsdk:"num_changed_chars_enabled"`
	NumChangedChars                  types.Int64  `tfsdk:"num_changed_chars"`
	NumReuseLimitEnabled             types.Bool   `tfsdk:"num_reuse_limit_enabled"`
	NumReuseLimit                    types.Int64  `tfsdk:"num_reuse_limit"`
	PasswordReuseDaysEnabled         types.Bool   `tfsdk:"password_reuse_days_enabled"`
	PasswordReuseDays                types.Int64  `tfsdk:"password_reuse_days"`
	FailedLoginsBeforeLockoutEnabled types.Bool   `tfsdk:"failed_logins_before_lockout_enabled"`
	FailedLoginsBeforeLockout        types.Int64  `tfsdk:"failed_logins_before_lockout"`
	LockOutUserTimeEnabled           types.Bool   `tfsdk:"lock_out_user_time_enabled"`
	LockOutUserTime                  types.Int64  `tfsdk:"lock_out_user_time"`
	PasswordExpiryDaysEnabled        types.Bool   `tfsdk:"password_expiry_days_enabled"`
	PasswordExpiryDays               types.Int64  `tfsdk:"password_expiry_days"`
	DaysForWarningEnabled            types.Bool   `tfsdk:"days_for_warning_enabled"`
	DaysForWarning                   types.Int64  `tfsdk:"days_for_warning"`
	ChangePasswordOnFirstLogin       types.Bool   `tfsdk:"change_password_on_first_login"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data AAAPasswordPolicy) getPath() string {
	return "/crosswork/aaa/v1/passwordPolicyConfig"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

func (data AAAPasswordPolicy) getPathDelete() string {
	return "/crosswork/aaa/v1/passwordPolicyConfig"
}

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data AAAPasswordPolicy) getPathGet() string {
	return "/crosswork/aaa/v1/passwordPolicyConfig"
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
func (data AAAPasswordPolicy) toBody(ctx context.Context, state AAAPasswordPolicy) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.MinPasswordLength.IsNull() {
		body, _ = sjson.Set(body, "MinPasswordLength", data.MinPasswordLength.ValueInt64())
	}
	if !data.NoUsername.IsNull() {
		body, _ = sjson.Set(body, "NoUsername", data.NoUsername.ValueBool())
	}
	if !data.NoCiscoVariant.IsNull() {
		body, _ = sjson.Set(body, "NoCiscoVariant", data.NoCiscoVariant.ValueBool())
	}
	if !data.NoCharRepetition.IsNull() {
		body, _ = sjson.Set(body, "NoCharRepetition", data.NoCharRepetition.ValueBool())
	}
	if !data.NumChangedCharsEnabled.IsNull() {
		body, _ = sjson.Set(body, "NumChangedCharsEnable", data.NumChangedCharsEnabled.ValueBool())
	}
	if !data.NumChangedChars.IsNull() {
		body, _ = sjson.Set(body, "NumChangedChars", data.NumChangedChars.ValueInt64())
	}
	if !data.NumReuseLimitEnabled.IsNull() {
		body, _ = sjson.Set(body, "NumReuseLimitEnable", data.NumReuseLimitEnabled.ValueBool())
	}
	if !data.NumReuseLimit.IsNull() {
		body, _ = sjson.Set(body, "NumReuseLimit", data.NumReuseLimit.ValueInt64())
	}
	if !data.PasswordReuseDaysEnabled.IsNull() {
		body, _ = sjson.Set(body, "PasswordReuseDaysEnable", data.PasswordReuseDaysEnabled.ValueBool())
	}
	if !data.PasswordReuseDays.IsNull() {
		body, _ = sjson.Set(body, "PasswordReuseDays", data.PasswordReuseDays.ValueInt64())
	}
	if !data.FailedLoginsBeforeLockoutEnabled.IsNull() {
		body, _ = sjson.Set(body, "FailedLoginsBefLoEnable", data.FailedLoginsBeforeLockoutEnabled.ValueBool())
	}
	if !data.FailedLoginsBeforeLockout.IsNull() {
		body, _ = sjson.Set(body, "FailedLoginsBeforeLockout", data.FailedLoginsBeforeLockout.ValueInt64())
	}
	if !data.LockOutUserTimeEnabled.IsNull() {
		body, _ = sjson.Set(body, "LockOutUserTimeEnable", data.LockOutUserTimeEnabled.ValueBool())
	}
	if !data.LockOutUserTime.IsNull() {
		body, _ = sjson.Set(body, "LockOutUserTime", data.LockOutUserTime.ValueInt64())
	}
	if !data.PasswordExpiryDaysEnabled.IsNull() {
		body, _ = sjson.Set(body, "PasswordExpiryDaysEnable", data.PasswordExpiryDaysEnabled.ValueBool())
	}
	if !data.PasswordExpiryDays.IsNull() {
		body, _ = sjson.Set(body, "PasswordExpiryDays", data.PasswordExpiryDays.ValueInt64())
	}
	if !data.DaysForWarningEnabled.IsNull() {
		body, _ = sjson.Set(body, "DaysForWarningEnable", data.DaysForWarningEnabled.ValueBool())
	}
	if !data.DaysForWarning.IsNull() {
		body, _ = sjson.Set(body, "DaysForWarning", data.DaysForWarning.ValueInt64())
	}
	if !data.ChangePasswordOnFirstLogin.IsNull() {
		body, _ = sjson.Set(body, "ChangePasswdOnFirstLogin", data.ChangePasswordOnFirstLogin.ValueBool())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *AAAPasswordPolicy) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("aaa_password_policy")
	if value := res.Get("MinPasswordLength"); value.Exists() {
		data.MinPasswordLength = types.Int64Value(value.Int())
	} else {
		data.MinPasswordLength = types.Int64Value(8)
	}
	if value := res.Get("NoUsername"); value.Exists() {
		data.NoUsername = types.BoolValue(value.Bool())
	} else {
		data.NoUsername = types.BoolValue(true)
	}
	if value := res.Get("NoCiscoVariant"); value.Exists() {
		data.NoCiscoVariant = types.BoolValue(value.Bool())
	} else {
		data.NoCiscoVariant = types.BoolValue(true)
	}
	if value := res.Get("NoCharRepetition"); value.Exists() {
		data.NoCharRepetition = types.BoolValue(value.Bool())
	} else {
		data.NoCharRepetition = types.BoolValue(true)
	}
	if value := res.Get("NumChangedCharsEnable"); value.Exists() {
		data.NumChangedCharsEnabled = types.BoolValue(value.Bool())
	} else {
		data.NumChangedCharsEnabled = types.BoolValue(true)
	}
	if value := res.Get("NumChangedChars"); value.Exists() {
		data.NumChangedChars = types.Int64Value(value.Int())
	} else {
		data.NumChangedChars = types.Int64Value(1)
	}
	if value := res.Get("NumReuseLimitEnable"); value.Exists() {
		data.NumReuseLimitEnabled = types.BoolValue(value.Bool())
	} else {
		data.NumReuseLimitEnabled = types.BoolValue(true)
	}
	if value := res.Get("NumReuseLimit"); value.Exists() {
		data.NumReuseLimit = types.Int64Value(value.Int())
	} else {
		data.NumReuseLimit = types.Int64Value(5)
	}
	if value := res.Get("PasswordReuseDaysEnable"); value.Exists() {
		data.PasswordReuseDaysEnabled = types.BoolValue(value.Bool())
	} else {
		data.PasswordReuseDaysEnabled = types.BoolValue(true)
	}
	if value := res.Get("PasswordReuseDays"); value.Exists() {
		data.PasswordReuseDays = types.Int64Value(value.Int())
	} else {
		data.PasswordReuseDays = types.Int64Value(15)
	}
	if value := res.Get("FailedLoginsBefLoEnable"); value.Exists() {
		data.FailedLoginsBeforeLockoutEnabled = types.BoolValue(value.Bool())
	} else {
		data.FailedLoginsBeforeLockoutEnabled = types.BoolValue(true)
	}
	if value := res.Get("FailedLoginsBeforeLockout"); value.Exists() {
		data.FailedLoginsBeforeLockout = types.Int64Value(value.Int())
	} else {
		data.FailedLoginsBeforeLockout = types.Int64Value(5)
	}
	if value := res.Get("LockOutUserTimeEnable"); value.Exists() {
		data.LockOutUserTimeEnabled = types.BoolValue(value.Bool())
	} else {
		data.LockOutUserTimeEnabled = types.BoolValue(true)
	}
	if value := res.Get("LockOutUserTime"); value.Exists() {
		data.LockOutUserTime = types.Int64Value(value.Int())
	} else {
		data.LockOutUserTime = types.Int64Value(6)
	}
	if value := res.Get("PasswordExpiryDaysEnable"); value.Exists() {
		data.PasswordExpiryDaysEnabled = types.BoolValue(value.Bool())
	} else {
		data.PasswordExpiryDaysEnabled = types.BoolValue(true)
	}
	if value := res.Get("PasswordExpiryDays"); value.Exists() {
		data.PasswordExpiryDays = types.Int64Value(value.Int())
	} else {
		data.PasswordExpiryDays = types.Int64Value(60)
	}
	if value := res.Get("DaysForWarningEnable"); value.Exists() {
		data.DaysForWarningEnabled = types.BoolValue(value.Bool())
	} else {
		data.DaysForWarningEnabled = types.BoolValue(true)
	}
	if value := res.Get("DaysForWarning"); value.Exists() {
		data.DaysForWarning = types.Int64Value(value.Int())
	} else {
		data.DaysForWarning = types.Int64Value(15)
	}
	if value := res.Get("ChangePasswdOnFirstLogin"); value.Exists() {
		data.ChangePasswordOnFirstLogin = types.BoolValue(value.Bool())
	} else {
		data.ChangePasswordOnFirstLogin = types.BoolValue(true)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *AAAPasswordPolicy) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("MinPasswordLength"); value.Exists() && !data.MinPasswordLength.IsNull() {
		data.MinPasswordLength = types.Int64Value(value.Int())
	} else if data.MinPasswordLength.ValueInt64() != 8 {
		data.MinPasswordLength = types.Int64Null()
	}
	if value := res.Get("NoUsername"); value.Exists() && !data.NoUsername.IsNull() {
		data.NoUsername = types.BoolValue(value.Bool())
	} else if data.NoUsername.ValueBool() != true {
		data.NoUsername = types.BoolNull()
	}
	if value := res.Get("NoCiscoVariant"); value.Exists() && !data.NoCiscoVariant.IsNull() {
		data.NoCiscoVariant = types.BoolValue(value.Bool())
	} else if data.NoCiscoVariant.ValueBool() != true {
		data.NoCiscoVariant = types.BoolNull()
	}
	if value := res.Get("NoCharRepetition"); value.Exists() && !data.NoCharRepetition.IsNull() {
		data.NoCharRepetition = types.BoolValue(value.Bool())
	} else if data.NoCharRepetition.ValueBool() != true {
		data.NoCharRepetition = types.BoolNull()
	}
	if value := res.Get("NumChangedCharsEnable"); value.Exists() && !data.NumChangedCharsEnabled.IsNull() {
		data.NumChangedCharsEnabled = types.BoolValue(value.Bool())
	} else if data.NumChangedCharsEnabled.ValueBool() != true {
		data.NumChangedCharsEnabled = types.BoolNull()
	}
	if value := res.Get("NumChangedChars"); value.Exists() && !data.NumChangedChars.IsNull() {
		data.NumChangedChars = types.Int64Value(value.Int())
	} else if data.NumChangedChars.ValueInt64() != 1 {
		data.NumChangedChars = types.Int64Null()
	}
	if value := res.Get("NumReuseLimitEnable"); value.Exists() && !data.NumReuseLimitEnabled.IsNull() {
		data.NumReuseLimitEnabled = types.BoolValue(value.Bool())
	} else if data.NumReuseLimitEnabled.ValueBool() != true {
		data.NumReuseLimitEnabled = types.BoolNull()
	}
	if value := res.Get("NumReuseLimit"); value.Exists() && !data.NumReuseLimit.IsNull() {
		data.NumReuseLimit = types.Int64Value(value.Int())
	} else if data.NumReuseLimit.ValueInt64() != 5 {
		data.NumReuseLimit = types.Int64Null()
	}
	if value := res.Get("PasswordReuseDaysEnable"); value.Exists() && !data.PasswordReuseDaysEnabled.IsNull() {
		data.PasswordReuseDaysEnabled = types.BoolValue(value.Bool())
	} else if data.PasswordReuseDaysEnabled.ValueBool() != true {
		data.PasswordReuseDaysEnabled = types.BoolNull()
	}
	if value := res.Get("PasswordReuseDays"); value.Exists() && !data.PasswordReuseDays.IsNull() {
		data.PasswordReuseDays = types.Int64Value(value.Int())
	} else if data.PasswordReuseDays.ValueInt64() != 15 {
		data.PasswordReuseDays = types.Int64Null()
	}
	if value := res.Get("FailedLoginsBefLoEnable"); value.Exists() && !data.FailedLoginsBeforeLockoutEnabled.IsNull() {
		data.FailedLoginsBeforeLockoutEnabled = types.BoolValue(value.Bool())
	} else if data.FailedLoginsBeforeLockoutEnabled.ValueBool() != true {
		data.FailedLoginsBeforeLockoutEnabled = types.BoolNull()
	}
	if value := res.Get("FailedLoginsBeforeLockout"); value.Exists() && !data.FailedLoginsBeforeLockout.IsNull() {
		data.FailedLoginsBeforeLockout = types.Int64Value(value.Int())
	} else if data.FailedLoginsBeforeLockout.ValueInt64() != 5 {
		data.FailedLoginsBeforeLockout = types.Int64Null()
	}
	if value := res.Get("LockOutUserTimeEnable"); value.Exists() && !data.LockOutUserTimeEnabled.IsNull() {
		data.LockOutUserTimeEnabled = types.BoolValue(value.Bool())
	} else if data.LockOutUserTimeEnabled.ValueBool() != true {
		data.LockOutUserTimeEnabled = types.BoolNull()
	}
	if value := res.Get("LockOutUserTime"); value.Exists() && !data.LockOutUserTime.IsNull() {
		data.LockOutUserTime = types.Int64Value(value.Int())
	} else if data.LockOutUserTime.ValueInt64() != 6 {
		data.LockOutUserTime = types.Int64Null()
	}
	if value := res.Get("PasswordExpiryDaysEnable"); value.Exists() && !data.PasswordExpiryDaysEnabled.IsNull() {
		data.PasswordExpiryDaysEnabled = types.BoolValue(value.Bool())
	} else if data.PasswordExpiryDaysEnabled.ValueBool() != true {
		data.PasswordExpiryDaysEnabled = types.BoolNull()
	}
	if value := res.Get("PasswordExpiryDays"); value.Exists() && !data.PasswordExpiryDays.IsNull() {
		data.PasswordExpiryDays = types.Int64Value(value.Int())
	} else if data.PasswordExpiryDays.ValueInt64() != 60 {
		data.PasswordExpiryDays = types.Int64Null()
	}
	if value := res.Get("DaysForWarningEnable"); value.Exists() && !data.DaysForWarningEnabled.IsNull() {
		data.DaysForWarningEnabled = types.BoolValue(value.Bool())
	} else if data.DaysForWarningEnabled.ValueBool() != true {
		data.DaysForWarningEnabled = types.BoolNull()
	}
	if value := res.Get("DaysForWarning"); value.Exists() && !data.DaysForWarning.IsNull() {
		data.DaysForWarning = types.Int64Value(value.Int())
	} else if data.DaysForWarning.ValueInt64() != 15 {
		data.DaysForWarning = types.Int64Null()
	}
	if value := res.Get("ChangePasswdOnFirstLogin"); value.Exists() && !data.ChangePasswordOnFirstLogin.IsNull() {
		data.ChangePasswordOnFirstLogin = types.BoolValue(value.Bool())
	} else if data.ChangePasswordOnFirstLogin.ValueBool() != true {
		data.ChangePasswordOnFirstLogin = types.BoolNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *AAAPasswordPolicy) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.MinPasswordLength.IsNull() {
		return false
	}
	if !data.NoUsername.IsNull() {
		return false
	}
	if !data.NoCiscoVariant.IsNull() {
		return false
	}
	if !data.NoCharRepetition.IsNull() {
		return false
	}
	if !data.NumChangedCharsEnabled.IsNull() {
		return false
	}
	if !data.NumChangedChars.IsNull() {
		return false
	}
	if !data.NumReuseLimitEnabled.IsNull() {
		return false
	}
	if !data.NumReuseLimit.IsNull() {
		return false
	}
	if !data.PasswordReuseDaysEnabled.IsNull() {
		return false
	}
	if !data.PasswordReuseDays.IsNull() {
		return false
	}
	if !data.FailedLoginsBeforeLockoutEnabled.IsNull() {
		return false
	}
	if !data.FailedLoginsBeforeLockout.IsNull() {
		return false
	}
	if !data.LockOutUserTimeEnabled.IsNull() {
		return false
	}
	if !data.LockOutUserTime.IsNull() {
		return false
	}
	if !data.PasswordExpiryDaysEnabled.IsNull() {
		return false
	}
	if !data.PasswordExpiryDays.IsNull() {
		return false
	}
	if !data.DaysForWarningEnabled.IsNull() {
		return false
	}
	if !data.DaysForWarning.IsNull() {
		return false
	}
	if !data.ChangePasswordOnFirstLogin.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

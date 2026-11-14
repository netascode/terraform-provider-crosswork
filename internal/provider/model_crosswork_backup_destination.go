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
type BackupDestination struct {
	Instance          types.String `tfsdk:"instance"`
	Id                types.String `tfsdk:"id"`
	HostName          types.String `tfsdk:"host_name"`
	Port              types.Int64  `tfsdk:"port"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	PasswordWo        types.String `tfsdk:"password_wo"`
	PasswordWoVersion types.Int64  `tfsdk:"password_wo_version"`
	RemotePath        types.String `tfsdk:"remote_path"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data BackupDestination) getPath() string {
	return "/crosswork/platform/v1/platform/backup/remote"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data BackupDestination) getPathGet() string {
	return "/crosswork/platform/v1/platform/backup/remote"
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
func (data BackupDestination) toBody(ctx context.Context, state BackupDestination) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.HostName.IsNull() {
		body, _ = sjson.Set(body, "hostName", data.HostName.ValueString())
	}
	if !data.Port.IsNull() {
		body, _ = sjson.Set(body, "port", data.Port.ValueInt64())
	}
	if !data.Username.IsNull() {
		body, _ = sjson.Set(body, "username", data.Username.ValueString())
	}
	if !data.Password.IsNull() {
		body, _ = sjson.Set(body, "password", data.Password.ValueString())
	}
	if !data.PasswordWo.IsNull() {
		body, _ = sjson.Set(body, "password", data.PasswordWo.ValueString())
	}
	if !data.RemotePath.IsNull() {
		body, _ = sjson.Set(body, "remotePath", data.RemotePath.ValueString())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *BackupDestination) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("backup_destination")
	if value := res.Get("hostName"); value.Exists() {
		data.HostName = types.StringValue(value.String())
	} else {
		data.HostName = types.StringNull()
	}
	if value := res.Get("port"); value.Exists() {
		data.Port = types.Int64Value(value.Int())
	} else {
		data.Port = types.Int64Null()
	}
	if value := res.Get("username"); value.Exists() {
		data.Username = types.StringValue(value.String())
	} else {
		data.Username = types.StringNull()
	}
	if value := res.Get("remotePath"); value.Exists() {
		data.RemotePath = types.StringValue(value.String())
	} else {
		data.RemotePath = types.StringNull()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *BackupDestination) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("hostName"); value.Exists() && !data.HostName.IsNull() {
		data.HostName = types.StringValue(value.String())
	} else {
		data.HostName = types.StringNull()
	}
	if value := res.Get("port"); value.Exists() && !data.Port.IsNull() {
		data.Port = types.Int64Value(value.Int())
	} else {
		data.Port = types.Int64Null()
	}
	if value := res.Get("username"); value.Exists() && !data.Username.IsNull() {
		data.Username = types.StringValue(value.String())
	} else {
		data.Username = types.StringNull()
	}
	if value := res.Get("remotePath"); value.Exists() && !data.RemotePath.IsNull() {
		data.RemotePath = types.StringValue(value.String())
	} else {
		data.RemotePath = types.StringNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *BackupDestination) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.HostName.IsNull() {
		return false
	}
	if !data.Port.IsNull() {
		return false
	}
	if !data.Username.IsNull() {
		return false
	}
	if !data.Password.IsNull() {
		return false
	}
	if !data.PasswordWo.IsNull() {
		return false
	}
	if !data.PasswordWoVersion.IsNull() {
		return false
	}
	if !data.RemotePath.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

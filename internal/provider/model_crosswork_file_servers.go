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
type FileServers struct {
	Instance          types.String `tfsdk:"instance"`
	Id                types.String `tfsdk:"id"`
	FtpEnabled        types.Bool   `tfsdk:"ftp_enabled"`
	SftpUploadEnabled types.Bool   `tfsdk:"sftp_upload_enabled"`
	TftpEnabled       types.Bool   `tfsdk:"tftp_enabled"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data FileServers) getPath() string {
	return "/crosswork/platform/v2/xftp/server_update"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data FileServers) getPathGet() string {
	return "/crosswork/platform/v2/xftp/server_status"
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
func (data FileServers) toBody(ctx context.Context, state FileServers) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.FtpEnabled.IsNull() {
		body, _ = sjson.Set(body, "ftpEnable", data.FtpEnabled.ValueBool())
	}
	if !data.SftpUploadEnabled.IsNull() {
		body, _ = sjson.Set(body, "sftpUploadEnable", data.SftpUploadEnabled.ValueBool())
	}
	if !data.TftpEnabled.IsNull() {
		body, _ = sjson.Set(body, "tftpEnable", data.TftpEnabled.ValueBool())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *FileServers) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("file_servers")
	if value := res.Get("ftpEnabled"); value.Exists() {
		data.FtpEnabled = types.BoolValue(value.Bool())
	} else {
		data.FtpEnabled = types.BoolValue(false)
	}
	if value := res.Get("sftpUploadEnabled"); value.Exists() {
		data.SftpUploadEnabled = types.BoolValue(value.Bool())
	} else {
		data.SftpUploadEnabled = types.BoolValue(false)
	}
	if value := res.Get("tftpEnabled"); value.Exists() {
		data.TftpEnabled = types.BoolValue(value.Bool())
	} else {
		data.TftpEnabled = types.BoolValue(false)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *FileServers) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("ftpEnabled"); value.Exists() && !data.FtpEnabled.IsNull() {
		data.FtpEnabled = types.BoolValue(value.Bool())
	} else if data.FtpEnabled.ValueBool() != false {
		data.FtpEnabled = types.BoolNull()
	}
	if value := res.Get("sftpUploadEnabled"); value.Exists() && !data.SftpUploadEnabled.IsNull() {
		data.SftpUploadEnabled = types.BoolValue(value.Bool())
	} else if data.SftpUploadEnabled.ValueBool() != false {
		data.SftpUploadEnabled = types.BoolNull()
	}
	if value := res.Get("tftpEnabled"); value.Exists() && !data.TftpEnabled.IsNull() {
		data.TftpEnabled = types.BoolValue(value.Bool())
	} else if data.TftpEnabled.ValueBool() != false {
		data.TftpEnabled = types.BoolNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *FileServers) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.FtpEnabled.IsNull() {
		return false
	}
	if !data.SftpUploadEnabled.IsNull() {
		return false
	}
	if !data.TftpEnabled.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

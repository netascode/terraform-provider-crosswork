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
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type SyslogDestination struct {
	Instance      types.String `tfsdk:"instance"`
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Address       types.String `tfsdk:"address"`
	Port          types.Int64  `tfsdk:"port"`
	Protocol      types.String `tfsdk:"protocol"`
	CertificateId types.String `tfsdk:"certificate_id"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data SyslogDestination) getPath() string {
	return "/crosswork/alarm/v1/destination/syslog"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data SyslogDestination) getPathGet() string {
	return "/crosswork/alarm/v1/destination/syslog"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Written by hand: the controller lists every syslog receiver, and the read returns the one with the name (a not found when it has none).
func (data SyslogDestination) readGet(c *cw.Client, path string) (cw.Res, error) {
	name := data.Name.ValueString()
	if name == "" {
		name = data.Id.ValueString()
	}
	return readAlarmReceiver(c, syslogDestinationPath, name)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data SyslogDestination) toBody(ctx context.Context, state SyslogDestination) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() && put == false {
		body, _ = sjson.Set(body, "notificationSyslog.name", data.Name.ValueString())
	}
	if !data.Address.IsNull() {
		body, _ = sjson.Set(body, "notificationSyslog.serverAddress", data.Address.ValueString())
	}
	if !data.Port.IsNull() {
		body, _ = sjson.Set(body, "notificationSyslog.port", data.Port.ValueInt64())
	}
	if !data.Protocol.IsNull() {
		body, _ = sjson.Set(body, "notificationSyslog.transportType", data.Protocol.ValueString())
	}
	if !data.CertificateId.IsNull() {
		body, _ = sjson.Set(body, "notificationSyslog.certificateId", data.CertificateId.ValueString())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *SyslogDestination) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("serverAddress"); value.Exists() {
		data.Address = types.StringValue(value.String())
	} else {
		data.Address = types.StringNull()
	}
	if value := res.Get("port"); value.Exists() {
		data.Port = types.Int64Value(value.Int())
	} else {
		data.Port = types.Int64Value(514)
	}
	if value := res.Get("transportType"); value.Exists() {
		data.Protocol = types.StringValue(value.String())
	} else {
		data.Protocol = types.StringValue("UDP")
	}
	if value := res.Get("certificateId"); value.Exists() && value.String() != "" {
		data.CertificateId = types.StringValue(value.String())
	} else {
		data.CertificateId = types.StringNull()
	}
}

// End of section. //template:end fromBody

// Written by hand: the address is never read back (the controller keeps the IP that a host name resolved to at the creation), so
// the value of the configuration stays; it is read only when the state has none.
func (data *SyslogDestination) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("serverAddress"); value.Exists() && data.Address.IsNull() {
		data.Address = types.StringValue(value.String())
	}
	if value := res.Get("port"); value.Exists() && !data.Port.IsNull() {
		data.Port = types.Int64Value(value.Int())
	} else if data.Port.ValueInt64() != 514 {
		data.Port = types.Int64Null()
	}
	if value := res.Get("transportType"); value.Exists() && !data.Protocol.IsNull() {
		data.Protocol = types.StringValue(value.String())
	} else if data.Protocol.ValueString() != "UDP" {
		data.Protocol = types.StringNull()
	}
	if value := res.Get("certificateId"); value.Exists() && value.String() != "" && !data.CertificateId.IsNull() {
		data.CertificateId = types.StringValue(value.String())
	} else {
		data.CertificateId = types.StringNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *SyslogDestination) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Address.IsNull() {
		return false
	}
	if !data.Port.IsNull() {
		return false
	}
	if !data.Protocol.IsNull() {
		return false
	}
	if !data.CertificateId.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

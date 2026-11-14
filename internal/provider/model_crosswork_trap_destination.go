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
type TrapDestination struct {
	Instance                 types.String `tfsdk:"instance"`
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Address                  types.String `tfsdk:"address"`
	AddressType              types.String `tfsdk:"address_type"`
	Port                     types.Int64  `tfsdk:"port"`
	TransportType            types.String `tfsdk:"transport_type"`
	SnmpVersion              types.String `tfsdk:"snmp_version"`
	Community                types.String `tfsdk:"community"`
	CommunityWo              types.String `tfsdk:"community_wo"`
	CommunityWoVersion       types.Int64  `tfsdk:"community_wo_version"`
	Username                 types.String `tfsdk:"username"`
	Mode                     types.String `tfsdk:"mode"`
	AuthType                 types.String `tfsdk:"auth_type"`
	AuthPassword             types.String `tfsdk:"auth_password"`
	AuthPasswordWo           types.String `tfsdk:"auth_password_wo"`
	AuthPasswordWoVersion    types.Int64  `tfsdk:"auth_password_wo_version"`
	PrivacyType              types.String `tfsdk:"privacy_type"`
	PrivacyPassword          types.String `tfsdk:"privacy_password"`
	PrivacyPasswordWo        types.String `tfsdk:"privacy_password_wo"`
	PrivacyPasswordWoVersion types.Int64  `tfsdk:"privacy_password_wo_version"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data TrapDestination) getPath() string {
	return "/crosswork/alarm/v1/destination/trap"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data TrapDestination) getPathGet() string {
	return "/crosswork/alarm/v1/destination/trap"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Written by hand: the controller lists every trap receiver, and the read returns the one with the name (a not found when it has none).
// A receiver created with a host name as its address reads back without a server address: the host is its name.
func (data TrapDestination) readGet(c *cw.Client, path string) (cw.Res, error) {
	name := data.Name.ValueString()
	if name == "" {
		name = data.Id.ValueString()
	}
	res, err := readAlarmReceiver(c, trapDestinationPath, name)
	if err != nil {
		return res, err
	}
	return trapReceiverEntry(res), nil
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data TrapDestination) toBody(ctx context.Context, state TrapDestination) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "name", data.Name.ValueString())
	}
	if !data.Address.IsNull() {
		body, _ = sjson.Set(body, "serverAddress", data.Address.ValueString())
	}
	if !data.AddressType.IsNull() {
		body, _ = sjson.Set(body, "notifyType", data.AddressType.ValueString())
	}
	if !data.Port.IsNull() {
		body, _ = sjson.Set(body, "port", data.Port.ValueInt64())
	}
	if !data.TransportType.IsNull() {
		body, _ = sjson.Set(body, "transportType", data.TransportType.ValueString())
	}
	if !data.SnmpVersion.IsNull() {
		body, _ = sjson.Set(body, "snmpVersion", data.SnmpVersion.ValueString())
	}
	if !data.Community.IsNull() {
		body, _ = sjson.Set(body, "community", data.Community.ValueString())
	}
	if !data.CommunityWo.IsNull() {
		body, _ = sjson.Set(body, "community", data.CommunityWo.ValueString())
	}
	if !data.Username.IsNull() {
		body, _ = sjson.Set(body, "snmpV3UserName", data.Username.ValueString())
	}
	if !data.Mode.IsNull() {
		body, _ = sjson.Set(body, "snmpV3AuthMode", data.Mode.ValueString())
	}
	if !data.AuthType.IsNull() {
		body, _ = sjson.Set(body, "snmpV3AuthType", data.AuthType.ValueString())
	}
	if !data.AuthPassword.IsNull() {
		body, _ = sjson.Set(body, "snmpV3AuthPassphrase", data.AuthPassword.ValueString())
	}
	if !data.AuthPasswordWo.IsNull() {
		body, _ = sjson.Set(body, "snmpV3AuthPassphrase", data.AuthPasswordWo.ValueString())
	}
	if !data.PrivacyType.IsNull() {
		body, _ = sjson.Set(body, "snmpV3PrivType", data.PrivacyType.ValueString())
	}
	if !data.PrivacyPassword.IsNull() {
		body, _ = sjson.Set(body, "snmpV3PrivPassphrase", data.PrivacyPassword.ValueString())
	}
	if !data.PrivacyPasswordWo.IsNull() {
		body, _ = sjson.Set(body, "snmpV3PrivPassphrase", data.PrivacyPasswordWo.ValueString())
	}
	body, _ = sjson.Set(body, "serverType", "NORTH_BOUND")
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *TrapDestination) fromBody(ctx context.Context, res gjson.Result) {
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
	if value := res.Get("notifyType"); value.Exists() {
		data.AddressType = types.StringValue(value.String())
	} else {
		data.AddressType = types.StringValue("IPADDRESS")
	}
	if value := res.Get("port"); value.Exists() {
		data.Port = types.Int64Value(value.Int())
	} else {
		data.Port = types.Int64Value(162)
	}
	if value := res.Get("transportType"); value.Exists() {
		data.TransportType = types.StringValue(value.String())
	} else {
		data.TransportType = types.StringValue("UDP")
	}
	if value := res.Get("snmpVersion"); value.Exists() {
		data.SnmpVersion = types.StringValue(value.String())
	} else {
		data.SnmpVersion = types.StringValue("VERSION_2C")
	}
	if value := res.Get("snmpV3UserName"); value.Exists() {
		data.Username = types.StringValue(value.String())
	} else {
		data.Username = types.StringNull()
	}
	if value := res.Get("snmpV3AuthMode"); value.Exists() {
		data.Mode = types.StringValue(value.String())
	} else {
		data.Mode = types.StringNull()
	}
	if value := res.Get("snmpV3AuthType"); value.Exists() {
		data.AuthType = types.StringValue(value.String())
	} else {
		data.AuthType = types.StringNull()
	}
	if value := res.Get("snmpV3PrivType"); value.Exists() {
		data.PrivacyType = types.StringValue(value.String())
	} else {
		data.PrivacyType = types.StringNull()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *TrapDestination) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("serverAddress"); value.Exists() && !data.Address.IsNull() {
		data.Address = types.StringValue(value.String())
	} else {
		data.Address = types.StringNull()
	}
	if value := res.Get("notifyType"); value.Exists() && !data.AddressType.IsNull() {
		data.AddressType = types.StringValue(value.String())
	} else if data.AddressType.ValueString() != "IPADDRESS" {
		data.AddressType = types.StringNull()
	}
	if value := res.Get("port"); value.Exists() && !data.Port.IsNull() {
		data.Port = types.Int64Value(value.Int())
	} else if data.Port.ValueInt64() != 162 {
		data.Port = types.Int64Null()
	}
	if value := res.Get("transportType"); value.Exists() && !data.TransportType.IsNull() {
		data.TransportType = types.StringValue(value.String())
	} else if data.TransportType.ValueString() != "UDP" {
		data.TransportType = types.StringNull()
	}
	if value := res.Get("snmpVersion"); value.Exists() && !data.SnmpVersion.IsNull() {
		data.SnmpVersion = types.StringValue(value.String())
	} else if data.SnmpVersion.ValueString() != "VERSION_2C" {
		data.SnmpVersion = types.StringNull()
	}
	if value := res.Get("snmpV3UserName"); value.Exists() && !data.Username.IsNull() {
		data.Username = types.StringValue(value.String())
	} else {
		data.Username = types.StringNull()
	}
	if value := res.Get("snmpV3AuthMode"); value.Exists() && !data.Mode.IsNull() {
		data.Mode = types.StringValue(value.String())
	} else {
		data.Mode = types.StringNull()
	}
	if value := res.Get("snmpV3AuthType"); value.Exists() && !data.AuthType.IsNull() {
		data.AuthType = types.StringValue(value.String())
	} else {
		data.AuthType = types.StringNull()
	}
	if value := res.Get("snmpV3PrivType"); value.Exists() && !data.PrivacyType.IsNull() {
		data.PrivacyType = types.StringValue(value.String())
	} else {
		data.PrivacyType = types.StringNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *TrapDestination) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Address.IsNull() {
		return false
	}
	if !data.AddressType.IsNull() {
		return false
	}
	if !data.Port.IsNull() {
		return false
	}
	if !data.TransportType.IsNull() {
		return false
	}
	if !data.SnmpVersion.IsNull() {
		return false
	}
	if !data.Community.IsNull() {
		return false
	}
	if !data.CommunityWo.IsNull() {
		return false
	}
	if !data.CommunityWoVersion.IsNull() {
		return false
	}
	if !data.Username.IsNull() {
		return false
	}
	if !data.Mode.IsNull() {
		return false
	}
	if !data.AuthType.IsNull() {
		return false
	}
	if !data.AuthPassword.IsNull() {
		return false
	}
	if !data.AuthPasswordWo.IsNull() {
		return false
	}
	if !data.AuthPasswordWoVersion.IsNull() {
		return false
	}
	if !data.PrivacyType.IsNull() {
		return false
	}
	if !data.PrivacyPassword.IsNull() {
		return false
	}
	if !data.PrivacyPasswordWo.IsNull() {
		return false
	}
	if !data.PrivacyPasswordWoVersion.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

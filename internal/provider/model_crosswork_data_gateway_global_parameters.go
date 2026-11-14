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

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type DataGatewayGlobalParameters struct {
	Instance            types.String `tfsdk:"instance"`
	Id                  types.String `tfsdk:"id"`
	CliSessions         types.Int64  `tfsdk:"cli_sessions"`
	SshSessionTimeout   types.Int64  `tfsdk:"ssh_session_timeout"`
	SnmpTrapPort        types.Int64  `tfsdk:"snmp_trap_port"`
	SyslogUdpPort       types.Int64  `tfsdk:"syslog_udp_port"`
	SyslogTcpPort       types.Int64  `tfsdk:"syslog_tcp_port"`
	SyslogTlsPort       types.Int64  `tfsdk:"syslog_tls_port"`
	ResyncSnmpv3Details types.Bool   `tfsdk:"resync_snmpv3_details"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data DataGatewayGlobalParameters) getPath() string {
	return "/crosswork/dg-manager/v1/command/global-parameter"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data DataGatewayGlobalParameters) getPathGet() string {
	return "/crosswork/dg-manager/v1/command/global-parameter/query"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Written by hand: the controller lists the parameters as keys with typed values, which are read as one object keyed by the attributes.
func (data DataGatewayGlobalParameters) readGet(c *cw.Client, path string) (cw.Res, error) {
	res, err := c.Post(path, `{}`, cw.ReadOnly)
	if err != nil {
		return res, err
	}
	return flattenGlobalParameters(res), nil
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data DataGatewayGlobalParameters) toBody(ctx context.Context, state DataGatewayGlobalParameters) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if data.CliSessions.ValueInt64() != 0 && !data.CliSessions.IsNull() {
		body, _ = sjson.Set(body, "cli_sessions", data.CliSessions.ValueInt64())
	}
	if data.SshSessionTimeout.ValueInt64() != 0 && !data.SshSessionTimeout.IsNull() {
		body, _ = sjson.Set(body, "ssh_session_timeout", data.SshSessionTimeout.ValueInt64())
	}
	if data.SnmpTrapPort.ValueInt64() != 0 && !data.SnmpTrapPort.IsNull() {
		body, _ = sjson.Set(body, "snmp_trap_port", data.SnmpTrapPort.ValueInt64())
	}
	if data.SyslogUdpPort.ValueInt64() != 0 && !data.SyslogUdpPort.IsNull() {
		body, _ = sjson.Set(body, "syslog_udp_port", data.SyslogUdpPort.ValueInt64())
	}
	if data.SyslogTcpPort.ValueInt64() != 0 && !data.SyslogTcpPort.IsNull() {
		body, _ = sjson.Set(body, "syslog_tcp_port", data.SyslogTcpPort.ValueInt64())
	}
	if data.SyslogTlsPort.ValueInt64() != 0 && !data.SyslogTlsPort.IsNull() {
		body, _ = sjson.Set(body, "syslog_tls_port", data.SyslogTlsPort.ValueInt64())
	}
	if !data.ResyncSnmpv3Details.IsUnknown() && !data.ResyncSnmpv3Details.IsNull() {
		body, _ = sjson.Set(body, "resync_snmpv3_details", data.ResyncSnmpv3Details.ValueBool())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *DataGatewayGlobalParameters) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("data_gateway_global_parameters")
	if value := res.Get("cli_sessions"); value.Exists() {
		data.CliSessions = types.Int64Value(value.Int())
	} else {
		data.CliSessions = types.Int64Null()
	}
	if value := res.Get("ssh_session_timeout"); value.Exists() {
		data.SshSessionTimeout = types.Int64Value(value.Int())
	} else {
		data.SshSessionTimeout = types.Int64Null()
	}
	if value := res.Get("snmp_trap_port"); value.Exists() {
		data.SnmpTrapPort = types.Int64Value(value.Int())
	} else {
		data.SnmpTrapPort = types.Int64Null()
	}
	if value := res.Get("syslog_udp_port"); value.Exists() {
		data.SyslogUdpPort = types.Int64Value(value.Int())
	} else {
		data.SyslogUdpPort = types.Int64Null()
	}
	if value := res.Get("syslog_tcp_port"); value.Exists() {
		data.SyslogTcpPort = types.Int64Value(value.Int())
	} else {
		data.SyslogTcpPort = types.Int64Null()
	}
	if value := res.Get("syslog_tls_port"); value.Exists() {
		data.SyslogTlsPort = types.Int64Value(value.Int())
	} else {
		data.SyslogTlsPort = types.Int64Null()
	}
	if value := res.Get("resync_snmpv3_details"); value.Exists() {
		data.ResyncSnmpv3Details = types.BoolValue(value.Bool())
	} else {
		data.ResyncSnmpv3Details = types.BoolNull()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *DataGatewayGlobalParameters) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cli_sessions"); value.Exists() && !data.CliSessions.IsNull() {
		data.CliSessions = types.Int64Value(value.Int())
	} else {
		data.CliSessions = types.Int64Null()
	}
	if value := res.Get("ssh_session_timeout"); value.Exists() && !data.SshSessionTimeout.IsNull() {
		data.SshSessionTimeout = types.Int64Value(value.Int())
	} else {
		data.SshSessionTimeout = types.Int64Null()
	}
	if value := res.Get("snmp_trap_port"); value.Exists() && !data.SnmpTrapPort.IsNull() {
		data.SnmpTrapPort = types.Int64Value(value.Int())
	} else {
		data.SnmpTrapPort = types.Int64Null()
	}
	if value := res.Get("syslog_udp_port"); value.Exists() && !data.SyslogUdpPort.IsNull() {
		data.SyslogUdpPort = types.Int64Value(value.Int())
	} else {
		data.SyslogUdpPort = types.Int64Null()
	}
	if value := res.Get("syslog_tcp_port"); value.Exists() && !data.SyslogTcpPort.IsNull() {
		data.SyslogTcpPort = types.Int64Value(value.Int())
	} else {
		data.SyslogTcpPort = types.Int64Null()
	}
	if value := res.Get("syslog_tls_port"); value.Exists() && !data.SyslogTlsPort.IsNull() {
		data.SyslogTlsPort = types.Int64Value(value.Int())
	} else {
		data.SyslogTlsPort = types.Int64Null()
	}
	if value := res.Get("resync_snmpv3_details"); value.Exists() && !data.ResyncSnmpv3Details.IsNull() {
		data.ResyncSnmpv3Details = types.BoolValue(value.Bool())
	} else {
		data.ResyncSnmpv3Details = types.BoolNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *DataGatewayGlobalParameters) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.CliSessions.IsNull() {
		return false
	}
	if !data.SshSessionTimeout.IsNull() {
		return false
	}
	if !data.SnmpTrapPort.IsNull() {
		return false
	}
	if !data.SyslogUdpPort.IsNull() {
		return false
	}
	if !data.SyslogTcpPort.IsNull() {
		return false
	}
	if !data.SyslogTlsPort.IsNull() {
		return false
	}
	if !data.ResyncSnmpv3Details.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

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
type SmartLicensingTransport struct {
	Instance  types.String `tfsdk:"instance"`
	Id        types.String `tfsdk:"id"`
	Mode      types.String `tfsdk:"mode"`
	Url       types.String `tfsdk:"url"`
	ProxyHost types.String `tfsdk:"proxy_host"`
	ProxyPort types.Int64  `tfsdk:"proxy_port"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data SmartLicensingTransport) getPath() string {
	return "/crosswork/clms/v1/transport"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data SmartLicensingTransport) getPathGet() string {
	return "/crosswork/clms/v1/get-transport"
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
func (data SmartLicensingTransport) toBody(ctx context.Context, state SmartLicensingTransport) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Mode.IsNull() {
		body, _ = sjson.Set(body, "mode", data.Mode.ValueString())
	}
	if !data.Url.IsNull() {
		body, _ = sjson.Set(body, "url", data.Url.ValueString())
	}
	if !data.ProxyHost.IsNull() {
		body, _ = sjson.Set(body, "proxy.host", data.ProxyHost.ValueString())
	}
	if !data.ProxyPort.IsNull() {
		body, _ = sjson.Set(body, "proxy.port", data.ProxyPort.ValueInt64())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *SmartLicensingTransport) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("smart_licensing_transport")
	if value := res.Get("mode"); value.Exists() {
		data.Mode = types.StringValue(value.String())
	} else {
		data.Mode = types.StringValue("DIRECT")
	}
	if value := res.Get("url"); value.Exists() {
		data.Url = types.StringValue(value.String())
	} else {
		data.Url = types.StringValue("https://smartreceiver.cisco.com/licservice/license")
	}
	if value := res.Get("proxy.host"); value.Exists() {
		data.ProxyHost = types.StringValue(value.String())
	} else {
		data.ProxyHost = types.StringNull()
	}
	if value := res.Get("proxy.port"); value.Exists() {
		data.ProxyPort = types.Int64Value(value.Int())
	} else {
		data.ProxyPort = types.Int64Null()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *SmartLicensingTransport) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("mode"); value.Exists() && !data.Mode.IsNull() {
		data.Mode = types.StringValue(value.String())
	} else if data.Mode.ValueString() != "DIRECT" {
		data.Mode = types.StringNull()
	}
	if value := res.Get("url"); value.Exists() && !data.Url.IsNull() {
		data.Url = types.StringValue(value.String())
	} else if data.Url.ValueString() != "https://smartreceiver.cisco.com/licservice/license" {
		data.Url = types.StringNull()
	}
	if value := res.Get("proxy.host"); value.Exists() && !data.ProxyHost.IsNull() {
		data.ProxyHost = types.StringValue(value.String())
	} else {
		data.ProxyHost = types.StringNull()
	}
	if value := res.Get("proxy.port"); value.Exists() && !data.ProxyPort.IsNull() {
		data.ProxyPort = types.Int64Value(value.Int())
	} else {
		data.ProxyPort = types.Int64Null()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *SmartLicensingTransport) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Mode.IsNull() {
		return false
	}
	if !data.Url.IsNull() {
		return false
	}
	if !data.ProxyHost.IsNull() {
		return false
	}
	if !data.ProxyPort.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

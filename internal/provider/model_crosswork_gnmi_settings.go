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
type GNMISettings struct {
	Instance        types.String `tfsdk:"instance"`
	Id              types.String `tfsdk:"id"`
	CiscoSystems    types.Bool   `tfsdk:"cisco_systems"`
	JuniperNetworks types.Bool   `tfsdk:"juniper_networks"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data GNMISettings) getPath() string {
	return "/crosswork/alarm/v1/gnmi/settings"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data GNMISettings) getPathGet() string {
	return "/crosswork/alarm/v1/gnmi/settings"
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
func (data GNMISettings) toBody(ctx context.Context, state GNMISettings) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.CiscoSystems.IsNull() {
		body, _ = sjson.Set(body, "Cisco Systems", data.CiscoSystems.ValueBool())
	}
	if !data.JuniperNetworks.IsNull() {
		body, _ = sjson.Set(body, "Juniper Networks", data.JuniperNetworks.ValueBool())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *GNMISettings) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("gnmi_settings")
	if value := res.Get("Cisco Systems"); value.Exists() {
		data.CiscoSystems = types.BoolValue(value.Bool())
	} else {
		data.CiscoSystems = types.BoolValue(false)
	}
	if value := res.Get("Juniper Networks"); value.Exists() {
		data.JuniperNetworks = types.BoolValue(value.Bool())
	} else {
		data.JuniperNetworks = types.BoolValue(true)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *GNMISettings) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("Cisco Systems"); value.Exists() && !data.CiscoSystems.IsNull() {
		data.CiscoSystems = types.BoolValue(value.Bool())
	} else if data.CiscoSystems.ValueBool() != false {
		data.CiscoSystems = types.BoolNull()
	}
	if value := res.Get("Juniper Networks"); value.Exists() && !data.JuniperNetworks.IsNull() {
		data.JuniperNetworks = types.BoolValue(value.Bool())
	} else if data.JuniperNetworks.ValueBool() != true {
		data.JuniperNetworks = types.BoolNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *GNMISettings) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.CiscoSystems.IsNull() {
		return false
	}
	if !data.JuniperNetworks.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

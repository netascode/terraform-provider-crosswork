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
type TopologyDiscoverySettings struct {
	Instance types.String `tfsdk:"instance"`
	Id       types.String `tfsdk:"id"`
	Lldp     types.Bool   `tfsdk:"lldp"`
	Cdp      types.Bool   `tfsdk:"cdp"`
	Lag      types.Bool   `tfsdk:"lag"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data TopologyDiscoverySettings) getPath() string {
	return "/crosswork/topology/v1/topology-service/topology/provider/settings/protocols"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

func (data TopologyDiscoverySettings) getPathDelete() string {
	return "/crosswork/topology/v1/topology-service/topology/provider/settings/protocols/defaults"
}

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data TopologyDiscoverySettings) getPathGet() string {
	return "/crosswork/topology/v1/topology-service/topology/provider/settings/protocols"
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
func (data TopologyDiscoverySettings) toBody(ctx context.Context, state TopologyDiscoverySettings) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Lldp.IsNull() {
		body, _ = sjson.Set(body, "lldpConfig.enabled", data.Lldp.ValueBool())
	}
	if !data.Cdp.IsNull() {
		body, _ = sjson.Set(body, "cdpConfig.enabled", data.Cdp.ValueBool())
	}
	if !data.Lag.IsNull() {
		body, _ = sjson.Set(body, "lagConfig.enabled", data.Lag.ValueBool())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *TopologyDiscoverySettings) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("topology_discovery_settings")
	if value := res.Get("providerConfig.lldpConfig.enabled"); value.Exists() {
		data.Lldp = types.BoolValue(value.Bool())
	} else {
		data.Lldp = types.BoolValue(false)
	}
	if value := res.Get("providerConfig.cdpConfig.enabled"); value.Exists() {
		data.Cdp = types.BoolValue(value.Bool())
	} else {
		data.Cdp = types.BoolValue(false)
	}
	if value := res.Get("providerConfig.lagConfig.enabled"); value.Exists() {
		data.Lag = types.BoolValue(value.Bool())
	} else {
		data.Lag = types.BoolValue(false)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *TopologyDiscoverySettings) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("providerConfig.lldpConfig.enabled"); value.Exists() && !data.Lldp.IsNull() {
		data.Lldp = types.BoolValue(value.Bool())
	} else if data.Lldp.ValueBool() != false {
		data.Lldp = types.BoolNull()
	}
	if value := res.Get("providerConfig.cdpConfig.enabled"); value.Exists() && !data.Cdp.IsNull() {
		data.Cdp = types.BoolValue(value.Bool())
	} else if data.Cdp.ValueBool() != false {
		data.Cdp = types.BoolNull()
	}
	if value := res.Get("providerConfig.lagConfig.enabled"); value.Exists() && !data.Lag.IsNull() {
		data.Lag = types.BoolValue(value.Bool())
	} else if data.Lag.ValueBool() != false {
		data.Lag = types.BoolNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *TopologyDiscoverySettings) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Lldp.IsNull() {
		return false
	}
	if !data.Cdp.IsNull() {
		return false
	}
	if !data.Lag.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

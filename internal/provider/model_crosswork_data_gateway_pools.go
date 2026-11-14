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
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type DataGatewayPools struct {
	Instance            types.String `tfsdk:"instance"`
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Description         types.String `tfsdk:"description"`
	VirtualIps          types.List   `tfsdk:"virtual_ips"`
	PrefixLength        types.Int64  `tfsdk:"prefix_length"`
	Spares              types.Int64  `tfsdk:"spares"`
	Instances           types.List   `tfsdk:"instances"`
	CdgReadinessTimeout types.Int64  `tfsdk:"cdg_readiness_timeout"`
	DataGatewayIds      types.List   `tfsdk:"data_gateway_ids"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data DataGatewayPools) getPath() string {
	return "/crosswork/dg-manager/v2/hapool"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data DataGatewayPools) getPathGet() string {
	return "/crosswork/dg-manager/v2/hapool/query"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// Escapes the id so it cannot change the structure of a request body.
func (data DataGatewayPools) idJSON() string {
	id := data.Id.ValueString()
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}

func (data DataGatewayPools) readGet(c *cw.Client, path string) (cw.Res, error) {
	return c.Post(path, strings.ReplaceAll(`{}`, "{id}", data.idJSON()), cw.ReadOnly)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Written by hand: the body mixes key casings and sends numbers as strings.
func (data DataGatewayPools) toBody(ctx context.Context, state DataGatewayPools) string {
	return poolRequestBody(data, state.Id.ValueString())
}

// End of section. //template:end toBody

// Written by hand: the numbers come back as strings, the spares are absent at 0, and the addresses are one per entry.
func (data *DataGatewayPools) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("puuid"); value.Exists() {
		data.Id = types.StringValue(value.String())
	}
	data.Name = types.StringNull()
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	}
	data.Description = types.StringNull()
	if value := res.Get("description"); value.Exists() && value.String() != "" {
		data.Description = types.StringValue(value.String())
	}
	addresses := poolVirtualIps(res)
	for i := range addresses {
		addresses[i] = normalizeAddress(addresses[i])
	}
	data.VirtualIps = stringList(addresses)
	data.PrefixLength = types.Int64Value(poolPrefixLength(res))
	data.Spares = types.Int64Value(res.Get("spares").Int())
	data.Instances = stringList(poolInstanceIds(res))
	if data.CdgReadinessTimeout.IsNull() || data.CdgReadinessTimeout.IsUnknown() {
		data.CdgReadinessTimeout = types.Int64Value(defaultCdgReadinessTimeout)
	}
}

// End of section. //template:end fromBody

// Written by hand: an address is kept as the user wrote it when it is the same address, the spares are absent at 0, and the
// instances and addresses are the whole lists the controller holds.
func (data *DataGatewayPools) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	}
	if value := res.Get("description"); value.Exists() && !data.Description.IsNull() {
		data.Description = types.StringValue(value.String())
	} else if !value.Exists() && !data.Description.IsNull() && data.Description.ValueString() != "" {
		data.Description = types.StringNull()
	}
	if addresses := poolVirtualIps(res); !data.VirtualIps.IsNull() && !sameAddresses(listOfStrings(data.VirtualIps), addresses) {
		for i := range addresses {
			addresses[i] = normalizeAddress(addresses[i])
		}
		data.VirtualIps = stringList(addresses)
	}
	if !data.PrefixLength.IsNull() {
		data.PrefixLength = types.Int64Value(poolPrefixLength(res))
	}
	if !data.Spares.IsNull() {
		data.Spares = types.Int64Value(res.Get("spares").Int())
	}
	if !data.Instances.IsNull() {
		data.Instances = stringList(poolInstanceIds(res))
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *DataGatewayPools) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Name.IsNull() {
		return false
	}
	if !data.Description.IsNull() {
		return false
	}
	if !data.VirtualIps.IsNull() {
		return false
	}
	if !data.PrefixLength.IsNull() {
		return false
	}
	if !data.Spares.IsNull() {
		return false
	}
	if !data.Instances.IsNull() {
		return false
	}
	if !data.CdgReadinessTimeout.IsNull() {
		return false
	}
	if !data.DataGatewayIds.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

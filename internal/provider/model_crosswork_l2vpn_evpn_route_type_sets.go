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
	"fmt"
	"net/url"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type L2VPNEVPNRouteTypeSets struct {
	Instance      types.String `tfsdk:"instance"`
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	EvpnRouteType types.Set    `tfsdk:"evpn_route_type"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data L2VPNEVPNRouteTypeSets) getPath() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/cisco-l2vpn-routing-policy:l2vpn-routing-policy/defined-sets/evpn-route-type-sets/evpn-route-type-set=%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

func (data L2VPNEVPNRouteTypeSets) getPathPost() string {
	return "/crosswork/proxy/nso/restconf/data/cisco-l2vpn-routing-policy:l2vpn-routing-policy/defined-sets/evpn-route-type-sets"
}

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
func (data L2VPNEVPNRouteTypeSets) toBody(ctx context.Context, state L2VPNEVPNRouteTypeSets) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-routing-policy:evpn-route-type-set.0.name", data.Name.ValueString())
	}
	if !data.EvpnRouteType.IsNull() {
		var values []int64
		data.EvpnRouteType.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "cisco-l2vpn-routing-policy:evpn-route-type-set.0.evpn-route-type", values)
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *L2VPNEVPNRouteTypeSets) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-l2vpn-routing-policy:evpn-route-type-set.0.evpn-route-type"); value.Exists() && len(value.Array()) > 0 {
		data.EvpnRouteType = helpers.GetInt64Set(value.Array())
	} else {
		data.EvpnRouteType = types.SetNull(types.Int64Type)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *L2VPNEVPNRouteTypeSets) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-l2vpn-routing-policy:evpn-route-type-set.0.evpn-route-type"); value.Exists() && !data.EvpnRouteType.IsNull() {
		data.EvpnRouteType = helpers.GetInt64Set(value.Array())
	} else {
		data.EvpnRouteType = types.SetNull(types.Int64Type)
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *L2VPNEVPNRouteTypeSets) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.EvpnRouteType.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

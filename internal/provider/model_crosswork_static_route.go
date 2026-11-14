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

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type StaticRoute struct {
	Instance     types.String `tfsdk:"instance"`
	Id           types.String `tfsdk:"id"`
	Address      types.String `tfsdk:"address"`
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data StaticRoute) getPath() string {
	return "/crosswork/ztp/v1/staticroutes"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data StaticRoute) getPathGet() string {
	return "/crosswork/ztp/v1/staticroutes/query"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Written by hand: the controller lists all routes, and the read returns the one of this resource (a not found when it has none).
func (data StaticRoute) readGet(c *cw.Client, path string) (cw.Res, error) {
	return data.readRoute(c)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data StaticRoute) toBody(ctx context.Context, state StaticRoute) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Address.IsNull() {
		body, _ = sjson.Set(body, "subnet", data.Address.ValueString())
	}
	if !data.PrefixLength.IsNull() {
		body, _ = sjson.Set(body, "mask", fmt.Sprint(data.PrefixLength.ValueInt64()))
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *StaticRoute) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("subnet"); value.Exists() {
		data.Address = types.StringValue(value.String())
	} else {
		data.Address = types.StringNull()
	}
	if value := res.Get("mask"); value.Exists() {
		data.PrefixLength = types.Int64Value(value.Int())
	} else {
		data.PrefixLength = types.Int64Null()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *StaticRoute) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("subnet"); value.Exists() && !data.Address.IsNull() {
		data.Address = types.StringValue(value.String())
	} else {
		data.Address = types.StringNull()
	}
	if value := res.Get("mask"); value.Exists() && !data.PrefixLength.IsNull() {
		data.PrefixLength = types.Int64Value(value.Int())
	} else {
		data.PrefixLength = types.Int64Null()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *StaticRoute) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Address.IsNull() {
		return false
	}
	if !data.PrefixLength.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

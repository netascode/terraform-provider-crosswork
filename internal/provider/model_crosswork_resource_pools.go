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

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type ResourcePools struct {
	Instance   types.String `tfsdk:"instance"`
	Id         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	RangeStart types.Int64  `tfsdk:"range_start"`
	RangeEnd   types.Int64  `tfsdk:"range_end"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data ResourcePools) getPath() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/resource-allocator:resource-pools/id-allocator:id-pool=%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

func (data ResourcePools) getPathPost() string {
	return "/crosswork/proxy/nso/restconf/data/resource-allocator:resource-pools"
}

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

func (data ResourcePools) getPathPut() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/resource-allocator:resource-pools/id-allocator:id-pool=%v/range", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data ResourcePools) toBody(ctx context.Context, state ResourcePools) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() && put == false {
		body, _ = sjson.Set(body, "id-allocator:id-pool.0.name", data.Name.ValueString())
	}
	if !data.RangeStart.IsNull() {
		if put {
			body, _ = sjson.Set(body, "id-allocator:range.start", data.RangeStart.ValueInt64())
		} else {
			body, _ = sjson.Set(body, "id-allocator:id-pool.0.range.start", data.RangeStart.ValueInt64())
		}
	}
	if !data.RangeEnd.IsNull() {
		if put {
			body, _ = sjson.Set(body, "id-allocator:range.end", data.RangeEnd.ValueInt64())
		} else {
			body, _ = sjson.Set(body, "id-allocator:id-pool.0.range.end", data.RangeEnd.ValueInt64())
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *ResourcePools) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("id-allocator:id-pool.0.range.start"); value.Exists() {
		data.RangeStart = types.Int64Value(value.Int())
	} else {
		data.RangeStart = types.Int64Null()
	}
	if value := res.Get("id-allocator:id-pool.0.range.end"); value.Exists() {
		data.RangeEnd = types.Int64Value(value.Int())
	} else {
		data.RangeEnd = types.Int64Null()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *ResourcePools) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("id-allocator:id-pool.0.range.start"); value.Exists() && !data.RangeStart.IsNull() {
		data.RangeStart = types.Int64Value(value.Int())
	} else {
		data.RangeStart = types.Int64Null()
	}
	if value := res.Get("id-allocator:id-pool.0.range.end"); value.Exists() && !data.RangeEnd.IsNull() {
		data.RangeEnd = types.Int64Value(value.Int())
	} else {
		data.RangeEnd = types.Int64Null()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *ResourcePools) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.RangeStart.IsNull() {
		return false
	}
	if !data.RangeEnd.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

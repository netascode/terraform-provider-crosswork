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
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type Tags struct {
	Instance types.String `tfsdk:"instance"`
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Category types.String `tfsdk:"category"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data Tags) getPath() string {
	return "/crosswork/inventory/v1/tags"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data Tags) getPathGet() string {
	return "/crosswork/inventory/v1/tags/query"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// Escapes the id so it cannot change the structure of a request body.
func (data Tags) idJSON() string {
	id := data.Id.ValueString()
	if id == "" {
		id = fmt.Sprint(data.Name.ValueString())
	}
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}

func (data Tags) readGet(c *cw.Client, path string) (cw.Res, error) {
	return c.Post(path, strings.ReplaceAll(`{"filter":{"name":"{id}"},"next_from":"0"}`, "{id}", data.idJSON()), cw.ReadOnly)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data Tags) toBody(ctx context.Context, state Tags) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "tags.0.name", data.Name.ValueString())
	}
	if !data.Category.IsNull() {
		body, _ = sjson.Set(body, "tags.0.category", data.Category.ValueString())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *Tags) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("tags.0.name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("tags.0.category"); value.Exists() {
		data.Category = types.StringValue(value.String())
	} else {
		data.Category = types.StringNull()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *Tags) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("tags.0.name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("tags.0.category"); value.Exists() && !data.Category.IsNull() {
		data.Category = types.StringValue(value.String())
	} else {
		data.Category = types.StringNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *Tags) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Category.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

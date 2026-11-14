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
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type LockUsers struct {
	Instance types.String `tfsdk:"instance"`
	Id       types.String `tfsdk:"id"`
	Username types.String `tfsdk:"username"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data LockUsers) getPath() string {
	return "/crosswork/aaa/v1/lockUsers"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

func (data LockUsers) getPathDelete() string {
	return "/crosswork/aaa/v1/unlockUsers"
}

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data LockUsers) getPathGet() string {
	return "/crosswork/aaa/v1/user"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// Escapes the id so it cannot change the structure of a request body.
func (data LockUsers) idJSON() string {
	id := data.Id.ValueString()
	if id == "" {
		id = fmt.Sprint(data.Username.ValueString())
	}
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}

// postBody places the generated create body in the request template, in one pass so a value that
// looks like a placeholder (a password such as "p{id}w") is never expanded twice.
func (data LockUsers) postBody(body string) string {
	return strings.NewReplacer("{id}", data.idJSON(), "{body}", body).Replace(`{"{id}": {body}}`)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data LockUsers) toBody(ctx context.Context, state LockUsers) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Username.IsNull() {
		body, _ = sjson.Set(body, "userName", data.Username.ValueString())
	}
	body, _ = sjson.Set(body, "Status", "Locked")
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *LockUsers) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("Username"); value.Exists() {
		data.Username = types.StringValue(value.String())
	} else {
		data.Username = types.StringNull()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *LockUsers) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("Username"); value.Exists() && !data.Username.IsNull() {
		data.Username = types.StringValue(value.String())
	} else {
		data.Username = types.StringNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *LockUsers) isNull(ctx context.Context, res gjson.Result) bool {
	return true
}

// End of section. //template:end isNull

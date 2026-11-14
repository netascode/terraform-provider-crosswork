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
type Users struct {
	Instance           types.String              `tfsdk:"instance"`
	Id                 types.String              `tfsdk:"id"`
	Username           types.String              `tfsdk:"username"`
	Password           types.String              `tfsdk:"password"`
	PasswordWo         types.String              `tfsdk:"password_wo"`
	PasswordWoVersion  types.Int64               `tfsdk:"password_wo_version"`
	PolicyId           types.String              `tfsdk:"policy_id"`
	FirstName          types.String              `tfsdk:"first_name"`
	LastName           types.String              `tfsdk:"last_name"`
	DeviceAccessGroups []UsersDeviceAccessGroups `tfsdk:"device_access_groups"`
}

type UsersDeviceAccessGroups struct {
	Uuid       types.String `tfsdk:"uuid"`
	DomainName types.String `tfsdk:"domain_name"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data Users) getPath() string {
	return "/crosswork/aaa/v1/user"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// Escapes the id so it cannot change the structure of a request body.
func (data Users) idJSON() string {
	id := data.Id.ValueString()
	if id == "" {
		id = fmt.Sprint(data.Username.ValueString())
	}
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}

// postBody places the generated create body in the request template, in one pass so a value that
// looks like a placeholder (a password such as "p{id}w") is never expanded twice.
func (data Users) postBody(body string) string {
	return strings.NewReplacer("{id}", data.idJSON(), "{body}", body).Replace(`{"{id}": {body}}`)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data Users) toBody(ctx context.Context, state Users) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Username.IsNull() {
		body, _ = sjson.Set(body, "Username", data.Username.ValueString())
	}
	if !data.Password.IsNull() {
		body, _ = sjson.Set(body, "Password", data.Password.ValueString())
	}
	if !data.PasswordWo.IsNull() {
		body, _ = sjson.Set(body, "Password", data.PasswordWo.ValueString())
	}
	if !data.PolicyId.IsNull() {
		body, _ = sjson.Set(body, "PolicyId", data.PolicyId.ValueString())
	}
	if !data.FirstName.IsNull() {
		body, _ = sjson.Set(body, "FirstName", data.FirstName.ValueString())
	}
	if !data.LastName.IsNull() {
		body, _ = sjson.Set(body, "LastName", data.LastName.ValueString())
	}
	if len(data.DeviceAccessGroups) > 0 {
		body, _ = sjson.Set(body, "DeviceAccessGroups", []interface{}{})
		for _, item := range data.DeviceAccessGroups {
			itemBody := ""
			if !item.Uuid.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "Uuid", item.Uuid.ValueString())
			}
			if !item.DomainName.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "DomainName", item.DomainName.ValueString())
			}
			body, _ = sjson.SetRaw(body, "DeviceAccessGroups.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *Users) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("Username"); value.Exists() {
		data.Username = types.StringValue(value.String())
	} else {
		data.Username = types.StringNull()
	}
	if value := res.Get("PolicyId"); value.Exists() {
		data.PolicyId = types.StringValue(value.String())
	} else {
		data.PolicyId = types.StringNull()
	}
	if value := res.Get("FirstName"); value.Exists() {
		data.FirstName = types.StringValue(value.String())
	} else {
		data.FirstName = types.StringNull()
	}
	if value := res.Get("LastName"); value.Exists() {
		data.LastName = types.StringValue(value.String())
	} else {
		data.LastName = types.StringNull()
	}
	if value := res.Get("DeviceAccessGroups"); value.Exists() && len(value.Array()) > 0 {
		data.DeviceAccessGroups = make([]UsersDeviceAccessGroups, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := UsersDeviceAccessGroups{}
			if cValue := v.Get("Uuid"); cValue.Exists() {
				item.Uuid = types.StringValue(cValue.String())
			} else {
				item.Uuid = types.StringNull()
			}
			if cValue := v.Get("DomainName"); cValue.Exists() {
				item.DomainName = types.StringValue(cValue.String())
			} else {
				item.DomainName = types.StringNull()
			}
			data.DeviceAccessGroups = append(data.DeviceAccessGroups, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *Users) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("Username"); value.Exists() && !data.Username.IsNull() {
		data.Username = types.StringValue(value.String())
	} else {
		data.Username = types.StringNull()
	}
	if value := res.Get("PolicyId"); value.Exists() && !data.PolicyId.IsNull() {
		data.PolicyId = types.StringValue(value.String())
	} else {
		data.PolicyId = types.StringNull()
	}
	if value := res.Get("FirstName"); value.Exists() && !data.FirstName.IsNull() {
		data.FirstName = types.StringValue(value.String())
	} else {
		data.FirstName = types.StringNull()
	}
	if value := res.Get("LastName"); value.Exists() && !data.LastName.IsNull() {
		data.LastName = types.StringValue(value.String())
	} else {
		data.LastName = types.StringNull()
	}
	for i := range data.DeviceAccessGroups {
		keys := [...]string{"Uuid"}
		keyValues := [...]string{data.DeviceAccessGroups[i].Uuid.ValueString()}

		var r gjson.Result
		res.Get("DeviceAccessGroups").ForEach(
			func(_, v gjson.Result) bool {
				found := false
				for ik := range keys {
					if v.Get(keys[ik]).String() == keyValues[ik] {
						found = true
						continue
					}
					found = false
					break
				}
				if found {
					r = v
					return false
				}
				return true
			},
		)
		if value := r.Get("Uuid"); value.Exists() && !data.DeviceAccessGroups[i].Uuid.IsNull() {
			data.DeviceAccessGroups[i].Uuid = types.StringValue(value.String())
		} else {
			data.DeviceAccessGroups[i].Uuid = types.StringNull()
		}
		if value := r.Get("DomainName"); value.Exists() && !data.DeviceAccessGroups[i].DomainName.IsNull() {
			data.DeviceAccessGroups[i].DomainName = types.StringValue(value.String())
		} else {
			data.DeviceAccessGroups[i].DomainName = types.StringNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *Users) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Password.IsNull() {
		return false
	}
	if !data.PasswordWo.IsNull() {
		return false
	}
	if !data.PasswordWoVersion.IsNull() {
		return false
	}
	if !data.PolicyId.IsNull() {
		return false
	}
	if !data.FirstName.IsNull() {
		return false
	}
	if !data.LastName.IsNull() {
		return false
	}
	if len(data.DeviceAccessGroups) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

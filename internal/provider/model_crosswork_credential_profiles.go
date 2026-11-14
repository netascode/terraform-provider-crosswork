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
type CredentialProfiles struct {
	Instance                      types.String                    `tfsdk:"instance"`
	Id                            types.String                    `tfsdk:"id"`
	Profile                       types.String                    `tfsdk:"profile"`
	Credentials                   []CredentialProfilesCredentials `tfsdk:"credentials"`
	Snmpv2ReadCommunity           types.String                    `tfsdk:"snmpv2_read_community"`
	Snmpv2ReadCommunityWo         types.String                    `tfsdk:"snmpv2_read_community_wo"`
	Snmpv2ReadCommunityWoVersion  types.Int64                     `tfsdk:"snmpv2_read_community_wo_version"`
	Snmpv2WriteCommunity          types.String                    `tfsdk:"snmpv2_write_community"`
	Snmpv2WriteCommunityWo        types.String                    `tfsdk:"snmpv2_write_community_wo"`
	Snmpv2WriteCommunityWoVersion types.Int64                     `tfsdk:"snmpv2_write_community_wo_version"`
	Snmpv3UserName                types.String                    `tfsdk:"snmpv3_user_name"`
	Snmpv3SecurityLevel           types.String                    `tfsdk:"snmpv3_security_level"`
	Snmpv3AuthType                types.String                    `tfsdk:"snmpv3_auth_type"`
	Snmpv3AuthPassword            types.String                    `tfsdk:"snmpv3_auth_password"`
	Snmpv3AuthPasswordWo          types.String                    `tfsdk:"snmpv3_auth_password_wo"`
	Snmpv3AuthPasswordWoVersion   types.Int64                     `tfsdk:"snmpv3_auth_password_wo_version"`
	Snmpv3PrivType                types.String                    `tfsdk:"snmpv3_priv_type"`
	Snmpv3PrivPassword            types.String                    `tfsdk:"snmpv3_priv_password"`
	Snmpv3PrivPasswordWo          types.String                    `tfsdk:"snmpv3_priv_password_wo"`
	Snmpv3PrivPasswordWoVersion   types.Int64                     `tfsdk:"snmpv3_priv_password_wo_version"`
}

type CredentialProfilesCredentials struct {
	Type                         types.String `tfsdk:"type"`
	UserName                     types.String `tfsdk:"user_name"`
	Password                     types.String `tfsdk:"password"`
	PasswordWo                   types.String `tfsdk:"password_wo"`
	PasswordWoVersion            types.Int64  `tfsdk:"password_wo_version"`
	EnablePassword               types.String `tfsdk:"enable_password"`
	EnablePasswordValue          types.String `tfsdk:"enable_password_value"`
	EnablePasswordValueWo        types.String `tfsdk:"enable_password_value_wo"`
	EnablePasswordValueWoVersion types.Int64  `tfsdk:"enable_password_value_wo_version"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data CredentialProfiles) getPath() string {
	return "/crosswork/inventory/v1/credentials"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data CredentialProfiles) getPathGet() string {
	return "/crosswork/inventory/v1/credentials/query"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// Escapes the id so it cannot change the structure of a request body.
func (data CredentialProfiles) idJSON() string {
	id := data.Id.ValueString()
	if id == "" {
		id = fmt.Sprint(data.Profile.ValueString())
	}
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}

func (data CredentialProfiles) readGet(c *cw.Client, path string) (cw.Res, error) {
	return c.Post(path, strings.ReplaceAll(`{"filter":{"profile":"{id}"}}`, "{id}", data.idJSON()), cw.ReadOnly)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data CredentialProfiles) toBody(ctx context.Context, state CredentialProfiles) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Profile.IsNull() {
		body, _ = sjson.Set(body, "data.0.profile", data.Profile.ValueString())
	}
	if len(data.Credentials) > 0 {
		body, _ = sjson.Set(body, "data.0.user_pass", []interface{}{})
		for _, item := range data.Credentials {
			itemBody := ""
			if !item.Type.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "type", item.Type.ValueString())
			}
			if !item.UserName.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "user_name", item.UserName.ValueString())
			}
			if !item.Password.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "password", item.Password.ValueString())
			}
			if !item.PasswordWo.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "password", item.PasswordWo.ValueString())
			}
			if !item.EnablePassword.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "enable_password", item.EnablePassword.ValueString())
			}
			if !item.EnablePasswordValue.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "enable_password_data", item.EnablePasswordValue.ValueString())
			}
			if !item.EnablePasswordValueWo.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "enable_password_data", item.EnablePasswordValueWo.ValueString())
			}
			body, _ = sjson.SetRaw(body, "data.0.user_pass.-1", itemBody)
		}
	}
	if !data.Snmpv2ReadCommunity.IsNull() {
		body, _ = sjson.Set(body, "data.0.v2_info.read_community", data.Snmpv2ReadCommunity.ValueString())
	}
	if !data.Snmpv2ReadCommunityWo.IsNull() {
		body, _ = sjson.Set(body, "data.0.v2_info.read_community", data.Snmpv2ReadCommunityWo.ValueString())
	}
	if !data.Snmpv2WriteCommunity.IsNull() {
		body, _ = sjson.Set(body, "data.0.v2_info.write_community", data.Snmpv2WriteCommunity.ValueString())
	}
	if !data.Snmpv2WriteCommunityWo.IsNull() {
		body, _ = sjson.Set(body, "data.0.v2_info.write_community", data.Snmpv2WriteCommunityWo.ValueString())
	}
	if !data.Snmpv3UserName.IsNull() {
		body, _ = sjson.Set(body, "data.0.v3_info.user_name", data.Snmpv3UserName.ValueString())
	}
	if !data.Snmpv3SecurityLevel.IsNull() {
		body, _ = sjson.Set(body, "data.0.v3_info.security_level", data.Snmpv3SecurityLevel.ValueString())
	}
	if !data.Snmpv3AuthType.IsNull() {
		body, _ = sjson.Set(body, "data.0.v3_info.auth_type", data.Snmpv3AuthType.ValueString())
	}
	if !data.Snmpv3AuthPassword.IsNull() {
		body, _ = sjson.Set(body, "data.0.v3_info.auth_password", data.Snmpv3AuthPassword.ValueString())
	}
	if !data.Snmpv3AuthPasswordWo.IsNull() {
		body, _ = sjson.Set(body, "data.0.v3_info.auth_password", data.Snmpv3AuthPasswordWo.ValueString())
	}
	if !data.Snmpv3PrivType.IsNull() {
		body, _ = sjson.Set(body, "data.0.v3_info.priv_type", data.Snmpv3PrivType.ValueString())
	}
	if !data.Snmpv3PrivPassword.IsNull() {
		body, _ = sjson.Set(body, "data.0.v3_info.priv_password", data.Snmpv3PrivPassword.ValueString())
	}
	if !data.Snmpv3PrivPasswordWo.IsNull() {
		body, _ = sjson.Set(body, "data.0.v3_info.priv_password", data.Snmpv3PrivPasswordWo.ValueString())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *CredentialProfiles) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("data.0.profile"); value.Exists() {
		data.Profile = types.StringValue(value.String())
	} else {
		data.Profile = types.StringNull()
	}
	if value := res.Get("data.0.user_pass"); value.Exists() && len(value.Array()) > 0 {
		data.Credentials = make([]CredentialProfilesCredentials, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := CredentialProfilesCredentials{}
			if cValue := v.Get("type"); cValue.Exists() {
				item.Type = types.StringValue(cValue.String())
			} else {
				item.Type = types.StringNull()
			}
			if cValue := v.Get("user_name"); cValue.Exists() {
				item.UserName = types.StringValue(cValue.String())
			} else {
				item.UserName = types.StringNull()
			}
			if cValue := v.Get("enable_password"); cValue.Exists() {
				item.EnablePassword = types.StringValue(cValue.String())
			} else {
				item.EnablePassword = types.StringNull()
			}
			data.Credentials = append(data.Credentials, item)
			return true
		})
	}
	if value := res.Get("data.0.v3_info.user_name"); value.Exists() {
		data.Snmpv3UserName = types.StringValue(value.String())
	} else {
		data.Snmpv3UserName = types.StringNull()
	}
	if value := res.Get("data.0.v3_info.security_level"); value.Exists() {
		data.Snmpv3SecurityLevel = types.StringValue(value.String())
	} else {
		data.Snmpv3SecurityLevel = types.StringNull()
	}
	if value := res.Get("data.0.v3_info.auth_type"); value.Exists() {
		data.Snmpv3AuthType = types.StringValue(value.String())
	} else {
		data.Snmpv3AuthType = types.StringNull()
	}
	if value := res.Get("data.0.v3_info.priv_type"); value.Exists() {
		data.Snmpv3PrivType = types.StringValue(value.String())
	} else {
		data.Snmpv3PrivType = types.StringNull()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *CredentialProfiles) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("data.0.profile"); value.Exists() && !data.Profile.IsNull() {
		data.Profile = types.StringValue(value.String())
	} else {
		data.Profile = types.StringNull()
	}
	for i := range data.Credentials {
		keys := [...]string{"type"}
		keyValues := [...]string{data.Credentials[i].Type.ValueString()}

		var r gjson.Result
		res.Get("data.0.user_pass").ForEach(
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
		if value := r.Get("type"); value.Exists() && !data.Credentials[i].Type.IsNull() {
			data.Credentials[i].Type = types.StringValue(value.String())
		} else {
			data.Credentials[i].Type = types.StringNull()
		}
		if value := r.Get("user_name"); value.Exists() && !data.Credentials[i].UserName.IsNull() {
			data.Credentials[i].UserName = types.StringValue(value.String())
		} else {
			data.Credentials[i].UserName = types.StringNull()
		}
		if value := r.Get("enable_password"); value.Exists() && !data.Credentials[i].EnablePassword.IsNull() {
			data.Credentials[i].EnablePassword = types.StringValue(value.String())
		} else {
			data.Credentials[i].EnablePassword = types.StringNull()
		}
	}
	if value := res.Get("data.0.v3_info.user_name"); value.Exists() && !data.Snmpv3UserName.IsNull() {
		data.Snmpv3UserName = types.StringValue(value.String())
	} else {
		data.Snmpv3UserName = types.StringNull()
	}
	if value := res.Get("data.0.v3_info.security_level"); value.Exists() && !data.Snmpv3SecurityLevel.IsNull() {
		data.Snmpv3SecurityLevel = types.StringValue(value.String())
	} else {
		data.Snmpv3SecurityLevel = types.StringNull()
	}
	if value := res.Get("data.0.v3_info.auth_type"); value.Exists() && !data.Snmpv3AuthType.IsNull() {
		data.Snmpv3AuthType = types.StringValue(value.String())
	} else {
		data.Snmpv3AuthType = types.StringNull()
	}
	if value := res.Get("data.0.v3_info.priv_type"); value.Exists() && !data.Snmpv3PrivType.IsNull() {
		data.Snmpv3PrivType = types.StringValue(value.String())
	} else {
		data.Snmpv3PrivType = types.StringNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *CredentialProfiles) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.Credentials) > 0 {
		return false
	}
	if !data.Snmpv2ReadCommunity.IsNull() {
		return false
	}
	if !data.Snmpv2ReadCommunityWo.IsNull() {
		return false
	}
	if !data.Snmpv2ReadCommunityWoVersion.IsNull() {
		return false
	}
	if !data.Snmpv2WriteCommunity.IsNull() {
		return false
	}
	if !data.Snmpv2WriteCommunityWo.IsNull() {
		return false
	}
	if !data.Snmpv2WriteCommunityWoVersion.IsNull() {
		return false
	}
	if !data.Snmpv3UserName.IsNull() {
		return false
	}
	if !data.Snmpv3SecurityLevel.IsNull() {
		return false
	}
	if !data.Snmpv3AuthType.IsNull() {
		return false
	}
	if !data.Snmpv3AuthPassword.IsNull() {
		return false
	}
	if !data.Snmpv3AuthPasswordWo.IsNull() {
		return false
	}
	if !data.Snmpv3AuthPasswordWoVersion.IsNull() {
		return false
	}
	if !data.Snmpv3PrivType.IsNull() {
		return false
	}
	if !data.Snmpv3PrivPassword.IsNull() {
		return false
	}
	if !data.Snmpv3PrivPasswordWo.IsNull() {
		return false
	}
	if !data.Snmpv3PrivPasswordWoVersion.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

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

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type AAALDAPServers struct {
	Instance         types.String                     `tfsdk:"instance"`
	Id               types.String                     `tfsdk:"id"`
	Servers          []AAALDAPServersServers          `tfsdk:"servers"`
	AttributeServers []AAALDAPServersAttributeServers `tfsdk:"attribute_servers"`
}

type AAALDAPServersServers struct {
	Name                    types.String `tfsdk:"name"`
	LdapUrl                 types.String `tfsdk:"ldap_url"`
	BindDn                  types.String `tfsdk:"bind_dn"`
	BindCredential          types.String `tfsdk:"bind_credential"`
	BindCredentialWo        types.String `tfsdk:"bind_credential_wo"`
	BindCredentialWoVersion types.Int64  `tfsdk:"bind_credential_wo_version"`
	BaseDn                  types.String `tfsdk:"base_dn"`
	UserFilter              types.String `tfsdk:"user_filter"`
	DnFormat                types.String `tfsdk:"dn_format"`
	PrincipalAttributeId    types.String `tfsdk:"principal_attribute_id"`
	PolicyId                types.String `tfsdk:"policy_id"`
	Order                   types.Int64  `tfsdk:"order"`
	ConnectTimeout          types.Int64  `tfsdk:"connect_timeout"`
}

type AAALDAPServersAttributeServers struct {
	LdapUrl                 types.String `tfsdk:"ldap_url"`
	BaseDn                  types.String `tfsdk:"base_dn"`
	BindDn                  types.String `tfsdk:"bind_dn"`
	BindCredential          types.String `tfsdk:"bind_credential"`
	BindCredentialWo        types.String `tfsdk:"bind_credential_wo"`
	BindCredentialWoVersion types.Int64  `tfsdk:"bind_credential_wo_version"`
	UserFilter              types.String `tfsdk:"user_filter"`
	ConnectTimeout          types.Int64  `tfsdk:"connect_timeout"`
	FailFast                types.Bool   `tfsdk:"fail_fast"`
	Attributes              types.Map    `tfsdk:"attributes"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data AAALDAPServers) getPath() string {
	return "/crosswork/authconfig/v1/remote/ldap"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

func (data AAALDAPServers) getPathDelete() string {
	return "/crosswork/authconfig/v1/remote/ldap"
}

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data AAALDAPServers) getPathGet() string {
	return "/crosswork/authconfig/v1/remote/ldap"
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
func (data AAALDAPServers) toBody(ctx context.Context, state AAALDAPServers) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if true {
		body, _ = sjson.Set(body, "ldap_server", []interface{}{})
		for _, item := range data.Servers {
			itemBody := ""
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			if !item.LdapUrl.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "ldapUrl", item.LdapUrl.ValueString())
			}
			if !item.BindDn.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "bindDn", item.BindDn.ValueString())
			}
			if !item.BindCredential.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "bindCredential", item.BindCredential.ValueString())
			}
			if !item.BindCredentialWo.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "bindCredential", item.BindCredentialWo.ValueString())
			}
			if !item.BaseDn.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "baseDn", item.BaseDn.ValueString())
			}
			if !item.UserFilter.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "userFilter", item.UserFilter.ValueString())
			}
			if !item.DnFormat.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "dnFormat", item.DnFormat.ValueString())
			}
			if !item.PrincipalAttributeId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "principalAttributeId", item.PrincipalAttributeId.ValueString())
			}
			if !item.PolicyId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "policyId", item.PolicyId.ValueString())
			}
			if !item.Order.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "order", item.Order.ValueInt64())
			}
			if !item.ConnectTimeout.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "connectTimeout", item.ConnectTimeout.ValueInt64())
			}
			body, _ = sjson.SetRaw(body, "ldap_server.-1", itemBody)
		}
	}
	if true {
		body, _ = sjson.Set(body, "ldap_attr_servers.ldap_attr_server", []interface{}{})
		for _, item := range data.AttributeServers {
			itemBody := ""
			if !item.LdapUrl.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "ldapUrl", item.LdapUrl.ValueString())
			}
			if !item.BaseDn.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "baseDn", item.BaseDn.ValueString())
			}
			if !item.BindDn.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "bindDn", item.BindDn.ValueString())
			}
			if !item.BindCredential.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "bindCredential", item.BindCredential.ValueString())
			}
			if !item.BindCredentialWo.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "bindCredential", item.BindCredentialWo.ValueString())
			}
			if !item.UserFilter.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "userFilter", item.UserFilter.ValueString())
			}
			if !item.ConnectTimeout.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "connectTimeout", item.ConnectTimeout.ValueInt64())
			}
			if !item.FailFast.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "failFast", item.FailFast.ValueBool())
			}
			if !item.Attributes.IsNull() {
				var values map[string]string
				item.Attributes.ElementsAs(ctx, &values, false)
				itemBody, _ = sjson.Set(itemBody, "attributes", values)
			}
			body, _ = sjson.SetRaw(body, "ldap_attr_servers.ldap_attr_server.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *AAALDAPServers) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("aaa_ldap_servers")
	if value := res.Get("ldap_server"); value.Exists() && len(value.Array()) > 0 {
		data.Servers = make([]AAALDAPServersServers, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := AAALDAPServersServers{}
			if cValue := v.Get("name"); cValue.Exists() {
				item.Name = types.StringValue(cValue.String())
			} else {
				item.Name = types.StringNull()
			}
			if cValue := v.Get("ldapUrl"); cValue.Exists() {
				item.LdapUrl = types.StringValue(cValue.String())
			} else {
				item.LdapUrl = types.StringNull()
			}
			if cValue := v.Get("bindDn"); cValue.Exists() {
				item.BindDn = types.StringValue(cValue.String())
			} else {
				item.BindDn = types.StringNull()
			}
			item.BindCredentialWo = types.StringValue("")
			if cValue := v.Get("baseDn"); cValue.Exists() {
				item.BaseDn = types.StringValue(cValue.String())
			} else {
				item.BaseDn = types.StringNull()
			}
			if cValue := v.Get("userFilter"); cValue.Exists() {
				item.UserFilter = types.StringValue(cValue.String())
			} else {
				item.UserFilter = types.StringNull()
			}
			if cValue := v.Get("dnFormat"); cValue.Exists() {
				item.DnFormat = types.StringValue(cValue.String())
			} else {
				item.DnFormat = types.StringNull()
			}
			if cValue := v.Get("principalAttributeId"); cValue.Exists() {
				item.PrincipalAttributeId = types.StringValue(cValue.String())
			} else {
				item.PrincipalAttributeId = types.StringNull()
			}
			if cValue := v.Get("policyId"); cValue.Exists() {
				item.PolicyId = types.StringValue(cValue.String())
			} else {
				item.PolicyId = types.StringNull()
			}
			if cValue := v.Get("order"); cValue.Exists() {
				item.Order = types.Int64Value(cValue.Int())
			} else {
				item.Order = types.Int64Null()
			}
			if cValue := v.Get("connectTimeout"); cValue.Exists() {
				item.ConnectTimeout = types.Int64Value(cValue.Int())
			} else {
				item.ConnectTimeout = types.Int64Null()
			}
			data.Servers = append(data.Servers, item)
			return true
		})
	}
	if value := res.Get("ldap_attr_servers.ldap_attr_server"); value.Exists() && len(value.Array()) > 0 {
		data.AttributeServers = make([]AAALDAPServersAttributeServers, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := AAALDAPServersAttributeServers{}
			if cValue := v.Get("ldapUrl"); cValue.Exists() {
				item.LdapUrl = types.StringValue(cValue.String())
			} else {
				item.LdapUrl = types.StringNull()
			}
			if cValue := v.Get("baseDn"); cValue.Exists() {
				item.BaseDn = types.StringValue(cValue.String())
			} else {
				item.BaseDn = types.StringNull()
			}
			if cValue := v.Get("bindDn"); cValue.Exists() {
				item.BindDn = types.StringValue(cValue.String())
			} else {
				item.BindDn = types.StringNull()
			}
			item.BindCredentialWo = types.StringValue("")
			if cValue := v.Get("userFilter"); cValue.Exists() {
				item.UserFilter = types.StringValue(cValue.String())
			} else {
				item.UserFilter = types.StringNull()
			}
			if cValue := v.Get("connectTimeout"); cValue.Exists() {
				item.ConnectTimeout = types.Int64Value(cValue.Int())
			} else {
				item.ConnectTimeout = types.Int64Null()
			}
			if cValue := v.Get("failFast"); cValue.Exists() {
				item.FailFast = types.BoolValue(cValue.Bool())
			} else {
				item.FailFast = types.BoolValue(false)
			}
			if cValue := v.Get("attributes"); cValue.Exists() && len(cValue.Map()) > 0 {
				item.Attributes = helpers.GetStringMap(cValue.Map())
			} else {
				item.Attributes = types.MapNull(types.StringType)
			}
			data.AttributeServers = append(data.AttributeServers, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *AAALDAPServers) updateFromBody(ctx context.Context, res gjson.Result) {
	for i := range data.Servers {
		keys := [...]string{"name"}
		keyValues := [...]string{data.Servers[i].Name.ValueString()}

		var r gjson.Result
		res.Get("ldap_server").ForEach(
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
		if value := r.Get("name"); value.Exists() && !data.Servers[i].Name.IsNull() {
			data.Servers[i].Name = types.StringValue(value.String())
		} else {
			data.Servers[i].Name = types.StringNull()
		}
		if value := r.Get("ldapUrl"); value.Exists() && !data.Servers[i].LdapUrl.IsNull() {
			data.Servers[i].LdapUrl = types.StringValue(value.String())
		} else {
			data.Servers[i].LdapUrl = types.StringNull()
		}
		if value := r.Get("bindDn"); value.Exists() && !data.Servers[i].BindDn.IsNull() {
			data.Servers[i].BindDn = types.StringValue(value.String())
		} else {
			data.Servers[i].BindDn = types.StringNull()
		}
		if value := r.Get("baseDn"); value.Exists() && !data.Servers[i].BaseDn.IsNull() {
			data.Servers[i].BaseDn = types.StringValue(value.String())
		} else {
			data.Servers[i].BaseDn = types.StringNull()
		}
		if value := r.Get("userFilter"); value.Exists() && !data.Servers[i].UserFilter.IsNull() {
			data.Servers[i].UserFilter = types.StringValue(value.String())
		} else {
			data.Servers[i].UserFilter = types.StringNull()
		}
		if value := r.Get("dnFormat"); value.Exists() && !data.Servers[i].DnFormat.IsNull() {
			data.Servers[i].DnFormat = types.StringValue(value.String())
		} else {
			data.Servers[i].DnFormat = types.StringNull()
		}
		if value := r.Get("principalAttributeId"); value.Exists() && !data.Servers[i].PrincipalAttributeId.IsNull() {
			data.Servers[i].PrincipalAttributeId = types.StringValue(value.String())
		} else {
			data.Servers[i].PrincipalAttributeId = types.StringNull()
		}
		if value := r.Get("policyId"); value.Exists() && !data.Servers[i].PolicyId.IsNull() {
			data.Servers[i].PolicyId = types.StringValue(value.String())
		} else {
			data.Servers[i].PolicyId = types.StringNull()
		}
		if value := r.Get("order"); value.Exists() && !data.Servers[i].Order.IsNull() {
			data.Servers[i].Order = types.Int64Value(value.Int())
		} else {
			data.Servers[i].Order = types.Int64Null()
		}
		if value := r.Get("connectTimeout"); value.Exists() && !data.Servers[i].ConnectTimeout.IsNull() {
			data.Servers[i].ConnectTimeout = types.Int64Value(value.Int())
		} else {
			data.Servers[i].ConnectTimeout = types.Int64Null()
		}
	}
	for i := range data.AttributeServers {
		keys := [...]string{"ldapUrl"}
		keyValues := [...]string{data.AttributeServers[i].LdapUrl.ValueString()}

		var r gjson.Result
		res.Get("ldap_attr_servers.ldap_attr_server").ForEach(
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
		if value := r.Get("ldapUrl"); value.Exists() && !data.AttributeServers[i].LdapUrl.IsNull() {
			data.AttributeServers[i].LdapUrl = types.StringValue(value.String())
		} else {
			data.AttributeServers[i].LdapUrl = types.StringNull()
		}
		if value := r.Get("baseDn"); value.Exists() && !data.AttributeServers[i].BaseDn.IsNull() {
			data.AttributeServers[i].BaseDn = types.StringValue(value.String())
		} else {
			data.AttributeServers[i].BaseDn = types.StringNull()
		}
		if value := r.Get("bindDn"); value.Exists() && !data.AttributeServers[i].BindDn.IsNull() {
			data.AttributeServers[i].BindDn = types.StringValue(value.String())
		} else {
			data.AttributeServers[i].BindDn = types.StringNull()
		}
		if value := r.Get("userFilter"); value.Exists() && !data.AttributeServers[i].UserFilter.IsNull() {
			data.AttributeServers[i].UserFilter = types.StringValue(value.String())
		} else {
			data.AttributeServers[i].UserFilter = types.StringNull()
		}
		if value := r.Get("connectTimeout"); value.Exists() && !data.AttributeServers[i].ConnectTimeout.IsNull() {
			data.AttributeServers[i].ConnectTimeout = types.Int64Value(value.Int())
		} else {
			data.AttributeServers[i].ConnectTimeout = types.Int64Null()
		}
		if value := r.Get("failFast"); value.Exists() && !data.AttributeServers[i].FailFast.IsNull() {
			data.AttributeServers[i].FailFast = types.BoolValue(value.Bool())
		} else if data.AttributeServers[i].FailFast.ValueBool() != false {
			data.AttributeServers[i].FailFast = types.BoolNull()
		}
		if value := r.Get("attributes"); value.Exists() && !data.AttributeServers[i].Attributes.IsNull() {
			data.AttributeServers[i].Attributes = helpers.GetStringMap(value.Map())
		} else {
			data.AttributeServers[i].Attributes = types.MapNull(types.StringType)
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *AAALDAPServers) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.Servers) > 0 {
		return false
	}
	if len(data.AttributeServers) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

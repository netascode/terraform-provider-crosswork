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
type AAATACACSServers struct {
	Instance types.String              `tfsdk:"instance"`
	Id       types.String              `tfsdk:"id"`
	Servers  []AAATACACSServersServers `tfsdk:"servers"`
}

type AAATACACSServersServers struct {
	Host              types.String `tfsdk:"host"`
	Port              types.Int64  `tfsdk:"port"`
	Secret            types.String `tfsdk:"secret"`
	SecretWo          types.String `tfsdk:"secret_wo"`
	SecretWoVersion   types.Int64  `tfsdk:"secret_wo_version"`
	SecretFormat      types.String `tfsdk:"secret_format"`
	Service           types.String `tfsdk:"service"`
	Priority          types.Int64  `tfsdk:"priority"`
	RetransmitTimeout types.Int64  `tfsdk:"retransmit_timeout"`
	Retries           types.Int64  `tfsdk:"retries"`
	AuthType          types.String `tfsdk:"auth_type"`
	Interface         types.String `tfsdk:"interface"`
	PolicyId          types.String `tfsdk:"policy_id"`
	VirtualDomain     types.String `tfsdk:"virtual_domain"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data AAATACACSServers) getPath() string {
	return "/crosswork/authconfig/v1/remote/tacacs"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

func (data AAATACACSServers) getPathDelete() string {
	return "/crosswork/authconfig/v1/remote/tacacs"
}

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data AAATACACSServers) getPathGet() string {
	return "/crosswork/authconfig/v1/remote/tacacs"
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
func (data AAATACACSServers) toBody(ctx context.Context, state AAATACACSServers) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if true {
		body, _ = sjson.Set(body, "tacacs_servers", []interface{}{})
		for _, item := range data.Servers {
			itemBody := ""
			if !item.Host.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "host", item.Host.ValueString())
			}
			if !item.Port.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "port", item.Port.ValueInt64())
			}
			if !item.Secret.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "secret", item.Secret.ValueString())
			}
			if !item.SecretWo.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "secret", item.SecretWo.ValueString())
			}
			if !item.SecretFormat.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "secretFormat", item.SecretFormat.ValueString())
			}
			if !item.Service.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "service", item.Service.ValueString())
			}
			if !item.Priority.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "priority", item.Priority.ValueInt64())
			}
			if !item.RetransmitTimeout.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "timeout", item.RetransmitTimeout.ValueInt64())
			}
			if !item.Retries.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "retries", item.Retries.ValueInt64())
			}
			if !item.AuthType.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "authType", item.AuthType.ValueString())
			}
			if !item.Interface.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "interface", item.Interface.ValueString())
			}
			if !item.PolicyId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "policy-id", item.PolicyId.ValueString())
			}
			if !item.VirtualDomain.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "virtualDomain", item.VirtualDomain.ValueString())
			}
			body, _ = sjson.SetRaw(body, "tacacs_servers.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *AAATACACSServers) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("aaa_tacacs_servers")
	if value := res.Get("tacacs_servers"); value.Exists() && len(value.Array()) > 0 {
		data.Servers = make([]AAATACACSServersServers, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := AAATACACSServersServers{}
			if cValue := v.Get("host"); cValue.Exists() {
				item.Host = types.StringValue(cValue.String())
			} else {
				item.Host = types.StringNull()
			}
			if cValue := v.Get("port"); cValue.Exists() {
				item.Port = types.Int64Value(cValue.Int())
			} else {
				item.Port = types.Int64Value(49)
			}
			item.SecretWo = types.StringValue("")
			if cValue := v.Get("secretFormat"); cValue.Exists() {
				item.SecretFormat = types.StringValue(cValue.String())
			} else {
				item.SecretFormat = types.StringValue("ascii")
			}
			if cValue := v.Get("service"); cValue.Exists() {
				item.Service = types.StringValue(cValue.String())
			} else {
				item.Service = types.StringNull()
			}
			if cValue := v.Get("priority"); cValue.Exists() {
				item.Priority = types.Int64Value(cValue.Int())
			} else {
				item.Priority = types.Int64Null()
			}
			if cValue := v.Get("timeout"); cValue.Exists() {
				item.RetransmitTimeout = types.Int64Value(cValue.Int())
			} else {
				item.RetransmitTimeout = types.Int64Null()
			}
			if cValue := v.Get("retries"); cValue.Exists() {
				item.Retries = types.Int64Value(cValue.Int())
			} else {
				item.Retries = types.Int64Value(0)
			}
			if cValue := v.Get("authType"); cValue.Exists() {
				item.AuthType = types.StringValue(cValue.String())
			} else {
				item.AuthType = types.StringValue("pap")
			}
			if cValue := v.Get("interface"); cValue.Exists() {
				item.Interface = types.StringValue(cValue.String())
			} else {
				item.Interface = types.StringNull()
			}
			if cValue := v.Get("policy-id"); cValue.Exists() {
				item.PolicyId = types.StringValue(cValue.String())
			} else {
				item.PolicyId = types.StringNull()
			}
			if cValue := v.Get("virtualDomain"); cValue.Exists() {
				item.VirtualDomain = types.StringValue(cValue.String())
			} else {
				item.VirtualDomain = types.StringValue("deviceAccessGroups")
			}
			data.Servers = append(data.Servers, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *AAATACACSServers) updateFromBody(ctx context.Context, res gjson.Result) {
	for i := range data.Servers {
		keys := [...]string{"host"}
		keyValues := [...]string{data.Servers[i].Host.ValueString()}

		var r gjson.Result
		res.Get("tacacs_servers").ForEach(
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
		if value := r.Get("host"); value.Exists() && !data.Servers[i].Host.IsNull() {
			data.Servers[i].Host = types.StringValue(value.String())
		} else {
			data.Servers[i].Host = types.StringNull()
		}
		if value := r.Get("port"); value.Exists() && !data.Servers[i].Port.IsNull() {
			data.Servers[i].Port = types.Int64Value(value.Int())
		} else if data.Servers[i].Port.ValueInt64() != 49 {
			data.Servers[i].Port = types.Int64Null()
		}
		if value := r.Get("secretFormat"); value.Exists() && !data.Servers[i].SecretFormat.IsNull() {
			data.Servers[i].SecretFormat = types.StringValue(value.String())
		} else if data.Servers[i].SecretFormat.ValueString() != "ascii" {
			data.Servers[i].SecretFormat = types.StringNull()
		}
		if value := r.Get("service"); value.Exists() && !data.Servers[i].Service.IsNull() {
			data.Servers[i].Service = types.StringValue(value.String())
		} else {
			data.Servers[i].Service = types.StringNull()
		}
		if value := r.Get("priority"); value.Exists() && !data.Servers[i].Priority.IsNull() {
			data.Servers[i].Priority = types.Int64Value(value.Int())
		} else {
			data.Servers[i].Priority = types.Int64Null()
		}
		if value := r.Get("timeout"); value.Exists() && !data.Servers[i].RetransmitTimeout.IsNull() {
			data.Servers[i].RetransmitTimeout = types.Int64Value(value.Int())
		} else {
			data.Servers[i].RetransmitTimeout = types.Int64Null()
		}
		if value := r.Get("retries"); value.Exists() && !data.Servers[i].Retries.IsNull() {
			data.Servers[i].Retries = types.Int64Value(value.Int())
		} else if data.Servers[i].Retries.ValueInt64() != 0 {
			data.Servers[i].Retries = types.Int64Null()
		}
		if value := r.Get("authType"); value.Exists() && !data.Servers[i].AuthType.IsNull() {
			data.Servers[i].AuthType = types.StringValue(value.String())
		} else if data.Servers[i].AuthType.ValueString() != "pap" {
			data.Servers[i].AuthType = types.StringNull()
		}
		if value := r.Get("interface"); value.Exists() && !data.Servers[i].Interface.IsNull() {
			data.Servers[i].Interface = types.StringValue(value.String())
		} else {
			data.Servers[i].Interface = types.StringNull()
		}
		if value := r.Get("policy-id"); value.Exists() && !data.Servers[i].PolicyId.IsNull() {
			data.Servers[i].PolicyId = types.StringValue(value.String())
		} else {
			data.Servers[i].PolicyId = types.StringNull()
		}
		if value := r.Get("virtualDomain"); value.Exists() && !data.Servers[i].VirtualDomain.IsNull() {
			data.Servers[i].VirtualDomain = types.StringValue(value.String())
		} else if data.Servers[i].VirtualDomain.ValueString() != "deviceAccessGroups" {
			data.Servers[i].VirtualDomain = types.StringNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *AAATACACSServers) isNull(ctx context.Context, res gjson.Result) bool {
	if len(data.Servers) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

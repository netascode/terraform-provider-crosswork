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

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type ProviderProfiles struct {
	Instance        types.String                      `tfsdk:"instance"`
	Id              types.String                      `tfsdk:"id"`
	Name            types.String                      `tfsdk:"name"`
	Family          types.String                      `tfsdk:"family"`
	Profile         types.String                      `tfsdk:"profile"`
	Connectivity    []ProviderProfilesConnectivity    `tfsdk:"connectivity"`
	CertificateName types.String                      `tfsdk:"certificate_name"`
	SupportedModels []ProviderProfilesSupportedModels `tfsdk:"supported_models"`
	Properties      types.Map                         `tfsdk:"properties"`
}

type ProviderProfilesConnectivity struct {
	Protocol       types.String `tfsdk:"protocol"`
	IpAddress      types.String `tfsdk:"ip_address"`
	FqdnHostName   types.String `tfsdk:"fqdn_host_name"`
	FqdnDomainName types.String `tfsdk:"fqdn_domain_name"`
	Port           types.Int64  `tfsdk:"port"`
	Timeout        types.Int64  `tfsdk:"timeout"`
	EncodingType   types.String `tfsdk:"encoding_type"`
}

type ProviderProfilesSupportedModels struct {
	ModelPrefix  types.String `tfsdk:"model_prefix"`
	ModelVersion types.String `tfsdk:"model_version"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data ProviderProfiles) getPath() string {
	return "/crosswork/inventory/v1/providers"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data ProviderProfiles) getPathGet() string {
	return "/crosswork/inventory/v1/providers/query"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// Escapes the id so it cannot change the structure of a request body.
func (data ProviderProfiles) idJSON() string {
	id := data.Id.ValueString()
	if id == "" {
		id = fmt.Sprint(data.Name.ValueString())
	}
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}

func (data ProviderProfiles) readGet(c *cw.Client, path string) (cw.Res, error) {
	return c.Post(path, strings.ReplaceAll(`{"filter":{"name":"{id}"},"filterData":{"PageSize":10,"Criteria":""}}`, "{id}", data.idJSON()), cw.ReadOnly)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data ProviderProfiles) toBody(ctx context.Context, state ProviderProfiles) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "providers.0.name", data.Name.ValueString())
	}
	if !data.Family.IsNull() {
		body, _ = sjson.Set(body, "providers.0.family", data.Family.ValueString())
	}
	if !data.Profile.IsNull() {
		body, _ = sjson.Set(body, "providers.0.profile", data.Profile.ValueString())
	}
	if len(data.Connectivity) > 0 {
		body, _ = sjson.Set(body, "providers.0.connectivity_info", []interface{}{})
		for _, item := range data.Connectivity {
			itemBody := ""
			if !item.Protocol.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "type", item.Protocol.ValueString())
			}
			if !item.IpAddress.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "ipaddrs.0.inet_addr", item.IpAddress.ValueString())
			}
			if !item.FqdnHostName.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "fqdn.host_name", item.FqdnHostName.ValueString())
			}
			if !item.FqdnDomainName.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "fqdn.domain_name", item.FqdnDomainName.ValueString())
			}
			if !item.Port.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "port", item.Port.ValueInt64())
			}
			if !item.Timeout.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "timeout", fmt.Sprint(item.Timeout.ValueInt64()))
			}
			if !item.EncodingType.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "encoding_type", item.EncodingType.ValueString())
			}
			body, _ = sjson.SetRaw(body, "providers.0.connectivity_info.-1", itemBody)
		}
	}
	if !data.CertificateName.IsNull() {
		body, _ = sjson.Set(body, "providers.0.certificate_name", data.CertificateName.ValueString())
	}
	if len(data.SupportedModels) > 0 {
		body, _ = sjson.Set(body, "providers.0.supported_models", []interface{}{})
		for _, item := range data.SupportedModels {
			itemBody := ""
			if !item.ModelPrefix.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "model_prefix", item.ModelPrefix.ValueString())
			}
			if !item.ModelVersion.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "model_version", item.ModelVersion.ValueString())
			}
			body, _ = sjson.SetRaw(body, "providers.0.supported_models.-1", itemBody)
		}
	}
	if !data.Properties.IsNull() {
		var values map[string]string
		data.Properties.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "providers.0.properties", values)
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *ProviderProfiles) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("data.0.name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("data.0.family"); value.Exists() {
		data.Family = types.StringValue(value.String())
	} else {
		data.Family = types.StringNull()
	}
	if value := res.Get("data.0.profile"); value.Exists() {
		data.Profile = types.StringValue(value.String())
	} else {
		data.Profile = types.StringNull()
	}
	if value := res.Get("data.0.connectivity_info"); value.Exists() && len(value.Array()) > 0 {
		data.Connectivity = make([]ProviderProfilesConnectivity, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := ProviderProfilesConnectivity{}
			if cValue := v.Get("type"); cValue.Exists() {
				item.Protocol = types.StringValue(cValue.String())
			} else {
				item.Protocol = types.StringNull()
			}
			if cValue := v.Get("ipaddrs.0.inet_addr"); cValue.Exists() {
				item.IpAddress = types.StringValue(cValue.String())
			} else {
				item.IpAddress = types.StringNull()
			}
			if cValue := v.Get("fqdn.host_name"); cValue.Exists() {
				item.FqdnHostName = types.StringValue(cValue.String())
			} else {
				item.FqdnHostName = types.StringNull()
			}
			if cValue := v.Get("fqdn.domain_name"); cValue.Exists() {
				item.FqdnDomainName = types.StringValue(cValue.String())
			} else {
				item.FqdnDomainName = types.StringNull()
			}
			if cValue := v.Get("port"); cValue.Exists() {
				item.Port = types.Int64Value(cValue.Int())
			} else {
				item.Port = types.Int64Null()
			}
			if cValue := v.Get("timeout"); cValue.Exists() {
				item.Timeout = types.Int64Value(cValue.Int())
			} else {
				item.Timeout = types.Int64Null()
			}
			if cValue := v.Get("encoding_type"); cValue.Exists() {
				item.EncodingType = types.StringValue(cValue.String())
			} else {
				item.EncodingType = types.StringNull()
			}
			data.Connectivity = append(data.Connectivity, item)
			return true
		})
	}
	if value := res.Get("data.0.certificate_name"); value.Exists() {
		data.CertificateName = types.StringValue(value.String())
	} else {
		data.CertificateName = types.StringNull()
	}
	if value := res.Get("data.0.supported_models"); value.Exists() && len(value.Array()) > 0 {
		data.SupportedModels = make([]ProviderProfilesSupportedModels, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := ProviderProfilesSupportedModels{}
			if cValue := v.Get("model_prefix"); cValue.Exists() {
				item.ModelPrefix = types.StringValue(cValue.String())
			} else {
				item.ModelPrefix = types.StringNull()
			}
			if cValue := v.Get("model_version"); cValue.Exists() {
				item.ModelVersion = types.StringValue(cValue.String())
			} else {
				item.ModelVersion = types.StringNull()
			}
			data.SupportedModels = append(data.SupportedModels, item)
			return true
		})
	}
	if value := res.Get("data.0.properties"); value.Exists() && len(value.Map()) > 0 {
		data.Properties = helpers.GetStringMap(value.Map())
	} else {
		data.Properties = types.MapNull(types.StringType)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *ProviderProfiles) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("data.0.name"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("data.0.family"); value.Exists() && !data.Family.IsNull() {
		data.Family = types.StringValue(value.String())
	} else {
		data.Family = types.StringNull()
	}
	if value := res.Get("data.0.profile"); value.Exists() && !data.Profile.IsNull() {
		data.Profile = types.StringValue(value.String())
	} else {
		data.Profile = types.StringNull()
	}
	for i := range data.Connectivity {
		keys := [...]string{"type"}
		keyValues := [...]string{data.Connectivity[i].Protocol.ValueString()}

		var r gjson.Result
		res.Get("data.0.connectivity_info").ForEach(
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
		if value := r.Get("type"); value.Exists() && !data.Connectivity[i].Protocol.IsNull() {
			data.Connectivity[i].Protocol = types.StringValue(value.String())
		} else {
			data.Connectivity[i].Protocol = types.StringNull()
		}
		if value := r.Get("ipaddrs.0.inet_addr"); value.Exists() && !data.Connectivity[i].IpAddress.IsNull() {
			data.Connectivity[i].IpAddress = types.StringValue(value.String())
		} else {
			data.Connectivity[i].IpAddress = types.StringNull()
		}
		if value := r.Get("fqdn.host_name"); value.Exists() && !data.Connectivity[i].FqdnHostName.IsNull() {
			data.Connectivity[i].FqdnHostName = types.StringValue(value.String())
		} else {
			data.Connectivity[i].FqdnHostName = types.StringNull()
		}
		if value := r.Get("fqdn.domain_name"); value.Exists() && !data.Connectivity[i].FqdnDomainName.IsNull() {
			data.Connectivity[i].FqdnDomainName = types.StringValue(value.String())
		} else {
			data.Connectivity[i].FqdnDomainName = types.StringNull()
		}
		if value := r.Get("port"); value.Exists() && !data.Connectivity[i].Port.IsNull() {
			data.Connectivity[i].Port = types.Int64Value(value.Int())
		} else {
			data.Connectivity[i].Port = types.Int64Null()
		}
		if value := r.Get("timeout"); value.Exists() && !data.Connectivity[i].Timeout.IsNull() {
			data.Connectivity[i].Timeout = types.Int64Value(value.Int())
		} else {
			data.Connectivity[i].Timeout = types.Int64Null()
		}
		if value := r.Get("encoding_type"); value.Exists() && !data.Connectivity[i].EncodingType.IsNull() {
			data.Connectivity[i].EncodingType = types.StringValue(value.String())
		} else {
			data.Connectivity[i].EncodingType = types.StringNull()
		}
	}
	if value := res.Get("data.0.certificate_name"); value.Exists() && !data.CertificateName.IsNull() {
		data.CertificateName = types.StringValue(value.String())
	} else {
		data.CertificateName = types.StringNull()
	}
	for i := range data.SupportedModels {
		keys := [...]string{"model_prefix"}
		keyValues := [...]string{data.SupportedModels[i].ModelPrefix.ValueString()}

		var r gjson.Result
		res.Get("data.0.supported_models").ForEach(
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
		if value := r.Get("model_prefix"); value.Exists() && !data.SupportedModels[i].ModelPrefix.IsNull() {
			data.SupportedModels[i].ModelPrefix = types.StringValue(value.String())
		} else {
			data.SupportedModels[i].ModelPrefix = types.StringNull()
		}
		if value := r.Get("model_version"); value.Exists() && !data.SupportedModels[i].ModelVersion.IsNull() {
			data.SupportedModels[i].ModelVersion = types.StringValue(value.String())
		} else {
			data.SupportedModels[i].ModelVersion = types.StringNull()
		}
	}
	if value := res.Get("data.0.properties"); value.Exists() && !data.Properties.IsNull() {
		data.Properties = helpers.GetStringMap(value.Map())
	} else {
		data.Properties = types.MapNull(types.StringType)
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *ProviderProfiles) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Family.IsNull() {
		return false
	}
	if !data.Profile.IsNull() {
		return false
	}
	if len(data.Connectivity) > 0 {
		return false
	}
	if !data.CertificateName.IsNull() {
		return false
	}
	if len(data.SupportedModels) > 0 {
		return false
	}
	if !data.Properties.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

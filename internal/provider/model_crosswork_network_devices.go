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
type NetworkDevices struct {
	Instance                   types.String                 `tfsdk:"instance"`
	Id                         types.String                 `tfsdk:"id"`
	Name                       types.String                 `tfsdk:"name"`
	HostName                   types.String                 `tfsdk:"host_name"`
	CredentialProfile          types.String                 `tfsdk:"credential_profile"`
	AdminState                 types.String                 `tfsdk:"admin_state"`
	ReachabilityCheck          types.String                 `tfsdk:"reachability_check"`
	SerialNumber               types.String                 `tfsdk:"serial_number"`
	MacAddress                 types.String                 `tfsdk:"mac_address"`
	Encryption                 types.String                 `tfsdk:"encryption"`
	DeviceType                 types.String                 `tfsdk:"device_type"`
	SoftwareType               types.String                 `tfsdk:"software_type"`
	SoftwareVersion            types.String                 `tfsdk:"software_version"`
	ProductType                types.String                 `tfsdk:"product_type"`
	SyslogFormat               types.String                 `tfsdk:"syslog_format"`
	Capabilities               types.Set                    `tfsdk:"capabilities"`
	CliCacheEnabled            types.Bool                   `tfsdk:"cli_cache_enabled"`
	ProviderProfile            types.String                 `tfsdk:"provider_profile"`
	Connectivity               []NetworkDevicesConnectivity `tfsdk:"connectivity"`
	IsisSystemId               types.String                 `tfsdk:"isis_system_id"`
	OspfRouterId               types.String                 `tfsdk:"ospf_router_id"`
	TeRouterId                 types.String                 `tfsdk:"te_router_id"`
	Ipv6RouterId               types.String                 `tfsdk:"ipv6_router_id"`
	TelemetryVrf               types.String                 `tfsdk:"telemetry_vrf"`
	TelemetrySourceInterface   types.String                 `tfsdk:"telemetry_source_interface"`
	TelemetryOptOutOfMdtConfig types.Bool                   `tfsdk:"telemetry_opt_out_of_mdt_config"`
	LocationBuilding           types.String                 `tfsdk:"location_building"`
	LocationStreet             types.String                 `tfsdk:"location_street"`
	LocationCity               types.String                 `tfsdk:"location_city"`
	LocationState              types.String                 `tfsdk:"location_state"`
	LocationCountry            types.String                 `tfsdk:"location_country"`
	LocationRegion             types.String                 `tfsdk:"location_region"`
	LocationZip                types.String                 `tfsdk:"location_zip"`
	LocationLatitude           types.Float64                `tfsdk:"location_latitude"`
	LocationLongitude          types.Float64                `tfsdk:"location_longitude"`
	LocationAltitude           types.Float64                `tfsdk:"location_altitude"`
	Tags                       []NetworkDevicesTags         `tfsdk:"tags"`
}

type NetworkDevicesConnectivity struct {
	Protocol                  types.String `tfsdk:"protocol"`
	IpAddress                 types.String `tfsdk:"ip_address"`
	PrefixLength              types.Int64  `tfsdk:"prefix_length"`
	Port                      types.Int64  `tfsdk:"port"`
	Timeout                   types.Int64  `tfsdk:"timeout"`
	EncodingType              types.String `tfsdk:"encoding_type"`
	TrapSourceIp              types.String `tfsdk:"trap_source_ip"`
	DisableTrapCommunityCheck types.Bool   `tfsdk:"disable_trap_community_check"`
}

type NetworkDevicesTags struct {
	Name types.String `tfsdk:"name"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data NetworkDevices) getPath() string {
	return "/crosswork/inventory/v1/nodes"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data NetworkDevices) getPathGet() string {
	return "/crosswork/inventory/v1/nodes/query"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// Escapes the id so it cannot change the structure of a request body.
func (data NetworkDevices) idJSON() string {
	id := data.Id.ValueString()
	if id == "" {
		id = fmt.Sprint(data.Name.ValueString())
	}
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}

func (data NetworkDevices) readGet(c *cw.Client, path string) (cw.Res, error) {
	return c.Post(path, strings.ReplaceAll(`{"filter":{"inventory_id":"{id}"}}`, "{id}", data.idJSON()), cw.ReadOnly)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data NetworkDevices) toBody(ctx context.Context, state NetworkDevices) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "data.0.inventory_id", data.Name.ValueString())
	}
	if data.HostName.ValueString() != "" && !data.HostName.IsNull() {
		body, _ = sjson.Set(body, "data.0.host_name", data.HostName.ValueString())
	}
	if !data.CredentialProfile.IsNull() {
		body, _ = sjson.Set(body, "data.0.profile", data.CredentialProfile.ValueString())
	}
	if !data.AdminState.IsNull() {
		body, _ = sjson.Set(body, "data.0.admin_state", data.AdminState.ValueString())
	}
	if !data.ReachabilityCheck.IsNull() {
		body, _ = sjson.Set(body, "data.0.reachability_check", data.ReachabilityCheck.ValueString())
	}
	if !data.SerialNumber.IsNull() {
		body, _ = sjson.Set(body, "data.0.serial_number", data.SerialNumber.ValueString())
	}
	if !data.MacAddress.IsNull() {
		body, _ = sjson.Set(body, "data.0.mac", data.MacAddress.ValueString())
	}
	if !data.Encryption.IsNull() {
		body, _ = sjson.Set(body, "data.0.encryption", data.Encryption.ValueString())
	}
	if !data.DeviceType.IsNull() {
		body, _ = sjson.Set(body, "data.0.product_info.device_type", data.DeviceType.ValueString())
	}
	if !data.SoftwareType.IsNull() {
		body, _ = sjson.Set(body, "data.0.product_info.software_type", data.SoftwareType.ValueString())
	}
	if !data.SoftwareVersion.IsNull() {
		body, _ = sjson.Set(body, "data.0.product_info.software_version", data.SoftwareVersion.ValueString())
	}
	if !data.ProductType.IsNull() {
		body, _ = sjson.Set(body, "data.0.product_info.product_type", data.ProductType.ValueString())
	}
	if !data.SyslogFormat.IsNull() {
		body, _ = sjson.Set(body, "data.0.product_info.syslog_format", data.SyslogFormat.ValueString())
	}
	if !data.Capabilities.IsNull() {
		var values []string
		data.Capabilities.ElementsAs(ctx, &values, false)
		body, _ = sjson.Set(body, "data.0.product_info.capability", values)
	}
	if !data.CliCacheEnabled.IsNull() {
		body, _ = sjson.Set(body, "data.0.Inventory.cli_cache_enabled", data.CliCacheEnabled.ValueBool())
	}
	if !data.ProviderProfile.IsNull() {
		body, _ = sjson.Set(body, "data.0.providers_family.ROBOT_PROVIDER_NSO.providers.provider_name", data.ProviderProfile.ValueString())
	}
	if len(data.Connectivity) > 0 {
		body, _ = sjson.Set(body, "data.0.connectivity_info", []interface{}{})
		for _, item := range data.Connectivity {
			itemBody := ""
			if !item.Protocol.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "type", item.Protocol.ValueString())
			}
			if !item.IpAddress.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "ipaddrs.0.inet_addr", item.IpAddress.ValueString())
			}
			itemBody, _ = sjson.Set(itemBody, "ipaddrs.0.inet_af", 0)
			if !item.PrefixLength.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "ipaddrs.0.mask", fmt.Sprint(item.PrefixLength.ValueInt64()))
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
			if !item.TrapSourceIp.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "trap_source_ip", item.TrapSourceIp.ValueString())
			}
			if !item.DisableTrapCommunityCheck.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "disable_trap_community_check", item.DisableTrapCommunityCheck.ValueBool())
			}
			body, _ = sjson.SetRaw(body, "data.0.connectivity_info.-1", itemBody)
		}
	}
	if !data.IsisSystemId.IsNull() {
		body, _ = sjson.Set(body, "data.0.routing_info.global_isis_system_id", data.IsisSystemId.ValueString())
	}
	if !data.OspfRouterId.IsNull() {
		body, _ = sjson.Set(body, "data.0.routing_info.global_ospf_router_id", data.OspfRouterId.ValueString())
	}
	if !data.TeRouterId.IsNull() {
		body, _ = sjson.Set(body, "data.0.routing_info.te_router_id", data.TeRouterId.ValueString())
	}
	if !data.Ipv6RouterId.IsNull() {
		body, _ = sjson.Set(body, "data.0.routing_info.ipv6_router_id", data.Ipv6RouterId.ValueString())
	}
	if !data.TelemetryVrf.IsNull() {
		body, _ = sjson.Set(body, "data.0.routing_info.telemetry_intf_src_vrf", data.TelemetryVrf.ValueString())
	}
	if !data.TelemetrySourceInterface.IsNull() {
		body, _ = sjson.Set(body, "data.0.routing_info.telemetry_intf_src_name", data.TelemetrySourceInterface.ValueString())
	}
	if !data.TelemetryOptOutOfMdtConfig.IsNull() {
		body, _ = sjson.Set(body, "data.0.routing_info.is_opt_out_of_mdt_collection_config_enabled", data.TelemetryOptOutOfMdtConfig.ValueBool())
	}
	if !data.LocationBuilding.IsNull() {
		body, _ = sjson.Set(body, "data.0.geo_info.building", data.LocationBuilding.ValueString())
	}
	if !data.LocationStreet.IsNull() {
		body, _ = sjson.Set(body, "data.0.geo_info.street", data.LocationStreet.ValueString())
	}
	if !data.LocationCity.IsNull() {
		body, _ = sjson.Set(body, "data.0.geo_info.city", data.LocationCity.ValueString())
	}
	if !data.LocationState.IsNull() {
		body, _ = sjson.Set(body, "data.0.geo_info.state", data.LocationState.ValueString())
	}
	if !data.LocationCountry.IsNull() {
		body, _ = sjson.Set(body, "data.0.geo_info.country", data.LocationCountry.ValueString())
	}
	if !data.LocationRegion.IsNull() {
		body, _ = sjson.Set(body, "data.0.geo_info.region", data.LocationRegion.ValueString())
	}
	if !data.LocationZip.IsNull() {
		body, _ = sjson.Set(body, "data.0.geo_info.zip", data.LocationZip.ValueString())
	}
	if !data.LocationLatitude.IsNull() {
		body, _ = sjson.Set(body, "data.0.geo_info.coordinates.latitude.value", data.LocationLatitude.ValueFloat64())
	}
	if !data.LocationLongitude.IsNull() {
		body, _ = sjson.Set(body, "data.0.geo_info.coordinates.longitude.value", data.LocationLongitude.ValueFloat64())
	}
	if !data.LocationAltitude.IsNull() {
		body, _ = sjson.Set(body, "data.0.geo_info.coordinates.altitude.value", data.LocationAltitude.ValueFloat64())
	}
	if len(data.Tags) > 0 {
		body, _ = sjson.Set(body, "data.0.tags", []interface{}{})
		for _, item := range data.Tags {
			itemBody := ""
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			body, _ = sjson.SetRaw(body, "data.0.tags.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *NetworkDevices) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("data.0.inventory_id"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("data.0.host_name"); value.Exists() {
		data.HostName = types.StringValue(value.String())
	} else {
		data.HostName = types.StringNull()
	}
	if value := res.Get("data.0.profile"); value.Exists() {
		data.CredentialProfile = types.StringValue(value.String())
	} else {
		data.CredentialProfile = types.StringNull()
	}
	if value := res.Get("data.0.admin_state"); value.Exists() {
		data.AdminState = types.StringValue(value.String())
	} else {
		data.AdminState = types.StringValue("ROBOT_ADMIN_STATE_DOWN")
	}
	if value := res.Get("data.0.reachability_check"); value.Exists() {
		data.ReachabilityCheck = types.StringValue(value.String())
	} else {
		data.ReachabilityCheck = types.StringValue("REACH_CHECK_ENABLE")
	}
	if value := res.Get("data.0.serial_number"); value.Exists() {
		data.SerialNumber = types.StringValue(value.String())
	} else {
		data.SerialNumber = types.StringNull()
	}
	if value := res.Get("data.0.mac"); value.Exists() {
		data.MacAddress = types.StringValue(value.String())
	} else {
		data.MacAddress = types.StringNull()
	}
	if value := res.Get("data.0.encryption"); value.Exists() {
		data.Encryption = types.StringValue(value.String())
	} else {
		data.Encryption = types.StringNull()
	}
	if value := res.Get("data.0.product_info.device_type"); value.Exists() {
		data.DeviceType = types.StringValue(value.String())
	} else {
		data.DeviceType = types.StringValue("NODE_TYPE_ROUTER")
	}
	if value := res.Get("data.0.product_info.software_type"); value.Exists() {
		data.SoftwareType = types.StringValue(value.String())
	} else {
		data.SoftwareType = types.StringNull()
	}
	if value := res.Get("data.0.product_info.software_version"); value.Exists() {
		data.SoftwareVersion = types.StringValue(value.String())
	} else {
		data.SoftwareVersion = types.StringNull()
	}
	if value := res.Get("data.0.product_info.product_type"); value.Exists() {
		data.ProductType = types.StringValue(value.String())
	} else {
		data.ProductType = types.StringNull()
	}
	if value := res.Get("data.0.product_info.syslog_format"); value.Exists() {
		data.SyslogFormat = types.StringValue(value.String())
	} else {
		data.SyslogFormat = types.StringValue("RFC_5424")
	}
	if value := res.Get("data.0.product_info.capability"); value.Exists() && len(value.Array()) > 0 {
		data.Capabilities = helpers.GetStringSet(value.Array())
	} else {
		data.Capabilities = types.SetNull(types.StringType)
	}
	if value := res.Get("data.0.Inventory.cli_cache_enabled"); value.Exists() {
		data.CliCacheEnabled = types.BoolValue(value.Bool())
	} else {
		data.CliCacheEnabled = types.BoolValue(false)
	}
	if value := res.Get("data.0.providers_family.ROBOT_PROVIDER_NSO.providers.provider_name"); value.Exists() {
		data.ProviderProfile = types.StringValue(value.String())
	} else {
		data.ProviderProfile = types.StringNull()
	}
	if value := res.Get("data.0.connectivity_info"); value.Exists() && len(value.Array()) > 0 {
		data.Connectivity = make([]NetworkDevicesConnectivity, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := NetworkDevicesConnectivity{}
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
			if cValue := v.Get("ipaddrs.0.mask"); cValue.Exists() {
				item.PrefixLength = types.Int64Value(cValue.Int())
			} else {
				item.PrefixLength = types.Int64Null()
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
			if cValue := v.Get("trap_source_ip"); cValue.Exists() {
				item.TrapSourceIp = types.StringValue(cValue.String())
			} else {
				item.TrapSourceIp = types.StringNull()
			}
			if cValue := v.Get("disable_trap_community_check"); cValue.Exists() {
				item.DisableTrapCommunityCheck = types.BoolValue(cValue.Bool())
			} else {
				item.DisableTrapCommunityCheck = types.BoolValue(false)
			}
			data.Connectivity = append(data.Connectivity, item)
			return true
		})
	}
	if value := res.Get("data.0.routing_info.global_isis_system_id"); value.Exists() {
		data.IsisSystemId = types.StringValue(value.String())
	} else {
		data.IsisSystemId = types.StringNull()
	}
	if value := res.Get("data.0.routing_info.global_ospf_router_id"); value.Exists() {
		data.OspfRouterId = types.StringValue(value.String())
	} else {
		data.OspfRouterId = types.StringNull()
	}
	if value := res.Get("data.0.routing_info.te_router_id"); value.Exists() {
		data.TeRouterId = types.StringValue(value.String())
	} else {
		data.TeRouterId = types.StringNull()
	}
	if value := res.Get("data.0.routing_info.ipv6_router_id"); value.Exists() {
		data.Ipv6RouterId = types.StringValue(value.String())
	} else {
		data.Ipv6RouterId = types.StringNull()
	}
	if value := res.Get("data.0.routing_info.telemetry_intf_src_vrf"); value.Exists() {
		data.TelemetryVrf = types.StringValue(value.String())
	} else {
		data.TelemetryVrf = types.StringNull()
	}
	if value := res.Get("data.0.routing_info.telemetry_intf_src_name"); value.Exists() {
		data.TelemetrySourceInterface = types.StringValue(value.String())
	} else {
		data.TelemetrySourceInterface = types.StringNull()
	}
	if value := res.Get("data.0.routing_info.is_opt_out_of_mdt_collection_config_enabled"); value.Exists() {
		data.TelemetryOptOutOfMdtConfig = types.BoolValue(value.Bool())
	} else {
		data.TelemetryOptOutOfMdtConfig = types.BoolValue(false)
	}
	if value := res.Get("data.0.geo_info.building"); value.Exists() {
		data.LocationBuilding = types.StringValue(value.String())
	} else {
		data.LocationBuilding = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.street"); value.Exists() {
		data.LocationStreet = types.StringValue(value.String())
	} else {
		data.LocationStreet = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.city"); value.Exists() {
		data.LocationCity = types.StringValue(value.String())
	} else {
		data.LocationCity = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.state"); value.Exists() {
		data.LocationState = types.StringValue(value.String())
	} else {
		data.LocationState = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.country"); value.Exists() {
		data.LocationCountry = types.StringValue(value.String())
	} else {
		data.LocationCountry = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.region"); value.Exists() {
		data.LocationRegion = types.StringValue(value.String())
	} else {
		data.LocationRegion = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.zip"); value.Exists() {
		data.LocationZip = types.StringValue(value.String())
	} else {
		data.LocationZip = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.coordinates.latitude.value"); value.Exists() {
		data.LocationLatitude = types.Float64Value(value.Float())
	} else {
		data.LocationLatitude = types.Float64Null()
	}
	if value := res.Get("data.0.geo_info.coordinates.longitude.value"); value.Exists() {
		data.LocationLongitude = types.Float64Value(value.Float())
	} else {
		data.LocationLongitude = types.Float64Null()
	}
	if value := res.Get("data.0.geo_info.coordinates.altitude.value"); value.Exists() {
		data.LocationAltitude = types.Float64Value(value.Float())
	} else {
		data.LocationAltitude = types.Float64Null()
	}
	if value := res.Get("data.0.tags"); value.Exists() && len(value.Array()) > 0 {
		data.Tags = make([]NetworkDevicesTags, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := NetworkDevicesTags{}
			if cValue := v.Get("name"); cValue.Exists() {
				item.Name = types.StringValue(cValue.String())
			} else {
				item.Name = types.StringNull()
			}
			data.Tags = append(data.Tags, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *NetworkDevices) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("data.0.inventory_id"); value.Exists() && !data.Name.IsNull() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("data.0.host_name"); value.Exists() && !data.HostName.IsNull() {
		data.HostName = types.StringValue(value.String())
	} else {
		data.HostName = types.StringNull()
	}
	if value := res.Get("data.0.profile"); value.Exists() && !data.CredentialProfile.IsNull() {
		data.CredentialProfile = types.StringValue(value.String())
	} else {
		data.CredentialProfile = types.StringNull()
	}
	if value := res.Get("data.0.admin_state"); value.Exists() && !data.AdminState.IsNull() {
		data.AdminState = types.StringValue(value.String())
	} else if data.AdminState.ValueString() != "ROBOT_ADMIN_STATE_DOWN" {
		data.AdminState = types.StringNull()
	}
	if value := res.Get("data.0.reachability_check"); value.Exists() && !data.ReachabilityCheck.IsNull() {
		data.ReachabilityCheck = types.StringValue(value.String())
	} else if data.ReachabilityCheck.ValueString() != "REACH_CHECK_ENABLE" {
		data.ReachabilityCheck = types.StringNull()
	}
	if value := res.Get("data.0.serial_number"); value.Exists() && !data.SerialNumber.IsNull() {
		data.SerialNumber = types.StringValue(value.String())
	} else {
		data.SerialNumber = types.StringNull()
	}
	if value := res.Get("data.0.mac"); value.Exists() && !data.MacAddress.IsNull() {
		data.MacAddress = types.StringValue(value.String())
	} else {
		data.MacAddress = types.StringNull()
	}
	if value := res.Get("data.0.encryption"); value.Exists() && !data.Encryption.IsNull() {
		data.Encryption = types.StringValue(value.String())
	} else {
		data.Encryption = types.StringNull()
	}
	if value := res.Get("data.0.product_info.device_type"); value.Exists() && !data.DeviceType.IsNull() {
		data.DeviceType = types.StringValue(value.String())
	} else if data.DeviceType.ValueString() != "NODE_TYPE_ROUTER" {
		data.DeviceType = types.StringNull()
	}
	if value := res.Get("data.0.product_info.software_type"); value.Exists() && !data.SoftwareType.IsNull() {
		data.SoftwareType = types.StringValue(value.String())
	} else {
		data.SoftwareType = types.StringNull()
	}
	if value := res.Get("data.0.product_info.software_version"); value.Exists() && !data.SoftwareVersion.IsNull() {
		data.SoftwareVersion = types.StringValue(value.String())
	} else {
		data.SoftwareVersion = types.StringNull()
	}
	if value := res.Get("data.0.product_info.product_type"); value.Exists() && !data.ProductType.IsNull() {
		data.ProductType = types.StringValue(value.String())
	} else {
		data.ProductType = types.StringNull()
	}
	if value := res.Get("data.0.product_info.syslog_format"); value.Exists() && !data.SyslogFormat.IsNull() {
		data.SyslogFormat = types.StringValue(value.String())
	} else if data.SyslogFormat.ValueString() != "RFC_5424" {
		data.SyslogFormat = types.StringNull()
	}
	if value := res.Get("data.0.product_info.capability"); value.Exists() && !data.Capabilities.IsNull() {
		data.Capabilities = helpers.GetStringSet(value.Array())
	} else {
		data.Capabilities = types.SetNull(types.StringType)
	}
	if value := res.Get("data.0.Inventory.cli_cache_enabled"); value.Exists() && !data.CliCacheEnabled.IsNull() {
		data.CliCacheEnabled = types.BoolValue(value.Bool())
	} else if data.CliCacheEnabled.ValueBool() != false {
		data.CliCacheEnabled = types.BoolNull()
	}
	if value := res.Get("data.0.providers_family.ROBOT_PROVIDER_NSO.providers.provider_name"); value.Exists() && !data.ProviderProfile.IsNull() {
		data.ProviderProfile = types.StringValue(value.String())
	} else {
		data.ProviderProfile = types.StringNull()
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
		if value := r.Get("ipaddrs.0.mask"); value.Exists() && !data.Connectivity[i].PrefixLength.IsNull() {
			data.Connectivity[i].PrefixLength = types.Int64Value(value.Int())
		} else {
			data.Connectivity[i].PrefixLength = types.Int64Null()
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
		if value := r.Get("trap_source_ip"); value.Exists() && !data.Connectivity[i].TrapSourceIp.IsNull() {
			data.Connectivity[i].TrapSourceIp = types.StringValue(value.String())
		} else {
			data.Connectivity[i].TrapSourceIp = types.StringNull()
		}
		if value := r.Get("disable_trap_community_check"); value.Exists() && !data.Connectivity[i].DisableTrapCommunityCheck.IsNull() {
			data.Connectivity[i].DisableTrapCommunityCheck = types.BoolValue(value.Bool())
		} else if data.Connectivity[i].DisableTrapCommunityCheck.ValueBool() != false {
			data.Connectivity[i].DisableTrapCommunityCheck = types.BoolNull()
		}
	}
	if value := res.Get("data.0.routing_info.global_isis_system_id"); value.Exists() && !data.IsisSystemId.IsNull() {
		data.IsisSystemId = types.StringValue(value.String())
	} else {
		data.IsisSystemId = types.StringNull()
	}
	if value := res.Get("data.0.routing_info.global_ospf_router_id"); value.Exists() && !data.OspfRouterId.IsNull() {
		data.OspfRouterId = types.StringValue(value.String())
	} else {
		data.OspfRouterId = types.StringNull()
	}
	if value := res.Get("data.0.routing_info.te_router_id"); value.Exists() && !data.TeRouterId.IsNull() {
		data.TeRouterId = types.StringValue(value.String())
	} else {
		data.TeRouterId = types.StringNull()
	}
	if value := res.Get("data.0.routing_info.ipv6_router_id"); value.Exists() && !data.Ipv6RouterId.IsNull() {
		data.Ipv6RouterId = types.StringValue(value.String())
	} else {
		data.Ipv6RouterId = types.StringNull()
	}
	if value := res.Get("data.0.routing_info.telemetry_intf_src_vrf"); value.Exists() && !data.TelemetryVrf.IsNull() {
		data.TelemetryVrf = types.StringValue(value.String())
	} else {
		data.TelemetryVrf = types.StringNull()
	}
	if value := res.Get("data.0.routing_info.telemetry_intf_src_name"); value.Exists() && !data.TelemetrySourceInterface.IsNull() {
		data.TelemetrySourceInterface = types.StringValue(value.String())
	} else {
		data.TelemetrySourceInterface = types.StringNull()
	}
	if value := res.Get("data.0.routing_info.is_opt_out_of_mdt_collection_config_enabled"); value.Exists() && !data.TelemetryOptOutOfMdtConfig.IsNull() {
		data.TelemetryOptOutOfMdtConfig = types.BoolValue(value.Bool())
	} else if data.TelemetryOptOutOfMdtConfig.ValueBool() != false {
		data.TelemetryOptOutOfMdtConfig = types.BoolNull()
	}
	if value := res.Get("data.0.geo_info.building"); value.Exists() && !data.LocationBuilding.IsNull() {
		data.LocationBuilding = types.StringValue(value.String())
	} else {
		data.LocationBuilding = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.street"); value.Exists() && !data.LocationStreet.IsNull() {
		data.LocationStreet = types.StringValue(value.String())
	} else {
		data.LocationStreet = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.city"); value.Exists() && !data.LocationCity.IsNull() {
		data.LocationCity = types.StringValue(value.String())
	} else {
		data.LocationCity = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.state"); value.Exists() && !data.LocationState.IsNull() {
		data.LocationState = types.StringValue(value.String())
	} else {
		data.LocationState = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.country"); value.Exists() && !data.LocationCountry.IsNull() {
		data.LocationCountry = types.StringValue(value.String())
	} else {
		data.LocationCountry = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.region"); value.Exists() && !data.LocationRegion.IsNull() {
		data.LocationRegion = types.StringValue(value.String())
	} else {
		data.LocationRegion = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.zip"); value.Exists() && !data.LocationZip.IsNull() {
		data.LocationZip = types.StringValue(value.String())
	} else {
		data.LocationZip = types.StringNull()
	}
	if value := res.Get("data.0.geo_info.coordinates.latitude.value"); value.Exists() && !data.LocationLatitude.IsNull() {
		data.LocationLatitude = types.Float64Value(value.Float())
	} else {
		data.LocationLatitude = types.Float64Null()
	}
	if value := res.Get("data.0.geo_info.coordinates.longitude.value"); value.Exists() && !data.LocationLongitude.IsNull() {
		data.LocationLongitude = types.Float64Value(value.Float())
	} else {
		data.LocationLongitude = types.Float64Null()
	}
	if value := res.Get("data.0.geo_info.coordinates.altitude.value"); value.Exists() && !data.LocationAltitude.IsNull() {
		data.LocationAltitude = types.Float64Value(value.Float())
	} else {
		data.LocationAltitude = types.Float64Null()
	}
	for i := range data.Tags {
		keys := [...]string{"name"}
		keyValues := [...]string{data.Tags[i].Name.ValueString()}

		var r gjson.Result
		res.Get("data.0.tags").ForEach(
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
		if value := r.Get("name"); value.Exists() && !data.Tags[i].Name.IsNull() {
			data.Tags[i].Name = types.StringValue(value.String())
		} else {
			data.Tags[i].Name = types.StringNull()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *NetworkDevices) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.HostName.IsNull() {
		return false
	}
	if !data.CredentialProfile.IsNull() {
		return false
	}
	if !data.AdminState.IsNull() {
		return false
	}
	if !data.ReachabilityCheck.IsNull() {
		return false
	}
	if !data.SerialNumber.IsNull() {
		return false
	}
	if !data.MacAddress.IsNull() {
		return false
	}
	if !data.Encryption.IsNull() {
		return false
	}
	if !data.DeviceType.IsNull() {
		return false
	}
	if !data.SoftwareType.IsNull() {
		return false
	}
	if !data.SoftwareVersion.IsNull() {
		return false
	}
	if !data.ProductType.IsNull() {
		return false
	}
	if !data.SyslogFormat.IsNull() {
		return false
	}
	if !data.Capabilities.IsNull() {
		return false
	}
	if !data.CliCacheEnabled.IsNull() {
		return false
	}
	if !data.ProviderProfile.IsNull() {
		return false
	}
	if len(data.Connectivity) > 0 {
		return false
	}
	if !data.IsisSystemId.IsNull() {
		return false
	}
	if !data.OspfRouterId.IsNull() {
		return false
	}
	if !data.TeRouterId.IsNull() {
		return false
	}
	if !data.Ipv6RouterId.IsNull() {
		return false
	}
	if !data.TelemetryVrf.IsNull() {
		return false
	}
	if !data.TelemetrySourceInterface.IsNull() {
		return false
	}
	if !data.TelemetryOptOutOfMdtConfig.IsNull() {
		return false
	}
	if !data.LocationBuilding.IsNull() {
		return false
	}
	if !data.LocationStreet.IsNull() {
		return false
	}
	if !data.LocationCity.IsNull() {
		return false
	}
	if !data.LocationState.IsNull() {
		return false
	}
	if !data.LocationCountry.IsNull() {
		return false
	}
	if !data.LocationRegion.IsNull() {
		return false
	}
	if !data.LocationZip.IsNull() {
		return false
	}
	if !data.LocationLatitude.IsNull() {
		return false
	}
	if !data.LocationLongitude.IsNull() {
		return false
	}
	if !data.LocationAltitude.IsNull() {
		return false
	}
	if len(data.Tags) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

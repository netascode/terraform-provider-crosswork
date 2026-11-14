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
	"strconv"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type L3VPNVPNServices struct {
	Instance                        types.String                          `tfsdk:"instance"`
	Id                              types.String                          `tfsdk:"id"`
	VpnId                           types.String                          `tfsdk:"vpn_id"`
	VpnServiceTopology              types.String                          `tfsdk:"vpn_service_topology"`
	AutoHubRt                       types.Bool                            `tfsdk:"auto_hub_rt"`
	AutoSpokeRt                     types.Bool                            `tfsdk:"auto_spoke_rt"`
	ServiceAssuranceProfileName     types.String                          `tfsdk:"service_assurance_profile_name"`
	ServiceAssuranceRuleName        types.String                          `tfsdk:"service_assurance_rule_name"`
	ServiceAssuranceMonitoringState types.String                          `tfsdk:"service_assurance_monitoring_state"`
	ServiceAssurancePreservation    types.String                          `tfsdk:"service_assurance_preservation"`
	VpnInstanceProfiles             []L3VPNVPNServicesVpnInstanceProfiles `tfsdk:"vpn_instance_profiles"`
	VpnNodes                        []L3VPNVPNServicesVpnNodes            `tfsdk:"vpn_nodes"`
}

type L3VPNVPNServicesVpnInstanceProfiles struct {
	ProfileId       types.String                                         `tfsdk:"profile_id"`
	Role            types.String                                         `tfsdk:"role"`
	Rd              types.String                                         `tfsdk:"rd"`
	LocalAs         types.Int64                                          `tfsdk:"local_as"`
	AddressFamilies []L3VPNVPNServicesVpnInstanceProfilesAddressFamilies `tfsdk:"address_families"`
}

type L3VPNVPNServicesVpnNodes struct {
	VpnNodeId                 types.String                                        `tfsdk:"vpn_node_id"`
	LocalAs                   types.Int64                                         `tfsdk:"local_as"`
	ActiveVpnInstanceProfiles []L3VPNVPNServicesVpnNodesActiveVpnInstanceProfiles `tfsdk:"active_vpn_instance_profiles"`
	VpnNetworkAccesses        []L3VPNVPNServicesVpnNodesVpnNetworkAccesses        `tfsdk:"vpn_network_accesses"`
}

type L3VPNVPNServicesVpnInstanceProfilesAddressFamilies struct {
	AddressFamily types.String                                                   `tfsdk:"address_family"`
	VpnTargets    []L3VPNVPNServicesVpnInstanceProfilesAddressFamiliesVpnTargets `tfsdk:"vpn_targets"`
	ImportPolicy  types.String                                                   `tfsdk:"import_policy"`
	ExportPolicy  types.String                                                   `tfsdk:"export_policy"`
}

type L3VPNVPNServicesVpnNodesActiveVpnInstanceProfiles struct {
	ProfileId       types.String                                                       `tfsdk:"profile_id"`
	Rd              types.String                                                       `tfsdk:"rd"`
	AddressFamilies []L3VPNVPNServicesVpnNodesActiveVpnInstanceProfilesAddressFamilies `tfsdk:"address_families"`
}
type L3VPNVPNServicesVpnNodesVpnNetworkAccesses struct {
	Id                types.String                                                 `tfsdk:"id"`
	InterfaceId       types.String                                                 `tfsdk:"interface_id"`
	EncapsulationType types.String                                                 `tfsdk:"encapsulation_type"`
	CvlanId           types.Int64                                                  `tfsdk:"cvlan_id"`
	Ipv4Address       types.String                                                 `tfsdk:"ipv4_address"`
	Ipv4PrefixLength  types.Int64                                                  `tfsdk:"ipv4_prefix_length"`
	Ipv6Address       types.String                                                 `tfsdk:"ipv6_address"`
	Ipv6PrefixLength  types.Int64                                                  `tfsdk:"ipv6_prefix_length"`
	Mtu               types.Int64                                                  `tfsdk:"mtu"`
	QosProfiles       []L3VPNVPNServicesVpnNodesVpnNetworkAccessesQosProfiles      `tfsdk:"qos_profiles"`
	RoutingProtocols  []L3VPNVPNServicesVpnNodesVpnNetworkAccessesRoutingProtocols `tfsdk:"routing_protocols"`
}

type L3VPNVPNServicesVpnInstanceProfilesAddressFamiliesVpnTargets struct {
	Id              types.Int64                                                                `tfsdk:"id"`
	RouteTargetType types.String                                                               `tfsdk:"route_target_type"`
	RouteTargets    []L3VPNVPNServicesVpnInstanceProfilesAddressFamiliesVpnTargetsRouteTargets `tfsdk:"route_targets"`
}

type L3VPNVPNServicesVpnNodesActiveVpnInstanceProfilesAddressFamilies struct {
	AddressFamily types.String                                                                 `tfsdk:"address_family"`
	VpnTargets    []L3VPNVPNServicesVpnNodesActiveVpnInstanceProfilesAddressFamiliesVpnTargets `tfsdk:"vpn_targets"`
	ImportPolicy  types.String                                                                 `tfsdk:"import_policy"`
	ExportPolicy  types.String                                                                 `tfsdk:"export_policy"`
}

type L3VPNVPNServicesVpnNodesVpnNetworkAccessesQosProfiles struct {
	Profile   types.String `tfsdk:"profile"`
	Direction types.String `tfsdk:"direction"`
}
type L3VPNVPNServicesVpnNodesVpnNetworkAccessesRoutingProtocols struct {
	Id                       types.String                                                                         `tfsdk:"id"`
	Type                     types.String                                                                         `tfsdk:"type"`
	BgpPeerAs                types.Int64                                                                          `tfsdk:"bgp_peer_as"`
	BgpAddressFamily         types.String                                                                         `tfsdk:"bgp_address_family"`
	BgpNeighbors             types.Set                                                                            `tfsdk:"bgp_neighbors"`
	BgpMultihop              types.Int64                                                                          `tfsdk:"bgp_multihop"`
	BgpRedistributeConnected []L3VPNVPNServicesVpnNodesVpnNetworkAccessesRoutingProtocolsBgpRedistributeConnected `tfsdk:"bgp_redistribute_connected"`
	BgpUpdateSourceIfType    types.String                                                                         `tfsdk:"bgp_update_source_if_type"`
	BgpUpdateSourceIfId      types.String                                                                         `tfsdk:"bgp_update_source_if_id"`
	BgpUpdateSourceSubIfId   types.Int64                                                                          `tfsdk:"bgp_update_source_sub_if_id"`
}

type L3VPNVPNServicesVpnInstanceProfilesAddressFamiliesVpnTargetsRouteTargets struct {
	RouteTarget types.String `tfsdk:"route_target"`
}

type L3VPNVPNServicesVpnNodesActiveVpnInstanceProfilesAddressFamiliesVpnTargets struct {
	Id              types.Int64                                                                              `tfsdk:"id"`
	RouteTargetType types.String                                                                             `tfsdk:"route_target_type"`
	RouteTargets    []L3VPNVPNServicesVpnNodesActiveVpnInstanceProfilesAddressFamiliesVpnTargetsRouteTargets `tfsdk:"route_targets"`
}

type L3VPNVPNServicesVpnNodesVpnNetworkAccessesRoutingProtocolsBgpRedistributeConnected struct {
	AddressFamily types.String `tfsdk:"address_family"`
	Enable        types.Bool   `tfsdk:"enable"`
	Metric        types.Int64  `tfsdk:"metric"`
}

type L3VPNVPNServicesVpnNodesActiveVpnInstanceProfilesAddressFamiliesVpnTargetsRouteTargets struct {
	RouteTarget types.String `tfsdk:"route_target"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data L3VPNVPNServices) getPath() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/ietf-l3vpn-ntw:l3vpn-ntw/vpn-services/vpn-service=%v", url.PathEscape(data.VpnId.ValueString()))
}

func (data L3VPNVPNServices) nsoDevices() nsoDeviceSet {
	devices := nsoDeviceSet{}
	for _, i1 := range data.VpnNodes {
		nsoAddDevice(devices, i1.VpnNodeId)
	}
	return devices
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

func (data L3VPNVPNServices) getPathPost() string {
	return "/crosswork/proxy/nso/restconf/data/ietf-l3vpn-ntw:l3vpn-ntw/vpn-services"
}

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

func (data L3VPNVPNServices) getPathPlan() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/ietf-l3vpn-ntw:l3vpn-ntw/vpn-services/cisco-l3vpn-ntw:vpn-service-plan=%v", url.PathEscape(data.VpnId.ValueString()))
}

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data L3VPNVPNServices) toBody(ctx context.Context, state L3VPNVPNServices) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.VpnId.IsNull() {
		body, _ = sjson.Set(body, "ietf-l3vpn-ntw:vpn-service.0.vpn-id", data.VpnId.ValueString())
	}
	if !data.VpnServiceTopology.IsNull() {
		body, _ = sjson.Set(body, "ietf-l3vpn-ntw:vpn-service.0.vpn-service-topology", data.VpnServiceTopology.ValueString())
	}
	if data.AutoHubRt.ValueBool() {
		body, _ = sjson.SetRaw(body, "ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:vpn-target.auto-hub-rt", "[null]")
	}
	if data.AutoSpokeRt.ValueBool() {
		body, _ = sjson.SetRaw(body, "ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:vpn-target.auto-spoke-rt", "[null]")
	}
	if !data.ServiceAssuranceProfileName.IsNull() {
		body, _ = sjson.Set(body, "ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:service-assurance.profile-name", data.ServiceAssuranceProfileName.ValueString())
	}
	if !data.ServiceAssuranceRuleName.IsNull() {
		body, _ = sjson.Set(body, "ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:service-assurance.rule-name", data.ServiceAssuranceRuleName.ValueString())
	}
	if !data.ServiceAssuranceMonitoringState.IsNull() {
		body, _ = sjson.Set(body, "ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:service-assurance.monitoring-state", data.ServiceAssuranceMonitoringState.ValueString())
	}
	if !data.ServiceAssurancePreservation.IsNull() {
		body, _ = sjson.Set(body, "ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:service-assurance.preservation", data.ServiceAssurancePreservation.ValueString())
	}
	if len(data.VpnInstanceProfiles) > 0 {
		body, _ = sjson.Set(body, "ietf-l3vpn-ntw:vpn-service.0.vpn-instance-profiles.vpn-instance-profile", []interface{}{})
		for _, item := range data.VpnInstanceProfiles {
			itemBody := ""
			if !item.ProfileId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "profile-id", item.ProfileId.ValueString())
			}
			if !item.Role.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "role", item.Role.ValueString())
			}
			if !item.Rd.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "rd", item.Rd.ValueString())
			}
			if !item.LocalAs.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "local-as", item.LocalAs.ValueInt64())
			}
			if len(item.AddressFamilies) > 0 {
				itemBody, _ = sjson.Set(itemBody, "address-family", []interface{}{})
				for _, childItem := range item.AddressFamilies {
					itemChildBody := ""
					if !childItem.AddressFamily.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "address-family", childItem.AddressFamily.ValueString())
					}
					if len(childItem.VpnTargets) > 0 {
						itemChildBody, _ = sjson.Set(itemChildBody, "vpn-targets.vpn-target", []interface{}{})
						for _, childChildItem := range childItem.VpnTargets {
							itemChildChildBody := ""
							if !childChildItem.Id.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "id", childChildItem.Id.ValueInt64())
							}
							if !childChildItem.RouteTargetType.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "route-target-type", childChildItem.RouteTargetType.ValueString())
							}
							if len(childChildItem.RouteTargets) > 0 {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "route-targets", []interface{}{})
								for _, childChildChildItem := range childChildItem.RouteTargets {
									itemChildChildChildBody := ""
									if !childChildChildItem.RouteTarget.IsNull() {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "route-target", childChildChildItem.RouteTarget.ValueString())
									}
									itemChildChildBody, _ = sjson.SetRaw(itemChildChildBody, "route-targets.-1", itemChildChildChildBody)
								}
							}
							itemChildBody, _ = sjson.SetRaw(itemChildBody, "vpn-targets.vpn-target.-1", itemChildChildBody)
						}
					}
					if !childItem.ImportPolicy.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "vpn-targets.vpn-policies.import-policy", childItem.ImportPolicy.ValueString())
					}
					if !childItem.ExportPolicy.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "vpn-targets.vpn-policies.export-policy", childItem.ExportPolicy.ValueString())
					}
					itemBody, _ = sjson.SetRaw(itemBody, "address-family.-1", itemChildBody)
				}
			}
			body, _ = sjson.SetRaw(body, "ietf-l3vpn-ntw:vpn-service.0.vpn-instance-profiles.vpn-instance-profile.-1", itemBody)
		}
	}
	if len(data.VpnNodes) > 0 {
		body, _ = sjson.Set(body, "ietf-l3vpn-ntw:vpn-service.0.vpn-nodes.vpn-node", []interface{}{})
		for _, item := range data.VpnNodes {
			itemBody := ""
			if !item.VpnNodeId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "vpn-node-id", item.VpnNodeId.ValueString())
			}
			if !item.LocalAs.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "local-as", item.LocalAs.ValueInt64())
			}
			if len(item.ActiveVpnInstanceProfiles) > 0 {
				itemBody, _ = sjson.Set(itemBody, "active-vpn-instance-profiles.vpn-instance-profile", []interface{}{})
				for _, childItem := range item.ActiveVpnInstanceProfiles {
					itemChildBody := ""
					if !childItem.ProfileId.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "profile-id", childItem.ProfileId.ValueString())
					}
					if !childItem.Rd.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "rd", childItem.Rd.ValueString())
					}
					if len(childItem.AddressFamilies) > 0 {
						itemChildBody, _ = sjson.Set(itemChildBody, "address-family", []interface{}{})
						for _, childChildItem := range childItem.AddressFamilies {
							itemChildChildBody := ""
							if !childChildItem.AddressFamily.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "address-family", childChildItem.AddressFamily.ValueString())
							}
							if len(childChildItem.VpnTargets) > 0 {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "vpn-targets.vpn-target", []interface{}{})
								for _, childChildChildItem := range childChildItem.VpnTargets {
									itemChildChildChildBody := ""
									if !childChildChildItem.Id.IsNull() {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "id", childChildChildItem.Id.ValueInt64())
									}
									if !childChildChildItem.RouteTargetType.IsNull() {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "route-target-type", childChildChildItem.RouteTargetType.ValueString())
									}
									if len(childChildChildItem.RouteTargets) > 0 {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "route-targets", []interface{}{})
										for _, childChildChildChildItem := range childChildChildItem.RouteTargets {
											itemChildChildChildChildBody := ""
											if !childChildChildChildItem.RouteTarget.IsNull() {
												itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "route-target", childChildChildChildItem.RouteTarget.ValueString())
											}
											itemChildChildChildBody, _ = sjson.SetRaw(itemChildChildChildBody, "route-targets.-1", itemChildChildChildChildBody)
										}
									}
									itemChildChildBody, _ = sjson.SetRaw(itemChildChildBody, "vpn-targets.vpn-target.-1", itemChildChildChildBody)
								}
							}
							if !childChildItem.ImportPolicy.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "vpn-targets.vpn-policies.import-policy", childChildItem.ImportPolicy.ValueString())
							}
							if !childChildItem.ExportPolicy.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "vpn-targets.vpn-policies.export-policy", childChildItem.ExportPolicy.ValueString())
							}
							itemChildBody, _ = sjson.SetRaw(itemChildBody, "address-family.-1", itemChildChildBody)
						}
					}
					itemBody, _ = sjson.SetRaw(itemBody, "active-vpn-instance-profiles.vpn-instance-profile.-1", itemChildBody)
				}
			}
			if len(item.VpnNetworkAccesses) > 0 {
				itemBody, _ = sjson.Set(itemBody, "vpn-network-accesses.vpn-network-access", []interface{}{})
				for _, childItem := range item.VpnNetworkAccesses {
					itemChildBody := ""
					if !childItem.Id.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "id", childItem.Id.ValueString())
					}
					if !childItem.InterfaceId.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "interface-id", childItem.InterfaceId.ValueString())
					}
					if !childItem.EncapsulationType.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "connection.encapsulation.type", childItem.EncapsulationType.ValueString())
					}
					if !childItem.CvlanId.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "connection.encapsulation.dot1q.cvlan-id", childItem.CvlanId.ValueInt64())
					}
					if !childItem.Ipv4Address.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "ip-connection.ipv4.local-address", childItem.Ipv4Address.ValueString())
					}
					if !childItem.Ipv4PrefixLength.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "ip-connection.ipv4.prefix-length", childItem.Ipv4PrefixLength.ValueInt64())
					}
					if !childItem.Ipv6Address.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "ip-connection.ipv6.local-address", childItem.Ipv6Address.ValueString())
					}
					if !childItem.Ipv6PrefixLength.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "ip-connection.ipv6.prefix-length", childItem.Ipv6PrefixLength.ValueInt64())
					}
					if !childItem.Mtu.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "service.mtu", childItem.Mtu.ValueInt64())
					}
					if len(childItem.QosProfiles) > 0 {
						itemChildBody, _ = sjson.Set(itemChildBody, "service.qos.qos-profile.qos-profile", []interface{}{})
						for _, childChildItem := range childItem.QosProfiles {
							itemChildChildBody := ""
							if !childChildItem.Profile.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "profile", childChildItem.Profile.ValueString())
							}
							if !childChildItem.Direction.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "direction", childChildItem.Direction.ValueString())
							}
							itemChildBody, _ = sjson.SetRaw(itemChildBody, "service.qos.qos-profile.qos-profile.-1", itemChildChildBody)
						}
					}
					if len(childItem.RoutingProtocols) > 0 {
						itemChildBody, _ = sjson.Set(itemChildBody, "routing-protocols.routing-protocol", []interface{}{})
						for _, childChildItem := range childItem.RoutingProtocols {
							itemChildChildBody := ""
							if !childChildItem.Id.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "id", childChildItem.Id.ValueString())
							}
							if !childChildItem.Type.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "type", childChildItem.Type.ValueString())
							}
							if !childChildItem.BgpPeerAs.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "bgp.peer-as", childChildItem.BgpPeerAs.ValueInt64())
							}
							if !childChildItem.BgpAddressFamily.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "bgp.address-family", childChildItem.BgpAddressFamily.ValueString())
							}
							if !childChildItem.BgpNeighbors.IsNull() {
								var values []string
								childChildItem.BgpNeighbors.ElementsAs(ctx, &values, false)
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "bgp.neighbor", values)
							}
							if !childChildItem.BgpMultihop.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "bgp.multihop", childChildItem.BgpMultihop.ValueInt64())
							}
							if len(childChildItem.BgpRedistributeConnected) > 0 {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "bgp.redistribute-connected", []interface{}{})
								for _, childChildChildItem := range childChildItem.BgpRedistributeConnected {
									itemChildChildChildBody := ""
									if !childChildChildItem.AddressFamily.IsNull() {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "address-family", childChildChildItem.AddressFamily.ValueString())
									}
									if !childChildChildItem.Enable.IsNull() {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "enable", childChildChildItem.Enable.ValueBool())
									}
									if !childChildChildItem.Metric.IsNull() {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "cisco-l3vpn-ntw:metric", childChildChildItem.Metric.ValueInt64())
									}
									itemChildChildBody, _ = sjson.SetRaw(itemChildChildBody, "bgp.redistribute-connected.-1", itemChildChildChildBody)
								}
							}
							if !childChildItem.BgpUpdateSourceIfType.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "bgp.cisco-l3vpn-ntw:update-source.if-type", childChildItem.BgpUpdateSourceIfType.ValueString())
							}
							if !childChildItem.BgpUpdateSourceIfId.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "bgp.cisco-l3vpn-ntw:update-source.if-id", childChildItem.BgpUpdateSourceIfId.ValueString())
							}
							if !childChildItem.BgpUpdateSourceSubIfId.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "bgp.cisco-l3vpn-ntw:update-source.sub-if-id", childChildItem.BgpUpdateSourceSubIfId.ValueInt64())
							}
							itemChildBody, _ = sjson.SetRaw(itemChildBody, "routing-protocols.routing-protocol.-1", itemChildChildBody)
						}
					}
					itemBody, _ = sjson.SetRaw(itemBody, "vpn-network-accesses.vpn-network-access.-1", itemChildBody)
				}
			}
			body, _ = sjson.SetRaw(body, "ietf-l3vpn-ntw:vpn-service.0.vpn-nodes.vpn-node.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *L3VPNVPNServices) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.vpn-service-topology"); value.Exists() {
		data.VpnServiceTopology = types.StringValue(value.String())
	} else {
		data.VpnServiceTopology = types.StringNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:vpn-target.auto-hub-rt"); value.Exists() {
		data.AutoHubRt = types.BoolValue(true)
	} else {
		data.AutoHubRt = types.BoolNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:vpn-target.auto-spoke-rt"); value.Exists() {
		data.AutoSpokeRt = types.BoolValue(true)
	} else {
		data.AutoSpokeRt = types.BoolNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:service-assurance.profile-name"); value.Exists() {
		data.ServiceAssuranceProfileName = types.StringValue(value.String())
	} else {
		data.ServiceAssuranceProfileName = types.StringNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:service-assurance.rule-name"); value.Exists() {
		data.ServiceAssuranceRuleName = types.StringValue(value.String())
	} else {
		data.ServiceAssuranceRuleName = types.StringNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:service-assurance.monitoring-state"); value.Exists() {
		data.ServiceAssuranceMonitoringState = types.StringValue(value.String())
	} else {
		data.ServiceAssuranceMonitoringState = types.StringNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:service-assurance.preservation"); value.Exists() {
		data.ServiceAssurancePreservation = types.StringValue(value.String())
	} else {
		data.ServiceAssurancePreservation = types.StringNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.vpn-instance-profiles.vpn-instance-profile"); value.Exists() && len(value.Array()) > 0 {
		data.VpnInstanceProfiles = make([]L3VPNVPNServicesVpnInstanceProfiles, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := L3VPNVPNServicesVpnInstanceProfiles{}
			if cValue := v.Get("profile-id"); cValue.Exists() {
				item.ProfileId = types.StringValue(cValue.String())
			} else {
				item.ProfileId = types.StringNull()
			}
			if cValue := v.Get("role"); cValue.Exists() {
				item.Role = types.StringValue(cValue.String())
			} else {
				item.Role = types.StringNull()
			}
			if cValue := v.Get("rd"); cValue.Exists() {
				item.Rd = types.StringValue(cValue.String())
			} else {
				item.Rd = types.StringNull()
			}
			if cValue := v.Get("local-as"); cValue.Exists() {
				item.LocalAs = types.Int64Value(cValue.Int())
			} else {
				item.LocalAs = types.Int64Null()
			}
			if cValue := v.Get("address-family"); cValue.Exists() && len(cValue.Array()) > 0 {
				item.AddressFamilies = make([]L3VPNVPNServicesVpnInstanceProfilesAddressFamilies, 0)
				cValue.ForEach(func(ck, cv gjson.Result) bool {
					cItem := L3VPNVPNServicesVpnInstanceProfilesAddressFamilies{}
					if cwValue := cv.Get("address-family"); cwValue.Exists() {
						cItem.AddressFamily = types.StringValue(cwValue.String())
					} else {
						cItem.AddressFamily = types.StringNull()
					}
					if cwValue := cv.Get("vpn-targets.vpn-target"); cwValue.Exists() && len(cwValue.Array()) > 0 {
						cItem.VpnTargets = make([]L3VPNVPNServicesVpnInstanceProfilesAddressFamiliesVpnTargets, 0)
						cwValue.ForEach(func(cck, ccv gjson.Result) bool {
							cwItem := L3VPNVPNServicesVpnInstanceProfilesAddressFamiliesVpnTargets{}
							if ccwValue := ccv.Get("id"); ccwValue.Exists() {
								cwItem.Id = types.Int64Value(ccwValue.Int())
							} else {
								cwItem.Id = types.Int64Null()
							}
							if ccwValue := ccv.Get("route-target-type"); ccwValue.Exists() {
								cwItem.RouteTargetType = types.StringValue(ccwValue.String())
							} else {
								cwItem.RouteTargetType = types.StringNull()
							}
							if ccwValue := ccv.Get("route-targets"); ccwValue.Exists() && len(ccwValue.Array()) > 0 {
								cwItem.RouteTargets = make([]L3VPNVPNServicesVpnInstanceProfilesAddressFamiliesVpnTargetsRouteTargets, 0)
								ccwValue.ForEach(func(c3k, c3v gjson.Result) bool {
									c3Item := L3VPNVPNServicesVpnInstanceProfilesAddressFamiliesVpnTargetsRouteTargets{}
									if c3Value := c3v.Get("route-target"); c3Value.Exists() {
										c3Item.RouteTarget = types.StringValue(c3Value.String())
									} else {
										c3Item.RouteTarget = types.StringNull()
									}
									cwItem.RouteTargets = append(cwItem.RouteTargets, c3Item)
									return true
								})
							}
							cItem.VpnTargets = append(cItem.VpnTargets, cwItem)
							return true
						})
					}
					if cwValue := cv.Get("vpn-targets.vpn-policies.import-policy"); cwValue.Exists() {
						cItem.ImportPolicy = types.StringValue(cwValue.String())
					} else {
						cItem.ImportPolicy = types.StringNull()
					}
					if cwValue := cv.Get("vpn-targets.vpn-policies.export-policy"); cwValue.Exists() {
						cItem.ExportPolicy = types.StringValue(cwValue.String())
					} else {
						cItem.ExportPolicy = types.StringNull()
					}
					item.AddressFamilies = append(item.AddressFamilies, cItem)
					return true
				})
			}
			data.VpnInstanceProfiles = append(data.VpnInstanceProfiles, item)
			return true
		})
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.vpn-nodes.vpn-node"); value.Exists() && len(value.Array()) > 0 {
		data.VpnNodes = make([]L3VPNVPNServicesVpnNodes, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := L3VPNVPNServicesVpnNodes{}
			if cValue := v.Get("vpn-node-id"); cValue.Exists() {
				item.VpnNodeId = types.StringValue(cValue.String())
			} else {
				item.VpnNodeId = types.StringNull()
			}
			if cValue := v.Get("local-as"); cValue.Exists() {
				item.LocalAs = types.Int64Value(cValue.Int())
			} else {
				item.LocalAs = types.Int64Null()
			}
			if cValue := v.Get("active-vpn-instance-profiles.vpn-instance-profile"); cValue.Exists() && len(cValue.Array()) > 0 {
				item.ActiveVpnInstanceProfiles = make([]L3VPNVPNServicesVpnNodesActiveVpnInstanceProfiles, 0)
				cValue.ForEach(func(ck, cv gjson.Result) bool {
					cItem := L3VPNVPNServicesVpnNodesActiveVpnInstanceProfiles{}
					if cwValue := cv.Get("profile-id"); cwValue.Exists() {
						cItem.ProfileId = types.StringValue(cwValue.String())
					} else {
						cItem.ProfileId = types.StringNull()
					}
					if cwValue := cv.Get("rd"); cwValue.Exists() {
						cItem.Rd = types.StringValue(cwValue.String())
					} else {
						cItem.Rd = types.StringNull()
					}
					if cwValue := cv.Get("address-family"); cwValue.Exists() && len(cwValue.Array()) > 0 {
						cItem.AddressFamilies = make([]L3VPNVPNServicesVpnNodesActiveVpnInstanceProfilesAddressFamilies, 0)
						cwValue.ForEach(func(cck, ccv gjson.Result) bool {
							cwItem := L3VPNVPNServicesVpnNodesActiveVpnInstanceProfilesAddressFamilies{}
							if ccwValue := ccv.Get("address-family"); ccwValue.Exists() {
								cwItem.AddressFamily = types.StringValue(ccwValue.String())
							} else {
								cwItem.AddressFamily = types.StringNull()
							}
							if ccwValue := ccv.Get("vpn-targets.vpn-target"); ccwValue.Exists() && len(ccwValue.Array()) > 0 {
								cwItem.VpnTargets = make([]L3VPNVPNServicesVpnNodesActiveVpnInstanceProfilesAddressFamiliesVpnTargets, 0)
								ccwValue.ForEach(func(c3k, c3v gjson.Result) bool {
									c3Item := L3VPNVPNServicesVpnNodesActiveVpnInstanceProfilesAddressFamiliesVpnTargets{}
									if c3Value := c3v.Get("id"); c3Value.Exists() {
										c3Item.Id = types.Int64Value(c3Value.Int())
									} else {
										c3Item.Id = types.Int64Null()
									}
									if c3Value := c3v.Get("route-target-type"); c3Value.Exists() {
										c3Item.RouteTargetType = types.StringValue(c3Value.String())
									} else {
										c3Item.RouteTargetType = types.StringNull()
									}
									if c3Value := c3v.Get("route-targets"); c3Value.Exists() && len(c3Value.Array()) > 0 {
										c3Item.RouteTargets = make([]L3VPNVPNServicesVpnNodesActiveVpnInstanceProfilesAddressFamiliesVpnTargetsRouteTargets, 0)
										c3Value.ForEach(func(c4k, c4v gjson.Result) bool {
											c4Item := L3VPNVPNServicesVpnNodesActiveVpnInstanceProfilesAddressFamiliesVpnTargetsRouteTargets{}
											if c4Value := c4v.Get("route-target"); c4Value.Exists() {
												c4Item.RouteTarget = types.StringValue(c4Value.String())
											} else {
												c4Item.RouteTarget = types.StringNull()
											}
											c3Item.RouteTargets = append(c3Item.RouteTargets, c4Item)
											return true
										})
									}
									cwItem.VpnTargets = append(cwItem.VpnTargets, c3Item)
									return true
								})
							}
							if ccwValue := ccv.Get("vpn-targets.vpn-policies.import-policy"); ccwValue.Exists() {
								cwItem.ImportPolicy = types.StringValue(ccwValue.String())
							} else {
								cwItem.ImportPolicy = types.StringNull()
							}
							if ccwValue := ccv.Get("vpn-targets.vpn-policies.export-policy"); ccwValue.Exists() {
								cwItem.ExportPolicy = types.StringValue(ccwValue.String())
							} else {
								cwItem.ExportPolicy = types.StringNull()
							}
							cItem.AddressFamilies = append(cItem.AddressFamilies, cwItem)
							return true
						})
					}
					item.ActiveVpnInstanceProfiles = append(item.ActiveVpnInstanceProfiles, cItem)
					return true
				})
			}
			if cValue := v.Get("vpn-network-accesses.vpn-network-access"); cValue.Exists() && len(cValue.Array()) > 0 {
				item.VpnNetworkAccesses = make([]L3VPNVPNServicesVpnNodesVpnNetworkAccesses, 0)
				cValue.ForEach(func(ck, cv gjson.Result) bool {
					cItem := L3VPNVPNServicesVpnNodesVpnNetworkAccesses{}
					if cwValue := cv.Get("id"); cwValue.Exists() {
						cItem.Id = types.StringValue(cwValue.String())
					} else {
						cItem.Id = types.StringNull()
					}
					if cwValue := cv.Get("interface-id"); cwValue.Exists() {
						cItem.InterfaceId = types.StringValue(cwValue.String())
					} else {
						cItem.InterfaceId = types.StringNull()
					}
					if cwValue := cv.Get("connection.encapsulation.type"); cwValue.Exists() {
						cItem.EncapsulationType = types.StringValue(cwValue.String())
					} else {
						cItem.EncapsulationType = types.StringNull()
					}
					if cwValue := cv.Get("connection.encapsulation.dot1q.cvlan-id"); cwValue.Exists() {
						cItem.CvlanId = types.Int64Value(cwValue.Int())
					} else {
						cItem.CvlanId = types.Int64Null()
					}
					if cwValue := cv.Get("ip-connection.ipv4.local-address"); cwValue.Exists() {
						cItem.Ipv4Address = types.StringValue(cwValue.String())
					} else {
						cItem.Ipv4Address = types.StringNull()
					}
					if cwValue := cv.Get("ip-connection.ipv4.prefix-length"); cwValue.Exists() {
						cItem.Ipv4PrefixLength = types.Int64Value(cwValue.Int())
					} else {
						cItem.Ipv4PrefixLength = types.Int64Null()
					}
					if cwValue := cv.Get("ip-connection.ipv6.local-address"); cwValue.Exists() {
						cItem.Ipv6Address = types.StringValue(cwValue.String())
					} else {
						cItem.Ipv6Address = types.StringNull()
					}
					if cwValue := cv.Get("ip-connection.ipv6.prefix-length"); cwValue.Exists() {
						cItem.Ipv6PrefixLength = types.Int64Value(cwValue.Int())
					} else {
						cItem.Ipv6PrefixLength = types.Int64Null()
					}
					if cwValue := cv.Get("service.mtu"); cwValue.Exists() {
						cItem.Mtu = types.Int64Value(cwValue.Int())
					} else {
						cItem.Mtu = types.Int64Null()
					}
					if cwValue := cv.Get("service.qos.qos-profile.qos-profile"); cwValue.Exists() && len(cwValue.Array()) > 0 {
						cItem.QosProfiles = make([]L3VPNVPNServicesVpnNodesVpnNetworkAccessesQosProfiles, 0)
						cwValue.ForEach(func(cck, ccv gjson.Result) bool {
							cwItem := L3VPNVPNServicesVpnNodesVpnNetworkAccessesQosProfiles{}
							if ccwValue := ccv.Get("profile"); ccwValue.Exists() {
								cwItem.Profile = types.StringValue(ccwValue.String())
							} else {
								cwItem.Profile = types.StringNull()
							}
							if ccwValue := ccv.Get("direction"); ccwValue.Exists() {
								cwItem.Direction = types.StringValue(ccwValue.String())
							} else {
								cwItem.Direction = types.StringNull()
							}
							cItem.QosProfiles = append(cItem.QosProfiles, cwItem)
							return true
						})
					}
					if cwValue := cv.Get("routing-protocols.routing-protocol"); cwValue.Exists() && len(cwValue.Array()) > 0 {
						cItem.RoutingProtocols = make([]L3VPNVPNServicesVpnNodesVpnNetworkAccessesRoutingProtocols, 0)
						cwValue.ForEach(func(cck, ccv gjson.Result) bool {
							cwItem := L3VPNVPNServicesVpnNodesVpnNetworkAccessesRoutingProtocols{}
							if ccwValue := ccv.Get("id"); ccwValue.Exists() {
								cwItem.Id = types.StringValue(ccwValue.String())
							} else {
								cwItem.Id = types.StringNull()
							}
							if ccwValue := ccv.Get("type"); ccwValue.Exists() {
								cwItem.Type = types.StringValue(ccwValue.String())
							} else {
								cwItem.Type = types.StringNull()
							}
							if ccwValue := ccv.Get("bgp.peer-as"); ccwValue.Exists() {
								cwItem.BgpPeerAs = types.Int64Value(ccwValue.Int())
							} else {
								cwItem.BgpPeerAs = types.Int64Null()
							}
							if ccwValue := ccv.Get("bgp.address-family"); ccwValue.Exists() {
								cwItem.BgpAddressFamily = types.StringValue(ccwValue.String())
							} else {
								cwItem.BgpAddressFamily = types.StringNull()
							}
							if ccwValue := ccv.Get("bgp.neighbor"); ccwValue.Exists() && len(ccwValue.Array()) > 0 {
								cwItem.BgpNeighbors = helpers.GetStringSet(ccwValue.Array())
							} else {
								cwItem.BgpNeighbors = types.SetNull(types.StringType)
							}
							if ccwValue := ccv.Get("bgp.multihop"); ccwValue.Exists() {
								cwItem.BgpMultihop = types.Int64Value(ccwValue.Int())
							} else {
								cwItem.BgpMultihop = types.Int64Null()
							}
							if ccwValue := ccv.Get("bgp.redistribute-connected"); ccwValue.Exists() && len(ccwValue.Array()) > 0 {
								cwItem.BgpRedistributeConnected = make([]L3VPNVPNServicesVpnNodesVpnNetworkAccessesRoutingProtocolsBgpRedistributeConnected, 0)
								ccwValue.ForEach(func(c3k, c3v gjson.Result) bool {
									c3Item := L3VPNVPNServicesVpnNodesVpnNetworkAccessesRoutingProtocolsBgpRedistributeConnected{}
									if c3Value := c3v.Get("address-family"); c3Value.Exists() {
										c3Item.AddressFamily = types.StringValue(c3Value.String())
									} else {
										c3Item.AddressFamily = types.StringNull()
									}
									if c3Value := c3v.Get("enable"); c3Value.Exists() {
										c3Item.Enable = types.BoolValue(c3Value.Bool())
									} else {
										c3Item.Enable = types.BoolNull()
									}
									if c3Value := c3v.Get("cisco-l3vpn-ntw:metric"); c3Value.Exists() {
										c3Item.Metric = types.Int64Value(c3Value.Int())
									} else {
										c3Item.Metric = types.Int64Null()
									}
									cwItem.BgpRedistributeConnected = append(cwItem.BgpRedistributeConnected, c3Item)
									return true
								})
							}
							if ccwValue := ccv.Get("bgp.cisco-l3vpn-ntw:update-source.if-type"); ccwValue.Exists() {
								cwItem.BgpUpdateSourceIfType = types.StringValue(ccwValue.String())
							} else {
								cwItem.BgpUpdateSourceIfType = types.StringNull()
							}
							if ccwValue := ccv.Get("bgp.cisco-l3vpn-ntw:update-source.if-id"); ccwValue.Exists() {
								cwItem.BgpUpdateSourceIfId = types.StringValue(ccwValue.String())
							} else {
								cwItem.BgpUpdateSourceIfId = types.StringNull()
							}
							if ccwValue := ccv.Get("bgp.cisco-l3vpn-ntw:update-source.sub-if-id"); ccwValue.Exists() {
								cwItem.BgpUpdateSourceSubIfId = types.Int64Value(ccwValue.Int())
							} else {
								cwItem.BgpUpdateSourceSubIfId = types.Int64Null()
							}
							cItem.RoutingProtocols = append(cItem.RoutingProtocols, cwItem)
							return true
						})
					}
					item.VpnNetworkAccesses = append(item.VpnNetworkAccesses, cItem)
					return true
				})
			}
			data.VpnNodes = append(data.VpnNodes, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *L3VPNVPNServices) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.vpn-service-topology"); value.Exists() && !data.VpnServiceTopology.IsNull() {
		data.VpnServiceTopology = types.StringValue(value.String())
	} else {
		data.VpnServiceTopology = types.StringNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:vpn-target.auto-hub-rt"); value.Exists() && !data.AutoHubRt.IsNull() {
		data.AutoHubRt = types.BoolValue(true)
	} else {
		data.AutoHubRt = types.BoolNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:vpn-target.auto-spoke-rt"); value.Exists() && !data.AutoSpokeRt.IsNull() {
		data.AutoSpokeRt = types.BoolValue(true)
	} else {
		data.AutoSpokeRt = types.BoolNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:service-assurance.profile-name"); value.Exists() && !data.ServiceAssuranceProfileName.IsNull() {
		data.ServiceAssuranceProfileName = types.StringValue(value.String())
	} else {
		data.ServiceAssuranceProfileName = types.StringNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:service-assurance.rule-name"); value.Exists() && !data.ServiceAssuranceRuleName.IsNull() {
		data.ServiceAssuranceRuleName = types.StringValue(value.String())
	} else {
		data.ServiceAssuranceRuleName = types.StringNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:service-assurance.monitoring-state"); value.Exists() && !data.ServiceAssuranceMonitoringState.IsNull() {
		data.ServiceAssuranceMonitoringState = types.StringValue(value.String())
	} else {
		data.ServiceAssuranceMonitoringState = types.StringNull()
	}
	if value := res.Get("ietf-l3vpn-ntw:vpn-service.0.cisco-l3vpn-ntw:service-assurance.preservation"); value.Exists() && !data.ServiceAssurancePreservation.IsNull() {
		data.ServiceAssurancePreservation = types.StringValue(value.String())
	} else {
		data.ServiceAssurancePreservation = types.StringNull()
	}
	for i := range data.VpnInstanceProfiles {
		keys := [...]string{"profile-id"}
		keyValues := [...]string{data.VpnInstanceProfiles[i].ProfileId.ValueString()}

		var r gjson.Result
		res.Get("ietf-l3vpn-ntw:vpn-service.0.vpn-instance-profiles.vpn-instance-profile").ForEach(
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
		if value := r.Get("profile-id"); value.Exists() && !data.VpnInstanceProfiles[i].ProfileId.IsNull() {
			data.VpnInstanceProfiles[i].ProfileId = types.StringValue(value.String())
		} else {
			data.VpnInstanceProfiles[i].ProfileId = types.StringNull()
		}
		if value := r.Get("role"); value.Exists() && !data.VpnInstanceProfiles[i].Role.IsNull() {
			data.VpnInstanceProfiles[i].Role = types.StringValue(value.String())
		} else {
			data.VpnInstanceProfiles[i].Role = types.StringNull()
		}
		if value := r.Get("rd"); value.Exists() && !data.VpnInstanceProfiles[i].Rd.IsNull() {
			data.VpnInstanceProfiles[i].Rd = types.StringValue(value.String())
		} else {
			data.VpnInstanceProfiles[i].Rd = types.StringNull()
		}
		if value := r.Get("local-as"); value.Exists() && !data.VpnInstanceProfiles[i].LocalAs.IsNull() {
			data.VpnInstanceProfiles[i].LocalAs = types.Int64Value(value.Int())
		} else {
			data.VpnInstanceProfiles[i].LocalAs = types.Int64Null()
		}
		for ci := range data.VpnInstanceProfiles[i].AddressFamilies {
			keys := [...]string{"address-family"}
			keyValues := [...]string{data.VpnInstanceProfiles[i].AddressFamilies[ci].AddressFamily.ValueString()}

			var cr gjson.Result
			r.Get("address-family").ForEach(
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
						cr = v
						return false
					}
					return true
				},
			)
			if value := cr.Get("address-family"); value.Exists() && !data.VpnInstanceProfiles[i].AddressFamilies[ci].AddressFamily.IsNull() {
				data.VpnInstanceProfiles[i].AddressFamilies[ci].AddressFamily = types.StringValue(value.String())
			} else {
				data.VpnInstanceProfiles[i].AddressFamilies[ci].AddressFamily = types.StringNull()
			}
			for cci := range data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets {
				keys := [...]string{"id"}
				keyValues := [...]string{strconv.FormatInt(data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets[cci].Id.ValueInt64(), 10)}

				var ccr gjson.Result
				cr.Get("vpn-targets.vpn-target").ForEach(
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
							ccr = v
							return false
						}
						return true
					},
				)
				if value := ccr.Get("id"); value.Exists() && !data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets[cci].Id.IsNull() {
					data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets[cci].Id = types.Int64Value(value.Int())
				} else {
					data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets[cci].Id = types.Int64Null()
				}
				if value := ccr.Get("route-target-type"); value.Exists() && !data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets[cci].RouteTargetType.IsNull() {
					data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets[cci].RouteTargetType = types.StringValue(value.String())
				} else {
					data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets[cci].RouteTargetType = types.StringNull()
				}
				for c3i := range data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets[cci].RouteTargets {
					keys := [...]string{"route-target"}
					keyValues := [...]string{data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets[cci].RouteTargets[c3i].RouteTarget.ValueString()}

					var c3r gjson.Result
					ccr.Get("route-targets").ForEach(
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
								c3r = v
								return false
							}
							return true
						},
					)
					if value := c3r.Get("route-target"); value.Exists() && !data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets[cci].RouteTargets[c3i].RouteTarget.IsNull() {
						data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets[cci].RouteTargets[c3i].RouteTarget = types.StringValue(value.String())
					} else {
						data.VpnInstanceProfiles[i].AddressFamilies[ci].VpnTargets[cci].RouteTargets[c3i].RouteTarget = types.StringNull()
					}
				}
			}
			if value := cr.Get("vpn-targets.vpn-policies.import-policy"); value.Exists() && !data.VpnInstanceProfiles[i].AddressFamilies[ci].ImportPolicy.IsNull() {
				data.VpnInstanceProfiles[i].AddressFamilies[ci].ImportPolicy = types.StringValue(value.String())
			} else {
				data.VpnInstanceProfiles[i].AddressFamilies[ci].ImportPolicy = types.StringNull()
			}
			if value := cr.Get("vpn-targets.vpn-policies.export-policy"); value.Exists() && !data.VpnInstanceProfiles[i].AddressFamilies[ci].ExportPolicy.IsNull() {
				data.VpnInstanceProfiles[i].AddressFamilies[ci].ExportPolicy = types.StringValue(value.String())
			} else {
				data.VpnInstanceProfiles[i].AddressFamilies[ci].ExportPolicy = types.StringNull()
			}
		}
	}
	for i := range data.VpnNodes {
		keys := [...]string{"vpn-node-id"}
		keyValues := [...]string{data.VpnNodes[i].VpnNodeId.ValueString()}

		var r gjson.Result
		res.Get("ietf-l3vpn-ntw:vpn-service.0.vpn-nodes.vpn-node").ForEach(
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
		if value := r.Get("vpn-node-id"); value.Exists() && !data.VpnNodes[i].VpnNodeId.IsNull() {
			data.VpnNodes[i].VpnNodeId = types.StringValue(value.String())
		} else {
			data.VpnNodes[i].VpnNodeId = types.StringNull()
		}
		if value := r.Get("local-as"); value.Exists() && !data.VpnNodes[i].LocalAs.IsNull() {
			data.VpnNodes[i].LocalAs = types.Int64Value(value.Int())
		} else {
			data.VpnNodes[i].LocalAs = types.Int64Null()
		}
		for ci := range data.VpnNodes[i].ActiveVpnInstanceProfiles {
			keys := [...]string{"profile-id"}
			keyValues := [...]string{data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].ProfileId.ValueString()}

			var cr gjson.Result
			r.Get("active-vpn-instance-profiles.vpn-instance-profile").ForEach(
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
						cr = v
						return false
					}
					return true
				},
			)
			if value := cr.Get("profile-id"); value.Exists() && !data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].ProfileId.IsNull() {
				data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].ProfileId = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].ProfileId = types.StringNull()
			}
			if value := cr.Get("rd"); value.Exists() && !data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].Rd.IsNull() {
				data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].Rd = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].Rd = types.StringNull()
			}
			for cci := range data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies {
				keys := [...]string{"address-family"}
				keyValues := [...]string{data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].AddressFamily.ValueString()}

				var ccr gjson.Result
				cr.Get("address-family").ForEach(
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
							ccr = v
							return false
						}
						return true
					},
				)
				if value := ccr.Get("address-family"); value.Exists() && !data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].AddressFamily.IsNull() {
					data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].AddressFamily = types.StringValue(value.String())
				} else {
					data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].AddressFamily = types.StringNull()
				}
				for c3i := range data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets {
					keys := [...]string{"id"}
					keyValues := [...]string{strconv.FormatInt(data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets[c3i].Id.ValueInt64(), 10)}

					var c3r gjson.Result
					ccr.Get("vpn-targets.vpn-target").ForEach(
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
								c3r = v
								return false
							}
							return true
						},
					)
					if value := c3r.Get("id"); value.Exists() && !data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets[c3i].Id.IsNull() {
						data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets[c3i].Id = types.Int64Value(value.Int())
					} else {
						data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets[c3i].Id = types.Int64Null()
					}
					if value := c3r.Get("route-target-type"); value.Exists() && !data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets[c3i].RouteTargetType.IsNull() {
						data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets[c3i].RouteTargetType = types.StringValue(value.String())
					} else {
						data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets[c3i].RouteTargetType = types.StringNull()
					}
					for c4i := range data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets[c3i].RouteTargets {
						keys := [...]string{"route-target"}
						keyValues := [...]string{data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets[c3i].RouteTargets[c4i].RouteTarget.ValueString()}

						var c4r gjson.Result
						c3r.Get("route-targets").ForEach(
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
									c4r = v
									return false
								}
								return true
							},
						)
						if value := c4r.Get("route-target"); value.Exists() && !data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets[c3i].RouteTargets[c4i].RouteTarget.IsNull() {
							data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets[c3i].RouteTargets[c4i].RouteTarget = types.StringValue(value.String())
						} else {
							data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].VpnTargets[c3i].RouteTargets[c4i].RouteTarget = types.StringNull()
						}
					}
				}
				if value := ccr.Get("vpn-targets.vpn-policies.import-policy"); value.Exists() && !data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].ImportPolicy.IsNull() {
					data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].ImportPolicy = types.StringValue(value.String())
				} else {
					data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].ImportPolicy = types.StringNull()
				}
				if value := ccr.Get("vpn-targets.vpn-policies.export-policy"); value.Exists() && !data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].ExportPolicy.IsNull() {
					data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].ExportPolicy = types.StringValue(value.String())
				} else {
					data.VpnNodes[i].ActiveVpnInstanceProfiles[ci].AddressFamilies[cci].ExportPolicy = types.StringNull()
				}
			}
		}
		for ci := range data.VpnNodes[i].VpnNetworkAccesses {
			keys := [...]string{"id"}
			keyValues := [...]string{data.VpnNodes[i].VpnNetworkAccesses[ci].Id.ValueString()}

			var cr gjson.Result
			r.Get("vpn-network-accesses.vpn-network-access").ForEach(
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
						cr = v
						return false
					}
					return true
				},
			)
			if value := cr.Get("id"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].Id.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Id = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Id = types.StringNull()
			}
			if value := cr.Get("interface-id"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].InterfaceId.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].InterfaceId = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].InterfaceId = types.StringNull()
			}
			if value := cr.Get("connection.encapsulation.type"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].EncapsulationType.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].EncapsulationType = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].EncapsulationType = types.StringNull()
			}
			if value := cr.Get("connection.encapsulation.dot1q.cvlan-id"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].CvlanId.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].CvlanId = types.Int64Value(value.Int())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].CvlanId = types.Int64Null()
			}
			if value := cr.Get("ip-connection.ipv4.local-address"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].Ipv4Address.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Ipv4Address = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Ipv4Address = types.StringNull()
			}
			if value := cr.Get("ip-connection.ipv4.prefix-length"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].Ipv4PrefixLength.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Ipv4PrefixLength = types.Int64Value(value.Int())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Ipv4PrefixLength = types.Int64Null()
			}
			if value := cr.Get("ip-connection.ipv6.local-address"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].Ipv6Address.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Ipv6Address = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Ipv6Address = types.StringNull()
			}
			if value := cr.Get("ip-connection.ipv6.prefix-length"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].Ipv6PrefixLength.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Ipv6PrefixLength = types.Int64Value(value.Int())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Ipv6PrefixLength = types.Int64Null()
			}
			if value := cr.Get("service.mtu"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].Mtu.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Mtu = types.Int64Value(value.Int())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Mtu = types.Int64Null()
			}
			for cci := range data.VpnNodes[i].VpnNetworkAccesses[ci].QosProfiles {
				keys := [...]string{"profile"}
				keyValues := [...]string{data.VpnNodes[i].VpnNetworkAccesses[ci].QosProfiles[cci].Profile.ValueString()}

				var ccr gjson.Result
				cr.Get("service.qos.qos-profile.qos-profile").ForEach(
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
							ccr = v
							return false
						}
						return true
					},
				)
				if value := ccr.Get("profile"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].QosProfiles[cci].Profile.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].QosProfiles[cci].Profile = types.StringValue(value.String())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].QosProfiles[cci].Profile = types.StringNull()
				}
				if value := ccr.Get("direction"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].QosProfiles[cci].Direction.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].QosProfiles[cci].Direction = types.StringValue(value.String())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].QosProfiles[cci].Direction = types.StringNull()
				}
			}
			for cci := range data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols {
				keys := [...]string{"id"}
				keyValues := [...]string{data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].Id.ValueString()}

				var ccr gjson.Result
				cr.Get("routing-protocols.routing-protocol").ForEach(
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
							ccr = v
							return false
						}
						return true
					},
				)
				if value := ccr.Get("id"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].Id.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].Id = types.StringValue(value.String())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].Id = types.StringNull()
				}
				if value := ccr.Get("type"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].Type.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].Type = types.StringValue(value.String())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].Type = types.StringNull()
				}
				if value := ccr.Get("bgp.peer-as"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpPeerAs.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpPeerAs = types.Int64Value(value.Int())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpPeerAs = types.Int64Null()
				}
				if value := ccr.Get("bgp.address-family"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpAddressFamily.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpAddressFamily = types.StringValue(value.String())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpAddressFamily = types.StringNull()
				}
				if value := ccr.Get("bgp.neighbor"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpNeighbors.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpNeighbors = helpers.GetStringSet(value.Array())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpNeighbors = types.SetNull(types.StringType)
				}
				if value := ccr.Get("bgp.multihop"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpMultihop.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpMultihop = types.Int64Value(value.Int())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpMultihop = types.Int64Null()
				}
				for c3i := range data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpRedistributeConnected {
					keys := [...]string{"address-family"}
					keyValues := [...]string{data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpRedistributeConnected[c3i].AddressFamily.ValueString()}

					var c3r gjson.Result
					ccr.Get("bgp.redistribute-connected").ForEach(
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
								c3r = v
								return false
							}
							return true
						},
					)
					if value := c3r.Get("address-family"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpRedistributeConnected[c3i].AddressFamily.IsNull() {
						data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpRedistributeConnected[c3i].AddressFamily = types.StringValue(value.String())
					} else {
						data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpRedistributeConnected[c3i].AddressFamily = types.StringNull()
					}
					if value := c3r.Get("enable"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpRedistributeConnected[c3i].Enable.IsNull() {
						data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpRedistributeConnected[c3i].Enable = types.BoolValue(value.Bool())
					} else {
						data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpRedistributeConnected[c3i].Enable = types.BoolNull()
					}
					if value := c3r.Get("cisco-l3vpn-ntw:metric"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpRedistributeConnected[c3i].Metric.IsNull() {
						data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpRedistributeConnected[c3i].Metric = types.Int64Value(value.Int())
					} else {
						data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpRedistributeConnected[c3i].Metric = types.Int64Null()
					}
				}
				if value := ccr.Get("bgp.cisco-l3vpn-ntw:update-source.if-type"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpUpdateSourceIfType.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpUpdateSourceIfType = types.StringValue(value.String())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpUpdateSourceIfType = types.StringNull()
				}
				if value := ccr.Get("bgp.cisco-l3vpn-ntw:update-source.if-id"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpUpdateSourceIfId.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpUpdateSourceIfId = types.StringValue(value.String())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpUpdateSourceIfId = types.StringNull()
				}
				if value := ccr.Get("bgp.cisco-l3vpn-ntw:update-source.sub-if-id"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpUpdateSourceSubIfId.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpUpdateSourceSubIfId = types.Int64Value(value.Int())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].RoutingProtocols[cci].BgpUpdateSourceSubIfId = types.Int64Null()
				}
			}
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *L3VPNVPNServices) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.VpnServiceTopology.IsNull() {
		return false
	}
	if !data.AutoHubRt.IsNull() {
		return false
	}
	if !data.AutoSpokeRt.IsNull() {
		return false
	}
	if !data.ServiceAssuranceProfileName.IsNull() {
		return false
	}
	if !data.ServiceAssuranceRuleName.IsNull() {
		return false
	}
	if !data.ServiceAssuranceMonitoringState.IsNull() {
		return false
	}
	if !data.ServiceAssurancePreservation.IsNull() {
		return false
	}
	if len(data.VpnInstanceProfiles) > 0 {
		return false
	}
	if len(data.VpnNodes) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

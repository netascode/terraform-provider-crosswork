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

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type L2VPNVPNServices struct {
	Instance                        types.String               `tfsdk:"instance"`
	Id                              types.String               `tfsdk:"id"`
	VpnId                           types.String               `tfsdk:"vpn_id"`
	VpnType                         types.String               `tfsdk:"vpn_type"`
	EviId                           types.Int64                `tfsdk:"evi_id"`
	AutoEviId                       types.Bool                 `tfsdk:"auto_evi_id"`
	EviSource                       types.Int64                `tfsdk:"evi_source"`
	AutoEviSource                   types.Bool                 `tfsdk:"auto_evi_source"`
	EviTarget                       types.Int64                `tfsdk:"evi_target"`
	AutoEviTarget                   types.Bool                 `tfsdk:"auto_evi_target"`
	ServiceAssuranceProfileName     types.String               `tfsdk:"service_assurance_profile_name"`
	ServiceAssuranceRuleName        types.String               `tfsdk:"service_assurance_rule_name"`
	ServiceAssuranceMonitoringState types.String               `tfsdk:"service_assurance_monitoring_state"`
	ServiceAssurancePreservation    types.String               `tfsdk:"service_assurance_preservation"`
	VpnNodes                        []L2VPNVPNServicesVpnNodes `tfsdk:"vpn_nodes"`
}

type L2VPNVPNServicesVpnNodes struct {
	VpnNodeId                          types.String                                 `tfsdk:"vpn_node_id"`
	MultiHomeEthernetSegmentIdentifier types.String                                 `tfsdk:"multi_home_ethernet_segment_identifier"`
	SrPolicy                           types.String                                 `tfsdk:"sr_policy"`
	SrPolicyFallback                   types.String                                 `tfsdk:"sr_policy_fallback"`
	OdnRoutePolicies                   types.Set                                    `tfsdk:"odn_route_policies"`
	VpnNetworkAccesses                 []L2VPNVPNServicesVpnNodesVpnNetworkAccesses `tfsdk:"vpn_network_accesses"`
}

type L2VPNVPNServicesVpnNodesVpnNetworkAccesses struct {
	Id                    types.String                                                 `tfsdk:"id"`
	Description           types.String                                                 `tfsdk:"description"`
	InterfaceId           types.String                                                 `tfsdk:"interface_id"`
	EncapsulationType     types.String                                                 `tfsdk:"encapsulation_type"`
	CvlanId               types.Int64                                                  `tfsdk:"cvlan_id"`
	TagOperationPop       types.String                                                 `tfsdk:"tag_operation_pop"`
	TagOperationPush      types.Bool                                                   `tfsdk:"tag_operation_push"`
	TagOperationTranslate types.String                                                 `tfsdk:"tag_operation_translate"`
	TagOperationTag       types.Int64                                                  `tfsdk:"tag_operation_tag"`
	Mtu                   types.Int64                                                  `tfsdk:"mtu"`
	QosProfiles           []L2VPNVPNServicesVpnNodesVpnNetworkAccessesQosProfiles      `tfsdk:"qos_profiles"`
	OamMdName             types.String                                                 `tfsdk:"oam_md_name"`
	OamMdLevel            types.Int64                                                  `tfsdk:"oam_md_level"`
	OamMaid               types.String                                                 `tfsdk:"oam_maid"`
	OamMepId              types.Int64                                                  `tfsdk:"oam_mep_id"`
	OamMessagePeriod      types.String                                                 `tfsdk:"oam_message_period"`
	OamIdType             types.String                                                 `tfsdk:"oam_id_type"`
	OamY1731Profiles      []L2VPNVPNServicesVpnNodesVpnNetworkAccessesOamY1731Profiles `tfsdk:"oam_y1731_profiles"`
}

type L2VPNVPNServicesVpnNodesVpnNetworkAccessesQosProfiles struct {
	Profile   types.String `tfsdk:"profile"`
	Direction types.String `tfsdk:"direction"`
}
type L2VPNVPNServicesVpnNodesVpnNetworkAccessesOamY1731Profiles struct {
	Name             types.String `tfsdk:"name"`
	ScheduleInterval types.Int64  `tfsdk:"schedule_interval"`
	ScheduleDuration types.String `tfsdk:"schedule_duration"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data L2VPNVPNServices) getPath() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/ietf-l2vpn-ntw:l2vpn-ntw/vpn-services/vpn-service=%v", url.PathEscape(data.VpnId.ValueString()))
}

func (data L2VPNVPNServices) nsoDevices() nsoDeviceSet {
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

func (data L2VPNVPNServices) getPathPost() string {
	return "/crosswork/proxy/nso/restconf/data/ietf-l2vpn-ntw:l2vpn-ntw/vpn-services"
}

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

func (data L2VPNVPNServices) getPathPlan() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/ietf-l2vpn-ntw:l2vpn-ntw/vpn-services/cisco-l2vpn-ntw:vpn-service-plan=%v", url.PathEscape(data.VpnId.ValueString()))
}

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data L2VPNVPNServices) toBody(ctx context.Context, state L2VPNVPNServices) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.VpnId.IsNull() {
		body, _ = sjson.Set(body, "ietf-l2vpn-ntw:vpn-service.0.vpn-id", data.VpnId.ValueString())
	}
	if !data.VpnType.IsNull() {
		body, _ = sjson.Set(body, "ietf-l2vpn-ntw:vpn-service.0.vpn-type", data.VpnType.ValueString())
	}
	if !data.EviId.IsNull() {
		body, _ = sjson.Set(body, "ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:evi-id", data.EviId.ValueInt64())
	}
	if data.AutoEviId.ValueBool() {
		body, _ = sjson.SetRaw(body, "ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:auto-evi-id", "[null]")
	}
	if !data.EviSource.IsNull() {
		body, _ = sjson.Set(body, "ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:evi-source", data.EviSource.ValueInt64())
	}
	if data.AutoEviSource.ValueBool() {
		body, _ = sjson.SetRaw(body, "ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:auto-evi-source", "[null]")
	}
	if !data.EviTarget.IsNull() {
		body, _ = sjson.Set(body, "ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:evi-target", data.EviTarget.ValueInt64())
	}
	if data.AutoEviTarget.ValueBool() {
		body, _ = sjson.SetRaw(body, "ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:auto-evi-target", "[null]")
	}
	if !data.ServiceAssuranceProfileName.IsNull() {
		body, _ = sjson.Set(body, "ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:service-assurance.profile-name", data.ServiceAssuranceProfileName.ValueString())
	}
	if !data.ServiceAssuranceRuleName.IsNull() {
		body, _ = sjson.Set(body, "ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:service-assurance.rule-name", data.ServiceAssuranceRuleName.ValueString())
	}
	if !data.ServiceAssuranceMonitoringState.IsNull() {
		body, _ = sjson.Set(body, "ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:service-assurance.monitoring-state", data.ServiceAssuranceMonitoringState.ValueString())
	}
	if !data.ServiceAssurancePreservation.IsNull() {
		body, _ = sjson.Set(body, "ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:service-assurance.preservation", data.ServiceAssurancePreservation.ValueString())
	}
	if len(data.VpnNodes) > 0 {
		body, _ = sjson.Set(body, "ietf-l2vpn-ntw:vpn-service.0.vpn-nodes.vpn-node", []interface{}{})
		for _, item := range data.VpnNodes {
			itemBody := ""
			if !item.VpnNodeId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "vpn-node-id", item.VpnNodeId.ValueString())
			}
			if !item.MultiHomeEthernetSegmentIdentifier.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "cisco-l2vpn-ntw:multi-home.ethernet-segment-identifier", item.MultiHomeEthernetSegmentIdentifier.ValueString())
			}
			if !item.SrPolicy.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "cisco-l2vpn-ntw:te-service-mapping.te-mapping.sr-policy.policy", item.SrPolicy.ValueString())
			}
			if !item.SrPolicyFallback.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "cisco-l2vpn-ntw:te-service-mapping.te-mapping.sr-policy.fallback", item.SrPolicyFallback.ValueString())
			}
			if !item.OdnRoutePolicies.IsNull() {
				var values []string
				item.OdnRoutePolicies.ElementsAs(ctx, &values, false)
				itemBody, _ = sjson.Set(itemBody, "cisco-l2vpn-ntw:te-service-mapping.te-mapping.odn.route-policy", values)
			}
			if len(item.VpnNetworkAccesses) > 0 {
				itemBody, _ = sjson.Set(itemBody, "vpn-network-accesses.vpn-network-access", []interface{}{})
				for _, childItem := range item.VpnNetworkAccesses {
					itemChildBody := ""
					if !childItem.Id.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "id", childItem.Id.ValueString())
					}
					if !childItem.Description.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "description", childItem.Description.ValueString())
					}
					if !childItem.InterfaceId.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "interface-id", childItem.InterfaceId.ValueString())
					}
					if !childItem.EncapsulationType.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "connection.encapsulation.encap-type", childItem.EncapsulationType.ValueString())
					}
					if !childItem.CvlanId.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "connection.encapsulation.dot1q.cvlan-id", childItem.CvlanId.ValueInt64())
					}
					if !childItem.TagOperationPop.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "connection.encapsulation.dot1q.tag-operations.pop", childItem.TagOperationPop.ValueString())
					}
					if childItem.TagOperationPush.ValueBool() {
						itemChildBody, _ = sjson.SetRaw(itemChildBody, "connection.encapsulation.dot1q.tag-operations.push", "[null]")
					}
					if !childItem.TagOperationTranslate.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "connection.encapsulation.dot1q.tag-operations.translate", childItem.TagOperationTranslate.ValueString())
					}
					if !childItem.TagOperationTag.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "connection.encapsulation.dot1q.tag-operations.tag-1", childItem.TagOperationTag.ValueInt64())
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
					if !childItem.OamMdName.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "ethernet-service-oam.md-name", childItem.OamMdName.ValueString())
					}
					if !childItem.OamMdLevel.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "ethernet-service-oam.md-level", childItem.OamMdLevel.ValueInt64())
					}
					if !childItem.OamMaid.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "ethernet-service-oam.y-1731.0.maid", childItem.OamMaid.ValueString())
					}
					if !childItem.OamMepId.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "ethernet-service-oam.y-1731.0.mep-id", childItem.OamMepId.ValueInt64())
					}
					if !childItem.OamMessagePeriod.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "ethernet-service-oam.y-1731.0.message-period", childItem.OamMessagePeriod.ValueString())
					}
					if !childItem.OamIdType.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "ethernet-service-oam.y-1731.0.cisco-l2vpn-ntw:id-type", childItem.OamIdType.ValueString())
					}
					if len(childItem.OamY1731Profiles) > 0 {
						itemChildBody, _ = sjson.Set(itemChildBody, "ethernet-service-oam.y-1731.0.cisco-l2vpn-ntw:y-1731-profile", []interface{}{})
						for _, childChildItem := range childItem.OamY1731Profiles {
							itemChildChildBody := ""
							if !childChildItem.Name.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "name", childChildItem.Name.ValueString())
							}
							if !childChildItem.ScheduleInterval.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "schedule.interval", childChildItem.ScheduleInterval.ValueInt64())
							}
							if !childChildItem.ScheduleDuration.IsNull() {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "schedule.duration", childChildItem.ScheduleDuration.ValueString())
							}
							itemChildBody, _ = sjson.SetRaw(itemChildBody, "ethernet-service-oam.y-1731.0.cisco-l2vpn-ntw:y-1731-profile.-1", itemChildChildBody)
						}
					}
					itemBody, _ = sjson.SetRaw(itemBody, "vpn-network-accesses.vpn-network-access.-1", itemChildBody)
				}
			}
			body, _ = sjson.SetRaw(body, "ietf-l2vpn-ntw:vpn-service.0.vpn-nodes.vpn-node.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *L2VPNVPNServices) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.vpn-type"); value.Exists() {
		data.VpnType = types.StringValue(value.String())
	} else {
		data.VpnType = types.StringNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:evi-id"); value.Exists() {
		data.EviId = types.Int64Value(value.Int())
	} else {
		data.EviId = types.Int64Null()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:auto-evi-id"); value.Exists() {
		data.AutoEviId = types.BoolValue(true)
	} else {
		data.AutoEviId = types.BoolNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:evi-source"); value.Exists() {
		data.EviSource = types.Int64Value(value.Int())
	} else {
		data.EviSource = types.Int64Null()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:auto-evi-source"); value.Exists() {
		data.AutoEviSource = types.BoolValue(true)
	} else {
		data.AutoEviSource = types.BoolNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:evi-target"); value.Exists() {
		data.EviTarget = types.Int64Value(value.Int())
	} else {
		data.EviTarget = types.Int64Null()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:auto-evi-target"); value.Exists() {
		data.AutoEviTarget = types.BoolValue(true)
	} else {
		data.AutoEviTarget = types.BoolNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:service-assurance.profile-name"); value.Exists() {
		data.ServiceAssuranceProfileName = types.StringValue(value.String())
	} else {
		data.ServiceAssuranceProfileName = types.StringNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:service-assurance.rule-name"); value.Exists() {
		data.ServiceAssuranceRuleName = types.StringValue(value.String())
	} else {
		data.ServiceAssuranceRuleName = types.StringNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:service-assurance.monitoring-state"); value.Exists() {
		data.ServiceAssuranceMonitoringState = types.StringValue(value.String())
	} else {
		data.ServiceAssuranceMonitoringState = types.StringNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:service-assurance.preservation"); value.Exists() {
		data.ServiceAssurancePreservation = types.StringValue(value.String())
	} else {
		data.ServiceAssurancePreservation = types.StringNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.vpn-nodes.vpn-node"); value.Exists() && len(value.Array()) > 0 {
		data.VpnNodes = make([]L2VPNVPNServicesVpnNodes, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := L2VPNVPNServicesVpnNodes{}
			if cValue := v.Get("vpn-node-id"); cValue.Exists() {
				item.VpnNodeId = types.StringValue(cValue.String())
			} else {
				item.VpnNodeId = types.StringNull()
			}
			if cValue := v.Get("cisco-l2vpn-ntw:multi-home.ethernet-segment-identifier"); cValue.Exists() {
				item.MultiHomeEthernetSegmentIdentifier = types.StringValue(cValue.String())
			} else {
				item.MultiHomeEthernetSegmentIdentifier = types.StringNull()
			}
			if cValue := v.Get("cisco-l2vpn-ntw:te-service-mapping.te-mapping.sr-policy.policy"); cValue.Exists() {
				item.SrPolicy = types.StringValue(cValue.String())
			} else {
				item.SrPolicy = types.StringNull()
			}
			if cValue := v.Get("cisco-l2vpn-ntw:te-service-mapping.te-mapping.sr-policy.fallback"); cValue.Exists() {
				item.SrPolicyFallback = types.StringValue(cValue.String())
			} else {
				item.SrPolicyFallback = types.StringNull()
			}
			if cValue := v.Get("cisco-l2vpn-ntw:te-service-mapping.te-mapping.odn.route-policy"); cValue.Exists() && len(cValue.Array()) > 0 {
				item.OdnRoutePolicies = helpers.GetStringSet(cValue.Array())
			} else {
				item.OdnRoutePolicies = types.SetNull(types.StringType)
			}
			if cValue := v.Get("vpn-network-accesses.vpn-network-access"); cValue.Exists() && len(cValue.Array()) > 0 {
				item.VpnNetworkAccesses = make([]L2VPNVPNServicesVpnNodesVpnNetworkAccesses, 0)
				cValue.ForEach(func(ck, cv gjson.Result) bool {
					cItem := L2VPNVPNServicesVpnNodesVpnNetworkAccesses{}
					if cwValue := cv.Get("id"); cwValue.Exists() {
						cItem.Id = types.StringValue(cwValue.String())
					} else {
						cItem.Id = types.StringNull()
					}
					if cwValue := cv.Get("description"); cwValue.Exists() {
						cItem.Description = types.StringValue(cwValue.String())
					} else {
						cItem.Description = types.StringNull()
					}
					if cwValue := cv.Get("interface-id"); cwValue.Exists() {
						cItem.InterfaceId = types.StringValue(cwValue.String())
					} else {
						cItem.InterfaceId = types.StringNull()
					}
					if cwValue := cv.Get("connection.encapsulation.encap-type"); cwValue.Exists() {
						cItem.EncapsulationType = types.StringValue(cwValue.String())
					} else {
						cItem.EncapsulationType = types.StringNull()
					}
					if cwValue := cv.Get("connection.encapsulation.dot1q.cvlan-id"); cwValue.Exists() {
						cItem.CvlanId = types.Int64Value(cwValue.Int())
					} else {
						cItem.CvlanId = types.Int64Null()
					}
					if cwValue := cv.Get("connection.encapsulation.dot1q.tag-operations.pop"); cwValue.Exists() {
						cItem.TagOperationPop = types.StringValue(cwValue.String())
					} else {
						cItem.TagOperationPop = types.StringNull()
					}
					if cwValue := cv.Get("connection.encapsulation.dot1q.tag-operations.push"); cwValue.Exists() {
						cItem.TagOperationPush = types.BoolValue(true)
					} else {
						cItem.TagOperationPush = types.BoolNull()
					}
					if cwValue := cv.Get("connection.encapsulation.dot1q.tag-operations.translate"); cwValue.Exists() {
						cItem.TagOperationTranslate = types.StringValue(cwValue.String())
					} else {
						cItem.TagOperationTranslate = types.StringNull()
					}
					if cwValue := cv.Get("connection.encapsulation.dot1q.tag-operations.tag-1"); cwValue.Exists() {
						cItem.TagOperationTag = types.Int64Value(cwValue.Int())
					} else {
						cItem.TagOperationTag = types.Int64Null()
					}
					if cwValue := cv.Get("service.mtu"); cwValue.Exists() {
						cItem.Mtu = types.Int64Value(cwValue.Int())
					} else {
						cItem.Mtu = types.Int64Null()
					}
					if cwValue := cv.Get("service.qos.qos-profile.qos-profile"); cwValue.Exists() && len(cwValue.Array()) > 0 {
						cItem.QosProfiles = make([]L2VPNVPNServicesVpnNodesVpnNetworkAccessesQosProfiles, 0)
						cwValue.ForEach(func(cck, ccv gjson.Result) bool {
							cwItem := L2VPNVPNServicesVpnNodesVpnNetworkAccessesQosProfiles{}
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
					if cwValue := cv.Get("ethernet-service-oam.md-name"); cwValue.Exists() {
						cItem.OamMdName = types.StringValue(cwValue.String())
					} else {
						cItem.OamMdName = types.StringNull()
					}
					if cwValue := cv.Get("ethernet-service-oam.md-level"); cwValue.Exists() {
						cItem.OamMdLevel = types.Int64Value(cwValue.Int())
					} else {
						cItem.OamMdLevel = types.Int64Null()
					}
					if cwValue := cv.Get("ethernet-service-oam.y-1731.0.maid"); cwValue.Exists() {
						cItem.OamMaid = types.StringValue(cwValue.String())
					} else {
						cItem.OamMaid = types.StringNull()
					}
					if cwValue := cv.Get("ethernet-service-oam.y-1731.0.mep-id"); cwValue.Exists() {
						cItem.OamMepId = types.Int64Value(cwValue.Int())
					} else {
						cItem.OamMepId = types.Int64Null()
					}
					if cwValue := cv.Get("ethernet-service-oam.y-1731.0.message-period"); cwValue.Exists() {
						cItem.OamMessagePeriod = types.StringValue(cwValue.String())
					} else {
						cItem.OamMessagePeriod = types.StringNull()
					}
					if cwValue := cv.Get("ethernet-service-oam.y-1731.0.cisco-l2vpn-ntw:id-type"); cwValue.Exists() {
						cItem.OamIdType = types.StringValue(cwValue.String())
					} else {
						cItem.OamIdType = types.StringNull()
					}
					if cwValue := cv.Get("ethernet-service-oam.y-1731.0.cisco-l2vpn-ntw:y-1731-profile"); cwValue.Exists() && len(cwValue.Array()) > 0 {
						cItem.OamY1731Profiles = make([]L2VPNVPNServicesVpnNodesVpnNetworkAccessesOamY1731Profiles, 0)
						cwValue.ForEach(func(cck, ccv gjson.Result) bool {
							cwItem := L2VPNVPNServicesVpnNodesVpnNetworkAccessesOamY1731Profiles{}
							if ccwValue := ccv.Get("name"); ccwValue.Exists() {
								cwItem.Name = types.StringValue(ccwValue.String())
							} else {
								cwItem.Name = types.StringNull()
							}
							if ccwValue := ccv.Get("schedule.interval"); ccwValue.Exists() {
								cwItem.ScheduleInterval = types.Int64Value(ccwValue.Int())
							} else {
								cwItem.ScheduleInterval = types.Int64Null()
							}
							if ccwValue := ccv.Get("schedule.duration"); ccwValue.Exists() {
								cwItem.ScheduleDuration = types.StringValue(ccwValue.String())
							} else {
								cwItem.ScheduleDuration = types.StringNull()
							}
							cItem.OamY1731Profiles = append(cItem.OamY1731Profiles, cwItem)
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
func (data *L2VPNVPNServices) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.vpn-type"); value.Exists() && !data.VpnType.IsNull() {
		data.VpnType = types.StringValue(value.String())
	} else {
		data.VpnType = types.StringNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:evi-id"); value.Exists() && !data.EviId.IsNull() {
		data.EviId = types.Int64Value(value.Int())
	} else {
		data.EviId = types.Int64Null()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:auto-evi-id"); value.Exists() && !data.AutoEviId.IsNull() {
		data.AutoEviId = types.BoolValue(true)
	} else {
		data.AutoEviId = types.BoolNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:evi-source"); value.Exists() && !data.EviSource.IsNull() {
		data.EviSource = types.Int64Value(value.Int())
	} else {
		data.EviSource = types.Int64Null()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:auto-evi-source"); value.Exists() && !data.AutoEviSource.IsNull() {
		data.AutoEviSource = types.BoolValue(true)
	} else {
		data.AutoEviSource = types.BoolNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:evi-target"); value.Exists() && !data.EviTarget.IsNull() {
		data.EviTarget = types.Int64Value(value.Int())
	} else {
		data.EviTarget = types.Int64Null()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:auto-evi-target"); value.Exists() && !data.AutoEviTarget.IsNull() {
		data.AutoEviTarget = types.BoolValue(true)
	} else {
		data.AutoEviTarget = types.BoolNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:service-assurance.profile-name"); value.Exists() && !data.ServiceAssuranceProfileName.IsNull() {
		data.ServiceAssuranceProfileName = types.StringValue(value.String())
	} else {
		data.ServiceAssuranceProfileName = types.StringNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:service-assurance.rule-name"); value.Exists() && !data.ServiceAssuranceRuleName.IsNull() {
		data.ServiceAssuranceRuleName = types.StringValue(value.String())
	} else {
		data.ServiceAssuranceRuleName = types.StringNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:service-assurance.monitoring-state"); value.Exists() && !data.ServiceAssuranceMonitoringState.IsNull() {
		data.ServiceAssuranceMonitoringState = types.StringValue(value.String())
	} else {
		data.ServiceAssuranceMonitoringState = types.StringNull()
	}
	if value := res.Get("ietf-l2vpn-ntw:vpn-service.0.cisco-l2vpn-ntw:service-assurance.preservation"); value.Exists() && !data.ServiceAssurancePreservation.IsNull() {
		data.ServiceAssurancePreservation = types.StringValue(value.String())
	} else {
		data.ServiceAssurancePreservation = types.StringNull()
	}
	for i := range data.VpnNodes {
		keys := [...]string{"vpn-node-id"}
		keyValues := [...]string{data.VpnNodes[i].VpnNodeId.ValueString()}

		var r gjson.Result
		res.Get("ietf-l2vpn-ntw:vpn-service.0.vpn-nodes.vpn-node").ForEach(
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
		if value := r.Get("cisco-l2vpn-ntw:multi-home.ethernet-segment-identifier"); value.Exists() && !data.VpnNodes[i].MultiHomeEthernetSegmentIdentifier.IsNull() {
			data.VpnNodes[i].MultiHomeEthernetSegmentIdentifier = types.StringValue(value.String())
		} else {
			data.VpnNodes[i].MultiHomeEthernetSegmentIdentifier = types.StringNull()
		}
		if value := r.Get("cisco-l2vpn-ntw:te-service-mapping.te-mapping.sr-policy.policy"); value.Exists() && !data.VpnNodes[i].SrPolicy.IsNull() {
			data.VpnNodes[i].SrPolicy = types.StringValue(value.String())
		} else {
			data.VpnNodes[i].SrPolicy = types.StringNull()
		}
		if value := r.Get("cisco-l2vpn-ntw:te-service-mapping.te-mapping.sr-policy.fallback"); value.Exists() && !data.VpnNodes[i].SrPolicyFallback.IsNull() {
			data.VpnNodes[i].SrPolicyFallback = types.StringValue(value.String())
		} else {
			data.VpnNodes[i].SrPolicyFallback = types.StringNull()
		}
		if value := r.Get("cisco-l2vpn-ntw:te-service-mapping.te-mapping.odn.route-policy"); value.Exists() && !data.VpnNodes[i].OdnRoutePolicies.IsNull() {
			data.VpnNodes[i].OdnRoutePolicies = helpers.GetStringSet(value.Array())
		} else {
			data.VpnNodes[i].OdnRoutePolicies = types.SetNull(types.StringType)
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
			if value := cr.Get("description"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].Description.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Description = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].Description = types.StringNull()
			}
			if value := cr.Get("interface-id"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].InterfaceId.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].InterfaceId = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].InterfaceId = types.StringNull()
			}
			if value := cr.Get("connection.encapsulation.encap-type"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].EncapsulationType.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].EncapsulationType = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].EncapsulationType = types.StringNull()
			}
			if value := cr.Get("connection.encapsulation.dot1q.cvlan-id"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].CvlanId.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].CvlanId = types.Int64Value(value.Int())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].CvlanId = types.Int64Null()
			}
			if value := cr.Get("connection.encapsulation.dot1q.tag-operations.pop"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].TagOperationPop.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].TagOperationPop = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].TagOperationPop = types.StringNull()
			}
			if value := cr.Get("connection.encapsulation.dot1q.tag-operations.push"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].TagOperationPush.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].TagOperationPush = types.BoolValue(true)
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].TagOperationPush = types.BoolNull()
			}
			if value := cr.Get("connection.encapsulation.dot1q.tag-operations.translate"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].TagOperationTranslate.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].TagOperationTranslate = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].TagOperationTranslate = types.StringNull()
			}
			if value := cr.Get("connection.encapsulation.dot1q.tag-operations.tag-1"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].TagOperationTag.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].TagOperationTag = types.Int64Value(value.Int())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].TagOperationTag = types.Int64Null()
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
			if value := cr.Get("ethernet-service-oam.md-name"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].OamMdName.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].OamMdName = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].OamMdName = types.StringNull()
			}
			if value := cr.Get("ethernet-service-oam.md-level"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].OamMdLevel.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].OamMdLevel = types.Int64Value(value.Int())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].OamMdLevel = types.Int64Null()
			}
			if value := cr.Get("ethernet-service-oam.y-1731.0.maid"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].OamMaid.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].OamMaid = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].OamMaid = types.StringNull()
			}
			if value := cr.Get("ethernet-service-oam.y-1731.0.mep-id"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].OamMepId.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].OamMepId = types.Int64Value(value.Int())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].OamMepId = types.Int64Null()
			}
			if value := cr.Get("ethernet-service-oam.y-1731.0.message-period"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].OamMessagePeriod.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].OamMessagePeriod = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].OamMessagePeriod = types.StringNull()
			}
			if value := cr.Get("ethernet-service-oam.y-1731.0.cisco-l2vpn-ntw:id-type"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].OamIdType.IsNull() {
				data.VpnNodes[i].VpnNetworkAccesses[ci].OamIdType = types.StringValue(value.String())
			} else {
				data.VpnNodes[i].VpnNetworkAccesses[ci].OamIdType = types.StringNull()
			}
			for cci := range data.VpnNodes[i].VpnNetworkAccesses[ci].OamY1731Profiles {
				keys := [...]string{"name"}
				keyValues := [...]string{data.VpnNodes[i].VpnNetworkAccesses[ci].OamY1731Profiles[cci].Name.ValueString()}

				var ccr gjson.Result
				cr.Get("ethernet-service-oam.y-1731.0.cisco-l2vpn-ntw:y-1731-profile").ForEach(
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
				if value := ccr.Get("name"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].OamY1731Profiles[cci].Name.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].OamY1731Profiles[cci].Name = types.StringValue(value.String())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].OamY1731Profiles[cci].Name = types.StringNull()
				}
				if value := ccr.Get("schedule.interval"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].OamY1731Profiles[cci].ScheduleInterval.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].OamY1731Profiles[cci].ScheduleInterval = types.Int64Value(value.Int())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].OamY1731Profiles[cci].ScheduleInterval = types.Int64Null()
				}
				if value := ccr.Get("schedule.duration"); value.Exists() && !data.VpnNodes[i].VpnNetworkAccesses[ci].OamY1731Profiles[cci].ScheduleDuration.IsNull() {
					data.VpnNodes[i].VpnNetworkAccesses[ci].OamY1731Profiles[cci].ScheduleDuration = types.StringValue(value.String())
				} else {
					data.VpnNodes[i].VpnNetworkAccesses[ci].OamY1731Profiles[cci].ScheduleDuration = types.StringNull()
				}
			}
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *L2VPNVPNServices) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.VpnType.IsNull() {
		return false
	}
	if !data.EviId.IsNull() {
		return false
	}
	if !data.AutoEviId.IsNull() {
		return false
	}
	if !data.EviSource.IsNull() {
		return false
	}
	if !data.AutoEviSource.IsNull() {
		return false
	}
	if !data.EviTarget.IsNull() {
		return false
	}
	if !data.AutoEviTarget.IsNull() {
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
	if len(data.VpnNodes) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

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
type SRTEPolicies struct {
	Instance                                   types.String                           `tfsdk:"instance"`
	Id                                         types.String                           `tfsdk:"id"`
	Name                                       types.String                           `tfsdk:"name"`
	HeadEnd                                    types.String                           `tfsdk:"head_end"`
	TailEnd                                    types.String                           `tfsdk:"tail_end"`
	Color                                      types.Int64                            `tfsdk:"color"`
	BindingSid                                 types.Int64                            `tfsdk:"binding_sid"`
	Bandwidth                                  types.Int64                            `tfsdk:"bandwidth"`
	SourceAddress                              types.String                           `tfsdk:"source_address"`
	Srv6LocatorName                            types.String                           `tfsdk:"srv6_locator_name"`
	Srv6LocatorBehavior                        types.String                           `tfsdk:"srv6_locator_behavior"`
	Srv6LocatorBindingSidType                  types.String                           `tfsdk:"srv6_locator_binding_sid_type"`
	DelayMeasurementProfile                    types.String                           `tfsdk:"delay_measurement_profile"`
	DelayMeasurementLoggingDelayExceeded       types.Bool                             `tfsdk:"delay_measurement_logging_delay_exceeded"`
	LivenessDetectionProfile                   types.String                           `tfsdk:"liveness_detection_profile"`
	LivenessDetectionBackup                    types.String                           `tfsdk:"liveness_detection_backup"`
	LivenessDetectionInvalidationAction        types.String                           `tfsdk:"liveness_detection_invalidation_action"`
	LivenessDetectionLoggingSessionStateChange types.Bool                             `tfsdk:"liveness_detection_logging_session_state_change"`
	PerformanceMeasurementReversePathLabel     types.Int64                            `tfsdk:"performance_measurement_reverse_path_label"`
	AutoRouteMetricRelative                    types.Int64                            `tfsdk:"auto_route_metric_relative"`
	AutoRouteMetricConstant                    types.Int64                            `tfsdk:"auto_route_metric_constant"`
	AutoRouteIncludePrefixes                   []SRTEPoliciesAutoRouteIncludePrefixes `tfsdk:"auto_route_include_prefixes"`
	AutoRouteForceSrInclude                    types.Bool                             `tfsdk:"auto_route_force_sr_include"`
	AutoRouteForwardClass                      types.Int64                            `tfsdk:"auto_route_forward_class"`
	Paths                                      []SRTEPoliciesPaths                    `tfsdk:"paths"`
}

type SRTEPoliciesAutoRouteIncludePrefixes struct {
	PrefixAddress types.String `tfsdk:"prefix_address"`
}

type SRTEPoliciesPaths struct {
	Preference                   types.Int64                              `tfsdk:"preference"`
	ExplicitSidLists             []SRTEPoliciesPathsExplicitSidLists      `tfsdk:"explicit_sid_lists"`
	ExplicitAffinityRules        []SRTEPoliciesPathsExplicitAffinityRules `tfsdk:"explicit_affinity_rules"`
	ExplicitSegmentsSidAlgorithm types.Int64                              `tfsdk:"explicit_segments_sid_algorithm"`
	ExplicitDisjointPathType     types.String                             `tfsdk:"explicit_disjoint_path_type"`
	ExplicitDisjointPathGroupId  types.Int64                              `tfsdk:"explicit_disjoint_path_group_id"`
	ExplicitDisjointPathSubId    types.Int64                              `tfsdk:"explicit_disjoint_path_sub_id"`
	DynamicPce                   types.Bool                               `tfsdk:"dynamic_pce"`
	DynamicMetricType            types.String                             `tfsdk:"dynamic_metric_type"`
	DynamicMetricMarginAbsolute  types.Int64                              `tfsdk:"dynamic_metric_margin_absolute"`
	DynamicMetricMarginRelative  types.Int64                              `tfsdk:"dynamic_metric_margin_relative"`
	DynamicSidLimit              types.Int64                              `tfsdk:"dynamic_sid_limit"`
	DynamicAffinityRules         []SRTEPoliciesPathsDynamicAffinityRules  `tfsdk:"dynamic_affinity_rules"`
	DynamicSegmentsSidAlgorithm  types.Int64                              `tfsdk:"dynamic_segments_sid_algorithm"`
	DynamicDisjointPathType      types.String                             `tfsdk:"dynamic_disjoint_path_type"`
	DynamicDisjointPathGroupId   types.Int64                              `tfsdk:"dynamic_disjoint_path_group_id"`
	DynamicDisjointPathSubId     types.Int64                              `tfsdk:"dynamic_disjoint_path_sub_id"`
}

type SRTEPoliciesPathsExplicitSidLists struct {
	Name   types.String `tfsdk:"name"`
	Weight types.Int64  `tfsdk:"weight"`
}
type SRTEPoliciesPathsExplicitAffinityRules struct {
	Action types.String `tfsdk:"action"`
	Colors types.Set    `tfsdk:"colors"`
}
type SRTEPoliciesPathsDynamicAffinityRules struct {
	Action types.String `tfsdk:"action"`
	Colors types.Set    `tfsdk:"colors"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data SRTEPolicies) getPath() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/cisco-sr-te-cfp:sr-te/cisco-sr-te-cfp-sr-policies:policies/policy=%v", url.PathEscape(data.Name.ValueString()))
}

func (data SRTEPolicies) nsoDevices() nsoDeviceSet {
	devices := nsoDeviceSet{}
	nsoAddDevice(devices, data.HeadEnd)
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

func (data SRTEPolicies) getPathPost() string {
	return "/crosswork/proxy/nso/restconf/data/cisco-sr-te-cfp:sr-te/cisco-sr-te-cfp-sr-policies:policies"
}

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

func (data SRTEPolicies) getPathPlan() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/cisco-sr-te-cfp:sr-te/cisco-sr-te-cfp-sr-policies:policies/policy-plan=%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data SRTEPolicies) toBody(ctx context.Context, state SRTEPolicies) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.name", data.Name.ValueString())
	}
	if !data.HeadEnd.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.head-end.0.name", data.HeadEnd.ValueString())
	}
	if !data.TailEnd.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.tail-end", data.TailEnd.ValueString())
	}
	if !data.Color.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.color", data.Color.ValueInt64())
	}
	if !data.BindingSid.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.binding-sid", data.BindingSid.ValueInt64())
	}
	if !data.Bandwidth.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.bandwidth", data.Bandwidth.ValueInt64())
	}
	if !data.SourceAddress.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.source-address", data.SourceAddress.ValueString())
	}
	if !data.Srv6LocatorName.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.srv6.locator.locator-name", data.Srv6LocatorName.ValueString())
	}
	if !data.Srv6LocatorBehavior.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.srv6.locator.behavior", data.Srv6LocatorBehavior.ValueString())
	}
	if !data.Srv6LocatorBindingSidType.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.srv6.locator.binding-sid-type", data.Srv6LocatorBindingSidType.ValueString())
	}
	if !data.DelayMeasurementProfile.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.delay-measurement.profile", data.DelayMeasurementProfile.ValueString())
	}
	if data.DelayMeasurementLoggingDelayExceeded.ValueBool() {
		body, _ = sjson.SetRaw(body, "cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.delay-measurement.logging.delay-exceeded", "[null]")
	}
	if !data.LivenessDetectionProfile.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.liveness-detection.profile", data.LivenessDetectionProfile.ValueString())
	}
	if !data.LivenessDetectionBackup.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.liveness-detection.backup", data.LivenessDetectionBackup.ValueString())
	}
	if !data.LivenessDetectionInvalidationAction.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.liveness-detection.invalidation-action", data.LivenessDetectionInvalidationAction.ValueString())
	}
	if data.LivenessDetectionLoggingSessionStateChange.ValueBool() {
		body, _ = sjson.SetRaw(body, "cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.liveness-detection.logging.session-state-change", "[null]")
	}
	if !data.PerformanceMeasurementReversePathLabel.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.reverse-path-label", data.PerformanceMeasurementReversePathLabel.ValueInt64())
	}
	if !data.AutoRouteMetricRelative.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.auto-route.auto-route-metric.metric-relative-value", data.AutoRouteMetricRelative.ValueInt64())
	}
	if !data.AutoRouteMetricConstant.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.auto-route.auto-route-metric.metric-constant-value", data.AutoRouteMetricConstant.ValueInt64())
	}
	if len(data.AutoRouteIncludePrefixes) > 0 {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.auto-route.include-prefixes.include-prefix", []interface{}{})
		for _, item := range data.AutoRouteIncludePrefixes {
			itemBody := ""
			if !item.PrefixAddress.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "prefix-address", item.PrefixAddress.ValueString())
			}
			body, _ = sjson.SetRaw(body, "cisco-sr-te-cfp-sr-policies:policy.0.auto-route.include-prefixes.include-prefix.-1", itemBody)
		}
	}
	if data.AutoRouteForceSrInclude.ValueBool() {
		body, _ = sjson.SetRaw(body, "cisco-sr-te-cfp-sr-policies:policy.0.auto-route.force-sr-include", "[null]")
	}
	if !data.AutoRouteForwardClass.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.auto-route.forward-class", data.AutoRouteForwardClass.ValueInt64())
	}
	if len(data.Paths) > 0 {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-policies:policy.0.path", []interface{}{})
		for _, item := range data.Paths {
			itemBody := ""
			if !item.Preference.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "preference", item.Preference.ValueInt64())
			}
			if len(item.ExplicitSidLists) > 0 {
				itemBody, _ = sjson.Set(itemBody, "explicit.sid-list", []interface{}{})
				for _, childItem := range item.ExplicitSidLists {
					itemChildBody := ""
					if !childItem.Name.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "name", childItem.Name.ValueString())
					}
					if !childItem.Weight.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "weight", childItem.Weight.ValueInt64())
					}
					itemBody, _ = sjson.SetRaw(itemBody, "explicit.sid-list.-1", itemChildBody)
				}
			}
			if len(item.ExplicitAffinityRules) > 0 {
				itemBody, _ = sjson.Set(itemBody, "explicit.constraints.affinity.rule", []interface{}{})
				for _, childItem := range item.ExplicitAffinityRules {
					itemChildBody := ""
					if !childItem.Action.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "action", childItem.Action.ValueString())
					}
					if !childItem.Colors.IsNull() {
						var values []string
						childItem.Colors.ElementsAs(ctx, &values, false)
						itemChildBody, _ = sjson.Set(itemChildBody, "color", values)
					}
					itemBody, _ = sjson.SetRaw(itemBody, "explicit.constraints.affinity.rule.-1", itemChildBody)
				}
			}
			if !item.ExplicitSegmentsSidAlgorithm.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "explicit.constraints.segments.sid-algorithm", item.ExplicitSegmentsSidAlgorithm.ValueInt64())
			}
			if !item.ExplicitDisjointPathType.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "explicit.constraints.disjoint-path.type", item.ExplicitDisjointPathType.ValueString())
			}
			if !item.ExplicitDisjointPathGroupId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "explicit.constraints.disjoint-path.group-id", item.ExplicitDisjointPathGroupId.ValueInt64())
			}
			if !item.ExplicitDisjointPathSubId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "explicit.constraints.disjoint-path.sub-id", item.ExplicitDisjointPathSubId.ValueInt64())
			}
			if item.DynamicPce.ValueBool() {
				itemBody, _ = sjson.SetRaw(itemBody, "dynamic.pce", "{}")
			}
			if !item.DynamicMetricType.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "dynamic.metric-type", item.DynamicMetricType.ValueString())
			}
			if !item.DynamicMetricMarginAbsolute.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "dynamic.metric-margin.absolute", item.DynamicMetricMarginAbsolute.ValueInt64())
			}
			if !item.DynamicMetricMarginRelative.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "dynamic.metric-margin.relative", item.DynamicMetricMarginRelative.ValueInt64())
			}
			if !item.DynamicSidLimit.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "dynamic.constraints.sid-limit", item.DynamicSidLimit.ValueInt64())
			}
			if len(item.DynamicAffinityRules) > 0 {
				itemBody, _ = sjson.Set(itemBody, "dynamic.constraints.affinity.rule", []interface{}{})
				for _, childItem := range item.DynamicAffinityRules {
					itemChildBody := ""
					if !childItem.Action.IsNull() {
						itemChildBody, _ = sjson.Set(itemChildBody, "action", childItem.Action.ValueString())
					}
					if !childItem.Colors.IsNull() {
						var values []string
						childItem.Colors.ElementsAs(ctx, &values, false)
						itemChildBody, _ = sjson.Set(itemChildBody, "color", values)
					}
					itemBody, _ = sjson.SetRaw(itemBody, "dynamic.constraints.affinity.rule.-1", itemChildBody)
				}
			}
			if !item.DynamicSegmentsSidAlgorithm.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "dynamic.constraints.segments.sid-algorithm", item.DynamicSegmentsSidAlgorithm.ValueInt64())
			}
			if !item.DynamicDisjointPathType.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "dynamic.constraints.disjoint-path.type", item.DynamicDisjointPathType.ValueString())
			}
			if !item.DynamicDisjointPathGroupId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "dynamic.constraints.disjoint-path.group-id", item.DynamicDisjointPathGroupId.ValueInt64())
			}
			if !item.DynamicDisjointPathSubId.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "dynamic.constraints.disjoint-path.sub-id", item.DynamicDisjointPathSubId.ValueInt64())
			}
			body, _ = sjson.SetRaw(body, "cisco-sr-te-cfp-sr-policies:policy.0.path.-1", itemBody)
		}
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *SRTEPolicies) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.head-end.0.name"); value.Exists() {
		data.HeadEnd = types.StringValue(value.String())
	} else {
		data.HeadEnd = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.tail-end"); value.Exists() {
		data.TailEnd = types.StringValue(value.String())
	} else {
		data.TailEnd = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.color"); value.Exists() {
		data.Color = types.Int64Value(value.Int())
	} else {
		data.Color = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.binding-sid"); value.Exists() {
		data.BindingSid = types.Int64Value(value.Int())
	} else {
		data.BindingSid = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.bandwidth"); value.Exists() {
		data.Bandwidth = types.Int64Value(value.Int())
	} else {
		data.Bandwidth = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.source-address"); value.Exists() {
		data.SourceAddress = types.StringValue(value.String())
	} else {
		data.SourceAddress = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.srv6.locator.locator-name"); value.Exists() {
		data.Srv6LocatorName = types.StringValue(value.String())
	} else {
		data.Srv6LocatorName = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.srv6.locator.behavior"); value.Exists() {
		data.Srv6LocatorBehavior = types.StringValue(value.String())
	} else {
		data.Srv6LocatorBehavior = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.srv6.locator.binding-sid-type"); value.Exists() {
		data.Srv6LocatorBindingSidType = types.StringValue(value.String())
	} else {
		data.Srv6LocatorBindingSidType = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.delay-measurement.profile"); value.Exists() {
		data.DelayMeasurementProfile = types.StringValue(value.String())
	} else {
		data.DelayMeasurementProfile = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.delay-measurement.logging.delay-exceeded"); value.Exists() {
		data.DelayMeasurementLoggingDelayExceeded = types.BoolValue(true)
	} else {
		data.DelayMeasurementLoggingDelayExceeded = types.BoolNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.liveness-detection.profile"); value.Exists() {
		data.LivenessDetectionProfile = types.StringValue(value.String())
	} else {
		data.LivenessDetectionProfile = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.liveness-detection.backup"); value.Exists() {
		data.LivenessDetectionBackup = types.StringValue(value.String())
	} else {
		data.LivenessDetectionBackup = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.liveness-detection.invalidation-action"); value.Exists() {
		data.LivenessDetectionInvalidationAction = types.StringValue(value.String())
	} else {
		data.LivenessDetectionInvalidationAction = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.liveness-detection.logging.session-state-change"); value.Exists() {
		data.LivenessDetectionLoggingSessionStateChange = types.BoolValue(true)
	} else {
		data.LivenessDetectionLoggingSessionStateChange = types.BoolNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.reverse-path-label"); value.Exists() {
		data.PerformanceMeasurementReversePathLabel = types.Int64Value(value.Int())
	} else {
		data.PerformanceMeasurementReversePathLabel = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.auto-route.auto-route-metric.metric-relative-value"); value.Exists() {
		data.AutoRouteMetricRelative = types.Int64Value(value.Int())
	} else {
		data.AutoRouteMetricRelative = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.auto-route.auto-route-metric.metric-constant-value"); value.Exists() {
		data.AutoRouteMetricConstant = types.Int64Value(value.Int())
	} else {
		data.AutoRouteMetricConstant = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.auto-route.include-prefixes.include-prefix"); value.Exists() && len(value.Array()) > 0 {
		data.AutoRouteIncludePrefixes = make([]SRTEPoliciesAutoRouteIncludePrefixes, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := SRTEPoliciesAutoRouteIncludePrefixes{}
			if cValue := v.Get("prefix-address"); cValue.Exists() {
				item.PrefixAddress = types.StringValue(cValue.String())
			} else {
				item.PrefixAddress = types.StringNull()
			}
			data.AutoRouteIncludePrefixes = append(data.AutoRouteIncludePrefixes, item)
			return true
		})
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.auto-route.force-sr-include"); value.Exists() {
		data.AutoRouteForceSrInclude = types.BoolValue(true)
	} else {
		data.AutoRouteForceSrInclude = types.BoolNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.auto-route.forward-class"); value.Exists() {
		data.AutoRouteForwardClass = types.Int64Value(value.Int())
	} else {
		data.AutoRouteForwardClass = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.path"); value.Exists() && len(value.Array()) > 0 {
		data.Paths = make([]SRTEPoliciesPaths, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := SRTEPoliciesPaths{}
			if cValue := v.Get("preference"); cValue.Exists() {
				item.Preference = types.Int64Value(cValue.Int())
			} else {
				item.Preference = types.Int64Null()
			}
			if cValue := v.Get("explicit.sid-list"); cValue.Exists() && len(cValue.Array()) > 0 {
				item.ExplicitSidLists = make([]SRTEPoliciesPathsExplicitSidLists, 0)
				cValue.ForEach(func(ck, cv gjson.Result) bool {
					cItem := SRTEPoliciesPathsExplicitSidLists{}
					if cwValue := cv.Get("name"); cwValue.Exists() {
						cItem.Name = types.StringValue(cwValue.String())
					} else {
						cItem.Name = types.StringNull()
					}
					if cwValue := cv.Get("weight"); cwValue.Exists() {
						cItem.Weight = types.Int64Value(cwValue.Int())
					} else {
						cItem.Weight = types.Int64Null()
					}
					item.ExplicitSidLists = append(item.ExplicitSidLists, cItem)
					return true
				})
			}
			if cValue := v.Get("explicit.constraints.affinity.rule"); cValue.Exists() && len(cValue.Array()) > 0 {
				item.ExplicitAffinityRules = make([]SRTEPoliciesPathsExplicitAffinityRules, 0)
				cValue.ForEach(func(ck, cv gjson.Result) bool {
					cItem := SRTEPoliciesPathsExplicitAffinityRules{}
					if cwValue := cv.Get("action"); cwValue.Exists() {
						cItem.Action = types.StringValue(cwValue.String())
					} else {
						cItem.Action = types.StringNull()
					}
					if cwValue := cv.Get("color"); cwValue.Exists() && len(cwValue.Array()) > 0 {
						cItem.Colors = helpers.GetStringSet(cwValue.Array())
					} else {
						cItem.Colors = types.SetNull(types.StringType)
					}
					item.ExplicitAffinityRules = append(item.ExplicitAffinityRules, cItem)
					return true
				})
			}
			if cValue := v.Get("explicit.constraints.segments.sid-algorithm"); cValue.Exists() {
				item.ExplicitSegmentsSidAlgorithm = types.Int64Value(cValue.Int())
			} else {
				item.ExplicitSegmentsSidAlgorithm = types.Int64Null()
			}
			if cValue := v.Get("explicit.constraints.disjoint-path.type"); cValue.Exists() {
				item.ExplicitDisjointPathType = types.StringValue(cValue.String())
			} else {
				item.ExplicitDisjointPathType = types.StringNull()
			}
			if cValue := v.Get("explicit.constraints.disjoint-path.group-id"); cValue.Exists() {
				item.ExplicitDisjointPathGroupId = types.Int64Value(cValue.Int())
			} else {
				item.ExplicitDisjointPathGroupId = types.Int64Null()
			}
			if cValue := v.Get("explicit.constraints.disjoint-path.sub-id"); cValue.Exists() {
				item.ExplicitDisjointPathSubId = types.Int64Value(cValue.Int())
			} else {
				item.ExplicitDisjointPathSubId = types.Int64Null()
			}
			if cValue := v.Get("dynamic.pce"); cValue.Exists() {
				item.DynamicPce = types.BoolValue(true)
			} else {
				item.DynamicPce = types.BoolNull()
			}
			if cValue := v.Get("dynamic.metric-type"); cValue.Exists() {
				item.DynamicMetricType = types.StringValue(cValue.String())
			} else {
				item.DynamicMetricType = types.StringNull()
			}
			if cValue := v.Get("dynamic.metric-margin.absolute"); cValue.Exists() {
				item.DynamicMetricMarginAbsolute = types.Int64Value(cValue.Int())
			} else {
				item.DynamicMetricMarginAbsolute = types.Int64Null()
			}
			if cValue := v.Get("dynamic.metric-margin.relative"); cValue.Exists() {
				item.DynamicMetricMarginRelative = types.Int64Value(cValue.Int())
			} else {
				item.DynamicMetricMarginRelative = types.Int64Null()
			}
			if cValue := v.Get("dynamic.constraints.sid-limit"); cValue.Exists() {
				item.DynamicSidLimit = types.Int64Value(cValue.Int())
			} else {
				item.DynamicSidLimit = types.Int64Null()
			}
			if cValue := v.Get("dynamic.constraints.affinity.rule"); cValue.Exists() && len(cValue.Array()) > 0 {
				item.DynamicAffinityRules = make([]SRTEPoliciesPathsDynamicAffinityRules, 0)
				cValue.ForEach(func(ck, cv gjson.Result) bool {
					cItem := SRTEPoliciesPathsDynamicAffinityRules{}
					if cwValue := cv.Get("action"); cwValue.Exists() {
						cItem.Action = types.StringValue(cwValue.String())
					} else {
						cItem.Action = types.StringNull()
					}
					if cwValue := cv.Get("color"); cwValue.Exists() && len(cwValue.Array()) > 0 {
						cItem.Colors = helpers.GetStringSet(cwValue.Array())
					} else {
						cItem.Colors = types.SetNull(types.StringType)
					}
					item.DynamicAffinityRules = append(item.DynamicAffinityRules, cItem)
					return true
				})
			}
			if cValue := v.Get("dynamic.constraints.segments.sid-algorithm"); cValue.Exists() {
				item.DynamicSegmentsSidAlgorithm = types.Int64Value(cValue.Int())
			} else {
				item.DynamicSegmentsSidAlgorithm = types.Int64Null()
			}
			if cValue := v.Get("dynamic.constraints.disjoint-path.type"); cValue.Exists() {
				item.DynamicDisjointPathType = types.StringValue(cValue.String())
			} else {
				item.DynamicDisjointPathType = types.StringNull()
			}
			if cValue := v.Get("dynamic.constraints.disjoint-path.group-id"); cValue.Exists() {
				item.DynamicDisjointPathGroupId = types.Int64Value(cValue.Int())
			} else {
				item.DynamicDisjointPathGroupId = types.Int64Null()
			}
			if cValue := v.Get("dynamic.constraints.disjoint-path.sub-id"); cValue.Exists() {
				item.DynamicDisjointPathSubId = types.Int64Value(cValue.Int())
			} else {
				item.DynamicDisjointPathSubId = types.Int64Null()
			}
			data.Paths = append(data.Paths, item)
			return true
		})
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *SRTEPolicies) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.head-end.0.name"); value.Exists() && !data.HeadEnd.IsNull() {
		data.HeadEnd = types.StringValue(value.String())
	} else {
		data.HeadEnd = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.tail-end"); value.Exists() && !data.TailEnd.IsNull() {
		data.TailEnd = types.StringValue(value.String())
	} else {
		data.TailEnd = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.color"); value.Exists() && !data.Color.IsNull() {
		data.Color = types.Int64Value(value.Int())
	} else {
		data.Color = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.binding-sid"); value.Exists() && !data.BindingSid.IsNull() {
		data.BindingSid = types.Int64Value(value.Int())
	} else {
		data.BindingSid = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.bandwidth"); value.Exists() && !data.Bandwidth.IsNull() {
		data.Bandwidth = types.Int64Value(value.Int())
	} else {
		data.Bandwidth = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.source-address"); value.Exists() && !data.SourceAddress.IsNull() {
		data.SourceAddress = types.StringValue(value.String())
	} else {
		data.SourceAddress = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.srv6.locator.locator-name"); value.Exists() && !data.Srv6LocatorName.IsNull() {
		data.Srv6LocatorName = types.StringValue(value.String())
	} else {
		data.Srv6LocatorName = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.srv6.locator.behavior"); value.Exists() && !data.Srv6LocatorBehavior.IsNull() {
		data.Srv6LocatorBehavior = types.StringValue(value.String())
	} else {
		data.Srv6LocatorBehavior = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.srv6.locator.binding-sid-type"); value.Exists() && !data.Srv6LocatorBindingSidType.IsNull() {
		data.Srv6LocatorBindingSidType = types.StringValue(value.String())
	} else {
		data.Srv6LocatorBindingSidType = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.delay-measurement.profile"); value.Exists() && !data.DelayMeasurementProfile.IsNull() {
		data.DelayMeasurementProfile = types.StringValue(value.String())
	} else {
		data.DelayMeasurementProfile = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.delay-measurement.logging.delay-exceeded"); value.Exists() && !data.DelayMeasurementLoggingDelayExceeded.IsNull() {
		data.DelayMeasurementLoggingDelayExceeded = types.BoolValue(true)
	} else {
		data.DelayMeasurementLoggingDelayExceeded = types.BoolNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.liveness-detection.profile"); value.Exists() && !data.LivenessDetectionProfile.IsNull() {
		data.LivenessDetectionProfile = types.StringValue(value.String())
	} else {
		data.LivenessDetectionProfile = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.liveness-detection.backup"); value.Exists() && !data.LivenessDetectionBackup.IsNull() {
		data.LivenessDetectionBackup = types.StringValue(value.String())
	} else {
		data.LivenessDetectionBackup = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.liveness-detection.invalidation-action"); value.Exists() && !data.LivenessDetectionInvalidationAction.IsNull() {
		data.LivenessDetectionInvalidationAction = types.StringValue(value.String())
	} else {
		data.LivenessDetectionInvalidationAction = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.liveness-detection.logging.session-state-change"); value.Exists() && !data.LivenessDetectionLoggingSessionStateChange.IsNull() {
		data.LivenessDetectionLoggingSessionStateChange = types.BoolValue(true)
	} else {
		data.LivenessDetectionLoggingSessionStateChange = types.BoolNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.performance-measurement.reverse-path-label"); value.Exists() && !data.PerformanceMeasurementReversePathLabel.IsNull() {
		data.PerformanceMeasurementReversePathLabel = types.Int64Value(value.Int())
	} else {
		data.PerformanceMeasurementReversePathLabel = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.auto-route.auto-route-metric.metric-relative-value"); value.Exists() && !data.AutoRouteMetricRelative.IsNull() {
		data.AutoRouteMetricRelative = types.Int64Value(value.Int())
	} else {
		data.AutoRouteMetricRelative = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.auto-route.auto-route-metric.metric-constant-value"); value.Exists() && !data.AutoRouteMetricConstant.IsNull() {
		data.AutoRouteMetricConstant = types.Int64Value(value.Int())
	} else {
		data.AutoRouteMetricConstant = types.Int64Null()
	}
	for i := range data.AutoRouteIncludePrefixes {
		keys := [...]string{"prefix-address"}
		keyValues := [...]string{data.AutoRouteIncludePrefixes[i].PrefixAddress.ValueString()}

		var r gjson.Result
		res.Get("cisco-sr-te-cfp-sr-policies:policy.0.auto-route.include-prefixes.include-prefix").ForEach(
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
		if value := r.Get("prefix-address"); value.Exists() && !data.AutoRouteIncludePrefixes[i].PrefixAddress.IsNull() {
			data.AutoRouteIncludePrefixes[i].PrefixAddress = types.StringValue(value.String())
		} else {
			data.AutoRouteIncludePrefixes[i].PrefixAddress = types.StringNull()
		}
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.auto-route.force-sr-include"); value.Exists() && !data.AutoRouteForceSrInclude.IsNull() {
		data.AutoRouteForceSrInclude = types.BoolValue(true)
	} else {
		data.AutoRouteForceSrInclude = types.BoolNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-policies:policy.0.auto-route.forward-class"); value.Exists() && !data.AutoRouteForwardClass.IsNull() {
		data.AutoRouteForwardClass = types.Int64Value(value.Int())
	} else {
		data.AutoRouteForwardClass = types.Int64Null()
	}
	for i := range data.Paths {
		keys := [...]string{"preference"}
		keyValues := [...]string{strconv.FormatInt(data.Paths[i].Preference.ValueInt64(), 10)}

		var r gjson.Result
		res.Get("cisco-sr-te-cfp-sr-policies:policy.0.path").ForEach(
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
		if value := r.Get("preference"); value.Exists() && !data.Paths[i].Preference.IsNull() {
			data.Paths[i].Preference = types.Int64Value(value.Int())
		} else {
			data.Paths[i].Preference = types.Int64Null()
		}
		for ci := range data.Paths[i].ExplicitSidLists {
			keys := [...]string{"name"}
			keyValues := [...]string{data.Paths[i].ExplicitSidLists[ci].Name.ValueString()}

			var cr gjson.Result
			r.Get("explicit.sid-list").ForEach(
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
			if value := cr.Get("name"); value.Exists() && !data.Paths[i].ExplicitSidLists[ci].Name.IsNull() {
				data.Paths[i].ExplicitSidLists[ci].Name = types.StringValue(value.String())
			} else {
				data.Paths[i].ExplicitSidLists[ci].Name = types.StringNull()
			}
			if value := cr.Get("weight"); value.Exists() && !data.Paths[i].ExplicitSidLists[ci].Weight.IsNull() {
				data.Paths[i].ExplicitSidLists[ci].Weight = types.Int64Value(value.Int())
			} else {
				data.Paths[i].ExplicitSidLists[ci].Weight = types.Int64Null()
			}
		}
		for ci := range data.Paths[i].ExplicitAffinityRules {
			keys := [...]string{"action"}
			keyValues := [...]string{data.Paths[i].ExplicitAffinityRules[ci].Action.ValueString()}

			var cr gjson.Result
			r.Get("explicit.constraints.affinity.rule").ForEach(
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
			if value := cr.Get("action"); value.Exists() && !data.Paths[i].ExplicitAffinityRules[ci].Action.IsNull() {
				data.Paths[i].ExplicitAffinityRules[ci].Action = types.StringValue(value.String())
			} else {
				data.Paths[i].ExplicitAffinityRules[ci].Action = types.StringNull()
			}
			if value := cr.Get("color"); value.Exists() && !data.Paths[i].ExplicitAffinityRules[ci].Colors.IsNull() {
				data.Paths[i].ExplicitAffinityRules[ci].Colors = helpers.GetStringSet(value.Array())
			} else {
				data.Paths[i].ExplicitAffinityRules[ci].Colors = types.SetNull(types.StringType)
			}
		}
		if value := r.Get("explicit.constraints.segments.sid-algorithm"); value.Exists() && !data.Paths[i].ExplicitSegmentsSidAlgorithm.IsNull() {
			data.Paths[i].ExplicitSegmentsSidAlgorithm = types.Int64Value(value.Int())
		} else {
			data.Paths[i].ExplicitSegmentsSidAlgorithm = types.Int64Null()
		}
		if value := r.Get("explicit.constraints.disjoint-path.type"); value.Exists() && !data.Paths[i].ExplicitDisjointPathType.IsNull() {
			data.Paths[i].ExplicitDisjointPathType = types.StringValue(value.String())
		} else {
			data.Paths[i].ExplicitDisjointPathType = types.StringNull()
		}
		if value := r.Get("explicit.constraints.disjoint-path.group-id"); value.Exists() && !data.Paths[i].ExplicitDisjointPathGroupId.IsNull() {
			data.Paths[i].ExplicitDisjointPathGroupId = types.Int64Value(value.Int())
		} else {
			data.Paths[i].ExplicitDisjointPathGroupId = types.Int64Null()
		}
		if value := r.Get("explicit.constraints.disjoint-path.sub-id"); value.Exists() && !data.Paths[i].ExplicitDisjointPathSubId.IsNull() {
			data.Paths[i].ExplicitDisjointPathSubId = types.Int64Value(value.Int())
		} else {
			data.Paths[i].ExplicitDisjointPathSubId = types.Int64Null()
		}
		if value := r.Get("dynamic.pce"); value.Exists() && !data.Paths[i].DynamicPce.IsNull() {
			data.Paths[i].DynamicPce = types.BoolValue(true)
		} else {
			data.Paths[i].DynamicPce = types.BoolNull()
		}
		if value := r.Get("dynamic.metric-type"); value.Exists() && !data.Paths[i].DynamicMetricType.IsNull() {
			data.Paths[i].DynamicMetricType = types.StringValue(value.String())
		} else {
			data.Paths[i].DynamicMetricType = types.StringNull()
		}
		if value := r.Get("dynamic.metric-margin.absolute"); value.Exists() && !data.Paths[i].DynamicMetricMarginAbsolute.IsNull() {
			data.Paths[i].DynamicMetricMarginAbsolute = types.Int64Value(value.Int())
		} else {
			data.Paths[i].DynamicMetricMarginAbsolute = types.Int64Null()
		}
		if value := r.Get("dynamic.metric-margin.relative"); value.Exists() && !data.Paths[i].DynamicMetricMarginRelative.IsNull() {
			data.Paths[i].DynamicMetricMarginRelative = types.Int64Value(value.Int())
		} else {
			data.Paths[i].DynamicMetricMarginRelative = types.Int64Null()
		}
		if value := r.Get("dynamic.constraints.sid-limit"); value.Exists() && !data.Paths[i].DynamicSidLimit.IsNull() {
			data.Paths[i].DynamicSidLimit = types.Int64Value(value.Int())
		} else {
			data.Paths[i].DynamicSidLimit = types.Int64Null()
		}
		for ci := range data.Paths[i].DynamicAffinityRules {
			keys := [...]string{"action"}
			keyValues := [...]string{data.Paths[i].DynamicAffinityRules[ci].Action.ValueString()}

			var cr gjson.Result
			r.Get("dynamic.constraints.affinity.rule").ForEach(
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
			if value := cr.Get("action"); value.Exists() && !data.Paths[i].DynamicAffinityRules[ci].Action.IsNull() {
				data.Paths[i].DynamicAffinityRules[ci].Action = types.StringValue(value.String())
			} else {
				data.Paths[i].DynamicAffinityRules[ci].Action = types.StringNull()
			}
			if value := cr.Get("color"); value.Exists() && !data.Paths[i].DynamicAffinityRules[ci].Colors.IsNull() {
				data.Paths[i].DynamicAffinityRules[ci].Colors = helpers.GetStringSet(value.Array())
			} else {
				data.Paths[i].DynamicAffinityRules[ci].Colors = types.SetNull(types.StringType)
			}
		}
		if value := r.Get("dynamic.constraints.segments.sid-algorithm"); value.Exists() && !data.Paths[i].DynamicSegmentsSidAlgorithm.IsNull() {
			data.Paths[i].DynamicSegmentsSidAlgorithm = types.Int64Value(value.Int())
		} else {
			data.Paths[i].DynamicSegmentsSidAlgorithm = types.Int64Null()
		}
		if value := r.Get("dynamic.constraints.disjoint-path.type"); value.Exists() && !data.Paths[i].DynamicDisjointPathType.IsNull() {
			data.Paths[i].DynamicDisjointPathType = types.StringValue(value.String())
		} else {
			data.Paths[i].DynamicDisjointPathType = types.StringNull()
		}
		if value := r.Get("dynamic.constraints.disjoint-path.group-id"); value.Exists() && !data.Paths[i].DynamicDisjointPathGroupId.IsNull() {
			data.Paths[i].DynamicDisjointPathGroupId = types.Int64Value(value.Int())
		} else {
			data.Paths[i].DynamicDisjointPathGroupId = types.Int64Null()
		}
		if value := r.Get("dynamic.constraints.disjoint-path.sub-id"); value.Exists() && !data.Paths[i].DynamicDisjointPathSubId.IsNull() {
			data.Paths[i].DynamicDisjointPathSubId = types.Int64Value(value.Int())
		} else {
			data.Paths[i].DynamicDisjointPathSubId = types.Int64Null()
		}
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *SRTEPolicies) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.HeadEnd.IsNull() {
		return false
	}
	if !data.TailEnd.IsNull() {
		return false
	}
	if !data.Color.IsNull() {
		return false
	}
	if !data.BindingSid.IsNull() {
		return false
	}
	if !data.Bandwidth.IsNull() {
		return false
	}
	if !data.SourceAddress.IsNull() {
		return false
	}
	if !data.Srv6LocatorName.IsNull() {
		return false
	}
	if !data.Srv6LocatorBehavior.IsNull() {
		return false
	}
	if !data.Srv6LocatorBindingSidType.IsNull() {
		return false
	}
	if !data.DelayMeasurementProfile.IsNull() {
		return false
	}
	if !data.DelayMeasurementLoggingDelayExceeded.IsNull() {
		return false
	}
	if !data.LivenessDetectionProfile.IsNull() {
		return false
	}
	if !data.LivenessDetectionBackup.IsNull() {
		return false
	}
	if !data.LivenessDetectionInvalidationAction.IsNull() {
		return false
	}
	if !data.LivenessDetectionLoggingSessionStateChange.IsNull() {
		return false
	}
	if !data.PerformanceMeasurementReversePathLabel.IsNull() {
		return false
	}
	if !data.AutoRouteMetricRelative.IsNull() {
		return false
	}
	if !data.AutoRouteMetricConstant.IsNull() {
		return false
	}
	if len(data.AutoRouteIncludePrefixes) > 0 {
		return false
	}
	if !data.AutoRouteForceSrInclude.IsNull() {
		return false
	}
	if !data.AutoRouteForwardClass.IsNull() {
		return false
	}
	if len(data.Paths) > 0 {
		return false
	}
	return true
}

// End of section. //template:end isNull

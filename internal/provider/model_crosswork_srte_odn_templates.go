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
type SRTEODNTemplates struct {
	Instance                                   types.String                           `tfsdk:"instance"`
	Id                                         types.String                           `tfsdk:"id"`
	Name                                       types.String                           `tfsdk:"name"`
	Color                                      types.Int64                            `tfsdk:"color"`
	HeadEnd                                    []SRTEODNTemplatesHeadEnd              `tfsdk:"head_end"`
	MaximumSidDepth                            types.Int64                            `tfsdk:"maximum_sid_depth"`
	PceGroup                                   types.String                           `tfsdk:"pce_group"`
	Bandwidth                                  types.Int64                            `tfsdk:"bandwidth"`
	SourceAddress                              types.String                           `tfsdk:"source_address"`
	Srv6LocatorName                            types.String                           `tfsdk:"srv6_locator_name"`
	Srv6LocatorBehavior                        types.String                           `tfsdk:"srv6_locator_behavior"`
	Srv6LocatorBindingSidType                  types.String                           `tfsdk:"srv6_locator_binding_sid_type"`
	DynamicPce                                 types.Bool                             `tfsdk:"dynamic_pce"`
	DynamicFlexAlg                             types.Int64                            `tfsdk:"dynamic_flex_alg"`
	DynamicMetricType                          types.String                           `tfsdk:"dynamic_metric_type"`
	DynamicMetricMarginAbsolute                types.Int64                            `tfsdk:"dynamic_metric_margin_absolute"`
	DynamicMetricMarginRelative                types.Int64                            `tfsdk:"dynamic_metric_margin_relative"`
	DynamicAffinityRules                       []SRTEODNTemplatesDynamicAffinityRules `tfsdk:"dynamic_affinity_rules"`
	DynamicSegmentsSidAlgorithm                types.Int64                            `tfsdk:"dynamic_segments_sid_algorithm"`
	DynamicDisjointPathType                    types.String                           `tfsdk:"dynamic_disjoint_path_type"`
	DynamicDisjointPathGroupId                 types.Int64                            `tfsdk:"dynamic_disjoint_path_group_id"`
	DynamicDisjointPathSubId                   types.Int64                            `tfsdk:"dynamic_disjoint_path_sub_id"`
	DelayMeasurementProfile                    types.String                           `tfsdk:"delay_measurement_profile"`
	DelayMeasurementLoggingDelayExceeded       types.Bool                             `tfsdk:"delay_measurement_logging_delay_exceeded"`
	LivenessDetectionProfile                   types.String                           `tfsdk:"liveness_detection_profile"`
	LivenessDetectionBackup                    types.String                           `tfsdk:"liveness_detection_backup"`
	LivenessDetectionInvalidationAction        types.String                           `tfsdk:"liveness_detection_invalidation_action"`
	LivenessDetectionLoggingSessionStateChange types.Bool                             `tfsdk:"liveness_detection_logging_session_state_change"`
	PerformanceMeasurementReversePathLabel     types.Int64                            `tfsdk:"performance_measurement_reverse_path_label"`
}

type SRTEODNTemplatesHeadEnd struct {
	Name types.String `tfsdk:"name"`
}

type SRTEODNTemplatesDynamicAffinityRules struct {
	Action types.String `tfsdk:"action"`
	Color  types.Set    `tfsdk:"color"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data SRTEODNTemplates) getPath() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/cisco-sr-te-cfp:sr-te/cisco-sr-te-cfp-sr-odn:odn/odn-template=%v", url.PathEscape(data.Name.ValueString()))
}

func (data SRTEODNTemplates) nsoDevices() nsoDeviceSet {
	devices := nsoDeviceSet{}
	for _, i1 := range data.HeadEnd {
		nsoAddDevice(devices, i1.Name)
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

func (data SRTEODNTemplates) getPathPost() string {
	return "/crosswork/proxy/nso/restconf/data/cisco-sr-te-cfp:sr-te/cisco-sr-te-cfp-sr-odn:odn"
}

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

func (data SRTEODNTemplates) getPathPlan() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/cisco-sr-te-cfp:sr-te/cisco-sr-te-cfp-sr-odn:odn/odn-template-plan=%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data SRTEODNTemplates) toBody(ctx context.Context, state SRTEODNTemplates) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.name", data.Name.ValueString())
	}
	if !data.Color.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.color", data.Color.ValueInt64())
	}
	if len(data.HeadEnd) > 0 {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.head-end", []interface{}{})
		for _, item := range data.HeadEnd {
			itemBody := ""
			if !item.Name.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "name", item.Name.ValueString())
			}
			body, _ = sjson.SetRaw(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.head-end.-1", itemBody)
		}
	}
	if !data.MaximumSidDepth.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.maximum-sid-depth", data.MaximumSidDepth.ValueInt64())
	}
	if !data.PceGroup.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.pce-group", data.PceGroup.ValueString())
	}
	if !data.Bandwidth.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.bandwidth", data.Bandwidth.ValueInt64())
	}
	if !data.SourceAddress.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.source-address", data.SourceAddress.ValueString())
	}
	if !data.Srv6LocatorName.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.srv6.locator.locator-name", data.Srv6LocatorName.ValueString())
	}
	if !data.Srv6LocatorBehavior.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.srv6.locator.behavior", data.Srv6LocatorBehavior.ValueString())
	}
	if !data.Srv6LocatorBindingSidType.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.srv6.locator.binding-sid-type", data.Srv6LocatorBindingSidType.ValueString())
	}
	if data.DynamicPce.ValueBool() {
		body, _ = sjson.SetRaw(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.pce", "{}")
	}
	if !data.DynamicFlexAlg.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.flex-alg", data.DynamicFlexAlg.ValueInt64())
	}
	if !data.DynamicMetricType.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.metric-type", data.DynamicMetricType.ValueString())
	}
	if !data.DynamicMetricMarginAbsolute.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.metric-margin.absolute", data.DynamicMetricMarginAbsolute.ValueInt64())
	}
	if !data.DynamicMetricMarginRelative.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.metric-margin.relative", data.DynamicMetricMarginRelative.ValueInt64())
	}
	if len(data.DynamicAffinityRules) > 0 {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.affinity.rule", []interface{}{})
		for _, item := range data.DynamicAffinityRules {
			itemBody := ""
			if !item.Action.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "action", item.Action.ValueString())
			}
			if !item.Color.IsNull() {
				var values []string
				item.Color.ElementsAs(ctx, &values, false)
				itemBody, _ = sjson.Set(itemBody, "color", values)
			}
			body, _ = sjson.SetRaw(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.affinity.rule.-1", itemBody)
		}
	}
	if !data.DynamicSegmentsSidAlgorithm.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.segments.sid-algorithm", data.DynamicSegmentsSidAlgorithm.ValueInt64())
	}
	if !data.DynamicDisjointPathType.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.disjoint-path.type", data.DynamicDisjointPathType.ValueString())
	}
	if !data.DynamicDisjointPathGroupId.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.disjoint-path.group-id", data.DynamicDisjointPathGroupId.ValueInt64())
	}
	if !data.DynamicDisjointPathSubId.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.disjoint-path.sub-id", data.DynamicDisjointPathSubId.ValueInt64())
	}
	if !data.DelayMeasurementProfile.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.delay-measurement.profile", data.DelayMeasurementProfile.ValueString())
	}
	if data.DelayMeasurementLoggingDelayExceeded.ValueBool() {
		body, _ = sjson.SetRaw(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.delay-measurement.logging.delay-exceeded", "[null]")
	}
	if !data.LivenessDetectionProfile.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.liveness-detection.profile", data.LivenessDetectionProfile.ValueString())
	}
	if !data.LivenessDetectionBackup.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.liveness-detection.backup", data.LivenessDetectionBackup.ValueString())
	}
	if !data.LivenessDetectionInvalidationAction.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.liveness-detection.invalidation-action", data.LivenessDetectionInvalidationAction.ValueString())
	}
	if data.LivenessDetectionLoggingSessionStateChange.ValueBool() {
		body, _ = sjson.SetRaw(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.liveness-detection.logging.session-state-change", "[null]")
	}
	if !data.PerformanceMeasurementReversePathLabel.IsNull() {
		body, _ = sjson.Set(body, "cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.reverse-path-label", data.PerformanceMeasurementReversePathLabel.ValueInt64())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *SRTEODNTemplates) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.color"); value.Exists() {
		data.Color = types.Int64Value(value.Int())
	} else {
		data.Color = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.head-end"); value.Exists() && len(value.Array()) > 0 {
		data.HeadEnd = make([]SRTEODNTemplatesHeadEnd, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := SRTEODNTemplatesHeadEnd{}
			if cValue := v.Get("name"); cValue.Exists() {
				item.Name = types.StringValue(cValue.String())
			} else {
				item.Name = types.StringNull()
			}
			data.HeadEnd = append(data.HeadEnd, item)
			return true
		})
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.maximum-sid-depth"); value.Exists() {
		data.MaximumSidDepth = types.Int64Value(value.Int())
	} else {
		data.MaximumSidDepth = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.pce-group"); value.Exists() {
		data.PceGroup = types.StringValue(value.String())
	} else {
		data.PceGroup = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.bandwidth"); value.Exists() {
		data.Bandwidth = types.Int64Value(value.Int())
	} else {
		data.Bandwidth = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.source-address"); value.Exists() {
		data.SourceAddress = types.StringValue(value.String())
	} else {
		data.SourceAddress = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.srv6.locator.locator-name"); value.Exists() {
		data.Srv6LocatorName = types.StringValue(value.String())
	} else {
		data.Srv6LocatorName = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.srv6.locator.behavior"); value.Exists() {
		data.Srv6LocatorBehavior = types.StringValue(value.String())
	} else {
		data.Srv6LocatorBehavior = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.srv6.locator.binding-sid-type"); value.Exists() {
		data.Srv6LocatorBindingSidType = types.StringValue(value.String())
	} else {
		data.Srv6LocatorBindingSidType = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.pce"); value.Exists() {
		data.DynamicPce = types.BoolValue(true)
	} else {
		data.DynamicPce = types.BoolNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.flex-alg"); value.Exists() {
		data.DynamicFlexAlg = types.Int64Value(value.Int())
	} else {
		data.DynamicFlexAlg = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.metric-type"); value.Exists() {
		data.DynamicMetricType = types.StringValue(value.String())
	} else {
		data.DynamicMetricType = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.metric-margin.absolute"); value.Exists() {
		data.DynamicMetricMarginAbsolute = types.Int64Value(value.Int())
	} else {
		data.DynamicMetricMarginAbsolute = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.metric-margin.relative"); value.Exists() {
		data.DynamicMetricMarginRelative = types.Int64Value(value.Int())
	} else {
		data.DynamicMetricMarginRelative = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.affinity.rule"); value.Exists() && len(value.Array()) > 0 {
		data.DynamicAffinityRules = make([]SRTEODNTemplatesDynamicAffinityRules, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := SRTEODNTemplatesDynamicAffinityRules{}
			if cValue := v.Get("action"); cValue.Exists() {
				item.Action = types.StringValue(cValue.String())
			} else {
				item.Action = types.StringNull()
			}
			if cValue := v.Get("color"); cValue.Exists() && len(cValue.Array()) > 0 {
				item.Color = helpers.GetStringSet(cValue.Array())
			} else {
				item.Color = types.SetNull(types.StringType)
			}
			data.DynamicAffinityRules = append(data.DynamicAffinityRules, item)
			return true
		})
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.segments.sid-algorithm"); value.Exists() {
		data.DynamicSegmentsSidAlgorithm = types.Int64Value(value.Int())
	} else {
		data.DynamicSegmentsSidAlgorithm = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.disjoint-path.type"); value.Exists() {
		data.DynamicDisjointPathType = types.StringValue(value.String())
	} else {
		data.DynamicDisjointPathType = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.disjoint-path.group-id"); value.Exists() {
		data.DynamicDisjointPathGroupId = types.Int64Value(value.Int())
	} else {
		data.DynamicDisjointPathGroupId = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.disjoint-path.sub-id"); value.Exists() {
		data.DynamicDisjointPathSubId = types.Int64Value(value.Int())
	} else {
		data.DynamicDisjointPathSubId = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.delay-measurement.profile"); value.Exists() {
		data.DelayMeasurementProfile = types.StringValue(value.String())
	} else {
		data.DelayMeasurementProfile = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.delay-measurement.logging.delay-exceeded"); value.Exists() {
		data.DelayMeasurementLoggingDelayExceeded = types.BoolValue(true)
	} else {
		data.DelayMeasurementLoggingDelayExceeded = types.BoolNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.liveness-detection.profile"); value.Exists() {
		data.LivenessDetectionProfile = types.StringValue(value.String())
	} else {
		data.LivenessDetectionProfile = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.liveness-detection.backup"); value.Exists() {
		data.LivenessDetectionBackup = types.StringValue(value.String())
	} else {
		data.LivenessDetectionBackup = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.liveness-detection.invalidation-action"); value.Exists() {
		data.LivenessDetectionInvalidationAction = types.StringValue(value.String())
	} else {
		data.LivenessDetectionInvalidationAction = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.liveness-detection.logging.session-state-change"); value.Exists() {
		data.LivenessDetectionLoggingSessionStateChange = types.BoolValue(true)
	} else {
		data.LivenessDetectionLoggingSessionStateChange = types.BoolNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.reverse-path-label"); value.Exists() {
		data.PerformanceMeasurementReversePathLabel = types.Int64Value(value.Int())
	} else {
		data.PerformanceMeasurementReversePathLabel = types.Int64Null()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *SRTEODNTemplates) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.color"); value.Exists() && !data.Color.IsNull() {
		data.Color = types.Int64Value(value.Int())
	} else {
		data.Color = types.Int64Null()
	}
	for i := range data.HeadEnd {
		keys := [...]string{"name"}
		keyValues := [...]string{data.HeadEnd[i].Name.ValueString()}

		var r gjson.Result
		res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.head-end").ForEach(
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
		if value := r.Get("name"); value.Exists() && !data.HeadEnd[i].Name.IsNull() {
			data.HeadEnd[i].Name = types.StringValue(value.String())
		} else {
			data.HeadEnd[i].Name = types.StringNull()
		}
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.maximum-sid-depth"); value.Exists() && !data.MaximumSidDepth.IsNull() {
		data.MaximumSidDepth = types.Int64Value(value.Int())
	} else {
		data.MaximumSidDepth = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.pce-group"); value.Exists() && !data.PceGroup.IsNull() {
		data.PceGroup = types.StringValue(value.String())
	} else {
		data.PceGroup = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.bandwidth"); value.Exists() && !data.Bandwidth.IsNull() {
		data.Bandwidth = types.Int64Value(value.Int())
	} else {
		data.Bandwidth = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.source-address"); value.Exists() && !data.SourceAddress.IsNull() {
		data.SourceAddress = types.StringValue(value.String())
	} else {
		data.SourceAddress = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.srv6.locator.locator-name"); value.Exists() && !data.Srv6LocatorName.IsNull() {
		data.Srv6LocatorName = types.StringValue(value.String())
	} else {
		data.Srv6LocatorName = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.srv6.locator.behavior"); value.Exists() && !data.Srv6LocatorBehavior.IsNull() {
		data.Srv6LocatorBehavior = types.StringValue(value.String())
	} else {
		data.Srv6LocatorBehavior = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.srv6.locator.binding-sid-type"); value.Exists() && !data.Srv6LocatorBindingSidType.IsNull() {
		data.Srv6LocatorBindingSidType = types.StringValue(value.String())
	} else {
		data.Srv6LocatorBindingSidType = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.pce"); value.Exists() && !data.DynamicPce.IsNull() {
		data.DynamicPce = types.BoolValue(true)
	} else {
		data.DynamicPce = types.BoolNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.flex-alg"); value.Exists() && !data.DynamicFlexAlg.IsNull() {
		data.DynamicFlexAlg = types.Int64Value(value.Int())
	} else {
		data.DynamicFlexAlg = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.metric-type"); value.Exists() && !data.DynamicMetricType.IsNull() {
		data.DynamicMetricType = types.StringValue(value.String())
	} else {
		data.DynamicMetricType = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.metric-margin.absolute"); value.Exists() && !data.DynamicMetricMarginAbsolute.IsNull() {
		data.DynamicMetricMarginAbsolute = types.Int64Value(value.Int())
	} else {
		data.DynamicMetricMarginAbsolute = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.metric-margin.relative"); value.Exists() && !data.DynamicMetricMarginRelative.IsNull() {
		data.DynamicMetricMarginRelative = types.Int64Value(value.Int())
	} else {
		data.DynamicMetricMarginRelative = types.Int64Null()
	}
	for i := range data.DynamicAffinityRules {
		keys := [...]string{"action"}
		keyValues := [...]string{data.DynamicAffinityRules[i].Action.ValueString()}

		var r gjson.Result
		res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.affinity.rule").ForEach(
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
		if value := r.Get("action"); value.Exists() && !data.DynamicAffinityRules[i].Action.IsNull() {
			data.DynamicAffinityRules[i].Action = types.StringValue(value.String())
		} else {
			data.DynamicAffinityRules[i].Action = types.StringNull()
		}
		if value := r.Get("color"); value.Exists() && !data.DynamicAffinityRules[i].Color.IsNull() {
			data.DynamicAffinityRules[i].Color = helpers.GetStringSet(value.Array())
		} else {
			data.DynamicAffinityRules[i].Color = types.SetNull(types.StringType)
		}
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.segments.sid-algorithm"); value.Exists() && !data.DynamicSegmentsSidAlgorithm.IsNull() {
		data.DynamicSegmentsSidAlgorithm = types.Int64Value(value.Int())
	} else {
		data.DynamicSegmentsSidAlgorithm = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.disjoint-path.type"); value.Exists() && !data.DynamicDisjointPathType.IsNull() {
		data.DynamicDisjointPathType = types.StringValue(value.String())
	} else {
		data.DynamicDisjointPathType = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.disjoint-path.group-id"); value.Exists() && !data.DynamicDisjointPathGroupId.IsNull() {
		data.DynamicDisjointPathGroupId = types.Int64Value(value.Int())
	} else {
		data.DynamicDisjointPathGroupId = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.dynamic.disjoint-path.sub-id"); value.Exists() && !data.DynamicDisjointPathSubId.IsNull() {
		data.DynamicDisjointPathSubId = types.Int64Value(value.Int())
	} else {
		data.DynamicDisjointPathSubId = types.Int64Null()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.delay-measurement.profile"); value.Exists() && !data.DelayMeasurementProfile.IsNull() {
		data.DelayMeasurementProfile = types.StringValue(value.String())
	} else {
		data.DelayMeasurementProfile = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.delay-measurement.logging.delay-exceeded"); value.Exists() && !data.DelayMeasurementLoggingDelayExceeded.IsNull() {
		data.DelayMeasurementLoggingDelayExceeded = types.BoolValue(true)
	} else {
		data.DelayMeasurementLoggingDelayExceeded = types.BoolNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.liveness-detection.profile"); value.Exists() && !data.LivenessDetectionProfile.IsNull() {
		data.LivenessDetectionProfile = types.StringValue(value.String())
	} else {
		data.LivenessDetectionProfile = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.liveness-detection.backup"); value.Exists() && !data.LivenessDetectionBackup.IsNull() {
		data.LivenessDetectionBackup = types.StringValue(value.String())
	} else {
		data.LivenessDetectionBackup = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.liveness-detection.invalidation-action"); value.Exists() && !data.LivenessDetectionInvalidationAction.IsNull() {
		data.LivenessDetectionInvalidationAction = types.StringValue(value.String())
	} else {
		data.LivenessDetectionInvalidationAction = types.StringNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.liveness-detection.logging.session-state-change"); value.Exists() && !data.LivenessDetectionLoggingSessionStateChange.IsNull() {
		data.LivenessDetectionLoggingSessionStateChange = types.BoolValue(true)
	} else {
		data.LivenessDetectionLoggingSessionStateChange = types.BoolNull()
	}
	if value := res.Get("cisco-sr-te-cfp-sr-odn:odn-template.0.performance-measurement.reverse-path-label"); value.Exists() && !data.PerformanceMeasurementReversePathLabel.IsNull() {
		data.PerformanceMeasurementReversePathLabel = types.Int64Value(value.Int())
	} else {
		data.PerformanceMeasurementReversePathLabel = types.Int64Null()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *SRTEODNTemplates) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Color.IsNull() {
		return false
	}
	if len(data.HeadEnd) > 0 {
		return false
	}
	if !data.MaximumSidDepth.IsNull() {
		return false
	}
	if !data.PceGroup.IsNull() {
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
	if !data.DynamicPce.IsNull() {
		return false
	}
	if !data.DynamicFlexAlg.IsNull() {
		return false
	}
	if !data.DynamicMetricType.IsNull() {
		return false
	}
	if !data.DynamicMetricMarginAbsolute.IsNull() {
		return false
	}
	if !data.DynamicMetricMarginRelative.IsNull() {
		return false
	}
	if len(data.DynamicAffinityRules) > 0 {
		return false
	}
	if !data.DynamicSegmentsSidAlgorithm.IsNull() {
		return false
	}
	if !data.DynamicDisjointPathType.IsNull() {
		return false
	}
	if !data.DynamicDisjointPathGroupId.IsNull() {
		return false
	}
	if !data.DynamicDisjointPathSubId.IsNull() {
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
	return true
}

// End of section. //template:end isNull

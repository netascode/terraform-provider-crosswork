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

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type L2VPNY1731Profiles struct {
	Instance                 types.String                        `tfsdk:"instance"`
	Id                       types.String                        `tfsdk:"id"`
	Name                     types.String                        `tfsdk:"name"`
	Type                     types.String                        `tfsdk:"type"`
	ScheduleInterval         types.Int64                         `tfsdk:"schedule_interval"`
	ScheduleDuration         types.String                        `tfsdk:"schedule_duration"`
	ProbeType                types.String                        `tfsdk:"probe_type"`
	ProbeMessageCount        types.Int64                         `tfsdk:"probe_message_count"`
	ProbeMessagePeriod       types.Int64                         `tfsdk:"probe_message_period"`
	ProbeMeasurementInterval types.String                        `tfsdk:"probe_measurement_interval"`
	ProbeFrameSize           types.Int64                         `tfsdk:"probe_frame_size"`
	ProbePriority            types.Int64                         `tfsdk:"probe_priority"`
	DelayStatistics          []L2VPNY1731ProfilesDelayStatistics `tfsdk:"delay_statistics"`
	LossStatistics           []L2VPNY1731ProfilesLossStatistics  `tfsdk:"loss_statistics"`
	BucketSize               types.Int64                         `tfsdk:"bucket_size"`
	BucketArchive            types.Int64                         `tfsdk:"bucket_archive"`
}

type L2VPNY1731ProfilesDelayStatistics struct {
	Type types.String `tfsdk:"type"`
}

type L2VPNY1731ProfilesLossStatistics struct {
	Type types.String `tfsdk:"type"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data L2VPNY1731Profiles) getPath() string {
	return fmt.Sprintf("/crosswork/proxy/nso/restconf/data/ietf-l2vpn-ntw:l2vpn-ntw/cisco-l2vpn-ntw:y-1731-profile=%v", url.PathEscape(data.Name.ValueString()))
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

func (data L2VPNY1731Profiles) getPathPost() string {
	return "/crosswork/proxy/nso/restconf/data/ietf-l2vpn-ntw:l2vpn-ntw"
}

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
func (data L2VPNY1731Profiles) toBody(ctx context.Context, state L2VPNY1731Profiles) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.Name.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.name", data.Name.ValueString())
	}
	if !data.Type.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.type", data.Type.ValueString())
	}
	if !data.ScheduleInterval.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.schedule.interval", data.ScheduleInterval.ValueInt64())
	}
	if !data.ScheduleDuration.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.schedule.duration", data.ScheduleDuration.ValueString())
	}
	if !data.ProbeType.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.probe.type", data.ProbeType.ValueString())
	}
	if !data.ProbeMessageCount.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.probe.burst.message-count", data.ProbeMessageCount.ValueInt64())
	}
	if !data.ProbeMessagePeriod.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.probe.burst.message-period", data.ProbeMessagePeriod.ValueInt64())
	}
	if !data.ProbeMeasurementInterval.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.probe.measurement-interval", data.ProbeMeasurementInterval.ValueString())
	}
	if !data.ProbeFrameSize.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.probe.frame-size", data.ProbeFrameSize.ValueInt64())
	}
	if !data.ProbePriority.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.probe.priority", data.ProbePriority.ValueInt64())
	}
	if len(data.DelayStatistics) > 0 {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.delay-params.statistic", []interface{}{})
		for _, item := range data.DelayStatistics {
			itemBody := ""
			if !item.Type.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "type", item.Type.ValueString())
			}
			body, _ = sjson.SetRaw(body, "cisco-l2vpn-ntw:y-1731-profile.0.delay-params.statistic.-1", itemBody)
		}
	}
	if len(data.LossStatistics) > 0 {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.loss-params.statistic", []interface{}{})
		for _, item := range data.LossStatistics {
			itemBody := ""
			if !item.Type.IsNull() {
				itemBody, _ = sjson.Set(itemBody, "type", item.Type.ValueString())
			}
			body, _ = sjson.SetRaw(body, "cisco-l2vpn-ntw:y-1731-profile.0.loss-params.statistic.-1", itemBody)
		}
	}
	if !data.BucketSize.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.bucket-details.bucket-size", data.BucketSize.ValueInt64())
	}
	if !data.BucketArchive.IsNull() {
		body, _ = sjson.Set(body, "cisco-l2vpn-ntw:y-1731-profile.0.bucket-details.bucket-archive", data.BucketArchive.ValueInt64())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *L2VPNY1731Profiles) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.type"); value.Exists() {
		data.Type = types.StringValue(value.String())
	} else {
		data.Type = types.StringNull()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.schedule.interval"); value.Exists() {
		data.ScheduleInterval = types.Int64Value(value.Int())
	} else {
		data.ScheduleInterval = types.Int64Null()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.schedule.duration"); value.Exists() {
		data.ScheduleDuration = types.StringValue(value.String())
	} else {
		data.ScheduleDuration = types.StringNull()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.probe.type"); value.Exists() {
		data.ProbeType = types.StringValue(value.String())
	} else {
		data.ProbeType = types.StringNull()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.probe.burst.message-count"); value.Exists() {
		data.ProbeMessageCount = types.Int64Value(value.Int())
	} else {
		data.ProbeMessageCount = types.Int64Null()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.probe.burst.message-period"); value.Exists() {
		data.ProbeMessagePeriod = types.Int64Value(value.Int())
	} else {
		data.ProbeMessagePeriod = types.Int64Null()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.probe.measurement-interval"); value.Exists() {
		data.ProbeMeasurementInterval = types.StringValue(value.String())
	} else {
		data.ProbeMeasurementInterval = types.StringNull()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.probe.frame-size"); value.Exists() {
		data.ProbeFrameSize = types.Int64Value(value.Int())
	} else {
		data.ProbeFrameSize = types.Int64Null()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.probe.priority"); value.Exists() {
		data.ProbePriority = types.Int64Value(value.Int())
	} else {
		data.ProbePriority = types.Int64Null()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.delay-params.statistic"); value.Exists() && len(value.Array()) > 0 {
		data.DelayStatistics = make([]L2VPNY1731ProfilesDelayStatistics, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := L2VPNY1731ProfilesDelayStatistics{}
			if cValue := v.Get("type"); cValue.Exists() {
				item.Type = types.StringValue(cValue.String())
			} else {
				item.Type = types.StringNull()
			}
			data.DelayStatistics = append(data.DelayStatistics, item)
			return true
		})
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.loss-params.statistic"); value.Exists() && len(value.Array()) > 0 {
		data.LossStatistics = make([]L2VPNY1731ProfilesLossStatistics, 0)
		value.ForEach(func(k, v gjson.Result) bool {
			item := L2VPNY1731ProfilesLossStatistics{}
			if cValue := v.Get("type"); cValue.Exists() {
				item.Type = types.StringValue(cValue.String())
			} else {
				item.Type = types.StringNull()
			}
			data.LossStatistics = append(data.LossStatistics, item)
			return true
		})
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.bucket-details.bucket-size"); value.Exists() {
		data.BucketSize = types.Int64Value(value.Int())
	} else {
		data.BucketSize = types.Int64Null()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.bucket-details.bucket-archive"); value.Exists() {
		data.BucketArchive = types.Int64Value(value.Int())
	} else {
		data.BucketArchive = types.Int64Null()
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *L2VPNY1731Profiles) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.type"); value.Exists() && !data.Type.IsNull() {
		data.Type = types.StringValue(value.String())
	} else {
		data.Type = types.StringNull()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.schedule.interval"); value.Exists() && !data.ScheduleInterval.IsNull() {
		data.ScheduleInterval = types.Int64Value(value.Int())
	} else {
		data.ScheduleInterval = types.Int64Null()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.schedule.duration"); value.Exists() && !data.ScheduleDuration.IsNull() {
		data.ScheduleDuration = types.StringValue(value.String())
	} else {
		data.ScheduleDuration = types.StringNull()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.probe.type"); value.Exists() && !data.ProbeType.IsNull() {
		data.ProbeType = types.StringValue(value.String())
	} else {
		data.ProbeType = types.StringNull()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.probe.burst.message-count"); value.Exists() && !data.ProbeMessageCount.IsNull() {
		data.ProbeMessageCount = types.Int64Value(value.Int())
	} else {
		data.ProbeMessageCount = types.Int64Null()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.probe.burst.message-period"); value.Exists() && !data.ProbeMessagePeriod.IsNull() {
		data.ProbeMessagePeriod = types.Int64Value(value.Int())
	} else {
		data.ProbeMessagePeriod = types.Int64Null()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.probe.measurement-interval"); value.Exists() && !data.ProbeMeasurementInterval.IsNull() {
		data.ProbeMeasurementInterval = types.StringValue(value.String())
	} else {
		data.ProbeMeasurementInterval = types.StringNull()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.probe.frame-size"); value.Exists() && !data.ProbeFrameSize.IsNull() {
		data.ProbeFrameSize = types.Int64Value(value.Int())
	} else {
		data.ProbeFrameSize = types.Int64Null()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.probe.priority"); value.Exists() && !data.ProbePriority.IsNull() {
		data.ProbePriority = types.Int64Value(value.Int())
	} else {
		data.ProbePriority = types.Int64Null()
	}
	for i := range data.DelayStatistics {
		keys := [...]string{"type"}
		keyValues := [...]string{data.DelayStatistics[i].Type.ValueString()}

		var r gjson.Result
		res.Get("cisco-l2vpn-ntw:y-1731-profile.0.delay-params.statistic").ForEach(
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
		if value := r.Get("type"); value.Exists() && !data.DelayStatistics[i].Type.IsNull() {
			data.DelayStatistics[i].Type = types.StringValue(value.String())
		} else {
			data.DelayStatistics[i].Type = types.StringNull()
		}
	}
	for i := range data.LossStatistics {
		keys := [...]string{"type"}
		keyValues := [...]string{data.LossStatistics[i].Type.ValueString()}

		var r gjson.Result
		res.Get("cisco-l2vpn-ntw:y-1731-profile.0.loss-params.statistic").ForEach(
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
		if value := r.Get("type"); value.Exists() && !data.LossStatistics[i].Type.IsNull() {
			data.LossStatistics[i].Type = types.StringValue(value.String())
		} else {
			data.LossStatistics[i].Type = types.StringNull()
		}
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.bucket-details.bucket-size"); value.Exists() && !data.BucketSize.IsNull() {
		data.BucketSize = types.Int64Value(value.Int())
	} else {
		data.BucketSize = types.Int64Null()
	}
	if value := res.Get("cisco-l2vpn-ntw:y-1731-profile.0.bucket-details.bucket-archive"); value.Exists() && !data.BucketArchive.IsNull() {
		data.BucketArchive = types.Int64Value(value.Int())
	} else {
		data.BucketArchive = types.Int64Null()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *L2VPNY1731Profiles) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Type.IsNull() {
		return false
	}
	if !data.ScheduleInterval.IsNull() {
		return false
	}
	if !data.ScheduleDuration.IsNull() {
		return false
	}
	if !data.ProbeType.IsNull() {
		return false
	}
	if !data.ProbeMessageCount.IsNull() {
		return false
	}
	if !data.ProbeMessagePeriod.IsNull() {
		return false
	}
	if !data.ProbeMeasurementInterval.IsNull() {
		return false
	}
	if !data.ProbeFrameSize.IsNull() {
		return false
	}
	if !data.ProbePriority.IsNull() {
		return false
	}
	if len(data.DelayStatistics) > 0 {
		return false
	}
	if len(data.LossStatistics) > 0 {
		return false
	}
	if !data.BucketSize.IsNull() {
		return false
	}
	if !data.BucketArchive.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

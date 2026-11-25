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
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type DeviceAlarmManager struct {
	Instance types.String `tfsdk:"instance"`
	Id       types.String `tfsdk:"id"`
	Models   types.Map    `tfsdk:"models"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data DeviceAlarmManager) getPath() string {
	return "/crosswork/alarm/v1/manager/settings"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data DeviceAlarmManager) getPathGet() string {
	return "/crosswork/alarm/v1/manager/settings"
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

// Written by hand: a write reads the device models of the controller first (deviceAlarmManagerWriteBody), so the model has no body of its own.
func (data DeviceAlarmManager) toBody(ctx context.Context, state DeviceAlarmManager) string {
	return "{}"
}

// End of section. //template:end toBody

// Written by hand: an import and the data source read every device model of the controller, without the prefix of the key.
func (data *DeviceAlarmManager) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("device_alarm_manager")
	if models := deviceAlarmManagerModels(res); len(models) > 0 {
		data.Models = deviceAlarmManagerMap(models)
	} else {
		data.Models = types.MapNull(types.BoolType)
	}
}

// End of section. //template:end fromBody

// Written by hand: only the models of the configuration are read back, so a model that is not listed never shows as a difference.
func (data *DeviceAlarmManager) updateFromBody(ctx context.Context, res gjson.Result) {
	if data.Models.IsNull() {
		return
	}
	known := deviceAlarmManagerModels(res)
	kept := map[string]bool{}
	for name := range deviceAlarmManagerPlanned(ctx, data.Models) {
		if enabled, ok := known[name]; ok {
			kept[name] = enabled
		}
	}
	data.Models = deviceAlarmManagerMap(kept)
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *DeviceAlarmManager) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Models.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

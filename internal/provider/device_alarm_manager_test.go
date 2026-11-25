// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License Version 2.0 (the "License"),
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

// The page as the controller answered when probed: a flat object, one key for each device model, with the prefix.
const deviceAlarmManagerAnswer = `{"alarmManager/Cisco ASR 9010 Router":false,"alarmManager/Cisco ASR 9006 Router":true,"alarmManager/Cisco NCS 5501":false}`

func TestDeviceAlarmManagerModelsDropThePrefix(t *testing.T) {
	models := deviceAlarmManagerModels(gjson.Parse(deviceAlarmManagerAnswer))
	assert.Equal(t, 3, len(models))
	assert.True(t, models["Cisco ASR 9006 Router"], "the state is read as sent")
	assert.False(t, models["Cisco ASR 9010 Router"])
	assert.Equal(t, 0, len(deviceAlarmManagerModels(gjson.Parse(`{"other":true,"alarmManager/x":"yes"}`))), "a key without the prefix and a value that is no bool are not models")
}

func TestDeviceAlarmManagerWriteBodyAddsThePrefixAndLeavesTheRestAlone(t *testing.T) {
	known := deviceAlarmManagerModels(gjson.Parse(deviceAlarmManagerAnswer))
	body, sent, err := deviceAlarmManagerWriteBody(known, map[string]bool{"Cisco ASR 9010 Router": true}, nil)
	assert.True(t, err == nil, err)
	assert.Equal(t, `{"alarmManager/Cisco ASR 9010 Router":true}`, body, "only the planned model is sent, with the prefix")
	assert.Equal(t, 1, len(sent))
}

func TestDeviceAlarmManagerWriteBodySwitchesOffWhatLeftTheConfiguration(t *testing.T) {
	known := deviceAlarmManagerModels(gjson.Parse(deviceAlarmManagerAnswer))
	body, _, err := deviceAlarmManagerWriteBody(known, map[string]bool{"Cisco ASR 9010 Router": true}, map[string]bool{"Cisco ASR 9010 Router": false, "Cisco ASR 9006 Router": true, "Cisco gone": true})
	assert.True(t, err == nil, err)
	assert.Equal(t, `{"alarmManager/Cisco ASR 9006 Router":false,"alarmManager/Cisco ASR 9010 Router":true}`, body, "a model that left the plan is set to false; one the controller dropped is not sent")
}

func TestDeviceAlarmManagerWriteBodyRefusesAnUnknownName(t *testing.T) {
	known := deviceAlarmManagerModels(gjson.Parse(deviceAlarmManagerAnswer))
	for _, name := range []string{"Cisco asr 9010 Router", "Cisco ASR 9010 Router ", "alarmManager/Cisco ASR 9010 Router", "nope"} {
		_, _, err := deviceAlarmManagerWriteBody(known, map[string]bool{name: true}, nil)
		assert.True(t, err != nil, "the name "+name+" is refused before anything is written")
		assert.Contains(t, err.Error(), "the controller has no device model")
	}
	body, _, err := deviceAlarmManagerWriteBody(known, nil, nil)
	assert.True(t, err == nil && body == "", "nothing planned and nothing before sends nothing")
}

func TestDeviceAlarmManagerCheckAnswer(t *testing.T) {
	sent := map[string]bool{"alarmManager/Cisco ASR 9010 Router": true}
	assert.True(t, deviceAlarmManagerCheckAnswer(sent, gjson.Parse(`{"alarmManager/Cisco ASR 9010 Router":true}`)) == nil)
	err := deviceAlarmManagerCheckAnswer(sent, gjson.Parse(`{"alarmManager/Cisco ASR 9010 Router":false}`))
	assert.True(t, err != nil, "a model that came back false after true was ignored")
	assert.Contains(t, err.Error(), "Cisco ASR 9010 Router")
	assert.True(t, deviceAlarmManagerCheckAnswer(sent, gjson.Parse(`{}`)) != nil, "a model missing from the answer is an error")
}

func TestDeviceAlarmManagerReadsTheModelsOfTheConfigurationAndAllOnImport(t *testing.T) {
	ctx := context.Background()
	state := DeviceAlarmManager{Models: deviceAlarmManagerMap(map[string]bool{"Cisco ASR 9010 Router": true, "Cisco gone": true})}
	state.updateFromBody(ctx, gjson.Parse(deviceAlarmManagerAnswer))
	got := deviceAlarmManagerPlanned(ctx, state.Models)
	assert.Equal(t, 1, len(got), "only the models of the state are read back, and a model the controller dropped is removed")
	assert.False(t, got["Cisco ASR 9010 Router"], "the value of the controller wins")

	var data DeviceAlarmManager
	data.fromBody(ctx, gjson.Parse(deviceAlarmManagerAnswer))
	assert.Equal(t, 3, len(data.Models.Elements()), "an import and the data source read every model")
	assert.Equal(t, "device_alarm_manager", data.Id.ValueString())

	empty := DeviceAlarmManager{Models: types.MapNull(types.BoolType)}
	empty.updateFromBody(ctx, gjson.Parse(deviceAlarmManagerAnswer))
	assert.True(t, empty.Models.IsNull(), "a null stays null")
}

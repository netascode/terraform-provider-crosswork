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
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"

	cw "github.com/netascode/go-crosswork"
)

const deviceAlarmManagerPath = "/crosswork/alarm/v1/manager/settings"

// The device alarm manager is one flat object with a key "alarmManager/<device model>" for each device model of the controller. As
// probed: a write merges the keys it sends, a key that the controller does not know (a misspelled model, a model without the
// prefix, another case) answers 200 and is ignored, a value that is not a bool is ignored, and the answer of a write holds only the
// keys that were sent, with the value the controller holds now. So a write reads the models first, refuses a name that is not one,
// and compares the answer with what it sent.

const deviceAlarmManagerPrefix = "alarmManager/"

// deviceAlarmManagerModels reads the answer: the device model (without the prefix) and its state, in the order of the controller.
func deviceAlarmManagerModels(res gjson.Result) map[string]bool {
	models := map[string]bool{}
	res.ForEach(func(key, value gjson.Result) bool {
		if name, ok := strings.CutPrefix(key.String(), deviceAlarmManagerPrefix); ok && (value.Type == gjson.True || value.Type == gjson.False) {
			models[name] = value.Bool()
		}
		return true
	})
	return models
}

func deviceAlarmManagerNames(models map[string]bool) []string {
	names := make([]string, 0, len(models))
	for name := range models {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// deviceAlarmManagerWriteBody builds the body of a write: the planned models with their state, and every model of the earlier state
// that the plan no longer names set to false. A planned model that the controller does not have is an error, and so is a model
// that the plan names with an unknown value.
func deviceAlarmManagerWriteBody(known, planned, previous map[string]bool) (string, map[string]bool, error) {
	var unknown []string
	for _, name := range deviceAlarmManagerNames(planned) {
		if _, ok := known[name]; !ok {
			unknown = append(unknown, fmt.Sprintf("%q", name))
		}
	}
	if len(unknown) > 0 {
		return "", nil, fmt.Errorf("the controller has no device model %s; use the name as the Device alarm manager table or the product type of a device shows it, spelled exactly", strings.Join(unknown, ", "))
	}
	sent := map[string]bool{}
	for name, enabled := range planned {
		sent[deviceAlarmManagerPrefix+name] = enabled
	}
	for name := range previous {
		if _, named := planned[name]; named {
			continue
		}
		if _, ok := known[name]; ok {
			sent[deviceAlarmManagerPrefix+name] = false
		}
	}
	if len(sent) == 0 {
		return "", sent, nil
	}
	body, _ := json.Marshal(sent)
	return string(body), sent, nil
}

// deviceAlarmManagerCheckAnswer compares the answer of a write with what was sent: the controller answers the value it holds for
// each key that was sent, and false for a key it ignored.
func deviceAlarmManagerCheckAnswer(sent map[string]bool, answer gjson.Result) error {
	var differ []string
	for _, key := range deviceAlarmManagerNames(sent) {
		got := answer.Get(escapeJSONKey(key))
		if !got.Exists() || got.Bool() != sent[key] {
			differ = append(differ, fmt.Sprintf("%q", strings.TrimPrefix(key, deviceAlarmManagerPrefix)))
		}
	}
	if len(differ) > 0 {
		return fmt.Errorf("the controller did not keep the setting of the device model %s", strings.Join(differ, ", "))
	}
	return nil
}

// deviceAlarmManagerPlanned returns the models of a map attribute, or none when it is null or unknown.
func deviceAlarmManagerPlanned(ctx context.Context, models types.Map) map[string]bool {
	planned := map[string]bool{}
	if models.IsNull() || models.IsUnknown() {
		return planned
	}
	values := map[string]types.Bool{}
	models.ElementsAs(ctx, &values, false)
	for name, value := range values {
		planned[name] = value.ValueBool()
	}
	return planned
}

func deviceAlarmManagerMap(models map[string]bool) types.Map {
	elements := map[string]attr.Value{}
	for name, enabled := range models {
		elements[name] = types.BoolValue(enabled)
	}
	value, _ := types.MapValue(types.BoolType, elements)
	return value
}

// applyDeviceAlarmManager writes the planned models and reads them back into the plan. previous holds the models of the state.
func applyDeviceAlarmManager(ctx context.Context, c *cw.Client, plan *DeviceAlarmManager, previous types.Map) error {
	current, err := c.Get(deviceAlarmManagerPath)
	if err != nil {
		return fmt.Errorf("%s, %s", err, current.String())
	}
	body, sent, err := deviceAlarmManagerWriteBody(deviceAlarmManagerModels(current), deviceAlarmManagerPlanned(ctx, plan.Models), deviceAlarmManagerPlanned(ctx, previous))
	if err != nil {
		return err
	}
	if body != "" {
		answer, err := c.Post(deviceAlarmManagerPath, body, cw.UseMutex)
		if err != nil {
			return fmt.Errorf("%s, %s", err, answer.String())
		}
		if err := deviceAlarmManagerCheckAnswer(sent, answer); err != nil {
			return err
		}
	}
	read, err := c.Get(deviceAlarmManagerPath)
	if err != nil {
		return fmt.Errorf("%s, %s", err, read.String())
	}
	plan.updateFromBody(ctx, read)
	return nil
}

// resetDeviceAlarmManager sets every model of the controller to false, the state it has by default.
func resetDeviceAlarmManager(c *cw.Client) error {
	current, err := c.Get(deviceAlarmManagerPath)
	if err != nil {
		return fmt.Errorf("%s, %s", err, current.String())
	}
	sent := map[string]bool{}
	for name := range deviceAlarmManagerModels(current) {
		sent[deviceAlarmManagerPrefix+name] = false
	}
	if len(sent) == 0 {
		return nil
	}
	body, _ := json.Marshal(sent)
	if res, err := c.Post(deviceAlarmManagerPath, string(body), cw.UseMutex); err != nil {
		return fmt.Errorf("%s, %s", err, res.String())
	}
	return nil
}

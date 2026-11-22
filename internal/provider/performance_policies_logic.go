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
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cw "github.com/netascode/go-crosswork"
)

const performanceMatchCos = "matchCos"

// The controller keeps a comma separated target list as it was sent (no trimming, no sorting, no duplicates removed), so the
// provider writes the members sorted and reads them back split, trimmed and without empty members.

func setStrings(ctx context.Context, set types.Set) []string {
	var values []string
	if !set.IsNull() && !set.IsUnknown() {
		set.ElementsAs(ctx, &values, false)
	}
	sort.Strings(values)
	return values
}

func splitTargets(list string) []string {
	var values []string
	seen := map[string]bool{}
	for _, part := range strings.Split(list, ",") {
		part = strings.TrimSpace(part)
		if part != "" && !seen[part] {
			seen[part] = true
			values = append(values, part)
		}
	}
	return values
}

func targetSet(values []string) types.Set {
	if len(values) == 0 {
		return types.SetNull(types.StringType)
	}
	elements, _ := types.SetValueFrom(context.Background(), types.StringType, values)
	return elements
}

// The contact is the only part of the tag that the provider owns: the UI writes {"contact":"<text>"} or nothing.
func tagFromContact(contact types.String) string {
	if contact.IsNull() || contact.ValueString() == "" {
		return ""
	}
	tag, _ := sjson.Set("", "contact", contact.ValueString())
	return tag
}

func contactFromTag(tag string) types.String {
	if strings.HasPrefix(strings.TrimSpace(tag), "{") && gjson.Valid(tag) {
		if contact := gjson.Get(tag, "contact"); contact.Type == gjson.String && contact.String() != "" {
			return types.StringValue(contact.String())
		}
	}
	return types.StringNull()
}

func conditionBody(ctx context.Context, id int, condition PerformancePoliciesThresholdsConditions) string {
	item, _ := sjson.Set("", "conditionID", id)
	item, _ = sjson.Set(item, "relation", condition.Relation.ValueString())
	item, _ = sjson.Set(item, "value", condition.Value.ValueString())
	item, _ = sjson.Set(item, "count", condition.Count.ValueInt64())
	item, _ = sjson.Set(item, "severity", condition.Severity.ValueString())
	if !condition.MatchCos.IsNull() && len(condition.MatchCos.Elements()) > 0 {
		var classes []int64
		condition.MatchCos.ElementsAs(ctx, &classes, false)
		sort.Slice(classes, func(i, j int) bool { return classes[i] < classes[j] })
		values := make([]string, 0, len(classes))
		for _, class := range classes {
			values = append(values, strconv.FormatInt(class, 10))
		}
		field, _ := sjson.Set("", "name", performanceMatchCos)
		field, _ = sjson.Set(field, "values", values)
		item, _ = sjson.SetRaw(item, "extraConditionFields", "["+field+"]")
	}
	return item
}

func conditionFromBody(condition gjson.Result) PerformancePoliciesThresholdsConditions {
	item := PerformancePoliciesThresholdsConditions{
		Relation: types.StringValue(condition.Get("relation").String()),
		Value:    types.StringValue(condition.Get("value").String()),
		Count:    types.Int64Value(condition.Get("count").Int()),
		Severity: types.StringValue(condition.Get("severity").String()),
		MatchCos: types.SetNull(types.Int64Type),
	}
	for _, field := range condition.Get("extraConditionFields").Array() {
		if field.Get("name").String() == performanceMatchCos {
			if values := field.Get("values").Array(); len(values) > 0 {
				classes := make([]int64, 0, len(values))
				for _, value := range values {
					classes = append(classes, value.Int())
				}
				item.MatchCos, _ = types.SetValueFrom(context.Background(), types.Int64Type, classes)
			}
		}
	}
	return item
}

// The controller names a device by uuid and the module by inventory id (the name of the device), as for the device group
// assignments, so the devices of a plan are turned into uuids before a write and back into names after a read. A uuid that has
// no device any more is dropped: the next plan then shows the difference and the apply names the missing device.
func (data *PerformancePolicies) devicesToUUIDs(ctx context.Context, c *cw.Client) error {
	names := setStrings(ctx, data.Devices)
	if len(names) == 0 {
		return nil
	}
	byName, err := deviceUUIDs(c, names)
	if err != nil {
		return err
	}
	uuids := make([]string, 0, len(names))
	for _, name := range names {
		uuids = append(uuids, byName[name])
	}
	data.Devices = targetSet(uuids)
	return nil
}

func (data *PerformancePolicies) devicesToNames(ctx context.Context, c *cw.Client) error {
	uuids := setStrings(ctx, data.Devices)
	if len(uuids) == 0 {
		return nil
	}
	byUUID, err := inventoryIDs(c, uuids)
	if err != nil {
		return err
	}
	names := make([]string, 0, len(uuids))
	for _, uuid := range uuids {
		if name := byUUID[uuid]; name != "" {
			names = append(names, name)
		}
	}
	data.Devices = targetSet(names)
	return nil
}

// A policy is known by its name to the user, and by a number to the controller.
func policyIDByName(c *cw.Client, path, name string) (string, error) {
	res, err := c.Get(path, cw.ReadOnly)
	if err != nil {
		return "", fmt.Errorf("listing the performance policies: %w", err)
	}
	for _, item := range res.Array() {
		if item.Get("monitoringPolicy.name").String() == name {
			return item.Get("monitoringPolicy.id").String(), nil
		}
	}
	return "", fmt.Errorf("the performance policy %q does not exist", name)
}

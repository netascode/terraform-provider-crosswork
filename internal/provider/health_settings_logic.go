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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cw "github.com/netascode/go-crosswork"
)

const (
	healthSettingsPath      = "/crosswork/performance/v1/dashboards/healthsettings"
	healthSettingsAllPath   = "/crosswork/performance/v1/dashboards/healthsettings/all"
	healthSettingsResetPath = "/crosswork/performance/v1/dashboards/healthsettings/reset"
)

// The health settings are one object by policy and metric key. The controller accepts a row only as it read it (a flat map, a list
// or a bare row is refused), so a write takes the row from the controller, replaces its categories and its unit, and sends it back.

func healthCategoriesBody(categories []HealthSettingsRowsCategories) string {
	body := "[]"
	for _, category := range categories {
		item, _ := sjson.Set("", "level", category.Level.ValueString())
		if !category.Min.IsNull() {
			item, _ = sjson.Set(item, "min", category.Min.ValueFloat64())
		}
		if !category.Max.IsNull() {
			item, _ = sjson.Set(item, "max", category.Max.ValueFloat64())
		}
		if !category.Values.IsNull() {
			var values []int64
			category.Values.ElementsAs(context.Background(), &values, false)
			item, _ = sjson.Set(item, "values", values)
		}
		body, _ = sjson.SetRaw(body, "-1", item)
	}
	return body
}

func healthCategoriesFromBody(categories gjson.Result) []HealthSettingsRowsCategories {
	var out []HealthSettingsRowsCategories
	categories.ForEach(func(_, category gjson.Result) bool {
		item := HealthSettingsRowsCategories{Level: types.StringValue(category.Get("level").String()), Min: types.Float64Null(), Max: types.Float64Null(), Values: types.ListNull(types.Int64Type)}
		if value := category.Get("min"); value.Exists() {
			item.Min = types.Float64Value(value.Float())
		}
		if value := category.Get("max"); value.Exists() {
			item.Max = types.Float64Value(value.Float())
		}
		if value := category.Get("values"); value.Exists() {
			elements := []attr.Value{}
			value.ForEach(func(_, v gjson.Result) bool {
				elements = append(elements, types.Int64Value(v.Int()))
				return true
			})
			item.Values, _ = types.ListValue(types.Int64Type, elements)
		}
		out = append(out, item)
		return true
	})
	return out
}

func healthRow(all gjson.Result, key string) (policy string, row gjson.Result, found bool) {
	all.ForEach(func(name, rows gjson.Result) bool {
		if value := rows.Get(escapeJSONKey(key)); value.IsObject() {
			policy, row, found = name.String(), value, true
			return false
		}
		return true
	})
	return
}

func healthWriteBody(all gjson.Result, rows map[string]HealthSettingsRows) (string, error) {
	body := "{}"
	for key, planned := range rows {
		policy, current, found := healthRow(all, key)
		if !found {
			return "", fmt.Errorf("the controller has no health setting %q", key)
		}
		if !current.Get("editable").Bool() {
			return "", fmt.Errorf("the health setting %q cannot be changed", key)
		}
		raw, _ := sjson.SetRaw(current.Raw, "categories", healthCategoriesBody(planned.Categories))
		if !planned.Unit.IsNull() && planned.Unit.ValueString() != "" {
			raw, _ = sjson.Set(raw, "unit", planned.Unit.ValueString())
		}
		body, _ = sjson.SetRaw(body, escapeJSONKey(policy)+"."+escapeJSONKey(key), raw)
	}
	return body, nil
}

// The unit is read only when the configuration names it.
func healthRowFromBody(row gjson.Result, withUnit bool) HealthSettingsRows {
	out := HealthSettingsRows{Unit: types.StringNull(), Categories: healthCategoriesFromBody(row.Get("categories"))}
	if withUnit {
		out.Unit = types.StringValue(row.Get("unit").String())
	}
	return out
}

func applyHealthSettings(c *cw.Client, plan *HealthSettings) error {
	all, err := c.Get(healthSettingsAllPath)
	if err != nil {
		return fmt.Errorf("%s, %s", err, all.String())
	}
	body, err := healthWriteBody(all, plan.Rows)
	if err != nil {
		return err
	}
	if res, err := c.Put(healthSettingsPath, body, cw.UseMutex); err != nil {
		return fmt.Errorf("%s, %s", err, res.String())
	}
	read, err := c.Get(healthSettingsAllPath)
	if err != nil {
		return fmt.Errorf("%s, %s", err, read.String())
	}
	plan.updateFromBody(context.Background(), read)
	return nil
}

func resetHealthSettings(c *cw.Client) error {
	if res, err := c.Post(healthSettingsResetPath, `{}`, cw.UseMutex); err != nil {
		return fmt.Errorf("%s, %s", err, res.String())
	}
	return nil
}

// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
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
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

type GoldenConfigVariable struct {
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	Value          types.String `tfsdk:"value"`
	ValueWo        types.String `tfsdk:"value_wo"`
	ValueWoVersion types.Int64  `tfsdk:"value_wo_version"`
}

var goldenConfigVariableAttrTypes = map[string]attr.Type{
	"name":             types.StringType,
	"type":             types.StringType,
	"value":            types.StringType,
	"value_wo":         types.StringType,
	"value_wo_version": types.Int64Type,
}

var goldenConfigVariableType = types.ObjectType{AttrTypes: goldenConfigVariableAttrTypes}

type GoldenConfigTemplates struct {
	Id            types.String `tfsdk:"id"`
	Instance      types.String `tfsdk:"instance"`
	Name          types.String `tfsdk:"name"`
	TemplateId    types.Int64  `tfsdk:"template_id"`
	ConfigId      types.Int64  `tfsdk:"config_id"`
	Description   types.String `tfsdk:"description"`
	Tags          types.Set    `tfsdk:"tags"`
	Mode          types.String `tfsdk:"mode"`
	ConfigType    types.String `tfsdk:"config_type"`
	Vendor        types.String `tfsdk:"vendor"`
	SoftwareType  types.String `tfsdk:"software_type"`
	ProductSeries types.Set    `tfsdk:"product_series"`
	Audit         types.Bool   `tfsdk:"audit"`
	Config        types.String `tfsdk:"config"`
	Variables     types.List   `tfsdk:"variables"`
	ForceDelete   types.Bool   `tfsdk:"force_delete"`
}

// The secret value comes from the configuration (cfg), the only place a write-only value exists; pass the same list for cfg when the variables
// come from a state.
func gcVarsOf(ctx context.Context, list, cfg types.List) ([]gcVar, diag.Diagnostics) {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return nil, diags
	}
	var items, cfgItems []GoldenConfigVariable
	diags.Append(list.ElementsAs(ctx, &items, false)...)
	diags.Append(cfg.ElementsAs(ctx, &cfgItems, false)...)
	if diags.HasError() {
		return nil, diags
	}
	out := make([]gcVar, 0, len(items))
	for i, it := range items {
		v := gcVar{Name: it.Name.ValueString(), Type: it.Type.ValueString(), Value: it.Value.ValueString()}
		if v.Type == "" {
			v.Type = "STRING"
		}
		if !it.ValueWoVersion.IsNull() {
			v.Secret, v.Version = true, it.ValueWoVersion.ValueInt64()
			v.Value = ""
			if i < len(cfgItems) {
				v.Value = cfgItems[i].ValueWo.ValueString()
			}
		}
		out = append(out, v)
	}
	return out, diags
}

// A secret keeps its version and no value.
func gcVarsList(ctx context.Context, vars []gcVar) (types.List, diag.Diagnostics) {
	if len(vars) == 0 {
		return types.ListNull(goldenConfigVariableType), nil
	}
	items := make([]GoldenConfigVariable, 0, len(vars))
	for _, v := range vars {
		it := GoldenConfigVariable{Name: types.StringValue(v.Name), Type: types.StringValue(v.Type), Value: types.StringNull(),
			ValueWo: types.StringNull(), ValueWoVersion: types.Int64Null()}
		if v.Secret {
			it.ValueWoVersion = types.Int64Value(v.Version)
		} else {
			it.Value = types.StringValue(v.Value)
		}
		items = append(items, it)
	}
	return types.ListValueFrom(ctx, goldenConfigVariableType, items)
}

func validateGcVariables(ctx context.Context, list types.List, where string) diag.Diagnostics {
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return diags
	}
	var items []GoldenConfigVariable
	diags.Append(list.ElementsAs(ctx, &items, false)...)
	seen := map[string]bool{}
	for _, it := range items {
		if it.Name.IsUnknown() || it.Value.IsUnknown() || it.ValueWo.IsUnknown() || it.Type.IsUnknown() {
			continue
		}
		name := it.Name.ValueString()
		if seen[name] {
			diags.AddError("Duplicate variable", where+": the variable `"+name+"` is listed twice.")
		}
		seen[name] = true
		switch {
		case !it.Value.IsNull() && !it.ValueWo.IsNull():
			diags.AddError("Invalid Attribute Combination", where+": variable `"+name+"` sets both `value` and `value_wo`.")
		case it.Value.IsNull() && it.ValueWo.IsNull():
			diags.AddError("Invalid Attribute Combination", where+": variable `"+name+"` needs one of `value` and `value_wo`.")
		}
		if !it.ValueWo.IsNull() && it.ValueWoVersion.IsNull() {
			diags.AddError("Missing Attribute", where+": variable `"+name+"` needs `value_wo_version` with `value_wo`.")
		}
		if it.Type.ValueString() == "JSON" {
			for _, v := range []types.String{it.Value, it.ValueWo} {
				if !v.IsNull() && !gcValidJSON(v.ValueString()) {
					diags.AddError("Invalid JSON", where+": variable `"+name+"` has the type `JSON` and needs a JSON document as its value.")
				}
			}
		}
	}
	return diags
}

func stringSet(ctx context.Context, values []string) (types.Set, diag.Diagnostics) {
	return types.SetValueFrom(ctx, types.StringType, values)
}

func stringsOf(ctx context.Context, set types.Set) ([]string, diag.Diagnostics) {
	var out []string
	if set.IsNull() || set.IsUnknown() {
		return nil, nil
	}
	d := set.ElementsAs(ctx, &out, false)
	sort.Strings(out)
	return out, d
}

func gcPick(r gjson.Result, paths ...string) gjson.Result {
	for _, p := range paths {
		if v := r.Get(p); v.Exists() {
			return v
		}
	}
	return gjson.Result{}
}

func gcStrings(r gjson.Result) []string {
	var out []string
	for _, v := range r.Array() {
		out = append(out, v.String())
	}
	return out
}

// The controller sorts the variables and the series and keeps the tags in an order of its own, so the variables are put back in the order of the
// model, matched by name, and the tags and series are sets. A secret keeps its value out of the model.
func (m *GoldenConfigTemplates) fromAPI(ctx context.Context, r gjson.Result) diag.Diagnostics {
	var diags diag.Diagnostics
	m.Name = types.StringValue(r.Get("name").String())
	m.Id = m.Name
	if v := r.Get("versions.0.tcId"); v.Exists() {
		m.ConfigId = types.Int64Value(v.Int())
	}
	if d := r.Get("description").String(); d != "" || !m.Description.IsNull() {
		m.Description = types.StringValue(d)
	}
	m.Vendor = types.StringValue(r.Get("vendor").String())
	m.SoftwareType = types.StringValue(r.Get("softwareType").String())
	m.Audit = types.BoolValue(gcPick(r, "audit", "versions.0.audit").Bool())
	m.Mode = types.StringValue(gcPick(r, "versions.0.mode", "mode").String())
	m.ConfigType = types.StringValue(gcPick(r, "versions.0.configType", "configType").String())
	m.Config = types.StringValue(gcPick(r, "versions.0.config", "config").String())

	series, d := stringSet(ctx, gcStrings(r.Get("productSeries")))
	diags.Append(d...)
	m.ProductSeries = series
	if tags := gcStrings(r.Get("tags")); len(tags) > 0 || !m.Tags.IsNull() {
		set, d := stringSet(ctx, tags)
		diags.Append(d...)
		m.Tags = set
	}

	// Variables: the model's order, then any the model does not know, by name.
	api := map[string]gjson.Result{}
	var apiNames []string
	for _, v := range gcPick(r, "versions.0.variables", "variables").Array() {
		api[v.Get("name").String()] = v
		apiNames = append(apiNames, v.Get("name").String())
	}
	sort.Strings(apiNames)
	var current []GoldenConfigVariable
	if !m.Variables.IsNull() && !m.Variables.IsUnknown() {
		diags.Append(m.Variables.ElementsAs(ctx, &current, false)...)
	}
	var out []gcVar
	used := map[string]bool{}
	for _, cur := range current {
		v, ok := api[cur.Name.ValueString()]
		if !ok {
			continue
		}
		used[cur.Name.ValueString()] = true
		typ := v.Get("type").String()
		gv := gcVar{Name: cur.Name.ValueString(), Type: typ}
		if !cur.ValueWoVersion.IsNull() {
			gv.Secret, gv.Version = true, cur.ValueWoVersion.ValueInt64()
		} else {
			gv.Value = gcValueString(typ, v.Get("value"))
			if gcSameValue(typ, cur.Value.ValueString(), gv.Value) {
				gv.Value = cur.Value.ValueString()
			}
		}
		out = append(out, gv)
	}
	for _, name := range apiNames {
		if used[name] {
			continue
		}
		v := api[name]
		typ := v.Get("type").String()
		out = append(out, gcVar{Name: name, Type: typ, Value: gcValueString(typ, v.Get("value"))})
	}
	list, d := gcVarsList(ctx, out)
	diags.Append(d...)
	m.Variables = list
	return diags
}

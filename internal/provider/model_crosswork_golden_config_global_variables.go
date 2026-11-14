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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

type GoldenConfigGlobalVariables struct {
	Id             types.String `tfsdk:"id"`
	Instance       types.String `tfsdk:"instance"`
	Name           types.String `tfsdk:"name"`
	VariableId     types.String `tfsdk:"variable_id"`
	Type           types.String `tfsdk:"type"`
	Value          types.String `tfsdk:"value"`
	ValueWo        types.String `tfsdk:"value_wo"`
	ValueWoVersion types.Int64  `tfsdk:"value_wo_version"`
}

func (m GoldenConfigGlobalVariables) secret() bool { return !m.ValueWoVersion.IsNull() }

// A secret takes its value from the write-only configuration value.
func (m GoldenConfigGlobalVariables) body(cfg GoldenConfigGlobalVariables) string {
	value := m.Value.ValueString()
	if m.secret() {
		value = cfg.ValueWo.ValueString()
	}
	typ := m.Type.ValueString()
	return gcJSON([]any{map[string]any{"name": m.Name.ValueString(), "type": typ, "value": gcBodyValue(typ, value)}})
}

// A secret keeps its value out of the state; any other value keeps the text of the state when it means the same, so a JSON document is not rewritten.
func (m *GoldenConfigGlobalVariables) fromRow(row gjson.Result) {
	m.Name = types.StringValue(row.Get("name").String())
	m.Id = m.Name
	m.VariableId = types.StringValue(row.Get("id").String())
	typ := row.Get("type").String()
	if typ == "" {
		typ = "STRING"
	}
	m.Type = types.StringValue(typ)
	if m.secret() {
		return
	}
	api := gcValueString(typ, row.Get("value"))
	if m.Value.IsNull() || m.Value.IsUnknown() || !gcSameValue(typ, m.Value.ValueString(), api) {
		m.Value = types.StringValue(api)
	}
}

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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GoldenConfigJobDevice struct {
	Variables types.List `tfsdk:"variables"`
}

var goldenConfigJobDeviceType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"variables": types.ListType{ElemType: goldenConfigVariableType},
}}

type GoldenConfigJobs struct {
	Id                types.String `tfsdk:"id"`
	Instance          types.String `tfsdk:"instance"`
	TemplateName      types.String `tfsdk:"template_name"`
	TemplateConfigId  types.Int64  `tfsdk:"template_config_id"`
	Vendor            types.String `tfsdk:"vendor"`
	ProductSeries     types.String `tfsdk:"product_series"`
	Mop               types.String `tfsdk:"mop"`
	BatchSize         types.Int64  `tfsdk:"batch_size"`
	ErrorBudget       types.Int64  `tfsdk:"error_budget"`
	Tags              types.Set    `tfsdk:"tags"`
	DeploymentTimeout types.Int64  `tfsdk:"deployment_timeout"`
	Trigger           types.String `tfsdk:"trigger"`
	Devices           types.Map    `tfsdk:"devices"`
}

// Secrets come from cfg, the only place a write-only value exists; for a state pass the same map twice (a secret then has no value).
func gcDeviceVars(ctx context.Context, devices, cfg types.Map) (map[string][]gcVar, diag.Diagnostics) {
	var diags diag.Diagnostics
	out := map[string][]gcVar{}
	if devices.IsNull() || devices.IsUnknown() {
		return out, diags
	}
	var plan, conf map[string]GoldenConfigJobDevice
	diags.Append(devices.ElementsAs(ctx, &plan, false)...)
	diags.Append(cfg.ElementsAs(ctx, &conf, false)...)
	if diags.HasError() {
		return out, diags
	}
	for name, d := range plan {
		vars, vd := gcVarsOf(ctx, d.Variables, conf[name].Variables)
		diags.Append(vd...)
		out[name] = vars
	}
	return out, diags
}

func gcSubsetDevices(names []string, from ...types.Map) types.Map {
	elems := map[string]attr.Value{}
	for _, name := range names {
		for _, m := range from {
			if m.IsNull() || m.IsUnknown() {
				continue
			}
			if v, ok := m.Elements()[name]; ok {
				elems[name] = v
				break
			}
		}
	}
	return types.MapValueMust(goldenConfigJobDeviceType, elems)
}

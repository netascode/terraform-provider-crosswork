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
	"fmt"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

var _ datasource.DataSource = &GoldenConfigGlobalVariablesDataSource{}
var _ datasource.DataSourceWithConfigure = &GoldenConfigGlobalVariablesDataSource{}

func NewGoldenConfigGlobalVariablesDataSource() datasource.DataSource {
	return &GoldenConfigGlobalVariablesDataSource{}
}

type GoldenConfigGlobalVariablesDataSource struct {
	clients map[string]*cw.Client
}

type goldenConfigGlobalVariablesData struct {
	Id         types.String `tfsdk:"id"`
	Instance   types.String `tfsdk:"instance"`
	Name       types.String `tfsdk:"name"`
	VariableId types.String `tfsdk:"variable_id"`
	Type       types.String `tfsdk:"type"`
	Value      types.String `tfsdk:"value"`
	InUse      types.Bool   `tfsdk:"in_use"`
}

func (d *GoldenConfigGlobalVariablesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_golden_config_global_variables"
}

func (d *GoldenConfigGlobalVariablesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This data source can read a global variable of Golden Configuration.").String,

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object, the variable name",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Variable name").String,
				Computed:            true,
			},
			"variable_id": schema.StringAttribute{
				MarkdownDescription: "The id the controller gave the variable",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Type of the value").String,
				Computed:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Value of the variable; the controller returns it in clear text").String,
				Computed:            true,
				Sensitive:           true,
			},
			"in_use": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the controller reports the variable as used; it is not reliable, a variable a template uses can show false").String,
				Computed:            true,
			},
		},
	}
}

func (d *GoldenConfigGlobalVariablesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*CwProviderData).Clients
}

func (d *GoldenConfigGlobalVariablesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

func (d *GoldenConfigGlobalVariablesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config goldenConfigGlobalVariablesData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := d.clientFor(ctx, config.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", config.Id.ValueString()))

	row, found, err := gcFindGlobalVariable(&c, config.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (POST), got error: %s", err))
		return
	}
	if !found {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("The global variable %q does not exist", config.Id.ValueString()))
		return
	}
	typ := row.Get("type").String()
	config.Name = types.StringValue(row.Get("name").String())
	config.VariableId = types.StringValue(row.Get("id").String())
	config.Type = types.StringValue(typ)
	config.Value = types.StringValue(gcValueString(typ, row.Get("value")))
	config.InUse = types.BoolValue(row.Get("inUse").Bool())
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))
}

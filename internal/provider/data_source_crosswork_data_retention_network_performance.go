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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &DataRetentionNetworkPerformanceDataSource{}
	_ datasource.DataSourceWithConfigure = &DataRetentionNetworkPerformanceDataSource{}
)

func NewDataRetentionNetworkPerformanceDataSource() datasource.DataSource {
	return &DataRetentionNetworkPerformanceDataSource{}
}

type DataRetentionNetworkPerformanceDataSource struct {
	clients map[string]*cw.Client
}

func (d *DataRetentionNetworkPerformanceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_retention_network_performance"
}

func (d *DataRetentionNetworkPerformanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the network performance data retention.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"collect_lsp_state_change": schema.BoolAttribute{
				MarkdownDescription: "Collect the LSP state change events",
				Computed:            true,
			},
			"collect_lsp_path_change": schema.BoolAttribute{
				MarkdownDescription: "Collect the LSP path change events",
				Computed:            true,
			},
			"collect_lsp_utilization": schema.BoolAttribute{
				MarkdownDescription: "Collect the LSP utilization metrics",
				Computed:            true,
			},
			"collect_lsp_pm": schema.BoolAttribute{
				MarkdownDescription: "Collect the LSP PM metrics",
				Computed:            true,
			},
			"collect_link_pm": schema.BoolAttribute{
				MarkdownDescription: "Collect the link PM metrics",
				Computed:            true,
			},
			"raw_hours": schema.Int64Attribute{
				MarkdownDescription: "Hours that raw data is kept",
				Computed:            true,
			},
			"hourly_days": schema.Int64Attribute{
				MarkdownDescription: "Days that hourly data is kept",
				Computed:            true,
			},
			"daily_days": schema.Int64Attribute{
				MarkdownDescription: "Days that daily data is kept",
				Computed:            true,
			},
			"weekly_weeks": schema.Int64Attribute{
				MarkdownDescription: "Weeks that weekly data is kept. Each period must not be longer than the next",
				Computed:            true,
			},
		},
	}
}

func (d *DataRetentionNetworkPerformanceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *DataRetentionNetworkPerformanceDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Written by hand: the settings are two calls, which readGet makes and joins.
func (d *DataRetentionNetworkPerformanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config DataRetentionNetworkPerformance

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := d.clientFor(ctx, config.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", config.Id.String()))

	params := ""
	res, err := config.readGet(&c, "/crosswork/optima-analytics/api/v1/config/settings"+params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}

	config.fromBody(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

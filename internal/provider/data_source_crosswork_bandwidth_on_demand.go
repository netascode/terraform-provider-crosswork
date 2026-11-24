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
	_ datasource.DataSource              = &BandwidthOnDemandDataSource{}
	_ datasource.DataSourceWithConfigure = &BandwidthOnDemandDataSource{}
)

func NewBandwidthOnDemandDataSource() datasource.DataSource {
	return &BandwidthOnDemandDataSource{}
}

type BandwidthOnDemandDataSource struct {
	clients map[string]*cw.Client
}

func (d *BandwidthOnDemandDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bandwidth_on_demand"
}

func (d *BandwidthOnDemandDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the bandwidth on demand function pack of traffic engineering.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enables the bandwidth on demand function pack",
				Computed:            true,
			},
			"primary_objective": schema.StringAttribute{
				MarkdownDescription: "Primary objective when optimizing policies",
				Computed:            true,
			},
			"link_utilization_threshold": schema.Int64Attribute{
				MarkdownDescription: "Interfaces whose utilization percentage exceeds this threshold are optimized",
				Computed:            true,
			},
			"reoptimization_interval": schema.Int64Attribute{
				MarkdownDescription: "Time in seconds before policies are available to be re-optimized since the last path change",
				Computed:            true,
			},
			"metric_reoptimization_hours": schema.Int64Attribute{
				MarkdownDescription: "Hour of the day (0 to 23) to run the metric re-optimization on all BWoD provisioned policies",
				Computed:            true,
			},
			"metric_reoptimization_minutes": schema.Int64Attribute{
				MarkdownDescription: "Minute of the hour (0 to 59) to run the metric re-optimization on all BWoD provisioned policies",
				Computed:            true,
			},
			"sr_policy_traffic": schema.StringAttribute{
				MarkdownDescription: "Method to determine the bandwidth associated with each SR policy when optimization is performed",
				Computed:            true,
			},
			"model_update_throttle": schema.Int64Attribute{
				MarkdownDescription: "Time in seconds to wait between model updates; 0 disables the throttle",
				Computed:            true,
			},
			"optimizer_event_threshold": schema.Int64Attribute{
				MarkdownDescription: "If the optimizer runs longer than this number of seconds, an event is sent to the UI",
				Computed:            true,
			},
			"optimizer_run_time_limit": schema.Int64Attribute{
				MarkdownDescription: "Maximum run time in seconds of the BWoD optimizer",
				Computed:            true,
			},
			"allowed_policy_violations": schema.StringAttribute{
				MarkdownDescription: "How BWoD acts when new policies are provisioned and policy violations are detected",
				Computed:            true,
			},
			"prefer_strict_sids": schema.BoolAttribute{
				MarkdownDescription: "Prefer strict SIDs when deriving segment lists; required for compatibility with LCM",
				Computed:            true,
			},
			"debug_optimizer": schema.BoolAttribute{
				MarkdownDescription: "Save debug plan files to the filesystem",
				Computed:            true,
			},
			"debug_optimizer_max_files": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of debug plan file sets to save",
				Computed:            true,
			},
		},
	}
}

func (d *BandwidthOnDemandDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *BandwidthOnDemandDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *BandwidthOnDemandDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config BandwidthOnDemand

	// Read config
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
	res, err := c.Get("/crosswork/optima/api/v1/functionPacks/bwod/request" + params)
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

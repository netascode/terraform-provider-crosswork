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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &LocalCongestionMitigationDataSource{}
	_ datasource.DataSourceWithConfigure = &LocalCongestionMitigationDataSource{}
)

func NewLocalCongestionMitigationDataSource() datasource.DataSource {
	return &LocalCongestionMitigationDataSource{}
}

type LocalCongestionMitigationDataSource struct {
	clients map[string]*cw.Client
}

func (d *LocalCongestionMitigationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_congestion_mitigation"
}

func (d *LocalCongestionMitigationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the local congestion mitigation configuration of one LCM domain.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Required:            true,
			},
			"domain_id": schema.StringAttribute{
				MarkdownDescription: "LCM domain identifier",
				Computed:            true,
			},
			"enable": schema.BoolAttribute{
				MarkdownDescription: "Enables LCM for the domain",
				Computed:            true,
			},
			"color": schema.Int64Attribute{
				MarkdownDescription: "Color assigned sequentially to tactical SR policies starting with this value",
				Computed:            true,
			},
			"utilization_threshold": schema.Int64Attribute{
				MarkdownDescription: "Utilization percentage at which an interface is considered congested",
				Computed:            true,
			},
			"utilization_hold_margin": schema.Int64Attribute{
				MarkdownDescription: "Percentage that dampens the removal of deployed tactical SR policies",
				Computed:            true,
			},
			"delete_lsps": schema.BoolAttribute{
				MarkdownDescription: "Delete all deployed tactical SR policies when LCM is disabled",
				Computed:            true,
			},
			"profile_id": schema.Int64Attribute{
				MarkdownDescription: "Profile ID assigned to all tactical SR policies; 0 assigns none",
				Computed:            true,
			},
			"congestion_check_interval": schema.Int64Attribute{
				MarkdownDescription: "Interval in seconds used to evaluate the network for congestion",
				Computed:            true,
			},
			"max_tte_group": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of tactical policies used to mitigate a single interface",
				Computed:            true,
			},
			"include_all_interfaces": schema.BoolAttribute{
				MarkdownDescription: "Monitor congestion on all interfaces; false monitors the selected interfaces",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the domain configuration",
				Computed:            true,
			},
			"auto_repair_solution": schema.BoolAttribute{
				MarkdownDescription: "Detect and deploy changes automatically after network failures",
				Computed:            true,
			},
			"stay_in_area": schema.BoolAttribute{
				MarkdownDescription: "Restrict bypass paths to the mitigated OSPF area or ISIS level",
				Computed:            true,
			},
			"adjacency_hop_type": schema.StringAttribute{
				MarkdownDescription: "Type of adjacency SID used when creating an SR policy",
				Computed:            true,
			},
			"operation_mode": schema.StringAttribute{
				MarkdownDescription: "Operation mode of LCM; AUTOMATED needs gRPC on every device of the domain",
				Computed:            true,
			},
			"objective": schema.StringAttribute{
				MarkdownDescription: "Path computation objective of the tactical SR policies",
				Computed:            true,
			},
			"pce_dispatch_timeout": schema.Int64Attribute{
				MarkdownDescription: "Maximum time in seconds to confirm the deployment of tactical SR policies",
				Computed:            true,
			},
			"congestion_check_suspension_interval": schema.Int64Attribute{
				MarkdownDescription: "Seconds to suspend congestion detection after any tactical SR policy change",
				Computed:            true,
			},
			"overprovision_factor": schema.Int64Attribute{
				MarkdownDescription: "Percentage deducted from interface thresholds to derive the effective threshold",
				Computed:            true,
			},
			"uneven_ecmp_traffic_threshold": schema.Int64Attribute{
				MarkdownDescription: "Sensitivity percentage to detect uneven traffic across solution bypass tunnels",
				Computed:            true,
			},
			"throttle_mode_threshold": schema.Int64Attribute{
				MarkdownDescription: "Solution oscillations per hour that put the solution in the PAUSED state; 0 disables detection",
				Computed:            true,
			},
			"history_retention_time": schema.Int64Attribute{
				MarkdownDescription: "Days the collected data is retained",
				Computed:            true,
			},
			"debug_opt": schema.BoolAttribute{
				MarkdownDescription: "Save debug plan files and set the log level to debug",
				Computed:            true,
			},
			"debug_opt_max_files": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of debug plan file sets to save",
				Computed:            true,
			},
			"affinity_include_any": schema.ListAttribute{
				MarkdownDescription: "Bit positions of the link affinities of the include any constraint",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"affinity_exclude_any": schema.ListAttribute{
				MarkdownDescription: "Bit positions of the link affinities of the exclude any constraint",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"affinity_include_all": schema.ListAttribute{
				MarkdownDescription: "Bit positions of the link affinities of the include all constraint",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"node_tag_max_sid_depth": schema.MapAttribute{
				MarkdownDescription: "Maximum SID depth (0 to 16) by device tag, written as a text number",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *LocalCongestionMitigationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *LocalCongestionMitigationDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *LocalCongestionMitigationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config LocalCongestionMitigation

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
	res, err := config.readGet(&c, "/crosswork/optima/api/v1/functionPacks/lcm/domain/config"+params)
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

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
	_ datasource.DataSource              = &SRTEODNTemplatesDataSource{}
	_ datasource.DataSourceWithConfigure = &SRTEODNTemplatesDataSource{}
)

func NewSRTEODNTemplatesDataSource() datasource.DataSource {
	return &SRTEODNTemplatesDataSource{}
}

type SRTEODNTemplatesDataSource struct {
	clients map[string]*cw.Client
}

func (d *SRTEODNTemplatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_srte_odn_templates"
}

func (d *SRTEODNTemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read an SR-TE on-demand next-hop (ODN) template.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Template name",
				Required:            true,
			},
			"color": schema.Int64Attribute{
				MarkdownDescription: "SR policy color. Traffic whose route carries this color is steered into the policy that the template creates.",
				Computed:            true,
			},
			"head_end": schema.SetNestedAttribute{
				MarkdownDescription: "Head ends of the template.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of a device in the Crosswork inventory that is onboarded to NSO",
							Computed:            true,
						},
					},
				},
			},
			"maximum_sid_depth": schema.Int64Attribute{
				MarkdownDescription: "Maximum SID depth.",
				Computed:            true,
			},
			"pce_group": schema.StringAttribute{
				MarkdownDescription: "PCE group",
				Computed:            true,
			},
			"bandwidth": schema.Int64Attribute{
				MarkdownDescription: "Requested bandwidth in kbps. Not supported together with SRv6.",
				Computed:            true,
			},
			"source_address": schema.StringAttribute{
				MarkdownDescription: "Source address of the policies. An IPv6 address needs SRv6.",
				Computed:            true,
			},
			"srv6_locator_name": schema.StringAttribute{
				MarkdownDescription: "SRv6 locator name. Setting any SRv6 attribute turns on SRv6 for the template.",
				Computed:            true,
			},
			"srv6_locator_behavior": schema.StringAttribute{
				MarkdownDescription: "SRv6 behavior of the locator",
				Computed:            true,
			},
			"srv6_locator_binding_sid_type": schema.StringAttribute{
				MarkdownDescription: "SRv6 binding SID type",
				Computed:            true,
			},
			"dynamic_pce": schema.BoolAttribute{
				MarkdownDescription: "Delegate the path computation to the PCE. Set to `true` or leave unset.",
				Computed:            true,
			},
			"dynamic_flex_alg": schema.Int64Attribute{
				MarkdownDescription: "Flexible algorithm (prefix-SID algorithm)",
				Computed:            true,
			},
			"dynamic_metric_type": schema.StringAttribute{
				MarkdownDescription: "Metric type of the path computation",
				Computed:            true,
			},
			"dynamic_metric_margin_absolute": schema.Int64Attribute{
				MarkdownDescription: "Absolute metric margin. Use this or the relative margin.",
				Computed:            true,
			},
			"dynamic_metric_margin_relative": schema.Int64Attribute{
				MarkdownDescription: "Relative metric margin in percent. Use this or the absolute margin.",
				Computed:            true,
			},
			"dynamic_affinity_rules": schema.SetNestedAttribute{
				MarkdownDescription: "Affinity rules of the path computation",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							MarkdownDescription: "Rule action",
							Computed:            true,
						},
						"color": schema.SetAttribute{
							MarkdownDescription: "Names of the affinity colors of the rule. NSO returns them sorted.",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
			},
			"dynamic_segments_sid_algorithm": schema.Int64Attribute{
				MarkdownDescription: "SID algorithm constraint",
				Computed:            true,
			},
			"dynamic_disjoint_path_type": schema.StringAttribute{
				MarkdownDescription: "Disjointness type. With `dynamic_disjoint_path_group_id` it makes the paths of the group disjoint.",
				Computed:            true,
			},
			"dynamic_disjoint_path_group_id": schema.Int64Attribute{
				MarkdownDescription: "Disjoint group id. Required with the type.",
				Computed:            true,
			},
			"dynamic_disjoint_path_sub_id": schema.Int64Attribute{
				MarkdownDescription: "Disjoint sub group id",
				Computed:            true,
			},
			"delay_measurement_profile": schema.StringAttribute{
				MarkdownDescription: "Delay measurement profile of the head end. Not `default`; leave it unset for the default profile.",
				Computed:            true,
			},
			"delay_measurement_logging_delay_exceeded": schema.BoolAttribute{
				MarkdownDescription: "Log when the delay is exceeded. Set to `true` or leave unset.",
				Computed:            true,
			},
			"liveness_detection_profile": schema.StringAttribute{
				MarkdownDescription: "Liveness detection profile. Not `default`. Use this or delay measurement.",
				Computed:            true,
			},
			"liveness_detection_backup": schema.StringAttribute{
				MarkdownDescription: "Backup liveness detection profile. Needs a profile.",
				Computed:            true,
			},
			"liveness_detection_invalidation_action": schema.StringAttribute{
				MarkdownDescription: "Action when the liveness session invalidates the path",
				Computed:            true,
			},
			"liveness_detection_logging_session_state_change": schema.BoolAttribute{
				MarkdownDescription: "Log a change of the liveness session state. Set to `true` or leave unset.",
				Computed:            true,
			},
			"performance_measurement_reverse_path_label": schema.Int64Attribute{
				MarkdownDescription: "Reverse path label of the performance measurement",
				Computed:            true,
			},
		},
	}
}

func (d *SRTEODNTemplatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *SRTEODNTemplatesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *SRTEODNTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config SRTEODNTemplates

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
	params += "?content=config"
	res, err := c.Get(config.getPath() + params)
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

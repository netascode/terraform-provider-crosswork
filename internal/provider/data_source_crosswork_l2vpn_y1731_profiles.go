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
	_ datasource.DataSource              = &L2VPNY1731ProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &L2VPNY1731ProfilesDataSource{}
)

func NewL2VPNY1731ProfilesDataSource() datasource.DataSource {
	return &L2VPNY1731ProfilesDataSource{}
}

type L2VPNY1731ProfilesDataSource struct {
	clients map[string]*cw.Client
}

func (d *L2VPNY1731ProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2vpn_y1731_profiles"
}

func (d *L2VPNY1731ProfilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read a Y.1731 profile of the L2VPN.",

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
				MarkdownDescription: "Name of the profile. Any ASCII character except `%`, `?` and `#`.",
				Required:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Kind of measurement. `delay_statistics` belong to `delay`, `loss_statistics` to `loss` and `synthetic-loss`.",
				Computed:            true,
			},
			"schedule_interval": schema.Int64Attribute{
				MarkdownDescription: "Minutes between the operations. It has to be a factor of 1440 (1, 2, 3, 4, 5, 6, 8, 9, 10, 12, 15, 16, 18, 20, 24, 30, 32, 36, 40, 45, 48, 60, 80 or 90).",
				Computed:            true,
			},
			"schedule_duration": schema.StringAttribute{
				MarkdownDescription: "Minutes that the operations last, 1 to 1440, or `forever`",
				Computed:            true,
			},
			"probe_type": schema.StringAttribute{
				MarkdownDescription: "Kind of probe, `burst` when left out. With `packet` NSO drops the message count and the message period without an error.",
				Computed:            true,
			},
			"probe_message_count": schema.Int64Attribute{
				MarkdownDescription: "Messages of a burst. Message count times message period divided by 1000 must not exceed the measurement interval.",
				Computed:            true,
			},
			"probe_message_period": schema.Int64Attribute{
				MarkdownDescription: "Milliseconds between the messages of a burst",
				Computed:            true,
			},
			"probe_measurement_interval": schema.StringAttribute{
				MarkdownDescription: "Seconds of a measurement, 1 to 3600, or `once` (only with the probe type `packet`)",
				Computed:            true,
			},
			"probe_frame_size": schema.Int64Attribute{
				MarkdownDescription: "Size of the frame in bytes.",
				Computed:            true,
			},
			"probe_priority": schema.Int64Attribute{
				MarkdownDescription: "Priority of the frames",
				Computed:            true,
			},
			"delay_statistics": schema.SetNestedAttribute{
				MarkdownDescription: "Delay statistics, only for the type `delay`",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Statistic",
							Computed:            true,
						},
					},
				},
			},
			"loss_statistics": schema.SetNestedAttribute{
				MarkdownDescription: "Loss statistics, for the types `loss` and `synthetic-loss`",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Statistic",
							Computed:            true,
						},
					},
				},
			},
			"bucket_size": schema.Int64Attribute{
				MarkdownDescription: "Number of measurements in a bucket",
				Computed:            true,
			},
			"bucket_archive": schema.Int64Attribute{
				MarkdownDescription: "Number of buckets kept",
				Computed:            true,
			},
		},
	}
}

func (d *L2VPNY1731ProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *L2VPNY1731ProfilesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *L2VPNY1731ProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config L2VPNY1731Profiles

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

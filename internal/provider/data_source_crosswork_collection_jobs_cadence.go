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

import (
	"context"
	"fmt"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

var _ datasource.DataSource = &CollectionJobsCadenceDataSource{}
var _ datasource.DataSourceWithConfigure = &CollectionJobsCadenceDataSource{}

func NewCollectionJobsCadenceDataSource() datasource.DataSource {
	return &CollectionJobsCadenceDataSource{}
}

type CollectionJobsCadenceDataSource struct {
	clients map[string]*cw.Client
}

func (d *CollectionJobsCadenceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection_jobs_cadence"
}

func (d *CollectionJobsCadenceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This data source can read the cadence of the Crosswork collection jobs that Device Lifecycle Management runs itself, the settings of the Collection Jobs Cadence (DLM) panel of the Collection Jobs page.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"reachability_check": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Seconds between two reachability checks of the devices").String,
				Computed:            true,
			},
			"snmp_check": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Seconds between two SNMP checks of the devices").String,
				Computed:            true,
			},
			"te_tunnel_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Seconds between two collections of the TE tunnel ids").String,
				Computed:            true,
			},
			"show_clock": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Seconds between two collections of the clock of the devices (`show clock`)").String,
				Computed:            true,
			},
		},
	}
}

func (d *CollectionJobsCadenceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*CwProviderData).Clients
}

func (d *CollectionJobsCadenceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config CollectionJobsCadence
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := clientFor(ctx, d.clients, config.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, "Beginning Read")

	current, err := readCadence(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (POST), got error: %s", err))
		return
	}
	out := config.with(current)
	out.Id = types.StringValue(cadenceId)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", out.Id.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &out)...)
}

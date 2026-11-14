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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

// End of section. //template:end imports

// Written by hand: update_timeout is how the resource waits for a change, so the data source has no such attribute, and the
// collectors are the ones the cluster lists, not the ones to manage.

var (
	_ datasource.DataSource              = &DataGatewayResourceConfigurationDataSource{}
	_ datasource.DataSourceWithConfigure = &DataGatewayResourceConfigurationDataSource{}
)

func NewDataGatewayResourceConfigurationDataSource() datasource.DataSource {
	return &DataGatewayResourceConfigurationDataSource{}
}

type DataGatewayResourceConfigurationDataSource struct {
	clients map[string]*cw.Client
}

func (d *DataGatewayResourceConfigurationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_gateway_resource_configuration"
}

func (d *DataGatewayResourceConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This data source can read the collectors of the Crosswork Data Gateways with their memory and state, the settings of the Resource configuration page. A cluster lists only the collectors it has.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"collectors": schema.ListNestedAttribute{
				MarkdownDescription: "The collectors that the cluster lists, with their memory and state",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the collector, as the Resource configuration page shows it",
							Computed:            true,
						},
						"memory": schema.Int64Attribute{
							MarkdownDescription: "Memory of the collector in MB",
							Computed:            true,
						},
						"enabled": schema.BoolAttribute{
							MarkdownDescription: "Whether the collector is enabled, the Enable collector checkbox",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DataGatewayResourceConfigurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

func (d *DataGatewayResourceConfigurationDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Written by hand: the data source has no update_timeout, so it reads into its own object and the collectors come from the
// model's reader; the settings have no id of their own, so it gives the one of the resource.
type dataGatewayCollectorsData struct {
	Instance   types.String                                 `tfsdk:"instance"`
	Id         types.String                                 `tfsdk:"id"`
	Collectors []DataGatewayResourceConfigurationCollectors `tfsdk:"collectors"`
}

func (d *DataGatewayResourceConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config dataGatewayCollectorsData

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

	model := DataGatewayResourceConfiguration{}
	res, err := model.readGet(&c, resourcesQueryPath)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}

	model.fromBody(ctx, res)
	config.Collectors = model.Collectors
	config.Id = types.StringValue("data_gateway_resource_configuration")

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

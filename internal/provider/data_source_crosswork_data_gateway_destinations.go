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

// Written by hand: the query of the controller takes a limit and no id, so the destination is found by its id or its name in the
// list of all destinations, the system-defined ones included.

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

var _ datasource.DataSource = &DataGatewayDestinationsDataSource{}
var _ datasource.DataSourceWithConfigure = &DataGatewayDestinationsDataSource{}
var _ datasource.DataSourceWithConfigValidators = &DataGatewayDestinationsDataSource{}

func NewDataGatewayDestinationsDataSource() datasource.DataSource {
	return &DataGatewayDestinationsDataSource{}
}

type DataGatewayDestinationsDataSource struct {
	clients map[string]*cw.Client
}

func (d *DataGatewayDestinationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_gateway_destinations"
}

func (d *DataGatewayDestinationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This data source can read a data destination of the Crosswork Data Gateways by its name or its id, including the system-defined destinations.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the destination; set either this or `name`",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the destination; set either this or `id`",
				Optional:            true,
				Computed:            true,
			},
			"data_source": schema.StringAttribute{
				MarkdownDescription: "Source of the data sent to the destination",
				Computed:            true,
			},
			"server_type": schema.StringAttribute{
				MarkdownDescription: "Type of the destination server",
				Computed:            true,
			},
			"encoding": schema.StringAttribute{
				MarkdownDescription: "Encoding of the data, for the data source `data-gateway`",
				Computed:            true,
			},
			"compression_method": schema.StringAttribute{
				MarkdownDescription: "Compression of the data",
				Computed:            true,
			},
			"dispatch_type": schema.StringAttribute{
				MarkdownDescription: "Dispatch type of a gRPC destination",
				Computed:            true,
			},
			"max_message_size": schema.Int64Attribute{
				MarkdownDescription: "Maximum message size of a Kafka destination in bytes",
				Computed:            true,
			},
			"buffer_memory": schema.Int64Attribute{
				MarkdownDescription: "Buffer memory of a Kafka destination in bytes",
				Computed:            true,
			},
			"batch_size": schema.Int64Attribute{
				MarkdownDescription: "Batch size of a Kafka destination in bytes",
				Computed:            true,
			},
			"linger_time": schema.Int64Attribute{
				MarkdownDescription: "Linger time of a Kafka destination in milliseconds",
				Computed:            true,
			},
			"request_timeout": schema.Int64Attribute{
				MarkdownDescription: "Request timeout of a Kafka destination in seconds",
				Computed:            true,
			},
			"collectors": schema.ListNestedAttribute{
				MarkdownDescription: "Values of a Kafka destination that replace the global ones for one collector",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Collector",
							Computed:            true,
						},
						"buffer_memory": schema.Int64Attribute{
							MarkdownDescription: "Buffer memory in bytes",
							Computed:            true,
						},
						"batch_size": schema.Int64Attribute{
							MarkdownDescription: "Batch size in bytes",
							Computed:            true,
						},
						"linger_time": schema.Int64Attribute{
							MarkdownDescription: "Linger time in milliseconds",
							Computed:            true,
						},
						"request_timeout": schema.Int64Attribute{
							MarkdownDescription: "Request timeout in seconds",
							Computed:            true,
						},
					},
				},
			},
			"connections": schema.ListNestedAttribute{
				MarkdownDescription: "Addresses of the destination server, each an IPv4 address with a prefix length or a host name with a domain name",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							MarkdownDescription: "IPv4 address",
							Computed:            true,
						},
						"prefix_length": schema.Int64Attribute{
							MarkdownDescription: "Prefix length of the IPv4 address",
							Computed:            true,
						},
						"host_name": schema.StringAttribute{
							MarkdownDescription: "Host name",
							Computed:            true,
						},
						"domain_name": schema.StringAttribute{
							MarkdownDescription: "Domain name",
							Computed:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "Port of the server",
							Computed:            true,
						},
					},
				},
			},
			"secure": schema.BoolAttribute{
				MarkdownDescription: "Enable secure communication",
				Computed:            true,
			},
			"auth_type": schema.StringAttribute{
				MarkdownDescription: "Authentication of the secure communication; a gRPC destination takes only `mutual`",
				Computed:            true,
			},
			"system_defined": schema.BoolAttribute{
				MarkdownDescription: "Whether Crosswork defined the destination, which makes it read-only",
				Computed:            true,
			},
		},
	}
}

func (d *DataGatewayDestinationsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
	}
}

func (d *DataGatewayDestinationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*CwProviderData).Clients
}

func (d *DataGatewayDestinationsDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

func (d *DataGatewayDestinationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config DataGatewayDestinations
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := d.clientFor(ctx, config.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s%s: Beginning Read", config.Id.ValueString(), config.Name.ValueString()))

	res, err := listDestinations(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the destinations, got error: %s, %s", err, res.String()))
		return
	}
	row := findDestination(res, config.Id.ValueString(), config.Name.ValueString())
	if !row.Exists() {
		resp.Diagnostics.AddError("Destination not found", fmt.Sprintf("There is no destination with the id %q or the name %q on the cluster", config.Id.ValueString(), config.Name.ValueString()))
		return
	}
	state := DataGatewayDestinations{Instance: config.Instance}
	state.fromBody(ctx, row)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))
}

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

// Written by hand: the query of the controller takes no filter, so the pool is found by its id or its name in the list of all pools.

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

var _ datasource.DataSource = &DataGatewayPoolsDataSource{}
var _ datasource.DataSourceWithConfigure = &DataGatewayPoolsDataSource{}
var _ datasource.DataSourceWithConfigValidators = &DataGatewayPoolsDataSource{}

func NewDataGatewayPoolsDataSource() datasource.DataSource {
	return &DataGatewayPoolsDataSource{}
}

type DataGatewayPoolsDataSource struct {
	clients map[string]*cw.Client
}

type dataGatewayPoolLookup struct {
	Instance         types.String `tfsdk:"instance"`
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	VirtualIps       types.List   `tfsdk:"virtual_ips"`
	PrefixLength     types.Int64  `tfsdk:"prefix_length"`
	Spares           types.Int64  `tfsdk:"spares"`
	Instances        types.List   `tfsdk:"instances"`
	DataGatewayIds   types.List   `tfsdk:"data_gateway_ids"`
	ProtectionStatus types.String `tfsdk:"protection_status"`
}

func (d *DataGatewayPoolsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_gateway_pools"
}

func (d *DataGatewayPoolsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This data source can read a pool of Crosswork Data Gateway instances by its name or its id.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Id of the pool; set either this or `name`",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the pool; set either this or `id`",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the pool",
				Computed:            true,
			},
			"virtual_ips": schema.ListAttribute{
				MarkdownDescription: "Virtual IP addresses of the pool, in the order of its data gateways",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"prefix_length": schema.Int64Attribute{
				MarkdownDescription: "Prefix length of the virtual IP addresses",
				Computed:            true,
			},
			"spares": schema.Int64Attribute{
				MarkdownDescription: "Number of spare instances of the pool",
				Computed:            true,
			},
			"instances": schema.ListAttribute{
				MarkdownDescription: "Ids of the data gateway instances of the pool",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"data_gateway_ids": schema.ListAttribute{
				MarkdownDescription: "Ids of the data gateways of the pool, in the order of the virtual IPs; the gateways are named `<pool name>-<n>`",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"protection_status": schema.StringAttribute{
				MarkdownDescription: "High availability status of the pool, such as `PROTECTED`, `NOT_PROTECTED` or `NOT_PLANNED`",
				Computed:            true,
			},
		},
	}
}

func (d *DataGatewayPoolsDataSource) ConfigValidators(_ context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("id"), path.MatchRoot("name")),
	}
}

func (d *DataGatewayPoolsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*CwProviderData).Clients
}

func (d *DataGatewayPoolsDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

func (d *DataGatewayPoolsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config dataGatewayPoolLookup
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := d.clientFor(ctx, config.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s%s: Beginning Read", config.Id.ValueString(), config.Name.ValueString()))

	res, err := listPools(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the pools, got error: %s, %s", err, res.String()))
		return
	}
	pool := findPool(res, config.Id.ValueString(), config.Name.ValueString())
	if !pool.Exists() {
		resp.Diagnostics.AddError("Pool not found", fmt.Sprintf("There is no pool with the id %q or the name %q on the cluster", config.Id.ValueString(), config.Name.ValueString()))
		return
	}
	var state DataGatewayPools
	state.fromBody(ctx, pool)
	gateways, err := listPoolGateways(&c, pool.Get("puuid").String())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the data gateways of the pool, got error: %s", err))
		return
	}
	config.Id = state.Id
	config.Name = state.Name
	config.Description = state.Description
	config.VirtualIps = state.VirtualIps
	config.PrefixLength = state.PrefixLength
	config.Spares = state.Spares
	config.Instances = state.Instances
	config.DataGatewayIds = stringList(gateways)
	config.ProtectionStatus = optionalString(pool.Get("protectionStatus").String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))
}

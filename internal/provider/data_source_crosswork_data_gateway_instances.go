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

var _ datasource.DataSource = &DataGatewayInstancesDataSource{}
var _ datasource.DataSourceWithConfigure = &DataGatewayInstancesDataSource{}

func NewDataGatewayInstancesDataSource() datasource.DataSource {
	return &DataGatewayInstancesDataSource{}
}

type DataGatewayInstancesDataSource struct {
	clients map[string]*cw.Client
}

func (d *DataGatewayInstancesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_gateway_instances"
}

func (d *DataGatewayInstancesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This data source fetches the Crosswork Data Gateway instances that are enrolled on the cluster, with the pool, the data gateway and the state of each. The instances are virtual machines that enrol themselves; Terraform cannot create them.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"items": schema.SetNestedAttribute{
				MarkdownDescription: "List of the data gateway instances of the cluster",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Id of the instance",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the instance",
							Computed:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "Role of the instance in its pool, `ASSIGNED`, `SPARE` or `UNASSIGNED`",
							Computed:            true,
						},
						"pool_id": schema.StringAttribute{
							MarkdownDescription: "Id of the pool of the instance, not set for an unassigned instance",
							Computed:            true,
						},
						"vdg_id": schema.StringAttribute{
							MarkdownDescription: "Id of the data gateway the instance is the active instance of, not set for a spare or an unassigned instance",
							Computed:            true,
						},
						"oper_state": schema.StringAttribute{
							MarkdownDescription: "Operational state of the instance, such as `OS_UP`, `OS_NOT_READY` or `OS_ERROR`",
							Computed:            true,
						},
						"admin_state": schema.StringAttribute{
							MarkdownDescription: "Administrative state of the instance, such as `AS_UP`",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DataGatewayInstancesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*CwProviderData).Clients
}

func (d *DataGatewayInstancesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

func (d *DataGatewayInstancesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config DataGatewayInstances
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := d.clientFor(ctx, config.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, "data gateway instances: Beginning Read")

	res, err := c.Post(instanceQueryPath, `{}`, cw.ReadOnly)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the data gateway instances, got error: %s, %s", err, res.String()))
		return
	}
	config.Id = types.StringValue("data_gateway_instances")
	config.Items = parseInstances(res)
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	tflog.Debug(ctx, "data gateway instances: Read finished successfully")
}

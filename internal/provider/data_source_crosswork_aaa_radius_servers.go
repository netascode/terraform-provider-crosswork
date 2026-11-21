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
	_ datasource.DataSource              = &AAARADIUSServersDataSource{}
	_ datasource.DataSourceWithConfigure = &AAARADIUSServersDataSource{}
)

func NewAAARADIUSServersDataSource() datasource.DataSource {
	return &AAARADIUSServersDataSource{}
}

type AAARADIUSServersDataSource struct {
	clients map[string]*cw.Client
}

func (d *AAARADIUSServersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aaa_radius_servers"
}

func (d *AAARADIUSServersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the RADIUS servers used for authentication.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"servers": schema.ListNestedAttribute{
				MarkdownDescription: "RADIUS servers",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							MarkdownDescription: "IP address or host name of the server",
							Computed:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "Port of the server",
							Computed:            true,
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: "Shared secret",
							Computed:            true,
						},
						"secret_wo": schema.StringAttribute{
							MarkdownDescription: "Shared secret",
							Computed:            true,
						},
						"secret_wo_version": schema.Int64Attribute{
							MarkdownDescription: "Rotation trigger for `secret_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.",
							Computed:            true,
						},
						"secret_format": schema.StringAttribute{
							MarkdownDescription: "Format of the shared secret",
							Computed:            true,
						},
						"retransmit_timeout": schema.Int64Attribute{
							MarkdownDescription: "Retransmit timeout in seconds",
							Computed:            true,
						},
						"retries": schema.Int64Attribute{
							MarkdownDescription: "Number of retries before the request fails",
							Computed:            true,
						},
						"auth_type": schema.StringAttribute{
							MarkdownDescription: "Authentication type",
							Computed:            true,
						},
						"priority": schema.Int64Attribute{
							MarkdownDescription: "Order of the server among all remote authentication servers",
							Computed:            true,
						},
						"policy_id": schema.StringAttribute{
							MarkdownDescription: "Key name of the role on the server",
							Computed:            true,
						},
						"virtual_domain": schema.StringAttribute{
							MarkdownDescription: "Device access group attribute",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AAARADIUSServersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *AAARADIUSServersDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *AAARADIUSServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config AAARADIUSServers

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
	res, err := c.Get("/crosswork/authconfig/v1/remote/radius" + params)
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

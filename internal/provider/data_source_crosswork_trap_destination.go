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
	"net/url"

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
	_ datasource.DataSource              = &TrapDestinationDataSource{}
	_ datasource.DataSourceWithConfigure = &TrapDestinationDataSource{}
)

func NewTrapDestinationDataSource() datasource.DataSource {
	return &TrapDestinationDataSource{}
}

type TrapDestinationDataSource struct {
	clients map[string]*cw.Client
}

func (d *TrapDestinationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trap_destination"
}

func (d *TrapDestinationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read a trap receiver of the notification destinations.",

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
				MarkdownDescription: "Name of the receiver. A receiver whose address is a host name takes the host name as its name",
				Computed:            true,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "IP address or host name of the receiver",
				Computed:            true,
			},
			"address_type": schema.StringAttribute{
				MarkdownDescription: "Whether the address is an IP address or a host name",
				Computed:            true,
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: "Port of the receiver",
				Computed:            true,
			},
			"transport_type": schema.StringAttribute{
				MarkdownDescription: "Transport protocol",
				Computed:            true,
			},
			"snmp_version": schema.StringAttribute{
				MarkdownDescription: "SNMP version",
				Computed:            true,
			},
			"community": schema.StringAttribute{
				MarkdownDescription: "SNMP v2c community",
				Computed:            true,
			},
			"community_wo": schema.StringAttribute{
				MarkdownDescription: "SNMP v2c community",
				Computed:            true,
			},
			"community_wo_version": schema.Int64Attribute{
				MarkdownDescription: "Rotation trigger for `community_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.",
				Computed:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: "SNMP v3 user name",
				Computed:            true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: "SNMP v3 security mode",
				Computed:            true,
			},
			"auth_type": schema.StringAttribute{
				MarkdownDescription: "SNMP v3 authentication type",
				Computed:            true,
			},
			"auth_password": schema.StringAttribute{
				MarkdownDescription: "SNMP v3 authentication password",
				Computed:            true,
			},
			"auth_password_wo": schema.StringAttribute{
				MarkdownDescription: "SNMP v3 authentication password",
				Computed:            true,
			},
			"auth_password_wo_version": schema.Int64Attribute{
				MarkdownDescription: "Rotation trigger for `auth_password_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.",
				Computed:            true,
			},
			"privacy_type": schema.StringAttribute{
				MarkdownDescription: "SNMP v3 privacy type. `NONE` is for the mode `AuthNoPriv` and is refused with `AuthPriv`",
				Computed:            true,
			},
			"privacy_password": schema.StringAttribute{
				MarkdownDescription: "SNMP v3 privacy password",
				Computed:            true,
			},
			"privacy_password_wo": schema.StringAttribute{
				MarkdownDescription: "SNMP v3 privacy password",
				Computed:            true,
			},
			"privacy_password_wo_version": schema.Int64Attribute{
				MarkdownDescription: "Rotation trigger for `privacy_password_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.",
				Computed:            true,
			},
		},
	}
}

func (d *TrapDestinationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *TrapDestinationDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Written by hand: the list of the controller holds every receiver, which readGet searches.
func (d *TrapDestinationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config TrapDestination

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
	params += "/" + url.QueryEscape(config.Id.ValueString())
	res, err := config.readGet(&c, "/crosswork/alarm/v1/destination/trap"+params)
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

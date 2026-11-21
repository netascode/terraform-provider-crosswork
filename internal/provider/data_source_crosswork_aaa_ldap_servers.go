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
	_ datasource.DataSource              = &AAALDAPServersDataSource{}
	_ datasource.DataSourceWithConfigure = &AAALDAPServersDataSource{}
)

func NewAAALDAPServersDataSource() datasource.DataSource {
	return &AAALDAPServersDataSource{}
}

type AAALDAPServersDataSource struct {
	clients map[string]*cw.Client
}

func (d *AAALDAPServersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aaa_ldap_servers"
}

func (d *AAALDAPServersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the LDAP servers used for authentication.",

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
				MarkdownDescription: "LDAP servers",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the connection",
							Computed:            true,
						},
						"ldap_url": schema.StringAttribute{
							MarkdownDescription: "URL of the server, `ldap://` or `ldaps://`",
							Computed:            true,
						},
						"bind_dn": schema.StringAttribute{
							MarkdownDescription: "Distinguished name used to bind to the server",
							Computed:            true,
						},
						"bind_credential": schema.StringAttribute{
							MarkdownDescription: "Password used to bind to the server",
							Computed:            true,
						},
						"bind_credential_wo": schema.StringAttribute{
							MarkdownDescription: "Password used to bind to the server",
							Computed:            true,
						},
						"bind_credential_wo_version": schema.Int64Attribute{
							MarkdownDescription: "Rotation trigger for `bind_credential_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.",
							Computed:            true,
						},
						"base_dn": schema.StringAttribute{
							MarkdownDescription: "Distinguished name where the search for users starts",
							Computed:            true,
						},
						"user_filter": schema.StringAttribute{
							MarkdownDescription: "Filter that matches the users, with `{user}` for the user name",
							Computed:            true,
						},
						"dn_format": schema.StringAttribute{
							MarkdownDescription: "Format of the distinguished name of a user, with `%s` for the user name",
							Computed:            true,
						},
						"principal_attribute_id": schema.StringAttribute{
							MarkdownDescription: "Attribute of the user name in the directory",
							Computed:            true,
						},
						"policy_id": schema.StringAttribute{
							MarkdownDescription: "Attribute of the role mapping in the directory",
							Computed:            true,
						},
						"order": schema.Int64Attribute{
							MarkdownDescription: "Order of the server among all remote authentication servers",
							Computed:            true,
						},
						"connect_timeout": schema.Int64Attribute{
							MarkdownDescription: "Connection timeout in seconds",
							Computed:            true,
						},
					},
				},
			},
			"attribute_servers": schema.ListNestedAttribute{
				MarkdownDescription: "LDAP servers that look up the role of a user, one for each server",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ldap_url": schema.StringAttribute{
							MarkdownDescription: "URL of the server, `ldap://` followed by host and port",
							Computed:            true,
						},
						"base_dn": schema.StringAttribute{
							MarkdownDescription: "Distinguished name where the search for the role starts",
							Computed:            true,
						},
						"bind_dn": schema.StringAttribute{
							MarkdownDescription: "Distinguished name used to bind to the server",
							Computed:            true,
						},
						"bind_credential": schema.StringAttribute{
							MarkdownDescription: "Password used to bind to the server",
							Computed:            true,
						},
						"bind_credential_wo": schema.StringAttribute{
							MarkdownDescription: "Password used to bind to the server",
							Computed:            true,
						},
						"bind_credential_wo_version": schema.Int64Attribute{
							MarkdownDescription: "Rotation trigger for `bind_credential_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.",
							Computed:            true,
						},
						"user_filter": schema.StringAttribute{
							MarkdownDescription: "Attribute of the directory that maps the user to its role",
							Computed:            true,
						},
						"connect_timeout": schema.Int64Attribute{
							MarkdownDescription: "Connection timeout in seconds",
							Computed:            true,
						},
						"fail_fast": schema.BoolAttribute{
							MarkdownDescription: "Fail at once when the server cannot be reached",
							Computed:            true,
						},
						"attributes": schema.MapAttribute{
							MarkdownDescription: "Map of the directory key that carries the role to `policy_id`",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AAALDAPServersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *AAALDAPServersDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *AAALDAPServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config AAALDAPServers

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
	res, err := c.Get("/crosswork/authconfig/v1/remote/ldap" + params)
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

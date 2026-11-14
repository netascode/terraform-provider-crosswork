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
	_ datasource.DataSource              = &CredentialProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &CredentialProfilesDataSource{}
)

func NewCredentialProfilesDataSource() datasource.DataSource {
	return &CredentialProfilesDataSource{}
}

type CredentialProfilesDataSource struct {
	clients map[string]*cw.Client
}

func (d *CredentialProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_profiles"
}

func (d *CredentialProfilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the Credential Profiles.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Required:            true,
			},
			"profile": schema.StringAttribute{
				MarkdownDescription: "Credential profile name",
				Computed:            true,
			},
			"credentials": schema.ListNestedAttribute{
				MarkdownDescription: "List of credentials, one per access type",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: "Access type",
							Computed:            true,
						},
						"user_name": schema.StringAttribute{
							MarkdownDescription: "User name",
							Computed:            true,
						},
						"password": schema.StringAttribute{
							MarkdownDescription: "Password",
							Computed:            true,
						},
						"password_wo": schema.StringAttribute{
							MarkdownDescription: "Password",
							Computed:            true,
						},
						"password_wo_version": schema.Int64Attribute{
							MarkdownDescription: "Rotation trigger for `password_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.",
							Computed:            true,
						},
						"enable_password": schema.StringAttribute{
							MarkdownDescription: "Second level authentication, for devices that support it",
							Computed:            true,
						},
						"enable_password_value": schema.StringAttribute{
							MarkdownDescription: "Enable password",
							Computed:            true,
						},
						"enable_password_value_wo": schema.StringAttribute{
							MarkdownDescription: "Enable password",
							Computed:            true,
						},
						"enable_password_value_wo_version": schema.Int64Attribute{
							MarkdownDescription: "Rotation trigger for `enable_password_value_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.",
							Computed:            true,
						},
					},
				},
			},
			"snmpv2_read_community": schema.StringAttribute{
				MarkdownDescription: "SNMPv2 read community",
				Computed:            true,
			},
			"snmpv2_read_community_wo": schema.StringAttribute{
				MarkdownDescription: "SNMPv2 read community",
				Computed:            true,
			},
			"snmpv2_read_community_wo_version": schema.Int64Attribute{
				MarkdownDescription: "Rotation trigger for `snmpv2_read_community_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.",
				Computed:            true,
			},
			"snmpv2_write_community": schema.StringAttribute{
				MarkdownDescription: "SNMPv2 write community",
				Computed:            true,
			},
			"snmpv2_write_community_wo": schema.StringAttribute{
				MarkdownDescription: "SNMPv2 write community",
				Computed:            true,
			},
			"snmpv2_write_community_wo_version": schema.Int64Attribute{
				MarkdownDescription: "Rotation trigger for `snmpv2_write_community_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.",
				Computed:            true,
			},
			"snmpv3_user_name": schema.StringAttribute{
				MarkdownDescription: "SNMPv3 user name",
				Computed:            true,
			},
			"snmpv3_security_level": schema.StringAttribute{
				MarkdownDescription: "SNMPv3 security level",
				Computed:            true,
			},
			"snmpv3_auth_type": schema.StringAttribute{
				MarkdownDescription: "SNMPv3 authentication type",
				Computed:            true,
			},
			"snmpv3_auth_password": schema.StringAttribute{
				MarkdownDescription: "SNMPv3 authentication password",
				Computed:            true,
			},
			"snmpv3_auth_password_wo": schema.StringAttribute{
				MarkdownDescription: "SNMPv3 authentication password",
				Computed:            true,
			},
			"snmpv3_auth_password_wo_version": schema.Int64Attribute{
				MarkdownDescription: "Rotation trigger for `snmpv3_auth_password_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.",
				Computed:            true,
			},
			"snmpv3_priv_type": schema.StringAttribute{
				MarkdownDescription: "SNMPv3 privacy type",
				Computed:            true,
			},
			"snmpv3_priv_password": schema.StringAttribute{
				MarkdownDescription: "SNMPv3 privacy password",
				Computed:            true,
			},
			"snmpv3_priv_password_wo": schema.StringAttribute{
				MarkdownDescription: "SNMPv3 privacy password",
				Computed:            true,
			},
			"snmpv3_priv_password_wo_version": schema.Int64Attribute{
				MarkdownDescription: "Rotation trigger for `snmpv3_priv_password_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.",
				Computed:            true,
			},
		},
	}
}

func (d *CredentialProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *CredentialProfilesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *CredentialProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config CredentialProfiles

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
	res, err := config.readGet(&c, "/crosswork/inventory/v1/credentials/query"+params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}
	if data := res.Get("data"); !data.Exists() || len(data.Array()) == 0 {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No object found with id: %s", config.Id.ValueString()))
		return
	}

	config.fromBody(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

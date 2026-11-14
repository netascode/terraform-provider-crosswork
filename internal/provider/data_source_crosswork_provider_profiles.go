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
	_ datasource.DataSource              = &ProviderProfilesDataSource{}
	_ datasource.DataSourceWithConfigure = &ProviderProfilesDataSource{}
)

func NewProviderProfilesDataSource() datasource.DataSource {
	return &ProviderProfilesDataSource{}
}

type ProviderProfilesDataSource struct {
	clients map[string]*cw.Client
}

func (d *ProviderProfilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_profiles"
}

func (d *ProviderProfilesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the Provider Profiles.",

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
				MarkdownDescription: "Provider name, unique on the Crosswork cluster",
				Computed:            true,
			},
			"family": schema.StringAttribute{
				MarkdownDescription: "Provider family; the controller stores the SR PCE family of the Add Provider form as `ROBOT_PROVIDER_XTC`",
				Computed:            true,
			},
			"profile": schema.StringAttribute{
				MarkdownDescription: "Name of the credential profile, which must hold a username and password for every protocol the provider uses",
				Computed:            true,
			},
			"connectivity": schema.ListNestedAttribute{
				MarkdownDescription: "Connections to the provider, one per protocol; each uses an IP address or an FQDN",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"protocol": schema.StringAttribute{
							MarkdownDescription: "Protocol",
							Computed:            true,
						},
						"ip_address": schema.StringAttribute{
							MarkdownDescription: "IPv4 address of the provider; leave it unset when `fqdn_host_name` is used",
							Computed:            true,
						},
						"fqdn_host_name": schema.StringAttribute{
							MarkdownDescription: "Host name part of the FQDN of the provider; leave it unset when `ip_address` is used",
							Computed:            true,
						},
						"fqdn_domain_name": schema.StringAttribute{
							MarkdownDescription: "Domain name part of the FQDN of the provider",
							Computed:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "Port",
							Computed:            true,
						},
						"timeout": schema.Int64Attribute{
							MarkdownDescription: "Timeout in seconds, the controller uses 30 when it is not set",
							Computed:            true,
						},
						"encoding_type": schema.StringAttribute{
							MarkdownDescription: "Encoding of a gNMI or gRPC connection, required for the gNMI protocol",
							Computed:            true,
						},
					},
				},
			},
			"certificate_name": schema.StringAttribute{
				MarkdownDescription: "Name of the certificate profile of the provider, required by some families (Accedian proxy) and checked to exist",
				Computed:            true,
			},
			"supported_models": schema.ListNestedAttribute{
				MarkdownDescription: "Device models the provider supports, required for the NSO family",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"model_prefix": schema.StringAttribute{
							MarkdownDescription: "Model prefix",
							Computed:            true,
						},
						"model_version": schema.StringAttribute{
							MarkdownDescription: "Model version",
							Computed:            true,
						},
					},
				},
			},
			"properties": schema.MapAttribute{
				MarkdownDescription: "Provider properties such as `nso_crosslaunch_url`, `embedded_nso` or `webui-root-resource`, as key and value",
				ElementType:         types.StringType,
				Computed:            true,
			},
		},
	}
}

func (d *ProviderProfilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *ProviderProfilesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *ProviderProfilesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config ProviderProfiles

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
	res, err := config.readGet(&c, "/crosswork/inventory/v1/providers/query"+params)
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

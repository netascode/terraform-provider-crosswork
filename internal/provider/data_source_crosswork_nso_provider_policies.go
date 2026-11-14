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
	_ datasource.DataSource              = &NSOProviderPoliciesDataSource{}
	_ datasource.DataSourceWithConfigure = &NSOProviderPoliciesDataSource{}
)

func NewNSOProviderPoliciesDataSource() datasource.DataSource {
	return &NSOProviderPoliciesDataSource{}
}

type NSOProviderPoliciesDataSource struct {
	clients map[string]*cw.Client
}

func (d *NSOProviderPoliciesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nso_provider_policies"
}

func (d *NSOProviderPoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the device onboarding policy of an NSO provider.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Required:            true,
			},
			"provider_name": schema.StringAttribute{
				MarkdownDescription: "Name of the NSO provider the policy belongs to",
				Required:            true,
			},
			"match": schema.BoolAttribute{
				MarkdownDescription: "Match devices of Crosswork and of the NSO provider by their IP address",
				Computed:            true,
			},
			"match_rule": schema.StringAttribute{
				MarkdownDescription: "Expression that selects the devices that are matched, for example `*` or `host_name = 'build-*'`",
				Computed:            true,
			},
			"onboard_to_nso": schema.BoolAttribute{
				MarkdownDescription: "Add the devices that Crosswork has and the NSO provider is missing to the NSO provider",
				Computed:            true,
			},
			"onboard_to_nso_rule": schema.StringAttribute{
				MarkdownDescription: "Expression that selects the devices that are added to the NSO provider",
				Computed:            true,
			},
			"onboard_from_nso": schema.BoolAttribute{
				MarkdownDescription: "Add the devices that the NSO provider has and Crosswork is missing to Crosswork",
				Computed:            true,
			},
			"onboard_from_nso_rule": schema.StringAttribute{
				MarkdownDescription: "Expression that selects the devices that are added to Crosswork",
				Computed:            true,
			},
			"sync_from_nso": schema.BoolAttribute{
				MarkdownDescription: "Run a sync-from on the NSO provider after the devices are onboarded",
				Computed:            true,
			},
			"sync_from_nso_rule": schema.StringAttribute{
				MarkdownDescription: "Expression that selects the devices that are synchronised",
				Computed:            true,
			},
			"check_sync": schema.BoolAttribute{
				MarkdownDescription: "Check the sync status of the devices on the NSO provider",
				Computed:            true,
			},
			"check_sync_rule": schema.StringAttribute{
				MarkdownDescription: "Expression that selects the devices whose sync status is checked",
				Computed:            true,
			},
			"neds": schema.ListNestedAttribute{
				MarkdownDescription: "Rules that select the NED an onboarded device gets. Both the NED and the rule are required.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ned": schema.StringAttribute{
							MarkdownDescription: "Name of the NED, for example `cisco-iosxr-cli-7.76`",
							Computed:            true,
						},
						"rule": schema.StringAttribute{
							MarkdownDescription: "Expression that selects the devices that get the NED, for example `product_info.software_type = 'IOS XR*'`",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *NSOProviderPoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *NSOProviderPoliciesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *NSOProviderPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config NSOProviderPolicies

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
	res, err := config.readGet(&c, "/crosswork/inventory/v1/nso/policy/query"+params)
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

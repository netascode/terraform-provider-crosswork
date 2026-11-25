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
	_ datasource.DataSource              = &LocalCongestionMitigationInterfaceThresholdsDataSource{}
	_ datasource.DataSourceWithConfigure = &LocalCongestionMitigationInterfaceThresholdsDataSource{}
)

func NewLocalCongestionMitigationInterfaceThresholdsDataSource() datasource.DataSource {
	return &LocalCongestionMitigationInterfaceThresholdsDataSource{}
}

type LocalCongestionMitigationInterfaceThresholdsDataSource struct {
	clients map[string]*cw.Client
}

func (d *LocalCongestionMitigationInterfaceThresholdsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_congestion_mitigation_interface_thresholds"
}

func (d *LocalCongestionMitigationInterfaceThresholdsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the per-interface utilization thresholds of one LCM domain.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Required:            true,
			},
			"domain_id": schema.StringAttribute{
				MarkdownDescription: "LCM domain identifier",
				Computed:            true,
			},
			"interface_thresholds": schema.SetNestedAttribute{
				MarkdownDescription: "The interfaces with their own utilization threshold; an empty list removes them all. Only interfaces of the IGP links of the domain are accepted, by the full interface name.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"node": schema.StringAttribute{
							MarkdownDescription: "Host name of the device that has the interface",
							Computed:            true,
						},
						"interface": schema.StringAttribute{
							MarkdownDescription: "Full name of the interface, of an IGP link of the domain",
							Computed:            true,
						},
						"utilization_threshold": schema.Int64Attribute{
							MarkdownDescription: "Utilization percentage at which the interface is considered congested; not below the hold margin of the domain",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *LocalCongestionMitigationInterfaceThresholdsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *LocalCongestionMitigationInterfaceThresholdsDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *LocalCongestionMitigationInterfaceThresholdsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config LocalCongestionMitigationInterfaceThresholds

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
	res, err := config.readGet(&c, "/crosswork/optima/api/v1/functionPacks/lcm/domain/interfaceThresholds"+params)
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

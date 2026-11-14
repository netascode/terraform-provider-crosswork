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
	_ datasource.DataSource              = &L2VPNVPNServicesDataSource{}
	_ datasource.DataSourceWithConfigure = &L2VPNVPNServicesDataSource{}
)

func NewL2VPNVPNServicesDataSource() datasource.DataSource {
	return &L2VPNVPNServicesDataSource{}
}

type L2VPNVPNServicesDataSource struct {
	clients map[string]*cw.Client
}

func (d *L2VPNVPNServicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2vpn_vpn_services"
}

func (d *L2VPNVPNServicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read an L2VPN service.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Required:            true,
			},
			"vpn_id": schema.StringAttribute{
				MarkdownDescription: "Name of the service. Letters, digits, `_` and `-`, at most 32 characters.",
				Required:            true,
			},
			"vpn_type": schema.StringAttribute{
				MarkdownDescription: "Type of the service",
				Computed:            true,
			},
			"evi_id": schema.Int64Attribute{
				MarkdownDescription: "EVPN instance id. Use `auto_evi_id` for a number from the pool.",
				Computed:            true,
			},
			"auto_evi_id": schema.BoolAttribute{
				MarkdownDescription: "Takes the EVPN instance id from the pool of the L2VPN id pools. Set to `true` or leave unset; not with `evi_id`.",
				Computed:            true,
			},
			"evi_source": schema.Int64Attribute{
				MarkdownDescription: "EVPN VPWS source id",
				Computed:            true,
			},
			"auto_evi_source": schema.BoolAttribute{
				MarkdownDescription: "Takes the EVPN VPWS source id from the pool of the L2VPN id pools. Set to `true` or leave unset; not with `evi_source`.",
				Computed:            true,
			},
			"evi_target": schema.Int64Attribute{
				MarkdownDescription: "EVPN VPWS target id",
				Computed:            true,
			},
			"auto_evi_target": schema.BoolAttribute{
				MarkdownDescription: "Takes the EVPN VPWS target id from the pool of the L2VPN id pools. Set to `true` or leave unset; not with `evi_target`.",
				Computed:            true,
			},
			"service_assurance_profile_name": schema.StringAttribute{
				MarkdownDescription: "Service assurance profile, which has to exist in Crosswork (the name of a system profile ends in a space and `system`)",
				Computed:            true,
			},
			"service_assurance_rule_name": schema.StringAttribute{
				MarkdownDescription: "Service assurance rule, which has to exist in Crosswork",
				Computed:            true,
			},
			"service_assurance_monitoring_state": schema.StringAttribute{
				MarkdownDescription: "Monitoring state. `pause` can be set only while the state is `enable`.",
				Computed:            true,
			},
			"service_assurance_preservation": schema.StringAttribute{
				MarkdownDescription: "What happens to the monitoring data when the service is removed",
				Computed:            true,
			},
			"vpn_nodes": schema.ListNestedAttribute{
				MarkdownDescription: "The nodes of the service. A service needs two nodes, and NSO refuses one with fewer. The example shows only one node.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpn_node_id": schema.StringAttribute{
							MarkdownDescription: "Name of a device in the Crosswork inventory that is onboarded to NSO",
							Computed:            true,
						},
						"multi_home_ethernet_segment_identifier": schema.StringAttribute{
							MarkdownDescription: "Ethernet segment identifier of the access, nine bytes in hexadecimal separated by dots. The device gets it as an EVPN ethernet segment of the interface.",
							Computed:            true,
						},
						"sr_policy": schema.StringAttribute{
							MarkdownDescription: "Name of an SR policy that carries the traffic.",
							Computed:            true,
						},
						"sr_policy_fallback": schema.StringAttribute{
							MarkdownDescription: "What the pseudowire does when the policy is down",
							Computed:            true,
						},
						"odn_route_policies": schema.SetAttribute{
							MarkdownDescription: "L2VPN route policies of the on-demand mapping.",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"vpn_network_accesses": schema.ListNestedAttribute{
							MarkdownDescription: "Accesses of the node.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: "Identifier of the access",
										Computed:            true,
									},
									"description": schema.StringAttribute{
										MarkdownDescription: "Description of the access",
										Computed:            true,
									},
									"interface_id": schema.StringAttribute{
										MarkdownDescription: "Interface of the access. A dot1q access becomes a sub-interface.",
										Computed:            true,
									},
									"encapsulation_type": schema.StringAttribute{
										MarkdownDescription: "Encapsulation type, the same on both nodes",
										Computed:            true,
									},
									"cvlan_id": schema.Int64Attribute{
										MarkdownDescription: "VLAN of the dot1q sub-interface",
										Computed:            true,
									},
									"tag_operation_pop": schema.StringAttribute{
										MarkdownDescription: "Number of tags that the access removes. Only one of pop, push and translate.",
										Computed:            true,
									},
									"tag_operation_push": schema.BoolAttribute{
										MarkdownDescription: "The access pushes the tag of `tag_operation_tag`. Set to `true` or leave unset.",
										Computed:            true,
									},
									"tag_operation_translate": schema.StringAttribute{
										MarkdownDescription: "The access translates the tag into the tag of `tag_operation_tag`",
										Computed:            true,
									},
									"tag_operation_tag": schema.Int64Attribute{
										MarkdownDescription: "Tag that push or translate uses",
										Computed:            true,
									},
									"mtu": schema.Int64Attribute{
										MarkdownDescription: "MTU of the access.",
										Computed:            true,
									},
									"qos_profiles": schema.ListNestedAttribute{
										MarkdownDescription: "QoS profiles of the access",
										Computed:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"profile": schema.StringAttribute{
													MarkdownDescription: "Name of a QoS profile id; the device needs a policy-map of that name",
													Computed:            true,
												},
												"direction": schema.StringAttribute{
													MarkdownDescription: "wan-to-site is the input and site-to-wan the output of the interface",
													Computed:            true,
												},
											},
										},
									},
									"oam_md_name": schema.StringAttribute{
										MarkdownDescription: "Maintenance domain name of the Ethernet OAM. Needs the other OAM attributes, on both nodes.",
										Computed:            true,
									},
									"oam_md_level": schema.Int64Attribute{
										MarkdownDescription: "Maintenance domain level",
										Computed:            true,
									},
									"oam_maid": schema.StringAttribute{
										MarkdownDescription: "Maintenance association id, equal on both nodes",
										Computed:            true,
									},
									"oam_mep_id": schema.Int64Attribute{
										MarkdownDescription: "Maintenance end point id",
										Computed:            true,
									},
									"oam_message_period": schema.StringAttribute{
										MarkdownDescription: "Period of the continuity check messages.",
										Computed:            true,
									},
									"oam_id_type": schema.StringAttribute{
										MarkdownDescription: "Kind of the maintenance association id",
										Computed:            true,
									},
									"oam_y1731_profiles": schema.ListNestedAttribute{
										MarkdownDescription: "Y.1731 profiles of the access, at most one delay and one loss profile",
										Computed:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													MarkdownDescription: "Name of a Y.1731 profile",
													Computed:            true,
												},
												"schedule_interval": schema.Int64Attribute{
													MarkdownDescription: "Minutes between the operations; replaces the interval of the profile",
													Computed:            true,
												},
												"schedule_duration": schema.StringAttribute{
													MarkdownDescription: "Minutes that the operations last, 1 to 1440, or `forever`; replaces the duration of the profile",
													Computed:            true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *L2VPNVPNServicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *L2VPNVPNServicesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *L2VPNVPNServicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config L2VPNVPNServices

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
	params += "?content=config"
	res, err := c.Get(config.getPath() + params)
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

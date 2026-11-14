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
	_ datasource.DataSource              = &L3VPNVPNServicesDataSource{}
	_ datasource.DataSourceWithConfigure = &L3VPNVPNServicesDataSource{}
)

func NewL3VPNVPNServicesDataSource() datasource.DataSource {
	return &L3VPNVPNServicesDataSource{}
}

type L3VPNVPNServicesDataSource struct {
	clients map[string]*cw.Client
}

func (d *L3VPNVPNServicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l3vpn_vpn_services"
}

func (d *L3VPNVPNServicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read an L3VPN service.",

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
				MarkdownDescription: "Name of the service. Letters, digits and `-`, at most 32 characters.",
				Required:            true,
			},
			"vpn_service_topology": schema.StringAttribute{
				MarkdownDescription: "Topology of the service",
				Computed:            true,
			},
			"auto_hub_rt": schema.BoolAttribute{
				MarkdownDescription: "Marks the hub route target as taken from the NSO pool. Needs the topology any-to-any or hub-spoke. Set to `true` or leave unset.",
				Computed:            true,
			},
			"auto_spoke_rt": schema.BoolAttribute{
				MarkdownDescription: "Marks the spoke route target as taken from the NSO pool. Needs the topology hub-spoke. Set to `true` or leave unset.",
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
			"vpn_instance_profiles": schema.ListNestedAttribute{
				MarkdownDescription: "Profiles of the service",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"profile_id": schema.StringAttribute{
							MarkdownDescription: "Name of the profile",
							Computed:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "Role of the profile in a hub and spoke service",
							Computed:            true,
						},
						"rd": schema.StringAttribute{
							MarkdownDescription: "Route distinguisher as type:administrator:value. NSO uses rd auto when it is left out. Changing it replaces the service.",
							Computed:            true,
						},
						"local_as": schema.Int64Attribute{
							MarkdownDescription: "Autonomous system of the profile",
							Computed:            true,
						},
						"address_families": schema.ListNestedAttribute{
							MarkdownDescription: "Address families of the profile",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"address_family": schema.StringAttribute{
										MarkdownDescription: "Address family",
										Computed:            true,
									},
									"vpn_targets": schema.ListNestedAttribute{
										MarkdownDescription: "Route targets of the family",
										Computed:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"id": schema.Int64Attribute{
													MarkdownDescription: "Number of the group of route targets",
													Computed:            true,
												},
												"route_target_type": schema.StringAttribute{
													MarkdownDescription: "Whether the targets are imported, exported or both",
													Computed:            true,
												},
												"route_targets": schema.ListNestedAttribute{
													MarkdownDescription: "Route targets of the group",
													Computed:            true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"route_target": schema.StringAttribute{
																MarkdownDescription: "Route target as type:administrator:value",
																Computed:            true,
															},
														},
													},
												},
											},
										},
									},
									"import_policy": schema.StringAttribute{
										MarkdownDescription: "Name of the route policy that is applied to the routes that are imported",
										Computed:            true,
									},
									"export_policy": schema.StringAttribute{
										MarkdownDescription: "Name of the route policy that is applied to the routes that are exported",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
			"vpn_nodes": schema.ListNestedAttribute{
				MarkdownDescription: "Nodes of the service.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpn_node_id": schema.StringAttribute{
							MarkdownDescription: "Name of a device in the Crosswork inventory that is onboarded to NSO",
							Computed:            true,
						},
						"local_as": schema.Int64Attribute{
							MarkdownDescription: "Autonomous system of the node",
							Computed:            true,
						},
						"active_vpn_instance_profiles": schema.ListNestedAttribute{
							MarkdownDescription: "Profiles that the node uses, with the values that replace the profile on this node",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"profile_id": schema.StringAttribute{
										MarkdownDescription: "Name of a profile of the service",
										Computed:            true,
									},
									"rd": schema.StringAttribute{
										MarkdownDescription: "Route distinguisher that replaces the one of the profile on this node. Changing it replaces the service.",
										Computed:            true,
									},
									"address_families": schema.ListNestedAttribute{
										MarkdownDescription: "Address families that replace those of the profile on this node",
										Computed:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"address_family": schema.StringAttribute{
													MarkdownDescription: "Address family",
													Computed:            true,
												},
												"vpn_targets": schema.ListNestedAttribute{
													MarkdownDescription: "Route targets of the family",
													Computed:            true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"id": schema.Int64Attribute{
																MarkdownDescription: "Number of the group of route targets",
																Computed:            true,
															},
															"route_target_type": schema.StringAttribute{
																MarkdownDescription: "Whether the targets are imported, exported or both",
																Computed:            true,
															},
															"route_targets": schema.ListNestedAttribute{
																MarkdownDescription: "Route targets of the group",
																Computed:            true,
																NestedObject: schema.NestedAttributeObject{
																	Attributes: map[string]schema.Attribute{
																		"route_target": schema.StringAttribute{
																			MarkdownDescription: "Route target as type:administrator:value",
																			Computed:            true,
																		},
																	},
																},
															},
														},
													},
												},
												"import_policy": schema.StringAttribute{
													MarkdownDescription: "Name of the route policy that is applied to the routes that are imported",
													Computed:            true,
												},
												"export_policy": schema.StringAttribute{
													MarkdownDescription: "Name of the route policy that is applied to the routes that are exported",
													Computed:            true,
												},
											},
										},
									},
								},
							},
						},
						"vpn_network_accesses": schema.ListNestedAttribute{
							MarkdownDescription: "Accesses of the node",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: "Identifier of the access",
										Computed:            true,
									},
									"interface_id": schema.StringAttribute{
										MarkdownDescription: "Interface of the access. A dot1q access becomes a sub-interface.",
										Computed:            true,
									},
									"encapsulation_type": schema.StringAttribute{
										MarkdownDescription: "Encapsulation type",
										Computed:            true,
									},
									"cvlan_id": schema.Int64Attribute{
										MarkdownDescription: "VLAN of the dot1q sub-interface",
										Computed:            true,
									},
									"ipv4_address": schema.StringAttribute{
										MarkdownDescription: "IPv4 address of the access",
										Computed:            true,
									},
									"ipv4_prefix_length": schema.Int64Attribute{
										MarkdownDescription: "Prefix length of the IPv4 address",
										Computed:            true,
									},
									"ipv6_address": schema.StringAttribute{
										MarkdownDescription: "IPv6 address of the access",
										Computed:            true,
									},
									"ipv6_prefix_length": schema.Int64Attribute{
										MarkdownDescription: "Prefix length of the IPv6 address",
										Computed:            true,
									},
									"mtu": schema.Int64Attribute{
										MarkdownDescription: "MTU of the access. NSO refuses it on a Loopback.",
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
									"routing_protocols": schema.ListNestedAttribute{
										MarkdownDescription: "Routing protocols of the access. A route policy of the profile is pushed only when there is one.",
										Computed:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													MarkdownDescription: "Identifier of the routing protocol",
													Computed:            true,
												},
												"type": schema.StringAttribute{
													MarkdownDescription: "Routing protocol",
													Computed:            true,
												},
												"bgp_peer_as": schema.Int64Attribute{
													MarkdownDescription: "Autonomous system of the BGP neighbors",
													Computed:            true,
												},
												"bgp_address_family": schema.StringAttribute{
													MarkdownDescription: "Address family of the BGP neighbors",
													Computed:            true,
												},
												"bgp_neighbors": schema.SetAttribute{
													MarkdownDescription: "BGP neighbors, IPv4 or IPv6 addresses",
													ElementType:         types.StringType,
													Computed:            true,
												},
												"bgp_multihop": schema.Int64Attribute{
													MarkdownDescription: "Number of hops to the BGP neighbors",
													Computed:            true,
												},
												"bgp_redistribute_connected": schema.ListNestedAttribute{
													MarkdownDescription: "Connected routes that BGP redistributes. Only an entry with enable set to true is pushed to the device.",
													Computed:            true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"address_family": schema.StringAttribute{
																MarkdownDescription: "Address family",
																Computed:            true,
															},
															"enable": schema.BoolAttribute{
																MarkdownDescription: "Redistribute the connected routes of the family",
																Computed:            true,
															},
															"metric": schema.Int64Attribute{
																MarkdownDescription: "Metric of the redistributed routes",
																Computed:            true,
															},
														},
													},
												},
												"bgp_update_source_if_type": schema.StringAttribute{
													MarkdownDescription: "Interface type of the BGP update source",
													Computed:            true,
												},
												"bgp_update_source_if_id": schema.StringAttribute{
													MarkdownDescription: "Interface number of the BGP update source",
													Computed:            true,
												},
												"bgp_update_source_sub_if_id": schema.Int64Attribute{
													MarkdownDescription: "Sub-interface number of the BGP update source",
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

func (d *L3VPNVPNServicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *L3VPNVPNServicesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *L3VPNVPNServicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config L3VPNVPNServices

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

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
	"strings"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &L3VPNVPNServicesResource{}
var _ resource.ResourceWithImportState = &L3VPNVPNServicesResource{}

func NewL3VPNVPNServicesResource() resource.Resource {
	return &L3VPNVPNServicesResource{}
}

type L3VPNVPNServicesResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
	deviceReady           nsoDeviceConfig
}

func (r *L3VPNVPNServicesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l3vpn_vpn_services"
}

func (r *L3VPNVPNServicesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage an L3VPN service. Create and update wait until the devices named in `vpn_nodes.vpn_node_id` are ready in Crosswork, as set by the `device_ready_mode` and `device_ready_timeout` arguments of the provider configuration.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"vpn_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the service. Letters, digits and `-`, at most 32 characters.").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpn_service_topology": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Topology of the service").AddStringEnumDescription("ietf-vpn-common:any-to-any", "ietf-vpn-common:hub-spoke", "ietf-vpn-common:custom").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ietf-vpn-common:any-to-any", "ietf-vpn-common:hub-spoke", "ietf-vpn-common:custom"),
				},
			},
			"auto_hub_rt": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Marks the hub route target as taken from the NSO pool. Needs the topology any-to-any or hub-spoke. Set to `true` or leave unset.").String,
				Optional:            true,
			},
			"auto_spoke_rt": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Marks the spoke route target as taken from the NSO pool. Needs the topology hub-spoke. Set to `true` or leave unset.").String,
				Optional:            true,
			},
			"service_assurance_profile_name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Service assurance profile, which has to exist in Crosswork (the name of a system profile ends in a space and `system`)").String,
				Optional:            true,
			},
			"service_assurance_rule_name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Service assurance rule, which has to exist in Crosswork").String,
				Optional:            true,
			},
			"service_assurance_monitoring_state": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Monitoring state. `pause` can be set only while the state is `enable`.").AddStringEnumDescription("enable", "disable", "pause").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("enable", "disable", "pause"),
				},
			},
			"service_assurance_preservation": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("What happens to the monitoring data when the service is removed").AddStringEnumDescription("remove", "preserve").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("remove", "preserve"),
				},
			},
			"vpn_instance_profiles": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Profiles of the service").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"profile_id": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of the profile").String,
							Required:            true,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Role of the profile in a hub and spoke service").AddStringEnumDescription("ietf-vpn-common:hub-role", "ietf-vpn-common:spoke-role").String,
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("ietf-vpn-common:hub-role", "ietf-vpn-common:spoke-role"),
							},
						},
						"rd": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Route distinguisher as type:administrator:value. NSO uses rd auto when it is left out. Changing it replaces the service.").String,
							Optional:            true,
						},
						"local_as": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Autonomous system of the profile").AddIntegerRangeDescription(1, 4294967295).String,
							Optional:            true,
							Validators: []validator.Int64{
								int64validator.Between(1, 4294967295),
							},
						},
						"address_families": schema.ListNestedAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Address families of the profile").String,
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"address_family": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Address family").AddStringEnumDescription("ietf-vpn-common:ipv4", "ietf-vpn-common:ipv6").String,
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("ietf-vpn-common:ipv4", "ietf-vpn-common:ipv6"),
										},
									},
									"vpn_targets": schema.ListNestedAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Route targets of the family").String,
										Optional:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"id": schema.Int64Attribute{
													MarkdownDescription: helpers.NewAttributeDescription("Number of the group of route targets").AddIntegerRangeDescription(1, 255).String,
													Required:            true,
													Validators: []validator.Int64{
														int64validator.Between(1, 255),
													},
												},
												"route_target_type": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Whether the targets are imported, exported or both").AddStringEnumDescription("both", "import", "export").String,
													Optional:            true,
													Validators: []validator.String{
														stringvalidator.OneOf("both", "import", "export"),
													},
												},
												"route_targets": schema.ListNestedAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Route targets of the group").String,
													Optional:            true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"route_target": schema.StringAttribute{
																MarkdownDescription: helpers.NewAttributeDescription("Route target as type:administrator:value").String,
																Required:            true,
															},
														},
													},
												},
											},
										},
									},
									"import_policy": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Name of the route policy that is applied to the routes that are imported").String,
										Optional:            true,
									},
									"export_policy": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Name of the route policy that is applied to the routes that are exported").String,
										Optional:            true,
									},
								},
							},
						},
					},
				},
			},
			"vpn_nodes": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Nodes of the service.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpn_node_id": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of a device in the Crosswork inventory that is onboarded to NSO").String,
							Required:            true,
						},
						"local_as": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Autonomous system of the node").AddIntegerRangeDescription(1, 4294967295).String,
							Optional:            true,
							Validators: []validator.Int64{
								int64validator.Between(1, 4294967295),
							},
						},
						"active_vpn_instance_profiles": schema.ListNestedAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Profiles that the node uses, with the values that replace the profile on this node").String,
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"profile_id": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Name of a profile of the service").String,
										Required:            true,
									},
									"rd": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Route distinguisher that replaces the one of the profile on this node. Changing it replaces the service.").String,
										Optional:            true,
									},
									"address_families": schema.ListNestedAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Address families that replace those of the profile on this node").String,
										Optional:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"address_family": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Address family").AddStringEnumDescription("ietf-vpn-common:ipv4", "ietf-vpn-common:ipv6").String,
													Required:            true,
													Validators: []validator.String{
														stringvalidator.OneOf("ietf-vpn-common:ipv4", "ietf-vpn-common:ipv6"),
													},
												},
												"vpn_targets": schema.ListNestedAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Route targets of the family").String,
													Optional:            true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"id": schema.Int64Attribute{
																MarkdownDescription: helpers.NewAttributeDescription("Number of the group of route targets").AddIntegerRangeDescription(1, 255).String,
																Required:            true,
																Validators: []validator.Int64{
																	int64validator.Between(1, 255),
																},
															},
															"route_target_type": schema.StringAttribute{
																MarkdownDescription: helpers.NewAttributeDescription("Whether the targets are imported, exported or both").AddStringEnumDescription("both", "import", "export").String,
																Optional:            true,
																Validators: []validator.String{
																	stringvalidator.OneOf("both", "import", "export"),
																},
															},
															"route_targets": schema.ListNestedAttribute{
																MarkdownDescription: helpers.NewAttributeDescription("Route targets of the group").String,
																Optional:            true,
																NestedObject: schema.NestedAttributeObject{
																	Attributes: map[string]schema.Attribute{
																		"route_target": schema.StringAttribute{
																			MarkdownDescription: helpers.NewAttributeDescription("Route target as type:administrator:value").String,
																			Required:            true,
																		},
																	},
																},
															},
														},
													},
												},
												"import_policy": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Name of the route policy that is applied to the routes that are imported").String,
													Optional:            true,
												},
												"export_policy": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Name of the route policy that is applied to the routes that are exported").String,
													Optional:            true,
												},
											},
										},
									},
								},
							},
						},
						"vpn_network_accesses": schema.ListNestedAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Accesses of the node").String,
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Identifier of the access").String,
										Required:            true,
									},
									"interface_id": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Interface of the access. A dot1q access becomes a sub-interface.").String,
										Required:            true,
									},
									"encapsulation_type": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Encapsulation type").AddStringEnumDescription("ietf-vpn-common:dot1q").String,
										Optional:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("ietf-vpn-common:dot1q"),
										},
									},
									"cvlan_id": schema.Int64Attribute{
										MarkdownDescription: helpers.NewAttributeDescription("VLAN of the dot1q sub-interface").AddIntegerRangeDescription(1, 4094).String,
										Optional:            true,
										Validators: []validator.Int64{
											int64validator.Between(1, 4094),
										},
									},
									"ipv4_address": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("IPv4 address of the access").String,
										Optional:            true,
									},
									"ipv4_prefix_length": schema.Int64Attribute{
										MarkdownDescription: helpers.NewAttributeDescription("Prefix length of the IPv4 address").AddIntegerRangeDescription(1, 32).String,
										Optional:            true,
										Validators: []validator.Int64{
											int64validator.Between(1, 32),
										},
									},
									"ipv6_address": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("IPv6 address of the access").String,
										Optional:            true,
									},
									"ipv6_prefix_length": schema.Int64Attribute{
										MarkdownDescription: helpers.NewAttributeDescription("Prefix length of the IPv6 address").AddIntegerRangeDescription(1, 128).String,
										Optional:            true,
										Validators: []validator.Int64{
											int64validator.Between(1, 128),
										},
									},
									"mtu": schema.Int64Attribute{
										MarkdownDescription: helpers.NewAttributeDescription("MTU of the access. NSO refuses it on a Loopback.").AddIntegerRangeDescription(64, 65535).String,
										Optional:            true,
										Validators: []validator.Int64{
											int64validator.Between(64, 65535),
										},
									},
									"qos_profiles": schema.ListNestedAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("QoS profiles of the access").String,
										Optional:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"profile": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Name of a QoS profile id; the device needs a policy-map of that name").String,
													Required:            true,
												},
												"direction": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("wan-to-site is the input and site-to-wan the output of the interface").AddStringEnumDescription("ietf-vpn-common:both", "ietf-vpn-common:site-to-wan", "ietf-vpn-common:wan-to-site").String,
													Optional:            true,
													Validators: []validator.String{
														stringvalidator.OneOf("ietf-vpn-common:both", "ietf-vpn-common:site-to-wan", "ietf-vpn-common:wan-to-site"),
													},
												},
											},
										},
									},
									"routing_protocols": schema.ListNestedAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Routing protocols of the access. A route policy of the profile is pushed only when there is one.").String,
										Optional:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Identifier of the routing protocol").String,
													Required:            true,
												},
												"type": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Routing protocol").AddStringEnumDescription("ietf-vpn-common:bgp-routing", "ietf-vpn-common:direct-routing").String,
													Required:            true,
													Validators: []validator.String{
														stringvalidator.OneOf("ietf-vpn-common:bgp-routing", "ietf-vpn-common:direct-routing"),
													},
												},
												"bgp_peer_as": schema.Int64Attribute{
													MarkdownDescription: helpers.NewAttributeDescription("Autonomous system of the BGP neighbors").AddIntegerRangeDescription(1, 4294967295).String,
													Optional:            true,
													Validators: []validator.Int64{
														int64validator.Between(1, 4294967295),
													},
												},
												"bgp_address_family": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Address family of the BGP neighbors").AddStringEnumDescription("ipv4", "ipv6", "dual-stack").String,
													Optional:            true,
													Validators: []validator.String{
														stringvalidator.OneOf("ipv4", "ipv6", "dual-stack"),
													},
												},
												"bgp_neighbors": schema.SetAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("BGP neighbors, IPv4 or IPv6 addresses").String,
													ElementType:         types.StringType,
													Optional:            true,
												},
												"bgp_multihop": schema.Int64Attribute{
													MarkdownDescription: helpers.NewAttributeDescription("Number of hops to the BGP neighbors").AddIntegerRangeDescription(1, 255).String,
													Optional:            true,
													Validators: []validator.Int64{
														int64validator.Between(1, 255),
													},
												},
												"bgp_redistribute_connected": schema.ListNestedAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Connected routes that BGP redistributes. Only an entry with enable set to true is pushed to the device.").String,
													Optional:            true,
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															"address_family": schema.StringAttribute{
																MarkdownDescription: helpers.NewAttributeDescription("Address family").AddStringEnumDescription("ietf-vpn-common:ipv4", "ietf-vpn-common:ipv6").String,
																Required:            true,
																Validators: []validator.String{
																	stringvalidator.OneOf("ietf-vpn-common:ipv4", "ietf-vpn-common:ipv6"),
																},
															},
															"enable": schema.BoolAttribute{
																MarkdownDescription: helpers.NewAttributeDescription("Redistribute the connected routes of the family").String,
																Optional:            true,
															},
															"metric": schema.Int64Attribute{
																MarkdownDescription: helpers.NewAttributeDescription("Metric of the redistributed routes").AddIntegerRangeDescription(0, 4294967295).String,
																Optional:            true,
																Validators: []validator.Int64{
																	int64validator.Between(0, 4294967295),
																},
															},
														},
													},
												},
												"bgp_update_source_if_type": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Interface type of the BGP update source").AddStringEnumDescription("Loopback", "HundredGigE", "FourHundredGigE", "TenGigE", "GigabitEthernet", "Bundle-Ether", "TwentyFiveGigE", "FortyGigE").String,
													Optional:            true,
													Validators: []validator.String{
														stringvalidator.OneOf("Loopback", "HundredGigE", "FourHundredGigE", "TenGigE", "GigabitEthernet", "Bundle-Ether", "TwentyFiveGigE", "FortyGigE"),
													},
												},
												"bgp_update_source_if_id": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Interface number of the BGP update source").String,
													Optional:            true,
												},
												"bgp_update_source_sub_if_id": schema.Int64Attribute{
													MarkdownDescription: helpers.NewAttributeDescription("Sub-interface number of the BGP update source").AddIntegerRangeDescription(0, 2147483647).String,
													Optional:            true,
													Validators: []validator.Int64{
														int64validator.Between(0, 2147483647),
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
		},
	}
}

func (r *L3VPNVPNServicesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
	r.deviceReady = req.ProviderData.(*CwProviderData).DeviceReady
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *L3VPNVPNServicesResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create
func (r *L3VPNVPNServicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan L3VPNVPNServices

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	if !waitForNSODevices(ctx, &c, plan.Instance.ValueString(), r.deviceReady, plan.nsoDevices(), &resp.Diagnostics) {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, L3VPNVPNServices{})

	params := ""
	res, err := c.Post("/crosswork/proxy/nso/restconf/data/ietf-l3vpn-ntw:l3vpn-ntw/vpn-services"+params, body, cw.PlanPath(plan.getPathPlan()))
	if err != nil && (cw.IsPlanTimeout(err) || cw.IsPlanFailed(err)) {
		// NSO accepted the object, but its plan did not complete. Keep the object in the state, which the error marks
		// as tainted, so that the next apply replaces or removes it instead of failing with "object already exists".
		plan.Id = types.StringValue(fmt.Sprint(plan.VpnId.ValueString()))
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("The object was created, but its plan did not complete, so it is tracked as tainted: %s", err))
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "POST", err, res.String()))
		return
	}
	plan.Id = types.StringValue(fmt.Sprint(plan.VpnId.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end create

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (r *L3VPNVPNServicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state L3VPNVPNServices

	// Read state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", state.Id.String()))

	params := ""
	params += "?content=config"
	res, err := c.Get(state.getPath() + params)
	if err != nil && (strings.Contains(err.Error(), "StatusCode 404") || strings.Contains(err.Error(), "StatusCode 406") || strings.Contains(err.Error(), "StatusCode 500") || strings.Contains(err.Error(), "StatusCode 400") || cw.IsNotFound(err)) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
		return
	}

	// If every attribute is set to null we are dealing with an import operation and therefore reading all attributes
	if state.isNull(ctx, res) {
		state.fromBody(ctx, res)
	} else {
		state.updateFromBody(ctx, res)
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

// Section below is generated&owned by "gen/generator.go". //template:begin update
func (r *L3VPNVPNServicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state L3VPNVPNServices

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read state
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	if !waitForNSODevices(ctx, &c, plan.Instance.ValueString(), r.deviceReady, nsoDevicesAdded(plan.nsoDevices(), state.nsoDevices()), &resp.Diagnostics) {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))

	body := plan.toBody(ctx, state)
	params := ""
	res, err := c.Put(plan.getPath()+params, body, cw.PlanPath(plan.getPathPlan()))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s, %s", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end update

// Section below is generated&owned by "gen/generator.go". //template:begin delete
func (r *L3VPNVPNServicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state L3VPNVPNServices

	// Read state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))
	res, err := c.Delete(state.getPath(), cw.PlanPath(state.getPathPlan()))
	if err != nil && !strings.Contains(err.Error(), "StatusCode 404") && !cw.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s, %s", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Section below is generated&owned by "gen/generator.go". //template:begin import
func (r *L3VPNVPNServicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <vpn_id>. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vpn_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *L3VPNVPNServicesResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state L3VPNVPNServices, params string) (cw.Res, error) {
	var err error
	cacheKey := "L3VPNVPNServices::" + state.Instance.ValueString() + "::"

	_, cacheSuffix, found := strings.Cut(params, "?")
	queryPart, err := url.ParseQuery(cacheSuffix)
	if err == nil {
		delete(queryPart, "id")
		newQuery := queryPart.Encode()
		cacheSuffix = "?" + newQuery
		cacheKey += cacheSuffix
	}

	cachedValue, found := r.cache.Get(cacheKey)
	if found {
		tflog.Debug(ctx, fmt.Sprintf("hit cache for %s", cacheKey))
		cwRes, ok := cachedValue.(cw.Res)
		if ok {
			return cwRes, nil
		}
		tflog.Info(ctx, fmt.Sprintf("Invalid cache entry type for %s", cacheKey))
		r.cache.Delete(cacheKey)
	}
	res, err := c.Get(state.getPath() + params)
	singleRes := res
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}

// End of section. //template:end readcache

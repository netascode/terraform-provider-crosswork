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
var _ resource.Resource = &L2VPNVPNServicesResource{}
var _ resource.ResourceWithImportState = &L2VPNVPNServicesResource{}

func NewL2VPNVPNServicesResource() resource.Resource {
	return &L2VPNVPNServicesResource{}
}

type L2VPNVPNServicesResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
	deviceReady           nsoDeviceConfig
}

func (r *L2VPNVPNServicesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2vpn_vpn_services"
}

func (r *L2VPNVPNServicesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage an L2VPN service. Create and update wait until the devices named in `vpn_nodes.vpn_node_id` are ready in Crosswork, as set by the `device_ready_mode` and `device_ready_timeout` arguments of the provider configuration.").String,

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
				MarkdownDescription: helpers.NewAttributeDescription("Name of the service. Letters, digits, `_` and `-`, at most 32 characters.").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpn_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Type of the service").AddStringEnumDescription("ietf-vpn-common:vpws-evpn").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ietf-vpn-common:vpws-evpn"),
				},
			},
			"evi_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("EVPN instance id. Use `auto_evi_id` for a number from the pool.").AddIntegerRangeDescription(1, 65534).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65534),
				},
			},
			"auto_evi_id": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Takes the EVPN instance id from the pool of the L2VPN id pools. Set to `true` or leave unset; not with `evi_id`.").String,
				Optional:            true,
			},
			"evi_source": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("EVPN VPWS source id").AddIntegerRangeDescription(1, 65535).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"auto_evi_source": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Takes the EVPN VPWS source id from the pool of the L2VPN id pools. Set to `true` or leave unset; not with `evi_source`.").String,
				Optional:            true,
			},
			"evi_target": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("EVPN VPWS target id").AddIntegerRangeDescription(1, 65535).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"auto_evi_target": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Takes the EVPN VPWS target id from the pool of the L2VPN id pools. Set to `true` or leave unset; not with `evi_target`.").String,
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
			"vpn_nodes": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The nodes of the service. A service needs two nodes, and NSO refuses one with fewer. The example shows only one node.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"vpn_node_id": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of a device in the Crosswork inventory that is onboarded to NSO").String,
							Required:            true,
						},
						"multi_home_ethernet_segment_identifier": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Ethernet segment identifier of the access, nine bytes in hexadecimal separated by dots. The device gets it as an EVPN ethernet segment of the interface.").String,
							Optional:            true,
						},
						"sr_policy": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of an SR policy that carries the traffic.").String,
							Optional:            true,
						},
						"sr_policy_fallback": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("What the pseudowire does when the policy is down").AddStringEnumDescription("disable").String,
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("disable"),
							},
						},
						"odn_route_policies": schema.SetAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("L2VPN route policies of the on-demand mapping.").String,
							ElementType:         types.StringType,
							Optional:            true,
						},
						"vpn_network_accesses": schema.ListNestedAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Accesses of the node.").String,
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Identifier of the access").String,
										Required:            true,
									},
									"description": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Description of the access").String,
										Optional:            true,
									},
									"interface_id": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Interface of the access. A dot1q access becomes a sub-interface.").String,
										Required:            true,
									},
									"encapsulation_type": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Encapsulation type, the same on both nodes").AddStringEnumDescription("ietf-vpn-common:dot1q", "ietf-vpn-common:untagged-int").String,
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("ietf-vpn-common:dot1q", "ietf-vpn-common:untagged-int"),
										},
									},
									"cvlan_id": schema.Int64Attribute{
										MarkdownDescription: helpers.NewAttributeDescription("VLAN of the dot1q sub-interface").AddIntegerRangeDescription(1, 4094).String,
										Optional:            true,
										Validators: []validator.Int64{
											int64validator.Between(1, 4094),
										},
									},
									"tag_operation_pop": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Number of tags that the access removes. Only one of pop, push and translate.").AddStringEnumDescription("1", "2").String,
										Optional:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("1", "2"),
										},
									},
									"tag_operation_push": schema.BoolAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("The access pushes the tag of `tag_operation_tag`. Set to `true` or leave unset.").String,
										Optional:            true,
									},
									"tag_operation_translate": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("The access translates the tag into the tag of `tag_operation_tag`").AddStringEnumDescription("1-to-1").String,
										Optional:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("1-to-1"),
										},
									},
									"tag_operation_tag": schema.Int64Attribute{
										MarkdownDescription: helpers.NewAttributeDescription("Tag that push or translate uses").AddIntegerRangeDescription(1, 4094).String,
										Optional:            true,
										Validators: []validator.Int64{
											int64validator.Between(1, 4094),
										},
									},
									"mtu": schema.Int64Attribute{
										MarkdownDescription: helpers.NewAttributeDescription("MTU of the access.").AddIntegerRangeDescription(64, 65535).String,
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
									"oam_md_name": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Maintenance domain name of the Ethernet OAM. Needs the other OAM attributes, on both nodes.").String,
										Optional:            true,
									},
									"oam_md_level": schema.Int64Attribute{
										MarkdownDescription: helpers.NewAttributeDescription("Maintenance domain level").AddIntegerRangeDescription(0, 7).String,
										Optional:            true,
										Validators: []validator.Int64{
											int64validator.Between(0, 7),
										},
									},
									"oam_maid": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Maintenance association id, equal on both nodes").String,
										Optional:            true,
									},
									"oam_mep_id": schema.Int64Attribute{
										MarkdownDescription: helpers.NewAttributeDescription("Maintenance end point id").AddIntegerRangeDescription(1, 8191).String,
										Optional:            true,
										Validators: []validator.Int64{
											int64validator.Between(1, 8191),
										},
									},
									"oam_message_period": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Period of the continuity check messages.").AddStringEnumDescription("100ms", "1s", "10s").String,
										Optional:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("100ms", "1s", "10s"),
										},
									},
									"oam_id_type": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Kind of the maintenance association id").AddStringEnumDescription("number", "icc-based").String,
										Optional:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("number", "icc-based"),
										},
									},
									"oam_y1731_profiles": schema.ListNestedAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Y.1731 profiles of the access, at most one delay and one loss profile").String,
										Optional:            true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Name of a Y.1731 profile").String,
													Required:            true,
												},
												"schedule_interval": schema.Int64Attribute{
													MarkdownDescription: helpers.NewAttributeDescription("Minutes between the operations; replaces the interval of the profile").AddIntegerRangeDescription(1, 90).String,
													Optional:            true,
													Validators: []validator.Int64{
														int64validator.Between(1, 90),
													},
												},
												"schedule_duration": schema.StringAttribute{
													MarkdownDescription: helpers.NewAttributeDescription("Minutes that the operations last, 1 to 1440, or `forever`; replaces the duration of the profile").String,
													Optional:            true,
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

func (r *L2VPNVPNServicesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
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
func (r *L2VPNVPNServicesResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create
func (r *L2VPNVPNServicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan L2VPNVPNServices

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
	body := plan.toBody(ctx, L2VPNVPNServices{})

	params := ""
	res, err := c.Post("/crosswork/proxy/nso/restconf/data/ietf-l2vpn-ntw:l2vpn-ntw/vpn-services"+params, body, cw.PlanPath(plan.getPathPlan()))
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
func (r *L2VPNVPNServicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state L2VPNVPNServices

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
func (r *L2VPNVPNServicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state L2VPNVPNServices

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
func (r *L2VPNVPNServicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state L2VPNVPNServices

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
func (r *L2VPNVPNServicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
func (r *L2VPNVPNServicesResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state L2VPNVPNServices, params string) (cw.Res, error) {
	var err error
	cacheKey := "L2VPNVPNServices::" + state.Instance.ValueString() + "::"

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

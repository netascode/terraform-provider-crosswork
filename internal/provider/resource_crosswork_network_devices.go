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
	"regexp"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/float64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NetworkDevicesResource{}
var _ resource.ResourceWithImportState = &NetworkDevicesResource{}

func NewNetworkDevicesResource() resource.Resource {
	return &NetworkDevicesResource{}
}

type NetworkDevicesResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *NetworkDevicesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_devices"
}

func (r *NetworkDevicesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a network device in inventory.").String,

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
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Inventory id of the device. A change replaces the device").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"host_name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Host name of the device; needed when the device has no SNMP row. A device with an SNMP row reports its own host name, which replaces the value given here once Crosswork has collected the device; leave it unset to keep the host name the controller holds. A change needs the device admin down, which is done automatically").String,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"credential_profile": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the credential profile, which should hold a credential for every protocol the device uses").String,
				Required:            true,
			},
			"admin_state": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Admin state of the device. Crosswork keeps a device that is not attached to a data gateway admin down. Attaching a device to a data gateway brings it admin up. The default is `ROBOT_ADMIN_STATE_DOWN`; `ROBOT_ADMIN_STATE_UNMANAGED` cannot be used while a device is attached to a data gateway").AddStringEnumDescription("ROBOT_ADMIN_STATE_UP", "ROBOT_ADMIN_STATE_DOWN", "ROBOT_ADMIN_STATE_UNMANAGED").AddDefaultValueDescription("ROBOT_ADMIN_STATE_DOWN").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ROBOT_ADMIN_STATE_UP", "ROBOT_ADMIN_STATE_DOWN", "ROBOT_ADMIN_STATE_UNMANAGED"),
				},
				Default: stringdefault.StaticString("ROBOT_ADMIN_STATE_DOWN"),
			},
			"reachability_check": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether Crosswork checks that the device is reachable. A change needs the device admin down, which is done automatically").AddStringEnumDescription("REACH_CHECK_ENABLE", "REACH_CHECK_DISABLE").AddDefaultValueDescription("REACH_CHECK_ENABLE").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("REACH_CHECK_ENABLE", "REACH_CHECK_DISABLE"),
				},
				Default: stringdefault.StaticString("REACH_CHECK_ENABLE"),
			},
			"serial_number": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Serial number of the device").String,
				Optional:            true,
			},
			"mac_address": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("MAC address of the device").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2}$`), ""),
				},
			},
			"encryption": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP encryption of the device in the lower case form, for example `aes128-ctr`; leave it unset for none").String,
				Optional:            true,
			},
			"device_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Device type").AddStringEnumDescription("NODE_TYPE_ROUTER", "NODE_TYPE_SWITCH", "NODE_TYPE_UNKNOWN").AddDefaultValueDescription("NODE_TYPE_ROUTER").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("NODE_TYPE_ROUTER", "NODE_TYPE_SWITCH", "NODE_TYPE_UNKNOWN"),
				},
				Default: stringdefault.StaticString("NODE_TYPE_ROUTER"),
			},
			"software_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Software type, for example `IOS XR`; needed when the device has no SNMP row").String,
				Optional:            true,
			},
			"software_version": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Software version. Crosswork overwrites it after it discovers the device, so set it only for a device it cannot reach").String,
				Optional:            true,
			},
			"product_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Product type. Crosswork overwrites it after it discovers the device, so set it only for a device it cannot reach").String,
				Optional:            true,
			},
			"syslog_format": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Syslog format of the device").AddStringEnumDescription("RFC_5424", "RFC_3164", "RFC_UNKNOWN").AddDefaultValueDescription("RFC_5424").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("RFC_5424", "RFC_3164", "RFC_UNKNOWN"),
				},
				Default: stringdefault.StaticString("RFC_5424"),
			},
			"capabilities": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Capabilities of the device, any of `YANG_MDT`, `YANG_CLI`, `SNMP` and `GNMI`; `SNMP` needs an SNMP row. A change needs the device admin down, which is done automatically").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"cli_cache_enabled": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether Crosswork caches the CLI output of the device").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"provider_profile": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the NSO provider profile that manages the device").String,
				Optional:            true,
			},
			"connectivity": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Connectivity rows, one per protocol. A change to an address, a prefix length or a port needs the device admin down, which is done automatically").String,
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"protocol": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Protocol to connect to the device").AddStringEnumDescription("ROBOT_MSVC_TRANS_SSH", "ROBOT_MSVC_TRANS_SNMP", "ROBOT_MSVC_TRANS_NETCONF", "ROBOT_MSVC_TRANS_HTTP", "ROBOT_MSVC_TRANS_HTTPS", "ROBOT_MSVC_TRANS_TELNET", "ROBOT_MSVC_TRANS_GNMI", "ROBOT_MSVC_TRANS_GRPC").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("ROBOT_MSVC_TRANS_SSH", "ROBOT_MSVC_TRANS_SNMP", "ROBOT_MSVC_TRANS_NETCONF", "ROBOT_MSVC_TRANS_HTTP", "ROBOT_MSVC_TRANS_HTTPS", "ROBOT_MSVC_TRANS_TELNET", "ROBOT_MSVC_TRANS_GNMI", "ROBOT_MSVC_TRANS_GRPC"),
							},
						},
						"ip_address": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("IPv4 address of the device, unique on the Crosswork cluster").String,
							Required:            true,
						},
						"prefix_length": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Prefix length of the address").AddIntegerRangeDescription(1, 32).String,
							Required:            true,
							Validators: []validator.Int64{
								int64validator.Between(1, 32),
							},
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Port of the protocol, for example 22 for SSH and 161 for SNMP; Crosswork refuses a row without one").AddIntegerRangeDescription(1, 65535).String,
							Required:            true,
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"timeout": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Timeout in seconds").String,
							Optional:            true,
						},
						"encoding_type": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Encoding, required for GNMI and GRPC rows").AddStringEnumDescription("JSON", "BYTES", "PROTO", "ASCII", "JSON_IETF").String,
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("JSON", "BYTES", "PROTO", "ASCII", "JSON_IETF"),
							},
						},
						"trap_source_ip": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Trap source IPv4 address, for the SNMP row").String,
							Optional:            true,
						},
						"disable_trap_community_check": schema.BoolAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Whether to skip the trap community check, for the SNMP row").AddDefaultValueDescription("false").String,
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
					},
				},
			},
			"isis_system_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("ISIS system id").String,
				Optional:            true,
			},
			"ospf_router_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("OSPF router id").String,
				Optional:            true,
			},
			"te_router_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("TE router id. A change needs the device admin down, which is done automatically").String,
				Optional:            true,
			},
			"ipv6_router_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IPv6 router id").String,
				Optional:            true,
			},
			"telemetry_vrf": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("VRF of the streaming telemetry source interface").String,
				Optional:            true,
			},
			"telemetry_source_interface": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Source interface of the streaming telemetry").String,
				Optional:            true,
			},
			"telemetry_opt_out_of_mdt_config": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether to opt out of the MDT configuration").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"location_building": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Building").String,
				Optional:            true,
			},
			"location_street": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Street").String,
				Optional:            true,
			},
			"location_city": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("City").String,
				Optional:            true,
			},
			"location_state": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("State").String,
				Optional:            true,
			},
			"location_country": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Country").String,
				Optional:            true,
			},
			"location_region": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Region").String,
				Optional:            true,
			},
			"location_zip": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Zip code").String,
				Optional:            true,
			},
			"location_latitude": schema.Float64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Latitude; set it together with the longitude").AddFloatRangeDescription(-90, 90).String,
				Optional:            true,
				Validators: []validator.Float64{
					float64validator.Between(-90, 90),
				},
			},
			"location_longitude": schema.Float64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Longitude; set it together with the latitude").AddFloatRangeDescription(-180, 180).String,
				Optional:            true,
				Validators: []validator.Float64{
					float64validator.Between(-180, 180),
				},
			},
			"location_altitude": schema.Float64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Altitude").String,
				Optional:            true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("User-defined tags of the device, which must exist as tags. The tags that Crosswork adds itself are never managed").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Tag name").String,
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (r *NetworkDevicesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *NetworkDevicesResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// End of section. //template:end model

// Written by hand: a host name the configuration leaves out is the one the device reports, so it is read back after the add.
func (r *NetworkDevicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan NetworkDevices

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	body, err := linkProviderBody(&c, plan.toBody(ctx, NetworkDevices{}))
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s", "POST", err))
		return
	}

	params := ""
	res, err := c.Post(plan.getPath()+params, body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "POST", err, res.String()))
		return
	}
	plan.Id = types.StringValue(fmt.Sprint(plan.Name.ValueString()))
	plan.HostName = hostNameAfterCreate(&c, plan)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end create

func (r *NetworkDevicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state NetworkDevices

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
	res, err := state.readGet(&c, "/crosswork/inventory/v1/nodes/query"+params)
	if err != nil && (strings.Contains(err.Error(), "StatusCode 404") || strings.Contains(err.Error(), "StatusCode 406") || strings.Contains(err.Error(), "StatusCode 500") || strings.Contains(err.Error(), "StatusCode 400") || cw.IsNotFound(err)) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
		return
	}

	if data := res.Get("data"); !data.Exists() || len(data.Array()) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}

	// Manually managed: the provider profile is the key of the provider block, which the generated read cannot reach.
	previous, previousTags := state.ProviderProfile, state.Tags
	imported := state.isNull(ctx, res)
	if imported {
		state.fromBody(ctx, res)
		state.ProviderProfile = providerProfileName(res)
	} else {
		state.updateFromBody(ctx, res)
		state.ProviderProfile = keepProviderProfile(previous, res)
	}
	// Likewise the tags: the device lists them by name next to the ones Crosswork adds, so the user-defined ones are
	// picked out. They are read when the configuration sets them and on an import.
	if imported || previousTags != nil {
		user, err := userTagNames(&c)
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s", err))
			return
		}
		state.Tags = tagsOf(res, user, previousTags)
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

func (r *NetworkDevicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state NetworkDevices

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))

	// Manually managed: the node read back with the plan written over it, and the retry for a device that is admin up.
	if err := updateNetworkDevice(ctx, &c, plan, state); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end update

// Section below is generated&owned by "gen/generator.go". //template:begin delete
func (r *NetworkDevicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state NetworkDevices

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
	// A delete of a missing object may fail, and a missing object only needs to leave the state.
	verifyParams := ""
	verifyRes, verifyErr := state.readGet(&c, "/crosswork/inventory/v1/nodes/query"+verifyParams)
	if verifyErr != nil && !strings.Contains(verifyErr.Error(), "StatusCode 404") && !cw.IsNotFound(verifyErr) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET) before delete, got error: %s, %s", verifyErr, verifyRes.String()))
		return
	}
	if verifyErr != nil || !verifyRes.Get("data.0.inventory_id").Exists() || verifyRes.Get("data.0.inventory_id").String() != state.Id.ValueString() {
		tflog.Debug(ctx, fmt.Sprintf("%s: Resource not found or ID mismatch during delete verification, removing from state only", state.Id.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Resource verified, proceeding with delete", state.Id.ValueString()))
	res, err := c.DeleteBody(state.getPath(), strings.ReplaceAll(`{"data":[{"inventory_id":"{id}"}]}`, "{id}", state.idJSON()))
	if err != nil && !strings.Contains(err.Error(), "StatusCode 404") && !cw.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s, %s", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Section below is generated&owned by "gen/generator.go". //template:begin import
func (r *NetworkDevicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <name>. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *NetworkDevicesResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state NetworkDevices, params string) (cw.Res, error) {
	var err error
	cacheKey := "NetworkDevices::" + state.Instance.ValueString() + "::"

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
	res, err := state.readGet(c, "/crosswork/inventory/v1/nodes/query"+params)
	singleRes := res
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}

// End of section. //template:end readcache

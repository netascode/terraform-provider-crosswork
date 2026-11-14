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
var _ resource.Resource = &ProviderProfilesResource{}
var _ resource.ResourceWithImportState = &ProviderProfilesResource{}

func NewProviderProfilesResource() resource.Resource {
	return &ProviderProfilesResource{}
}

type ProviderProfilesResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *ProviderProfilesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_provider_profiles"
}

func (r *ProviderProfilesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a provider profile.").String,

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
				MarkdownDescription: helpers.NewAttributeDescription("Provider name, unique on the Crosswork cluster").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"family": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Provider family; the controller stores the SR PCE family of the Add Provider form as `ROBOT_PROVIDER_XTC`").AddStringEnumDescription("ROBOT_PROVIDER_NSO", "ROBOT_PROVIDER_EPNM", "ROBOT_PROVIDER_WAE", "ROBOT_PROVIDER_CSM", "ROBOT_PROVIDER_XTC", "ROBOT_PROVIDER_SYSLOG_STORAGE", "ROBOT_PROVIDER_SR_PCE", "ROBOT_PROVIDER_ALERT", "ROBOT_PROVIDER_DESTINATION", "ROBOT_PROVIDER_OPTIMA", "ROBOT_PROVIDER_PROXY", "ROBOT_PROVIDER_ONC", "ROBOT_PROVIDER_ACCEDIAN_PROXY", "ROBOT_PROVIDER_CROSSWORK_AI").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ROBOT_PROVIDER_NSO", "ROBOT_PROVIDER_EPNM", "ROBOT_PROVIDER_WAE", "ROBOT_PROVIDER_CSM", "ROBOT_PROVIDER_XTC", "ROBOT_PROVIDER_SYSLOG_STORAGE", "ROBOT_PROVIDER_SR_PCE", "ROBOT_PROVIDER_ALERT", "ROBOT_PROVIDER_DESTINATION", "ROBOT_PROVIDER_OPTIMA", "ROBOT_PROVIDER_PROXY", "ROBOT_PROVIDER_ONC", "ROBOT_PROVIDER_ACCEDIAN_PROXY", "ROBOT_PROVIDER_CROSSWORK_AI"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"profile": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the credential profile, which must hold a username and password for every protocol the provider uses").String,
				Required:            true,
			},
			"connectivity": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Connections to the provider, one per protocol; each uses an IP address or an FQDN").String,
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"protocol": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Protocol").AddStringEnumDescription("ROBOT_MSVC_TRANS_SSH", "ROBOT_MSVC_TRANS_NETCONF", "ROBOT_MSVC_TRANS_HTTP", "ROBOT_MSVC_TRANS_HTTPS", "ROBOT_MSVC_TRANS_TCP", "ROBOT_MSVC_TRANS_GRPC", "ROBOT_MSVC_TRANS_GNMI").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("ROBOT_MSVC_TRANS_SSH", "ROBOT_MSVC_TRANS_NETCONF", "ROBOT_MSVC_TRANS_HTTP", "ROBOT_MSVC_TRANS_HTTPS", "ROBOT_MSVC_TRANS_TCP", "ROBOT_MSVC_TRANS_GRPC", "ROBOT_MSVC_TRANS_GNMI"),
							},
						},
						"ip_address": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("IPv4 address of the provider; leave it unset when `fqdn_host_name` is used").String,
							Optional:            true,
						},
						"fqdn_host_name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Host name part of the FQDN of the provider; leave it unset when `ip_address` is used").String,
							Optional:            true,
						},
						"fqdn_domain_name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Domain name part of the FQDN of the provider").String,
							Optional:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Port").AddIntegerRangeDescription(1, 65535).String,
							Required:            true,
							Validators: []validator.Int64{
								int64validator.Between(1, 65535),
							},
						},
						"timeout": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Timeout in seconds, the controller uses 30 when it is not set").String,
							Optional:            true,
						},
						"encoding_type": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Encoding of a gNMI or gRPC connection, required for the gNMI protocol").AddStringEnumDescription("JSON", "BYTES", "PROTO", "ASCII", "JSON_IETF", "XML", "YANG").String,
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("JSON", "BYTES", "PROTO", "ASCII", "JSON_IETF", "XML", "YANG"),
							},
						},
					},
				},
			},
			"certificate_name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the certificate profile of the provider, required by some families (Accedian proxy) and checked to exist").String,
				Optional:            true,
			},
			"supported_models": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Device models the provider supports, required for the NSO family").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"model_prefix": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Model prefix").AddStringEnumDescription("Cisco-IOS-XR", "Cisco-NX-OS", "Cisco-IOS-XE").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("Cisco-IOS-XR", "Cisco-NX-OS", "Cisco-IOS-XE"),
							},
						},
						"model_version": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Model version").String,
							Required:            true,
						},
					},
				},
			},
			"properties": schema.MapAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Provider properties such as `nso_crosslaunch_url`, `embedded_nso` or `webui-root-resource`, as key and value").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
		},
	}
}

func (r *ProviderProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *ProviderProfilesResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create
func (r *ProviderProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan ProviderProfiles

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

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, ProviderProfiles{})

	params := ""
	res, err := c.Post(plan.getPath()+params, body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "POST", err, res.String()))
		return
	}
	plan.Id = types.StringValue(fmt.Sprint(plan.Name.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end create

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (r *ProviderProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state ProviderProfiles

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
	res, err := state.readGet(&c, "/crosswork/inventory/v1/providers/query"+params)
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
func (r *ProviderProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state ProviderProfiles

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

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))

	body := plan.toBody(ctx, state)
	params := ""
	res, err := c.Put(plan.getPath()+params, body)
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
func (r *ProviderProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state ProviderProfiles

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
	verifyRes, verifyErr := state.readGet(&c, "/crosswork/inventory/v1/providers/query"+verifyParams)
	if verifyErr != nil && !strings.Contains(verifyErr.Error(), "StatusCode 404") && !cw.IsNotFound(verifyErr) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET) before delete, got error: %s, %s", verifyErr, verifyRes.String()))
		return
	}
	if verifyErr != nil || !verifyRes.Get("data.0.name").Exists() || verifyRes.Get("data.0.name").String() != state.Id.ValueString() {
		tflog.Debug(ctx, fmt.Sprintf("%s: Resource not found or ID mismatch during delete verification, removing from state only", state.Id.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Resource verified, proceeding with delete", state.Id.ValueString()))
	res, err := c.DeleteBody(state.getPath(), strings.ReplaceAll(`{"providers":[{"name":"{id}"}]}`, "{id}", state.idJSON()))
	if err != nil && !strings.Contains(err.Error(), "StatusCode 404") && !cw.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s, %s", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Section below is generated&owned by "gen/generator.go". //template:begin import
func (r *ProviderProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
func (r *ProviderProfilesResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state ProviderProfiles, params string) (cw.Res, error) {
	var err error
	cacheKey := "ProviderProfiles::" + state.Instance.ValueString() + "::"

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
	res, err := state.readGet(c, "/crosswork/inventory/v1/providers/query"+params)
	singleRes := res
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}

// End of section. //template:end readcache

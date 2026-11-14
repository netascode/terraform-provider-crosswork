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
var _ resource.Resource = &NSOProviderPoliciesResource{}
var _ resource.ResourceWithImportState = &NSOProviderPoliciesResource{}

func NewNSOProviderPoliciesResource() resource.Resource {
	return &NSOProviderPoliciesResource{}
}

type NSOProviderPoliciesResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *NSOProviderPoliciesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nso_provider_policies"
}

func (r *NSOProviderPoliciesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage the device onboarding policy of an NSO provider.").String,

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
			"provider_name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the NSO provider the policy belongs to").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"match": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Match devices of Crosswork and of the NSO provider by their IP address").String,
				Optional:            true,
			},
			"match_rule": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Expression that selects the devices that are matched, for example `*` or `host_name = 'build-*'`").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 4096),
				},
			},
			"onboard_to_nso": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Add the devices that Crosswork has and the NSO provider is missing to the NSO provider").String,
				Optional:            true,
			},
			"onboard_to_nso_rule": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Expression that selects the devices that are added to the NSO provider").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 4096),
				},
			},
			"onboard_from_nso": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Add the devices that the NSO provider has and Crosswork is missing to Crosswork").String,
				Optional:            true,
			},
			"onboard_from_nso_rule": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Expression that selects the devices that are added to Crosswork").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 4096),
				},
			},
			"sync_from_nso": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Run a sync-from on the NSO provider after the devices are onboarded").String,
				Optional:            true,
			},
			"sync_from_nso_rule": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Expression that selects the devices that are synchronised").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 4096),
				},
			},
			"check_sync": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Check the sync status of the devices on the NSO provider").String,
				Optional:            true,
			},
			"check_sync_rule": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Expression that selects the devices whose sync status is checked").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 4096),
				},
			},
			"neds": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Rules that select the NED an onboarded device gets. Both the NED and the rule are required.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ned": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of the NED, for example `cisco-iosxr-cli-7.76`").String,
							Required:            true,
						},
						"rule": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Expression that selects the devices that get the NED, for example `product_info.software_type = 'IOS XR*'`").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 4096),
							},
						},
					},
				},
			},
		},
	}
}

func (r *NSOProviderPoliciesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *NSOProviderPoliciesResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// End of section. //template:end model

// Written by hand: the entry of the provider exists (the controller creates it with the provider), so a create writes the values of
// the configuration into it, the same as an update.
func (r *NSOProviderPoliciesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan NSOProviderPolicies

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.ProviderName.ValueString()))

	if err := applyNSOProviderPolicy(ctx, &c, &plan); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s", "policy", err))
		return
	}
	plan.Id = types.StringValue(plan.ProviderName.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end create

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (r *NSOProviderPoliciesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state NSOProviderPolicies

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
	res, err := state.readGet(&c, "/crosswork/inventory/v1/nso/policy/query"+params)
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

// Written by hand: the entry of the provider is replaced in the policy document, everything else in it is sent back as it was.
func (r *NSOProviderPoliciesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan NSOProviderPolicies

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.ProviderName.ValueString()))

	if err := applyNSOProviderPolicy(ctx, &c, &plan); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s", "policy", err))
		return
	}
	plan.Id = types.StringValue(plan.ProviderName.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end update

// Written by hand: the controller creates the policy with the NSO provider and removes it with the provider, so a destroy only
// removes the resource from the state and leaves the policy as it is.
func (r *NSOProviderPoliciesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state NSOProviderPolicies

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete leaves the policy of the controller as it is", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Section below is generated&owned by "gen/generator.go". //template:begin import
func (r *NSOProviderPoliciesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <provider_name>. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("provider_name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *NSOProviderPoliciesResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state NSOProviderPolicies, params string) (cw.Res, error) {
	var err error
	cacheKey := "NSOProviderPolicies::" + state.Instance.ValueString() + "::"

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
	res, err := state.readGet(c, "/crosswork/inventory/v1/nso/policy/query"+params)
	singleRes := res
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}

// End of section. //template:end readcache

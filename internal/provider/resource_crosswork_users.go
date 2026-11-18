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
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
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
var _ resource.Resource = &UsersResource{}
var _ resource.ResourceWithImportState = &UsersResource{}

func NewUsersResource() resource.Resource {
	return &UsersResource{}
}

type UsersResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *UsersResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

func (r *UsersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a user in the local authentication store.").String,

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
			"username": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Username, 2 to 50 characters (letters, numbers and . - _ @)").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Password, 8 to 72 characters with upper and lower case letters, a number and one of @!$%*?&").AddMutualExclusivityDescription("**Required**: exactly one of `password` and `password_wo` must be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `password_wo` together with `password_wo_version`, which keeps it out of state.").String,
				Sensitive:           true,
				Optional:            true,
			},
			"password_wo": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Password, 8 to 72 characters with upper and lower case letters, a number and one of @!$%*?&").AddMutualExclusivityDescription("**Required**: exactly one of `password` and `password_wo` must be set.").String,
				Optional:            true,
				WriteOnly:           true,
				Sensitive:           true,
			},
			"password_wo_version": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `password_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
				Optional:            true,
			},
			"policy_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Role of the user").String,
				Required:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("First name").String,
				Required:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Last name").String,
				Required:            true,
			},
			"device_access_groups": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("List of device access groups of the user").String,
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Device access group UUID").String,
							Required:            true,
						},
						"domain_name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Device access group name").String,
							Required:            true,
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
		},
	}
}

func (r *UsersResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *UsersResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// ValidateConfig enforces the relationship between a secret attribute, its write-only
// "_wo" counterpart and the "_wo_version" rotation trigger.
//
// The checks are not schema validators because Terraform renders an attribute-scoped
// diagnostic together with the offending configuration line, which for a secret prints
// the value into plan output and logs. A resource-scoped diagnostic names the attributes
// (and the list element index) instead.
func (r *UsersResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var legacyPassword types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password"), &legacyPassword)...)
	var woPassword types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo"), &woPassword)...)
	var woVersionPassword types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo_version"), &woVersionPassword)...)
	if !legacyPassword.IsUnknown() && !woPassword.IsUnknown() && !legacyPassword.IsNull() && !woPassword.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Only one of `password` and `password_wo` can be set.",
		)
	}
	if !legacyPassword.IsUnknown() && !woPassword.IsUnknown() && legacyPassword.IsNull() && woPassword.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Exactly one of `password` and `password_wo` must be set.",
		)
	}
	if !woPassword.IsUnknown() && !woVersionPassword.IsUnknown() && !woPassword.IsNull() && woVersionPassword.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"`password_wo_version` must be set when `password_wo` is used. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.",
		)
	}
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create
func (r *UsersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan Users

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
	// Write-only values are not in plan or state; read them from config.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo"), &plan.PasswordWo)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, Users{})
	body = plan.postBody(body)

	params := ""
	res, err := c.Post(plan.getPath()+params, body, cw.UseMutex)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "POST", err, res.String()))
		return
	}
	plan.Id = types.StringValue(fmt.Sprint(plan.Username.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end create

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (r *UsersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state Users

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
	params += "/" + url.QueryEscape(state.Id.ValueString())
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
func (r *UsersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state Users

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
	// Read on every Update: a full-object replace PUT is rejected when it omits the secret; the version companion only drives change detection.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("password_wo"), &plan.PasswordWo)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))

	body := plan.toBody(ctx, state)
	params := ""
	res, err := c.Put(plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString())+params, body, cw.UseMutex)
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
func (r *UsersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state Users

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
	res, err := c.Delete(state.getPath()+"/"+url.QueryEscape(state.Id.ValueString()), cw.UseMutex)
	if err != nil && !strings.Contains(err.Error(), "StatusCode 404") && !cw.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s, %s", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Section below is generated&owned by "gen/generator.go". //template:begin import
func (r *UsersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <username>. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("username"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *UsersResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state Users, params string) (cw.Res, error) {
	var err error
	cacheKey := "Users::" + state.Instance.ValueString() + "::"

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

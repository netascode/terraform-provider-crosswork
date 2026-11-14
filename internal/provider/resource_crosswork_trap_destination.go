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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
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
var _ resource.Resource = &TrapDestinationResource{}
var _ resource.ResourceWithImportState = &TrapDestinationResource{}

func NewTrapDestinationResource() resource.Resource {
	return &TrapDestinationResource{}
}

type TrapDestinationResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *TrapDestinationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_trap_destination"
}

func (r *TrapDestinationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a trap receiver for the notification destinations.").String,

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
				MarkdownDescription: helpers.NewAttributeDescription("Name of the receiver. A receiver whose address is a host name takes the host name as its name").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 128),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("IP address or host name of the receiver").String,
				Required:            true,
			},
			"address_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the address is an IP address or a host name").AddStringEnumDescription("IPADDRESS", "DNS").AddDefaultValueDescription("IPADDRESS").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("IPADDRESS", "DNS"),
				},
				Default: stringdefault.StaticString("IPADDRESS"),
			},
			"port": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Port of the receiver").AddIntegerRangeDescription(100, 65535).AddDefaultValueDescription("162").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(100, 65535),
				},
				Default: int64default.StaticInt64(162),
			},
			"transport_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Transport protocol").AddStringEnumDescription("UDP", "TCP").AddDefaultValueDescription("UDP").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("UDP", "TCP"),
				},
				Default: stringdefault.StaticString("UDP"),
			},
			"snmp_version": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP version").AddStringEnumDescription("VERSION_2C", "VERSION_3").AddDefaultValueDescription("VERSION_2C").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("VERSION_2C", "VERSION_3"),
				},
				Default: stringdefault.StaticString("VERSION_2C"),
			},
			"community": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP v2c community").AddMutualExclusivityDescription("Only one of `community` and `community_wo` can be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `community_wo` together with `community_wo_version`, which keeps it out of state.").String,
				Sensitive:           true,
				Optional:            true,
			},
			"community_wo": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP v2c community").AddMutualExclusivityDescription("Only one of `community` and `community_wo` can be set.").String,
				Optional:            true,
				WriteOnly:           true,
				Sensitive:           true,
			},
			"community_wo_version": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `community_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
				Optional:            true,
			},
			"username": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP v3 user name").String,
				Optional:            true,
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP v3 security mode").AddStringEnumDescription("AuthPriv", "AuthNoPriv", "NoAuthNoPriv").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("AuthPriv", "AuthNoPriv", "NoAuthNoPriv"),
				},
			},
			"auth_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP v3 authentication type").AddStringEnumDescription("HMACMD5", "HMACSHA").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("HMACMD5", "HMACSHA"),
				},
			},
			"auth_password": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP v3 authentication password").AddMutualExclusivityDescription("Only one of `auth_password` and `auth_password_wo` can be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `auth_password_wo` together with `auth_password_wo_version`, which keeps it out of state.").String,
				Sensitive:           true,
				Optional:            true,
			},
			"auth_password_wo": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP v3 authentication password").AddMutualExclusivityDescription("Only one of `auth_password` and `auth_password_wo` can be set.").String,
				Optional:            true,
				WriteOnly:           true,
				Sensitive:           true,
			},
			"auth_password_wo_version": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `auth_password_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
				Optional:            true,
			},
			"privacy_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP v3 privacy type. `NONE` is for the mode `AuthNoPriv` and is refused with `AuthPriv`").AddStringEnumDescription("NONE", "DES", "AES128", "AES192", "AES256", "CiscoAES192", "CiscoAES256").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("NONE", "DES", "AES128", "AES192", "AES256", "CiscoAES192", "CiscoAES256"),
				},
			},
			"privacy_password": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP v3 privacy password").AddMutualExclusivityDescription("Only one of `privacy_password` and `privacy_password_wo` can be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `privacy_password_wo` together with `privacy_password_wo_version`, which keeps it out of state.").String,
				Sensitive:           true,
				Optional:            true,
			},
			"privacy_password_wo": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP v3 privacy password").AddMutualExclusivityDescription("Only one of `privacy_password` and `privacy_password_wo` can be set.").String,
				Optional:            true,
				WriteOnly:           true,
				Sensitive:           true,
			},
			"privacy_password_wo_version": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `privacy_password_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
				Optional:            true,
			},
		},
	}
}

func (r *TrapDestinationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *TrapDestinationResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// ValidateConfig enforces the relationship between a secret attribute, its write-only
// "_wo" counterpart and the "_wo_version" rotation trigger.
//
// The checks are not schema validators because Terraform renders an attribute-scoped
// diagnostic together with the offending configuration line, which for a secret prints
// the value into plan output and logs. A resource-scoped diagnostic names the attributes
// (and the list element index) instead.
func (r *TrapDestinationResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var legacyCommunity types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("community"), &legacyCommunity)...)
	var woCommunity types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("community_wo"), &woCommunity)...)
	var woVersionCommunity types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("community_wo_version"), &woVersionCommunity)...)
	if !legacyCommunity.IsUnknown() && !woCommunity.IsUnknown() && !legacyCommunity.IsNull() && !woCommunity.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Only one of `community` and `community_wo` can be set.",
		)
	}
	if !woCommunity.IsUnknown() && !woVersionCommunity.IsUnknown() && !woCommunity.IsNull() && woVersionCommunity.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"`community_wo_version` must be set when `community_wo` is used. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.",
		)
	}
	var legacyAuthPassword types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("auth_password"), &legacyAuthPassword)...)
	var woAuthPassword types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("auth_password_wo"), &woAuthPassword)...)
	var woVersionAuthPassword types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("auth_password_wo_version"), &woVersionAuthPassword)...)
	if !legacyAuthPassword.IsUnknown() && !woAuthPassword.IsUnknown() && !legacyAuthPassword.IsNull() && !woAuthPassword.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Only one of `auth_password` and `auth_password_wo` can be set.",
		)
	}
	if !woAuthPassword.IsUnknown() && !woVersionAuthPassword.IsUnknown() && !woAuthPassword.IsNull() && woVersionAuthPassword.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"`auth_password_wo_version` must be set when `auth_password_wo` is used. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.",
		)
	}
	var legacyPrivacyPassword types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("privacy_password"), &legacyPrivacyPassword)...)
	var woPrivacyPassword types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("privacy_password_wo"), &woPrivacyPassword)...)
	var woVersionPrivacyPassword types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("privacy_password_wo_version"), &woVersionPrivacyPassword)...)
	if !legacyPrivacyPassword.IsUnknown() && !woPrivacyPassword.IsUnknown() && !legacyPrivacyPassword.IsNull() && !woPrivacyPassword.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Only one of `privacy_password` and `privacy_password_wo` can be set.",
		)
	}
	if !woPrivacyPassword.IsUnknown() && !woVersionPrivacyPassword.IsUnknown() && !woPrivacyPassword.IsNull() && woVersionPrivacyPassword.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"`privacy_password_wo_version` must be set when `privacy_password_wo` is used. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.",
		)
	}
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create
func (r *TrapDestinationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan TrapDestination

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
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("community_wo"), &plan.CommunityWo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Write-only values are not in plan or state; read them from config.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("auth_password_wo"), &plan.AuthPasswordWo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Write-only values are not in plan or state; read them from config.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("privacy_password_wo"), &plan.PrivacyPasswordWo)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, TrapDestination{})

	params := ""
	res, err := c.Post(plan.getPath()+params, body, cw.UseMutex)
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

// Written by hand: the list of the controller holds every receiver, which readGet searches.
func (r *TrapDestinationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state TrapDestination

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
	res, err := state.readGet(&c, "/crosswork/alarm/v1/destination/trap"+params)
	if err != nil && (strings.Contains(err.Error(), "StatusCode 404") || strings.Contains(err.Error(), "StatusCode 406") || strings.Contains(err.Error(), "StatusCode 500") || strings.Contains(err.Error(), "StatusCode 400") || cw.IsNotFound(err)) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
		return
	}

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

// Written by hand: the update (PUT) of the controller does not work, so the receiver is removed by its name and added again with the new values.
func (r *TrapDestinationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state TrapDestination

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
	// community_wo is write-only, so it is read from the configuration. Updates are full-object PUTs and the API requires the secret on every write, so it is read on every Update.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("community_wo"), &plan.CommunityWo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// auth_password_wo is write-only, so it is read from the configuration. Updates are full-object PUTs and the API requires the secret on every write, so it is read on every Update.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("auth_password_wo"), &plan.AuthPasswordWo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// privacy_password_wo is write-only, so it is read from the configuration. Updates are full-object PUTs and the API requires the secret on every write, so it is read on every Update.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("privacy_password_wo"), &plan.PrivacyPasswordWo)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))

	res, err := c.Delete(state.getPath()+"/"+url.QueryEscape(state.Id.ValueString()), cw.UseMutex)
	if err != nil && !strings.Contains(err.Error(), "StatusCode 404") && !cw.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to remove the receiver before it is added with the new values (DELETE), got error: %s, %s", err, res.String()))
		return
	}
	res, err = c.Post(plan.getPath(), plan.toBody(ctx, TrapDestination{}), cw.UseMutex)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("The receiver was removed and could not be added with the new values (POST), got error: %s, %s. The next apply adds it", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end update

// Section below is generated&owned by "gen/generator.go". //template:begin delete
func (r *TrapDestinationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state TrapDestination

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
func (r *TrapDestinationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
func (r *TrapDestinationResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state TrapDestination, params string) (cw.Res, error) {
	var err error
	cacheKey := "TrapDestination::" + state.Instance.ValueString() + "::"

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
	res, err := c.Get("/crosswork/alarm/v1/destination/trap" + params)
	singleRes := res
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}

// End of section. //template:end readcache

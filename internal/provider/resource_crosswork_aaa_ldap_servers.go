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
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &AAALDAPServersResource{}
var _ resource.ResourceWithImportState = &AAALDAPServersResource{}

func NewAAALDAPServersResource() resource.Resource {
	return &AAALDAPServersResource{}
}

type AAALDAPServersResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *AAALDAPServersResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aaa_ldap_servers"
}

func (r *AAALDAPServersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage the LDAP servers used for authentication.").String,

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
			"servers": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("LDAP servers").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of the connection").String,
							Required:            true,
						},
						"ldap_url": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("URL of the server, `ldap://` or `ldaps://`").String,
							Required:            true,
						},
						"bind_dn": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Distinguished name used to bind to the server").String,
							Required:            true,
						},
						"bind_credential": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Password used to bind to the server").AddMutualExclusivityDescription("**Required**: exactly one of `bind_credential` and `bind_credential_wo` must be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `bind_credential_wo` together with `bind_credential_wo_version`, which keeps it out of state.").String,
							Sensitive:           true,
							Optional:            true,
						},
						"bind_credential_wo": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Password used to bind to the server").AddMutualExclusivityDescription("**Required**: exactly one of `bind_credential` and `bind_credential_wo` must be set.").String,
							Optional:            true,
							WriteOnly:           true,
							Sensitive:           true,
						},
						"bind_credential_wo_version": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `bind_credential_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
							Optional:            true,
						},
						"base_dn": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Distinguished name where the search for users starts").String,
							Required:            true,
						},
						"user_filter": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Filter that matches the users, with `{user}` for the user name").String,
							Required:            true,
						},
						"dn_format": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Format of the distinguished name of a user, with `%s` for the user name").String,
							Required:            true,
						},
						"principal_attribute_id": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Attribute of the user name in the directory").String,
							Required:            true,
						},
						"policy_id": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Attribute of the role mapping in the directory").String,
							Required:            true,
						},
						"order": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Order of the server among all remote authentication servers").String,
							Required:            true,
						},
						"connect_timeout": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Connection timeout in seconds").String,
							Required:            true,
						},
					},
				},
			},
			"attribute_servers": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("LDAP servers that look up the role of a user, one for each server").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"ldap_url": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("URL of the server, `ldap://` followed by host and port").String,
							Required:            true,
						},
						"base_dn": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Distinguished name where the search for the role starts").String,
							Required:            true,
						},
						"bind_dn": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Distinguished name used to bind to the server").String,
							Required:            true,
						},
						"bind_credential": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Password used to bind to the server").AddMutualExclusivityDescription("**Required**: exactly one of `bind_credential` and `bind_credential_wo` must be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `bind_credential_wo` together with `bind_credential_wo_version`, which keeps it out of state.").String,
							Sensitive:           true,
							Optional:            true,
						},
						"bind_credential_wo": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Password used to bind to the server").AddMutualExclusivityDescription("**Required**: exactly one of `bind_credential` and `bind_credential_wo` must be set.").String,
							Optional:            true,
							WriteOnly:           true,
							Sensitive:           true,
						},
						"bind_credential_wo_version": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `bind_credential_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
							Optional:            true,
						},
						"user_filter": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Attribute of the directory that maps the user to its role").String,
							Required:            true,
						},
						"connect_timeout": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Connection timeout in seconds").String,
							Required:            true,
						},
						"fail_fast": schema.BoolAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Fail at once when the server cannot be reached").AddDefaultValueDescription("false").String,
							Optional:            true,
							Computed:            true,
							Default:             booldefault.StaticBool(false),
						},
						"attributes": schema.MapAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Map of the directory key that carries the role to `policy_id`").String,
							ElementType:         types.StringType,
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (r *AAALDAPServersResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *AAALDAPServersResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// ValidateConfig enforces the relationship between a secret attribute, its write-only
// "_wo" counterpart and the "_wo_version" rotation trigger.
//
// The checks are not schema validators because Terraform renders an attribute-scoped
// diagnostic together with the offending configuration line, which for a secret prints
// the value into plan output and logs. A resource-scoped diagnostic names the attributes
// (and the list element index) instead.
func (r *AAALDAPServersResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	{
		// Secrets nested in "servers" are validated per element. A list that cannot be
		// read as a whole - because it is still unknown at validation time - is skipped
		// rather than reported, since there is nothing to check yet.
		var cfgServers []AAALDAPServersServers
		if diags := req.Config.GetAttribute(ctx, path.Root("servers"), &cfgServers); !diags.HasError() {
			for i := range cfgServers {
				if !cfgServers[i].BindCredential.IsUnknown() && !cfgServers[i].BindCredentialWo.IsUnknown() && !cfgServers[i].BindCredential.IsNull() && !cfgServers[i].BindCredentialWo.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("Only one of `bind_credential` and `bind_credential_wo` can be set in `servers` element %d.", i),
					)
				}
				if !cfgServers[i].BindCredential.IsUnknown() && !cfgServers[i].BindCredentialWo.IsUnknown() && cfgServers[i].BindCredential.IsNull() && cfgServers[i].BindCredentialWo.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("Exactly one of `bind_credential` and `bind_credential_wo` must be set in `servers` element %d.", i),
					)
				}
				if !cfgServers[i].BindCredentialWo.IsUnknown() && !cfgServers[i].BindCredentialWoVersion.IsUnknown() && !cfgServers[i].BindCredentialWo.IsNull() && cfgServers[i].BindCredentialWoVersion.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("`bind_credential_wo_version` must be set when `bind_credential_wo` is used in `servers` element %d. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.", i),
					)
				}
			}
		}
	}
	{
		// Secrets nested in "attribute_servers" are validated per element. A list that cannot be
		// read as a whole - because it is still unknown at validation time - is skipped
		// rather than reported, since there is nothing to check yet.
		var cfgAttributeServers []AAALDAPServersAttributeServers
		if diags := req.Config.GetAttribute(ctx, path.Root("attribute_servers"), &cfgAttributeServers); !diags.HasError() {
			for i := range cfgAttributeServers {
				if !cfgAttributeServers[i].BindCredential.IsUnknown() && !cfgAttributeServers[i].BindCredentialWo.IsUnknown() && !cfgAttributeServers[i].BindCredential.IsNull() && !cfgAttributeServers[i].BindCredentialWo.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("Only one of `bind_credential` and `bind_credential_wo` can be set in `attribute_servers` element %d.", i),
					)
				}
				if !cfgAttributeServers[i].BindCredential.IsUnknown() && !cfgAttributeServers[i].BindCredentialWo.IsUnknown() && cfgAttributeServers[i].BindCredential.IsNull() && cfgAttributeServers[i].BindCredentialWo.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("Exactly one of `bind_credential` and `bind_credential_wo` must be set in `attribute_servers` element %d.", i),
					)
				}
				if !cfgAttributeServers[i].BindCredentialWo.IsUnknown() && !cfgAttributeServers[i].BindCredentialWoVersion.IsUnknown() && !cfgAttributeServers[i].BindCredentialWo.IsNull() && cfgAttributeServers[i].BindCredentialWoVersion.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("`bind_credential_wo_version` must be set when `bind_credential_wo` is used in `attribute_servers` element %d. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.", i),
					)
				}
			}
		}
	}
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create
func (r *AAALDAPServersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan AAALDAPServers

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
	// Write-only values in the list are not in plan or state; copy them from the config list into the plan elements.
	{
		var cfgServers []AAALDAPServersServers
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("servers"), &cfgServers)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for i := range plan.Servers {
			if i < len(cfgServers) {
				plan.Servers[i].BindCredentialWo = cfgServers[i].BindCredentialWo
			}
		}
	}
	// Write-only values in the list are not in plan or state; copy them from the config list into the plan elements.
	{
		var cfgAttributeServers []AAALDAPServersAttributeServers
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("attribute_servers"), &cfgAttributeServers)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for i := range plan.AttributeServers {
			if i < len(cfgAttributeServers) {
				plan.AttributeServers[i].BindCredentialWo = cfgAttributeServers[i].BindCredentialWo
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, AAALDAPServers{})

	params := ""
	res, err := c.Put(plan.getPath()+params, body, cw.UseMutex)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "PUT", err, res.String()))
		return
	}
	plan.Id = types.StringValue(res.Get("!\"aaa_ldap_servers\"").String())

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end create

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (r *AAALDAPServersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state AAALDAPServers

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
	res, err := c.Get("/crosswork/authconfig/v1/remote/ldap" + params)
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
func (r *AAALDAPServersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state AAALDAPServers

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
	// Read unconditionally, for the same reason as the top-level secrets above.
	{
		var cfgServers []AAALDAPServersServers
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("servers"), &cfgServers)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for i := range plan.Servers {
			if i < len(cfgServers) {
				plan.Servers[i].BindCredentialWo = cfgServers[i].BindCredentialWo
			}
		}
	}
	// Read unconditionally, for the same reason as the top-level secrets above.
	{
		var cfgAttributeServers []AAALDAPServersAttributeServers
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("attribute_servers"), &cfgAttributeServers)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for i := range plan.AttributeServers {
			if i < len(cfgAttributeServers) {
				plan.AttributeServers[i].BindCredentialWo = cfgAttributeServers[i].BindCredentialWo
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))

	body := plan.toBody(ctx, state)
	params := ""
	res, err := c.Put(plan.getPath()+params, body, cw.UseMutex)
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
func (r *AAALDAPServersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state AAALDAPServers

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
	res, err := c.Put(state.getPathDelete(), `{"ldap_server":[],"ldap_attr_servers":{"ldap_attr_server":[]}}`, cw.UseMutex)
	if err != nil && !strings.Contains(err.Error(), "StatusCode 404") && !cw.IsNotFound(err) {
		errorCode := res.Get("response.errorCode").String()
		if deleteWarnErrorPrefix != "" && strings.HasPrefix(errorCode, deleteWarnErrorPrefix) {
			failureReason := res.Get("response.failureReason").String()
			resp.Diagnostics.AddWarning("Empty input Warning", fmt.Sprintf("Empty input detected (error code: %s, reason %s).", errorCode, failureReason))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (%s), got error: %s, %s", "PUT", err, res.String()))
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Section below is generated&owned by "gen/generator.go". //template:begin import
func (r *AAALDAPServersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "aaa_ldap_servers" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: aaa_ldap_servers. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *AAALDAPServersResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state AAALDAPServers, params string) (cw.Res, error) {
	var err error
	cacheKey := "AAALDAPServers::" + state.Instance.ValueString() + "::"

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
	res, err := c.Get("/crosswork/authconfig/v1/remote/ldap" + params)
	singleRes := res
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}

// End of section. //template:end readcache

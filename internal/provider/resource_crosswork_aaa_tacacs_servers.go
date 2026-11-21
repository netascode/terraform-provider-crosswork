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
var _ resource.Resource = &AAATACACSServersResource{}
var _ resource.ResourceWithImportState = &AAATACACSServersResource{}

func NewAAATACACSServersResource() resource.Resource {
	return &AAATACACSServersResource{}
}

type AAATACACSServersResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *AAATACACSServersResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aaa_tacacs_servers"
}

func (r *AAATACACSServersResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage the TACACS+ servers used for authentication.").String,

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
				MarkdownDescription: helpers.NewAttributeDescription("TACACS+ servers").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("IP address or host name of the server").String,
							Required:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Port of the server").AddDefaultValueDescription("49").String,
							Optional:            true,
							Computed:            true,
							Default:             int64default.StaticInt64(49),
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Shared secret").AddMutualExclusivityDescription("**Required**: exactly one of `secret` and `secret_wo` must be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `secret_wo` together with `secret_wo_version`, which keeps it out of state.").String,
							Sensitive:           true,
							Optional:            true,
						},
						"secret_wo": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Shared secret").AddMutualExclusivityDescription("**Required**: exactly one of `secret` and `secret_wo` must be set.").String,
							Optional:            true,
							WriteOnly:           true,
							Sensitive:           true,
						},
						"secret_wo_version": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `secret_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
							Optional:            true,
						},
						"secret_format": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Format of the shared secret").AddStringEnumDescription("ascii", "hex").AddDefaultValueDescription("ascii").String,
							Optional:            true,
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("ascii", "hex"),
							},
							Default: stringdefault.StaticString("ascii"),
						},
						"service": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("TACACS+ service name to request").String,
							Required:            true,
						},
						"priority": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Order of the server among all remote authentication servers").String,
							Required:            true,
						},
						"retransmit_timeout": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Retransmit timeout in seconds").String,
							Optional:            true,
						},
						"retries": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Number of retries before the request fails").AddDefaultValueDescription("0").String,
							Optional:            true,
							Computed:            true,
							Default:             int64default.StaticInt64(0),
						},
						"auth_type": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Authentication type").AddStringEnumDescription("pap", "chap").AddDefaultValueDescription("pap").String,
							Optional:            true,
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("pap", "chap"),
							},
							Default: stringdefault.StaticString("pap"),
						},
						"interface": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Interface IP address Crosswork uses to reach the server").String,
							Optional:            true,
						},
						"policy_id": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Key name of the role on the server").String,
							Optional:            true,
						},
						"virtual_domain": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Device access group attribute").AddDefaultValueDescription("deviceAccessGroups").String,
							Optional:            true,
							Computed:            true,
							Default:             stringdefault.StaticString("deviceAccessGroups"),
						},
					},
				},
			},
		},
	}
}

func (r *AAATACACSServersResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *AAATACACSServersResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// ValidateConfig enforces the relationship between a secret attribute, its write-only
// "_wo" counterpart and the "_wo_version" rotation trigger.
//
// The checks are not schema validators because Terraform renders an attribute-scoped
// diagnostic together with the offending configuration line, which for a secret prints
// the value into plan output and logs. A resource-scoped diagnostic names the attributes
// (and the list element index) instead.
func (r *AAATACACSServersResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	{
		// Secrets nested in "servers" are validated per element. A list that cannot be
		// read as a whole - because it is still unknown at validation time - is skipped
		// rather than reported, since there is nothing to check yet.
		var cfgServers []AAATACACSServersServers
		if diags := req.Config.GetAttribute(ctx, path.Root("servers"), &cfgServers); !diags.HasError() {
			for i := range cfgServers {
				if !cfgServers[i].Secret.IsUnknown() && !cfgServers[i].SecretWo.IsUnknown() && !cfgServers[i].Secret.IsNull() && !cfgServers[i].SecretWo.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("Only one of `secret` and `secret_wo` can be set in `servers` element %d.", i),
					)
				}
				if !cfgServers[i].Secret.IsUnknown() && !cfgServers[i].SecretWo.IsUnknown() && cfgServers[i].Secret.IsNull() && cfgServers[i].SecretWo.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("Exactly one of `secret` and `secret_wo` must be set in `servers` element %d.", i),
					)
				}
				if !cfgServers[i].SecretWo.IsUnknown() && !cfgServers[i].SecretWoVersion.IsUnknown() && !cfgServers[i].SecretWo.IsNull() && cfgServers[i].SecretWoVersion.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("`secret_wo_version` must be set when `secret_wo` is used in `servers` element %d. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.", i),
					)
				}
			}
		}
	}
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create
func (r *AAATACACSServersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan AAATACACSServers

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
		var cfgServers []AAATACACSServersServers
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("servers"), &cfgServers)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for i := range plan.Servers {
			if i < len(cfgServers) {
				plan.Servers[i].SecretWo = cfgServers[i].SecretWo
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, AAATACACSServers{})

	params := ""
	res, err := c.Put(plan.getPath()+params, body, cw.UseMutex)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "PUT", err, res.String()))
		return
	}
	plan.Id = types.StringValue(res.Get("!\"aaa_tacacs_servers\"").String())

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end create

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (r *AAATACACSServersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state AAATACACSServers

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
	res, err := c.Get("/crosswork/authconfig/v1/remote/tacacs" + params)
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
func (r *AAATACACSServersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state AAATACACSServers

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
		var cfgServers []AAATACACSServersServers
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("servers"), &cfgServers)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for i := range plan.Servers {
			if i < len(cfgServers) {
				plan.Servers[i].SecretWo = cfgServers[i].SecretWo
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
func (r *AAATACACSServersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state AAATACACSServers

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
	res, err := c.Put(state.getPathDelete(), `{"tacacs_servers":[]}`, cw.UseMutex)
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
func (r *AAATACACSServersResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != "aaa_tacacs_servers" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: aaa_tacacs_servers. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *AAATACACSServersResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state AAATACACSServers, params string) (cw.Res, error) {
	var err error
	cacheKey := "AAATACACSServers::" + state.Instance.ValueString() + "::"

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
	res, err := c.Get("/crosswork/authconfig/v1/remote/tacacs" + params)
	singleRes := res
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}

// End of section. //template:end readcache

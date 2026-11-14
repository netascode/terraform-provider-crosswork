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
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

// End of section. //template:end imports

// Written by hand: data_gateway_ids is computed only, which the generator cannot express (it makes a computed attribute optional too).

var _ resource.Resource = &DataGatewayPoolsResource{}
var _ resource.ResourceWithImportState = &DataGatewayPoolsResource{}

func NewDataGatewayPoolsResource() resource.Resource {
	return &DataGatewayPoolsResource{}
}

type DataGatewayPoolsResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *DataGatewayPoolsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_gateway_pools"
}

func (r *DataGatewayPoolsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a pool of Crosswork Data Gateway instances.").String,

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
				MarkdownDescription: helpers.NewAttributeDescription("Pool name, made of letters, digits, underscores, hyphens and dots").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z0-9_.-]+$`), ""),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Description of the pool").String,
				Optional:            true,
			},
			"virtual_ips": schema.ListAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Virtual IP addresses of the pool, all IPv4 or all IPv6, one for each active instance; each one is the address of one data gateway in the order of this list").String,
				ElementType:         types.StringType,
				Required:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
			},
			"prefix_length": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Prefix length of the virtual IP addresses, 1 to 32 for IPv4 and 1 to 128 for IPv6").AddIntegerRangeDescription(1, 128).String,
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 128),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"spares": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Number of spare instances of the pool").AddDefaultValueDescription("0").String,
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(0),
			},
			"instances": schema.ListAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Ids of the data gateway instances of the pool, as `crosswork_data_gateway_instances` lists them; the first one becomes the active instance of the first data gateway").String,
				ElementType:         types.StringType,
				Required:            true,
			},
			"cdg_readiness_timeout": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Seconds to wait for the controller to accept a change of the pool and for its data gateways to appear; it only limits how long Terraform waits and is not part of the Crosswork configuration").AddIntegerRangeDescription(30, 1800).AddDefaultValueDescription("60").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(30, 1800),
				},
				Default: int64default.StaticInt64(60),
			},
			"data_gateway_ids": schema.ListAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Ids of the data gateways of the pool, in the order of the virtual IPs").String,
				ElementType:         types.StringType,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DataGatewayPoolsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

func (r *DataGatewayPoolsResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// End of section. //template:end model

// Written by hand: the id comes back as "<uuid> <name>", a change is repeated while the controller settles, and the
// data gateways of the pool are read once they are listed.
func (r *DataGatewayPoolsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan DataGatewayPools

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Name.ValueString()))

	id, gateways, err := createPool(ctx, &c, plan)
	if id == "" {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (POST), got error: %s", err))
		return
	}
	plan.Id = types.StringValue(id)
	plan.DataGatewayIds = stringList(gateways)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("The pool %s was created, but not completely: %s", plan.Name.ValueString(), err))
	}
}

// End of section. //template:end create

// Written by hand: the query takes no filter, so the pool is found by its id in the list, and its data gateways come from a second call.
func (r *DataGatewayPoolsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state DataGatewayPools

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

	res, err := listPools(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
		return
	}
	pool := findPool(res, state.Id.ValueString(), "")
	if !pool.Exists() {
		resp.State.RemoveResource(ctx)
		return
	}

	if state.isNull(ctx, pool) {
		state.fromBody(ctx, pool)
	} else {
		state.updateFromBody(ctx, pool)
	}
	gateways, err := listPoolGateways(&c, state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the data gateways of the pool, got error: %s", err))
		return
	}
	state.DataGatewayIds = stringList(gateways)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

// Written by hand: the whole pool is sent with its id, and a change is repeated while the controller settles.
func (r *DataGatewayPoolsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state DataGatewayPools

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

	if err := updatePool(ctx, &c, plan, state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s", err))
		return
	}
	plan.Id = state.Id
	// A change that does not replace the pool leaves its data gateways as they are.
	plan.DataGatewayIds = state.DataGatewayIds

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end update

// Written by hand: the pool is deleted by its name and id, the controller refuses while a data gateway of the pool has devices
// attached (its message is returned unchanged), and the answer is empty, so the resource waits until the pool is gone.
func (r *DataGatewayPoolsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state DataGatewayPools

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

	res, err := listPools(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET) before delete, got error: %s, %s", err, res.String()))
		return
	}
	if !findPool(res, state.Id.ValueString(), "").Exists() {
		tflog.Debug(ctx, fmt.Sprintf("%s: Resource not found during delete verification, removing from state only", state.Id.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err := deletePool(ctx, &c, state); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Written by hand: the import identifier is the name of the pool or its id, and the id of the state is the id the controller gave it.
func (r *DataGatewayPoolsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <name> or <id>. Got: %q", req.ID),
		)
		return
	}
	c, ok := r.clientFor(ctx, "", &resp.Diagnostics)
	if !ok {
		return
	}
	id, err := poolIdByNameOrId(&c, idParts[0])
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s", err))
		return
	}
	if id == "" {
		resp.Diagnostics.AddError("Pool not found", fmt.Sprintf("There is no pool named %s on the cluster", idParts[0]))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *DataGatewayPoolsResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state DataGatewayPools, params string) (cw.Res, error) {
	var err error
	cacheKey := "DataGatewayPools::" + state.Instance.ValueString() + "::"

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
	res, err := state.readGet(c, "/crosswork/dg-manager/v2/hapool/query"+params)
	singleRes := res
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}

// End of section. //template:end readcache

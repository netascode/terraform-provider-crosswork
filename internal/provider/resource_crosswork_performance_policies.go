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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
var _ resource.Resource = &PerformancePoliciesResource{}
var _ resource.ResourceWithImportState = &PerformancePoliciesResource{}

func NewPerformancePoliciesResource() resource.Resource {
	return &PerformancePoliciesResource{}
}

type PerformancePoliciesResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *PerformancePoliciesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_performance_policies"
}

func (r *PerformancePoliciesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a performance monitoring policy.").String,

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
				MarkdownDescription: helpers.NewAttributeDescription("Policy name, 1 to 255 characters").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"policy_template": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Policy type, for example `deviceHealth` or `INTERFACE`").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Policy description").String,
				Optional:            true,
			},
			"contact": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Contact written into the tag of the policy").String,
				Optional:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether the policy collects data").AddDefaultValueDescription("true").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
			},
			"devices": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Devices to poll, by their inventory id; set one of `devices`, `device_groups` and `port_groups`").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"device_groups": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Ids of the device groups to poll").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"port_groups": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Ids of the port groups to poll").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"schema_intervals": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Polling intervals by schema, one row for each schema of the policy type").String,
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"schema": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Schema of the policy type, for example `CPU`").String,
							Required:            true,
						},
						"seconds": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Polling interval in seconds, 0 for no polling").String,
							Required:            true,
						},
					},
				},
			},
			"thresholds": schema.MapNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Thresholds by metric, for example `cpuUtilization`").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"conditions": schema.ListNestedAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Conditions of the metric, in order").String,
							Required:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"relation": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Comparison of the metric with the value").AddStringEnumDescription("EQ", "GT", "GE", "LT", "LE").String,
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("EQ", "GT", "GE", "LT", "LE"),
										},
									},
									"value": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Value of the threshold").String,
										Required:            true,
									},
									"count": schema.Int64Attribute{
										MarkdownDescription: helpers.NewAttributeDescription("Number of polls that must match before the alarm is raised").String,
										Required:            true,
									},
									"severity": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Severity of the alarm, `INFORMATION` raises an event").AddStringEnumDescription("CRITICAL", "MAJOR", "MINOR", "WARNING", "INFORMATION").String,
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("CRITICAL", "MAJOR", "MINOR", "WARNING", "INFORMATION"),
										},
									},
									"match_cos": schema.SetAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Classes of service of a Quality of service drop metric, 0 to 7").String,
										ElementType:         types.Int64Type,
										Optional:            true,
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

func (r *PerformancePoliciesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *PerformancePoliciesResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// End of section. //template:end model

// Written by hand: the devices of the plan are named by inventory id and sent as uuids.
func (r *PerformancePoliciesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan PerformancePolicies

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
	wire := plan
	if err := wire.devicesToUUIDs(ctx, &c); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to resolve the devices, got error: %s", err))
		return
	}
	body := wire.toBody(ctx, PerformancePolicies{})

	params := ""
	res, err := c.Post(plan.getPath()+params, body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "POST", err, res.String()))
		return
	}
	plan.Id = types.StringValue(res.Get("monitoringPolicy.id").String())

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end create

// Written by hand: the devices of the answer are uuids and are turned into inventory ids.
func (r *PerformancePoliciesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state PerformancePolicies

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
	if err := state.devicesToNames(ctx, &c); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to resolve the devices, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

// Written by hand: the devices of the plan are named by inventory id and sent as uuids; one PUT also switches active.
func (r *PerformancePoliciesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state PerformancePolicies

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

	wire := plan
	if err := wire.devicesToUUIDs(ctx, &c); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to resolve the devices, got error: %s", err))
		return
	}
	body := wire.toBody(ctx, state)
	params := ""
	res, err := c.Put(plan.getPath()+"/"+url.QueryEscape(plan.Id.ValueString())+params, body)
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
func (r *PerformancePoliciesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state PerformancePolicies

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
	res, err := c.Delete(state.getPath() + "/" + url.QueryEscape(state.Id.ValueString()))
	if err != nil && !strings.Contains(err.Error(), "StatusCode 404") && !cw.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s, %s", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Written by hand: the user knows the name of a policy, the controller its number, so the import takes the name and falls back to
// the number.
func (r *PerformancePoliciesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID == "" || strings.Contains(req.ID, ",") {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <name>. Got: %q", req.ID),
		)
		return
	}
	c, ok := r.clientFor(ctx, "", &resp.Diagnostics)
	if !ok {
		return
	}
	id, err := policyIDByName(&c, (&PerformancePolicies{}).getPath(), req.ID)
	if err != nil {
		if _, getErr := c.Get((&PerformancePolicies{}).getPath()+"/"+url.QueryEscape(req.ID), cw.ReadOnly); getErr != nil {
			resp.Diagnostics.AddError("Import Error", err.Error())
			return
		}
		id = req.ID
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *PerformancePoliciesResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state PerformancePolicies, params string) (cw.Res, error) {
	var err error
	cacheKey := "PerformancePolicies::" + state.Instance.ValueString() + "::"

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

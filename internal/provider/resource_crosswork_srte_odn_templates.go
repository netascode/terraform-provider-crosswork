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
var _ resource.Resource = &SRTEODNTemplatesResource{}
var _ resource.ResourceWithImportState = &SRTEODNTemplatesResource{}

func NewSRTEODNTemplatesResource() resource.Resource {
	return &SRTEODNTemplatesResource{}
}

type SRTEODNTemplatesResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
	deviceReady           nsoDeviceConfig
}

func (r *SRTEODNTemplatesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_srte_odn_templates"
}

func (r *SRTEODNTemplatesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage an SR-TE on-demand next-hop (ODN) template. Create and update wait until the devices named in `head_end.name` are ready in Crosswork, as set by the `device_ready_mode` and `device_ready_timeout` arguments of the provider configuration.").String,

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
				MarkdownDescription: helpers.NewAttributeDescription("Template name").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"color": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("SR policy color. Traffic whose route carries this color is steered into the policy that the template creates.").AddIntegerRangeDescription(1, 4294967295).String,
				Required:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 4294967295),
				},
			},
			"head_end": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Head ends of the template.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Name of a device in the Crosswork inventory that is onboarded to NSO").String,
							Required:            true,
						},
					},
				},
			},
			"maximum_sid_depth": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Maximum SID depth.").AddIntegerRangeDescription(1, 255).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 255),
				},
			},
			"pce_group": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("PCE group").String,
				Optional:            true,
			},
			"bandwidth": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Requested bandwidth in kbps. Not supported together with SRv6.").AddIntegerRangeDescription(1, 4294967295).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 4294967295),
				},
			},
			"source_address": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Source address of the policies. An IPv6 address needs SRv6.").String,
				Optional:            true,
			},
			"srv6_locator_name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SRv6 locator name. Setting any SRv6 attribute turns on SRv6 for the template.").String,
				Optional:            true,
			},
			"srv6_locator_behavior": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SRv6 behavior of the locator").AddStringEnumDescription("ub6-insert-reduced").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("ub6-insert-reduced"),
				},
			},
			"srv6_locator_binding_sid_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SRv6 binding SID type").AddStringEnumDescription("srv6-dynamic").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("srv6-dynamic"),
				},
			},
			"dynamic_pce": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Delegate the path computation to the PCE. Set to `true` or leave unset.").String,
				Optional:            true,
			},
			"dynamic_flex_alg": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Flexible algorithm (prefix-SID algorithm)").AddIntegerRangeDescription(128, 255).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(128, 255),
				},
			},
			"dynamic_metric_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Metric type of the path computation").AddStringEnumDescription("hopcount", "igp", "latency", "te").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("hopcount", "igp", "latency", "te"),
				},
			},
			"dynamic_metric_margin_absolute": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Absolute metric margin. Use this or the relative margin.").AddIntegerRangeDescription(0, 2147483647).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 2147483647),
				},
			},
			"dynamic_metric_margin_relative": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Relative metric margin in percent. Use this or the absolute margin.").AddIntegerRangeDescription(0, 2147483647).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 2147483647),
				},
			},
			"dynamic_affinity_rules": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Affinity rules of the path computation").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"action": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Rule action").AddStringEnumDescription("exclude-any", "include-all", "include-any").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("exclude-any", "include-all", "include-any"),
							},
						},
						"color": schema.SetAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Names of the affinity colors of the rule. NSO returns them sorted.").String,
							ElementType:         types.StringType,
							Optional:            true,
						},
					},
				},
			},
			"dynamic_segments_sid_algorithm": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("SID algorithm constraint").AddIntegerRangeDescription(128, 255).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(128, 255),
				},
			},
			"dynamic_disjoint_path_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Disjointness type. With `dynamic_disjoint_path_group_id` it makes the paths of the group disjoint.").AddStringEnumDescription("link", "node", "srlg", "srlg-node").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("link", "node", "srlg", "srlg-node"),
				},
			},
			"dynamic_disjoint_path_group_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Disjoint group id. Required with the type.").AddIntegerRangeDescription(1, 65535).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"dynamic_disjoint_path_sub_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Disjoint sub group id").AddIntegerRangeDescription(1, 65535).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 65535),
				},
			},
			"delay_measurement_profile": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Delay measurement profile of the head end. Not `default`; leave it unset for the default profile.").String,
				Optional:            true,
			},
			"delay_measurement_logging_delay_exceeded": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Log when the delay is exceeded. Set to `true` or leave unset.").String,
				Optional:            true,
			},
			"liveness_detection_profile": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Liveness detection profile. Not `default`. Use this or delay measurement.").String,
				Optional:            true,
			},
			"liveness_detection_backup": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Backup liveness detection profile. Needs a profile.").String,
				Optional:            true,
			},
			"liveness_detection_invalidation_action": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Action when the liveness session invalidates the path").AddStringEnumDescription("down", "none").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("down", "none"),
				},
			},
			"liveness_detection_logging_session_state_change": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Log a change of the liveness session state. Set to `true` or leave unset.").String,
				Optional:            true,
			},
			"performance_measurement_reverse_path_label": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Reverse path label of the performance measurement").AddIntegerRangeDescription(16, 1048575).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(16, 1048575),
				},
			},
		},
	}
}

func (r *SRTEODNTemplatesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
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
func (r *SRTEODNTemplatesResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create
func (r *SRTEODNTemplatesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan SRTEODNTemplates

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
	body := plan.toBody(ctx, SRTEODNTemplates{})

	params := ""
	res, err := c.Post("/crosswork/proxy/nso/restconf/data/cisco-sr-te-cfp:sr-te/cisco-sr-te-cfp-sr-odn:odn"+params, body, cw.PlanPath(plan.getPathPlan()))
	if err != nil && (cw.IsPlanTimeout(err) || cw.IsPlanFailed(err)) {
		// NSO accepted the object, but its plan did not complete. Keep the object in the state, which the error marks
		// as tainted, so that the next apply replaces or removes it instead of failing with "object already exists".
		plan.Id = types.StringValue(fmt.Sprint(plan.Name.ValueString()))
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("The object was created, but its plan did not complete, so it is tracked as tainted: %s", err))
		return
	}
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
func (r *SRTEODNTemplatesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state SRTEODNTemplates

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
func (r *SRTEODNTemplatesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state SRTEODNTemplates

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
func (r *SRTEODNTemplatesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state SRTEODNTemplates

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
func (r *SRTEODNTemplatesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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
func (r *SRTEODNTemplatesResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state SRTEODNTemplates, params string) (cw.Res, error) {
	var err error
	cacheKey := "SRTEODNTemplates::" + state.Instance.ValueString() + "::"

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

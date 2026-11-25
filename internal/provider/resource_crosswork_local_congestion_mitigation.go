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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
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
var _ resource.Resource = &LocalCongestionMitigationResource{}
var _ resource.ResourceWithImportState = &LocalCongestionMitigationResource{}

func NewLocalCongestionMitigationResource() resource.Resource {
	return &LocalCongestionMitigationResource{}
}

type LocalCongestionMitigationResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *LocalCongestionMitigationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_local_congestion_mitigation"
}

func (r *LocalCongestionMitigationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage the local congestion mitigation configuration of one LCM domain.").String,

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
			"domain_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("LCM domain identifier").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Enables LCM for the domain").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"color": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Color assigned sequentially to tactical SR policies starting with this value").AddIntegerRangeDescription(1, 4294967294).AddDefaultValueDescription("2000").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 4294967294),
				},
				Default: int64default.StaticInt64(2000),
			},
			"utilization_threshold": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Utilization percentage at which an interface is considered congested").AddIntegerRangeDescription(0, 100).AddDefaultValueDescription("100").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
				Default: int64default.StaticInt64(100),
			},
			"utilization_hold_margin": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Percentage that dampens the removal of deployed tactical SR policies").AddIntegerRangeDescription(0, 100).AddDefaultValueDescription("5").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
				Default: int64default.StaticInt64(5),
			},
			"delete_lsps": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Delete all deployed tactical SR policies when LCM is disabled").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"profile_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Profile ID assigned to all tactical SR policies; 0 assigns none").AddIntegerRangeDescription(0, 65534).AddDefaultValueDescription("0").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 65534),
				},
				Default: int64default.StaticInt64(0),
			},
			"congestion_check_interval": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Interval in seconds used to evaluate the network for congestion").AddIntegerRangeDescription(60, 86400).AddDefaultValueDescription("900").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(60, 86400),
				},
				Default: int64default.StaticInt64(900),
			},
			"max_tte_group": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Maximum number of tactical policies used to mitigate a single interface").AddIntegerRangeDescription(1, 8).AddDefaultValueDescription("8").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 8),
				},
				Default: int64default.StaticInt64(8),
			},
			"include_all_interfaces": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Monitor congestion on all interfaces; false monitors the selected interfaces").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Description of the domain configuration").AddDefaultValueDescription("LCM startup config").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 256),
				},
				Default: stringdefault.StaticString("LCM startup config"),
			},
			"auto_repair_solution": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Detect and deploy changes automatically after network failures").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"stay_in_area": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Restrict bypass paths to the mitigated OSPF area or ISIS level").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"adjacency_hop_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Type of adjacency SID used when creating an SR policy").AddStringEnumDescription("UNPROTECTED", "PROTECTED", "UNPROTECTED_PREFERRED", "PROTECTED_PREFERRED").AddDefaultValueDescription("UNPROTECTED").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("UNPROTECTED", "PROTECTED", "UNPROTECTED_PREFERRED", "PROTECTED_PREFERRED"),
				},
				Default: stringdefault.StaticString("UNPROTECTED"),
			},
			"operation_mode": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Operation mode of LCM; AUTOMATED needs gRPC on every device of the domain").AddStringEnumDescription("MANUAL", "AUTOMATED").AddDefaultValueDescription("MANUAL").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("MANUAL", "AUTOMATED"),
				},
				Default: stringdefault.StaticString("MANUAL"),
			},
			"objective": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Path computation objective of the tactical SR policies").AddStringEnumDescription("MIN_IGP_METRIC", "MIN_TE_METRIC", "MIN_DELAY").AddDefaultValueDescription("MIN_IGP_METRIC").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("MIN_IGP_METRIC", "MIN_TE_METRIC", "MIN_DELAY"),
				},
				Default: stringdefault.StaticString("MIN_IGP_METRIC"),
			},
			"pce_dispatch_timeout": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Maximum time in seconds to confirm the deployment of tactical SR policies").AddIntegerRangeDescription(10, 300).AddDefaultValueDescription("180").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(10, 300),
				},
				Default: int64default.StaticInt64(180),
			},
			"congestion_check_suspension_interval": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Seconds to suspend congestion detection after any tactical SR policy change").AddIntegerRangeDescription(600, 3600).AddDefaultValueDescription("600").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(600, 3600),
				},
				Default: int64default.StaticInt64(600),
			},
			"overprovision_factor": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Percentage deducted from interface thresholds to derive the effective threshold").AddIntegerRangeDescription(0, 10).AddDefaultValueDescription("0").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 10),
				},
				Default: int64default.StaticInt64(0),
			},
			"uneven_ecmp_traffic_threshold": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Sensitivity percentage to detect uneven traffic across solution bypass tunnels").AddIntegerRangeDescription(0, 100).AddDefaultValueDescription("0").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 100),
				},
				Default: int64default.StaticInt64(0),
			},
			"throttle_mode_threshold": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Solution oscillations per hour that put the solution in the PAUSED state; 0 disables detection").AddIntegerRangeDescription(0, 10).AddDefaultValueDescription("5").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 10),
				},
				Default: int64default.StaticInt64(5),
			},
			"history_retention_time": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Days the collected data is retained").AddIntegerRangeDescription(1, 90).AddDefaultValueDescription("30").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 90),
				},
				Default: int64default.StaticInt64(30),
			},
			"debug_opt": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Save debug plan files and set the log level to debug").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"debug_opt_max_files": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Maximum number of debug plan file sets to save").AddIntegerRangeDescription(1, 512).AddDefaultValueDescription("40").String,
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 512),
				},
				Default: int64default.StaticInt64(40),
			},
			"affinity_include_any": schema.ListAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Bit positions of the link affinities of the include any constraint").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"affinity_exclude_any": schema.ListAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Bit positions of the link affinities of the exclude any constraint").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"affinity_include_all": schema.ListAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Bit positions of the link affinities of the include all constraint").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"node_tag_max_sid_depth": schema.MapAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Maximum SID depth (0 to 16) by device tag, written as a text number").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
		},
	}
}

func (r *LocalCongestionMitigationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *LocalCongestionMitigationResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create
func (r *LocalCongestionMitigationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan LocalCongestionMitigation

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
	body := plan.toBody(ctx, LocalCongestionMitigation{})

	params := ""
	res, err := c.Put(plan.getPath()+params, body, cw.UseMutex)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "PUT", err, res.String()))
		return
	}
	plan.Id = types.StringValue(fmt.Sprint(plan.DomainId.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end create

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (r *LocalCongestionMitigationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state LocalCongestionMitigation

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
	res, err := state.readGet(&c, "/crosswork/optima/api/v1/functionPacks/lcm/domain/config"+params)
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
func (r *LocalCongestionMitigationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state LocalCongestionMitigation

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
func (r *LocalCongestionMitigationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state LocalCongestionMitigation

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
	res, err := c.Put(state.getPathDelete(), strings.ReplaceAll(`{"id":"{id}","config":{"enable":false,"autoRepairSolution":false,"stayInArea":false,"adjacencyHopType":"UNPROTECTED","utilThreshold":100,"utilHoldMargin":5,"color":2000,"profileId":0,"maxTteGroup":8,"delLsps":false,"includeAllInterfaces":false,"pceDispatchTimeout":180,"overprovisionFactor":0,"unevenEcmpTrafficThreshold":0,"throttleModeThreshold":5,"objective":"MIN_IGP_METRIC","operationMode":"MANUAL","congestionCheckInterval":900,"congestionCheckSuspensionInterval":600,"description":"LCM startup config","debugOpt":false,"debugOptMaxFiles":40,"nodeTagMaxSidDepth":[],"affinityConstraint":{"excludeAny":[],"includeAny":[],"includeAll":[]},"historyRetentionTime":30}}`, "{id}", state.Id.ValueString()), cw.UseMutex)
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
func (r *LocalCongestionMitigationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <domain_id>. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("domain_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *LocalCongestionMitigationResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state LocalCongestionMitigation, params string) (cw.Res, error) {
	var err error
	cacheKey := "LocalCongestionMitigation::" + state.Instance.ValueString() + "::"

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
	res, err := state.readGet(c, "/crosswork/optima/api/v1/functionPacks/lcm/domain/config"+params)
	singleRes := res
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}

// End of section. //template:end readcache

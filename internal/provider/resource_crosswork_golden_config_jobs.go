// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

var _ resource.Resource = &GoldenConfigJobsResource{}
var _ resource.ResourceWithValidateConfig = &GoldenConfigJobsResource{}

func NewGoldenConfigJobsResource() resource.Resource {
	return &GoldenConfigJobsResource{}
}

type GoldenConfigJobsResource struct {
	clients     map[string]*cw.Client
	deviceReady nsoDeviceConfig
}

func (r *GoldenConfigJobsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_golden_config_jobs"
}

func (r *GoldenConfigJobsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This resource applies a Golden Configuration template to devices and keeps it applied.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"template_name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the template to apply").String,
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"template_config_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Integer id of the configuration of the template (`config_id` of the template resource).").String,
				Required:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.RequiresReplace()},
			},
			"vendor": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Vendor sent with the job, as the template has it").String,
				Required:            true,
			},
			"product_series": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Product series sent with the job").String,
				Required:            true,
			},
			"mop": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the MOP that applies the template, for example `Default Golden Config - Cisco`").String,
				Required:            true,
			},
			"batch_size": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Number of devices the controller configures in parallel; the controller uses 1 when it is not set").String,
				Optional:            true,
				Validators:          []validator.Int64{int64validator.AtLeast(1)},
			},
			"error_budget": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Number of failures the controller tolerates before it stops a job; unset behaves like 1").String,
				Optional:            true,
				Validators:          []validator.Int64{int64validator.AtLeast(0)},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Tags of the jobs").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
			"deployment_timeout": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Seconds to wait for each job to finish, counted from when the job starts. It only limits how long Terraform waits and is not part of the Crosswork configuration. When it passes the apply fails and the job keeps running on the controller; the next apply runs it again").AddIntegerRangeDescription(10, 86400).AddDefaultValueDescription("900").String,
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(900),
				Validators:          []validator.Int64{int64validator.Between(10, 86400)},
			},
			"trigger": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("A value, such as a hash of the template and its global variables, that triggers applying the template again to all devices when it changes.").String,
				Optional:            true,
			},
			"devices": schema.MapNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The devices the template is applied to, by inventory id, each with its associated values. A device that leaves the map has the configuration removed").String,
				Required:            true,
				Validators:          []validator.Map{mapvalidator.SizeAtLeast(1)},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"variables": schema.ListNestedAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Device level variables, which are preferred over the template and the global variables").String,
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Variable name").String,
										Required:            true,
										Validators:          []validator.String{stringvalidator.RegexMatches(gcVariableName, "must be 1 to 64 letters, digits and underscores and must not start with a digit")},
									},
									"type": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Type of the value").AddStringEnumDescription("STRING", "JSON").AddDefaultValueDescription("STRING").String,
										Optional:            true,
										Computed:            true,
										Default:             stringdefault.StaticString("STRING"),
										Validators:          []validator.String{stringvalidator.OneOf("STRING", "JSON")},
									},
									"value": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Value; for the type `JSON` a JSON document").AddMutualExclusivityDescription("**Required**: exactly one of `value` and `value_wo` must be set.").String,
										Optional:            true,
									},
									"value_wo": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("Value").AddMutualExclusivityDescription("**Required**: exactly one of `value` and `value_wo` must be set.").String,
										Optional:            true,
										WriteOnly:           true,
										Sensitive:           true,
									},
									"value_wo_version": schema.Int64Attribute{
										MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `value_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new value.").String,
										Optional:            true,
										Validators:          []validator.Int64{int64validator.AtLeast(1)},
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

func (r *GoldenConfigJobsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.deviceReady = req.ProviderData.(*CwProviderData).DeviceReady
}

func (r *GoldenConfigJobsResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

func (r *GoldenConfigJobsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var cfg GoldenConfigJobs
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() || cfg.Devices.IsNull() || cfg.Devices.IsUnknown() {
		return
	}
	var devices map[string]GoldenConfigJobDevice
	resp.Diagnostics.Append(cfg.Devices.ElementsAs(ctx, &devices, false)...)
	for _, name := range sortedKeys(devices) {
		resp.Diagnostics.Append(validateGcVariables(ctx, devices[name].Variables, "device "+name)...)
	}
}

type gcOutcome struct {
	state    GoldenConfigJobs
	failures map[string]string // device name to the reason, for the devices that did not take the template
	removals map[string]error  // device name to the error, for the devices whose configuration could not be removed
}

func (o gcOutcome) failed() bool { return len(o.failures) > 0 || len(o.removals) > 0 }

func (o gcOutcome) message() string {
	var parts []string
	if len(o.failures) > 0 {
		parts = append(parts, fmt.Sprintf("%d device(s) did not take the template:\n%s", len(o.failures), gcFailureMessage(o.failures)))
	}
	if len(o.removals) > 0 {
		lines := make([]string, 0, len(o.removals))
		for _, d := range sortedKeys(o.removals) {
			lines = append(lines, "  "+d+": "+o.removals[d].Error())
		}
		parts = append(parts, fmt.Sprintf("the configuration could not be removed from %d device(s):\n%s", len(o.removals), strings.Join(lines, "\n")))
	}
	return strings.Join(parts, "\n")
}

// The devices that are new, or whose values changed, get jobs, one per group of equal values; when the trigger changed, the devices that did not
// change get a re-apply, and a device that left has its application deleted. The state it returns holds the devices that are in place: all of
// them when nothing failed, otherwise those that succeeded plus those that were applied before.
func (r *GoldenConfigJobsResource) apply(ctx context.Context, c *cw.Client, plan GoldenConfigJobs, cfgDevices types.Map, old *GoldenConfigJobs) (gcOutcome, diag.Diagnostics) {
	var diags diag.Diagnostics
	out := gcOutcome{state: plan, failures: map[string]string{}, removals: map[string]error{}}

	planVars, d := gcDeviceVars(ctx, plan.Devices, cfgDevices)
	diags.Append(d...)
	oldVars := map[string][]gcVar{}
	if old != nil {
		ov, d := gcDeviceVars(ctx, old.Devices, old.Devices)
		diags.Append(d...)
		oldVars = ov
	}
	if diags.HasError() {
		return out, diags
	}

	var changed, unchanged, gone []string
	for name, vars := range planVars {
		if prev, ok := oldVars[name]; ok && gcVarsNameKey(prev) == gcVarsNameKey(vars) {
			unchanged = append(unchanged, name)
		} else {
			changed = append(changed, name)
		}
	}
	for name := range oldVars {
		if _, ok := planVars[name]; !ok {
			gone = append(gone, name)
		}
	}
	sort.Strings(changed)
	sort.Strings(unchanged)
	sort.Strings(gone)
	triggerChanged := old != nil && !plan.Trigger.Equal(old.Trigger)
	timeout := time.Duration(plan.DeploymentTimeout.ValueInt64()) * time.Second

	var apps []gcApplication
	appsRead := false
	applications := func() ([]gcApplication, error) {
		if !appsRead {
			var err error
			// The applications are listed under the id of the template; the job names the id of its configuration.
			id, found, err := gcFindTemplate(c, plan.TemplateName.ValueString())
			if err != nil {
				return nil, err
			}
			if !found {
				return nil, fmt.Errorf("the template %q does not exist", plan.TemplateName.ValueString())
			}
			if apps, err = gcTemplateApplications(c, id); err != nil {
				return nil, err
			}
			appsRead = true
		}
		return apps, nil
	}

	// A device that left: its configuration goes first, so a failing job later does not hide it.
	if len(gone) > 0 {
		list, err := applications()
		if err != nil {
			diags.AddError("Client Error", fmt.Sprintf("Failed to read the applications of the template (POST), got error: %s", err))
			return out, diags
		}
		var leaving []gcApplication
		for _, a := range list {
			if contains(gone, a.DeviceName) {
				leaving = append(leaving, a)
			}
		}
		out.removals = gcDeleteApplications(c, leaving)
	}

	// The devices must be ready before a job runs on them, as the provider setting device_ready_mode says.
	touched := append([]string(nil), changed...)
	if triggerChanged {
		touched = append(touched, unchanged...)
	}
	var nodes map[string]nsoNode
	if len(touched) > 0 {
		var notReady []string
		var err error
		cfg := r.deviceReady
		if cfg.Mode == nsoModeNone {
			// A job cannot run without the uuid of the device, so the devices must still be in the inventory, but nothing waits.
			cfg.Mode, cfg.Timeout = nsoModePresent, 0
		}
		devices := nsoDeviceSet{}
		for _, name := range touched {
			devices[name] = struct{}{}
		}
		nodes, notReady, err = waitNSODevices(ctx, c, plan.Instance.ValueString(), cfg, devices)
		if err != nil {
			diags.AddError("Client Error", fmt.Sprintf("Failed to read the inventory (POST), got error: %s", err))
			return out, diags
		}
		for _, line := range notReady {
			name := strings.SplitN(line, ":", 2)[0]
			out.failures[name] = "not ready: " + strings.TrimPrefix(line, name+": ")
		}
	}
	ready := func(names []string) []string {
		var ok []string
		for _, n := range names {
			if _, bad := out.failures[n]; !bad {
				ok = append(ok, n)
			}
		}
		return ok
	}

	if todo := ready(changed); len(todo) > 0 {
		mopID, err := gcMopID(c, plan.Mop.ValueString())
		if err != nil {
			diags.AddError("Client Error", fmt.Sprintf("Failed to resolve the MOP (POST), got error: %s", err))
			return out, diags
		}
		tags, d := stringsOf(ctx, plan.Tags)
		diags.Append(d...)
		p := gcJobParams{ConfigID: plan.TemplateConfigId.ValueInt64(), TemplateName: plan.TemplateName.ValueString(), Vendor: plan.Vendor.ValueString(),
			Series: plan.ProductSeries.ValueString(), Mop: plan.Mop.ValueString(), MopID: mopID, Tags: tags}
		if !plan.BatchSize.IsNull() {
			v := plan.BatchSize.ValueInt64()
			p.Batch = &v
		}
		if !plan.ErrorBudget.IsNull() {
			v := plan.ErrorBudget.ValueInt64()
			p.ErrorBudget = &v
		}
		results := gcRunJobs(ctx, c, p, gcMakeGroups(plan.TemplateName.ValueString(), todo, planVars), nodes, timeout)
		for _, res := range results {
			for dev, why := range res.Failed {
				out.failures[dev] = why
			}
		}
	}

	reapplyFailed := false
	if triggerChanged {
		if again := ready(unchanged); len(again) > 0 {
			list, err := applications()
			if err != nil {
				diags.AddError("Client Error", fmt.Sprintf("Failed to read the applications of the template (POST), got error: %s", err))
				return out, diags
			}
			var targets []gcApplication
			byUUID := map[string]string{}
			for _, a := range list {
				if contains(again, a.DeviceName) {
					targets = append(targets, a)
					byUUID[a.DeviceUUID] = a.DeviceName
				}
			}
			for uuid, why := range gcReapply(c, targets, timeout) {
				name := byUUID[uuid]
				if name == "" {
					name = uuid
				}
				out.failures[name] = "re-apply: " + why
				reapplyFailed = true
			}
		}
	}

	// A device that did not change and was not re-applied (not ready, or its re-apply failed) still lacks the new
	// trigger: keep the old one, so the next plan shows the change again.
	if triggerChanged {
		for _, n := range unchanged {
			if _, bad := out.failures[n]; bad {
				reapplyFailed = true
			}
		}
	}

	var fromPlan, fromOld []string
	for name := range planVars {
		if _, bad := out.failures[name]; !bad {
			fromPlan = append(fromPlan, name)
		} else if _, was := oldVars[name]; was {
			fromOld = append(fromOld, name)
		}
	}
	for name := range out.removals {
		fromOld = append(fromOld, name)
	}
	sort.Strings(fromPlan)
	sort.Strings(fromOld)
	merged := gcSubsetDevices(fromPlan, plan.Devices)
	if len(fromOld) > 0 && old != nil {
		all := merged.Elements()
		for k, v := range gcSubsetDevices(fromOld, old.Devices).Elements() {
			all[k] = v
		}
		merged = types.MapValueMust(goldenConfigJobDeviceType, all)
	}
	out.state.Devices = merged
	if reapplyFailed && old != nil {
		out.state.Trigger = old.Trigger
	}
	return out, diags
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func (r *GoldenConfigJobsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, cfg GoldenConfigJobs
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.TemplateName.ValueString()))

	out, diags := r.apply(ctx, &c, plan, cfg.Devices, nil)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if out.failed() {
		// No state is saved: a resource saved with an error is tainted, and the next apply would destroy it and
		// remove the configuration from the devices that took it. With no state the next apply runs the jobs again,
		// which leaves those devices alone and retries the others.
		resp.Diagnostics.AddError("Golden Config job failed",
			out.message()+"\nNo state was saved; the devices that took the template keep it, and the next apply runs the jobs again. The controller keeps the applications, so destroying now would not remove them.")
		return
	}
	out.state.Id = out.state.TemplateName
	resp.Diagnostics.Append(resp.State.Set(ctx, &out.state)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.TemplateName.ValueString()))
}

func (r *GoldenConfigJobsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state GoldenConfigJobs
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	name := state.Id.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", name))

	id, found, err := gcFindTemplate(&c, state.TemplateName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (POST), got error: %s", err))
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	apps, err := gcTemplateApplications(&c, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the applications (POST), got error: %s", err))
		return
	}
	present := map[string]bool{}
	for _, a := range apps {
		present[a.DeviceName] = true
	}
	var keep []string
	for name := range state.Devices.Elements() {
		if present[name] {
			keep = append(keep, name)
		}
	}
	if len(keep) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}
	sort.Strings(keep)
	// A device whose application is gone leaves the state, so the plan applies the template to it again.
	state.Devices = gcSubsetDevices(keep, state.Devices)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", name))
}

func (r *GoldenConfigJobsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state, cfg GoldenConfigJobs
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))

	out, diags := r.apply(ctx, &c, plan, cfg.Devices, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	out.state.Id = state.Id
	// An update that fails still saves what is in place: an update is never tainted, and the next apply
	// retries only the devices that failed.
	resp.Diagnostics.Append(resp.State.Set(ctx, &out.state)...)
	if out.failed() {
		resp.Diagnostics.AddError("Golden Config job failed", out.message())
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))
}

func (r *GoldenConfigJobsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state GoldenConfigJobs
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))

	// The applications are found through the template, so one left by a device that was already deleted from
	// the inventory is removed too. Without the template there is nothing applied.
	id, found, err := gcFindTemplate(&c, state.TemplateName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (POST), got error: %s", err))
		return
	}
	if !found {
		return
	}
	apps, err := gcTemplateApplications(&c, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the applications (POST), got error: %s", err))
		return
	}
	var mine []gcApplication
	for _, a := range apps {
		if _, ok := state.Devices.Elements()[a.DeviceName]; ok {
			mine = append(mine, a)
		}
	}
	if failed := gcDeleteApplications(&c, mine); len(failed) > 0 {
		lines := make([]string, 0, len(failed))
		for _, d := range sortedKeys(failed) {
			lines = append(lines, "  "+d+": "+failed[d].Error())
		}
		resp.Diagnostics.AddError("Client Error", "Failed to remove the configuration from the devices (DELETE):\n"+strings.Join(lines, "\n"))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))
}

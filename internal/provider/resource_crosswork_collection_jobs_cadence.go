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

import (
	"context"
	"fmt"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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

var _ resource.Resource = &CollectionJobsCadenceResource{}
var _ resource.ResourceWithImportState = &CollectionJobsCadenceResource{}

func NewCollectionJobsCadenceResource() resource.Resource {
	return &CollectionJobsCadenceResource{}
}

type CollectionJobsCadenceResource struct {
	clients map[string]*cw.Client
}

func (r *CollectionJobsCadenceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection_jobs_cadence"
}

func (r *CollectionJobsCadenceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage the cadence of the Crosswork collection jobs.").String,

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
			"reachability_check": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Seconds between two reachability checks of the devices, 300 (5 minutes) to 86400 (24 hours). Leave it unset to keep the value the cluster holds").AddIntegerRangeDescription(300, 86400).String,
				Optional:            true,
				Validators:          []validator.Int64{int64validator.Between(300, 86400)},
			},
			"snmp_check": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Seconds between two SNMP checks of the devices, 600 (10 minutes) to 86400 (24 hours). Leave it unset to keep the value the cluster holds").AddIntegerRangeDescription(600, 86400).String,
				Optional:            true,
				Validators:          []validator.Int64{int64validator.Between(600, 86400)},
			},
			"te_tunnel_id": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Seconds between two collections of the TE tunnel ids, 600 (10 minutes) to 86400 (24 hours). Leave it unset to keep the value the cluster holds").AddIntegerRangeDescription(600, 86400).String,
				Optional:            true,
				Validators:          []validator.Int64{int64validator.Between(600, 86400)},
			},
			"show_clock": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Seconds between two collections of the clock of the devices (`show clock`), 600 (10 minutes) to 86400 (24 hours). Leave it unset to keep the value the cluster holds").AddIntegerRangeDescription(600, 86400).String,
				Optional:            true,
				Validators:          []validator.Int64{int64validator.Between(600, 86400)},
			},
		},
	}
}

func (r *CollectionJobsCadenceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clients = req.ProviderData.(*CwProviderData).Clients
}

func (r *CollectionJobsCadenceResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// write sends the plan with the values it leaves out taken from the cluster, which the controller needs, and waits until the
// controller shows them.
func (r *CollectionJobsCadenceResource) write(ctx context.Context, plan CollectionJobsCadence, verb string, diags *diag.Diagnostics) bool {
	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), diags)
	if !ok {
		return false
	}
	current, err := readCadence(&c)
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (POST) before %s, got error: %s", verb, err))
		return false
	}
	body, err := cadenceBody(cadenceFromModel(plan), current)
	if err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s", err))
		return false
	}
	if err := putCadence(&c, body); err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s", err))
		return false
	}
	return true
}

// Create and Update are the same call: the setting always exists.
func (r *CollectionJobsCadenceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan CollectionJobsCadence
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "Beginning Create")
	if !r.write(ctx, plan, "create", &resp.Diagnostics) {
		return
	}
	plan.Id = types.StringValue(cadenceId)
	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// After an import, which names no values, all four are kept.
func (r *CollectionJobsCadenceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state CollectionJobsCadence
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", state.Id.String()))

	current, err := readCadence(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (POST), got error: %s", err))
		return
	}
	state = state.with(cadenceKeep(cadenceFromModel(state), current))

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CollectionJobsCadenceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state CollectionJobsCadence
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))
	if !r.write(ctx, plan, "update", &resp.Diagnostics) {
		return
	}
	plan.Id = state.Id
	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete only leaves the state: Crosswork has no call that removes the setting, and the values stay as they are.
func (r *CollectionJobsCadenceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state CollectionJobsCadence
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Delete leaves the values on the cluster and removes the state", state.Id.ValueString()))
	resp.State.RemoveResource(ctx)
}

func (r *CollectionJobsCadenceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	if req.ID != cadenceId {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: %s. Got: %q", cadenceId, req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

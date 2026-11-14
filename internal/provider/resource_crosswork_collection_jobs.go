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
	"strings"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

var _ resource.Resource = &CollectionJobsResource{}
var _ resource.ResourceWithImportState = &CollectionJobsResource{}

func NewCollectionJobsResource() resource.Resource {
	return &CollectionJobsResource{}
}

type CollectionJobsResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
}

func (r *CollectionJobsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection_jobs"
}

func (r *CollectionJobsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a bulk collection job of Crosswork.").String,

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
			"application_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Application id of the job, which with the context id identifies it. A change replaces the job").String,
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"context_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Context id of the job, unique within the application id. A change replaces the job").String,
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"collector_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Collector type of the job. A job holds only the sensors of its collector type, `cli_sensors` for `CLI` and `snmp_sensors` for `SNMP`. A change replaces the job").AddStringEnumDescription("CLI", "SNMP").String,
				Required:            true,
				Validators:          []validator.String{stringvalidator.OneOf("CLI", "SNMP")},
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"devices": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Inventory ids of the devices the job collects from, in any order. Set either this or `device_tag`. Crosswork does not accept a change of the devices of a job, so a change replaces the job").String,
				ElementType:         types.StringType,
				Optional:            true,
				PlanModifiers:       []planmodifier.Set{setplanmodifier.RequiresReplace()},
			},
			"device_tag": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Tag whose devices the job collects from, resolved by Crosswork when the job runs. Set either this or `devices`. A change replaces the job").String,
				Optional:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"destination": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the data destination that receives the output, as `crosswork_data_gateway_destinations` names it").String,
				Required:            true,
			},
			"cli_sensors": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("CLI commands the job runs on the devices, for the collector type `CLI`. Each command has its own cadence and topic").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"command": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("CLI command to run, for example `show version`").String,
							Required:            true,
						},
						"cadence": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Seconds between two runs of the command, at least 1").String,
							Required:            true,
							Validators:          []validator.Int64{int64validator.AtLeast(1)},
						},
						"topic": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Topic the output of the command is sent to on the destination").String,
							Required:            true,
						},
					},
				},
			},
			"snmp_sensors": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP objects the job collects from the devices, for the collector type `SNMP`. Each object has its own cadence and topic").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"oid": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("OID to collect, one for each sensor").String,
							Required:            true,
						},
						"operation": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("SNMP operation").AddStringEnumDescription("WALK", "SCALAR", "TABLE", "COLUMN").String,
							Required:            true,
							Validators:          []validator.String{stringvalidator.OneOf("WALK", "SCALAR", "TABLE", "COLUMN")},
						},
						"cadence": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Seconds between two collections of the OID, at least 1").String,
							Required:            true,
							Validators:          []validator.Int64{int64validator.AtLeast(1)},
						},
						"topic": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Topic the output of the object is sent to on the destination").String,
							Required:            true,
						},
					},
				},
			},
		},
	}
}

func (r *CollectionJobsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
}

func (r *CollectionJobsResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// The same PUT replaces a job that exists, so a job found under the id is left alone unless allow_existing_on_create is set; the status query is
// asked, because it stops listing a deleted job before the configuration query does.
func (r *CollectionJobsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan CollectionJobs
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	id := collectionJobId(plan.ApplicationId.ValueString(), plan.ContextId.ValueString())
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", id))

	listed, err := collectionJobListed(&c, plan.ApplicationId.ValueString(), plan.ContextId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to check for the job (POST), got error: %s", err))
		return
	}
	if listed && !r.AllowExistingOnCreate {
		resp.Diagnostics.AddError("Collection job exists", fmt.Sprintf("A collection job %s already exists. Import it with `terraform import`, or set allow_existing_on_create on the provider to replace its sensors and destination.", id))
		return
	}

	names, uuids, destinationId, diags := collectionResolve(ctx, &c, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := putCollectionJob(&c, collectionJobBody(plan, names, uuids, destinationId)); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s", err))
		return
	}
	plan.Id = types.StringValue(id)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", id))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CollectionJobsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state CollectionJobs
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", state.Id.String()))

	row, found, err := readCollectionJob(&c, state.ApplicationId.ValueString(), state.ContextId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (POST), got error: %s", err))
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	resp.Diagnostics.Append(collectionJobFromRow(ctx, &c, row, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// The controller replaces the sensors and the destination; the attributes it refuses to change replace the resource and never reach here.
func (r *CollectionJobsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state CollectionJobs
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))

	names, uuids, destinationId, diags := collectionResolve(ctx, &c, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := putCollectionJob(&c, collectionJobBody(plan, names, uuids, destinationId)); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s", err))
		return
	}
	plan.Id = state.Id

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *CollectionJobsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state CollectionJobs
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))

	gone, err := deleteCollectionJob(ctx, &c, state.ApplicationId.ValueString(), state.ContextId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s", err))
		return
	}
	if !gone {
		resp.Diagnostics.AddWarning("Collection job still listed", fmt.Sprintf("Crosswork accepted the delete of the collection job %s, but it still lists the job. It is removed from the state.", state.Id.ValueString()))
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))
	resp.State.RemoveResource(ctx)
}

func (r *CollectionJobsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.SplitN(req.ID, ",", 2)
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <application_id>,<context_id>. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("application_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("context_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
}

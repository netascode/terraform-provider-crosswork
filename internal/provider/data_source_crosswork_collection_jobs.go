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
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

var _ datasource.DataSource = &CollectionJobsDataSource{}
var _ datasource.DataSourceWithConfigure = &CollectionJobsDataSource{}

func NewCollectionJobsDataSource() datasource.DataSource {
	return &CollectionJobsDataSource{}
}

type CollectionJobsDataSource struct {
	clients map[string]*cw.Client
}

func (d *CollectionJobsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_collection_jobs"
}

func (d *CollectionJobsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This data source can read a bulk collection job of Crosswork by its application id and context id, with its status.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"application_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Application id of the job, which with the context id identifies it").String,
				Required:            true,
			},
			"context_id": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Context id of the job, unique within the application id").String,
				Required:            true,
			},
			"collector_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Collector type of the job, `CLI` or `SNMP` for a job this provider can create; another type is shown as the controller names it").String,
				Computed:            true,
			},
			"devices": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Inventory ids of the devices the job collects from").String,
				ElementType:         types.StringType,
				Computed:            true,
			},
			"device_tag": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Tag whose devices the job collects from").String,
				Computed:            true,
			},
			"destination": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Name of the data destination that receives the output").String,
				Computed:            true,
			},
			"cli_sensors": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("CLI commands the job runs on the devices").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"command": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("CLI command to run").String,
							Computed:            true,
						},
						"cadence": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Seconds between two runs of the command").String,
							Computed:            true,
						},
						"topic": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Topic the output of the command is sent to on the destination").String,
							Computed:            true,
						},
					},
				},
			},
			"snmp_sensors": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMP objects the job collects from the devices").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"oid": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("OID to collect").String,
							Computed:            true,
						},
						"operation": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("SNMP operation").String,
							Computed:            true,
						},
						"cadence": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Seconds between two collections of the OID").String,
							Computed:            true,
						},
						"topic": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Topic the output of the object is sent to on the destination").String,
							Computed:            true,
						},
					},
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Status of the job, for example `READY`").String,
				Computed:            true,
			},
			"phase": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Phase of the job, for example `ACTIVE`").String,
				Computed:            true,
			},
			"life_cycle_state": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Life cycle state of the job, for example `SUCCESS_LIFE_CYCLE_STATE`").String,
				Computed:            true,
			},
		},
	}
}

func (d *CollectionJobsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*CwProviderData).Clients
}

func (d *CollectionJobsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config CollectionJobsData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := clientFor(ctx, d.clients, config.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	application, contextID := config.ApplicationId.ValueString(), config.ContextId.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", collectionJobId(application, contextID)))

	row, found, err := readCollectionJob(&c, application, contextID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (POST), got error: %s", err))
		return
	}
	if !found {
		resp.Diagnostics.AddError("Collection job not found", fmt.Sprintf("There is no collection job with the application id %s and the context id %s on the cluster", application, contextID))
		return
	}
	job := CollectionJobs{Instance: config.Instance}
	resp.Diagnostics.Append(collectionJobFromRow(ctx, &c, row, &job)...)
	status, phase, state, err := collectionJobStatus(&c, application, contextID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the status of the job (POST), got error: %s", err))
	}
	if resp.Diagnostics.HasError() {
		return
	}
	out := CollectionJobsData{
		Id: job.Id, Instance: job.Instance, ApplicationId: job.ApplicationId, ContextId: job.ContextId, CollectorType: job.CollectorType,
		Devices: job.Devices, DeviceTag: job.DeviceTag, Destination: job.Destination, CliSensors: job.CliSensors, SnmpSensors: job.SnmpSensors,
		Status: types.StringValue(status), Phase: types.StringValue(phase), LifeCycleState: types.StringValue(state),
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", out.Id.ValueString()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &out)...)
}

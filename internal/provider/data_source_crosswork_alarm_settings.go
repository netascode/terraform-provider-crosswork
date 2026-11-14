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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &AlarmSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &AlarmSettingsDataSource{}
)

func NewAlarmSettingsDataSource() datasource.DataSource {
	return &AlarmSettingsDataSource{}
}

type AlarmSettingsDataSource struct {
	clients map[string]*cw.Client
}

func (d *AlarmSettingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alarm_settings"
}

func (d *AlarmSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the cleanup options and the trap and syslog processing of alarms and events.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"security_device_alarms_cleared_days": schema.Int64Attribute{
				MarkdownDescription: "Days after which cleared security device alarms are deleted",
				Computed:            true,
			},
			"non_security_device_alarms_cleared_days": schema.Int64Attribute{
				MarkdownDescription: "Days after which cleared non-security device alarms are deleted",
				Computed:            true,
			},
			"all_events_days": schema.Int64Attribute{
				MarkdownDescription: "Days after which all events are deleted",
				Computed:            true,
			},
			"device_alarms_days": schema.Int64Attribute{
				MarkdownDescription: "Days after which all device alarms, active and cleared, are deleted",
				Computed:            true,
			},
			"system_alarms_days": schema.Int64Attribute{
				MarkdownDescription: "Days after which all system alarms, active and cleared, are deleted",
				Computed:            true,
			},
			"network_alarms_days": schema.Int64Attribute{
				MarkdownDescription: "Days after which all network alarms, active and cleared, are deleted",
				Computed:            true,
			},
			"audit_logs_days": schema.Int64Attribute{
				MarkdownDescription: "Days after which audit logs are deleted",
				Computed:            true,
			},
			"trap_processing": schema.BoolAttribute{
				MarkdownDescription: "Process the traps of the devices",
				Computed:            true,
			},
			"syslog_processing": schema.BoolAttribute{
				MarkdownDescription: "Process the syslog messages of the devices",
				Computed:            true,
			},
		},
	}
}

func (d *AlarmSettingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *AlarmSettingsDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *AlarmSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config AlarmSettings

	// Read config
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := d.clientFor(ctx, config.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", config.Id.String()))

	params := ""
	res, err := c.Get("/crosswork/alarm/v1/settings" + params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}

	config.fromBody(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &DeviceNotificationPoliciesDataSource{}
	_ datasource.DataSourceWithConfigure = &DeviceNotificationPoliciesDataSource{}
)

func NewDeviceNotificationPoliciesDataSource() datasource.DataSource {
	return &DeviceNotificationPoliciesDataSource{}
}

type DeviceNotificationPoliciesDataSource struct {
	clients map[string]*cw.Client
}

func (d *DeviceNotificationPoliciesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_notification_policies"
}

func (d *DeviceNotificationPoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read a notification policy for device alarms and events.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the policy, unique on the Crosswork cluster (also against system policies). The controller accepts letters, digits, spaces and the characters `- . _ :`",
				Required:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the policy. The controller accepts letters, digits, spaces, line breaks and the characters `, - . _ :`",
				Computed:            true,
			},
			"device_group": schema.StringAttribute{
				MarkdownDescription: "Id of the device group whose devices the policy applies to. Without it the policy applies to all devices",
				Computed:            true,
			},
			"filters": schema.ListNestedAttribute{
				MarkdownDescription: "Filters of the alarms and events, an alarm or event is forwarded when it matches any filter. Each filter holds at least one of `severities`, `category` and `events`. A `category` needs `severities`",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"severities": schema.SetAttribute{
							MarkdownDescription: "Severities of the filter, each one of `information`, `warning`, `minor`, `major` and `critical`",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"category": schema.StringAttribute{
							MarkdownDescription: "Category of the filter, one name. It needs `severities` and cannot hold a single quote or a backslash",
							Computed:            true,
						},
						"events": schema.SetAttribute{
							MarkdownDescription: "Event types of the filter (the `eventTypeName` of the catalog). A value cannot hold a single quote or a backslash",
							ElementType:         types.StringType,
							Computed:            true,
						},
					},
				},
			},
			"trap_receivers": schema.SetNestedAttribute{
				MarkdownDescription: "Trap receivers that get the notifications. A policy needs at least one trap or syslog receiver",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the trap receiver, as in `crosswork_trap_destination`. It cannot hold a single quote or a backslash",
							Computed:            true,
						},
					},
				},
			},
			"syslog_receivers": schema.SetNestedAttribute{
				MarkdownDescription: "Syslog receivers that get the notifications. A policy needs at least one trap or syslog receiver",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the syslog receiver, as in `crosswork_syslog_destination`. It cannot hold a single quote or a backslash",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DeviceNotificationPoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *DeviceNotificationPoliciesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Written by hand: the read refuses a policy of another type, as the resource does.
func (d *DeviceNotificationPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config DeviceNotificationPolicies

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
	res, err := c.Get("/crosswork/alarm-forwarding/v1/alarm-forwarding/policy/Alarms" + params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}
	res = res.Get("groups.#(groupName==\"" + config.Id.ValueString() + "\")")
	if !res.Exists() {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No object found with id: %s", config.Id.ValueString()))
		return
	}

	if alarmType := res.Get("alarmType").String(); alarmType != "Device" {
		resp.Diagnostics.AddError("Not a Device policy", fmt.Sprintf("The policy %q exists but its alarm type is %q, not Device. Use crosswork_system_notification_policies for a System or Network policy.", config.Id.ValueString(), alarmType))
		return
	}

	config.fromBody(ctx, res)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

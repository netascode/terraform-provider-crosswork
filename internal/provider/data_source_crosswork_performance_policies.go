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

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &PerformancePoliciesDataSource{}
	_ datasource.DataSourceWithConfigure = &PerformancePoliciesDataSource{}
)

func NewPerformancePoliciesDataSource() datasource.DataSource {
	return &PerformancePoliciesDataSource{}
}

type PerformancePoliciesDataSource struct {
	clients map[string]*cw.Client
}

func (d *PerformancePoliciesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_performance_policies"
}

func (d *PerformancePoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read a performance monitoring policy.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Policy name, 1 to 255 characters",
				Optional:            true,
				Computed:            true,
			},
			"policy_template": schema.StringAttribute{
				MarkdownDescription: "Policy type, for example `deviceHealth` or `INTERFACE`",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Policy description",
				Computed:            true,
			},
			"contact": schema.StringAttribute{
				MarkdownDescription: "Contact written into the tag of the policy",
				Computed:            true,
			},
			"active": schema.BoolAttribute{
				MarkdownDescription: "Whether the policy collects data",
				Computed:            true,
			},
			"devices": schema.SetAttribute{
				MarkdownDescription: "Devices to poll, by their inventory id; set one of `devices`, `device_groups` and `port_groups`",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"device_groups": schema.SetAttribute{
				MarkdownDescription: "Ids of the device groups to poll",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"port_groups": schema.SetAttribute{
				MarkdownDescription: "Ids of the port groups to poll",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"schema_intervals": schema.SetNestedAttribute{
				MarkdownDescription: "Polling intervals by schema, one row for each schema of the policy type",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"schema": schema.StringAttribute{
							MarkdownDescription: "Schema of the policy type, for example `CPU`",
							Computed:            true,
						},
						"seconds": schema.Int64Attribute{
							MarkdownDescription: "Polling interval in seconds, 0 for no polling",
							Computed:            true,
						},
					},
				},
			},
			"thresholds": schema.MapNestedAttribute{
				MarkdownDescription: "Thresholds by metric, for example `cpuUtilization`",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"conditions": schema.ListNestedAttribute{
							MarkdownDescription: "Conditions of the metric, in order",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"relation": schema.StringAttribute{
										MarkdownDescription: "Comparison of the metric with the value",
										Computed:            true,
									},
									"value": schema.StringAttribute{
										MarkdownDescription: "Value of the threshold",
										Computed:            true,
									},
									"count": schema.Int64Attribute{
										MarkdownDescription: "Number of polls that must match before the alarm is raised",
										Computed:            true,
									},
									"severity": schema.StringAttribute{
										MarkdownDescription: "Severity of the alarm, `INFORMATION` raises an event",
										Computed:            true,
									},
									"match_cos": schema.SetAttribute{
										MarkdownDescription: "Classes of service of a Quality of service drop metric, 0 to 7",
										ElementType:         types.Int64Type,
										Computed:            true,
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
func (d *PerformancePoliciesDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(
			path.MatchRoot("id"),
			path.MatchRoot("name"),
		),
	}
}

func (d *PerformancePoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *PerformancePoliciesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Written by hand: the list holds the policies in "monitoringPolicy", and the devices are turned into inventory ids.
func (d *PerformancePoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config PerformancePolicies

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
	if config.Id.IsNull() && !config.Name.IsNull() {
		res, err := c.Get(config.getPath())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve objects, got error: %s", err))
			return
		}
		if value := res; len(value.Array()) > 0 {
			value.ForEach(func(k, v gjson.Result) bool {
				if config.Name.ValueString() == v.Get("monitoringPolicy.name").String() {
					config.Id = types.StringValue(v.Get("monitoringPolicy.id").String())
					tflog.Debug(ctx, fmt.Sprintf("%s: Found object with name '%v', id: %v", config.Id.String(), config.Name.ValueString(), config.Id.String()))
					return false
				}
				return true
			})
		}

		if config.Id.IsNull() {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to find object with name: %s", config.Name.ValueString()))
			return
		}
	}

	params := ""
	params += "/" + url.QueryEscape(config.Id.ValueString())
	res, err := c.Get(config.getPath() + params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}

	config.fromBody(ctx, res)
	if err := config.devicesToNames(ctx, &c); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to resolve the devices, got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

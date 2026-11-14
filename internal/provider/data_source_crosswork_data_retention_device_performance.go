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
	_ datasource.DataSource              = &DataRetentionDevicePerformanceDataSource{}
	_ datasource.DataSourceWithConfigure = &DataRetentionDevicePerformanceDataSource{}
)

func NewDataRetentionDevicePerformanceDataSource() datasource.DataSource {
	return &DataRetentionDevicePerformanceDataSource{}
}

type DataRetentionDevicePerformanceDataSource struct {
	clients map[string]*cw.Client
}

func (d *DataRetentionDevicePerformanceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_retention_device_performance"
}

func (d *DataRetentionDevicePerformanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the data retention of the device performance data.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"rows": schema.MapNestedAttribute{
				MarkdownDescription: "Retention periods by policy type, for example `CEPM_INTERFACE`, in hours. A policy type that is not listed keeps its periods",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"raw_hours": schema.Int64Attribute{
							MarkdownDescription: "Hours that raw data is kept",
							Computed:            true,
						},
						"hourly_hours": schema.Int64Attribute{
							MarkdownDescription: "Hours that hourly data is kept, 72 to 720, or 0 for a policy type without aggregation",
							Computed:            true,
						},
						"daily_hours": schema.Int64Attribute{
							MarkdownDescription: "Hours that daily data is kept, 336 to 8760, or 0 for a policy type without aggregation",
							Computed:            true,
						},
						"weekly_hours": schema.Int64Attribute{
							MarkdownDescription: "Hours that weekly data is kept, 336 to 17472, or 0 for a policy type without aggregation",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DataRetentionDevicePerformanceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *DataRetentionDevicePerformanceDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *DataRetentionDevicePerformanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config DataRetentionDevicePerformance

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
	res, err := c.Get("/crosswork/performance/v1/dataretention/all" + params)
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

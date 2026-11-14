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

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

var _ datasource.DataSource = &GoldenConfigJobsDataSource{}
var _ datasource.DataSourceWithConfigure = &GoldenConfigJobsDataSource{}

func NewGoldenConfigJobsDataSource() datasource.DataSource {
	return &GoldenConfigJobsDataSource{}
}

type GoldenConfigJobsDataSource struct {
	clients map[string]*cw.Client
}

type goldenConfigApplicationData struct {
	DeviceName      types.String `tfsdk:"device_name"`
	DeviceUuid      types.String `tfsdk:"device_uuid"`
	ApplicationName types.String `tfsdk:"application_name"`
	Conformant      types.Bool   `tfsdk:"conformant"`
}

var goldenConfigApplicationType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"device_name": types.StringType, "device_uuid": types.StringType, "application_name": types.StringType, "conformant": types.BoolType,
}}

func sortApplications(rows []goldenConfigApplicationData) {
	sort.Slice(rows, func(i, j int) bool { return rows[i].DeviceName.ValueString() < rows[j].DeviceName.ValueString() })
}

type goldenConfigJobsData struct {
	Id         types.String `tfsdk:"id"`
	Instance   types.String `tfsdk:"instance"`
	Device     types.String `tfsdk:"device"`
	TemplateId types.Int64  `tfsdk:"template_id"`
	Devices    types.List   `tfsdk:"devices"`
}

func (d *GoldenConfigJobsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_golden_config_jobs"
}

func (d *GoldenConfigJobsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This data source can read the devices a Golden Configuration template is applied to, with their conformance.").String,

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object, the template name",
				Required:            true,
			},
			"device": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Only this device, by inventory id").String,
				Optional:            true,
			},
			"template_id": schema.Int64Attribute{MarkdownDescription: "The integer id the controller gave the template", Computed: true},
			"devices": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("The devices the template is applied to, sorted by name; conformance is the controller's own check, and turns true about 45 seconds after a job").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"device_name":      schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("Inventory id of the device").String, Computed: true},
						"device_uuid":      schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("Uuid of the device").String, Computed: true},
						"application_name": schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("Name of the application, the template on the device").String, Computed: true},
						"conformant":       schema.BoolAttribute{MarkdownDescription: helpers.NewAttributeDescription("Whether the device matches the template").String, Computed: true},
					},
				},
			},
		},
	}
}

func (d *GoldenConfigJobsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*CwProviderData).Clients
}

func (d *GoldenConfigJobsDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

func (d *GoldenConfigJobsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config goldenConfigJobsData
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := d.clientFor(ctx, config.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	name := config.Id.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", name))

	id, found, err := gcFindTemplate(&c, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (POST), got error: %s", err))
		return
	}
	if !found {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("The template %q does not exist", name))
		return
	}
	apps, err := gcTemplateApplications(&c, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the applications (POST), got error: %s", err))
		return
	}
	rows := []goldenConfigApplicationData{}
	for _, a := range apps {
		if !config.Device.IsNull() && a.DeviceName != config.Device.ValueString() {
			continue
		}
		rows = append(rows, goldenConfigApplicationData{DeviceName: types.StringValue(a.DeviceName), DeviceUuid: types.StringValue(a.DeviceUUID),
			ApplicationName: types.StringValue(a.ApplicationName), Conformant: types.BoolValue(a.Conformant)})
	}
	sortApplications(rows)
	list, d2 := types.ListValueFrom(ctx, goldenConfigApplicationType, rows)
	resp.Diagnostics.Append(d2...)
	config.TemplateId = types.Int64Value(id)
	config.Devices = list
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", name))
}

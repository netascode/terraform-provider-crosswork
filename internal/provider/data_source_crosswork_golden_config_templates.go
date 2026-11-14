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

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

var _ datasource.DataSource = &GoldenConfigTemplatesDataSource{}
var _ datasource.DataSourceWithConfigure = &GoldenConfigTemplatesDataSource{}

func NewGoldenConfigTemplatesDataSource() datasource.DataSource {
	return &GoldenConfigTemplatesDataSource{}
}

type GoldenConfigTemplatesDataSource struct {
	clients map[string]*cw.Client
}

type goldenConfigTemplateVariableData struct {
	Name  types.String `tfsdk:"name"`
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

var goldenConfigTemplateVariableDataType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"name": types.StringType, "type": types.StringType, "value": types.StringType,
}}

type goldenConfigTemplatesData struct {
	Id            types.String `tfsdk:"id"`
	Instance      types.String `tfsdk:"instance"`
	Name          types.String `tfsdk:"name"`
	TemplateId    types.Int64  `tfsdk:"template_id"`
	ConfigId      types.Int64  `tfsdk:"config_id"`
	Description   types.String `tfsdk:"description"`
	Tags          types.Set    `tfsdk:"tags"`
	Mode          types.String `tfsdk:"mode"`
	ConfigType    types.String `tfsdk:"config_type"`
	Vendor        types.String `tfsdk:"vendor"`
	SoftwareType  types.String `tfsdk:"software_type"`
	ProductSeries types.Set    `tfsdk:"product_series"`
	Audit         types.Bool   `tfsdk:"audit"`
	Config        types.String `tfsdk:"config"`
	Variables     types.List   `tfsdk:"variables"`
}

func (d *GoldenConfigTemplatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_golden_config_templates"
}

func (d *GoldenConfigTemplatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This data source can read a template of Golden Configuration.").String,

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object, the template name",
				Required:            true,
			},
			"name":           schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("Template name").String, Computed: true},
			"template_id":    schema.Int64Attribute{MarkdownDescription: "The integer id the controller gave the template", Computed: true},
			"config_id":      schema.Int64Attribute{MarkdownDescription: "The integer id of the configuration of the template, which a job refers to", Computed: true},
			"description":    schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("Description").String, Computed: true},
			"tags":           schema.SetAttribute{MarkdownDescription: helpers.NewAttributeDescription("Tags of the template").String, ElementType: types.StringType, Computed: true},
			"mode":           schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("How the configuration is applied").String, Computed: true},
			"config_type":    schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("Format of the configuration text").String, Computed: true},
			"vendor":         schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("Vendor of the devices").String, Computed: true},
			"software_type":  schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("Software type of the devices").String, Computed: true},
			"product_series": schema.SetAttribute{MarkdownDescription: helpers.NewAttributeDescription("Product series of the devices").String, ElementType: types.StringType, Computed: true},
			"audit":          schema.BoolAttribute{MarkdownDescription: helpers.NewAttributeDescription("Whether the template is part of the global conformance schedule").String, Computed: true},
			"config":         schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("Configuration text with Jinja2 variables").String, Computed: true},
			"variables": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Variables of the template, sorted by name; the controller returns the values in clear text").String,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name":  schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("Variable name").String, Computed: true},
						"type":  schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("Type of the value").String, Computed: true},
						"value": schema.StringAttribute{MarkdownDescription: helpers.NewAttributeDescription("Value").String, Computed: true, Sensitive: true},
					},
				},
			},
		},
	}
}

func (d *GoldenConfigTemplatesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*CwProviderData).Clients
}

func (d *GoldenConfigTemplatesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

func (d *GoldenConfigTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config goldenConfigTemplatesData
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
	res, err := gcGetTemplate(&c, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s", err))
		return
	}
	m := GoldenConfigTemplates{Variables: types.ListNull(goldenConfigVariableType), Tags: types.SetNull(types.StringType), Description: types.StringNull()}
	resp.Diagnostics.Append(m.fromAPI(ctx, res)...)
	vars, _ := gcVarsOfValues(ctx, m.Variables)
	config.Name, config.TemplateId, config.Description = m.Name, types.Int64Value(id), types.StringValue(res.Get("description").String())
	config.ConfigId = m.ConfigId
	config.Mode, config.ConfigType, config.Vendor, config.SoftwareType = m.Mode, m.ConfigType, m.Vendor, m.SoftwareType
	config.ProductSeries, config.Audit, config.Config = m.ProductSeries, m.Audit, m.Config
	tags, d2 := stringSet(ctx, gcStrings(res.Get("tags")))
	resp.Diagnostics.Append(d2...)
	config.Tags = tags
	list, d3 := types.ListValueFrom(ctx, goldenConfigTemplateVariableDataType, vars)
	resp.Diagnostics.Append(d3...)
	config.Variables = list
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", name))
}

func gcVarsOfValues(ctx context.Context, list types.List) ([]goldenConfigTemplateVariableData, diag.Diagnostics) {
	var items []GoldenConfigVariable
	var diags diag.Diagnostics
	if list.IsNull() || list.IsUnknown() {
		return []goldenConfigTemplateVariableData{}, diags
	}
	diags.Append(list.ElementsAs(ctx, &items, false)...)
	out := make([]goldenConfigTemplateVariableData, 0, len(items))
	for _, it := range items {
		out = append(out, goldenConfigTemplateVariableData{Name: it.Name, Type: it.Type, Value: it.Value})
	}
	return out, diags
}

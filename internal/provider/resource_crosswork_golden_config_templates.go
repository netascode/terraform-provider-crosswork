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
	"regexp"
	"strconv"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

var _ resource.Resource = &GoldenConfigTemplatesResource{}
var _ resource.ResourceWithImportState = &GoldenConfigTemplatesResource{}
var _ resource.ResourceWithValidateConfig = &GoldenConfigTemplatesResource{}

func NewGoldenConfigTemplatesResource() resource.Resource {
	return &GoldenConfigTemplatesResource{}
}

type GoldenConfigTemplatesResource struct {
	clients map[string]*cw.Client
}

var gcTemplateName = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_-]{0,63}$`)

func (r *GoldenConfigTemplatesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_golden_config_templates"
}

func (r *GoldenConfigTemplatesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a template of Golden Configuration.").String,

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
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Template name, 1 to 64 letters, digits, underscores and hyphens, starting with a letter").String,
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators:          []validator.String{stringvalidator.RegexMatches(gcTemplateName, "must be 1 to 64 letters, digits, underscores and hyphens and start with a letter")},
			},
			"template_id": schema.Int64Attribute{
				MarkdownDescription: "The integer id the controller gave the template",
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"config_id": schema.Int64Attribute{
				MarkdownDescription: "The integer id of the configuration of the template, which a job refers to",
				Computed:            true,
				PlanModifiers:       []planmodifier.Int64{int64planmodifier.UseStateForUnknown()},
			},
			"description": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Description, up to 255 characters").String,
				Optional:            true,
				Validators:          []validator.String{stringvalidator.LengthAtMost(255)},
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Tags of the template, up to 64 characters each").String,
				ElementType:         types.StringType,
				Optional:            true,
				Validators:          []validator.Set{setvalidator.ValueStringsAre(stringvalidator.LengthAtMost(64))},
			},
			"mode": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("How the configuration is applied: merged into the device or replacing what it covers").AddStringEnumDescription("merge", "replace").AddDefaultValueDescription("merge").String,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("merge"),
				Validators:          []validator.String{stringvalidator.OneOf("merge", "replace")},
			},
			"config_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Format of the configuration text").AddStringEnumDescription("native", "xml", "json", "c-style", "j-style").AddDefaultValueDescription("native").String,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("native"),
				Validators:          []validator.String{stringvalidator.OneOf("native", "xml", "json", "c-style", "j-style")},
			},
			"vendor": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Vendor of the devices, as the inventory spells it; the controller refuses a vendor the inventory does not know").String,
				Required:            true,
			},
			"software_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Software type of the devices, for example `IOS XR`").String,
				Required:            true,
			},
			"product_series": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Product series of the devices; the controller does not check them against a device").String,
				ElementType:         types.StringType,
				Required:            true,
				Validators:          []validator.Set{setvalidator.SizeAtLeast(1)},
			},
			"audit": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Include the template in the global conformance schedule").AddDefaultValueDescription("true").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(true),
				PlanModifiers:       []planmodifier.Bool{boolplanmodifier.RequiresReplace()},
			},
			"config": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Configuration text with Jinja2 variables").String,
				Required:            true,
			},
			"variables": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Variables of the template. Preferred over global variables.").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Variable name, 1 to 64 letters, digits and underscores, not starting with a digit").String,
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
							MarkdownDescription: helpers.NewAttributeDescription("Value; for the type `JSON` a JSON document, which the controller refuses to be empty").AddMutualExclusivityDescription("**Required**: exactly one of `value` and `value_wo` must be set.").AddCoexistenceNote("This attribute stores the value in Terraform state. Use `value_wo` together with `value_wo_version` for a secret.").String,
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
			"force_delete": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Delete the template with force, which also removes it from devices it is applied to. Without it the controller refuses to delete an applied template.").AddDefaultValueDescription("false").String,
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
		},
	}
}

func (r *GoldenConfigTemplatesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clients = req.ProviderData.(*CwProviderData).Clients
}

func (r *GoldenConfigTemplatesResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

func (r *GoldenConfigTemplatesResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var cfg GoldenConfigTemplates
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(validateGcVariables(ctx, cfg.Variables, "variables")...)
}

func (m GoldenConfigTemplates) content(ctx context.Context, cfg GoldenConfigTemplates) (gcTemplate, diag.Diagnostics) {
	var diags diag.Diagnostics
	tags, d := stringsOf(ctx, m.Tags)
	diags.Append(d...)
	series, d := stringsOf(ctx, m.ProductSeries)
	diags.Append(d...)
	vars, d := gcVarsOf(ctx, m.Variables, cfg.Variables)
	diags.Append(d...)
	return gcTemplate{Name: m.Name.ValueString(), Description: m.Description.ValueString(), Tags: tags,
		Mode: m.Mode.ValueString(), ConfigType: m.ConfigType.ValueString(), Vendor: m.Vendor.ValueString(),
		SoftwareType: m.SoftwareType.ValueString(), ProductSeries: series, Audit: m.Audit.ValueBool(),
		Config: m.Config.ValueString(), Variables: vars}, diags
}

func (r *GoldenConfigTemplatesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, cfg GoldenConfigTemplates
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Name.ValueString()))

	t, d := plan.content(ctx, cfg)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	res, err := c.Post(gcPath+"/templates", t.body(), cw.NoLogPayload)
	if err != nil {
		if cw.IsConflict(err) {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("A template named %q exists already (POST): import it with terraform import, got error: %s", t.Name, err))
			return
		}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (POST), got error: %s", err))
		return
	}
	id, perr := strconv.ParseInt(res.Get("key").String(), 10, 64)
	if perr != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("The controller answered without the id of the template: %s", res.Raw))
		return
	}
	plan.Id = plan.Name
	plan.TemplateId = types.Int64Value(id)
	// The id a job needs is the one of the configuration, which only a read of the template gives.
	created, err := gcGetTemplate(&c, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("The template was created, but reading it failed (GET), got error: %s", err))
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}
	plan.ConfigId = types.Int64Value(created.Get("versions.0.tcId").Int())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Name.ValueString()))
}

func (r *GoldenConfigTemplatesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state GoldenConfigTemplates
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

	id, found, err := gcFindTemplate(&c, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (POST), got error: %s", err))
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	res, err := gcGetTemplate(&c, id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s", err))
		return
	}
	state.TemplateId = types.Int64Value(id)
	if state.ForceDelete.IsNull() {
		state.ForceDelete = types.BoolValue(false)
	}
	resp.Diagnostics.Append(state.fromAPI(ctx, res)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", name))
}

func (r *GoldenConfigTemplatesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state, cfg GoldenConfigTemplates
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

	t, d := plan.content(ctx, cfg)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	// The update is in place: the template keeps its one version and its id.
	if _, err := c.Put(fmt.Sprintf("%s/templates/%d", gcPath, state.TemplateId.ValueInt64()), t.body(), cw.NoLogPayload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s", err))
		return
	}
	plan.TemplateId = state.TemplateId
	plan.ConfigId = state.ConfigId
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))
}

func (r *GoldenConfigTemplatesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state GoldenConfigTemplates
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))

	// Without force the controller refuses a template that is applied (409); with force it also removes the
	// configuration from the devices.
	_, err := c.Delete(gcPath+"/templates", cw.Query("ids", strconv.FormatInt(state.TemplateId.ValueInt64(), 10)),
		cw.Query("force", strconv.FormatBool(state.ForceDelete.ValueBool())), cw.NoLogPayload)
	if err != nil && !cw.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))
}

func (r *GoldenConfigTemplatesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <name>. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
}

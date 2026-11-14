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
	"strings"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

var _ resource.Resource = &GoldenConfigGlobalVariablesResource{}
var _ resource.ResourceWithImportState = &GoldenConfigGlobalVariablesResource{}
var _ resource.ResourceWithValidateConfig = &GoldenConfigGlobalVariablesResource{}

func NewGoldenConfigGlobalVariablesResource() resource.Resource {
	return &GoldenConfigGlobalVariablesResource{}
}

type GoldenConfigGlobalVariablesResource struct {
	clients map[string]*cw.Client
}

var gcVariableName = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]{0,63}$`)

func (r *GoldenConfigGlobalVariablesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_golden_config_global_variables"
}

func (r *GoldenConfigGlobalVariablesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a global variable of Golden Configuration.").String,

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
				MarkdownDescription: helpers.NewAttributeDescription("Variable name, 1 to 64 letters, digits and underscores, not starting with a digit").String,
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators:          []validator.String{stringvalidator.RegexMatches(gcVariableName, "must be 1 to 64 letters, digits and underscores and must not start with a digit")},
			},
			"variable_id": schema.StringAttribute{
				MarkdownDescription: "The id the controller gave the variable",
				Computed:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Type of the value").AddStringEnumDescription("STRING", "JSON").AddDefaultValueDescription("STRING").String,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("STRING"),
				Validators:          []validator.String{stringvalidator.OneOf("STRING", "JSON")},
			},
			"value": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Value of the variable; for the type `JSON` a JSON document").AddMutualExclusivityDescription("**Required**: exactly one of `value` and `value_wo` must be set.").AddCoexistenceNote("This attribute stores the value in Terraform state. Use `value_wo` together with `value_wo_version` for a secret: the controller returns the value in clear text on a read, but it stays out of the plan and the state.").String,
				Optional:            true,
			},
			"value_wo": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Value of the variable").AddMutualExclusivityDescription("**Required**: exactly one of `value` and `value_wo` must be set.").String,
				Optional:            true,
				WriteOnly:           true,
				Sensitive:           true,
			},
			"value_wo_version": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `value_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new value. The version is stored in state; the value is not.").String,
				Optional:            true,
			},
		},
	}
}

func (r *GoldenConfigGlobalVariablesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clients = req.ProviderData.(*CwProviderData).Clients
}

func (r *GoldenConfigGlobalVariablesResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

func (r *GoldenConfigGlobalVariablesResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var cfg GoldenConfigGlobalVariables
	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if cfg.Value.IsUnknown() || cfg.ValueWo.IsUnknown() || cfg.Type.IsUnknown() {
		return
	}
	switch {
	case !cfg.Value.IsNull() && !cfg.ValueWo.IsNull():
		resp.Diagnostics.AddError("Invalid Attribute Combination", "Only one of `value` and `value_wo` can be set.")
	case cfg.Value.IsNull() && cfg.ValueWo.IsNull():
		resp.Diagnostics.AddError("Invalid Attribute Combination", "Exactly one of `value` and `value_wo` must be set.")
	}
	if !cfg.ValueWo.IsNull() && cfg.ValueWoVersion.IsNull() {
		resp.Diagnostics.AddError("Missing Attribute", "`value_wo_version` must be set when `value_wo` is used. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.")
	}
	if cfg.Type.ValueString() == "JSON" {
		for _, v := range []types.String{cfg.Value, cfg.ValueWo} {
			if !v.IsNull() && !gcValidJSON(v.ValueString()) {
				resp.Diagnostics.AddAttributeError(path.Root("value"), "Invalid JSON", "A variable of the type `JSON` needs a JSON document as its value.")
			}
		}
	}
}

func (r *GoldenConfigGlobalVariablesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, cfg GoldenConfigGlobalVariables
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

	res, err := c.Post(gcPath+"/gvars", plan.body(cfg), cw.NoLogPayload)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (POST), got error: %s", err))
		return
	}
	plan.Id = plan.Name
	plan.VariableId = types.StringValue(res.Get("0.key").String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Name.ValueString()))
}

func (r *GoldenConfigGlobalVariablesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state GoldenConfigGlobalVariables
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

	row, found, err := gcFindGlobalVariable(&c, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (POST), got error: %s", err))
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	state.fromRow(row)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", name))
}

func (r *GoldenConfigGlobalVariablesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state, cfg GoldenConfigGlobalVariables
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

	// The update is a PATCH by name; the write-only value is in the configuration on every update.
	if _, err := c.Patch(gcPath+"/gvars", plan.body(cfg), cw.NoLogPayload); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PATCH), got error: %s", err))
		return
	}
	plan.VariableId = state.VariableId
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))
}

func (r *GoldenConfigGlobalVariablesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state GoldenConfigGlobalVariables
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))

	// The controller does not refuse to delete a variable a template uses, so the destroy has to be ordered with depends_on.
	_, err := c.Delete(gcPath+"/gvars", cw.Query("ids", state.VariableId.ValueString()), cw.NoLogPayload)
	if err != nil && !cw.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))
}

func (r *GoldenConfigGlobalVariablesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

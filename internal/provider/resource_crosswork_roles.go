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
	"net/url"
	"sort"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

var _ resource.Resource = &RolesResource{}
var _ resource.ResourceWithImportState = &RolesResource{}

func NewRolesResource() resource.Resource {
	return &RolesResource{}
}

type RolesResource struct {
	clients map[string]*cw.Client
}

func (r *RolesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_roles"
}

func (r *RolesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a role in the local authentication store.").String,

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
				MarkdownDescription: helpers.NewAttributeDescription("Role name").String,
				Required:            true,
				PlanModifiers:       []planmodifier.String{stringplanmodifier.RequiresReplace()},
				Validators:          []validator.String{stringvalidator.LengthAtLeast(1)},
			},
			"api_permissions": schema.MapAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("API id with its access level: `R` read, `W` write, `D` delete. The ids are listed by the `crosswork_role_catalog` data source").AddStringEnumDescription("R", "W", "D", "RW", "RD", "WD", "RWD").String,
				ElementType:         types.StringType,
				Required:            true,
				Validators: []validator.Map{
					mapvalidator.SizeAtLeast(1),
					mapvalidator.ValueStringsAre(stringvalidator.OneOf("R", "W", "D", "RW", "RD", "WD", "RWD")),
				},
			},
			"task_permissions": schema.SetAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Task ids of the role. The ids are listed by the `crosswork_role_catalog` data source").String,
				ElementType:         types.StringType,
				Optional:            true,
			},
		},
	}
}

func (r *RolesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.clients = req.ProviderData.(*CwProviderData).Clients
}

func (r *RolesResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

func roleAPINames(c *cw.Client) (map[string]string, error) {
	res, err := c.Get(apiPath)
	if err != nil {
		return nil, err
	}
	names := map[string]string{}
	res.ForEach(func(_, v gjson.Result) bool {
		names[v.Get("api_id").String()] = v.Get("name").String()
		return true
	})
	return names, nil
}

// Crosswork has no call for a single role.
func getRole(c *cw.Client, name string) (gjson.Result, bool, error) {
	res, err := c.Get(rolePath)
	if err != nil {
		return gjson.Result{}, false, err
	}
	role, found := res.Map()[name]
	return role, found, nil
}

func getTasks(c *cw.Client, name string) ([]string, error) {
	res, err := c.Get(taskPath + "/" + url.PathEscape(name))
	if err != nil {
		return nil, err
	}
	enabled := []string{}
	res.ForEach(func(_, group gjson.Result) bool {
		group.Get("items").ForEach(func(_, item gjson.Result) bool {
			if item.Get("enabled").Bool() {
				enabled = append(enabled, item.Get("id").String())
			}
			return true
		})
		return true
	})
	sort.Strings(enabled)
	return enabled, nil
}

// Crosswork rejects a task whose APIs the role does not hold, and accepts an empty list.
func setTasks(c *cw.Client, name string, ids []string) error {
	_, err := c.Post(taskPath+"/"+url.PathEscape(name), tasksBody(ids))
	return err
}

func putRole(c *cw.Client, name string, perms, names map[string]string) error {
	role, err := roleBody(perms, names)
	if err != nil {
		return err
	}
	body, err := flatBody(role)
	if err != nil {
		return err
	}
	_, err = c.Put(rolePath+"/"+url.PathEscape(name), body)
	return err
}

func (r *RolesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan Roles
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	name := plan.Name.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", name))

	perms, d := plan.permissions(ctx)
	resp.Diagnostics.Append(d...)
	tasks, d := plan.taskIds(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	names, err := roleAPINames(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to read the API list (GET), got error: %s", err))
		return
	}
	role, err := roleBody(perms, names)
	if err != nil {
		resp.Diagnostics.AddError("Invalid api_permissions", err.Error())
		return
	}
	body, err := createBody(name, role)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to render the request body, got error: %s", err))
		return
	}
	if _, err := c.Post(rolePath, body); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (POST), got error: %s", err))
		return
	}
	plan.Id = types.StringValue(name)
	// The role exists from here on, so it is saved before the tasks are set: if that fails the role stays tracked.
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if len(tasks) > 0 {
		if err := setTasks(&c, name, tasks); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("The role was created, but setting its tasks failed (POST), got error: %s", err))
			return
		}
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", name))
}

func (r *RolesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state Roles
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

	role, found, err := getRole(&c, name)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s", err))
		return
	}
	if !found {
		resp.State.RemoveResource(ctx)
		return
	}
	var enabled []string
	if state.ApiPermissions.IsNull() || !state.TaskPermissions.IsNull() {
		if enabled, err = getTasks(&c, name); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the tasks (GET), got error: %s", err))
			return
		}
	}
	state.Name = types.StringValue(name)
	resp.Diagnostics.Append(state.fromRole(ctx, role, enabled, true)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", name))
}

// Update keeps the order Crosswork needs: a task is rejected unless the role already holds the APIs it names,
// so an API is granted before the tasks that need it are set. The reverse is not enforced (Crosswork accepts
// taking an API away while a task still needs it and leaves the role inconsistent), so the tasks are removed
// before the API is revoked to avoid that. A change that does both goes through the union of the old and
// the new permissions.
func (r *RolesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state Roles
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	name := plan.Id.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", name))

	newPerms, d := plan.permissions(ctx)
	resp.Diagnostics.Append(d...)
	oldPerms, d := state.permissions(ctx)
	resp.Diagnostics.Append(d...)
	newTasks, d := plan.taskIds(ctx)
	resp.Diagnostics.Append(d...)
	oldTasks, d := state.taskIds(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	permsChanged := !equalPermissions(oldPerms, newPerms)
	tasksChanged := !sameStrings(oldTasks, newTasks)
	union := unionLevels(oldPerms, newPerms)
	var names map[string]string
	if permsChanged {
		var err error
		if names, err = roleAPINames(&c); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to read the API list (GET), got error: %s", err))
			return
		}
		if _, err := roleBody(newPerms, names); err != nil {
			resp.Diagnostics.AddError("Invalid api_permissions", err.Error())
			return
		}
		if !equalPermissions(union, oldPerms) {
			if err := putRole(&c, name, union, names); err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s", err))
				return
			}
		}
	}
	if tasksChanged {
		if err := setTasks(&c, name, newTasks); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to set the tasks (POST), got error: %s", err))
			return
		}
	}
	if permsChanged && !equalPermissions(union, newPerms) {
		if err := putRole(&c, name, newPerms, names); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s", err))
			return
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", name))
}

func (r *RolesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state Roles
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	name := state.Id.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", name))

	// Crosswork reports a role that is already gone as HTTP 500 "Role doesn't exist".
	_, err := c.Delete(rolePath + "/" + url.PathEscape(name))
	if err != nil && !cw.IsNotFound(err) && !strings.Contains(err.Error(), "doesn't exist") {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s", err))
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", name))
}

func (r *RolesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

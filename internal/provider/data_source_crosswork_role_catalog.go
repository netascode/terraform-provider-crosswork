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

// Written by hand: the catalog is read from three calls and has no object of its own.

import (
	"context"
	"fmt"
	"net/url"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

const (
	apiCategoriesPath = "/crosswork/aaa/v2/api"
	taskAPIPath       = "/crosswork/aaa/v1/taskAPIPermission"
	// catalogRole is the role the task calls are made for. Its tasks are all enabled, so the task list and
	// the task-to-API map are complete; the other roles list only what they have enabled.
	catalogRole = "admin"
)

var _ datasource.DataSource = &RoleCatalogDataSource{}
var _ datasource.DataSourceWithConfigure = &RoleCatalogDataSource{}

func NewRoleCatalogDataSource() datasource.DataSource {
	return &RoleCatalogDataSource{}
}

type RoleCatalogDataSource struct {
	clients map[string]*cw.Client
}

type catalogEntry struct {
	Id       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Category types.String `tfsdk:"category"`
}

type catalogRequirement struct {
	TaskId types.String `tfsdk:"task_id"`
	ApiId  types.String `tfsdk:"api_id"`
	Access types.String `tfsdk:"access"`
}

func (d *RoleCatalogDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role_catalog"
}

func (d *RoleCatalogDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	entry := func(idNote, nameNote, categoryNote string) map[string]schema.Attribute {
		return map[string]schema.Attribute{
			"id":       schema.StringAttribute{MarkdownDescription: idNote, Computed: true},
			"name":     schema.StringAttribute{MarkdownDescription: nameNote, Computed: true},
			"category": schema.StringAttribute{MarkdownDescription: categoryNote, Computed: true},
		}
	}
	resp.Schema = schema.Schema{
		MarkdownDescription: "This data source fetches all APIs and tasks that a role can be granted on the Crosswork cluster.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"apis": schema.SetNestedAttribute{
				MarkdownDescription: "List of APIs on the Crosswork cluster.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: entry("API id", "API name", "API category"),
				},
			},
			"tasks": schema.SetNestedAttribute{
				MarkdownDescription: "List of tasks on the Crosswork cluster.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: entry("Task id", "Task name", "Task category"),
				},
			},
			"task_api_permissions": schema.SetNestedAttribute{
				MarkdownDescription: "List of APIs that tasks need.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"task_id": schema.StringAttribute{MarkdownDescription: "Task id", Computed: true},
						"api_id":  schema.StringAttribute{MarkdownDescription: "API id", Computed: true},
						"access":  schema.StringAttribute{MarkdownDescription: "Access level, such as `R` or `RW`", Computed: true},
					},
				},
			},
		},
	}
}

func (d *RoleCatalogDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*CwProviderData).Clients
}

func (d *RoleCatalogDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// An object of category name to a list of {api_id, name}.
func catalogAPIs(c *cw.Client) ([]catalogEntry, error) {
	res, err := c.Get(apiCategoriesPath)
	if err != nil {
		return nil, err
	}
	apis := []catalogEntry{}
	res.ForEach(func(category, items gjson.Result) bool {
		items.ForEach(func(_, item gjson.Result) bool {
			apis = append(apis, catalogEntry{
				Id:       types.StringValue(item.Get("api_id").String()),
				Name:     types.StringValue(item.Get("name").String()),
				Category: types.StringValue(category.String()),
			})
			return true
		})
		return true
	})
	sort.Slice(apis, func(i, j int) bool { return apis[i].Id.ValueString() < apis[j].Id.ValueString() })
	return apis, nil
}

// A list of task groups, each with its items.
func catalogTasks(c *cw.Client) ([]catalogEntry, error) {
	res, err := c.Get(taskPath + "/" + url.PathEscape(catalogRole))
	if err != nil {
		return nil, err
	}
	tasks := []catalogEntry{}
	res.ForEach(func(_, group gjson.Result) bool {
		group.Get("items").ForEach(func(_, item gjson.Result) bool {
			tasks = append(tasks, catalogEntry{
				Id:       types.StringValue(item.Get("id").String()),
				Name:     types.StringValue(item.Get("name").String()),
				Category: types.StringValue(group.Get("name").String()),
			})
			return true
		})
		return true
	})
	sort.Slice(tasks, func(i, j int) bool { return tasks[i].Id.ValueString() < tasks[j].Id.ValueString() })
	return tasks, nil
}

// The letters of the result come in the order R, W, D.
func levelFromLetters(letters gjson.Result) string {
	have := map[string]bool{}
	letters.ForEach(func(_, l gjson.Result) bool {
		have[l.String()] = true
		return true
	})
	level := ""
	for _, l := range []string{"R", "W", "D"} {
		if have[l] {
			level += l
		}
	}
	return level
}

func catalogTaskAPIs(c *cw.Client) ([]catalogRequirement, error) {
	res, err := c.Get(taskAPIPath + "/" + url.PathEscape(catalogRole))
	if err != nil {
		return nil, err
	}
	needs := []catalogRequirement{}
	res.ForEach(func(task, entry gjson.Result) bool {
		entry.Get("apiIds").ForEach(func(api, letters gjson.Result) bool {
			needs = append(needs, catalogRequirement{
				TaskId: types.StringValue(task.String()),
				ApiId:  types.StringValue(api.String()),
				Access: types.StringValue(levelFromLetters(letters)),
			})
			return true
		})
		return true
	})
	sort.Slice(needs, func(i, j int) bool {
		if needs[i].TaskId.ValueString() != needs[j].TaskId.ValueString() {
			return needs[i].TaskId.ValueString() < needs[j].TaskId.ValueString()
		}
		return needs[i].ApiId.ValueString() < needs[j].ApiId.ValueString()
	})
	return needs, nil
}

func (d *RoleCatalogDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config RoleCatalog
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := d.clientFor(ctx, config.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	tflog.Debug(ctx, "Role catalog: Beginning Read")

	apis, err := catalogAPIs(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the APIs, got error: %s", err))
		return
	}
	tasks, err := catalogTasks(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the tasks, got error: %s", err))
		return
	}
	needs, err := catalogTaskAPIs(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the task API permissions, got error: %s", err))
		return
	}

	config.Apis = make([]RoleCatalogApis, len(apis))
	for i, e := range apis {
		config.Apis[i] = RoleCatalogApis(e)
	}
	config.Tasks = make([]RoleCatalogTasks, len(tasks))
	for i, e := range tasks {
		config.Tasks[i] = RoleCatalogTasks(e)
	}
	config.TaskApiPermissions = make([]RoleCatalogTaskApiPermissions, len(needs))
	for i, e := range needs {
		config.TaskApiPermissions[i] = RoleCatalogTaskApiPermissions(e)
	}
	config.Id = types.StringValue("role_catalog")
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	tflog.Debug(ctx, "Role catalog: Read finished successfully")
}

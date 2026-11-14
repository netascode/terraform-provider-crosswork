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

// Written by hand: the groups of a classifier form a tree, which the data source flattens with the path of names down to each group.

import (
	"context"
	"fmt"
	"net/url"
	"sort"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

const groupingPath = "/crosswork/grouping/v1/grouping"

var _ datasource.DataSource = &GroupingsDataSource{}
var _ datasource.DataSourceWithConfigure = &GroupingsDataSource{}

func NewGroupingsDataSource() datasource.DataSource {
	return &GroupingsDataSource{}
}

type GroupingsDataSource struct {
	clients map[string]*cw.Client
}

func (d *GroupingsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groupings"
}

func (d *GroupingsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This data source fetches all groups of a classifier, such as the device access groups",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"classifier": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Group classifier").AddStringEnumDescription("DeviceAccess", "LocationDevices", "UserDefinedDevices", "TopologyTypeDevices", "DeviceType", "UserDefinedPorts", "PortType").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("DeviceAccess", "LocationDevices", "UserDefinedDevices", "TopologyTypeDevices", "DeviceType", "UserDefinedPorts", "PortType"),
				},
			},
			"items": schema.SetNestedAttribute{
				MarkdownDescription: "List of the groups of the classifier",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Group id",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Group name",
							Computed:            true,
						},
						"parent_id": schema.StringAttribute{
							MarkdownDescription: "Id of the parent group, not set for a root group",
							Computed:            true,
						},
						"path": schema.StringAttribute{
							MarkdownDescription: "Names of the group and its parents, separated by `/`, for example `Location/All Locations/AWS`",
							Computed:            true,
						},
						"discovery_type": schema.StringAttribute{
							MarkdownDescription: "How the group came to exist, for example `StaticSystem` for a built-in group, `Static` for a group created by a user and `Dynamic` for a group the system maintains",
							Computed:            true,
						},
						"subgroups_allowed": schema.BoolAttribute{
							MarkdownDescription: "Whether subgroups can be created under the group",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *GroupingsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.clients = req.ProviderData.(*CwProviderData).Clients
}

func (d *GroupingsDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// One call for each root group returns the whole hierarchy below it (direct=false), so the number of calls does not grow with the number of groups.
func listGroups(c *cw.Client, classifier string) ([]GroupingsItems, error) {
	roots, err := c.Get(groupingPath + "/group/root/" + url.PathEscape(classifier) + "/uuid")
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	items := []GroupingsItems{}
	var walk func(group gjson.Result, parentID, parentPath string)
	walk = func(group gjson.Result, parentID, parentPath string) {
		id := group.Get("uuid").String()
		if id == "" || seen[id] {
			return
		}
		seen[id] = true
		name := group.Get("name").String()
		path := name
		if parentPath != "" {
			path = parentPath + "/" + name
		}
		parent := types.StringNull()
		if parentID != "" {
			parent = types.StringValue(parentID)
		}
		items = append(items, GroupingsItems{
			Id:               types.StringValue(id),
			Name:             types.StringValue(name),
			ParentId:         parent,
			Path:             types.StringValue(path),
			DiscoveryType:    types.StringValue(group.Get("discoveryType").String()),
			SubgroupsAllowed: types.BoolValue(group.Get("operations.subGrp").Bool()),
		})
		group.Get("children").ForEach(func(_, child gjson.Result) bool {
			walk(child, id, path)
			return true
		})
	}

	for _, root := range roots.Array() {
		res, err := c.Get(groupingPath + "/groups/" + url.PathEscape(root.String()) + "?direct=false")
		if err != nil {
			return nil, err
		}
		res.ForEach(func(_, group gjson.Result) bool {
			if group.Get("uuid").String() == root.String() {
				walk(group, "", "")
			}
			return true
		})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Path.ValueString() != items[j].Path.ValueString() {
			return items[i].Path.ValueString() < items[j].Path.ValueString()
		}
		return items[i].Id.ValueString() < items[j].Id.ValueString()
	})
	return items, nil
}

func (d *GroupingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config Groupings
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	c, ok := d.clientFor(ctx, config.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	classifier := config.Classifier.ValueString()
	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", classifier))

	items, err := listGroups(&c, classifier)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve the groups, got error: %s", err))
		return
	}
	config.Id = types.StringValue(classifier)
	config.Items = items
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", classifier))
}

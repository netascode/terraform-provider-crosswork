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
	"regexp"
	"strings"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

// End of section. //template:end imports

// Written by hand: system_defined is computed only, which the generator cannot express (it makes a computed attribute optional
// too).

var _ resource.Resource = &DataGatewayDestinationsResource{}
var _ resource.ResourceWithImportState = &DataGatewayDestinationsResource{}

func NewDataGatewayDestinationsResource() resource.Resource {
	return &DataGatewayDestinationsResource{}
}

type DataGatewayDestinationsResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *DataGatewayDestinationsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_gateway_destinations"
}

func (r *DataGatewayDestinationsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a data destination of the Crosswork Data Gateways.").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Destination name, made of letters, digits, underscores, hyphens and dots").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z0-9_.-]+$`), ""),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"data_source": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Source of the data sent to the destination").AddStringEnumDescription("data-gateway", "application").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("data-gateway", "application"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Type of the destination server").AddStringEnumDescription("kafka", "grpc").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("kafka", "grpc"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"encoding": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Encoding of the data, for the data source `data-gateway`").AddStringEnumDescription("gpbkv", "json").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("gpbkv", "json"),
				},
			},
			"compression_method": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Compression of the data, `snappy`, `gzip`, `zstd` or `none` for Kafka and `gzip`, `deflate` or `snappy` for gRPC").AddStringEnumDescription("snappy", "gzip", "zstd", "none", "deflate").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("snappy", "gzip", "zstd", "none", "deflate"),
				},
			},
			"dispatch_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Dispatch type of a gRPC destination").AddStringEnumDescription("stream", "unary").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("stream", "unary"),
				},
			},
			"max_message_size": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Maximum message size of a Kafka destination in bytes").AddIntegerRangeDescription(1000000, 100000000).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1000000, 100000000),
				},
			},
			"buffer_memory": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Buffer memory of a Kafka destination in bytes").AddIntegerRangeDescription(52428800, 314572800).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(52428800, 314572800),
				},
			},
			"batch_size": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Batch size of a Kafka destination in bytes, at most the buffer memory").AddIntegerRangeDescription(16384, 52428800).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(16384, 52428800),
				},
			},
			"linger_time": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Linger time of a Kafka destination in milliseconds").AddIntegerRangeDescription(0, 5000).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 5000),
				},
			},
			"request_timeout": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Request timeout of a Kafka destination in seconds").AddIntegerRangeDescription(30, 60).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(30, 60),
				},
			},
			"collectors": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Values of a Kafka destination that replace the global ones for one collector").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Collector").AddStringEnumDescription("CLI", "GNMI", "MDT", "SNMP", "SYSLOG").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("CLI", "GNMI", "MDT", "SNMP", "SYSLOG"),
							},
						},
						"buffer_memory": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Buffer memory in bytes").AddIntegerRangeDescription(52428800, 314572800).String,
							Optional:            true,
							Validators: []validator.Int64{
								int64validator.Between(52428800, 314572800),
							},
						},
						"batch_size": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Batch size in bytes").AddIntegerRangeDescription(16384, 52428800).String,
							Optional:            true,
							Validators: []validator.Int64{
								int64validator.Between(16384, 52428800),
							},
						},
						"linger_time": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Linger time in milliseconds").AddIntegerRangeDescription(0, 5000).String,
							Optional:            true,
							Validators: []validator.Int64{
								int64validator.Between(0, 5000),
							},
						},
						"request_timeout": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Request timeout in seconds").AddIntegerRangeDescription(30, 60).String,
							Optional:            true,
							Validators: []validator.Int64{
								int64validator.Between(30, 60),
							},
						},
					},
				},
			},
			"connections": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Addresses of the destination server, each an IPv4 address with a prefix length or a host name with a domain name").String,
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("IPv4 address").String,
							Optional:            true,
						},
						"prefix_length": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Prefix length of the IPv4 address").AddIntegerRangeDescription(1, 32).String,
							Optional:            true,
							Validators: []validator.Int64{
								int64validator.Between(1, 32),
							},
						},
						"host_name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Host name").String,
							Optional:            true,
						},
						"domain_name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Domain name").String,
							Optional:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Port of the server").AddIntegerRangeDescription(1024, 65535).String,
							Required:            true,
							Validators: []validator.Int64{
								int64validator.Between(1024, 65535),
							},
						},
					},
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"secure": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Enable secure communication").String,
				Optional:            true,
			},
			"auth_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Authentication of the secure communication; a gRPC destination takes only `mutual`").AddStringEnumDescription("mutual", "server").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("mutual", "server"),
				},
			},
			"system_defined": schema.BoolAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Whether Crosswork defined the destination, which makes it read-only").String,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DataGatewayDestinationsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

func (r *DataGatewayDestinationsResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// End of section. //template:end model

// Written by hand: the answer of a create is a job with no id, so the destination is found by its name in the list.
func (r *DataGatewayDestinationsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan DataGatewayDestinations

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Name.ValueString()))

	id, err := createDestination(ctx, &c, plan)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (POST), got error: %s", err))
		return
	}
	plan.Id = types.StringValue(id)
	plan.SystemDefined = types.BoolValue(false)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end create

// Written by hand: the query takes a limit and no id, so the destination is found by its id in the list.
func (r *DataGatewayDestinationsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state DataGatewayDestinations

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", state.Id.String()))

	res, err := listDestinations(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
		return
	}
	row := findDestination(res, state.Id.ValueString(), "")
	if !row.Exists() {
		resp.State.RemoveResource(ctx)
		return
	}

	if state.isNull(ctx, row) {
		state.fromBody(ctx, row)
	} else {
		state.updateFromBody(ctx, row)
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

// Written by hand: the whole destination is sent with its id; the controller replaces the properties and the connections.
func (r *DataGatewayDestinationsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state DataGatewayDestinations

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))

	if err := updateDestination(&c, plan, state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s", err))
		return
	}
	plan.Id = state.Id
	plan.SystemDefined = types.BoolValue(false)

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end update

// Written by hand: the destination is deleted by its id; the controller refuses a system-defined destination and its message is
// returned unchanged.
func (r *DataGatewayDestinationsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state DataGatewayDestinations

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))

	res, err := listDestinations(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET) before delete, got error: %s, %s", err, res.String()))
		return
	}
	if !findDestination(res, state.Id.ValueString(), "").Exists() {
		tflog.Debug(ctx, fmt.Sprintf("%s: Resource not found during delete verification, removing from state only", state.Id.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}
	if err := deleteDestination(&c, state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s", err))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Written by hand: the import identifier is the name of the destination or its id, and the id of the state is its uuid.
func (r *DataGatewayDestinationsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <name> or <id>. Got: %q", req.ID),
		)
		return
	}
	c, ok := r.clientFor(ctx, "", &resp.Diagnostics)
	if !ok {
		return
	}
	id, err := destinationIdByNameOrId(&c, idParts[0])
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s", err))
		return
	}
	if id == "" {
		resp.Diagnostics.AddError("Destination not found", fmt.Sprintf("There is no destination named %s on the cluster", idParts[0]))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *DataGatewayDestinationsResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state DataGatewayDestinations, params string) (cw.Res, error) {
	var err error
	cacheKey := "DataGatewayDestinations::" + state.Instance.ValueString() + "::"

	_, cacheSuffix, found := strings.Cut(params, "?")
	queryPart, err := url.ParseQuery(cacheSuffix)
	if err == nil {
		delete(queryPart, "id")
		newQuery := queryPart.Encode()
		cacheSuffix = "?" + newQuery
		cacheKey += cacheSuffix
	}

	cachedValue, found := r.cache.Get(cacheKey)
	if found {
		tflog.Debug(ctx, fmt.Sprintf("hit cache for %s", cacheKey))
		cwRes, ok := cachedValue.(cw.Res)
		if ok {
			return cwRes, nil
		}
		tflog.Info(ctx, fmt.Sprintf("Invalid cache entry type for %s", cacheKey))
		r.cache.Delete(cacheKey)
	}
	res, err := state.readGet(c, "/crosswork/dg-manager/v1/destinations/query"+params)
	singleRes := res
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}

// End of section. //template:end readcache

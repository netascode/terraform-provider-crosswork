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
	"strings"

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &L2VPNY1731ProfilesResource{}
var _ resource.ResourceWithImportState = &L2VPNY1731ProfilesResource{}

func NewL2VPNY1731ProfilesResource() resource.Resource {
	return &L2VPNY1731ProfilesResource{}
}

type L2VPNY1731ProfilesResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *L2VPNY1731ProfilesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_l2vpn_y1731_profiles"
}

func (r *L2VPNY1731ProfilesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a Y.1731 profile of the L2VPN.").String,

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
				MarkdownDescription: helpers.NewAttributeDescription("Name of the profile. Any ASCII character except `%`, `?` and `#`.").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Kind of measurement. `delay_statistics` belong to `delay`, `loss_statistics` to `loss` and `synthetic-loss`.").AddStringEnumDescription("delay", "loss", "synthetic-loss").String,
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("delay", "loss", "synthetic-loss"),
				},
			},
			"schedule_interval": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Minutes between the operations. It has to be a factor of 1440 (1, 2, 3, 4, 5, 6, 8, 9, 10, 12, 15, 16, 18, 20, 24, 30, 32, 36, 40, 45, 48, 60, 80 or 90).").AddIntegerRangeDescription(1, 90).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 90),
				},
			},
			"schedule_duration": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Minutes that the operations last, 1 to 1440, or `forever`").String,
				Optional:            true,
			},
			"probe_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Kind of probe, `burst` when left out. With `packet` NSO drops the message count and the message period without an error.").AddStringEnumDescription("burst", "packet").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("burst", "packet"),
				},
			},
			"probe_message_count": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Messages of a burst. Message count times message period divided by 1000 must not exceed the measurement interval.").AddIntegerRangeDescription(2, 1200).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(2, 1200),
				},
			},
			"probe_message_period": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Milliseconds between the messages of a burst").AddIntegerRangeDescription(50, 30000).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(50, 30000),
				},
			},
			"probe_measurement_interval": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Seconds of a measurement, 1 to 3600, or `once` (only with the probe type `packet`)").String,
				Optional:            true,
			},
			"probe_frame_size": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Size of the frame in bytes.").AddIntegerRangeDescription(1, 9000).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 9000),
				},
			},
			"probe_priority": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Priority of the frames").AddIntegerRangeDescription(0, 7).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 7),
				},
			},
			"delay_statistics": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Delay statistics, only for the type `delay`").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Statistic").AddStringEnumDescription("delay-ds", "delay-sd", "delay-two-way", "jitter-ds", "jitter-sd", "jitter-two-way").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("delay-ds", "delay-sd", "delay-two-way", "jitter-ds", "jitter-sd", "jitter-two-way"),
							},
						},
					},
				},
			},
			"loss_statistics": schema.SetNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Loss statistics, for the types `loss` and `synthetic-loss`").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Statistic").AddStringEnumDescription("loss-ds", "loss-sd").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("loss-ds", "loss-sd"),
							},
						},
					},
				},
			},
			"bucket_size": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Number of measurements in a bucket").AddIntegerRangeDescription(1, 100).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 100),
				},
			},
			"bucket_archive": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Number of buckets kept").AddIntegerRangeDescription(1, 100).String,
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(1, 100),
				},
			},
		},
	}
}

func (r *L2VPNY1731ProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *L2VPNY1731ProfilesResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create
func (r *L2VPNY1731ProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan L2VPNY1731Profiles

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, L2VPNY1731Profiles{})

	params := ""
	res, err := c.Post("/crosswork/proxy/nso/restconf/data/ietf-l2vpn-ntw:l2vpn-ntw"+params, body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "POST", err, res.String()))
		return
	}
	plan.Id = types.StringValue(fmt.Sprint(plan.Name.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end create

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (r *L2VPNY1731ProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state L2VPNY1731Profiles

	// Read state
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

	params := ""
	params += "?content=config"
	res, err := c.Get(state.getPath() + params)
	if err != nil && (strings.Contains(err.Error(), "StatusCode 404") || strings.Contains(err.Error(), "StatusCode 406") || strings.Contains(err.Error(), "StatusCode 500") || strings.Contains(err.Error(), "StatusCode 400") || cw.IsNotFound(err)) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
		return
	}

	// If every attribute is set to null we are dealing with an import operation and therefore reading all attributes
	if state.isNull(ctx, res) {
		state.fromBody(ctx, res)
	} else {
		state.updateFromBody(ctx, res)
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

// Section below is generated&owned by "gen/generator.go". //template:begin update
func (r *L2VPNY1731ProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state L2VPNY1731Profiles

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read state
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

	body := plan.toBody(ctx, state)
	params := ""
	res, err := c.Put(plan.getPath()+params, body)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s, %s", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end update

// Section below is generated&owned by "gen/generator.go". //template:begin delete
func (r *L2VPNY1731ProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state L2VPNY1731Profiles

	// Read state
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
	res, err := c.Delete(state.getPath())
	if err != nil && !strings.Contains(err.Error(), "StatusCode 404") && !cw.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s, %s", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Section below is generated&owned by "gen/generator.go". //template:begin import
func (r *L2VPNY1731ProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
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

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *L2VPNY1731ProfilesResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state L2VPNY1731Profiles, params string) (cw.Res, error) {
	var err error
	cacheKey := "L2VPNY1731Profiles::" + state.Instance.ValueString() + "::"

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
	res, err := c.Get(state.getPath() + params)
	singleRes := res
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}

// End of section. //template:end readcache

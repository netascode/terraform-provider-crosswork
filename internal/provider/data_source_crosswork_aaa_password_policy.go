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
	_ datasource.DataSource              = &AAAPasswordPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &AAAPasswordPolicyDataSource{}
)

func NewAAAPasswordPolicyDataSource() datasource.DataSource {
	return &AAAPasswordPolicyDataSource{}
}

type AAAPasswordPolicyDataSource struct {
	clients map[string]*cw.Client
}

func (d *AAAPasswordPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_aaa_password_policy"
}

func (d *AAAPasswordPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the AAA local password policy.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"min_password_length": schema.Int64Attribute{
				MarkdownDescription: "Minimum number of characters of a password",
				Computed:            true,
			},
			"no_username": schema.BoolAttribute{
				MarkdownDescription: "Reject a password that contains the user name or its reverse",
				Computed:            true,
			},
			"no_cisco_variant": schema.BoolAttribute{
				MarkdownDescription: "Reject a password that contains `cisco`, `ocsic` or a capitalized variant of them",
				Computed:            true,
			},
			"no_char_repetition": schema.BoolAttribute{
				MarkdownDescription: "Reject a password that repeats a character more than three times in a row",
				Computed:            true,
			},
			"num_changed_chars_enabled": schema.BoolAttribute{
				MarkdownDescription: "Require a number of changed characters in a new password",
				Computed:            true,
			},
			"num_changed_chars": schema.Int64Attribute{
				MarkdownDescription: "Number of characters a new password must change",
				Computed:            true,
			},
			"num_reuse_limit_enabled": schema.BoolAttribute{
				MarkdownDescription: "Limit the reuse of the last passwords",
				Computed:            true,
			},
			"num_reuse_limit": schema.Int64Attribute{
				MarkdownDescription: "Number of last passwords that cannot be reused",
				Computed:            true,
			},
			"password_reuse_days_enabled": schema.BoolAttribute{
				MarkdownDescription: "Limit the reuse of passwords used in the last days",
				Computed:            true,
			},
			"password_reuse_days": schema.Int64Attribute{
				MarkdownDescription: "Number of days a used password cannot be reused",
				Computed:            true,
			},
			"failed_logins_before_lockout_enabled": schema.BoolAttribute{
				MarkdownDescription: "Lock a user out after failed login attempts",
				Computed:            true,
			},
			"failed_logins_before_lockout": schema.Int64Attribute{
				MarkdownDescription: "Number of failed login attempts before the lockout",
				Computed:            true,
			},
			"lock_out_user_time_enabled": schema.BoolAttribute{
				MarkdownDescription: "Lock a user out for a time",
				Computed:            true,
			},
			"lock_out_user_time": schema.Int64Attribute{
				MarkdownDescription: "Minutes a user stays locked out",
				Computed:            true,
			},
			"password_expiry_days_enabled": schema.BoolAttribute{
				MarkdownDescription: "Expire passwords after a number of days",
				Computed:            true,
			},
			"password_expiry_days": schema.Int64Attribute{
				MarkdownDescription: "Days after which a password expires",
				Computed:            true,
			},
			"days_for_warning_enabled": schema.BoolAttribute{
				MarkdownDescription: "Warn users before their password expires",
				Computed:            true,
			},
			"days_for_warning": schema.Int64Attribute{
				MarkdownDescription: "Days before the expiry that a warning is shown",
				Computed:            true,
			},
			"change_password_on_first_login": schema.BoolAttribute{
				MarkdownDescription: "Require a new password at the first login",
				Computed:            true,
			},
		},
	}
}

func (d *AAAPasswordPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *AAAPasswordPolicyDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (d *AAAPasswordPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config AAAPasswordPolicy

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
	res, err := c.Get("/crosswork/aaa/v1/passwordPolicyConfig" + params)
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

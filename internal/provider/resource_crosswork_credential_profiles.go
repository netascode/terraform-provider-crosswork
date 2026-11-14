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
var _ resource.Resource = &CredentialProfilesResource{}
var _ resource.ResourceWithImportState = &CredentialProfilesResource{}

func NewCredentialProfilesResource() resource.Resource {
	return &CredentialProfilesResource{}
}

type CredentialProfilesResource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
}

func (r *CredentialProfilesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_profiles"
}

func (r *CredentialProfilesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This resource can manage a credential profile.").String,

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
			"profile": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Credential profile name").String,
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"credentials": schema.ListNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("List of credentials, one per access type").String,
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Access type").AddStringEnumDescription("ROBOT_USERPASS_TCP", "ROBOT_USERPASS_UDP", "ROBOT_USERPASS_HTTP", "ROBOT_USERPASS_HTTPS", "ROBOT_USERPASS_GRPC", "ROBOT_USERPASS_SSH", "ROBOT_USERPASS_NETCONF", "ROBOT_USERPASS_TELNET", "ROBOT_USERPASS_SNMP", "ROBOT_USERPASS_TL1", "ROBOT_USERPASS_TL1_SECURE", "ROBOT_USERPASS_ADMIN", "ROBOT_USERPASS_ENABLE", "ROBOT_USERPASS_GNMI").String,
							Required:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("ROBOT_USERPASS_TCP", "ROBOT_USERPASS_UDP", "ROBOT_USERPASS_HTTP", "ROBOT_USERPASS_HTTPS", "ROBOT_USERPASS_GRPC", "ROBOT_USERPASS_SSH", "ROBOT_USERPASS_NETCONF", "ROBOT_USERPASS_TELNET", "ROBOT_USERPASS_SNMP", "ROBOT_USERPASS_TL1", "ROBOT_USERPASS_TL1_SECURE", "ROBOT_USERPASS_ADMIN", "ROBOT_USERPASS_ENABLE", "ROBOT_USERPASS_GNMI"),
							},
						},
						"user_name": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("User name").String,
							Required:            true,
						},
						"password": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Password").AddMutualExclusivityDescription("Only one of `password` and `password_wo` can be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `password_wo` together with `password_wo_version`, which keeps it out of state.").String,
							Sensitive:           true,
							Optional:            true,
						},
						"password_wo": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Password").AddMutualExclusivityDescription("Only one of `password` and `password_wo` can be set.").String,
							Optional:            true,
							WriteOnly:           true,
							Sensitive:           true,
						},
						"password_wo_version": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `password_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
							Optional:            true,
						},
						"enable_password": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Second level authentication, for devices that support it").AddStringEnumDescription("ROBOT_USERPASSEN_ENABLE", "ROBOT_USERPASSEN_DISABLE").String,
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("ROBOT_USERPASSEN_ENABLE", "ROBOT_USERPASSEN_DISABLE"),
							},
						},
						"enable_password_value": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Enable password").AddMutualExclusivityDescription("Only one of `enable_password_value` and `enable_password_value_wo` can be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `enable_password_value_wo` together with `enable_password_value_wo_version`, which keeps it out of state.").String,
							Sensitive:           true,
							Optional:            true,
						},
						"enable_password_value_wo": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Enable password").AddMutualExclusivityDescription("Only one of `enable_password_value` and `enable_password_value_wo` can be set.").String,
							Optional:            true,
							WriteOnly:           true,
							Sensitive:           true,
						},
						"enable_password_value_wo_version": schema.Int64Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `enable_password_value_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
							Optional:            true,
						},
					},
				},
			},
			"snmpv2_read_community": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMPv2 read community").AddMutualExclusivityDescription("Only one of `snmpv2_read_community` and `snmpv2_read_community_wo` can be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `snmpv2_read_community_wo` together with `snmpv2_read_community_wo_version`, which keeps it out of state.").String,
				Sensitive:           true,
				Optional:            true,
			},
			"snmpv2_read_community_wo": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMPv2 read community").AddMutualExclusivityDescription("Only one of `snmpv2_read_community` and `snmpv2_read_community_wo` can be set.").String,
				Optional:            true,
				WriteOnly:           true,
				Sensitive:           true,
			},
			"snmpv2_read_community_wo_version": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `snmpv2_read_community_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
				Optional:            true,
			},
			"snmpv2_write_community": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMPv2 write community").AddMutualExclusivityDescription("Only one of `snmpv2_write_community` and `snmpv2_write_community_wo` can be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `snmpv2_write_community_wo` together with `snmpv2_write_community_wo_version`, which keeps it out of state.").String,
				Sensitive:           true,
				Optional:            true,
			},
			"snmpv2_write_community_wo": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMPv2 write community").AddMutualExclusivityDescription("Only one of `snmpv2_write_community` and `snmpv2_write_community_wo` can be set.").String,
				Optional:            true,
				WriteOnly:           true,
				Sensitive:           true,
			},
			"snmpv2_write_community_wo_version": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `snmpv2_write_community_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
				Optional:            true,
			},
			"snmpv3_user_name": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMPv3 user name").String,
				Optional:            true,
			},
			"snmpv3_security_level": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMPv3 security level").AddStringEnumDescription("SL_NO_AUTH_NO_PRIV", "SL_AUTH_NO_PRIV", "SL_AUTH_PRIV").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("SL_NO_AUTH_NO_PRIV", "SL_AUTH_NO_PRIV", "SL_AUTH_PRIV"),
				},
			},
			"snmpv3_auth_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMPv3 authentication type").AddStringEnumDescription("AT_HMAC_MD5", "AT_HMAC_SHA", "AT_HMAC_SHA2_512", "AT_HMAC_SHA2_384", "AT_HMAC_SHA2_256", "AT_HMAC_SHA2_224").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("AT_HMAC_MD5", "AT_HMAC_SHA", "AT_HMAC_SHA2_512", "AT_HMAC_SHA2_384", "AT_HMAC_SHA2_256", "AT_HMAC_SHA2_224"),
				},
			},
			"snmpv3_auth_password": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMPv3 authentication password").AddMutualExclusivityDescription("Only one of `snmpv3_auth_password` and `snmpv3_auth_password_wo` can be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `snmpv3_auth_password_wo` together with `snmpv3_auth_password_wo_version`, which keeps it out of state.").String,
				Sensitive:           true,
				Optional:            true,
			},
			"snmpv3_auth_password_wo": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMPv3 authentication password").AddMutualExclusivityDescription("Only one of `snmpv3_auth_password` and `snmpv3_auth_password_wo` can be set.").String,
				Optional:            true,
				WriteOnly:           true,
				Sensitive:           true,
			},
			"snmpv3_auth_password_wo_version": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `snmpv3_auth_password_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
				Optional:            true,
			},
			"snmpv3_priv_type": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMPv3 privacy type").AddStringEnumDescription("PT_CFB_AES_128", "PT_CBC_DES_56", "PT_CFB_AES_192", "PT_CFB_AES_256", "PT_CBC_3_DES").String,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("PT_CFB_AES_128", "PT_CBC_DES_56", "PT_CFB_AES_192", "PT_CFB_AES_256", "PT_CBC_3_DES"),
				},
			},
			"snmpv3_priv_password": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMPv3 privacy password").AddMutualExclusivityDescription("Only one of `snmpv3_priv_password` and `snmpv3_priv_password_wo` can be set.").AddCoexistenceNote("This attribute stores the secret in Terraform state. Prefer `snmpv3_priv_password_wo` together with `snmpv3_priv_password_wo_version`, which keeps it out of state.").String,
				Sensitive:           true,
				Optional:            true,
			},
			"snmpv3_priv_password_wo": schema.StringAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("SNMPv3 privacy password").AddMutualExclusivityDescription("Only one of `snmpv3_priv_password` and `snmpv3_priv_password_wo` can be set.").String,
				Optional:            true,
				WriteOnly:           true,
				Sensitive:           true,
			},
			"snmpv3_priv_password_wo_version": schema.Int64Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("Rotation trigger for `snmpv3_priv_password_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.").String,
				Optional:            true,
			},
		},
	}
}

func (r *CredentialProfilesResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *CredentialProfilesResource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}

// ValidateConfig enforces the relationship between a secret attribute, its write-only
// "_wo" counterpart and the "_wo_version" rotation trigger.
//
// The checks are not schema validators because Terraform renders an attribute-scoped
// diagnostic together with the offending configuration line, which for a secret prints
// the value into plan output and logs. A resource-scoped diagnostic names the attributes
// (and the list element index) instead.
func (r *CredentialProfilesResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var legacySnmpv2ReadCommunity types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv2_read_community"), &legacySnmpv2ReadCommunity)...)
	var woSnmpv2ReadCommunity types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv2_read_community_wo"), &woSnmpv2ReadCommunity)...)
	var woVersionSnmpv2ReadCommunity types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv2_read_community_wo_version"), &woVersionSnmpv2ReadCommunity)...)
	if !legacySnmpv2ReadCommunity.IsUnknown() && !woSnmpv2ReadCommunity.IsUnknown() && !legacySnmpv2ReadCommunity.IsNull() && !woSnmpv2ReadCommunity.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Only one of `snmpv2_read_community` and `snmpv2_read_community_wo` can be set.",
		)
	}
	if !woSnmpv2ReadCommunity.IsUnknown() && !woVersionSnmpv2ReadCommunity.IsUnknown() && !woSnmpv2ReadCommunity.IsNull() && woVersionSnmpv2ReadCommunity.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"`snmpv2_read_community_wo_version` must be set when `snmpv2_read_community_wo` is used. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.",
		)
	}
	var legacySnmpv2WriteCommunity types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv2_write_community"), &legacySnmpv2WriteCommunity)...)
	var woSnmpv2WriteCommunity types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv2_write_community_wo"), &woSnmpv2WriteCommunity)...)
	var woVersionSnmpv2WriteCommunity types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv2_write_community_wo_version"), &woVersionSnmpv2WriteCommunity)...)
	if !legacySnmpv2WriteCommunity.IsUnknown() && !woSnmpv2WriteCommunity.IsUnknown() && !legacySnmpv2WriteCommunity.IsNull() && !woSnmpv2WriteCommunity.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Only one of `snmpv2_write_community` and `snmpv2_write_community_wo` can be set.",
		)
	}
	if !woSnmpv2WriteCommunity.IsUnknown() && !woVersionSnmpv2WriteCommunity.IsUnknown() && !woSnmpv2WriteCommunity.IsNull() && woVersionSnmpv2WriteCommunity.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"`snmpv2_write_community_wo_version` must be set when `snmpv2_write_community_wo` is used. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.",
		)
	}
	var legacySnmpv3AuthPassword types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv3_auth_password"), &legacySnmpv3AuthPassword)...)
	var woSnmpv3AuthPassword types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv3_auth_password_wo"), &woSnmpv3AuthPassword)...)
	var woVersionSnmpv3AuthPassword types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv3_auth_password_wo_version"), &woVersionSnmpv3AuthPassword)...)
	if !legacySnmpv3AuthPassword.IsUnknown() && !woSnmpv3AuthPassword.IsUnknown() && !legacySnmpv3AuthPassword.IsNull() && !woSnmpv3AuthPassword.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Only one of `snmpv3_auth_password` and `snmpv3_auth_password_wo` can be set.",
		)
	}
	if !woSnmpv3AuthPassword.IsUnknown() && !woVersionSnmpv3AuthPassword.IsUnknown() && !woSnmpv3AuthPassword.IsNull() && woVersionSnmpv3AuthPassword.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"`snmpv3_auth_password_wo_version` must be set when `snmpv3_auth_password_wo` is used. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.",
		)
	}
	var legacySnmpv3PrivPassword types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv3_priv_password"), &legacySnmpv3PrivPassword)...)
	var woSnmpv3PrivPassword types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv3_priv_password_wo"), &woSnmpv3PrivPassword)...)
	var woVersionSnmpv3PrivPassword types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv3_priv_password_wo_version"), &woVersionSnmpv3PrivPassword)...)
	if !legacySnmpv3PrivPassword.IsUnknown() && !woSnmpv3PrivPassword.IsUnknown() && !legacySnmpv3PrivPassword.IsNull() && !woSnmpv3PrivPassword.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Only one of `snmpv3_priv_password` and `snmpv3_priv_password_wo` can be set.",
		)
	}
	if !woSnmpv3PrivPassword.IsUnknown() && !woVersionSnmpv3PrivPassword.IsUnknown() && !woSnmpv3PrivPassword.IsNull() && woVersionSnmpv3PrivPassword.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"`snmpv3_priv_password_wo_version` must be set when `snmpv3_priv_password_wo` is used. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.",
		)
	}
	{
		// Secrets nested in "credentials" are validated per element. A list that cannot be
		// read as a whole - because it is still unknown at validation time - is skipped
		// rather than reported, since there is nothing to check yet.
		var cfgCredentials []CredentialProfilesCredentials
		if diags := req.Config.GetAttribute(ctx, path.Root("credentials"), &cfgCredentials); !diags.HasError() {
			for i := range cfgCredentials {
				if !cfgCredentials[i].Password.IsUnknown() && !cfgCredentials[i].PasswordWo.IsUnknown() && !cfgCredentials[i].Password.IsNull() && !cfgCredentials[i].PasswordWo.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("Only one of `password` and `password_wo` can be set in `credentials` element %d.", i),
					)
				}
				if !cfgCredentials[i].PasswordWo.IsUnknown() && !cfgCredentials[i].PasswordWoVersion.IsUnknown() && !cfgCredentials[i].PasswordWo.IsNull() && cfgCredentials[i].PasswordWoVersion.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("`password_wo_version` must be set when `password_wo` is used in `credentials` element %d. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.", i),
					)
				}
				if !cfgCredentials[i].EnablePasswordValue.IsUnknown() && !cfgCredentials[i].EnablePasswordValueWo.IsUnknown() && !cfgCredentials[i].EnablePasswordValue.IsNull() && !cfgCredentials[i].EnablePasswordValueWo.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("Only one of `enable_password_value` and `enable_password_value_wo` can be set in `credentials` element %d.", i),
					)
				}
				if !cfgCredentials[i].EnablePasswordValueWo.IsUnknown() && !cfgCredentials[i].EnablePasswordValueWoVersion.IsUnknown() && !cfgCredentials[i].EnablePasswordValueWo.IsNull() && cfgCredentials[i].EnablePasswordValueWoVersion.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("`enable_password_value_wo_version` must be set when `enable_password_value_wo` is used in `credentials` element %d. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.", i),
					)
				}
			}
		}
	}
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create
func (r *CredentialProfilesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan CredentialProfiles

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
	// Write-only values are not in plan or state; read them from config.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv2_read_community_wo"), &plan.Snmpv2ReadCommunityWo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Write-only values are not in plan or state; read them from config.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv2_write_community_wo"), &plan.Snmpv2WriteCommunityWo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Write-only values are not in plan or state; read them from config.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv3_auth_password_wo"), &plan.Snmpv3AuthPasswordWo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Write-only values are not in plan or state; read them from config.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv3_priv_password_wo"), &plan.Snmpv3PrivPasswordWo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Write-only values in the list are not in plan or state; copy them from the config list into the plan elements.
	{
		var cfgCredentials []CredentialProfilesCredentials
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("credentials"), &cfgCredentials)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for i := range plan.Credentials {
			if i < len(cfgCredentials) {
				plan.Credentials[i].PasswordWo = cfgCredentials[i].PasswordWo
				plan.Credentials[i].EnablePasswordValueWo = cfgCredentials[i].EnablePasswordValueWo
			}
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, CredentialProfiles{})

	params := ""
	res, err := c.Post(plan.getPath()+params, body)
	if err != nil {
		if r.AllowExistingOnCreate {
			tflog.Info(ctx, fmt.Sprintf("Failed to configure object (%s), got error: %s, %s. allow_existing_on_create is true, beginning update", "POST", err, res.String()))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "POST", err, res.String()))
			return
		}
	}
	plan.Id = types.StringValue(fmt.Sprint(plan.Profile.ValueString()))
	if !r.AllowExistingOnCreate {
		tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))
	} else {
		params = ""
		body = plan.toBody(ctx, CredentialProfiles{Id: plan.Id})
		res, err = c.Put(plan.getPath()+params, body)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "PUT", err, res.String()))
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("%s: Fallback to update existing resource finished successfully", plan.Id.ValueString()))
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end create

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (r *CredentialProfilesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state CredentialProfiles

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
	res, err := state.readGet(&c, "/crosswork/inventory/v1/credentials/query"+params)
	if err != nil && (strings.Contains(err.Error(), "StatusCode 404") || strings.Contains(err.Error(), "StatusCode 406") || strings.Contains(err.Error(), "StatusCode 500") || strings.Contains(err.Error(), "StatusCode 400") || cw.IsNotFound(err)) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
		return
	}

	if data := res.Get("data"); !data.Exists() || len(data.Array()) == 0 {
		resp.State.RemoveResource(ctx)
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
func (r *CredentialProfilesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state CredentialProfiles

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
	// Read on every Update: a full-object replace PUT is rejected when it omits the secret; the version companion only drives change detection.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv2_read_community_wo"), &plan.Snmpv2ReadCommunityWo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read on every Update: a full-object replace PUT is rejected when it omits the secret; the version companion only drives change detection.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv2_write_community_wo"), &plan.Snmpv2WriteCommunityWo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read on every Update: a full-object replace PUT is rejected when it omits the secret; the version companion only drives change detection.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv3_auth_password_wo"), &plan.Snmpv3AuthPasswordWo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read on every Update: a full-object replace PUT is rejected when it omits the secret; the version companion only drives change detection.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("snmpv3_priv_password_wo"), &plan.Snmpv3PrivPasswordWo)...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read unconditionally, for the same reason as the top-level secrets above.
	{
		var cfgCredentials []CredentialProfilesCredentials
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("credentials"), &cfgCredentials)...)
		if resp.Diagnostics.HasError() {
			return
		}
		for i := range plan.Credentials {
			if i < len(cfgCredentials) {
				plan.Credentials[i].PasswordWo = cfgCredentials[i].PasswordWo
				plan.Credentials[i].EnablePasswordValueWo = cfgCredentials[i].EnablePasswordValueWo
			}
		}
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
func (r *CredentialProfilesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state CredentialProfiles

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
	// A delete of a missing object may fail, and a missing object only needs to leave the state.
	verifyParams := ""
	verifyRes, verifyErr := state.readGet(&c, "/crosswork/inventory/v1/credentials/query"+verifyParams)
	if verifyErr != nil && !strings.Contains(verifyErr.Error(), "StatusCode 404") && !cw.IsNotFound(verifyErr) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET) before delete, got error: %s, %s", verifyErr, verifyRes.String()))
		return
	}
	if verifyErr != nil || !verifyRes.Get("data.0.profile").Exists() || verifyRes.Get("data.0.profile").String() != state.Id.ValueString() {
		tflog.Debug(ctx, fmt.Sprintf("%s: Resource not found or ID mismatch during delete verification, removing from state only", state.Id.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Resource verified, proceeding with delete", state.Id.ValueString()))
	res, err := c.DeleteBody(state.getPath(), strings.ReplaceAll(`{"data":[{"profile":"{id}"}]}`, "{id}", state.idJSON()))
	if err != nil && !strings.Contains(err.Error(), "StatusCode 404") && !cw.IsNotFound(err) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s, %s", err, res.String()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Section below is generated&owned by "gen/generator.go". //template:begin import
func (r *CredentialProfilesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")

	if len(idParts) != 1 || idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: <profile>. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("profile"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *CredentialProfilesResource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state CredentialProfiles, params string) (cw.Res, error) {
	var err error
	cacheKey := "CredentialProfiles::" + state.Instance.ValueString() + "::"

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
	res, err := state.readGet(c, "/crosswork/inventory/v1/credentials/query"+params)
	singleRes := res
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}

// End of section. //template:end readcache

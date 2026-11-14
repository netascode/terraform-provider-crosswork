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

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &DataGatewayGlobalParametersDataSource{}
	_ datasource.DataSourceWithConfigure = &DataGatewayGlobalParametersDataSource{}
)

func NewDataGatewayGlobalParametersDataSource() datasource.DataSource {
	return &DataGatewayGlobalParametersDataSource{}
}

type DataGatewayGlobalParametersDataSource struct {
	clients map[string]*cw.Client
}

func (d *DataGatewayGlobalParametersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_data_gateway_global_parameters"
}

func (d *DataGatewayGlobalParametersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the global parameters of the Crosswork Data Gateways, the settings of the Data Collector(s) Global Settings page.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
			},
			"cli_sessions": schema.Int64Attribute{
				MarkdownDescription: "Number of CLI sessions",
				Computed:            true,
			},
			"ssh_session_timeout": schema.Int64Attribute{
				MarkdownDescription: "SSH session timeout in seconds",
				Computed:            true,
			},
			"snmp_trap_port": schema.Int64Attribute{
				MarkdownDescription: "SNMP trap port",
				Computed:            true,
			},
			"syslog_udp_port": schema.Int64Attribute{
				MarkdownDescription: "Syslog UDP port",
				Computed:            true,
			},
			"syslog_tcp_port": schema.Int64Attribute{
				MarkdownDescription: "Syslog TCP port",
				Computed:            true,
			},
			"syslog_tls_port": schema.Int64Attribute{
				MarkdownDescription: "Syslog TLS port",
				Computed:            true,
			},
			"resync_snmpv3_details": schema.BoolAttribute{
				MarkdownDescription: "Whether the SNMPv3 details of a device are synchronised again",
				Computed:            true,
			},
		},
	}
}

func (d *DataGatewayGlobalParametersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *DataGatewayGlobalParametersDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

// Written by hand: the settings have no id of their own, so the data source gives the one of the resource.
func (d *DataGatewayGlobalParametersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config DataGatewayGlobalParameters

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

	res, err := config.readGet(&c, globalParametersQueryPath)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}

	config.fromBody(ctx, res)
	config.Id = types.StringValue("data_gateway_global_parameters")

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

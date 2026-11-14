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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &NetworkDevicesDataSource{}
	_ datasource.DataSourceWithConfigure = &NetworkDevicesDataSource{}
)

func NewNetworkDevicesDataSource() datasource.DataSource {
	return &NetworkDevicesDataSource{}
}

type NetworkDevicesDataSource struct {
	clients map[string]*cw.Client
}

func (d *NetworkDevicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_devices"
}

func (d *NetworkDevicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "This data source can read the Network Devices.",

		Attributes: map[string]schema.Attribute{
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Inventory id of the device. A change replaces the device",
				Computed:            true,
			},
			"host_name": schema.StringAttribute{
				MarkdownDescription: "Host name of the device; needed when the device has no SNMP row. A device with an SNMP row reports its own host name, which replaces the value given here once Crosswork has collected the device; leave it unset to keep the host name the controller holds. A change needs the device admin down, which is done automatically",
				Computed:            true,
			},
			"credential_profile": schema.StringAttribute{
				MarkdownDescription: "Name of the credential profile, which should hold a credential for every protocol the device uses",
				Computed:            true,
			},
			"admin_state": schema.StringAttribute{
				MarkdownDescription: "Admin state of the device. Crosswork keeps a device that is not attached to a data gateway admin down. Attaching a device to a data gateway brings it admin up. The default is `ROBOT_ADMIN_STATE_DOWN`; `ROBOT_ADMIN_STATE_UNMANAGED` cannot be used while a device is attached to a data gateway",
				Computed:            true,
			},
			"reachability_check": schema.StringAttribute{
				MarkdownDescription: "Whether Crosswork checks that the device is reachable. A change needs the device admin down, which is done automatically",
				Computed:            true,
			},
			"serial_number": schema.StringAttribute{
				MarkdownDescription: "Serial number of the device",
				Computed:            true,
			},
			"mac_address": schema.StringAttribute{
				MarkdownDescription: "MAC address of the device",
				Computed:            true,
			},
			"encryption": schema.StringAttribute{
				MarkdownDescription: "SNMP encryption of the device in the lower case form, for example `aes128-ctr`; leave it unset for none",
				Computed:            true,
			},
			"device_type": schema.StringAttribute{
				MarkdownDescription: "Device type",
				Computed:            true,
			},
			"software_type": schema.StringAttribute{
				MarkdownDescription: "Software type, for example `IOS XR`; needed when the device has no SNMP row",
				Computed:            true,
			},
			"software_version": schema.StringAttribute{
				MarkdownDescription: "Software version. Crosswork overwrites it after it discovers the device, so set it only for a device it cannot reach",
				Computed:            true,
			},
			"product_type": schema.StringAttribute{
				MarkdownDescription: "Product type. Crosswork overwrites it after it discovers the device, so set it only for a device it cannot reach",
				Computed:            true,
			},
			"syslog_format": schema.StringAttribute{
				MarkdownDescription: "Syslog format of the device",
				Computed:            true,
			},
			"capabilities": schema.SetAttribute{
				MarkdownDescription: "Capabilities of the device, any of `YANG_MDT`, `YANG_CLI`, `SNMP` and `GNMI`; `SNMP` needs an SNMP row. A change needs the device admin down, which is done automatically",
				ElementType:         types.StringType,
				Computed:            true,
			},
			"cli_cache_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether Crosswork caches the CLI output of the device",
				Computed:            true,
			},
			"provider_profile": schema.StringAttribute{
				MarkdownDescription: "Name of the NSO provider profile that manages the device",
				Computed:            true,
			},
			"connectivity": schema.ListNestedAttribute{
				MarkdownDescription: "Connectivity rows, one per protocol. A change to an address, a prefix length or a port needs the device admin down, which is done automatically",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"protocol": schema.StringAttribute{
							MarkdownDescription: "Protocol to connect to the device",
							Computed:            true,
						},
						"ip_address": schema.StringAttribute{
							MarkdownDescription: "IPv4 address of the device, unique on the Crosswork cluster",
							Computed:            true,
						},
						"prefix_length": schema.Int64Attribute{
							MarkdownDescription: "Prefix length of the address",
							Computed:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "Port of the protocol, for example 22 for SSH and 161 for SNMP; Crosswork refuses a row without one",
							Computed:            true,
						},
						"timeout": schema.Int64Attribute{
							MarkdownDescription: "Timeout in seconds",
							Computed:            true,
						},
						"encoding_type": schema.StringAttribute{
							MarkdownDescription: "Encoding, required for GNMI and GRPC rows",
							Computed:            true,
						},
						"trap_source_ip": schema.StringAttribute{
							MarkdownDescription: "Trap source IPv4 address, for the SNMP row",
							Computed:            true,
						},
						"disable_trap_community_check": schema.BoolAttribute{
							MarkdownDescription: "Whether to skip the trap community check, for the SNMP row",
							Computed:            true,
						},
					},
				},
			},
			"isis_system_id": schema.StringAttribute{
				MarkdownDescription: "ISIS system id",
				Computed:            true,
			},
			"ospf_router_id": schema.StringAttribute{
				MarkdownDescription: "OSPF router id",
				Computed:            true,
			},
			"te_router_id": schema.StringAttribute{
				MarkdownDescription: "TE router id. A change needs the device admin down, which is done automatically",
				Computed:            true,
			},
			"ipv6_router_id": schema.StringAttribute{
				MarkdownDescription: "IPv6 router id",
				Computed:            true,
			},
			"telemetry_vrf": schema.StringAttribute{
				MarkdownDescription: "VRF of the streaming telemetry source interface",
				Computed:            true,
			},
			"telemetry_source_interface": schema.StringAttribute{
				MarkdownDescription: "Source interface of the streaming telemetry",
				Computed:            true,
			},
			"telemetry_opt_out_of_mdt_config": schema.BoolAttribute{
				MarkdownDescription: "Whether to opt out of the MDT configuration",
				Computed:            true,
			},
			"location_building": schema.StringAttribute{
				MarkdownDescription: "Building",
				Computed:            true,
			},
			"location_street": schema.StringAttribute{
				MarkdownDescription: "Street",
				Computed:            true,
			},
			"location_city": schema.StringAttribute{
				MarkdownDescription: "City",
				Computed:            true,
			},
			"location_state": schema.StringAttribute{
				MarkdownDescription: "State",
				Computed:            true,
			},
			"location_country": schema.StringAttribute{
				MarkdownDescription: "Country",
				Computed:            true,
			},
			"location_region": schema.StringAttribute{
				MarkdownDescription: "Region",
				Computed:            true,
			},
			"location_zip": schema.StringAttribute{
				MarkdownDescription: "Zip code",
				Computed:            true,
			},
			"location_latitude": schema.Float64Attribute{
				MarkdownDescription: "Latitude; set it together with the longitude",
				Computed:            true,
			},
			"location_longitude": schema.Float64Attribute{
				MarkdownDescription: "Longitude; set it together with the latitude",
				Computed:            true,
			},
			"location_altitude": schema.Float64Attribute{
				MarkdownDescription: "Altitude",
				Computed:            true,
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "User-defined tags of the device, which must exist as tags. The tags that Crosswork adds itself are never managed",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Tag name",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *NetworkDevicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.clients = req.ProviderData.(*CwProviderData).Clients
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (d *NetworkDevicesDataSource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, d.clients, instance, diags)
}

// End of section. //template:end model

func (d *NetworkDevicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	applyProviderMeta(d.clients, ctx, req.ProviderMeta)
	var config NetworkDevices

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
	res, err := config.readGet(&c, "/crosswork/inventory/v1/nodes/query"+params)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}
	if data := res.Get("data"); !data.Exists() || len(data.Array()) == 0 {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("No object found with id: %s", config.Id.ValueString()))
		return
	}

	config.fromBody(ctx, res)
	// Manually managed: the provider profile is the key of the provider block, which the generated read cannot reach.
	config.ProviderProfile = providerProfileName(res)
	user, err := userTagNames(&c)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object, got error: %s", err))
		return
	}
	config.Tags = tagsOf(res, user, nil)

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", config.Id.ValueString()))

	diags = resp.State.Set(ctx, &config)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end read

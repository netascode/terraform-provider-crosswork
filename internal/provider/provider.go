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

// Section below is generated&owned by "gen/generator.go". //template:begin provider
import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/metaschema"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
)

// Bounds the retry loop of the definition key retry_on_error_codes.
const retryMaxWaitSeconds = 180

// When not empty, a failed delete whose response error code starts with it is only a warning.
const deleteWarnErrorPrefix = ""

type CwProvider struct {
	// "dev" for a local build, "test" in acceptance tests.
	version string
}

type CwProviderModel struct {
	Username              types.String              `tfsdk:"username"`
	Password              types.String              `tfsdk:"password"`
	URL                   types.String              `tfsdk:"url"`
	Insecure              types.Bool                `tfsdk:"insecure"`
	Retries               types.Int64               `tfsdk:"retries"`
	MaxTimeout            types.Int64               `tfsdk:"max_timeout"`
	DeviceReadyMode       types.String              `tfsdk:"device_ready_mode"`
	DeviceReadyTimeout    types.Int64               `tfsdk:"device_ready_timeout"`
	AllowExistingOnCreate types.Bool                `tfsdk:"allow_existing_on_create"`
	Instances             []CwProviderModelInstance `tfsdk:"instances"`
}

type CwProviderModelInstance struct {
	Name types.String `tfsdk:"name"`
	URL  types.String `tfsdk:"url"`
}

type CwProviderData struct {
	// The key "" is the default client.
	Clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	Cache                 *ThreadSafeCache
	DeviceReady           nsoDeviceConfig
}

func (p *CwProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "crosswork"
	resp.Version = p.version
}

func (p *CwProvider) MetaSchema(_ context.Context, _ provider.MetaSchemaRequest, resp *provider.MetaSchemaResponse) {
	resp.Schema = metaschema.Schema{
		Attributes: map[string]metaschema.Attribute{
			"module_name": metaschema.StringAttribute{
				MarkdownDescription: "Identifier a Network-as-Code module sets via `provider_meta` so Crosswork can attribute API traffic to the module. Appended to the HTTP User-Agent. Not intended to be set directly by users.",
				Optional:            true,
			},
		},
	}
}

func (p *CwProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				MarkdownDescription: "Username for the Crosswork instance. This can also be set as the CW_USERNAME environment variable.",
				Optional:            true,
			},
			"password": schema.StringAttribute{
				MarkdownDescription: "Password for the Crosswork instance. This can also be set as the CW_PASSWORD environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "URL of the Crosswork instance. This can also be set as the CW_URL environment variable. If it is not set, the URL of the first entry of `instances` is used.",
				Optional:            true,
			},
			"insecure": schema.BoolAttribute{
				MarkdownDescription: "Allow insecure HTTPS client. This can also be set as the CW_INSECURE environment variable. Defaults to `true`.",
				Optional:            true,
			},
			"retries": schema.Int64Attribute{
				MarkdownDescription: "Number of retries for REST API calls. This can also be set as the CW_RETRIES environment variable. Defaults to `3`.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 9),
				},
			},
			"max_timeout": schema.Int64Attribute{
				MarkdownDescription: "Longest time in seconds to wait for the NSO plan of a service after it was written (the Services and Traffic Engineering resources), unless the resource sets its own. This can also be set as the CW_MAX_TIMEOUT environment variable. Defaults to `300`.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 3600),
				},
			},
			"device_ready_mode": schema.StringAttribute{
				MarkdownDescription: "What a device needs in Crosswork before a resource that names it is created or updated (the Services and Traffic Engineering resources and the Golden Configuration jobs): `none` makes no wait, `synced` waits until the device is synced with NSO, `full` also waits until it is reachable and admin up. With `none` NSO refuses a reference to a device that it does not know yet. This can also be set as the CW_DEVICE_READY_MODE environment variable. Defaults to `synced`.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.OneOf(nsoModeNone, nsoModeSynced, nsoModeFull),
				},
			},
			"device_ready_timeout": schema.Int64Attribute{
				MarkdownDescription: "Longest time in seconds to wait for the devices that a resource names to be ready as `device_ready_mode` says; 0 checks once. This can also be set as the CW_DEVICE_READY_TIMEOUT environment variable. Defaults to `300`.",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.Between(0, 86400),
				},
			},
			"instances": schema.ListNestedAttribute{
				MarkdownDescription: "This can be used to manage a list of Crosswork instances from a single provider. All instances must use the same credentials. Each resource and data source has an optional attribute named `instance`, which selects an instance by its name from this list. Without `instance`, the default instance is used: the cluster of `url`.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Instance name.",
							Required:            true,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "URL of the Crosswork instance.",
							Required:            true,
						},
					},
				},
			},
			"allow_existing_on_create": schema.BoolAttribute{
				MarkdownDescription: "**Experimental (use at your own risk).** Allow existing objects in Crosswork to be managed. If `true`, a POST request that attempts to create an already existing resource will trigger an update instead of failing. This behavior is not guaranteed and may not work for all resources. This can also be set as the CW_ALLOW_EXISTING_ON_CREATE environment variable. Defaults to `false`.",
				Optional:            true,
			},
		},
	}
}

func (p *CwProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config CwProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var username string
	if config.Username.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as username",
		)
		return
	}

	if config.Username.IsNull() {
		username = os.Getenv("CW_USERNAME")
	} else {
		username = config.Username.ValueString()
	}

	if username == "" {
		resp.Diagnostics.AddError(
			"Unable to find username",
			"Username cannot be an empty string",
		)
		return
	}

	var password string
	if config.Password.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as password",
		)
		return
	}

	if config.Password.IsNull() {
		password = os.Getenv("CW_PASSWORD")
	} else {
		password = config.Password.ValueString()
	}

	if password == "" {
		resp.Diagnostics.AddError(
			"Unable to find password",
			"Password cannot be an empty string",
		)
		return
	}

	var url string
	if config.URL.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as url",
		)
		return
	}

	if config.URL.IsNull() {
		url = os.Getenv("CW_URL")
		if url == "" && len(config.Instances) > 0 {
			url = config.Instances[0].URL.ValueString()
		}
	} else {
		url = config.URL.ValueString()
	}

	if url == "" {
		resp.Diagnostics.AddError(
			"Unable to find url",
			"URL cannot be an empty string",
		)
		return
	}

	var insecure bool
	if config.Insecure.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as insecure",
		)
		return
	}

	if config.Insecure.IsNull() {
		insecureStr := os.Getenv("CW_INSECURE")
		if insecureStr == "" {
			insecure = true
		} else {
			insecure, _ = strconv.ParseBool(insecureStr)
		}
	} else {
		insecure = config.Insecure.ValueBool()
	}

	var allow_existing_on_create bool
	if config.AllowExistingOnCreate.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as allow_existing_on_create",
		)
		return
	}

	if config.AllowExistingOnCreate.IsNull() {
		allow_existing_on_createStr := os.Getenv("CW_ALLOW_EXISTING_ON_CREATE")
		if allow_existing_on_createStr == "" {
			allow_existing_on_create = false
		} else {
			allow_existing_on_create, _ = strconv.ParseBool(allow_existing_on_createStr)
		}
	} else {
		allow_existing_on_create = config.AllowExistingOnCreate.ValueBool()
	}

	var retries int64
	if config.Retries.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as retries",
		)
		return
	}

	if config.Retries.IsNull() {
		retriesStr := os.Getenv("CW_RETRIES")
		if retriesStr == "" {
			retries = 3
		} else {
			retries, _ = strconv.ParseInt(retriesStr, 0, 64)
		}
	} else {
		retries = config.Retries.ValueInt64()
	}

	var maxTimeout int64
	if config.MaxTimeout.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as max_timeout",
		)
		return
	}

	if config.MaxTimeout.IsNull() {
		maxTimeoutStr := os.Getenv("CW_MAX_TIMEOUT")
		if maxTimeoutStr == "" {
			maxTimeout = 300
		} else {
			maxTimeout, _ = strconv.ParseInt(maxTimeoutStr, 0, 64)
		}
	} else {
		maxTimeout = config.MaxTimeout.ValueInt64()
	}

	if config.DeviceReadyMode.IsUnknown() || config.DeviceReadyTimeout.IsUnknown() {
		resp.Diagnostics.AddWarning(
			"Unable to create client",
			"Cannot use unknown value as device_ready_mode or device_ready_timeout",
		)
		return
	}
	deviceReady := nsoDeviceConfig{Mode: nsoModeSynced, Timeout: nsoDefaultTimeout, Store: newNsoSnapshots()}
	if !config.DeviceReadyMode.IsNull() {
		deviceReady.Mode = config.DeviceReadyMode.ValueString()
	} else if mode := os.Getenv("CW_DEVICE_READY_MODE"); mode != "" {
		deviceReady.Mode = mode
	}
	if deviceReady.Mode != nsoModeNone && deviceReady.Mode != nsoModeSynced && deviceReady.Mode != nsoModeFull {
		resp.Diagnostics.AddError(
			"Invalid device_ready_mode",
			fmt.Sprintf("device_ready_mode (or CW_DEVICE_READY_MODE) must be %s, %s or %s, got %q", nsoModeNone, nsoModeSynced, nsoModeFull, deviceReady.Mode),
		)
		return
	}
	if !config.DeviceReadyTimeout.IsNull() {
		deviceReady.Timeout = time.Duration(config.DeviceReadyTimeout.ValueInt64()) * time.Second
	} else if s := os.Getenv("CW_DEVICE_READY_TIMEOUT"); s != "" {
		seconds, err := strconv.ParseInt(s, 0, 64)
		if err != nil || seconds < 0 || seconds > 86400 {
			resp.Diagnostics.AddError(
				"Invalid device_ready_timeout",
				fmt.Sprintf("device_ready_timeout (or CW_DEVICE_READY_TIMEOUT) must be a number of seconds from 0 to 86400, got %q", s),
			)
			return
		}
		deviceReady.Timeout = time.Duration(seconds) * time.Second
	}

	// The provider needs a full URL: the client adds no scheme
	instanceURLs := []string{url}
	seen := map[string]bool{}
	for _, instance := range config.Instances {
		if instance.Name.IsUnknown() || instance.URL.IsUnknown() {
			resp.Diagnostics.AddWarning(
				"Unable to create client",
				"Cannot use unknown value as instance name or url",
			)
			return
		}
		name := instance.Name.ValueString()
		if name == "" || seen[name] {
			resp.Diagnostics.AddError(
				"Invalid instances",
				"Every instance needs a unique, non-empty name. Offending name: \""+name+"\"",
			)
			return
		}
		seen[name] = true
		instanceURLs = append(instanceURLs, instance.URL.ValueString())
	}
	for _, u := range instanceURLs {
		if !strings.HasPrefix(u, "https://") && !strings.HasPrefix(u, "http://") {
			resp.Diagnostics.AddError(
				"Invalid url",
				"The url must start with https:// (or http://), for example https://10.1.1.1. Got: "+u,
			)
			return
		}
	}

	// One client (one Crosswork session) per cluster and process; it is logged out by Shutdown.
	// The default client and an instance with the same URL share one registered client.
	clients := make(map[string]*cw.Client, len(config.Instances)+1)
	for i, u := range instanceURLs {
		name := ""
		if i > 0 {
			name = config.Instances[i-1].Name.ValueString()
		}
		c, err := registeredClientWithTimeout(u, username, password, insecure, int(retries), int(maxTimeout), p.version)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to create client",
				"Unable to create Crosswork client:\n\n"+err.Error(),
			)
			return
		}
		// Log in now so bad credentials or a session limit surface here.
		if err := c.Authenticate(); err != nil {
			resp.Diagnostics.AddError("Unable to log in to Crosswork", loginErrorDetail(err))
			return
		}
		clients[name] = c
	}

	data := CwProviderData{Clients: clients, AllowExistingOnCreate: allow_existing_on_create, Cache: NewThreadSafeCache(), DeviceReady: deviceReady}
	resp.DataSourceData = &data
	resp.ResourceData = &data
}

func (p *CwProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewAAALDAPServersResource,
		NewAAAPasswordPolicyResource,
		NewAAARADIUSServersResource,
		NewAAASessionSettingsResource,
		NewAAASSOServicesResource,
		NewAAATACACSServersResource,
		NewAlarmSettingsResource,
		NewAssignDevicesToAccessGroupResource,
		NewAssignDevicesToDataGatewayResource,
		NewAssignDevicesToDeviceGroupResource,
		NewAssignPortsToPortGroupResource,
		NewBackupDestinationResource,
		NewCredentialProfilesResource,
		NewDataGatewayDestinationsResource,
		NewDataGatewayGlobalParametersResource,
		NewDataGatewayPoolsResource,
		NewDataGatewayResourceConfigurationResource,
		NewDataRetentionDevicePerformanceResource,
		NewDataRetentionNetworkPerformanceResource,
		NewDeviceAccessGroupsResource,
		NewDeviceGroupsResource,
		NewDeviceNotificationPoliciesResource,
		NewFileServersResource,
		NewGNMISettingsResource,
		NewHealthSettingsResource,
		NewL2VPNEtagSetsResource,
		NewL2VPNEVPNRouteTypeSetsResource,
		NewL2VPNIDPoolsResource,
		NewL2VPNQoSProfileIdsResource,
		NewL2VPNRDSetsResource,
		NewL2VPNRoutePoliciesResource,
		NewL2VPNTagSetsResource,
		NewL2VPNVPNServicesResource,
		NewL2VPNY1731ProfilesResource,
		NewL3VPNDestinationPrefixSetsResource,
		NewL3VPNQoSProfileIdsResource,
		NewL3VPNRoutePoliciesResource,
		NewL3VPNSourcePrefixSetsResource,
		NewL3VPNTagSetsResource,
		NewL3VPNVPNServicesResource,
		NewLockUsersResource,
		NewLoginBannerResource,
		NewMetricThresholdsResource,
		NewNetworkDevicesResource,
		NewPerformancePoliciesResource,
		NewPortGroupsResource,
		NewPortGroupsRulesResource,
		NewProviderNSOPoliciesResource,
		NewProviderProfilesResource,
		NewResourcePoolsResource,
		NewSmartLicensingTransportResource,
		NewSRTEODNTemplatesResource,
		NewSRTEPoliciesResource,
		NewSRTESIDListsResource,
		NewStaticRouteResource,
		NewSuppressionPoliciesResource,
		NewSurveySettingsResource,
		NewSyslogDestinationResource,
		NewSystemNotificationPoliciesResource,
		NewTagsResource,
		NewTEFlexAlgoAffinityResource,
		NewTEGeneralSettingsResource,
		NewTELinkAffinityResource,
		NewTopologyDiscoverySettingsResource,
		NewTrapDestinationResource,
		NewUsersResource,
		NewRolesResource,
		NewGoldenConfigGlobalVariablesResource,
		NewGoldenConfigTemplatesResource,
		NewGoldenConfigJobsResource,
		NewCollectionJobsResource,
		NewCollectionJobsCadenceResource,
	}
}

func (p *CwProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAAALDAPServersDataSource,
		NewAAAPasswordPolicyDataSource,
		NewAAARADIUSServersDataSource,
		NewAAASessionSettingsDataSource,
		NewAAASSOServicesDataSource,
		NewAAATACACSServersDataSource,
		NewAlarmSettingsDataSource,
		NewBackupDestinationDataSource,
		NewCredentialProfilesDataSource,
		NewDataGatewayGlobalParametersDataSource,
		NewDataGatewayResourceConfigurationDataSource,
		NewDataRetentionDevicePerformanceDataSource,
		NewDataRetentionNetworkPerformanceDataSource,
		NewDeviceAccessGroupsDataSource,
		NewDeviceGroupsDataSource,
		NewDeviceNotificationPoliciesDataSource,
		NewFileServersDataSource,
		NewGNMISettingsDataSource,
		NewHealthSettingsDataSource,
		NewL2VPNEtagSetsDataSource,
		NewL2VPNEVPNRouteTypeSetsDataSource,
		NewL2VPNIDPoolsDataSource,
		NewL2VPNQoSProfileIdsDataSource,
		NewL2VPNRDSetsDataSource,
		NewL2VPNRoutePoliciesDataSource,
		NewL2VPNTagSetsDataSource,
		NewL2VPNVPNServicesDataSource,
		NewL2VPNY1731ProfilesDataSource,
		NewL3VPNDestinationPrefixSetsDataSource,
		NewL3VPNQoSProfileIdsDataSource,
		NewL3VPNRoutePoliciesDataSource,
		NewL3VPNSourcePrefixSetsDataSource,
		NewL3VPNTagSetsDataSource,
		NewL3VPNVPNServicesDataSource,
		NewLockUsersDataSource,
		NewLoginBannerDataSource,
		NewMetricThresholdsDataSource,
		NewNetworkDevicesDataSource,
		NewPerformancePoliciesDataSource,
		NewPortGroupsDataSource,
		NewProviderNSOPoliciesDataSource,
		NewProviderProfilesDataSource,
		NewResourcePoolsDataSource,
		NewSmartLicensingTransportDataSource,
		NewSRTEODNTemplatesDataSource,
		NewSRTEPoliciesDataSource,
		NewSRTESIDListsDataSource,
		NewStaticRouteDataSource,
		NewSuppressionPoliciesDataSource,
		NewSurveySettingsDataSource,
		NewSyslogDestinationDataSource,
		NewSystemNotificationPoliciesDataSource,
		NewTagsDataSource,
		NewTEFlexAlgoAffinityDataSource,
		NewTEGeneralSettingsDataSource,
		NewTELinkAffinityDataSource,
		NewTopologyDiscoverySettingsDataSource,
		NewTrapDestinationDataSource,
		NewUsersDataSource,
		NewRolesDataSource,
		NewRoleCatalogDataSource,
		NewGroupingsDataSource,
		NewGoldenConfigGlobalVariablesDataSource,
		NewGoldenConfigTemplatesDataSource,
		NewGoldenConfigJobsDataSource,
		NewDataGatewayPoolsDataSource,
		NewDataGatewayInstancesDataSource,
		NewDataGatewayDestinationsDataSource,
		NewCollectionJobsDataSource,
		NewCollectionJobsCadenceDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CwProvider{
			version: version,
		}
	}
}

type ThreadSafeCache struct {
	mu    sync.RWMutex
	items map[string]interface{}
}

func NewThreadSafeCache() *ThreadSafeCache {
	return &ThreadSafeCache{
		items: make(map[string]interface{}),
	}
}

func (c *ThreadSafeCache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

func (c *ThreadSafeCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	return item, found
}

func (c *ThreadSafeCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *ThreadSafeCache) DeletePattern(pattern string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key := range c.items {
		if strings.Contains(key, pattern) {
			delete(c.items, key)
		}
	}
}

// End of section. //template:end provider

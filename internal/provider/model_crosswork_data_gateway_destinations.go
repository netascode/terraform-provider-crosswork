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
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type DataGatewayDestinations struct {
	Instance          types.String                         `tfsdk:"instance"`
	Id                types.String                         `tfsdk:"id"`
	Name              types.String                         `tfsdk:"name"`
	DataSource        types.String                         `tfsdk:"data_source"`
	ServerType        types.String                         `tfsdk:"server_type"`
	Encoding          types.String                         `tfsdk:"encoding"`
	CompressionMethod types.String                         `tfsdk:"compression_method"`
	DispatchType      types.String                         `tfsdk:"dispatch_type"`
	MaxMessageSize    types.Int64                          `tfsdk:"max_message_size"`
	BufferMemory      types.Int64                          `tfsdk:"buffer_memory"`
	BatchSize         types.Int64                          `tfsdk:"batch_size"`
	LingerTime        types.Int64                          `tfsdk:"linger_time"`
	RequestTimeout    types.Int64                          `tfsdk:"request_timeout"`
	Collectors        []DataGatewayDestinationsCollectors  `tfsdk:"collectors"`
	Connections       []DataGatewayDestinationsConnections `tfsdk:"connections"`
	Secure            types.Bool                           `tfsdk:"secure"`
	AuthType          types.String                         `tfsdk:"auth_type"`
	SystemDefined     types.Bool                           `tfsdk:"system_defined"`
}

type DataGatewayDestinationsCollectors struct {
	Name           types.String `tfsdk:"name"`
	BufferMemory   types.Int64  `tfsdk:"buffer_memory"`
	BatchSize      types.Int64  `tfsdk:"batch_size"`
	LingerTime     types.Int64  `tfsdk:"linger_time"`
	RequestTimeout types.Int64  `tfsdk:"request_timeout"`
}

type DataGatewayDestinationsConnections struct {
	Address      types.String `tfsdk:"address"`
	PrefixLength types.Int64  `tfsdk:"prefix_length"`
	HostName     types.String `tfsdk:"host_name"`
	DomainName   types.String `tfsdk:"domain_name"`
	Port         types.Int64  `tfsdk:"port"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data DataGatewayDestinations) getPath() string {
	return "/crosswork/dg-manager/v1/destinations"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data DataGatewayDestinations) getPathGet() string {
	return "/crosswork/dg-manager/v1/destinations/query"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// Escapes the id so it cannot change the structure of a request body.
func (data DataGatewayDestinations) idJSON() string {
	id := data.Id.ValueString()
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}

func (data DataGatewayDestinations) readGet(c *cw.Client, path string) (cw.Res, error) {
	return c.Post(path, strings.ReplaceAll(`{"limit":1000,"filter":{}}`, "{id}", data.idJSON()), cw.ReadOnly)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Written by hand: the properties are a flat map of text values and a connection carries its mask as a string.
func (data DataGatewayDestinations) toBody(ctx context.Context, state DataGatewayDestinations) string {
	return destinationRequestBody(data, state.Id.ValueString())
}

// End of section. //template:end toBody

// Written by hand: the numbers come back as text properties, the collector values as keys like CLI.LINGER_MS_CONFIG, and the
// data source and the server type as the wire values of the controller.
func (data *DataGatewayDestinations) fromBody(ctx context.Context, res gjson.Result) {
	props := res.Get("properties").Map()
	if value := res.Get("uuid"); value.Exists() {
		data.Id = types.StringValue(value.String())
	}
	data.Name = types.StringValue(res.Get("name").String())
	data.DataSource = types.StringValue(destinationKey(destinationDataSources, props["DISPATCH_SOURCE"].String()))
	data.ServerType = types.StringValue(destinationKey(destinationServerTypes, props["DESTINATION_TYPE"].String()))
	data.Encoding = destinationStringProperty(props, "ENCODING")
	data.CompressionMethod = destinationStringProperty(props, "COMPRESSION_TYPE_CONFIG")
	data.DispatchType = destinationStringProperty(props, "GRPC_DISPATCH_TYPE")
	data.MaxMessageSize = destinationIntProperty(props, "MAX_REQUEST_SIZE_CONFIG")
	data.BufferMemory = destinationIntProperty(props, "BUFFER_MEMORY_CONFIG")
	data.BatchSize = destinationIntProperty(props, "BATCH_SIZE_CONFIG")
	data.LingerTime = destinationIntProperty(props, "LINGER_MS_CONFIG")
	data.RequestTimeout = destinationIntProperty(props, "REQUEST_TIMEOUT_CONFIG")
	data.Collectors = destinationCollectorsFrom(props, nil)
	data.Connections = destinationConnections(res)
	secure := props["IS_SECURITY_ENABLED"].String() == "true"
	data.Secure = types.BoolValue(secure)
	data.AuthType = types.StringNull()
	if secure {
		data.AuthType = destinationStringProperty(props, "AUTH_TYPE")
	}
	data.SystemDefined = types.BoolValue(props["IS_SYSTEM_DEFINED"].String() == "true")
}

// End of section. //template:end fromBody

// Written by hand: only the optional values that the configuration sets are refreshed, so the values that belong to the other
// server type stay null; the collectors keep the order of the configuration and the connections are the whole list.
func (data *DataGatewayDestinations) updateFromBody(ctx context.Context, res gjson.Result) {
	props := res.Get("properties").Map()
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	}
	data.DataSource = types.StringValue(destinationKey(destinationDataSources, props["DISPATCH_SOURCE"].String()))
	data.ServerType = types.StringValue(destinationKey(destinationServerTypes, props["DESTINATION_TYPE"].String()))
	data.CompressionMethod = destinationStringProperty(props, "COMPRESSION_TYPE_CONFIG")
	if !data.Encoding.IsNull() {
		data.Encoding = destinationStringProperty(props, "ENCODING")
	}
	if !data.DispatchType.IsNull() {
		data.DispatchType = destinationStringProperty(props, "GRPC_DISPATCH_TYPE")
	}
	if !data.MaxMessageSize.IsNull() {
		data.MaxMessageSize = destinationIntProperty(props, "MAX_REQUEST_SIZE_CONFIG")
	}
	if !data.BufferMemory.IsNull() {
		data.BufferMemory = destinationIntProperty(props, "BUFFER_MEMORY_CONFIG")
	}
	if !data.BatchSize.IsNull() {
		data.BatchSize = destinationIntProperty(props, "BATCH_SIZE_CONFIG")
	}
	if !data.LingerTime.IsNull() {
		data.LingerTime = destinationIntProperty(props, "LINGER_MS_CONFIG")
	}
	if !data.RequestTimeout.IsNull() {
		data.RequestTimeout = destinationIntProperty(props, "REQUEST_TIMEOUT_CONFIG")
	}
	if len(data.Collectors) > 0 {
		data.Collectors = destinationCollectorsFrom(props, data.Collectors)
	}
	data.Connections = destinationConnections(res)
	if !data.Secure.IsNull() {
		data.Secure = types.BoolValue(props["IS_SECURITY_ENABLED"].String() == "true")
	}
	if !data.AuthType.IsNull() {
		data.AuthType = destinationStringProperty(props, "AUTH_TYPE")
	}
	data.SystemDefined = types.BoolValue(props["IS_SYSTEM_DEFINED"].String() == "true")
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *DataGatewayDestinations) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Name.IsNull() {
		return false
	}
	if !data.DataSource.IsNull() {
		return false
	}
	if !data.ServerType.IsNull() {
		return false
	}
	if !data.Encoding.IsNull() {
		return false
	}
	if !data.CompressionMethod.IsNull() {
		return false
	}
	if !data.DispatchType.IsNull() {
		return false
	}
	if !data.MaxMessageSize.IsNull() {
		return false
	}
	if !data.BufferMemory.IsNull() {
		return false
	}
	if !data.BatchSize.IsNull() {
		return false
	}
	if !data.LingerTime.IsNull() {
		return false
	}
	if !data.RequestTimeout.IsNull() {
		return false
	}
	if len(data.Collectors) > 0 {
		return false
	}
	if len(data.Connections) > 0 {
		return false
	}
	if !data.Secure.IsNull() {
		return false
	}
	if !data.AuthType.IsNull() {
		return false
	}
	if !data.SystemDefined.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

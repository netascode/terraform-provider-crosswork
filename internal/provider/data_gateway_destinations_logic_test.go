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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

func testKafkaDestination() DataGatewayDestinations {
	return DataGatewayDestinations{
		Name:              types.StringValue("cwtest-kafka"),
		DataSource:        types.StringValue("data-gateway"),
		ServerType:        types.StringValue("kafka"),
		Encoding:          types.StringValue("gpbkv"),
		CompressionMethod: types.StringValue("snappy"),
		MaxMessageSize:    types.Int64Value(100000000),
		BufferMemory:      types.Int64Value(52428800),
		BatchSize:         types.Int64Value(1048576),
		LingerTime:        types.Int64Value(2000),
		RequestTimeout:    types.Int64Value(30),
		Collectors: []DataGatewayDestinationsCollectors{{
			Name:           types.StringValue("CLI"),
			BufferMemory:   types.Int64Null(),
			BatchSize:      types.Int64Value(500000),
			LingerTime:     types.Int64Value(1000),
			RequestTimeout: types.Int64Null(),
		}},
		Connections: []DataGatewayDestinationsConnections{{
			Address:      types.StringValue("192.0.2.150"),
			PrefixLength: types.Int64Value(20),
			HostName:     types.StringNull(),
			DomainName:   types.StringNull(),
			Port:         types.Int64Value(9192),
		}},
		Secure:   types.BoolValue(true),
		AuthType: types.StringValue("mutual"),
	}
}

func TestDestinationRequestBodyKafka(t *testing.T) {
	body := gjson.Parse(destinationRequestBody(testKafkaDestination(), ""))
	dest := body.Get("providers.0")
	assert.Equal(t, "cwtest-kafka", dest.Get("name").String())
	assert.Equal(t, "ROBOT_PROVIDER_DESTINATION", dest.Get("family").String())
	assert.False(t, dest.Get("uuid").Exists(), "a create has no id")
	props := dest.Get("properties").Map()
	assert.Equal(t, "destination_type_kafka", props["DESTINATION_TYPE"].String())
	assert.Equal(t, "datagateway", props["DISPATCH_SOURCE"].String())
	assert.Equal(t, "gpbkv", props["ENCODING"].String())
	assert.Equal(t, "snappy", props["COMPRESSION_TYPE_CONFIG"].String())
	// the controller wants every number as a string
	assert.Equal(t, gjson.String, props["MAX_REQUEST_SIZE_CONFIG"].Type)
	assert.Equal(t, "100000000", props["MAX_REQUEST_SIZE_CONFIG"].String())
	assert.Equal(t, "52428800", props["BUFFER_MEMORY_CONFIG"].String())
	assert.Equal(t, "1048576", props["BATCH_SIZE_CONFIG"].String())
	assert.Equal(t, "2000", props["LINGER_MS_CONFIG"].String())
	assert.Equal(t, "30", props["REQUEST_TIMEOUT_CONFIG"].String())
	// a collector override is a flat key with a dot, and a value that is not set is left out
	assert.Equal(t, "500000", props["CLI.BATCH_SIZE_CONFIG"].String())
	assert.Equal(t, "1000", props["CLI.LINGER_MS_CONFIG"].String())
	_, ok := props["CLI.BUFFER_MEMORY_CONFIG"]
	assert.False(t, ok, "an unset collector value is not sent")
	// the controller accepts true and then refuses every later change, so it is always false
	assert.Equal(t, "false", props["IS_SYSTEM_DEFINED"].String())
	assert.Equal(t, "true", props["IS_SECURITY_ENABLED"].String())
	assert.Equal(t, "mutual", props["AUTH_TYPE"].String())
	assert.False(t, props["GRPC_DISPATCH_TYPE"].Exists(), "no gRPC value on a Kafka destination")
	conn := dest.Get("connectivity_info.0")
	assert.Equal(t, "ROBOT_MSVC_TRANS_KAFKA", conn.Get("type").String())
	assert.Equal(t, "ROBOT_INET_ADDR_TYPE_v4", conn.Get("ipaddrs.0.inet_af").String())
	assert.Equal(t, "192.0.2.150", conn.Get("ipaddrs.0.inet_addr").String())
	assert.Equal(t, gjson.String, conn.Get("ipaddrs.0.mask").Type)
	assert.Equal(t, "20", conn.Get("ipaddrs.0.mask").String())
	assert.Equal(t, int64(9192), conn.Get("port").Int())
	assert.False(t, conn.Get("fqdn").Exists())
}

func TestDestinationRequestBodyUpdateCarriesTheId(t *testing.T) {
	assert.Equal(t, "abc-1", gjson.Get(destinationRequestBody(testKafkaDestination(), "abc-1"), "providers.0.uuid").String())
}

func TestDestinationRequestBodyGrpcByName(t *testing.T) {
	d := DataGatewayDestinations{
		Name:              types.StringValue("cwtest-grpc"),
		DataSource:        types.StringValue("application"),
		ServerType:        types.StringValue("grpc"),
		Encoding:          types.StringNull(),
		CompressionMethod: types.StringValue("deflate"),
		DispatchType:      types.StringValue("unary"),
		LingerTime:        types.Int64Value(1000),
		Connections: []DataGatewayDestinationsConnections{{
			Address: types.StringNull(), PrefixLength: types.Int64Null(),
			HostName: types.StringValue("collector"), DomainName: types.StringValue("example.com"), Port: types.Int64Value(50052),
		}},
		Secure:   types.BoolNull(),
		AuthType: types.StringNull(),
	}
	dest := gjson.Parse(destinationRequestBody(d, "")).Get("providers.0")
	props := dest.Get("properties").Map()
	assert.Equal(t, "destination_type_grpc", props["DESTINATION_TYPE"].String())
	assert.Equal(t, "application", props["DISPATCH_SOURCE"].String())
	assert.Equal(t, "unary", props["GRPC_DISPATCH_TYPE"].String())
	assert.Equal(t, "deflate", props["COMPRESSION_TYPE_CONFIG"].String())
	assert.Equal(t, "false", props["IS_SECURITY_ENABLED"].String())
	assert.Equal(t, "1000", props["LINGER_MS_CONFIG"].String())
	for _, k := range []string{"ENCODING", "AUTH_TYPE"} {
		_, ok := props[k]
		assert.False(t, ok, k+" is not sent when the model leaves it out")
	}
	conn := dest.Get("connectivity_info.0")
	assert.Equal(t, "ROBOT_MSVC_TRANS_GRPC", conn.Get("type").String())
	assert.Equal(t, "collector", conn.Get("fqdn.host_name").String())
	assert.Equal(t, "example.com", conn.Get("fqdn.domain_name").String())
	assert.False(t, conn.Get("ipaddrs").Exists())
}

func TestDestinationSendsTheValuesTheModelSetsWhateverTheServerType(t *testing.T) {
	d := testKafkaDestination()
	d.Secure = types.BoolValue(false)
	d.DispatchType = types.StringValue("stream")
	props := gjson.Parse(destinationRequestBody(d, "")).Get("providers.0.properties").Map()
	// the controller keeps both and reads them back, so the state matches the configuration
	assert.Equal(t, "mutual", props["AUTH_TYPE"].String())
	assert.Equal(t, "stream", props["GRPC_DISPATCH_TYPE"].String())
	assert.Equal(t, "false", props["IS_SECURITY_ENABLED"].String())
}

const kafkaRow = `{"uuid":"id-1","name":"cwtest-kafka","family":"ROBOT_PROVIDER_DESTINATION",
 "connectivity_info":[{"type":"ROBOT_MSVC_TRANS_KAFKA","port":9192,"ipaddrs":[{"inet_af":"ROBOT_INET_ADDR_TYPE_v4","inet_addr":"192.0.2.150","mask":"20"}]},
                      {"type":"ROBOT_MSVC_TRANS_KAFKA","port":9193,"fqdn":{"host_name":"kafka1","domain_name":"example.com"}}],
 "properties":{"DESTINATION_TYPE":"destination_type_kafka","DISPATCH_SOURCE":"datagateway","ENCODING":"json","COMPRESSION_TYPE_CONFIG":"zstd",
  "MAX_REQUEST_SIZE_CONFIG":"100000000","BUFFER_MEMORY_CONFIG":"52428800","BATCH_SIZE_CONFIG":"1048576","LINGER_MS_CONFIG":"2000","REQUEST_TIMEOUT_CONFIG":"30",
  "CLI.LINGER_MS_CONFIG":"100","SNMP.BATCH_SIZE_CONFIG":"20000","IS_SECURITY_ENABLED":"true","AUTH_TYPE":"server","IS_SYSTEM_DEFINED":"false"}}`

func TestDestinationFromBody(t *testing.T) {
	var d DataGatewayDestinations
	d.fromBody(context.Background(), gjson.Parse(kafkaRow))
	assert.Equal(t, "id-1", d.Id.ValueString())
	assert.Equal(t, "data-gateway", d.DataSource.ValueString())
	assert.Equal(t, "kafka", d.ServerType.ValueString())
	assert.Equal(t, "json", d.Encoding.ValueString())
	assert.Equal(t, "zstd", d.CompressionMethod.ValueString())
	assert.Equal(t, int64(100000000), d.MaxMessageSize.ValueInt64())
	assert.Equal(t, int64(30), d.RequestTimeout.ValueInt64())
	assert.True(t, d.DispatchType.IsNull(), "a Kafka destination has no dispatch type")
	assert.True(t, d.Secure.ValueBool())
	assert.Equal(t, "server", d.AuthType.ValueString())
	assert.False(t, d.SystemDefined.ValueBool())
	assert.Len(t, d.Collectors, 2)
	assert.Equal(t, "CLI", d.Collectors[0].Name.ValueString())
	assert.Equal(t, int64(100), d.Collectors[0].LingerTime.ValueInt64())
	assert.True(t, d.Collectors[0].BatchSize.IsNull())
	assert.Equal(t, "SNMP", d.Collectors[1].Name.ValueString())
	assert.Equal(t, int64(20000), d.Collectors[1].BatchSize.ValueInt64())
	assert.Len(t, d.Connections, 2)
	assert.Equal(t, "192.0.2.150", d.Connections[0].Address.ValueString())
	assert.Equal(t, int64(20), d.Connections[0].PrefixLength.ValueInt64())
	assert.True(t, d.Connections[0].HostName.IsNull())
	assert.Equal(t, "kafka1", d.Connections[1].HostName.ValueString())
	assert.Equal(t, "example.com", d.Connections[1].DomainName.ValueString())
	assert.True(t, d.Connections[1].Address.IsNull())
}

func TestDestinationFromBodyReadsASystemDefinedGrpcDestination(t *testing.T) {
	row := `{"uuid":"id-2","name":"cdg-common-pipeline","connectivity_info":[{"type":"ROBOT_MSVC_TRANS_GRPC","port":8005,"ipaddrs":[{"inet_af":"ROBOT_INET_ADDR_TYPE_v4","inet_addr":"127.0.0.1","mask":"8"}]}],
	 "properties":{"DESTINATION_TYPE":"destination_type_grpc","DISPATCH_SOURCE":"datagateway","ENCODING":"gpbkv","COMPRESSION_TYPE_CONFIG":"gzip","GRPC_DISPATCH_TYPE":"stream","SERVICE_NAME":"common-pipeline","TYPE":"INTERNAL","IS_SYSTEM_DEFINED":"true"}}`
	var d DataGatewayDestinations
	d.fromBody(context.Background(), gjson.Parse(row))
	assert.Equal(t, "grpc", d.ServerType.ValueString())
	assert.Equal(t, "stream", d.DispatchType.ValueString())
	assert.True(t, d.SystemDefined.ValueBool())
	assert.False(t, d.Secure.ValueBool())
	assert.True(t, d.AuthType.IsNull())
	assert.True(t, d.MaxMessageSize.IsNull())
	assert.Len(t, d.Collectors, 0)
}

func TestDestinationUpdateFromBodyRefreshesOnlyWhatIsSet(t *testing.T) {
	state := DataGatewayDestinations{
		Name: types.StringValue("cwtest-kafka"), Encoding: types.StringNull(), LingerTime: types.Int64Value(1),
		DispatchType: types.StringNull(), MaxMessageSize: types.Int64Null(), BufferMemory: types.Int64Null(),
		BatchSize: types.Int64Null(), RequestTimeout: types.Int64Null(), Secure: types.BoolNull(), AuthType: types.StringNull(),
		// the order of the configuration is kept, and a collector the controller lost is dropped
		Collectors: []DataGatewayDestinationsCollectors{
			{Name: types.StringValue("SNMP"), BatchSize: types.Int64Value(1)},
			{Name: types.StringValue("MDT"), BatchSize: types.Int64Value(1)},
		},
	}
	state.updateFromBody(context.Background(), gjson.Parse(kafkaRow))
	assert.Equal(t, int64(2000), state.LingerTime.ValueInt64())
	assert.True(t, state.Encoding.IsNull(), "an encoding that the configuration leaves out stays null")
	assert.True(t, state.Secure.IsNull(), "a secure flag that the configuration leaves out stays null")
	assert.True(t, state.BatchSize.IsNull())
	assert.Len(t, state.Collectors, 2)
	assert.Equal(t, "SNMP", state.Collectors[0].Name.ValueString())
	assert.Equal(t, "CLI", state.Collectors[1].Name.ValueString())
	assert.Len(t, state.Connections, 2)
}

func TestDestinationFindAndKey(t *testing.T) {
	res := gjson.Parse(`{"data":[` + kafkaRow + `],"total_count":1}`)
	assert.Equal(t, "id-1", findDestination(res, "", "cwtest-kafka").Get("uuid").String())
	assert.Equal(t, "cwtest-kafka", findDestination(res, "id-1", "").Get("name").String())
	assert.False(t, findDestination(res, "", "").Exists(), "an empty id and an empty name match nothing")
	assert.False(t, findDestination(res, "other", "").Exists())
	assert.Equal(t, "data-gateway", destinationKey(destinationDataSources, "datagateway"))
	assert.Equal(t, "any", destinationKey(destinationDataSources, "any"))
}

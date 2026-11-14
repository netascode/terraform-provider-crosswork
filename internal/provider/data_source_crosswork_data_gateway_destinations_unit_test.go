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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func seedDestinations(s *destinationStub) {
	s.seed(`{"name":"CW_KAFKA_DESTINATION","family":"ROBOT_PROVIDER_DESTINATION",
	  "connectivity_info":[{"type":"ROBOT_MSVC_TRANS_KAFKA","port":30993,"ipaddrs":[{"inet_af":"ROBOT_INET_ADDR_TYPE_v4","inet_addr":"192.0.2.189","mask":"24"}]}],
	  "properties":{"DESTINATION_TYPE":"destination_type_kafka","DISPATCH_SOURCE":"datagateway","ENCODING":"gpbkv","COMPRESSION_TYPE_CONFIG":"snappy",
	   "MAX_REQUEST_SIZE_CONFIG":"100000000","BUFFER_MEMORY_CONFIG":"314572800","BATCH_SIZE_CONFIG":"1048576","LINGER_MS_CONFIG":"1000","REQUEST_TIMEOUT_CONFIG":"30",
	   "IS_SECURITY_ENABLED":"true","AUTH_TYPE":"mutual","IS_SYSTEM_DEFINED":"true"}}`)
	s.seed(`{"name":"cdg-common-pipeline","family":"ROBOT_PROVIDER_DESTINATION",
	  "connectivity_info":[{"type":"ROBOT_MSVC_TRANS_GRPC","port":8005,"ipaddrs":[{"inet_af":"ROBOT_INET_ADDR_TYPE_v4","inet_addr":"127.0.0.1","mask":"8"}]}],
	  "properties":{"DESTINATION_TYPE":"destination_type_grpc","DISPATCH_SOURCE":"datagateway","ENCODING":"gpbkv","COMPRESSION_TYPE_CONFIG":"gzip",
	   "GRPC_DISPATCH_TYPE":"stream","SERVICE_NAME":"common-pipeline","TYPE":"INTERNAL","IS_SYSTEM_DEFINED":"true"}}`)
}

func TestDataGatewayDestinationsDataSource(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newDestinationStub(t)
	seedDestinations(s)
	const ds = "data.crosswork_data_gateway_destinations.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: s.providerBlock() + "\ndata \"crosswork_data_gateway_destinations\" \"test\" {\n  name = \"CW_KAFKA_DESTINATION\"\n}\n",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ds, "id", "dest-id-1"),
					resource.TestCheckResourceAttr(ds, "system_defined", "true"),
					resource.TestCheckResourceAttr(ds, "data_source", "data-gateway"),
					resource.TestCheckResourceAttr(ds, "server_type", "kafka"),
					resource.TestCheckResourceAttr(ds, "encoding", "gpbkv"),
					resource.TestCheckResourceAttr(ds, "compression_method", "snappy"),
					resource.TestCheckResourceAttr(ds, "buffer_memory", "314572800"),
					resource.TestCheckResourceAttr(ds, "connections.0.address", "192.0.2.189"),
					resource.TestCheckResourceAttr(ds, "connections.0.port", "30993"),
					resource.TestCheckResourceAttr(ds, "secure", "true"),
					resource.TestCheckResourceAttr(ds, "auth_type", "mutual"),
					resource.TestCheckResourceAttr(ds, "collectors.#", "0"),
				),
			},
			{
				Config: s.providerBlock() + "\ndata \"crosswork_data_gateway_destinations\" \"test\" {\n  id = \"dest-id-2\"\n}\n",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ds, "name", "cdg-common-pipeline"),
					resource.TestCheckResourceAttr(ds, "server_type", "grpc"),
					resource.TestCheckResourceAttr(ds, "dispatch_type", "stream"),
					resource.TestCheckResourceAttr(ds, "secure", "false"),
				),
			},
			{
				Config:      s.providerBlock() + "\ndata \"crosswork_data_gateway_destinations\" \"test\" {\n  name = \"no-such-destination\"\n}\n",
				ExpectError: regexp.MustCompile(`(?s)There is no destination`),
			},
			{
				Config:      s.providerBlock() + "\ndata \"crosswork_data_gateway_destinations\" \"test\" {\n  id = \"dest-id-1\"\n  name = \"CW_KAFKA_DESTINATION\"\n}\n",
				ExpectError: regexp.MustCompile(`(?is)(only one|exactly one|Invalid Attribute Combination)`),
			},
			{
				Config:      s.providerBlock() + "\ndata \"crosswork_data_gateway_destinations\" \"test\" {\n}\n",
				ExpectError: regexp.MustCompile(`(?is)(exactly one|No attribute specified|Invalid Attribute Combination)`),
			},
			{
				Config: s.providerBlock() + "\ndata \"crosswork_data_gateway_destinations\" \"test\" {\n  id = \"dest-id-1\"\n}\n",
				Check:  resource.TestCheckResourceAttr(ds, "name", "CW_KAFKA_DESTINATION"),
			},
		},
	})
}

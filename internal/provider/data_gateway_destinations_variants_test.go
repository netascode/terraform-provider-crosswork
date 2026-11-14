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

// Covers the variants the generated acceptance test leaves out: a host name with a domain, secure communication with server
// authentication, the data source `application` without an encoding, and a gRPC destination.

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const dataGatewayDestinationsVariantsConfig = `
resource "crosswork_data_gateway_destinations" "kafka_name" {
  name               = "cwtest-gen-destination-kafka-name"
  data_source        = "data-gateway"
  server_type        = "kafka"
  encoding           = "json"
  compression_method = "zstd"
  max_message_size   = 50000000
  buffer_memory      = 104857600
  batch_size         = 524288
  linger_time        = 1000
  request_timeout    = 45
  collectors = [
    { name = "GNMI", linger_time = 3000 },
    { name = "SYSLOG", request_timeout = 40 },
  ]
  connections = [
    { host_name = "kafka1", domain_name = "example.com", port = 9093 },
  ]
  secure    = true
  auth_type = "server"
}

resource "crosswork_data_gateway_destinations" "kafka_app" {
  name               = "cwtest-gen-destination-kafka-app"
  data_source        = "application"
  server_type        = "kafka"
  compression_method = "none"
  max_message_size   = 100000000
  buffer_memory      = 52428800
  batch_size         = 1048576
  linger_time        = 2000
  request_timeout    = 30
  connections = [
    { address = "10.14.0.151", prefix_length = 24, port = 9192 },
    { address = "10.14.0.152", prefix_length = 24, port = 9193 },
  ]
}

resource "crosswork_data_gateway_destinations" "grpc_app" {
  name               = "cwtest-gen-destination-grpc-app"
  data_source        = "application"
  server_type        = "grpc"
  compression_method = "deflate"
  dispatch_type      = "unary"
  connections = [
    { address = "10.14.0.153", prefix_length = 24, port = 50053 },
  ]
}
`

func TestAccCwDataGatewayDestinationsVariants(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: dataGatewayDestinationsVariantsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.kafka_name", "connections.0.host_name", "kafka1"),
					resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.kafka_name", "connections.0.domain_name", "example.com"),
					resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.kafka_name", "collectors.1.name", "SYSLOG"),
					resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.kafka_name", "auth_type", "server"),
					resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.kafka_app", "data_source", "application"),
					resource.TestCheckNoResourceAttr("crosswork_data_gateway_destinations.kafka_app", "encoding"),
					resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.kafka_app", "connections.#", "2"),
					resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.grpc_app", "dispatch_type", "unary"),
					resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.grpc_app", "compression_method", "deflate"),
				),
			},
			{
				ResourceName:      "crosswork_data_gateway_destinations.kafka_name",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

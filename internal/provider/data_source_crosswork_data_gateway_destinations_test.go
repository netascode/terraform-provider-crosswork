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

// Needs a Crosswork with Crosswork Data Gateway.

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const dataGatewayDestinationsLookupConfig = `
resource "crosswork_data_gateway_destinations" "lookup" {
  name               = "cwtest-gen-destination-lookup"
  data_source        = "data-gateway"
  server_type        = "grpc"
  encoding           = "gpbkv"
  compression_method = "gzip"
  dispatch_type      = "stream"
  connections = [{
    address       = "10.14.0.150"
    prefix_length = 20
    port          = 50052
  }]
}

data "crosswork_data_gateway_destinations" "by_name" {
  name = crosswork_data_gateway_destinations.lookup.name
}

data "crosswork_data_gateway_destinations" "by_id" {
  id = crosswork_data_gateway_destinations.lookup.id
}
`

func TestAccDataSourceCwDataGatewayDestinations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: dataGatewayDestinationsLookupConfig,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrPair("data.crosswork_data_gateway_destinations.by_name", "id", "crosswork_data_gateway_destinations.lookup", "id"),
				resource.TestCheckResourceAttrPair("data.crosswork_data_gateway_destinations.by_id", "name", "crosswork_data_gateway_destinations.lookup", "name"),
				resource.TestCheckResourceAttr("data.crosswork_data_gateway_destinations.by_name", "server_type", "grpc"),
				resource.TestCheckResourceAttr("data.crosswork_data_gateway_destinations.by_name", "dispatch_type", "stream"),
				resource.TestCheckResourceAttr("data.crosswork_data_gateway_destinations.by_name", "system_defined", "false"),
				resource.TestCheckResourceAttr("data.crosswork_data_gateway_destinations.by_id", "connections.0.address", "10.14.0.150"),
			),
		}},
	})
}

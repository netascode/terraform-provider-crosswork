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

// Needs a Crosswork with no pool, a free instance (the `CLUSTER` tag) and a virtual IP of its subnet.

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const dataGatewayPoolsLookupConfig = `
data "crosswork_data_gateway_instances" "lookup" {}

resource "crosswork_data_gateway_pools" "lookup" {
  name          = "cwtest-gen-pool-lookup"
  virtual_ips   = ["172.16.34.198"]
  prefix_length = 24
  instances = [[for instance in data.crosswork_data_gateway_instances.lookup.items : instance.id if instance.name == sort([for i in data.crosswork_data_gateway_instances.lookup.items : i.name])[0]][0]]
}

data "crosswork_data_gateway_pools" "by_name" {
  name = crosswork_data_gateway_pools.lookup.name
}

data "crosswork_data_gateway_pools" "by_id" {
  id = crosswork_data_gateway_pools.lookup.id
}
`

func TestAccDataSourceCwDataGatewayPools(t *testing.T) {
	if os.Getenv("CLUSTER") == "" {
		t.Skip("skipping test, set environment variable CLUSTER")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: dataGatewayPoolsLookupConfig,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrPair("data.crosswork_data_gateway_pools.by_name", "id", "crosswork_data_gateway_pools.lookup", "id"),
				resource.TestCheckResourceAttrPair("data.crosswork_data_gateway_pools.by_id", "name", "crosswork_data_gateway_pools.lookup", "name"),
				resource.TestCheckResourceAttr("data.crosswork_data_gateway_pools.by_name", "virtual_ips.0", "172.16.34.198"),
				resource.TestCheckResourceAttr("data.crosswork_data_gateway_pools.by_name", "spares", "0"),
				resource.TestCheckResourceAttr("data.crosswork_data_gateway_pools.by_name", "data_gateway_ids.#", "1"),
				resource.TestCheckResourceAttrPair("data.crosswork_data_gateway_pools.by_name", "data_gateway_ids.0", "crosswork_data_gateway_pools.lookup", "data_gateway_ids.0"),
			),
		}},
	})
}

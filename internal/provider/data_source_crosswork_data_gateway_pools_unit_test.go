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
	"github.com/tidwall/gjson"
)

func seedPool(s *poolStub) {
	s.pools = []*stubPool{{id: "pool-id-1", body: gjson.Parse(`{"name":"cw-dg-pool-1","description":"Example pool","spares":"1","subnetMask":"24",
	  "ipaddrs":[{"inetaddrs":[{"inet_af":"ROBOT_INET_ADDR_TYPE_v4","inet_addr":"192.0.2.195","mask":"24"}]}],"pdgUuids":["dg-1","dg-2"]}`)}}
}

func TestDataGatewayPoolsDataSource(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newPoolStub(t)
	seedPool(s)
	const ds = "data.crosswork_data_gateway_pools.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: s.providerBlock() + "\ndata \"crosswork_data_gateway_pools\" \"test\" {\n  name = \"cw-dg-pool-1\"\n}\n",
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ds, "id", "pool-id-1"),
					resource.TestCheckResourceAttr(ds, "description", "Example pool"),
					resource.TestCheckResourceAttr(ds, "virtual_ips.0", "192.0.2.195"),
					resource.TestCheckResourceAttr(ds, "prefix_length", "24"),
					resource.TestCheckResourceAttr(ds, "spares", "1"),
					resource.TestCheckResourceAttr(ds, "instances.#", "2"),
					resource.TestCheckResourceAttr(ds, "data_gateway_ids.0", "vdg-pool-id-1-1"),
					resource.TestCheckResourceAttr(ds, "protection_status", "PROTECTED"),
				),
			},
			{
				Config:      s.providerBlock() + "\ndata \"crosswork_data_gateway_pools\" \"test\" {\n  name = \"no-such-pool\"\n}\n",
				ExpectError: regexp.MustCompile(`(?s)There is no pool`),
			},
			{
				Config:      s.providerBlock() + "\ndata \"crosswork_data_gateway_pools\" \"test\" {\n  id = \"pool-id-1\"\n  name = \"cw-dg-pool-1\"\n}\n",
				ExpectError: regexp.MustCompile(`(?is)(only one|exactly one|Invalid Attribute Combination)`),
			},
			{
				Config:      s.providerBlock() + "\ndata \"crosswork_data_gateway_pools\" \"test\" {\n}\n",
				ExpectError: regexp.MustCompile(`(?is)(exactly one|No attribute specified|Invalid Attribute Combination)`),
			},
			{
				Config: s.providerBlock() + "\ndata \"crosswork_data_gateway_pools\" \"test\" {\n  id = \"pool-id-1\"\n}\n",
				Check:  resource.TestCheckResourceAttr(ds, "name", "cw-dg-pool-1"),
			},
		},
	})
}

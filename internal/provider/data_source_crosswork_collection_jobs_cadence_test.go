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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataSourceCwCollectionJobsCadence(t *testing.T) {
	const ds = "data.crosswork_collection_jobs_cadence.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "crosswork_collection_jobs_cadence" "test" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(ds, "id", "collection_jobs_cadence"),
					resource.TestCheckResourceAttrSet(ds, "reachability_check"),
					resource.TestCheckResourceAttrSet(ds, "snmp_check"),
					resource.TestCheckResourceAttrSet(ds, "te_tunnel_id"),
					resource.TestCheckResourceAttrSet(ds, "show_clock"),
					func(st *terraform.State) error {
						got, err := readCadence(testDeviceClientNoT())
						if err != nil {
							return err
						}
						a := st.RootModule().Resources[ds].Primary.Attributes
						for attr, want := range map[string]int64{"reachability_check": got.ReachabilityCheck.ValueInt64(), "snmp_check": got.SnmpCheck.ValueInt64(), "te_tunnel_id": got.TeTunnelId.ValueInt64(), "show_clock": got.ShowClock.ValueInt64()} {
							if a[attr] != fmt.Sprint(want) {
								return fmt.Errorf("%s is %s, the controller reports %d", attr, a[attr], want)
							}
						}
						return nil
					},
				),
			},
		},
	})
}

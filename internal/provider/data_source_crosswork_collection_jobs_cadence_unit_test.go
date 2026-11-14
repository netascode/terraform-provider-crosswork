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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCollectionJobsCadenceDataSource(t *testing.T) {
	s := newCadenceStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: s.providerBlock() + `
data "crosswork_collection_jobs_cadence" "test" {
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.crosswork_collection_jobs_cadence.test", "id", "collection_jobs_cadence"),
					resource.TestCheckResourceAttr("data.crosswork_collection_jobs_cadence.test", "reachability_check", "600"),
					resource.TestCheckResourceAttr("data.crosswork_collection_jobs_cadence.test", "snmp_check", "1200"),
					resource.TestCheckResourceAttr("data.crosswork_collection_jobs_cadence.test", "te_tunnel_id", "1200"),
					resource.TestCheckResourceAttr("data.crosswork_collection_jobs_cadence.test", "show_clock", "1800"),
				),
			},
		},
	})
}

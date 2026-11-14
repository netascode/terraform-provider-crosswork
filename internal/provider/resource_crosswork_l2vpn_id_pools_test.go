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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc
func TestAccCwL2VPNIDPools(t *testing.T) {
	var checks []resource.TestCheckFunc

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwL2VPNIDPoolsPrerequisitesConfig + testAccCwL2VPNIDPoolsConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_l2vpn_id_pools.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{},
		ImportStateId:           "l2vpn_id_pools",
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// End of section. //template:end testAcc

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccCwL2VPNIDPoolsPrerequisitesConfig = `
resource "crosswork_resource_pools" "test_evi" {
  name        = "cwtest-l2-evi"
  range_start = 3000
  range_end   = 3100
}

resource "crosswork_resource_pools" "test_st" {
  name        = "cwtest-l2-evi-st"
  range_start = 5000
  range_end   = 5100
}

resource "crosswork_resource_pools" "test_rt" {
  name        = "cwtest-l2-rt"
  range_start = 7000
  range_end   = 7100
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwL2VPNIDPoolsConfig_minimum() string {
	config := `resource "crosswork_l2vpn_id_pools" "test" {` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwL2VPNIDPoolsConfig_all() string {
	config := `resource "crosswork_l2vpn_id_pools" "test" {` + "\n"
	config += `	evi_id_pool = crosswork_resource_pools.test_evi.name` + "\n"
	config += `	evi_source_target_pool = crosswork_resource_pools.test_st.name` + "\n"
	config += `	rt_pool = crosswork_resource_pools.test_rt.name` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc
func TestAccCwL2VPNRoutePolicies(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_l2vpn_route_policies.test", "name", "cwtest-l2-policy"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_l2vpn_route_policies.test", "statements.0.name", "1"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccCwL2VPNRoutePoliciesPrerequisitesConfig + testAccCwL2VPNRoutePoliciesConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccCwL2VPNRoutePoliciesPrerequisitesConfig + testAccCwL2VPNRoutePoliciesConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// End of section. //template:end testAcc

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccCwL2VPNRoutePoliciesPrerequisitesConfig = `
resource "crosswork_l2vpn_rd_sets" "test" {
  name = "cwtest-l2-rds"
  rd   = ["65000:1"]
}

resource "crosswork_l2vpn_tag_sets" "test" {
  name      = "cwtest-l2-tags"
  tag_value = [4812]
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwL2VPNRoutePoliciesConfig_minimum() string {
	config := `resource "crosswork_l2vpn_route_policies" "test" {` + "\n"
	config += `	name = "cwtest-l2-policy"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwL2VPNRoutePoliciesConfig_all() string {
	config := `resource "crosswork_l2vpn_route_policies" "test" {` + "\n"
	config += `	name = "cwtest-l2-policy"` + "\n"
	config += `	statements = [{` + "\n"
	config += `	  name = "1"` + "\n"
	config += `	  match_rd_set = crosswork_l2vpn_rd_sets.test.name` + "\n"
	config += `	  set_ext_community = crosswork_l2vpn_tag_sets.test.name` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

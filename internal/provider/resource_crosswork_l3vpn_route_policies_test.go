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
func TestAccCwL3VPNRoutePolicies(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_l3vpn_route_policies.test", "name", "cwtest-policy"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_l3vpn_route_policies.test", "statements.0.name", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_l3vpn_route_policies.test", "statements.0.policy_result", "accept-route"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccCwL3VPNRoutePoliciesPrerequisitesConfig + testAccCwL3VPNRoutePoliciesConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccCwL3VPNRoutePoliciesPrerequisitesConfig + testAccCwL3VPNRoutePoliciesConfig_all(),
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
const testAccCwL3VPNRoutePoliciesPrerequisitesConfig = `
resource "crosswork_l3vpn_source_prefix_sets" "test" {
  name = "cwtest-source-prefixes"
  mode = "ipv4"
  prefix_list = [{ ip_prefix = "10.250.1.0/24" }]
}

resource "crosswork_l3vpn_destination_prefix_sets" "test" {
  name = "cwtest-dest-prefixes"
  mode = "ipv4"
  prefix_list = [{ ip_prefix = "10.250.2.0/24" }]
}

resource "crosswork_l3vpn_tag_sets" "test" {
  name      = "cwtest-tags"
  tag_value = [4242]
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwL3VPNRoutePoliciesConfig_minimum() string {
	config := `resource "crosswork_l3vpn_route_policies" "test" {` + "\n"
	config += `	name = "cwtest-policy"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwL3VPNRoutePoliciesConfig_all() string {
	config := `resource "crosswork_l3vpn_route_policies" "test" {` + "\n"
	config += `	name = "cwtest-policy"` + "\n"
	config += `	statements = [{` + "\n"
	config += `	  name = "1"` + "\n"
	config += `	  match_source_prefix_set = crosswork_l3vpn_source_prefix_sets.test.name` + "\n"
	config += `	  match_dest_prefix_set = crosswork_l3vpn_destination_prefix_sets.test.name` + "\n"
	config += `	  policy_result = "accept-route"` + "\n"
	config += `	  set_ext_community = crosswork_l3vpn_tag_sets.test.name` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

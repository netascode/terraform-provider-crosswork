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
func TestAccCwDataGatewayPools(t *testing.T) {
	if os.Getenv("CLUSTER") == "" {
		t.Skip("skipping test, set environment variable CLUSTER")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_pools.test", "name", "cwtest-gen-pool"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_pools.test", "description", "Pool of the production gateways"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_pools.test", "virtual_ips.0", "172.16.34.198"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_pools.test", "prefix_length", "24"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_pools.test", "spares", "0"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_pools.test", "cdg_readiness_timeout", "60"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwDataGatewayPoolsPrerequisitesConfig + testAccCwDataGatewayPoolsConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_data_gateway_pools.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// End of section. //template:end testAcc

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccCwDataGatewayPoolsPrerequisitesConfig = `
data "crosswork_data_gateway_instances" "pool_test" {}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwDataGatewayPoolsConfig_minimum() string {
	config := `resource "crosswork_data_gateway_pools" "test" {` + "\n"
	config += `	name = "cwtest-gen-pool"` + "\n"
	config += `	virtual_ips = ["172.16.34.198"]` + "\n"
	config += `	prefix_length = 24` + "\n"
	config += `	instances = [[for instance in data.crosswork_data_gateway_instances.pool_test.items : instance.id if instance.name == sort([for i in data.crosswork_data_gateway_instances.pool_test.items : i.name])[0]][0]]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwDataGatewayPoolsConfig_all() string {
	config := `resource "crosswork_data_gateway_pools" "test" {` + "\n"
	config += `	name = "cwtest-gen-pool"` + "\n"
	config += `	description = "Pool of the production gateways"` + "\n"
	config += `	virtual_ips = ["172.16.34.198"]` + "\n"
	config += `	prefix_length = 24` + "\n"
	config += `	spares = 0` + "\n"
	config += `	instances = [[for instance in data.crosswork_data_gateway_instances.pool_test.items : instance.id if instance.name == sort([for i in data.crosswork_data_gateway_instances.pool_test.items : i.name])[0]][0]]` + "\n"
	config += `	cdg_readiness_timeout = 60` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

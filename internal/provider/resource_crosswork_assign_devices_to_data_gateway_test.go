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
func TestAccCwAssignDevicesToDataGateway(t *testing.T) {
	if os.Getenv("CLUSTER") == "" {
		t.Skip("skipping test, set environment variable CLUSTER")
	}
	var checks []resource.TestCheckFunc

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwAssignDevicesToDataGatewayPrerequisitesConfig + testAccCwAssignDevicesToDataGatewayConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_assign_devices_to_data_gateway.test",
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
const testAccCwAssignDevicesToDataGatewayPrerequisitesConfig = `
data "crosswork_data_gateway_instances" "assign_test" {}

resource "crosswork_data_gateway_pools" "assign_test" {
  name          = "cwtest-gen-assign-pool"
  virtual_ips   = ["172.16.34.198"]
  prefix_length = 24
  instances = [[for instance in data.crosswork_data_gateway_instances.assign_test.items : instance.id if instance.name == sort([for i in data.crosswork_data_gateway_instances.assign_test.items : i.name])[0]][0]]
}

data "crosswork_network_devices" "assign_test" {
  id = "build-router-1"
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwAssignDevicesToDataGatewayConfig_minimum() string {
	config := `resource "crosswork_assign_devices_to_data_gateway" "test" {` + "\n"
	config += `	data_gateway_id = crosswork_data_gateway_pools.assign_test.data_gateway_ids[0]` + "\n"
	config += `	devices = [data.crosswork_network_devices.assign_test.id]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwAssignDevicesToDataGatewayConfig_all() string {
	config := `resource "crosswork_assign_devices_to_data_gateway" "test" {` + "\n"
	config += `	data_gateway_id = crosswork_data_gateway_pools.assign_test.data_gateway_ids[0]` + "\n"
	config += `	devices = [data.crosswork_network_devices.assign_test.id]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

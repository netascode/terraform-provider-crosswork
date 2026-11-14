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
func TestAccCwDataGatewayDestinations(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "name", "cwtest-gen-destination"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "data_source", "data-gateway"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "server_type", "kafka"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "encoding", "gpbkv"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "compression_method", "snappy"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "max_message_size", "100000000"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "buffer_memory", "52428800"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "batch_size", "1048576"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "linger_time", "2000"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "request_timeout", "30"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "collectors.0.name", "CLI"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "collectors.0.buffer_memory", "104857600"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "collectors.0.batch_size", "500000"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "collectors.0.linger_time", "1000"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "collectors.0.request_timeout", "45"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "connections.0.address", "10.14.0.150"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "connections.0.prefix_length", "20"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "connections.0.port", "9192"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_data_gateway_destinations.test", "secure", "false"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwDataGatewayDestinationsConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_data_gateway_destinations.test",
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
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwDataGatewayDestinationsConfig_minimum() string {
	config := `resource "crosswork_data_gateway_destinations" "test" {` + "\n"
	config += `	name = "cwtest-gen-destination"` + "\n"
	config += `	data_source = "data-gateway"` + "\n"
	config += `	server_type = "kafka"` + "\n"
	config += `	compression_method = "snappy"` + "\n"
	config += `	connections = [{` + "\n"
	config += `	  port = 9192` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwDataGatewayDestinationsConfig_all() string {
	config := `resource "crosswork_data_gateway_destinations" "test" {` + "\n"
	config += `	name = "cwtest-gen-destination"` + "\n"
	config += `	data_source = "data-gateway"` + "\n"
	config += `	server_type = "kafka"` + "\n"
	config += `	encoding = "gpbkv"` + "\n"
	config += `	compression_method = "snappy"` + "\n"
	config += `	max_message_size = 100000000` + "\n"
	config += `	buffer_memory = 52428800` + "\n"
	config += `	batch_size = 1048576` + "\n"
	config += `	linger_time = 2000` + "\n"
	config += `	request_timeout = 30` + "\n"
	config += `	collectors = [{` + "\n"
	config += `	  name = "CLI"` + "\n"
	config += `	  buffer_memory = 104857600` + "\n"
	config += `	  batch_size = 500000` + "\n"
	config += `	  linger_time = 1000` + "\n"
	config += `	  request_timeout = 45` + "\n"
	config += `	}]` + "\n"
	config += `	connections = [{` + "\n"
	config += `	  address = "10.14.0.150"` + "\n"
	config += `	  prefix_length = 20` + "\n"
	config += `	  port = 9192` + "\n"
	config += `	}]` + "\n"
	config += `	secure = false` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

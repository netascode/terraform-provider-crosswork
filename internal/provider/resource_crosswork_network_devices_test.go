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
func TestAccCwNetworkDevices(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "name", "cwtest-gen-device"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "host_name", "cwtest-device"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "admin_state", "ROBOT_ADMIN_STATE_UP"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "reachability_check", "REACH_CHECK_ENABLE"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "serial_number", "FOC1234ABCD"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "mac_address", "00:1b:44:11:3a:b7"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "encryption", "aes128-ctr"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "device_type", "NODE_TYPE_ROUTER"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "software_type", "IOS XR"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "syslog_format", "RFC_5424"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "cli_cache_enabled", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "provider_profile", "nso"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "connectivity.0.protocol", "ROBOT_MSVC_TRANS_SSH"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "connectivity.0.ip_address", "192.0.2.231"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "connectivity.0.prefix_length", "24"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "connectivity.0.port", "22"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "connectivity.0.timeout", "60"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "isis_system_id", "0000.0000.0001"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "ospf_router_id", "192.0.2.1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "te_router_id", "192.0.2.2"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "ipv6_router_id", "2001:db8::1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "telemetry_vrf", "default"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "telemetry_source_interface", "Loopback0"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "telemetry_opt_out_of_mdt_config", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "location_building", "Building 1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "location_street", "Main Street"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "location_city", "Amsterdam"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "location_state", "North Holland"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "location_country", "Netherlands"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "location_region", "EMEA"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "location_zip", "1011"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "location_latitude", "52.37"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "location_longitude", "4.89"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_network_devices.test", "location_altitude", "10"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwNetworkDevicesPrerequisitesConfig + testAccCwNetworkDevicesConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_network_devices.test",
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
const testAccCwNetworkDevicesPrerequisitesConfig = `
resource "crosswork_tags" "device_test" {
  name     = "cwtest-gen-device-tag"
  category = "cwtest-category"
}

resource "crosswork_credential_profiles" "device_test" {
  profile = "cwtest-gen-device-profile"
  credentials = [
    {
      type      = "ROBOT_USERPASS_SSH"
      user_name = "cwtest"
      password  = "Cwtest-Password-1"
    }
  ]
  snmpv2_read_community_wo          = "public"
  snmpv2_read_community_wo_version  = 1
  snmpv2_write_community_wo         = "private"
  snmpv2_write_community_wo_version = 1
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwNetworkDevicesConfig_minimum() string {
	config := `resource "crosswork_network_devices" "test" {` + "\n"
	config += `	name = "cwtest-gen-device"` + "\n"
	config += `	credential_profile = crosswork_credential_profiles.device_test.id` + "\n"
	config += `	connectivity = [{` + "\n"
	config += `	  protocol = "ROBOT_MSVC_TRANS_SSH"` + "\n"
	config += `	  ip_address = "192.0.2.231"` + "\n"
	config += `	  prefix_length = 24` + "\n"
	config += `	  port = 22` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwNetworkDevicesConfig_all() string {
	config := `resource "crosswork_network_devices" "test" {` + "\n"
	config += `	name = "cwtest-gen-device"` + "\n"
	config += `	host_name = "cwtest-device"` + "\n"
	config += `	credential_profile = crosswork_credential_profiles.device_test.id` + "\n"
	config += `	admin_state = "ROBOT_ADMIN_STATE_UP"` + "\n"
	config += `	reachability_check = "REACH_CHECK_ENABLE"` + "\n"
	config += `	serial_number = "FOC1234ABCD"` + "\n"
	config += `	mac_address = "00:1b:44:11:3a:b7"` + "\n"
	config += `	encryption = "aes128-ctr"` + "\n"
	config += `	device_type = "NODE_TYPE_ROUTER"` + "\n"
	config += `	software_type = "IOS XR"` + "\n"
	config += `	syslog_format = "RFC_5424"` + "\n"
	config += `	capabilities = ["YANG_CLI"]` + "\n"
	config += `	cli_cache_enabled = false` + "\n"
	config += `	provider_profile = "nso"` + "\n"
	config += `	connectivity = [{` + "\n"
	config += `	  protocol = "ROBOT_MSVC_TRANS_SSH"` + "\n"
	config += `	  ip_address = "192.0.2.231"` + "\n"
	config += `	  prefix_length = 24` + "\n"
	config += `	  port = 22` + "\n"
	config += `	  timeout = 60` + "\n"
	config += `	}]` + "\n"
	config += `	isis_system_id = "0000.0000.0001"` + "\n"
	config += `	ospf_router_id = "192.0.2.1"` + "\n"
	config += `	te_router_id = "192.0.2.2"` + "\n"
	config += `	ipv6_router_id = "2001:db8::1"` + "\n"
	config += `	telemetry_vrf = "default"` + "\n"
	config += `	telemetry_source_interface = "Loopback0"` + "\n"
	config += `	telemetry_opt_out_of_mdt_config = false` + "\n"
	config += `	location_building = "Building 1"` + "\n"
	config += `	location_street = "Main Street"` + "\n"
	config += `	location_city = "Amsterdam"` + "\n"
	config += `	location_state = "North Holland"` + "\n"
	config += `	location_country = "Netherlands"` + "\n"
	config += `	location_region = "EMEA"` + "\n"
	config += `	location_zip = "1011"` + "\n"
	config += `	location_latitude = 52.37` + "\n"
	config += `	location_longitude = 4.89` + "\n"
	config += `	location_altitude = 10` + "\n"
	config += `	tags = [{` + "\n"
	config += `	  name = crosswork_tags.device_test.id` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

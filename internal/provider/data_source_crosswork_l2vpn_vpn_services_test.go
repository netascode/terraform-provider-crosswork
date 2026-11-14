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

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSource
func TestAccDataSourceCwL2VPNVPNServices(t *testing.T) {
	if os.Getenv("L2VPN") == "" {
		t.Skip("skipping test, set environment variable L2VPN")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_id", "cwtest-l2"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_type", "ietf-vpn-common:vpws-evpn"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "evi_id", "3961"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "evi_source", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "evi_target", "2"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "service_assurance_profile_name", "Gold_L2VPN_ConfigProfile system"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "service_assurance_rule_name", "Rule-L2VPN-NM system"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "service_assurance_monitoring_state", "enable"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "service_assurance_preservation", "remove"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.multi_home_ethernet_segment_identifier", "00.11.22.33.44.55.66.77.88"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.sr_policy", "cwtest-policy"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.sr_policy_fallback", "disable"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.id", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.description", "cwtest access"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.interface_id", "GigabitEthernet0/0/0/9"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.encapsulation_type", "ietf-vpn-common:dot1q"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.cvlan_id", "3961"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.mtu", "1500"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.oam_md_name", "EVC3961"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.oam_md_level", "4"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.oam_maid", "3961"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.oam_mep_id", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.oam_message_period", "1s"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.oam_id_type", "number"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.oam_y1731_profiles.0.schedule_interval", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l2vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.oam_y1731_profiles.0.schedule_duration", "5"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwL2VPNVPNServicesPrerequisitesConfig + testAccDataSourceCwL2VPNVPNServicesConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccDataSourceCwL2VPNVPNServicesPrerequisitesConfig = `
data "crosswork_network_devices" "vpn_test" {
  id = "build-router-1"
}

resource "crosswork_l2vpn_y1731_profiles" "vpn_test" {
  name = "cwtest-y1731"
  type = "delay"
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwL2VPNVPNServicesConfig() string {
	config := `resource "crosswork_l2vpn_vpn_services" "test" {` + "\n"
	config += `	vpn_id = "cwtest-l2"` + "\n"
	config += `	vpn_type = "ietf-vpn-common:vpws-evpn"` + "\n"
	config += `	evi_id = 3961` + "\n"
	config += `	evi_source = 1` + "\n"
	config += `	evi_target = 2` + "\n"
	config += `	service_assurance_profile_name = "Gold_L2VPN_ConfigProfile system"` + "\n"
	config += `	service_assurance_rule_name = "Rule-L2VPN-NM system"` + "\n"
	config += `	service_assurance_monitoring_state = "enable"` + "\n"
	config += `	service_assurance_preservation = "remove"` + "\n"
	config += `	vpn_nodes = [{` + "\n"
	config += `	  vpn_node_id = data.crosswork_network_devices.vpn_test.id` + "\n"
	config += `	  multi_home_ethernet_segment_identifier = "00.11.22.33.44.55.66.77.88"` + "\n"
	config += `	  sr_policy = "cwtest-policy"` + "\n"
	config += `	  sr_policy_fallback = "disable"` + "\n"
	config += `	  vpn_network_accesses = [{` + "\n"
	config += `		id = "1"` + "\n"
	config += `		description = "cwtest access"` + "\n"
	config += `		interface_id = "GigabitEthernet0/0/0/9"` + "\n"
	config += `		encapsulation_type = "ietf-vpn-common:dot1q"` + "\n"
	config += `		cvlan_id = 3961` + "\n"
	config += `		mtu = 1500` + "\n"
	config += `		oam_md_name = "EVC3961"` + "\n"
	config += `		oam_md_level = 4` + "\n"
	config += `		oam_maid = "3961"` + "\n"
	config += `		oam_mep_id = 1` + "\n"
	config += `		oam_message_period = "1s"` + "\n"
	config += `		oam_id_type = "number"` + "\n"
	config += `      oam_y1731_profiles = [{` + "\n"
	config += `			name = crosswork_l2vpn_y1731_profiles.vpn_test.name` + "\n"
	config += `			schedule_interval = 5` + "\n"
	config += `			schedule_duration = "5"` + "\n"
	config += `		}]` + "\n"
	config += `	}]` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_l2vpn_vpn_services" "test" {
			id = crosswork_l2vpn_vpn_services.test.id
			vpn_id = "cwtest-l2"
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

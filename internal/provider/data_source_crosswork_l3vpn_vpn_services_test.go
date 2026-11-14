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

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSource
func TestAccDataSourceCwL3VPNVPNServices(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_id", "cwtest-svc"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_service_topology", "ietf-vpn-common:any-to-any"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "service_assurance_profile_name", "Gold_L3VPN_ConfigProfile system"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "service_assurance_rule_name", "Rule-L3VPN-NM system"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "service_assurance_monitoring_state", "enable"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "service_assurance_preservation", "preserve"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_instance_profiles.0.profile_id", "SITE1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_instance_profiles.0.rd", "0:65001:101"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_instance_profiles.0.local_as", "65001"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_instance_profiles.0.address_families.0.address_family", "ietf-vpn-common:ipv4"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_instance_profiles.0.address_families.0.vpn_targets.0.id", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_instance_profiles.0.address_families.0.vpn_targets.0.route_target_type", "both"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_instance_profiles.0.address_families.0.vpn_targets.0.route_targets.0.route_target", "0:65001:4242"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.local_as", "65001"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.active_vpn_instance_profiles.0.profile_id", "SITE1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.active_vpn_instance_profiles.0.rd", "0:65001:202"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.id", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.interface_id", "GigabitEthernet0/0/0/5"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.encapsulation_type", "ietf-vpn-common:dot1q"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.cvlan_id", "100"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.ipv4_address", "10.10.10.1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.ipv4_prefix_length", "30"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.ipv6_address", "2001:db8:10:10::1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.ipv6_prefix_length", "64"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.mtu", "1600"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.qos_profiles.0.direction", "ietf-vpn-common:wan-to-site"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.id", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.type", "ietf-vpn-common:bgp-routing"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.bgp_peer_as", "65100"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.bgp_address_family", "ipv4"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.bgp_multihop", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.bgp_redistribute_connected.0.address_family", "ietf-vpn-common:ipv4"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.bgp_redistribute_connected.0.enable", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.bgp_redistribute_connected.0.metric", "10"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.bgp_update_source_if_type", "GigabitEthernet"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.bgp_update_source_if_id", "0/0/0/5"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_l3vpn_vpn_services.test", "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.bgp_update_source_sub_if_id", "100"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwL3VPNVPNServicesPrerequisitesConfig + testAccDataSourceCwL3VPNVPNServicesConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccDataSourceCwL3VPNVPNServicesPrerequisitesConfig = `
data "crosswork_network_devices" "vpn_test" {
  id = "build-router-1"
}

resource "crosswork_l3vpn_source_prefix_sets" "vpn_test" {
  name        = "cwtest-vpn-source"
  mode        = "ipv4"
  prefix_list = [{ ip_prefix = "10.10.10.0/30" }]
}

resource "crosswork_l3vpn_destination_prefix_sets" "vpn_test" {
  name        = "cwtest-vpn-dest"
  mode        = "ipv4"
  prefix_list = [{ ip_prefix = "10.10.20.0/30" }]
}

resource "crosswork_l3vpn_tag_sets" "vpn_test" {
  name      = "cwtest-vpn-tags"
  tag_value = [4242]
}

resource "crosswork_l3vpn_route_policies" "vpn_test" {
  name = "cwtest-vpn-policy"
  statements = [{
    name                    = "1"
    match_source_prefix_set = crosswork_l3vpn_source_prefix_sets.vpn_test.name
    match_dest_prefix_set   = crosswork_l3vpn_destination_prefix_sets.vpn_test.name
    policy_result           = "accept-route"
    set_ext_community       = crosswork_l3vpn_tag_sets.vpn_test.name
  }]
}

# the policy-map GOLD-IN is on build-router-1, the profile id has to carry its name
resource "crosswork_l3vpn_qos_profile_ids" "vpn_test" {
  name = "GOLD-IN"
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwL3VPNVPNServicesConfig() string {
	config := `resource "crosswork_l3vpn_vpn_services" "test" {` + "\n"
	config += `	vpn_id = "cwtest-svc"` + "\n"
	config += `	vpn_service_topology = "ietf-vpn-common:any-to-any"` + "\n"
	config += `	service_assurance_profile_name = "Gold_L3VPN_ConfigProfile system"` + "\n"
	config += `	service_assurance_rule_name = "Rule-L3VPN-NM system"` + "\n"
	config += `	service_assurance_monitoring_state = "enable"` + "\n"
	config += `	service_assurance_preservation = "preserve"` + "\n"
	config += `	vpn_instance_profiles = [{` + "\n"
	config += `	  profile_id = "SITE1"` + "\n"
	config += `	  rd = "0:65001:101"` + "\n"
	config += `	  local_as = 65001` + "\n"
	config += `	  address_families = [{` + "\n"
	config += `		address_family = "ietf-vpn-common:ipv4"` + "\n"
	config += `      vpn_targets = [{` + "\n"
	config += `			id = 1` + "\n"
	config += `			route_target_type = "both"` + "\n"
	config += `      route_targets = [{` + "\n"
	config += `			route_target = "0:65001:4242"` + "\n"
	config += `		}]` + "\n"
	config += `		}]` + "\n"
	config += `		import_policy = crosswork_l3vpn_route_policies.vpn_test.name` + "\n"
	config += `		export_policy = crosswork_l3vpn_route_policies.vpn_test.name` + "\n"
	config += `	}]` + "\n"
	config += `	}]` + "\n"
	config += `	vpn_nodes = [{` + "\n"
	config += `	  vpn_node_id = data.crosswork_network_devices.vpn_test.id` + "\n"
	config += `	  local_as = 65001` + "\n"
	config += `	  active_vpn_instance_profiles = [{` + "\n"
	config += `		profile_id = "SITE1"` + "\n"
	config += `		rd = "0:65001:202"` + "\n"
	config += `	}]` + "\n"
	config += `	  vpn_network_accesses = [{` + "\n"
	config += `		id = "1"` + "\n"
	config += `		interface_id = "GigabitEthernet0/0/0/5"` + "\n"
	config += `		encapsulation_type = "ietf-vpn-common:dot1q"` + "\n"
	config += `		cvlan_id = 100` + "\n"
	config += `		ipv4_address = "10.10.10.1"` + "\n"
	config += `		ipv4_prefix_length = 30` + "\n"
	config += `		ipv6_address = "2001:db8:10:10::1"` + "\n"
	config += `		ipv6_prefix_length = 64` + "\n"
	config += `		mtu = 1600` + "\n"
	config += `      qos_profiles = [{` + "\n"
	config += `			profile = crosswork_l3vpn_qos_profile_ids.vpn_test.name` + "\n"
	config += `			direction = "ietf-vpn-common:wan-to-site"` + "\n"
	config += `		}]` + "\n"
	config += `      routing_protocols = [{` + "\n"
	config += `			id = "1"` + "\n"
	config += `			type = "ietf-vpn-common:bgp-routing"` + "\n"
	config += `			bgp_peer_as = 65100` + "\n"
	config += `			bgp_address_family = "ipv4"` + "\n"
	config += `			bgp_neighbors = ["10.10.10.2"]` + "\n"
	config += `			bgp_multihop = 5` + "\n"
	config += `      bgp_redistribute_connected = [{` + "\n"
	config += `			address_family = "ietf-vpn-common:ipv4"` + "\n"
	config += `			enable = true` + "\n"
	config += `			metric = 10` + "\n"
	config += `		}]` + "\n"
	config += `			bgp_update_source_if_type = "GigabitEthernet"` + "\n"
	config += `			bgp_update_source_if_id = "0/0/0/5"` + "\n"
	config += `			bgp_update_source_sub_if_id = 100` + "\n"
	config += `		}]` + "\n"
	config += `	}]` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_l3vpn_vpn_services" "test" {
			id = crosswork_l3vpn_vpn_services.test.id
			vpn_id = "cwtest-svc"
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

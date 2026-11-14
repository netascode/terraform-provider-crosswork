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

// Life cycle of the service of the generated test (the example values and prerequisites of the definition, on build-router-1)
// that the generated test cannot express: a change in place, a changed RD (which replaces the service), the removal
// of the node, which finishes in the background so that the update waits for the plan, and the node added again at once, which
// NSO refuses (zombie) unless the removal waited. The device is left as it was.

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

const vpnRes = "crosswork_l3vpn_vpn_services.test"

func TestAccCwL3VPNVPNServicesLifecycle(t *testing.T) {
	full := testAccCwL3VPNVPNServicesPrerequisitesConfig + testAccCwL3VPNVPNServicesConfig_all()
	// a route target and the BGP neighbor change: one PUT of the whole service
	changed := strings.NewReplacer("0:65001:4242", "0:65001:4243", "10.10.10.2", "10.10.10.3").Replace(full)
	// the RD of the node changes: the service is replaced
	rdChanged := strings.ReplaceAll(changed, "0:65001:202", "0:65001:203")
	noNode := changed[:strings.Index(changed, "\tvpn_nodes = [{")] + "}\n"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: full, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(vpnRes, "vpn_nodes.#", "1"),
				resource.TestCheckResourceAttr(vpnRes, "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.bgp_neighbors.#", "1"),
			)},
			{Config: changed, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(vpnRes, "vpn_instance_profiles.0.address_families.0.vpn_targets.0.route_targets.0.route_target", "0:65001:4243"),
				resource.TestCheckResourceAttr(vpnRes, "vpn_nodes.0.vpn_network_accesses.0.routing_protocols.0.bgp_neighbors.0", "10.10.10.3"),
			)},
			{Config: noNode, Check: resource.TestCheckResourceAttr(vpnRes, "vpn_nodes.#", "0")},
			// a node that is added never replaces the service
			{Config: changed, ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(vpnRes, plancheck.ResourceActionUpdate)}}, Check: resource.TestCheckResourceAttr(vpnRes, "vpn_nodes.#", "1")},
			// NSO refuses a changed RD of an active node, so the service is replaced
			{Config: rdChanged, ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(vpnRes, plancheck.ResourceActionDestroyBeforeCreate)}}, Check: resource.TestCheckResourceAttr(vpnRes, "vpn_nodes.0.active_vpn_instance_profiles.0.rd", "0:65001:203")},
		},
	})
}

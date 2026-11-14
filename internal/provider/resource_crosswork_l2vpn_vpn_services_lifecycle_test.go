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

// NSO refuses an L2VPN (EVPN VPWS) service without its two nodes, so the generated tests of the service are skipped (set L2VPN
// to run them) and the life cycle runs here: the plain service on build-router-1 and build-router-2, which are always onboarded,
// an Ethernet OAM service (set L2VPN_OAM; the XRv9k refuses Y.1731, so it needs the two 8201s BC01-INV-8201-02 and
// BC02-INV-8201-04), and a service whose EVI numbers come from the pools of the L2VPN id pools.

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

const l2Res = "crosswork_l2vpn_vpn_services.test"

func l2Devices(a, b string) string {
	return fmt.Sprintf(`
data "crosswork_network_devices" "a" {
  id = "%s"
}

data "crosswork_network_devices" "b" {
  id = "%s"
}
`, a, b)
}

func l2Service(evi string, access func(node string) string, depends string) string {
	return fmt.Sprintf(`
resource "crosswork_l2vpn_vpn_services" "test" {
  vpn_id   = "cwtest-l2-life"
  vpn_type = "ietf-vpn-common:vpws-evpn"
%s
  vpn_nodes = [
    {
      vpn_node_id = data.crosswork_network_devices.a.id
      vpn_network_accesses = [{
%s
      }]
    },
    {
      vpn_node_id = data.crosswork_network_devices.b.id
      vpn_network_accesses = [{
%s
      }]
    }
  ]
%s
}
`, evi, access("a"), access("b"), depends)
}

func l2Plain(iface string, vlan int, extra string) func(string) string {
	return func(string) string {
		return fmt.Sprintf(`        id                  = "1"
        interface_id        = "%s"
        encapsulation_type  = "ietf-vpn-common:dot1q"
        cvlan_id            = %d
%s`, iface, vlan, extra)
	}
}

func TestAccCwL2VPNVPNServicesLifecycle(t *testing.T) {
	devices := l2Devices("build-router-1", "build-router-2")
	evi := "  evi_id     = 3961\n  evi_source = 1\n  evi_target = 2\n"
	created := devices + l2Service(evi, l2Plain("GigabitEthernet0/0/0/9", 3961, "        description = \"cwtest\"\n        mtu = 1500\n"), "")
	// an EVI number, the VLAN, the MTU and a tag operation change: one PUT of the whole service
	changed := devices + l2Service("  evi_id     = 3962\n  evi_source = 3\n  evi_target = 4\n",
		l2Plain("GigabitEthernet0/0/0/9", 3962, "        mtu = 1400\n        tag_operation_translate = \"1-to-1\"\n        tag_operation_tag = 200\n"), "")
	untagged := devices + l2Service("  evi_id     = 3962\n  evi_source = 3\n  evi_target = 4\n", func(string) string {
		return "        id                 = \"1\"\n        interface_id       = \"GigabitEthernet0/0/0/9\"\n        encapsulation_type = \"ietf-vpn-common:untagged-int\"\n"
	}, "")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: created, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(l2Res, "vpn_nodes.#", "2"),
				resource.TestCheckResourceAttr(l2Res, "evi_id", "3961"),
				resource.TestCheckResourceAttr(l2Res, "vpn_nodes.0.vpn_network_accesses.0.cvlan_id", "3961"),
			)},
			{Config: changed, ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(l2Res, plancheck.ResourceActionUpdate)}},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(l2Res, "evi_id", "3962"),
					resource.TestCheckResourceAttr(l2Res, "vpn_nodes.0.vpn_network_accesses.0.tag_operation_tag", "200"),
					resource.TestCheckResourceAttr(l2Res, "vpn_nodes.0.vpn_network_accesses.0.mtu", "1400"),
				)},
			{Config: untagged, ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(l2Res, plancheck.ResourceActionUpdate)}},
				Check: resource.TestCheckResourceAttr(l2Res, "vpn_nodes.0.vpn_network_accesses.0.encapsulation_type", "ietf-vpn-common:untagged-int")},
		},
	})
}

func TestAccCwL2VPNVPNServicesOAM(t *testing.T) {
	if os.Getenv("L2VPN_OAM") == "" {
		t.Skip("skipping test, set environment variable L2VPN_OAM (needs the devices BC01-INV-8201-02 and BC02-INV-8201-04)")
	}
	profiles := `
resource "crosswork_l2vpn_y1731_profiles" "delay" {
  name                       = "cwtest-l2-delay"
  type                       = "delay"
  schedule_interval          = 5
  schedule_duration          = "5"
  probe_message_count        = 10
  probe_message_period       = 100
  probe_measurement_interval = "60"
  probe_frame_size           = 1000
  probe_priority             = 5
  delay_statistics           = [{ type = "jitter-two-way" }, { type = "delay-two-way" }]
  bucket_size                = 1
  bucket_archive             = 20
}

resource "crosswork_l2vpn_y1731_profiles" "loss" {
  name                       = "cwtest-l2-loss"
  type                       = "synthetic-loss"
  loss_statistics            = [{ type = "loss-sd" }, { type = "loss-ds" }]
  schedule_interval          = 1
  schedule_duration          = "1"
  probe_message_count        = 10
  probe_message_period       = 100
  probe_measurement_interval = "60"
  probe_frame_size           = 1000
  probe_priority             = 5
  bucket_size                = 1
  bucket_archive             = 20
}
`
	oam := func(md string, mep int, period string) func(node string) string {
		return func(node string) string {
			m := mep
			if node == "b" {
				m = mep + 1
			}
			return fmt.Sprintf(`        id                 = "1"
        interface_id       = "HundredGigE0/0/0/31"
        encapsulation_type = "ietf-vpn-common:dot1q"
        cvlan_id           = 3963
        oam_md_name        = "%s"
        oam_md_level       = 4
        oam_maid           = "3963"
        oam_mep_id         = %d
        oam_message_period = "%s"
        oam_id_type        = "number"
        oam_y1731_profiles = [
          { name = crosswork_l2vpn_y1731_profiles.delay.name, schedule_interval = 5, schedule_duration = "5" },
          { name = crosswork_l2vpn_y1731_profiles.loss.name },
        ]
`, md, m, period)
		}
	}
	devices := l2Devices("BC01-INV-8201-02", "BC02-INV-8201-04")
	evi := "  evi_id     = 3963\n  evi_source = 1\n  evi_target = 2\n"
	created := profiles + devices + l2Service(evi, oam("EVC3963", 1, "1s"), "")
	changed := profiles + devices + l2Service(evi, oam("EVC3964", 11, "10s"), "")
	noOAM := profiles + devices + l2Service(evi, l2Plain("HundredGigE0/0/0/31", 3963, ""), "")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: created, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(l2Res, "vpn_nodes.0.vpn_network_accesses.0.oam_md_name", "EVC3963"),
				resource.TestCheckResourceAttr(l2Res, "vpn_nodes.0.vpn_network_accesses.0.oam_y1731_profiles.#", "2"),
			)},
			{Config: changed, ConfigPlanChecks: resource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(l2Res, plancheck.ResourceActionUpdate)}},
				Check: resource.TestCheckResourceAttr(l2Res, "vpn_nodes.0.vpn_network_accesses.0.oam_mep_id", "11")},
			{Config: noOAM, Check: resource.TestCheckNoResourceAttr(l2Res, "vpn_nodes.0.vpn_network_accesses.0.oam_md_name")},
		},
	})
}

func TestAccCwL2VPNVPNServicesAutoEVI(t *testing.T) {
	pools := `
resource "crosswork_resource_pools" "evi" {
  name        = "cwtest-l2-life-evi"
  range_start = 3200
  range_end   = 3299
}

resource "crosswork_resource_pools" "st" {
  name        = "cwtest-l2-life-st"
  range_start = 5200
  range_end   = 5299
}

resource "crosswork_l2vpn_id_pools" "test" {
  evi_id_pool            = crosswork_resource_pools.evi.name
  evi_source_target_pool = crosswork_resource_pools.st.name
}
`
	auto := "  auto_evi_id     = true\n  auto_evi_source = true\n  auto_evi_target = true\n"
	config := pools + l2Devices("build-router-1", "build-router-2") + l2Service(auto,
		l2Plain("GigabitEthernet0/0/0/9", 3964, ""), "  depends_on = [crosswork_l2vpn_id_pools.test]")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(l2Res, "auto_evi_id", "true"),
				resource.TestCheckResourceAttr(l2Res, "auto_evi_source", "true"),
				resource.TestCheckResourceAttr(l2Res, "auto_evi_target", "true"),
			)},
		},
	})
}

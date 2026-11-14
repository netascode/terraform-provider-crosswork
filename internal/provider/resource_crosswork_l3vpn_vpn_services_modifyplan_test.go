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

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func vpnNode(id string, rds map[string]string) L3VPNVPNServicesVpnNodes {
	node := L3VPNVPNServicesVpnNodes{VpnNodeId: types.StringValue(id)}
	for profile, rd := range rds {
		p := L3VPNVPNServicesVpnNodesActiveVpnInstanceProfiles{ProfileId: types.StringValue(profile), Rd: types.StringValue(rd)}
		if rd == "" {
			p.Rd = types.StringNull()
		}
		node.ActiveVpnInstanceProfiles = append(node.ActiveVpnInstanceProfiles, p)
	}
	return node
}

func vpnProfile(id, rd string) L3VPNVPNServicesVpnInstanceProfiles {
	p := L3VPNVPNServicesVpnInstanceProfiles{ProfileId: types.StringValue(id), Rd: types.StringValue(rd)}
	if rd == "" {
		p.Rd = types.StringNull()
	}
	return p
}

func TestL3VPNRdReplacePaths(t *testing.T) {
	base := L3VPNVPNServices{
		VpnInstanceProfiles: []L3VPNVPNServicesVpnInstanceProfiles{vpnProfile("A", "0:65001:1"), vpnProfile("B", "")},
		VpnNodes: []L3VPNVPNServicesVpnNodes{
			vpnNode("r1", map[string]string{"A": "0:65001:11"}),
			vpnNode("r2", map[string]string{"A": "0:65001:12", "B": ""}),
		},
	}
	with := func(f func(s *L3VPNVPNServices)) L3VPNVPNServices {
		s := L3VPNVPNServices{VpnInstanceProfiles: append([]L3VPNVPNServicesVpnInstanceProfiles{}, base.VpnInstanceProfiles...)}
		for _, n := range base.VpnNodes {
			n.ActiveVpnInstanceProfiles = append([]L3VPNVPNServicesVpnNodesActiveVpnInstanceProfiles{}, n.ActiveVpnInstanceProfiles...)
			s.VpnNodes = append(s.VpnNodes, n)
		}
		f(&s)
		return s
	}
	cases := []struct {
		name string
		plan L3VPNVPNServices
		want []string
	}{
		{"no change", with(func(s *L3VPNVPNServices) {}), nil},
		{"profile rd changed", with(func(s *L3VPNVPNServices) { s.VpnInstanceProfiles[0] = vpnProfile("A", "0:65001:99") }), []string{"vpn_instance_profiles"}},
		{"profile rd set where there was none", with(func(s *L3VPNVPNServices) { s.VpnInstanceProfiles[1] = vpnProfile("B", "0:65001:5") }), []string{"vpn_instance_profiles"}},
		{"profile rd removed", with(func(s *L3VPNVPNServices) { s.VpnInstanceProfiles[0] = vpnProfile("A", "") }), []string{"vpn_instance_profiles"}},
		{"node rd changed", with(func(s *L3VPNVPNServices) { s.VpnNodes[0] = vpnNode("r1", map[string]string{"A": "0:65001:99"}) }), []string{"vpn_nodes"}},
		{"node rd set where there was none", with(func(s *L3VPNVPNServices) {
			s.VpnNodes[1] = vpnNode("r2", map[string]string{"A": "0:65001:12", "B": "0:65001:7"})
		}), []string{"vpn_nodes"}},
		{"profile added", with(func(s *L3VPNVPNServices) {
			s.VpnInstanceProfiles = append(s.VpnInstanceProfiles, vpnProfile("C", "0:65001:3"))
		}), nil},
		{"profile removed", with(func(s *L3VPNVPNServices) { s.VpnInstanceProfiles = s.VpnInstanceProfiles[:1] }), nil},
		{"node added with an rd", with(func(s *L3VPNVPNServices) {
			s.VpnNodes = append(s.VpnNodes, vpnNode("r3", map[string]string{"A": "0:65001:13"}))
		}), nil},
		{"node removed", with(func(s *L3VPNVPNServices) { s.VpnNodes = s.VpnNodes[:1] }), nil},
		{"node gets a profile it did not have", with(func(s *L3VPNVPNServices) {
			s.VpnNodes[0] = vpnNode("r1", map[string]string{"A": "0:65001:11", "B": "0:65001:8"})
		}), nil},
		{"node added and a profile rd changed", with(func(s *L3VPNVPNServices) {
			s.VpnNodes = append(s.VpnNodes, vpnNode("r3", map[string]string{"A": "0:65001:13"}))
			s.VpnInstanceProfiles[0] = vpnProfile("A", "0:65001:99")
		}), []string{"vpn_instance_profiles"}},
		{"rd not known yet", with(func(s *L3VPNVPNServices) {
			s.VpnInstanceProfiles[0].Rd = types.StringUnknown()
			s.VpnNodes[0].ActiveVpnInstanceProfiles[0].Rd = types.StringUnknown()
		}), nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := l3vpnRdReplacePaths(c.plan, base)
			var names []string
			for _, p := range got {
				names = append(names, p.String())
			}
			if len(names) != len(c.want) {
				t.Fatalf("paths %v, want %v", names, c.want)
			}
			for i := range names {
				if names[i] != c.want[i] {
					t.Fatalf("paths %v, want %v", names, c.want)
				}
			}
		})
	}
}

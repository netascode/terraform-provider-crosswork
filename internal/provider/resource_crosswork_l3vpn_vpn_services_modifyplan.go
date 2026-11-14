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

// NSO refuses to change the RD of a node while its VRF address family is active
// (TSDN-L3VPN-420), and it accepts a change of the RD of a profile without pushing it to the device, so a change of either
// RD of a profile or node that is already in the state means a new service. The RDs sit inside lists, where a plan modifier
// of the attribute compares by index (a node that is added or removed before another one shifts the others and would replace
// the service), so the elements are matched by their keys here: vpn_instance_profiles by profile_id, and the profiles of a node by vpn_node_id
// and profile_id. A profile or node that is new or removed never forces a replacement.

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.ResourceWithModifyPlan = (*L3VPNVPNServicesResource)(nil)

func (r *L3VPNVPNServicesResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// a create or a destroy has nothing to compare
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}
	var plan, state L3VPNVPNServices
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.RequiresReplace = append(resp.RequiresReplace, l3vpnRdReplacePaths(plan, state)...)
}

func l3vpnRdReplacePaths(plan, state L3VPNVPNServices) []path.Path {
	var paths []path.Path
	stateProfiles := map[string]types.String{}
	for _, p := range state.VpnInstanceProfiles {
		stateProfiles[p.ProfileId.ValueString()] = p.Rd
	}
	for _, p := range plan.VpnInstanceProfiles {
		if old, ok := stateProfiles[p.ProfileId.ValueString()]; ok && rdChanged(old, p.Rd) {
			paths = append(paths, path.Root("vpn_instance_profiles"))
			break
		}
	}

	stateNodes := map[string]types.String{}
	for _, n := range state.VpnNodes {
		for _, p := range n.ActiveVpnInstanceProfiles {
			stateNodes[n.VpnNodeId.ValueString()+"\x00"+p.ProfileId.ValueString()] = p.Rd
		}
	}
	for _, n := range plan.VpnNodes {
		for _, p := range n.ActiveVpnInstanceProfiles {
			if old, ok := stateNodes[n.VpnNodeId.ValueString()+"\x00"+p.ProfileId.ValueString()]; ok && rdChanged(old, p.Rd) {
				return append(paths, path.Root("vpn_nodes"))
			}
		}
	}
	return paths
}

// An RD that is not known yet is left to the apply.
func rdChanged(old, planned types.String) bool {
	if planned.IsUnknown() || old.IsUnknown() {
		return false
	}
	return old.IsNull() != planned.IsNull() || old.ValueString() != planned.ValueString()
}

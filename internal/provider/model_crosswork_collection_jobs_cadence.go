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

import "github.com/hashicorp/terraform-plugin-framework/types"

type CollectionJobsCadence struct {
	Id                types.String `tfsdk:"id"`
	Instance          types.String `tfsdk:"instance"`
	ReachabilityCheck types.Int64  `tfsdk:"reachability_check"`
	SnmpCheck         types.Int64  `tfsdk:"snmp_check"`
	TeTunnelId        types.Int64  `tfsdk:"te_tunnel_id"`
	ShowClock         types.Int64  `tfsdk:"show_clock"`
}

func (m CollectionJobsCadence) with(v cadenceValues) CollectionJobsCadence {
	m.ReachabilityCheck, m.SnmpCheck, m.TeTunnelId, m.ShowClock = v.ReachabilityCheck, v.SnmpCheck, v.TeTunnelId, v.ShowClock
	return m
}

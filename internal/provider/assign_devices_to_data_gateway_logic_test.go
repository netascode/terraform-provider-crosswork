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

import "testing"

func desiredDevices(devices ...gatewayDevice) map[string]gatewayDevice {
	out := map[string]gatewayDevice{}
	for _, d := range devices {
		out[d.Name] = d
	}
	return out
}

func TestPlanMapping(t *testing.T) {
	const gw = "gw-1"
	desired := desiredDevices(
		gatewayDevice{Name: "on", Uuid: "u-on", Gateway: gw},
		gatewayDevice{Name: "free", Uuid: "u-free"},
		gatewayDevice{Name: "other", Uuid: "u-other", Gateway: "gw-2"},
	)
	members := map[string]string{"on": "u-on", "left": "u-left", "foreign": "u-foreign"}
	m := planMapping(gw, desired, members, []string{"on", "left"})

	assert.Equal(t, []string{"u-free"}, m.Add, "a free device is added")
	assert.Equal(t, []string{"u-other"}, m.Move, "a device of another gateway is moved, not added")
	assert.Equal(t, []string{"u-left"}, m.Remove, "a device of the resource that left the plan is removed")
	assert.False(t, m.empty())
}

func TestPlanMappingLeavesWhatIsNotTheResourcesAlone(t *testing.T) {
	m := planMapping("gw-1", nil, map[string]string{"foreign": "u-foreign"}, nil)
	assert.True(t, m.empty(), "a device that someone else attached is not removed")

	m = planMapping("gw-1", desiredDevices(gatewayDevice{Name: "on", Uuid: "u-on", Gateway: "gw-1"}), map[string]string{"on": "u-on"}, nil)
	assert.True(t, m.empty(), "a device that is on the gateway needs no call, which an ADD would refuse")

	// a device that left in the plan and is no longer on the gateway needs no REMOVE, which would be refused
	m = planMapping("gw-1", nil, map[string]string{}, []string{"gone"})
	assert.True(t, m.empty())
}

func TestReached(t *testing.T) {
	members := map[string]string{"a": "1", "b": "2"}
	assert.True(t, reached([]string{"a", "b"}, nil, members))
	assert.False(t, reached([]string{"a", "c"}, nil, members), "a wanted device that is not on the gateway")
	assert.False(t, reached([]string{"a"}, []string{"a", "b"}, members), "a device that had to leave is still there")
	assert.True(t, reached([]string{"a", "b"}, []string{"a", "b"}, members))
	assert.True(t, reached(nil, []string{"x"}, members))
}

func TestGatewayIdPattern(t *testing.T) {
	assert.True(t, gatewayIdPattern.MatchString("00000000-0000-4000-8000-000000000004"))
	assert.False(t, gatewayIdPattern.MatchString("x' or 1=1"), "the id goes into the criteria of a query")
	assert.False(t, gatewayIdPattern.MatchString(""))
}

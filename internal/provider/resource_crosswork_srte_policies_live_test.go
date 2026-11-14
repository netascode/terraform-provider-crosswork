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

// Steps a definition cannot express: an explicit path that becomes a dynamic one, paths that are written in another order than
// NSO returns them (it sorts them by preference), a change of the color and of the tail end (the device policy is named after
// both, so NSO removes the old policy and creates the new one in the same call), a move to the other head end, and a new name,
// which replaces the policy: the old one is deleted and the new one, with the same head end, color and tail end, created in the
// same apply, which NSO refuses ("Service still in zombie state") unless the delete waited for the plan to go.
// The head ends are build-router-1 and build-router-2, always onboarded and synced with NSO; the test leaves them as they were.
// Color 999 is the color of the policy tests.

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const policyRes = "crosswork_srte_policies.test"

func TestAccCwSRTEPoliciesLifecycle(t *testing.T) {
	const first, second = "build-router-1", "build-router-2"
	const sidList = `
resource "crosswork_srte_sid_lists" "test" {
  name = "cwtest-policy-life-sl"
  sids = [{ index = 1, label = 16001 }, { index = 2, address = "192.0.2.10" }]
}
`
	explicit := `{ preference = 100, explicit_sid_lists = [{ name = crosswork_srte_sid_lists.test.name, weight = 10 }] }`
	dynamic := `{ preference = 200, dynamic_pce = true, dynamic_metric_type = "igp" }`
	config := func(name, head, tail string, color int, paths ...string) string {
		list := ""
		for _, p := range paths {
			list += "    " + p + ",\n"
		}
		return sidList + fmt.Sprintf("resource \"crosswork_srte_policies\" \"test\" {\n  name     = %q\n  head_end = %q\n  tail_end = %q\n  color    = %d\n  paths = [\n%s  ]\n}\n", name, head, tail, color, list)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config("cwtest-policy-life", first, "192.0.2.61", 999, explicit), Check: resource.TestCheckResourceAttr(policyRes, "paths.0.explicit_sid_lists.#", "1")},
			{Config: config("cwtest-policy-life", first, "192.0.2.61", 999, `{ preference = 100, dynamic_pce = true, dynamic_metric_type = "igp" }`), Check: resource.TestCheckResourceAttr(policyRes, "paths.0.dynamic_metric_type", "igp")},
			// two paths, written in the opposite order of the preference: NSO returns them sorted and the plan must stay empty
			{Config: config("cwtest-policy-life", first, "192.0.2.61", 999, dynamic, explicit), Check: resource.TestCheckResourceAttr(policyRes, "paths.#", "2")},
			{Config: config("cwtest-policy-life", first, "192.0.2.62", 998, dynamic, explicit), Check: resource.TestCheckResourceAttr(policyRes, "color", "998")},
			// the other head end: NSO removes the policy from the first and pushes it to the second
			{Config: config("cwtest-policy-life", second, "192.0.2.62", 998, dynamic, explicit), Check: resource.TestCheckResourceAttr(policyRes, "head_end", second)},
			{Config: config("cwtest-policy-life-b", second, "192.0.2.62", 998, dynamic, explicit), Check: resource.TestCheckResourceAttr(policyRes, "name", "cwtest-policy-life-b")},
		},
	})
}

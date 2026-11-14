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

// The entries change in place (a label, an added entry, an IPv4 entry, a removed entry) and a new name replaces the list, which
// NSO accepts at once because a SID list has no plan and no zombie state. A SID list is only NSO configuration until a policy
// uses it, so no device is involved.

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCwSRTESIDListsLifecycle(t *testing.T) {
	const res = "crosswork_srte_sid_lists.test"
	config := func(name, sids string) string {
		return `resource "crosswork_srte_sid_lists" "test" {
  name = "` + name + `"
  sids = [` + sids + `]
}
`
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config("cwtest-sid-life", `{ index = 1, label = 16001 }, { index = 2, label = 16002 }`), Check: resource.TestCheckResourceAttr(res, "sids.#", "2")},
			{Config: config("cwtest-sid-life", `{ index = 1, label = 16001 }, { index = 2, label = 16009 }`), Check: resource.TestCheckResourceAttr(res, "sids.1.label", "16009")},
			{Config: config("cwtest-sid-life", `{ index = 1, label = 16001 }, { index = 2, label = 16009 }, { index = 3, address = "192.0.2.10" }`), Check: resource.TestCheckResourceAttr(res, "sids.2.address", "192.0.2.10")},
			{Config: config("cwtest-sid-life", `{ index = 1, label = 16001 }, { index = 2, address = "192.0.2.10" }`), Check: resource.TestCheckResourceAttr(res, "sids.#", "2")},
			{Config: config("cwtest-sid-life-b", `{ index = 1, label = 16001 }, { index = 2, address = "192.0.2.10" }`), Check: resource.TestCheckResourceAttr(res, "name", "cwtest-sid-life-b")},
		},
	})
}

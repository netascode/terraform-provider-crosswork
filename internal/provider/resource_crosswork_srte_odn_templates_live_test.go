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

// Steps a definition cannot express: dropping the head end, adding it back, and a change of the name, which replaces the
// template (the old one is deleted and the new one created in the same apply, which NSO refuses ("Service still in zombie
// state") unless the delete waited for the plan to go). The head end is build-router-1, always onboarded and synced with NSO;
// the test leaves it as it was.

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCwSRTEODNTemplatesHeadEnd(t *testing.T) {
	const device = "build-router-1"
	config := func(name string, color int, headEnd bool) string {
		he := ""
		if headEnd {
			he = fmt.Sprintf("  head_end = [{ name = %q }]\n", device)
		}
		return fmt.Sprintf("resource \"crosswork_srte_odn_templates\" \"test\" {\n  name  = %q\n  color = %d\n%s  dynamic_pce = true\n  dynamic_metric_type = \"igp\"\n}\n", name, color, he)
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config("cwtest-odn-a", 4701, true), Check: resource.TestCheckResourceAttr(odnRes, "head_end.#", "1")},
			// the head end is dropped: the device is cleaned in the background and the update waits for it
			{Config: config("cwtest-odn-a", 4702, false), Check: resource.TestCheckResourceAttr(odnRes, "color", "4702")},
			{Config: config("cwtest-odn-a", 4703, true), Check: resource.TestCheckResourceAttr(odnRes, "head_end.#", "1")},
			{Config: config("cwtest-odn-b", 4703, true), Check: resource.TestCheckResourceAttr(odnRes, "name", "cwtest-odn-b")},
		},
	})
}

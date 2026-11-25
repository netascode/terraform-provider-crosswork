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

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAlarmSeveritiesDataSourceReadsEveryType(t *testing.T) {
	s := newAlarmStub(t)
	const ds = "data.crosswork_alarm_severities.test"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: s.providerBlock() + `
data "crosswork_alarm_severities" "test" {
}
`,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(ds, "id", "alarm_severities"),
				resource.TestCheckResourceAttr(ds, "severities.#", "5"),
				resource.TestCheckResourceAttr(ds, "severities.0.name", "AUTO-AMPLI-MISMATCH"),
				resource.TestCheckResourceAttr(ds, "severities.0.severity", "information"),
				resource.TestCheckNoResourceAttr(ds, "severities.0.auto_clear_minutes"),
				resource.TestCheckResourceAttr(ds, "severities.1.severity", "minor"),
				resource.TestCheckResourceAttr(ds, "severities.1.auto_clear_minutes", "1440"),
				resource.TestCheckResourceAttr(ds, "severities.4.name", "ADHOC_ROGUE_AUTO_CONTAINED"),
			),
		}},
	})
}

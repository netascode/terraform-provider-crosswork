// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://mozilla.org/MPL/2.0/
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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestGoldenConfigGlobalVariables(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	cfg := func(typ, value string) string {
		return s.providerBlock() + fmt.Sprintf(`
resource "crosswork_golden_config_global_variables" "v" {
  name  = "cwtest_gc_v"
  type  = %q
  value = %q
}

data "crosswork_golden_config_global_variables" "v" {
  id         = crosswork_golden_config_global_variables.v.id
  depends_on = [crosswork_golden_config_global_variables.v]
}
`, typ, value)
	}
	const res = "crosswork_golden_config_global_variables.v"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{
			{
				Config: cfg("STRING", "10.1.1.1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", "cwtest_gc_v"),
					resource.TestCheckResourceAttr(res, "variable_id", "10"),
					resource.TestCheckResourceAttr("data.crosswork_golden_config_global_variables.v", "value", "10.1.1.1"),
					resource.TestCheckResourceAttr("data.crosswork_golden_config_global_variables.v", "type", "STRING"),
				),
				ConfigPlanChecks: emptyPlanCheck,
			},
			{ // an update is a PATCH by name
				PreConfig:        s.resetCalls,
				Config:           cfg("STRING", "10.9.9.9"),
				Check:            s.callsAre("PATCH gvar cwtest_gc_v"),
				ConfigPlanChecks: emptyPlanCheck,
			},
			{ // a JSON value is written as the user writes it and the controller's sorted form causes no difference
				Config:           cfg("JSON", `{ "b": 1, "a": 2 }`),
				Check:            resource.TestCheckResourceAttr(res, "value", `{ "b": 1, "a": 2 }`),
				ConfigPlanChecks: emptyPlanCheck,
			},
			{
				ResourceName:  res,
				ImportState:   true,
				ImportStateId: "cwtest_gc_v",
				ImportStateCheck: func(st []*terraform.InstanceState) error {
					if a := st[0].Attributes; a["type"] != "JSON" || a["value"] != `{"a":2,"b":1}` {
						return fmt.Errorf("unexpected imported state: %v", a)
					}
					return nil
				},
			},
			{ // deleted on the controller: the plan creates it again
				PreConfig: func() {
					s.mu.Lock()
					delete(s.gvars, "cwtest_gc_v")
					s.mu.Unlock()
				},
				Config:           cfg("JSON", `{ "b": 1, "a": 2 }`),
				ConfigPlanChecks: emptyPlanCheck,
			},
		},
	})
}

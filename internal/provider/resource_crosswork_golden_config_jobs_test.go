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

// Applies a template to build-router-1 and build-router-2 (IOS XR, NSO provider, synced with NSO), which are always onboarded;
// the test neither creates nor deletes them.

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func testAccCwGoldenConfigConfig(device1 string, trigger string) string {
	return fmt.Sprintf(`
resource "crosswork_golden_config_global_variables" "test" {
  name  = "cwtest_gc_desc"
  value = "global-value"
}

resource "crosswork_golden_config_templates" "test" {
  name           = "cwtest-gc-acc"
  vendor         = "Cisco Systems"
  software_type  = "IOS XR"
  product_series = ["Cisco XRv 9000 Series Virtual Routers"]
  tags           = ["cwtest"]
  config         = "interface Loopback99\n description {{ cwtest_gc_desc }}\n"
  depends_on     = [crosswork_golden_config_global_variables.test]
}

resource "crosswork_golden_config_jobs" "test" {
  template_name  = crosswork_golden_config_templates.test.name
  template_config_id = crosswork_golden_config_templates.test.config_id
  vendor         = "Cisco Systems"
  product_series = "Cisco XRv 9000 Series Virtual Routers"
  mop            = "Default Golden Config - Cisco"
  tags           = ["cwtest"]
  trigger        = %q
  devices = {
    "build-router-1" = {}
    "build-router-2" = %s
  }
}

data "crosswork_golden_config_jobs" "test" {
  id         = crosswork_golden_config_templates.test.name
  depends_on = [crosswork_golden_config_jobs.test]
}
`, trigger, device1)
}

func TestAccCwGoldenConfig(t *testing.T) {
	const job = "crosswork_golden_config_jobs.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCwGoldenConfigConfig("{}", "1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(job, "devices.%", "2"),
					resource.TestCheckResourceAttr("data.crosswork_golden_config_jobs.test", "devices.#", "2"),
				),
			},
			{
				Config: testAccCwGoldenConfigConfig(`{ variables = [{ name = "cwtest_gc_desc", value = "device-value" }] }`, "2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(job, "trigger", "2"),
					resource.TestCheckResourceAttr(job, "devices.build-router-2.variables.0.value", "device-value"),
				),
			},
			{
				Config: testAccCwGoldenConfigConfig("{}", "2"),
				Check:  resource.TestCheckResourceAttr(job, "devices.build-router-2.variables.#", "0"),
			},
		},
	})
}

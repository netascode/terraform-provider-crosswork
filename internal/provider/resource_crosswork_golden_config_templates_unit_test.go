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
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestGoldenConfigTemplates(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	cfg := func(extra string) string {
		return s.providerBlock() + fmt.Sprintf(`
resource "crosswork_golden_config_templates" "t" {
  name           = "cwtest-banner"
  vendor         = "Cisco Systems"
  software_type  = "IOS XR"
  product_series = ["Cisco 8000 Series Routers"]
  config         = "banner motd #hello#"
  %s
}

data "crosswork_golden_config_templates" "t" {
  id         = crosswork_golden_config_templates.t.id
  depends_on = [crosswork_golden_config_templates.t]
}

data "crosswork_golden_config_jobs" "t" {
  id         = crosswork_golden_config_templates.t.id
  depends_on = [crosswork_golden_config_templates.t]
}
`, extra)
	}
	const res, ds = "crosswork_golden_config_templates.t", "data.crosswork_golden_config_templates.t"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{
			{ // defaults: merge, native, audit on, no force; nothing is applied
				Config: cfg(""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "mode", "merge"),
					resource.TestCheckResourceAttr(res, "config_type", "native"),
					resource.TestCheckResourceAttr(res, "audit", "true"),
					resource.TestCheckResourceAttr(res, "force_delete", "false"),
					resource.TestCheckResourceAttr(res, "template_id", "1"),
					resource.TestCheckResourceAttr(res, "config_id", "101"),
					resource.TestCheckNoResourceAttr(res, "description"),
					resource.TestCheckResourceAttr(ds, "config", "banner motd #hello#"),
					resource.TestCheckResourceAttr(ds, "template_id", "1"),
					resource.TestCheckResourceAttr(ds, "config_id", "101"),
					resource.TestCheckResourceAttr("data.crosswork_golden_config_jobs.t", "devices.#", "0"),
				),
				ConfigPlanChecks: emptyPlanCheck,
			},
			{
				PreConfig:        s.resetCalls,
				Config:           cfg(`description = "Login banner"`),
				Check:            resource.ComposeTestCheckFunc(s.callsAre("PUT template 1"), resource.TestCheckResourceAttr(res, "description", "Login banner")),
				ConfigPlanChecks: emptyPlanCheck,
			},
			{ // the secret of a variable is in the body and not in the state
				Config: cfg(`variables = [{ name = "k", value_wo = "s3cr3t-banner", value_wo_version = 1 }]`),
				Check:  resource.TestCheckResourceAttr(res, "variables.0.value_wo_version", "1"),
			},
			{
				Config:      cfg(`variables = [{ name = "k", value = "a", value_wo = "b", value_wo_version = 1 }]`),
				ExpectError: regexp.MustCompile(`sets both`),
			},
			{
				Config:      cfg(`variables = [{ name = "k", type = "JSON", value = "{" }]`),
				ExpectError: regexp.MustCompile(`needs a JSON document`),
			},
			{ // force_delete is applied at destroy; an existing application of the template is removed with it
				Config: cfg(`force_delete = true`),
				PreConfig: func() {
					s.mu.Lock()
					s.apps["cwtest-banner-uuid-d1"] = &stubApp{Name: "cwtest-banner-uuid-d1", Template: 1, Device: "d1", UUID: "uuid-d1"}
					s.mu.Unlock()
				},
				Check: resource.TestCheckResourceAttr(res, "force_delete", "true"),
			},
		},
	})
	if len(s.callsWith("DELETE template 1 force=true")) != 1 {
		t.Errorf("the template was not deleted with force: %s", s.callLog())
	}
}

// The controller takes the audit flag only when a template is created, so a change replaces the template.
func TestGoldenConfigTemplatesAuditReplaces(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	cfg := func(audit string) string {
		return s.providerBlock() + fmt.Sprintf(`
resource "crosswork_golden_config_templates" "t" {
  name           = "cwtest-audit"
  vendor         = "Cisco Systems"
  software_type  = "IOS XR"
  product_series = ["Cisco 8000 Series Routers"]
  config         = "banner motd #hello#"
  %s
}
`, audit)
	}
	const res = "crosswork_golden_config_templates.t"
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:           cfg(""),
				Check:            resource.TestCheckResourceAttr(res, "audit", "true"),
				ConfigPlanChecks: emptyPlanCheck,
			},
			{
				Config: cfg(`audit = false`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply:             []plancheck.PlanCheck{plancheck.ExpectResourceAction(res, plancheck.ResourceActionDestroyBeforeCreate)},
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.TestCheckResourceAttr(res, "audit", "false"),
			},
			{
				Config: cfg(`audit = true`),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(res, plancheck.ResourceActionDestroyBeforeCreate)},
				},
				Check: resource.TestCheckResourceAttr(res, "audit", "true"),
			},
		},
	})
}

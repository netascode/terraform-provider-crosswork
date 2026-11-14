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

// Section below is generated&owned by "gen/generator.go". //template:begin imports
import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc
func TestAccCwLoginBanner(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_login_banner.test", "show_message", "true"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_login_banner.test", "user_ack", "false"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_login_banner.test", "message", "Authorized use only"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_login_banner.test", "icon", "BANNER_ICON_IMPORTANT"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_login_banner.test", "title", "Notice"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwLoginBannerConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_login_banner.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{},
		ImportStateId:           "login_banner",
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// End of section. //template:end testAcc

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwLoginBannerConfig_minimum() string {
	config := `resource "crosswork_login_banner" "test" {` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwLoginBannerConfig_all() string {
	config := `resource "crosswork_login_banner" "test" {` + "\n"
	config += `	show_message = true` + "\n"
	config += `	user_ack = false` + "\n"
	config += `	message = "Authorized use only"` + "\n"
	config += `	icon = "BANNER_ICON_IMPORTANT"` + "\n"
	config += `	title = "Notice"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

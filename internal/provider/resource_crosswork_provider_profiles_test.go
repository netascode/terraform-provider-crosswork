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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin testAcc
func TestAccCwProviderProfiles(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_provider_profiles.test", "name", "cwtest-provider-profiles"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_provider_profiles.test", "connectivity.0.protocol", "ROBOT_MSVC_TRANS_SSH"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_provider_profiles.test", "connectivity.0.ip_address", "192.0.2.50"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_provider_profiles.test", "connectivity.0.port", "22"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_provider_profiles.test", "connectivity.0.timeout", "60"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_provider_profiles.test", "supported_models.0.model_prefix", "Cisco-IOS-XR"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_provider_profiles.test", "supported_models.0.model_version", "7.8.2"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwProviderProfilesPrerequisitesConfig + testAccCwProviderProfilesConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_provider_profiles.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    steps,
	})
}

// End of section. //template:end testAcc

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccCwProviderProfilesPrerequisitesConfig = `
resource "crosswork_credential_profiles" "provider_test" {
  profile = "cwtest-provider-profile"
  credentials = [{
    type      = "ROBOT_USERPASS_SSH"
    user_name = "cwtest"
    password  = "Cwtest-Password-1"
  }]
}

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigMinimal
func testAccCwProviderProfilesConfig_minimum() string {
	config := `resource "crosswork_provider_profiles" "test" {` + "\n"
	config += `	name = "cwtest-provider-profiles"` + "\n"
	config += `	family = "ROBOT_PROVIDER_EPNM"` + "\n"
	config += `	profile = crosswork_credential_profiles.provider_test.id` + "\n"
	config += `	connectivity = [{` + "\n"
	config += `	  protocol = "ROBOT_MSVC_TRANS_SSH"` + "\n"
	config += `	  port = 22` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwProviderProfilesConfig_all() string {
	config := `resource "crosswork_provider_profiles" "test" {` + "\n"
	config += `	name = "cwtest-provider-profiles"` + "\n"
	config += `	family = "ROBOT_PROVIDER_EPNM"` + "\n"
	config += `	profile = crosswork_credential_profiles.provider_test.id` + "\n"
	config += `	connectivity = [{` + "\n"
	config += `	  protocol = "ROBOT_MSVC_TRANS_SSH"` + "\n"
	config += `	  ip_address = "192.0.2.50"` + "\n"
	config += `	  port = 22` + "\n"
	config += `	  timeout = 60` + "\n"
	config += `	}]` + "\n"
	config += `	supported_models = [{` + "\n"
	config += `	  model_prefix = "Cisco-IOS-XR"` + "\n"
	config += `	  model_version = "7.8.2"` + "\n"
	config += `	}]` + "\n"
	config += `	properties = { nso_crosslaunch_url = "https://nso.example.invalid/nso" }` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

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
func TestAccCwCredentialProfiles(t *testing.T) {
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "profile", "cwtest-profile"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "credentials.0.type", "ROBOT_USERPASS_SSH"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "credentials.0.user_name", "cwtest-user"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "credentials.0.password_wo_version", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "credentials.0.enable_password", "ROBOT_USERPASSEN_ENABLE"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "credentials.0.enable_password_value_wo_version", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "snmpv2_read_community_wo_version", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "snmpv2_write_community_wo_version", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "snmpv3_user_name", "cwtest-snmp"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "snmpv3_security_level", "SL_AUTH_PRIV"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "snmpv3_auth_type", "AT_HMAC_SHA"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "snmpv3_auth_password_wo_version", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "snmpv3_priv_type", "PT_CFB_AES_128"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_credential_profiles.test", "snmpv3_priv_password_wo_version", "1"))

	var steps []resource.TestStep
	if os.Getenv("SKIP_MINIMUM_TEST") == "" {
		steps = append(steps, resource.TestStep{
			Config: testAccCwCredentialProfilesConfig_minimum(),
		})
	}
	steps = append(steps, resource.TestStep{
		Config: testAccCwCredentialProfilesConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_credential_profiles.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"credentials.0.password", "credentials.0.password_wo", "credentials.0.password_wo_version", "credentials.0.enable_password_value", "credentials.0.enable_password_value_wo", "credentials.0.enable_password_value_wo_version", "snmpv2_read_community", "snmpv2_read_community_wo", "snmpv2_read_community_wo_version", "snmpv2_write_community", "snmpv2_write_community_wo", "snmpv2_write_community_wo_version", "snmpv3_auth_password", "snmpv3_auth_password_wo", "snmpv3_auth_password_wo_version", "snmpv3_priv_password", "snmpv3_priv_password_wo", "snmpv3_priv_password_wo_version"},
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
func testAccCwCredentialProfilesConfig_minimum() string {
	config := `resource "crosswork_credential_profiles" "test" {` + "\n"
	config += `	profile = "cwtest-profile"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwCredentialProfilesConfig_all() string {
	config := `resource "crosswork_credential_profiles" "test" {` + "\n"
	config += `	profile = "cwtest-profile"` + "\n"
	config += `	credentials = [{` + "\n"
	config += `	  type = "ROBOT_USERPASS_SSH"` + "\n"
	config += `	  user_name = "cwtest-user"` + "\n"
	config += `	  password_wo = "Cwtest-Password-1"` + "\n"
	config += `	  password_wo_version = 1` + "\n"
	config += `	  enable_password = "ROBOT_USERPASSEN_ENABLE"` + "\n"
	config += `	  enable_password_value_wo = "Cwtest-Enable-1"` + "\n"
	config += `	  enable_password_value_wo_version = 1` + "\n"
	config += `	}]` + "\n"
	config += `	snmpv2_read_community_wo = "cwtest-read"` + "\n"
	config += `	snmpv2_read_community_wo_version = 1` + "\n"
	config += `	snmpv2_write_community_wo = "cwtest-write"` + "\n"
	config += `	snmpv2_write_community_wo_version = 1` + "\n"
	config += `	snmpv3_user_name = "cwtest-snmp"` + "\n"
	config += `	snmpv3_security_level = "SL_AUTH_PRIV"` + "\n"
	config += `	snmpv3_auth_type = "AT_HMAC_SHA"` + "\n"
	config += `	snmpv3_auth_password_wo = "Cwtest-Auth-1"` + "\n"
	config += `	snmpv3_auth_password_wo_version = 1` + "\n"
	config += `	snmpv3_priv_type = "PT_CFB_AES_128"` + "\n"
	config += `	snmpv3_priv_password_wo = "Cwtest-Priv-1"` + "\n"
	config += `	snmpv3_priv_password_wo_version = 1` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

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
func TestAccCwAAALDAPServers(t *testing.T) {
	if os.Getenv("AAA") == "" {
		t.Skip("skipping test, set environment variable AAA")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "servers.0.name", "cwtest-ldap"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "servers.0.ldap_url", "ldap://192.0.2.30:389"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "servers.0.bind_dn", "cn=admin,dc=example,dc=org"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "servers.0.bind_credential_wo_version", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "servers.0.base_dn", "dc=example,dc=org"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "servers.0.user_filter", "uid={user}"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "servers.0.dn_format", "uid=%s,dc=example,dc=org"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "servers.0.principal_attribute_id", "uid"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "servers.0.policy_id", "crossworkPolicyId"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "servers.0.order", "50"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "servers.0.connect_timeout", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "attribute_servers.0.ldap_url", "ldap://192.0.2.30:389"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "attribute_servers.0.base_dn", "ou=groups,dc=example,dc=org"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "attribute_servers.0.bind_dn", "cn=admin,dc=example,dc=org"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "attribute_servers.0.bind_credential_wo_version", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "attribute_servers.0.user_filter", "memberOf"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "attribute_servers.0.connect_timeout", "5"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_ldap_servers.test", "attribute_servers.0.fail_fast", "false"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwAAALDAPServersConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_aaa_ldap_servers.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"servers.0.bind_credential", "servers.0.bind_credential_wo", "servers.0.bind_credential_wo_version", "attribute_servers.0.bind_credential", "attribute_servers.0.bind_credential_wo", "attribute_servers.0.bind_credential_wo_version"},
		ImportStateId:           "aaa_ldap_servers",
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
func testAccCwAAALDAPServersConfig_minimum() string {
	config := `resource "crosswork_aaa_ldap_servers" "test" {` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwAAALDAPServersConfig_all() string {
	config := `resource "crosswork_aaa_ldap_servers" "test" {` + "\n"
	config += `	servers = [{` + "\n"
	config += `	  name = "cwtest-ldap"` + "\n"
	config += `	  ldap_url = "ldap://192.0.2.30:389"` + "\n"
	config += `	  bind_dn = "cn=admin,dc=example,dc=org"` + "\n"
	config += `	  bind_credential_wo = "AAASecret"` + "\n"
	config += `	  bind_credential_wo_version = 1` + "\n"
	config += `	  base_dn = "dc=example,dc=org"` + "\n"
	config += `	  user_filter = "uid={user}"` + "\n"
	config += `	  dn_format = "uid=%s,dc=example,dc=org"` + "\n"
	config += `	  principal_attribute_id = "uid"` + "\n"
	config += `	  policy_id = "crossworkPolicyId"` + "\n"
	config += `	  order = 50` + "\n"
	config += `	  connect_timeout = 5` + "\n"
	config += `	}]` + "\n"
	config += `	attribute_servers = [{` + "\n"
	config += `	  ldap_url = "ldap://192.0.2.30:389"` + "\n"
	config += `	  base_dn = "ou=groups,dc=example,dc=org"` + "\n"
	config += `	  bind_dn = "cn=admin,dc=example,dc=org"` + "\n"
	config += `	  bind_credential_wo = "AAASecret"` + "\n"
	config += `	  bind_credential_wo_version = 1` + "\n"
	config += `	  user_filter = "memberOf"` + "\n"
	config += `	  connect_timeout = 5` + "\n"
	config += `	  fail_fast = false` + "\n"
	config += `	  attributes = { crossworkPolicyId = "policy_id" }` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

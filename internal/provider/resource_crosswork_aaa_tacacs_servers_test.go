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
func TestAccCwAAATACACSServers(t *testing.T) {
	if os.Getenv("AAA") == "" {
		t.Skip("skipping test, set environment variable AAA")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_tacacs_servers.test", "servers.0.host", "192.0.2.10"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_tacacs_servers.test", "servers.0.port", "49"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_tacacs_servers.test", "servers.0.secret_wo_version", "1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_tacacs_servers.test", "servers.0.secret_format", "ascii"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_tacacs_servers.test", "servers.0.service", "crosswork"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_tacacs_servers.test", "servers.0.priority", "30"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_tacacs_servers.test", "servers.0.retransmit_timeout", "10"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_tacacs_servers.test", "servers.0.retries", "3"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_tacacs_servers.test", "servers.0.auth_type", "pap"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_tacacs_servers.test", "servers.0.interface", "192.0.2.1"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_tacacs_servers.test", "servers.0.policy_id", "role0"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_aaa_tacacs_servers.test", "servers.0.virtual_domain", "deviceAccessGroups"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwAAATACACSServersConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_aaa_tacacs_servers.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"servers.0.secret", "servers.0.secret_wo", "servers.0.secret_wo_version"},
		ImportStateId:           "aaa_tacacs_servers",
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
func testAccCwAAATACACSServersConfig_minimum() string {
	config := `resource "crosswork_aaa_tacacs_servers" "test" {` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwAAATACACSServersConfig_all() string {
	config := `resource "crosswork_aaa_tacacs_servers" "test" {` + "\n"
	config += `	servers = [{` + "\n"
	config += `	  host = "192.0.2.10"` + "\n"
	config += `	  port = 49` + "\n"
	config += `	  secret_wo = "AAASecret"` + "\n"
	config += `	  secret_wo_version = 1` + "\n"
	config += `	  secret_format = "ascii"` + "\n"
	config += `	  service = "crosswork"` + "\n"
	config += `	  priority = 30` + "\n"
	config += `	  retransmit_timeout = 10` + "\n"
	config += `	  retries = 3` + "\n"
	config += `	  auth_type = "pap"` + "\n"
	config += `	  interface = "192.0.2.1"` + "\n"
	config += `	  policy_id = "role0"` + "\n"
	config += `	  virtual_domain = "deviceAccessGroups"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

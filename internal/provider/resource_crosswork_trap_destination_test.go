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
func TestAccCwTrapDestination(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_trap_destination.test", "name", "cwtest-trap"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_trap_destination.test", "address", "192.0.2.30"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_trap_destination.test", "address_type", "IPADDRESS"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_trap_destination.test", "port", "162"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_trap_destination.test", "transport_type", "UDP"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_trap_destination.test", "snmp_version", "VERSION_2C"))
	checks = append(checks, resource.TestCheckResourceAttr("crosswork_trap_destination.test", "community_wo_version", "1"))

	var steps []resource.TestStep
	steps = append(steps, resource.TestStep{
		Config: testAccCwTrapDestinationConfig_all(),
		Check:  resource.ComposeTestCheckFunc(checks...),
	})
	steps = append(steps, resource.TestStep{
		ResourceName:            "crosswork_trap_destination.test",
		ImportState:             true,
		ImportStateVerify:       true,
		ImportStateVerifyIgnore: []string{"community", "community_wo", "community_wo_version", "auth_password", "auth_password_wo", "auth_password_wo_version", "privacy_password", "privacy_password_wo", "privacy_password_wo_version"},
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
func testAccCwTrapDestinationConfig_minimum() string {
	config := `resource "crosswork_trap_destination" "test" {` + "\n"
	config += `	name = "cwtest-trap"` + "\n"
	config += `	address = "192.0.2.30"` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigMinimal

// Section below is generated&owned by "gen/generator.go". //template:begin testAccConfigAll
func testAccCwTrapDestinationConfig_all() string {
	config := `resource "crosswork_trap_destination" "test" {` + "\n"
	config += `	name = "cwtest-trap"` + "\n"
	config += `	address = "192.0.2.30"` + "\n"
	config += `	address_type = "IPADDRESS"` + "\n"
	config += `	port = 162` + "\n"
	config += `	transport_type = "UDP"` + "\n"
	config += `	snmp_version = "VERSION_2C"` + "\n"
	config += `	community_wo = "cwtest-community"` + "\n"
	config += `	community_wo_version = 1` + "\n"
	config += `}` + "\n"
	return config
}

// End of section. //template:end testAccConfigAll

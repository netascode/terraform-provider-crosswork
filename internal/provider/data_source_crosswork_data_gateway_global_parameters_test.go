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

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSource
func TestAccDataSourceCwDataGatewayGlobalParameters(t *testing.T) {
	if os.Getenv("CLUSTER") == "" {
		t.Skip("skipping test, set environment variable CLUSTER")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_gateway_global_parameters.test", "cli_sessions", "3"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_gateway_global_parameters.test", "ssh_session_timeout", "120"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_gateway_global_parameters.test", "snmp_trap_port", "1062"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_gateway_global_parameters.test", "syslog_udp_port", "9514"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_gateway_global_parameters.test", "syslog_tcp_port", "9898"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_gateway_global_parameters.test", "syslog_tls_port", "6514"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_data_gateway_global_parameters.test", "resync_snmpv3_details", "false"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwDataGatewayGlobalParametersConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwDataGatewayGlobalParametersConfig() string {
	config := `resource "crosswork_data_gateway_global_parameters" "test" {` + "\n"
	config += `	cli_sessions = 3` + "\n"
	config += `	ssh_session_timeout = 120` + "\n"
	config += `	snmp_trap_port = 1062` + "\n"
	config += `	syslog_udp_port = 9514` + "\n"
	config += `	syslog_tcp_port = 9898` + "\n"
	config += `	syslog_tls_port = 6514` + "\n"
	config += `	resync_snmpv3_details = false` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_data_gateway_global_parameters" "test" {
			depends_on = [crosswork_data_gateway_global_parameters.test]
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

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
func TestAccDataSourceCwSmartLicensingTransport(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_smart_licensing_transport.test", "mode", "DIRECT"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_smart_licensing_transport.test", "url", "https://smartreceiver.cisco.com/licservice/license"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_smart_licensing_transport.test", "proxy_host", "192.0.2.1"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_smart_licensing_transport.test", "proxy_port", "8080"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwSmartLicensingTransportConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwSmartLicensingTransportConfig() string {
	config := `resource "crosswork_smart_licensing_transport" "test" {` + "\n"
	config += `	mode = "DIRECT"` + "\n"
	config += `	url = "https://smartreceiver.cisco.com/licservice/license"` + "\n"
	config += `	proxy_host = "192.0.2.1"` + "\n"
	config += `	proxy_port = 8080` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_smart_licensing_transport" "test" {
			depends_on = [crosswork_smart_licensing_transport.test]
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

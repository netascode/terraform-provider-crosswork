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
func TestAccDataSourceCwAAASSOServices(t *testing.T) {
	if os.Getenv("AAA") == "" {
		t.Skip("skipping test, set environment variable AAA")
	}
	var checks []resource.TestCheckFunc
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_sso_services.test", "services.0.name", "cwtestsso"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_sso_services.test", "services.0.evaluation_order", "3"))
	checks = append(checks, resource.TestCheckResourceAttr("data.crosswork_aaa_sso_services.test", "services.0.metadata_location", "cwtest-metadata.xml"))
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwAAASSOServicesConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwAAASSOServicesConfig() string {
	config := `resource "crosswork_aaa_sso_services" "test" {` + "\n"
	config += `	services = [{` + "\n"
	config += `	  name = "cwtestsso"` + "\n"
	config += `	  evaluation_order = 3` + "\n"
	config += `	  metadata = "PE1ldGFkYXRhLz4="` + "\n"
	config += `	  metadata_location = "cwtest-metadata.xml"` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_aaa_sso_services" "test" {
			depends_on = [crosswork_aaa_sso_services.test]
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

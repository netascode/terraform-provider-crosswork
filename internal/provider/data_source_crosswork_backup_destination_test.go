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
func TestAccDataSourceCwBackupDestination(t *testing.T) {
	if os.Getenv("SFTP") == "" {
		t.Skip("skipping test, set environment variable SFTP")
	}
	var checks []resource.TestCheckFunc
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwBackupDestinationPrerequisitesConfig + testAccDataSourceCwBackupDestinationConfig(),
				Check:  resource.ComposeTestCheckFunc(checks...),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
const testAccDataSourceCwBackupDestinationPrerequisitesConfig = `
variable "sftp_host" { type = string }
variable "sftp_port" { type = number }
variable "sftp_username" { type = string }
variable "sftp_password" {
  type      = string
  sensitive = true
}
variable "sftp_remote_path" { type = string }

`

// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwBackupDestinationConfig() string {
	config := `resource "crosswork_backup_destination" "test" {` + "\n"
	config += `	host_name = var.sftp_host` + "\n"
	config += `	port = var.sftp_port` + "\n"
	config += `	username = var.sftp_username` + "\n"
	config += `	password_wo = var.sftp_password` + "\n"
	config += `	password_wo_version = 1` + "\n"
	config += `	remote_path = var.sftp_remote_path` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_backup_destination" "test" {
			depends_on = [crosswork_backup_destination.test]
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

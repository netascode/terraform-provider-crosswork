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

// Written by hand: the test of the data source looks a row up by name.
import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// End of section. //template:end imports

// Written by hand: the types come in the order of Crosswork, so the row of the example is looked up by name.
func TestAccDataSourceCwAlarmSeverities(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	const ds = "data.crosswork_alarm_severities.test"
	row := func(state *terraform.State) error {
		attrs := state.RootModule().Resources[ds].Primary.Attributes
		count, _ := strconv.Atoi(attrs["severities.#"])
		if count < 1000 {
			return fmt.Errorf("the data source reads every type of Crosswork, got %d", count)
		}
		for i := 0; i < count; i++ {
			if attrs[fmt.Sprintf("severities.%d.name", i)] == "L2-ELMI-4-AIPC" {
				if got := attrs[fmt.Sprintf("severities.%d.severity", i)]; got != "minor" {
					return fmt.Errorf("severity is %q, want minor", got)
				}
				if got := attrs[fmt.Sprintf("severities.%d.auto_clear_minutes", i)]; got != "1440" {
					return fmt.Errorf("auto_clear_minutes is %q, want 1440", got)
				}
				return nil
			}
		}
		return fmt.Errorf("the data source has no type L2-ELMI-4-AIPC")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCwAlarmSeveritiesConfig(),
				Check:  resource.ComposeTestCheckFunc(row),
			},
		},
	})
}

// End of section. //template:end testAccDataSource

// Section below is generated&owned by "gen/generator.go". //template:begin testPrerequisites
// End of section. //template:end testPrerequisites

// Section below is generated&owned by "gen/generator.go". //template:begin testAccDataSourceConfig
func testAccDataSourceCwAlarmSeveritiesConfig() string {
	config := `resource "crosswork_alarm_severities" "test" {` + "\n"
	config += `	severities = [{` + "\n"
	config += `	  name = "L2-ELMI-4-AIPC"` + "\n"
	config += `	  severity = "minor"` + "\n"
	config += `	  auto_clear_minutes = 1440` + "\n"
	config += `	}]` + "\n"
	config += `}` + "\n"

	config += `
		data "crosswork_alarm_severities" "test" {
			depends_on = [crosswork_alarm_severities.test]
		}
	`
	return config
}

// End of section. //template:end testAccDataSourceConfig

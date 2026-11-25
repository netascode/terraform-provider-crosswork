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

// Changes on a lab: a row of every kind is taken through a change, a drop and an import, and every row is written back to the value it
// had before the test (the destroy of the resource leaves the values, and a revert does not bring back the built-in 1440).

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

var alarmLifecycleNames = []string{
	"AUTO-AMPLI-MISMATCH",          // plain
	"AUTO-LASER-SHUT",              // plain, both fields
	"DWDM-5-G709ALARM description", // a name that is not its event type and has a space
	"L2-ELMI-4-AIPC",               // carries the built-in 1440
}

func alarmLabClient(t *testing.T) *cw.Client {
	t.Helper()
	c, err := registeredClient(os.Getenv("CW_URL"), os.Getenv("CW_USERNAME"), os.Getenv("CW_PASSWORD"), os.Getenv("CW_INSECURE") == "true", 3, "test")
	if err != nil {
		t.Fatalf("login: %s", err)
	}
	return c
}

func alarmLabRows(t *testing.T, c *cw.Client) map[string]gjson.Result {
	t.Helper()
	res, err := c.Get(alarmSeveritiesReadPath)
	if err != nil {
		t.Fatalf("read the alarm types: %s, %s", err, res.String())
	}
	return alarmSeverityIndex(res)
}

// alarmRestore writes the rows back to the snapshot: the severity by `scc` and the auto clear by `aac`, or by a revert for a type that
// had none.
func alarmRestore(t *testing.T, c *cw.Client, snapshot map[string]gjson.Result) {
	for _, name := range alarmLifecycleNames {
		row, found := snapshot[name]
		if !found {
			continue
		}
		write := func(path, body string) {
			if res, err := c.Post(path, body, cw.UseMutex); err != nil {
				t.Errorf("restore %q: %s, %s", name, err, res.String())
			}
		}
		write(alarmSeveritiesPath, alarmWrite{sourceType: "scc", sourceValue: strings.ToLower(row.Get("severity").String()), names: []string{name}}.body())
		if row.Get("revert").Exists() {
			write(alarmSeveritiesPath, alarmWrite{sourceType: "aac", sourceValue: row.Get("revert").String(), names: []string{name}}.body())
		} else {
			write(alarmSeveritiesPath+"/revert", fmt.Sprintf(`{"eventTypes":[%q]}`, name))
		}
	}
}

func alarmLabCheck(t *testing.T, c *cw.Client, name, severity string, minutes int64) resource.TestCheckFunc {
	return func(*terraform.State) error {
		row, found := alarmLabRows(t, c)[name]
		if !found {
			return fmt.Errorf("%q is not on the cluster", name)
		}
		if severity != "" && row.Get("severity").String() != severity {
			return fmt.Errorf("%q has the severity %s, want %s", name, row.Get("severity").String(), severity)
		}
		if minutes != 0 && row.Get("revert").String() != strconv.FormatInt(minutes, 10) {
			return fmt.Errorf("%q has the auto clear %q, want %d", name, row.Get("revert").String(), minutes)
		}
		return nil
	}
}

func TestAccCwAlarmSeveritiesLifecycle(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	c := alarmLabClient(t)
	snapshot := alarmLabRows(t, c)
	for _, name := range alarmLifecycleNames {
		if _, found := snapshot[name]; !found {
			t.Skipf("the cluster has no alarm type %q", name)
		}
	}
	t.Cleanup(func() { alarmRestore(t, c, snapshot) })

	const res = "crosswork_alarm_severities.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: alarmConfig(`    { name = "AUTO-AMPLI-MISMATCH", severity = "minor" },
    { name = "AUTO-LASER-SHUT", severity = "major", auto_clear_minutes = 120 },
    { name = "DWDM-5-G709ALARM description", auto_clear_minutes = 60 },
    { name = "L2-ELMI-4-AIPC", auto_clear_minutes = 300 }`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "severities.#", "4"),
					alarmLabCheck(t, c, "AUTO-AMPLI-MISMATCH", "Minor", 0),
					alarmLabCheck(t, c, "AUTO-LASER-SHUT", "Major", 120),
					alarmLabCheck(t, c, "DWDM-5-G709ALARM description", "", 60),
					alarmLabCheck(t, c, "L2-ELMI-4-AIPC", "", 300),
				),
			},
			{
				// a change, a dropped row and the built-in value written back
				Config: alarmConfig(`    { name = "AUTO-AMPLI-MISMATCH", severity = "warning" },
    { name = "L2-ELMI-4-AIPC", auto_clear_minutes = 1440 }`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "severities.#", "2"),
					alarmLabCheck(t, c, "AUTO-AMPLI-MISMATCH", "Warning", 0),
					alarmLabCheck(t, c, "L2-ELMI-4-AIPC", "", 1440),
					alarmLabCheck(t, c, "AUTO-LASER-SHUT", "Major", 120),
				),
			},
			{
				ResourceName:            res,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"severities"},
				ImportStateId:           "alarm_severities",
			},
		},
	})
}

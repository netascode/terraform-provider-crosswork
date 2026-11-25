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

import (
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const alarmRes = "crosswork_alarm_severities.test"

func alarmConfig(rows string) string {
	return `
resource "crosswork_alarm_severities" "test" {
  severities = [
` + rows + `
  ]
}
`
}

func alarmWant(t *testing.T, s *alarmStub, want ...string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		if got := s.posts(); !reflect.DeepEqual(got, want) {
			t.Errorf("writes: got %q, want %q", got, want)
		}
		return nil
	}
}

func TestAlarmSeveritiesLifecycle(t *testing.T) {
	s := newAlarmStub(t)
	first := `    { name = "AUTO-AMPLI-MISMATCH", severity = "minor" },
    { name = "DWDM-5-G709ALARM description", auto_clear_minutes = 120 },
    { name = "ADHOC_ROGUE_AUTO_CONTAINED", auto_clear_minutes = 60 }`

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// one call per value and one list of names; the severity calls come first; a hidden type is read and its auto clear written
				Config: s.providerBlock() + alarmConfig(first),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(alarmRes, "id", "alarm_severities"),
					resource.TestCheckResourceAttr(alarmRes, "severities.#", "3"),
					resource.TestCheckResourceAttr(alarmRes, "severities.0.severity", "minor"),
					resource.TestCheckNoResourceAttr(alarmRes, "severities.0.auto_clear_minutes"),
					resource.TestCheckNoResourceAttr(alarmRes, "severities.1.severity"),
					resource.TestCheckResourceAttr(alarmRes, "severities.1.auto_clear_minutes", "120"),
					resource.TestCheckNoResourceAttr(alarmRes, "severities.2.severity"),
					resource.TestCheckResourceAttr(alarmRes, "severities.2.auto_clear_minutes", "60"),
					alarmWant(t, s,
						"POST scc minor [AUTO-AMPLI-MISMATCH]",
						"POST aac 60 [ADHOC_ROGUE_AUTO_CONTAINED]",
						"POST aac 120 [DWDM-5-G709ALARM description]",
					),
				),
			},
			{
				// a field that a row leaves out and a type that the configuration does not list are not compared: no difference, no write
				PreConfig: func() {
					s.set("AUTO-AMPLI-MISMATCH", "Minor", "30")
					s.set("DWDM-5-G709ALARM description", "Critical", "120")
					s.set("L2-ELMI-4-AIPC", "Warning", "")
					s.resetCalls()
				},
				Config: s.providerBlock() + alarmConfig(first),
				Check:  alarmWant(t, s),
			},
			{
				// a declared value that was changed outside Terraform is a difference and is written back; nothing else is written
				PreConfig: func() { s.set("AUTO-AMPLI-MISMATCH", "Warning", "30"); s.resetCalls() },
				Config:    s.providerBlock() + alarmConfig(first),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(alarmRes, "severities.0.severity", "minor"),
					alarmWant(t, s, "POST scc minor [AUTO-AMPLI-MISMATCH]"),
				),
			},
			{
				// a changed value of one row and an auto clear value of another; a value that already matches is not written
				PreConfig: func() { s.resetCalls() },
				Config: s.providerBlock() + alarmConfig(`    { name = "AUTO-AMPLI-MISMATCH", severity = "major" },
    { name = "Link down(admin down)", severity = "major", auto_clear_minutes = 55 },
    { name = "DWDM-5-G709ALARM description", auto_clear_minutes = 120 },
    { name = "ADHOC_ROGUE_AUTO_CONTAINED", auto_clear_minutes = 60 }`),
				Check: alarmWant(t, s, "POST scc major [AUTO-AMPLI-MISMATCH]", "POST aac 55 [Link down(admin down)]"),
			},
			{
				// a row that is dropped from the configuration leaves its values on the controller
				PreConfig: func() { s.resetCalls() },
				Config:    s.providerBlock() + alarmConfig(`    { name = "Link down(admin down)", severity = "major" }`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(alarmRes, "severities.#", "1"),
					alarmWant(t, s),
					func(*terraform.State) error {
						if got := s.row("ADHOC_ROGUE_AUTO_CONTAINED"); got.severity != "Major" || got.revert != "60" {
							t.Errorf("a dropped row keeps its values, got %+v", got)
						}
						return nil
					},
				),
			},
			{
				// import reads no row: the configuration names them
				ResourceName:  alarmRes,
				ImportState:   true,
				ImportStateId: "alarm_severities",
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if len(states) != 1 || states[0].Attributes["severities.#"] != "0" {
						t.Errorf("import has no rows, got %v", states)
					}
					return nil
				},
			},
		},
	})
	// the destroy at the end of the test changes nothing
	if got := s.row("Link down(admin down)"); got.severity != "Major" || got.revert != "55" {
		t.Errorf("destroy leaves the values, got %+v", got)
	}
	if got := s.row("ADHOC_ROGUE_AUTO_CONTAINED"); got.severity != "Major" || got.revert != "60" {
		t.Errorf("destroy leaves the values, got %+v", got)
	}
}

func TestAlarmSeveritiesUnknownNameIsRefusedBeforeAnyWrite(t *testing.T) {
	s := newAlarmStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: s.providerBlock() + alarmConfig(`    { name = "AUTO-AMPLI-MISMATCH", severity = "minor" },
    { name = "cwtest-unknown-b", severity = "minor" },
    { name = "cwtest-unknown-a", auto_clear_minutes = 60 },
    { name = "LINK_DOWN_ADMIN", severity = "minor" }`),
			ExpectError: regexp.MustCompile(`(?s)unknown\s+alarm\s+or\s+event\s+types:\s+LINK_DOWN_ADMIN,\s+cwtest-unknown-a,\s+cwtest-unknown-b`),
		}},
	})
	if posts := s.posts(); len(posts) != 0 {
		t.Errorf("no write when a name is unknown, got %q", posts)
	}
	if got := s.row("AUTO-AMPLI-MISMATCH"); got.severity != "Information" {
		t.Errorf("a valid row next to an unknown one is not applied, got %+v", got)
	}
}

func TestAlarmSeveritiesControllerRefusalShowsItsText(t *testing.T) {
	s := newAlarmStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      s.providerBlock() + alarmConfig(`    { name = "AUTO-AMPLI-MISMATCH", auto_clear_minutes = 58 }`),
			ExpectError: regexp.MustCompile(`(?s)the auto clear minutes.*multiples of 60`),
		}},
	})
}

func TestAlarmSeveritiesSendsNothingWhenTheValuesAlreadyMatch(t *testing.T) {
	s := newAlarmStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			// the severity is read capitalised and compared as written; 1440 is the built-in value
			Config: s.providerBlock() + alarmConfig(`    { name = "L2-ELMI-4-AIPC", severity = "minor", auto_clear_minutes = 1440 }`),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(alarmRes, "severities.0.severity", "minor"),
				resource.TestCheckResourceAttr(alarmRes, "severities.0.auto_clear_minutes", "1440"),
				alarmWant(t, s),
			),
		}},
	})
}

func TestAlarmSeveritiesReadsHiddenTypesWithTheCompleteQuery(t *testing.T) {
	s := newAlarmStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: s.providerBlock() + alarmConfig(`    { name = "ADHOC_ROGUE_AUTO_CONTAINED", severity = "major" }`),
			Check: func(*terraform.State) error {
				for _, call := range s.callLog() {
					if call == "GET" {
						t.Errorf("every read asks for the complete list, got %q", s.callLog())
					}
				}
				return nil
			},
		}},
	})
}

func TestAlarmSeveritiesRefusesASeverityOfAHiddenTypeBeforeAnyWrite(t *testing.T) {
	s := newAlarmStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: s.providerBlock() + alarmConfig(`    { name = "AUTO-AMPLI-MISMATCH", severity = "minor" },
    { name = "ADHOC_ROGUE_AUTO_CONTAINED", severity = "critical", auto_clear_minutes = 60 }`),
			ExpectError: regexp.MustCompile(`(?s)severity\s+of\s+ADHOC_ROGUE_AUTO_CONTAINED\s+cannot\s+be\s+changed`),
		}},
	})
	if posts := s.posts(); len(posts) != 0 {
		t.Errorf("no write when a severity is refused, got %q", posts)
	}
	// the severity that the type already has is not a write, so it is not refused
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: s.providerBlock() + alarmConfig(`    { name = "ADHOC_ROGUE_AUTO_CONTAINED", severity = "major", auto_clear_minutes = 60 }`),
			Check:  alarmWant(t, s, "POST aac 60 [ADHOC_ROGUE_AUTO_CONTAINED]"),
		}},
	})
}

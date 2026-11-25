// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License Version 2.0 (the "License"),
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://mozilla.org/MPL/2.0/
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
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const deviceAlarmManagerRes = "crosswork_device_alarm_manager.test"

func deviceAlarmManagerConfig(models map[string]bool) string {
	names := make([]string, 0, len(models))
	for name := range models {
		names = append(names, name)
	}
	sort.Strings(names)
	var b strings.Builder
	for _, name := range names {
		fmt.Fprintf(&b, "    %q = %t\n", name, models[name])
	}
	return fmt.Sprintf("\nresource \"crosswork_device_alarm_manager\" \"test\" {\n  models = {\n%s  }\n}\n", b.String())
}

func deviceAlarmManagerCheckHeld(s *deviceAlarmManagerStub, enabled ...string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		want := map[string]bool{}
		for _, name := range enabled {
			want[name] = true
		}
		for name, on := range s.state() {
			if on != want[name] {
				return fmt.Errorf("the controller holds %v, want only %v enabled", s.state(), enabled)
			}
		}
		return nil
	}
}

func deviceAlarmManagerCheckPosts(s *deviceAlarmManagerStub, want int) resource.TestCheckFunc {
	return func(*terraform.State) error {
		got := 0
		for _, call := range s.callLog() {
			if strings.HasPrefix(call, "POST") {
				got++
			}
		}
		if got != want {
			return fmt.Errorf("%d writes so far, want %d: %v", got, want, s.callLog())
		}
		return nil
	}
}

const (
	dam9010 = "Cisco ASR 9010 Router"
	dam9006 = "Cisco ASR 9006 Router"
	dam5501 = "Cisco NCS 5501"
	dam8201 = "Cisco 8201 Router"
)

func TestDeviceAlarmManagerLifecycle(t *testing.T) {
	s := newDeviceAlarmManagerStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             deviceAlarmManagerCheckHeld(s),
		Steps: []resource.TestStep{
			{
				// create: only the listed models are sent, with the prefix
				Config: s.providerBlock() + deviceAlarmManagerConfig(map[string]bool{dam9010: true, dam9006: false}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(deviceAlarmManagerRes, "id", "device_alarm_manager"),
					resource.TestCheckResourceAttr(deviceAlarmManagerRes, "models.%", "2"),
					deviceAlarmManagerCheckHeld(s, dam9010),
					func(*terraform.State) error {
						if got := s.callLog()[:3]; got[0] != "GET" || got[2] != "GET" || !strings.HasPrefix(got[1], "POST ") || !strings.Contains(got[1], `"alarmManager/Cisco ASR 9010 Router":true`) {
							return fmt.Errorf("the calls were %v, want a read, a write and a read", got)
						}
						return nil
					},
				),
			},
			{
				// a model that is not listed is changed outside Terraform: no difference, and it is left alone
				PreConfig: func() { s.set(dam5501, true) },
				Config:    s.providerBlock() + deviceAlarmManagerConfig(map[string]bool{dam9010: true, dam9006: false}),
				Check:     deviceAlarmManagerCheckHeld(s, dam9010, dam5501),
			},
			{
				// a listed model is changed outside Terraform: one apply restores it
				PreConfig: func() { s.set(dam9010, false) },
				Config:    s.providerBlock() + deviceAlarmManagerConfig(map[string]bool{dam9010: true, dam9006: false}),
				Check:     deviceAlarmManagerCheckHeld(s, dam9010, dam5501),
			},
			{
				// a model that leaves the configuration is switched off; one that was never listed keeps its state
				Config: s.providerBlock() + deviceAlarmManagerConfig(map[string]bool{dam9006: true, dam8201: false}),
				Check: resource.ComposeTestCheckFunc(
					deviceAlarmManagerCheckHeld(s, dam9006, dam5501),
					resource.TestCheckResourceAttr(deviceAlarmManagerRes, "models.%", "2"),
				),
			},
			{
				// an unknown name is refused before any write
				Config:      s.providerBlock() + deviceAlarmManagerConfig(map[string]bool{"Cisco asr 9010 Router": true}),
				ExpectError: regexp.MustCompile(`the controller has no device model`),
			},
			{
				Config: s.providerBlock() + deviceAlarmManagerConfig(map[string]bool{dam9006: true, dam8201: false}),
				Check: resource.ComposeTestCheckFunc(
					deviceAlarmManagerCheckHeld(s, dam9006, dam5501),
					func(*terraform.State) error {
						for _, call := range s.callLog() {
							if strings.Contains(call, "asr 9010") {
								return fmt.Errorf("the refused name reached the controller: %v", s.callLog())
							}
						}
						return nil
					},
				),
			},
			{
				// import reads every model of the controller
				ResourceName:            deviceAlarmManagerRes,
				ImportState:             true,
				ImportStateId:           "device_alarm_manager",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"models"},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					if got := states[0].Attributes["models.%"]; got != fmt.Sprint(len(deviceAlarmManagerStubModels)) {
						return fmt.Errorf("the import holds %s models, want %d", got, len(deviceAlarmManagerStubModels))
					}
					if states[0].Attributes["models."+dam5501] != "true" || states[0].Attributes["models."+dam9010] != "false" {
						return fmt.Errorf("the import does not hold the state of the controller: %v", states[0].Attributes)
					}
					return nil
				},
			},
		},
	})
}

func TestDeviceAlarmManagerImportThenAShorterConfigurationSwitchesTheOthersOff(t *testing.T) {
	s := newDeviceAlarmManagerStub(t)
	s.set(dam9010, true)
	s.set(dam5501, true)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             deviceAlarmManagerCheckHeld(s),
		Steps: []resource.TestStep{
			{
				Config:             s.providerBlock() + deviceAlarmManagerConfig(map[string]bool{dam9010: true}),
				ResourceName:       deviceAlarmManagerRes,
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateId:      "device_alarm_manager",
			},
			{
				// the state holds every model, so the shorter list plans the removal of the others and the apply switches them off
				Config: s.providerBlock() + deviceAlarmManagerConfig(map[string]bool{dam9010: true}),
				Check:  deviceAlarmManagerCheckHeld(s, dam9010),
			},
		},
	})
}

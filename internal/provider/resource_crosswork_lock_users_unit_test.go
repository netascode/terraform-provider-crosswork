// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
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
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func (s *userStub) lockConfig() string {
	return s.config("Test", "Secret-1", 1) + `
resource "crosswork_lock_users" "test" {
  username = crosswork_users.test.username
}
`
}

// The lock is a separate resource: creating it locks the user, destroying it unlocks the user
// first, and the user itself is then deleted.
func TestLockUsersLifecycle(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newUserStub(t)
	locked := func(want bool) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if got := s.user("cwtest-unit")["Locked"] == true; got != want {
				return fmt.Errorf("user locked = %v, want %v", got, want)
			}
			return nil
		}
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             locked(false),
		Steps: []resource.TestStep{
			{
				Config: s.lockConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("crosswork_lock_users.test", "id", "cwtest-unit"),
					locked(true),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
			{
				ResourceName:  "crosswork_lock_users.test",
				ImportState:   true,
				ImportStateId: "cwtest-unit",
			},
		},
	})

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.locks != 2 {
		t.Errorf("stub saw %d successful lock and unlock calls; want 2", s.locks)
	}
	var lockBody, unlockBody string
	for _, b := range s.bodies {
		if strings.Contains(b, "/lockUsers") {
			lockBody = b
		}
		if strings.Contains(b, "/unlockUsers") {
			unlockBody = b
		}
	}
	if !strings.Contains(lockBody, `{"cwtest-unit": {"userName":"cwtest-unit"`) {
		t.Errorf("the lock body is not wrapped under the user name: %q", lockBody)
	}
	if !strings.Contains(unlockBody, `{"cwtest-unit": {"userName": "cwtest-unit"}}`) {
		t.Errorf("the unlock body is wrong: %q", unlockBody)
	}
}

func TestLockUsersAlreadyLocked(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newUserStub(t)
	s.users["cwtest-unit"] = map[string]any{"Username": "cwtest-unit", "PolicyId": "admin", "FirstName": "A", "LastName": "B", "Password": "x", "Locked": true}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: s.providerBlock() + `
resource "crosswork_lock_users" "test" {
  username = "cwtest-unit"
}
`,
			ExpectError: regexp.MustCompile(`At least one active user`),
		}},
	})
}

// A user unlocked on the controller is locked again by the next apply: the read does not find it among the locked ones.
func TestLockUsersDrift(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newUserStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: s.lockConfig()},
			{
				PreConfig: func() {
					s.mu.Lock()
					s.users["cwtest-unit"]["Locked"] = false
					s.mu.Unlock()
				},
				Config: s.lockConfig(),
				Check: func(*terraform.State) error {
					if s.user("cwtest-unit")["Locked"] != true {
						return fmt.Errorf("the user was not locked again")
					}
					return nil
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
		},
	})
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.locks != 3 {
		t.Errorf("stub saw %d successful lock and unlock calls; want 3", s.locks)
	}
}

func TestLockUsersDataSource(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newUserStub(t)
	s.users["cwtest-locked"] = map[string]any{"Username": "cwtest-locked", "PolicyId": "admin", "FirstName": "A", "LastName": "B", "Password": "x", "Locked": true}
	s.users["cwtest-open"] = map[string]any{"Username": "cwtest-open", "PolicyId": "admin", "FirstName": "A", "LastName": "B", "Password": "x"}
	config := func(id string) string {
		return s.providerBlock() + fmt.Sprintf("\ndata \"crosswork_lock_users\" \"test\" {\n  id = %q\n}\n", id)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config("cwtest-locked"),
				Check:  resource.TestCheckResourceAttr("data.crosswork_lock_users.test", "username", "cwtest-locked"),
			},
			{Config: config("cwtest-open"), ExpectError: regexp.MustCompile(`No object found with id: cwtest-open`)},
			{Config: config("cwtest-missing"), ExpectError: regexp.MustCompile(`No object found with id: cwtest-missing`)},
		},
	})
}

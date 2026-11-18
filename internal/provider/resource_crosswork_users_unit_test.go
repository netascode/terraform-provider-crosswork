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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// Writes need Content-Type "application/json; charset=UTF-8" (405 otherwise); the create body is wrapped under the user name,
// the update body is flat and repeats Username; a missing user is HTTP 500; a user created without device access groups gets the default group.
type userStub struct {
	*ssoStub

	mu      sync.Mutex
	users   map[string]map[string]any
	bodies  []string
	creates int
	updates int
	deletes int
	locks   int

	// deadlockOnOverlap makes a write that overlaps another write fail like Crosswork does (500, pq: deadlock detected).
	deadlockOnOverlap bool
	inflight          int32
	overlaps          int32
}

func newUserStub(t *testing.T) *userStub {
	t.Helper()
	s := &userStub{ssoStub: newSSOStub(t), users: map[string]map[string]any{}}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.HasPrefix(r.URL.Path, "/crosswork/aaa/v1/user"):
			s.serve(w, r)
			return
		case r.URL.Path == "/crosswork/aaa/v1/lockUsers" || r.URL.Path == "/crosswork/aaa/v1/unlockUsers":
			s.lock(w, r)
			return
		}
		sso.ServeHTTP(w, r)
	})
	t.Cleanup(Shutdown)
	return s
}

func userFail(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, `{"error":%q,"code":%d}`, msg, code)
}

func (s *userStub) serve(w http.ResponseWriter, r *http.Request) {
	if s.deadlockOnOverlap && r.Method != http.MethodGet {
		atomic.AddInt32(&s.inflight, 1)
		defer atomic.AddInt32(&s.inflight, -1)
		time.Sleep(100 * time.Millisecond) // long enough for parallel writes to overlap
		if atomic.LoadInt32(&s.inflight) > 1 {
			atomic.AddInt32(&s.overlaps, 1)
			userFail(w, 500, "Failed to create user: pq: deadlock detected")
			return
		}
	}
	raw, _ := io.ReadAll(r.Body)
	s.mu.Lock()
	defer s.mu.Unlock()
	name := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/crosswork/aaa/v1/user"), "/")

	if r.Method != http.MethodGet {
		s.bodies = append(s.bodies, r.Method+" "+r.URL.Path+" "+string(raw))
		if r.Method != http.MethodDelete && r.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	}
	_, exists := s.users[name]

	switch {
	case r.Method == http.MethodGet && name == "":
		all := map[string]any{}
		for n, u := range s.users {
			status := "Active"
			if u["Locked"] == true {
				status = "Locked"
			}
			all[n] = map[string]any{"Username": n, "Password": cwtestMask, "PolicyId": u["PolicyId"], "FirstName": u["FirstName"], "LastName": u["LastName"], "Status": status}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(all)
	case r.Method == http.MethodGet && name != "":
		if !exists {
			userFail(w, 500, "Invalid Username")
			return
		}
		u := s.users[name]
		status := "Active"
		if u["Locked"] == true {
			status = "Locked"
		}
		out := map[string]any{"Username": name, "Password": cwtestMask, "PolicyId": u["PolicyId"], "FirstName": u["FirstName"],
			"LastName": u["LastName"], "Status": status, "DeviceAccessGroups": []any{map[string]any{"Uuid": "default-uuid", "DomainName": "ALL-ACCESS"}}}
		if g, ok := u["DeviceAccessGroups"]; ok {
			out["DeviceAccessGroups"] = g
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(out)
	case r.Method == http.MethodPost && name == "":
		s.creates++
		var wrapped map[string]map[string]any
		if err := json.Unmarshal(raw, &wrapped); err != nil || len(wrapped) != 1 {
			userFail(w, 422, "Failed to create user: json: cannot unmarshal string into Go value of type user.Policy")
			return
		}
		for key, u := range wrapped {
			switch {
			case u["Username"] != key:
				userFail(w, 500, "Failed to create user: Invalid Username")
			case u["FirstName"] == nil || u["FirstName"] == "":
				userFail(w, 500, "Failed to create user: First Name cannot be empty.")
			case u["Password"] == nil:
				userFail(w, 500, "Failed to create user: Password cannot be empty.")
			case exists || s.users[key] != nil:
				userFail(w, 500, "Failed to create user: cannot add user. User already exist")
			default:
				s.users[key] = u
				w.WriteHeader(http.StatusCreated)
			}
		}
	case r.Method == http.MethodPut && name != "":
		s.updates++
		var u map[string]any
		_ = json.Unmarshal(raw, &u)
		switch {
		case !exists:
			userFail(w, 500, "Failed to update user: cannot update user as it does not exist")
		case u["Username"] != name:
			userFail(w, 500, "Failed to update user: Invalid user name.")
		default:
			if _, ok := u["Password"]; !ok {
				u["Password"] = s.users[name]["Password"]
			}
			s.users[name] = u
			w.WriteHeader(http.StatusNoContent)
		}
	case r.Method == http.MethodDelete && name != "":
		s.deletes++
		if !exists {
			userFail(w, 500, "Failed to delete user: cannot delete user as it does not exist")
			return
		}
		delete(s.users, name)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// An unknown user, locking a locked user and unlocking an active one are 500s.
func (s *userStub) lock(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bodies = append(s.bodies, r.Method+" "+r.URL.Path+" "+string(raw))
	if r.Method != http.MethodPut || r.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var wrapped map[string]map[string]any
	if err := json.Unmarshal(raw, &wrapped); err != nil || len(wrapped) != 1 {
		userFail(w, 422, "json: cannot unmarshal string into Go value of type user.Policy")
		return
	}
	locking := r.URL.Path == "/crosswork/aaa/v1/lockUsers"
	for key, body := range wrapped {
		u, exists := s.users[key]
		switch {
		case !exists || body["userName"] != key:
			userFail(w, 500, "Invalid Username: "+key)
		case locking && u["Locked"] == true:
			userFail(w, 500, "Cannot Lock User(s): At least one active user should be provided to lock users.")
		case !locking && u["Locked"] != true:
			userFail(w, 500, "Cannot Unlock User(s): At least one locked user should be provided to unlock users.")
		default:
			u["Locked"] = locking
			s.locks++
			w.WriteHeader(http.StatusNoContent)
		}
	}
}

func (s *userStub) user(name string) map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.users[name]
}

func (s *userStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

func (s *userStub) config(first, secret string, version int) string {
	return s.providerBlock() + fmt.Sprintf(`
resource "crosswork_users" "test" {
  username            = "cwtest-unit"
  policy_id           = "admin"
  first_name          = %q
  last_name           = "User"
  password_wo         = %q
  password_wo_version = %d
  device_access_groups = [{
    uuid        = "group-1"
    domain_name = "ALL-ACCESS"
  }]
}
`, first, secret, version)
}

func TestUsersLifecycle(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newUserStub(t)
	const res = "crosswork_users.test"
	emptyPlan := resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(*terraform.State) error {
			if s.user("cwtest-unit") != nil {
				return fmt.Errorf("user cwtest-unit still exists on the stub")
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: s.config("Test", "Secret-1", 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", "cwtest-unit"),
					resource.TestCheckResourceAttr(res, "first_name", "Test"),
					resource.TestCheckResourceAttr(res, "device_access_groups.0.uuid", "group-1"),
					noSecretInState("Secret-1", cwtestMask),
					func(*terraform.State) error {
						if u := s.user("cwtest-unit"); u["Password"] != "Secret-1" || u["FirstName"] != "Test" {
							return fmt.Errorf("the stub holds %v", u)
						}
						return nil
					},
				),
				ConfigPlanChecks: emptyPlan,
			},
			{
				Config: s.config("Changed", "Secret-2", 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "first_name", "Changed"),
					noSecretInState("Secret-1", "Secret-2", cwtestMask),
					func(*terraform.State) error {
						if u := s.user("cwtest-unit"); u["Password"] != "Secret-2" || u["FirstName"] != "Changed" {
							return fmt.Errorf("the stub holds %v", u)
						}
						return nil
					},
				),
				ConfigPlanChecks: emptyPlan,
			},
			{ // import by user name; the masked password is not read back
				ResourceName:            res,
				ImportState:             true,
				ImportStateId:           "cwtest-unit",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password_wo_version"},
			},
			{ // deleted on the controller: the read treats the 500 as gone, so the apply recreates the user
				PreConfig: func() {
					s.mu.Lock()
					delete(s.users, "cwtest-unit")
					s.mu.Unlock()
				},
				Config:           s.config("Changed", "Secret-2", 2),
				ConfigPlanChecks: emptyPlan,
			},
		},
	})

	s.mu.Lock()
	defer s.mu.Unlock()
	if s.creates != 2 || s.updates != 1 || s.deletes != 1 {
		t.Errorf("stub saw %d creates, %d updates and %d deletes; want 2, 1 and 1 (the second create follows the out of band delete)", s.creates, s.updates, s.deletes)
	}
	if len(s.bodies) < 1 || !strings.HasPrefix(s.bodies[0], `POST /crosswork/aaa/v1/user {"cwtest-unit": {`) {
		t.Errorf("the create body is not wrapped under the user name: %v", s.bodies)
	}
	if len(s.bodies) < 2 || !strings.HasPrefix(s.bodies[1], `PUT /crosswork/aaa/v1/user/cwtest-unit {"Username":"cwtest-unit"`) {
		t.Errorf("the update body is not flat: %v", s.bodies)
	}
}

func TestUsersCreateExisting(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newUserStub(t)
	s.users["cwtest-unit"] = map[string]any{"Username": "cwtest-unit", "PolicyId": "admin", "FirstName": "Other", "LastName": "X", "Password": "x"}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      s.config("Test", "Secret-1", 1),
			ExpectError: regexp.MustCompile(`already exist`),
		}},
	})
	if got := s.user("cwtest-unit"); got["FirstName"] != "Other" {
		t.Errorf("the existing user was changed: %v", got)
	}
}

// The users definition has mutex: true, so parallel creates and deletes of many users are sent one at a time; without it
// Crosswork answers some of them with "pq: deadlock detected".
func TestUsersWritesAreSerialized(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newUserStub(t)
	s.deadlockOnOverlap = true

	cfg := s.providerBlock() + `
resource "crosswork_users" "test" {
  count               = 8
  username            = "cwtest-par-${count.index}"
  policy_id           = "admin"
  first_name          = "Par"
  last_name           = "User"
  password_wo         = "Secret-1"
  password_wo_version = 1
  device_access_groups = [{
    uuid        = "group-1"
    domain_name = "ALL-ACCESS"
  }]
}
`
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps:                    []resource.TestStep{{Config: cfg}},
	})
	if n := atomic.LoadInt32(&s.overlaps); n != 0 {
		t.Errorf("%d writes overlapped; the users writes are not serialized", n)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.creates != 8 || s.deletes != 8 {
		t.Errorf("stub saw %d creates and %d deletes; want 8 and 8", s.creates, s.deletes)
	}
}

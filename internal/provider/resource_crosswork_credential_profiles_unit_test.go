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
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const cwtestMask = "******"

// Every failure is HTTP 200 with a failed job; a missing profile has no "data" member; passwords come back masked.
type credentialStub struct {
	*ssoStub

	mu       sync.Mutex
	profiles map[string][]map[string]any
	creates  int
	updates  int
	deletes  int
	queries  int
	vanish   int // when > 0, the profile vanishes after this many more queries
	vanishOf string
	bodies   []string
}

func newCredentialStub(t *testing.T) *credentialStub {
	t.Helper()
	s := &credentialStub{ssoStub: newSSOStub(t), profiles: map[string][]map[string]any{}}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/crosswork/inventory/v1/credentials") {
			s.serve(w, r)
			return
		}
		sso.ServeHTTP(w, r)
	})
	// Registered after the SSO stub, so it runs first and logs out while the server is still up.
	t.Cleanup(Shutdown)
	return s
}

func jobReply(w http.ResponseWriter, failure string) {
	w.Header().Set("Content-Type", "application/json")
	if failure != "" {
		fmt.Fprintf(w, `{"job_id":"job-1","state":"JOB_FAILED","error":%q}`, failure)
		return
	}
	fmt.Fprint(w, `{"job_id":"job-1","state":"JOB_COMPLETED","impacted_count":1}`)
}

type credentialBody struct {
	Data []struct {
		Profile  string           `json:"profile"`
		UserPass []map[string]any `json:"user_pass"`
	} `json:"data"`
	Filter struct {
		Profile string `json:"profile"`
	} `json:"filter"`
}

func (s *credentialStub) serve(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	var body credentialBody
	_ = json.Unmarshal(raw, &body)

	s.mu.Lock()
	defer s.mu.Unlock()

	if r.Method == http.MethodPost && r.URL.Path == "/crosswork/inventory/v1/credentials/query" {
		s.queries++
		w.Header().Set("Content-Type", "application/json")
		list, ok := s.profiles[body.Filter.Profile]
		if !ok {
			fmt.Fprint(w, `{"total_count":0}`)
		} else {
			masked := []map[string]any{}
			for _, up := range list {
				m := map[string]any{"type": up["type"], "user_name": up["user_name"], "password": cwtestMask}
				if _, has := up["enable_password_data"]; has {
					m["enable_password_data"] = cwtestMask
				}
				masked = append(masked, m)
			}
			out, _ := json.Marshal(map[string]any{
				"data":        []any{map[string]any{"profile": body.Filter.Profile, "user_pass": masked}},
				"total_count": 1,
			})
			_, _ = w.Write(out)
		}
		if s.vanish > 0 {
			if s.vanish--; s.vanish == 0 {
				delete(s.profiles, s.vanishOf)
			}
		}
		return
	}

	s.bodies = append(s.bodies, r.Method+" "+string(raw))
	if len(body.Data) != 1 {
		jobReply(w, "exactly one profile expected")
		return
	}
	name, list := body.Data[0].Profile, body.Data[0].UserPass
	_, exists := s.profiles[name]
	switch r.Method {
	case http.MethodPost:
		s.creates++
		if exists {
			jobReply(w, fmt.Sprintf("credential profile %s already exists", name))
			return
		}
		s.profiles[name] = list
	case http.MethodPut:
		s.updates++
		if !exists {
			jobReply(w, fmt.Sprintf("credential profile %s does not exist", name))
			return
		}
		s.profiles[name] = list
	case http.MethodDelete:
		s.deletes++
		if !exists {
			jobReply(w, "None or more than one credentials was found")
			return
		}
		delete(s.profiles, name)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	jobReply(w, "")
}

func (s *credentialStub) userPass(profile string) []map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.profiles[profile]
}

func (s *credentialStub) counts() (creates, updates, deletes int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.creates, s.updates, s.deletes
}

func (s *credentialStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

const cwtestProfile = "cwtest-unit"

func (s *credentialStub) profileConfig(entries ...string) string {
	return s.providerBlock() + fmt.Sprintf(`
resource "crosswork_credential_profiles" "test" {
  profile     = %q
  credentials = [%s]
}
`, cwtestProfile, strings.Join(entries, ", "))
}

func sshEntry(user, secret string, version int) string {
	return fmt.Sprintf(`{
    type                = "ROBOT_USERPASS_SSH"
    user_name           = %q
    password_wo         = %q
    password_wo_version = %d
  }`, user, secret, version)
}

func netconfEntry(user, secret string, version int) string {
	return fmt.Sprintf(`{
    type                = "ROBOT_USERPASS_NETCONF"
    user_name           = %q
    password_wo         = %q
    password_wo_version = %d
  }`, user, secret, version)
}

func noSecretInState(secrets ...string) resource.TestCheckFunc {
	return func(st *terraform.State) error {
		for _, m := range st.Modules {
			for name, rs := range m.Resources {
				for k, v := range rs.Primary.Attributes {
					for _, secret := range secrets {
						if strings.Contains(v, secret) {
							return fmt.Errorf("%s: attribute %s holds a secret", name, k)
						}
					}
				}
			}
		}
		return nil
	}
}

func (s *credentialStub) stubHas(want ...[3]string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		got := s.userPass(cwtestProfile)
		if len(got) != len(want) {
			return fmt.Errorf("stub holds %d user_pass entries, want %d: %v", len(got), len(want), got)
		}
		for i, w := range want {
			if got[i]["type"] != w[0] || got[i]["user_name"] != w[1] || got[i]["password"] != w[2] {
				return fmt.Errorf("stub entry %d is %v, want %v", i, got[i], w)
			}
		}
		return nil
	}
}

func (s *credentialStub) checkGone() resource.TestCheckFunc {
	return func(*terraform.State) error {
		if got := s.userPass(cwtestProfile); got != nil {
			return fmt.Errorf("profile %s still exists on the stub", cwtestProfile)
		}
		return nil
	}
}

func TestCredentialProfilesLifecycle(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newCredentialStub(t)
	const res = "crosswork_credential_profiles.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             func(*terraform.State) error { return s.checkGone()(nil) },
		Steps: []resource.TestStep{
			{
				Config: s.profileConfig(sshEntry("user-a", "Secret-A", 1)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", cwtestProfile),
					resource.TestCheckResourceAttr(res, "credentials.#", "1"),
					resource.TestCheckResourceAttr(res, "credentials.0.type", "ROBOT_USERPASS_SSH"),
					resource.TestCheckResourceAttr(res, "credentials.0.user_name", "user-a"),
					s.stubHas([3]string{"ROBOT_USERPASS_SSH", "user-a", "Secret-A"}),
					noSecretInState("Secret-A", cwtestMask),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
			{
				Config:   s.profileConfig(sshEntry("user-a", "Secret-A", 1)),
				PlanOnly: true,
			},
			{ // update: the credentials list is replaced, not merged
				Config: s.profileConfig(netconfEntry("user-b", "Secret-B", 1), sshEntry("user-a", "Secret-A2", 2)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "credentials.#", "2"),
					resource.TestCheckResourceAttr(res, "credentials.0.type", "ROBOT_USERPASS_NETCONF"),
					s.stubHas([3]string{"ROBOT_USERPASS_NETCONF", "user-b", "Secret-B"}, [3]string{"ROBOT_USERPASS_SSH", "user-a", "Secret-A2"}),
					noSecretInState("Secret-A", "Secret-A2", "Secret-B", cwtestMask),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
			},
			{ // import by profile name: the masked passwords are not read back into the state
				ResourceName:      res,
				ImportState:       true,
				ImportStateId:     cwtestProfile,
				ImportStateVerify: true,
				// the write-only version counters live only in the configuration and the state
				ImportStateVerifyIgnore: []string{"credentials.0.password_wo_version", "credentials.1.password_wo_version"},
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					for _, st := range states {
						for k, v := range st.Attributes {
							if strings.Contains(v, cwtestMask) || strings.Contains(v, "Secret") {
								return fmt.Errorf("imported state attribute %s holds a secret: %q", k, v)
							}
						}
					}
					return nil
				},
			},
		},
	})
	if creates, updates, deletes := s.counts(); creates != 1 || updates != 1 || deletes != 1 {
		t.Errorf("stub saw %d creates, %d updates, %d deletes; want 1, 1, 1", creates, updates, deletes)
	}
}

func TestCredentialProfilesCreateExisting(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newCredentialStub(t)
	s.profiles[cwtestProfile] = []map[string]any{{"type": "ROBOT_USERPASS_SSH", "user_name": "other", "password": "x"}}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      s.profileConfig(sshEntry("user-a", "Secret-A", 1)),
			ExpectError: regexp.MustCompile(`already exists`),
		}},
	})
	if got := s.userPass(cwtestProfile); got[0]["user_name"] != "other" {
		t.Errorf("the existing profile was changed: %v", got)
	}
}

// The delete reads first and only leaves the state when the profile is gone, so the failing DELETE is never sent.
func TestCredentialProfilesDestroyAfterOutOfBandDelete(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newCredentialStub(t)
	cfg := s.profileConfig(sshEntry("user-a", "Secret-A", 1))

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: cfg},
			{
				// The refresh before the destroy still finds the profile; it is gone when the delete reads it.
				PreConfig: func() {
					s.mu.Lock()
					s.vanish, s.vanishOf = 1, cwtestProfile
					s.mu.Unlock()
				},
				Config:  cfg,
				Destroy: true,
			},
		},
	})
	if _, _, deletes := s.counts(); deletes != 0 {
		t.Errorf("the provider sent %d DELETE calls for a profile that was already gone", deletes)
	}
	if s.userPass(cwtestProfile) != nil {
		t.Error("the profile is still on the stub")
	}
}

// Crosswork answers a read of a missing profile with HTTP 200 and no "data" member.
func TestCredentialProfilesDataSource(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newCredentialStub(t)
	s.profiles[cwtestProfile] = []map[string]any{{"type": "ROBOT_USERPASS_SSH", "user_name": "user-a", "password": "x"}}
	config := func(id string) string {
		return s.providerBlock() + fmt.Sprintf("\ndata \"crosswork_credential_profiles\" \"test\" {\n  id = %q\n}\n", id)
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config(cwtestProfile),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.crosswork_credential_profiles.test", "profile", cwtestProfile),
					resource.TestCheckResourceAttr("data.crosswork_credential_profiles.test", "credentials.0.user_name", "user-a"),
				),
			},
			{Config: config("cwtest-missing"), ExpectError: regexp.MustCompile(`No object found with id: cwtest-missing`)},
		},
	})
}

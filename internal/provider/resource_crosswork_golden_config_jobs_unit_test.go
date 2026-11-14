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
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

const (
	gcJob = "crosswork_golden_config_jobs.j"
	gcTpl = "crosswork_golden_config_templates.t"
)

func (s *gcStub) config(devices, trigger string) string {
	return s.providerBlock() + fmt.Sprintf(`
resource "crosswork_golden_config_global_variables" "g" {
  name  = "cwtest_gc_g"
  value = "g1"
}

resource "crosswork_golden_config_templates" "t" {
  name           = "cwtest-ntp"
  vendor         = "Cisco Systems"
  software_type  = "IOS XR"
  product_series = ["Cisco 8000 Series Routers", "Cisco ASR 9000 Series Aggregation Services Routers"]
  tags           = ["z-tag", "a-tag"]
  config         = "ntp server {{ ntp_server }}\n"
  variables = [
    { name = "zeta", value = "z" },
    { name = "alpha", value = "a" },
    { name = "obj", type = "JSON", value = jsonencode({ b = 1, a = 2 }) },
  ]
}

resource "crosswork_golden_config_jobs" "j" {
  template_name            = crosswork_golden_config_templates.t.name
  template_config_id       = crosswork_golden_config_templates.t.config_id
  vendor                   = "Cisco Systems"
  product_series           = "Cisco 8000 Series Routers"
  mop                      = "Default Golden Config - Cisco"
  trigger                  = %q
  deployment_timeout       = 10
  devices                  = {%s}
}
`, trigger, devices)
}

func (s *gcStub) callsAre(want string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		if got := s.callLog(); got != want {
			return fmt.Errorf("calls were\n  %s\nwant\n  %s", got, want)
		}
		return nil
	}
}

func (s *gcStub) appsAre(want ...string) resource.TestCheckFunc {
	return func(*terraform.State) error {
		got := strings.Join(s.appNames(), ",")
		if got != strings.Join(want, ",") {
			return fmt.Errorf("applications on the stub are %s, want %s", got, strings.Join(want, ","))
		}
		return nil
	}
}

func (s *gcStub) resetCalls() {
	s.mu.Lock()
	s.calls = nil
	s.mu.Unlock()
}

var emptyPlanCheck = resource.ConfigPlanChecks{PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}}

func (s *gcStub) checkDestroy(*terraform.State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.apps) != 0 || len(s.templates) != 0 || len(s.gvars) != 0 {
		return fmt.Errorf("the stub still holds %d applications, %d templates and %d variables", len(s.apps), len(s.templates), len(s.gvars))
	}
	return nil
}

func TestGoldenConfigLifecycle(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	app := func(d string) string { return "cwtest-ntp-uuid-" + d }

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{
			{ // the controller's reordering of variables, tags and series causes no difference
				Config: s.config("d1 = {}, d2 = {}", "1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(gcJob, "devices.%", "2"),
					resource.TestCheckResourceAttr(gcJob, "id", "cwtest-ntp"),
					resource.TestCheckResourceAttr(gcTpl, "variables.0.name", "zeta"),
					resource.TestCheckResourceAttr(gcTpl, "variables.1.name", "alpha"),
					resource.TestCheckResourceAttr(gcTpl, "tags.#", "2"),
					s.appsAre(app("d1"), app("d2")),
					s.callsAre("POST gvar cwtest_gc_g | POST template cwtest-ntp | job cwtest-ntp d1,d2"),
				),
				ConfigPlanChecks: emptyPlanCheck,
			},
			{ // a new trigger and nothing else: the devices that did not change are re-applied, no job runs
				PreConfig:        s.resetCalls,
				Config:           s.config("d1 = {}, d2 = {}", "2"),
				Check:            s.callsAre("re-apply d1,d2"),
				ConfigPlanChecks: emptyPlanCheck,
			},
			{
				PreConfig: s.resetCalls,
				Config:    s.config(`d1 = {}, d2 = {}, d3 = { variables = [{ name = "ntp_server", value = "10.2.2.2" }] }`, "2"),
				Check: resource.ComposeTestCheckFunc(
					s.callsAre("job cwtest-ntp d3"),
					resource.TestCheckResourceAttr(gcJob, "devices.%", "3"),
					func(*terraform.State) error {
						s.mu.Lock()
						defer s.mu.Unlock()
						if got := s.apps[app("d3")].Variables["ntp_server"]; got != "10.2.2.2" {
							return fmt.Errorf("d3 got the value %q, want its own 10.2.2.2", got)
						}
						if len(s.apps[app("d1")].Variables) != 0 {
							return fmt.Errorf("d1 got job values %v, want none", s.apps[app("d1")].Variables)
						}
						return nil
					},
				),
				ConfigPlanChecks: emptyPlanCheck,
			},
			{
				PreConfig: s.resetCalls,
				Config:    s.config(`d1 = {}, d3 = { variables = [{ name = "ntp_server", value = "10.2.2.2" }] }`, "2"),
				Check: resource.ComposeTestCheckFunc(
					s.callsAre("DELETE app "+app("d2")),
					s.appsAre(app("d1"), app("d3")),
				),
				ConfigPlanChecks: emptyPlanCheck,
			},
			{ // the application of d1 is removed on the controller: the read drops the device and the plan applies it again
				PreConfig: func() {
					s.mu.Lock()
					delete(s.apps, app("d1"))
					s.calls = nil
					s.mu.Unlock()
				},
				Config:           s.config(`d1 = {}, d3 = { variables = [{ name = "ntp_server", value = "10.2.2.2" }] }`, "2"),
				Check:            s.callsAre("job cwtest-ntp d1"),
				ConfigPlanChecks: emptyPlanCheck,
			},
			{
				ResourceName:  gcTpl,
				ImportState:   true,
				ImportStateId: "cwtest-ntp",
				ImportStateCheck: func(states []*terraform.InstanceState) error {
					a := states[0].Attributes
					if a["name"] != "cwtest-ntp" || a["variables.#"] != "3" || a["force_delete"] != "false" {
						return fmt.Errorf("unexpected imported state: %v", a)
					}
					return nil
				},
			},
		},
	})
	// The destroy removed the applications before the template, and deleted the template without force.
	calls := s.callLog()
	if !strings.Contains(calls, "DELETE template 1 force=false") {
		t.Errorf("the template was not deleted without force: %s", calls)
	}
	if i, j := strings.LastIndex(calls, "DELETE app"), strings.Index(calls, "DELETE template"); i < 0 || j < 0 || i > j {
		t.Errorf("the applications must be deleted before the template: %s", calls)
	}
}

// A create in which a device fails saves no state: a resource saved with an error is tainted, and the next apply
// would destroy it and remove the configuration from the devices that took it.
func TestGoldenConfigCreateFailure(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	app := func(d string) string { return "cwtest-ntp-uuid-" + d }
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() { s.failDevices["d2"] = "NSO error: Failed to connect to device d2" },
				Config:      s.config("d1 = {}, d2 = {}", "1"),
				ExpectError: regexp.MustCompile(`(?s)did not take the template.*d2: job 1 \(cwtest-ntp\).*Failed to connect to device d2.*No state was saved`),
			},
			{
				PreConfig: func() {
					s.mu.Lock()
					delete(s.failDevices, "d2")
					s.mu.Unlock()
				},
				Config: s.config("d1 = {}, d2 = {}", "1"),
				Check: resource.ComposeTestCheckFunc(
					s.appsAre(app("d1"), app("d2")),
					resource.TestCheckResourceAttr(gcJob, "devices.%", "2"),
					func(*terraform.State) error {
						if n := len(s.callsWith("job ")); n != 2 {
							return fmt.Errorf("%d jobs were started, want 2 (one failed, one that fixed it)", n)
						}
						return nil
					},
				),
				ConfigPlanChecks: emptyPlanCheck,
			},
		},
	})
}

// An update in which a device fails saves what is in place and is not tainted: the next apply retries only the failed device.
func TestGoldenConfigUpdateFailure(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	app := func(d string) string { return "cwtest-ntp-uuid-" + d }
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{
			{Config: s.config("d1 = {}", "1"), ConfigPlanChecks: emptyPlanCheck},
			{
				PreConfig:   func() { s.resetCalls(); s.failDevices["d2"] = "NSO error: bad line" },
				Config:      s.config("d1 = {}, d2 = {}", "1"),
				ExpectError: regexp.MustCompile(`(?s)1 device\(s\) did not take the template.*d2`),
			},
			{
				PreConfig: func() {
					s.mu.Lock()
					delete(s.failDevices, "d2")
					s.calls = nil
					s.mu.Unlock()
				},
				Config: s.config("d1 = {}, d2 = {}", "1"),
				Check: resource.ComposeTestCheckFunc(
					s.callsAre("job cwtest-ntp d2"),
					s.appsAre(app("d1"), app("d2")),
				),
				ConfigPlanChecks: emptyPlanCheck,
			},
		},
	})
}

// A re-apply that fails for a device keeps the old trigger, so the next plan still shows the change and tries again.
func TestGoldenConfigReapplyFailure(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{
			{Config: s.config("d1 = {}, d2 = {}", "1")},
			{
				PreConfig:   func() { s.resetCalls(); s.failReapply["d2"] = "Failed to re-apply template: skipped 1 line" },
				Config:      s.config("d1 = {}, d2 = {}", "2"),
				ExpectError: regexp.MustCompile(`(?s)did not take the template.*d2: re-apply: Failed to re-apply template`),
			},
			{
				PreConfig: func() {
					s.mu.Lock()
					delete(s.failReapply, "d2")
					s.calls = nil
					s.mu.Unlock()
				},
				Config:           s.config("d1 = {}, d2 = {}", "2"),
				Check:            resource.ComposeTestCheckFunc(s.callsAre("re-apply d1,d2"), resource.TestCheckResourceAttr(gcJob, "trigger", "2")),
				ConfigPlanChecks: emptyPlanCheck,
			},
		},
	})
}

// A device that did not change and is not ready when the trigger changes is not re-applied; the state keeps the old
// trigger, so the next plan shows the change again and the device gets it once it is ready.
func TestGoldenConfigReapplyNotReady(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{
			{Config: s.config("d1 = {}, d2 = {}", "1")},
			{
				PreConfig:   func() { s.resetCalls(); s.mu.Lock(); s.nodes["d2"].NSO = "NOT_SYNCED"; s.mu.Unlock() },
				Config:      s.config("d1 = {}, d2 = {}", "2"),
				ExpectError: regexp.MustCompile(`(?s)did not take the template.*d2: not ready`),
			},
			{
				PreConfig: func() {
					s.mu.Lock()
					s.nodes["d2"].NSO = nsoSynced
					s.calls = nil
					s.mu.Unlock()
				},
				Config:           s.config("d1 = {}, d2 = {}", "2"),
				Check:            resource.ComposeTestCheckFunc(s.callsAre("re-apply d1,d2"), resource.TestCheckResourceAttr(gcJob, "trigger", "2")),
				ConfigPlanChecks: emptyPlanCheck,
			},
		},
	})
}

func TestGoldenConfigDeviceNotReady(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					s.nodes["d2"].NSO = "FETCH_SSH_KEYS_FAILED"
					s.nodes["d3"].NSO = "NOT_SYNCED"
				},
				Config:      s.config("d1 = {}, d2 = {}, d3 = {}", "1"),
				ExpectError: regexp.MustCompile(`(?s)d2: not ready: reachability CONN_STATE_REACHABLE, nso FETCH_SSH_KEYS_FAILED, admin ROBOT_ADMIN_STATE_UP, waited 0s.*d3: not ready: reachability CONN_STATE_REACHABLE, nso NOT_SYNCED`),
			},
			{
				PreConfig: func() {
					s.mu.Lock()
					s.nodes["d2"].NSO, s.nodes["d3"].NSO = nsoSynced, nsoSynced
					s.mu.Unlock()
				},
				Config: s.config("d1 = {}, d2 = {}, d3 = {}", "1"),
				Check: func(*terraform.State) error {
					// The ready device had its job in the failed apply, and the others got theirs afterwards.
					jobs := s.callsWith("job ")
					if len(jobs) != 2 || jobs[0] != "job cwtest-ntp d1" || jobs[1] != "job cwtest-ntp d1,d2,d3" {
						return fmt.Errorf("jobs were %v", jobs)
					}
					return nil
				},
			},
		},
	})
}

// The readiness mode decides which states a device needs: synced (the default) takes a device that is admin down and
// not reachable as long as NSO has it synced, full needs it reachable and admin up as well.
func TestGoldenConfigReadinessMode(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	s.nodes["d1"].Reach, s.nodes["d1"].Admin = "CONN_STATE_UNKNOWN", "ROBOT_ADMIN_STATE_DOWN"
	s.readyMode = "full"
	full := s.config("d1 = {}", "1")
	s.readyMode = ""
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      full,
				ExpectError: regexp.MustCompile(`d1: not ready: reachability CONN_STATE_UNKNOWN, nso SYNCED, admin ROBOT_ADMIN_STATE_DOWN, waited 0s`),
			},
			{
				Config: s.config("d1 = {}", "1"),
				Check: resource.ComposeTestCheckFunc(
					func(*terraform.State) error {
						if jobs := s.callsWith("job "); len(jobs) != 1 || jobs[0] != "job cwtest-ntp d1" {
							return fmt.Errorf("jobs were %v", jobs)
						}
						return nil
					},
				),
			},
		},
	})
}

func TestGoldenConfigDeviceBecomesReady(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	s.nodes["d2"].NSO = "NOT_SYNCED"
	go func() {
		for i := 0; i < 200; i++ {
			if len(s.callsWith("POST template")) > 0 {
				break
			}
			sleepMs(5)
		}
		sleepMs(30)
		s.mu.Lock()
		s.nodes["d2"].NSO = nsoSynced
		s.mu.Unlock()
	}()
	s.readyTimeout = 5
	cfg := s.config("d1 = {}, d2 = {}", "1")
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{{
			Config: cfg,
			Check:  s.callsAre("POST gvar cwtest_gc_g | POST template cwtest-ntp | job cwtest-ntp d1,d2"),
		}},
	})
}

func TestGoldenConfigSecrets(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	const secret = "s3cr3t-value-1"
	cfg := s.providerBlock() + fmt.Sprintf(`
resource "crosswork_golden_config_global_variables" "g" {
  name             = "cwtest_gc_key"
  value_wo         = %q
  value_wo_version = 1
}

resource "crosswork_golden_config_templates" "t" {
  name           = "cwtest-key"
  vendor         = "Cisco Systems"
  software_type  = "IOS XR"
  product_series = ["Cisco 8000 Series Routers"]
  config         = "key {{ ntp_key }}"
  variables      = [{ name = "ntp_key", value_wo = %q, value_wo_version = 1 }]
}

resource "crosswork_golden_config_jobs" "j" {
  template_name            = crosswork_golden_config_templates.t.name
  template_config_id       = crosswork_golden_config_templates.t.config_id
  vendor                   = "Cisco Systems"
  product_series           = "Cisco 8000 Series Routers"
  mop                      = "Default Golden Config - Cisco"
  devices = {
    d1 = { variables = [{ name = "ntp_key", value_wo = %q, value_wo_version = 1 }] }
    d2 = { variables = [{ name = "ntp_key", value_wo = %q, value_wo_version = 1 }] }
  }
}
`, secret, secret+"-template", secret+"-d1", secret+"-d2")
	noSecret := func(st *terraform.State) error {
		b := fmt.Sprint(st.RootModule().Resources)
		if strings.Contains(b, "s3cr3t") {
			return fmt.Errorf("a secret is in the state")
		}
		return nil
	}
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{{
			Config: cfg,
			Check: resource.ComposeTestCheckFunc(noSecret,
				func(*terraform.State) error {
					s.mu.Lock()
					got := strings.Join(s.secrets, " ")
					s.mu.Unlock()
					for _, want := range []string{secret, secret + "-template", secret + "-d1", secret + "-d2"} {
						if !strings.Contains(got, want) {
							return fmt.Errorf("the controller never received %q", want)
						}
					}
					jobs := s.callsWith("job ")
					if len(jobs) != 2 {
						return fmt.Errorf("jobs were %v, want one per secret", jobs)
					}
					return nil
				}),
			ConfigPlanChecks: emptyPlanCheck,
		}},
	})
}

func sleepMs(n int) { time.Sleep(time.Duration(n) * time.Millisecond) }

func (s *gcStub) twoJobsConfig() string {
	first := s.config("d1 = {}", "1")
	body := strings.TrimPrefix(first, s.providerBlock())
	second := strings.NewReplacer(`"t"`, `"t2"`, `"j"`, `"j2"`, "templates.t.", "templates.t2.", "cwtest-ntp", "cwtest-ntp2", "d1 = {}", "d2 = {}").Replace(body)
	i := strings.Index(second, `resource "crosswork_golden_config_global_variables"`)
	j := strings.Index(second[i:], "\n}\n") + i + len("\n}\n")
	return first + second[:i] + second[j:]
}

// Golden Configuration jobs take one lock: two jobs of one apply never run at the same time, because NSO refuses two
// overlapping pushes of the templates and global variables. The control turns the lock off and sees them overlap.
func TestGoldenConfigJobsRunOneAtATime(t *testing.T) {
	for name, serialize := range map[string]bool{"with the lock": true, "control without the lock": false} {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			Shutdown()
			s := newGCStub(t)
			s.finishAfter = 3
			s.postDelay = 30 * time.Millisecond
			old := gcSerialize
			gcSerialize = serialize
			t.Cleanup(func() { gcSerialize = old })
			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				CheckDestroy:             s.checkDestroy,
				Steps: []resource.TestStep{{
					Config: s.twoJobsConfig(),
					Check: func(*terraform.State) error {
						s.mu.Lock()
						defer s.mu.Unlock()
						if len(s.jobs) != 2 {
							return fmt.Errorf("the stub ran %d jobs, want 2", len(s.jobs))
						}
						if serialize && s.maxRunning != 1 {
							return fmt.Errorf("%d jobs ran at once, want 1", s.maxRunning)
						}
						if !serialize && s.maxRunning < 2 {
							return fmt.Errorf("the control did not overlap the jobs (max %d), so the test cannot show the lock", s.maxRunning)
						}
						return nil
					},
				}},
			})
		})
	}
}

// A job that does not finish in its time (deployment_timeout, at least 10 seconds) fails its devices with the
// message that it keeps running, releases the lock, and the next job starts.
func TestGoldenConfigJobTimeoutReleasesTheLock(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newGCStub(t)
	s.stuckJob = 1
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             s.checkDestroy,
		Steps: []resource.TestStep{
			{
				Config:      s.twoJobsConfig(),
				ExpectError: regexp.MustCompile(`(?s)did not take the template.*job 1 \(cwtest-ntp2?\).*it keeps running on the controller`),
			},
			{ // the first job timed out and released the lock: the second one ran
				PreConfig: func() {
					s.mu.Lock()
					defer s.mu.Unlock()
					if len(s.jobs) != 2 {
						t.Errorf("%d jobs were started, want 2: the job behind the stuck one did not start", len(s.jobs))
					}
					s.stuckJob = 0
				},
				Config: s.twoJobsConfig(),
				Check:  resource.TestCheckResourceAttr(gcJob, "devices.%", "1"),
			},
		},
	})
}

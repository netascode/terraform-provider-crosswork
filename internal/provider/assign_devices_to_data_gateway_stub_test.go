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

// A failed ADD, REMOVE or UPDATE answers HTTP 200 with a failed job, which may have applied part of the request.

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type mappedNode struct{ uuid, gateway string }

type mappingStub struct {
	*ssoStub
	mu         sync.Mutex
	nodes      map[string]*mappedNode
	calls      []string // "<operation> <gateway> <uuids>"
	failAfter  bool     // apply the next mapping request, then answer a failed job
	failBefore string   // refuse the next mapping request with this message, without applying it
}

var criteriaGateway = regexp.MustCompile(`where dg_uuid=(\S+)$`)

func newMappingStub(t *testing.T, nodes map[string]string) *mappingStub {
	t.Helper()
	s := &mappingStub{ssoStub: newSSOStub(t), nodes: map[string]*mappedNode{}}
	for name, gateway := range nodes {
		s.nodes[name] = &mappedNode{uuid: "u-" + name, gateway: gateway}
	}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/crosswork/inventory/v1/nodes/query":
			s.mu.Lock()
			defer s.mu.Unlock()
			raw, _ := io.ReadAll(r.Body)
			fmt.Fprint(w, s.nodeList(gjson.ParseBytes(raw)))
		case r.Method == http.MethodPut && r.URL.Path == "/crosswork/inventory/v1/dg/devicemapping":
			s.mu.Lock()
			defer s.mu.Unlock()
			raw, _ := io.ReadAll(r.Body)
			fmt.Fprint(w, s.mapping(gjson.ParseBytes(raw).Get("dgDeviceMappings.0")))
		default:
			sso.ServeHTTP(w, r)
		}
	})
	t.Cleanup(Shutdown)
	return s
}

func (s *mappingStub) names() []string {
	out := make([]string, 0, len(s.nodes))
	for n := range s.nodes {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

func (s *mappingStub) nodeList(body gjson.Result) string {
	var matches []string
	for _, n := range s.names() {
		node := s.nodes[n]
		switch {
		case criteriaGateway.MatchString(body.Get("filterData.Criteria").String()):
			if node.gateway == criteriaGateway.FindStringSubmatch(body.Get("filterData.Criteria").String())[1] {
				matches = append(matches, n)
			}
		case body.Get("filter.inventory_id").Exists():
			if strings.EqualFold(n, body.Get("filter.inventory_id").String()) {
				matches = append(matches, n)
			}
		default:
			matches = append(matches, n)
		}
	}
	size, page := int(body.Get("filterData.PageSize").Int()), int(body.Get("filterData.PageNum").Int())
	if size == 0 {
		size = len(matches) + 1
	}
	start := min(page*size, len(matches))
	end := min(start+size, len(matches))
	out := `{"data":[]}`
	for _, n := range matches[start:end] {
		item, _ := sjson.Set(`{}`, "inventory_id", n)
		item, _ = sjson.Set(item, "uuid", s.nodes[n].uuid)
		if g := s.nodes[n].gateway; g != "" {
			item, _ = sjson.Set(item, "dg_uuid", g)
		}
		out, _ = sjson.SetRaw(out, "data.-1", item)
	}
	out, _ = sjson.Set(out, "total_count", len(matches))
	return out
}

func (s *mappingStub) byUuid(uuid string) *mappedNode {
	for _, n := range s.nodes {
		if n.uuid == uuid {
			return n
		}
	}
	return nil
}

func (s *mappingStub) mapping(m gjson.Result) string {
	op, gateway := m.Get("mapping_oper").String(), m.Get("cdg_duuid").String()
	var uuids []string
	for _, u := range m.Get("device_uuid").Array() {
		uuids = append(uuids, u.String())
	}
	sort.Strings(uuids)
	s.calls = append(s.calls, fmt.Sprintf("%s %s %s", op, gateway, strings.Join(uuids, ",")))
	failed := func(message string) string {
		job, _ := sjson.Set(`{"job_id":"j","state":"JOB_FAILED","type":"Update mappings for 1 data gateway."}`, "error", message)
		return job
	}
	if s.failBefore != "" {
		message := s.failBefore
		s.failBefore = ""
		return failed(message)
	}
	var problem string
	for _, u := range uuids {
		n := s.byUuid(u)
		switch {
		case n == nil:
			problem = "Device with ID " + u + " does not exist or is Invalid!"
		case op == "ADD_OPER" && n.gateway != "":
			problem = "Mapping operation : ADD_OPER cannot be performed on Device with ID : " + u + " because invalid dg ID is requested"
		case op == "REMOVE_OPER" && n.gateway != gateway:
			problem = "Mapping operation : REMOVE_OPER cannot be performed on Device with ID : " + u + " because requested dg is not the same as current one"
		case op == "UPDATE_OPER" && n.gateway == gateway:
			problem = "Mapping operation : UPDATE_OPER cannot be performed on Device with ID : " + u + " because source and target dg is same"
		case op == "REMOVE_OPER":
			n.gateway = ""
		default:
			n.gateway = gateway
		}
	}
	if s.failAfter {
		s.failAfter = false
		return failed("the job ended with an error after it changed the devices")
	}
	if problem != "" {
		return failed(problem)
	}
	return fmt.Sprintf(`{"job_id":"j","state":"JOB_COMPLETED","type":"Update mappings for 1 data gateway.","impacted":[%q]}`, gateway)
}

func (s *mappingStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

func (s *mappingStub) config(devices ...string) string {
	quoted := make([]string, len(devices))
	for i, d := range devices {
		quoted[i] = fmt.Sprintf("%q", d)
	}
	return s.providerBlock() + fmt.Sprintf("\nresource \"crosswork_assign_devices_to_data_gateway\" \"test\" {\n  data_gateway_id = \"gw-12345678\"\n  devices         = [%s]\n}\n", strings.Join(quoted, ", "))
}

func TestAssignDevicesToDataGatewayLifecycle(t *testing.T) {
	clearEnv(t)
	Shutdown()
	// a is on the gateway already, d is on another gateway, e is on the gateway and belongs to someone else
	s := newMappingStub(t, map[string]string{"a": "gw-12345678", "b": "", "c": "", "d": "gw-87654321", "e": "gw-12345678", "f": "", "g": ""})
	const res = "crosswork_assign_devices_to_data_gateway.test"
	last := func() string { return s.calls[len(s.calls)-1] }
	callsBefore := 0

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: s.config("a", "b"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", "gw-12345678"),
					resource.TestCheckResourceAttr(res, "devices.#", "2"),
					func(*terraform.State) error {
						// a was on the gateway already, so only b is sent, as an ADD (an ADD of a would be refused)
						if strings.Join(s.calls, "|") != "ADD_OPER gw-12345678 u-b" {
							return fmt.Errorf("calls: %v", s.calls)
						}
						return nil
					},
				),
			},
			{
				Config: s.config("b", "c"),
				Check: func(*terraform.State) error {
					if got := strings.Join(s.calls[1:], "|"); got != "REMOVE_OPER gw-12345678 u-a|ADD_OPER gw-12345678 u-c" {
						return fmt.Errorf("calls: %v", s.calls)
					}
					return nil
				},
			},
			{
				// d is on another gateway: it moves with an UPDATE, not an ADD
				Config: s.config("b", "c", "d"),
				Check: func(*terraform.State) error {
					if last() != "UPDATE_OPER gw-12345678 u-d" {
						return fmt.Errorf("last call: %s", last())
					}
					return nil
				},
			},
			{
				// the job fails after it attached f: the members are read again and the result is accepted
				PreConfig: func() { s.mu.Lock(); s.failAfter = true; s.mu.Unlock() },
				Config:    s.config("b", "c", "d", "f"),
				Check:     resource.TestCheckResourceAttr(res, "devices.#", "4"),
			},
			{
				PreConfig: func() {
					s.mu.Lock()
					s.failBefore = "Data Gateway gw-12345678 does not exist! Cannot update Device Mapping."
					s.mu.Unlock()
				},
				Config:      s.config("b", "c", "d", "f", "g"),
				ExpectError: regexp.MustCompile(`(?s)Data\s+Gateway\s+gw-12345678\s+does\s+not\s+exist!\s+Cannot\s+update\s+Device\s+Mapping`),
			},
			{
				Config:      s.config("b", "c", "d", "f", "nosuch"),
				ExpectError: regexp.MustCompile(`(?s)the\s+device\s+nosuch\s+does\s+not\s+exist`),
			},
			{
				PreConfig: func() { s.mu.Lock(); callsBefore = len(s.calls); s.mu.Unlock() },
				Config:    s.config("b", "c", "d", "f"),
				Check: func(*terraform.State) error {
					if len(s.calls) != callsBefore {
						return fmt.Errorf("a plan with no change sent %d calls", len(s.calls)-callsBefore)
					}
					return nil
				},
			},
		},
	})
	// the destroy detached the devices of the resource and left e, which someone else attached
	for _, n := range []string{"b", "c", "d", "f", "g"} {
		if s.nodes[n].gateway != "" {
			t.Errorf("%s is still attached to %s", n, s.nodes[n].gateway)
		}
	}
	if s.nodes["e"].gateway != "gw-12345678" {
		t.Errorf("e, attached by someone else, was detached")
	}
}

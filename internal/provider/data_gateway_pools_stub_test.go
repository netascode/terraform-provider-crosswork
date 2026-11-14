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

// The id is "<uuid> <name>", spares at 0 are left out, update and delete answer 204, and data gateways are listed after
// the create.

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type stubPool struct {
	id   string
	body gjson.Result
}

type poolStub struct {
	*ssoStub
	mu         sync.Mutex
	pools      []*stubPool
	next       int
	failCreate int  // creates answered "is currently locked" before one succeeds
	attached   bool // a data gateway has devices: delete and a change of the addresses are refused
	vdgDelay   int  // gateway queries answered empty after a create
	vdgLeft    int
	posts      int
	charset    int // requests that sent a charset in the content type
}

func newPoolStub(t *testing.T) *poolStub {
	t.Helper()
	s := &poolStub{ssoStub: newSSOStub(t)}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "/crosswork/dg-manager/"
		if !strings.HasPrefix(r.URL.Path, prefix) {
			sso.ServeHTTP(w, r)
			return
		}
		s.mu.Lock()
		defer s.mu.Unlock()
		if strings.Contains(r.Header.Get("Content-Type"), "charset") {
			s.charset++
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error":"received unexpected Content-Type: `+r.Header.Get("Content-Type")+`"}`)
			return
		}
		raw, _ := io.ReadAll(r.Body)
		body := gjson.ParseBytes(raw)
		refuse := func(message string) {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error":%q}`, message)
		}
		switch r.Method + " " + strings.TrimPrefix(r.URL.Path, prefix) {
		case "POST v2/hapool":
			s.posts++
			if s.failCreate > 0 {
				s.failCreate--
				refuse("update HA Pool cannot be processed since pool: x is currently locked")
				return
			}
			pool := body.Get("pools.0")
			if want := len(pool.Get("ipaddrs").Array()) + int(pool.Get("spares").Int()); len(pool.Get("pdgUuids").Array()) != want {
				refuse(fmt.Sprintf("The number of data gateway instance(s) {%d} added to the pool does not match the number of virtual IPs {%d} + number of standby data gateways {%s} specified.",
					len(pool.Get("pdgUuids").Array()), len(pool.Get("ipaddrs").Array()), pool.Get("spares").String()))
				return
			}
			for _, p := range s.pools {
				if p.body.Get("name").String() == pool.Get("name").String() {
					refuse("HA Pool with name " + pool.Get("name").String() + " already exists")
					return
				}
			}
			s.next++
			id := fmt.Sprintf("pool-id-%d", s.next)
			s.pools = append(s.pools, &stubPool{id: id, body: pool})
			s.vdgLeft = s.vdgDelay
			fmt.Fprintf(w, `{"puuid":[%q]}`, id+" "+pool.Get("name").String())
		case "PUT v2/hapool":
			pool := body.Get("pools.0")
			for _, p := range s.pools {
				if p.id == pool.Get("puuid").String() {
					if s.attached && poolAddresses(p.body) != poolAddresses(pool) {
						refuse("southbound ip address is being used by Virtual Data Gateway g which has devices attached. Please detach devices before deleting southbound ip.")
						return
					}
					p.body = pool
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
			refuse("pool not found")
		case "DELETE v2/hapool":
			id := body.Get("pools.0.puuid").String()
			if s.attached {
				refuse("Virtual Data Gateway g in the pool p has devices attached to it. Please detach devices before deleting the pool.")
				return
			}
			for i, p := range s.pools {
				if p.id == id {
					s.pools = append(s.pools[:i], s.pools[i+1:]...)
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
			refuse("pool not found")
		case "POST v2/hapool/query":
			fmt.Fprint(w, s.poolList())
		case "POST v2/vdg/query":
			fmt.Fprint(w, s.gatewayList())
		case "POST v1/dg/query":
			fmt.Fprint(w, `{"data":[{"duuid":"dg-1","name":"cw72-dg-1a","configData":{"role":"UNASSIGNED","adminState":"AS_UP"},"operationalData":{"operState":"OS_NOT_READY"}},`+
				`{"duuid":"dg-2","name":"cw72-dg-2a","configData":{"role":"ASSIGNED","poolId":"pool-id-1","vdgUuid":"vdg-pool-id-1-1","adminState":"AS_UP"},"operationalData":{"operState":"OS_UP"}}],"totalCount":2}`)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	t.Cleanup(Shutdown)
	return s
}

func poolAddresses(pool gjson.Result) string {
	return strings.Join(poolVirtualIps(pool), ",")
}

func (s *poolStub) poolList() string {
	out := `{"data":[]}`
	for _, p := range s.pools {
		pool, _ := sjson.Set(`{}`, "puuid", p.id)
		pool, _ = sjson.Set(pool, "name", p.body.Get("name").String())
		if d := p.body.Get("description"); d.Exists() {
			pool, _ = sjson.Set(pool, "description", d.String())
		}
		for i, a := range p.body.Get("ipaddrs").Array() {
			pool, _ = sjson.Set(pool, fmt.Sprintf("ipaddrs.%d.inetaddrs.0.inetAf", i), a.Get("inetaddrs.0.inet_af").String())
			pool, _ = sjson.Set(pool, fmt.Sprintf("ipaddrs.%d.inetaddrs.0.inetAddr", i), a.Get("inetaddrs.0.inet_addr").String())
			pool, _ = sjson.Set(pool, fmt.Sprintf("ipaddrs.%d.inetaddrs.0.mask", i), a.Get("inetaddrs.0.mask").String())
		}
		pool, _ = sjson.SetRaw(pool, "pdgUuids", p.body.Get("pdgUuids").Raw)
		if sp := p.body.Get("spares").String(); sp != "0" {
			pool, _ = sjson.Set(pool, "spares", sp)
		}
		pool, _ = sjson.Set(pool, "subnetMask", p.body.Get("subnetMask").String())
		pool, _ = sjson.Set(pool, "haStrategy", "VIP_BASED")
		pool, _ = sjson.Set(pool, "protectionStatus", "PROTECTED")
		out, _ = sjson.SetRaw(out, "data.-1", pool)
	}
	return out
}

// One gateway per virtual IP of a pool, named <pool>-<n>; the first queries after a create answer an empty list.
func (s *poolStub) gatewayList() string {
	out := `{"data":[]}`
	if s.vdgLeft > 0 {
		s.vdgLeft--
		return out
	}
	for _, p := range s.pools {
		for i := range p.body.Get("ipaddrs").Array() {
			g, _ := sjson.Set(`{}`, "name", fmt.Sprintf("%s-%d", p.body.Get("name").String(), i+1))
			g, _ = sjson.Set(g, "vdgUuid", fmt.Sprintf("vdg-%s-%d", p.id, i+1))
			g, _ = sjson.Set(g, "poolId", p.id)
			out, _ = sjson.SetRaw(out, "data.-1", g)
		}
	}
	return out
}

func (s *poolStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

func fastRetry(t *testing.T) {
	old := poolRetryInterval
	poolRetryInterval = time.Millisecond
	t.Cleanup(func() { poolRetryInterval = old })
}

const poolConfig = `
resource "crosswork_data_gateway_pools" "test" {
  name          = "cwtest-pool"
  description   = %q
  virtual_ips   = [%q]
  prefix_length = %d
  spares        = 1
  instances = ["dg-1", "dg-2"]
}
`

func poolStep(s *poolStub, description, address string, prefix int) string {
	return s.providerBlock() + fmt.Sprintf(poolConfig, description, address, prefix)
}

func TestDataGatewayPoolsLifecycle(t *testing.T) {
	clearEnv(t)
	Shutdown()
	fastRetry(t)
	s := newPoolStub(t)
	s.vdgDelay = 2
	const res = "crosswork_data_gateway_pools.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: poolStep(s, "first", "192.0.2.198", 24),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", "pool-id-1"),
					resource.TestCheckResourceAttr(res, "spares", "1"),
					resource.TestCheckResourceAttr(res, "cdg_readiness_timeout", "60"),
					resource.TestCheckResourceAttr(res, "instances.#", "2"),
					resource.TestCheckResourceAttr(res, "data_gateway_ids.#", "1"),
					resource.TestCheckResourceAttr(res, "data_gateway_ids.0", "vdg-pool-id-1-1"),
				),
			},
			{
				Config: poolStep(s, "second", "192.0.2.198", 24),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", "pool-id-1"),
					resource.TestCheckResourceAttr(res, "description", "second"),
					resource.TestCheckResourceAttr(res, "data_gateway_ids.0", "vdg-pool-id-1-1"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateIdFunc: func(*terraform.State) (string, error) { return "cwtest-pool", nil },
				ImportStateVerify: true,
			},
			{
				// another prefix length replaces the pool, which gets a new id and new gateways
				Config: poolStep(s, "second", "192.0.2.198", 25),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", "pool-id-2"),
					resource.TestCheckResourceAttr(res, "data_gateway_ids.0", "vdg-pool-id-2-1"),
				),
			},
		},
	})
	if s.charset != 0 {
		t.Errorf("%d requests to the data gateway manager sent a charset in the content type", s.charset)
	}
	if len(s.pools) != 0 {
		t.Errorf("%d pools are left on the controller", len(s.pools))
	}
}

func TestDataGatewayPoolsRepeatsATransientRefusal(t *testing.T) {
	clearEnv(t)
	Shutdown()
	fastRetry(t)
	s := newPoolStub(t)
	s.failCreate = 2

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: poolStep(s, "first", "192.0.2.198", 24),
			Check:  resource.TestCheckResourceAttr("crosswork_data_gateway_pools.test", "id", "pool-id-1"),
		}},
	})
	if s.posts != 3 {
		t.Errorf("the create was sent %d times, want 3 (two refusals and one success)", s.posts)
	}
}

func TestDataGatewayPoolsReturnsTheControllerMessageUnchanged(t *testing.T) {
	clearEnv(t)
	Shutdown()
	s := newPoolStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			// two instances, one virtual IP and no spare: the controller refuses, and the message is final
			Config:      s.providerBlock() + strings.Replace(fmt.Sprintf(poolConfig, "x", "192.0.2.198", 24), `spares        = 1`, `spares        = 0`, 1),
			ExpectError: regexp.MustCompile(`(?s)does not match the number of virtual IPs`),
		}},
	})
	if s.posts != 1 {
		t.Errorf("a final refusal was sent %d times, want 1", s.posts)
	}
}

func TestDataGatewayPoolsRefusedDeleteIsReturnedUnchanged(t *testing.T) {
	clearEnv(t)
	Shutdown()
	fastRetry(t)
	s := newPoolStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: poolStep(s, "first", "192.0.2.198", 24)},
			{
				PreConfig:   func() { s.mu.Lock(); s.attached = true; s.mu.Unlock() },
				Config:      s.providerBlock(),
				ExpectError: regexp.MustCompile(`(?s)has devices attached to it`),
			},
			{
				PreConfig: func() { s.mu.Lock(); s.attached = false; s.mu.Unlock() },
				Config:    s.providerBlock(),
			},
		},
	})
	if len(s.pools) != 0 {
		t.Errorf("%d pools are left on the controller", len(s.pools))
	}
}

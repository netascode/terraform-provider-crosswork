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

// A create answers with a job and no id; a failure is a job in HTTP 200; name and server type cannot change; a
// system-defined destination cannot be updated or deleted.

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

type stubDestination struct {
	id   string
	body gjson.Result
}

type destinationStub struct {
	*ssoStub
	mu      sync.Mutex
	rows    []*stubDestination
	next    int
	posts   int
	charset int // requests that sent a charset in the content type
}

var destinationNamePattern = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

func newDestinationStub(t *testing.T) *destinationStub {
	t.Helper()
	s := &destinationStub{ssoStub: newSSOStub(t)}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix = "/crosswork/dg-manager/v1/"
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
		dest := gjson.ParseBytes(raw).Get("providers.0")
		name := dest.Get("name").String()
		job := func(kind string, errs ...string) {
			if len(errs) == 0 {
				fmt.Fprintf(w, `{"job_id":"job-1","state":"JOB_COMPLETED","type":%q,"created_by":"admin"}`, kind)
				return
			}
			out, _ := sjson.Set(`{"job_id":"job-1","state":"JOB_FAILED","created_by":"admin"}`, "type", kind)
			out, _ = sjson.Set(out, "errors", errs)
			fmt.Fprint(w, out)
		}
		switch r.Method + " " + strings.TrimPrefix(r.URL.Path, prefix) {
		case "POST destinations/query":
			fmt.Fprint(w, s.list())
		case "POST destinations":
			s.posts++
			switch {
			case !destinationNamePattern.MatchString(name):
				job("Insert 1 Destination(s)", "Destination Name: "+name+" is not valid, only a-zA-Z0-9_.- are allowed.")
			case s.find(name) != nil:
				job("Insert 1 Destination(s)", "destination already present for UUID:  or Name: "+name)
			default:
				s.next++
				s.rows = append(s.rows, &stubDestination{id: fmt.Sprintf("dest-id-%d", s.next), body: dest})
				job("Insert 1 Destination(s)")
			}
		case "PUT destinations":
			row := s.byId(dest.Get("uuid").String())
			switch {
			case dest.Get("uuid").String() == "":
				job("Update 1 Destination(s)", "UUID missing in Update Destination request for : "+name)
			case row == nil:
				job("Update 1 Destination(s)", "No destination found for uuid : "+dest.Get("uuid").String())
			case row.body.Get("properties.IS_SYSTEM_DEFINED").String() == "true":
				job("Update 1 Destination(s)", "Cannot Update a System Defined Destination : "+name)
			case row.body.Get("name").String() != name:
				job("Update 1 Destination(s)", "Destination Name cannot be changed")
			case row.body.Get("properties.DESTINATION_TYPE").String() != dest.Get("properties.DESTINATION_TYPE").String():
				job("Update 1 Destination(s)", "Destination:"+name+", Destination Type cannot be changed ")
			default:
				row.body = withoutId(dest)
				job("Update 1 Destination(s)")
			}
		case "DELETE destinations":
			row := s.byId(dest.Get("uuid").String())
			switch {
			case row == nil:
				job("Delete 1 Destination(s)", "destination with Name  not found in cache")
			case row.body.Get("properties.IS_SYSTEM_DEFINED").String() == "true":
				job("Delete 1 Destination(s)", "Cannot delete a System Defined Destination : "+row.body.Get("name").String())
			default:
				for i, p := range s.rows {
					if p == row {
						s.rows = append(s.rows[:i], s.rows[i+1:]...)
						break
					}
				}
				job("Delete 1 Destination(s)")
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
	t.Cleanup(Shutdown)
	return s
}

func withoutId(dest gjson.Result) gjson.Result {
	out, _ := sjson.Delete(dest.Raw, "uuid")
	return gjson.Parse(out)
}

func (s *destinationStub) find(name string) *stubDestination {
	for _, r := range s.rows {
		if r.body.Get("name").String() == name {
			return r
		}
	}
	return nil
}

func (s *destinationStub) byId(id string) *stubDestination {
	for _, r := range s.rows {
		if r.id == id {
			return r
		}
	}
	return nil
}

func (s *destinationStub) seed(body string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	s.rows = append(s.rows, &stubDestination{id: fmt.Sprintf("dest-id-%d", s.next), body: gjson.Parse(body)})
}

func (s *destinationStub) setSystemDefined(name string, on bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	row := s.find(name)
	out, _ := sjson.Set(row.body.Raw, "properties.IS_SYSTEM_DEFINED", fmt.Sprint(on))
	row.body = gjson.Parse(out)
}

func (s *destinationStub) list() string {
	out := `{"data":[]}`
	for _, r := range s.rows {
		row, _ := sjson.Set(r.body.Raw, "uuid", r.id)
		out, _ = sjson.SetRaw(out, "data.-1", row)
	}
	out, _ = sjson.Set(out, "total_count", len(s.rows))
	return out
}

func (s *destinationStub) providerBlock() string {
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
}
`, s.URL)
}

func fastDestinationWait(t *testing.T) {
	old := destinationRetryInterval
	destinationRetryInterval = time.Millisecond
	t.Cleanup(func() { destinationRetryInterval = old })
}

const destinationKafkaConfig = `
resource "crosswork_data_gateway_destinations" "test" {
  name               = "cwtest-dest"
  data_source        = "data-gateway"
  server_type        = "kafka"
  encoding           = "gpbkv"
  compression_method = %q
  max_message_size   = 100000000
  buffer_memory      = 52428800
  batch_size         = 1048576
  linger_time        = 2000
  request_timeout    = 30
  collectors = [%s]
  connections = [{
    address       = "192.0.2.150"
    prefix_length = 20
    port          = 9192
  }]
  secure    = true
  auth_type = %q
}
`

func destinationKafkaStep(s *destinationStub, compression, collectors, auth string) string {
	return s.providerBlock() + fmt.Sprintf(destinationKafkaConfig, compression, collectors, auth)
}

const (
	destinationCliCollector  = `{ name = "CLI", linger_time = 1000 }`
	destinationSnmpCollector = `{ name = "SNMP", batch_size = 20000 }`
)

const destinationGrpcConfig = `
resource "crosswork_data_gateway_destinations" "test" {
  name               = "cwtest-dest"
  data_source        = "data-gateway"
  server_type        = "grpc"
  encoding           = "json"
  compression_method = "snappy"
  dispatch_type      = "unary"
  connections = [{
    host_name   = "collector"
    domain_name = "example.com"
    port        = 50052
  }]
}
`

func TestDataGatewayDestinationsLifecycle(t *testing.T) {
	clearEnv(t)
	Shutdown()
	fastDestinationWait(t)
	s := newDestinationStub(t)
	const res = "crosswork_data_gateway_destinations.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: destinationKafkaStep(s, "snappy", destinationCliCollector+", "+destinationSnmpCollector, "mutual"),
				Check: resource.ComposeTestCheckFunc(
					// the answer of the create has no id: the resource found it by the name
					resource.TestCheckResourceAttr(res, "id", "dest-id-1"),
					resource.TestCheckResourceAttr(res, "system_defined", "false"),
					resource.TestCheckResourceAttr(res, "collectors.#", "2"),
					resource.TestCheckResourceAttr(res, "collectors.1.name", "SNMP"),
					resource.TestCheckResourceAttr(res, "collectors.1.batch_size", "20000"),
					resource.TestCheckResourceAttr(res, "connections.0.address", "192.0.2.150"),
					resource.TestCheckResourceAttr(res, "connections.0.prefix_length", "20"),
					resource.TestCheckResourceAttr(res, "secure", "true"),
					resource.TestCheckResourceAttr(res, "auth_type", "mutual"),
				),
			},
			{
				Config: destinationKafkaStep(s, "zstd", destinationCliCollector, "server"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", "dest-id-1"),
					resource.TestCheckResourceAttr(res, "compression_method", "zstd"),
					resource.TestCheckResourceAttr(res, "auth_type", "server"),
					resource.TestCheckResourceAttr(res, "collectors.#", "1"),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateIdFunc: func(*terraform.State) (string, error) { return "cwtest-dest", nil },
				ImportStateVerify: true,
			},
			{
				// another server type replaces the destination, which gets a new id
				Config: s.providerBlock() + destinationGrpcConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", "dest-id-2"),
					resource.TestCheckResourceAttr(res, "server_type", "grpc"),
					resource.TestCheckResourceAttr(res, "connections.0.host_name", "collector"),
					resource.TestCheckNoResourceAttr(res, "secure"),
				),
			},
		},
	})
	if s.charset != 0 {
		t.Errorf("%d requests to the data gateway manager sent a charset in the content type", s.charset)
	}
	if len(s.rows) != 0 {
		t.Errorf("%d destinations are left on the controller", len(s.rows))
	}
}

func TestDataGatewayDestinationsReturnsTheControllerMessageUnchanged(t *testing.T) {
	clearEnv(t)
	Shutdown()
	fastDestinationWait(t)
	s := newDestinationStub(t)
	s.seed(`{"name":"cwtest-dest","properties":{"DESTINATION_TYPE":"destination_type_kafka","IS_SYSTEM_DEFINED":"false"},"connectivity_info":[]}`)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config:      destinationKafkaStep(s, "snappy", destinationCliCollector, "mutual"),
			ExpectError: regexp.MustCompile(`(?s)destination already present for UUID:  or Name: cwtest-dest`),
		}},
	})
	if s.posts != 1 {
		t.Errorf("a refused create was sent %d times, want 1", s.posts)
	}
}

func TestDataGatewayDestinationsSystemDefinedIsRefusedUnchanged(t *testing.T) {
	clearEnv(t)
	Shutdown()
	fastDestinationWait(t)
	s := newDestinationStub(t)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: destinationKafkaStep(s, "snappy", destinationCliCollector, "mutual")},
			{
				PreConfig:   func() { s.setSystemDefined("cwtest-dest", true) },
				Config:      destinationKafkaStep(s, "gzip", destinationCliCollector, "mutual"),
				ExpectError: regexp.MustCompile(`(?s)Cannot Update a System Defined Destination : cwtest-dest`),
			},
			{
				Config:      s.providerBlock(),
				ExpectError: regexp.MustCompile(`(?s)Cannot delete a System Defined Destination : cwtest-dest`),
			},
			{
				PreConfig: func() { s.setSystemDefined("cwtest-dest", false) },
				Config:    s.providerBlock(),
			},
		},
	})
	if len(s.rows) != 0 {
		t.Errorf("%d destinations are left on the controller", len(s.rows))
	}
}

func TestDataGatewayDestinationsKeepsAValueOfAnotherServerType(t *testing.T) {
	clearEnv(t)
	Shutdown()
	fastDestinationWait(t)
	s := newDestinationStub(t)
	const res = "crosswork_data_gateway_destinations.test"

	// The provider sends what the model sets and the controller keeps it, so a Kafka value on a gRPC destination is read back and
	// the next plan has no difference.
	config := s.providerBlock() + strings.Replace(destinationGrpcConfig, `dispatch_type      = "unary"`, "dispatch_type      = \"unary\"\n  linger_time        = 2000", 1)
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{{
			Config: config,
			Check:  resource.TestCheckResourceAttr(res, "linger_time", "2000"),
		}},
	})
	if len(s.rows) != 0 {
		t.Errorf("%d destinations are left on the controller", len(s.rows))
	}
}

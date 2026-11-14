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

// A failure is request_result REJECTED in HTTP 200; collector type and device set cannot change; a deleted job may still
// be returned by the configuration query.

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type collectionStubJob struct {
	body       gjson.Result // collection_job as it was sent
	staleReads int          // reads of the configuration query that still return a deleted job
	deleted    bool
}

type collectionStub struct {
	*ssoStub
	mu      sync.Mutex
	jobs    map[string]*collectionStubJob
	calls   []string // "PUT app/ctx", "DELETE app/ctx" and so on, in order
	devices map[string]string
	tags    map[string]bool
	stale   int // reads of a deleted job that still answer
}

var collectionDestinations = map[string]string{"cwtest-dest-kafka": "dest-id-1", "cwtest-dest-plain": "dest-id-2"}

func newCollectionStub(t *testing.T) *collectionStub {
	t.Helper()
	clearEnv(t)
	Shutdown()
	s := &collectionStub{ssoStub: newSSOStub(t), jobs: map[string]*collectionStubJob{}, stale: 2,
		devices: map[string]string{"dev-1": "uuid-1", "dev-2": "uuid-2", "dev-3": "uuid-3"}, tags: map[string]bool{"cwtest-tag": true}}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		raw, _ := io.ReadAll(r.Body)
		body := gjson.ParseBytes(raw)
		switch r.URL.Path {
		case nodesQueryPath:
			s.nodes(w, body)
		case destinationQueryPath:
			out := `{"data":[]}`
			names := make([]string, 0, len(collectionDestinations))
			for n := range collectionDestinations {
				names = append(names, n)
			}
			sort.Strings(names)
			for _, n := range names {
				row, _ := sjson.Set(`{}`, "name", n)
				row, _ = sjson.Set(row, "uuid", collectionDestinations[n])
				out, _ = sjson.SetRaw(out, "data.-1", row)
			}
			out, _ = sjson.Set(out, "total_count", len(names))
			fmt.Fprint(w, out)
		case collectionJobPath:
			s.change(w, r.Method, body)
		case collectionConfigPath, collectionStatusPath, collectionStatePath:
			s.query(w, r.URL.Path, body)
		default:
			sso.ServeHTTP(w, r)
		}
	})
	t.Cleanup(Shutdown)
	return s
}

func (s *collectionStub) nodes(w http.ResponseWriter, body gjson.Result) {
	out := `{"data":[]}`
	add := func(name, uuid string) {
		row, _ := sjson.Set(`{}`, "inventory_id", name)
		row, _ = sjson.Set(row, "uuid", uuid)
		out, _ = sjson.SetRaw(out, "data.-1", row)
	}
	for name, uuid := range s.devices {
		if f := body.Get("filter.inventory_id"); f.Exists() && !strings.EqualFold(f.String(), name) {
			continue
		}
		if f := body.Get("filter.uuid"); f.Exists() && f.String() != uuid {
			continue
		}
		add(name, uuid)
	}
	fmt.Fprint(w, out)
}

func collectionAccepted() string {
	return `{"result":{"request_result":"ACCEPTED","error":{"error":""}}}`
}

func collectionRejected(msg string) string {
	out, _ := sjson.Set(`{"result":{"request_result":"REJECTED"}}`, "result.error.error", msg)
	return out
}

func collectionNotFoundText(app, ctx string) string {
	return fmt.Sprintf(`application_id:%q context_id:%q is not found in db`, app, ctx)
}

func collectionKey(app, ctx string) string { return app + "/" + ctx }

func collectionDeviceSet(job gjson.Result) string {
	var ids []string
	for _, d := range job.Get("job_device_set.device_set.devices.device_ids").Array() {
		ids = append(ids, d.String())
	}
	sort.Strings(ids)
	return strings.Join(ids, ",") + "|" + job.Get("job_device_set.device_set.device_group").String()
}

func (s *collectionStub) change(w http.ResponseWriter, method string, body gjson.Result) {
	switch method {
	case http.MethodPut:
		job := body.Get("collection_job")
		app, ctx := job.Get("application_context.application_id").String(), job.Get("application_context.context_id").String()
		key := collectionKey(app, ctx)
		s.calls = append(s.calls, "PUT "+key)
		existing := s.jobs[key]
		if existing != nil && existing.deleted {
			existing = nil
		}
		wire := job.Get("collection_mode.collector_type").String()
		for _, in := range job.Get("sensor_input_configs").Array() {
			if in.Get("sensor_data.cli_sensor").Exists() && wire != "CLI_COLLECTOR" {
				fmt.Fprint(w, collectionRejected("*common_collection_data.SensorData_CliSensor paths need to be of CLI_COLLECTOR"))
				return
			}
			if in.Get("sensor_data.snmp_sensor").Exists() && wire != "SNMP_COLLECTOR" {
				fmt.Fprint(w, collectionRejected("*common_collection_data.SensorData_SnmpSensor paths need to be of SNMP_COLLECTOR"))
				return
			}
			if len(in.Get("sensor_data.snmp_sensor.snmp_mib.oids").Array()) > 1 && in.Get("sensor_data.snmp_sensor.snmp_mib.snmp_operation").String() == "WALK" {
				fmt.Fprint(w, collectionRejected("snmp TABLE/WALK operation only supports an OID list of size 1"))
				return
			}
		}
		if tag := job.Get("job_device_set.device_set.device_group").String(); tag != "" && !s.tags[tag] {
			fmt.Fprint(w, collectionRejected("group/tag "+tag+" not found in inventory or invalid"))
			return
		}
		if existing != nil {
			old := existing.body
			switch {
			case old.Get("collection_mode.collector_type").String() != wire:
				fmt.Fprint(w, collectionRejected(fmt.Sprintf("application_id:%q context_id:%q  collection mode and device set  is immutable", app, ctx)))
				return
			case collectionDeviceSet(old) != collectionDeviceSet(job):
				fmt.Fprint(w, collectionRejected(fmt.Sprintf("application_id:%q context_id:%q  device set is immutable", app, ctx)))
				return
			}
		}
		// the controller keeps the devices in the order it was sent; the test reads them back reversed to prove that the order is no change
		s.jobs[key] = &collectionStubJob{body: job}
		fmt.Fprint(w, collectionAccepted())
	case http.MethodDelete:
		app, ctx := body.Get("application_context.application_id").String(), body.Get("application_context.context_id").String()
		key := collectionKey(app, ctx)
		s.calls = append(s.calls, "DELETE "+key)
		job := s.jobs[key]
		if job == nil || job.deleted {
			fmt.Fprint(w, collectionRejected(collectionNotFoundText(app, ctx)))
			return
		}
		job.deleted, job.staleReads = true, s.stale
		fmt.Fprint(w, collectionAccepted())
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *collectionStub) query(w http.ResponseWriter, path string, body gjson.Result) {
	app, ctx := body.Get("application_context.application_id").String(), body.Get("application_context.context_id").String()
	key := collectionKey(app, ctx)
	job := s.jobs[key]
	list := map[string]string{collectionConfigPath: "collection_job_configuration_list", collectionStatusPath: "collection_job_status_list"}[path]
	empty := func(listKey string) string {
		return collectionRejectedList(listKey, collectionNotFoundText(app, ctx))
	}
	if path == collectionStatePath {
		if job == nil || job.deleted {
			fmt.Fprint(w, `{"collection_life_cycle_states":[],"query_options":{"page_token":"","page_size":500},"result":{"request_result":"ACCEPTED","error":{"error":""}}}`)
			return
		}
		row, _ := sjson.Set(`{"life_cycle_state":"SUCCESS_LIFE_CYCLE_STATE"}`, "application_context.application_id", app)
		row, _ = sjson.Set(row, "application_context.context_id", ctx)
		out, _ := sjson.SetRaw(`{"collection_life_cycle_states":[],"result":{"request_result":"ACCEPTED","error":{"error":""}}}`, "collection_life_cycle_states.-1", row)
		fmt.Fprint(w, out)
		return
	}
	if job == nil {
		fmt.Fprint(w, empty(list))
		return
	}
	if job.deleted {
		// the status query stops listing a deleted job at once; the configuration query takes more reads
		if path == collectionStatusPath || job.staleReads == 0 {
			fmt.Fprint(w, empty(list))
			return
		}
		job.staleReads--
	}
	if path == collectionStatusPath {
		row, _ := sjson.Set(`{"status":"READY","phase":"ACTIVE","progress":100,"job_error":{"error":""}}`, "application_context.application_id", app)
		row, _ = sjson.Set(row, "application_context.context_id", ctx)
		out, _ := sjson.SetRaw(`{"collection_job_status_list":[],"result":{"request_result":"ACCEPTED","error":{"error":""}}}`, "collection_job_status_list.-1", row)
		fmt.Fprint(w, out)
		return
	}
	cfg := job.body.Raw
	var ids []string
	for _, d := range job.body.Get("job_device_set.device_set.devices.device_ids").Array() {
		ids = append([]string{d.String()}, ids...)
	}
	if len(ids) > 0 {
		cfg, _ = sjson.Set(cfg, "job_device_set.device_set.devices.device_ids", ids)
	}
	out, _ := sjson.SetRaw(`{"collection_job_configuration_list":[],"query_options":null,"result":{"request_result":"ACCEPTED","error":{"error":""}}}`, "collection_job_configuration_list.-1", cfg)
	fmt.Fprint(w, out)
}

func collectionRejectedList(listKey, msg string) string {
	out, _ := sjson.SetRaw(`{"query_options":null}`, listKey, `[]`)
	out, _ = sjson.SetRaw(out, "result", collectionRejected(msg)[len(`{"result":`):len(collectionRejected(msg))-1])
	return out
}

func (s *collectionStub) seed(app, ctx, wire, devices, sensor string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	job := `{"collection_mode":{"lifetime_type":"APPLICATION_MANAGED","schedule_id":"","n_collections":0,"track_collect_n_job":false}}`
	job, _ = sjson.Set(job, "application_context.application_id", app)
	job, _ = sjson.Set(job, "application_context.context_id", ctx)
	job, _ = sjson.Set(job, "collection_mode.collector_type", wire)
	job, _ = sjson.SetRaw(job, "job_device_set", devices)
	in, _ := sjson.SetRaw(`{"cadence_in_millisec":"60000"}`, "sensor_data", sensor)
	job, _ = sjson.SetRaw(job, "sensor_input_configs", "["+in+"]")
	out, _ := sjson.SetRaw(`{"filter_profile_id":""}`, "sensor_data", sensor)
	out, _ = sjson.Set(out, "destination.destination_id", "dest-id-1")
	out, _ = sjson.Set(out, "destination.destination_name", "cwtest-dest-kafka")
	out, _ = sjson.Set(out, "destination.context_id", "topic-seed")
	job, _ = sjson.SetRaw(job, "sensor_output_configs", "["+out+"]")
	s.jobs[collectionKey(app, ctx)] = &collectionStubJob{body: gjson.Parse(job)}
}

func (s *collectionStub) remove(app, ctx string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.jobs, collectionKey(app, ctx))
}

func (s *collectionStub) callLog() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.calls...)
}

func (s *collectionStub) resetCalls() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = nil
}

func (s *collectionStub) providerBlock(allowExisting bool) string {
	return fmt.Sprintf(`provider "crosswork" {
  url                      = %q
  username                 = "admin"
  password                 = "stub-password"
  insecure                 = true
  retries                  = 0
  allow_existing_on_create = %t
}
`, s.URL, allowExisting)
}

func fastCollectionWait(t *testing.T) {
	old := collectionGoneInterval
	collectionGoneInterval = time.Millisecond
	t.Cleanup(func() { collectionGoneInterval = old })
}

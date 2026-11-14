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
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// The template id is an integer and the controller sorts variables and series and shuffles tags; a job answers with an id
// and ends later, Success even when a device failed inside the error budget; a re-apply and a delete of an application are
// synchronous (an unknown one is a 404 entry inside a 207); an applied template is refused with a 409 unless deleted with force.
type gcStub struct {
	*ssoStub

	mu        sync.Mutex
	templates map[int]map[string]any
	nextID    int
	gvars     map[string]map[string]any
	nodes     map[string]*stubNode
	apps      map[string]*stubApp
	jobs      map[string]*stubJob
	calls     []string
	// failDevices makes the job fail for a device, with the message; failReapply does the same for a re-apply.
	failDevices map[string]string
	// readyMode and readyTimeout are the provider settings device_ready_mode and device_ready_timeout (seconds); the
	// timeout is 0 by default, so a device that is not ready fails at once.
	readyMode    string
	readyTimeout int
	// postDelay is how long the job POST takes; maxRunning is the most jobs that were running at once (a job runs from
	// its POST until the status read after finishAfter).
	postDelay  time.Duration
	maxRunning int
	// stuckJob is the number (1, 2, ...) of the job that never finishes; 0 for none.
	stuckJob    int
	failReapply map[string]string
	// finishAfter is the number of status reads before a job ends.
	finishAfter int
	jobSeq      int
	// secrets records every variable value a body carried, to prove a write-only value was sent.
	secrets []string
}

type stubNode struct {
	Name, UUID, Host, Reach, NSO, Admin string
}

type stubApp struct {
	Name       string
	Template   int
	Device     string
	UUID       string
	Variables  map[string]string // job values, as the last job sent them
	Applied    int               // jobs and re-applies it took
	Conformant bool
}

type stubJob struct {
	ID      string
	Name    string
	Devices []string
	Result  map[string]string // device name to Success or Failed
	Message map[string]string // device name to the NSO message
	Reads   int
}

func newGCStub(t *testing.T) *gcStub {
	t.Helper()
	s := &gcStub{ssoStub: newSSOStub(t), templates: map[int]map[string]any{}, nextID: 1, gvars: map[string]map[string]any{},
		nodes: map[string]*stubNode{}, apps: map[string]*stubApp{}, jobs: map[string]*stubJob{}, failDevices: map[string]string{},
		failReapply: map[string]string{}, finishAfter: 1}
	for _, n := range []string{"d1", "d2", "d3"} {
		s.nodes[n] = &stubNode{Name: n, UUID: "uuid-" + n, Host: "host-" + n, Reach: nsoReady, NSO: nsoSynced, Admin: nsoAdminUp}
	}
	sso := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/crosswork/cwms/") || r.URL.Path == nsoNodesQuery {
			s.serve(w, r)
			return
		}
		sso.ServeHTTP(w, r)
	})
	oldPoll, oldReady, oldWorkers := gcJobPollInterval, nsoReadinessInterval, gcDeleteWorkers
	gcJobPollInterval, nsoReadinessInterval, gcDeleteWorkers = 5*time.Millisecond, 5*time.Millisecond, 2
	t.Cleanup(func() {
		gcJobPollInterval, nsoReadinessInterval, gcDeleteWorkers = oldPoll, oldReady, oldWorkers
		Shutdown()
	})
	return s
}

func (s *gcStub) log(format string, a ...any) { s.calls = append(s.calls, fmt.Sprintf(format, a...)) }

func (s *gcStub) callLog() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return strings.Join(s.calls, " | ")
}

func (s *gcStub) callsWith(prefix string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []string
	for _, c := range s.calls {
		if strings.HasPrefix(c, prefix) {
			out = append(out, c)
		}
	}
	return out
}

func (s *gcStub) appNames() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var out []string
	for n := range s.apps {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

func (s *gcStub) providerBlock() string {
	mode := ""
	if s.readyMode != "" {
		mode = fmt.Sprintf("\n  device_ready_mode    = %q", s.readyMode)
	}
	return fmt.Sprintf(`provider "crosswork" {
  url      = %q
  username = "admin"
  password = "stub-password"
  insecure = true
  retries  = 0
  device_ready_timeout = %d%s
}
`, s.URL, s.readyTimeout, mode)
}

func write(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func entry(code int, message, key string) map[string]any {
	status := "success"
	if code >= 400 {
		status = "failed"
	}
	return map[string]any{"code": code, "message": message, "status": status, "key": key}
}

func page(rows any) map[string]any {
	return map[string]any{"data": rows, "pagination": map[string]any{"pageNumber": 0, "pageSize": 100, "total": 0, "totalPages": 1, "first": true, "last": true}}
}

func (s *gcStub) serve(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	s.mu.Lock()
	defer s.mu.Unlock()
	p := r.URL.Path
	var body map[string]any
	var list []any
	_ = json.Unmarshal(raw, &body)
	_ = json.Unmarshal(raw, &list)
	switch {
	case p == nsoNodesQuery:
		var rows []any
		for _, n := range sortedKeys(s.nodes) {
			x := s.nodes[n]
			rows = append(rows, map[string]any{"inventory_id": x.Name, "uuid": x.UUID, "host_name": x.Host, "reachability_state": x.Reach, "nso_state": x.NSO, "admin_state": x.Admin})
		}
		write(w, 200, map[string]any{"data": rows})

	case p == gcPath+"/templates" && r.Method == http.MethodPost:
		name, _ := body["name"].(string)
		for _, t := range s.templates {
			if t["name"] == name {
				write(w, 409, entry(409, "Template already exists", ""))
				return
			}
		}
		id := s.nextID
		s.nextID++
		s.templates[id] = body
		s.recordVars(body)
		s.log("POST template %s", name)
		write(w, 200, entry(200, "Template created", strconv.Itoa(id)))
	case p == gcPath+"/templates/query":
		var rows []any
		for _, id := range sortedInts(s.templates) {
			t := s.templates[id]
			if n, _ := body["name"].(string); n != "" && !strings.Contains(t["name"].(string), n) {
				continue
			}
			rows = append(rows, map[string]any{"id": id, "tcId": tcIDOf(id), "name": t["name"]})
		}
		write(w, 200, page(rows))
	case strings.HasPrefix(p, gcPath+"/templates/") && r.Method == http.MethodGet:
		id, _ := strconv.Atoi(strings.TrimPrefix(p, gcPath+"/templates/"))
		t, ok := s.templates[id]
		if !ok {
			write(w, 404, entry(404, "Template not found", ""))
			return
		}
		write(w, 200, s.templateAnswer(id, t))
	case strings.HasPrefix(p, gcPath+"/templates/") && r.Method == http.MethodPut:
		id, _ := strconv.Atoi(strings.TrimPrefix(p, gcPath+"/templates/"))
		if _, ok := s.templates[id]; !ok {
			write(w, 404, entry(404, "Template not found", ""))
			return
		}
		s.templates[id] = body
		s.recordVars(body)
		s.log("PUT template %d", id)
		write(w, 200, entry(200, "Template updated", strconv.Itoa(id)))
	case p == gcPath+"/templates" && r.Method == http.MethodDelete:
		id, _ := strconv.Atoi(r.URL.Query().Get("ids"))
		force := r.URL.Query().Get("force") == "true"
		s.log("DELETE template %d force=%v", id, force)
		if _, ok := s.templates[id]; !ok {
			write(w, 207, []any{entry(404, "Template not found", strconv.Itoa(id))})
			return
		}
		var using []string
		for _, a := range s.apps {
			if a.Template == id {
				using = append(using, a.Device)
			}
		}
		if len(using) > 0 && !force {
			sort.Strings(using)
			write(w, 207, []any{entry(409, "Cannot delete template as it is associated with devices: "+strings.Join(using, ", "), strconv.Itoa(id))})
			return
		}
		for n, a := range s.apps {
			if a.Template == id {
				delete(s.apps, n)
			}
		}
		delete(s.templates, id)
		write(w, 207, []any{entry(200, "Template deleted", strconv.Itoa(id))})

	case p == gcPath+"/gvars" && r.Method == http.MethodPost:
		v := list[0].(map[string]any)
		name := v["name"].(string)
		if _, ok := s.gvars[name]; ok {
			write(w, 207, []any{entry(409, "Global variable already exists", name)})
			return
		}
		v["id"] = strconv.Itoa(len(s.gvars) + 10)
		s.gvars[name] = v
		s.record(v["value"])
		s.log("POST gvar %s", name)
		write(w, 207, []any{entry(200, "Global variable "+name+" created", v["id"].(string))})
	case p == gcPath+"/gvars" && r.Method == http.MethodPatch:
		v := list[0].(map[string]any)
		name := v["name"].(string)
		cur, ok := s.gvars[name]
		if !ok {
			write(w, 207, []any{entry(404, "Global variable not found", name)})
			return
		}
		cur["value"], cur["type"] = v["value"], v["type"]
		s.record(v["value"])
		s.log("PATCH gvar %s", name)
		write(w, 207, []any{entry(200, "Global variable updated", name)})
	case p == gcPath+"/gvars/query":
		var rows []any
		for _, n := range sortedKeys(s.gvars) {
			if want, _ := body["name"].(string); want != "" && !strings.Contains(n, want) {
				continue
			}
			g := map[string]any{}
			for k, v := range s.gvars[n] {
				g[k] = v
			}
			g["inUse"] = false
			rows = append(rows, g)
		}
		write(w, 200, page(rows))
	case p == gcPath+"/gvars" && r.Method == http.MethodDelete:
		id := r.URL.Query().Get("ids")
		for n, g := range s.gvars {
			if g["id"] == id {
				delete(s.gvars, n)
				s.log("DELETE gvar %s", n)
				write(w, 207, []any{entry(200, "Global Variable deleted", id)})
				return
			}
		}
		write(w, 207, []any{entry(404, "Global variable not found", id)})

	case p == gcMopPath+"/summary":
		write(w, 200, page([]any{map[string]any{"id": 4, "name": "Default Golden Config - Cisco", "app": gcApp, "vendor": "Cisco Systems"}}))

	case p == gcPath+"/applications/job":
		s.newJob(w, body)
	case strings.HasPrefix(p, gcPath+"/template/") && strings.HasSuffix(p, "/devices/query"):
		id, _ := strconv.Atoi(strings.Split(strings.TrimPrefix(p, gcPath+"/template/"), "/")[0])
		var rows []any
		var names []string
		for n, a := range s.apps {
			if a.Template == id {
				names = append(names, n)
			}
		}
		sort.Strings(names)
		for _, n := range names {
			a := s.apps[n]
			rows = append(rows, map[string]any{"templateName": s.templates[id]["name"], "applicationName": a.Name, "deviceName": a.Device, "deviceUuid": a.UUID,
				"hostName": "host-" + a.Device, "conformant": a.Conformant})
		}
		write(w, 200, page(rows))
	case p == gcPath+"/applications/re-apply":
		var names []string
		for _, n := range list {
			names = append(names, n.(string))
		}
		s.log("re-apply %s", strings.Join(s.deviceNamesOf(names), ","))
		var out []any
		failed := false
		for _, n := range names {
			a, ok := s.apps[n]
			switch {
			case !ok:
				out = append(out, entry(404, "Application not found", n))
				failed = true
			case s.failReapply[a.Device] != "":
				out = append(out, entry(400, s.failReapply[a.Device], a.UUID))
				failed = true
			default:
				a.Applied++
				out = append(out, entry(200, "Template re-applied", a.UUID))
			}
		}
		_ = failed
		write(w, 207, out)
	case strings.HasPrefix(p, gcPath+"/applications/") && r.Method == http.MethodDelete:
		name := strings.TrimPrefix(p, gcPath+"/applications/")
		s.log("DELETE app %s", name)
		if _, ok := s.apps[name]; !ok {
			write(w, 207, []any{entry(404, "Application not found", name)})
			return
		}
		delete(s.apps, name)
		write(w, 207, []any{entry(200, "Application deleted", name)})

	case strings.HasPrefix(p, mopStatus):
		s.jobStatus(w, strings.TrimPrefix(p, mopStatus+"/"))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

const mopStatus = "/crosswork/cwms/automations/v1/mop/status/job"

func sortedInts[V any](m map[int]V) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Ints(out)
	return out
}

// tcIDOf differs from the template id, as when the two id sequences have drifted apart.
func tcIDOf(id int) int { return id + 100 }

func templateOfTcID(s *gcStub, tc int) int {
	for id := range s.templates {
		if tcIDOf(id) == tc {
			return id
		}
	}
	return 0
}

func (s *gcStub) record(v any) {
	if str, ok := v.(string); ok {
		s.secrets = append(s.secrets, str)
	}
}

func (s *gcStub) recordVars(body map[string]any) {
	vars, _ := body["variables"].([]any)
	for _, v := range vars {
		s.record(v.(map[string]any)["value"])
	}
}

func (s *gcStub) deviceNamesOf(apps []string) []string {
	var out []string
	for _, n := range apps {
		if a, ok := s.apps[n]; ok {
			out = append(out, a.Device)
		}
	}
	sort.Strings(out)
	return out
}

// templateAnswer: the controller sorts variables by name and series alphabetically, and keeps tags in an order of its own
// (here reversed).
func (s *gcStub) templateAnswer(id int, t map[string]any) map[string]any {
	vars, _ := t["variables"].([]any)
	sort.Slice(vars, func(i, j int) bool {
		return vars[i].(map[string]any)["name"].(string) < vars[j].(map[string]any)["name"].(string)
	})
	for _, v := range vars {
		m := v.(map[string]any)
		s.record(m["value"])
		if str, ok := m["value"].(string); ok && m["type"] == "JSON" {
			var x any
			if json.Unmarshal([]byte(str), &x) == nil {
				m["value"] = x
			}
		}
	}
	series := toStrings(t["productSeries"])
	sort.Strings(series)
	tags := toStrings(t["tags"])
	for i, j := 0, len(tags)-1; i < j; i, j = i+1, j-1 {
		tags[i], tags[j] = tags[j], tags[i]
	}
	return map[string]any{"id": id, "name": t["name"], "description": t["description"], "tags": tags, "vendor": t["vendor"],
		"softwareType": t["softwareType"], "productSeries": series, "audit": t["audit"],
		"versions": []any{map[string]any{"tcId": tcIDOf(id), "version": "1.0", "configType": t["configType"], "config": t["config"], "mode": t["mode"], "variables": vars}}}
}

func toStrings(v any) []string {
	var out []string
	for _, x := range v.([]any) {
		out = append(out, x.(string))
	}
	return out
}

func (s *gcStub) newJob(w http.ResponseWriter, body map[string]any) {
	time.Sleep(s.postDelay) // the lock of the stub is held, so a second POST waits behind this one
	s.jobSeq++
	j := &stubJob{ID: strconv.Itoa(s.jobSeq), Name: body["name"].(string), Result: map[string]string{}, Message: map[string]string{}}
	tpl := body["data"].(map[string]any)["templates"].([]any)[0].(map[string]any)
	// A job names the id of the configuration of the template (tcId), which is not the id of the template; a tcId the
	// controller does not know ends in an Optional.get() on an empty result, an HTTP 500.
	tid := templateOfTcID(s, int(tpl["tcId"].(float64)))
	if tid == 0 {
		s.jobSeq--
		write(w, 500, map[string]any{"code": 500, "message": "An internal server error occurred", "status": "failed"})
		return
	}
	vars := map[string]string{}
	for _, v := range tpl["variables"].([]any) {
		m := v.(map[string]any)
		val := fmt.Sprint(m["value"])
		if o, ok := m["value"].(map[string]any); ok {
			b, _ := json.Marshal(o)
			val = string(b)
		}
		vars[m["name"].(string)] = val
		s.record(val)
	}
	byUUID := map[string]*stubNode{}
	for _, n := range s.nodes {
		byUUID[n.UUID] = n
	}
	var names []string
	for _, ri := range body["resourceInfo"].([]any) {
		n := byUUID[ri.(map[string]any)["uuid"].(string)]
		j.Devices = append(j.Devices, n.Name)
		names = append(names, n.Name)
		if msg := s.failDevices[n.Name]; msg != "" {
			j.Result[n.Name], j.Message[n.Name] = "Failed", msg
			continue
		}
		j.Result[n.Name] = "Success"
		name := fmt.Sprintf("%s-%s", s.templates[tid]["name"], n.UUID)
		a, ok := s.apps[name]
		if !ok {
			a = &stubApp{Name: name, Template: tid, Device: n.Name, UUID: n.UUID}
			s.apps[name] = a
		}
		a.Variables = vars
		a.Applied++
		a.Conformant = true
	}
	sort.Strings(names)
	s.log("job %s %s", j.Name, strings.Join(names, ","))
	s.jobs[j.ID] = j
	s.noteRunning()
	write(w, 200, entry(200, "Application JOB created", j.ID))
}

func (j *stubJob) status() (string, string) {
	failed, ok := 0, 0
	for _, r := range j.Result {
		if r == "Failed" {
			failed++
		} else {
			ok++
		}
	}
	if ok == 0 {
		return "Failed", "The Golden Config job has failed because all the devices in the job have failed."
	}
	return "Success", ""
}

func (s *gcStub) noteRunning() {
	running := 0
	for _, j := range s.jobs {
		if j.Reads <= s.finishAfter {
			running++
		}
	}
	if running > s.maxRunning {
		s.maxRunning = running
	}
}

func (s *gcStub) jobStatus(w http.ResponseWriter, rest string) {
	parts := strings.Split(rest, "/")
	j, ok := s.jobs[parts[0]]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch {
	case len(parts) == 1:
		j.Reads++
		st, msg := j.status()
		if j.Reads <= s.finishAfter || (s.stuckJob != 0 && j.ID == strconv.Itoa(s.stuckJob)) {
			st, msg = "Running", ""
		}
		s.noteRunning()
		write(w, 200, map[string]any{"jobId": j.ID, "jobName": j.Name, "jobStatus": st, "message": msg,
			"jobInput": map[string]any{"data": "echoed with the values"}})
	case len(parts) == 2 && parts[1] == "resource":
		var rows []any
		for i, d := range j.Devices {
			// The name of a result is the host name of the device, not its inventory id.
			rows = append(rows, map[string]any{"name": s.nodes[d].Host, "result": j.Result[d], "deviceUuid": s.nodes[d].UUID,
				"describeURI": fmt.Sprintf("http://cwm-solutions-automation:8080/v1/mop/status/job/%s/resource/%d", j.ID, i)})
		}
		write(w, 200, map[string]any{"jobId": j.ID, "resources": rows, "pagination": map[string]any{"total": len(rows), "pageNumber": 1, "pageSize": 100}})
	case len(parts) == 3 && parts[1] == "resource":
		i, _ := strconv.Atoi(parts[2])
		d := j.Devices[i]
		msg := "Apply templates successful"
		if j.Result[d] == "Failed" {
			b, _ := json.Marshal([]any{map[string]any{"template": "t", "status": "failed", "message": j.Message[d]}})
			msg = string(b)
		}
		write(w, 200, map[string]any{"name": s.nodes[d].Host, "deviceUuid": s.nodes[d].UUID, "result": j.Result[d], "message": "",
			"activities": []any{map[string]any{"name": "Apply Config", "stage": "configuration", "result": j.Result[d], "message": msg}}})
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

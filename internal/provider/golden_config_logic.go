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
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

const (
	gcPath    = "/crosswork/cwms/gc/v1"
	gcMopPath = "/crosswork/cwms/mop/v1"
	gcApp     = "Golden Config"

	// gcLockKey names the lock that the Golden Configuration writes share: NSO refuses two overlapping pushes of
	// the templates and global variables.
	gcLockKey = "golden-config"

	// The longest job name the controller accepts, and the room of the hash suffix.
	gcJobNameMax = 64
	gcHashLen    = 8
)

// The waits are variables so the tests can shorten them.
var (
	// gcJobPostTimeout is the request timeout of the call that starts a job; its body lists every device.
	gcJobPostTimeout  = 180 * time.Second
	gcJobPollInterval = 5 * time.Second
	gcDeleteWorkers   = 5
	// gcSerialize is a variable so the test of the lock can turn it off.
	gcSerialize = true
)

// A secret has a version and its value is known only from the configuration.
type gcVar struct {
	Name    string
	Type    string
	Value   string
	Secret  bool
	Version int64
}

// The controller keeps a JSON value with sorted keys and no spaces; any other value is unchanged.
func gcCanonical(typ, value string) string {
	if typ != "JSON" {
		return value
	}
	var v any
	if err := json.Unmarshal([]byte(value), &v); err != nil {
		return value
	}
	b, err := json.Marshal(v)
	if err != nil {
		return value
	}
	return string(b)
}

func gcSameValue(typ, a, b string) bool {
	return gcCanonical(typ, a) == gcCanonical(typ, b)
}

func gcBodyValue(typ, value string) any {
	if typ == "JSON" {
		var v any
		if err := json.Unmarshal([]byte(value), &v); err == nil {
			return v
		}
	}
	return value
}

func gcValidJSON(value string) bool {
	var v any
	return json.Unmarshal([]byte(value), &v) == nil
}

func gcValueString(typ string, v gjson.Result) string {
	if typ == "JSON" && v.Type != gjson.String {
		return gcCanonical(typ, v.Raw)
	}
	return v.String()
}

// The version stands for the value of a secret.
func (v gcVar) nameKey() string {
	if v.Secret {
		return fmt.Sprintf("%s|%s|#%d", v.Name, v.Type, v.Version)
	}
	return fmt.Sprintf("%s|%s|=%s", v.Name, v.Type, gcCanonical(v.Type, v.Value))
}

// groupKey identifies a variable with its value, a secret included, so two devices with different secrets
// of the same version never share a job.
func (v gcVar) groupKey() string {
	if v.Secret {
		return fmt.Sprintf("%s|%s|#%d|=%s", v.Name, v.Type, v.Version, v.Value)
	}
	return v.nameKey()
}

func gcSortedVars(vars []gcVar) []gcVar {
	out := append([]gcVar(nil), vars...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

func gcVarsNameKey(vars []gcVar) string {
	var keys []string
	for _, v := range gcSortedVars(vars) {
		keys = append(keys, v.nameKey())
	}
	return strings.Join(keys, "\n")
}

func gcVarsGroupKey(vars []gcVar) string {
	var keys []string
	for _, v := range gcSortedVars(vars) {
		keys = append(keys, v.groupKey())
	}
	return strings.Join(keys, "\n")
}

func gcVarsBody(vars []gcVar) []any {
	out := make([]any, 0, len(vars))
	for _, v := range vars {
		out = append(out, map[string]any{"name": v.Name, "type": v.Type, "value": gcBodyValue(v.Type, v.Value)})
	}
	return out
}

func gcJSON(v any) string {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(v)
	return strings.TrimSuffix(buf.String(), "\n")
}

type gcTemplate struct {
	Name          string
	Description   string
	Tags          []string
	Mode          string
	ConfigType    string
	Vendor        string
	SoftwareType  string
	ProductSeries []string
	Audit         bool
	Config        string
	Variables     []gcVar
}

func (t gcTemplate) body() string {
	m := map[string]any{
		"name": t.Name, "description": t.Description, "config": t.Config, "vendor": t.Vendor,
		"softwareType": t.SoftwareType, "productSeries": nonNil(t.ProductSeries), "tags": nonNil(t.Tags),
		"configType": t.ConfigType, "mode": t.Mode, "audit": t.Audit, "variables": gcVarsBody(t.Variables),
	}
	return gcJSON(m)
}

func nonNil(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}

func gcFindTemplate(c *cw.Client, name string) (int64, bool, error) {
	body := gcJSON(map[string]any{"name": name})
	rows, err := c.Paginate(&cw.BodyPager{Path: gcPath + "/templates/query", Body: body, Size: 100}, cw.NoLogPayload)
	if err != nil {
		return 0, false, err
	}
	for _, r := range rows {
		if r.Get("name").String() == name {
			return r.Get("id").Int(), true, nil
		}
	}
	return 0, false, nil
}

// The answer holds variable values in clear text.
func gcGetTemplate(c *cw.Client, id int64) (gjson.Result, error) {
	return c.Get(gcPath+"/templates/"+strconv.FormatInt(id, 10), cw.NoLogPayload)
}

func gcFindGlobalVariable(c *cw.Client, name string) (gjson.Result, bool, error) {
	body := gcJSON(map[string]any{"name": name})
	rows, err := c.Paginate(&cw.BodyPager{Path: gcPath + "/gvars/query", Body: body, Size: 100}, cw.NoLogPayload)
	if err != nil {
		return gjson.Result{}, false, err
	}
	for _, r := range rows {
		if r.Get("name").String() == name {
			return r, true, nil
		}
	}
	return gjson.Result{}, false, nil
}

type gcApplication struct {
	DeviceName      string
	DeviceUUID      string
	HostName        string
	ApplicationName string
	Conformant      bool
}

func gcTemplateApplications(c *cw.Client, templateID int64) ([]gcApplication, error) {
	rows, err := c.Paginate(&cw.BodyPager{Path: fmt.Sprintf("%s/template/%d/devices/query", gcPath, templateID), Body: "{}", Size: 100}, cw.NoLogPayload)
	if err != nil {
		return nil, err
	}
	out := make([]gcApplication, 0, len(rows))
	for _, r := range rows {
		out = append(out, gcApplication{DeviceName: r.Get("deviceName").String(), DeviceUUID: r.Get("deviceUuid").String(),
			HostName: r.Get("hostName").String(), ApplicationName: r.Get("applicationName").String(), Conformant: r.Get("conformant").Bool()})
	}
	return out, nil
}

// gcMopID resolves the name of a MOP to its id. The MOP lists count pages from 1.
func gcMopID(c *cw.Client, name string) (int64, error) {
	rows, err := c.Paginate(&cw.BodyPager{Path: gcMopPath + "/summary", Size: 100, FirstPage: 1})
	if err != nil {
		return 0, err
	}
	var names []string
	for _, r := range rows {
		if r.Get("name").String() == name {
			return r.Get("id").Int(), nil
		}
		names = append(names, r.Get("name").String())
	}
	sort.Strings(names)
	return 0, fmt.Errorf("no MOP named %q, the controller has: %s", name, strings.Join(names, ", "))
}

type gcGroup struct {
	Name    string
	Devices []string
	Vars    []gcVar
}

// The first group keeps the template name; the others get the name plus a short hash of the non-secret values and the
// secret versions. Secrets are compared by their value.
func gcMakeGroups(template string, devices []string, vars map[string][]gcVar) []gcGroup {
	names := append([]string(nil), devices...)
	sort.Strings(names)
	byKey := map[string]*gcGroup{}
	var order []string
	for _, d := range names {
		k := gcVarsGroupKey(vars[d])
		g, ok := byKey[k]
		if !ok {
			g = &gcGroup{Vars: vars[d]}
			byKey[k] = g
			order = append(order, k)
		}
		g.Devices = append(g.Devices, d)
	}
	groups := make([]gcGroup, 0, len(order))
	for i, k := range order {
		g := *byKey[k]
		if i == 0 {
			g.Name = template
		} else {
			g.Name = gcJobName(template, gcVarsNameKey(g.Vars))
		}
		groups = append(groups, g)
	}
	return groups
}

func gcJobName(template, key string) string {
	sum := sha256.Sum256([]byte(key))
	suffix := hex.EncodeToString(sum[:])[:gcHashLen]
	room := gcJobNameMax - gcHashLen - 1
	if len(template) > room {
		template = template[:room]
	}
	return template + "-" + suffix
}

type gcJobParams struct {
	// ConfigID is the tcId of the job body: the id of the configuration of the template, not of the template.
	ConfigID     int64
	TemplateName string
	Vendor       string
	Series       string
	Mop          string
	MopID        int64
	Batch        *int64
	ErrorBudget  *int64
	Tags         []string
}

// The device values go as job level variables, without a source, so they beat the template and the global values.
func gcJobBody(p gcJobParams, g gcGroup, nodes map[string]nsoNode) string {
	info := make([]any, 0, len(g.Devices))
	for _, d := range g.Devices {
		n := nodes[d]
		// hostName is the inventory id.
		info = append(info, map[string]any{"uuid": n.UUID, "hostName": n.Name})
	}
	m := map[string]any{
		"name": g.Name, "app": gcApp, "mopName": p.Mop, "mopId": p.MopID, "vendor": p.Vendor, "productSeries": p.Series,
		"resourceInfo": info, "tags": nonNil(p.Tags),
		"data": map[string]any{"templates": []any{map[string]any{"tcId": p.ConfigID, "templateName": p.TemplateName, "variables": gcVarsBody(g.Vars)}}},
	}
	if p.Batch != nil {
		m["batchSize"] = *p.Batch
	}
	if p.ErrorBudget != nil {
		m["errorBudget"] = *p.ErrorBudget
	}
	return gcJSON(m)
}

type gcJobResult struct {
	Group  gcGroup
	JobID  string
	Failed map[string]string
}

// A device whose result is Success has no entry, whatever the job status says for the others.
//
// A job holds the lock (gcLockKey) from its POST until it has finished, so the timeout of each job starts when
// the job holds the lock.
func gcRunJobs(ctx context.Context, c *cw.Client, p gcJobParams, groups []gcGroup, nodes map[string]nsoNode, timeout time.Duration) []gcJobResult {
	results := make([]gcJobResult, 0, len(groups))
	for _, g := range groups {
		r := gcJobResult{Group: g, Failed: map[string]string{}}
		mods := append([]func(*cw.Req){cw.NoLogPayload, cw.Timeout(int(gcJobPostTimeout.Seconds())), c.WaitForMopJob(gcJobPollInterval, timeout)}, gcLock()...)
		res, err := c.Post(gcPath+"/applications/job", gcJobBody(p, g, nodes), mods...)
		r.JobID = res.Get("key").String()
		var waitErr *cw.WaitError
		switch {
		case err == nil:
		case errors.As(err, &waitErr) && r.JobID != "":
			// The job was started and did not finish: it keeps running on the controller.
			for _, d := range g.Devices {
				r.Failed[d] = fmt.Sprintf("job %s (%s): %s; it keeps running on the controller", r.JobID, g.Name, waitErr)
			}
		case errors.As(err, &waitErr):
			for _, d := range g.Devices {
				r.Failed[d] = fmt.Sprintf("the job %q answered without a job id", g.Name)
			}
		default:
			for _, d := range g.Devices {
				r.Failed[d] = fmt.Sprintf("the job %q could not be started: %s", g.Name, err)
			}
		}
		if err == nil {
			gcJobDevices(c, &r, nodes)
		}
		results = append(results, r)
	}
	return results
}

func gcJobDevices(c *cw.Client, r *gcJobResult, nodes map[string]nsoNode) {
	job, err := c.MopJob(r.JobID, cw.NoLogPayload)
	if err != nil {
		for _, d := range r.Group.Devices {
			r.Failed[d] = fmt.Sprintf("job %s (%s): the result could not be read: %s", r.JobID, r.Group.Name, err)
		}
		return
	}
	out, err := c.MopJobOutcome(r.JobID, cw.NoLogPayload)
	if err != nil {
		for _, d := range r.Group.Devices {
			r.Failed[d] = fmt.Sprintf("job %s (%s): the result could not be read: %s", r.JobID, r.Group.Name, err)
		}
		return
	}
	// A device is found by its uuid: the name of a result is the host name, which is not the inventory id.
	ok := map[string]bool{}
	for _, res := range out.Resources {
		if res.Result == cw.MopSuccess {
			ok[res.DeviceUUID] = true
		}
	}
	byDevice := map[string]cw.MopJobFailure{}
	for _, f := range out.Failures {
		byDevice[f.DeviceUUID] = f
	}
	for _, d := range r.Group.Devices {
		uuid := nodes[d].UUID
		if ok[uuid] {
			continue
		}
		reason := fmt.Sprintf("job %s (%s) ended %s", r.JobID, r.Group.Name, job.Status)
		if job.Message != "" {
			reason += ": " + job.Message
		}
		if f, found := byDevice[uuid]; found {
			reason = fmt.Sprintf("job %s (%s): %s", r.JobID, r.Group.Name, f.String())
		}
		r.Failed[d] = reason
	}
}

func gcLock() []func(*cw.Req) {
	if !gcSerialize {
		return nil
	}
	return []func(*cw.Req){cw.UseMutexKey(gcLockKey)}
}

func gcReapply(c *cw.Client, applications []gcApplication, timeout time.Duration) map[string]string {
	failed := map[string]string{}
	if len(applications) == 0 {
		return failed
	}
	names := make([]string, 0, len(applications))
	for _, a := range applications {
		names = append(names, a.ApplicationName)
	}
	res, err := c.Post(gcPath+"/applications/re-apply", gcJSON(names), append([]func(*cw.Req){cw.NoLogPayload, cw.Timeout(int(timeout.Seconds()))}, gcLock()...)...)
	if err == nil {
		return failed
	}
	var ae *cw.APIError
	if asAPIError(err, &ae) && len(ae.EntryStatus) > 0 {
		res.ForEach(func(_, e gjson.Result) bool {
			if e.Get("code").Int() >= 400 {
				failed[e.Get("key").String()] = e.Get("message").String()
			}
			return true
		})
		if len(failed) > 0 {
			return failed
		}
	}
	for _, a := range applications {
		failed[a.DeviceUUID] = fmt.Sprintf("the re-apply failed: %s", err)
	}
	return failed
}

// An application that is already gone is a success.
func gcDeleteApplications(c *cw.Client, applications []gcApplication) map[string]error {
	failed := map[string]error{}
	var mu sync.Mutex
	work := make(chan gcApplication)
	var wg sync.WaitGroup
	workers := gcDeleteWorkers
	if workers < 1 {
		workers = 1
	}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for a := range work {
				if _, err := c.Delete(gcPath+"/applications/"+url.PathEscape(a.ApplicationName), append([]func(*cw.Req){cw.NoLogPayload}, gcLock()...)...); err != nil && !cw.IsNotFound(err) {
					mu.Lock()
					failed[a.DeviceName] = err
					mu.Unlock()
				}
			}
		}()
	}
	for _, a := range applications {
		work <- a
	}
	close(work)
	wg.Wait()
	return failed
}

func asAPIError(err error, target **cw.APIError) bool {
	for err != nil {
		if ae, ok := err.(*cw.APIError); ok {
			*target = ae
			return true
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return false
		}
		err = u.Unwrap()
	}
	return false
}

func sortedKeys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func gcFailureMessage(failed map[string]string) string {
	lines := make([]string, 0, len(failed))
	for _, d := range sortedKeys(failed) {
		lines = append(lines, "  "+d+": "+failed[d])
	}
	return strings.Join(lines, "\n")
}

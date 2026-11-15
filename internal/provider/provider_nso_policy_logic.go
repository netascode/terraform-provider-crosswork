// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License Version 2.0 (the "License"),
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
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cw "github.com/netascode/go-crosswork"
)

// The device onboarding policy of the NSO providers (Manage Provider Access, Edit Policy Details). The controller holds one
// policy document for all NSO providers, and a PUT replaces all of it: an entry that the body leaves out is removed, and so are
// the providers criteria and the LSA policy block. A resource therefore reads the document, replaces the entry of its provider
// and sends everything else back as it was.

const (
	providerNSOPolicyQueryPath = "/crosswork/inventory/v1/nso/policy/query"
	providerNSOPolicyPath      = "/crosswork/inventory/v1/nso/policy"
	nsoJobsQueryPath           = "/crosswork/inventory/v1/jobs/query"
)

// providerNSOPolicyFlags are the flags of an entry. The controller stores a false flag but reads it back as missing.
var providerNSOPolicyFlags = []string{"match", "onboardTo", "onboardFrom", "syncFrom", "checkSync"}

// A variable so that tests can shorten it.
var providerNSOPolicyJobTimeout = 60 * time.Second

func providerNSOPolicyEntry(doc gjson.Result, name string) (gjson.Result, bool) {
	entry, ok := doc.Get("provider_policy").Map()[name]
	return entry, ok && entry.IsObject()
}

// The error is a not found when the document has no entry for the provider, which is how the controller reports a
// provider that was removed (the policy goes with it).
func flattenProviderNSOPolicyEntry(doc gjson.Result, name string) (gjson.Result, error) {
	entry, ok := providerNSOPolicyEntry(doc, name)
	if !ok {
		return gjson.Result{}, fmt.Errorf("StatusCode 404: the policy has no entry for the NSO provider %q", name)
	}
	body := entry.Raw
	for _, flag := range providerNSOPolicyFlags {
		if !entry.Get(flag).Exists() {
			body, _ = sjson.Set(body, flag, false)
		}
	}
	body, _ = sjson.Set(body, "provider_name", name)
	return gjson.Parse(body), nil
}

func providerNSOPolicyDocumentWith(doc gjson.Result, name, entry string) (string, error) {
	var root map[string]json.RawMessage
	if err := json.Unmarshal([]byte(doc.Raw), &root); err != nil || root == nil {
		return "", fmt.Errorf("the policy of the controller is not a JSON object: %s", doc.Raw)
	}
	entries := map[string]json.RawMessage{}
	if raw, ok := root["provider_policy"]; ok {
		if err := json.Unmarshal(raw, &entries); err != nil || entries == nil {
			entries = map[string]json.RawMessage{}
		}
	}
	if entry == "" {
		entry = "{}"
	}
	entries[name] = json.RawMessage(entry)
	raw, err := json.Marshal(entries)
	if err != nil {
		return "", err
	}
	root["provider_policy"] = raw
	out, err := json.Marshal(root)
	return string(out), err
}

func providerNSOPolicyJobFailure(job gjson.Result) string {
	if job.Get("state").String() != "JOB_FAILED" {
		return ""
	}
	message := job.Get("error").String()
	if errors := job.Get("errors"); errors.Exists() && errors.Raw != "[]" && errors.Raw != "null" {
		message = strings.TrimSpace(message + " " + errors.Raw)
	}
	if message == "" {
		message = job.Raw
	}
	return message
}

// A PUT of the policy is normally finished when it answers (state JOB_COMPLETED), so this is only a fallback.
func waitProviderNSOPolicyJob(ctx context.Context, c *cw.Client, jobId string) error {
	deadline := time.Now().Add(providerNSOPolicyJobTimeout)
	body, _ := sjson.Set("{}", "filter.job_id", jobId)
	for {
		res, err := c.Post(nsoJobsQueryPath, body, cw.ReadOnly)
		if err == nil {
			job := res.Get("data.0")
			if !job.Exists() {
				job = res.Get("jobs.0")
			}
			if failure := providerNSOPolicyJobFailure(job); failure != "" {
				return fmt.Errorf("the controller did not apply the policy: %s", failure)
			}
			if job.Get("state").String() == "JOB_COMPLETED" {
				return nil
			}
		}
		if !time.Now().Add(settingsRetryInterval).Before(deadline) || !sleepContext(ctx, settingsRetryInterval) {
			return fmt.Errorf("the job %s of the policy did not finish in time", jobId)
		}
	}
}

func applyProviderNSOPolicy(ctx context.Context, c *cw.Client, plan *ProviderNSOPolicies) error {
	name := plan.ProviderName.ValueString()
	doc, err := c.Post(providerNSOPolicyQueryPath, `{}`, cw.ReadOnly)
	if err != nil {
		return fmt.Errorf("%s, %s", err, doc.String())
	}
	if _, ok := providerNSOPolicyEntry(doc, name); !ok {
		return fmt.Errorf("the controller has no policy for the NSO provider %q: the controller creates it with the provider, so create the provider first", name)
	}
	body, err := providerNSOPolicyDocumentWith(doc, name, plan.toBody(ctx, ProviderNSOPolicies{}))
	if err != nil {
		return err
	}
	res, err := c.Put(providerNSOPolicyPath, body, cw.UseMutex)
	if err != nil {
		return fmt.Errorf("%s, %s", err, res.String())
	}
	// The controller answers 200 for a refusal and names it in the state of the job; nothing is applied then.
	if failure := providerNSOPolicyJobFailure(res); failure != "" {
		return fmt.Errorf("the controller refused the policy: %s", failure)
	}
	if state := res.Get("state").String(); state != "JOB_COMPLETED" && res.Get("job_id").Exists() {
		tflog.Debug(ctx, fmt.Sprintf("provider_nso_policies: job %s is %q, waiting for it", res.Get("job_id").String(), state))
		if err := waitProviderNSOPolicyJob(ctx, c, res.Get("job_id").String()); err != nil {
			return err
		}
	}
	read, err := plan.readGet(c, providerNSOPolicyQueryPath)
	if err != nil {
		return err
	}
	plan.updateFromBody(ctx, read)
	return nil
}

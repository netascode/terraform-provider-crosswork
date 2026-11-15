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
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/tidwall/gjson"
)

const providerNSOPolicyAnswer = `{"name":"default","providers_criteria":"*","provider_policy":{"nso":{"match":true,"matchRule":"*","onboardTo":true,"onboardToRule":"*","onboardFromRule":"*","syncFrom":true,"syncFromRule":"*","checkSync":true,"checkSyncRule":"*"},"cnc.one":{"onboardFrom":true,"neds":[{"ned":"a","rule":"b"}]}},"policy":{"auto_onboard_rfs":true,"rfs_spread_method":"ROUND_ROBIN","rfs_spread_value":150}}`

func TestFlattenProviderNSOPolicyEntry(t *testing.T) {
	flat, err := flattenProviderNSOPolicyEntry(gjson.Parse(providerNSOPolicyAnswer), "nso")
	assert.True(t, err == nil, err)
	assert.Equal(t, "nso", flat.Get("provider_name").String())
	assert.True(t, flat.Get("match").Bool())
	for _, flag := range []string{"onboardFrom"} {
		assert.True(t, flat.Get(flag).Exists(), "a flag the controller leaves out is read as false")
		assert.False(t, flat.Get(flag).Bool())
	}
	assert.False(t, flat.Get("neds").Exists(), "no NED rows stay absent")

	// a provider name with dots is one key
	flat, err = flattenProviderNSOPolicyEntry(gjson.Parse(providerNSOPolicyAnswer), "cnc.one")
	assert.True(t, err == nil, err)
	assert.True(t, flat.Get("onboardFrom").Bool())
	assert.False(t, flat.Get("syncFrom").Bool())
	assert.Equal(t, "a", flat.Get("neds.0.ned").String())
}

func TestFlattenProviderNSOPolicyEntryMissing(t *testing.T) {
	_, err := flattenProviderNSOPolicyEntry(gjson.Parse(providerNSOPolicyAnswer), "gone")
	assert.True(t, err != nil)
	assert.Contains(t, err.Error(), "StatusCode 404")
	_, err = flattenProviderNSOPolicyEntry(gjson.Parse(`{"name":"default"}`), "nso")
	assert.True(t, err != nil)
}

func TestProviderNSOPolicyDocumentKeepsTheRest(t *testing.T) {
	doc := gjson.Parse(providerNSOPolicyAnswer)
	plan := ProviderNSOPolicies{ProviderName: types.StringValue("nso"), Match: types.BoolValue(false), MatchRule: types.StringValue("host_name = 'x*'"),
		Neds: []ProviderNSOPoliciesNeds{{Ned: types.StringValue("n1"), Rule: types.StringValue("r1")}, {Ned: types.StringValue("n0"), Rule: types.StringValue("r0")}}}
	out, err := providerNSOPolicyDocumentWith(doc, "nso", plan.toBody(context.Background(), ProviderNSOPolicies{}))
	assert.True(t, err == nil, err)
	got := gjson.Parse(out)
	// the entry is the one of the plan, as a whole: a key the plan leaves out is gone
	assert.False(t, got.Get("provider_policy.nso.match").Bool())
	assert.Equal(t, "host_name = 'x*'", got.Get("provider_policy.nso.matchRule").String())
	assert.False(t, got.Get("provider_policy.nso.syncFrom").Exists())
	assert.Equal(t, "n1", got.Get("provider_policy.nso.neds.0.ned").String(), "the order of the NED rows is kept")
	assert.Equal(t, "n0", got.Get("provider_policy.nso.neds.1.ned").String())
	assert.Equal(t, "default", got.Get("name").String())
	assert.Equal(t, "*", got.Get("providers_criteria").String())
	assert.Equal(t, doc.Get("policy").Raw, got.Get("policy").Raw)
	other, _ := providerNSOPolicyEntry(got, "cnc.one")
	original, _ := providerNSOPolicyEntry(doc, "cnc.one")
	assert.Equal(t, original.Raw, other.Raw)
	assert.Equal(t, int64(150), got.Get("policy.rfs_spread_value").Int())
}

func TestProviderNSOPolicyJobFailure(t *testing.T) {
	assert.Equal(t, "", providerNSOPolicyJobFailure(gjson.Parse(`{"state":"JOB_COMPLETED","job_id":"x"}`)))
	assert.Equal(t, "", providerNSOPolicyJobFailure(gjson.Parse(`{"state":"JOB_RUNNING"}`)))
	failure := providerNSOPolicyJobFailure(gjson.Parse(`{"state":"JOB_FAILED","error":"Cannot find comparator operation","errors":["bad rule"]}`))
	assert.Contains(t, failure, "Cannot find comparator operation")
	assert.Contains(t, failure, "bad rule")
	assert.NotEqual(t, "", providerNSOPolicyJobFailure(gjson.Parse(`{"state":"JOB_FAILED"}`)))
}

// The live tests change the policy of the NSO provider named "nso" and put the entry back.

func providerNSOPolicyTestProvider() string {
	return "nso"
}

func providerNSOPolicyLiveEntry(t *testing.T) (gjson.Result, gjson.Result) {
	t.Helper()
	if os.Getenv("TF_ACC") == "" {
		t.Skip("skipping test, set environment variable TF_ACC")
	}
	insecure, _ := strconv.ParseBool(os.Getenv("CW_INSECURE"))
	c, err := registeredClient(os.Getenv("CW_URL"), os.Getenv("CW_USERNAME"), os.Getenv("CW_PASSWORD"), insecure, 3, "test")
	if err != nil {
		t.Fatalf("login: %s", err)
	}
	doc, err := c.Post(providerNSOPolicyQueryPath, `{}`)
	if err != nil {
		t.Fatalf("policy: %s %s", err, doc.String())
	}
	entry, ok := providerNSOPolicyEntry(doc, providerNSOPolicyTestProvider())
	if !ok {
		t.Skipf("the controller has no policy for the NSO provider %q", providerNSOPolicyTestProvider())
	}
	return doc, entry
}

func providerNSOPolicyRestore(t *testing.T, original gjson.Result) {
	insecure, _ := strconv.ParseBool(os.Getenv("CW_INSECURE"))
	c, err := registeredClient(os.Getenv("CW_URL"), os.Getenv("CW_USERNAME"), os.Getenv("CW_PASSWORD"), insecure, 3, "test")
	if err != nil {
		t.Errorf("restore, login: %s", err)
		return
	}
	doc, err := c.Post(providerNSOPolicyQueryPath, `{}`)
	if err != nil {
		t.Errorf("restore, read: %s", err)
		return
	}
	body, err := providerNSOPolicyDocumentWith(doc, providerNSOPolicyTestProvider(), original.Raw)
	if err != nil {
		t.Errorf("restore: %s", err)
		return
	}
	if res, err := c.Put(providerNSOPolicyPath, body); err != nil {
		t.Errorf("restore: %s %s", err, res.String())
	}
}

func providerNSOPolicyConfig(extra string) string {
	return fmt.Sprintf("resource \"crosswork_provider_nso_policies\" \"test\" {\n  provider_name = %q\n%s}\n", providerNSOPolicyTestProvider(), extra)
}

func providerNSOPolicyHolds(t *testing.T, before gjson.Result) resource.TestCheckFunc {
	return func(*terraform.State) error {
		after, _ := providerNSOPolicyLiveEntry(t)
		for _, path := range []string{"name", "providers_criteria", "policy"} {
			if after.Get(path).Raw != before.Get(path).Raw {
				return fmt.Errorf("%s changed: %s, was %s", path, after.Get(path).Raw, before.Get(path).Raw)
			}
		}
		var a, b map[string]json.RawMessage
		_ = json.Unmarshal([]byte(after.Get("provider_policy").Raw), &a)
		_ = json.Unmarshal([]byte(before.Get("provider_policy").Raw), &b)
		for name, raw := range b {
			if name == providerNSOPolicyTestProvider() {
				continue
			}
			if string(a[name]) != string(raw) {
				return fmt.Errorf("the entry of %q changed: %s, was %s", name, a[name], raw)
			}
		}
		return nil
	}
}

func TestAccCwProviderNSOPoliciesKeepsTheRest(t *testing.T) {
	before, original := providerNSOPolicyLiveEntry(t)
	t.Cleanup(func() { providerNSOPolicyRestore(t, original) })
	const res = "crosswork_provider_nso_policies.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// the values of the configuration are the entry; a key left out is off
			{Config: providerNSOPolicyConfig("  match = true\n  match_rule = \"*\"\n  onboard_to_nso = false\n  onboard_to_nso_rule = \"*\"\n  neds = [{ ned = \"cwtest-ned\", rule = \"host_name = 'cwtest-*'\" }]\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "match", "true"),
					resource.TestCheckResourceAttr(res, "onboard_to_nso", "false"),
					resource.TestCheckResourceAttr(res, "neds.0.ned", "cwtest-ned"),
					resource.TestCheckNoResourceAttr(res, "sync_from_nso"),
					providerNSOPolicyHolds(t, before))},
			{Config: providerNSOPolicyConfig("  match = true\n  match_rule = \"host_name = 'cwtest-*'\"\n  sync_from_nso = true\n  sync_from_nso_rule = \"*\"\n"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "match_rule", "host_name = 'cwtest-*'"),
					resource.TestCheckResourceAttr(res, "sync_from_nso", "true"),
					resource.TestCheckResourceAttr(res, "neds.#", "0"),
					providerNSOPolicyHolds(t, before))},
			// import by provider name; a flag the configuration leaves out is null in the state and false after an import
			{ResourceName: res, ImportState: true, ImportStateId: providerNSOPolicyTestProvider(), ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{"onboard_to_nso", "onboard_from_nso", "check_sync"}},
		},
	})
}

func TestAccCwProviderNSOPoliciesRefused(t *testing.T) {
	_, original := providerNSOPolicyLiveEntry(t)
	t.Cleanup(func() { providerNSOPolicyRestore(t, original) })
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: providerNSOPolicyConfig("  match = true\n  match_rule = \"cwtest-*\"\n"), ExpectError: regexp.MustCompile("Cannot find comparator")},
		},
	})
}

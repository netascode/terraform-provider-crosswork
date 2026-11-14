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

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type ssoStub struct {
	*httptest.Server
	logins, logouts atomic.Int32
	status          int // when not zero, the ticket call answers this status
}

func newSSOStub(t *testing.T) *ssoStub {
	s := &ssoStub{}
	s.Server = httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/crosswork/sso/v1/tickets":
			if s.status != 0 {
				w.WriteHeader(s.status)
				fmt.Fprint(w, `{"error":"Per user session limit reached. Close unused sessions or try after sometime."}`)
				return
			}
			s.logins.Add(1)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprint(w, "TGT-1-abc-cas")
		case r.Method == http.MethodPost && r.URL.Path == "/crosswork/sso/v2/tickets/jwt":
			fmt.Fprint(w, "aaa.bbb.ccc")
		case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/crosswork/sso/v1/tickets/"):
			s.logouts.Add(1)
		case r.Method == http.MethodGet && r.URL.Path == "/crosswork/test/cluster":
			fmt.Fprintf(w, `{"cluster":%q}`, s.URL)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(s.Close)
	return s
}

func configure(t *testing.T, url, username, password string) *provider.ConfigureResponse {
	t.Helper()
	return configureWith(t, url, username, password, nil)
}

func configureWith(t *testing.T, url, username, password string, extra map[string]any) *provider.ConfigureResponse {
	t.Helper()
	ctx := context.Background()
	p := New("test")()

	var schemaResp provider.SchemaResponse
	p.Schema(ctx, provider.SchemaRequest{}, &schemaResp)

	objType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	values := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		values[name] = tftypes.NewValue(typ, nil)
	}
	set := func(name, v string) {
		if v != "" {
			values[name] = tftypes.NewValue(tftypes.String, v)
		}
	}
	set("url", url)
	set("username", username)
	set("password", password)
	for name, v := range extra {
		switch x := v.(type) {
		case string:
			values[name] = tftypes.NewValue(tftypes.String, x)
		case int64:
			values[name] = tftypes.NewValue(tftypes.Number, x)
		}
	}

	resp := &provider.ConfigureResponse{}
	p.Configure(ctx, provider.ConfigureRequest{
		Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, values)},
	}, resp)
	return resp
}

func clearEnv(t *testing.T) {
	for _, v := range []string{"CW_URL", "CW_USERNAME", "CW_PASSWORD", "CW_INSECURE", "CW_RETRIES", "CW_MAX_TIMEOUT", "CW_ALLOW_EXISTING_ON_CREATE", "CW_DEVICE_READY_MODE", "CW_DEVICE_READY_TIMEOUT"} {
		t.Setenv(v, "")
	}
}

func TestConfigureOneLoginPerCluster(t *testing.T) {
	clearEnv(t)
	Shutdown()
	stub := newSSOStub(t)

	for i := 0; i < 3; i++ {
		resp := configure(t, stub.URL, "admin", "pw")
		if resp.Diagnostics.HasError() {
			t.Fatalf("configure %d: %v", i, resp.Diagnostics)
		}
	}
	if got := stub.logins.Load(); got != 1 {
		t.Errorf("three Configure calls made %d logins, want 1", got)
	}

	Shutdown()
	if got := stub.logouts.Load(); got != 1 {
		t.Errorf("Shutdown made %d logouts, want 1", got)
	}
	Shutdown()
	if got := stub.logouts.Load(); got != 1 {
		t.Errorf("a second Shutdown made %d logouts in total, want 1", got)
	}
}

func TestConfigureSessionLimit(t *testing.T) {
	clearEnv(t)
	t.Setenv("CW_RETRIES", "0") // the client retries a 503 login with backoff
	Shutdown()
	stub := newSSOStub(t)
	stub.status = http.StatusServiceUnavailable

	resp := configure(t, stub.URL, "admin", "pw")
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error when the login hits the session limit")
	}
	if detail := resp.Diagnostics.Errors()[0].Detail(); !strings.Contains(detail, "NumParallelSessionsPerUser") {
		t.Errorf("diagnostic gives no advice: %s", detail)
	}
	Shutdown()
}

func TestConfigureBadURL(t *testing.T) {
	clearEnv(t)
	resp := configure(t, "10.1.1.1", "admin", "pw")
	if !resp.Diagnostics.HasError() {
		t.Fatal("expected an error for a url without a scheme")
	}
}

func TestConfigureMissingSettings(t *testing.T) {
	clearEnv(t)
	for name, args := range map[string][3]string{
		"url":      {"", "admin", "pw"},
		"username": {"https://10.1.1.1", "", "pw"},
		"password": {"https://10.1.1.1", "admin", ""},
	} {
		if !configure(t, args[0], args[1], args[2]).Diagnostics.HasError() {
			t.Errorf("expected an error when %s is missing", name)
		}
	}
}

func configureInstances(t *testing.T, url string, instances map[string]string, order []string) *provider.ConfigureResponse {
	t.Helper()
	ctx := context.Background()
	p := New("test")()

	var schemaResp provider.SchemaResponse
	p.Schema(ctx, provider.SchemaRequest{}, &schemaResp)

	objType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
	values := map[string]tftypes.Value{}
	for name, typ := range objType.AttributeTypes {
		values[name] = tftypes.NewValue(typ, nil)
	}
	values["username"] = tftypes.NewValue(tftypes.String, "admin")
	values["password"] = tftypes.NewValue(tftypes.String, "pw")
	if url != "" {
		values["url"] = tftypes.NewValue(tftypes.String, url)
	}
	if len(order) > 0 {
		listType := objType.AttributeTypes["instances"].(tftypes.List)
		var items []tftypes.Value
		for _, name := range order {
			items = append(items, tftypes.NewValue(listType.ElementType, map[string]tftypes.Value{
				"name": tftypes.NewValue(tftypes.String, name),
				"url":  tftypes.NewValue(tftypes.String, instances[name]),
			}))
		}
		values["instances"] = tftypes.NewValue(listType, items)
	}

	resp := &provider.ConfigureResponse{}
	p.Configure(ctx, provider.ConfigureRequest{
		Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objType, values)},
	}, resp)
	return resp
}

func clusterOf(t *testing.T, resp *provider.ConfigureResponse, instance string) (string, diag.Diagnostics) {
	t.Helper()
	data := resp.ResourceData.(*CwProviderData)
	var diags diag.Diagnostics
	c, ok := clientFor(context.Background(), data.Clients, instance, &diags)
	if !ok {
		return "", diags
	}
	res, err := c.Get("/crosswork/test/cluster")
	if err != nil {
		t.Fatalf("instance %q: %v", instance, err)
	}
	return res.Get("cluster").String(), diags
}

func TestInstancesSelectCluster(t *testing.T) {
	clearEnv(t)
	Shutdown()
	a, b := newSSOStub(t), newSSOStub(t)

	resp := configureInstances(t, a.URL, map[string]string{"east": a.URL, "west": b.URL}, []string{"east", "west"})
	if resp.Diagnostics.HasError() {
		t.Fatalf("configure: %v", resp.Diagnostics)
	}

	for instance, want := range map[string]string{"": a.URL, "east": a.URL, "west": b.URL} {
		got, diags := clusterOf(t, resp, instance)
		if diags.HasError() || got != want {
			t.Errorf("instance %q answered %q (%v), want %q", instance, got, diags, want)
		}
	}

	// the default instance and "east" share one URL, so they share one session
	if got := a.logins.Load(); got != 1 {
		t.Errorf("cluster a saw %d logins, want 1", got)
	}
	if got := b.logins.Load(); got != 1 {
		t.Errorf("cluster b saw %d logins, want 1", got)
	}

	Shutdown()
	if a.logouts.Load() != 1 || b.logouts.Load() != 1 {
		t.Errorf("Shutdown logouts: a=%d b=%d, want 1 each", a.logouts.Load(), b.logouts.Load())
	}
}

func TestInstancesDefaultFromFirstInstance(t *testing.T) {
	clearEnv(t)
	Shutdown()
	a, b := newSSOStub(t), newSSOStub(t)

	resp := configureInstances(t, "", map[string]string{"east": a.URL, "west": b.URL}, []string{"east", "west"})
	if resp.Diagnostics.HasError() {
		t.Fatalf("configure: %v", resp.Diagnostics)
	}
	if got, _ := clusterOf(t, resp, ""); got != a.URL {
		t.Errorf("default instance answered %q, want the first instance %q", got, a.URL)
	}
	Shutdown()
}

func TestInstancesUnknownInstance(t *testing.T) {
	clearEnv(t)
	Shutdown()
	a := newSSOStub(t)

	resp := configureInstances(t, a.URL, nil, nil)
	if resp.Diagnostics.HasError() {
		t.Fatalf("configure: %v", resp.Diagnostics)
	}
	_, diags := clusterOf(t, resp, "nope")
	if !diags.HasError() {
		t.Fatal("expected an error for an unknown instance")
	}
	e := diags.Errors()[0]
	if e.Summary() != "Instance does not exist in provider configuration" || !strings.Contains(e.Detail(), `"nope"`) {
		t.Errorf("unclear diagnostic: %s: %s", e.Summary(), e.Detail())
	}
	Shutdown()
}

func TestInstancesInvalidList(t *testing.T) {
	clearEnv(t)
	Shutdown()
	a := newSSOStub(t)
	for name, c := range map[string]struct {
		urls  map[string]string
		order []string
	}{
		"duplicate name": {map[string]string{"x": a.URL}, []string{"x", "x"}},
		"empty name":     {map[string]string{"": a.URL}, []string{""}},
		"url no scheme":  {map[string]string{"x": "10.1.1.1"}, []string{"x"}},
	} {
		if !configureInstances(t, a.URL, c.urls, c.order).Diagnostics.HasError() {
			t.Errorf("%s: expected an error", name)
		}
	}
	Shutdown()
}

func registeredWaits() map[string]int {
	registry.Lock()
	defer registry.Unlock()
	out := map[string]int{}
	for k, c := range registry.clients {
		out[k.url] = c.DefaultMaxAsyncWaitTime
	}
	return out
}

func TestConfigureMaxTimeout(t *testing.T) {
	for name, tc := range map[string]struct {
		env  string
		want int
	}{"default": {"", 180}, "environment": {"45", 45}} {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			Shutdown()
			stub := newSSOStub(t)
			t.Setenv("CW_MAX_TIMEOUT", tc.env)
			if resp := configure(t, stub.URL, "admin", "pw"); resp.Diagnostics.HasError() {
				t.Fatalf("configure failed: %v", resp.Diagnostics)
			}
			if got := registeredWaits()[stub.URL]; got != tc.want {
				t.Errorf("the client waits %d s for a plan, want %d", got, tc.want)
			}
		})
	}
}

func TestConfigureDeviceReady(t *testing.T) {
	for name, tc := range map[string]struct {
		env     map[string]string
		extra   map[string]any
		mode    string
		timeout time.Duration
		errText string
	}{
		"defaults":         {nil, nil, nsoModeSynced, 300 * time.Second, ""},
		"environment":      {map[string]string{"CW_DEVICE_READY_MODE": "full", "CW_DEVICE_READY_TIMEOUT": "45"}, nil, nsoModeFull, 45 * time.Second, ""},
		"provider block":   {map[string]string{"CW_DEVICE_READY_MODE": "full", "CW_DEVICE_READY_TIMEOUT": "45"}, map[string]any{"device_ready_mode": "none", "device_ready_timeout": int64(0)}, nsoModeNone, 0, ""},
		"mode none by env": {map[string]string{"CW_DEVICE_READY_MODE": "none"}, nil, nsoModeNone, 300 * time.Second, ""},
		"bad mode":         {map[string]string{"CW_DEVICE_READY_MODE": "fast"}, nil, "", 0, "must be none, synced or full"},
		"bad timeout":      {map[string]string{"CW_DEVICE_READY_TIMEOUT": "soon"}, nil, "", 0, "number of seconds from 0 to 86400"},
		"timeout too long": {map[string]string{"CW_DEVICE_READY_TIMEOUT": "86401"}, nil, "", 0, "number of seconds from 0 to 86400"},
	} {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			Shutdown()
			stub := newSSOStub(t)
			for k, v := range tc.env {
				t.Setenv(k, v)
			}
			resp := configureWith(t, stub.URL, "admin", "pw", tc.extra)
			if tc.errText != "" {
				if !resp.Diagnostics.HasError() || !strings.Contains(resp.Diagnostics.Errors()[0].Detail(), tc.errText) {
					t.Fatalf("want an error with %q, got %v", tc.errText, resp.Diagnostics)
				}
				return
			}
			if resp.Diagnostics.HasError() {
				t.Fatalf("configure failed: %v", resp.Diagnostics)
			}
			got := resp.ResourceData.(*CwProviderData).DeviceReady
			if got.Mode != tc.mode || got.Timeout != tc.timeout || got.Store == nil {
				t.Errorf("device ready is mode %q, timeout %v, store %v; want %q, %v and a store", got.Mode, got.Timeout, got.Store, tc.mode, tc.timeout)
			}
		})
	}
}

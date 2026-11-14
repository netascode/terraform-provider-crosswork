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

// Crosswork limits the sessions per user, so the provider keeps one client per cluster and process and logs every client out
// when the process ends.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	cw "github.com/netascode/go-crosswork"
)

// Terraform stops waiting for a provider process shortly after asking it to stop.
const shutdownTimeout = 1500 * time.Millisecond

type clientKey struct {
	url          string
	username     string
	passwordHash string
	insecure     bool
	retries      int
	maxTimeout   int
}

// Zero keeps the library default. Only the tests change it.
var planPollInterval time.Duration

var registry = struct {
	sync.Mutex
	clients map[clientKey]*cw.Client
}{clients: map[clientKey]*cw.Client{}}

// Terraform calls Configure once per run and the acceptance tests call it for every step; the registry avoids a login per step.
func registeredClient(url, username, password string, insecure bool, retries int, version string) (*cw.Client, error) {
	return registeredClientWithTimeout(url, username, password, insecure, retries, cw.DefaultDefaultMaxAsyncWaitTime, version)
}

// maxTimeout is in seconds (the provider attribute max_timeout).
func registeredClientWithTimeout(url, username, password string, insecure bool, retries, maxTimeout int, version string) (*cw.Client, error) {
	sum := sha256.Sum256([]byte(password))
	key := clientKey{url: url, username: username, passwordHash: hex.EncodeToString(sum[:]), insecure: insecure, retries: retries, maxTimeout: maxTimeout}

	registry.Lock()
	defer registry.Unlock()
	if c, ok := registry.clients[key]; ok {
		return c, nil
	}
	opts := []func(*cw.Client){cw.Insecure(insecure), cw.MaxRetries(retries), cw.DefaultMaxAsyncWaitTime(maxTimeout), cw.UserAgent(fmt.Sprintf("CrossworkTerraform/%s Cisco", version))}
	if planPollInterval > 0 {
		opts = append(opts, cw.PlanPollInterval(planPollInterval))
	}
	c, err := cw.NewClient(url, username, password, opts...)
	if err != nil {
		return nil, err
	}
	registry.clients[key] = &c
	return &c, nil
}

// Shutdown logs every registered client out, in parallel and within shutdownTimeout.
// Errors are ignored: a session that cannot be released ends with the idle timeout.
func Shutdown() {
	registry.Lock()
	clients := make([]*cw.Client, 0, len(registry.clients))
	for _, c := range registry.clients {
		clients = append(clients, c)
	}
	registry.clients = map[clientKey]*cw.Client{}
	registry.Unlock()

	// A fresh context: a request context may already be cancelled.
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	var wg sync.WaitGroup
	for _, c := range clients {
		wg.Add(1)
		go func(c *cw.Client) {
			defer wg.Done()
			bound := c.WithContext(ctx)
			_ = bound.Logout()
		}(c)
	}
	wg.Wait()
}

func loginErrorDetail(err error) string {
	if cw.IsSessionLimit(err) {
		return "Crosswork refused the login because the user has reached the per-user session limit.\n\n" +
			"Terminate unused sessions (Crosswork UI, Active Sessions), raise NumParallelSessionsPerUser " +
			"(GET/PUT /crosswork/aaa/v1/sessionconfig), or use a dedicated automation user.\n\n" + err.Error()
	}
	return err.Error()
}

// An empty name selects the default client; an unknown name adds an error to diags and returns false.
func clientFor(ctx context.Context, clients map[string]*cw.Client, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	c, ok := clients[instance]
	if !ok || c == nil {
		if instance == "" {
			diags.AddError("Client Error", "The provider has no default instance")
		} else {
			diags.AddAttributeError(path.Root("instance"), "Instance does not exist in provider configuration",
				fmt.Sprintf("Instance %q is not listed in the provider's instances attribute.", instance))
		}
		return cw.Client{}, false
	}
	return c.WithContext(ctx), true
}

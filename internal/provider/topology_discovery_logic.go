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
	"fmt"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	cw "github.com/netascode/go-crosswork"
)

const (
	discoveryPath         = "/crosswork/topology/v1/topology-service/topology/provider/settings/protocols"
	discoveryDefaultsPath = "/crosswork/topology/v1/topology-service/topology/provider/settings/protocols/defaults"
)

// The topology applies a change of the discovery protocols in the background: a protocol has the status EnableInProgress or
// DisableInProgress until it is Ready, a write while a change runs is refused, and a change that the controller did apply can be
// answered with the error PROVIDER_NOT_AVAILABLE. A write therefore waits for Ready, sends the values, and judges by reading them back.

var (
	discoveryWait = 180 * time.Second
	discoveryPoll = 3 * time.Second
)

var discoveryProtocols = []string{"lldpConfig", "cdpConfig", "lagConfig"}

func discoveryReady(res gjson.Result) bool {
	for _, protocol := range discoveryProtocols {
		if res.Get("providerConfig."+protocol+".status").String() != "Ready" {
			return false
		}
	}
	return true
}

func discoveryMatches(res gjson.Result, want map[string]bool) bool {
	for protocol, enabled := range want {
		if res.Get("providerConfig."+protocol+".enabled").Bool() != enabled {
			return false
		}
	}
	return true
}

func waitDiscovery(ctx context.Context, c *cw.Client, want map[string]bool) (gjson.Result, error) {
	deadline := time.Now().Add(discoveryWait)
	for {
		res, err := c.Get(discoveryPath)
		if err == nil && discoveryReady(res) && discoveryMatches(res, want) {
			return res, nil
		}
		if !time.Now().Add(discoveryPoll).Before(deadline) || !sleepContext(ctx, discoveryPoll) {
			if err != nil {
				return res, err
			}
			return res, fmt.Errorf("the topology did not finish the change of the discovery protocols in time: %s", res.Raw)
		}
	}
}

// The library turns the answer PROVIDER_NOT_AVAILABLE, which comes with a change that the controller did apply, into an
// error, so the state that is read back decides.
func discoveryWriteRefused(err error) bool {
	return err != nil && !strings.Contains(err.Error(), "PROVIDER_NOT_AVAILABLE")
}

func discoveryRefusal(res gjson.Result) string {
	if message := res.Get("error").String(); message != "" && message != "PROVIDER_NOT_AVAILABLE" {
		return message
	}
	return ""
}

func discoveryWanted(plan TopologyDiscoverySettings) map[string]bool {
	return map[string]bool{"lldpConfig": plan.Lldp.ValueBool(), "cdpConfig": plan.Cdp.ValueBool(), "lagConfig": plan.Lag.ValueBool()}
}

func applyDiscovery(ctx context.Context, c *cw.Client, plan TopologyDiscoverySettings) error {
	if _, err := waitDiscovery(ctx, c, nil); err != nil {
		return err
	}
	res, err := c.Post(discoveryPath, plan.toBody(ctx, TopologyDiscoverySettings{}), cw.UseMutex)
	if discoveryWriteRefused(err) {
		return fmt.Errorf("%s, %s", err, res.String())
	}
	if message := discoveryRefusal(res); message != "" {
		return fmt.Errorf("the controller refused the change: %s", message)
	}
	if _, err := waitDiscovery(ctx, c, discoveryWanted(plan)); err != nil {
		return fmt.Errorf("%s (the answer of the write was %q)", err, res.Raw)
	}
	return nil
}

func resetDiscovery(ctx context.Context, c *cw.Client) error {
	if _, err := waitDiscovery(ctx, c, nil); err != nil {
		return err
	}
	res, err := c.Post(discoveryDefaultsPath, `{}`, cw.UseMutex)
	if discoveryWriteRefused(err) {
		return fmt.Errorf("%s, %s", err, res.String())
	}
	if message := discoveryRefusal(res); message != "" {
		return fmt.Errorf("the controller refused the reset: %s", message)
	}
	if _, err := waitDiscovery(ctx, c, map[string]bool{"lldpConfig": false, "cdpConfig": false, "lagConfig": false}); err != nil {
		return fmt.Errorf("%s (the answer of the reset was %q)", err, res.Raw)
	}
	return nil
}

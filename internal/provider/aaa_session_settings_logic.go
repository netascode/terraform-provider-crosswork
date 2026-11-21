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
	"fmt"

	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// The source IP setting of the AAA settings page is not part of the session configuration (the controller
// refuses an unknown member in that body); it is read and written through the cluster manager.
const (
	sourceIPGetPath = "/crosswork/platform/v2/cluster/ip-traceability/get"
	sourceIPSetPath = "/crosswork/platform/v2/cluster/ip-traceability/set"
)

// sessionPutBody removes the source IP member, which the model keeps in the session body, from the body of the session call.
func sessionPutBody(body string) string {
	body, _ = sjson.Delete(body, "IsEnabled")
	return body
}

// sourceIPSetBody is the body of the call that sets the source IP setting.
func sourceIPSetBody(enabled bool) string {
	return fmt.Sprintf(`{"IsEnabled":%t,"ScaleUp":false}`, enabled)
}

// sourceIPFromAnswer reads the setting from the answer of the cluster manager, which names it Enabled.
func sourceIPFromAnswer(res gjson.Result) (bool, error) {
	if msg := res.Get("Error").String(); msg != "" {
		return false, fmt.Errorf("the controller reports an error for the source IP setting: %s", msg)
	}
	if !res.Get("Enabled").Exists() {
		return false, fmt.Errorf("the answer for the source IP setting has no Enabled member: %s", res.Raw)
	}
	return res.Get("Enabled").Bool(), nil
}

// sourceIPSetResult checks the answer of the call that sets the source IP setting.
func sourceIPSetResult(res gjson.Result) error {
	if v := res.Get("resp_value").String(); v != "R_SUCCESS" {
		return fmt.Errorf("the controller did not accept the source IP setting: %s %s", v, res.Get("resp_error").String())
	}
	return nil
}

// withSourceIP returns the session answer with the source IP setting added as the member the model reads (IsEnabled).
func withSourceIP(res gjson.Result, enabled bool) gjson.Result {
	raw, err := sjson.Set(res.Raw, "IsEnabled", enabled)
	if err != nil {
		return res
	}
	return gjson.Parse(raw)
}

// readSourceIP reads the source IP setting and adds it to the session answer.
func readSourceIP(c cw.Client, session gjson.Result) (gjson.Result, error) {
	res, err := c.Get(sourceIPGetPath)
	if err != nil {
		return session, err
	}
	enabled, err := sourceIPFromAnswer(res)
	if err != nil {
		return session, err
	}
	return withSourceIP(session, enabled), nil
}

// writeSourceIP sets the source IP setting.
func writeSourceIP(c cw.Client, enabled bool) error {
	res, err := c.Post(sourceIPSetPath, sourceIPSetBody(enabled), cw.UseMutex)
	if err != nil {
		return err
	}
	return sourceIPSetResult(res)
}

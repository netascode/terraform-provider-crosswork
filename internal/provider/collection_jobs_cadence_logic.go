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

// Written by hand: the setting is one value for the whole cluster with no id, the body must hold all four values, a refused
// value is a job failing inside HTTP 200, and there is no delete.

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	cadencePath      = "/crosswork/inventory/v1/devicepackage/cadence"
	cadenceQueryPath = cadencePath + "/query"

	cadenceId = "collection_jobs_cadence"
)

// Bound the wait for the controller to show an accepted value; variables so that tests can shorten them.
const cadenceShownAttempts = 20

var cadenceShownInterval = 250 * time.Millisecond

var cadenceKeys = []struct{ attribute, key string }{
	{"reachability_check", "reach-check"},
	{"snmp_check", "snmp"},
	{"te_tunnel_id", "te-tunnel-id"},
	{"show_clock", "show-clock"},
}

type cadenceValues struct {
	ReachabilityCheck types.Int64
	SnmpCheck         types.Int64
	TeTunnelId        types.Int64
	ShowClock         types.Int64
}

func (v *cadenceValues) fields() [4]*types.Int64 {
	return [4]*types.Int64{&v.ReachabilityCheck, &v.SnmpCheck, &v.TeTunnelId, &v.ShowClock}
}

func cadenceFromModel(m CollectionJobsCadence) cadenceValues {
	return cadenceValues{m.ReachabilityCheck, m.SnmpCheck, m.TeTunnelId, m.ShowClock}
}

func readCadence(c *cw.Client) (cadenceValues, error) {
	res, err := c.Post(cadenceQueryPath, `{}`, cw.ReadOnly)
	if err != nil {
		return cadenceValues{}, fmt.Errorf("%w, %s", err, res.String())
	}
	return cadenceFromBody(res), nil
}

func cadenceFromBody(res gjson.Result) cadenceValues {
	var v cadenceValues
	for i, f := range v.fields() {
		if n := res.Get("JobToCadence." + cadenceKeys[i].key); n.Exists() {
			*f = types.Int64Value(n.Int())
		} else {
			*f = types.Int64Null()
		}
	}
	return v
}

// The controller refuses a body that does not name all four jobs, so unset values are filled from the cluster.
func cadenceBody(plan, current cadenceValues) (string, error) {
	body := `{}`
	want, have := plan.fields(), current.fields()
	for i, f := range want {
		v := *f
		if v.IsNull() || v.IsUnknown() {
			v = *have[i]
		}
		if v.IsNull() {
			return "", fmt.Errorf("the controller does not report %s, so a value for %s is needed", cadenceKeys[i].key, cadenceKeys[i].attribute)
		}
		var err error
		if body, err = sjson.Set(body, "JobToCadence."+cadenceKeys[i].key, v.ValueInt64()); err != nil {
			return "", err
		}
	}
	return body, nil
}

func putCadence(c *cw.Client, body string) error {
	if res, err := c.Put(cadencePath, body); err != nil {
		return fmt.Errorf("%w, %s", err, res.String())
	}
	sent := cadenceFromBody(gjson.Parse(body))
	var last cadenceValues
	for i := 0; i < cadenceShownAttempts; i++ {
		got, err := readCadence(c)
		if err != nil {
			return err
		}
		last = got
		if got == sent {
			return nil
		}
		time.Sleep(cadenceShownInterval)
	}
	return fmt.Errorf("Crosswork accepted the cadence but still shows %s, not the %s that was sent", cadenceBodyText(last), cadenceBodyText(sent))
}

func cadenceBodyText(v cadenceValues) string {
	text := ""
	for i, f := range v.fields() {
		if i > 0 {
			text += ", "
		}
		if f.IsNull() {
			text += cadenceKeys[i].key + "=unset"
		} else {
			text += fmt.Sprintf("%s=%d", cadenceKeys[i].key, f.ValueInt64())
		}
	}
	return text
}

// Keeps only the values the state names, so a value the configuration leaves out shows no difference; a state that
// names none is an import and takes all four.
func cadenceKeep(state, current cadenceValues) cadenceValues {
	named := false
	for _, f := range state.fields() {
		if !f.IsNull() {
			named = true
		}
	}
	if !named {
		return current
	}
	out := state
	have, set := current.fields(), out.fields()
	for i, f := range set {
		if !f.IsNull() {
			*f = *have[i]
		}
	}
	return out
}

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
	"math/rand"
	"reflect"
	"testing"

	"github.com/tidwall/gjson"
)

const uiDelivery = " (or (notificationDeliveryMechanismNE 8 9))"

func TestDeviceCriteriaGolden(t *testing.T) {
	cases := []struct {
		name    string
		filters []deviceFilter
		group   string
		want    string
	}{
		{"severity only", []deviceFilter{{Severities: []string{"critical"}}}, "",
			"(and (or (and (severity 'CRITICAL'))) (or (notificationDeliveryMechanismNE 8 9)))"},
		{"category with severities, high to low", []deviceFilter{{Category: "BGP", Severities: []string{"warning", "major", "minor"}}}, "",
			"(and (or (and (category 'BGP') (severity 'MAJOR' 'MINOR' 'WARNING'))) (or (notificationDeliveryMechanismNE 8 9)))"},
		{"events sorted with the space of the UI", []deviceFilter{{Events: []string{"b", "a"}}}, "",
			"(and (or (and ( eventType 'a' 'b'))) (or (notificationDeliveryMechanismNE 8 9)))"},
		{"group has no space before it", []deviceFilter{{Severities: []string{"information"}}}, "g-1",
			"(and (or (and (severity 'INFORMATION'))) (or (notificationDeliveryMechanismNE 8 9))(deviceGroup 'g-1'))"},
		{"filters sorted by their text", []deviceFilter{{Severities: []string{"critical"}}, {Category: "BGP", Severities: []string{"major"}}}, "",
			"(and (or (and (category 'BGP') (severity 'MAJOR')) (and (severity 'CRITICAL'))) (or (notificationDeliveryMechanismNE 8 9)))"},
	}
	for _, c := range cases {
		if got := deviceCriteria(c.filters, c.group); got != c.want {
			t.Errorf("%s:\n got %s\nwant %s", c.name, got, c.want)
		}
	}
}

func TestParseDeviceCriteriaUI(t *testing.T) {
	// The strings the UI stored (probe log 27, 32): severity-only filters first, a space after the opening parenthesis of events.
	text := "(and (or (and (severity 'CRITICAL')) (and (category 'BGP') (severity 'MAJOR' 'MINOR' 'WARNING')) (and ( eventType 'x' 'y'))) (or (notificationDeliveryMechanismNE 8 9))(deviceGroup 'abc'))"
	filters, group, ok := parseDeviceCriteria(text)
	if !ok || group != "abc" {
		t.Fatalf("ok %v group %q", ok, group)
	}
	want := []deviceFilter{
		{Severities: []string{"critical"}},
		{Category: "BGP", Severities: []string{"major", "minor", "warning"}},
		{Events: []string{"x", "y"}},
	}
	if !sameDeviceFilters(filters, want) {
		t.Errorf("got %+v", filters)
	}
	// Merged and reordered terms read as the same filter.
	merged, _, ok := parseDeviceCriteria("(and (or (and (severity 'MINOR') (severity 'CRITICAL'))) " + uiDelivery[1:] + ")")
	if !ok || !sameDeviceFilters(merged, []deviceFilter{{Severities: []string{"critical", "minor"}}}) {
		t.Errorf("merged %+v %v", merged, ok)
	}
}

func TestParseDeviceCriteriaFailures(t *testing.T) {
	cases := map[string]string{
		"cleared":          "(and (or (and (severity 'CLEARED'))) (or (notificationDeliveryMechanismNE 8 9)))",
		"invalid severity": "(and (or (and (severity 'INVALIDSEVERITY'))) (or (notificationDeliveryMechanismNE 8 9)))",
		"unknown term":     "(and (or (and (deviceName 'x'))) (or (notificationDeliveryMechanismNE 8 9)))",
		"no delivery":      "(and (or (and (severity 'MAJOR'))))",
		"wrong delivery":   "(and (or (and (severity 'MAJOR'))) (or (notificationDeliveryMechanismNE 1 2)))",
		"two categories":   "(and (or (and (category 'BGP' 'MPLS') (severity 'MAJOR'))) (or (notificationDeliveryMechanismNE 8 9)))",
		"empty filter":     "(and (or (and )) (or (notificationDeliveryMechanismNE 8 9)))",
		"unbalanced":       "(and (or (and (severity 'MAJOR'))",
		"not a list":       "severity",
		"trailing text":    "(and (or (and (severity 'MAJOR'))) (or (notificationDeliveryMechanismNE 8 9))) x",
		"two groups":       "(and (or (and (severity 'MAJOR'))) (or (notificationDeliveryMechanismNE 8 9))(deviceGroup 'a')(deviceGroup 'b'))",
	}
	for name, text := range cases {
		if filters, _, ok := parseDeviceCriteria(text); ok || filters != nil {
			t.Errorf("%s: parsed %+v", name, filters)
		}
	}
	// The group is still found when the filters cannot be read.
	if _, group, ok := parseDeviceCriteria("(and (or (and (severity 'CLEARED'))) (or (notificationDeliveryMechanismNE 8 9))(deviceGroup 'g'))"); ok || group != "g" {
		t.Errorf("group %q ok %v", group, ok)
	}
}

func TestDeviceCriteriaRoundTrip(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	severities := []string{"information", "warning", "minor", "major", "critical"}
	words := []string{"BGP", "Switches and Routers", "Power Module Error (PM_X)", "a:b", "x.y-z_1"}
	pick := func(from []string, max int) []string {
		var out []string
		for _, i := range rng.Perm(len(from))[:1+rng.Intn(max)] {
			out = append(out, from[i])
		}
		return out
	}
	for i := 0; i < 500; i++ {
		var filters []deviceFilter
		for n := 0; n <= rng.Intn(4); n++ {
			var filter deviceFilter
			switch rng.Intn(3) {
			case 0:
				filter.Severities = pick(severities, 5)
			case 1:
				filter.Category = words[rng.Intn(len(words))]
				filter.Severities = pick(severities, 5)
			default:
				filter.Events = pick(words, 5)
			}
			filters = append(filters, filter)
		}
		group := ""
		if rng.Intn(2) == 0 {
			group = "2661ed2d-11fb-49c4-9f9f-f71f1c12767b"
		}
		text := deviceCriteria(filters, group)
		got, gotGroup, ok := parseDeviceCriteria(text)
		if !ok || gotGroup != group || !sameDeviceFilters(got, filters) {
			t.Fatalf("round trip of %s: %+v %q %v", text, got, gotGroup, ok)
		}
		if again := deviceCriteria(got, gotGroup); again != text {
			t.Fatalf("not stable:\n%s\n%s", text, again)
		}
	}
}

func TestDeviceRulesAndReceivers(t *testing.T) {
	rules := deviceRules("p 1", "C", []string{"t1", "t2"}, []string{"s1"})
	want := []map[string]interface{}{
		{"actions": []string{"(sendTrap 't1')"}, "ruleName": "p 1_RESTCONF_RULE_0", "criteria": "C"},
		{"actions": []string{"(sendTrap 't2')"}, "ruleName": "p 1_RESTCONF_RULE_1", "criteria": "C"},
		{"actions": []string{"(sendSyslog 's1')"}, "ruleName": "p 1_RESTCONF_RULE_2", "criteria": "C"},
	}
	if !reflect.DeepEqual(rules, want) {
		t.Errorf("got %+v", rules)
	}
	traps, syslogs := receiversFromActions([]string{"(sendTrap 'b')", "(sendTrap 'a')", "(sendSyslog 's')", "(sendEmail 'x')", "(sendTrap 'a')"})
	if !reflect.DeepEqual(traps, []string{"a", "b"}) || !reflect.DeepEqual(syslogs, []string{"s"}) {
		t.Errorf("traps %v syslogs %v", traps, syslogs)
	}
}

func TestDeviceFromBody(t *testing.T) {
	criteria := "(and (or (and (severity 'CRITICAL'))) (or (notificationDeliveryMechanismNE 8 9))(deviceGroup 'g'))"
	body := `{"groupName":"p","alarmType":"Device","rules":[` +
		`{"ruleName":"p_RESTCONF_RULE_0","actions":["(sendTrap 't')"],"criteria":"` + criteria + `"},` +
		`{"ruleName":"p_RESTCONF_RULE_1","actions":["(sendSyslog 's')"],"criteria":"` + criteria + `"}]}`
	var data DeviceNotificationPolicies
	data.fromBody(context.Background(), gjson.Parse(body))
	if data.DeviceGroup.ValueString() != "g" || len(data.Filters) != 1 || len(data.TrapReceivers) != 1 || len(data.SyslogReceivers) != 1 || !data.Description.IsNull() {
		t.Fatalf("got %+v", data)
	}
	// Rules with different criteria: the filters read as null.
	differing := `{"groupName":"p","policyRules":[` +
		`{"actions":["(sendTrap 't')"],"criteria":"` + criteria + `"},` +
		`{"actions":["(sendSyslog 's')"],"criteria":"(and (or (and (severity 'MAJOR'))) (or (notificationDeliveryMechanismNE 8 9)))"}]}`
	var other DeviceNotificationPolicies
	other.fromBody(context.Background(), gjson.Parse(differing))
	if other.Filters != nil {
		t.Errorf("filters %+v", other.Filters)
	}
	// The body written by toBody is read back to the same attributes, with the order of the state kept.
	written := data.toBody(context.Background(), DeviceNotificationPolicies{})
	var back DeviceNotificationPolicies
	back.Name = data.Name
	back.updateFromBody(context.Background(), gjson.Parse(written))
	if back.DeviceGroup != data.DeviceGroup || len(back.Filters) != 1 || len(back.TrapReceivers) != 1 {
		t.Errorf("read back %+v from %s", back, written)
	}
}

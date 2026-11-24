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
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// A Device policy is stored by the controller as a Lisp criteria. The UI of the Device wizard writes
//
//	(and (or <filter> ...) (or (notificationDeliveryMechanismNE 8 9))(deviceGroup '<id>'))
//
// and one rule per receiver with the same criteria. The builder writes exactly that form (the golden tests hold the strings of the
// UI) and the parser reads it back, so that a change made in the UI shows as a difference.

const deviceDeliveryTerm = "(or (notificationDeliveryMechanismNE 8 9))"

// The severities as the controller names them, from high to low (the order the UI writes).
var deviceSeverities = []string{"CRITICAL", "MAJOR", "MINOR", "WARNING", "INFORMATION"}

type deviceFilter struct {
	Severities []string // lowercase names
	Category   string
	Events     []string
}

func severityRank(name string) int {
	for i, severity := range deviceSeverities {
		if severity == name {
			return i
		}
	}
	return len(deviceSeverities)
}

func quoteList(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, "'"+value+"'")
	}
	return strings.Join(quoted, " ")
}

func (filter deviceFilter) lisp() string {
	var terms []string
	if filter.Category != "" {
		terms = append(terms, "(category '"+filter.Category+"')")
	}
	if len(filter.Severities) > 0 {
		severities := make([]string, 0, len(filter.Severities))
		for _, severity := range filter.Severities {
			severities = append(severities, strings.ToUpper(severity))
		}
		sort.SliceStable(severities, func(i, j int) bool {
			ri, rj := severityRank(severities[i]), severityRank(severities[j])
			if ri != rj {
				return ri < rj
			}
			return severities[i] < severities[j]
		})
		terms = append(terms, "(severity "+quoteList(severities)+")")
	}
	if len(filter.Events) > 0 {
		events := append([]string(nil), filter.Events...)
		sort.Strings(events)
		terms = append(terms, "( eventType "+quoteList(events)+")")
	}
	return "(and " + strings.Join(terms, " ") + ")"
}

// deviceCriteria builds the criteria of a policy. The filters are sorted by their own text, so that the order in the
// configuration does not change the result. A group id is left out when empty.
func deviceCriteria(filters []deviceFilter, group string) string {
	lisps := make([]string, 0, len(filters))
	for _, filter := range filters {
		lisps = append(lisps, filter.lisp())
	}
	sort.Strings(lisps)
	criteria := "(and (or " + strings.Join(lisps, " ") + ") " + deviceDeliveryTerm
	if group != "" {
		criteria += "(deviceGroup '" + group + "')"
	}
	return criteria + ")"
}

func deviceFiltersOf(ctx context.Context, filters []DeviceNotificationPoliciesFilters) []deviceFilter {
	result := make([]deviceFilter, 0, len(filters))
	for _, item := range filters {
		filter := deviceFilter{
			Severities: setStrings(ctx, item.Severities),
			Events:     setStrings(ctx, item.Events),
		}
		if !item.Category.IsNull() && !item.Category.IsUnknown() {
			filter.Category = item.Category.ValueString()
		}
		result = append(result, filter)
	}
	return result
}

func deviceFiltersToModel(filters []deviceFilter) []DeviceNotificationPoliciesFilters {
	result := make([]DeviceNotificationPoliciesFilters, 0, len(filters))
	for _, filter := range filters {
		item := DeviceNotificationPoliciesFilters{
			Severities: targetSet(uniqueSorted(filter.Severities)),
			Events:     targetSet(uniqueSorted(filter.Events)),
			Category:   types.StringNull(),
		}
		if filter.Category != "" {
			item.Category = types.StringValue(filter.Category)
		}
		result = append(result, item)
	}
	return result
}

// sameDeviceFilters is true when both lists hold the same filters in any order. The schema holds the filters as a list (a set of
// objects with sets is not allowed), so the order of the state is kept when the controller holds the same filters.
func sameDeviceFilters(a, b []deviceFilter) bool {
	if len(a) != len(b) {
		return false
	}
	left, right := make([]string, 0, len(a)), make([]string, 0, len(b))
	for i := range a {
		left = append(left, a[i].lisp())
		right = append(right, b[i].lisp())
	}
	sort.Strings(left)
	sort.Strings(right)
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

// A node of the parsed Lisp: an atom (text, quoted or not) or a list.
type lispNode struct {
	atom   string
	quoted bool
	list   []lispNode
	isList bool
}

func (node lispNode) head() string {
	if node.isList && len(node.list) > 0 && !node.list[0].isList && !node.list[0].quoted {
		return node.list[0].atom
	}
	return ""
}

// parseLisp reads atoms, 'quoted' strings (no escapes) and lists; whitespace is not significant. It returns false on anything else.
func parseLisp(text string) (lispNode, bool) {
	pos := 0
	var read func() (lispNode, bool)
	read = func() (lispNode, bool) {
		for pos < len(text) && strings.ContainsRune(" \t\r\n", rune(text[pos])) {
			pos++
		}
		if pos >= len(text) {
			return lispNode{}, false
		}
		switch text[pos] {
		case '(':
			pos++
			node := lispNode{isList: true}
			for {
				for pos < len(text) && strings.ContainsRune(" \t\r\n", rune(text[pos])) {
					pos++
				}
				if pos >= len(text) {
					return lispNode{}, false
				}
				if text[pos] == ')' {
					pos++
					return node, true
				}
				child, ok := read()
				if !ok {
					return lispNode{}, false
				}
				node.list = append(node.list, child)
			}
		case ')':
			return lispNode{}, false
		case '\'':
			end := strings.IndexByte(text[pos+1:], '\'')
			if end < 0 {
				return lispNode{}, false
			}
			node := lispNode{atom: text[pos+1 : pos+1+end], quoted: true}
			pos += end + 2
			return node, true
		default:
			start := pos
			for pos < len(text) && !strings.ContainsRune(" \t\r\n()'", rune(text[pos])) {
				pos++
			}
			return lispNode{atom: text[start:pos]}, true
		}
	}
	node, ok := read()
	if !ok {
		return lispNode{}, false
	}
	for pos < len(text) && strings.ContainsRune(" \t\r\n", rune(text[pos])) {
		pos++
	}
	return node, pos == len(text)
}

func quotedArgs(node lispNode) ([]string, bool) {
	var values []string
	for _, arg := range node.list[1:] {
		if arg.isList || !arg.quoted {
			return nil, false
		}
		values = append(values, arg.atom)
	}
	return values, len(values) > 0
}

func parseDeviceFilter(node lispNode) (deviceFilter, bool) {
	var filter deviceFilter
	if node.head() != "and" || len(node.list) < 2 {
		return filter, false
	}
	seen := map[string]bool{}
	for _, term := range node.list[1:] {
		kind := term.head()
		values, ok := quotedArgs(term)
		if !ok {
			return filter, false
		}
		switch kind {
		case "category":
			if len(values) != 1 || seen[kind] {
				return filter, false
			}
			filter.Category = values[0]
		case "severity":
			for _, value := range values {
				if severityRank(value) == len(deviceSeverities) {
					return filter, false
				}
				filter.Severities = append(filter.Severities, strings.ToLower(value))
			}
		case "eventType":
			filter.Events = append(filter.Events, values...)
		default:
			return filter, false
		}
		seen[kind] = true
	}
	return filter, true
}

// parseDeviceCriteria reads a criteria back. The group is returned whenever its term is found, also when the filters cannot be read.
// ok is false when the filters cannot be read: a severity that the provider does not write (CLEARED, INVALIDSEVERITY), an unknown
// term, a missing delivery term, a category with several values (the UI cannot read it) or any other shape.
func parseDeviceCriteria(text string) (filters []deviceFilter, group string, ok bool) {
	root, valid := parseLisp(text)
	if !valid || root.head() != "and" {
		return nil, "", false
	}
	var filterOr *lispNode
	delivery := false
	ok = true
	for i := range root.list[1:] {
		child := root.list[1+i]
		switch {
		case child.head() == "deviceGroup":
			values, valid := quotedArgs(child)
			if !valid || len(values) != 1 || group != "" {
				ok = false
			} else {
				group = values[0]
			}
		case child.head() == "or" && len(child.list) == 2 && child.list[1].head() == "notificationDeliveryMechanismNE":
			args := child.list[1].list[1:]
			if delivery || len(args) != 2 || args[0].atom != "8" || args[1].atom != "9" {
				ok = false
			}
			delivery = true
		case child.head() == "or" && filterOr == nil:
			filterOr = &root.list[1+i]
		default:
			ok = false
		}
	}
	if !ok || !delivery || filterOr == nil || len(filterOr.list) < 2 {
		return nil, group, false
	}
	for _, node := range filterOr.list[1:] {
		filter, valid := parseDeviceFilter(node)
		if !valid {
			return nil, group, false
		}
		filters = append(filters, filter)
	}
	return filters, group, true
}

func uniqueSorted(values []string) []string {
	seen := map[string]bool{}
	var result []string
	for _, value := range values {
		if !seen[value] {
			seen[value] = true
			result = append(result, value)
		}
	}
	sort.Strings(result)
	return result
}

// deviceRules builds one rule per receiver: the trap receivers by name first, then the syslog receivers by name. A rule name is
// unique on the cluster, so it carries the name of the policy.
func deviceRules(name, criteria string, traps, syslogs []string) []map[string]interface{} {
	rules := []map[string]interface{}{}
	add := func(action string) {
		rules = append(rules, map[string]interface{}{
			"actions":  []string{action},
			"ruleName": name + "_RESTCONF_RULE_" + strconv.Itoa(len(rules)),
			"criteria": criteria,
		})
	}
	for _, trap := range traps {
		add("(sendTrap '" + trap + "')")
	}
	for _, syslog := range syslogs {
		add("(sendSyslog '" + syslog + "')")
	}
	return rules
}

// receiversFromActions reads the receivers from the actions of every rule; any other action is ignored.
func receiversFromActions(actions []string) (traps, syslogs []string) {
	for _, action := range actions {
		switch {
		case strings.HasPrefix(action, "(sendTrap '") && strings.HasSuffix(action, "')"):
			traps = append(traps, strings.TrimSuffix(strings.TrimPrefix(action, "(sendTrap '"), "')"))
		case strings.HasPrefix(action, "(sendSyslog '") && strings.HasSuffix(action, "')"):
			syslogs = append(syslogs, strings.TrimSuffix(strings.TrimPrefix(action, "(sendSyslog '"), "')"))
		}
	}
	return uniqueSorted(traps), uniqueSorted(syslogs)
}

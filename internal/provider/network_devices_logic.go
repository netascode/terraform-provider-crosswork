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
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	nodesPath          = "/crosswork/inventory/v1/nodes"
	nodesQueryPath     = nodesPath + "/query"
	tagsQueryPath      = "/crosswork/inventory/v1/tags/query"
	providersQueryPath = "/crosswork/inventory/v1/providers/query"

	adminDown = "ROBOT_ADMIN_STATE_DOWN"
)

// The controller fills in these keys and refuses them when they are sent back.
var nodeReadOnly = []string{"node_ip", "nso_state", "operational_state", "creation_time", "tag_names"}

var rowReadOnly = []string{"reachability_state", "reachability_state_upd_time", "error"}

// The controller has two wordings ("attached to a VDG and in admin up state", "while the node is in Admin State UP"); the
// retry goes by the message, not by a list of fields, so every field is covered.
func refusedWhileUp(message string) bool {
	message = strings.ToLower(message)
	return strings.Contains(message, "admin up state") || strings.Contains(message, "admin state up")
}

func nodeBody(node string) string {
	return `{"data":[` + node + `]}`
}

// Writes the address family as the number an add uses.
func cleanNode(node string) string {
	for _, k := range nodeReadOnly {
		node, _ = sjson.Delete(node, k)
	}
	rows := gjson.Get(node, "connectivity_info").Array()
	for i := range rows {
		for _, k := range rowReadOnly {
			node, _ = sjson.Delete(node, fmt.Sprintf("connectivity_info.%d.%s", i, k))
		}
		for j := range rows[i].Get("ipaddrs").Array() {
			node, _ = sjson.Set(node, fmt.Sprintf("connectivity_info.%d.ipaddrs.%d.inet_af", i, j), 0)
		}
	}
	return node
}

// Objects are merged, connectivity rows are merged by protocol, and every other value, an array included, is replaced.
func overlay(dst, prefix string, src gjson.Result) string {
	src.ForEach(func(k, v gjson.Result) bool {
		path := prefix + k.String()
		switch {
		case prefix == "" && k.String() == "connectivity_info":
			rows := []string{}
			for _, row := range v.Array() {
				base := "{}"
				gjson.Get(dst, "connectivity_info").ForEach(func(_, old gjson.Result) bool {
					if old.Get("type").String() == row.Get("type").String() {
						base = old.Raw
						return false
					}
					return true
				})
				rows = append(rows, overlay(base, "", row))
			}
			dst, _ = sjson.SetRaw(dst, "connectivity_info", "["+strings.Join(rows, ",")+"]")
		case prefix == "" && k.String() == "providers_family":
			dst, _ = sjson.SetRaw(dst, path, v.Raw)
		case v.IsObject():
			dst = overlay(dst, path+".", v)
		default:
			dst, _ = sjson.SetRaw(dst, path, v.Raw)
		}
		return true
	})
	return dst
}

func removed(dst, prefix string, state, plan gjson.Result) string {
	state.ForEach(func(k, v gjson.Result) bool {
		key := k.String()
		if prefix == "" && (key == "tags" || key == "connectivity_info") {
			return true
		}
		path := prefix + key
		switch p := plan.Get(key); {
		case !p.Exists():
			dst, _ = sjson.Delete(dst, path)
		case v.IsObject() && p.IsObject():
			dst = removed(dst, path+".", v, p)
		}
		return true
	})
	return dst
}

// userTagNames returns the names of the user-defined tags on the cluster. The tag_names of a device also hold the
// tags Crosswork adds itself, which are not managed.
func userTagNames(c *cw.Client) (map[string]bool, error) {
	names := map[string]bool{}
	for {
		body := fmt.Sprintf(`{"filter":{"tag_type":"TAG_TYPE_USER_DEFINED"},"next_from":"%d"}`, len(names))
		res, err := c.Post(tagsQueryPath, body, cw.ReadOnly)
		if err != nil {
			return nil, fmt.Errorf("listing the tags: %w", err)
		}
		before := len(names)
		for _, t := range res.Get("tags").Array() {
			names[t.Get("name").String()] = true
		}
		if len(names) == before || int64(len(names)) >= res.Get("total_count").Int() {
			return names, nil
		}
	}
}

func userTagObjects(node gjson.Result, user map[string]bool) string {
	tags := []string{}
	for _, name := range node.Get("tag_names").Array() {
		if user[name.String()] {
			tags = append(tags, fmt.Sprintf(`{"name":%q}`, name.String()))
		}
	}
	return "[" + strings.Join(tags, ",") + "]"
}

// The order of previous is kept, then the others, so only a real change shows.
func tagsOf(res gjson.Result, user map[string]bool, previous []NetworkDevicesTags) []NetworkDevicesTags {
	on := map[string]bool{}
	for _, name := range res.Get("data.0.tag_names").Array() {
		if user[name.String()] {
			on[name.String()] = true
		}
	}
	tags := make([]NetworkDevicesTags, 0, len(on))
	done := map[string]bool{}
	for _, t := range previous {
		if name := t.Name.ValueString(); on[name] && !done[name] {
			tags = append(tags, NetworkDevicesTags{Name: types.StringValue(name)})
			done[name] = true
		}
	}
	for _, name := range res.Get("data.0.tag_names").Array() {
		if on[name.String()] && !done[name.String()] {
			tags = append(tags, NetworkDevicesTags{Name: types.StringValue(name.String())})
			done[name.String()] = true
		}
	}
	return tags
}

// The edit replaces every field it leaves out. A change that the controller refuses while the device is admin up is
// retried as down, change, up.
func updateNetworkDevice(ctx context.Context, c *cw.Client, plan, state NetworkDevices) error {
	res, err := state.readGet(c, nodesQueryPath)
	if err != nil {
		return fmt.Errorf("reading the device: %w", err)
	}
	node := res.Get("data.0")
	if !node.Exists() {
		return fmt.Errorf("the device %s does not exist", state.Id.ValueString())
	}
	current := cleanNode(node.Raw)

	plannedNode, err := linkProvider(c, gjson.Parse(plan.toBody(ctx, state)).Get("data.0").Raw)
	if err != nil {
		return err
	}
	planned := gjson.Parse(plannedNode)
	previous := gjson.Parse(state.toBody(ctx, state)).Get("data.0")
	next := removed(overlay(current, "", planned), "", previous, planned)
	if plan.ProviderProfile.IsNull() && !state.ProviderProfile.IsNull() {
		next, _ = sjson.Delete(next, "providers_family")
	}
	// Tags are managed only when the configuration sets them (an empty list removes them all); otherwise the user
	// tags of the device are sent back unchanged. The API wants objects and the node lists names.
	switch {
	case plan.Tags != nil:
	case state.Tags != nil:
		next, _ = sjson.SetRaw(next, "tags", "[]")
	default:
		user, err := userTagNames(c)
		if err != nil {
			return err
		}
		next, _ = sjson.SetRaw(next, "tags", userTagObjects(node, user))
	}

	res, err = c.Put(nodesPath, nodeBody(next))
	if err == nil {
		return nil
	}
	if !refusedWhileUp(err.Error() + res.String()) {
		return err
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: the device is admin up, changing it as down, change, up", state.Id.ValueString()))
	downCurrent, _ := sjson.Set(current, "admin_state", adminDown)
	if _, err = c.Put(nodesPath, nodeBody(downCurrent)); err != nil {
		return err
	}
	restore := func() { _, _ = c.Put(nodesPath, nodeBody(current)) }
	downNext, _ := sjson.Set(next, "admin_state", adminDown)
	if _, err = c.Put(nodesPath, nodeBody(downNext)); err != nil {
		restore()
		return err
	}
	if gjson.Get(next, "admin_state").String() == adminDown || gjson.Get(next, "admin_state").String() == "ROBOT_ADMIN_STATE_UNMANAGED" {
		return nil
	}
	if _, err = c.Put(nodesPath, nodeBody(next)); err != nil {
		restore()
		return err
	}
	return nil
}

// providerPlaceholder is where the generated body writes the provider profile; the wire key is the provider name.
const providerPlaceholder = "providers_family.ROBOT_PROVIDER_NSO.providers.provider_name"

var keyEscaper = strings.NewReplacer(`\`, `\\`, ".", `\.`, "*", `\*`, "?", `\?`)

// linkProvider moves the provider profile from the placeholder to the block the API wants: keyed by the provider
// name, with the provider uuid and the device key (the inventory id).
func linkProvider(c *cw.Client, node string) (string, error) {
	name := gjson.Get(node, providerPlaceholder)
	if !name.Exists() {
		return node, nil
	}
	node, _ = sjson.Delete(node, "providers_family")
	body := fmt.Sprintf(`{"filter":{"name":%q},"filterData":{"PageSize":10,"Criteria":""}}`, name.String())
	res, err := c.Post(providersQueryPath, body, cw.ReadOnly)
	if err != nil {
		return node, fmt.Errorf("looking up the provider profile %s: %w", name.String(), err)
	}
	uuid := res.Get("data.0.uuid").String()
	if uuid == "" {
		return node, fmt.Errorf("the provider profile %s does not exist", name.String())
	}
	block := fmt.Sprintf(`{"provider_name":%q,"provider_node_id":%q,"provider_uuid":%q}`, name.String(), gjson.Get(node, "inventory_id").String(), uuid)
	node, _ = sjson.SetRaw(node, "providers_family.ROBOT_PROVIDER_NSO.providers."+keyEscaper.Replace(name.String()), block)
	return node, nil
}

func linkProviderBody(c *cw.Client, body string) (string, error) {
	node, err := linkProvider(c, gjson.Get(body, "data.0").Raw)
	if err != nil {
		return body, err
	}
	body, _ = sjson.SetRaw(body, "data.0", node)
	return body, nil
}

func providerProfileName(res gjson.Result) types.String {
	name := types.StringNull()
	res.Get("data.0.providers_family.ROBOT_PROVIDER_NSO.providers").ForEach(func(k, _ gjson.Result) bool {
		name = types.StringValue(k.String())
		return false
	})
	return name
}

// keepProviderProfile refreshes the provider profile only when the configuration sets it: a device that Crosswork
// attached to the built-in provider itself is not a difference.
func keepProviderProfile(previous types.String, res gjson.Result) types.String {
	if previous.IsNull() {
		return previous
	}
	return providerProfileName(res)
}

func hostNameAfterCreate(c *cw.Client, plan NetworkDevices) types.String {
	if !plan.HostName.IsUnknown() {
		return plan.HostName
	}
	res, err := plan.readGet(c, "/crosswork/inventory/v1/nodes/query")
	if err != nil {
		return types.StringNull()
	}
	return hostNameOf(plan.HostName, res)
}

// hostNameOf returns the host name for the state: the configured one, or, when the configuration leaves it out,
// the one the controller holds (null if none yet), so the plan never shows a change for a value the device owns.
func hostNameOf(planned types.String, node gjson.Result) types.String {
	if !planned.IsUnknown() {
		return planned
	}
	if value := node.Get("data.0.host_name"); value.Exists() && value.String() != "" {
		return types.StringValue(value.String())
	}
	return types.StringNull()
}

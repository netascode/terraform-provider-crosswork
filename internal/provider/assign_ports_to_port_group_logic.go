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
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const portPage = 500

// Two devices can share a host name and a port name, so a port is keyed by its device uuid and name.
type portKey struct{ device, name string }

// Add lists every port Crosswork has collected, members included; Manual lists the ports added by hand.
func listPorts(c *cw.Client, group, mode string) ([]gjson.Result, error) {
	var ports []gjson.Result
	for {
		res, err := c.Get(fmt.Sprintf("%s/port/%s?start=%d&end=%d&mode=%s", groupingPath, group, len(ports), len(ports)+portPage, mode), cw.ReadOnly)
		if err != nil {
			return nil, err
		}
		page := res.Get("ports").Array()
		ports = append(ports, page...)
		if len(page) == 0 || int64(len(ports)) >= res.Get("total").Int() {
			return ports, nil
		}
	}
}

func portsByKey(ports []gjson.Result) map[portKey]gjson.Result {
	m := map[portKey]gjson.Result{}
	for _, p := range ports {
		m[portKey{p.Get("deviceUuid").String(), p.Get("portName").String()}] = p
	}
	return m
}

func portDevices(c *cw.Client, ports []AssignPortsToPortGroupPorts) (map[string]string, error) {
	seen := map[string]bool{}
	var ids []string
	for _, p := range ports {
		if id := p.Device.ValueString(); !seen[id] {
			seen[id] = true
			ids = append(ids, id)
		}
	}
	return deviceUUIDs(c, ids)
}

func portMembers(c *cw.Client, group string) ([]AssignPortsToPortGroupPorts, error) {
	members, err := listPorts(c, group, "Manual")
	if err != nil {
		return nil, err
	}
	uuids := map[string]bool{}
	for _, m := range members {
		uuids[m.Get("deviceUuid").String()] = true
	}
	list := make([]string, 0, len(uuids))
	for u := range uuids {
		list = append(list, u)
	}
	ids, err := inventoryIDs(c, list)
	if err != nil {
		return nil, err
	}
	// An empty, not a null, set when the group has no ports, to match an empty list in the configuration.
	out := []AssignPortsToPortGroupPorts{}
	for _, m := range members {
		id, ok := ids[m.Get("deviceUuid").String()]
		if !ok {
			continue
		}
		out = append(out, AssignPortsToPortGroupPorts{Device: types.StringValue(id), PortName: types.StringValue(m.Get("portName").String())})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Device.ValueString() != out[j].Device.ValueString() {
			return out[i].Device.ValueString() < out[j].Device.ValueString()
		}
		return out[i].PortName.ValueString() < out[j].PortName.ValueString()
	})
	return out, nil
}

// The members are read back, because a port id that Crosswork does not know is answered with success and adds nothing.
func setGroupPorts(c *cw.Client, group string, wanted []AssignPortsToPortGroupPorts) error {
	devices, err := portDevices(c, wanted)
	if err != nil {
		return err
	}
	want := map[portKey]AssignPortsToPortGroupPorts{}
	for _, p := range wanted {
		want[portKey{devices[p.Device.ValueString()], p.PortName.ValueString()}] = p
	}
	manual, err := listPorts(c, group, "Manual")
	if err != nil {
		return err
	}
	members := portsByKey(manual)

	var remove []interface{}
	for key, obj := range members {
		if _, ok := want[key]; !ok {
			remove = append(remove, obj.Value())
		}
	}
	if len(remove) > 0 {
		body, _ := sjson.Set("", "ports", remove)
		if res, err := c.DeleteBody(fmt.Sprintf("%s/port/%s", groupingPath, group), body); err != nil {
			return fmt.Errorf("removing ports from the group: %w, %s", err, res.String())
		}
	}

	var add []map[string]string
	var missing []string
	var candidates map[portKey]gjson.Result
	for key, p := range want {
		if _, ok := members[key]; ok {
			continue
		}
		if candidates == nil {
			all, err := listPorts(c, group, "Add")
			if err != nil {
				return err
			}
			candidates = portsByKey(all)
		}
		cand, ok := candidates[key]
		if !ok {
			missing = append(missing, fmt.Sprintf("%s on %s", p.PortName.ValueString(), p.Device.ValueString()))
			continue
		}
		add = append(add, map[string]string{"portUuid": cand.Get("portUuid").String(), "deviceUuid": key.device})
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("Crosswork does not know the port(s) %v: ports exist once Crosswork has collected the device", missing)
	}
	if len(add) > 0 {
		body, _ := sjson.Set("", "ports", add)
		if res, err := c.Post(fmt.Sprintf("%s/port/%s", groupingPath, group), body); err != nil {
			return fmt.Errorf("adding ports to the group: %w, %s", err, res.String())
		}
	}

	after, err := listPorts(c, group, "Manual")
	if err != nil {
		return err
	}
	have := portsByKey(after)
	for key, p := range want {
		if _, ok := have[key]; !ok {
			return fmt.Errorf("the port %s of %s is not a member of the group after it was added", p.PortName.ValueString(), p.Device.ValueString())
		}
	}
	return nil
}

func unsetGroupPorts(c *cw.Client, group string, listed []AssignPortsToPortGroupPorts) error {
	manual, err := listPorts(c, group, "Manual")
	if err != nil {
		return err
	}
	uuids := map[string]bool{}
	for _, m := range manual {
		uuids[m.Get("deviceUuid").String()] = true
	}
	list := make([]string, 0, len(uuids))
	for u := range uuids {
		list = append(list, u)
	}
	ids, err := inventoryIDs(c, list)
	if err != nil {
		return err
	}
	want := map[[2]string]bool{}
	for _, p := range listed {
		want[[2]string{p.Device.ValueString(), p.PortName.ValueString()}] = true
	}
	var remove []interface{}
	for _, m := range manual {
		if want[[2]string{ids[m.Get("deviceUuid").String()], m.Get("portName").String()}] {
			remove = append(remove, m.Value())
		}
	}
	if len(remove) == 0 {
		return nil
	}
	body, _ := sjson.Set("", "ports", remove)
	if res, err := c.DeleteBody(fmt.Sprintf("%s/port/%s", groupingPath, group), body); err != nil {
		return fmt.Errorf("removing ports from the group: %w, %s", err, res.String())
	}
	return nil
}

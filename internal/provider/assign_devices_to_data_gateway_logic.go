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

// Written by hand: a device joins a data gateway by an ADD, leaves it by a REMOVE and moves to another by an UPDATE, and none of
// them can be repeated. The members are read from the node list, and a job that fails can have changed part of the request, so
// every failure is followed by a read of the members.

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/sjson"
)

const (
	deviceMappingPath = "/crosswork/inventory/v1/dg/devicemapping"

	mappingAdd    = "ADD_OPER"
	mappingRemove = "REMOVE_OPER"
	mappingUpdate = "UPDATE_OPER"

	listAllThreshold = 25
)

// The id goes into the criteria of a node query, so nothing else is accepted.
var gatewayIdPattern = regexp.MustCompile(`^[0-9A-Za-z-]{8,64}$`)

type gatewayDevice struct {
	Name    string
	Uuid    string
	Gateway string
}

type mappingChanges struct {
	Add    []string
	Move   []string
	Remove []string
}

func (m mappingChanges) empty() bool { return len(m.Add)+len(m.Move)+len(m.Remove) == 0 }

// A member that is not in the previous list belongs to someone else and is left alone.
func planMapping(gatewayId string, desired map[string]gatewayDevice, members map[string]string, previous []string) mappingChanges {
	var m mappingChanges
	for _, d := range desired {
		switch {
		case d.Gateway == gatewayId:
		case d.Gateway == "":
			m.Add = append(m.Add, d.Uuid)
		default:
			m.Move = append(m.Move, d.Uuid)
		}
	}
	for _, name := range previous {
		if _, wanted := desired[name]; wanted {
			continue
		}
		if uuid, member := members[name]; member {
			m.Remove = append(m.Remove, uuid)
		}
	}
	sort.Strings(m.Add)
	sort.Strings(m.Move)
	sort.Strings(m.Remove)
	return m
}

func attachedDevices(c *cw.Client, gatewayId string) (map[string]string, error) {
	if !gatewayIdPattern.MatchString(gatewayId) {
		return nil, fmt.Errorf("%q is not the id of a data gateway", gatewayId)
	}
	items, err := c.Paginate(&cw.FilterDataPager{Path: nodesQueryPath, Size: 500,
		Criteria: "select uuid from robotnodedata where dg_uuid=" + gatewayId})
	if err != nil {
		return nil, fmt.Errorf("listing the devices of the data gateway: %w", err)
	}
	members := map[string]string{}
	for _, item := range items {
		if name := item.Get("inventory_id").String(); name != "" {
			members[name] = item.Get("uuid").String()
		}
	}
	return members, nil
}

// The match ignores case, as Crosswork does for an inventory id.
func lookupDevices(c *cw.Client, names []string) (map[string]gatewayDevice, error) {
	found := map[string]gatewayDevice{}
	byLower := map[string]string{}
	for _, n := range names {
		byLower[strings.ToLower(n)] = n
	}
	take := func(inventoryId, uuid, gateway string) {
		if name, ok := byLower[strings.ToLower(inventoryId)]; ok {
			found[name] = gatewayDevice{Name: name, Uuid: uuid, Gateway: gateway}
		}
	}
	if len(names) > listAllThreshold {
		items, err := c.Paginate(&cw.FilterDataPager{Path: nodesQueryPath, Size: 500})
		if err != nil {
			return nil, fmt.Errorf("listing the devices: %w", err)
		}
		for _, item := range items {
			take(item.Get("inventory_id").String(), item.Get("uuid").String(), item.Get("dg_uuid").String())
		}
	} else {
		for _, n := range names {
			body, _ := sjson.Set(`{"filter":{}}`, "filter.inventory_id", n)
			res, err := c.Post(nodesQueryPath, body, cw.ReadOnly)
			if err != nil {
				return nil, fmt.Errorf("looking up the device %s: %w", n, err)
			}
			for _, item := range res.Get("data").Array() {
				take(item.Get("inventory_id").String(), item.Get("uuid").String(), item.Get("dg_uuid").String())
			}
		}
	}
	for _, n := range names {
		if _, ok := found[n]; !ok {
			return nil, fmt.Errorf("the device %s does not exist", n)
		}
	}
	return found, nil
}

func mapDevices(c *cw.Client, gatewayId, operation string, uuids []string) error {
	body, _ := sjson.Set(`{"dgDeviceMappings":[{}]}`, "dgDeviceMappings.0.cdg_duuid", gatewayId)
	body, _ = sjson.Set(body, "dgDeviceMappings.0.mapping_oper", operation)
	body, _ = sjson.Set(body, "dgDeviceMappings.0.device_uuid", uuids)
	res, err := c.Put(deviceMappingPath, body)
	if err != nil {
		return fmt.Errorf("%w, %s", err, res.String())
	}
	return nil
}

func applyMapping(c *cw.Client, gatewayId string, m mappingChanges) error {
	for _, step := range []struct {
		operation string
		uuids     []string
	}{{mappingRemove, m.Remove}, {mappingAdd, m.Add}, {mappingUpdate, m.Move}} {
		if len(step.uuids) == 0 {
			continue
		}
		if err := mapDevices(c, gatewayId, step.operation, step.uuids); err != nil {
			return err
		}
	}
	return nil
}

func reached(desired, previous []string, members map[string]string) bool {
	want := map[string]bool{}
	for _, n := range desired {
		want[n] = true
		if _, ok := members[n]; !ok {
			return false
		}
	}
	for _, n := range previous {
		if _, ok := members[n]; ok && !want[n] {
			return false
		}
	}
	return true
}

// A job that fails can have changed part of the request, so after a failure the members are read again.
func syncDevices(c *cw.Client, gatewayId string, desired, previous []string) error {
	want := map[string]gatewayDevice{}
	if len(desired) > 0 {
		var err error
		if want, err = lookupDevices(c, desired); err != nil {
			return err
		}
	}
	members, err := attachedDevices(c, gatewayId)
	if err != nil {
		return err
	}
	changes := planMapping(gatewayId, want, members, previous)
	if changes.empty() {
		return nil
	}
	err = applyMapping(c, gatewayId, changes)
	if err == nil {
		return nil
	}
	if after, readErr := attachedDevices(c, gatewayId); readErr == nil && reached(desired, previous, after) {
		return nil
	}
	return err
}

func detachDevices(c *cw.Client, gatewayId string, previous []string) error {
	members, err := attachedDevices(c, gatewayId)
	if err != nil {
		return err
	}
	changes := planMapping(gatewayId, nil, members, previous)
	if changes.empty() {
		return nil
	}
	err = applyMapping(c, gatewayId, changes)
	if err == nil {
		return nil
	}
	if after, readErr := attachedDevices(c, gatewayId); readErr == nil && reached(nil, previous, after) {
		return nil
	}
	return err
}

func currentDevices(c *cw.Client, gatewayId string, previous []string) ([]string, error) {
	members, err := attachedDevices(c, gatewayId)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, n := range previous {
		if _, ok := members[n]; ok {
			out = append(out, n)
		}
	}
	sort.Strings(out)
	return out, nil
}

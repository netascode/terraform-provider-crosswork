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
	"time"

	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	memberPage      = 100
	sourceSearches  = 4
	sourceRetryWait = 2 * time.Second
)

func groupMembers(c *cw.Client, group string) ([]string, error) {
	var uuids []string
	for {
		res, err := c.Get(fmt.Sprintf("%s/device/%s?start=%d&end=%d", groupingPath, group, len(uuids), len(uuids)+memberPage), cw.ReadOnly)
		if err != nil {
			return nil, err
		}
		before := len(uuids)
		for _, d := range res.Get("devices").Array() {
			uuids = append(uuids, d.Get("uuid").String())
		}
		if len(uuids) == before || int64(len(uuids)) >= res.Get("total").Int() {
			return uuids, nil
		}
	}
}

// A group that does not exist answers GROUP_NOT_EXIST, which the client recognises as not found.
func groupExists(c *cw.Client, group string) (bool, error) {
	if _, err := c.Get(fmt.Sprintf("%s/group/%s/details", groupingPath, group), cw.ReadOnly); err != nil {
		if cw.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Below this count, one call per device is cheaper than paging the inventory list.
var bulkLookupMin = 8

// The controller returns at most 1000 nodes per page.
const nodePage = 1000

// A partial list is reported so the caller can look up the devices it did not find one by one.
func listNodeKeys(c *cw.Client, maxPages int) (pairs [][2]string, complete bool, err error) {
	for page := 0; page < maxPages; page++ {
		body := fmt.Sprintf(`{"limit":%d,"filter":{},"filterData":{"PageSize":%d,"PageNum":%d}}`, nodePage, nodePage, page)
		res, err := c.Post(nodesQueryPath, body, cw.ReadOnly)
		if err != nil {
			return nil, false, fmt.Errorf("listing the devices: %w", err)
		}
		nodes := res.Get("data").Array()
		for _, n := range nodes {
			pairs = append(pairs, [2]string{n.Get("uuid").String(), n.Get("inventory_id").String()})
		}
		if len(nodes) < nodePage || int64(len(pairs)) >= res.Get("total_count").Int() {
			return pairs, true, nil
		}
	}
	return pairs, false, nil
}

func bulkPages(count int) int {
	if count < bulkLookupMin {
		return 0
	}
	return count / bulkLookupMin
}

func deviceUUIDs(c *cw.Client, inventoryIDs []string) (map[string]string, error) {
	uuids := map[string]string{}
	if pages := bulkPages(len(inventoryIDs)); pages > 0 {
		pairs, _, err := listNodeKeys(c, pages)
		if err != nil {
			return nil, err
		}
		byID := map[string]string{}
		for _, p := range pairs {
			byID[p[1]] = p[0]
		}
		for _, id := range inventoryIDs {
			if uuid, ok := byID[id]; ok {
				uuids[id] = uuid
			}
		}
	}
	// Looking up one by one also decides whether the device exists.
	for _, id := range inventoryIDs {
		if _, ok := uuids[id]; ok {
			continue
		}
		body, _ := sjson.Set(`{"filter":{}}`, "filter.inventory_id", id)
		res, err := c.Post(nodesQueryPath, body, cw.ReadOnly)
		if err != nil {
			return nil, fmt.Errorf("looking up the device %s: %w", id, err)
		}
		uuid := res.Get("data.0.uuid").String()
		if uuid == "" {
			return nil, fmt.Errorf("the device %s does not exist", id)
		}
		uuids[id] = uuid
	}
	return uuids, nil
}

func inventoryIDs(c *cw.Client, uuids []string) (map[string]string, error) {
	ids := map[string]string{}
	if pages := bulkPages(len(uuids)); pages > 0 {
		pairs, complete, err := listNodeKeys(c, pages)
		if err != nil {
			return nil, err
		}
		byUUID := map[string]string{}
		for _, p := range pairs {
			byUUID[p[0]] = p[1]
		}
		for _, uuid := range uuids {
			if id, ok := byUUID[uuid]; ok && id != "" {
				ids[uuid] = id
			}
		}
		if complete {
			// every node was listed, so a uuid that is not in the list has no node
			return ids, nil
		}
	}
	for _, uuid := range uuids {
		if _, ok := ids[uuid]; ok {
			continue
		}
		body, _ := sjson.Set(`{"filter":{}}`, "filter.uuid", uuid)
		res, err := c.Post(nodesQueryPath, body, cw.ReadOnly)
		if err != nil {
			return nil, fmt.Errorf("looking up the device %s: %w", uuid, err)
		}
		if id := res.Get("data.0.inventory_id").String(); id != "" {
			ids[uuid] = id
		}
	}
	return ids, nil
}

// A device sits in exactly one location group, and the copy into an access group needs that group as its source
// (an ancestor is refused). Groups that report operations.showMem hold no devices; so do the root and All Locations.
func locationGroups(c *cw.Client) ([]string, error) {
	roots, err := c.Get(groupingPath+"/group/root/LocationDevices/uuid", cw.ReadOnly)
	if err != nil {
		return nil, err
	}
	var groups []string
	for _, root := range roots.Array() {
		tree, err := c.Get(fmt.Sprintf("%s/groups/%s?brief=true&direct=false", groupingPath, root.String()), cw.ReadOnly)
		if err != nil {
			return nil, err
		}
		var walk func(n gjson.Result)
		walk = func(n gjson.Result) {
			if uuid := n.Get("uuid").String(); uuid != "" && (!n.Get("operations").Exists() || n.Get("operations.showMem").Bool()) {
				groups = append(groups, uuid)
			}
			for _, child := range n.Get("children").Array() {
				walk(child)
			}
		}
		for _, n := range tree.Array() {
			walk(n)
		}
	}
	return groups, nil
}

// The groups are read one after another, so a device being moved at that moment can be missed; a device is always in
// some location group, so the search is repeated for the ones not found.
func sourceGroups(c *cw.Client, uuids []string) (map[string]string, error) {
	found := map[string]string{}
	for attempt := 0; attempt < sourceSearches; attempt++ {
		if attempt > 0 {
			time.Sleep(sourceRetryWait)
		}
		var missing []string
		for _, u := range uuids {
			if _, ok := found[u]; !ok {
				missing = append(missing, u)
			}
		}
		if len(missing) == 0 {
			break
		}
		got, err := searchSourceGroups(c, missing)
		if err != nil {
			return nil, err
		}
		for u, g := range got {
			found[u] = g
		}
	}
	return found, nil
}

func searchSourceGroups(c *cw.Client, uuids []string) (map[string]string, error) {
	want := map[string]bool{}
	for _, u := range uuids {
		want[u] = true
	}
	groups, err := locationGroups(c)
	if err != nil {
		return nil, err
	}
	found := map[string]string{}
	for _, group := range groups {
		if len(found) == len(want) {
			break
		}
		members, err := groupMembers(c, group)
		if err != nil {
			return nil, err
		}
		for _, m := range members {
			if want[m] {
				found[m] = group
			}
		}
	}
	return found, nil
}

func addMembers(c *cw.Client, group string, uuids []string) error {
	sources, err := sourceGroups(c, uuids)
	if err != nil {
		return fmt.Errorf("finding the location group of the devices: %w", err)
	}
	bySource := map[string][]string{}
	for _, u := range uuids {
		source, ok := sources[u]
		if !ok {
			return fmt.Errorf("the device %s is in no location group", u)
		}
		bySource[source] = append(bySource[source], u)
	}
	for source, refs := range bySource {
		body, _ := sjson.Set("", "groupUuid", source)
		body, _ = sjson.Set(body, "newGroupUuid", group)
		body, _ = sjson.Set(body, "references", refs)
		if res, err := c.Post(groupingPath+"/group/member/copy", body); err != nil {
			return fmt.Errorf("adding devices to the group: %w, %s", err, res.String())
		}
	}
	return nil
}

func removeMembers(c *cw.Client, group string, uuids []string) error {
	body, _ := sjson.Set("", "references", uuids)
	if res, err := c.Put(fmt.Sprintf("%s/group/%s/members", groupingPath, group), body); err != nil {
		return fmt.Errorf("removing devices from the group: %w, %s", err, res.String())
	}
	return nil
}

func groupDevices(c *cw.Client, group string) ([]string, error) {
	members, err := groupMembers(c, group)
	if err != nil {
		return nil, err
	}
	ids, err := inventoryIDs(c, members)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(members))
	for _, m := range members {
		if id, ok := ids[m]; ok {
			names = append(names, id)
		}
	}
	sort.Strings(names)
	return names, nil
}

// The members are read back, because a copy of a device that cannot be copied still answers success.
func setGroupDevices(c *cw.Client, group string, names []string) error {
	wanted, err := deviceUUIDs(c, names)
	if err != nil {
		return err
	}
	current, err := groupMembers(c, group)
	if err != nil {
		return err
	}
	have := map[string]bool{}
	for _, m := range current {
		have[m] = true
	}
	isWanted := map[string]bool{}
	var add []string
	for _, uuid := range wanted {
		isWanted[uuid] = true
		if !have[uuid] {
			add = append(add, uuid)
		}
	}
	var remove []string
	for _, m := range current {
		if !isWanted[m] {
			remove = append(remove, m)
		}
	}
	sort.Strings(add)
	sort.Strings(remove)
	if len(remove) > 0 {
		if err := removeMembers(c, group, remove); err != nil {
			return err
		}
	}
	if len(add) > 0 {
		if err := addMembers(c, group, add); err != nil {
			return err
		}
	}
	after, err := groupMembers(c, group)
	if err != nil {
		return err
	}
	member := map[string]bool{}
	for _, m := range after {
		member[m] = true
	}
	for name, uuid := range wanted {
		if !member[uuid] {
			return fmt.Errorf("the device %s is not a member of the group after it was added", name)
		}
	}
	return nil
}

func unsetGroupDevices(c *cw.Client, group string, names []string) error {
	current, err := groupMembers(c, group)
	if err != nil {
		return err
	}
	ids, err := inventoryIDs(c, current)
	if err != nil {
		return err
	}
	listed := map[string]bool{}
	for _, n := range names {
		listed[n] = true
	}
	var remove []string
	for _, m := range current {
		if listed[ids[m]] {
			remove = append(remove, m)
		}
	}
	if len(remove) == 0 {
		return nil
	}
	return removeMembers(c, group, remove)
}

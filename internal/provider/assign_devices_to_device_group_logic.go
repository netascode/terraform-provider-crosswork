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

	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/sjson"
)

// A device is in exactly one location group, so moving it takes it out of the group it was in.
func moveMembers(c *cw.Client, group string, uuids []string) error {
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
		if res, err := c.Post(groupingPath+"/group/member", body); err != nil {
			return fmt.Errorf("moving devices into the group: %w, %s", err, res.String())
		}
	}
	return nil
}

// Removed devices go to Unassigned Devices. The members are read back, because a move that cannot be done may still
// answer success.
func setLocationGroupDevices(c *cw.Client, group string, names []string) error {
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
		if err := moveMembers(c, group, add); err != nil {
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
			return fmt.Errorf("the device %s is not a member of the group after it was moved", name)
		}
	}
	return nil
}

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
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Written by hand because the write sends constants the generator has no key for (merge, and the errors of every
// row), because the controller refuses an empty list (it is cleared by deleting the rows that were read), and because a
// delete that names one row that is gone is refused whole.

// lcmThreshold is one row of the list of a domain, keyed by its node and interface.
type lcmThreshold struct {
	node, iface string
	threshold   int64
}

func (data LocalCongestionMitigationInterfaceThresholds) domain() string {
	if id := data.DomainId.ValueString(); id != "" {
		return id
	}
	return data.Id.ValueString()
}

// wanted returns the rows of the configuration.
func (data LocalCongestionMitigationInterfaceThresholds) wanted() []lcmThreshold {
	rows := make([]lcmThreshold, 0, len(data.InterfaceThresholds))
	for _, item := range data.InterfaceThresholds {
		rows = append(rows, lcmThreshold{item.Node.ValueString(), item.Interface.ValueString(), item.UtilizationThreshold.ValueInt64()})
	}
	return rows
}

// lcmThresholdRows builds the rows of a request; the controller expects an errors list in every row.
func lcmThresholdRows(rows []lcmThreshold) []interface{} {
	list := make([]interface{}, 0, len(rows))
	for _, r := range rows {
		list = append(list, map[string]interface{}{"node": r.node, "interface": r.iface, "threshold": r.threshold, "errors": []interface{}{}})
	}
	return list
}

// lcmThresholdsPutBody replaces the whole list: merge is always sent, the controller fails without it.
func lcmThresholdsPutBody(domain string, rows []lcmThreshold) string {
	body, _ := sjson.Set("", "id", domain)
	body, _ = sjson.Set(body, "merge", false)
	body, _ = sjson.Set(body, "addInts", lcmThresholdRows(rows))
	return body
}

func lcmThresholdsDeleteBody(domain string, rows []lcmThreshold) string {
	body, _ := sjson.Set("", "id", domain)
	body, _ = sjson.Set(body, "data", lcmThresholdRows(rows))
	return body
}

// lcmThresholdsFromBody reads the answer of the list call: a decimal string per row, kept as its whole part.
func lcmThresholdsFromBody(res gjson.Result) []lcmThreshold {
	var rows []lcmThreshold
	for _, v := range res.Get("data").Array() {
		rows = append(rows, lcmThreshold{v.Get("node").String(), v.Get("interface").String(), v.Get("threshold").Int()})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].node != rows[j].node {
			return rows[i].node < rows[j].node
		}
		return rows[i].iface < rows[j].iface
	})
	return rows
}

// readLcmThresholds returns the rows the controller holds for the domain.
func readLcmThresholds(c *cw.Client, data LocalCongestionMitigationInterfaceThresholds) ([]lcmThreshold, cw.Res, error) {
	res, err := data.readGet(c, data.getPathGet())
	if err != nil {
		return nil, res, err
	}
	return lcmThresholdsFromBody(res), res, nil
}

// setLcmThresholds makes the list of the domain equal to wanted. An empty list is cleared by deleting the rows read just
// before, because the controller refuses to write an empty list.
func setLcmThresholds(c *cw.Client, data LocalCongestionMitigationInterfaceThresholds, wanted []lcmThreshold) error {
	if len(wanted) == 0 {
		return clearLcmThresholds(c, data)
	}
	_, err := c.Put(data.getPath(), lcmThresholdsPutBody(data.domain(), wanted), cw.UseMutex)
	return err
}

func clearLcmThresholds(c *cw.Client, data LocalCongestionMitigationInterfaceThresholds) error {
	current, _, err := readLcmThresholds(c, data)
	if err != nil || len(current) == 0 {
		return err
	}
	_, err = c.DeleteBody(data.getPathDelete(), lcmThresholdsDeleteBody(data.domain(), current), cw.UseMutex)
	return err
}

// lcmThresholdItems converts rows to the attribute. An empty list stays an empty set when the state already held one,
// so a configuration of `[]` does not show a difference.
func lcmThresholdItems(rows []lcmThreshold, prior []LocalCongestionMitigationInterfaceThresholdsInterfaceThresholds) []LocalCongestionMitigationInterfaceThresholdsInterfaceThresholds {
	if len(rows) == 0 {
		if prior != nil {
			return []LocalCongestionMitigationInterfaceThresholdsInterfaceThresholds{}
		}
		return nil
	}
	items := make([]LocalCongestionMitigationInterfaceThresholdsInterfaceThresholds, 0, len(rows))
	for _, r := range rows {
		items = append(items, LocalCongestionMitigationInterfaceThresholdsInterfaceThresholds{
			Node:                 types.StringValue(r.node),
			Interface:            types.StringValue(r.iface),
			UtilizationThreshold: types.Int64Value(r.threshold),
		})
	}
	return items
}

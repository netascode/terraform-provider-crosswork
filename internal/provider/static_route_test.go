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
	"testing"

	"github.com/tidwall/gjson"
)

const routesAnswer = `{"ztpStaticRoutes":[{"uuid":"b19e","subnet":"192.0.2.0","mask":"25","status":"success"},{"uuid":"ecdc","subnet":"192.0.2.0","mask":"24","status":"add-inprogress"}],"code":200}`

func TestStaticRouteBodyIsOneRouteInTheList(t *testing.T) {
	add := gjson.Parse(routeBody("192.0.2.0", "24", ""))
	assert.Equal(t, "192.0.2.0", add.Get("staticroutes.0.subnet").String())
	assert.Equal(t, "24", add.Get("staticroutes.0.mask").String(), "the prefix length is a string")
	assert.False(t, add.Get("staticroutes.0.uuid").Exists())
	remove := gjson.Parse(routeBody("192.0.2.0", "24", "ecdc"))
	assert.Equal(t, "ecdc", remove.Get("staticroutes.0.uuid").String(), "a removal names the uuid")
}

func TestStaticRouteIsFoundByIdOrByAddress(t *testing.T) {
	res := gjson.Parse(routesAnswer)
	route, ok := findStaticRoute(res, "b19e", "", "")
	assert.True(t, ok)
	assert.Equal(t, "25", route.Get("mask").String())
	route, ok = findStaticRoute(res, "", "192.0.2.0", "24")
	assert.True(t, ok)
	assert.Equal(t, "ecdc", route.Get("uuid").String(), "the same subnet with another mask is another route")
	_, ok = findStaticRoute(res, "gone", "", "")
	assert.False(t, ok)
	_, ok = findStaticRoute(gjson.Parse(`{"code":200}`), "", "192.0.2.0", "24")
	assert.False(t, ok, "an answer with no routes has none")
}

func TestStaticRouteRefusalIsACodeInTheBody(t *testing.T) {
	assert.Equal(t, "", staticRouteAnswer(gjson.Parse(`{"code":201,"message":"Add static route is initiated."}`)))
	assert.Contains(t, staticRouteAnswer(gjson.Parse(`{"code":400,"message":"192.0.2.0/33 : invalid subnet mask range for ipv4 address"}`)), "invalid subnet mask range")
}

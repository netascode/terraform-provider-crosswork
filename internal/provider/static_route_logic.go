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
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cw "github.com/netascode/go-crosswork"
)

const (
	staticRoutesPath      = "/crosswork/ztp/v1/staticroutes"
	staticRoutesQueryPath = "/crosswork/ztp/v1/staticroutes/query"
)

// The controller adds and removes a static route in the background: the answer is `{"code":201,...}` and the status of the route
// goes from add-inprogress to success (and from delete-inprogress to gone). A refusal is an HTTP 200 with a code of 400 and a message.

// A variable so that tests can shorten it.
var (
	staticRouteWait = 120 * time.Second
	staticRoutePoll = 2 * time.Second
)

func staticRouteAnswer(res gjson.Result) string {
	if code := res.Get("code"); code.Exists() && code.Int() >= 400 {
		return res.Get("message").String()
	}
	return ""
}

// Falls back to the subnet and prefix length when the id is not known.
func findStaticRoute(res gjson.Result, id, subnet, mask string) (gjson.Result, bool) {
	var found gjson.Result
	ok := false
	res.Get("ztpStaticRoutes").ForEach(func(_, route gjson.Result) bool {
		if id != "" && route.Get("uuid").String() == id || id == "" && route.Get("subnet").String() == subnet && route.Get("mask").String() == mask {
			found, ok = route, true
			return false
		}
		return true
	})
	return found, ok
}

func (data StaticRoute) readRoute(c *cw.Client) (cw.Res, error) {
	res, err := c.Post(staticRoutesQueryPath, `{"filter":{}}`, cw.ReadOnly)
	if err != nil {
		return res, err
	}
	route, ok := findStaticRoute(res, data.Id.ValueString(), data.Address.ValueString(), strconv.FormatInt(data.PrefixLength.ValueInt64(), 10))
	if !ok || route.Get("status").String() == "delete-inprogress" {
		return res, fmt.Errorf("StatusCode 404: the controller has no static route %s/%d", data.Address.ValueString(), data.PrefixLength.ValueInt64())
	}
	return route, nil
}

// One route in the list `staticroutes`.
func routeBody(subnet, mask, id string) string {
	route := ""
	if id != "" {
		route, _ = sjson.Set(route, "uuid", id)
	}
	route, _ = sjson.Set(route, "subnet", subnet)
	route, _ = sjson.Set(route, "mask", mask)
	body, _ := sjson.SetRaw(`{"staticroutes":[]}`, "staticroutes.-1", route)
	return body
}

func createStaticRoute(ctx context.Context, c *cw.Client, plan *StaticRoute) error {
	subnet, mask := plan.Address.ValueString(), strconv.FormatInt(plan.PrefixLength.ValueInt64(), 10)
	res, err := c.Post(staticRoutesPath, routeBody(subnet, mask, ""), cw.UseMutex)
	if err != nil {
		return fmt.Errorf("%s, %s", err, res.String())
	}
	if message := staticRouteAnswer(res); message != "" {
		return fmt.Errorf("the controller refused the route: %s", message)
	}
	deadline := time.Now().Add(staticRouteWait)
	for {
		list, err := c.Post(staticRoutesQueryPath, `{"filter":{}}`, cw.ReadOnly)
		if err == nil {
			if route, ok := findStaticRoute(list, "", subnet, mask); ok {
				switch status := route.Get("status").String(); status {
				case "success":
					plan.Id = types.StringValue(route.Get("uuid").String())
					return nil
				case "add-inprogress", "":
				default:
					return fmt.Errorf("the controller could not add the route (%s): %s", status, route.Get("message").String())
				}
			}
		}
		if !time.Now().Add(staticRoutePoll).Before(deadline) || !sleepContext(ctx, staticRoutePoll) {
			return fmt.Errorf("the route %s/%s was not added in time", subnet, mask)
		}
	}
}

func deleteStaticRoute(ctx context.Context, c *cw.Client, state StaticRoute) error {
	subnet, mask := state.Address.ValueString(), strconv.FormatInt(state.PrefixLength.ValueInt64(), 10)
	res, err := c.DeleteBody(staticRoutesPath, routeBody(subnet, mask, state.Id.ValueString()), cw.UseMutex)
	if err != nil {
		return fmt.Errorf("%s, %s", err, res.String())
	}
	if message := staticRouteAnswer(res); message != "" {
		return fmt.Errorf("the controller refused the removal: %s", message)
	}
	deadline := time.Now().Add(staticRouteWait)
	for {
		list, err := c.Post(staticRoutesQueryPath, `{"filter":{}}`, cw.ReadOnly)
		if err == nil {
			if _, ok := findStaticRoute(list, state.Id.ValueString(), "", ""); !ok {
				return nil
			}
		}
		if !time.Now().Add(staticRoutePoll).Before(deadline) || !sleepContext(ctx, staticRoutePoll) {
			return fmt.Errorf("the route %s/%s was not removed in time", subnet, mask)
		}
	}
}

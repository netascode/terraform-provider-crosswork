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
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cw "github.com/netascode/go-crosswork"
)

// The syslog and the trap receivers are two lists of the alarm service: a GET answers an array of all receivers of a kind, and a
// receiver is addressed by its name (`.../syslog/{name}`, `.../trap/{name}`).

const (
	syslogDestinationPath = "/crosswork/alarm/v1/destination/syslog"
	trapDestinationPath   = "/crosswork/alarm/v1/destination/trap"
)

func findAlarmReceiver(res gjson.Result, name string) (gjson.Result, bool) {
	var found gjson.Result
	ok := false
	res.ForEach(func(_, receiver gjson.Result) bool {
		if receiver.Get("name").String() == name {
			found, ok = receiver, true
			return false
		}
		return true
	})
	return found, ok
}

func readAlarmReceiver(c *cw.Client, path, name string) (cw.Res, error) {
	res, err := c.Get(path)
	if err != nil {
		return res, err
	}
	receiver, ok := findAlarmReceiver(res, name)
	if !ok {
		return res, fmt.Errorf("StatusCode 404: the controller has no receiver %q", name)
	}
	return receiver, nil
}

// A receiver of the DNS type lacks the server address; the host is the name.
func trapReceiverEntry(receiver gjson.Result) gjson.Result {
	if receiver.Get("notifyType").String() == "DNS" && !receiver.Get("serverAddress").Exists() {
		raw, _ := sjson.Set(receiver.Raw, "serverAddress", receiver.Get("name").String())
		return gjson.Parse(raw)
	}
	return receiver
}

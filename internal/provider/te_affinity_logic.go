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
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cw "github.com/netascode/go-crosswork"
)

// The TE link affinities and the Flex-Algo affinities are two lists of the optima service, read and written through the same calls.
// The controller assigns the id of an affinity, has no update (a PUT is a 405), and answers a refusal with HTTP 200 and a status.

const affinityPath = "/crosswork/optima/api/v1/affinityMap"

type affinityKind struct {
	list string
	kind string
}

var (
	teLinkAffinity     = affinityKind{"teLink", "REGULAR"}
	teFlexAlgoAffinity = affinityKind{"flexAlgo", "FLEX_ALGO"}
)

func (k affinityKind) listPath() string { return affinityPath + "/" + k.list }

func affinityRefusal(res gjson.Result) string {
	if res.Get("data.status").String() != "failure" {
		return ""
	}
	errors := res.Get("data.errors")
	if errors.IsArray() {
		var messages []string
		errors.ForEach(func(_, e gjson.Result) bool {
			messages = append(messages, strings.TrimSpace(e.Get("error").String()))
			return true
		})
		return strings.Join(messages, "; ")
	}
	if errors.Exists() {
		return errors.String()
	}
	return res.Raw
}

// Falls back to the name when the id is not known.
func findAffinity(res gjson.Result, id, name string) (gjson.Result, bool) {
	var found gjson.Result
	ok := false
	res.Get("data.affinityMappings").ForEach(func(_, entry gjson.Result) bool {
		if id != "" && entry.Get("id").String() == id || id == "" && entry.Get("name").String() == name {
			found, ok = entry, true
			return false
		}
		return true
	})
	return found, ok
}

func affinityCreateBody(k affinityKind, name string, bit int64) string {
	item, _ := sjson.Set("", "name", name)
	item, _ = sjson.Set(item, "bitPosition", bit)
	item, _ = sjson.Set(item, "type", k.kind)
	body, _ := sjson.SetRaw("[]", "-1", item)
	return body
}

func readAffinity(c *cw.Client, k affinityKind, id, name string) (cw.Res, error) {
	res, err := c.Get(k.listPath())
	if err != nil {
		return res, err
	}
	entry, ok := findAffinity(res, id, name)
	if !ok {
		return res, fmt.Errorf("StatusCode 404: the controller has no affinity %q", name)
	}
	return entry, nil
}

func createAffinity(c *cw.Client, k affinityKind, name string, bit int64) (string, error) {
	res, err := c.Post(affinityPath, affinityCreateBody(k, name, bit), cw.UseMutex)
	if err != nil {
		return "", fmt.Errorf("%s, %s", err, res.String())
	}
	if message := affinityRefusal(res); message != "" {
		return "", fmt.Errorf("the controller refused the affinity: %s", message)
	}
	entry, err := readAffinity(c, k, "", name)
	if err != nil {
		return "", err
	}
	return strconv.FormatInt(entry.Get("id").Int(), 10), nil
}

func deleteAffinity(c *cw.Client, k affinityKind, id string) error {
	number, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return fmt.Errorf("the id %q of the affinity is not a number", id)
	}
	res, err := c.DeleteBody(k.listPath(), "["+strconv.FormatInt(number, 10)+"]", cw.UseMutex)
	if err != nil {
		return fmt.Errorf("%s, %s", err, res.String())
	}
	if message := affinityRefusal(res); message != "" {
		return fmt.Errorf("the controller refused the removal: %s", message)
	}
	return nil
}

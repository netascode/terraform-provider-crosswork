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

const affinitiesAnswer = `{"data":{"count":2,"affinityMappings":[{"id":2,"name":"gold","bitPosition":5},{"id":1,"name":"silver","bitPosition":128}]},"attributes":{}}`

func TestAffinityCreateBodyIsAListWithTheKind(t *testing.T) {
	link := gjson.Parse(affinityCreateBody(teLinkAffinity, "gold", 5))
	assert.Equal(t, "gold", link.Get("0.name").String())
	assert.Equal(t, int64(5), link.Get("0.bitPosition").Int())
	assert.Equal(t, "REGULAR", link.Get("0.type").String())
	assert.Equal(t, "FLEX_ALGO", gjson.Parse(affinityCreateBody(teFlexAlgoAffinity, "x", 200)).Get("0.type").String())
	assert.Equal(t, "/crosswork/optima/api/v1/affinityMap/flexAlgo", teFlexAlgoAffinity.listPath())
}

func TestAffinityIsFoundByIdOrByName(t *testing.T) {
	res := gjson.Parse(affinitiesAnswer)
	entry, ok := findAffinity(res, "1", "")
	assert.True(t, ok)
	assert.Equal(t, "silver", entry.Get("name").String())
	entry, ok = findAffinity(res, "", "gold")
	assert.True(t, ok)
	assert.Equal(t, int64(2), entry.Get("id").Int())
	_, ok = findAffinity(res, "9", "")
	assert.False(t, ok)
}

func TestAffinityRefusalIsAStatusInTheBody(t *testing.T) {
	assert.Equal(t, "", affinityRefusal(gjson.Parse(`{"data":{"status":"success"},"attributes":{}}`)))
	assert.Contains(t, affinityRefusal(gjson.Parse(`{"data":{"errors":[{"id":0,"error":"DETAIL:  Key (\"Name\")=(gold) already exists.\n"}],"status":"failure"}}`)), "already exists")
	assert.Contains(t, affinityRefusal(gjson.Parse(`{"data":{"errors":"Out of range uint32 value: -1","status":"failure"}}`)), "Out of range")
}

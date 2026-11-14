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
	"testing"

	"github.com/tidwall/gjson"
)

const discoveryReadyAnswer = `{"providerConfig":{"cdpConfig":{"enabled":true,"status":"Ready","allowEdit":true},"lldpConfig":{"enabled":false,"status":"Ready","allowEdit":true},"lagConfig":{"enabled":true,"status":"Ready","allowEdit":true}}}`
const discoveryBusyAnswer = `{"providerConfig":{"cdpConfig":{"enabled":true,"status":"DisableInProgress","allowEdit":false},"lldpConfig":{"enabled":false,"status":"Ready","allowEdit":true},"lagConfig":{"enabled":true,"status":"Ready","allowEdit":true}}}`

func TestDiscoveryReadyNeedsEveryProtocolReady(t *testing.T) {
	assert.True(t, discoveryReady(gjson.Parse(discoveryReadyAnswer)))
	assert.False(t, discoveryReady(gjson.Parse(discoveryBusyAnswer)), "a protocol that is still changing is not ready")
	assert.False(t, discoveryReady(gjson.Parse(`{}`)))
}

func TestDiscoveryMatchesTheWantedFlags(t *testing.T) {
	res := gjson.Parse(discoveryReadyAnswer)
	assert.True(t, discoveryMatches(res, map[string]bool{"lldpConfig": false, "cdpConfig": true, "lagConfig": true}))
	assert.False(t, discoveryMatches(res, map[string]bool{"lldpConfig": true}))
}

func TestDiscoveryAnswerThatWasAppliedIsNotARefusal(t *testing.T) {
	assert.Equal(t, "", discoveryRefusal(gjson.Parse(`{"error":"PROVIDER_NOT_AVAILABLE"}`)))
	assert.Equal(t, "", discoveryRefusal(gjson.Parse(`{}`)))
	assert.Equal(t, "SOMETHING_ELSE", discoveryRefusal(gjson.Parse(`{"error":"SOMETHING_ELSE"}`)))
}

func TestDiscoveryWriteFailureOfAnAppliedChangeIsNotRefused(t *testing.T) {
	assert.False(t, discoveryWriteRefused(nil))
	assert.False(t, discoveryWriteRefused(fmt.Errorf("HTTP Request failed: StatusCode 200: PROVIDER_NOT_AVAILABLE")), "the state that is read back decides")
	assert.True(t, discoveryWriteRefused(fmt.Errorf("HTTP Request failed: StatusCode 500: boom")))
}

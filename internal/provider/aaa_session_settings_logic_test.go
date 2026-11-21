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
	"testing"

	"github.com/tidwall/gjson"
)

func TestAAASessionPutBodyLeavesTheSourceIPOut(t *testing.T) {
	body := sessionPutBody(`{"FallbackType":"NO_FALLBACK","IdleSessionTimeout":30,"NumParallelSessions":100,"NumParallelSessionsPerUser":100,"IsEnabled":true}`)
	assert.False(t, gjson.Get(body, "IsEnabled").Exists(), "the session call refuses members it does not know")
	assert.Equal(t, int64(100), gjson.Get(body, "NumParallelSessions").Int())
}

func TestAAASourceIPSetBody(t *testing.T) {
	assert.Equal(t, `{"IsEnabled":true,"ScaleUp":false}`, sourceIPSetBody(true))
	assert.Equal(t, `{"IsEnabled":false,"ScaleUp":false}`, sourceIPSetBody(false))
}

func TestAAASourceIPFromAnswer(t *testing.T) {
	on, err := sourceIPFromAnswer(gjson.Parse(`{"Enabled":true,"Error":""}`))
	assert.True(t, err == nil, err)
	assert.True(t, on)
	off, err := sourceIPFromAnswer(gjson.Parse(`{"Enabled":false,"Error":""}`))
	assert.True(t, err == nil, err)
	assert.False(t, off)
	_, err = sourceIPFromAnswer(gjson.Parse(`{"Enabled":false,"Error":"cluster busy"}`))
	assert.True(t, err != nil)
	assert.Contains(t, err.Error(), "cluster busy")
	_, err = sourceIPFromAnswer(gjson.Parse(`{}`))
	assert.True(t, err != nil)
	assert.Contains(t, err.Error(), "no Enabled member")
}

func TestAAASourceIPSetResult(t *testing.T) {
	err := sourceIPSetResult(gjson.Parse(`{"resp_value":"R_SUCCESS","resp_error":"","description":""}`))
	assert.True(t, err == nil, err)
	err = sourceIPSetResult(gjson.Parse(`{"resp_value":"R_FAILURE","resp_error":"no quorum"}`))
	assert.True(t, err != nil)
	assert.Contains(t, err.Error(), "no quorum")
}

func TestAAAWithSourceIPAddsTheMemberTheModelReads(t *testing.T) {
	res := withSourceIP(gjson.Parse(`{"IdleSessionTimeout":30,"enableDAG":true}`), true)
	assert.True(t, res.Get("IsEnabled").Bool())
	assert.Equal(t, int64(30), res.Get("IdleSessionTimeout").Int(), "the session members stay")
}

func TestAAASessionSettingsModelReadsTheSourceIPFromIsEnabled(t *testing.T) {
	var data AAASessionSettings
	data.fromBody(context.Background(), withSourceIP(gjson.Parse(`{"FallbackType":"FB_ON_NO_RESPONSE","IdleSessionTimeout":60,"NumParallelSessions":100,"NumParallelSessionsPerUser":50}`), true))
	assert.True(t, data.SourceIpAuditing.ValueBool())
	assert.Equal(t, "FB_ON_NO_RESPONSE", data.FallbackType.ValueString())
	assert.Equal(t, int64(50), data.ParallelSessionsPerUser.ValueInt64())
}

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

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

func TestAAARADIUSReadOfAnEmptyAnswerIsAnEmptyList(t *testing.T) {
	var data AAARADIUSServers
	data.fromBody(context.Background(), gjson.Parse(`{}`))
	assert.Equal(t, 0, len(data.Servers), "an empty RADIUS configuration has no list member")
	assert.Equal(t, "aaa_radius_servers", data.Id.ValueString())
}

func TestAAARADIUSBody(t *testing.T) {
	data := AAARADIUSServers{Servers: []AAARADIUSServersServers{{
		Host: types.StringValue("192.0.2.20"), Port: types.Int64Value(1812), Secret: types.StringValue("AAASecret"), SecretWo: types.StringNull(),
		SecretFormat: types.StringValue("ascii"), RetransmitTimeout: types.Int64Value(5), Retries: types.Int64Value(3), AuthType: types.StringValue("chap"),
		Priority: types.Int64Value(40), PolicyId: types.StringValue("role0"), VirtualDomain: types.StringValue("deviceAccessGroups"),
	}}}
	e := gjson.Parse(data.toBody(context.Background(), AAARADIUSServers{})).Get("radius_servers.0")
	assert.Equal(t, int64(1812), e.Get("port").Int())
	assert.Equal(t, "chap", e.Get("authType").String())
	assert.Equal(t, "role0", e.Get("policyId").String(), "RADIUS names it policyId, TACACS policy-id")
	assert.Equal(t, int64(40), e.Get("priority").Int())
	assert.Equal(t, "deviceAccessGroups", e.Get("virtualDomain").String(), "the device access group attribute")
}

func TestAAARADIUSServersEmptyListIsSentSoThatADestroyOrAnUpdateCanClearIt(t *testing.T) {
	assert.Equal(t, `{"radius_servers":[]}`, AAARADIUSServers{}.toBody(context.Background(), AAARADIUSServers{}))
}

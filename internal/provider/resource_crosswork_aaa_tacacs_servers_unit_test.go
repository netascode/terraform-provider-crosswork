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

func tacacsEntry() AAATACACSServersServers {
	return AAATACACSServersServers{
		Host: types.StringValue("192.0.2.10"), Port: types.Int64Value(49), Secret: types.StringValue("AAASecret"), SecretWo: types.StringNull(),
		SecretFormat: types.StringValue("ascii"), Service: types.StringValue("crosswork"), Priority: types.Int64Value(30),
		RetransmitTimeout: types.Int64Value(10), Retries: types.Int64Value(3), AuthType: types.StringValue("pap"),
		Interface: types.StringNull(), PolicyId: types.StringValue("role0"), VirtualDomain: types.StringValue("deviceAccessGroups"),
	}
}

func TestAAATACACSBodyUsesTheNamesOfTheController(t *testing.T) {
	body := gjson.Parse(AAATACACSServers{Servers: []AAATACACSServersServers{tacacsEntry()}}.toBody(context.Background(), AAATACACSServers{}))
	e := body.Get("tacacs_servers.0")
	assert.Equal(t, "192.0.2.10", e.Get("host").String())
	assert.Equal(t, "AAASecret", e.Get("secret").String())
	assert.Equal(t, "ascii", e.Get("secretFormat").String())
	assert.Equal(t, "pap", e.Get("authType").String())
	assert.Equal(t, "role0", e.Get("policy-id").String(), "the name has a dash")
	assert.Equal(t, "deviceAccessGroups", e.Get("virtualDomain").String())
	assert.False(t, e.Get("interface").Exists(), "an attribute that is not set is not sent")
}

func TestAAATACACSSecretOfTheWriteOnlyFormIsSentInTheSecretMember(t *testing.T) {
	e := tacacsEntry()
	e.Secret, e.SecretWo = types.StringNull(), types.StringValue("AAASecret")
	body := gjson.Parse(AAATACACSServers{Servers: []AAATACACSServersServers{e}}.toBody(context.Background(), AAATACACSServers{}))
	assert.Equal(t, "AAASecret", body.Get("tacacs_servers.0.secret").String())
}

func TestAAATACACSReadKeepsTheSecretOfTheStateWhenTheControllerMasksIt(t *testing.T) {
	data := AAATACACSServers{Id: types.StringValue("aaa_tacacs_servers"), Servers: []AAATACACSServersServers{tacacsEntry()}}
	data.updateFromBody(context.Background(), gjson.Parse(`{"tacacs_servers":[{"host":"192.0.2.10","port":49,"secret":"******","secretFormat":"ascii","service":"crosswork","priority":30,"timeout":10,"retries":3,"authType":"pap","policy-id":"role0","virtualDomain":"deviceAccessGroups"}]}`))
	assert.Equal(t, "AAASecret", data.Servers[0].Secret.ValueString(), "the masked value is never stored")
	assert.Equal(t, int64(30), data.Servers[0].Priority.ValueInt64())
}

func TestAAATACACSReadOfAServerThatIsGoneClearsItsAttributes(t *testing.T) {
	data := AAATACACSServers{Id: types.StringValue("aaa_tacacs_servers"), Servers: []AAATACACSServersServers{tacacsEntry()}}
	data.updateFromBody(context.Background(), gjson.Parse(`{"tacacs_servers":[]}`))
	assert.True(t, data.Servers[0].Priority.IsNull(), "the plan then shows the server as a change")
}

func TestAAATACACSImportReadsTheServersInTheOrderOfTheController(t *testing.T) {
	var data AAATACACSServers
	data.fromBody(context.Background(), gjson.Parse(`{"tacacs_servers":[{"host":"192.0.2.10","port":49,"priority":10},{"host":"192.0.2.11","port":49,"priority":20}]}`))
	assert.Equal(t, 2, len(data.Servers))
	assert.Equal(t, "192.0.2.11", data.Servers[1].Host.ValueString())
	assert.Equal(t, "aaa_tacacs_servers", data.Id.ValueString())
}

func TestAAATACACSServersEmptyListIsSentSoThatADestroyOrAnUpdateCanClearIt(t *testing.T) {
	assert.Equal(t, `{"tacacs_servers":[]}`, AAATACACSServers{}.toBody(context.Background(), AAATACACSServers{}))
}

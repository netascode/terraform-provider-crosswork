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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/tidwall/gjson"
)

const syslogAnswer = `[{"serverAddress":"9.0.1.88","transportType":"UDP","name":"noc","port":514,"certificateId":""},{"serverAddress":"192.0.2.1","transportType":"TLS","name":"tls","port":6514,"certificateId":"c1"}]`

func TestSyslogReceiverIsFoundByName(t *testing.T) {
	receiver, ok := findAlarmReceiver(gjson.Parse(syslogAnswer), "tls")
	assert.True(t, ok)
	assert.Equal(t, "c1", receiver.Get("certificateId").String())
	_, ok = findAlarmReceiver(gjson.Parse(syslogAnswer), "gone")
	assert.False(t, ok)
}

func TestSyslogAddressIsNeverReadBack(t *testing.T) {
	receiver, _ := findAlarmReceiver(gjson.Parse(syslogAnswer), "noc")
	state := SyslogDestination{Name: types.StringValue("noc"), Address: types.StringValue("as-linux1"), Port: types.Int64Value(1), Protocol: types.StringValue("TLS"), CertificateId: types.StringValue("x")}
	state.updateFromBody(context.Background(), receiver)
	assert.Equal(t, "as-linux1", state.Address.ValueString(), "the host name of the configuration stays, whatever IP the controller keeps")
	assert.Equal(t, int64(514), state.Port.ValueInt64(), "the other values are read")
	assert.Equal(t, "UDP", state.Protocol.ValueString())
	assert.True(t, state.CertificateId.IsNull(), "an empty certificate id is read as none")
	var imported SyslogDestination
	imported.Name = types.StringValue("noc")
	imported.updateFromBody(context.Background(), receiver)
	assert.Equal(t, "9.0.1.88", imported.Address.ValueString(), "a state without an address takes the one of the controller")
}

// The addresses are from the documentation range 192.0.2.0/24 (RFC 5737): no DNS or reachable receiver is needed.
func TestAccCwSyslogDestinationAddress(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	first, second := "192.0.2.1", "192.0.2.2"
	config := func(host string) string {
		return `resource "crosswork_syslog_destination" "test" {` + "\n" + `	name = "cwtest-syslog-host"` + "\n" + `	address = "` + host + `"` + "\n" + `}` + "\n"
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config(first), Check: resource.TestCheckResourceAttr("crosswork_syslog_destination.test", "address", first)},
			{Config: config(second), Check: resource.TestCheckResourceAttr("crosswork_syslog_destination.test", "address", second)},
		},
	})
}

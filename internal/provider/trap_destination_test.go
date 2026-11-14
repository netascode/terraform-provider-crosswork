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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/tidwall/gjson"
)

func TestTrapReceiverOfADnsTypeHasItsHostAsAddress(t *testing.T) {
	made := gjson.Parse(`{"name":"as-linux1","notifyType":"DNS","port":162}`)
	assert.Equal(t, "as-linux1", trapReceiverEntry(made).Get("serverAddress").String(), "the host of a receiver that the page made is its name")
	api := gjson.Parse(`{"name":"x","notifyType":"DNS","serverAddress":"192.0.2.2"}`)
	assert.Equal(t, "192.0.2.2", trapReceiverEntry(api).Get("serverAddress").String(), "a receiver with a server address is left as it is")
	ip := gjson.Parse(`{"name":"y","notifyType":"IPADDRESS","serverAddress":"192.0.2.30"}`)
	assert.Equal(t, "192.0.2.30", trapReceiverEntry(ip).Get("serverAddress").String())
}

func TestAccCwTrapDestinationV3(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	config := `resource "crosswork_trap_destination" "test" {
	name = "cwtest-trap-v3"
	address = "192.0.2.41"
	snmp_version = "VERSION_3"
	username = "cwtest"
	mode = "AuthPriv"
	auth_type = "HMACSHA"
	auth_password = "cwtest-auth-1234"
	privacy_type = "AES128"
	privacy_password = "cwtest-priv-1234"
}
`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("crosswork_trap_destination.test", "mode", "AuthPriv"),
				resource.TestCheckResourceAttr("crosswork_trap_destination.test", "privacy_type", "AES128"))},
		},
	})
}

func TestAccCwTrapDestinationHostName(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	config := func(host, addressType string) string {
		return `resource "crosswork_trap_destination" "test" {
	name = "` + host + `"
	address = "` + host + `"
	address_type = "` + addressType + `"
	community = "cwtest-community"
}
`
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config("as-linux1", "DNS"), Check: resource.TestCheckResourceAttr("crosswork_trap_destination.test", "address", "as-linux1")},
			{Config: config("192.0.2.2", "IPADDRESS"), Check: resource.TestCheckResourceAttr("crosswork_trap_destination.test", "address", "192.0.2.2")},
		},
	})
}

// The Update deletes and creates the receiver because the controller's update does not work.
func TestAccCwTrapDestinationUpdateAndRotation(t *testing.T) {
	if os.Getenv("SETTINGS") == "" {
		t.Skip("skipping test, set environment variable SETTINGS")
	}
	config := func(port, version, password string) string {
		return `resource "crosswork_trap_destination" "test" {
	name = "cwtest-trap-rotation"
	address = "192.0.2.42"
	port = ` + port + `
	snmp_version = "VERSION_3"
	username = "cwtest"
	mode = "AuthNoPriv"
	auth_type = "HMACSHA"
	auth_password_wo = "` + password + `"
	auth_password_wo_version = ` + version + `
	privacy_type = "NONE"
}
`
	}
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config("162", "1", "cwtest-auth-1111"), Check: resource.TestCheckResourceAttr("crosswork_trap_destination.test", "port", "162")},
			{Config: config("162", "2", "cwtest-auth-2222"), Check: resource.TestCheckResourceAttr("crosswork_trap_destination.test", "auth_password_wo_version", "2")},
			{Config: config("163", "2", "cwtest-auth-2222"), Check: resource.TestCheckResourceAttr("crosswork_trap_destination.test", "port", "163")},
		},
	})
}

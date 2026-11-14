// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://mozilla.org/MPL/2.0/
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
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	updateDevice = "cwtest-gen-device-update"
	tagFirst     = "cwtest-gen-device-tag"
	tagSecond    = "cwtest-gen-device-tag-2"
)

// Reuses the provider's client so no second session is opened.
func testDeviceClient(t *testing.T) *cw.Client {
	insecure, retries := true, 3
	if v, err := strconv.ParseBool(os.Getenv("CW_INSECURE")); err == nil {
		insecure = v
	}
	if v, err := strconv.Atoi(os.Getenv("CW_RETRIES")); err == nil {
		retries = v
	}
	c, err := registeredClient(os.Getenv("CW_URL"), os.Getenv("CW_USERNAME"), os.Getenv("CW_PASSWORD"), insecure, retries, "test")
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func setDeviceTags(t *testing.T, names ...string) func() {
	return func() {
		c := testDeviceClient(t)
		res, err := c.Post(nodesQueryPath, fmt.Sprintf(`{"filter":{"inventory_id":%q}}`, updateDevice), cw.ReadOnly)
		if err != nil || !res.Get("data.0").Exists() {
			t.Fatalf("reading the device: %v", err)
		}
		tags := []string{}
		for _, n := range names {
			tags = append(tags, fmt.Sprintf(`{"name":%q}`, n))
		}
		node, _ := sjson.SetRaw(cleanNode(res.Get("data.0").Raw), "tags", "["+strings.Join(tags, ",")+"]")
		if _, err := c.Put(nodesPath, nodeBody(node)); err != nil {
			t.Fatalf("changing the tags: %v", err)
		}
	}
}

func checkDeviceTag(t *testing.T, name string, want bool) resource.TestCheckFunc {
	return func(*terraform.State) error {
		res, err := testDeviceClient(t).Post(nodesQueryPath, fmt.Sprintf(`{"filter":{"inventory_id":%q}}`, updateDevice), cw.ReadOnly)
		if err != nil {
			return err
		}
		has := false
		for _, n := range gjson.Get(res.Raw, "data.0.tag_names").Array() {
			has = has || n.String() == name
		}
		if has != want {
			return fmt.Errorf("tag %s on the device: got %v, want %v", name, has, want)
		}
		return nil
	}
}

const testAccCwNetworkDevicesUpdateTagsConfig = `
resource "crosswork_tags" "update_second" {
  name     = "cwtest-gen-device-tag-2"
  category = "cwtest-category"
}
`

func testAccCwNetworkDevicesUpdateConfig(hostName, serial, tags, provider string) string {
	config := `resource "crosswork_network_devices" "test" {` + "\n"
	config += `	name = "cwtest-gen-device-update"` + "\n"
	config += fmt.Sprintf(`	host_name = %q`+"\n", hostName)
	config += fmt.Sprintf(`	serial_number = %q`+"\n", serial)
	config += `	credential_profile = crosswork_credential_profiles.device_test.id` + "\n"
	config += `	software_type = "IOS XR"` + "\n"
	config += `	capabilities = ["SNMP", "YANG_CLI"]` + "\n"
	config += `	connectivity = [` + "\n"
	config += `	  { protocol = "ROBOT_MSVC_TRANS_SSH", ip_address = "192.0.2.232", prefix_length = 24, port = 22, timeout = 60 },` + "\n"
	config += `	  { protocol = "ROBOT_MSVC_TRANS_SNMP", ip_address = "192.0.2.232", prefix_length = 24, port = 161, timeout = 60 },` + "\n"
	config += `	]` + "\n"
	if provider != "" {
		config += fmt.Sprintf(`	provider_profile = %q`+"\n", provider)
	}
	if tags != "" {
		config += `	tags = [` + tags + `]` + "\n"
	}
	// Crosswork can move the admin state itself (an embedded data gateway pool brings a new device up); ignoring it keeps the update on an admin up device.
	config += "\n\tlifecycle {\n\t  ignore_changes = [admin_state]\n\t}\n"
	config += `}` + "\n"
	return config
}

func TestAccCwNetworkDevicesUpdate(t *testing.T) {
	const res = "crosswork_network_devices.test"
	prereq := testAccCwNetworkDevicesPrerequisitesConfig + testAccCwNetworkDevicesUpdateTagsConfig
	first := `{ name = crosswork_tags.device_test.id }`
	second := `{ name = crosswork_tags.update_second.id }`
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: prereq + testAccCwNetworkDevicesUpdateConfig("cwtest-update-a", "SN-A", first, "nso"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "host_name", "cwtest-update-a"),
					resource.TestCheckResourceAttr(res, "serial_number", "SN-A"),
					resource.TestCheckResourceAttr(res, "tags.#", "1"),
					resource.TestCheckResourceAttr(res, "provider_profile", "nso"),
				),
			},
			{
				// A field the controller accepts while the device is admin up: one plain PUT.
				Config: prereq + testAccCwNetworkDevicesUpdateConfig("cwtest-update-a", "SN-B", first, "nso"),
				Check:  resource.TestCheckResourceAttr(res, "serial_number", "SN-B"),
			},
			{
				// A field the controller refuses while the device is admin up: down, change, up.
				Config: prereq + testAccCwNetworkDevicesUpdateConfig("cwtest-update-b", "SN-B", first, "nso"),
				Check:  resource.TestCheckResourceAttr(res, "host_name", "cwtest-update-b"),
			},
			{
				Config: prereq + testAccCwNetworkDevicesUpdateConfig("cwtest-update-b", "SN-B", second, "nso"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "tags.#", "1"),
					resource.TestCheckResourceAttrPair(res, "tags.0.name", "crosswork_tags.update_second", "id"),
				),
			},
			{
				PreConfig:          setDeviceTags(t, tagSecond, tagFirst),
				Config:             prereq + testAccCwNetworkDevicesUpdateConfig("cwtest-update-b", "SN-B", second, "nso"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: prereq + testAccCwNetworkDevicesUpdateConfig("cwtest-update-b", "SN-B", second, "nso"),
				Check:  resource.ComposeTestCheckFunc(checkDeviceTag(t, tagSecond, true), checkDeviceTag(t, tagFirst, false)),
			},
			{
				PreConfig:          setDeviceTags(t),
				Config:             prereq + testAccCwNetworkDevicesUpdateConfig("cwtest-update-b", "SN-B", second, "nso"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
			{
				Config: prereq + testAccCwNetworkDevicesUpdateConfig("cwtest-update-b", "SN-B", second, "nso"),
				Check:  checkDeviceTag(t, tagSecond, true),
			},
			{
				Config: prereq + testAccCwNetworkDevicesUpdateConfig("cwtest-update-c", "SN-C", "", ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "host_name", "cwtest-update-c"),
					resource.TestCheckResourceAttr(res, "serial_number", "SN-C"),
					resource.TestCheckResourceAttr(res, "tags.#", "0"),
					resource.TestCheckNoResourceAttr(res, "provider_profile"),
				),
			},
			{
				// With tags unset, a tag added on the controller is not a difference ...
				PreConfig: setDeviceTags(t, tagFirst),
				Config:    prereq + testAccCwNetworkDevicesUpdateConfig("cwtest-update-c", "SN-C", "", ""),
				PlanOnly:  true,
			},
			{
				// ... and an update of another field keeps it.
				Config: prereq + testAccCwNetworkDevicesUpdateConfig("cwtest-update-c", "SN-D", "", ""),
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr(res, "serial_number", "SN-D"), checkDeviceTag(t, tagFirst, true)),
			},
			{
				// Remove the added tag so the destroy can delete the test's tags.
				PreConfig: setDeviceTags(t),
				Config:    prereq + testAccCwNetworkDevicesUpdateConfig("cwtest-update-c", "SN-D", "", ""),
				PlanOnly:  true,
			},
		},
	})
}

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
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"crosswork": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CW_USERNAME"); v == "" {
		t.Fatal("CW_USERNAME env variable must be set for acceptance tests")
	}
	if v := os.Getenv("CW_PASSWORD"); v == "" {
		t.Fatal("CW_PASSWORD env variable must be set for acceptance tests")
	}
	if v := os.Getenv("CW_URL"); v == "" {
		t.Fatal("CW_URL env variable must be set for acceptance tests")
	}
}

// The test process shares one session per cluster; release them on exit.
func TestMain(m *testing.M) {
	code := m.Run()
	Shutdown()
	os.Exit(code)
}

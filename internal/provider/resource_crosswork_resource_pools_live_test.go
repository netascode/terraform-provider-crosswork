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

// A pool that services use holds allocations, and an update of the range must not delete them. A PUT of the whole pool does
// delete them; the update of the resource is a PUT of the range only.

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	cw "github.com/netascode/go-crosswork"
)

const poolsAllocID = "cwtest-alloc-1"

func poolsClientPath(name string) string {
	return ResourcePools{Name: types.StringValue(name)}.getPath()
}

func poolsTestClient(t *testing.T) *cw.Client {
	t.Helper()
	insecure, _ := strconv.ParseBool(os.Getenv("CW_INSECURE"))
	c, err := registeredClient(os.Getenv("CW_URL"), os.Getenv("CW_USERNAME"), os.Getenv("CW_PASSWORD"), insecure, 3, "test")
	if err != nil {
		t.Fatalf("login: %s", err)
	}
	return c
}

func poolsConfig(name string, start, end int) string {
	return fmt.Sprintf("resource \"crosswork_resource_pools\" \"test\" {\n  name        = %q\n  range_start = %d\n  range_end   = %d\n}\n", name, start, end)
}

func poolsCreate(t *testing.T, name string, start, end int) func() {
	return func() {
		c := poolsTestClient(t)
		body := fmt.Sprintf(`{"id-allocator:id-pool":[{"name":%q,"range":{"start":%d,"end":%d}}]}`, name, start, end)
		if res, err := c.Post("/crosswork/proxy/nso/restconf/data/resource-allocator:resource-pools", body); err != nil {
			t.Fatalf("pool: %s %s", err, res.String())
		}
		poolsAllocate(t, name)()
	}
}

func poolsAllocate(t *testing.T, name string) func() {
	return func() {
		body := fmt.Sprintf(`{"id-allocator:allocation":[{"id":%q,"username":"cwtest","request":{"sync":false}}]}`, poolsAllocID)
		if res, err := poolsTestClient(t).Post(poolsClientPath(name), body); err != nil {
			t.Fatalf("allocation: %s %s", err, res.String())
		}
	}
}

func poolsHolds(t *testing.T, name string, end int) resource.TestCheckFunc {
	return func(*terraform.State) error {
		res, err := poolsTestClient(t).Get(poolsClientPath(name))
		if err != nil {
			return err
		}
		if got := res.Get("id-allocator:id-pool.0.allocation.#").Int(); got != 1 {
			return fmt.Errorf("the pool holds %d allocations, want 1: %s", got, res.String())
		}
		if got := res.Get("id-allocator:id-pool.0.allocation.0.response.id").Int(); got < 1 {
			return fmt.Errorf("the allocation has no id: %s", res.String())
		}
		if got := res.Get("id-allocator:id-pool.0.range.end").Int(); got != int64(end) {
			return fmt.Errorf("range end is %d, want %d", got, end)
		}
		return nil
	}
}

func TestAccCwResourcePoolsKeepAllocations(t *testing.T) {
	const res = "crosswork_resource_pools.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: poolsConfig("cwtest-pool-live", 1, 100), Check: resource.TestCheckResourceAttr(res, "range_end", "100")},
			{PreConfig: poolsAllocate(t, "cwtest-pool-live"), Config: poolsConfig("cwtest-pool-live", 1, 200), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(res, "range_end", "200"), poolsHolds(t, "cwtest-pool-live", 200))},
			{Config: poolsConfig("cwtest-pool-live", 5, 300), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(res, "range_start", "5"), poolsHolds(t, "cwtest-pool-live", 300))},
			// a new name: the old pool (and its allocation) is deleted and a new one created
			{Config: poolsConfig("cwtest-pool-live2", 5, 300), Check: resource.TestCheckResourceAttr(res, "name", "cwtest-pool-live2")},
			{ResourceName: res, ImportState: true, ImportStateId: "cwtest-pool-live2", ImportStateVerify: true},
		},
	})
}

// A pool that exists before Terraform (made by the install or by an administrator, and in use by a service) is taken over when the
// provider setting allow_existing_on_create is on: its range becomes the configured one and its allocations stay.
func TestAccCwResourcePoolsAdoptExisting(t *testing.T) {
	t.Setenv("CW_ALLOW_EXISTING_ON_CREATE", "true")
	const res = "crosswork_resource_pools.test"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{PreConfig: poolsCreate(t, "cwtest-pool-adopt", 1, 50), Config: poolsConfig("cwtest-pool-adopt", 1, 100), Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr(res, "range_end", "100"), poolsHolds(t, "cwtest-pool-adopt", 100))},
		},
	})
}

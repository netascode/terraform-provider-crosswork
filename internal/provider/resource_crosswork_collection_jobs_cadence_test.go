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

// The cadence is cluster-wide: the test restores the original values when it ends, also on failure.

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func otherCadence(v int64) int64 {
	if v+600 <= 86400 {
		return v + 600
	}
	return v - 600
}

func restoreCadence(t *testing.T) cadenceValues {
	t.Helper()
	c := testDeviceClient(t)
	original, err := readCadence(c)
	if err != nil {
		t.Fatalf("reading the cadence before the test: %v", err)
	}
	body, err := cadenceBody(original, original)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := putCadence(c, body); err != nil {
			t.Errorf("RESTORE THE CADENCE BY HAND, %s: %v", cadenceBodyText(original), err)
			return
		}
		if got, err := readCadence(c); err != nil || got != original {
			t.Errorf("RESTORE THE CADENCE BY HAND, %s: the cluster shows %s (%v)", cadenceBodyText(original), cadenceBodyText(got), err)
		}
	})
	return original
}

func TestAccCwCollectionJobsCadence(t *testing.T) {
	const res, ds = "crosswork_collection_jobs_cadence.test", "data.crosswork_collection_jobs_cadence.test"
	if os.Getenv("TF_ACC") == "" {
		t.Skip("skipping acceptance test, set TF_ACC")
	}
	testAccPreCheck(t)
	original := restoreCadence(t)
	reach, snmp := otherCadence(original.ReachabilityCheck.ValueInt64()), otherCadence(original.SnmpCheck.ValueInt64())
	tunnel, clock := otherCadence(original.TeTunnelId.ValueInt64()), otherCadence(original.ShowClock.ValueInt64())
	s := func(v int64) string { return strconv.FormatInt(v, 10) }

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// the two values that are not set stay as the cluster holds them
				Config: fmt.Sprintf("resource \"crosswork_collection_jobs_cadence\" \"test\" {\n  reachability_check = %d\n  show_clock = %d\n}\n", reach, clock),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "id", "collection_jobs_cadence"),
					resource.TestCheckResourceAttr(res, "reachability_check", s(reach)),
					resource.TestCheckResourceAttr(res, "show_clock", s(clock)),
					resource.TestCheckNoResourceAttr(res, "snmp_check"),
					resource.TestCheckNoResourceAttr(res, "te_tunnel_id"),
					func(*terraform.State) error {
						got, err := readCadence(testDeviceClientNoT())
						if err != nil {
							return err
						}
						if got.SnmpCheck != original.SnmpCheck || got.TeTunnelId != original.TeTunnelId {
							return fmt.Errorf("a value the configuration leaves out changed: %s", cadenceBodyText(got))
						}
						return nil
					},
				),
			},
			{
				Config: fmt.Sprintf(`
resource "crosswork_collection_jobs_cadence" "test" {
  reachability_check = %d
  snmp_check         = %d
  te_tunnel_id       = %d
  show_clock         = %d
}

data "crosswork_collection_jobs_cadence" "test" {
  depends_on = [crosswork_collection_jobs_cadence.test]
}
`, reach, snmp, tunnel, clock),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(res, "snmp_check", s(snmp)),
					resource.TestCheckResourceAttr(res, "te_tunnel_id", s(tunnel)),
					resource.TestCheckResourceAttr(ds, "id", "collection_jobs_cadence"),
					resource.TestCheckResourceAttr(ds, "reachability_check", s(reach)),
					resource.TestCheckResourceAttr(ds, "snmp_check", s(snmp)),
					resource.TestCheckResourceAttr(ds, "te_tunnel_id", s(tunnel)),
					resource.TestCheckResourceAttr(ds, "show_clock", s(clock)),
				),
			},
			{
				ResourceName:      res,
				ImportState:       true,
				ImportStateId:     "collection_jobs_cadence",
				ImportStateVerify: true,
			},
		},
	})
}

// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
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
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Membership writes (the copy that adds devices, the PUT that removes them) take cw.UseMutex, so they run one at a time. The
// stub answers an overlapping write like Crosswork answers a deadlock; without the lock the parallel calls below fail.
func TestMembershipWritesAreSerialized(t *testing.T) {
	s := newLookupStub(t, 8)
	var inflight, overlaps, writes atomic.Int32
	inner := s.ssoStub.Config.Handler
	s.ssoStub.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, groupingPath+"/device/g-unassigned"):
			var devices []string
			for n := 0; n < 8; n++ {
				devices = append(devices, fmt.Sprintf(`{"uuid":"uuid-%d"}`, n))
			}
			fmt.Fprintf(w, `{"devices":[%s],"total":8}`, strings.Join(devices, ","))
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, groupingPath+"/device/"):
			fmt.Fprint(w, `{"devices":[],"total":0}`)
		case r.Method != http.MethodGet && strings.HasPrefix(r.URL.Path, groupingPath+"/group/"):
			writes.Add(1)
			inflight.Add(1)
			defer inflight.Add(-1)
			time.Sleep(100 * time.Millisecond) // long enough for parallel writes to overlap
			if inflight.Load() > 1 {
				overlaps.Add(1)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, `{"error":"pq: deadlock detected","code":500}`)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			inner.ServeHTTP(w, r)
		}
	})
	c, _ := registeredClient(s.URL, "admin", "stub-password", true, 0, "test")

	var wg sync.WaitGroup
	errs := make(chan error, 8)
	for i := 0; i < 4; i++ {
		wg.Add(2)
		go func(i int) {
			defer wg.Done()
			errs <- addMembers(c, fmt.Sprintf("group-%d", i), []string{fmt.Sprintf("uuid-%d", i)})
		}(i)
		go func(i int) {
			defer wg.Done()
			errs <- removeMembers(c, fmt.Sprintf("group-%d", i), []string{fmt.Sprintf("uuid-%d", i+4)})
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Errorf("a membership write failed: %v", err)
		}
	}
	if overlaps.Load() != 0 || writes.Load() != 8 {
		t.Errorf("%d of %d writes overlapped; want 0 of 8", overlaps.Load(), writes.Load())
	}
}

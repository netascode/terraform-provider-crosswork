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

// A Terraform module identifies itself in the User-Agent of the controller request logs, which separates module-driven traffic from direct provider use.

import (
	"context"
	"strings"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
)

var uaApplied sync.Map

type cwProviderMeta struct {
	ModuleName types.String `tfsdk:"module_name"`
}

func applyProviderMeta(clients map[string]*cw.Client, ctx context.Context, meta tfsdk.Config) {
	// meta.Raw is null when the module declares no provider_meta block for us.
	if meta.Raw.IsNull() {
		return
	}

	var m cwProviderMeta
	if diags := meta.Get(ctx, &m); diags.HasError() {
		// Never fail an operation over telemetry; just skip the marker.
		return
	}

	module := m.ModuleName.ValueString()
	if module == "" {
		return
	}

	for _, client := range clients {
		if client == nil {
			continue
		}
		if _, done := uaApplied.LoadOrStore(client, true); done {
			continue
		}
		ua := client.UserAgent
		if ua == "" {
			client.UserAgent = module
			continue
		}
		const vendor = " Cisco"
		if strings.HasSuffix(ua, vendor) {
			client.UserAgent = ua[:len(ua)-len(vendor)] + " " + module + vendor
		} else {
			client.UserAgent = ua + " " + module
		}
	}
}

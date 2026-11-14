// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

func TestHostNameOf(t *testing.T) {
	node := gjson.Parse(`{"data":[{"host_name":"ios"}]}`)
	empty := gjson.Parse(`{"data":[{"inventory_id":"a"}]}`)
	blank := gjson.Parse(`{"data":[{"host_name":""}]}`)

	if got := hostNameOf(types.StringValue("mine"), node); got.ValueString() != "mine" {
		t.Errorf("a configured host name is kept, got %v", got)
	}
	if got := hostNameOf(types.StringUnknown(), node); got.ValueString() != "ios" {
		t.Errorf("an unset host name is the one the controller holds, got %v", got)
	}
	for name, res := range map[string]gjson.Result{"absent": empty, "blank": blank} {
		if got := hostNameOf(types.StringUnknown(), res); !got.IsNull() {
			t.Errorf("%s: an unset host name the node does not have is null, got %v", name, got)
		}
	}
}

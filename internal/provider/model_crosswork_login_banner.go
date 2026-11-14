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

// Section below is generated&owned by "gen/generator.go". //template:begin imports
import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type LoginBanner struct {
	Instance    types.String `tfsdk:"instance"`
	Id          types.String `tfsdk:"id"`
	ShowMessage types.Bool   `tfsdk:"show_message"`
	UserAck     types.Bool   `tfsdk:"user_ack"`
	Message     types.String `tfsdk:"message"`
	Icon        types.String `tfsdk:"icon"`
	Title       types.String `tfsdk:"title"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data LoginBanner) getPath() string {
	return "/crosswork/platform/v2/cluster/banner/set"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data LoginBanner) getPathGet() string {
	return "/crosswork/platform/v2/cluster/banner/get"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// Escapes the id so it cannot change the structure of a request body.
func (data LoginBanner) idJSON() string {
	id := data.Id.ValueString()
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}

func (data LoginBanner) readGet(c *cw.Client, path string) (cw.Res, error) {
	return c.Post(path, strings.ReplaceAll(`{}`, "{id}", data.idJSON()), cw.ReadOnly)
}

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data LoginBanner) toBody(ctx context.Context, state LoginBanner) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.ShowMessage.IsNull() {
		body, _ = sjson.Set(body, "ShowMessage", data.ShowMessage.ValueBool())
	}
	if !data.UserAck.IsNull() {
		body, _ = sjson.Set(body, "UserAck", data.UserAck.ValueBool())
	}
	if !data.Message.IsNull() {
		body, _ = sjson.Set(body, "Message", data.Message.ValueString())
	}
	if !data.Icon.IsNull() {
		body, _ = sjson.Set(body, "Icon", data.Icon.ValueString())
	}
	if !data.Title.IsNull() {
		body, _ = sjson.Set(body, "Title", data.Title.ValueString())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *LoginBanner) fromBody(ctx context.Context, res gjson.Result) {
	data.Id = types.StringValue("login_banner")
	if value := res.Get("ShowMessage"); value.Exists() {
		data.ShowMessage = types.BoolValue(value.Bool())
	} else {
		data.ShowMessage = types.BoolValue(false)
	}
	if value := res.Get("UserAck"); value.Exists() {
		data.UserAck = types.BoolValue(value.Bool())
	} else {
		data.UserAck = types.BoolValue(false)
	}
	if value := res.Get("Message"); value.Exists() {
		data.Message = types.StringValue(value.String())
	} else {
		data.Message = types.StringValue("Welcome to Cisco Crosswork. Logging into the system means that you acknowledge and agree to comply with Cisco Crosswork terms and conditions.")
	}
	if value := res.Get("Icon"); value.Exists() {
		data.Icon = types.StringValue(value.String())
	} else {
		data.Icon = types.StringValue("BANNER_ICON_INFO")
	}
	if value := res.Get("Title"); value.Exists() {
		data.Title = types.StringValue(value.String())
	} else {
		data.Title = types.StringValue("Crosswork Legal Disclaimer")
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *LoginBanner) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("ShowMessage"); value.Exists() && !data.ShowMessage.IsNull() {
		data.ShowMessage = types.BoolValue(value.Bool())
	} else if data.ShowMessage.ValueBool() != false {
		data.ShowMessage = types.BoolNull()
	}
	if value := res.Get("UserAck"); value.Exists() && !data.UserAck.IsNull() {
		data.UserAck = types.BoolValue(value.Bool())
	} else if data.UserAck.ValueBool() != false {
		data.UserAck = types.BoolNull()
	}
	if value := res.Get("Message"); value.Exists() && !data.Message.IsNull() {
		data.Message = types.StringValue(value.String())
	} else if data.Message.ValueString() != "Welcome to Cisco Crosswork. Logging into the system means that you acknowledge and agree to comply with Cisco Crosswork terms and conditions." {
		data.Message = types.StringNull()
	}
	if value := res.Get("Icon"); value.Exists() && !data.Icon.IsNull() {
		data.Icon = types.StringValue(value.String())
	} else if data.Icon.ValueString() != "BANNER_ICON_INFO" {
		data.Icon = types.StringNull()
	}
	if value := res.Get("Title"); value.Exists() && !data.Title.IsNull() {
		data.Title = types.StringValue(value.String())
	} else if data.Title.ValueString() != "Crosswork Legal Disclaimer" {
		data.Title = types.StringNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *LoginBanner) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.ShowMessage.IsNull() {
		return false
	}
	if !data.UserAck.IsNull() {
		return false
	}
	if !data.Message.IsNull() {
		return false
	}
	if !data.Icon.IsNull() {
		return false
	}
	if !data.Title.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

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

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type DemandMatrixDomain struct {
	Instance types.String `tfsdk:"instance"`
	Id       types.String `tfsdk:"id"`
	DomainId types.String `tfsdk:"domain_id"`
	Enable   types.Bool   `tfsdk:"enable"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data DemandMatrixDomain) getPath() string {
	return "/crosswork/optima/api/v1/enableDdm"
}

// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath

// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete

func (data DemandMatrixDomain) getPathDelete() string {
	return "/crosswork/optima/api/v1/enableDdm"
}

// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet

func (data DemandMatrixDomain) getPathGet() string {
	return "/crosswork/optima/api/v1/getDdmDomains"
}

// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost

// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut

// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan

// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet

// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery

// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data DemandMatrixDomain) toBody(ctx context.Context, state DemandMatrixDomain) string {
	body := ""
	put := false
	if state.Id.ValueString() != "" {
		put = true
	}
	_ = put
	if !data.DomainId.IsNull() {
		body, _ = sjson.Set(body, "domainId", data.DomainId.ValueString())
	}
	if !data.Enable.IsNull() {
		body, _ = sjson.Set(body, "enable", data.Enable.ValueBool())
	}
	return body
}

// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *DemandMatrixDomain) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("domainId"); value.Exists() {
		data.DomainId = types.StringValue(value.String())
	} else {
		data.DomainId = types.StringNull()
	}
	if value := res.Get("enabled"); value.Exists() {
		data.Enable = types.BoolValue(value.Bool())
	} else {
		data.Enable = types.BoolValue(false)
	}
}

// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *DemandMatrixDomain) updateFromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("domainId"); value.Exists() && !data.DomainId.IsNull() {
		data.DomainId = types.StringValue(value.String())
	} else {
		data.DomainId = types.StringNull()
	}
	if value := res.Get("enabled"); value.Exists() && !data.Enable.IsNull() {
		data.Enable = types.BoolValue(value.Bool())
	} else if data.Enable.ValueBool() != false {
		data.Enable = types.BoolNull()
	}
}

// End of section. //template:end updateFromBody

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *DemandMatrixDomain) isNull(ctx context.Context, res gjson.Result) bool {
	if !data.Enable.IsNull() {
		return false
	}
	return true
}

// End of section. //template:end isNull

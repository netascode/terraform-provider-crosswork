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
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

func testPolicy() PerformancePolicies {
	return PerformancePolicies{
		Name:           types.StringValue("cwtest-policy"),
		PolicyTemplate: types.StringValue("INTERFACE"),
		Description:    types.StringValue("Test policy"),
		Contact:        types.StringValue("noc@example.com"),
		Active:         types.BoolValue(true),
		Devices:        targetSet([]string{"u-2", "u-1"}),
		DeviceGroups:   types.SetNull(types.StringType),
		PortGroups:     types.SetNull(types.StringType),
		SchemaIntervals: []PerformancePoliciesSchemaIntervals{
			{Schema: types.StringValue("CPU"), Seconds: types.Int64Value(600)},
			{Schema: types.StringValue("QOS"), Seconds: types.Int64Value(0)},
		},
		Thresholds: map[string]PerformancePoliciesThresholds{
			"qosDrop": {Conditions: []PerformancePoliciesThresholdsConditions{
				{
					Relation: types.StringValue("GE"), Value: types.StringValue("1e-6"), Count: types.Int64Value(2),
					Severity: types.StringValue("MAJOR"), MatchCos: mustInt64Set(3, 1),
				},
				{
					Relation: types.StringValue("GT"), Value: types.StringValue("90"), Count: types.Int64Value(1),
					Severity: types.StringValue("CRITICAL"), MatchCos: types.SetNull(types.Int64Type),
				},
			}},
		},
	}
}

func mustInt64Set(values ...int64) types.Set {
	set, _ := types.SetValueFrom(context.Background(), types.Int64Type, values)
	return set
}

func TestPolicyBody(t *testing.T) {
	ctx := context.Background()
	body := gjson.Parse(testPolicy().toBody(ctx, PerformancePolicies{}))
	assert.Equal(t, "cwtest-policy", body.Get("name").String())
	assert.Equal(t, "INTERFACE", body.Get("policyTemplate").String())
	assert.Equal(t, `{"contact":"noc@example.com"}`, body.Get("tag").String())
	assert.True(t, body.Get("active").Bool())
	assert.Equal(t, "u-1,u-2", body.Get("devices").String(), "sorted, comma joined, no spaces")
	assert.Equal(t, "", body.Get("deviceGroups").String(), "an unused target kind is an empty string")
	assert.Equal(t, int64(600), body.Get("schemasInterval.CPU").Int())
	assert.Equal(t, int64(0), body.Get("schemasInterval.QOS").Int())
	assert.False(t, body.Get("id").Exists(), "a create has no id")

	first := body.Get("thresholds.qosDrop.0")
	assert.Equal(t, int64(1), first.Get("conditionID").Int())
	assert.Equal(t, gjson.String, first.Get("value").Type, "the value is sent as a string")
	assert.Equal(t, "1e-6", first.Get("value").String())
	assert.Equal(t, "matchCos", first.Get("extraConditionFields.0.name").String())
	assert.Equal(t, `["1","3"]`, first.Get("extraConditionFields.0.values").Raw, "classes sorted, as strings")
	second := body.Get("thresholds.qosDrop.1")
	assert.Equal(t, int64(2), second.Get("conditionID").Int(), "conditions are numbered from 1 for each metric")
	assert.False(t, second.Get("extraConditionFields").Exists())

	update := gjson.Parse(testPolicy().toBody(ctx, PerformancePolicies{Id: types.StringValue("7")}))
	assert.Equal(t, gjson.String, update.Get("id").Type, "an update sends the id as a string")
	assert.Equal(t, "7", update.Get("id").String())
}

func TestPolicyBodyEmptyAndNoContact(t *testing.T) {
	data := testPolicy()
	data.Contact, data.Description, data.Thresholds = types.StringNull(), types.StringNull(), nil
	body := gjson.Parse(data.toBody(context.Background(), PerformancePolicies{}))
	assert.Equal(t, "", body.Get("tag").String())
	assert.Equal(t, "", body.Get("description").String())
	assert.Equal(t, "{}", body.Get("thresholds").Raw, "an empty thresholds map clears the thresholds on an update")
}

func TestPolicyRoundTrip(t *testing.T) {
	ctx := context.Background()
	want := testPolicy()
	wire := want
	var got PerformancePolicies
	got.fromBody(ctx, gjson.Parse(`{"monitoringPolicy":`+wire.toBody(ctx, PerformancePolicies{Id: types.StringValue("7")})+`}`))
	assert.Equal(t, "7", got.Id.ValueString())
	assert.Equal(t, "noc@example.com", got.Contact.ValueString())
	assert.Equal(t, []string{"u-1", "u-2"}, setStrings(ctx, got.Devices))
	assert.True(t, got.DeviceGroups.IsNull())
	assert.Equal(t, 2, len(got.SchemaIntervals))
	conditions := got.Thresholds["qosDrop"].Conditions
	assert.Equal(t, 2, len(conditions))
	assert.Equal(t, "1e-6", conditions[0].Value.ValueString())
	assert.Equal(t, int64(2), conditions[0].Count.ValueInt64())
	assert.Equal(t, 2, len(conditions[0].MatchCos.Elements()))
	assert.True(t, conditions[1].MatchCos.IsNull(), "no extra field reads as an unset set")
}

func TestPolicyFromBodyTags(t *testing.T) {
	assert.True(t, contactFromTag("").IsNull())
	assert.True(t, contactFromTag("plain text").IsNull(), "a tag the UI did not write reads as unset")
	assert.True(t, contactFromTag(`{"contact":""}`).IsNull())
	assert.Equal(t, "noc", contactFromTag(`{"contact":"noc"}`).ValueString())
	assert.Equal(t, []string{"a", "b"}, splitTargets(" a, b,,a "))
	assert.Equal(t, 0, len(splitTargets("")))
}

func TestPolicyUpdateFromBodyKeepsEmptySets(t *testing.T) {
	ctx := context.Background()
	res := gjson.Parse(`{"monitoringPolicy":{"id":9,"name":"p","policyTemplate":"INTERFACE","active":false,
	  "devices":"","deviceGroups":"g-1","portGroups":"","schemasInterval":{"CPU":0},"thresholds":{}}}`)
	empty, _ := types.SetValueFrom(ctx, types.StringType, []string{})
	state := PerformancePolicies{Devices: empty, DeviceGroups: empty, PortGroups: types.SetNull(types.StringType), Thresholds: map[string]PerformancePoliciesThresholds{}}
	state.updateFromBody(ctx, res)
	assert.Equal(t, 0, len(state.Devices.Elements()))
	assert.False(t, state.Devices.IsNull(), "an empty set in the state stays empty")
	assert.Equal(t, []string{"g-1"}, setStrings(ctx, state.DeviceGroups), "a target set outside Terraform shows as a difference")
	assert.True(t, state.PortGroups.IsNull())
	assert.True(t, state.Thresholds != nil)
	assert.False(t, state.Active.ValueBool())
	assert.Equal(t, "9", state.Id.ValueString())
}

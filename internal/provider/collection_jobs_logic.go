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

// Written by hand: the body has an input and an output list that are paired by position, device names and a destination name
// become ids, a cadence in seconds is sent in milliseconds, one PUT creates and updates, a missing job is a failure in HTTP 200,
// and a deleted job stays listed for a while.

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

const (
	collectionJobPath    = "/crosswork/collection/v1/collectionjob"
	collectionConfigPath = collectionJobPath + "/config/query"
	collectionStatusPath = collectionJobPath + "/summary/query"
	collectionStatePath  = collectionJobPath + "/state/query"

	// Failure text the controller answers, in HTTP 200, for a job it does not hold.
	collectionJobMissing = "is not found in db"
)

// collectionGoneAttempts and collectionGoneInterval bound the wait for a deleted job to leave the controller's lists;
// the interval is a variable so the tests can shorten it.
const collectionGoneAttempts = 30

var collectionGoneInterval = time.Second

var collectionCollectors = map[string]string{"CLI": "CLI_COLLECTOR", "SNMP": "SNMP_COLLECTOR"}

func collectionJobId(application, contextID string) string { return application + "," + contextID }

func collectionJobKey(application, contextID string) string {
	body, _ := sjson.Set(`{}`, "application_context.application_id", application)
	body, _ = sjson.Set(body, "application_context.context_id", contextID)
	return body
}

func collectionMissing(err error, res cw.Res) bool {
	if err == nil {
		return false
	}
	return strings.Contains(res.Get("result.error.error").String()+" "+err.Error(), collectionJobMissing)
}

type sensorWire struct {
	data    string
	cadence int64 // seconds
	topic   string
}

func collectionSensors(plan CollectionJobs) []sensorWire {
	var out []sensorWire
	for _, s := range plan.CliSensors {
		data, _ := sjson.Set(`{}`, "cli_sensor.command", s.Command.ValueString())
		out = append(out, sensorWire{data, s.Cadence.ValueInt64(), s.Topic.ValueString()})
	}
	for _, s := range plan.SnmpSensors {
		data, _ := sjson.Set(`{}`, "snmp_sensor.snmp_mib.snmp_operation", s.Operation.ValueString())
		data, _ = sjson.Set(data, "snmp_sensor.snmp_mib.oids", []string{s.Oid.ValueString()})
		out = append(out, sensorWire{data, s.Cadence.ValueInt64(), s.Topic.ValueString()})
	}
	return out
}

func collectionDeviceNames(ctx context.Context, devices types.Set) ([]string, diag.Diagnostics) {
	var names []string
	if devices.IsNull() || devices.IsUnknown() {
		return nil, nil
	}
	diags := devices.ElementsAs(ctx, &names, false)
	sort.Strings(names)
	return names, diags
}

// The input and the output list go out in the same order, because the controller pairs them by position. The controller
// refuses what does not fit together with its own message, so nothing is checked here.
func collectionJobBody(plan CollectionJobs, names []string, uuids map[string]string, destinationId string) string {
	body := `{"collection_job":{}}`
	body, _ = sjson.Set(body, "collection_job.application_context.application_id", plan.ApplicationId.ValueString())
	body, _ = sjson.Set(body, "collection_job.application_context.context_id", plan.ContextId.ValueString())
	body, _ = sjson.Set(body, "collection_job.collection_mode.lifetime_type", "APPLICATION_MANAGED")
	body, _ = sjson.Set(body, "collection_job.collection_mode.collector_type", collectionCollectors[plan.CollectorType.ValueString()])
	if len(names) > 0 {
		ids := make([]string, len(names))
		for i, n := range names {
			ids[i] = uuids[n]
		}
		body, _ = sjson.Set(body, "collection_job.job_device_set.device_set.devices.device_ids", ids)
	}
	if !plan.DeviceTag.IsNull() && !plan.DeviceTag.IsUnknown() && plan.DeviceTag.ValueString() != "" {
		body, _ = sjson.Set(body, "collection_job.job_device_set.device_set.device_group", plan.DeviceTag.ValueString())
	}
	body, _ = sjson.SetRaw(body, "collection_job.sensor_input_configs", `[]`)
	body, _ = sjson.SetRaw(body, "collection_job.sensor_output_configs", `[]`)
	for _, s := range collectionSensors(plan) {
		in, _ := sjson.SetRaw(`{}`, "sensor_data", s.data)
		in, _ = sjson.Set(in, "cadence_in_millisec", strconv.FormatInt(s.cadence*1000, 10))
		body, _ = sjson.SetRaw(body, "collection_job.sensor_input_configs.-1", in)

		out, _ := sjson.SetRaw(`{}`, "sensor_data", s.data)
		out, _ = sjson.Set(out, "destination.destination_id", destinationId)
		out, _ = sjson.Set(out, "destination.destination_name", plan.Destination.ValueString())
		out, _ = sjson.Set(out, "destination.context_id", s.topic)
		body, _ = sjson.SetRaw(body, "collection_job.sensor_output_configs.-1", out)
	}
	return body
}

func collectionResolve(ctx context.Context, c *cw.Client, plan CollectionJobs) (names []string, uuids map[string]string, destinationId string, diags diag.Diagnostics) {
	names, diags = collectionDeviceNames(ctx, plan.Devices)
	if diags.HasError() {
		return
	}
	if len(names) > 0 {
		var err error
		if uuids, err = deviceUUIDs(c, names); err != nil {
			diags.AddError("Client Error", fmt.Sprintf("Failed to look up the devices of the job, got error: %s", err))
			return
		}
	}
	var err error
	if destinationId, err = destinationIdByNameOrId(c, plan.Destination.ValueString()); err != nil {
		diags.AddError("Client Error", fmt.Sprintf("Failed to look up the destination of the job, got error: %s", err))
	} else if destinationId == "" {
		diags.AddError("Destination not found", fmt.Sprintf("There is no destination named %s on the cluster", plan.Destination.ValueString()))
	}
	return
}

// One PUT creates the job and, when it exists, replaces its sensors and destination.
func putCollectionJob(c *cw.Client, body string) error {
	if res, err := c.Put(collectionJobPath, body); err != nil {
		return fmt.Errorf("%w, %s", err, res.String())
	}
	return nil
}

func readCollectionJob(c *cw.Client, application, contextID string) (row gjson.Result, found bool, err error) {
	res, err := c.Post(collectionConfigPath, collectionJobKey(application, contextID), cw.ReadOnly)
	if err != nil {
		if collectionMissing(err, res) {
			return gjson.Result{}, false, nil
		}
		return gjson.Result{}, false, fmt.Errorf("%w, %s", err, res.String())
	}
	row = res.Get("collection_job_configuration_list.0")
	return row, row.Exists(), nil
}

// After a delete the status query drops the job before the configuration query does, so a create looks here.
func collectionJobListed(c *cw.Client, application, contextID string) (bool, error) {
	res, err := c.Post(collectionStatusPath, collectionJobKey(application, contextID), cw.ReadOnly)
	if err != nil {
		if collectionMissing(err, res) {
			return false, nil
		}
		return false, fmt.Errorf("%w, %s", err, res.String())
	}
	row := res.Get("collection_job_status_list.0")
	return row.Exists() && row.Get("phase").String() != "DELETED", nil
}

func collectionJobStatus(c *cw.Client, application, contextID string) (status, phase, state string, err error) {
	res, err := c.Post(collectionStatusPath, collectionJobKey(application, contextID), cw.ReadOnly)
	if err != nil {
		return "", "", "", fmt.Errorf("%w, %s", err, res.String())
	}
	row := res.Get("collection_job_status_list.0")
	status, phase = row.Get("status").String(), row.Get("phase").String()
	res, err = c.Post(collectionStatePath, collectionJobKey(application, contextID), cw.ReadOnly)
	if err != nil {
		return "", "", "", fmt.Errorf("%w, %s", err, res.String())
	}
	return status, phase, res.Get("collection_life_cycle_states.0.life_cycle_state").String(), nil
}

// Waits until the controller no longer lists the job, so that a create that follows (the replace of a job) does not
// meet the old one.
func deleteCollectionJob(ctx context.Context, c *cw.Client, application, contextID string) (gone bool, err error) {
	res, err := c.DeleteBody(collectionJobPath, collectionJobKey(application, contextID))
	if err != nil && !collectionMissing(err, res) {
		return false, fmt.Errorf("%w, %s", err, res.String())
	}
	for i := 0; i < collectionGoneAttempts; i++ {
		listed, err := collectionJobListed(c, application, contextID)
		if err != nil {
			return false, err
		}
		_, found, err := readCollectionJob(c, application, contextID)
		if err != nil {
			return false, err
		}
		if !listed && !found {
			return true, nil
		}
		if !sleepContext(ctx, collectionGoneInterval) {
			break
		}
	}
	return false, nil
}

// Devices are given back by inventory id, in the spelling the model had when it is the same name; a uuid with no
// device is left out, which shows as a difference.
func collectionJobFromRow(ctx context.Context, c *cw.Client, row gjson.Result, m *CollectionJobs) diag.Diagnostics {
	var diags diag.Diagnostics
	m.ApplicationId = types.StringValue(row.Get("application_context.application_id").String())
	m.ContextId = types.StringValue(row.Get("application_context.context_id").String())
	m.Id = types.StringValue(collectionJobId(m.ApplicationId.ValueString(), m.ContextId.ValueString()))
	collector := row.Get("collection_mode.collector_type").String()
	m.CollectorType = types.StringValue(collector)
	for name, wire := range collectionCollectors {
		if wire == collector {
			m.CollectorType = types.StringValue(name)
		}
	}

	m.DeviceTag = types.StringNull()
	if tag := row.Get("job_device_set.device_set.device_group").String(); tag != "" {
		m.DeviceTag = types.StringValue(tag)
	}
	var uuids []string
	for _, d := range row.Get("job_device_set.device_set.devices.device_ids").Array() {
		uuids = append(uuids, d.String())
	}
	spelling := map[string]string{}
	if known, d := collectionDeviceNames(ctx, m.Devices); !d.HasError() {
		for _, n := range known {
			spelling[strings.ToLower(n)] = n
		}
	}
	m.Devices = types.SetNull(types.StringType)
	if len(uuids) > 0 {
		ids, err := inventoryIDs(c, uuids)
		if err != nil {
			diags.AddError("Client Error", fmt.Sprintf("Failed to look up the devices of the job, got error: %s", err))
			return diags
		}
		var elems []attr.Value
		for _, uuid := range uuids {
			if id, ok := ids[uuid]; ok {
				if kept, ok := spelling[strings.ToLower(id)]; ok {
					id = kept
				}
				elems = append(elems, types.StringValue(id))
			}
		}
		m.Devices = types.SetValueMust(types.StringType, elems)
	}

	inputs, outputs := row.Get("sensor_input_configs").Array(), row.Get("sensor_output_configs").Array()
	m.Destination = types.StringNull()
	m.CliSensors, m.SnmpSensors = nil, nil
	for i, in := range inputs {
		var out gjson.Result
		if i < len(outputs) {
			out = outputs[i]
		}
		if name := out.Get("destination.destination_name").String(); name != "" && m.Destination.IsNull() {
			m.Destination = types.StringValue(name)
		}
		cadence := types.Int64Value(in.Get("cadence_in_millisec").Int() / 1000)
		topic := types.StringValue(out.Get("destination.context_id").String())
		switch {
		case in.Get("sensor_data.cli_sensor").Exists():
			m.CliSensors = append(m.CliSensors, CollectionJobsCliSensors{
				Command: types.StringValue(in.Get("sensor_data.cli_sensor.command").String()), Cadence: cadence, Topic: topic})
		case in.Get("sensor_data.snmp_sensor").Exists():
			m.SnmpSensors = append(m.SnmpSensors, CollectionJobsSnmpSensors{
				Oid:       types.StringValue(in.Get("sensor_data.snmp_sensor.snmp_mib.oids.0").String()),
				Operation: types.StringValue(in.Get("sensor_data.snmp_sensor.snmp_mib.snmp_operation").String()),
				Cadence:   cadence, Topic: topic})
		}
	}
	return diags
}

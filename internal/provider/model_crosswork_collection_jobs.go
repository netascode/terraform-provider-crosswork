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

import "github.com/hashicorp/terraform-plugin-framework/types"

type CollectionJobs struct {
	Id            types.String                `tfsdk:"id"`
	Instance      types.String                `tfsdk:"instance"`
	ApplicationId types.String                `tfsdk:"application_id"`
	ContextId     types.String                `tfsdk:"context_id"`
	CollectorType types.String                `tfsdk:"collector_type"`
	Devices       types.Set                   `tfsdk:"devices"`
	DeviceTag     types.String                `tfsdk:"device_tag"`
	Destination   types.String                `tfsdk:"destination"`
	CliSensors    []CollectionJobsCliSensors  `tfsdk:"cli_sensors"`
	SnmpSensors   []CollectionJobsSnmpSensors `tfsdk:"snmp_sensors"`
}

type CollectionJobsCliSensors struct {
	Command types.String `tfsdk:"command"`
	Cadence types.Int64  `tfsdk:"cadence"`
	Topic   types.String `tfsdk:"topic"`
}

type CollectionJobsSnmpSensors struct {
	Oid       types.String `tfsdk:"oid"`
	Operation types.String `tfsdk:"operation"`
	Cadence   types.Int64  `tfsdk:"cadence"`
	Topic     types.String `tfsdk:"topic"`
}

// The resource leaves out the status fields, which change while the job runs.
type CollectionJobsData struct {
	Id             types.String                `tfsdk:"id"`
	Instance       types.String                `tfsdk:"instance"`
	ApplicationId  types.String                `tfsdk:"application_id"`
	ContextId      types.String                `tfsdk:"context_id"`
	CollectorType  types.String                `tfsdk:"collector_type"`
	Devices        types.Set                   `tfsdk:"devices"`
	DeviceTag      types.String                `tfsdk:"device_tag"`
	Destination    types.String                `tfsdk:"destination"`
	CliSensors     []CollectionJobsCliSensors  `tfsdk:"cli_sensors"`
	SnmpSensors    []CollectionJobsSnmpSensors `tfsdk:"snmp_sensors"`
	Status         types.String                `tfsdk:"status"`
	Phase          types.String                `tfsdk:"phase"`
	LifeCycleState types.String                `tfsdk:"life_cycle_state"`
}

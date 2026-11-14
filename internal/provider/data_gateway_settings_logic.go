// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License Version 2.0 (the "License"),
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
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cw "github.com/netascode/go-crosswork"
)

// The global parameters and the resource configuration are each one set for the cluster, so the resources change them and
// never delete them.

const (
	globalParametersQueryPath = "/crosswork/dg-manager/v1/command/global-parameter/query"
	resourcesQueryPath        = "/crosswork/dg-manager/v2/oam/container/resources/query"
)

var globalParameterKeys = []struct {
	api, attribute string
	flag           bool
}{
	{"CLI_SESSIONS", "cli_sessions", false},
	{"SSH_SESSION_TIMEOUT_IN_SECS", "ssh_session_timeout", false},
	{"SNMP_TRAP_PORT", "snmp_trap_port", false},
	{"SYSLOG_UDP_PORT", "syslog_udp_port", false},
	{"SYSLOG_TCP_PORT", "syslog_tcp_port", false},
	{"SYSLOG_TLS_PORT", "syslog_tls_port", false},
	{"RESYNC_ENGINE_DETAILS", "resync_snmpv3_details", true},
}

func flattenGlobalParameters(res gjson.Result) gjson.Result {
	body := "{}"
	for _, parameter := range res.Get("globalParameters").Array() {
		for _, key := range globalParameterKeys {
			if parameter.Get("key").String() != key.api {
				continue
			}
			if key.flag {
				body, _ = sjson.Set(body, key.attribute, parameter.Get("value.boolValue").Bool())
			} else {
				body, _ = sjson.Set(body, key.attribute, parameter.Get("value.uint32Value").Int())
			}
		}
	}
	return gjson.Parse(body)
}

func newCommandId() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func globalParametersCommand(flat gjson.Result) (string, bool) {
	list := "[]"
	count := 0
	for _, key := range globalParameterKeys {
		value := flat.Get(key.attribute)
		if !value.Exists() {
			continue
		}
		item, _ := sjson.Set("{}", "key", key.api)
		if key.flag {
			item, _ = sjson.Set(item, "value.boolValue", value.Bool())
		} else {
			item, _ = sjson.Set(item, "value.uint32Value", value.Int())
		}
		list, _ = sjson.SetRaw(list, "-1", item)
		count++
	}
	if count == 0 {
		return "", false
	}
	command, _ := sjson.Set("{}", "id", newCommandId())
	command, _ = sjson.Set(command, "commandName", "GLOBAL_PARAMETERS")
	command, _ = sjson.SetRaw(command, "args", `[{"key":"globalParams","value":{"globalParameterRequest":{}}}]`)
	command, _ = sjson.SetRaw(command, "args.0.value.globalParameterRequest.globalParameters", list)
	return command, true
}

func globalParametersFailed(res cw.Res) string {
	failed := res.Get("response.value.globalParameterResponse.failedParamList")
	if failed.IsArray() && len(failed.Array()) > 0 || failed.IsObject() && len(failed.Map()) > 0 {
		return failed.Raw
	}
	return ""
}

func globalParametersHeld(flat, current gjson.Result) bool {
	held := true
	flat.ForEach(func(key, value gjson.Result) bool {
		if current.Get(key.String()).Raw != value.Raw {
			held = false
		}
		return held
	})
	return held
}

func changedParameters(flat, current gjson.Result) gjson.Result {
	changed := "{}"
	flat.ForEach(func(key, value gjson.Result) bool {
		if current.Get(key.String()).Raw != value.Raw {
			changed, _ = sjson.SetRaw(changed, key.String(), value.Raw)
		}
		return true
	})
	return gjson.Parse(changed)
}

// Bounds the check of a write against the held values after the controller lists a parameter as failed; a variable so that
// tests can shorten it.
var settingsConfirmTimeout = 30 * time.Second

// The controller can list a parameter as failed that it has applied, so a failed list is checked against the values it
// holds before it is an error.
func patchGlobalParameters(ctx context.Context, c *cw.Client, path string, flat gjson.Result) error {
	// The controller refuses a request that changes nothing (500, "no changes found in global parameter to update"), which a new
	// resource sends when the controller already holds the values (a destroy leaves them), so only differing parameters are sent.
	if current, readErr := (DataGatewayGlobalParameters{}).readGet(c, globalParametersQueryPath); readErr == nil {
		flat = changedParameters(flat, current)
	}
	command, ok := globalParametersCommand(flat)
	if !ok {
		tflog.Debug(ctx, "data_gateway_global_parameters: No changes detected")
		return nil
	}
	res, err := c.Patch(path, command, cw.UseMutex)
	if err != nil {
		return fmt.Errorf("%s, %s", err, res.String())
	}
	failed := globalParametersFailed(res)
	if failed == "" {
		return nil
	}
	deadline := time.Now().Add(settingsConfirmTimeout)
	for {
		current, readErr := (DataGatewayGlobalParameters{}).readGet(c, globalParametersQueryPath)
		if readErr == nil && globalParametersHeld(flat, current) {
			return nil
		}
		if !time.Now().Add(settingsRetryInterval).Before(deadline) || !sleepContext(ctx, settingsRetryInterval) {
			return fmt.Errorf("the controller did not apply these parameters: %s", failed)
		}
	}
}

// Reads all parameters back, so the ones the configuration leaves out take the value of the controller.
func applyGlobalParameters(ctx context.Context, c *cw.Client, plan *DataGatewayGlobalParameters) error {
	flat := gjson.Parse(plan.toBody(ctx, DataGatewayGlobalParameters{}))
	if err := patchGlobalParameters(ctx, c, plan.getPath(), flat); err != nil {
		return err
	}
	res, err := plan.readGet(c, globalParametersQueryPath)
	if err != nil {
		return err
	}
	plan.updateFromBody(ctx, res)
	return nil
}

var collectorContainers = map[string]string{
	"CLI": "cli-collector", "GNMI": "gnmi-collector", "MDT": "mdt-collector", "NETCONF": "netconf-collector", "SNMP": "snmp-collector", "SYSLOG": "syslog-collector",
}

// A fixed order, so that a read gives the same list every time.
var collectorOrder = []string{"CLI", "GNMI", "MDT", "NETCONF", "SNMP", "SYSLOG"}

// A limit looks like 9216m or 9g.
func memoryMegabytes(limit string) (int64, bool) {
	limit = strings.ToLower(strings.TrimSpace(limit))
	if limit == "" {
		return 0, false
	}
	unit := limit[len(limit)-1]
	number := limit
	if unit < '0' || unit > '9' {
		number = limit[:len(limit)-1]
	} else {
		unit = 'b'
	}
	n, err := strconv.ParseFloat(number, 64)
	if err != nil {
		return 0, false
	}
	switch unit {
	case 'b':
		n = n / (1024 * 1024)
	case 'k':
		n = n / 1024
	case 'm':
	case 'g':
		n = n * 1024
	default:
		return 0, false
	}
	return int64(math.Round(n)), true
}

func flattenCollectors(res gjson.Result) gjson.Result {
	body := `{"collectors":[]}`
	for _, name := range collectorOrder {
		container := collectorContainers[name]
		for _, item := range res.Get("containers").Array() {
			if item.Get("containerName").String() != container {
				continue
			}
			entry, _ := sjson.Set("{}", "name", name)
			if mb, ok := memoryMegabytes(item.Get("memoryLimit").String()); ok {
				entry, _ = sjson.Set(entry, "memory", mb)
			}
			entry, _ = sjson.Set(entry, "enabled", !item.Get("notDeployed").Bool())
			body, _ = sjson.SetRaw(body, "collectors.-1", entry)
		}
	}
	return gjson.Parse(body)
}

func resourcesRequest(current gjson.Result, collectors []DataGatewayResourceConfigurationCollectors) (string, bool, error) {
	containers := current.Get("containers").Array()
	request := `{"containers":[]}`
	changed := false
	wanted := map[string]DataGatewayResourceConfigurationCollectors{}
	for _, collector := range collectors {
		name := collector.Name.ValueString()
		container, known := collectorContainers[name]
		if !known {
			return "", false, fmt.Errorf("collector %s is not one of CLI, GNMI, MDT, NETCONF, SNMP or SYSLOG", name)
		}
		found := false
		for _, item := range containers {
			if item.Get("containerName").String() == container {
				found = true
			}
		}
		if !found {
			return "", false, fmt.Errorf("the cluster has no %s collector", name)
		}
		wanted[container] = collector
	}
	for _, item := range containers {
		raw := item.Raw
		if collector, ok := wanted[item.Get("containerName").String()]; ok {
			if !collector.Memory.IsNull() && !collector.Memory.IsUnknown() {
				if mb, ok := memoryMegabytes(item.Get("memoryLimit").String()); !ok || mb != collector.Memory.ValueInt64() {
					raw, _ = sjson.Set(raw, "memoryLimit", fmt.Sprintf("%dm", collector.Memory.ValueInt64()))
					changed = true
				}
			}
			if !collector.Enabled.IsNull() && !collector.Enabled.IsUnknown() && item.Get("notDeployed").Bool() == collector.Enabled.ValueBool() {
				raw, _ = sjson.Set(raw, "notDeployed", !collector.Enabled.ValueBool())
				changed = true
			}
		}
		request, _ = sjson.SetRaw(request, "containers.-1", raw)
	}
	return request, changed, nil
}

// A variable so that tests can shorten it.
var settingsRetryInterval = 10 * time.Second

func resourcesBusy(message string) bool {
	return strings.Contains(message, "429") || strings.Contains(strings.ToLower(message), "too many requests")
}

// Repeats while the controller is busy with the previous change, up to the timeout.
func postResources(ctx context.Context, c *cw.Client, path, request string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		res, err := c.Post(path, request, cw.UseMutex)
		if err == nil {
			return nil
		}
		if !resourcesBusy(err.Error()+res.String()) || !time.Now().Add(settingsRetryInterval).Before(deadline) {
			return fmt.Errorf("%s, %s", err, res.String())
		}
		if !sleepContext(ctx, settingsRetryInterval) {
			return fmt.Errorf("%s, %s", err, res.String())
		}
	}
}

func applyResources(ctx context.Context, c *cw.Client, plan *DataGatewayResourceConfiguration) error {
	current, err := c.Post(resourcesQueryPath, "{}", cw.ReadOnly)
	if err != nil {
		return fmt.Errorf("%s, %s", err, current.String())
	}
	request, changed, err := resourcesRequest(current, plan.Collectors)
	if err != nil {
		return err
	}
	if !changed {
		return nil
	}
	return postResources(ctx, c, plan.getPath(), request, time.Duration(plan.UpdateTimeout.ValueInt64())*time.Second)
}

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

// Hand-written operations of crosswork_data_gateway_destinations, which the generator cannot express: a flat map
// of upper-case text properties with the numbers as strings, collector overrides as keys like CLI.LINGER_MS_CONFIG,
// an id that is found by name after the create, and every change a job that can fail inside HTTP 200.

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
)

const (
	destinationPath      = "/crosswork/dg-manager/v1/destinations"
	destinationQueryPath = "/crosswork/dg-manager/v1/destinations/query"
	destinationQueryBody = `{"limit":1000,"filter":{}}`

	destinationFamily  = "ROBOT_PROVIDER_DESTINATION"
	destinationInetV4  = "ROBOT_INET_ADDR_TYPE_v4"
	destinationKafka   = "destination_type_kafka"
	destinationGrpc    = "destination_type_grpc"
	destinationTransKa = "ROBOT_MSVC_TRANS_KAFKA"
	destinationTransGr = "ROBOT_MSVC_TRANS_GRPC"
)

// destinationListedAttempts and destinationRetryInterval bound the wait for a new destination to be listed; the
// interval is a variable so the tests can shorten it.
const destinationListedAttempts = 10

var destinationRetryInterval = time.Second

var destinationCollectors = []string{"CLI", "GNMI", "MDT", "SNMP", "SYSLOG"}

var destinationCollectorKeys = map[string]string{
	"buffer_memory":   "BUFFER_MEMORY_CONFIG",
	"batch_size":      "BATCH_SIZE_CONFIG",
	"linger_time":     "LINGER_MS_CONFIG",
	"request_timeout": "REQUEST_TIMEOUT_CONFIG",
}

var destinationDataSources = map[string]string{"data-gateway": "datagateway", "application": "application"}
var destinationServerTypes = map[string]string{"kafka": destinationKafka, "grpc": destinationGrpc}

func destinationKey(m map[string]string, wire string) string {
	for k, v := range m {
		if v == wire {
			return k
		}
	}
	return wire
}

type destinationWire struct {
	Uuid             string                    `json:"uuid,omitempty"`
	Name             string                    `json:"name"`
	Family           string                    `json:"family"`
	Properties       map[string]string         `json:"properties"`
	ConnectivityInfo []destinationConnectivity `json:"connectivity_info"`
}

type destinationConnectivity struct {
	Type    string               `json:"type"`
	Ipaddrs []destinationAddress `json:"ipaddrs,omitempty"`
	Fqdn    *destinationFqdn     `json:"fqdn,omitempty"`
	Port    int64                `json:"port"`
}

type destinationAddress struct {
	InetAf   string `json:"inet_af"`
	InetAddr string `json:"inet_addr"`
	Mask     string `json:"mask"`
}

type destinationFqdn struct {
	HostName   string `json:"host_name"`
	DomainName string `json:"domain_name"`
}

func setIntProperty(props map[string]string, key string, v types.Int64) {
	if !v.IsNull() && !v.IsUnknown() {
		props[key] = strconv.FormatInt(v.ValueInt64(), 10)
	}
}

// destinationProperties builds the flat property map. The controller keeps every property it is given, the Kafka values of a
// gRPC destination too, so nothing here depends on the server type. IS_SYSTEM_DEFINED is always false: the controller accepts true and then refuses every later change and delete.
func destinationProperties(data DataGatewayDestinations) map[string]string {
	props := map[string]string{
		"DESTINATION_TYPE":        destinationServerTypes[data.ServerType.ValueString()],
		"DISPATCH_SOURCE":         destinationDataSources[data.DataSource.ValueString()],
		"COMPRESSION_TYPE_CONFIG": data.CompressionMethod.ValueString(),
		"IS_SYSTEM_DEFINED":       "false",
		"IS_SECURITY_ENABLED":     strconv.FormatBool(data.Secure.ValueBool()),
	}
	if !data.Encoding.IsNull() && !data.Encoding.IsUnknown() {
		props["ENCODING"] = data.Encoding.ValueString()
	}
	if !data.AuthType.IsNull() && !data.AuthType.IsUnknown() {
		props["AUTH_TYPE"] = data.AuthType.ValueString()
	}
	if !data.DispatchType.IsNull() && !data.DispatchType.IsUnknown() {
		props["GRPC_DISPATCH_TYPE"] = data.DispatchType.ValueString()
	}
	setIntProperty(props, "MAX_REQUEST_SIZE_CONFIG", data.MaxMessageSize)
	setIntProperty(props, "BUFFER_MEMORY_CONFIG", data.BufferMemory)
	setIntProperty(props, "BATCH_SIZE_CONFIG", data.BatchSize)
	setIntProperty(props, "LINGER_MS_CONFIG", data.LingerTime)
	setIntProperty(props, "REQUEST_TIMEOUT_CONFIG", data.RequestTimeout)
	for _, c := range data.Collectors {
		prefix := c.Name.ValueString() + "."
		setIntProperty(props, prefix+destinationCollectorKeys["buffer_memory"], c.BufferMemory)
		setIntProperty(props, prefix+destinationCollectorKeys["batch_size"], c.BatchSize)
		setIntProperty(props, prefix+destinationCollectorKeys["linger_time"], c.LingerTime)
		setIntProperty(props, prefix+destinationCollectorKeys["request_timeout"], c.RequestTimeout)
	}
	return props
}

func destinationConnectivityInfo(data DataGatewayDestinations) []destinationConnectivity {
	transport := destinationTransKa
	if data.ServerType.ValueString() == "grpc" {
		transport = destinationTransGr
	}
	out := make([]destinationConnectivity, 0, len(data.Connections))
	for _, c := range data.Connections {
		info := destinationConnectivity{Type: transport, Port: c.Port.ValueInt64()}
		if !c.Address.IsNull() && c.Address.ValueString() != "" {
			info.Ipaddrs = []destinationAddress{{
				InetAf:   destinationInetV4,
				InetAddr: strings.TrimSpace(c.Address.ValueString()),
				Mask:     strconv.FormatInt(c.PrefixLength.ValueInt64(), 10),
			}}
		} else {
			info.Fqdn = &destinationFqdn{HostName: c.HostName.ValueString(), DomainName: c.DomainName.ValueString()}
		}
		out = append(out, info)
	}
	return out
}

// destinationRequestBody builds the body of a create (empty id) or an update; the controller takes a list.
func destinationRequestBody(data DataGatewayDestinations, id string) string {
	body, _ := json.Marshal(map[string][]destinationWire{"providers": {{
		Uuid:             id,
		Name:             data.Name.ValueString(),
		Family:           destinationFamily,
		Properties:       destinationProperties(data),
		ConnectivityInfo: destinationConnectivityInfo(data),
	}}})
	return string(body)
}

func destinationDeleteBody(id string) string {
	body, _ := json.Marshal(map[string][]map[string]string{"providers": {{"uuid": id}}})
	return string(body)
}

// It fails when the controller holds more destinations than the one query returns, so a destination is never reported
// missing because it is on a later page.
func listDestinations(c *cw.Client) (cw.Res, error) {
	res, err := c.Post(destinationQueryPath, destinationQueryBody, cw.ReadOnly)
	if err != nil {
		return res, err
	}
	if total, rows := res.Get("total_count").Int(), int64(len(res.Get("data").Array())); total > rows {
		return res, fmt.Errorf("the controller holds %d destinations and one query returned %d", total, rows)
	}
	return res, nil
}

func findDestination(res cw.Res, id, name string) gjson.Result {
	var found gjson.Result
	res.Get("data").ForEach(func(_, row gjson.Result) bool {
		if (id != "" && row.Get("uuid").String() == id) || (id == "" && name != "" && row.Get("name").String() == name) {
			found = row
			return false
		}
		return true
	})
	return found
}

func destinationIntProperty(props map[string]gjson.Result, key string) types.Int64 {
	v, ok := props[key]
	if !ok {
		return types.Int64Null()
	}
	n, err := strconv.ParseInt(strings.TrimSpace(v.String()), 10, 64)
	if err != nil {
		return types.Int64Null()
	}
	return types.Int64Value(n)
}

func destinationStringProperty(props map[string]gjson.Result, key string) types.String {
	if v, ok := props[key]; ok && v.String() != "" {
		return types.StringValue(v.String())
	}
	return types.StringNull()
}

func destinationConnections(row gjson.Result) []DataGatewayDestinationsConnections {
	var out []DataGatewayDestinationsConnections
	row.Get("connectivity_info").ForEach(func(_, info gjson.Result) bool {
		conn := DataGatewayDestinationsConnections{
			Address:      types.StringNull(),
			PrefixLength: types.Int64Null(),
			HostName:     types.StringNull(),
			DomainName:   types.StringNull(),
			Port:         types.Int64Value(info.Get("port").Int()),
		}
		if info.Get("fqdn").Exists() {
			conn.HostName = types.StringValue(info.Get("fqdn.host_name").String())
			conn.DomainName = types.StringValue(info.Get("fqdn.domain_name").String())
		} else {
			conn.Address = types.StringValue(info.Get("ipaddrs.0.inet_addr").String())
			if n, err := strconv.ParseInt(info.Get("ipaddrs.0.mask").String(), 10, 64); err == nil {
				conn.PrefixLength = types.Int64Value(n)
			}
		}
		out = append(out, conn)
		return true
	})
	return out
}

func destinationCollector(name string, props map[string]gjson.Result) (DataGatewayDestinationsCollectors, bool) {
	c := DataGatewayDestinationsCollectors{
		Name:           types.StringValue(name),
		BufferMemory:   destinationIntProperty(props, name+"."+destinationCollectorKeys["buffer_memory"]),
		BatchSize:      destinationIntProperty(props, name+"."+destinationCollectorKeys["batch_size"]),
		LingerTime:     destinationIntProperty(props, name+"."+destinationCollectorKeys["linger_time"]),
		RequestTimeout: destinationIntProperty(props, name+"."+destinationCollectorKeys["request_timeout"]),
	}
	ok := !c.BufferMemory.IsNull() || !c.BatchSize.IsNull() || !c.LingerTime.IsNull() || !c.RequestTimeout.IsNull()
	return c, ok
}

// destinationCollectorsFrom lists the collectors of known first, so the order of the configuration is kept.
func destinationCollectorsFrom(props map[string]gjson.Result, known []DataGatewayDestinationsCollectors) []DataGatewayDestinationsCollectors {
	var out []DataGatewayDestinationsCollectors
	seen := map[string]bool{}
	for _, k := range known {
		name := k.Name.ValueString()
		if c, ok := destinationCollector(name, props); ok && !seen[name] {
			out = append(out, c)
			seen[name] = true
		}
	}
	for _, name := range destinationCollectors {
		if c, ok := destinationCollector(name, props); ok && !seen[name] {
			out = append(out, c)
			seen[name] = true
		}
	}
	return out
}

// destinationIdFromList finds the id of a new destination in the list: the answer of a create has no id.
func destinationIdFromList(ctx context.Context, c *cw.Client, name string) (string, error) {
	var last error
	for i := 0; i < destinationListedAttempts; i++ {
		res, err := listDestinations(c)
		if err == nil {
			if id := findDestination(res, "", name).Get("uuid").String(); id != "" {
				return id, nil
			}
			last = fmt.Errorf("the destination %s is not listed", name)
		} else {
			last = err
		}
		if !sleepContext(ctx, destinationRetryInterval) {
			break
		}
	}
	return "", last
}

// createDestination: a refused change comes back as a failed job, which the library turns into an error.
func createDestination(ctx context.Context, c *cw.Client, plan DataGatewayDestinations) (string, error) {
	res, err := c.Post(destinationPath, destinationRequestBody(plan, ""), cw.UseMutex)
	if err != nil {
		return "", fmt.Errorf("%w, %s", err, res.String())
	}
	id, err := destinationIdFromList(ctx, c, plan.Name.ValueString())
	if err != nil {
		return "", fmt.Errorf("the destination was created, but its id could not be read: %w", err)
	}
	return id, nil
}

// updateDestination: the controller replaces the properties and the connections.
func updateDestination(c *cw.Client, plan DataGatewayDestinations, id string) error {
	res, err := c.Put(destinationPath, destinationRequestBody(plan, id), cw.UseMutex)
	if err != nil {
		return fmt.Errorf("%w, %s", err, res.String())
	}
	return nil
}

// deleteDestination: the controller refuses a system-defined destination.
func deleteDestination(c *cw.Client, id string) error {
	res, err := c.DeleteBody(destinationPath, destinationDeleteBody(id), cw.UseMutex)
	if err != nil {
		return fmt.Errorf("%w, %s", err, res.String())
	}
	return nil
}

func destinationIdByNameOrId(c *cw.Client, key string) (string, error) {
	res, err := listDestinations(c)
	if err != nil {
		return "", err
	}
	if row := findDestination(res, key, ""); row.Exists() {
		return row.Get("uuid").String(), nil
	}
	return findDestination(res, "", key).Get("uuid").String(), nil
}

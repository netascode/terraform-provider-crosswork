// Copyright © 2026 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

const globalParametersAnswer = `{"globalParameters":[{"key":"SNMP_TRAP_PORT","value":{"uint32Value":1062}},{"key":"CLI_SESSIONS","value":{"uint32Value":3}},{"key":"RESYNC_ENGINE_DETAILS","value":{"boolValue":false}},{"key":"SSH_SESSION_TIMEOUT_IN_SECS","value":{"uint32Value":120}}]}`

const containersAnswer = `{"containers":[{"serviceName":"dg-cli-collector-service","containerName":"cli-collector","memoryLimit":"9216m","cpuSet":"","notDeployed":false,"imageTag":"x"},{"serviceName":"dg-mdt-collector-service","containerName":"mdt-collector","memoryLimit":"5g","cpuSet":"","notDeployed":true}]}`

func TestFlattenGlobalParameters(t *testing.T) {
	flat := flattenGlobalParameters(gjson.Parse(globalParametersAnswer))
	assert.Equal(t, int64(1062), flat.Get("snmp_trap_port").Int())
	assert.Equal(t, int64(3), flat.Get("cli_sessions").Int())
	assert.Equal(t, int64(120), flat.Get("ssh_session_timeout").Int())
	assert.True(t, flat.Get("resync_snmpv3_details").Exists())
	assert.False(t, flat.Get("resync_snmpv3_details").Bool())
	assert.False(t, flat.Get("syslog_udp_port").Exists(), "a parameter the controller does not list is absent")
}

func TestGlobalParametersCommandSendsOnlyTheValuesSet(t *testing.T) {
	data := DataGatewayGlobalParameters{CliSessions: types.Int64Value(4), SyslogUdpPort: types.Int64Unknown(), SnmpTrapPort: types.Int64Null(), ResyncSnmpv3Details: types.BoolValue(true)}
	command, ok := globalParametersCommand(gjson.Parse(data.toBody(context.Background(), DataGatewayGlobalParameters{})))
	assert.True(t, ok)
	list := gjson.Get(command, "args.0.value.globalParameterRequest.globalParameters").Array()
	assert.Len(t, list, 2)
	assert.Equal(t, "GLOBAL_PARAMETERS", gjson.Get(command, "commandName").String())
	assert.Equal(t, "CLI_SESSIONS", list[0].Get("key").String())
	assert.Equal(t, int64(4), list[0].Get("value.uint32Value").Int())
	assert.Equal(t, "RESYNC_ENGINE_DETAILS", list[1].Get("key").String())
	assert.True(t, list[1].Get("value.boolValue").Bool())
	assert.Len(t, gjson.Get(command, "id").String(), 36)

	_, ok = globalParametersCommand(gjson.Parse(DataGatewayGlobalParameters{CliSessions: types.Int64Unknown()}.toBody(context.Background(), DataGatewayGlobalParameters{})))
	assert.False(t, ok, "nothing set sends nothing")
}

func TestGlobalParametersFailed(t *testing.T) {
	assert.Equal(t, "", globalParametersFailed(gjson.Parse(`{"response":{"value":{"globalParameterResponse":{"failedParamList":{}}}}}`)))
	assert.Equal(t, "", globalParametersFailed(gjson.Parse(`{"response":{"value":{"globalParameterResponse":{"failedParamList":[]}}}}`)))
	assert.NotEqual(t, "", globalParametersFailed(gjson.Parse(`{"response":{"value":{"globalParameterResponse":{"failedParamList":["CLI_SESSIONS"]}}}}`)))
}

func TestMemoryMegabytes(t *testing.T) {
	for limit, want := range map[string]int64{"9216m": 9216, "9g": 9216, "10240M": 10240, "1048576k": 1024, "1073741824": 1024, "500m": 500} {
		got, ok := memoryMegabytes(limit)
		assert.True(t, ok, limit)
		assert.Equal(t, want, got, limit)
	}
	for _, limit := range []string{"", "lots", "m"} {
		_, ok := memoryMegabytes(limit)
		assert.False(t, ok, limit)
	}
}

func TestFlattenCollectors(t *testing.T) {
	flat := flattenCollectors(gjson.Parse(containersAnswer))
	collectors := flat.Get("collectors").Array()
	assert.Len(t, collectors, 2)
	assert.Equal(t, "CLI", collectors[0].Get("name").String(), "the collectors come in a fixed order")
	assert.Equal(t, "MDT", collectors[1].Get("name").String())
	byName := map[string]gjson.Result{}
	for _, collector := range collectors {
		byName[collector.Get("name").String()] = collector
	}
	assert.Equal(t, int64(9216), byName["CLI"].Get("memory").Int())
	assert.True(t, byName["CLI"].Get("enabled").Bool())
	assert.Equal(t, int64(5120), byName["MDT"].Get("memory").Int())
	assert.False(t, byName["MDT"].Get("enabled").Bool(), "notDeployed is not enabled")
}

func collector(name string, memory int64, enabled any) DataGatewayResourceConfigurationCollectors {
	c := DataGatewayResourceConfigurationCollectors{Name: types.StringValue(name), Memory: types.Int64Null(), Enabled: types.BoolNull()}
	if memory != 0 {
		c.Memory = types.Int64Value(memory)
	}
	if e, ok := enabled.(bool); ok {
		c.Enabled = types.BoolValue(e)
	}
	return c
}

func TestResourcesRequest(t *testing.T) {
	current := gjson.Parse(containersAnswer)

	_, changed, err := resourcesRequest(current, []DataGatewayResourceConfigurationCollectors{collector("CLI", 9216, true), collector("MDT", 5120, false)})
	assert.True(t, err == nil, "no error")
	assert.False(t, changed, "the same values change nothing, and 5g is the same as 5120 MB")

	request, changed, err := resourcesRequest(current, []DataGatewayResourceConfigurationCollectors{collector("CLI", 9300, nil), collector("MDT", 0, true)})
	assert.True(t, err == nil, "no error")
	assert.True(t, changed)
	containers := gjson.Get(request, "containers").Array()
	assert.Len(t, containers, 2, "the whole list is sent")
	assert.Equal(t, "9300m", containers[0].Get("memoryLimit").String())
	assert.Equal(t, "x", containers[0].Get("imageTag").String(), "what the controller returned is sent back")
	assert.False(t, containers[1].Get("notDeployed").Bool())
	assert.Equal(t, "5g", containers[1].Get("memoryLimit").String(), "a memory the configuration leaves out is not touched")

	_, _, err = resourcesRequest(current, []DataGatewayResourceConfigurationCollectors{collector("NETCONF", 1000, nil)})
	assert.Contains(t, err.Error(), "no NETCONF collector")
}

func TestCollectorsAreReadByName(t *testing.T) {
	data := DataGatewayResourceConfiguration{Collectors: []DataGatewayResourceConfigurationCollectors{collector("MDT", 1, nil), collector("CLI", 1, true)}}
	data.updateFromBody(context.Background(), flattenCollectors(gjson.Parse(containersAnswer)))
	assert.Equal(t, int64(5120), data.Collectors[0].Memory.ValueInt64())
	assert.True(t, data.Collectors[0].Enabled.IsNull(), "a state the configuration leaves out stays unset")
	assert.Equal(t, int64(9216), data.Collectors[1].Memory.ValueInt64())
	assert.True(t, data.Collectors[1].Enabled.ValueBool())
}

func TestResourcesBusy(t *testing.T) {
	assert.True(t, resourcesBusy("StatusCode 429: Too many requests, please wait for the previous request"))
	assert.True(t, resourcesBusy("Too many requests"))
	assert.False(t, resourcesBusy("Validation Error: cli-collector memLimit 100m is invalid"))
}

func TestGlobalParametersHeld(t *testing.T) {
	current := flattenGlobalParameters(gjson.Parse(globalParametersAnswer))
	assert.True(t, globalParametersHeld(gjson.Parse(`{"cli_sessions":3,"snmp_trap_port":1062,"resync_snmpv3_details":false}`), current))
	assert.False(t, globalParametersHeld(gjson.Parse(`{"cli_sessions":4}`), current), "a value the controller does not hold")
	assert.False(t, globalParametersHeld(gjson.Parse(`{"syslog_udp_port":9514}`), current), "a value the controller does not list")
}

func TestChangedParametersKeepsOnlyTheValuesTheControllerDoesNotHold(t *testing.T) {
	flat := gjson.Parse(`{"cli_sessions":6,"snmp_trap_port":1063,"resync_snmpv3_details":true}`)
	current := gjson.Parse(`{"cli_sessions":6,"snmp_trap_port":1062,"resync_snmpv3_details":true,"syslog_udp_port":9514}`)
	changed := changedParameters(flat, current)
	assert.Equal(t, `{"snmp_trap_port":1063}`, changed.Raw)
	_, ok := globalParametersCommand(changedParameters(flat, flat))
	assert.False(t, ok, "nothing is sent when the controller holds every value")
}

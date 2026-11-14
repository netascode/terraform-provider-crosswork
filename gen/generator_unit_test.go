//go:build ignore

package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Run with: go test gen/generator.go gen/generator_unit_test.go
func TestValidatePostBody(t *testing.T) {
	tests := []struct {
		name    string
		config  YamlConfig
		wantErr bool
	}{
		{"unset", YamlConfig{}, false},
		{"wrapped under the id", YamlConfig{PostBody: `{"{id}": {body}}`, IdFromAttribute: true}, false},
		{"envelope without the id", YamlConfig{PostBody: `{"data": [{body}]}`}, false},
		{"missing body", YamlConfig{PostBody: `{"{id}": {}}`, IdFromAttribute: true}, true},
		{"id without id_from_attribute", YamlConfig{PostBody: `{"{id}": {body}}`}, true},
		{"not valid json", YamlConfig{PostBody: `{"{id}": {body}`, IdFromAttribute: true}, true},
		{"root list", YamlConfig{PostBody: `{"data": [{body}]}`, RootList: true}, true},
		{"chunked root list", YamlConfig{PostBody: `{"data": [{body}]}`, Attributes: []YamlConfigAttribute{{MaxElementsInRootList: 10}}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePostBody(tt.config); (err != nil) != tt.wantErr {
				t.Fatalf("validatePostBody() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGoEscape(t *testing.T) {
	tests := map[string]string{
		"plain text":         "plain text",
		`has "quotes"`:       `has \"quotes\"`,
		`back\slash`:         `back\\slash`,
		"two\nlines":         `two\nlines`,
		"keeps ©, — and `x`": "keeps ©, — and `x`",
	}
	for in, want := range tests {
		if got := goEscape(in); got != want {
			t.Errorf("goEscape(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestEscapeGoStringsCoversNestedAttributes(t *testing.T) {
	config := YamlConfig{
		ResDescription: `a "res"`,
		Attributes: []YamlConfigAttribute{{
			Description: `top "a"`,
			EnumValues:  []string{`x"y`},
			Attributes:  []YamlConfigAttribute{{Description: `nested "b"`, CoexistenceNote: `note "c"`}},
		}},
	}
	escapeGoStrings(&config)
	if config.ResDescription != `a \"res\"` || config.Attributes[0].Description != `top \"a\"` ||
		config.Attributes[0].EnumValues[0] != `x\"y` || config.Attributes[0].Attributes[0].Description != `nested \"b\"` ||
		config.Attributes[0].Attributes[0].CoexistenceNote != `note \"c\"` {
		t.Fatalf("not escaped: %+v", config)
	}
}

func TestSnakeFromModelName(t *testing.T) {
	tests := map[string]string{
		"profile":              "profile",
		"user_pass":            "user_pass",
		"enable_password_data": "enable_password_data",
		"v2_info":              "v2_info",
		"userName":             "user_name",
		"FirstName":            "first_name",
		"PolicyId":             "policy_id",
		"DeviceAccessGroups":   "device_access_groups",
		"Uuid":                 "uuid",
		"UUID":                 "uuid",
		"userID":               "user_id",
		"DNSServers":           "dns_servers",
		"ID":                   "id",
		"nodeIPAddress":        "node_ip_address",
		"Base64Value":          "base64_value",
		"":                     "",
	}
	for in, want := range tests {
		if got := snakeFromModelName(in); got != want {
			t.Errorf("snakeFromModelName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestValidatePlanKeys(t *testing.T) {
	ref := YamlConfigAttribute{ModelName: "name", Reference: true, Id: true}
	tests := []struct {
		name    string
		config  YamlConfig
		wantErr bool
	}{
		{"unset", YamlConfig{}, false},
		{"plan and one reference", YamlConfig{IdFromAttribute: true, PlanEndpoint: "/x/plan=%v", Attributes: []YamlConfigAttribute{ref}}, false},
		{"plan with a wait time", YamlConfig{IdFromAttribute: true, PlanEndpoint: "/x/plan=%v", MaxAsyncWaitTime: 60, Attributes: []YamlConfigAttribute{ref}}, false},
		{"plan with no_wait", YamlConfig{IdFromAttribute: true, PlanEndpoint: "/x/plan=%v", NoWait: true, Attributes: []YamlConfigAttribute{ref}}, false},
		{"wait time without a plan", YamlConfig{MaxAsyncWaitTime: 60}, true},
		{"no_wait without a plan", YamlConfig{NoWait: true}, true},
		{"no_wait and a wait time", YamlConfig{IdFromAttribute: true, PlanEndpoint: "/x/plan=%v", NoWait: true, MaxAsyncWaitTime: 60, Attributes: []YamlConfigAttribute{ref}}, true},
		{"negative wait time", YamlConfig{IdFromAttribute: true, PlanEndpoint: "/x/plan=%v", MaxAsyncWaitTime: -1, Attributes: []YamlConfigAttribute{ref}}, true},
		{"too few placeholders", YamlConfig{PlanEndpoint: "/x/plan", Attributes: []YamlConfigAttribute{ref}}, true},
		{"too many placeholders", YamlConfig{PlanEndpoint: "/x/%v/plan=%v", Attributes: []YamlConfigAttribute{ref}}, true},
		{"plan with a fallback endpoint", YamlConfig{IdFromAttribute: true, PlanEndpoint: "/x/plan=%v", FallbackRestEndpoint: "/y", Attributes: []YamlConfigAttribute{ref}}, true},
		{"plan with a root list", YamlConfig{IdFromAttribute: true, PlanEndpoint: "/x/plan=%v", RootList: true, Attributes: []YamlConfigAttribute{ref}}, true},
		{"plan with delete_iterative", YamlConfig{IdFromAttribute: true, PlanEndpoint: "/x/plan=%v", DeleteIterative: true, Attributes: []YamlConfigAttribute{ref}}, true},
		{"plan with a put delete", YamlConfig{IdFromAttribute: true, PlanEndpoint: "/x/plan=%v", PutDelete: true, Attributes: []YamlConfigAttribute{ref}}, true},
		{"plan with chunked root list", YamlConfig{IdFromAttribute: true, PlanEndpoint: "/x/plan=%v", Attributes: []YamlConfigAttribute{ref, {ModelName: "l", MaxElementsInRootList: 5}}}, true},
		{"a computed attribute", YamlConfig{IdFromAttribute: true, PlanEndpoint: "/x/plan=%v", Attributes: []YamlConfigAttribute{ref, {ModelName: "uuid", Computed: true}}}, true},
		{"the id is not from an attribute", YamlConfig{PlanEndpoint: "/x/plan=%v", Attributes: []YamlConfigAttribute{ref}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePlanKeys(tt.config); (err != nil) != tt.wantErr {
				t.Fatalf("validatePlanKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateAttributeKeys(t *testing.T) {
	tests := []struct {
		name    string
		attrs   []YamlConfigAttribute
		wantErr bool
	}{
		{"none", nil, false},
		{"key in the body", []YamlConfigAttribute{{ModelName: "name", Reference: true, ReferenceInBody: true}}, false},
		{"key in the body without reference", []YamlConfigAttribute{{ModelName: "name", ReferenceInBody: true}}, true},
		{"escape", []YamlConfigAttribute{{ModelName: "name", Reference: true, EscapeReferencePath: true}}, false},
		{"escape without reference", []YamlConfigAttribute{{ModelName: "name", EscapeReferencePath: true}}, true},
		{"empty", []YamlConfigAttribute{{ModelName: "pce", Type: "Bool", TypeYangBool: "empty"}}, false},
		{"presence", []YamlConfigAttribute{{ModelName: "pce", Type: "Bool", TypeYangBool: "presence"}}, false},
		{"boolean", []YamlConfigAttribute{{ModelName: "x", Type: "Bool", TypeYangBool: "boolean"}}, false},
		{"unknown encoding", []YamlConfigAttribute{{ModelName: "x", Type: "Bool", TypeYangBool: "flag"}}, true},
		{"encoding on a string", []YamlConfigAttribute{{ModelName: "x", Type: "String", TypeYangBool: "empty"}}, true},
		{"nested error is found", []YamlConfigAttribute{{ModelName: "list", Type: "List", Attributes: []YamlConfigAttribute{{ModelName: "x", Type: "Int64", TypeYangBool: "empty"}}}}, true},
		{"ignore extra values on a String Set", []YamlConfigAttribute{{ModelName: "x", Type: "Set", ElementType: "String", IgnoreExtraResponseValues: true}}, false},
		{"ignore extra values on a List", []YamlConfigAttribute{{ModelName: "x", Type: "List", ElementType: "String", IgnoreExtraResponseValues: true}}, true},
		{"write only tf native", []YamlConfigAttribute{{ModelName: "x", Type: "String", WriteOnly: true, WriteOnlyTFNative: true}}, false},
		{"write only tf native without write only", []YamlConfigAttribute{{ModelName: "x", Type: "String", WriteOnlyTFNative: true}}, true},
		{"write only tf native with write only tf", []YamlConfigAttribute{{ModelName: "x", Type: "String", WriteOnly: true, WriteOnlyTF: true, WriteOnlyTFNative: true}}, true},
		{"write only tf", []YamlConfigAttribute{{ModelName: "x", Type: "String", WriteOnly: true, WriteOnlyTF: true}}, false},
		{"write only tf on an Int64", []YamlConfigAttribute{{ModelName: "x", Type: "Int64", WriteOnly: true, WriteOnlyTF: true}}, true},
		{"write only tf without write only", []YamlConfigAttribute{{ModelName: "x", Type: "String", WriteOnlyTF: true}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAttributeKeys(tt.attrs, ""); (err != nil) != tt.wantErr {
				t.Fatalf("validateAttributeKeys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNsoDevicesCode(t *testing.T) {
	nested := []YamlConfigAttribute{
		{TfName: "vpn_nodes", Type: "List", Attributes: []YamlConfigAttribute{
			{TfName: "local_as", Type: "Int64"},
			{TfName: "vpn_node_id", Type: "String", ReferenceNsoDevice: true},
			{TfName: "accesses", Type: "Map", Attributes: []YamlConfigAttribute{{TfName: "peer", Type: "String", ReferenceNsoDevice: true}}},
		}},
		{TfName: "name", Type: "String"},
	}
	config := YamlConfig{Attributes: nested}
	if !HasNsoDevices(config) {
		t.Fatal("a flagged leaf in a nested list was not found")
	}
	want := "\tfor _, i1 := range data.VpnNodes {\n" +
		"\t\tnsoAddDevice(devices, i1.VpnNodeId)\n" +
		"\t\tfor _, i2 := range i1.Accesses {\n" +
		"\t\t\tnsoAddDevice(devices, i2.Peer)\n" +
		"\t\t}\n" +
		"\t}\n"
	if got := NsoDevicesCode(config); got != want {
		t.Errorf("NsoDevicesCode() =\n%s\nwant\n%s", got, want)
	}
	if got := nsoDevicesDescription(config); !strings.Contains(got, "`vpn_nodes.vpn_node_id`, `vpn_nodes.accesses.peer`") {
		t.Errorf("description %q does not name the attributes", got)
	}

	top := YamlConfig{Attributes: []YamlConfigAttribute{{TfName: "device", Type: "String", ReferenceNsoDevice: true}}}
	if got, want := NsoDevicesCode(top), "\tnsoAddDevice(devices, data.Device)\n"; got != want {
		t.Errorf("top level: %q, want %q", got, want)
	}

	none := YamlConfig{Attributes: []YamlConfigAttribute{{TfName: "name", Type: "String"}, {TfName: "items", Type: "List", Attributes: []YamlConfigAttribute{{TfName: "x", Type: "String"}}}}}
	if HasNsoDevices(none) || NsoDevicesCode(none) != "" || nsoDevicesDescription(none) != "" {
		t.Error("a definition without the flag must generate nothing")
	}
}

func TestReferenceNsoDeviceNeedsAString(t *testing.T) {
	for name, tc := range map[string]struct {
		attr    YamlConfigAttribute
		wantErr bool
	}{
		"string":       {YamlConfigAttribute{ModelName: "n", Type: "String", ReferenceNsoDevice: true}, false},
		"list":         {YamlConfigAttribute{ModelName: "n", Type: "List", ReferenceNsoDevice: true}, true},
		"fixed value":  {YamlConfigAttribute{ModelName: "n", Type: "String", Value: "x", ReferenceNsoDevice: true}, true},
		"write only":   {YamlConfigAttribute{ModelName: "n", Type: "String", WriteOnly: true, ReferenceNsoDevice: true}, true},
		"not a device": {YamlConfigAttribute{ModelName: "n", Type: "Int64"}, false},
	} {
		t.Run(name, func(t *testing.T) {
			if err := validateAttributeKeys([]YamlConfigAttribute{tc.attr}, ""); (err != nil) != tc.wantErr {
				t.Fatalf("validateAttributeKeys() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestConstantId(t *testing.T) {
	for name, tc := range map[string]struct {
		config YamlConfig
		want   string
	}{
		"literal":            {YamlConfig{IdPath: `!\"alarm_settings\"`}, "alarm_settings"},
		"path":               {YamlConfig{IdPath: "group.uuid"}, ""},
		"rule path":          {YamlConfig{IdPath: "rule.uuid"}, ""},
		"empty":              {YamlConfig{}, ""},
		"id from attribute":  {YamlConfig{IdPath: `!\"x\"`, IdFromAttribute: true}, ""},
		"embedded quote":     {YamlConfig{IdPath: `!\"a"b\"`}, ""},
		"embedded backslash": {YamlConfig{IdPath: `!\"a\b\"`}, ""},
	} {
		t.Run(name, func(t *testing.T) {
			if got := ConstantId(tc.config); got != tc.want {
				t.Errorf("ConstantId() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestImportVerifyExcludes(t *testing.T) {
	config := YamlConfig{Attributes: []YamlConfigAttribute{
		{ModelName: "name", Type: "String"},
		{ModelName: "secret", Type: "String", WriteOnly: true, WriteOnlyTF: true},
		{ModelName: "token", Type: "String", WriteOnly: true, WriteOnlyTF: true, RequiresReplace: true},
		{ModelName: "rows", Type: "Map", AllowImportChanges: true, Attributes: []YamlConfigAttribute{
			{ModelName: "inner", Type: "String", WriteOnly: true},
		}},
		{ModelName: "credentials", Type: "List", Attributes: []YamlConfigAttribute{
			{ModelName: "access", Type: "String"},
			{ModelName: "password", Type: "String", WriteOnly: true, WriteOnlyTF: true},
			{ModelName: "encoding", Type: "String", AllowImportChanges: true},
		}},
		{ModelName: "plain", Type: "List", Attributes: []YamlConfigAttribute{{ModelName: "x", Type: "String"}}},
	}}
	augmentConfig(&config)
	got, err := ImportVerifyExcludes(config.Attributes)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"secret", "secret_wo", "secret_wo_version",
		"token", "token_wo",
		"rows",
		"credentials.0.password", "credentials.0.password_wo", "credentials.0.password_wo_version",
		"credentials.0.encoding",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ImportVerifyExcludes() = %v, want %v", got, want)
	}
	for _, e := range got {
		if e == "" {
			t.Error("an empty entry would ignore every attribute")
		}
	}

	got, err = ImportVerifyExcludes([]YamlConfigAttribute{{TfName: "name", Type: "String"}})
	if err != nil || len(got) != 0 {
		t.Errorf("no flagged attribute: got %v, %v", got, err)
	}

	for _, typ := range []string{"Set", "Map"} {
		bad := []YamlConfigAttribute{{TfName: "rows", Type: typ, Attributes: []YamlConfigAttribute{{TfName: "pw", Type: "String", WriteOnly: true}}}}
		if _, err := ImportVerifyExcludes(bad); err == nil {
			t.Errorf("a nested secret in a %s must be an error", typ)
		}
		bad[0].Attributes[0] = YamlConfigAttribute{TfName: "x", Type: "String", AllowImportChanges: true}
		if _, err := ImportVerifyExcludes(bad); err == nil {
			t.Errorf("a nested allow_import_changes in a %s must be an error", typ)
		}
	}
}

func TestImportTemplatesForConstantId(t *testing.T) {
	root, err := filepath.Abs(".")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "gen/templates/model.go")); err != nil {
		root = filepath.Dir(root)
	}
	newConfig := func(idPath string, keyed bool) YamlConfig {
		config := YamlConfig{Name: "Import Fixture", IdPath: idPath, IdFromAttribute: keyed, DataSourceNoId: !keyed, Attributes: []YamlConfigAttribute{
			{ModelName: "name", Type: "String", Id: keyed, Example: "x"},
			{ModelName: "secret", Type: "String", WriteOnly: true, WriteOnlyTF: true, Example: "s"},
			{ModelName: "rows", Type: "Map", AllowImportChanges: true, Attributes: []YamlConfigAttribute{{ModelName: "v", Type: "String", Example: "v"}}},
		}}
		augmentConfig(&config)
		return config
	}
	render := func(name string, config YamlConfig) string {
		out := filepath.Join(t.TempDir(), name)
		renderTemplate(filepath.Join(root, "gen/templates", name), out, config)
		b, err := os.ReadFile(out)
		if err != nil {
			t.Fatal(err)
		}
		return string(b)
	}

	constant := newConfig(`!\"import_fixture\"`, false)
	if got := strings.TrimSpace(render("import.sh", constant)); got != `terraform import crosswork_import_fixture.example "import_fixture"` {
		t.Errorf("import.sh = %q", got)
	}
	res := getTemplateSection(render("resource.go", constant), "import")
	for _, want := range []string{`req.ID != "import_fixture"`, "Unexpected Import Identifier", "Expected import identifier with format: import_fixture. Got: %q", `path.Root("id"), req.ID`} {
		if !strings.Contains(res, want) {
			t.Errorf("ImportState lacks %q:\n%s", want, res)
		}
	}
	if strings.Contains(res, "strings.Split") || strings.Contains(res, "ImportStatePassthroughID") {
		t.Errorf("ImportState of a constant id must not use the keyed or passthrough form:\n%s", res)
	}
	tst := render("resource_test.go", constant)
	for _, want := range []string{`ImportStateId: "import_fixture"`, "ImportStateVerify: true", `"secret", "secret_wo", "secret_wo_version", "rows"`} {
		if !strings.Contains(strings.Join(strings.Fields(tst), " "), strings.Join(strings.Fields(want), " ")) {
			t.Errorf("test step lacks %q", want)
		}
	}
	model := getTemplateSection(render("model.go", constant), "fromBody")
	if !strings.Contains(model, `data.Id = types.StringValue("import_fixture")`) || strings.Contains(model, `res.Get("")`) {
		t.Errorf("fromBody must set Id to the constant:\n%s", model)
	}

	keyed := newConfig("", true)
	if got := strings.TrimSpace(render("import.sh", keyed)); got != `terraform import crosswork_import_fixture.example "<name>"` {
		t.Errorf("keyed import.sh = %q", got)
	}
	if res := getTemplateSection(render("resource.go", keyed), "import"); !strings.Contains(res, "strings.Split") {
		t.Errorf("keyed ImportState lost its form:\n%s", res)
	}

	plain := newConfig("group.uuid", false)
	if got := strings.TrimSpace(render("import.sh", plain)); got != `terraform import crosswork_import_fixture.example "<id>"` {
		t.Errorf("an id path without a constant must keep the <id> placeholder, got %q", got)
	}
	noId := newConfig("", false)
	noId.ImportNoId = true
	if res := getTemplateSection(render("resource.go", noId), "import"); !strings.Contains(res, "ImportStatePassthroughID") {
		t.Errorf("import_no_id must keep the passthrough form:\n%s", res)
	}
}

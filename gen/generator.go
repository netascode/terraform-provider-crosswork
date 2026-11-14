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

//go:build ignore

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"gopkg.in/yaml.v3"
)

const (
	definitionsPath   = "./gen/definitions/"
	providerTemplate  = "./gen/templates/provider.go"
	providerLocation  = "./internal/provider/provider.go"
	changelogTemplate = "./gen/templates/changelog.md.tmpl"
	changelogLocation = "./templates/guides/changelog.md.tmpl"
	changelogOriginal = "./CHANGELOG.md"
)

type t struct {
	path   string
	prefix string
	suffix string
}

var templates = []t{
	{
		path:   "./gen/templates/model.go",
		prefix: "./internal/provider/model_crosswork_",
		suffix: ".go",
	},
	{
		path:   "./gen/templates/data_source.go",
		prefix: "./internal/provider/data_source_crosswork_",
		suffix: ".go",
	},
	{
		path:   "./gen/templates/data_source_test.go",
		prefix: "./internal/provider/data_source_crosswork_",
		suffix: "_test.go",
	},
	{
		path:   "./gen/templates/resource.go",
		prefix: "./internal/provider/resource_crosswork_",
		suffix: ".go",
	},
	{
		path:   "./gen/templates/resource_test.go",
		prefix: "./internal/provider/resource_crosswork_",
		suffix: "_test.go",
	},
	{
		path:   "./gen/templates/data-source.tf",
		prefix: "./examples/data-sources/crosswork_",
		suffix: "/data-source.tf",
	},
	{
		path:   "./gen/templates/resource.tf",
		prefix: "./examples/resources/crosswork_",
		suffix: "/resource.tf",
	},
	{
		path:   "./gen/templates/import.sh",
		prefix: "./examples/resources/crosswork_",
		suffix: "/import.sh",
	},
}

type YamlConfig struct {
	Name                          string                `yaml:"name"`
	NoResource                    bool                  `yaml:"no_resource"`
	NoDataSource                  bool                  `yaml:"no_data_source"`
	RestEndpoint                  string                `yaml:"rest_endpoint"`
	GetRestEndpoint               string                `yaml:"get_rest_endpoint"`
	PostRestEndpoint              string                `yaml:"post_rest_endpoint"`
	PutRestEndpoint               string                `yaml:"put_rest_endpoint"`
	DeleteRestEndpoint            string                `yaml:"delete_rest_endpoint"`
	FallbackRestEndpoint          string                `yaml:"fallback_rest_endpoint"`
	GetNoId                       bool                  `yaml:"get_no_id"`
	GetFromAll                    bool                  `yaml:"get_from_all"`
	GetRequiresId                 bool                  `yaml:"get_requires_id"`
	GetExtraQueryParams           string                `yaml:"get_extra_query_params"`
	NoDelete                      bool                  `yaml:"no_delete"`
	DataSourceNoId                bool                  `yaml:"data_source_no_id"`
	DataSourceNoInputs            bool                  `yaml:"data_source_no_inputs"`
	DeleteNoId                    bool                  `yaml:"delete_no_id"`
	NoUpdate                      bool                  `yaml:"no_update"`
	NoRead                        bool                  `yaml:"no_read"`
	NoImport                      bool                  `yaml:"no_import"`
	ImportNoId                    bool                  `yaml:"import_no_id"`
	PostUpdate                    bool                  `yaml:"post_update"`
	PutCreate                     bool                  `yaml:"put_create"`
	PutDelete                     bool                  `yaml:"put_delete"`
	PostDelete                    bool                  `yaml:"post_delete"`
	UpdateComputed                bool                  `yaml:"update_computed"`
	DeleteIterative               bool                  `yaml:"delete_iterative"`
	RootList                      bool                  `yaml:"root_list"`
	NoReadPrefix                  bool                  `yaml:"no_read_prefix"`
	NoWait                        bool                  `yaml:"no_wait"`
	IdPath                        string                `yaml:"id_path"`
	IdFromQueryPath               string                `yaml:"id_from_query_path"`
	IdFromQueryPathAttribute      string                `yaml:"id_from_query_path_attribute"`
	IdQueryRestEndpoint           string                `yaml:"id_query_rest_endpoint"`
	IdQueryParam                  string                `yaml:"id_query_param"`
	IdFromAttribute               bool                  `yaml:"id_from_attribute"`
	DeviceUnreachabilityWarning   bool                  `yaml:"device_unreachability_warning"`
	AllowExistingOnCreate         bool                  `yaml:"allow_existing_on_create"`
	SkipUnmatchedRows             bool                  `yaml:"skip_unmatched_rows"`
	RemoveResourceOnEmptyResponse bool                  `yaml:"remove_resource_on_empty_response"`
	PutIdIncludePath              string                `yaml:"put_id_include_path"`
	PutIdQueryParam               string                `yaml:"put_id_query_param"`
	PutNoId                       bool                  `yaml:"put_no_id"`
	PutId                         bool                  `yaml:"put_id"`
	PutUpdateId                   bool                  `yaml:"put_update_id"`
	DeleteIdQueryParam            string                `yaml:"delete_id_query_param"`
	DeleteBody                    string                `yaml:"delete_body"`
	MinimumVersion                string                `yaml:"minimum_version"`
	DsDescription                 string                `yaml:"ds_description"`
	ResDescription                string                `yaml:"res_description"`
	DocCategory                   string                `yaml:"doc_category"`
	ExcludeTest                   bool                  `yaml:"exclude_test"`
	SkipMinimumTest               bool                  `yaml:"skip_minimum_test"`
	Attributes                    []YamlConfigAttribute `yaml:"attributes"`
	TestTags                      []string              `yaml:"test_tags"`
	TestPrerequisites             string                `yaml:"test_prerequisites"`
	MaxAsyncWaitTime              int64                 `yaml:"max_async_wait_time"`
	RetryOnErrorCodes             []string              `yaml:"retry_on_error_codes"`
	Mutex                         bool                  `yaml:"mutex"`
	UseCache                      bool                  `yaml:"use_cache"`
	CacheRestEndpoint             string                `yaml:"cache_rest_endpoint"`
	CacheFilterAttributes         []string              `yaml:"cache_filter_attributes"`
	SkipDeleteOnEmptyId           bool                  `yaml:"skip_delete_on_empty_id"`
	FailCreateOnEmptyId           bool                  `yaml:"fail_create_on_empty_id"`
	GetBeforeDelete               bool                  `yaml:"get_before_delete"`
	GetBody                       string                `yaml:"get_body"`
	PostBody                      string                `yaml:"post_body"`
	PlanEndpoint                  string                `yaml:"plan_endpoint"`
}

type YamlConfigAttribute struct {
	ModelName                 string                `yaml:"model_name"`
	ResponseModelName         string                `yaml:"response_model_name"`
	FallbackResponseModelName string                `yaml:"fallback_response_model_name"`
	TfName                    string                `yaml:"tf_name"`
	Type                      string                `yaml:"type"`
	ElementType               string                `yaml:"element_type"`
	DataPath                  string                `yaml:"data_path"`
	ResponseDataPath          string                `yaml:"response_data_path"`
	FallbackResponseDataPath  string                `yaml:"fallback_response_data_path"`
	PutDataPath               string                `yaml:"put_data_path"`
	Id                        bool                  `yaml:"id"`
	MatchId                   bool                  `yaml:"match_id"`
	Reference                 bool                  `yaml:"reference"`
	ReferenceInBody           bool                  `yaml:"reference_in_body"`
	ReferenceNsoDevice        bool                  `yaml:"reference_nso_device"`
	EscapeReferencePath       bool                  `yaml:"escape_reference_path"`
	TypeYangBool              string                `yaml:"type_yang_bool"`
	RequiresReplace           bool                  `yaml:"requires_replace"`
	QueryParam                bool                  `yaml:"query_param"`
	DeleteQueryParam          bool                  `yaml:"delete_query_param"`
	GetQueryParam             bool                  `yaml:"get_query_param"`
	PutQueryParam             bool                  `yaml:"put_query_param"`
	PostQueryParam            bool                  `yaml:"post_query_param"`
	QueryParamName            string                `yaml:"query_param_name"`
	DeleteQueryParamName      string                `yaml:"delete_query_param_name"`
	GetQueryParamName         string                `yaml:"get_query_param_name"`
	PutQueryParamName         string                `yaml:"put_query_param_name"`
	PostQueryParamName        string                `yaml:"post_query_param_name"`
	CreateQueryPath           bool                  `yaml:"create_query_path"`
	DataSourceQuery           bool                  `yaml:"data_source_query"`
	QueryParamNoBody          bool                  `yaml:"query_param_no_body"`
	Mandatory                 bool                  `yaml:"mandatory"`
	Computed                  bool                  `yaml:"computed"`
	ComputedRefreshValue      bool                  `yaml:"computed_refresh_value"`
	NoUseStateForUnknown      bool                  `yaml:"no_use_state_for_unknown"`
	WriteOnly                 bool                  `yaml:"write_only"`
	WriteOnlyTF               bool                  `yaml:"write_only_tf"`
	ModelTypeString           bool                  `yaml:"model_type_string"`
	WriteOnlyTFNative         bool                  `yaml:"write_only_tf_native"` // Native schema write-only input without legacy or generated version companions.
	WoVersion                 bool                  `yaml:"-"`                    // Internal: marks a generated "<attr>_wo_version" companion attribute (state-only rotation trigger)
	CoexistingSecret          bool                  `yaml:"-"`                    // Internal: marks the legacy state-storing twin of a "<attr>_wo" write-only attribute
	WoBaseName                string                `yaml:"-"`                    // Internal: on a "<attr>_wo" attribute, the name of its legacy twin
	MutualExclusivityNote     string                `yaml:"-"`                    // Internal: documentation note carried by both halves of a write_only_tf pair
	WoPairMandatory           bool                  `yaml:"-"`                    // Internal: on a legacy twin, whether the pair must supply the secret through one of its halves
	WoPairHasVersion          bool                  `yaml:"-"`                    // Internal: on a legacy twin, whether a "_wo_version" companion was generated
	CoexistenceNote           string                `yaml:"coexistence_note"`
	ExcludeFromPut            bool                  `yaml:"exclude_from_put"`
	ExcludeTest               bool                  `yaml:"exclude_test"`
	AllowImportChanges        bool                  `yaml:"allow_import_changes"`
	ExcludeExample            bool                  `yaml:"exclude_example"`
	Description               string                `yaml:"description"`
	Example                   string                `yaml:"example"`
	EnumValues                []string              `yaml:"enum_values"`
	MinList                   int64                 `yaml:"min_list"`
	MaxList                   int64                 `yaml:"max_list"`
	MinInt                    int64                 `yaml:"min_int"`
	MaxInt                    int64                 `yaml:"max_int"`
	MinFloat                  float64               `yaml:"min_float"`
	MaxFloat                  float64               `yaml:"max_float"`
	DecimalPlaces             int64                 `yaml:"decimal_places"`
	StringPatterns            []string              `yaml:"string_patterns"`
	StringMinLength           int64                 `yaml:"string_min_length"`
	StringMaxLength           int64                 `yaml:"string_max_length"`
	DefaultValue              string                `yaml:"default_value"`
	Value                     string                `yaml:"value"`
	ValueCondition            string                `yaml:"value_condition"`
	TestValue                 string                `yaml:"test_value"`
	MinimumTestValue          string                `yaml:"minimum_test_value"`
	TestTags                  []string              `yaml:"test_tags"`
	Attributes                []YamlConfigAttribute `yaml:"attributes"`
	MaxElementsInRootList     int64                 `yaml:"max_elements_in_root_list"`
	MapKeyExample             string                `yaml:"map_key_example"`
	KeyPart                   bool                  `yaml:"key_part"`
	NoPut                     bool                  `yaml:"no_put"`
	NullOnEmpty               bool                  `yaml:"null_on_empty"`
	CustomModifier            string                `yaml:"custom_modifier"`
	AlwaysInclude             bool                  `yaml:"always_include"`
	IgnoreExtraResponseValues bool                  `yaml:"ignore_extra_response_values"`
	GetIfUnsetOnUpdate        bool                  `yaml:"get_if_unset_on_update"`
}

func ToGoName(s string) string {
	var g []string

	p := strings.Split(s, "_")

	for _, value := range p {
		g = append(g, strings.Title(value))
	}
	s = strings.Join(g, "")
	return s
}

func CamelCase(s string) string {
	var g []string

	s = strings.ReplaceAll(s, "-", " ")
	p := strings.Fields(s)

	for _, value := range p {
		g = append(g, strings.Title(value))
	}
	return strings.Join(g, "")
}

func SnakeCase(s string) string {
	var g []string

	s = strings.ReplaceAll(s, "-", " ")
	p := strings.Fields(s)

	for _, value := range p {
		g = append(g, strings.ToLower(value))
	}
	return strings.Join(g, "_")
}

func BuildPath(s []string) string {
	return strings.Join(s, ".")
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func HasId(attributes []YamlConfigAttribute) bool {
	for _, attr := range attributes {
		if attr.Id {
			return true
		}
	}
	return false
}

// The top-level member of the response body: the first segment of response_data_path, or of data_path when that is unset.
// Falls back to "response" when the first segment is empty or a list index.
func ResponseEnvelope(attributes []YamlConfigAttribute) string {
	for _, attr := range attributes {
		path := attr.ResponseDataPath
		if path == "" {
			path = attr.DataPath
		}
		if path == "" {
			continue
		}
		first := strings.Split(path, ".")[0]
		if _, err := strconv.Atoi(first); err == nil {
			continue
		}
		return first
	}
	return "response"
}

func ReadGet(config YamlConfig, receiver string, ptr bool) string {
	if config.GetBody == "" {
		return "c.Get("
	}
	if ptr {
		return receiver + ".readGet(c, "
	}
	return receiver + ".readGet(&c, "
}

func HasReference(attributes []YamlConfigAttribute) bool {
	for _, attr := range attributes {
		if attr.Reference {
			return true
		}
	}
	return false
}

func HasQueryParam(attributes []YamlConfigAttribute) bool {
	for _, attr := range attributes {
		if attr.QueryParam {
			return true
		}
	}
	return false
}

func HasGetQueryParam(attributes []YamlConfigAttribute) bool {
	for _, attr := range attributes {
		if attr.GetQueryParam {
			return true
		}
	}
	return false
}

func HasDeleteQueryParam(attributes []YamlConfigAttribute) bool {
	for _, attr := range attributes {
		if attr.DeleteQueryParam {
			return true
		}
	}
	return false
}

func HasComputedRefreshValue(attributes []YamlConfigAttribute) bool {
	for _, attr := range attributes {
		if attr.ComputedRefreshValue {
			return true
		}
		if len(attr.Attributes) > 0 {
			if HasComputedRefreshValue(attr.Attributes) {
				return true
			}
		}
	}
	return false
}

func HasComputedRefreshValueInKeys(attributes []YamlConfigAttribute) bool {
	noId := !HasId(attributes)
	for _, attr := range attributes {
		if attr.ComputedRefreshValue && (attr.Id || noId) {
			return true
		}
	}
	return false
}

// A nested collection whose children are all plain attributes contributes no statements there,
// so its loop must be skipped or the generated code declares `r` without using it.
func HasUnknownsChildren(attributes []YamlConfigAttribute) bool {
	for _, attr := range attributes {
		if attr.Value == "" && !attr.WriteOnly && !attr.Reference && attr.Computed {
			return true
		}
	}
	return false
}

// Values the controller assigns itself are null in the plan; Update fetches each one and
// puts it in the PUT body when the plan leaves it empty, so a full replace does not strip it.
func GetIfUnsetOnUpdateAttributes(config YamlConfig) []YamlConfigAttribute {
	r := []YamlConfigAttribute{}
	for _, attr := range config.Attributes {
		if attr.GetIfUnsetOnUpdate {
			r = append(r, attr)
		}
	}
	return r
}

func GetId(attributes []YamlConfigAttribute) YamlConfigAttribute {
	for _, attr := range attributes {
		if attr.Id {
			return attr
		}
	}
	return YamlConfigAttribute{}
}

func GetMatchId(attributes []YamlConfigAttribute) YamlConfigAttribute {
	for _, attr := range attributes {
		if attr.MatchId {
			return attr
		}
	}
	return YamlConfigAttribute{}
}

// Nested match_id attributes identify rows within a list, not the resource.
func HasMatchId(attributes []YamlConfigAttribute) bool {
	for _, attr := range attributes {
		if attr.MatchId {
			return true
		}
	}
	return false
}

func CacheFilterGjsonExpr(attributes []YamlConfigAttribute, filterNames []string) string {
	if len(filterNames) == 0 {
		return ""
	}
	tfNameByModel := make(map[string]string)
	for _, attr := range attributes {
		tfNameByModel[attr.ModelName] = attr.TfName
	}

	parts := make([]string, 0, len(filterNames))
	for _, modelName := range filterNames {
		tfName := tfNameByModel[modelName]
		goName := ToGoName(tfName)
		parts = append(parts, modelName+`==\"" + state.`+goName+`.ValueString() + "\"`)
	}

	if len(parts) == 1 {
		return `"response.#(` + parts[0] + `)"`
	}
	expr := `"response.#(` + parts[0] + `)#`
	for _, p := range parts[1:] {
		expr += `|#(` + p + `)#`
	}
	expr = expr[:len(expr)-1]
	expr += `"`
	return expr
}

func GetQueryParam(attributes []YamlConfigAttribute) YamlConfigAttribute {
	for _, attr := range attributes {
		if attr.QueryParam {
			return attr
		}
	}
	return YamlConfigAttribute{}
}

func GetDeleteQueryParam(attributes []YamlConfigAttribute) YamlConfigAttribute {
	for _, attr := range attributes {
		if attr.DeleteQueryParam {
			return attr
		}
	}
	return YamlConfigAttribute{}
}

func HasDataSourceQuery(attributes []YamlConfigAttribute) bool {
	for _, attr := range attributes {
		if attr.DataSourceQuery {
			return true
		}
	}
	return false
}

func FirstPathElement(path, getFromAllPath string) string {
	if getFromAllPath != "" {
		return getFromAllPath
	} else if strings.HasPrefix(path, "response.") {
		return strings.Split(path, ".")[0]
	}
	return ""
}

func RemainingPathElements(path, getFromAllPath string) string {
	if getFromAllPath != "" {
		return path
	} else if strings.HasPrefix(path, "response.0.") {
		return strings.Join(strings.Split(path, ".")[2:], ".")
	} else if strings.HasPrefix(path, "response.") {
		return strings.Join(strings.Split(path, ".")[1:], ".")
	}
	return path
}

func GetFromAllPath(config YamlConfig) string {
	if config.GetFromAll {
		if config.IdFromQueryPath != "" {
			return config.IdFromQueryPath
		}
	}
	return ""
}

func IsListSet(attribute YamlConfigAttribute) bool {
	if (attribute.Type == "List" || attribute.Type == "Set") && attribute.ElementType != "" {
		return true
	}
	return false
}

func IsList(attribute YamlConfigAttribute) bool {
	if attribute.Type == "List" && attribute.ElementType != "" {
		return true
	}
	return false
}

func IsSet(attribute YamlConfigAttribute) bool {
	if attribute.Type == "Set" && attribute.ElementType != "" {
		return true
	}
	return false
}

// IsYangFlag reports whether a Bool attribute stands for a YANG empty leaf or presence container:
// the node is sent when the value is true, left out otherwise, and read as true when it exists.
func IsYangFlag(attribute YamlConfigAttribute) bool {
	return attribute.Type == "Bool" && (attribute.TypeYangBool == "empty" || attribute.TypeYangBool == "presence")
}

func IsStringListSet(attribute YamlConfigAttribute) bool {
	if (attribute.Type == "List" || attribute.Type == "Set") && attribute.ElementType == "String" {
		return true
	}
	return false
}

func IsInt64ListSet(attribute YamlConfigAttribute) bool {
	if (attribute.Type == "List" || attribute.Type == "Set") && attribute.ElementType == "Int64" {
		return true
	}
	return false
}

func IsNestedListSet(attribute YamlConfigAttribute) bool {
	if (attribute.Type == "List" || attribute.Type == "Set") && attribute.ElementType == "" {
		return true
	}
	return false
}

func IsNestedList(attribute YamlConfigAttribute) bool {
	if attribute.Type == "List" && attribute.ElementType == "" {
		return true
	}
	return false
}

func IsNestedSet(attribute YamlConfigAttribute) bool {
	if attribute.Type == "Set" && attribute.ElementType == "" {
		return true
	}
	return false
}

func IsNestedMap(attribute YamlConfigAttribute) bool {
	if attribute.Type == "Map" && attribute.ElementType == "" && len(attribute.Attributes) > 0 {
		return true
	}
	return false
}

func IsNestedListSetMap(attribute YamlConfigAttribute) bool {
	return IsNestedListSet(attribute) || IsNestedMap(attribute)
}

func HasKeyPart(attributes []YamlConfigAttribute) bool {
	for _, attr := range attributes {
		if attr.KeyPart {
			return true
		}
	}
	return false
}

func GetKeyPartAttributes(attributes []YamlConfigAttribute) []YamlConfigAttribute {
	var result []YamlConfigAttribute
	for _, attr := range attributes {
		if attr.KeyPart {
			result = append(result, attr)
		}
	}
	return result
}

func HasCreateQueryPath(attributes []YamlConfigAttribute) bool {
	for _, attr := range attributes {
		if attr.CreateQueryPath {
			return true
		}
	}
	return false
}

func GetCreateQueryPath(attributes []YamlConfigAttribute) YamlConfigAttribute {
	for _, attr := range attributes {
		if attr.CreateQueryPath {
			return attr
		}
	}
	return YamlConfigAttribute{}
}

func GenerateQueryParamString(method string, inputSource string, attributes []YamlConfigAttribute) string {
	var params []string
	first := true

	for _, attr := range attributes {
		var queryParamName string
		includeParam := false

		switch method {
		case "DELETE":
			if attr.DeleteQueryParam {
				queryParamName = attr.DeleteQueryParamName
				includeParam = true
			}
		case "GET":
			if attr.GetQueryParam {
				queryParamName = attr.GetQueryParamName
				includeParam = true
			}
		case "POST":
			if attr.PostQueryParam {
				queryParamName = attr.PostQueryParamName
				includeParam = true
			}
		case "PUT":
			if attr.PutQueryParam {
				queryParamName = attr.PutQueryParamName
				includeParam = true
			}
		}

		if !includeParam && attr.QueryParam {
			queryParamName = attr.QueryParamName
			includeParam = true
		}

		if queryParamName == "" {
			queryParamName = attr.ModelName
		}

		if includeParam {
			var valueExpr string
			switch attr.Type {
			case "Int64":
				valueExpr = `strconv.FormatInt(` + inputSource + `.` + ToGoName(attr.TfName) + `.Value` + attr.Type + `(), 10)`
			case "Bool":
				valueExpr = `strconv.FormatBool(` + inputSource + `.` + ToGoName(attr.TfName) + `.Value` + attr.Type + `())`
			default:
				valueExpr = inputSource + `.` + ToGoName(attr.TfName) + `.Value` + attr.Type + `()`
			}

			if first {
				params = append(params, `"?`+queryParamName+`=" + url.QueryEscape(`+valueExpr+`)`)
				first = false
			} else {
				params = append(params, `"&`+queryParamName+`=" + url.QueryEscape(`+valueExpr+`)`)
			}
		}
	}

	if len(params) == 0 {
		return ""
	} else {
		return strings.Join(params, "+")
	}
}

func GenerateDeleteOnlyQueryParamString(inputSource string, attributes []YamlConfigAttribute) string {
	var params []string
	first := true

	for _, attr := range attributes {
		if !attr.DeleteQueryParam {
			continue
		}

		queryParamName := attr.DeleteQueryParamName
		if queryParamName == "" {
			queryParamName = attr.ModelName
		}

		var valueExpr string
		switch attr.Type {
		case "Int64":
			valueExpr = `strconv.FormatInt(` + inputSource + `.` + ToGoName(attr.TfName) + `.Value` + attr.Type + `(), 10)`
		case "Bool":
			valueExpr = `strconv.FormatBool(` + inputSource + `.` + ToGoName(attr.TfName) + `.Value` + attr.Type + `())`
		default:
			valueExpr = inputSource + `.` + ToGoName(attr.TfName) + `.Value` + attr.Type + `()`
		}

		if first {
			params = append(params, `"?`+queryParamName+`=" + url.QueryEscape(`+valueExpr+`)`)
			first = false
		} else {
			params = append(params, `"&`+queryParamName+`=" + url.QueryEscape(`+valueExpr+`)`)
		}
	}

	if len(params) == 0 {
		return ""
	}
	return strings.Join(params, "+")
}

func ImportAttributes(config YamlConfig) []YamlConfigAttribute {
	r := []YamlConfigAttribute{}
	for _, attr := range config.Attributes {
		if attr.Reference || attr.QueryParam || attr.GetQueryParam || attr.Id {
			r = append(r, attr)
		}
	}
	if !config.IdFromAttribute && !config.ImportNoId {
		attr := YamlConfigAttribute{}
		attr.ModelName = "id"
		attr.TfName = "id"
		attr.Type = "String"
		attr.Example = "4b0b7a80-44c0-4bf2-bab5-fc24b4e0a17e"
		r = append(r, attr)
	}
	return r
}

var constantIdRegexp = regexp.MustCompile(`^!\\"(.+)\\"$`)

// ConstantId returns the constant id of a singleton whose id_path is a gjson literal such as
// !\"alarm_settings\" (the value holds literal backslash-quote pairs), else an empty string.
// A definition with id_from_attribute is keyed and never constant.
func ConstantId(config YamlConfig) string {
	if config.IdFromAttribute {
		return ""
	}
	m := constantIdRegexp.FindStringSubmatch(config.IdPath)
	if m == nil || strings.ContainsAny(m[1], "\"\\") {
		return ""
	}
	return m[1]
}

// ImportVerifyExcludes returns the ImportStateVerifyIgnore entries of the generated import step:
// every attribute that is write-only, a write-only version companion or marked
// allow_import_changes. The test library matches by prefix, so a name also covers its children.
// Flagged attributes nested in a List get the index form (credentials.0.password); a flag nested
// in a Set or Map that is not covered by a flagged parent is an error.
func ImportVerifyExcludes(attributes []YamlConfigAttribute) ([]string, error) {
	return importVerifyExcludes(attributes, "")
}

func importVerifyFlagged(attr YamlConfigAttribute) bool {
	return attr.WriteOnly || attr.WoVersion || attr.AllowImportChanges
}

func importVerifyAnyFlagged(attrs []YamlConfigAttribute) bool {
	for _, a := range attrs {
		if importVerifyFlagged(a) || importVerifyAnyFlagged(a.Attributes) {
			return true
		}
	}
	return false
}

func importVerifyExcludes(attrs []YamlConfigAttribute, prefix string) ([]string, error) {
	r := []string{}
	for _, attr := range attrs {
		name := prefix + attr.TfName
		if importVerifyFlagged(attr) {
			r = append(r, name)
			continue
		}
		if len(attr.Attributes) == 0 {
			continue
		}
		switch attr.Type {
		case "List":
			sub, err := importVerifyExcludes(attr.Attributes, name+".0.")
			if err != nil {
				return nil, err
			}
			r = append(r, sub...)
		default:
			if importVerifyAnyFlagged(attr.Attributes) {
				return nil, fmt.Errorf("attribute %q is a %s with nested write-only or allow_import_changes attributes; mark %q itself with allow_import_changes", name, attr.Type, name)
			}
		}
	}
	return r, nil
}

// Rounds to decimal_places to prevent precision drift (an API returning 5 decimals against a
// 6-decimal config).
func FloatReadExpr(varName string, attr YamlConfigAttribute) string {
	if attr.DecimalPlaces > 0 {
		factor := "1"
		for i := int64(0); i < attr.DecimalPlaces; i++ {
			factor += "0"
		}
		return fmt.Sprintf("math.Round(%s.Float()*%s) / %s", varName, factor, factor)
	}
	return varName + ".Float()"
}

func Subtract(a, b int) int {
	return a - b
}

func HasWriteOnlyTFChildren(attr YamlConfigAttribute) bool {
	for _, child := range attr.Attributes {
		if child.WriteOnlyTF {
			return true
		}
	}
	return false
}

func WriteOnlyTFChildren(attr YamlConfigAttribute) []YamlConfigAttribute {
	r := []YamlConfigAttribute{}
	for _, child := range attr.Attributes {
		if child.WriteOnlyTF {
			r = append(r, child)
		}
	}
	return r
}

// Nested secrets are excluded: they would need a per-element walk of the enclosing list.
func CoexistingSecretAttributes(config YamlConfig) []YamlConfigAttribute {
	r := []YamlConfigAttribute{}
	for _, attr := range config.Attributes {
		if attr.CoexistingSecret {
			r = append(r, attr)
		}
	}
	return r
}

func CoexistingSecretChildren(attr YamlConfigAttribute) []YamlConfigAttribute {
	r := []YamlConfigAttribute{}
	for _, child := range attr.Attributes {
		if child.CoexistingSecret {
			r = append(r, child)
		}
	}
	return r
}

func CoexistingSecretParentLists(config YamlConfig) []YamlConfigAttribute {
	r := []YamlConfigAttribute{}
	for _, attr := range config.Attributes {
		if len(CoexistingSecretChildren(attr)) > 0 {
			r = append(r, attr)
		}
	}
	return r
}

func WriteOnlyTFParentLists(config YamlConfig) []YamlConfigAttribute {
	r := []YamlConfigAttribute{}
	for _, attr := range config.Attributes {
		if HasWriteOnlyTFChildren(attr) {
			r = append(r, attr)
		}
	}
	return r
}

func nsoDevicePaths(attrs []YamlConfigAttribute, prefix string) []string {
	var paths []string
	for _, attr := range attrs {
		if attr.Value != "" {
			continue
		}
		if attr.ReferenceNsoDevice {
			paths = append(paths, prefix+attr.TfName)
		}
		paths = append(paths, nsoDevicePaths(attr.Attributes, prefix+attr.TfName+".")...)
	}
	return paths
}

func HasNsoDevices(config YamlConfig) bool {
	return len(nsoDevicePaths(config.Attributes, "")) > 0
}

func NsoDevicesCode(config YamlConfig) string {
	var b strings.Builder
	emitNsoDevices(&b, config.Attributes, "data", 1)
	return b.String()
}

func emitNsoDevices(b *strings.Builder, attrs []YamlConfigAttribute, expr string, depth int) {
	indent := strings.Repeat("\t", depth)
	for _, attr := range attrs {
		if attr.Value != "" {
			continue
		}
		field := expr + "." + ToGoName(attr.TfName)
		switch {
		case attr.ReferenceNsoDevice:
			fmt.Fprintf(b, "%snsoAddDevice(devices, %s)\n", indent, field)
		case len(nsoDevicePaths(attr.Attributes, "")) > 0:
			item := fmt.Sprintf("i%d", depth)
			fmt.Fprintf(b, "%sfor _, %s := range %s {\n", indent, item, field)
			emitNsoDevices(b, attr.Attributes, item, depth+1)
			fmt.Fprintf(b, "%s}\n", indent)
		}
	}
}

func nsoDevicesDescription(config YamlConfig) string {
	paths := nsoDevicePaths(config.Attributes, "")
	if len(paths) == 0 {
		return ""
	}
	for i := range paths {
		paths[i] = "`" + paths[i] + "`"
	}
	return fmt.Sprintf(" Create and update wait until the devices named in %s are ready in Crosswork, as set by the `device_ready_mode` and `device_ready_timeout` arguments of the provider configuration.", strings.Join(paths, ", "))
}

var functions = template.FuncMap{
	"toGoName":                           ToGoName,
	"camelCase":                          CamelCase,
	"strContains":                        strings.Contains,
	"snakeCase":                          SnakeCase,
	"sprintf":                            fmt.Sprintf,
	"toLower":                            strings.ToLower,
	"path":                               BuildPath,
	"hasId":                              HasId,
	"hasReference":                       HasReference,
	"responseEnvelope":                   ResponseEnvelope,
	"readGet":                            ReadGet,
	"hasQueryParam":                      HasQueryParam,
	"hasGetQueryParam":                   HasGetQueryParam,
	"hasDeleteQueryParam":                HasDeleteQueryParam,
	"generateQueryParamString":           GenerateQueryParamString,
	"generateDeleteOnlyQueryParamString": GenerateDeleteOnlyQueryParamString,
	"getId":                              GetId,
	"getMatchId":                         GetMatchId,
	"hasMatchId":                         HasMatchId,
	"cacheFilterGjsonExpr":               CacheFilterGjsonExpr,
	"hasCreateQueryPath":                 HasCreateQueryPath,
	"getCreateQueryPath":                 GetCreateQueryPath,
	"getQueryParam":                      GetQueryParam,
	"getDeleteQueryParam":                GetDeleteQueryParam,
	"hasDataSourceQuery":                 HasDataSourceQuery,
	"hasComputedRefreshValue":            HasComputedRefreshValue,
	"hasComputedRefreshValueInKeys":      HasComputedRefreshValueInKeys,
	"hasUnknownsChildren":                HasUnknownsChildren,
	"firstPathElement":                   FirstPathElement,
	"remainingPathElements":              RemainingPathElements,
	"getFromAllPath":                     GetFromAllPath,
	"isListSet":                          IsListSet,
	"isList":                             IsList,
	"isSet":                              IsSet,
	"isYangFlag":                         IsYangFlag,
	"isStringListSet":                    IsStringListSet,
	"isInt64ListSet":                     IsInt64ListSet,
	"isNestedListSet":                    IsNestedListSet,
	"isNestedList":                       IsNestedList,
	"isNestedSet":                        IsNestedSet,
	"isNestedMap":                        IsNestedMap,
	"isNestedListSetMap":                 IsNestedListSetMap,
	"hasKeyPart":                         HasKeyPart,
	"getKeyPartAttributes":               GetKeyPartAttributes,
	"importAttributes":                   ImportAttributes,
	"constantId":                         ConstantId,
	"importVerifyExcludes":               ImportVerifyExcludes,
	"floatReadExpr":                      FloatReadExpr,
	"subtract":                           Subtract,
	"hasWriteOnlyTFChildren":             HasWriteOnlyTFChildren,
	"writeOnlyTFChildren":                WriteOnlyTFChildren,
	"writeOnlyTFParentLists":             WriteOnlyTFParentLists,
	"coexistingSecretAttributes":         CoexistingSecretAttributes,
	"coexistingSecretChildren":           CoexistingSecretChildren,
	"coexistingSecretParentLists":        CoexistingSecretParentLists,
	"getIfUnsetOnUpdateAttributes":       GetIfUnsetOnUpdateAttributes,
	"hasNsoDevices":                      HasNsoDevices,
	"nsoDevicesCode":                     NsoDevicesCode,
}

// A run of capitals is one acronym (userID becomes user_id, UUID becomes uuid, DNSServers becomes dns_servers). Names that mix
// cases inside an acronym (IPv4Address) are not split correctly; set tf_name for those.
func snakeFromModelName(s string) string {
	r := []rune(s)
	var words []string
	start := 0
	for i := 1; i < len(r); i++ {
		if !unicode.IsUpper(r[i]) {
			continue
		}
		prev := r[i-1]
		if unicode.IsLower(prev) || unicode.IsDigit(prev) ||
			(unicode.IsUpper(prev) && i+1 < len(r) && unicode.IsLower(r[i+1])) {
			words = append(words, strings.ToLower(string(r[start:i])))
			start = i
		}
	}
	words = append(words, strings.ToLower(string(r[start:])))
	return strings.Join(words, "_")
}

func augmentAttribute(attr *YamlConfigAttribute) {
	if attr.TfName == "" {
		attr.TfName = snakeFromModelName(attr.ModelName)
	}
	if attr.Type == "List" || attr.Type == "Set" || attr.Type == "Map" {
		for a := range attr.Attributes {
			augmentAttribute(&attr.Attributes[a])
		}
	}
	if attr.Type == "Map" && attr.MapKeyExample == "" && len(attr.Attributes) > 0 {
		var keyParts []string
		for _, child := range attr.Attributes {
			if child.KeyPart && child.Example != "" {
				keyParts = append(keyParts, child.Example)
			}
		}
		if len(keyParts) > 0 {
			attr.MapKeyExample = strings.Join(keyParts, "_")
		}
	}
}

func augmentConfig(config *YamlConfig) {
	for ia := range config.Attributes {
		augmentAttribute(&config.Attributes[ia])
	}
	augmentWriteOnlyTF(config)
	if config.DsDescription == "" {
		config.DsDescription = fmt.Sprintf("This data source can read the %s.", config.Name)
	}
	if config.ResDescription == "" {
		name := strings.ToLower(config.Name)
		if strings.HasPrefix(name, "a") || strings.HasPrefix(name, "e") || strings.HasPrefix(name, "i") || strings.HasPrefix(name, "o") || strings.HasPrefix(name, "u") {
			config.ResDescription = fmt.Sprintf("This resource can manage an %s.", config.Name)
		} else {
			config.ResDescription = fmt.Sprintf("This resource can manage a %s.", config.Name)
		}
	}
	config.ResDescription += nsoDevicesDescription(*config)
	escapeGoStrings(config)
}

func goEscape(s string) string {
	q := strconv.Quote(s)
	return q[1 : len(q)-1]
}

// Call once, after all attributes are final.
func escapeGoStrings(config *YamlConfig) {
	config.DsDescription = goEscape(config.DsDescription)
	config.ResDescription = goEscape(config.ResDescription)
	escapeAttributeStrings(config.Attributes)
}

func escapeAttributeStrings(attrs []YamlConfigAttribute) {
	for i := range attrs {
		attrs[i].Description = goEscape(attrs[i].Description)
		attrs[i].CoexistenceNote = goEscape(attrs[i].CoexistenceNote)
		for j := range attrs[i].EnumValues {
			attrs[i].EnumValues[j] = goEscape(attrs[i].EnumValues[j])
		}
		escapeAttributeStrings(attrs[i].Attributes)
	}
}

// augmentWriteOnlyTF expands a write_only_tf attribute into three attributes:
//
//	<tf_name>              the original, still Optional, same API path
//	<tf_name>_wo           write-only variant, same API field, never in plan or state
//	<tf_name>_wo_version   Int64, in state, never sent; the rotation trigger
//
// The first two are mutually exclusive: ExactlyOneOf when the secret is mandatory, ConflictsWith
// otherwise.
func augmentWriteOnlyTF(config *YamlConfig) {
	config.Attributes = rewriteWriteOnlyTF(config.Attributes)
}

// Nested secrets get a per-element version companion.
func rewriteWriteOnlyTF(attrs []YamlConfigAttribute) []YamlConfigAttribute {
	newAttrs := make([]YamlConfigAttribute, 0, len(attrs))
	for _, attr := range attrs {
		if len(attr.Attributes) > 0 {
			attr.Attributes = rewriteWriteOnlyTF(attr.Attributes)
		}
		if !attr.WriteOnlyTF {
			newAttrs = append(newAttrs, attr)
			continue
		}
		baseName := attr.TfName

		// tfplugindocs derives Required/Optional from the schema booleans alone, and
		// neither half of the pair can be Required without making that spelling the only
		// usable one, so a mandatory secret says so in the description. Computed before
		// Mandatory is cleared below.
		exclusivity := fmt.Sprintf("Only one of `%s` and `%s_wo` can be set.", baseName, baseName)
		if attr.Mandatory {
			exclusivity = fmt.Sprintf("**Required**: exactly one of `%s` and `%s_wo` must be set.", baseName, baseName)
		}

		// The legacy attribute is not Mandatory on its own, because the "_wo" variant is
		// an equally valid way to supply the secret; ExactlyOneOf carries that requirement.
		// It is excluded from the generated tests so they use the write-only path and do
		// not set both spellings.
		legacy := attr
		legacy.WriteOnlyTF = false
		legacy.CoexistingSecret = true
		legacy.Mandatory = false
		legacy.ExcludeTest = true
		legacy.ExcludeExample = true
		// The minimum test selects attributes on MinimumTestValue without consulting
		// ExcludeTest, so it has to be cleared too or the secret is emitted alongside its
		// "_wo" twin.
		legacy.MinimumTestValue = ""
		legacy.CoexistenceNote = fmt.Sprintf("This attribute stores the secret in Terraform state. Prefer `%s_wo` together with `%s_wo_version`, which keeps it out of state.", baseName, baseName)
		legacy.MutualExclusivityNote = exclusivity
		legacy.WoPairMandatory = attr.Mandatory
		legacy.WoPairHasVersion = !attr.RequiresReplace
		newAttrs = append(newAttrs, legacy)

		// The write-only variant is always Optional in the schema. Mandatory is kept so the
		// template picks ExactlyOneOf over ConflictsWith for the mutual-exclusion validator.
		wo := attr
		wo.TfName = baseName + "_wo"
		wo.WoBaseName = baseName
		wo.MutualExclusivityNote = exclusivity
		newAttrs = append(newAttrs, wo)

		// requires_replace secrets rotate by recreate, so a state-stored version int
		// (which can only drive an in-place update) is meaningless. Convert to write-only
		// but drop the version companion.
		if attr.RequiresReplace {
			continue
		}

		version := YamlConfigAttribute{
			TfName:      baseName + "_wo_version",
			Type:        "Int64",
			WoVersion:   true,
			Description: fmt.Sprintf("Rotation trigger for `%s_wo`. Increment this integer whenever the write-only value changes so Terraform sends the new secret. The value is stored in state; the secret is not.", baseName),
			Example:     "1",
			ExcludeTest: attr.ExcludeTest,
		}
		// The version belongs in the minimum test exactly when its "_wo" secret does; the
		// minimum template does not consult ExcludeTest, so both conditions are mirrored.
		if !attr.ExcludeTest && (attr.Mandatory || attr.MinimumTestValue != "") {
			version.MinimumTestValue = "1"
		}
		newAttrs = append(newAttrs, version)
	}
	return newAttrs
}

func getTemplateSection(content, name string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	result := ""
	foundSection := false
	beginRegex := regexp.MustCompile(`\/\/template:begin\s` + name + `$`)
	endRegex := regexp.MustCompile(`\/\/template:end\s` + name + `$`)
	for scanner.Scan() {
		line := scanner.Text()
		if !foundSection {
			match := beginRegex.MatchString(line)
			if match {
				foundSection = true
				result += line + "\n"
			}
		} else {
			result += line + "\n"
			match := endRegex.MatchString(line)
			if match {
				foundSection = false
			}
		}
	}
	return result
}

func renderTemplate(templatePath, outputPath string, config interface{}) {
	file, err := os.Open(templatePath)
	if err != nil {
		log.Fatalf("Error opening template: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if strings.HasSuffix(templatePath, ".go") {
		scanner.Scan()
	}
	var temp string
	for scanner.Scan() {
		temp = temp + scanner.Text() + "\n"
	}

	template, err := template.New(path.Base(templatePath)).Funcs(functions).Parse(temp)
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	output := new(bytes.Buffer)
	err = template.Execute(output, config)
	if err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	outputFile := filepath.Join(outputPath)
	existingFile, err := os.Open(outputPath)
	if err != nil {
		os.MkdirAll(filepath.Dir(outputFile), 0755)
	} else if strings.HasSuffix(templatePath, ".go") {
		existingScanner := bufio.NewScanner(existingFile)
		var newContent string
		currentSectionName := ""
		beginRegex := regexp.MustCompile(`\/\/template:begin\s(.*?)$`)
		endRegex := regexp.MustCompile(`\/\/template:end\s(.*?)$`)
		for existingScanner.Scan() {
			line := existingScanner.Text()
			if currentSectionName == "" {
				matches := beginRegex.FindStringSubmatch(line)
				if len(matches) > 1 && matches[1] != "" {
					currentSectionName = matches[1]
				} else {
					newContent += line + "\n"
				}
			} else {
				matches := endRegex.FindStringSubmatch(line)
				if len(matches) > 1 && matches[1] == currentSectionName {
					currentSectionName = ""
					newSection := getTemplateSection(string(output.Bytes()), matches[1])
					newContent += newSection
				}
			}
		}
		output = bytes.NewBufferString(newContent)
	}
	f, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	f.Write(output.Bytes())
}

func validatePostBody(config YamlConfig) error {
	if config.PostBody == "" {
		return nil
	}
	if !strings.Contains(config.PostBody, "{body}") {
		return fmt.Errorf("post_body must contain {body}, the generated request body")
	}
	if strings.Contains(config.PostBody, "{id}") && !config.IdFromAttribute {
		return fmt.Errorf("post_body uses {id}, which needs id_from_attribute: true because the id is not known yet when the object is created")
	}
	if config.RootList {
		return fmt.Errorf("post_body cannot be combined with root_list")
	}
	for _, attr := range config.Attributes {
		if attr.MaxElementsInRootList != 0 {
			return fmt.Errorf("post_body cannot be combined with max_elements_in_root_list")
		}
	}
	sample := strings.NewReplacer("{id}", "x", "{body}", "{}").Replace(config.PostBody)
	if !json.Valid([]byte(sample)) {
		return fmt.Errorf("post_body is not valid JSON once {id} and {body} are replaced: %s", sample)
	}
	return nil
}

// max_async_wait_time and no_wait only act on a resource that has a plan_endpoint.
func validatePlanKeys(config YamlConfig) error {
	if config.PlanEndpoint == "" {
		if config.NoWait || config.MaxAsyncWaitTime != 0 {
			return fmt.Errorf("no_wait and max_async_wait_time need plan_endpoint: without a plan there is nothing to wait for")
		}
		return nil
	}
	if config.NoWait && config.MaxAsyncWaitTime != 0 {
		return fmt.Errorf("no_wait and max_async_wait_time cannot be combined: the first skips the wait that the second limits")
	}
	if config.MaxAsyncWaitTime < 0 {
		return fmt.Errorf("max_async_wait_time must be positive")
	}
	for key, set := range map[string]bool{"fallback_rest_endpoint": config.FallbackRestEndpoint != "", "root_list": config.RootList,
		"delete_iterative": config.DeleteIterative, "put_delete": config.PutDelete, "post_delete": config.PostDelete} {
		if set {
			return fmt.Errorf("plan_endpoint cannot be combined with %s: the wait is built into the single-object create, update and delete", key)
		}
	}
	for _, attr := range config.Attributes {
		if attr.MaxElementsInRootList != 0 {
			return fmt.Errorf("plan_endpoint cannot be combined with max_elements_in_root_list")
		}
	}
	if !config.IdFromAttribute {
		return fmt.Errorf("plan_endpoint needs id_from_attribute: true, because the id of the object must be known before the response is read")
	}
	for _, attr := range config.Attributes {
		if attr.Computed {
			return fmt.Errorf("plan_endpoint cannot be combined with the computed attribute %q: an object whose plan fails is kept in the state, which must hold no unknown value", attr.ModelName)
		}
	}
	refs := 0
	for _, attr := range config.Attributes {
		if attr.Reference && !attr.Computed {
			refs++
		}
	}
	if got := strings.Count(config.PlanEndpoint, "%v"); got != refs {
		return fmt.Errorf("plan_endpoint has %d %%v placeholders but the definition has %d reference attributes; they are filled in the order of the attributes", got, refs)
	}
	return nil
}

func validateAttributeKeys(attrs []YamlConfigAttribute, parent string) error {
	for _, attr := range attrs {
		name := parent + attr.ModelName
		if attr.ReferenceInBody && !attr.Reference {
			return fmt.Errorf("attribute %q sets reference_in_body, which needs reference: true", name)
		}
		if attr.ReferenceNsoDevice && (attr.Type != "String" || attr.Value != "" || attr.WriteOnly) {
			return fmt.Errorf("attribute %q sets reference_nso_device, which needs a String attribute that is sent to the controller and has no fixed value", name)
		}
		// KeepStringSetIn handles only a String Set; on a List the generated call would use the wrong
		// helper or panic.
		if attr.IgnoreExtraResponseValues && (attr.Type != "Set" || attr.ElementType != "String") {
			return fmt.Errorf("attribute %q sets ignore_extra_response_values, which is only supported on String Set attributes, but it has type %q element_type %q", name, attr.Type, attr.ElementType)
		}
		if attr.WriteOnlyTFNative && (!attr.WriteOnly || attr.WriteOnlyTF || attr.Computed || attr.DefaultValue != "") {
			return fmt.Errorf("attribute %q sets write_only_tf_native, which requires write_only and cannot be combined with write_only_tf, computed, or default_value", name)
		}
		if attr.WriteOnlyTF {
			// The generated "_wo" schema entry emits []validator.String, so a non-String
			// secret would not compile.
			if attr.Type != "String" {
				return fmt.Errorf("attribute %q sets write_only_tf, which is only supported on String attributes, but it has type %q", name, attr.Type)
			}
			// fromBody and updateFromBody skip attributes on .WriteOnly alone, so a secret
			// flagged write_only_tf without it would have the API response written back into
			// its "_wo" field on every read.
			if !attr.WriteOnly {
				return fmt.Errorf("attribute %q sets write_only_tf, which requires write_only to be set as well", name)
			}
		}
		if attr.EscapeReferencePath && !attr.Reference {
			return fmt.Errorf("attribute %q sets escape_reference_path, which needs reference: true", name)
		}
		switch attr.TypeYangBool {
		case "":
		case "empty", "presence", "boolean":
			if attr.Type != "Bool" {
				return fmt.Errorf("attribute %q sets type_yang_bool, which needs type: Bool", name)
			}
			if attr.PutDataPath != "" && attr.TypeYangBool != "boolean" {
				return fmt.Errorf("attribute %q sets type_yang_bool %s together with put_data_path, which is not supported for it", name, attr.TypeYangBool)
			}
		default:
			return fmt.Errorf("attribute %q has type_yang_bool %q: use empty, presence or boolean", name, attr.TypeYangBool)
		}
		if err := validateAttributeKeys(attr.Attributes, name+"."); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	files, _ := os.ReadDir(definitionsPath)
	configs := make([]YamlConfig, 0, len(files))

	// Only *.yaml files are definitions (the folder also holds a README).
	for _, filename := range files {
		if filepath.Ext(filename.Name()) != ".yaml" {
			continue
		}
		yamlFile, err := os.ReadFile(filepath.Join(definitionsPath, filename.Name()))
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}

		config := YamlConfig{}
		err = yaml.Unmarshal(yamlFile, &config)
		if err != nil {
			log.Fatalf("Error parsing yaml: %v", err)
		}
		if err := validatePlanKeys(config); err != nil {
			log.Fatalf("definition %q: %v", config.Name, err)
		}
		if err := validateAttributeKeys(config.Attributes, ""); err != nil {
			log.Fatalf("definition %q: %v", config.Name, err)
		}
		if err := validatePostBody(config); err != nil {
			log.Fatalf("definition %q: %v", config.Name, err)
		}
		configs = append(configs, config)
	}

	for i := range configs {
		augmentConfig(&configs[i])

		for _, t := range templates {
			if (configs[i].NoImport && t.path == "./gen/templates/import.sh") ||
				(configs[i].NoDataSource && t.path == "./gen/templates/data_source.go") ||
				(configs[i].NoDataSource && t.path == "./gen/templates/data_source_test.go") ||
				(configs[i].NoDataSource && t.path == "./gen/templates/data-source.tf") ||
				(configs[i].NoResource && t.path == "./gen/templates/resource.go") ||
				(configs[i].NoResource && t.path == "./gen/templates/resource_test.go") ||
				(configs[i].NoResource && t.path == "./gen/templates/resource.tf") ||
				(configs[i].NoResource && t.path == "./gen/templates/import.sh") {
				continue
			}
			renderTemplate(t.path, t.prefix+SnakeCase(configs[i].Name)+t.suffix, configs[i])
		}
	}

	renderTemplate(providerTemplate, providerLocation, configs)

	changelog, err := os.ReadFile(changelogOriginal)
	if err != nil {
		log.Fatalf("Error reading changelog: %v", err)
	}
	renderTemplate(changelogTemplate, changelogLocation, string(changelog))
}

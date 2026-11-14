//go:build ignore
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
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)
// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
{{- $name := camelCase .Name}}
type {{camelCase .Name}} struct {
	Instance types.String `tfsdk:"instance"`
	Id types.String `tfsdk:"id"`
{{- range .Attributes}}
{{- if not .Value}}
{{- if isNestedMap .}}
	{{toGoName .TfName}} map[string]{{$name}}{{toGoName .TfName}} `tfsdk:"{{.TfName}}"`
{{- else if isNestedListSet .}}
	{{toGoName .TfName}} []{{$name}}{{toGoName .TfName}} `tfsdk:"{{.TfName}}"`
{{- else}}
	{{toGoName .TfName}} types.{{.Type}} `tfsdk:"{{.TfName}}"`
{{- end}}
{{- end}}
{{- end}}
}

{{ range .Attributes}}
{{- if not .Value}}
{{- $childName := toGoName .TfName}}
{{- if isNestedListSetMap .}}
type {{$name}}{{toGoName .TfName}} struct {
{{- range .Attributes}}
{{- if not .Value}}
	{{- if isNestedMap .}}
	{{toGoName .TfName}} map[string]{{$name}}{{$childName}}{{toGoName .TfName}} `tfsdk:"{{.TfName}}"`
	{{- else if isNestedListSetMap .}}
	{{toGoName .TfName}} []{{$name}}{{$childName}}{{toGoName .TfName}} `tfsdk:"{{.TfName}}"`
{{- else}}
	{{toGoName .TfName}} types.{{.Type}} `tfsdk:"{{.TfName}}"`
{{- end}}
{{- end}}
{{- end}}
}
{{- end}}
{{- end}}
{{ end}}

{{ range .Attributes}}
{{- if not .Value}}
{{- $childName := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $childChildName := toGoName .TfName}}
{{- if isNestedListSetMap .}}
type {{$name}}{{$childName}}{{toGoName .TfName}} struct {
{{- range .Attributes}}
{{- if not .Value}}
	{{- if isNestedMap .}}
	{{toGoName .TfName}} map[string]{{$name}}{{$childName}}{{$childChildName}}{{toGoName .TfName}} `tfsdk:"{{.TfName}}"`
	{{- else if isNestedListSetMap .}}
	{{toGoName .TfName}} []{{$name}}{{$childName}}{{$childChildName}}{{toGoName .TfName}} `tfsdk:"{{.TfName}}"`
{{- else}}
	{{toGoName .TfName}} types.{{.Type}} `tfsdk:"{{.TfName}}"`
{{- end}}
{{- end}}
{{- end}}
}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{ end}}

{{ range .Attributes}}
{{- if not .Value}}
{{- $childName := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $childChildName := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $c3Name := toGoName .TfName}}
{{- if isNestedListSetMap .}}
type {{$name}}{{$childName}}{{$childChildName}}{{toGoName .TfName}} struct {
{{- range .Attributes}}
{{- if not .Value}}
	{{- if isNestedMap .}}
	{{toGoName .TfName}} map[string]{{$name}}{{$childName}}{{$childChildName}}{{$c3Name}}{{toGoName .TfName}} `tfsdk:"{{.TfName}}"`
	{{- else if isNestedListSetMap .}}
	{{toGoName .TfName}} []{{$name}}{{$childName}}{{$childChildName}}{{$c3Name}}{{toGoName .TfName}} `tfsdk:"{{.TfName}}"`
{{- else}}
	{{toGoName .TfName}} types.{{.Type}} `tfsdk:"{{.TfName}}"`
{{- end}}
{{- end}}
{{- end}}
}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{ end}}

{{ range .Attributes}}
{{- if not .Value}}
{{- $childName := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $childChildName := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $c3Name := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $c4Name := toGoName .TfName}}
{{- if isNestedListSetMap .}}
type {{$name}}{{$childName}}{{$childChildName}}{{$c3Name}}{{toGoName .TfName}} struct {
{{- range .Attributes}}
{{- if not .Value}}
	{{- if isNestedMap .}}
	{{toGoName .TfName}} map[string]{{$name}}{{$childName}}{{$childChildName}}{{$c3Name}}{{$c4Name}}{{toGoName .TfName}} `tfsdk:"{{.TfName}}"`
	{{- else if isNestedListSetMap .}}
	{{toGoName .TfName}} []{{$name}}{{$childName}}{{$childChildName}}{{$c3Name}}{{$c4Name}}{{toGoName .TfName}} `tfsdk:"{{.TfName}}"`
{{- else}}
	{{toGoName .TfName}} types.{{.Type}} `tfsdk:"{{.TfName}}"`
{{- end}}
{{- end}}
{{- end}}
}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{ end}}

{{ range .Attributes}}
{{- if not .Value}}
{{- $childName := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $childChildName := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $c3Name := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $c4Name := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $c5Name := toGoName .TfName}}
{{- if isNestedListSetMap .}}
type {{$name}}{{$childName}}{{$childChildName}}{{$c3Name}}{{$c4Name}}{{toGoName .TfName}} struct {
{{- range .Attributes}}
{{- if not .Value}}
	{{- if isNestedMap .}}
	{{toGoName .TfName}} map[string]{{$name}}{{$childName}}{{$childChildName}}{{$c3Name}}{{$c4Name}}{{$c5Name}}{{toGoName .TfName}} `tfsdk:"{{.TfName}}"`
	{{- else if isNestedListSetMap .}}
	{{toGoName .TfName}} []{{$name}}{{$childName}}{{$childChildName}}{{$c3Name}}{{$c4Name}}{{$c5Name}}{{toGoName .TfName}} `tfsdk:"{{.TfName}}"`
{{- else}}
	{{toGoName .TfName}} types.{{.Type}} `tfsdk:"{{.TfName}}"`
{{- end}}
{{- end}}
{{- end}}
}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{ end}}

{{ range .Attributes}}
{{- if not .Value}}
{{- $childName := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $childChildName := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $c3Name := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $c4Name := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $c5Name := toGoName .TfName}}
{{- if isNestedListSetMap .}}
{{ range .Attributes}}
{{- if not .Value}}
{{- $c6Name := toGoName .TfName}}
{{- if isNestedListSetMap .}}
type {{$name}}{{$childName}}{{$childChildName}}{{$c3Name}}{{$c4Name}}{{$c5Name}}{{toGoName .TfName}} struct {
{{- range .Attributes}}
{{- if not .Value}}
	{{toGoName .TfName}} types.{{.Type}} `tfsdk:"{{.TfName}}"`
{{- end}}
{{- end}}
}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
{{ end}}
// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data {{camelCase .Name}}) getPath() string {
	{{- if hasReference .Attributes}}
		return fmt.Sprintf("{{.RestEndpoint}}"{{range .Attributes}}{{if and .Reference (not .Computed)}}, {{if .EscapeReferencePath}}url.PathEscape{{else}}url.QueryEscape{{end}}(data.{{toGoName .TfName}}.Value{{.Type}}()){{end}}{{end}})
	{{- else}}
		return "{{.RestEndpoint}}"
	{{- end}}
}
{{- if hasNsoDevices .}}

func (data {{camelCase .Name}}) nsoDevices() nsoDeviceSet {
	devices := nsoDeviceSet{}
{{nsoDevicesCode .}}	return devices
}
{{- end}}
// End of section. //template:end getPath

// Section below is generated&owned by "gen/generator.go". //template:begin getFallbackPath
{{if .FallbackRestEndpoint}}
func (data {{camelCase .Name}}) getFallbackPath() string {
	{{- if hasReference .Attributes}}
		return fmt.Sprintf("{{.FallbackRestEndpoint}}"{{range .Attributes}}{{if .Reference}}, {{if .EscapeReferencePath}}url.PathEscape{{else}}url.QueryEscape{{end}}(data.{{toGoName .TfName}}.Value{{.Type}}()){{end}}{{end}})
	{{- else}}
		return "{{.FallbackRestEndpoint}}"
	{{- end}}
}
{{- end}}
// End of section. //template:end getFallbackPath

// Section below is generated&owned by "gen/generator.go". //template:begin getPathDelete
{{if .DeleteRestEndpoint}}
func (data {{camelCase .Name}}) getPathDelete() string {
	{{- if and (hasReference .Attributes) (strContains .DeleteRestEndpoint "%v")}}
		return fmt.Sprintf("{{.DeleteRestEndpoint}}"{{range .Attributes}}{{if .Reference}}, {{if .EscapeReferencePath}}url.PathEscape{{else}}url.QueryEscape{{end}}(data.{{toGoName .TfName}}.Value{{.Type}}()){{end}}{{end}})
	{{- else}}
		return "{{.DeleteRestEndpoint}}"
	{{- end}}
}
{{end}}
// End of section. //template:end getPathDelete

// Section below is generated&owned by "gen/generator.go". //template:begin getPathGet
{{if .GetRestEndpoint}}
func (data {{camelCase .Name}}) getPathGet() string {
	{{- if and (hasReference .Attributes) (strContains .GetRestEndpoint "%v")}}
		return fmt.Sprintf("{{.GetRestEndpoint}}"{{range .Attributes}}{{if .Reference}}, {{if .EscapeReferencePath}}url.PathEscape{{else}}url.QueryEscape{{end}}(data.{{toGoName .TfName}}.Value{{.Type}}()){{end}}{{end}})
	{{- else}}
		return "{{.GetRestEndpoint}}"
	{{- end}}
}
{{end}}
// End of section. //template:end getPathGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPost
{{if .PostRestEndpoint}}
func (data {{camelCase .Name}}) getPathPost() string {
	{{- if and (hasReference .Attributes) (strContains .PostRestEndpoint "%v")}}
		return fmt.Sprintf("{{.PostRestEndpoint}}"{{range .Attributes}}{{if and .Reference (not .Computed)}}, {{if .EscapeReferencePath}}url.PathEscape{{else}}url.QueryEscape{{end}}(data.{{toGoName .TfName}}.Value{{.Type}}()){{end}}{{end}})
	{{- else}}
		return "{{.PostRestEndpoint}}"
	{{- end}}
}
{{end}}
// End of section. //template:end getPathPost

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPut
{{if .PutRestEndpoint}}
func (data {{camelCase .Name}}) getPathPut() string {
	{{- if and (hasReference .Attributes) (strContains .PutRestEndpoint "%v")}}
		return fmt.Sprintf("{{.PutRestEndpoint}}"{{range .Attributes}}{{if .Reference}}, {{if .EscapeReferencePath}}url.PathEscape{{else}}url.QueryEscape{{end}}(data.{{toGoName .TfName}}.Value{{.Type}}()){{end}}{{end}})
	{{- else}}
		return "{{.PutRestEndpoint}}"
	{{- end}}
}
{{end}}
// End of section. //template:end getPathPut

// Section below is generated&owned by "gen/generator.go". //template:begin getPathPlan
{{if .PlanEndpoint}}
func (data {{camelCase .Name}}) getPathPlan() string {
	{{- if hasReference .Attributes}}
		return fmt.Sprintf("{{.PlanEndpoint}}"{{range .Attributes}}{{if and .Reference (not .Computed)}}, {{if .EscapeReferencePath}}url.PathEscape{{else}}url.QueryEscape{{end}}(data.{{toGoName .TfName}}.Value{{.Type}}()){{end}}{{end}})
	{{- else}}
		return "{{.PlanEndpoint}}"
	{{- end}}
}
{{end}}
// End of section. //template:end getPathPlan

// Section below is generated&owned by "gen/generator.go". //template:begin readGet
{{if or .GetBody .PostBody (and .DeleteNoId .DeleteBody (not .PutDelete) (not .PostDelete))}}
// Escapes the id so it cannot change the structure of a request body.
func (data {{camelCase .Name}}) idJSON() string {
	{{- if .IdFromAttribute}}
	{{- $readId := getId .Attributes}}
	id := data.Id.ValueString()
	if id == "" {
		id = fmt.Sprint(data.{{toGoName $readId.TfName}}.Value{{$readId.Type}}())
	}
	{{- else}}
	id := data.Id.ValueString()
	{{- end}}
	b, _ := json.Marshal(id)
	return string(b[1 : len(b)-1])
}
{{end}}
{{- if .GetBody}}

func (data {{camelCase .Name}}) readGet(c *cw.Client, path string) (cw.Res, error) {
	return c.Post(path, strings.ReplaceAll(`{{.GetBody}}`, "{id}", data.idJSON()), cw.ReadOnly)
}
{{- end}}
{{- if .PostBody}}

// postBody places the generated create body in the request template, in one pass so a value that
// looks like a placeholder (a password such as "p{id}w") is never expanded twice.
func (data {{camelCase .Name}}) postBody(body string) string {
	return strings.NewReplacer("{id}", data.idJSON(), "{body}", body).Replace(`{{.PostBody}}`)
}
{{- end}}
// End of section. //template:end readGet

// Section below is generated&owned by "gen/generator.go". //template:begin getPathIdQuery
{{if .IdQueryRestEndpoint}}
func (data {{camelCase .Name}}) getPathIdQuery() string {
	{{- if and (hasReference .Attributes) (strContains .IdQueryRestEndpoint "%v")}}
		return fmt.Sprintf("{{.IdQueryRestEndpoint}}"{{range .Attributes}}{{if and .Reference (not .Computed)}}, {{if .EscapeReferencePath}}url.PathEscape{{else}}url.QueryEscape{{end}}(data.{{toGoName .TfName}}.Value{{.Type}}()){{end}}{{end}})
	{{- else}}
		return "{{.IdQueryRestEndpoint}}"
	{{- end}}
}
{{end}}
// End of section. //template:end getPathIdQuery

// Section below is generated&owned by "gen/generator.go". //template:begin toBody
func (data {{camelCase .Name}}) toBody(ctx context.Context, state {{camelCase .Name}}) string {
	{{- if .RootList}}
	body := "[]"
	{{- else}}
	body := ""
	{{- end}}
	put := false
	if state.Id.ValueString() != "" {
		put = true
		{{- if .PutIdIncludePath}}
		body, _ = sjson.Set(body, "{{ .PutIdIncludePath}}", state.Id.ValueString())
		{{- end}}
	}
	_ = put
	{{- range .Attributes}}
	{{- if .Value}}
	{{- if .ValueCondition}}
	if {{.ValueCondition}} {
		{{- if .PutDataPath}}
		if put {
			body, _ = sjson.Set(body, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
		} else {
			body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
		}
		{{- else}}
		body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
		{{- end}}
	}
	{{- else}}
	{{- if .PutDataPath}}
	if put {
		{{- if eq .PutDataPath "_"}}
		body, _ = sjson.Set(body, "{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
		{{- else}}
		body, _ = sjson.Set(body, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
		{{- end}}
	} else {
		body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
	}
	{{- else}}
	body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
	{{- end}}
	{{- end}}
	{{- else if and (or (not .Reference) .ReferenceInBody) (not .CreateQueryPath) (not .QueryParamNoBody) (not .WoVersion) (eq .Type "Bool") .TypeYangBool (ne .TypeYangBool "boolean")}}
	if data.{{toGoName .TfName}}.ValueBool() {
		body, _ = sjson.SetRaw(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", "{{if eq .TypeYangBool "empty"}}[null]{{else}}{}{{end}}")
	}
	{{- else if and (or (not .Reference) .ReferenceInBody) (not .CreateQueryPath) (not .QueryParamNoBody) (not .WoVersion)}}
	{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
	if {{- if .Computed}}{{if eq .Type "Bool"}} !data.{{toGoName .TfName}}.IsUnknown(){{else}} data.{{toGoName .TfName}}.Value{{if eq .Type "Int64"}}Int64() != 0{{else if eq .Type "Float64"}}Float64() != 0{{else}}String() != ""{{end}}{{end}} &&{{- end}} !data.{{toGoName .TfName}}.IsNull() {{if .ExcludeFromPut}}&& put == false{{end}} {
		{{- if .PutDataPath}}
		if put {
			{{- if eq .PutDataPath "_"}}
			body, _ = sjson.Set(body, "{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}data.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
			{{- else}}
			body, _ = sjson.Set(body, "{{.PutDataPath}}.{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}data.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
			{{- end}}
		} else {
			body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}data.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
		}
		{{- else}}
		body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}data.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
		{{- end}}
	}
	{{- else if isListSet .}}
	if !data.{{toGoName .TfName}}.IsNull() {
		var values []{{if isStringListSet .}}string{{else if isInt64ListSet .}}int64{{end}}
		data.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
		{{- if .PutDataPath}}
		if put {
			{{- if eq .PutDataPath "_"}}
			body, _ = sjson.Set(body, "{{.ModelName}}", values)
			{{- else}}
			body, _ = sjson.Set(body, "{{.PutDataPath}}.{{.ModelName}}", values)
			{{- end}}
		} else {
			body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
		}
		{{- else}}
		body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
		{{- end}}
	}
	{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
	if !data.{{toGoName .TfName}}.IsNull() {
		{{- if $.NoRead}}
		var listValues map[string]types.List
		data.{{toGoName .TfName}}.ElementsAs(ctx, &listValues, false)
		params := make(map[string]any)
		for k, v := range listValues {
			var strList []string
			v.ElementsAs(ctx, &strList, false)
			if len(strList) == 1 {
				params[k] = strList[0]
			} else {
				params[k] = strList
			}
		}
		{{- if .PutDataPath}}
		if put {
			{{- if eq .PutDataPath "_"}}
			body, _ = sjson.Set(body, "{{.ModelName}}", params)
			{{- else}}
			body, _ = sjson.Set(body, "{{.PutDataPath}}.{{.ModelName}}", params)
			{{- end}}
		} else {
			body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
		}
		{{- else}}
		body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
		{{- end}}
		{{- else}}
		var values map[string]string
		data.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
		{{- if .PutDataPath}}
		if put {
			{{- if eq .PutDataPath "_"}}
			body, _ = sjson.Set(body, "{{.ModelName}}", values)
			{{- else}}
			body, _ = sjson.Set(body, "{{.PutDataPath}}.{{.ModelName}}", values)
			{{- end}}
		} else {
			body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
		}
		{{- else}}
		body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
		{{- end}}
		{{- end}}
	}
	{{- else if isNestedListSetMap .}}
	{{- if .AlwaysInclude}}
	if true {
	{{- else}}
	if len(data.{{toGoName .TfName}}) > 0 {
	{{- end}}
		{{- if .PutDataPath}}
		if put {
			body, _ = sjson.Set(body, "{{.PutDataPath}}.{{.ModelName}}", []interface{}{})
		} else {
			body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", []interface{}{})
		}
		{{- else}}
		body, _ = sjson.Set(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", []interface{}{})
		{{- end}}
		for _, item := range data.{{toGoName .TfName}} {
			itemBody := ""
			{{- range .Attributes}}
			{{- if .Value}}
			{{- if .ValueCondition}}
			if {{.ValueCondition}} {
				{{- if .PutDataPath}}
				if put {
					itemBody, _ = sjson.Set(itemBody, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
				} else {
					itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
				}
				{{- else}}
				itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
				{{- end}}
			}
			{{- else}}
			{{- if .PutDataPath}}
			if put {
				itemBody, _ = sjson.Set(itemBody, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
			} else {
				itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
			}
			{{- else}}
			itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
			{{- end}}
			{{- end}}
			{{- else if and (or (not .Reference) .ReferenceInBody) (not .QueryParamNoBody) (not .WoVersion) (eq .Type "Bool") .TypeYangBool (ne .TypeYangBool "boolean")}}
			if item.{{toGoName .TfName}}.ValueBool() {
				itemBody, _ = sjson.SetRaw(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", "{{if eq .TypeYangBool "empty"}}[null]{{else}}{}{{end}}")
			}
			{{- else if and (or (not .Reference) .ReferenceInBody) (not .QueryParamNoBody) (not .WoVersion)}}
			{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
			if {{- if .Computed}}{{if eq .Type "Bool"}} !item.{{toGoName .TfName}}.IsUnknown(){{else}} item.{{toGoName .TfName}}.Value{{if eq .Type "Int64"}}Int64() != 0{{else if eq .Type "Float64"}}Float64() != 0{{else}}String() != ""{{end}}{{end}} &&{{- end}} !item.{{toGoName .TfName}}.IsNull() {{if .ExcludeFromPut}}&& put == false{{end}} {
				{{- if .PutDataPath}}
				if put {
					itemBody, _ = sjson.Set(itemBody, "{{.PutDataPath}}.{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}item.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
				} else {
					itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}item.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
				}
				{{- else}}
				itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}item.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
				{{- end}}
			}
			{{- else if isListSet .}}
			if !item.{{toGoName .TfName}}.IsNull() {
				var values []{{if isStringListSet .}}string{{else if isInt64ListSet .}}int64{{end}}
				item.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
				{{- if .PutDataPath}}
				if put {
					itemBody, _ = sjson.Set(itemBody, "{{.PutDataPath}}.{{.ModelName}}", values)
				} else {
					itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
				}
				{{- else}}
				itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
				{{- end}}
			}
			{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
			if !item.{{toGoName .TfName}}.IsNull() {
				{{- if $.NoRead}}
				var listValues map[string]types.List
				item.{{toGoName .TfName}}.ElementsAs(ctx, &listValues, false)
				params := make(map[string]any)
				for k, v := range listValues {
					var strList []string
					v.ElementsAs(ctx, &strList, false)
					if len(strList) == 1 {
						params[k] = strList[0]
					} else {
						params[k] = strList
					}
				}
				{{- if .PutDataPath}}
				if put {
					itemBody, _ = sjson.Set(itemBody, "{{.PutDataPath}}.{{.ModelName}}", params)
				} else {
					itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
				}
				{{- else}}
				itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
				{{- end}}
				{{- else}}
				var values map[string]string
				item.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
				{{- if .PutDataPath}}
				if put {
					itemBody, _ = sjson.Set(itemBody, "{{.PutDataPath}}.{{.ModelName}}", values)
				} else {
					itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
				}
				{{- else}}
				itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
				{{- end}}
				{{- end}}
			}
			{{- else if isNestedListSetMap .}}
			{{- if .AlwaysInclude}}
			if true {
			{{- else}}
			if len(item.{{toGoName .TfName}}) > 0 {
			{{- end}}
				{{- if .PutDataPath}}
				if put {
					itemBody, _ = sjson.Set(itemBody, "{{.PutDataPath}}.{{.ModelName}}", []interface{}{})
				} else {
					itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", []interface{}{})
				}
				{{- else}}
				itemBody, _ = sjson.Set(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", []interface{}{})
				{{- end}}
				for _, childItem := range item.{{toGoName .TfName}} {
					itemChildBody := ""
					{{- range .Attributes}}
					{{- if .Value}}
					{{- if .ValueCondition}}
					if {{.ValueCondition}} {
						{{- if .PutDataPath}}
						if put {
							itemChildBody, _ = sjson.Set(itemChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
						} else {
							itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
						}
						{{- else}}
						itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
						{{- end}}
					}
					{{- else}}
					{{- if .PutDataPath}}
					if put {
						itemChildBody, _ = sjson.Set(itemChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
					} else {
						itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
					}
					{{- else}}
					itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
					{{- end}}
					{{- end}}
					{{- else if and (or (not .Reference) .ReferenceInBody) (not .QueryParamNoBody) (not .WoVersion) (eq .Type "Bool") .TypeYangBool (ne .TypeYangBool "boolean")}}
					if childItem.{{toGoName .TfName}}.ValueBool() {
						itemChildBody, _ = sjson.SetRaw(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", "{{if eq .TypeYangBool "empty"}}[null]{{else}}{}{{end}}")
					}
					{{- else if and (or (not .Reference) .ReferenceInBody) (not .QueryParamNoBody) (not .WoVersion)}}
					{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
					if !childItem.{{toGoName .TfName}}.IsNull() {
						{{- if .PutDataPath}}
						if put {
							itemChildBody, _ = sjson.Set(itemChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
						} else {
							itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
						}
						{{- else}}
						itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
						{{- end}}
					}
					{{- else if isListSet .}}
					if !childItem.{{toGoName .TfName}}.IsNull() {
						var values []{{if isStringListSet .}}string{{else if isInt64ListSet .}}int64{{end}}
						childItem.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
						{{- if .PutDataPath}}
						if put {
							itemChildBody, _ = sjson.Set(itemChildBody, "{{.PutDataPath}}.{{.ModelName}}", values)
						} else {
							itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
						}
						{{- else}}
						itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
						{{- end}}
					}
					{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
					if !childItem.{{toGoName .TfName}}.IsNull() {
						{{- if $.NoRead}}
						var listValues map[string]types.List
						childItem.{{toGoName .TfName}}.ElementsAs(ctx, &listValues, false)
						params := make(map[string]any)
						for k, v := range listValues {
							var strList []string
							v.ElementsAs(ctx, &strList, false)
							if len(strList) == 1 {
								params[k] = strList[0]
							} else {
								params[k] = strList
							}
						}
						{{- if .PutDataPath}}
						if put {
							itemChildBody, _ = sjson.Set(itemChildBody, "{{.PutDataPath}}.{{.ModelName}}", params)
						} else {
							itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
						}
						{{- else}}
						itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
						{{- end}}
						{{- else}}
						var values map[string]string
						childItem.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
						{{- if .PutDataPath}}
						if put {
							itemChildBody, _ = sjson.Set(itemChildBody, "{{.PutDataPath}}.{{.ModelName}}", values)
						} else {
							itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
						}
						{{- else}}
						itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
						{{- end}}
						{{- end}}
					}
					{{- else if isNestedListSetMap .}}
					{{- if .AlwaysInclude}}
					if true {
					{{- else}}
					if len(childItem.{{toGoName .TfName}}) > 0 {
					{{- end}}
						{{- if .PutDataPath}}
						if put {
							itemChildBody, _ = sjson.Set(itemChildBody, "{{.PutDataPath}}.{{.ModelName}}", []interface{}{})
						} else {
							itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", []interface{}{})
						}
						{{- else}}
						itemChildBody, _ = sjson.Set(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", []interface{}{})
						{{- end}}
						for _, childChildItem := range childItem.{{toGoName .TfName}} {
							itemChildChildBody := ""
							{{- range .Attributes}}
							{{- if .Value}}
							{{- if .ValueCondition}}
							if {{.ValueCondition}} {
								{{- if .PutDataPath}}
								if put {
									itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
								} else {
									itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
								}
								{{- else}}
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
								{{- end}}
							}
							{{- else}}
							{{- if .PutDataPath}}
							if put {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
							} else {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
							}
							{{- else}}
							itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
							{{- end}}
							{{- end}}
		{{- else if and (or (not .Reference) .ReferenceInBody) (not .QueryParamNoBody) (not .WoVersion) (eq .Type "Bool") .TypeYangBool (ne .TypeYangBool "boolean")}}
		if childChildItem.{{toGoName .TfName}}.ValueBool() {
			itemChildChildBody, _ = sjson.SetRaw(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", "{{if eq .TypeYangBool "empty"}}[null]{{else}}{}{{end}}")
		}
		{{- else if and (or (not .Reference) .ReferenceInBody) (not .QueryParamNoBody) (not .WoVersion)}}
							{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
							if !childChildItem.{{toGoName .TfName}}.IsNull() {
								{{- if .PutDataPath}}
								if put {
									itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childChildItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
								} else {
									itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childChildItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
								}
								{{- else}}
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childChildItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
								{{- end}}
							}
							{{- else if isListSet .}}
							if !childChildItem.{{toGoName .TfName}}.IsNull() {
								var values []{{if isStringListSet .}}string{{else if isInt64ListSet .}}int64{{end}}
								childChildItem.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
								{{- if .PutDataPath}}
								if put {
									itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", values)
								} else {
									itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
								}
								{{- else}}
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
								{{- end}}
							}
							{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
							if !childChildItem.{{toGoName .TfName}}.IsNull() {
								{{- if $.NoRead}}
								var listValues map[string]types.List
								childChildItem.{{toGoName .TfName}}.ElementsAs(ctx, &listValues, false)
								params := make(map[string]any)
								for k, v := range listValues {
									var strList []string
									v.ElementsAs(ctx, &strList, false)
									if len(strList) == 1 {
										params[k] = strList[0]
									} else {
										params[k] = strList
									}
								}
								{{- if .PutDataPath}}
								if put {
									itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", params)
								} else {
									itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
								}
								{{- else}}
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
								{{- end}}
								{{- else}}
								var values map[string]string
								childChildItem.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
								{{- if .PutDataPath}}
								if put {
									itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", values)
								} else {
									itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
								}
								{{- else}}
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
								{{- end}}
								{{- end}}
							}
						{{- else if isNestedListSetMap .}}
						{{- if .AlwaysInclude}}
						if true {
						{{- else}}
						if len(childChildItem.{{toGoName .TfName}}) > 0 {
						{{- end}}
							{{- if .PutDataPath}}
							if put {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", []interface{}{})
							} else {
								itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", []interface{}{})
							}
							{{- else}}
							itemChildChildBody, _ = sjson.Set(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", []interface{}{})
							{{- end}}
							for _, childChildChildItem := range childChildItem.{{toGoName .TfName}} {
								itemChildChildChildBody := ""
								{{- range .Attributes}}
								{{- if .Value}}
								{{- if .ValueCondition}}
								if {{.ValueCondition}} {
									{{- if .PutDataPath}}
									if put {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
									} else {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
									}
									{{- else}}
									itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
									{{- end}}
								}
								{{- else}}
								{{- if .PutDataPath}}
								if put {
									itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
								} else {
									itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
								}
								{{- else}}
								itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
								{{- end}}
								{{- end}}
			{{- else if and (or (not .Reference) .ReferenceInBody) (not .QueryParamNoBody) (not .WoVersion) (eq .Type "Bool") .TypeYangBool (ne .TypeYangBool "boolean")}}
			if childChildChildItem.{{toGoName .TfName}}.ValueBool() {
				itemChildChildChildBody, _ = sjson.SetRaw(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", "{{if eq .TypeYangBool "empty"}}[null]{{else}}{}{{end}}")
			}
			{{- else if and (or (not .Reference) .ReferenceInBody) (not .QueryParamNoBody) (not .WoVersion)}}
								{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
								if !childChildChildItem.{{toGoName .TfName}}.IsNull() {
									{{- if .PutDataPath}}
									if put {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childChildChildItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
									} else {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childChildChildItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
									}
									{{- else}}
									itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childChildChildItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
									{{- end}}
								}
								{{- else if isListSet .}}
								if !childChildChildItem.{{toGoName .TfName}}.IsNull() {
									var values []{{if isStringListSet .}}string{{else if isInt64ListSet .}}int64{{end}}
									childChildChildItem.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
									{{- if .PutDataPath}}
									if put {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", values)
									} else {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
									}
									{{- else}}
									itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
									{{- end}}
								}
								{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
								if !childChildChildItem.{{toGoName .TfName}}.IsNull() {
									{{- if $.NoRead}}
									var listValues map[string]types.List
									childChildChildItem.{{toGoName .TfName}}.ElementsAs(ctx, &listValues, false)
									params := make(map[string]any)
									for k, v := range listValues {
										var strList []string
										v.ElementsAs(ctx, &strList, false)
										if len(strList) == 1 {
											params[k] = strList[0]
										} else {
											params[k] = strList
										}
									}
									{{- if .PutDataPath}}
									if put {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", params)
									} else {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
									}
									{{- else}}
									itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
									{{- end}}
									{{- else}}
									var values map[string]string
									childChildChildItem.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
									{{- if .PutDataPath}}
									if put {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", values)
									} else {
										itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
									}
									{{- else}}
									itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
									{{- end}}
									{{- end}}
								}
							{{- else if isNestedListSetMap .}}
							{{- if .AlwaysInclude}}
							if true {
							{{- else}}
							if len(childChildChildItem.{{toGoName .TfName}}) > 0 {
							{{- end}}
								{{- if .PutDataPath}}
								if put {
									itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", []interface{}{})
								} else {
									itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", []interface{}{})
								}
								{{- else}}
								itemChildChildChildBody, _ = sjson.Set(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", []interface{}{})
								{{- end}}
								for _, childChildChildChildItem := range childChildChildItem.{{toGoName .TfName}} {
									itemChildChildChildChildBody := ""
									{{- range .Attributes}}
									{{- if .Value}}
									{{- if .ValueCondition}}
									if {{.ValueCondition}} {
										{{- if .PutDataPath}}
										if put {
											itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
										} else {
											itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
										}
										{{- else}}
										itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
										{{- end}}
									}
									{{- else}}
									{{- if .PutDataPath}}
									if put {
										itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
									} else {
										itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
									}
									{{- else}}
									itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
									{{- end}}
									{{- end}}
				{{- else if and (or (not .Reference) .ReferenceInBody) (not .QueryParamNoBody) (not .WoVersion) (eq .Type "Bool") .TypeYangBool (ne .TypeYangBool "boolean")}}
				if childChildChildChildItem.{{toGoName .TfName}}.ValueBool() {
					itemChildChildChildChildBody, _ = sjson.SetRaw(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", "{{if eq .TypeYangBool "empty"}}[null]{{else}}{}{{end}}")
				}
				{{- else if and (or (not .Reference) .ReferenceInBody) (not .QueryParamNoBody) (not .WoVersion)}}
									{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
									if !childChildChildChildItem.{{toGoName .TfName}}.IsNull() {
										{{- if .PutDataPath}}
										if put {
											itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childChildChildChildItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
										} else {
											itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childChildChildChildItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
										}
										{{- else}}
										itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childChildChildChildItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
										{{- end}}
									}
									{{- else if isListSet .}}
									if !childChildChildChildItem.{{toGoName .TfName}}.IsNull() {
										var values []{{if isStringListSet .}}string{{else if isInt64ListSet .}}int64{{end}}
										childChildChildChildItem.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
										{{- if .PutDataPath}}
										if put {
											itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", values)
										} else {
											itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
										}
										{{- else}}
										itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
										{{- end}}
									}
									{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
									if !childChildChildChildItem.{{toGoName .TfName}}.IsNull() {
										{{- if $.NoRead}}
										var listValues map[string]types.List
										childChildChildChildItem.{{toGoName .TfName}}.ElementsAs(ctx, &listValues, false)
										params := make(map[string]any)
										for k, v := range listValues {
											var strList []string
											v.ElementsAs(ctx, &strList, false)
											if len(strList) == 1 {
												params[k] = strList[0]
											} else {
												params[k] = strList
											}
										}
										{{- if .PutDataPath}}
										if put {
											itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", params)
										} else {
											itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
										}
										{{- else}}
										itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
										{{- end}}
										{{- else}}
										var values map[string]string
										childChildChildChildItem.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
										{{- if .PutDataPath}}
										if put {
											itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", values)
										} else {
											itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
										}
										{{- else}}
										itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
										{{- end}}
										{{- end}}
									}
								{{- else if isNestedListSetMap .}}
								{{- if .AlwaysInclude}}
								if true {
								{{- else}}
								if len(childChildChildChildItem.{{toGoName .TfName}}) > 0 {
								{{- end}}
									{{- if .PutDataPath}}
									if put {
										itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", []interface{}{})
									} else {
										itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", []interface{}{})
									}
									{{- else}}
									itemChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", []interface{}{})
									{{- end}}
									for _, childChildChildChildChildItem := range childChildChildChildItem.{{toGoName .TfName}} {
										itemChildChildChildChildChildBody := ""
										{{- range .Attributes}}
										{{- if .Value}}
										{{- if .ValueCondition}}
										if {{.ValueCondition}} {
											{{- if .PutDataPath}}
											if put {
												itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
											} else {
												itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
											}
											{{- else}}
											itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
											{{- end}}
										}
										{{- else}}
										{{- if .PutDataPath}}
										if put {
											itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
										} else {
											itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
										}
										{{- else}}
										itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if eq .Type "String"}}"{{end}}{{.Value}}{{if eq .Type "String"}}"{{end}})
										{{- end}}
										{{- end}}
					{{- else if and (or (not .Reference) .ReferenceInBody) (not .QueryParamNoBody) (not .WoVersion) (eq .Type "Bool") .TypeYangBool (ne .TypeYangBool "boolean")}}
					if childChildChildChildChildItem.{{toGoName .TfName}}.ValueBool() {
						itemChildChildChildChildChildBody, _ = sjson.SetRaw(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", "{{if eq .TypeYangBool "empty"}}[null]{{else}}{}{{end}}")
					}
					{{- else if and (or (not .Reference) .ReferenceInBody) (not .QueryParamNoBody) (not .WoVersion)}}
										{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
										if !childChildChildChildChildItem.{{toGoName .TfName}}.IsNull() {
											{{- if .PutDataPath}}
											if put {
												itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childChildChildChildChildItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
											} else {
												itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childChildChildChildChildItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
											}
											{{- else}}
											itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{if .ModelTypeString}}fmt.Sprint({{end}}childChildChildChildChildItem.{{toGoName .TfName}}.Value{{.Type}}(){{if .ModelTypeString}}){{end}})
											{{- end}}
										}
										{{- else if isListSet .}}
										if !childChildChildChildChildItem.{{toGoName .TfName}}.IsNull() {
											var values []{{if isStringListSet .}}string{{else if isInt64ListSet .}}int64{{end}}
											childChildChildChildChildItem.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
											{{- if .PutDataPath}}
											if put {
												itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", values)
											} else {
												itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
											}
											{{- else}}
											itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
											{{- end}}
										}
										{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
										if !childChildChildChildChildItem.{{toGoName .TfName}}.IsNull() {
											{{- if $.NoRead}}
											var listValues map[string]types.List
											childChildChildChildChildItem.{{toGoName .TfName}}.ElementsAs(ctx, &listValues, false)
											params := make(map[string]any)
											for k, v := range listValues {
												var strList []string
												v.ElementsAs(ctx, &strList, false)
												if len(strList) == 1 {
													params[k] = strList[0]
												} else {
													params[k] = strList
												}
											}
											{{- if .PutDataPath}}
											if put {
												itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", params)
											} else {
												itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
											}
											{{- else}}
											itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", params)
											{{- end}}
											{{- else}}
											var values map[string]string
											childChildChildChildChildItem.{{toGoName .TfName}}.ElementsAs(ctx, &values, false)
											{{- if .PutDataPath}}
											if put {
												itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}", values)
											} else {
												itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
											}
											{{- else}}
											itemChildChildChildChildChildBody, _ = sjson.Set(itemChildChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", values)
											{{- end}}
											{{- end}}
										}
										{{- end}}
										{{- end}}
										{{- end}}
										{{- if .PutDataPath}}
										if put {
											itemChildChildChildChildBody, _ = sjson.SetRaw(itemChildChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}.-1", itemChildChildChildChildChildBody)
										} else {
											itemChildChildChildChildBody, _ = sjson.SetRaw(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}.-1", itemChildChildChildChildChildBody)
										}
										{{- else}}
										itemChildChildChildChildBody, _ = sjson.SetRaw(itemChildChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}.-1", itemChildChildChildChildChildBody)
										{{- end}}
									}
								}
									{{- end}}
									{{- end}}
									{{- end}}
									{{- if .PutDataPath}}
									if put {
										itemChildChildChildBody, _ = sjson.SetRaw(itemChildChildChildBody, "{{.PutDataPath}}.{{.ModelName}}.-1", itemChildChildChildChildBody)
									} else {
										itemChildChildChildBody, _ = sjson.SetRaw(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}.-1", itemChildChildChildChildBody)
									}
									{{- else}}
									itemChildChildChildBody, _ = sjson.SetRaw(itemChildChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}.-1", itemChildChildChildChildBody)
									{{- end}}
								}
							}
								{{- end}}
								{{- end}}
								{{- end}}
								{{- if .PutDataPath}}
								if put {
									itemChildChildBody, _ = sjson.SetRaw(itemChildChildBody, "{{.PutDataPath}}.{{.ModelName}}.-1", itemChildChildChildBody)
								} else {
									itemChildChildBody, _ = sjson.SetRaw(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}.-1", itemChildChildChildBody)
								}
								{{- else}}
								itemChildChildBody, _ = sjson.SetRaw(itemChildChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}.-1", itemChildChildChildBody)
								{{- end}}
							}
						}
							{{- end}}
							{{- end}}
							{{- end}}
							{{- if .PutDataPath}}
							if put {
								itemChildBody, _ = sjson.SetRaw(itemChildBody, "{{.PutDataPath}}.{{.ModelName}}.-1", itemChildChildBody)
							} else {
								itemChildBody, _ = sjson.SetRaw(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}.-1", itemChildChildBody)
							}
							{{- else}}
							itemChildBody, _ = sjson.SetRaw(itemChildBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}.-1", itemChildChildBody)
							{{- end}}
						}
					}
					{{- end}}
					{{- end}}
					{{- end}}
					{{- if .PutDataPath}}
					if put {
						itemBody, _ = sjson.SetRaw(itemBody, "{{.PutDataPath}}.{{.ModelName}}.-1", itemChildBody)
					} else {
						itemBody, _ = sjson.SetRaw(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}.-1", itemChildBody)
					}
					{{- else}}
					itemBody, _ = sjson.SetRaw(itemBody, "{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}.-1", itemChildBody)
					{{- end}}
				}
			}
			{{- end}}
			{{- end}}
			{{- end}}
			{{- if .PutDataPath}}
			if put {
				body, _ = sjson.SetRaw(body, "{{.PutDataPath}}.{{if .ModelName}}{{.ModelName}}.{{end}}-1", itemBody)
			} else {
				body, _ = sjson.SetRaw(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{if .ModelName}}{{.ModelName}}.{{end}}-1", itemBody)
			}
			{{- else}}
			body, _ = sjson.SetRaw(body, "{{if .DataPath}}{{.DataPath}}.{{end}}{{if .ModelName}}{{.ModelName}}.{{end}}-1", itemBody)
			{{- end}}
		}
	}
	{{- end}}
	{{- end}}
	{{- end}}
	return body
}
// End of section. //template:end toBody

// Section below is generated&owned by "gen/generator.go". //template:begin fromBody
func (data *{{camelCase .Name}}) fromBody(ctx context.Context, res gjson.Result) {
	{{- if .RootList}}
	{{- range .Attributes}}
	{{if .ResponseDataPath}}
	res = res.Get("{{.ResponseDataPath}}")
	{{- end}}
	{{- end}}
	{{- end}}
	{{- if .DataSourceNoId}}
	{{- if hasId .Attributes}}
		{{- $id := getId .Attributes}}
	data.Id = types.StringValue(fmt.Sprint(data.{{toGoName $id.TfName}}.Value{{$id.Type}}()))
	{{- else if constantId .}}
	data.Id = types.StringValue("{{constantId .}}")
	{{- else }}
	if value := res.Get("{{if .IdFromQueryPath}}{{if eq .IdFromQueryPath "." }}{{else}}{{.IdFromQueryPath}}.{{end}}{{if .IdFromQueryPathAttribute}}{{.IdFromQueryPathAttribute}}{{else}}id{{end}}{{end}}"); value.Exists() {
		data.Id = types.StringValue(value.String())
	} else {
		data.Id = types.StringNull()
	}
	{{- end}}
	{{- end}}
	{{- range .Attributes}}
	{{- if and .WriteOnly .ExcludeTest (eq .Type "Bool") (not .QueryParamNoBody) }}
	if value := res.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"); value.Exists() && !data.{{toGoName .TfName}}.IsNull() {
		data.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
	} else {
		data.{{toGoName .TfName}} = types.{{.Type}}Value(false)
	}
	{{- else if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference) (not .CreateQueryPath) (not .QueryParamNoBody)}}
	{{- $cname := toGoName .TfName}}
	{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
	if value := res.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"); value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && value.String() != ""{{end}} {
		data.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
	}{{if or .FallbackResponseDataPath .FallbackResponseModelName}} else if value := res.Get("{{if .FallbackResponseDataPath}}{{.FallbackResponseDataPath}}{{else}}{{.FallbackResponseModelName}}{{end}}"); value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && value.String() != ""{{end}} {
		data.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
	}{{end}} else {
		{{- if .DefaultValue}}
		data.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}})
		{{- else}}
		data.{{toGoName .TfName}} = types.{{.Type}}Null()
		{{- end}}
	}
	{{- else if isListSet .}}
	if value := res.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"); value.Exists() && len(value.Array()) > 0 {
		data.{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
	}{{if or .FallbackResponseDataPath .FallbackResponseModelName}} else if value := res.Get("{{if .FallbackResponseDataPath}}{{.FallbackResponseDataPath}}{{else}}{{.FallbackResponseModelName}}{{end}}"); value.Exists() && len(value.Array()) > 0 {
		data.{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
	}{{end}} else {
		data.{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
	}
	{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
	if value := res.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"); value.Exists() && len(value.Map()) > 0 {
		data.{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
	} else {
		data.{{toGoName .TfName}} = types.MapNull(types.StringType)
	}
	{{- else if isNestedListSetMap .}}
	if value := res{{if .ModelName}}.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"){{end}}; value.Exists() && len(value.Array()) > 0 {
		{{- if isNestedMap .}}
		data.{{toGoName .TfName}} = make(map[string]{{$name}}{{toGoName .TfName}})
		{{- else}}
		data.{{toGoName .TfName}} = make([]{{$name}}{{toGoName .TfName}}, 0)
		{{- end}}
		value.ForEach(func(k, v gjson.Result) bool {
			item := {{$name}}{{toGoName .TfName}}{}
			{{- range .Attributes}}
			{{- $ccname := toGoName .TfName}}
			{{- if and .WriteOnly .Mandatory }}
			item.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Bool"}}false{{else if eq .Type "Int64"}}0{{else if eq .Type "String"}}""{{else}}{{.Type}}{{end}})
			{{- end}}
			{{- if and .WriteOnly .ExcludeTest (eq .Type "Bool") (not .QueryParamNoBody) }}
			if cValue := v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); cValue.Exists() {
				item.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "cValue" .}}{{else if isYangFlag .}}true{{else}}cValue.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
			} else {
				item.{{toGoName .TfName}} = types.{{.Type}}Value(false)
			}
			{{- else if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
			{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
			if cValue := v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); cValue.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && cValue.String() != ""{{end}} {
				item.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "cValue" .}}{{else if isYangFlag .}}true{{else}}cValue.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
			} else {
				{{- if .DefaultValue}}
				item.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}})
				{{- else}}
				item.{{toGoName .TfName}} = types.{{.Type}}Null()
				{{- end}}
			}
			{{- else if isListSet .}}
			if cValue := v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); cValue.Exists() && len(cValue.Array()) > 0 {
				item.{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(cValue.Array())
			} else {
				item.{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
			}
			{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
			if cValue := v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); cValue.Exists() && len(cValue.Map()) > 0 {
				item.{{toGoName .TfName}} = helpers.GetStringMap(cValue.Map())
			} else {
				item.{{toGoName .TfName}} = types.MapNull(types.StringType)
			}
			{{- else if isNestedListSetMap .}}
			if cValue := v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); cValue.Exists() && len(cValue.Array()) > 0 {
				item.{{toGoName .TfName}} = make([]{{$name}}{{$cname}}{{toGoName .TfName}}, 0)
				cValue.ForEach(func(ck, cv gjson.Result) bool {
					cItem := {{$name}}{{$cname}}{{toGoName .TfName}}{}
					{{- range .Attributes}}
					{{- if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
					{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
					if cwValue := cv.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); cwValue.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && cwValue.String() != ""{{end}} {
						cItem.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "cwValue" .}}{{else if isYangFlag .}}true{{else}}cwValue.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
					} else {
						{{- if .DefaultValue}}
						cItem.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}})
						{{- else}}
						cItem.{{toGoName .TfName}} = types.{{.Type}}Null()
						{{- end}}
					}
					{{- else if isListSet .}}
					if cwValue := cv.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); cwValue.Exists() && len(cwValue.Array()) > 0 {
						cItem.{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(cwValue.Array())
					} else {
						cItem.{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
					}
					{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
					if cwValue := cv.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); cwValue.Exists() && len(cwValue.Map()) > 0 {
						cItem.{{toGoName .TfName}} = helpers.GetStringMap(cwValue.Map())
					} else {
						cItem.{{toGoName .TfName}} = types.MapNull(types.StringType)
					}
					{{- else if isNestedListSetMap .}}
{{- $c3name := toGoName .TfName}}
					if cwValue := cv.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); cwValue.Exists() && len(cwValue.Array()) > 0 {
						cItem.{{toGoName .TfName}} = make([]{{$name}}{{$cname}}{{$ccname}}{{toGoName .TfName}}, 0)
						cwValue.ForEach(func(cck, ccv gjson.Result) bool {
							cwItem := {{$name}}{{$cname}}{{$ccname}}{{toGoName .TfName}}{}
							{{- range .Attributes}}
							{{- if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
							{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
							if ccwValue := ccv.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); ccwValue.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && ccwValue.String() != ""{{end}} {
								cwItem.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "ccwValue" .}}{{else if isYangFlag .}}true{{else}}ccwValue.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
							} else {
								{{- if .DefaultValue}}
								cwItem.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}})
								{{- else}}
								cwItem.{{toGoName .TfName}} = types.{{.Type}}Null()
								{{- end}}
							}
							{{- else if isListSet .}}
							if ccwValue := ccv.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); ccwValue.Exists() && len(ccwValue.Array()) > 0 {
								cwItem.{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(ccwValue.Array())
							} else {
								cwItem.{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
							}
							{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
							if ccwValue := ccv.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); ccwValue.Exists() && len(ccwValue.Map()) > 0 {
								cwItem.{{toGoName .TfName}} = helpers.GetStringMap(ccwValue.Map())
							} else {
								cwItem.{{toGoName .TfName}} = types.MapNull(types.StringType)
							}
						{{- else if isNestedListSetMap .}}
{{- $c4name := toGoName .TfName}}
						if ccwValue := ccv.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); ccwValue.Exists() && len(ccwValue.Array()) > 0 {
							cwItem.{{toGoName .TfName}} = make([]{{$name}}{{$cname}}{{$ccname}}{{$c3name}}{{toGoName .TfName}}, 0)
							ccwValue.ForEach(func(c3k, c3v gjson.Result) bool {
								c3Item := {{$name}}{{$cname}}{{$ccname}}{{$c3name}}{{toGoName .TfName}}{}
								{{- range .Attributes}}
								{{- if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
								{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
								if c3Value := c3v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); c3Value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && c3Value.String() != ""{{end}} {
									c3Item.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "c3Value" .}}{{else if isYangFlag .}}true{{else}}c3Value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
								} else {
									{{- if .DefaultValue}}
									c3Item.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}})
									{{- else}}
									c3Item.{{toGoName .TfName}} = types.{{.Type}}Null()
									{{- end}}
								}
								{{- else if isListSet .}}
								if c3Value := c3v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); c3Value.Exists() && len(c3Value.Array()) > 0 {
									c3Item.{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(c3Value.Array())
								} else {
									c3Item.{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
								}
								{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
								if c3Value := c3v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); c3Value.Exists() && len(c3Value.Map()) > 0 {
									c3Item.{{toGoName .TfName}} = helpers.GetStringMap(c3Value.Map())
								} else {
									c3Item.{{toGoName .TfName}} = types.MapNull(types.StringType)
								}
							{{- else if isNestedListSetMap .}}
{{- $c5name := toGoName .TfName}}
							if c3Value := c3v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); c3Value.Exists() && len(c3Value.Array()) > 0 {
								c3Item.{{toGoName .TfName}} = make([]{{$name}}{{$cname}}{{$ccname}}{{$c3name}}{{$c4name}}{{toGoName .TfName}}, 0)
								c3Value.ForEach(func(c4k, c4v gjson.Result) bool {
									c4Item := {{$name}}{{$cname}}{{$ccname}}{{$c3name}}{{$c4name}}{{toGoName .TfName}}{}
									{{- range .Attributes}}
									{{- if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
									{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
									if c4Value := c4v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); c4Value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && c4Value.String() != ""{{end}} {
										c4Item.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "c4Value" .}}{{else if isYangFlag .}}true{{else}}c4Value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
									} else {
										{{- if .DefaultValue}}
										c4Item.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}})
										{{- else}}
										c4Item.{{toGoName .TfName}} = types.{{.Type}}Null()
										{{- end}}
									}
									{{- else if isListSet .}}
									if c4Value := c4v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); c4Value.Exists() && len(c4Value.Array()) > 0 {
										c4Item.{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(c4Value.Array())
									} else {
										c4Item.{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
									}
									{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
									if c4Value := c4v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); c4Value.Exists() && len(c4Value.Map()) > 0 {
										c4Item.{{toGoName .TfName}} = helpers.GetStringMap(c4Value.Map())
									} else {
										c4Item.{{toGoName .TfName}} = types.MapNull(types.StringType)
									}
								{{- else if isNestedListSetMap .}}
{{- $c6name := toGoName .TfName}}
								if c4Value := c4v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); c4Value.Exists() && len(c4Value.Array()) > 0 {
									c4Item.{{toGoName .TfName}} = make([]{{$name}}{{$cname}}{{$ccname}}{{$c3name}}{{$c4name}}{{$c5name}}{{toGoName .TfName}}, 0)
									c4Value.ForEach(func(c5k, c5v gjson.Result) bool {
										c5Item := {{$name}}{{$cname}}{{$ccname}}{{$c3name}}{{$c4name}}{{$c5name}}{{toGoName .TfName}}{}
										{{- range .Attributes}}
										{{- if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
										{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
										if c5Value := c5v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); c5Value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && c5Value.String() != ""{{end}} {
											c5Item.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "c5Value" .}}{{else if isYangFlag .}}true{{else}}c5Value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
										} else {
											{{- if .DefaultValue}}
											c5Item.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}})
											{{- else}}
											c5Item.{{toGoName .TfName}} = types.{{.Type}}Null()
											{{- end}}
										}
										{{- else if isListSet .}}
										if c5Value := c5v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); c5Value.Exists() && len(c5Value.Array()) > 0 {
											c5Item.{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(c5Value.Array())
										} else {
											c5Item.{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
										}
										{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
										if c5Value := c5v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); c5Value.Exists() && len(c5Value.Map()) > 0 {
											c5Item.{{toGoName .TfName}} = helpers.GetStringMap(c5Value.Map())
										} else {
											c5Item.{{toGoName .TfName}} = types.MapNull(types.StringType)
										}
										{{- end}}
										{{- end}}
										{{- end}}
										c4Item.{{toGoName .TfName}} = append(c4Item.{{toGoName .TfName}}, c5Item)
										return true
									})
								}
									{{- end}}
									{{- end}}
									{{- end}}
									c3Item.{{toGoName .TfName}} = append(c3Item.{{toGoName .TfName}}, c4Item)
									return true
								})
							}
								{{- end}}
								{{- end}}
								{{- end}}
								cwItem.{{toGoName .TfName}} = append(cwItem.{{toGoName .TfName}}, c3Item)
								return true
							})
						}
							{{- end}}
							{{- end}}
							{{- end}}
							cItem.{{toGoName .TfName}} = append(cItem.{{toGoName .TfName}}, cwItem)
							return true
						})
					}
					{{- end}}
					{{- end}}
					{{- end}}
					item.{{toGoName .TfName}} = append(item.{{toGoName .TfName}}, cItem)
					return true
				})
			}
			{{- end}}
			{{- end}}
			{{- end}}
			{{- if isNestedMap .}}
			key := {{range $i, $kp := getKeyPartAttributes .Attributes}}{{if $i}} + "_" + {{end}}{{if eq $kp.Type "Int64"}}fmt.Sprintf("%v", item.{{toGoName $kp.TfName}}.ValueInt64()){{else if eq $kp.Type "Bool"}}fmt.Sprintf("%v", item.{{toGoName $kp.TfName}}.ValueBool()){{else}}item.{{toGoName $kp.TfName}}.ValueString(){{end}}{{end}}
			data.{{toGoName .TfName}}[key] = item
			{{- else}}
			data.{{toGoName .TfName}} = append(data.{{toGoName .TfName}}, item)
			{{- end}}
			return true
		})
	}
	{{- end}}
	{{- end}}
	{{- end}}
}
// End of section. //template:end fromBody

// Section below is generated&owned by "gen/generator.go". //template:begin updateFromBody
func (data *{{camelCase .Name}}) updateFromBody(ctx context.Context, res gjson.Result) {
	{{- $chunkElementName := "" }} 
	{{- $chunkIsNestedMap := false }}
	{{- range  .Attributes}}
	{{- if .MaxElementsInRootList}}
	{{- $chunkElementName = .TfName }}
	{{- $chunkIsNestedMap = isNestedMap . }}
	{{- if isNestedMap .}}
	final := make(map[string]{{$name}}{{toGoName .TfName}})
	{{- else}}
	var final []{{$name}}{{toGoName .TfName}}
	{{- end}}
	{{- end}}
	{{- end}}
	{{- if .RootList}}
	{{- range .Attributes}}
	{{if .ResponseDataPath}}
	res = res.Get("{{.ResponseDataPath}}")
	{{- end}}
	{{- end}}
	{{- end}}
	{{- range .Attributes}}
	{{- if and .WriteOnly .ExcludeTest (eq .Type "Bool") (not .QueryParamNoBody) }}
	if value := res.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"); value.Exists() && !data.{{toGoName .TfName}}.IsNull() {
		data.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
	} else {
		data.{{toGoName .TfName}} = types.{{.Type}}Value(false)
	}
	{{- else if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference) (not .CreateQueryPath) (not .QueryParamNoBody)}}
	{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
	if value := res.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"); value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && value.String() != ""{{end}} && !data.{{toGoName .TfName}}.IsNull() {
		data.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
	}{{if or .FallbackResponseDataPath .FallbackResponseModelName}} else if value := res.Get("{{if .FallbackResponseDataPath}}{{.FallbackResponseDataPath}}{{else}}{{.FallbackResponseModelName}}{{end}}"); value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && value.String() != ""{{end}} && !data.{{toGoName .TfName}}.IsNull() {
		data.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
	}{{end}} else {{if .DefaultValue}}if data.{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
		data.{{toGoName .TfName}} = types.{{.Type}}Null()
	}
	{{- else if isListSet .}}
	if value := res.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"); value.Exists() && !data.{{toGoName .TfName}}.IsNull() {
		{{- if and .IgnoreExtraResponseValues (eq .Type "Set") (eq .ElementType "String")}}
		data.{{toGoName .TfName}} = helpers.KeepStringSetIn(ctx, data.{{toGoName .TfName}}, value.Array())
		{{- else}}
		data.{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
		{{- end}}
	}{{if or .FallbackResponseDataPath .FallbackResponseModelName}} else if value := res.Get("{{if .FallbackResponseDataPath}}{{.FallbackResponseDataPath}}{{else}}{{.FallbackResponseModelName}}{{end}}"); value.Exists() && !data.{{toGoName .TfName}}.IsNull() {
		{{- if and .IgnoreExtraResponseValues (eq .Type "Set") (eq .ElementType "String")}}
		data.{{toGoName .TfName}} = helpers.KeepStringSetIn(ctx, data.{{toGoName .TfName}}, value.Array())
		{{- else}}
		data.{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
		{{- end}}
	}{{end}} else {
		data.{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
	}
	{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
	if value := res.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"); value.Exists() && !data.{{toGoName .TfName}}.IsNull() {
		data.{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
	}{{if or .FallbackResponseDataPath .FallbackResponseModelName}} else if value := res.Get("{{if .FallbackResponseDataPath}}{{.FallbackResponseDataPath}}{{else}}{{.FallbackResponseModelName}}{{end}}"); value.Exists() && !data.{{toGoName .TfName}}.IsNull() {
		data.{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
	}{{end}} else {
		data.{{toGoName .TfName}} = types.MapNull(types.StringType)
	}
	{{- else if isNestedMap .}}
	if value := res{{if .ModelName}}.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"){{end}}; value.Exists() {
		for mapKey, mapItem := range data.{{toGoName .TfName}} {
			value.ForEach(func(_, v gjson.Result) bool {
				computedKey := {{range $i, $kp := getKeyPartAttributes .Attributes}}{{if $i}} + "_" + {{end}}{{if eq $kp.Type "Int64"}}fmt.Sprintf("%v", v.Get("{{if $kp.DataPath}}{{$kp.DataPath}}.{{end}}{{$kp.ModelName}}").Int()){{else if eq $kp.Type "Bool"}}fmt.Sprintf("%v", v.Get("{{if $kp.DataPath}}{{$kp.DataPath}}.{{end}}{{$kp.ModelName}}").Bool()){{else}}v.Get("{{if $kp.DataPath}}{{$kp.DataPath}}.{{end}}{{$kp.ModelName}}").String(){{end}}{{end}}
				if computedKey != mapKey {
					return true
				}
				{{- range .Attributes}}
				{{- if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
				{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
				if fv := v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); fv.Exists() && !mapItem.{{toGoName .TfName}}.IsNull() {
					mapItem.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "fv" .}}{{else}}fv.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
				} else {{if .DefaultValue}}if mapItem.{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
					mapItem.{{toGoName .TfName}} = types.{{.Type}}Null()
				}
				{{- else if isListSet .}}
				if fv := v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); fv.Exists() && !mapItem.{{toGoName .TfName}}.IsNull() {
					mapItem.{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(fv.Array())
				} else {
					mapItem.{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
				}
				{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
				if fv := v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); fv.Exists() && !mapItem.{{toGoName .TfName}}.IsNull() {
					mapItem.{{toGoName .TfName}} = helpers.GetStringMap(fv.Map())
				} else {
					mapItem.{{toGoName .TfName}} = types.MapNull(types.StringType)
				}
				{{- end}}
				{{- end}}
				{{- end}}
				{{- if ne $chunkElementName ""}}
				if !mapItem.Id.IsNull() {
					final[mapKey] = mapItem
				}
				{{- else}}
				data.{{toGoName .TfName}}[mapKey] = mapItem
				{{- end}}
				return false
			})
		}
	}
	{{- else if isNestedListSetMap .}}
	{{- $list := (toGoName .TfName)}}
	for i := range data.{{toGoName .TfName}} {
		{{- if and (isNestedListSet .) (eq .MaxList 1)}}
		values := res{{if .ModelName}}.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"){{end}}.Array()
		if len(values) != 1 {
			data.{{$list}} = nil
			break
		}
		r := values[0]
		{{- else}}
		keys := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if or (eq .Type "Int64") (eq .Type "Bool") (eq .Type "String")}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}}{{end}} }
		keyValues := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(data.{{$list}}[i].{{toGoName .TfName}}.ValueInt64(), 10), {{else if eq .Type "Bool"}}strconv.FormatBool(data.{{$list}}[i].{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}data.{{$list}}[i].{{toGoName .TfName}}.Value{{.Type}}(), {{end}}{{end}}{{end}} }
		{{- if hasComputedRefreshValueInKeys .Attributes}}
		emptyKeys := [...]string{ {{range .Attributes}}{{if .ComputedRefreshValue}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}} }
		{{- end}}

		var r gjson.Result
		res.{{if .ModelName}}Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}").{{end}}ForEach(
			func(_, v gjson.Result) bool {
				found := false
				for ik := range keys {
					{{- if hasComputedRefreshValueInKeys .Attributes}}
					if v.Get(keys[ik]).String() == keyValues[ik] || slices.Contains(emptyKeys[:], keys[ik]) && keyValues[ik] == "" {
					{{- else}}
					if v.Get(keys[ik]).String() == keyValues[ik] {
					{{- end}}
						found = true
						continue
					}
					found = false
					break
				}
				if found {
					r = v
					return false
				}
				return true
			},
		)
		{{- end}}
		{{- if $.SkipUnmatchedRows}}
		if !r.Exists() {
			continue
		}
		{{- end}}

		{{- range .Attributes}}
		{{- if and .WriteOnly .ExcludeTest (eq .Type "Bool") }}
		if value := r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{toGoName .TfName}}.IsNull() {
			data.{{$list}}[i].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
		} else {{if .DefaultValue}}if data.{{$list}}[i].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
			data.{{$list}}[i].{{toGoName .TfName}} = types.{{.Type}}Value(false)
		}		
		{{- else if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
		{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
		if value := r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && value.String() != ""{{end}} && !data.{{$list}}[i].{{toGoName .TfName}}.IsNull() {
			data.{{$list}}[i].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
		} else {{if .DefaultValue}}if data.{{$list}}[i].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
			data.{{$list}}[i].{{toGoName .TfName}} = types.{{.Type}}Null()
		}
		{{- else if isListSet .}}
		if value := r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{toGoName .TfName}}.IsNull() {
			data.{{$list}}[i].{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
		} else {
			data.{{$list}}[i].{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
		}
		{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
		if value := r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{toGoName .TfName}}.IsNull() {
			data.{{$list}}[i].{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
		} else {
			data.{{$list}}[i].{{toGoName .TfName}} = types.MapNull(types.StringType)
		}
		{{- else if isNestedListSetMap .}}
		{{- $clist := (toGoName .TfName)}}
		for ci := range data.{{$list}}[i].{{toGoName .TfName}} {
			{{- if and (isNestedListSet .) (eq .MaxList 1)}}
			values := r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").Array()
			if len(values) != 1 {
				data.{{$list}}[i].{{$clist}} = nil
				break
			}
			cr := values[0]
			{{- else}}
			keys := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if or (eq .Type "Int64") (eq .Type "Bool") (eq .Type "String")}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}}{{end}} }
			keyValues := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.ValueInt64(), 10), {{else if eq .Type "Bool"}}strconv.FormatBool(data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.Value{{.Type}}(), {{end}}{{end}}{{end}} }
			{{- if hasComputedRefreshValueInKeys .Attributes}}
			emptyKeys := [...]string{ {{range .Attributes}}{{if .ComputedRefreshValue}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}} }
			{{- end}}

			var cr gjson.Result
			r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").ForEach(
				func(_, v gjson.Result) bool {
					found := false
					for ik := range keys {
						{{- if hasComputedRefreshValueInKeys .Attributes}}
						if v.Get(keys[ik]).String() == keyValues[ik] || slices.Contains(emptyKeys[:], keys[ik]) && keyValues[ik] == "" {
						{{- else}}
						if v.Get(keys[ik]).String() == keyValues[ik] {
						{{- end}}
							found = true
							continue
						}
						found = false
						break
					}
					if found {
						cr = v
						return false
					}
					return true
				},
			)
			{{- end}}

			{{- range .Attributes}}
			{{- if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
			{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
			if value := cr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && value.String() != ""{{end}} && !data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.IsNull() {
				data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
			} else {{if .DefaultValue}}if data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
				data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} = types.{{.Type}}Null()
			}
			{{- else if isListSet .}}
			if value := cr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.IsNull() {
				data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
			} else {
				data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
			}
			{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
			if value := cr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.IsNull() {
				data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
			} else {
				data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} = types.MapNull(types.StringType)
			}
			{{- else if isNestedListSetMap .}}
			{{- $cclist := (toGoName .TfName)}}
			for cci := range data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} {
				{{- if and (isNestedListSet .) (eq .MaxList 1)}}
				values := cr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").Array()
				if len(values) != 1 {
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}} = nil
					break
				}
				ccr := values[0]
				{{- else}}
				keys := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if or (eq .Type "Int64") (eq .Type "Bool") (eq .Type "String")}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}}{{end}} }
				keyValues := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.ValueInt64(), 10), {{else if eq .Type "Bool"}}strconv.FormatBool(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.Value{{.Type}}(), {{end}}{{end}}{{end}} }
				{{- if hasComputedRefreshValueInKeys .Attributes}}
				emptyKeys := [...]string{ {{range .Attributes}}{{if .ComputedRefreshValue}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}} }
				{{- end}}

				var ccr gjson.Result
				cr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").ForEach(
					func(_, v gjson.Result) bool {
						found := false
						for ik := range keys {
							{{- if hasComputedRefreshValueInKeys .Attributes}}
							if v.Get(keys[ik]).String() == keyValues[ik] || slices.Contains(emptyKeys[:], keys[ik]) && keyValues[ik] == "" {
							{{- else}}
							if v.Get(keys[ik]).String() == keyValues[ik] {
							{{- end}}
								found = true
								continue
							}
							found = false
							break
						}
						if found {
							ccr = v
							return false
						}
						return true
					},
				)
				{{- end}}

				{{- range .Attributes}}
				{{- if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
				{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
				if value := ccr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && value.String() != ""{{end}} && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.IsNull() {
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
				} else {{if .DefaultValue}}if data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} = types.{{.Type}}Null()
				}
				{{- else if isListSet .}}
				if value := ccr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.IsNull() {
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
				} else {
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
				}
				{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
				if value := ccr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.IsNull() {
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
				} else {
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} = types.MapNull(types.StringType)
				}
				{{- else if isNestedListSetMap .}}
				{{- $c3list := (toGoName .TfName)}}
				for c3i := range data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} {
					{{- if and (isNestedListSet .) (eq .MaxList 1)}}
					values := ccr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").Array()
					if len(values) != 1 {
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}} = nil
						break
					}
					c3r := values[0]
					{{- else}}
					keys := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if or (eq .Type "Int64") (eq .Type "Bool") (eq .Type "String")}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}}{{end}} }
					keyValues := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.ValueInt64(), 10), {{else if eq .Type "Bool"}}strconv.FormatBool(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.Value{{.Type}}(), {{end}}{{end}}{{end}} }
					{{- if hasComputedRefreshValueInKeys .Attributes}}
					emptyKeys := [...]string{ {{range .Attributes}}{{if .ComputedRefreshValue}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}} }
					{{- end}}

					var c3r gjson.Result
					ccr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").ForEach(
						func(_, v gjson.Result) bool {
							found := false
							for ik := range keys {
								{{- if hasComputedRefreshValueInKeys .Attributes}}
								if v.Get(keys[ik]).String() == keyValues[ik] || slices.Contains(emptyKeys[:], keys[ik]) && keyValues[ik] == "" {
								{{- else}}
								if v.Get(keys[ik]).String() == keyValues[ik] {
								{{- end}}
									found = true
									continue
								}
								found = false
								break
							}
							if found {
								c3r = v
								return false
							}
							return true
						},
					)
					{{- end}}

					{{- range .Attributes}}
					{{- if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
					{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
					if value := c3r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && value.String() != ""{{end}} && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.IsNull() {
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
					} else {{if .DefaultValue}}if data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} = types.{{.Type}}Null()
					}
					{{- else if isListSet .}}
					if value := c3r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.IsNull() {
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
					} else {
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
					}
					{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
					if value := c3r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.IsNull() {
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
					} else {
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} = types.MapNull(types.StringType)
					}
					{{- else if isNestedListSetMap .}}
					{{- $c4list := (toGoName .TfName)}}
					for c4i := range data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} {
						{{- if and (isNestedListSet .) (eq .MaxList 1)}}
						values := c3r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").Array()
						if len(values) != 1 {
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}} = nil
							break
						}
						c4r := values[0]
						{{- else}}
						keys := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if or (eq .Type "Int64") (eq .Type "Bool") (eq .Type "String")}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}}{{end}} }
						keyValues := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.ValueInt64(), 10), {{else if eq .Type "Bool"}}strconv.FormatBool(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.Value{{.Type}}(), {{end}}{{end}}{{end}} }
						{{- if hasComputedRefreshValueInKeys .Attributes}}
						emptyKeys := [...]string{ {{range .Attributes}}{{if .ComputedRefreshValue}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}} }
						{{- end}}

						var c4r gjson.Result
						c3r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").ForEach(
							func(_, v gjson.Result) bool {
								found := false
								for ik := range keys {
									{{- if hasComputedRefreshValueInKeys .Attributes}}
									if v.Get(keys[ik]).String() == keyValues[ik] || slices.Contains(emptyKeys[:], keys[ik]) && keyValues[ik] == "" {
									{{- else}}
									if v.Get(keys[ik]).String() == keyValues[ik] {
									{{- end}}
										found = true
										continue
									}
									found = false
									break
								}
								if found {
									c4r = v
									return false
								}
								return true
							},
						)
						{{- end}}

						{{- range .Attributes}}
						{{- if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
						{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
						if value := c4r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && value.String() != ""{{end}} && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.IsNull() {
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
						} else {{if .DefaultValue}}if data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} = types.{{.Type}}Null()
						}
						{{- else if isListSet .}}
						if value := c4r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.IsNull() {
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
						} else {
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
						}
						{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
						if value := c4r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.IsNull() {
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
						} else {
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} = types.MapNull(types.StringType)
						}
						{{- else if isNestedListSetMap .}}
						{{- $c5list := (toGoName .TfName)}}
						for c5i := range data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} {
							{{- if and (isNestedListSet .) (eq .MaxList 1)}}
							values := c4r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").Array()
							if len(values) != 1 {
								data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}} = nil
								break
							}
							c5r := values[0]
							{{- else}}
							keys := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if or (eq .Type "Int64") (eq .Type "Bool") (eq .Type "String")}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}}{{end}} }
							keyValues := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.ValueInt64(), 10), {{else if eq .Type "Bool"}}strconv.FormatBool(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.Value{{.Type}}(), {{end}}{{end}}{{end}} }
							{{- if hasComputedRefreshValueInKeys .Attributes}}
							emptyKeys := [...]string{ {{range .Attributes}}{{if .ComputedRefreshValue}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}} }
							{{- end}}

							var c5r gjson.Result
							c4r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").ForEach(
								func(_, v gjson.Result) bool {
									found := false
									for ik := range keys {
										{{- if hasComputedRefreshValueInKeys .Attributes}}
										if v.Get(keys[ik]).String() == keyValues[ik] || slices.Contains(emptyKeys[:], keys[ik]) && keyValues[ik] == "" {
										{{- else}}
										if v.Get(keys[ik]).String() == keyValues[ik] {
										{{- end}}
											found = true
											continue
										}
										found = false
										break
									}
									if found {
										c5r = v
										return false
									}
									return true
								},
							)
							{{- end}}

							{{- range .Attributes}}
							{{- if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference)}}
							{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
							if value := c5r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists(){{if and .NullOnEmpty (eq .Type "String")}} && value.String() != ""{{end}} && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.IsNull() {
								data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
							} else {{if .DefaultValue}}if data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
								data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}} = types.{{.Type}}Null()
							}
							{{- else if isListSet .}}
							if value := c5r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.IsNull() {
								data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
							} else {
								data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
							}
							{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
							if value := c5r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.IsNull() {
								data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
							} else {
								data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}} = types.MapNull(types.StringType)
							}
							{{- end}}
							{{- end}}
							{{- end}}
						}
						{{- end}}
						{{- end}}
						{{- end}}
					}
					{{- end}}
					{{- end}}
					{{- end}}
				}
				{{- end}}
				{{- end}}
				{{- end}}
			}

			{{- end}}
			{{- end}}
			{{- end}}
		}
		{{- end}}
		{{- end}}
		{{- end}}
		{{- if ne $chunkElementName "" }}
		{{- if $chunkIsNestedMap }}
		if data.{{toGoName $chunkElementName}}[i].Id != types.StringNull() {
			final[i] = data.{{toGoName $chunkElementName}}[i]
		}
		{{- else}}
		if data.{{toGoName $chunkElementName}}[i].Id != types.StringNull() {
			final = append(final, data.{{toGoName $chunkElementName}}[i])
		}
		{{- end}}
		{{- end}}
	}
	{{- end}}
	{{- end}}
	{{- end}}
	{{- if $chunkElementName }} 
	data.{{toGoName $chunkElementName}} = final
	{{- end}}
}
// End of section. //template:end updateFromBody

{{if .UpdateComputed}}
// Section below is generated&owned by "gen/generator.go". //template:begin fromBodyUnknowns

// Known values are not changed (usual for Computed attributes with UseStateForUnknown or a Default).
func (data *{{camelCase .Name}}) fromBodyUnknowns(ctx context.Context, res gjson.Result) {
{{- if .RootList}}
{{- range .Attributes}}
{{if .ResponseDataPath}}
res = res.Get("{{.ResponseDataPath}}")
{{- end}}
{{- end}}
{{- end}}
{{- range .Attributes}}
{{- if and (not .Value) (not .WriteOnly) (not .WoVersion) (not .Reference) (not .CreateQueryPath) (not .QueryParamNoBody)}}
{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
if data.{{toGoName .TfName}}.IsUnknown() {
	if value := res.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"); value.Exists() && !data.{{toGoName .TfName}}.IsNull() {
		data.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
	} else {{if .DefaultValue}}if data.{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
		data.{{toGoName .TfName}} = types.{{.Type}}Null()
	}
}
{{- else if isListSet .}}
if value := res.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"); value.Exists() && !data.{{toGoName .TfName}}.IsNull() {
	data.{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
} else {
	data.{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
}
{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
if value := res.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"); value.Exists() && !data.{{toGoName .TfName}}.IsNull() {
	data.{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
} else {
	data.{{toGoName .TfName}} = types.MapNull(types.StringType)
}
{{- else if isNestedMap .}}
{{- $list := (toGoName .TfName)}}
if value := res{{if .ModelName}}.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"){{end}}; value.Exists() {
	for mapKey, mapItem := range data.{{toGoName .TfName}} {
		value.ForEach(func(_, v gjson.Result) bool {
			computedKey := {{range $i, $kp := getKeyPartAttributes .Attributes}}{{if $i}} + "_" + {{end}}{{if eq $kp.Type "Int64"}}fmt.Sprintf("%v", v.Get("{{if $kp.DataPath}}{{$kp.DataPath}}.{{end}}{{$kp.ModelName}}").Int()){{else if eq $kp.Type "Bool"}}fmt.Sprintf("%v", v.Get("{{if $kp.DataPath}}{{$kp.DataPath}}.{{end}}{{$kp.ModelName}}").Bool()){{else}}v.Get("{{if $kp.DataPath}}{{$kp.DataPath}}.{{end}}{{$kp.ModelName}}").String(){{end}}{{end}}
			if computedKey != mapKey {
				return true
			}
			{{- range .Attributes}}
			{{- if and (not .Value) (not .WriteOnly) (not .Reference) .Computed}}
			{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
			if mapItem.{{toGoName .TfName}}.IsUnknown() {
				if fv := v.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); fv.Exists() {
					mapItem.{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "fv" .}}{{else}}fv.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
				} else {{if .DefaultValue}}if mapItem.{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
					mapItem.{{toGoName .TfName}} = types.{{.Type}}Null()
				}
			}
			{{- end}}
			{{- end}}
			{{- end}}
			data.{{$list}}[mapKey] = mapItem
			return false
		})
	}
}
{{- else if isNestedListSetMap .}}
{{- if hasUnknownsChildren .Attributes}}
{{- $list := (toGoName .TfName)}}
for i := range data.{{toGoName .TfName}} {
	{{- if and (isNestedListSet .) (eq .MaxList 1)}}
	values := res{{if .ModelName}}.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}"){{end}}.Array()
	if len(values) != 1 {
		data.{{$list}} = nil
		break
	}
	r := values[0]
	{{- else}}
	keys := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if not .Computed}}{{if or .Id $noId}}{{if or (eq .Type "Int64") (eq .Type "Bool") (eq .Type "String")}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}}{{end}}{{end}} }
	keyValues := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if not .Computed}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(data.{{$list}}[i].{{toGoName .TfName}}.ValueInt64(), 10), {{else if eq .Type "Bool"}}strconv.FormatBool(data.{{$list}}[i].{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}data.{{$list}}[i].{{toGoName .TfName}}.Value{{.Type}}(), {{end}}{{end}}{{end}}{{end}} }
	{{- if hasComputedRefreshValueInKeys .Attributes}}
	emptyKeys := [...]string{ {{range .Attributes}}{{if .ComputedRefreshValue}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}} }
	{{- end}}

	var r gjson.Result
	res.{{if .ModelName}}Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else}}{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}{{end}}").{{end}}ForEach(
		func(_, v gjson.Result) bool {
			found := false
			for ik := range keys {
				{{- if hasComputedRefreshValueInKeys .Attributes}}
				if v.Get(keys[ik]).String() == keyValues[ik] || slices.Contains(emptyKeys[:], keys[ik]) && keyValues[ik] == "" {
				{{- else}}
				if v.Get(keys[ik]).String() == keyValues[ik] {
				{{- end}}
					found = true
					continue
				}
				found = false
				break
			}
			if found {
				r = v
				return false
			}
			return true
		},
	)
	{{- end}}

	{{- range .Attributes}}
	{{- if and (not .Value) (not .WriteOnly) (not .Reference) .Computed}}
	{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
	if data.{{$list}}[i].{{toGoName .TfName}}.IsUnknown() {
		if value := r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{toGoName .TfName}}.IsNull() {
			data.{{$list}}[i].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
		} else {{if .DefaultValue}}if data.{{$list}}[i].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
			data.{{$list}}[i].{{toGoName .TfName}} = types.{{.Type}}Null()
			{{- if eq .ModelName "id"}}
			continue
			{{- end}}
		}
	}
	{{- else if isListSet .}}
	if value := r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{toGoName .TfName}}.IsNull() {
		data.{{$list}}[i].{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
	} else {
		data.{{$list}}[i].{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
	}
	{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
	if value := r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{toGoName .TfName}}.IsNull() {
		data.{{$list}}[i].{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
	} else {
		data.{{$list}}[i].{{toGoName .TfName}} = types.MapNull(types.StringType)
	}
	{{- else if isNestedListSetMap .}}
	{{- $clist := (toGoName .TfName)}}
	for ci := range data.{{$list}}[i].{{toGoName .TfName}} {
		{{- if and (isNestedListSet .) (eq .MaxList 1)}}
		values := r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").Array()
		if len(values) != 1 {
			data.{{$list}}[i].{{$clist}} = nil
			break
		}
		cr := values[0]
		{{- else}}
		keys := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if or (eq .Type "Int64") (eq .Type "Bool") (eq .Type "String")}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}}{{end}} }
		keyValues := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.ValueInt64(), 10), {{else if eq .Type "Bool"}}strconv.FormatBool(data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.Value{{.Type}}(), {{end}}{{end}}{{end}} }
		emptyKeys := [...]string{ {{range .Attributes}}{{if .ComputedRefreshValue}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}",{{end}}{{end}} }

		var cr gjson.Result
		r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").ForEach(
			func(_, v gjson.Result) bool {
				found := false
				for ik := range keys {
					if v.Get(keys[ik]).String() == keyValues[ik] || slices.Contains(emptyKeys[:], keys[ik]) &&  keyValues[ik] == "" {
						found = true
						continue
					}
					found = false
					break
				}
				if found {
					cr = v
					return false
				}
				return true
			},
		)
		{{- end}}

		{{- range .Attributes}}
		{{- if and (not .Value) (not .WriteOnly) (not .Reference)}}
		{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
		if value := cr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.IsNull() {
			data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
		} else {{if .DefaultValue}}if data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
			data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} = types.{{.Type}}Null()
		}
		{{- else if isListSet .}}
		if value := cr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.IsNull() {
			data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
		} else {
			data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
		}
		{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
		if value := cr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}}.IsNull() {
			data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
		} else {
			data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} = types.MapNull(types.StringType)
		}
		{{- else if isNestedListSetMap .}}
		{{- $cclist := (toGoName .TfName)}}
		for cci := range data.{{$list}}[i].{{$clist}}[ci].{{toGoName .TfName}} {
			{{- if and (isNestedListSet .) (eq .MaxList 1)}}
			values := cr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").Array()
			if len(values) != 1 {
				data.{{$list}}[i].{{$clist}}[ci].{{$cclist}} = nil
				break
			}
			ccr := values[0]
			{{- else}}
			keys := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if or (eq .Type "Int64") (eq .Type "Bool") (eq .Type "String")}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}}{{end}} }
			keyValues := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.ValueInt64(), 10), {{else if eq .Type "Bool"}}strconv.FormatBool(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.Value{{.Type}}(), {{end}}{{end}}{{end}} }
			{{- if hasComputedRefreshValueInKeys .Attributes}}
			emptyKeys := [...]string{ {{range .Attributes}}{{if .ComputedRefreshValue}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}} }
			{{- end}}

			var ccr gjson.Result
			cr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").ForEach(
				func(_, v gjson.Result) bool {
					found := false
					for ik := range keys {
						{{- if hasComputedRefreshValueInKeys .Attributes}}
						if v.Get(keys[ik]).String() == keyValues[ik] || slices.Contains(emptyKeys[:], keys[ik]) && keyValues[ik] == "" {
						{{- else}}
						if v.Get(keys[ik]).String() == keyValues[ik] {
						{{- end}}
							found = true
							continue
						}
						found = false
						break
					}
					if found {
						ccr = v
						return false
					}
					return true
				},
			)
			{{- end}}

			{{- range .Attributes}}
			{{- if and (not .Value) (not .WriteOnly) (not .Reference)}}
			{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
			if value := ccr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.IsNull() {
				data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
			} else {{if .DefaultValue}}if data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
				data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} = types.{{.Type}}Null()
			}
			{{- else if isListSet .}}
			if value := ccr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.IsNull() {
				data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
			} else {
				data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
			}
			{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
			if value := ccr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}}.IsNull() {
				data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
			} else {
				data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} = types.MapNull(types.StringType)
			}
			{{- else if isNestedListSetMap .}}
			{{- $c3list := (toGoName .TfName)}}
			for c3i := range data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{toGoName .TfName}} {
				{{- if and (isNestedListSet .) (eq .MaxList 1)}}
				values := ccr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").Array()
				if len(values) != 1 {
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}} = nil
					break
				}
				c3r := values[0]
				{{- else}}
				keys := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if or (eq .Type "Int64") (eq .Type "Bool") (eq .Type "String")}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}}{{end}} }
				keyValues := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.ValueInt64(), 10), {{else if eq .Type "Bool"}}strconv.FormatBool(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.Value{{.Type}}(), {{end}}{{end}}{{end}} }
				{{- if hasComputedRefreshValueInKeys .Attributes}}
				emptyKeys := [...]string{ {{range .Attributes}}{{if .ComputedRefreshValue}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}} }
				{{- end}}

				var c3r gjson.Result
				ccr.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").ForEach(
					func(_, v gjson.Result) bool {
						found := false
						for ik := range keys {
							{{- if hasComputedRefreshValueInKeys .Attributes}}
							if v.Get(keys[ik]).String() == keyValues[ik] || slices.Contains(emptyKeys[:], keys[ik]) && keyValues[ik] == "" {
							{{- else}}
							if v.Get(keys[ik]).String() == keyValues[ik] {
							{{- end}}
								found = true
								continue
							}
							found = false
							break
						}
						if found {
							c3r = v
							return false
						}
						return true
					},
				)
				{{- end}}

				{{- range .Attributes}}
				{{- if and (not .Value) (not .WriteOnly) (not .Reference)}}
				{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
				if value := c3r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.IsNull() {
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
				} else {{if .DefaultValue}}if data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} = types.{{.Type}}Null()
				}
				{{- else if isListSet .}}
				if value := c3r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.IsNull() {
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
				} else {
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
				}
				{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
				if value := c3r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}}.IsNull() {
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
				} else {
					data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} = types.MapNull(types.StringType)
				}
				{{- else if isNestedListSetMap .}}
				{{- $c4list := (toGoName .TfName)}}
				for c4i := range data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{toGoName .TfName}} {
					{{- if and (isNestedListSet .) (eq .MaxList 1)}}
					values := c3r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").Array()
					if len(values) != 1 {
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}} = nil
						break
					}
					c4r := values[0]
					{{- else}}
					keys := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if or (eq .Type "Int64") (eq .Type "Bool") (eq .Type "String")}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}}{{end}} }
					keyValues := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.ValueInt64(), 10), {{else if eq .Type "Bool"}}strconv.FormatBool(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.Value{{.Type}}(), {{end}}{{end}}{{end}} }
					{{- if hasComputedRefreshValueInKeys .Attributes}}
					emptyKeys := [...]string{ {{range .Attributes}}{{if .ComputedRefreshValue}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}} }
					{{- end}}

					var c4r gjson.Result
					c3r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").ForEach(
						func(_, v gjson.Result) bool {
							found := false
							for ik := range keys {
								{{- if hasComputedRefreshValueInKeys .Attributes}}
								if v.Get(keys[ik]).String() == keyValues[ik] || slices.Contains(emptyKeys[:], keys[ik]) && keyValues[ik] == "" {
								{{- else}}
								if v.Get(keys[ik]).String() == keyValues[ik] {
								{{- end}}
									found = true
									continue
								}
								found = false
								break
							}
							if found {
								c4r = v
								return false
							}
							return true
						},
					)
					{{- end}}

					{{- range .Attributes}}
					{{- if and (not .Value) (not .WriteOnly) (not .Reference)}}
					{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
					if value := c4r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.IsNull() {
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
					} else {{if .DefaultValue}}if data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} = types.{{.Type}}Null()
					}
					{{- else if isListSet .}}
					if value := c4r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.IsNull() {
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
					} else {
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
					}
					{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
					if value := c4r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}}.IsNull() {
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
					} else {
						data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} = types.MapNull(types.StringType)
					}
					{{- else if isNestedListSetMap .}}
					{{- $c5list := (toGoName .TfName)}}
					for c5i := range data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{toGoName .TfName}} {
						{{- if and (isNestedListSet .) (eq .MaxList 1)}}
						values := c4r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").Array()
						if len(values) != 1 {
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}} = nil
							break
						}
						c5r := values[0]
						{{- else}}
						keys := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if or (eq .Type "Int64") (eq .Type "Bool") (eq .Type "String")}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}}{{end}} }
						keyValues := [...]string{ {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.ValueInt64(), 10), {{else if eq .Type "Bool"}}strconv.FormatBool(data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.Value{{.Type}}(), {{end}}{{end}}{{end}} }
						{{- if hasComputedRefreshValueInKeys .Attributes}}
						emptyKeys := [...]string{ {{range .Attributes}}{{if .ComputedRefreshValue}}"{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}", {{end}}{{end}} }
						{{- end}}

						var c5r gjson.Result
						c4r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}").ForEach(
							func(_, v gjson.Result) bool {
								found := false
								for ik := range keys {
									{{- if hasComputedRefreshValueInKeys .Attributes}}
									if v.Get(keys[ik]).String() == keyValues[ik] || slices.Contains(emptyKeys[:], keys[ik]) && keyValues[ik] == "" {
									{{- else}}
									if v.Get(keys[ik]).String() == keyValues[ik] {
									{{- end}}
										found = true
										continue
									}
									found = false
									break
								}
								if found {
									c5r = v
									return false
								}
								return true
							},
						)
						{{- end}}

						{{- range .Attributes}}
						{{- if and (not .Value) (not .WriteOnly) (not .Reference)}}
						{{- if or (eq .Type "String") (eq .Type "Int64") (eq .Type "Float64") (eq .Type "Bool")}}
						if value := c5r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.IsNull() {
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}} = types.{{.Type}}Value({{if eq .Type "Float64"}}{{floatReadExpr "value" .}}{{else if isYangFlag .}}true{{else}}value.{{if eq .Type "Int64"}}Int{{else}}{{.Type}}{{end}}(){{end}})
						} else {{if .DefaultValue}}if data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.Value{{.Type}}() != {{if eq .Type "String"}}"{{end}}{{.DefaultValue}}{{if eq .Type "String"}}"{{end}} {{end}}{
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}} = types.{{.Type}}Null()
						}
						{{- else if isListSet .}}
						if value := c5r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.IsNull() {
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}} = helpers.Get{{.ElementType}}{{.Type}}(value.Array())
						} else {
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}} = types.{{.Type}}Null(types.{{.ElementType}}Type)
						}
						{{- else if and (eq .Type "Map") (not (isNestedMap .))}}
						if value := c5r.Get("{{if .DataPath}}{{.DataPath}}.{{end}}{{.ModelName}}"); value.Exists() && !data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}}.IsNull() {
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}} = helpers.GetStringMap(value.Map())
						} else {
							data.{{$list}}[i].{{$clist}}[ci].{{$cclist}}[cci].{{$c3list}}[c3i].{{$c4list}}[c4i].{{$c5list}}[c5i].{{toGoName .TfName}} = types.MapNull(types.StringType)
						}
						{{- end}}
						{{- end}}
						{{- end}}
					}
					{{- end}}
					{{- end}}
					{{- end}}
				}
				{{- end}}
				{{- end}}
				{{- end}}
			}
			{{- end}}
			{{- end}}
			{{- end}}
		}

		{{- end}}
		{{- end}}
		{{- end}}
	}
	{{- end}}
	{{- end}}
	{{- end}}
}
{{- end}}
{{- end}}
{{- end}}
{{- end}}
}
// End of section. //template:end fromBodyUnknowns
{{- end}}

// Section below is generated&owned by "gen/generator.go". //template:begin isNull
func (data *{{camelCase .Name}}) isNull(ctx context.Context, res gjson.Result) bool {
	{{- range .Attributes}}
	{{- if and (not .Value) (not .Id) (not .Reference) (not .QueryParam) (not .GetQueryParam) (not .QueryParamNoBody) }}
	{{- if isNestedListSetMap .}}
	if len(data.{{toGoName .TfName}}) > 0 {
		return false
	}
	{{- else}}
	if !data.{{toGoName .TfName}}.IsNull() {
		return false
	}
	{{- end}}
	{{- end}}
	{{- end}}
	return true
}
// End of section. //template:end isNull

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
	"fmt"
	"net/url"
	{{- if .DeviceUnreachabilityWarning}}
	"regexp"
	{{- end}}
	"strings"
	"sync"
	{{- if .RetryOnErrorCodes}}
	"time"
	{{- end}}

	"github.com/CiscoDevNet/terraform-provider-crosswork/internal/provider/helpers"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	cw "github.com/netascode/go-crosswork"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &{{camelCase .Name}}Resource{}
{{- if not .NoImport}}
var _ resource.ResourceWithImportState = &{{camelCase .Name}}Resource{}
{{- end}}

func New{{camelCase .Name}}Resource() resource.Resource {
	return &{{camelCase .Name}}Resource{}
}

type {{camelCase .Name}}Resource struct {
	clients               map[string]*cw.Client
	AllowExistingOnCreate bool
	cache                 *ThreadSafeCache
	{{- if hasNsoDevices .}}
	deviceReady           nsoDeviceConfig
	{{- end}}
}

func (r *{{camelCase .Name}}Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_{{snakeCase .Name}}"
}

func (r *{{camelCase .Name}}Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("{{.ResDescription}}").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
				{{- if not .PutUpdateId}}
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				{{- end}}
			},
			"instance": schema.StringAttribute{
				MarkdownDescription: "Select an instance of the provider configuration by its name. Leave it unset to use the default instance.",
				Optional:            true,
			},
			{{- $root := . }}
			{{- range  .Attributes}}
			{{- if not .Value}}
            "{{.TfName}}": schema.{{if isNestedListSetMap .}}{{.Type}}Nested{{else if isList .}}List{{else if isSet .}}Set{{else if eq .Type "Versions"}}List{{else if eq .Type "Version"}}Int64{{else}}{{.Type}}{{end}}Attribute{
				MarkdownDescription: helpers.NewAttributeDescription("{{.Description}}")
					{{- if len .EnumValues -}}
					.AddStringEnumDescription({{range .EnumValues}}"{{.}}", {{end}})
					{{- end -}}
					{{- if or (ne .MinInt 0) (ne .MaxInt 0) -}}
					.AddIntegerRangeDescription({{.MinInt}}, {{.MaxInt}})
					{{- end -}}
					{{- if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0) -}}
					.AddFloatRangeDescription({{.MinFloat}}, {{.MaxFloat}})
					{{- end -}}
					{{- if .DefaultValue -}}
					.AddDefaultValueDescription("{{.DefaultValue}}")
					{{- end -}}
					{{- if len .MutualExclusivityNote -}}
					.AddMutualExclusivityDescription("{{.MutualExclusivityNote}}")
					{{- end -}}
					{{- if len .CoexistenceNote -}}
					.AddCoexistenceNote("{{.CoexistenceNote}}")
					{{- end -}}
					.String,
				{{- if isListSet .}}
				ElementType:         types.{{.ElementType}}Type,
                {{- else if and (eq .Type "Map") (not (isNestedMap .))}}
				{{- if .ElementType}}
				ElementType: types.{{.ElementType}}Type,
				{{- else if $.NoRead}}
				ElementType:         types.ListType{ElemType: types.StringType},
				{{- else}}
				ElementType:         types.StringType,
				{{- end}}
				{{- end}}
				{{- if or .WriteOnlyTF .WriteOnlyTFNative}}
				Optional:            true,
				WriteOnly:           true,
				Sensitive:           true,
				{{- else}}
				{{- if .CoexistingSecret}}
				Sensitive:           true,
				{{- end}}
				{{- if or .Id .MatchId (and .Reference (not .Computed)) .Mandatory}}
				Required:            true,
				{{- else}}
				Optional:            true,
				{{- end}}
				{{- if or (len .DefaultValue) .Computed}}
				Computed:            true,
				{{- end}}
				{{- end}}
				{{- if len .EnumValues}}
				{{- if eq .Type "Int64" }}
		                Validators: []validator.Int64{
		                    int64validator.OneOf({{range .EnumValues}}{{.}}, {{end}}),
		                },
				{{- else }}
				Validators: []validator.String{
					stringvalidator.OneOf({{range .EnumValues}}"{{.}}", {{end}}),
				},
				{{- end }}
				{{- else if or (len .StringPatterns) (ne .StringMinLength 0) (ne .StringMaxLength 0) }}
				Validators: []validator.String{
					{{- if or (ne .StringMinLength 0) (ne .StringMaxLength 0)}}
					stringvalidator.LengthBetween({{.StringMinLength}}, {{.StringMaxLength}}),
					{{- end}}
					{{- range .StringPatterns}}
					stringvalidator.RegexMatches(regexp.MustCompile(`{{.}}`), ""),
					{{- end}}
				},
				{{- else if or (ne .MinInt 0) (ne .MaxInt 0)}}
				Validators: []validator.Int64{
					int64validator.Between({{.MinInt}}, {{.MaxInt}}),
				},
				{{- else if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0)}}
				Validators: []validator.Float64{
					float64validator.Between({{.MinFloat}}, {{.MaxFloat}}),
				},
				{{- end}}
				{{- if and (len .DefaultValue) (eq .Type "Int64")}}
				Default:             int64default.StaticInt64({{.DefaultValue}}),
				{{- else if and (len .DefaultValue) (eq .Type "Bool")}}
				Default:             booldefault.StaticBool({{.DefaultValue}}),
				{{- else if and (len .DefaultValue) (eq .Type "String")}}
				Default:             stringdefault.StaticString("{{.DefaultValue}}"),
				{{- end}}
				{{- if .CustomModifier}}
				PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
					{{.CustomModifier}}{},
				},
				{{- else}}
				{{- if or .Id .MatchId (and .Reference (not .Computed)) .RequiresReplace}}
				PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
					{{if eq .Type "StringList"}}list{{else}}{{snakeCase .Type}}{{end}}planmodifier.RequiresReplace(),
				},
				{{- else if and .Computed (not .ComputedRefreshValue) (not (hasComputedRefreshValue .Attributes)) (not .NoUseStateForUnknown)}}
				PlanModifiers: []planmodifier.{{.Type}}{
					{{snakeCase .Type}}planmodifier.UseStateForUnknown(),
				},
				{{- end}}
				{{- end}}
                {{- if isNestedListSetMap .}}
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						{{- range  .Attributes}}
						{{- if not .Value}}
                        "{{.TfName}}": schema.{{if isNestedListSetMap .}}{{.Type}}Nested{{else if isList .}}List{{else if isSet .}}Set{{else if eq .Type "Versions"}}List{{else if eq .Type "Version"}}Int64{{else}}{{.Type}}{{end}}Attribute{
							MarkdownDescription: helpers.NewAttributeDescription("{{.Description}}")
								{{- if len .EnumValues -}}
								.AddStringEnumDescription({{range .EnumValues}}"{{.}}", {{end}})
								{{- end -}}
								{{- if or (ne .MinInt 0) (ne .MaxInt 0) -}}
								.AddIntegerRangeDescription({{.MinInt}}, {{.MaxInt}})
								{{- end -}}
								{{- if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0) -}}
								.AddFloatRangeDescription({{.MinFloat}}, {{.MaxFloat}})
								{{- end -}}
								{{- if .DefaultValue -}}
								.AddDefaultValueDescription("{{.DefaultValue}}")
								{{- end -}}
								{{- if len .MutualExclusivityNote -}}
								.AddMutualExclusivityDescription("{{.MutualExclusivityNote}}")
								{{- end -}}
								{{- if len .CoexistenceNote -}}
								.AddCoexistenceNote("{{.CoexistenceNote}}")
								{{- end -}}
								.String,
							{{- if isListSet .}}
							ElementType:         types.{{.ElementType}}Type,
                            {{- else if and (eq .Type "Map") (not (isNestedMap .))}}
							{{- if .ElementType}}
							ElementType: types.{{.ElementType}}Type,
							{{- else if $.NoRead}}
							ElementType:         types.ListType{ElemType: types.StringType},
							{{- else}}
							ElementType:         types.StringType,
							{{- end}}
							{{- end}}
							{{- if or .WriteOnlyTF .WriteOnlyTFNative}}
							Optional:            true,
							WriteOnly:           true,
							Sensitive:           true,
							{{- else}}
							{{- if .CoexistingSecret}}
							Sensitive:           true,
							{{- end}}
							{{- if or (and .Id (not .ComputedRefreshValue)) .Reference .Mandatory }}
							Required:            true,
							{{- else if and (not .Computed) (not .ComputedRefreshValue) }}
							Optional:            true,
							{{- else if and .Computed (not .ComputedRefreshValue) }}
							Optional:            true,
							{{- else if and .ComputedRefreshValue (not .Computed) }}
							Optional:            true,
							{{- end}}
							{{- if or (len .DefaultValue) .Computed}}
							Computed:            true,
							{{- end}}
							{{- end}}
							{{- if len .EnumValues}}
							{{- if eq .Type "Int64" }}
							Validators: []validator.Int64{
								int64validator.OneOf({{range .EnumValues}}{{.}}, {{end}}),
							},
							{{- else }}
							Validators: []validator.String{
								stringvalidator.OneOf({{range .EnumValues}}"{{.}}", {{end}}),
							},
							{{- end }}
							{{- else if or (len .StringPatterns) (ne .StringMinLength 0) (ne .StringMaxLength 0) }}
							Validators: []validator.String{
								{{- if or (ne .StringMinLength 0) (ne .StringMaxLength 0)}}
								stringvalidator.LengthBetween({{.StringMinLength}}, {{.StringMaxLength}}),
								{{- end}}
								{{- range .StringPatterns}}
								stringvalidator.RegexMatches(regexp.MustCompile(`{{.}}`), ""),
								{{- end}}
							},
							{{- else if or (ne .MinInt 0) (ne .MaxInt 0)}}
							Validators: []validator.Int64{
								int64validator.Between({{.MinInt}}, {{.MaxInt}}),
							},
							{{- else if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0)}}
							Validators: []validator.Float64{
								float64validator.Between({{.MinFloat}}, {{.MaxFloat}}),
							},
							{{- end}}
							{{- if and (len .DefaultValue) (eq .Type "Int64")}}
							Default:             int64default.StaticInt64({{.DefaultValue}}),
							{{- else if and (len .DefaultValue) (eq .Type "Bool")}}
							Default:             booldefault.StaticBool({{.DefaultValue}}),
							{{- else if and (len .DefaultValue) (eq .Type "String")}}
							Default:             stringdefault.StaticString("{{.DefaultValue}}"),
							{{- end}}
							{{- if .CustomModifier}}
							PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
								{{.CustomModifier}}{},
							},
							{{- else}}
							{{- if and .RequiresReplace (not $root.RootList)}}
							PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
								{{if eq .Type "StringList"}}list{{else}}{{snakeCase .Type}}{{end}}planmodifier.RequiresReplace(),
							},
							{{- end}}
							{{- if and .Computed (not .ComputedRefreshValue) (not .NoUseStateForUnknown)}}
							PlanModifiers: []planmodifier.{{.Type}}{
								{{snakeCase .Type}}planmodifier.UseStateForUnknown(),
							},
							{{- end}}
							{{- end}}
                            {{- if isNestedListSetMap .}}
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									{{- range  .Attributes}}
									{{- if not .Value}}
                                    "{{.TfName}}": schema.{{if isNestedListSetMap .}}{{.Type}}Nested{{else if isList .}}List{{else if isSet .}}Set{{else if eq .Type "Versions"}}List{{else if eq .Type "Version"}}Int64{{else}}{{.Type}}{{end}}Attribute{
										MarkdownDescription: helpers.NewAttributeDescription("{{.Description}}")
											{{- if len .EnumValues -}}
											.AddStringEnumDescription({{range .EnumValues}}"{{.}}", {{end}})
											{{- end -}}
											{{- if or (ne .MinInt 0) (ne .MaxInt 0) -}}
											.AddIntegerRangeDescription({{.MinInt}}, {{.MaxInt}})
											{{- end -}}
											{{- if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0) -}}
											.AddFloatRangeDescription({{.MinFloat}}, {{.MaxFloat}})
											{{- end -}}
											{{- if .DefaultValue -}}
											.AddDefaultValueDescription("{{.DefaultValue}}")
											{{- end -}}
											{{- if len .MutualExclusivityNote -}}
											.AddMutualExclusivityDescription("{{.MutualExclusivityNote}}")
											{{- end -}}
											{{- if len .CoexistenceNote -}}
											.AddCoexistenceNote("{{.CoexistenceNote}}")
											{{- end -}}
											.String,
										{{- if isListSet .}}
										ElementType:         types.{{.ElementType}}Type,
                                        {{- else if and (eq .Type "Map") (not (isNestedMap .))}}
										{{- if .ElementType}}
										ElementType: types.{{.ElementType}}Type,
										{{- else if $.NoRead}}
										ElementType:         types.ListType{ElemType: types.StringType},
										{{- else}}
										ElementType:         types.StringType,
										{{- end}}
										{{- end}}
										{{- if or .WriteOnlyTF .WriteOnlyTFNative}}
										Optional:            true,
										WriteOnly:           true,
										Sensitive:           true,
										{{- else}}
										{{- if .CoexistingSecret}}
										Sensitive:           true,
										{{- end}}
										{{- if or .Id .Reference .Mandatory}}
										Required:            true,
										{{- else}}
										Optional:            true,
										{{- end}}
										{{- if or (len .DefaultValue) .Computed}}
										Computed:            true,
										{{- end}}
										{{- end}}
										{{- if len .EnumValues}}
										{{- if eq .Type "Int64" }}
										Validators: []validator.Int64{
											int64validator.OneOf({{range .EnumValues}}{{.}}, {{end}}),
										},
										{{- else }}
										Validators: []validator.String{
											stringvalidator.OneOf({{range .EnumValues}}"{{.}}", {{end}}),
										},
										{{- end }}
										{{- else if or (len .StringPatterns) (ne .StringMinLength 0) (ne .StringMaxLength 0) }}
										Validators: []validator.String{
											{{- if or (ne .StringMinLength 0) (ne .StringMaxLength 0)}}
											stringvalidator.LengthBetween({{.StringMinLength}}, {{.StringMaxLength}}),
											{{- end}}
											{{- range .StringPatterns}}
											stringvalidator.RegexMatches(regexp.MustCompile(`{{.}}`), ""),
											{{- end}}
										},
										{{- else if or (ne .MinInt 0) (ne .MaxInt 0)}}
										Validators: []validator.Int64{
											int64validator.Between({{.MinInt}}, {{.MaxInt}}),
										},
										{{- else if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0)}}
										Validators: []validator.Float64{
											float64validator.Between({{.MinFloat}}, {{.MaxFloat}}),
										},
										{{- end}}
										{{- if and (len .DefaultValue) (eq .Type "Int64")}}
										Default:             int64default.StaticInt64({{.DefaultValue}}),
										{{- else if and (len .DefaultValue) (eq .Type "Bool")}}
										Default:             booldefault.StaticBool({{.DefaultValue}}),
										{{- else if and (len .DefaultValue) (eq .Type "String")}}
										Default:             stringdefault.StaticString("{{.DefaultValue}}"),
										{{- end}}
										{{- if .CustomModifier}}
										PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
											{{.CustomModifier}}{},
										},
										{{- else}}
										{{- if .RequiresReplace}}
										PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
											{{if eq .Type "StringList"}}list{{else}}{{snakeCase .Type}}{{end}}planmodifier.RequiresReplace(),
										},
										{{- end}}
										{{- if and .Computed (not .ComputedRefreshValue) (not .NoUseStateForUnknown)}}
										PlanModifiers: []planmodifier.{{.Type}}{
											{{snakeCase .Type}}planmodifier.UseStateForUnknown(),
										},
										{{- end}}
										{{- end}}
                                        {{- if isNestedListSetMap .}}
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												{{- range  .Attributes}}
												{{- if not .Value}}
                                                "{{.TfName}}": schema.{{if isNestedListSetMap .}}{{.Type}}Nested{{else if isList .}}List{{else if isSet .}}Set{{else if eq .Type "Versions"}}List{{else if eq .Type "Version"}}Int64{{else}}{{.Type}}{{end}}Attribute{
													MarkdownDescription: helpers.NewAttributeDescription("{{.Description}}")
														{{- if len .EnumValues -}}
														.AddStringEnumDescription({{range .EnumValues}}"{{.}}", {{end}})
														{{- end -}}
														{{- if or (ne .MinInt 0) (ne .MaxInt 0) -}}
														.AddIntegerRangeDescription({{.MinInt}}, {{.MaxInt}})
														{{- end -}}
														{{- if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0) -}}
														.AddFloatRangeDescription({{.MinFloat}}, {{.MaxFloat}})
														{{- end -}}
														{{- if .DefaultValue -}}
														.AddDefaultValueDescription("{{.DefaultValue}}")
														{{- end -}}
														{{- if len .MutualExclusivityNote -}}
														.AddMutualExclusivityDescription("{{.MutualExclusivityNote}}")
														{{- end -}}
														{{- if len .CoexistenceNote -}}
														.AddCoexistenceNote("{{.CoexistenceNote}}")
														{{- end -}}
														.String,
													{{- if isListSet .}}
													ElementType:         types.{{.ElementType}}Type,
                                                    {{- else if and (eq .Type "Map") (not (isNestedMap .))}}
													{{- if .ElementType}}
													ElementType: types.{{.ElementType}}Type,
													{{- else if $.NoRead}}
													ElementType:         types.ListType{ElemType: types.StringType},
													{{- else}}
													ElementType:         types.StringType,
													{{- end}}
													{{- end}}
													{{- if or .WriteOnlyTF .WriteOnlyTFNative}}
													Optional:            true,
													WriteOnly:           true,
													Sensitive:           true,
													{{- else}}
													{{- if .CoexistingSecret}}
													Sensitive:           true,
													{{- end}}
													{{- if or .Id .Reference .Mandatory}}
													Required:            true,
													{{- else}}
													Optional:            true,
													{{- end}}
													{{- if or (len .DefaultValue) .Computed}}
													Computed:            true,
													{{- end}}
													{{- end}}
													{{- if len .EnumValues}}
													{{- if eq .Type "Int64" }}
													Validators: []validator.Int64{
														int64validator.OneOf({{range .EnumValues}}{{.}}, {{end}}),
													},
													{{- else }}
													Validators: []validator.String{
														stringvalidator.OneOf({{range .EnumValues}}"{{.}}", {{end}}),
													},
													{{- end }}
													{{- else if or (len .StringPatterns) (ne .StringMinLength 0) (ne .StringMaxLength 0) }}
													Validators: []validator.String{
														{{- if or (ne .StringMinLength 0) (ne .StringMaxLength 0)}}
														stringvalidator.LengthBetween({{.StringMinLength}}, {{.StringMaxLength}}),
														{{- end}}
														{{- range .StringPatterns}}
														stringvalidator.RegexMatches(regexp.MustCompile(`{{.}}`), ""),
														{{- end}}
													},
													{{- else if or (ne .MinInt 0) (ne .MaxInt 0)}}
													Validators: []validator.Int64{
														int64validator.Between({{.MinInt}}, {{.MaxInt}}),
													},
													{{- else if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0)}}
													Validators: []validator.Float64{
														float64validator.Between({{.MinFloat}}, {{.MaxFloat}}),
													},
													{{- end}}
													{{- if and (len .DefaultValue) (eq .Type "Int64")}}
													Default:             int64default.StaticInt64({{.DefaultValue}}),
													{{- else if and (len .DefaultValue) (eq .Type "Bool")}}
													Default:             booldefault.StaticBool({{.DefaultValue}}),
													{{- else if and (len .DefaultValue) (eq .Type "String")}}
													Default:             stringdefault.StaticString("{{.DefaultValue}}"),
													{{- end}}
													{{- if .CustomModifier}}
													PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
														{{.CustomModifier}}{},
													},
													{{- else}}
													{{- if .RequiresReplace}}
													PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
														{{if eq .Type "StringList"}}list{{else}}{{snakeCase .Type}}{{end}}planmodifier.RequiresReplace(),
													},
													{{- end}}
													{{- if and .Computed (not .ComputedRefreshValue) (not .NoUseStateForUnknown)}}
													PlanModifiers: []planmodifier.{{.Type}}{
														{{snakeCase .Type}}planmodifier.UseStateForUnknown(),
													},
													{{- end}}
													{{- end}}
	                                        {{- if isNestedListSetMap .}}
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													{{- range  .Attributes}}
													{{- if not .Value}}
	                                                "{{.TfName}}": schema.{{if isNestedListSetMap .}}{{.Type}}Nested{{else if isList .}}List{{else if isSet .}}Set{{else if eq .Type "Versions"}}List{{else if eq .Type "Version"}}Int64{{else}}{{.Type}}{{end}}Attribute{
														MarkdownDescription: helpers.NewAttributeDescription("{{.Description}}")
															{{- if len .EnumValues -}}
															.AddStringEnumDescription({{range .EnumValues}}"{{.}}", {{end}})
															{{- end -}}
															{{- if or (ne .MinInt 0) (ne .MaxInt 0) -}}
															.AddIntegerRangeDescription({{.MinInt}}, {{.MaxInt}})
															{{- end -}}
															{{- if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0) -}}
															.AddFloatRangeDescription({{.MinFloat}}, {{.MaxFloat}})
															{{- end -}}
															{{- if .DefaultValue -}}
															.AddDefaultValueDescription("{{.DefaultValue}}")
															{{- end -}}
															{{- if len .MutualExclusivityNote -}}
															.AddMutualExclusivityDescription("{{.MutualExclusivityNote}}")
															{{- end -}}
															{{- if len .CoexistenceNote -}}
															.AddCoexistenceNote("{{.CoexistenceNote}}")
															{{- end -}}
															.String,
														{{- if isListSet .}}
														ElementType:         types.{{.ElementType}}Type,
	                                                    {{- else if and (eq .Type "Map") (not (isNestedMap .))}}
														{{- if .ElementType}}
														ElementType: types.{{.ElementType}}Type,
														{{- else if $.NoRead}}
														ElementType:         types.ListType{ElemType: types.StringType},
														{{- else}}
														ElementType:         types.StringType,
														{{- end}}
														{{- end}}
														{{- if or .WriteOnlyTF .WriteOnlyTFNative}}
														Optional:            true,
														WriteOnly:           true,
														Sensitive:           true,
														{{- else}}
														{{- if .CoexistingSecret}}
														Sensitive:           true,
														{{- end}}
														{{- if or .Id .Reference .Mandatory}}
														Required:            true,
														{{- else}}
														Optional:            true,
														{{- end}}
														{{- if or (len .DefaultValue) .Computed}}
														Computed:            true,
														{{- end}}
														{{- end}}
														{{- if len .EnumValues}}
														{{- if eq .Type "Int64" }}
														Validators: []validator.Int64{
															int64validator.OneOf({{range .EnumValues}}{{.}}, {{end}}),
														},
														{{- else }}
														Validators: []validator.String{
															stringvalidator.OneOf({{range .EnumValues}}"{{.}}", {{end}}),
														},
														{{- end }}
														{{- else if or (len .StringPatterns) (ne .StringMinLength 0) (ne .StringMaxLength 0) }}
														Validators: []validator.String{
															{{- if or (ne .StringMinLength 0) (ne .StringMaxLength 0)}}
															stringvalidator.LengthBetween({{.StringMinLength}}, {{.StringMaxLength}}),
															{{- end}}
															{{- range .StringPatterns}}
															stringvalidator.RegexMatches(regexp.MustCompile(`{{.}}`), ""),
															{{- end}}
														},
														{{- else if or (ne .MinInt 0) (ne .MaxInt 0)}}
														Validators: []validator.Int64{
															int64validator.Between({{.MinInt}}, {{.MaxInt}}),
														},
														{{- else if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0)}}
														Validators: []validator.Float64{
															float64validator.Between({{.MinFloat}}, {{.MaxFloat}}),
														},
														{{- end}}
														{{- if and (len .DefaultValue) (eq .Type "Int64")}}
														Default:             int64default.StaticInt64({{.DefaultValue}}),
														{{- else if and (len .DefaultValue) (eq .Type "Bool")}}
														Default:             booldefault.StaticBool({{.DefaultValue}}),
														{{- else if and (len .DefaultValue) (eq .Type "String")}}
														Default:             stringdefault.StaticString("{{.DefaultValue}}"),
														{{- end}}
														{{- if .CustomModifier}}
														PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
															{{.CustomModifier}}{},
														},
														{{- else}}
														{{- if .RequiresReplace}}
														PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
															{{if eq .Type "StringList"}}list{{else}}{{snakeCase .Type}}{{end}}planmodifier.RequiresReplace(),
														},
														{{- end}}
														{{- if and .Computed (not .ComputedRefreshValue) (not .NoUseStateForUnknown)}}
														PlanModifiers: []planmodifier.{{.Type}}{
															{{snakeCase .Type}}planmodifier.UseStateForUnknown(),
														},
														{{- end}}
														{{- end}}
		                                        {{- if isNestedListSetMap .}}
												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														{{- range  .Attributes}}
														{{- if not .Value}}
		                                                "{{.TfName}}": schema.{{if isNestedListSetMap .}}{{.Type}}Nested{{else if isList .}}List{{else if isSet .}}Set{{else if eq .Type "Versions"}}List{{else if eq .Type "Version"}}Int64{{else}}{{.Type}}{{end}}Attribute{
															MarkdownDescription: helpers.NewAttributeDescription("{{.Description}}")
																{{- if len .EnumValues -}}
																.AddStringEnumDescription({{range .EnumValues}}"{{.}}", {{end}})
																{{- end -}}
																{{- if or (ne .MinInt 0) (ne .MaxInt 0) -}}
																.AddIntegerRangeDescription({{.MinInt}}, {{.MaxInt}})
																{{- end -}}
																{{- if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0) -}}
																.AddFloatRangeDescription({{.MinFloat}}, {{.MaxFloat}})
																{{- end -}}
																{{- if .DefaultValue -}}
																.AddDefaultValueDescription("{{.DefaultValue}}")
																{{- end -}}
																{{- if len .MutualExclusivityNote -}}
																.AddMutualExclusivityDescription("{{.MutualExclusivityNote}}")
																{{- end -}}
																{{- if len .CoexistenceNote -}}
																.AddCoexistenceNote("{{.CoexistenceNote}}")
																{{- end -}}
																.String,
															{{- if isListSet .}}
															ElementType:         types.{{.ElementType}}Type,
		                                                    {{- else if and (eq .Type "Map") (not (isNestedMap .))}}
															{{- if .ElementType}}
															ElementType: types.{{.ElementType}}Type,
															{{- else if $.NoRead}}
															ElementType:         types.ListType{ElemType: types.StringType},
															{{- else}}
															ElementType:         types.StringType,
															{{- end}}
															{{- end}}
															{{- if or .WriteOnlyTF .WriteOnlyTFNative}}
															Optional:            true,
															WriteOnly:           true,
															Sensitive:           true,
															{{- else}}
															{{- if .CoexistingSecret}}
															Sensitive:           true,
															{{- end}}
															{{- if or .Id .Reference .Mandatory}}
															Required:            true,
															{{- else}}
															Optional:            true,
															{{- end}}
															{{- if or (len .DefaultValue) .Computed}}
															Computed:            true,
															{{- end}}
															{{- end}}
															{{- if len .EnumValues}}
															{{- if eq .Type "Int64" }}
															Validators: []validator.Int64{
																int64validator.OneOf({{range .EnumValues}}{{.}}, {{end}}),
															},
															{{- else }}
															Validators: []validator.String{
																stringvalidator.OneOf({{range .EnumValues}}"{{.}}", {{end}}),
															},
															{{- end }}
															{{- else if or (len .StringPatterns) (ne .StringMinLength 0) (ne .StringMaxLength 0) }}
															Validators: []validator.String{
																{{- if or (ne .StringMinLength 0) (ne .StringMaxLength 0)}}
																stringvalidator.LengthBetween({{.StringMinLength}}, {{.StringMaxLength}}),
																{{- end}}
																{{- range .StringPatterns}}
																stringvalidator.RegexMatches(regexp.MustCompile(`{{.}}`), ""),
																{{- end}}
															},
															{{- else if or (ne .MinInt 0) (ne .MaxInt 0)}}
															Validators: []validator.Int64{
																int64validator.Between({{.MinInt}}, {{.MaxInt}}),
															},
															{{- else if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0)}}
															Validators: []validator.Float64{
																float64validator.Between({{.MinFloat}}, {{.MaxFloat}}),
															},
															{{- end}}
															{{- if and (len .DefaultValue) (eq .Type "Int64")}}
															Default:             int64default.StaticInt64({{.DefaultValue}}),
															{{- else if and (len .DefaultValue) (eq .Type "Bool")}}
															Default:             booldefault.StaticBool({{.DefaultValue}}),
															{{- else if and (len .DefaultValue) (eq .Type "String")}}
															Default:             stringdefault.StaticString("{{.DefaultValue}}"),
															{{- end}}
															{{- if .CustomModifier}}
															PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
																{{.CustomModifier}}{},
															},
															{{- else}}
															{{- if .RequiresReplace}}
															PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
																{{if eq .Type "StringList"}}list{{else}}{{snakeCase .Type}}{{end}}planmodifier.RequiresReplace(),
															},
															{{- end}}
															{{- if and .Computed (not .ComputedRefreshValue) (not .NoUseStateForUnknown)}}
															PlanModifiers: []planmodifier.{{.Type}}{
																{{snakeCase .Type}}planmodifier.UseStateForUnknown(),
															},
															{{- end}}
															{{- end}}
			                                        {{- if isNestedListSetMap .}}
													NestedObject: schema.NestedAttributeObject{
														Attributes: map[string]schema.Attribute{
															{{- range  .Attributes}}
															{{- if not .Value}}
			                                                "{{.TfName}}": schema.{{if isNestedListSetMap .}}{{.Type}}Nested{{else if isList .}}List{{else if isSet .}}Set{{else if eq .Type "Versions"}}List{{else if eq .Type "Version"}}Int64{{else}}{{.Type}}{{end}}Attribute{
																MarkdownDescription: helpers.NewAttributeDescription("{{.Description}}")
																	{{- if len .EnumValues -}}
																	.AddStringEnumDescription({{range .EnumValues}}"{{.}}", {{end}})
																	{{- end -}}
																	{{- if or (ne .MinInt 0) (ne .MaxInt 0) -}}
																	.AddIntegerRangeDescription({{.MinInt}}, {{.MaxInt}})
																	{{- end -}}
																	{{- if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0) -}}
																	.AddFloatRangeDescription({{.MinFloat}}, {{.MaxFloat}})
																	{{- end -}}
																	{{- if .DefaultValue -}}
																	.AddDefaultValueDescription("{{.DefaultValue}}")
																	{{- end -}}
																	{{- if len .MutualExclusivityNote -}}
																	.AddMutualExclusivityDescription("{{.MutualExclusivityNote}}")
																	{{- end -}}
																	{{- if len .CoexistenceNote -}}
																	.AddCoexistenceNote("{{.CoexistenceNote}}")
																	{{- end -}}
																	.String,
																{{- if isListSet .}}
																ElementType:         types.{{.ElementType}}Type,
			                                                    {{- else if and (eq .Type "Map") (not (isNestedMap .))}}
																{{- if .ElementType}}
																ElementType: types.{{.ElementType}}Type,
																{{- else if $.NoRead}}
																ElementType:         types.ListType{ElemType: types.StringType},
																{{- else}}
																ElementType:         types.StringType,
																{{- end}}
																{{- end}}
																{{- if or .WriteOnlyTF .WriteOnlyTFNative}}
																Optional:            true,
																WriteOnly:           true,
																Sensitive:           true,
																{{- else}}
																{{- if .CoexistingSecret}}
																Sensitive:           true,
																{{- end}}
																{{- if or .Id .Reference .Mandatory}}
																Required:            true,
																{{- else}}
																Optional:            true,
																{{- end}}
																{{- if or (len .DefaultValue) .Computed}}
																Computed:            true,
																{{- end}}
																{{- end}}
																{{- if len .EnumValues}}
																{{- if eq .Type "Int64" }}
																Validators: []validator.Int64{
																	int64validator.OneOf({{range .EnumValues}}{{.}}, {{end}}),
																},
																{{- else }}
																Validators: []validator.String{
																	stringvalidator.OneOf({{range .EnumValues}}"{{.}}", {{end}}),
																},
																{{- end }}
																{{- else if or (len .StringPatterns) (ne .StringMinLength 0) (ne .StringMaxLength 0) }}
																Validators: []validator.String{
																	{{- if or (ne .StringMinLength 0) (ne .StringMaxLength 0)}}
																	stringvalidator.LengthBetween({{.StringMinLength}}, {{.StringMaxLength}}),
																	{{- end}}
																	{{- range .StringPatterns}}
																	stringvalidator.RegexMatches(regexp.MustCompile(`{{.}}`), ""),
																	{{- end}}
																},
																{{- else if or (ne .MinInt 0) (ne .MaxInt 0)}}
																Validators: []validator.Int64{
																	int64validator.Between({{.MinInt}}, {{.MaxInt}}),
																},
																{{- else if or (ne .MinFloat 0.0) (ne .MaxFloat 0.0)}}
																Validators: []validator.Float64{
																	float64validator.Between({{.MinFloat}}, {{.MaxFloat}}),
																},
																{{- end}}
																{{- if and (len .DefaultValue) (eq .Type "Int64")}}
																Default:             int64default.StaticInt64({{.DefaultValue}}),
																{{- else if and (len .DefaultValue) (eq .Type "Bool")}}
																Default:             booldefault.StaticBool({{.DefaultValue}}),
																{{- else if and (len .DefaultValue) (eq .Type "String")}}
																Default:             stringdefault.StaticString("{{.DefaultValue}}"),
																{{- end}}
																{{- if .CustomModifier}}
																PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
																	{{.CustomModifier}}{},
																},
																{{- else}}
																{{- if .RequiresReplace}}
																PlanModifiers: []planmodifier.{{if eq .Type "StringList"}}List{{else}}{{.Type}}{{end}}{
																	{{if eq .Type "StringList"}}list{{else}}{{snakeCase .Type}}{{end}}planmodifier.RequiresReplace(),
																},
																{{- end}}
																{{- if and .Computed (not .ComputedRefreshValue) (not .NoUseStateForUnknown)}}
																PlanModifiers: []planmodifier.{{.Type}}{
																	{{snakeCase .Type}}planmodifier.UseStateForUnknown(),
																},
																{{- end}}
																{{- end}}
															},
															{{- end}}
															{{- end}}
														},
													},
													{{- if or (ne .MinList 0) (ne .MaxList 0)}}
													Validators: []validator.{{.Type}}{
														{{- if ne .MinList 0}}
														{{toLower .Type}}validator.SizeAtLeast({{.MinList}}),
														{{- end}}
														{{- if ne .MaxList 0}}
														{{toLower .Type}}validator.SizeAtMost({{.MaxList}}),
														{{- end}}
													},
													{{- end}}
													{{- end}}
														},
														{{- end}}
														{{- end}}
													},
												},
												{{- if or (ne .MinList 0) (ne .MaxList 0)}}
												Validators: []validator.{{.Type}}{
													{{- if ne .MinList 0}}
													{{toLower .Type}}validator.SizeAtLeast({{.MinList}}),
													{{- end}}
													{{- if ne .MaxList 0}}
													{{toLower .Type}}validator.SizeAtMost({{.MaxList}}),
													{{- end}}
												},
												{{- end}}
												{{- end}}
													},
													{{- end}}
													{{- end}}
												},
											},
											{{- if or (ne .MinList 0) (ne .MaxList 0)}}
											Validators: []validator.{{.Type}}{
												{{- if ne .MinList 0}}
												{{toLower .Type}}validator.SizeAtLeast({{.MinList}}),
												{{- end}}
												{{- if ne .MaxList 0}}
												{{toLower .Type}}validator.SizeAtMost({{.MaxList}}),
												{{- end}}
											},
											{{- end}}
											{{- end}}
												},
												{{- end}}
												{{- end}}
											},
										},
										{{- if or (ne .MinList 0) (ne .MaxList 0)}}
										Validators: []validator.{{.Type}}{
											{{- if ne .MinList 0}}
											{{toLower .Type}}validator.SizeAtLeast({{.MinList}}),
											{{- end}}
											{{- if ne .MaxList 0}}
											{{toLower .Type}}validator.SizeAtMost({{.MaxList}}),
											{{- end}}
										},
										{{- end}}
										{{- end}}
									},
									{{- end}}
									{{- end}}
								},
							},
							{{- if or (ne .MinList 0) (ne .MaxList 0)}}
							Validators: []validator.{{.Type}}{
								{{- if ne .MinList 0}}
								{{toLower .Type}}validator.SizeAtLeast({{.MinList}}),
								{{- end}}
								{{- if ne .MaxList 0}}
								{{toLower .Type}}validator.SizeAtMost({{.MaxList}}),
								{{- end}}
							},
							{{- end}}
							{{- end}}
						},
						{{- end}}
						{{- end}}
					},
				},
				{{- if or (ne .MinList 0) (ne .MaxList 0)}}
				Validators: []validator.{{.Type}}{
					{{- if ne .MinList 0}}
					{{toLower .Type}}validator.SizeAtLeast({{.MinList}}),
					{{- end}}
					{{- if ne .MaxList 0}}
					{{toLower .Type}}validator.SizeAtMost({{.MaxList}}),
					{{- end}}
				},
				{{- end}}
				{{- end}}
			},
			{{- end}}
			{{- end}}
		},
	}
}

func (r *{{camelCase .Name}}Resource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.clients = req.ProviderData.(*CwProviderData).Clients
	r.AllowExistingOnCreate = req.ProviderData.(*CwProviderData).AllowExistingOnCreate
	r.cache = req.ProviderData.(*CwProviderData).Cache
	{{- if hasNsoDevices .}}
	r.deviceReady = req.ProviderData.(*CwProviderData).DeviceReady
	{{- end}}
}

// clientFor returns the client of the selected instance, bound to ctx. An empty
// instance selects the default client. An unknown instance adds an error to diags.
func (r *{{camelCase .Name}}Resource) clientFor(ctx context.Context, instance string, diags *diag.Diagnostics) (cw.Client, bool) {
	return clientFor(ctx, r.clients, instance, diags)
}
{{- if or (coexistingSecretAttributes .) (coexistingSecretParentLists .)}}

// ValidateConfig enforces the relationship between a secret attribute, its write-only
// "_wo" counterpart and the "_wo_version" rotation trigger.
//
// The checks are not schema validators because Terraform renders an attribute-scoped
// diagnostic together with the offending configuration line, which for a secret prints
// the value into plan output and logs. A resource-scoped diagnostic names the attributes
// (and the list element index) instead.
func (r *{{camelCase .Name}}Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	{{- range coexistingSecretAttributes .}}
	{{- $go := toGoName .TfName}}
	var legacy{{$go}} types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("{{.TfName}}"), &legacy{{$go}})...)
	var wo{{$go}} types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("{{.TfName}}_wo"), &wo{{$go}})...)
	{{- if .WoPairHasVersion}}
	var woVersion{{$go}} types.Int64
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("{{.TfName}}_wo_version"), &woVersion{{$go}})...)
	{{- end}}
	if !legacy{{$go}}.IsUnknown() && !wo{{$go}}.IsUnknown() && !legacy{{$go}}.IsNull() && !wo{{$go}}.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Only one of `{{.TfName}}` and `{{.TfName}}_wo` can be set.",
		)
	}
	{{- if .WoPairMandatory}}
	if !legacy{{$go}}.IsUnknown() && !wo{{$go}}.IsUnknown() && legacy{{$go}}.IsNull() && wo{{$go}}.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"Exactly one of `{{.TfName}}` and `{{.TfName}}_wo` must be set.",
		)
	}
	{{- end}}
	{{- if .WoPairHasVersion}}
	if !wo{{$go}}.IsUnknown() && !woVersion{{$go}}.IsUnknown() && !wo{{$go}}.IsNull() && woVersion{{$go}}.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid Attribute Combination",
			"`{{.TfName}}_wo_version` must be set when `{{.TfName}}_wo` is used. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.",
		)
	}
	{{- end}}
	{{- end}}
	{{- range coexistingSecretParentLists .}}
	{{- $parentTf := .TfName}}
	{{- $parentGo := toGoName .TfName}}
	{
		// Secrets nested in "{{$parentTf}}" are validated per element. A list that cannot be
		// read as a whole - because it is still unknown at validation time - is skipped
		// rather than reported, since there is nothing to check yet.
		var cfg{{$parentGo}} []{{camelCase $.Name}}{{$parentGo}}
		if diags := req.Config.GetAttribute(ctx, path.Root("{{$parentTf}}"), &cfg{{$parentGo}}); !diags.HasError() {
			for i := range cfg{{$parentGo}} {
				{{- range coexistingSecretChildren .}}
				{{- $go := toGoName .TfName}}
				if !cfg{{$parentGo}}[i].{{$go}}.IsUnknown() && !cfg{{$parentGo}}[i].{{$go}}Wo.IsUnknown() && !cfg{{$parentGo}}[i].{{$go}}.IsNull() && !cfg{{$parentGo}}[i].{{$go}}Wo.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("Only one of `{{.TfName}}` and `{{.TfName}}_wo` can be set in `{{$parentTf}}` element %d.", i),
					)
				}
				{{- if .WoPairMandatory}}
				if !cfg{{$parentGo}}[i].{{$go}}.IsUnknown() && !cfg{{$parentGo}}[i].{{$go}}Wo.IsUnknown() && cfg{{$parentGo}}[i].{{$go}}.IsNull() && cfg{{$parentGo}}[i].{{$go}}Wo.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("Exactly one of `{{.TfName}}` and `{{.TfName}}_wo` must be set in `{{$parentTf}}` element %d.", i),
					)
				}
				{{- end}}
				{{- if .WoPairHasVersion}}
				if !cfg{{$parentGo}}[i].{{$go}}Wo.IsUnknown() && !cfg{{$parentGo}}[i].{{$go}}WoVersion.IsUnknown() && !cfg{{$parentGo}}[i].{{$go}}Wo.IsNull() && cfg{{$parentGo}}[i].{{$go}}WoVersion.IsNull() {
					resp.Diagnostics.AddError(
						"Invalid Attribute Combination",
						fmt.Sprintf("`{{.TfName}}_wo_version` must be set when `{{.TfName}}_wo` is used in `{{$parentTf}}` element %d. The write-only value is not stored in state, so Terraform can only detect a change to it through the version.", i),
					)
				}
				{{- end}}
				{{- end}}
			}
		}
	}
	{{- end}}
}
{{- end}}
// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create
func (r *{{camelCase .Name}}Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan {{camelCase .Name}}

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	{{- if hasNsoDevices .}}

	if !waitForNSODevices(ctx, &c, plan.Instance.ValueString(), r.deviceReady, plan.nsoDevices(), &resp.Diagnostics) {
		return
	}
	{{- end}}

	{{- if .UseCache}}
	cacheKey := "{{camelCase .Name}}::" + plan.Instance.ValueString() + "::"
	r.cache.DeletePattern(cacheKey)
	{{- end}}
	{{- range .Attributes}}
	{{- if .WriteOnlyTF}}
	// Write-only values are not in plan or state; read them from config.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("{{.TfName}}"), &plan.{{toGoName .TfName}})...)
	if resp.Diagnostics.HasError() {
		return
	}
	{{- end}}
	{{- end}}
	{{- range writeOnlyTFParentLists .}}
	{{- $parentTf := .TfName}}
	{{- $parentGo := toGoName .TfName}}
	// Write-only values in the list are not in plan or state; copy them from the config list into the plan elements.
	{
		var cfg{{$parentGo}} []{{camelCase $.Name}}{{$parentGo}}
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("{{$parentTf}}"), &cfg{{$parentGo}})...)
		if resp.Diagnostics.HasError() {
			return
		}
		for i := range plan.{{$parentGo}} {
			if i < len(cfg{{$parentGo}}) {
				{{- range writeOnlyTFChildren .}}
				plan.{{$parentGo}}[i].{{toGoName .TfName}} = cfg{{$parentGo}}[i].{{toGoName .TfName}}
				{{- end}}
			}
		}
	}
	{{- end}}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Create object
	body := plan.toBody(ctx, {{camelCase .Name}}{})
	{{- if .PostBody}}
	body = plan.postBody(body)
	{{- end}}

	{{- $resName := .Name }}
	{{- $MaxElementsInRootList := 0 }} 
	{{- range  .Attributes}}
	{{- if .MaxElementsInRootList}}
	var planList []{{camelCase $resName}}
	{{- $MaxElementsInRootList = .MaxElementsInRootList }} 
	maxElementsPerShard := {{.MaxElementsInRootList}}
	originalList := plan.{{toGoName .TfName}}
	{{- if isNestedMap .}}
	chunkMap := make(map[string]{{camelCase $resName}}{{toGoName .TfName}})
	count := 0
	for k, v := range originalList {
		chunkMap[k] = v
		count++
		if count >= maxElementsPerShard {
			currentPlanForShard := plan 
			currentPlanForShard.{{toGoName .TfName}} = chunkMap
			planList = append(planList, currentPlanForShard)
			chunkMap = make(map[string]{{camelCase $resName}}{{toGoName .TfName}})
			count = 0
		}
	}
	if count > 0 {
		currentPlanForShard := plan 
		currentPlanForShard.{{toGoName .TfName}} = chunkMap
		planList = append(planList, currentPlanForShard)
	}
	{{- else}}
	for i := 0; i < len(originalList); i += maxElementsPerShard {
		end := i+maxElementsPerShard
		if end > len(originalList){
			end = len(originalList)
		}
		chunk := originalList[i:end]
		currentPlanForShard := plan 
		currentPlanForShard.{{toGoName .TfName}} = chunk
		planList = append(planList, currentPlanForShard)

	}
	{{- end}}
	{{- end}}
	{{- end}}

	params := ""
	{{- if hasCreateQueryPath .Attributes}}
		{{- $createQueryPath := getCreateQueryPath .Attributes}}
	params += "/" + url.QueryEscape(plan.{{toGoName $createQueryPath.TfName}}.Value{{$createQueryPath.Type}}())
	{{- end}}

	{{- if $MaxElementsInRootList}}
	var err error
	var res gjson.Result
	for _, pl := range planList{
		body = pl.toBody(ctx, {{camelCase .Name}}{})
		{{- if .PutCreate}}
		res, err = c.Put(plan.getPath() + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
		{{- else}}
		res, err = c.Post(plan.getPath() + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
		{{- end}}
		if err != nil {
		{{- if .DeviceUnreachabilityWarning}}
			errorCode := res.Get("response.errorCode").String()
			failureReason := res.Get("response.failureReason").String()
			deviceFailureMatch, _ := regexp.MatchString(`(?i)Operation failed on '\d+' devices`, failureReason)
			if deviceFailureMatch {
				resp.Diagnostics.AddWarning("Device Unreachability Warning", fmt.Sprintf("Device unreachability detected (error code: %s, reason %s).", errorCode, failureReason))
			} else {
				{{- if and .AllowExistingOnCreate (not .PutCreate)}}	
				if r.AllowExistingOnCreate  {
					tflog.Info(ctx, fmt.Sprintf("Failed to configure object (%s), got error: %s, %s. allow_existing_on_create is true, beginning update", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
				} else{
					resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
					break
				}
				{{- else}}
				resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
				break
				{{- end}}
			}
		{{- else}}
			{{- if and .AllowExistingOnCreate (not .PutCreate)}}	
			if r.AllowExistingOnCreate  {
				tflog.Info(ctx, fmt.Sprintf("Failed to configure object (%s), got error: %s, %s. allow_existing_on_create is true, beginning update", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
			} else{
				resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
				break
			}
			{{- else}}
			resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
			break
			{{- end}}
		{{- end}}
		}
	}
	{{- else}}
	{{- if .PutCreate}}
	res, err := c.Put(plan.getPath() + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- else}}
	res, err := c.Post({{if .PostRestEndpoint}}{{if strContains .PostRestEndpoint "%v"}}plan.getPathPost(){{else}}"{{.PostRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- end}}
	{{- if .PlanEndpoint}}
	{{- $planId := getId .Attributes}}
	if err != nil && (cw.IsPlanTimeout(err) || cw.IsPlanFailed(err)) {
		// NSO accepted the object, but its plan did not complete. Keep the object in the state, which the error marks
		// as tainted, so that the next apply replaces or removes it instead of failing with "object already exists".
		plan.Id = types.StringValue(fmt.Sprint(plan.{{toGoName $planId.TfName}}.Value{{$planId.Type}}()))
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("The object was created, but its plan did not complete, so it is tracked as tainted: %s", err))
		return
	}
	{{- end}}
	{{- if .RetryOnErrorCodes}}
	if err != nil {
		retryErrorCodes := []string{ {{- range $i, $code := .RetryOnErrorCodes}}{{if $i}}, {{end}}"{{$code}}"{{- end}} }
		errorCode := res.Get("response.errorCode").String()
		
		shouldRetry := false
		for _, code := range retryErrorCodes {
			if errorCode == code {
				shouldRetry = true
				break
			}
		}
		
		if shouldRetry {
			maxWaitTime := time.Duration(retryMaxWaitSeconds) * time.Second
			startTime := time.Now()
			retryInterval := 15 * time.Second
			
			for shouldRetry && time.Since(startTime) < maxWaitTime {
				tflog.Warn(ctx, fmt.Sprintf("%s: Error code %s encountered, waiting %v before retry (elapsed: %v, max: %v)", 
					plan.Id.ValueString(), errorCode, retryInterval, time.Since(startTime), maxWaitTime))
				time.Sleep(retryInterval)
				
				{{- if .PutCreate}}
				res, err = c.Put(plan.getPath() + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
				{{- else}}
				res, err = c.Post(plan.getPath() + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
				{{- end}}
				
				if err == nil {
					shouldRetry = false
				} else {
					errorCode = res.Get("response.errorCode").String()
					shouldRetry = false
					for _, code := range retryErrorCodes {
						if errorCode == code {
							shouldRetry = true
							break
						}
					}
				}
			}
		}
	}
	{{- end}}
	if err != nil {
	{{- if .DeviceUnreachabilityWarning}}
		errorCode := res.Get("response.errorCode").String()
		failureReason := res.Get("response.failureReason").String()
		deviceFailureMatch, _ := regexp.MatchString(`(?i)Operation failed on '\d+' devices`, failureReason)
		if deviceFailureMatch {
			resp.Diagnostics.AddWarning("Device Unreachability Warning", fmt.Sprintf("Device unreachability detected (error code: %s, reason %s).", errorCode, failureReason))
		} else {
			{{- if and .AllowExistingOnCreate (not .PutCreate)}}	
			if r.AllowExistingOnCreate  {
				tflog.Info(ctx, fmt.Sprintf("Failed to configure object (%s), got error: %s, %s. allow_existing_on_create is true, beginning update", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
			} else{
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
				return
			}
			{{- else}}
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
			return
			{{- end}}
		}
	{{- else}}
		{{- if and .AllowExistingOnCreate (not .PutCreate)}}	
		if r.AllowExistingOnCreate  {
			tflog.Info(ctx, fmt.Sprintf("Failed to configure object (%s), got error: %s, %s. allow_existing_on_create is true, beginning update", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
		} else{
			{{- if .FallbackRestEndpoint }}
			if strings.Contains(err.Error(), "StatusCode 404") {
				{{- if .PutCreate}}
				res, err = c.Put(plan.getFallbackPath() + params, body {{- if .MaxAsyncWaitTime }}, func(r *cw.Req) { r.MaxAsyncWaitTime={{.MaxAsyncWaitTime}} }{{end}}{{- if .Mutex }}, cw.UseMutex{{- end}})
				{{- else}}
				res, err = c.Post(plan.getFallbackPath() + params, body {{- if .MaxAsyncWaitTime }}, func(r *cw.Req) { r.MaxAsyncWaitTime={{.MaxAsyncWaitTime}} }{{end}}{{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .NoWait }}, cw.NoWait{{- end}})
				{{- end}}
			}
			if err != nil {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
				return
			}
			{{- else}}
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
			return
			{{- end}}
		}
		{{- else}}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
		return
		{{- end}}
	{{- end}}
	}
	{{- end}}


	{{- /* Check if id can be resolved directly from response */}}
	{{- if .IdPath}}
	plan.Id = types.StringValue(res.Get("{{.IdPath}}").String())
	{{- /* Check if we need an extra query to resolve the id */}}
	{{- else if and .IdFromQueryPath (not .IdFromAttribute)}}
		{{- $id := getMatchId .Attributes}}
	params = ""
		{{- if hasQueryParam .Attributes}}
		{{- $queryParams := generateQueryParamString "GET" "plan" .Attributes }}
		{{- if $queryParams }}
	params += {{$queryParams}}
		{{- end}}
		{{- end}}
	{{- if .GetExtraQueryParams}}
	params += "{{.GetExtraQueryParams}}"
	{{- end}}
	{{- if and .FallbackRestEndpoint $id.FallbackResponseModelName }}
	useFallback := false
	{{- end}}
	res, err = {{readGet . "plan" false}}{{if .IdQueryRestEndpoint}}plan.getPathIdQuery(){{else if .GetRestEndpoint}}{{if strContains .GetRestEndpoint "%v"}}plan.getPathGet(){{else}}"{{.GetRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + params)
	{{- if .FallbackRestEndpoint }}

	if err != nil && (strings.Contains(err.Error(), "StatusCode 404") || strings.Contains(err.Error(), "StatusCode 500")) {
		tflog.Debug(ctx, fmt.Sprintf("%s: Primary endpoint returned 404 or 500, trying fallback endpoint", plan.Id.ValueString()))
		res, err = {{readGet . "plan" false}}plan.getFallbackPath() + params)
		{{- if $id.FallbackResponseModelName }}
		useFallback = true
		{{- end}}
	}
	{{- end}}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
		return
	}
	{{- if and .IdFromQueryPathAttribute .IdFromQueryPath (not .GetExtraQueryParams) (not .GetFromAll) (not (hasMatchId .Attributes)) }}
	plan.Id = types.StringValue(res.Get("{{if eq .IdFromQueryPath "." }}{{else}}{{.IdFromQueryPath}}.{{end}}{{.IdFromQueryPathAttribute}}").String())
	{{- else}}
	{{- if and .FallbackRestEndpoint $id.FallbackResponseModelName }}
	if useFallback {
		plan.Id = types.StringValue(res.Get("{{.IdFromQueryPath}}.#({{$id.FallbackResponseModelName}}==\""+ {{if eq $id.Type "Int64"}}strconv.FormatInt(plan.{{toGoName $id.TfName}}.ValueInt64(), 10){{else}}plan.{{toGoName $id.TfName}}.Value{{$id.Type}}(){{end}} +"\").{{if .IdFromQueryPathAttribute}}{{.IdFromQueryPathAttribute}}{{else}}id{{end}}").String())
	} else {
		plan.Id = types.StringValue(res.Get("{{.IdFromQueryPath}}.#({{if $id.ResponseModelName}}{{$id.ResponseModelName}}{{else}}{{$id.ModelName}}{{end}}==\""+ {{if eq $id.Type "Int64"}}strconv.FormatInt(plan.{{toGoName $id.TfName}}.ValueInt64(), 10){{else}}plan.{{toGoName $id.TfName}}.Value{{$id.Type}}(){{end}} +"\").{{if .IdFromQueryPathAttribute}}{{.IdFromQueryPathAttribute}}{{else}}id{{end}}").String())
	}
	{{- else}}
	plan.Id = types.StringValue(res.Get("{{.IdFromQueryPath}}.#({{if $id.ResponseModelName}}{{$id.ResponseModelName}}{{else}}{{$id.ModelName}}{{end}}==\""+ {{if eq $id.Type "Int64"}}strconv.FormatInt(plan.{{toGoName $id.TfName}}.ValueInt64(), 10){{else}}plan.{{toGoName $id.TfName}}.Value{{$id.Type}}(){{end}} +"\").{{if .IdFromQueryPathAttribute}}{{.IdFromQueryPathAttribute}}{{else}}id{{end}}").String())
	{{- end}}
	{{- end}}
	{{- /* If we have an id attribute we will use that as id */}}
	{{- else if hasId .Attributes}}
		{{- $id := getId .Attributes}}
	plan.Id = types.StringValue(fmt.Sprint(plan.{{toGoName $id.TfName}}.Value{{$id.Type}}()))
	{{- end}}
	{{- if .FailCreateOnEmptyId}}
	if plan.Id.ValueString() == "" {
		resp.Diagnostics.AddError("Client Error", "Failed to retrieve object ID: the object was not found by the ID query after creation")
		return
	}
	{{- end}}
	{{- if .IdFromQueryPathAttribute}}
	{{- range .Attributes}}
	{{- if and .Reference .Computed (eq .ModelName $.IdFromQueryPathAttribute)}}
	plan.{{toGoName .TfName}} = plan.Id
	{{- end}}
	{{- end}}
	{{- end}}

	{{- if and .UpdateComputed .RootList}}
	params = ""
	{{- $queryParams := generateQueryParamString "GET" "plan" .Attributes }}

	{{- if .IdQueryParam}}
	params += "?{{.IdQueryParam}}=" + url.QueryEscape(plan.Id.ValueString())
	{{- else if and (hasQueryParam .Attributes) (not .GetRequiresId)}}
	{{- if $queryParams }}
	params += {{$queryParams}}
	{{- end}}
	{{- else if and (not .GetNoId) (not .GetFromAll)}}
	params += "/" + url.QueryEscape(plan.Id.ValueString())
	{{- end}}
	{{- if .GetExtraQueryParams}}
	params += "{{.GetExtraQueryParams}}"
	{{- end}}
	res, err = {{readGet . "plan" false}}{{if .GetRestEndpoint}}{{if strContains .GetRestEndpoint "%v"}}plan.getPathGet(){{else}}"{{.GetRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + params)
	{{- if .FallbackRestEndpoint }}

	if err != nil && (strings.Contains(err.Error(), "StatusCode 404") || strings.Contains(err.Error(), "StatusCode 500")) {
		tflog.Debug(ctx, fmt.Sprintf("%s: Primary endpoint returned 404 or 500, trying fallback endpoint", plan.Id.ValueString()))
		res, err = {{readGet . "plan" false}}plan.getFallbackPath() + params)
	}
	{{- end}}
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
		return
	}

	{{- if .GetFromAll}}
		{{- if .IdFromAttribute}}
			{{- $id := getId .Attributes}}
	res = res.Get("{{.IdFromQueryPath}}.#({{if $id.ResponseModelName}}{{$id.ResponseModelName}}{{else}}{{$id.ModelName}}{{end}}==\"" + state.{{toGoName $id.TfName}}.Value{{$id.Type}}() + "\")")
		{{- else}}
	res = res.Get("{{.IdFromQueryPath}}.#({{if .IdFromQueryPathAttribute}}{{.IdFromQueryPathAttribute}}{{else}}id{{end}}==\"" + state.Id.ValueString() + "\")")
		{{- end}}
	{{- end}}
	plan.fromBodyUnknowns(ctx, res)
	{{- end}}

	{{- if and .UpdateComputed (not .RootList)}}
	plan.fromBodyUnknowns(ctx, res)
	{{- end}}
	


	{{- if and .AllowExistingOnCreate (not .PutCreate)}}	
	if !r.AllowExistingOnCreate  {
		tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))
	} else {
		params = ""
		{{- if hasCreateQueryPath .Attributes}}
		{{- $createQueryPath := getCreateQueryPath .Attributes}}
		params += "/" + url.QueryEscape(plan.{{toGoName $createQueryPath.TfName}}.Value{{$createQueryPath.Type}}())
		{{- end}}
		{{- if .PutIdQueryParam}}
		params += "?{{.PutIdQueryParam}}=" + url.QueryEscape(plan.Id.ValueString())
		{{- end}}
		body = plan.toBody(ctx, {{camelCase .Name}}{Id: plan.Id})
		{{- if .PostUpdate}}
		res, err = c.Post(plan.getPath() + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
		{{- else if .PutCreate}}
		res, err = c.Put(plan.getPath() + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
		{{- else if and (hasQueryParam .Attributes) (not .PutId)}}
		res, err = c.Put({{if .PutRestEndpoint}}{{if strContains .PutRestEndpoint "%v"}}plan.getPathPut(){{else}}"{{.PutRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
		{{- else if .PutNoId}}
		res, err = c.Put({{if .PutRestEndpoint}}{{if strContains .PutRestEndpoint "%v"}}plan.getPathPut(){{else}}"{{.PutRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
		{{- else}}
		res, err = c.Put({{if .PutRestEndpoint}}{{if strContains .PutRestEndpoint "%v"}}plan.getPathPut(){{else}}"{{.PutRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + "/" + url.QueryEscape(plan.Id.ValueString()) + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
		{{- end}}

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "PUT", err, res.String()))
			return
		}
		tflog.Debug(ctx, fmt.Sprintf("%s: Fallback to update existing resource finished successfully", plan.Id.ValueString()))
	}
	{{- end}}
	
	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}
// End of section. //template:end create

// Section below is generated&owned by "gen/generator.go". //template:begin read
func (r *{{camelCase .Name}}Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state {{camelCase .Name}}

	// Read state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

{{if .NoRead}}	_ = c
{{end}}	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", state.Id.String()))
	{{- if not .NoRead}}

	params := ""
	{{- $queryParams := generateQueryParamString "GET" "state" .Attributes }}

	{{- if .IdQueryParam}}
	params += "?{{.IdQueryParam}}=" + url.QueryEscape(state.Id.ValueString())
	{{- else if and (hasQueryParam .Attributes) (not .GetRequiresId)}}
	{{- if $queryParams }}
	params += {{$queryParams}}
	{{- end}}
	{{- else if and (not .GetNoId) (not .GetFromAll) (not (and .GetRestEndpoint (strContains .GetRestEndpoint "%v")))}}
	params += "/" + url.QueryEscape(state.Id.ValueString())
	{{- end}}
	{{- if hasGetQueryParam .Attributes }}
	params += {{$queryParams}}
	{{- end }}
	{{- if .GetExtraQueryParams}}
	params += "{{.GetExtraQueryParams}}"
	{{- end}}
	{{- if .UseCache}}
	res, err := r.ReadCache(ctx, &c, req, state, params)
	{{- else}}
	res, err := {{readGet . "state" false}}{{if .GetRestEndpoint}}{{if strContains .GetRestEndpoint "%v"}}state.getPathGet(){{else}}"{{.GetRestEndpoint}}"{{end}}{{else}}state.getPath(){{end}} + params)
	{{- end}}
	{{- if .FallbackRestEndpoint }}

	if err != nil && (strings.Contains(err.Error(), "StatusCode 404") || strings.Contains(err.Error(), "StatusCode 500")) {
		tflog.Debug(ctx, fmt.Sprintf("%s: Primary endpoint returned 404 or 500, trying fallback endpoint", state.Id.ValueString()))
		res, err = {{readGet . "state" false}}state.getFallbackPath() + params)
	}
	{{- end}}
	if err != nil && (strings.Contains(err.Error(), "StatusCode 404") || strings.Contains(err.Error(), "StatusCode 406") || strings.Contains(err.Error(), "StatusCode 500") || strings.Contains(err.Error(), "StatusCode 400") || cw.IsNotFound(err)) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
		return
	}

	{{- if .GetFromAll}}
		{{- if .IdFromAttribute}}
			{{- $id := getId .Attributes}}
	res = res.Get("{{.IdFromQueryPath}}.#({{if $id.ResponseModelName}}{{$id.ResponseModelName}}{{else}}{{$id.ModelName}}{{end}}==\"" + state.{{toGoName $id.TfName}}.Value{{$id.Type}}() + "\")")
		{{- else}}
	res = res.Get("{{.IdFromQueryPath}}.#({{if .IdFromQueryPathAttribute}}{{.IdFromQueryPathAttribute}}{{else}}id{{end}}==\"" + state.Id.ValueString() + "\")")
		{{- end}}
	if !res.Exists() {
		resp.State.RemoveResource(ctx)
		return
	}
	{{- end}}

	{{- if .RemoveResourceOnEmptyResponse}}

	if data := res.Get("{{responseEnvelope .Attributes}}"); !data.Exists() || len(data.Array()) == 0 {
		resp.State.RemoveResource(ctx)
		return
	}
	{{- end}}

	// If every attribute is set to null we are dealing with an import operation and therefore reading all attributes
	if state.isNull(ctx, res) {
		state.fromBody(ctx, res)
	} else {
		state.updateFromBody(ctx, res)
	}
	{{- end}}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
// End of section. //template:end read

// Section below is generated&owned by "gen/generator.go". //template:begin update
func (r *{{camelCase .Name}}Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var plan, state {{camelCase .Name}}
	{{- $name := camelCase .Name}}

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Read state
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, plan.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}
	{{- if hasNsoDevices .}}

	if !waitForNSODevices(ctx, &c, plan.Instance.ValueString(), r.deviceReady, nsoDevicesAdded(plan.nsoDevices(), state.nsoDevices()), &resp.Diagnostics) {
		return
	}
	{{- end}}

	{{- if .UseCache}}
	cacheKey := "{{camelCase .Name}}::" + plan.Instance.ValueString() + "::"
	r.cache.DeletePattern(cacheKey)
	{{- end}}
	{{- range .Attributes}}
	{{- if .WriteOnlyTF}}
	// Read on every Update: a full-object replace PUT is rejected when it omits the secret; the version companion only drives change detection.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("{{.TfName}}"), &plan.{{toGoName .TfName}})...)
	if resp.Diagnostics.HasError() {
		return
	}
	{{- end}}
	{{- end}}
	{{- range writeOnlyTFParentLists .}}
	{{- $parentTf := .TfName}}
	{{- $parentGo := toGoName .TfName}}
	// Read unconditionally, for the same reason as the top-level secrets above.
	{
		var cfg{{$parentGo}} []{{camelCase $.Name}}{{$parentGo}}
		resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("{{$parentTf}}"), &cfg{{$parentGo}})...)
		if resp.Diagnostics.HasError() {
			return
		}
		for i := range plan.{{$parentGo}} {
			if i < len(cfg{{$parentGo}}) {
				{{- range writeOnlyTFChildren .}}
				plan.{{$parentGo}}[i].{{toGoName .TfName}} = cfg{{$parentGo}}[i].{{toGoName .TfName}}
				{{- end}}
			}
		}
	}
	{{- end}}

{{if .NoUpdate}}	_ = c
{{end}}	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))
	{{- if not .NoUpdate}}
	{{- if or (not .UpdateComputed) (not .RootList)}}
	{{- if getIfUnsetOnUpdateAttributes .}}

	// An attribute the controller assigns itself is null in state and plan, and a full-object PUT that omits it counts as a removal:
	// fetch the current value and send it. A value that is non-null in state is left to the plan, so clearing it in the configuration still
	// removes it. Works on a copy of the plan.
	bodyPlan := plan
	getIfUnsetNeeded := false
	{{- range getIfUnsetOnUpdateAttributes .}}
	if bodyPlan.{{toGoName .TfName}}.{{if eq .Type "Int64"}}ValueInt64() == 0{{else}}ValueString() == ""{{end}} && state.{{toGoName .TfName}}.IsNull() {
		getIfUnsetNeeded = true
	}
	{{- end}}
	if getIfUnsetNeeded {
		getIfUnsetParams := ""
		{{- $queryParams := generateQueryParamString "GET" "state" .Attributes }}
		{{- if .IdQueryParam}}
		getIfUnsetParams += "?{{.IdQueryParam}}=" + url.QueryEscape(state.Id.ValueString())
		{{- else if and (hasQueryParam .Attributes) (not .GetRequiresId)}}
		{{- if $queryParams }}
		getIfUnsetParams += {{$queryParams}}
		{{- end}}
		{{- else if and (not .GetNoId) (not .GetFromAll) (not (and .GetRestEndpoint (strContains .GetRestEndpoint "%v")))}}
		getIfUnsetParams += "/" + url.QueryEscape(state.Id.ValueString())
		{{- end}}
		{{- if hasGetQueryParam .Attributes }}
		getIfUnsetParams += {{$queryParams}}
		{{- end }}
		{{- if .GetExtraQueryParams}}
		getIfUnsetParams += "{{.GetExtraQueryParams}}"
		{{- end}}
		{{- if .UseCache}}
		curRes, curErr := r.ReadCache(ctx, &c, resource.ReadRequest{}, state, getIfUnsetParams)
		{{- else}}
		curRes, curErr := {{readGet . "state" false}}{{if .GetRestEndpoint}}{{if strContains .GetRestEndpoint "%v"}}state.getPathGet(){{else}}"{{.GetRestEndpoint}}"{{end}}{{else}}state.getPath(){{end}} + getIfUnsetParams)
		{{- end}}
		{{- if .FallbackRestEndpoint }}
		if curErr != nil && (strings.Contains(curErr.Error(), "StatusCode 404") || strings.Contains(curErr.Error(), "StatusCode 500")) {
			curRes, curErr = {{readGet . "state" false}}state.getFallbackPath() + getIfUnsetParams)
		}
		{{- end}}
		if curErr == nil {
			{{- if .GetFromAll}}
			{{- if .IdFromAttribute}}
			{{- $id := getId .Attributes}}
			curRes = curRes.Get("{{.IdFromQueryPath}}.#({{if $id.ResponseModelName}}{{$id.ResponseModelName}}{{else}}{{$id.ModelName}}{{end}}==\"" + state.{{toGoName $id.TfName}}.Value{{$id.Type}}() + "\")")
			{{- else}}
			curRes = curRes.Get("{{.IdFromQueryPath}}.#({{if .IdFromQueryPathAttribute}}{{.IdFromQueryPathAttribute}}{{else}}id{{end}}==\"" + state.Id.ValueString() + "\")")
			{{- end}}
			{{- end}}
			if curRes.Exists() {
				{{- range getIfUnsetOnUpdateAttributes .}}
				if bodyPlan.{{toGoName .TfName}}.{{if eq .Type "Int64"}}ValueInt64() == 0{{else}}ValueString() == ""{{end}} && state.{{toGoName .TfName}}.IsNull() {
					if v := curRes.Get("{{if .ResponseDataPath}}{{.ResponseDataPath}}{{else if .DataPath}}{{.DataPath}}.{{.ModelName}}{{else}}{{.ModelName}}{{end}}"); v.Exists() && {{if eq .Type "Int64"}}v.Int() != 0{{else}}v.String() != ""{{end}} {
						bodyPlan.{{toGoName .TfName}} = types.{{if eq .Type "Int64"}}Int64Value(v.Int()){{else}}StringValue(v.String()){{end}}
					}
				}
				{{- end}}
			}
		} else {
			tflog.Warn(ctx, fmt.Sprintf("%s: Unable to fetch current object to preserve server-assigned attributes: %s", state.Id.ValueString(), curErr))
		}
	}
	body := bodyPlan.toBody(ctx, state)
	{{- else}}

	body := plan.toBody(ctx, state)
	{{- end}}
	params := ""
	{{- if hasCreateQueryPath .Attributes}}
		{{- $createQueryPath := getCreateQueryPath .Attributes}}
	params += "/" + url.QueryEscape(plan.{{toGoName $createQueryPath.TfName}}.Value{{$createQueryPath.Type}}())
	{{- end}}
	{{- if .PutIdQueryParam}}
	params += "?{{.PutIdQueryParam}}=" + url.QueryEscape(plan.Id.ValueString())
	{{- end}}
	{{- if .PostUpdate}}
	res, err := c.Post({{if .PutRestEndpoint}}{{if strContains .PutRestEndpoint "%v"}}plan.getPathPut(){{else}}"{{.PutRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- else if .PutCreate}}
	res, err := c.Put(plan.getPath() + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- else if and (hasQueryParam .Attributes) (not .PutId)}}
	res, err := c.Put({{if .PutRestEndpoint}}{{if strContains .PutRestEndpoint "%v"}}plan.getPathPut(){{else}}"{{.PutRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- else if .PutNoId}}
	res, err := c.Put({{if .PutRestEndpoint}}{{if strContains .PutRestEndpoint "%v"}}plan.getPathPut(){{else}}"{{.PutRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- else}}
	res, err := c.Put({{if .PutRestEndpoint}}{{if strContains .PutRestEndpoint "%v"}}plan.getPathPut(){{else}}"{{.PutRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + "/" + url.QueryEscape(plan.Id.ValueString()) + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- end}}
	{{- if .RetryOnErrorCodes}}
	if err != nil {
		retryErrorCodes := []string{ {{- range $i, $code := .RetryOnErrorCodes}}{{if $i}}, {{end}}"{{$code}}"{{- end}} }
		errorCode := res.Get("response.errorCode").String()

		shouldRetry := false
		for _, code := range retryErrorCodes {
			if errorCode == code {
				shouldRetry = true
				break
			}
		}

		if shouldRetry {
			maxWaitTime := time.Duration(retryMaxWaitSeconds) * time.Second
			startTime := time.Now()
			retryInterval := 15 * time.Second

			for shouldRetry && time.Since(startTime) < maxWaitTime {
				tflog.Warn(ctx, fmt.Sprintf("%s: Error code %s encountered, waiting %v before retry (elapsed: %v, max: %v)",
					plan.Id.ValueString(), errorCode, retryInterval, time.Since(startTime), maxWaitTime))
				time.Sleep(retryInterval)

				{{- if .PostUpdate}}
				res, err = c.Post({{if .PutRestEndpoint}}{{if strContains .PutRestEndpoint "%v"}}plan.getPathPut(){{else}}"{{.PutRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
				{{- else if .PutCreate}}
				res, err = c.Put(plan.getPath() + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
				{{- else if and (hasQueryParam .Attributes) (not .PutId)}}
				res, err = c.Put({{if .PutRestEndpoint}}{{if strContains .PutRestEndpoint "%v"}}plan.getPathPut(){{else}}"{{.PutRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
				{{- else if .PutNoId}}
				res, err = c.Put({{if .PutRestEndpoint}}{{if strContains .PutRestEndpoint "%v"}}plan.getPathPut(){{else}}"{{.PutRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
				{{- else}}
				res, err = c.Put({{if .PutRestEndpoint}}{{if strContains .PutRestEndpoint "%v"}}plan.getPathPut(){{else}}"{{.PutRestEndpoint}}"{{end}}{{else}}plan.getPath(){{end}} + "/" + url.QueryEscape(plan.Id.ValueString()) + params, body {{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(plan.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
				{{- end}}

				if err == nil {
					shouldRetry = false
				} else {
					errorCode = res.Get("response.errorCode").String()
					shouldRetry = false
					for _, code := range retryErrorCodes {
						if errorCode == code {
							shouldRetry = true
							break
						}
					}
				}
			}
		}
	}
	{{- end}}
	if err != nil {
	{{- if .DeviceUnreachabilityWarning}}
		errorCode := res.Get("response.errorCode").String()
		failureReason := res.Get("response.failureReason").String()
		deviceFailureMatch, _ := regexp.MatchString(`(?i)Operation failed on '\d+' devices`, failureReason)
		if deviceFailureMatch {
			resp.Diagnostics.AddWarning("Device Unreachability Warning", fmt.Sprintf("Device unreachability detected (error code: %s, reason %s).", errorCode, failureReason))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to update object (%s), got error: %s, %s", {{- if .PostUpdate }} "POST" {{- else }} "PUT" {{- end }}, err, res.String()))
			return
		}
	{{- else}}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s, %s", err, res.String()))
		return
	{{- end}}
	}

	{{- if and .IdPath .PutUpdateId}}
	plan.Id = types.StringValue(res.Get("{{.IdPath}}").String())
	{{- end}}
	{{- end}}
	{{- end}}

	{{- if and .RootList .UpdateComputed}}{{"\n"}}
	{{- $items := "" }}
	{{- $itemsIsMap := false }}
	{{- range .Attributes}}
    {{- if isNestedListSetMap .}}
	{{- $items = .TfName }}
	{{- $itemsIsMap = isNestedMap . }}
	{{- end}}
	{{- end}}

	{{- if $itemsIsMap}}
	var toDelete = {{camelCase .Name}}{
		{{toGoName $items}}: make(map[string]{{camelCase .Name}}{{toGoName $items}}),
	}
	var toCreate = {{camelCase .Name}}{
		{{toGoName $items}}: make(map[string]{{camelCase .Name}}{{toGoName $items}}),
	}
	var toUpdate = {{camelCase .Name}}{
		{{toGoName $items}}: make(map[string]{{camelCase .Name}}{{toGoName $items}}),
	}
	var toReplace = {{camelCase .Name}}{
		{{toGoName $items}}: make(map[string]{{camelCase .Name}}{{toGoName $items}}),
	}
	{{- else}}
	var toDelete = {{camelCase .Name}}{
		{{toGoName $items}}: []{{camelCase .Name}}{{toGoName $items}}{},
	}
	var toCreate = {{camelCase .Name}}{
		{{toGoName $items}}: []{{camelCase .Name}}{{toGoName $items}}{},
	}
	var toUpdate = {{camelCase .Name}}{
		{{toGoName $items}}: []{{camelCase .Name}}{{toGoName $items}}{},
	}
	var toReplace = {{camelCase .Name}}{
		{{toGoName $items}}: []{{camelCase .Name}}{{toGoName $items}}{},
	}
	{{- end}}

	planMap := make(map[string]{{camelCase .Name}}{{toGoName $items}})
	stateMap := make(map[string]{{camelCase .Name}}{{toGoName $items}})

	{{- if $itemsIsMap}}
	for k, v := range state.{{toGoName $items}} {
		stateMap[k] = v
	}
	{{- else}}
	for _, v := range state.{{toGoName $items}} {
		{{- range .Attributes}}
		{{- $id := getId .Attributes}}
		{{- if not (eq (toGoName $id.TfName) "") }}
		stateMap[{{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if not .Computed}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(v.{{toGoName .TfName}}.ValueInt64(), 10){{else if eq .Type "Bool"}}strconv.FormatBool(v.{{toGoName .TfName}}.ValueBool()){{else if eq .Type "String"}}v.{{toGoName .TfName}}.Value{{.Type}}(){{end}}{{end}}{{end}}{{end}}] = v
		{{- end}}
		{{- end}}
	}
	{{- end}}

	{{- if $itemsIsMap}}
	for k, v := range plan.{{toGoName $items}} {
		planMap[k] = v
	}
	{{- else}}
	for _, v := range plan.{{toGoName $items}} {
		{{- range .Attributes}}
		{{- $id := getId .Attributes}}
		{{- if not (eq (toGoName $id.TfName) "") }}
		planMap[{{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if not .Computed}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(v.{{toGoName .TfName}}.ValueInt64(), 10){{else if eq .Type "Bool"}}strconv.FormatBool(v.{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}v.{{toGoName .TfName}}.Value{{.Type}}(){{end}}{{end}}{{end}}{{end}}] = v
		{{- end}}
		{{- end}}
	}
	{{- end}}

	for stateKey, stateItem := range stateMap {
		if _, exists := planMap[stateKey]; !exists {
			{{- if $itemsIsMap}}
			toDelete.{{toGoName $items}}[stateKey] = stateItem
			{{- else}}
			toDelete.{{toGoName $items}} = append(toDelete.{{toGoName $items}}, stateItem)
			{{- end}}
		}
	}

	for planKey, planItem := range planMap {
		if stateItem, exists := stateMap[planKey]; exists {
			if !reflect.DeepEqual(planItem, stateItem) {
				planItem.Id = stateItem.Id
				planMap[planKey] = planItem // Store back in planMap
				{{- range .Attributes}}
				{{- range .Attributes}}
				{{- if .RequiresReplace }}
				{{- if eq .Type "Set"}}
				var planR []string
				var stateR []string
				stateItem.{{toGoName .TfName}}.ElementsAs(ctx, &stateR, false)
				planItem.{{toGoName .TfName}}.ElementsAs(ctx, &planR, false)
				sort.Strings(stateR)
				sort.Strings(planR)
				if !reflect.DeepEqual(stateR, planR) {
				{{- else}}
				{{- if .NoPut }}
				if planItem.{{toGoName .TfName}} != stateItem.{{toGoName .TfName}} {
				{{- else}}
				if planItem.{{toGoName .TfName}} != stateItem.{{toGoName .TfName}} {
				{{- end}}
				{{- end}}
					{{- if $itemsIsMap}}
					toReplace.{{toGoName $items}}[planKey] = planItem
					{{- else}}
					toReplace.{{toGoName $items}} = append(toReplace.{{toGoName $items}}, planItem)
					{{- end}}
					continue
				}
				{{- end}}
				{{- end}}
				{{- end}}
				{{- if $itemsIsMap}}
				toUpdate.{{toGoName $items}}[planKey] = planItem
				{{- else}}
				toUpdate.{{toGoName $items}} = append(toUpdate.{{toGoName $items}}, planItem)
				{{- end}}
			}
		} else {
			{{- if $itemsIsMap}}
			toCreate.{{toGoName $items}}[planKey] = planItem
			{{- else}}
			toCreate.{{toGoName $items}} = append(toCreate.{{toGoName $items}}, planItem)
			{{- end}}
		}
	}

	if len(toReplace.{{toGoName $items}}) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("%s: Number of items to replace: %d", state.Id.ValueString(), len(toReplace.{{toGoName $items}})))
		{{- if $itemsIsMap}}
		var toReplaceNoId = {{camelCase .Name}}{
			{{toGoName $items}}: make(map[string]{{camelCase .Name}}{{toGoName $items}}),
		}
		for replaceKey, item := range toReplace.{{toGoName $items}} {
		{{- else}}
		var toReplaceNoId = {{camelCase .Name}}{
			{{toGoName $items}}: []{{camelCase .Name}}{{toGoName $items}}{},
		}
		for _, item := range toReplace.{{toGoName $items}} {
		{{- end}}
			{{- range .Attributes}}
			{{- $Atts := .Attributes}}
			{{- range .Attributes}}
			{{- if and .NoPut (not .ExcludeTest)}}
				{{$putName := toGoName .TfName}}
				{{$putVal := .Type}} 
				{{- range $Atts}}
				{{$transformedName := toGoName .TfName}}

				{{- if and (strContains $transformedName $putName) (ne $transformedName $putName)  }}
				if !item.{{$transformedName}}.IsNull() && item.{{$transformedName}}.Value{{.Type}}() {
					item.{{$putName}} = types.{{$putVal}}Null()
				}
				{{- end}}
				{{- end}}
			{{- else if and .Computed .ExcludeTest  (ne .Type "Set") }}
			item.{{toGoName .TfName}} = types.{{.Type}}Null()
			{{- end}}
			{{- end}}
			{{- end}}
			{{- if $itemsIsMap}}
			toReplaceNoId.{{toGoName $items}}[replaceKey] = item
			{{- else}}
			toReplaceNoId.{{toGoName $items}} = append(toReplaceNoId.{{toGoName $items}}, item)
			{{- end}}
		}

		{{- if $itemsIsMap}}
		for k, v := range toReplace.{{toGoName $items}} {
			toDelete.{{toGoName $items}}[k] = v
		}
		for k, v := range toReplaceNoId.{{toGoName $items}} {
			toCreate.{{toGoName $items}}[k] = v
		}
		{{- else}}
		toDelete.{{toGoName $items}} = append(toDelete.{{toGoName $items}}, toReplace.{{toGoName $items}}...)
		toCreate.{{toGoName $items}} = append(toCreate.{{toGoName $items}}, toReplaceNoId.{{toGoName $items}}...)
		{{- end}}
	}

	if len(toDelete.{{toGoName $items}}) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("%s: Number of items to delete: %d", state.Id.ValueString(), len(toDelete.{{toGoName $items}})))
		for _, v := range toDelete.{{toGoName $items}} {
			{{- if .SkipDeleteOnEmptyId}}
			// An empty id would send DELETE to the base endpoint.
			if v.Id.IsNull() || v.Id.IsUnknown() || v.Id.ValueString() == "" {
				tflog.Debug(ctx, fmt.Sprintf("%s: Skipping delete for item - ID is empty or null", state.Id.ValueString()))
				continue
			}
			{{- else}}
			if v.Id.IsNull() {
				continue // Skip if id is null
			}
			{{- end}}
			res, err := c.Delete(plan.getPath() + "/" + url.QueryEscape(v.Id.ValueString()){{- if .Mutex }}, cw.UseMutex{{- end}})
			if err != nil {
			{{- if .DeviceUnreachabilityWarning}}
				errorCode := res.Get("response.errorCode").String()
				failureReason := res.Get("response.failureReason").String()
				deviceFailureMatch, _ := regexp.MatchString(`(?i)Operation failed on '\d+' devices`, failureReason)
				if deviceFailureMatch {
					resp.Diagnostics.AddWarning("Device Unreachability Warning", fmt.Sprintf("Device unreachability detected (error code: %s, reason %s).", errorCode, failureReason))
				} else {
					resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (%s), got error: %s, %s", "DELETE", err, res.String()))
					return
				}
			{{- else}}
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (%s), got error: %s, %s", "DELETE", err, res.String()))
				return
			{{- end}}
			}
		}
	}

	if len(toCreate.{{toGoName $items}}) > 0 {
		{{$chunks := false}}
		{{- $resName := .Name }}
		{{- $MaxElementsInRootList := 0 }} 
		{{- range  .Attributes}}
		{{- if .MaxElementsInRootList}}
		{{- $MaxElementsInRootList = .MaxElementsInRootList }} 
		maxElementsPerShard := {{$MaxElementsInRootList}}
		var createList []{{$name}}
		{{- if $itemsIsMap}}
		chunkMap := make(map[string]{{camelCase $.Name}}{{toGoName $items}})
		count := 0
		for k, v := range toCreate.{{toGoName $items}} {
			chunkMap[k] = v
			count++
			if count >= maxElementsPerShard {
				currentPlanForShard := plan
				currentPlanForShard.{{toGoName $items}} = chunkMap
				createList = append(createList, currentPlanForShard)
				chunkMap = make(map[string]{{camelCase $.Name}}{{toGoName $items}})
				count = 0
			}
		}
		if count > 0 {
			currentPlanForShard := plan
			currentPlanForShard.{{toGoName $items}} = chunkMap
			createList = append(createList, currentPlanForShard)
		}
		{{- else}}
		for i := 0; i < len(toCreate.{{toGoName $items}}); i += maxElementsPerShard {
			end := min(i+maxElementsPerShard, len(toCreate.{{toGoName $items}}))
			chunk := toCreate.{{toGoName $items}}[i:end]
			currentPlanForShard := plan 
			currentPlanForShard.{{toGoName $items}} = chunk
			createList = append(createList, currentPlanForShard)

			
		}
		{{- end}}
		{{$chunks = true}}
		{{- end}}
		{{- end}}

		tflog.Debug(ctx, fmt.Sprintf("%s: Number of items to create: %d", state.Id.ValueString(), len(toCreate.{{toGoName $items}})))

		{{- if $chunks}}
		var err error
		var res gjson.Result
		params := ""
		for _, pl := range createList {
			body := pl.toBody(ctx, {{$name}}{}) // Convert to request body
			res, err := c.Post(plan.getPath()+params, body, cw.UseMutex)
			if err != nil {
				errorCode := res.Get("response.errorCode").String()
				failureReason := res.Get("response.failureReason").String()
				deviceFailureMatch, _ := regexp.MatchString(`(?i)Operation failed on '\d+' devices`, failureReason)
				if deviceFailureMatch {
					resp.Diagnostics.AddWarning("Device Unreachability Warning", fmt.Sprintf("Device unreachability detected (error code: %s, reason %s).", errorCode, failureReason))
				} else {
					resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", "POST", err, res.String()))
					break
				}
			}
		}

		{{- else}}
		body := toCreate.toBody(ctx, {{camelCase .Name}}{}) // Convert to request body
		params := ""
		{{- if hasCreateQueryPath .Attributes}}
			{{- $createQueryPath := getCreateQueryPath .Attributes}}
		params += "/" + url.QueryEscape(plan.{{toGoName $createQueryPath.TfName}}.Value{{$createQueryPath.Type}}())
		{{- end}}
		res, err := c.Post(plan.getPath() + params, body {{- if .MaxAsyncWaitTime }}, func(r *cw.Req) { r.MaxAsyncWaitTime={{.MaxAsyncWaitTime}} }{{end}}{{- if .Mutex }}, cw.UseMutex{{- end}})
		if err != nil {
		{{- if .DeviceUnreachabilityWarning}}
			errorCode := res.Get("response.errorCode").String()
			failureReason := res.Get("response.failureReason").String()
			deviceFailureMatch, _ := regexp.MatchString(`(?i)Operation failed on '\d+' devices`, failureReason)
			if deviceFailureMatch {
				resp.Diagnostics.AddWarning("Device Unreachability Warning", fmt.Sprintf("Device unreachability detected (error code: %s, reason %s).", errorCode, failureReason))
			} else {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
				return
			}
		{{- else}}
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (%s), got error: %s, %s", {{- if .PutCreate }} "PUT" {{- else }} "POST" {{- end }}, err, res.String()))
			return
		{{- end}}
		}
		{{- end}}
		{{- $queryParams := generateQueryParamString "GET" "plan" .Attributes }}

		{{- if .IdQueryParam}}
		params += "?{{.IdQueryParam}}=" + url.QueryEscape(plan.Id.ValueString())
		{{- else if and (hasQueryParam .Attributes) (not .GetRequiresId)}}
		{{- if $queryParams }}
		params += {{$queryParams}}
		{{- end}}
		{{- else if and (not .GetNoId) (not .GetFromAll)}}
		params += "/" + url.QueryEscape(plan.Id.ValueString())
		{{- end}}
		{{- if .GetExtraQueryParams}}
		params += "{{.GetExtraQueryParams}}"
		{{- end}}
		res, err = {{readGet . "plan" false}}{{if .GetRestEndpoint}}"{{.GetRestEndpoint}}"{{else}}plan.getPath(){{end}} + params)
		{{- if .FallbackRestEndpoint }}

		if err != nil && (strings.Contains(err.Error(), "StatusCode 404") || strings.Contains(err.Error(), "StatusCode 500")) {
			tflog.Debug(ctx, fmt.Sprintf("%s: Primary endpoint returned 404 or 500, trying fallback endpoint", plan.Id.ValueString()))
			res, err = {{readGet . "plan" false}}plan.getFallbackPath() + params)
		}
		{{- end}}
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
			return
		}


		plan.fromBodyUnknowns(ctx, res)
	}

	if len(toUpdate.{{toGoName $items}}) > 0 {

		{{$chunks := false}}
		{{- $resName := .Name }}
		{{- $MaxElementsInRootList := 0 }} 
		{{- range  .Attributes}}{{- if .MaxElementsInRootList}}
		{{- $MaxElementsInRootList = .MaxElementsInRootList }} 
		maxElementsPerShard := {{$MaxElementsInRootList}}
		var updateList []{{$name}}
		{{- if $itemsIsMap}}
		chunkMap := make(map[string]{{camelCase $.Name}}{{toGoName $items}})
		count := 0
		for k, v := range toUpdate.{{toGoName $items}} {
			chunkMap[k] = v
			count++
			if count >= maxElementsPerShard {
				currentPlanForShard := plan
				currentPlanForShard.{{toGoName $items}} = chunkMap
				updateList = append(updateList, currentPlanForShard)
				chunkMap = make(map[string]{{camelCase $.Name}}{{toGoName $items}})
				count = 0
			}
		}
		if count > 0 {
			currentPlanForShard := plan
			currentPlanForShard.{{toGoName $items}} = chunkMap
			updateList = append(updateList, currentPlanForShard)
		}
		{{- else}}
		for i := 0; i < len(toUpdate.{{toGoName $items}}); i += maxElementsPerShard {
			end := min(i+maxElementsPerShard, len(toUpdate.{{toGoName $items}}))
			chunk := toUpdate.{{toGoName $items}}[i:end]
			currentPlanForShard := plan 
			currentPlanForShard.{{toGoName $items}} = chunk
			updateList = append(updateList, currentPlanForShard)

			
		}
		{{- end}}
		{{$chunks = true}}
		{{- end}}
		{{- end}}

		tflog.Debug(ctx, fmt.Sprintf("%s: Number of items to update: %d", state.Id.ValueString(), len(toUpdate.{{toGoName $items}})))
		{{- if not $itemsIsMap}}
		planIndexMap := make(map[string]int)
		{{- end}}


		{{- if $chunks}}
		{{$noPutAttr := ""}}
		{{$noPutAttrType := ""}}
		{{$idValue := ""}}
		{{- range  .Attributes}}
		{{- if eq .Type "Set"}}
		{{- range  .Attributes}}
		{{- if and .ExcludeTest .NoPut (eq $noPutAttr "")}}
		{{$noPutAttr = .ModelName}}
		{{$noPutAttrType = .Type}}
		{{- end}}
		{{- end}}
		{{- end}}

		{{- if ne $noPutAttr ""}}
		{{$noId := not (hasId .Attributes)}}
		{{- range .Attributes}}
			{{- if not .Computed}}
				{{- if or .Id $noId}}
					{{- if eq .Type "Int64"}}
						{{$idValue = printf "strconv.FormatInt(item.%s.ValueInt64(), 10)" (toGoName .TfName)}}
					{{- else if eq .Type "Bool"}}
						{{$idValue = printf "strconv.FormatBool(item.%s.ValueBool())" (toGoName .TfName)}}
					{{- else if eq .Type "String"}}
						{{$idValue = printf "item.%s.Value%s()" (toGoName .TfName) .Type}}
					{{- end}}
				{{- end}}
			{{- end}}
		{{- end}}

		{{- end}}
		{{- end}}

		{{- if ne $noPutAttr ""}}
		var getState {{toGoName $items}}
		getParams := ""
		{{- $queryParams := generateQueryParamString "GET" "state" .Attributes }}

		{{- if .IdQueryParam}}
		getParams += "?{{.IdQueryParam}}=" + url.QueryEscape(state.Id.ValueString())
		{{- else if and (hasQueryParam .Attributes) (not .GetRequiresId)}}
		{{- if $queryParams }}
		getParams += {{$queryParams}}
		{{- end}}
		{{- else if and (not .GetNoId) (not .GetFromAll)}}
		getParams += "/" + url.QueryEscape(state.Id.ValueString())
		{{- end}}
		{{- if .GetExtraQueryParams}}
		getParams += "{{.GetExtraQueryParams}}"
		{{- end}}
		getRes, err := {{readGet . "state" false}}{{if .GetRestEndpoint}}"{{.GetRestEndpoint}}"{{else}}state.getPath(){{end}} + getParams)
		{{- if .FallbackRestEndpoint }}

		if err != nil && (strings.Contains(err.Error(), "StatusCode 404") || strings.Contains(err.Error(), "StatusCode 500")) {
			tflog.Debug(ctx, fmt.Sprintf("%s: Primary endpoint returned 404 or 500, trying fallback endpoint", state.Id.ValueString()))
			getRes, err = {{readGet . "state" false}}state.getFallbackPath() + getParams)
		}
		{{- end}}
		if err != nil && (strings.Contains(err.Error(), "StatusCode 404") || strings.Contains(err.Error(), "StatusCode 406") || strings.Contains(err.Error(), "StatusCode 500") || strings.Contains(err.Error(), "StatusCode 400") || cw.IsNotFound(err)) {
			resp.State.RemoveResource(ctx)
			return
		} else if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, getRes.String()))
			return
		}

		if getState.isNull(ctx, getRes) {
			getState.fromBody(ctx, getRes)
		} else {
			getState.updateFromBody(ctx, getRes)
		}

		existingNoPut := make(map[string]types.{{$noPutAttrType}})
		{{- range .Attributes}}{{- if eq .Type "Set"}}{{- range .Attributes}}{{- if and .NoPut (ne .ModelName $noPutAttr)}}
		existing{{toGoName .ModelName}} := make(map[string]types.{{.Type}})
		{{- end}}{{- end}}{{- end}}{{- end}}
		for _, item := range getState.{{toGoName $items}} {
			updateKey := {{$idValue}}
			existingNoPut[updateKey] = item.{{toGoName $noPutAttr}}
			{{- range .Attributes}}{{- if eq .Type "Set"}}{{- range .Attributes}}{{- if and .NoPut (ne .ModelName $noPutAttr)}}
			existing{{toGoName .ModelName}}[updateKey] = item.{{toGoName .ModelName}}
			{{- end}}{{- end}}{{- end}}{{- end}}
		}
		{{- end}}

		for _, pl := range updateList {

			{{- range .Attributes}}
			{{- $id := getId .Attributes}}
			{{- if not (eq (toGoName $id.TfName) "") }}
			{{- if $itemsIsMap}}
			for toUpdateKey, _ := range pl.{{toGoName $items}} {
				if updatedItem, exists := planMap[toUpdateKey]; exists {
					{{- if ne $noPutAttr ""}}
					updatedItem.{{toGoName $noPutAttr}} = existingNoPut[toUpdateKey]
					{{- range $.Attributes}}{{- if eq .Type "Set"}}{{- range .Attributes}}{{- if and .NoPut (ne .ModelName $noPutAttr)}}
					if v, ok := existing{{toGoName .ModelName}}[toUpdateKey]; ok {
						updatedItem.{{toGoName .ModelName}} = v
					}
					{{- end}}{{- end}}{{- end}}{{- end}}
					{{- end}}
					plan.{{toGoName $items}}[toUpdateKey] = updatedItem
				}
			}
			{{- else}}
			for i, v := range plan.{{toGoName $items}} {
				planIndexMap[{{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if not .Computed}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(v.{{toGoName .TfName}}.ValueInt64(), 10){{else if eq .Type "Bool"}}strconv.FormatBool(v.{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}v.{{toGoName .TfName}}.Value{{.Type}}(){{end}}{{end}}{{end}}{{end}}] = i
			}
			{{- if ne $noPutAttr ""}}
			for itemInd, item := range pl.{{toGoName $items}} {
			{{- else}}
			for _, item := range pl.{{toGoName $items}} {
			{{- end}}
				toUpdateKey := {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if not .Computed}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(item.{{toGoName .TfName}}.ValueInt64(), 10){{else if eq .Type "Bool"}}strconv.FormatBool(v.{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}item.{{toGoName .TfName}}.Value{{.Type}}(){{end}}{{end}}{{end}}{{end}}
				if updatedItem, exists := planMap[toUpdateKey]; exists {
					if index, found := planIndexMap[toUpdateKey]; found {
						{{- if ne $noPutAttr ""}}
						if noPutValue, noPutOk := existingNoPut[toUpdateKey]; noPutOk {
							updatedItem.{{toGoName $noPutAttr}} = noPutValue
							pl.{{toGoName $items}}[itemInd].{{toGoName $noPutAttr}} = noPutValue
						}
						{{- range $.Attributes}}{{- if eq .Type "Set"}}{{- range .Attributes}}{{- if and .NoPut (ne .ModelName $noPutAttr)}}
						if v, ok := existing{{toGoName .ModelName}}[toUpdateKey]; ok {
							updatedItem.{{toGoName .ModelName}} = v
							pl.{{toGoName $items}}[itemInd].{{toGoName .ModelName}} = v
						}
						{{- end}}{{- end}}{{- end}}{{- end}}
						{{- end}}
						plan.{{toGoName $items}}[index] = updatedItem
					}
				}
			}
			{{- end}}
			{{- end}}
			{{- end}}

			body := pl.toBody(ctx, {{$name}}{})
			params := ""
			res, err := c.Put(plan.getPath()+params, body, cw.UseMutex)
			if err != nil {
				resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s, %s", err, res.String()))
				break
			}
		}
		
		{{- else}}

		{{- range .Attributes}}
		{{- $id := getId .Attributes}}
		{{- if not (eq (toGoName $id.TfName) "") }}
		{{- if $itemsIsMap}}
		for toUpdateKey, _ := range toUpdate.{{toGoName $items}} {
			if updatedItem, exists := planMap[toUpdateKey]; exists {
				plan.{{toGoName $items}}[toUpdateKey] = updatedItem
			}
		}
		{{- else}}
		for i, v := range plan.{{toGoName $items}} {
			planIndexMap[{{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if not .Computed}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(v.{{toGoName .TfName}}.ValueInt64(), 10){{else if eq .Type "Bool"}}strconv.FormatBool(v.{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}v.{{toGoName .TfName}}.Value{{.Type}}(){{end}}{{end}}{{end}}{{end}}] = i
		}
		for _, item := range toUpdate.{{toGoName $items}} {
			toUpdateKey := {{$noId := not (hasId .Attributes)}}{{range .Attributes}}{{if not .Computed}}{{if or .Id $noId}}{{if eq .Type "Int64"}}strconv.FormatInt(item.{{toGoName .TfName}}.ValueInt64(), 10){{else if eq .Type "Bool"}}strconv.FormatBool(v.{{toGoName .TfName}}.ValueBool()), {{else if eq .Type "String"}}item.{{toGoName .TfName}}.Value{{.Type}}(){{end}}{{end}}{{end}}{{end}}
			if updatedItem, exists := planMap[toUpdateKey]; exists {
				if index, found := planIndexMap[toUpdateKey]; found {
					plan.{{toGoName $items}}[index] = updatedItem
				}
			}
		}
		{{- end}}
		{{- end}}
		{{- end}}

		body := toUpdate.toBody(ctx, {{camelCase .Name}}{})
		params := ""
		res, err := c.Put(plan.getPath()+params, body{{- if .Mutex }}, cw.UseMutex{{- end}})
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to configure object (PUT), got error: %s, %s", err, res.String()))
			return
		}

		{{- end}}
	}
	{{- end}}

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}
// End of section. //template:end update

// Section below is generated&owned by "gen/generator.go". //template:begin delete
func (r *{{camelCase .Name}}Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	applyProviderMeta(r.clients, ctx, req.ProviderMeta)
	var state {{camelCase .Name}}

	// Read state
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, ok := r.clientFor(ctx, state.Instance.ValueString(), &resp.Diagnostics)
	if !ok {
		return
	}

	{{- if .UseCache}}
	cacheKey := "{{camelCase .Name}}::" + state.Instance.ValueString() + "::"
	r.cache.DeletePattern(cacheKey)
	{{- end}}

{{if .NoDelete}}	_ = c
{{end}}	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))

	{{- if not .NoDelete}}
	{{- if or .PutDelete .PostDelete}}
	{{- if .DeleteNoId}}
	res, err := c.{{if .PostDelete}}Post{{else}}Put{{end}}({{if .DeleteRestEndpoint}}state.getPathDelete(){{else}}state.getPath(){{end}}, {{if and .DeleteBody (strContains .DeleteBody "{id}")}}strings.ReplaceAll(`{{.DeleteBody}}`, "{id}", state.Id.ValueString()){{else}}`{{if .DeleteBody}}{{.DeleteBody}}{{else}}{}{{end}}`{{end}}{{- if .Mutex }}, cw.UseMutex{{- end}})
	{{- else if .DeleteIdQueryParam}}
	res, err := c.{{if .PostDelete}}Post{{else}}Put{{end}}({{if .DeleteRestEndpoint}}state.getPathDelete(){{else}}state.getPath(){{end}} + "?{{.DeleteIdQueryParam}}=" + url.QueryEscape(state.Id.ValueString()), {{if and .DeleteBody (strContains .DeleteBody "{id}")}}strings.ReplaceAll(`{{.DeleteBody}}`, "{id}", state.Id.ValueString()){{else}}`{{if .DeleteBody}}{{.DeleteBody}}{{else}}{}{{end}}`{{end}}{{- if .Mutex }}, cw.UseMutex{{- end}})
	{{- else if hasDeleteQueryParam .Attributes }}
	{{- $queryParams := generateQueryParamString "DELETE" "state" .Attributes }}
	{{- if $queryParams }}
	params := {{$queryParams}}
	{{- end}}
	res, err := c.{{if .PostDelete}}Post{{else}}Put{{end}}({{if .DeleteRestEndpoint}}state.getPathDelete(){{else}}state.getPath(){{end}} + params, {{if and .DeleteBody (strContains .DeleteBody "{id}")}}strings.ReplaceAll(`{{.DeleteBody}}`, "{id}", state.Id.ValueString()){{else}}`{{if .DeleteBody}}{{.DeleteBody}}{{else}}{}{{end}}`{{end}}{{- if .Mutex }}, cw.UseMutex{{- end}})
	{{- else}}
	res, err := c.{{if .PostDelete}}Post{{else}}Put{{end}}({{if .DeleteRestEndpoint}}state.getPathDelete(){{else}}state.getPath(){{end}} + "/" + url.QueryEscape(state.Id.ValueString()), {{if and .DeleteBody (strContains .DeleteBody "{id}")}}strings.ReplaceAll(`{{.DeleteBody}}`, "{id}", state.Id.ValueString()){{else}}`{{if .DeleteBody}}{{.DeleteBody}}{{else}}{}{{end}}`{{end}}){{- if .Mutex }}, cw.UseMutex{{- end}}
	{{- end}}
	{{- if and .RetryOnErrorCodes .DeleteNoId}}
	if err != nil {
		retryErrorCodes := []string{ {{- range $i, $code := .RetryOnErrorCodes}}{{if $i}}, {{end}}"{{$code}}"{{- end}} }
		errorCode := res.Get("response.errorCode").String()

		shouldRetry := false
		for _, code := range retryErrorCodes {
			if errorCode == code {
				shouldRetry = true
				break
			}
		}

		if shouldRetry {
			maxWaitTime := time.Duration(retryMaxWaitSeconds) * time.Second
			startTime := time.Now()
			retryInterval := 15 * time.Second

			for shouldRetry && time.Since(startTime) < maxWaitTime {
				tflog.Warn(ctx, fmt.Sprintf("%s: Error code %s encountered, waiting %v before retry (elapsed: %v, max: %v)",
					state.Id.ValueString(), errorCode, retryInterval, time.Since(startTime), maxWaitTime))
				time.Sleep(retryInterval)
				res, err = c.{{if .PostDelete}}Post{{else}}Put{{end}}({{if .DeleteRestEndpoint}}state.getPathDelete(){{else}}state.getPath(){{end}}, {{if and .DeleteBody (strContains .DeleteBody "{id}")}}strings.ReplaceAll(`{{.DeleteBody}}`, "{id}", state.Id.ValueString()){{else}}`{{if .DeleteBody}}{{.DeleteBody}}{{else}}{}{{end}}`{{end}}{{- if .Mutex }}, cw.UseMutex{{- end}})

				if err == nil {
					shouldRetry = false
				} else {
					errorCode = res.Get("response.errorCode").String()
					shouldRetry = false
					for _, code := range retryErrorCodes {
						if errorCode == code {
							shouldRetry = true
							break
						}
					}
				}
			}
		}
	}
	{{- end}}
	if err != nil && !strings.Contains(err.Error(), "StatusCode 404") && !cw.IsNotFound(err) {
	{{- if or .PutDelete .PostDelete}}
		errorCode := res.Get("response.errorCode").String()
		if deleteWarnErrorPrefix != "" && strings.HasPrefix(errorCode, deleteWarnErrorPrefix) {
			failureReason := res.Get("response.failureReason").String()
			resp.Diagnostics.AddWarning("Empty input Warning", fmt.Sprintf("Empty input detected (error code: %s, reason %s).", errorCode, failureReason))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (%s), got error: %s, %s", "{{if .PostDelete}}POST{{else}}PUT{{end}}", err, res.String()))
			return
		}
	{{- else}}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (PUT), got error: %s, %s", err, res.String()))
		return
	{{- end}}
	}
	{{- else}}
	{{- if and .DeleteIterative .RootList }}
	{{- if and .RootList .UpdateComputed}}
	{{- $items := "" }}
	{{- range .Attributes}}
    {{- if isNestedListSetMap .}}
	{{- $items = .TfName }}
	{{- end}}
	{{- end}}
	for _, v := range state.{{toGoName $items}} {
		{{- if $.SkipDeleteOnEmptyId}}
		// An empty id would send DELETE to the base endpoint.
		if v.Id.IsNull() || v.Id.IsUnknown() || v.Id.ValueString() == "" {
			tflog.Debug(ctx, fmt.Sprintf("%s: Skipping delete for item - ID is empty or null", state.Id.ValueString()))
			continue
		}
		{{- end}}
		res, err := c.Delete(state.getPath() + "/" + url.QueryEscape(v.Id.ValueString()){{- if .Mutex }}, cw.UseMutex{{- end}})
		if err != nil {
		{{- if .DeviceUnreachabilityWarning}}
			errorCode := res.Get("response.errorCode").String()
			failureReason := res.Get("response.failureReason").String()
			deviceFailureMatch, _ := regexp.MatchString(`(?i)Operation failed on '\d+' devices`, failureReason)
			if deviceFailureMatch {
				resp.Diagnostics.AddWarning("Device Unreachability Warning", fmt.Sprintf("Device unreachability detected (error code: %s, reason %s).", errorCode, failureReason))
			} else {
				resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (%s), got error: %s, %s", "DELETE", err, res.String()))
				return
			}
		{{- else}}
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (%s), got error: %s, %s", "DELETE", err, res.String()))
			return
		{{- end}}
		}
	}
	{{- end}}
	{{- else if .DeleteNoId}}
	{{- if .GetBeforeDelete}}
	// A delete of a missing object may fail, and a missing object only needs to leave the state.
	verifyParams := ""
	{{- $verifyQueryParams := generateQueryParamString "GET" "state" .Attributes }}
	{{- if hasQueryParam .Attributes }}
	{{- if $verifyQueryParams }}
	verifyParams = {{$verifyQueryParams}}
	{{- end}}
	{{- end}}
	verifyRes, verifyErr := {{readGet . "state" false}}{{if .GetRestEndpoint}}{{if strContains .GetRestEndpoint "%v"}}state.getPathGet(){{else}}"{{.GetRestEndpoint}}"{{end}}{{else}}state.getPath(){{end}} + verifyParams)
	{{- if .GetBody}}
	if verifyErr != nil && !strings.Contains(verifyErr.Error(), "StatusCode 404") && !cw.IsNotFound(verifyErr) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET) before delete, got error: %s, %s", verifyErr, verifyRes.String()))
		return
	}
	{{- end}}
	{{- if .IdFromAttribute}}
	{{- $verifyId := getId .Attributes}}
	if verifyErr != nil || !verifyRes.Get("{{responseEnvelope .Attributes}}.0.{{$verifyId.ModelName}}").Exists() || verifyRes.Get("{{responseEnvelope .Attributes}}.0.{{$verifyId.ModelName}}").String() != state.Id.ValueString() {
	{{- else}}
	if verifyErr != nil || !verifyRes.Get("{{responseEnvelope .Attributes}}.0.id").Exists() || verifyRes.Get("{{responseEnvelope .Attributes}}.0.id").String() != state.Id.ValueString() {
	{{- end}}
		tflog.Debug(ctx, fmt.Sprintf("%s: Resource not found or ID mismatch during delete verification, removing from state only", state.Id.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Resource verified, proceeding with delete", state.Id.ValueString()))
	{{- end}}
	{{- if .DeleteBody}}
	res, err := c.DeleteBody({{if .DeleteRestEndpoint}}state.getPathDelete(){{else}}state.getPath(){{end}}, {{if strContains .DeleteBody "{id}"}}strings.ReplaceAll(`{{.DeleteBody}}`, "{id}", state.idJSON()){{else}}`{{.DeleteBody}}`{{end}}{{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(state.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- else}}
	res, err := c.Delete({{if .DeleteRestEndpoint}}state.getPathDelete(){{else}}state.getPath(){{end}}{{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(state.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- end}}
	{{- else if .DeleteIdQueryParam}}
	res, err := c.Delete({{if .DeleteRestEndpoint}}state.getPathDelete(){{else}}state.getPath(){{end}} + "?{{.DeleteIdQueryParam}}=" + url.QueryEscape(state.Id.ValueString()){{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(state.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- else if and (hasDeleteQueryParam .Attributes) (not .SkipDeleteOnEmptyId) (not .GetBeforeDelete) }}
	{{- $queryParams := generateQueryParamString "DELETE" "state" .Attributes }}
	{{- if $queryParams }}
	params := {{$queryParams}}
	{{- end}}
	res, err := c.Delete({{if .DeleteRestEndpoint}}state.getPathDelete(){{else}}state.getPath(){{end}} + params{{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(state.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- else}}
	{{- if .SkipDeleteOnEmptyId}}
	// An empty id would send DELETE to the base endpoint.
	if state.Id.IsNull() || state.Id.IsUnknown() || state.Id.ValueString() == "" {
		tflog.Debug(ctx, fmt.Sprintf("%s: Skipping delete - ID is empty or null", state.Id.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}
	{{- end}}
	{{- if .GetBeforeDelete}}
	{{- if not .IdFromAttribute}}
	// Validate ID is a proper UUID to prevent path traversal attacks
	if err := uuid.Validate(state.Id.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object: ID is not a valid UUID: %s", state.Id.ValueString()))
		return
	}
	{{- end}}
	// Verify resource exists with GET before DELETE to prevent mass deletion from path traversal attacks
	verifyParams := ""
	{{- $verifyQueryParams := generateQueryParamString "GET" "state" .Attributes }}
	{{- if hasQueryParam .Attributes }}
	{{- if $verifyQueryParams }}
	verifyParams = {{$verifyQueryParams}}
	{{- end}}
	{{- end}}
	verifyRes, verifyErr := {{readGet . "state" false}}{{if .GetRestEndpoint}}{{if strContains .GetRestEndpoint "%v"}}state.getPathGet(){{else}}"{{.GetRestEndpoint}}"{{end}}{{else}}state.getPath(){{end}} + verifyParams)
	{{- if .GetBody}}
	if verifyErr != nil && !strings.Contains(verifyErr.Error(), "StatusCode 404") && !cw.IsNotFound(verifyErr) {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET) before delete, got error: %s, %s", verifyErr, verifyRes.String()))
		return
	}
	{{- end}}
	{{- if .IdFromAttribute}}
	{{- $verifyId := getId .Attributes}}
	if verifyErr != nil || !verifyRes.Get("{{responseEnvelope .Attributes}}.0.{{$verifyId.ModelName}}").Exists() || verifyRes.Get("{{responseEnvelope .Attributes}}.0.{{$verifyId.ModelName}}").String() != state.Id.ValueString() {
	{{- else}}
	if verifyErr != nil || !verifyRes.Get("{{responseEnvelope .Attributes}}.0.id").Exists() || verifyRes.Get("{{responseEnvelope .Attributes}}.0.id").String() != state.Id.ValueString() {
	{{- end}}
		tflog.Debug(ctx, fmt.Sprintf("%s: Resource not found or ID mismatch during delete verification, removing from state only", state.Id.ValueString()))
		resp.State.RemoveResource(ctx)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("%s: Resource verified, proceeding with delete", state.Id.ValueString()))
	{{- end}}
	{{- if hasDeleteQueryParam .Attributes }}
	{{- $deleteQueryParams := generateDeleteOnlyQueryParamString "state" .Attributes }}
	{{- if $deleteQueryParams }}
	deleteParams := {{$deleteQueryParams}}
	{{- end}}
	{{- end}}
	{{- if and .DeleteRestEndpoint (strContains .DeleteRestEndpoint "%v")}}
	res, err := c.Delete(state.getPathDelete(){{- if hasDeleteQueryParam .Attributes }} + deleteParams{{- end}}{{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(state.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- else}}
	res, err := c.Delete({{if .DeleteRestEndpoint}}state.getPathDelete(){{else}}state.getPath(){{end}} + "/" + url.QueryEscape(state.Id.ValueString()){{- if hasDeleteQueryParam .Attributes }} + deleteParams{{- end}}{{- if .Mutex }}, cw.UseMutex{{- end}}{{- if .PlanEndpoint}}, cw.PlanPath(state.getPathPlan()){{- if .NoWait}}, cw.NoWait{{- end}}{{- if .MaxAsyncWaitTime}}, cw.MaxAsyncWaitTime({{.MaxAsyncWaitTime}}){{- end}}{{- end}})
	{{- end}}
	{{- end}}
	{{- if not .DeleteIterative }}
	if err != nil && !strings.Contains(err.Error(), "StatusCode 404") && !cw.IsNotFound(err) {
	{{- if .DeviceUnreachabilityWarning}}
		errorCode := res.Get("response.errorCode").String()
		failureReason := res.Get("response.failureReason").String()
		deviceFailureMatch, _ := regexp.MatchString(`(?i)Operation failed on '\d+' devices`, failureReason)
		if deviceFailureMatch {
			resp.Diagnostics.AddWarning("Device Unreachability Warning", fmt.Sprintf("Device unreachability detected (error code: %s, reason %s).", errorCode, failureReason))
		} else {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (%s), got error: %s, %s", "DELETE", err, res.String()))
			return
		}
	{{- else}}
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to delete object (DELETE), got error: %s, %s", err, res.String()))
		return
	{{- end}}
	}
	{{- end}}
	{{- end}}
	{{- end}}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}
// End of section. //template:end delete

// Section below is generated&owned by "gen/generator.go". //template:begin import
{{- if not .NoImport}}
func (r *{{camelCase .Name}}Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	{{- if constantId .}}
	if req.ID != "{{constantId .}}" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: {{constantId .}}. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	{{- else if importAttributes .}}
	idParts := strings.Split(req.ID, ",")

    if len(idParts) != {{len (importAttributes .)}}{{range $index, $attr := importAttributes .}} || idParts[{{$index}}] == ""{{end}} {
        resp.Diagnostics.AddError(
            "Unexpected Import Identifier",
            fmt.Sprintf("Expected import identifier with format: {{range $index, $attr := importAttributes .}}{{if $index}},{{end}}<{{$attr.TfName}}>{{end}}. Got: %q", req.ID),
        )
        return
    }

	{{- $config := .}}
	{{- range $index, $attr := importAttributes .}}
    resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("{{$attr.TfName}}"), idParts[{{$index}}])...)
	{{- if and $config.IdFromAttribute $attr.Id}}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[{{$index}}])...)
	{{- end}}
	{{- end}}
	{{- else}}
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	{{- end}}
}
{{- end}}
// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin readcache
func (r *{{camelCase .Name}}Resource) ReadCache(ctx context.Context, c *cw.Client, req resource.ReadRequest, state {{camelCase .Name}}, params string) (cw.Res, error) {
	var err error
	cacheKey := "{{camelCase .Name}}::" + state.Instance.ValueString() + "::"

	_, cacheSuffix, found := strings.Cut(params, "?")
	queryPart, err := url.ParseQuery(cacheSuffix)
	if err == nil{
		delete(queryPart, "id")
		{{- range .CacheFilterAttributes}}
		delete(queryPart, "{{.}}")
		{{- end}}
		newQuery := queryPart.Encode()
		cacheSuffix  =  "?" + newQuery
		cacheKey += cacheSuffix
	}
	
	cachedValue, found := r.cache.Get(cacheKey)
	if found {
		tflog.Debug(ctx, fmt.Sprintf("hit cache for %s", cacheKey))
		{{- if .CacheRestEndpoint}}
		cwRes, ok := cachedValue.(cw.Res)
		if ok {
			{{- if .CacheFilterAttributes}}
			filteredValue := cwRes.Get({{cacheFilterGjsonExpr .Attributes .CacheFilterAttributes}})
			wrappedRes := cw.Body{}.SetRaw("response", "[" + filteredValue.Raw + "]").Res()
			{{- else}}
			filteredValue := cwRes.Get("response.#(id==\"" + state.Id.ValueString() + "\")")
			wrappedRes := cw.Body{}.SetRaw("response", filteredValue.Raw).Res()
			{{- end}}
			return wrappedRes, nil
		}
		tflog.Info(ctx, fmt.Sprintf("Invalid cache entry type for %s", cacheKey))
		r.cache.Delete(cacheKey) 
		{{- else}}
		cwRes, ok := cachedValue.(cw.Res)
		if ok {
			return cwRes, nil
		}
		tflog.Info(ctx, fmt.Sprintf("Invalid cache entry type for %s", cacheKey))
		r.cache.Delete(cacheKey) 
		{{- end}}
	}
	{{- if .CacheRestEndpoint}}
	{{- if strContains .CacheRestEndpoint "?"}}
	res, err := c.Get("{{.CacheRestEndpoint}}" + strings.Replace(cacheSuffix, "?", "&", 1))
	{{- else}}
	res, err := c.Get("{{.CacheRestEndpoint}}" + cacheSuffix)
	{{- end}}
	{{- if .CacheFilterAttributes}}
	singleRes := res.Get({{cacheFilterGjsonExpr .Attributes .CacheFilterAttributes}})
	singleRes = cw.Body{}.SetRaw("response", "[" + singleRes.Raw + "]").Res()
	{{- else}}
	singleRes := res.Get("response.#(id==\"" + state.Id.ValueString() + "\")")
	singleRes = cw.Body{}.SetRaw("response", singleRes.Raw).Res()
	{{- end}}
	{{- else}}
	res, err := {{readGet . "state" true}}{{if .GetRestEndpoint}}"{{.GetRestEndpoint}}"{{else}}state.getPath(){{end}} + params)
	singleRes := res
	{{- end}}
	if err == nil {
		tflog.Debug(ctx, fmt.Sprintf("set cache for %s", cacheKey))
		r.cache.Set(cacheKey, res)
	}
	return singleRes, err
}
// End of section. //template:end readcache

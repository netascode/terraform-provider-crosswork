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
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

func TestAAASSOBodySendsTheFourMembersAndNoId(t *testing.T) {
	data := AAASSOServices{Services: []AAASSOServicesServices{{
		Name: types.StringValue("cwtestsso"), EvaluationOrder: types.Int64Value(3), Metadata: types.StringValue("PE1ldGFkYXRhLz4="), MetadataLocation: types.StringValue("cwtest-metadata.xml"),
	}}}
	e := gjson.Parse(data.toBody(context.Background(), AAASSOServices{})).Get("registeredServiceProvider.0")
	assert.Equal(t, "cwtestsso", e.Get("name").String())
	assert.Equal(t, int64(3), e.Get("evaluationOrder").Int())
	assert.Equal(t, "PE1ldGFkYXRhLz4=", e.Get("metadata").String())
	assert.Equal(t, "cwtest-metadata.xml", e.Get("metadataLocation").String())
	assert.False(t, e.Get("id").Exists(), "the controller keeps the id of an existing service and numbers a new one")
}

func TestAAASSOReadKeepsTheMetadataOfTheStateBecauseTheControllerReturnsItEmpty(t *testing.T) {
	data := AAASSOServices{Id: types.StringValue("aaa_sso_services"), Services: []AAASSOServicesServices{{
		Name: types.StringValue("cwtestsso"), EvaluationOrder: types.Int64Value(3), Metadata: types.StringValue("PE1ldGFkYXRhLz4="), MetadataLocation: types.StringValue("cwtest-metadata.xml"),
	}}}
	data.updateFromBody(context.Background(), gjson.Parse(`{"registeredServiceProvider":[{"name":"cwtestsso","evaluationOrder":3,"metadata":"","metadataLocation":"cwtest-metadata.xml","id":2,"serviceId":"https://cwtest-sp.example.com/saml/metadata/"}]}`))
	assert.Equal(t, "PE1ldGFkYXRhLz4=", data.Services[0].Metadata.ValueString(), "metadata is write-only")
	assert.Equal(t, int64(3), data.Services[0].EvaluationOrder.ValueInt64())
}

func TestAAASSOServicesEmptyListIsSentSoThatADestroyOrAnUpdateCanClearIt(t *testing.T) {
	assert.Equal(t, `{"registeredServiceProvider":[]}`, AAASSOServices{}.toBody(context.Background(), AAASSOServices{}))
}

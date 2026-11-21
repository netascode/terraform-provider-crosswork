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

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
)

func TestAAALDAPBodySendsTheNumbersAsNumbersAndTheCredentialInItsMember(t *testing.T) {
	data := AAALDAPServers{Servers: []AAALDAPServersServers{{
		Name: types.StringValue("cwtest-ldap"), LdapUrl: types.StringValue("ldap://192.0.2.30:389"), BindDn: types.StringValue("cn=admin,dc=example,dc=org"),
		BindCredential: types.StringValue("AAASecret"), BindCredentialWo: types.StringNull(), BaseDn: types.StringValue("dc=example,dc=org"),
		UserFilter: types.StringValue("uid={user}"), DnFormat: types.StringValue("uid=%s,dc=example,dc=org"), PrincipalAttributeId: types.StringValue("uid"),
		PolicyId: types.StringValue("crossworkPolicyId"), Order: types.Int64Value(50), ConnectTimeout: types.Int64Value(5),
	}}}
	e := gjson.Parse(data.toBody(context.Background(), AAALDAPServers{})).Get("ldap_server.0")
	assert.Equal(t, gjson.Number, e.Get("order").Type, "the controller refuses the order as a string")
	assert.Equal(t, gjson.Number, e.Get("connectTimeout").Type)
	assert.Equal(t, "AAASecret", e.Get("bindCredential").String())
	assert.Equal(t, "ldap://192.0.2.30:389", e.Get("ldapUrl").String())
}

func TestAAALDAPServersEmptyListIsSentSoThatADestroyOrAnUpdateCanClearIt(t *testing.T) {
	assert.Equal(t, `{"ldap_server":[],"ldap_attr_servers":{"ldap_attr_server":[]}}`, AAALDAPServers{}.toBody(context.Background(), AAALDAPServers{}))
}

func TestAAALDAPBodyHasTheSecondListThatMapsTheRole(t *testing.T) {
	data := AAALDAPServers{AttributeServers: []AAALDAPServersAttributeServers{{
		LdapUrl: types.StringValue("ldap://192.0.2.30:389"), BaseDn: types.StringValue("ou=groups,dc=example,dc=org"),
		BindDn: types.StringValue("cn=admin,dc=example,dc=org"), BindCredential: types.StringValue("AAASecret"), BindCredentialWo: types.StringNull(),
		UserFilter: types.StringValue("memberOf"), ConnectTimeout: types.Int64Value(5), FailFast: types.BoolValue(false),
		Attributes: types.MapValueMust(types.StringType, map[string]attr.Value{"crossworkPolicyId": types.StringValue("policy_id")}),
	}}}
	e := gjson.Parse(data.toBody(context.Background(), AAALDAPServers{})).Get("ldap_attr_servers.ldap_attr_server.0")
	assert.Equal(t, "ldap://192.0.2.30:389", e.Get("ldapUrl").String())
	assert.Equal(t, "ou=groups,dc=example,dc=org", e.Get("baseDn").String())
	assert.Equal(t, "AAASecret", e.Get("bindCredential").String())
	assert.Equal(t, "memberOf", e.Get("userFilter").String())
	assert.Equal(t, gjson.False, e.Get("failFast").Type)
	assert.Equal(t, "policy_id", e.Get("attributes.crossworkPolicyId").String(), "the key is the directory attribute, the value is policy_id")
}

func TestAAALDAPReadOfTheSecondListKeepsTheCredentialOfTheState(t *testing.T) {
	data := AAALDAPServers{Id: types.StringValue("aaa_ldap_servers"), AttributeServers: []AAALDAPServersAttributeServers{{
		LdapUrl: types.StringValue("ldap://192.0.2.30:389"), BaseDn: types.StringValue("ou=groups,dc=example,dc=org"),
		BindDn: types.StringValue("cn=admin,dc=example,dc=org"), BindCredential: types.StringValue("AAASecret"),
		UserFilter: types.StringValue("memberOf"), ConnectTimeout: types.Int64Value(5), FailFast: types.BoolValue(false),
		Attributes: types.MapValueMust(types.StringType, map[string]attr.Value{"crossworkPolicyId": types.StringValue("policy_id")}),
	}}}
	data.updateFromBody(context.Background(), gjson.Parse(`{"ldap_server":[],"ldap_attr_servers":{"ldap_attr_server":[{"ldapUrl":"ldap://192.0.2.30:389","baseDn":"ou=groups,dc=example,dc=org","bindDn":"cn=admin,dc=example,dc=org","bindCredential":"******","userFilter":"memberOf","connectTimeout":5,"failFast":false,"attributes":{"crossworkPolicyId":"policy_id"}}]}}`))
	assert.Equal(t, "AAASecret", data.AttributeServers[0].BindCredential.ValueString(), "the masked value is never stored")
	assert.Equal(t, "ou=groups,dc=example,dc=org", data.AttributeServers[0].BaseDn.ValueString())
}

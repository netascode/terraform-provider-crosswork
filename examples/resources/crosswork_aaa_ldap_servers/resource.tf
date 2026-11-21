resource "crosswork_aaa_ldap_servers" "example" {
  servers = [
    {
      name                       = "cwtest-ldap"
      ldap_url                   = "ldap://192.0.2.30:389"
      bind_dn                    = "cn=admin,dc=example,dc=org"
      bind_credential_wo         = "AAASecret"
      bind_credential_wo_version = 1
      base_dn                    = "dc=example,dc=org"
      user_filter                = "uid={user}"
      dn_format                  = "uid=%s,dc=example,dc=org"
      principal_attribute_id     = "uid"
      policy_id                  = "crossworkPolicyId"
      order                      = 50
      connect_timeout            = 5
    }
  ]
  attribute_servers = [
    {
      ldap_url                   = "ldap://192.0.2.30:389"
      base_dn                    = "ou=groups,dc=example,dc=org"
      bind_dn                    = "cn=admin,dc=example,dc=org"
      bind_credential_wo         = "AAASecret"
      bind_credential_wo_version = 1
      user_filter                = "memberOf"
      connect_timeout            = 5
      fail_fast                  = false
      attributes                 = { crossworkPolicyId = "policy_id" }
    }
  ]
}

resource "crosswork_credential_profiles" "example" {
  profile = "cwtest-profile"
  credentials = [
    {
      type                             = "ROBOT_USERPASS_SSH"
      user_name                        = "cwtest-user"
      password_wo                      = "Cwtest-Password-1"
      password_wo_version              = 1
      enable_password                  = "ROBOT_USERPASSEN_ENABLE"
      enable_password_value_wo         = "Cwtest-Enable-1"
      enable_password_value_wo_version = 1
    }
  ]
  snmpv2_read_community_wo          = "cwtest-read"
  snmpv2_read_community_wo_version  = 1
  snmpv2_write_community_wo         = "cwtest-write"
  snmpv2_write_community_wo_version = 1
  snmpv3_user_name                  = "cwtest-snmp"
  snmpv3_security_level             = "SL_AUTH_PRIV"
  snmpv3_auth_type                  = "AT_HMAC_SHA"
  snmpv3_auth_password_wo           = "Cwtest-Auth-1"
  snmpv3_auth_password_wo_version   = 1
  snmpv3_priv_type                  = "PT_CFB_AES_128"
  snmpv3_priv_password_wo           = "Cwtest-Priv-1"
  snmpv3_priv_password_wo_version   = 1
}

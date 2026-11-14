resource "crosswork_data_gateway_global_parameters" "example" {
  cli_sessions          = 3
  ssh_session_timeout   = 120
  snmp_trap_port        = 1062
  syslog_udp_port       = 9514
  syslog_tcp_port       = 9898
  syslog_tls_port       = 6514
  resync_snmpv3_details = false
}

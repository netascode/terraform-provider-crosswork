resource "crosswork_trap_destination" "example" {
  name                 = "cwtest-trap"
  address              = "192.0.2.30"
  address_type         = "IPADDRESS"
  port                 = 162
  transport_type       = "UDP"
  snmp_version         = "VERSION_2C"
  community_wo         = "cwtest-community"
  community_wo_version = 1
}

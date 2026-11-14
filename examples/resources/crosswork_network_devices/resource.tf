resource "crosswork_network_devices" "example" {
  name               = "cwtest-gen-device"
  host_name          = "cwtest-device"
  credential_profile = "cwtest-device-profile"
  admin_state        = "ROBOT_ADMIN_STATE_UP"
  reachability_check = "REACH_CHECK_ENABLE"
  serial_number      = "FOC1234ABCD"
  mac_address        = "00:1b:44:11:3a:b7"
  encryption         = "aes128-ctr"
  device_type        = "NODE_TYPE_ROUTER"
  software_type      = "IOS XR"
  syslog_format      = "RFC_5424"
  capabilities       = ["YANG_CLI"]
  cli_cache_enabled  = false
  provider_profile   = "nso"
  connectivity = [
    {
      protocol      = "ROBOT_MSVC_TRANS_SSH"
      ip_address    = "192.0.2.231"
      prefix_length = 24
      port          = 22
      timeout       = 60
    }
  ]
  isis_system_id                  = "0000.0000.0001"
  ospf_router_id                  = "192.0.2.1"
  te_router_id                    = "192.0.2.2"
  ipv6_router_id                  = "2001:db8::1"
  telemetry_vrf                   = "default"
  telemetry_source_interface      = "Loopback0"
  telemetry_opt_out_of_mdt_config = false
  location_building               = "Building 1"
  location_street                 = "Main Street"
  location_city                   = "Amsterdam"
  location_state                  = "North Holland"
  location_country                = "Netherlands"
  location_region                 = "EMEA"
  location_zip                    = "1011"
  location_latitude               = 52.37
  location_longitude              = 4.89
  location_altitude               = 10
  tags = [
    {
      name = "cwtest-tag"
    }
  ]
}

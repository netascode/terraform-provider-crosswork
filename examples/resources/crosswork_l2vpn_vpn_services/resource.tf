resource "crosswork_l2vpn_vpn_services" "example" {
  vpn_id                             = "cwtest-l2"
  vpn_type                           = "ietf-vpn-common:vpws-evpn"
  evi_id                             = 3961
  evi_source                         = 1
  evi_target                         = 2
  service_assurance_profile_name     = "Gold_L2VPN_ConfigProfile system"
  service_assurance_rule_name        = "Rule-L2VPN-NM system"
  service_assurance_monitoring_state = "enable"
  service_assurance_preservation     = "remove"
  vpn_nodes = [
    {
      vpn_node_id                            = "build-router-1"
      multi_home_ethernet_segment_identifier = "00.11.22.33.44.55.66.77.88"
      sr_policy                              = "cwtest-policy"
      sr_policy_fallback                     = "disable"
      vpn_network_accesses = [
        {
          id                 = "1"
          description        = "cwtest access"
          interface_id       = "GigabitEthernet0/0/0/9"
          encapsulation_type = "ietf-vpn-common:dot1q"
          cvlan_id           = 3961
          mtu                = 1500
          oam_md_name        = "EVC3961"
          oam_md_level       = 4
          oam_maid           = "3961"
          oam_mep_id         = 1
          oam_message_period = "1s"
          oam_id_type        = "number"
          oam_y1731_profiles = [
            {
              name              = "cwtest-y1731"
              schedule_interval = 5
              schedule_duration = "5"
            }
          ]
        }
      ]
    }
  ]
}

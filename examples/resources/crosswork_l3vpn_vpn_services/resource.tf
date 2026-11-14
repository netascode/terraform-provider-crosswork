resource "crosswork_l3vpn_vpn_services" "example" {
  vpn_id                             = "cwtest-svc"
  vpn_service_topology               = "ietf-vpn-common:any-to-any"
  service_assurance_profile_name     = "Gold_L3VPN_ConfigProfile system"
  service_assurance_rule_name        = "Rule-L3VPN-NM system"
  service_assurance_monitoring_state = "enable"
  service_assurance_preservation     = "preserve"
  vpn_instance_profiles = [
    {
      profile_id = "SITE1"
      rd         = "0:65001:101"
      local_as   = 65001
      address_families = [
        {
          address_family = "ietf-vpn-common:ipv4"
          vpn_targets = [
            {
              id                = 1
              route_target_type = "both"
              route_targets = [
                {
                  route_target = "0:65001:4242"
                }
              ]
            }
          ]
          import_policy = "cwtest-policy"
          export_policy = "cwtest-policy"
        }
      ]
    }
  ]
  vpn_nodes = [
    {
      vpn_node_id = "build-router-1"
      local_as    = 65001
      active_vpn_instance_profiles = [
        {
          profile_id = "SITE1"
          rd         = "0:65001:202"
        }
      ]
      vpn_network_accesses = [
        {
          id                 = "1"
          interface_id       = "GigabitEthernet0/0/0/5"
          encapsulation_type = "ietf-vpn-common:dot1q"
          cvlan_id           = 100
          ipv4_address       = "10.10.10.1"
          ipv4_prefix_length = 30
          ipv6_address       = "2001:db8:10:10::1"
          ipv6_prefix_length = 64
          mtu                = 1600
          qos_profiles = [
            {
              profile   = "GOLD-IN"
              direction = "ietf-vpn-common:wan-to-site"
            }
          ]
          routing_protocols = [
            {
              id                 = "1"
              type               = "ietf-vpn-common:bgp-routing"
              bgp_peer_as        = 65100
              bgp_address_family = "ipv4"
              bgp_neighbors      = ["10.10.10.2"]
              bgp_multihop       = 5
              bgp_redistribute_connected = [
                {
                  address_family = "ietf-vpn-common:ipv4"
                  enable         = true
                  metric         = 10
                }
              ]
              bgp_update_source_if_type   = "GigabitEthernet"
              bgp_update_source_if_id     = "0/0/0/5"
              bgp_update_source_sub_if_id = 100
            }
          ]
        }
      ]
    }
  ]
}

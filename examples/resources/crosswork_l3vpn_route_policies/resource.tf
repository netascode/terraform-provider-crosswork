resource "crosswork_l3vpn_route_policies" "example" {
  name = "cwtest-policy"
  statements = [
    {
      name                    = "1"
      match_source_prefix_set = "cwtest-source-prefixes"
      match_dest_prefix_set   = "cwtest-dest-prefixes"
      policy_result           = "accept-route"
      set_ext_community       = "cwtest-tags"
    }
  ]
}

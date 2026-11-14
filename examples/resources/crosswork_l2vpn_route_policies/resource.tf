resource "crosswork_l2vpn_route_policies" "example" {
  name = "cwtest-l2-policy"
  statements = [
    {
      name              = "1"
      match_rd_set      = "cwtest-l2-rds"
      set_ext_community = "cwtest-l2-tags"
    }
  ]
}

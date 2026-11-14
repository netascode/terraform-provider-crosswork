resource "crosswork_l2vpn_evpn_route_type_sets" "example" {
  name            = "cwtest-l2-route-types"
  evpn_route_type = [1]
}

resource "crosswork_l2vpn_id_pools" "example" {
  evi_id_pool            = "cwtest-l2-evi"
  evi_source_target_pool = "cwtest-l2-evi-st"
  rt_pool                = "cwtest-l2-rt"
}

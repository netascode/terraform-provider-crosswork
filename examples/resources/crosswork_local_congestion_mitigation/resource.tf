resource "crosswork_local_congestion_mitigation" "example" {
  domain_id                            = "200"
  enable                               = false
  color                                = 2000
  utilization_threshold                = 100
  utilization_hold_margin              = 5
  delete_lsps                          = false
  profile_id                           = 0
  congestion_check_interval            = 900
  max_tte_group                        = 8
  include_all_interfaces               = false
  description                          = "LCM startup config"
  auto_repair_solution                 = false
  stay_in_area                         = false
  adjacency_hop_type                   = "UNPROTECTED"
  operation_mode                       = "MANUAL"
  objective                            = "MIN_IGP_METRIC"
  pce_dispatch_timeout                 = 180
  congestion_check_suspension_interval = 600
  overprovision_factor                 = 0
  uneven_ecmp_traffic_threshold        = 0
  throttle_mode_threshold              = 5
  history_retention_time               = 30
  debug_opt                            = false
  debug_opt_max_files                  = 40
  affinity_include_any                 = ["30"]
  affinity_exclude_any                 = ["30"]
  affinity_include_all                 = ["30"]
  node_tag_max_sid_depth               = { GRPC_MSL = "10" }
}

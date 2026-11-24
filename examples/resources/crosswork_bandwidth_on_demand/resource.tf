resource "crosswork_bandwidth_on_demand" "example" {
  enable                        = false
  primary_objective             = "MIN_METRIC"
  link_utilization_threshold    = 85
  reoptimization_interval       = 61
  metric_reoptimization_hours   = 2
  metric_reoptimization_minutes = 15
  sr_policy_traffic             = "MEASURED"
  model_update_throttle         = 5
  optimizer_event_threshold     = 50
  optimizer_run_time_limit      = 70
  allowed_policy_violations     = "STRICT"
  prefer_strict_sids            = true
  debug_optimizer               = true
  debug_optimizer_max_files     = 31
}

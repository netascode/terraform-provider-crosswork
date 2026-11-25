resource "crosswork_demand_matrix_global_configuration" "example" {
  computation_cadence    = 30
  minimum_demand_traffic = 100
  debug_mode             = true
  debug_max_plan_files   = 5
}

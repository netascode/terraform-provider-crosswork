resource "crosswork_data_retention_network_performance" "example" {
  collect_lsp_state_change = false
  collect_lsp_path_change  = false
  collect_lsp_utilization  = false
  collect_lsp_pm           = true
  collect_link_pm          = true
  raw_hours                = 12
  hourly_days              = 5
  daily_days               = 20
  weekly_weeks             = 30
}

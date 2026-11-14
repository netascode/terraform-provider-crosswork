resource "crosswork_data_retention_device_performance" "example" {
  rows = {
    "CEPM_INTERFACE" = {
      raw_hours    = 24
      hourly_hours = 168
      daily_hours  = 744
      weekly_hours = 9072
    }
  }
}

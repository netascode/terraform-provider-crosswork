resource "crosswork_alarm_settings" "example" {
  security_device_alarms_cleared_days     = 31
  non_security_device_alarms_cleared_days = 8
  all_events_days                         = 31
  device_alarms_days                      = 61
  system_alarms_days                      = 62
  network_alarms_days                     = 63
  audit_logs_days                         = 31
  trap_processing                         = false
  syslog_processing                       = false
}
